import { requireAdmin } from "./auth";
import { AppEnv, HttpError, json, readJson, todayKST } from "./http";
import { IMAGE_PATH, localImageKeys, rehostImages } from "./images";

const LIMIT = 40;
const MAX_THUMBNAIL = 15 * 1024 * 1024;

type PostRow = { id: string; title: string; subtitle: string; write_time: string; thumbnail: string | null; text?: string };
type Input = { title: string; subtitle: string; tags: string[]; text: string };

function toJson(row: PostRow, tags: string[]) {
  return {
    id: row.id,
    title: row.title,
    subtitle: row.subtitle,
    writeTime: row.write_time,
    tags,
    thumbnail: row.thumbnail,
    ...(row.text !== undefined ? { text: row.text } : {}),
  };
}

const charLength = (s: string) => Array.from(s).length;

// 화면과 같은 규칙으로 확인하고, 어긴 이유를 그대로 돌려준다.
function validate(body: Partial<Record<keyof Input, unknown>>): Input {
  const title = typeof body.title === "string" ? body.title.trim() : "";
  const subtitle = typeof body.subtitle === "string" ? body.subtitle.trim() : "";
  const text = typeof body.text === "string" ? body.text : "";
  if (!title) throw new HttpError(400, "제목을 적어 주세요.");
  if (charLength(title) > LIMIT) throw new HttpError(400, `제목은 ${LIMIT}자까지입니다.`);
  if (charLength(subtitle) > LIMIT) throw new HttpError(400, `부제는 ${LIMIT}자까지입니다.`);
  if (!text.trim()) throw new HttpError(400, "본문을 적어 주세요.");
  const seen = new Set<string>();
  const tags: string[] = [];
  for (const raw of Array.isArray(body.tags) ? body.tags : []) {
    const name = String(raw).trim().replace(/\s+/g, " ");
    if (!name || seen.has(name.toLowerCase())) continue;
    if (charLength(name) > LIMIT) throw new HttpError(400, `주제는 ${LIMIT}자까지입니다: ${name.slice(0, 20)}`);
    seen.add(name.toLowerCase());
    tags.push(name);
  }
  return { title, subtitle, tags, text };
}

export async function listPosts(env: AppEnv, tagFilter: string | null) {
  const [posts, tags] = await env.DB.batch<Record<string, string>>([
    env.DB.prepare("SELECT id, title, subtitle, write_time, thumbnail FROM posts ORDER BY write_time DESC, title"),
    env.DB.prepare("SELECT post_id, name FROM tags ORDER BY rowid"),
  ]);
  const tagsOf = new Map<string, string[]>();
  for (const t of tags.results) {
    if (!tagsOf.has(t.post_id)) tagsOf.set(t.post_id, []);
    (tagsOf.get(t.post_id) as string[]).push(t.name);
  }
  let list = (posts.results as unknown as PostRow[]).map((p) => toJson(p, tagsOf.get(p.id) || []));
  const wanted = (tagFilter || "")
    .split(",")
    .map((t) => t.trim().toLowerCase())
    .filter((t) => t && t !== "all");
  if (wanted.length) list = list.filter((p) => p.tags.some((t) => wanted.includes(t.toLowerCase())));
  return list;
}

async function findPost(env: AppEnv, id: string, withText = false): Promise<PostRow | null> {
  return env.DB.prepare(
    `SELECT id, title, subtitle, write_time, thumbnail${withText ? ", text" : ""} FROM posts WHERE id = ?1`
  )
    .bind(id)
    .first<PostRow>();
}

export async function getPost(env: AppEnv, id: string) {
  const post = await findPost(env, id, true);
  if (!post) throw new HttpError(404, "없는 글입니다.");
  const { results } = await env.DB.prepare("SELECT name FROM tags WHERE post_id = ?1 ORDER BY rowid").bind(id).all<{ name: string }>();
  return toJson(post, results.map((r) => r.name));
}

function tagStatements(env: AppEnv, id: string, tags: string[]) {
  return tags.map((name) => env.DB.prepare("INSERT INTO tags (post_id, name) VALUES (?1, ?2)").bind(id, name));
}

function imageStatements(env: AppEnv, id: string, created: { id: string; source: string; contentType: string }[]) {
  return created.map((img) =>
    env.DB.prepare("INSERT INTO images (id, post_id, source, content_type) VALUES (?1, ?2, ?3, ?4)").bind(img.id, id, img.source, img.contentType)
  );
}

export async function createPost(request: Request, env: AppEnv): Promise<Response> {
  await requireAdmin(request, env);
  const input = validate(await readJson(request));
  const id = crypto.randomUUID();
  const { text, created } = await rehostImages(input.text, env, new URL(request.url).origin);
  try {
    await env.DB.batch([
      env.DB.prepare("INSERT INTO posts (id, title, subtitle, text, write_time) VALUES (?1, ?2, ?3, ?4, ?5)").bind(
        id,
        input.title,
        input.subtitle,
        text,
        todayKST()
      ),
      ...tagStatements(env, id, input.tags),
      ...imageStatements(env, id, created),
    ]);
  } catch (error) {
    if (created.length) await env.IMAGES.delete(created.map((c) => c.id));
    throw error;
  }
  return json({ id }, 201);
}

