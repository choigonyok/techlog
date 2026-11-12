import { AppEnv, json, todayKST } from "./http";

// 화면을 처음 불러올 때 한 번 부른다. 오늘(한국 날짜)과 누적을 하나씩 올리고 새 값을 돌려준다.
export async function countVisit(env: AppEnv): Promise<Response> {
  const [today, total] = await env.DB.batch<{ count: number }>([
    env.DB.prepare(
      "INSERT INTO visitor (date, count) VALUES (?1, 1) ON CONFLICT (date) DO UPDATE SET count = count + 1 RETURNING count"
    ).bind(todayKST()),
    env.DB.prepare(
      "INSERT INTO visitor_total (id, count) VALUES (1, 1) ON CONFLICT (id) DO UPDATE SET count = count + 1 RETURNING count"
    ),
  ]);
  return json({ today: today.results[0].count, total: total.results[0].count });
}
