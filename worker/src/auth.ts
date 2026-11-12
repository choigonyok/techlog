import { AppEnv, HttpError, json, readCookie, readJson } from "./http";

// 관리자 로그인: GitHub OAuth 로 계정을 확인하고, 허가된 계정(ADMIN_LOGIN)이면 HS256 JWT 를 준다.
// 글을 쓰고 고치고 지우는 요청은 모두 이 토큰을 Authorization 헤더로 확인한다.

const TOKEN_TTL = 12 * 60 * 60;
const STATE_COOKIE = "techlog_oauth_state";
const UA = "techlog-worker";
const enc = new TextEncoder();
const dec = new TextDecoder();

function b64url(bytes: ArrayBuffer | Uint8Array): string {
  const view = bytes instanceof Uint8Array ? bytes : new Uint8Array(bytes);
  let s = "";
  for (const b of view) s += String.fromCharCode(b);
  return btoa(s).replace(/\+/g, "-").replace(/\//g, "_").replace(/=+$/, "");
}

function fromB64url(s: string): Uint8Array {
  const pad = s.replace(/-/g, "+").replace(/_/g, "/") + "===".slice((s.length + 3) % 4);
  const raw = atob(pad);
  const out = new Uint8Array(raw.length);
  for (let i = 0; i < raw.length; i += 1) out[i] = raw.charCodeAt(i);
  return out;
}

function hmacKey(secret: string): Promise<CryptoKey> {
  return crypto.subtle.importKey("raw", enc.encode(secret), { name: "HMAC", hash: "SHA-256" }, false, ["sign", "verify"]);
}

function secretOf(env: AppEnv): string {
  if (!env.JWT_SECRET) throw new HttpError(503, "관리자 로그인이 아직 설정되지 않았습니다.");
  return env.JWT_SECRET;
}

export async function signToken(sub: string, secret: string, ttl = TOKEN_TTL): Promise<string> {
  const now = Math.floor(Date.now() / 1000);
  const header = b64url(enc.encode(JSON.stringify({ alg: "HS256", typ: "JWT" })));
  const payload = b64url(enc.encode(JSON.stringify({ sub, iat: now, exp: now + ttl })));
  const sig = await crypto.subtle.sign("HMAC", await hmacKey(secret), enc.encode(`${header}.${payload}`));
  return `${header}.${payload}.${b64url(sig)}`;
}

type Claims = { sub: string; exp: number };

export async function verifyToken(token: string, secret: string): Promise<Claims | null> {
  const parts = token.split(".");
  if (parts.length !== 3) return null;
  try {
    const header = JSON.parse(dec.decode(fromB64url(parts[0])));
    if (header.alg !== "HS256") return null;
    const ok = await crypto.subtle.verify("HMAC", await hmacKey(secret), fromB64url(parts[2]), enc.encode(`${parts[0]}.${parts[1]}`));
    if (!ok) return null;
    const claims = JSON.parse(dec.decode(fromB64url(parts[1])));
    if (typeof claims.sub !== "string" || typeof claims.exp !== "number") return null;
    if (claims.exp < Math.floor(Date.now() / 1000)) return null;
    return claims;
  } catch {
    return null;
  }
}

const isAdmin = (env: AppEnv, login: string) => login.toLowerCase() === env.ADMIN_LOGIN.toLowerCase();

export async function requireAdmin(request: Request, env: AppEnv): Promise<void> {
  const match = /^Bearer\s+(\S+)$/i.exec(request.headers.get("Authorization") || "");
  if (!match) throw new HttpError(401, "로그인이 필요합니다.");
  const claims = await verifyToken(match[1], secretOf(env));
  if (!claims) throw new HttpError(401, "로그인한 지 오래되었습니다. 다시 로그인하세요.");
  if (!isAdmin(env, claims.sub)) throw new HttpError(403, "허가되지 않은 계정입니다.");
}

// 화면이 관리 화면에 들어올 때 토큰을 확인한다. 200 통과, 401 만료·잘못된 토큰, 403 다른 계정.
export async function checkToken(request: Request, env: AppEnv): Promise<Response> {
  const { token } = await readJson<{ token?: string }>(request);
  const claims = token ? await verifyToken(token, secretOf(env)) : null;
  if (!claims) return json({ error: "로그인한 지 오래되었습니다." }, 401);
  if (!isAdmin(env, claims.sub)) return json({ error: "허가되지 않은 계정입니다." }, 403);
  return json({ ok: true });
}

const stateCookie = (value: string, maxAge: number) =>
  `${STATE_COOKIE}=${value}; Path=/; HttpOnly; Secure; SameSite=Lax; Max-Age=${maxAge}`;

// GitHub 로그인 화면으로 보낸다. 돌아왔을 때 같은 브라우저인지 확인하도록 state 를 쿠키에 둔다.
export function startLogin(request: Request, env: AppEnv): Response {
  const origin = new URL(request.url).origin;
  if (!env.GITHUB_CLIENT_ID || !env.GITHUB_CLIENT_SECRET || !env.JWT_SECRET) {
    return Response.redirect(`${origin}/login?error=unavailable`, 302);
  }
  const state = b64url(crypto.getRandomValues(new Uint8Array(16)));
  const params = new URLSearchParams({
    client_id: env.GITHUB_CLIENT_ID,
    redirect_uri: `${origin}/github/callback`,
    state,
    allow_signup: "false",
  });
  return new Response(null, {
    status: 302,
    headers: { Location: `https://github.com/login/oauth/authorize?${params}`, "Set-Cookie": stateCookie(state, 600) },
  });
}

// 화면의 /github/callback 이 GitHub 가 준 code 를 넘기면 토큰으로 바꿔 준다.
export async function finishLogin(request: Request, env: AppEnv): Promise<Response> {
  const url = new URL(request.url);
  const code = url.searchParams.get("code");
  const state = url.searchParams.get("state");
  const expected = readCookie(request, STATE_COOKIE);
  if (!code || !state || !expected || state !== expected) {
    throw new HttpError(400, "로그인 요청이 맞지 않습니다. 처음부터 다시 로그인하세요.");
  }
  if (!env.GITHUB_CLIENT_ID || !env.GITHUB_CLIENT_SECRET) throw new HttpError(503, "관리자 로그인이 아직 설정되지 않았습니다.");

  const tokenRes = await fetch("https://github.com/login/oauth/access_token", {
    method: "POST",
    headers: { Accept: "application/json", "Content-Type": "application/json", "User-Agent": UA },
    body: JSON.stringify({
      client_id: env.GITHUB_CLIENT_ID,
      client_secret: env.GITHUB_CLIENT_SECRET,
      code,
      redirect_uri: `${url.origin}/github/callback`,
    }),
  });
  const { access_token: accessToken } = (await tokenRes.json()) as { access_token?: string };
  if (!accessToken) throw new HttpError(400, "GitHub 확인이 끝나지 않았습니다. 다시 로그인하세요.");

  const userRes = await fetch("https://api.github.com/user", {
    headers: { Authorization: `Bearer ${accessToken}`, Accept: "application/vnd.github+json", "User-Agent": UA },
  });
  const user = (await userRes.json()) as { login?: string };
  if (!user.login || !isAdmin(env, user.login)) {
    return json({ error: "허가되지 않은 계정입니다." }, 403, { "Set-Cookie": stateCookie("", 0) });
  }
  const token = await signToken(user.login, secretOf(env));
  return json({ token }, 200, { "Set-Cookie": stateCookie("", 0) });
}