export async function updatePost(request: Request, env: AppEnv, id: string): Promise<Response> {
  await requireAdmin(request, env);
  if (!(await findPost(env, id))) throw new HttpError(404, "없는 글입니다.");
  const input = validate(await readJson(request));
  const { text, created } = await rehostImages(input.text, env, new URL(request.url).origin);
  const kept = localImageKeys(text);
  const { results: old } = await env.DB.prepare("SELECT id FROM images WHERE post_id = ?1").bind(id).all<{ id: string }>();
  const removed = old.map((r) => r.id).filter((key) => !kept.has(key));
  try {
    await env.DB.batch([
      env.DB.prepare("UPDATE posts SET title = ?1, subtitle = ?2, text = ?3, updated_at = datetime('now') WHERE id = ?4").bind(
        input.title,
        input.subtitle,
        text,
        id
      ),
      env.DB.prepare("DELETE FROM tags WHERE post_id = ?1").bind(id),
      ...tagStatements(env, id, input.tags),
      ...imageStatements(env, id, created),
      ...removed.map((key) => env.DB.prepare("DELETE FROM images WHERE id = ?1").bind(key)),
    ]);
  } catch (error) {
    if (created.length) await env.IMAGES.delete(created.map((c) => c.id));
    throw error;
  }
  // 본문에서 빠진 그림은 저장이 끝난 뒤에 지운다.
  if (removed.length) await env.IMAGES.delete(removed);
  return json({ id });
}

// 썸네일은 파일을 본문 그대로 받아 R2 로 흘려 넣는다(multipart 를 풀지 않아 CPU 를 거의 쓰지 않는다).
export async function putThumbnail(request: Request, env: AppEnv, id: string): Promise<Response> {
  await requireAdmin(request, env);
  const post = await findPost(env, id);
  if (!post) throw new HttpError(404, "없는 글입니다.");
  const type = (request.headers.get("Content-Type") || "").split(";")[0].trim().toLowerCase();
  if (!type.startsWith("image/")) throw new HttpError(415, "이미지 파일만 썸네일로 쓸 수 있습니다.");
  const length = Number(request.headers.get("Content-Length") || "0");
  if (!length) throw new HttpError(411, "파일 크기를 알 수 없습니다.");
  if (length > MAX_THUMBNAIL) throw new HttpError(413, "썸네일은 15MB 까지 올릴 수 있습니다.");
  const key = crypto.randomUUID();
  await env.IMAGES.put(key, request.body, { httpMetadata: { contentType: type } });
  await env.DB.prepare("UPDATE posts SET thumbnail = ?1, updated_at = datetime('now') WHERE id = ?2").bind(key, id).run();
  if (post.thumbnail) await env.IMAGES.delete(post.thumbnail);
  return json({ thumbnail: key });
}

// 예전 화면이 쓰던 썸네일 주소는 이미지 주소로 넘긴다.
export async function thumbnailRedirect(request: Request, env: AppEnv, id: string): Promise<Response> {
  const post = await findPost(env, id);
  if (!post || !post.thumbnail) throw new HttpError(404, "썸네일이 없습니다.");
  return Response.redirect(new URL(IMAGE_PATH + post.thumbnail, request.url).toString(), 302);
}

export async function deletePost(request: Request, env: AppEnv, id: string): Promise<Response> {
  await requireAdmin(request, env);
  const post = await findPost(env, id);
  if (!post) throw new HttpError(404, "없는 글입니다.");
  const { results } = await env.DB.prepare("SELECT id FROM images WHERE post_id = ?1").bind(id).all<{ id: string }>();
  await env.DB.batch([
    env.DB.prepare("DELETE FROM tags WHERE post_id = ?1").bind(id),
    env.DB.prepare("DELETE FROM images WHERE post_id = ?1").bind(id),
    env.DB.prepare("DELETE FROM posts WHERE id = ?1").bind(id),
  ]);
  const keys = results.map((r) => r.id).concat(post.thumbnail ? [post.thumbnail] : []);
  if (keys.length) await env.IMAGES.delete(keys);
  return new Response(null, { status: 204 });
}

export async function tagCounts(env: AppEnv) {
  const { results } = await env.DB.prepare("SELECT name, COUNT(*) AS count FROM tags GROUP BY name ORDER BY count DESC").all<{
    name: string;
    count: number;
  }>();
  const total = await env.DB.prepare("SELECT COUNT(*) AS count FROM posts").first<number>("count");
  return [{ name: "ALL", count: total || 0 }, ...results];
}
