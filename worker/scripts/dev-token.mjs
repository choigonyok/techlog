// 로컬 시험용 관리자 토큰을 만든다. 기본으로 .dev.vars 의 JWT_SECRET 으로 서명한다.
//
//   node scripts/dev-token.mjs [login=choigonyok] [--ttl 초] [--secret 값]
//
// 운영 비밀 값으로는 만들지 않는다(운영 토큰은 GitHub 로그인으로만 받는다).
import crypto from "node:crypto";
import fs from "node:fs";

const args = process.argv.slice(2);
const opt = (name, fallback) => (args.includes(name) ? args[args.indexOf(name) + 1] : fallback);
const login = args[0] && !args[0].startsWith("--") ? args[0] : "choigonyok";
const ttl = Number(opt("--ttl", String(12 * 60 * 60)));
let secret = opt("--secret", null);
if (!secret) {
  const vars = fs.existsSync(".dev.vars") ? fs.readFileSync(".dev.vars", "utf8") : "";
  secret = (/^JWT_SECRET=(.*)$/m.exec(vars) || [])[1];
}
if (!secret) {
  console.error(".dev.vars 에 JWT_SECRET 이 없습니다.");
  process.exit(1);
}

const b64url = (buf) => Buffer.from(buf).toString("base64url");
const now = Math.floor(Date.now() / 1000);
const header = b64url(JSON.stringify({ alg: "HS256", typ: "JWT" }));
const payload = b64url(JSON.stringify({ sub: login, iat: now, exp: now + ttl }));
const sig = crypto.createHmac("sha256", secret).update(`${header}.${payload}`).digest("base64url");
process.stdout.write(`${header}.${payload}.${sig}\n`);
