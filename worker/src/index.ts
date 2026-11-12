import { checkToken, finishLogin, startLogin } from "./auth";
import { AppEnv, HttpError, json } from "./http";
import { serveImage } from "./images";
import { createPost, deletePost, getPost, listPosts, putThumbnail, tagCounts, thumbnailRedirect, updatePost } from "./posts";
import { countVisit } from "./visitor";

// techlog API. 화면(정적 자산)은 Cloudflare 가 바로 내보내고, /api/* 만 여기로 온다.
//
//   GET    /api/posts[?tag=a,b]        글 목록(본문 제외), 작성일 최근 순
//   POST   /api/posts                  새 글 {title, subtitle, tags[], text}      관리자
//   GET    /api/posts/:id              글 한 편
//   PUT    /api/posts/:id              글 고치기                                   관리자
//   DELETE /api/posts/:id              글 지우기(이미지까지)                       관리자
//   PUT    /api/posts/:id/thumbnail    썸네일 파일(본문 그대로, Content-Type: image/*)  관리자
//   GET    /api/images/:key            R2 의 그림
//   GET    /api/visitor                방문 수(+1)
//   GET    /api/tags                   주제별 글 수
//   GET    /api/login                  GitHub 로그인으로 보냄
//   GET    /api/github/callback        GitHub code → 관리자 토큰
//   POST   /api/token                  관리자 토큰 확인
async function route(request: Request, env: AppEnv, url: URL): Promise<Response> {
  const { pathname } = url;
  const method = request.method;
  let m: RegExpExecArray | null;

  if (pathname === "/api/posts") {
    if (method === "GET") return json(await listPosts(env, url.searchParams.get("tag")), 200, { "Cache-Control": "no-cache" });
    if (method === "POST") return createPost(request, env);
  } else if ((m = /^\/api\/posts\/([^/]+)$/.exec(pathname))) {
    const id = decodeURIComponent(m[1]);
    if (method === "GET") return json(await getPost(env, id), 200, { "Cache-Control": "no-cache" });
    if (method === "PUT") return updatePost(request, env, id);
    if (method === "DELETE") return deletePost(request, env, id);
  } else if ((m = /^\/api\/posts\/([^/]+)\/thumbnail$/.exec(pathname))) {
    const id = decodeURIComponent(m[1]);
    if (method === "PUT") return putThumbnail(request, env, id);
    if (method === "GET") return thumbnailRedirect(request, env, id);
  } else if ((m = /^\/api\/images\/([^/]+)$/.exec(pathname))) {
    if (method === "GET" || method === "HEAD") return serveImage(request, env, decodeURIComponent(m[1]));
  } else if (pathname === "/api/visitor") {
    if (method === "GET") return countVisit(env);
  } else if (pathname === "/api/tags") {
    if (method === "GET") return json(await tagCounts(env));
  } else if (pathname === "/api/login") {
    if (method === "GET") return startLogin(request, env);
  } else if (pathname === "/api/github/callback") {
    if (method === "GET") return finishLogin(request, env);
  } else if (pathname === "/api/token") {
    if (method === "POST") return checkToken(request, env);
  } else {
    throw new HttpError(404, "없는 주소입니다.");
  }
  throw new HttpError(405, "허용하지 않는 요청입니다.");
}

export default {
  async fetch(request, env): Promise<Response> {
    const url = new URL(request.url);
    if (!url.pathname.startsWith("/api/")) return env.ASSETS.fetch(request);
    try {
      return await route(request, env as AppEnv, url);
    } catch (error) {
      if (error instanceof HttpError) return json({ error: error.message }, error.status);
      console.error(error);
      return json({ error: "서버에서 오류가 났습니다." }, 500);
    }
  },
} satisfies ExportedHandler<Env>;
