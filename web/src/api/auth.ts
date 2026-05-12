import client from "./client";

export type SessionUser = {
  id: number;
  username: string;
};

export async function login(payload: { username: string; password: string }) {
  const response = await client.post("/auth/login", payload);
  return response.data as { user: SessionUser; expiresAt: string };
}

export async function logout() {
  await client.post("/auth/logout");
}

export async function fetchMe() {
  const response = await client.get("/auth/me");
  return response.data as SessionUser;
}
