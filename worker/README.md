# techlog worker

[tech.choigonyok.com](https://tech.choigonyok.com)을 내보내는 Cloudflare Worker 입니다. 예전 Go(Gin)+Postgres+S3 백엔드를 옮긴 것입니다.

- 화면: `../techlog-ui/build` 를 정적 자산으로 내보낸다(요청은 무료·무제한). 없는 주소는 `index.html` 로 돌린다.
- API: `/api/*` 만 Worker 가 받는다. 경로 목록은 `src/index.ts` 맨 위에 있다.
- 데이터: D1 `techlog`(글·태그·그림·방문 수), R2 `techlog-images`(본문 그림·썸네일).
- 관리자: GitHub OAuth 로 계정을 확인하고 `ADMIN_LOGIN`(choigonyok)이면 12시간짜리 JWT 를 준다. 글 쓰기·고치기·지우기는 이 토큰을 확인한다.

## 비밀 값

| 이름 | 내용 |
| --- | --- |
| `JWT_SECRET` | 관리자 토큰 서명 키(아무 긴 난수) |
| `GITHUB_CLIENT_ID`, `GITHUB_CLIENT_SECRET` | GitHub OAuth 앱. 콜백 주소는 `https://tech.choigonyok.com/github/callback` |

```bash
npx wrangler secret put JWT_SECRET
npx wrangler secret put GITHUB_CLIENT_ID
npx wrangler secret put GITHUB_CLIENT_SECRET
```

GitHub 앱 정보가 없으면 로그인 화면이 '준비되지 않았습니다'를 보이고, 글 읽기는 그대로 된다.

## 로컬

```bash
npm install
printf 'JWT_SECRET=local-dev-secret\n' > .dev.vars
npm run db:migrate:local
git clone https://github.com/choigonyok/techlog-posts /tmp/techlog-posts      # 글 원고(처음 한 번)
node scripts/import-posts.mjs /tmp/techlog-posts                              # import/ 에 시드 SQL·그림
npx wrangler d1 execute techlog --local --file import/seed.sql
node scripts/upload-images.mjs --local
(cd ../techlog-ui && npm install && npm run build)
npm run dev                                                                   # :8787
node scripts/dev-token.mjs        # 관리 화면 시험용 토큰(.dev.vars 의 키로 서명). 브라우저 localStorage 의 jwt_token 에 넣는다
```

`npm run check` 은 타입 검사다.

## 배포

```bash
(cd ../techlog-ui && npm ci && npm run build)
npx wrangler deploy
```

처음 한 번은 D1·R2 를 만들고(`wrangler d1 create techlog`, `wrangler r2 bucket create techlog-images`), `wrangler.jsonc` 의 `database_id` 를 채운 뒤 `npm run db:migrate:remote` 와 위의 `--remote` 들여오기를 한다.

## 무료 플랜에서 지키는 것

- 요청당 CPU 10ms: 썸네일은 multipart 를 풀지 않고 본문 그대로 R2 로 흘려 넣는다.
- 요청당 바깥 요청 50개·동시 연결 6개: 본문의 새 그림은 한 번에 40장까지, 셋씩 내려받는다.
- 방문 누적 수는 따로 한 행에 두어 방문마다 두 행만 쓴다.
