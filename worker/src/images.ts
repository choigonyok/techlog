import { AppEnv, HttpError } from "./http";

export const IMAGE_PATH = "/api/images/";
// 무료 플랜은 요청 하나에 바깥 요청 50개, 동시 연결 6개까지다.
const MAX_NEW = 40;
const PARALLEL = 3;
const MAX_BYTES = 15 * 1024 * 1024;
const KEY = /^[A-Za-z0-9][A-Za-z0-9._-]{7,99}$/;
const IMAGE_MD = /!\[([^\]]*)\]\(\s*<?([^)\s>]+)>?((?:\s+"[^"]*")?)\s*\)/g;

export type NewImage = { id: string; source: string; contentType: string };

export const isImageKey = (key: string) => KEY.test(key);

// 본문이 가리키는 이 서버의 이미지 키들
export function localImageKeys(text: string): Set<string> {
  const keys = new Set<string>();
  for (const m of text.matchAll(IMAGE_MD)) {
    if (m[2].startsWith(IMAGE_PATH)) keys.add(decodeURIComponent(m[2].slice(IMAGE_PATH.length)));
  }
  return keys;
}

async function download(url: string): Promise<{ body: ArrayBuffer; type: string }> {
  const res = await fetch(url, { signal: AbortSignal.timeout(15_000), headers: { "User-Agent": "techlog-worker" } }).catch(
    () => null
  );
  if (!res || !res.ok) throw new HttpError(422, `이미지를 내려받지 못했습니다: ${url.slice(0, 120)}`);
  const type = (res.headers.get("Content-Type") || "").split(";")[0].trim().toLowerCase();
  if (!type.startsWith("image/")) throw new HttpError(422, `이미지가 아닌 주소입니다: ${url.slice(0, 120)}`);
  const body = await res.arrayBuffer();
  if (body.byteLength > MAX_BYTES) throw new HttpError(413, `15MB 를 넘는 이미지입니다: ${url.slice(0, 120)}`);
  return { body, type };
}

// 본문의 바깥 이미지(http/https)를 내려받아 R2 에 올리고 주소를 /api/images/<키> 로 바꾼다.
// 바깥 주소가 사라져도 글의 그림이 남게 하려는 것이다. 하나라도 실패하면 이번에 올린 것을 지우고 멈춘다.
export async function rehostImages(text: string, env: AppEnv, origin: string): Promise<{ text: string; created: NewImage[] }> {
  const ours = `${origin}${IMAGE_PATH}`;
  const external = new Set<string>();
  for (const m of text.matchAll(IMAGE_MD)) {
    const url = m[2];
    if (url.startsWith(IMAGE_PATH) || url.startsWith(ours)) continue;
    if (!/^https?:\/\//i.test(url)) throw new HttpError(400, `이미지 주소는 http:// 나 https:// 로 시작해야 합니다: ${url.slice(0, 80)}`);
    external.add(url);
  }
  if (external.size > MAX_NEW) {
    throw new HttpError(413, `새 이미지는 한 번에 ${MAX_NEW}장까지 올릴 수 있습니다. 나눠서 저장하세요.`);
  }

  const created: NewImage[] = [];
  const keyOf = new Map<string, string>();
  const queue = [...external];
  let failure: unknown = null;
  const worker = async () => {
    while (!failure && queue.length) {
      const url = queue.shift() as string;
      try {
        const { body, type } = await download(url);
        const id = crypto.randomUUID();
        await env.IMAGES.put(id, body, { httpMetadata: { contentType: type } });
        created.push({ id, source: url, contentType: type });
        keyOf.set(url, id);
      } catch (error) {
        failure = failure || error;
      }
    }
  };
  await Promise.all(Array.from({ length: Math.min(PARALLEL, queue.length) }, worker));
  if (failure) {
    if (created.length) await env.IMAGES.delete(created.map((c) => c.id));
    throw failure;
  }

  const rewritten = text.replace(IMAGE_MD, (all, alt: string, url: string, title: string) => {
    const key = keyOf.get(url);
    const target = key ? IMAGE_PATH + key : url.startsWith(ours) ? url.slice(origin.length) : url;
    return target === url ? all : `![${alt}](${target}${title})`;
  });
  return { text: rewritten, created };
}

// 키는 올릴 때마다 새로 만들므로 한 번 받은 그림은 바뀌지 않는다(immutable).
export async function serveImage(request: Request, env: AppEnv, key: string): Promise<Response> {
  if (!isImageKey(key)) throw new HttpError(404, "없는 이미지입니다.");
  const object = await env.IMAGES.get(key, { onlyIf: request.headers });
  if (!object) throw new HttpError(404, "없는 이미지입니다.");
  const headers = new Headers();
  object.writeHttpMetadata(headers);
  headers.set("ETag", object.httpEtag);
  headers.set("Cache-Control", "public, max-age=31536000, immutable");
  headers.set("X-Content-Type-Options", "nosniff");
  // 주소를 직접 열어도 그림 안의 스크립트(SVG)가 돌지 않게 한다.
  headers.set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox");
  if (!("body" in object) || !object.body) return new Response(null, { status: 304, headers });
  return new Response(request.method === "HEAD" ? null : object.body, { headers });
}
