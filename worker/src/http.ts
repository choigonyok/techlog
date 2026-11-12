// 비밀 값은 wrangler.jsonc 에 없으므로 생성된 Env 에 더해 둔다.
export type AppEnv = Env & {
  JWT_SECRET?: string;
  GITHUB_CLIENT_ID?: string;
  GITHUB_CLIENT_SECRET?: string;
};

export class HttpError extends Error {
  constructor(
    readonly status: number,
    message: string
  ) {
    super(message);
  }
}

export function json(data: unknown, status = 200, headers: Record<string, string> = {}): Response {
  return new Response(JSON.stringify(data), {
    status,
    headers: { "Content-Type": "application/json; charset=utf-8", "Cache-Control": "no-store", ...headers },
  });
}

export async function readJson<T>(request: Request, limit = 2_000_000): Promise<T> {
  const text = await request.text();
  if (text.length > limit) throw new HttpError(413, "요청이 너무 큽니다.");
  try {
    return JSON.parse(text) as T;
  } catch {
    throw new HttpError(400, "요청 본문이 JSON 이 아닙니다.");
  }
}

export function readCookie(request: Request, name: string): string | null {
  const header = request.headers.get("Cookie") || "";
  for (const part of header.split(";")) {
    const [key, ...rest] = part.trim().split("=");
    if (key === name) return rest.join("=");
  }
  return null;
}

// 글 작성일과 방문 수는 한국 날짜로 센다.
export function todayKST(now = new Date()): string {
  return new Intl.DateTimeFormat("en-CA", { timeZone: "Asia/Seoul", year: "numeric", month: "2-digit", day: "2-digit" }).format(now);
}
