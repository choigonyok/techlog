// import-posts.mjs 가 받아 둔 그림을 R2(techlog-images)에 올린다.
//
//   node scripts/upload-images.mjs [--out import] (--local | --remote) [--parallel n]
import { spawn } from "node:child_process";
import fs from "node:fs";
import path from "node:path";

const args = process.argv.slice(2);
const opt = (name, fallback) => (args.includes(name) ? args[args.indexOf(name) + 1] : fallback);
const out = opt("--out", "import");
const where = args.includes("--remote") ? "--remote" : args.includes("--local") ? "--local" : null;
// 로컬 저장소는 여러 프로세스가 한꺼번에 쓰면 잠겨 실패하므로 하나씩 올린다.
const parallel = Number(opt("--parallel", where === "--local" ? "1" : "6"));
if (!where) {
  console.error("--local 또는 --remote 를 고르세요.");
  process.exit(1);
}

const rows = fs
  .readFileSync(path.join(out, "images.tsv"), "utf8")
  .split("\n")
  .filter(Boolean)
  .map((line) => line.split("\t"));

// VS Code 터미널은 ELECTRON_RUN_AS_NODE 를 켜 두어 wrangler 가 띄우는 프로그램이 깨질 수 있어 뺀다.
const env = { ...process.env };
delete env.ELECTRON_RUN_AS_NODE;

function put([key, file, type]) {
  return new Promise((resolve) => {
    const child = spawn("npx", ["wrangler", "r2", "object", "put", `techlog-images/${key}`, "--file", file, "--content-type", type, where], {
      env,
      stdio: ["ignore", "ignore", "pipe"],
    });
    let err = "";
    child.stderr.on("data", (d) => (err += d));
    child.on("close", (code) => resolve(code === 0 ? null : `${key}: ${err.trim().split("\n").pop()}`));
  });
}

const queue = [...rows];
const failures = [];
let done = 0;
await Promise.all(
  Array.from({ length: parallel }, async () => {
    while (queue.length) {
      const failure = await put(queue.shift());
      if (failure) failures.push(failure);
      done += 1;
      if (done % 20 === 0 || done === rows.length) console.log(`${done}/${rows.length}`);
    }
  })
);
if (failures.length) {
  console.error(`실패 ${failures.length}장:\n  ${failures.join("\n  ")}`);
  process.exit(1);
}
console.log(`R2(${where.slice(2)})에 ${rows.length}장을 올렸습니다.`);
