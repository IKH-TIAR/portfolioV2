import { redirect } from "next/navigation";
import { cookies } from "next/headers";

const API = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function requireAuth(): Promise<{ token: string }> {
  const token = (await cookies()).get("admin_token")?.value || "";
  if (!token) {
    redirect("/admin/login");
  }
  const res = await fetch(`${API}/api/auth/me`, {
    cache: "no-store",
    headers: { Authorization: `Bearer ${token}` },
  });
  if (!res.ok) {
    redirect("/admin/login");
  }
  return { token };
}

export async function login(email: string, password: string): Promise<string> {
  const res = await fetch(`${API}/api/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  const json = await res.json().catch(() => null);
  if (!res.ok) {
    throw new Error(json?.error || "Invalid credentials");
  }
  return json.token as string;
}

export async function logoutAction() {
  "use server";
  (await cookies()).set("admin_token", "", {
    httpOnly: true,
    sameSite: "lax",
    path: "/admin",
    maxAge: 0,
  });
}

export async function adminFetch<T>(path: string, init: RequestInit = {}): Promise<T> {
  const { token } = await requireAuth();
  const res = await fetch(`${API}${path}`, {
    cache: "no-store",
    ...init,
    headers: {
      Accept: "application/json",
      ...(init.headers || {}),
      Authorization: `Bearer ${token}`,
    },
  });
  const json = await res.json().catch(() => null);
  if (!res.ok) {
    const msg = json?.error || `Request failed (${res.status})`;
    throw new Error(msg);
  }
  return json as T;
}
