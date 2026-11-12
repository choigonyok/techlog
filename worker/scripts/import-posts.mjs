// 글 원고 저장소(choigonyok/techlog-posts)를 D1 시드 SQL 과 R2 에 올릴 그림 파일로 바꾼다.
//
//   node scripts/import-posts.mjs <원고 폴더> [--out import] [--skip-images]
//
// 결과
//   <out>/seed.sql      글·태그·그림 행. 같은 원고를 다시 들여오면 같은 행을 덮어쓴다.
//   <out>/images/<키>   본문 그림 원본(--skip-images 면 받지 않고 원래 주소를 둔다).
//                       사라진 그림은 /api/images/missing-… 으로 바꾼다.
//   <out>/images.tsv    키, 파일, Content-Type  (scripts/upload-images.mjs 가 읽는다)
import crypto from "node:crypto";
import fs from "node:fs";
import path from "node:path";

const args = process.argv.slice(2);
const src = args[0];
const opt = (name, fallback) => (args.includes(name) ? args[args.indexOf(name) + 1] : fallback);
const out = opt("--out", "import");
const skipImages = args.includes("--skip-images");
if (!src) {
  console.error("사용: node scripts/import-posts.mjs <원고 폴더> [--out import] [--skip-images]");
  process.exit(1);
}

const sha1 = (s) => crypto.createHash("sha1").update(s).digest("hex");
const idOf = (file) => {
  const h = sha1(`techlog-posts/${file}`);
  return `${h.slice(0, 8)}-${h.slice(8, 12)}-5${h.slice(13, 16)}-8${h.slice(17, 20)}-${h.slice(20, 32)}`;
};
const q = (v) => (v === null ? "NULL" : `'${String(v).replace(/'/g, "''")}'`);

function normalizeDate(raw) {
  const m = /^(\d{4})[-/.](\d{1,2})[-/.](\d{1,2})/.exec(String(raw || "").trim());
  return m ? `${m[1]}-${m[2].padStart(2, "0")}-${m[3].padStart(2, "0")}` : "";
}

// 원고 머리: [ID: ..] [Tags: ..] [Title: ..] [WriteTime: ..] [ImageNames: ..] [Subtitle: ..]
function parse(file) {
  const lines = fs.readFileSync(path.join(src, file), "utf8").replace(/\r\n/g, "\n").split("\n");
  const meta = {};
  let i = 0;
  for (; i < Math.min(lines.length, 12); i += 1) {
    const m = /^\s*\[(ID|Tags|Title|WriteTime|ImageNames|Subtitle):\s?(.*)\]\s*$/.exec(lines[i]);
    if (m) meta[m[1]] = m[2].trim();
    else if (lines[i].trim() === "" && Object.keys(meta).length) continue;
    else if (Object.keys(meta).length) break;
  }
  if (!meta.Title) return null;
  return {
    id: idOf(file),
    title: meta.Title,
    subtitle: meta.Subtitle || "",
    writeTime: normalizeDate(meta.WriteTime),
    tags: [...new Set((meta.Tags || "").split(/[\s,]+/).filter(Boolean))],
    text: lines.slice(i).join("\n").replace(/^\s+/, ""),
  };
}

const IMAGE_MD = /!\[([^\]]*)\]\(\s*<?([^)\s>]+)>?((?:\s+"[^"]*")?)\s*\)/g;

async function download(url) {
  for (let attempt = 1; attempt <= 3; attempt += 1) {
    try {
      const res = await fetch(url, { signal: AbortSignal.timeout(20000) });
      const type = (res.headers.get("content-type") || "").split(";")[0].trim().toLowerCase();
      if (res.ok && type.startsWith("image/")) return { body: Buffer.from(await res.arrayBuffer()), type };
      if (res.status === 404 || res.status === 403) return null;
    } catch {
      // 잠시 뒤 다시 시도한다.
    }
    await new Promise((r) => setTimeout(r, 800 * attempt));
  }
  return null;
}

const posts = fs
  .readdirSync(src)
  .filter((f) => f.endsWith(".md") && f !== "README.md")
  .sort()
  .map((f) => parse(f))
  .filter(Boolean);

fs.mkdirSync(path.join(out, "images"), { recursive: true });
const manifest = [];
const failed = [];
const sql = [];
for (const post of posts) {
  const images = [];
  const urls = [...new Set([...post.text.matchAll(IMAGE_MD)].map((m) => m[2]).filter((u) => /^https?:\/\//i.test(u)))];
  const keyOf = new Map();
  if (!skipImages) {
    for (const url of urls) {
      const key = `imp-${sha1(url).slice(0, 24)}`;
      const file = path.join(out, "images", key);
      let type = manifest.find((m) => m.key === key)?.type;
      if (!type) {
        const got = await download(url);
        if (!got) {
          // 사라진 그림은 없는 그림 주소로 바꿔 둔다. 원래 주소로 두면 글을 고칠 때 서버가
          // 다시 내려받으려다 실패해 저장이 막힌다. 화면은 '그림을 불러오지 못했습니다'를 보인다.
          failed.push(`${post.title}: ${url}`);
          keyOf.set(url, `missing-${sha1(url).slice(0, 12)}`);
          continue;
        }
        fs.writeFileSync(file, got.body);
        type = got.type;
        manifest.push({ key, file, type });
      }
      keyOf.set(url, key);
      images.push({ key, url, type });
    }
  }
  const text = post.text.replace(IMAGE_MD, (all, alt, url, title) => (keyOf.has(url) ? `![${alt}](/api/images/${keyOf.get(url)}${title})` : all));

  sql.push(`DELETE FROM tags WHERE post_id = ${q(post.id)};`);
  sql.push(`DELETE FROM images WHERE post_id = ${q(post.id)};`);
  sql.push(
    `INSERT OR REPLACE INTO posts (id, title, subtitle, text, write_time) VALUES (${[post.id, post.title, post.subtitle, text, post.writeTime].map(q).join(", ")});`
  );
  for (const tag of post.tags) sql.push(`INSERT INTO tags (post_id, name) VALUES (${q(post.id)}, ${q(tag)});`);
  for (const img of images) {
    sql.push(`INSERT OR REPLACE INTO images (id, post_id, source, content_type) VALUES (${[img.key, post.id, img.url, img.type].map(q).join(", ")});`);
  }
  console.log(`${post.writeTime || "날짜 없음"}  ${post.title}  (그림 ${images.length}/${urls.length})`);
}

// D1 은 SQL 문 하나가 100KB 를 넘으면 받지 않는다.
const tooLong = sql.filter((s) => Buffer.byteLength(s) > 100_000);
if (tooLong.length) {
  console.error(`100KB 를 넘는 SQL 문 ${tooLong.length}개`);
  process.exit(1);
}
fs.writeFileSync(path.join(out, "seed.sql"), `${sql.join("\n")}\n`);
fs.writeFileSync(path.join(out, "images.tsv"), manifest.map((m) => `${m.key}\t${m.file}\t${m.type}`).join("\n") + (manifest.length ? "\n" : ""));
console.log(`\n글 ${posts.length}편, 그림 ${manifest.length}장 → ${out}/`);
if (failed.length) console.log(`받지 못한 그림 ${failed.length}장(없는 그림 주소로 바꿈):\n  ${failed.join("\n  ")}`);
