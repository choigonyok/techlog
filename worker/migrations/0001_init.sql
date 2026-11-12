-- 글. 새로 쓰거나 고칠 때의 40자 제한은 API 가 지킨다. 옛 원고에는 40자를 넘는 제목이 있어
-- 테이블은 100자까지 받는다(SQLite length 는 글자 수를 센다).
CREATE TABLE posts (
  id TEXT PRIMARY KEY,
  title TEXT NOT NULL CHECK (length(title) BETWEEN 1 AND 100),
  subtitle TEXT NOT NULL DEFAULT '' CHECK (length(subtitle) <= 100),
  text TEXT NOT NULL,
  write_time TEXT NOT NULL DEFAULT '',           -- YYYY-MM-DD(한국 시간), 옛 글 중 날짜 없는 글은 ''
  thumbnail TEXT,                                -- R2 키, 없으면 NULL
  created_at TEXT NOT NULL DEFAULT (datetime('now')),
  updated_at TEXT NOT NULL DEFAULT (datetime('now'))
);
CREATE INDEX posts_write_time ON posts (write_time DESC);

-- 태그는 적은 대로 둔다. 대소문자 묶기는 화면이 한다. rowid 순서가 적은 순서다.
CREATE TABLE tags (
  post_id TEXT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  PRIMARY KEY (post_id, name)
);

-- 본문 이미지. R2 키와 원래 주소를 남긴다.
CREATE TABLE images (
  id TEXT PRIMARY KEY,
  post_id TEXT NOT NULL REFERENCES posts (id) ON DELETE CASCADE,
  source TEXT NOT NULL DEFAULT '',
  content_type TEXT NOT NULL DEFAULT ''
);
CREATE INDEX images_post ON images (post_id);

-- 방문 수. 날짜별 행과 누적 한 행을 따로 두어 누적을 셀 때 모든 날을 읽지 않는다.
CREATE TABLE visitor (
  date TEXT PRIMARY KEY,
  count INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE visitor_total (
  id INTEGER PRIMARY KEY CHECK (id = 1),
  count INTEGER NOT NULL DEFAULT 0
);
