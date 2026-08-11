"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";
import { cookies } from "next/headers";

const API = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function token(): Promise<string> {
  return (await cookies()).get("admin_token")?.value || "";
}

async function mutate(path: string, method: string, body?: unknown) {
  const res = await fetch(`${API}${path}`, {
    method,
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${await token()}`,
    },
    body: body ? JSON.stringify(body) : undefined,
  });
  const json = await res.json().catch(() => null);
  if (!res.ok) {
    const msg = json?.error || `Request failed (${res.status})`;
    throw new Error(msg);
  }
  return json;
}

async function del(path: string) {
  const res = await fetch(`${API}${path}`, {
    method: "DELETE",
    headers: { Authorization: `Bearer ${await token()}` },
  });
  if (!res.ok) {
    const json = await res.json().catch(() => null);
    throw new Error(json?.error || `Request failed (${res.status})`);
  }
}

export async function createProject(formData: FormData): Promise<{ error?: string; ok?: boolean }> {
  try {
    const payload = projectPayload(formData);
    await mutate("/api/admin/projects", "POST", payload);
    revalidatePath("/admin/projects");
    revalidatePath("/projects");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function updateProject(id: number, formData: FormData): Promise<{ error?: string; ok?: boolean }> {
  try {
    const payload = projectPayload(formData);
    await mutate(`/api/admin/projects/${id}`, "PUT", payload);
    revalidatePath("/admin/projects");
    revalidatePath("/projects");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function deleteProject(id: number): Promise<{ error?: string; ok?: boolean }> {
  try {
    await del(`/api/admin/projects/${id}`);
    revalidatePath("/admin/projects");
    revalidatePath("/projects");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

function projectPayload(f: FormData) {
  return {
    slug: String(f.get("slug") || "").trim(),
    title: String(f.get("title") || "").trim(),
    tagline: String(f.get("tagline") || "").trim(),
    description: String(f.get("description") || ""),
    tech: String(f.get("tech") || "")
      .split(",")
      .map((s) => s.trim())
      .filter(Boolean),
    repoUrl: String(f.get("repoUrl") || "").trim(),
    demoUrl: String(f.get("demoUrl") || "").trim(),
    imageUrl: String(f.get("imageUrl") || "").trim(),
    featured: f.get("featured") === "on",
    published: f.get("published") === "on",
    sortOrder: Number(f.get("sortOrder") || 0),
  };
}

export async function createPost(formData: FormData): Promise<{ error?: string; ok?: boolean }> {
  try {
    const payload = postPayload(formData);
    await mutate("/api/admin/posts", "POST", payload);
    revalidatePath("/admin/posts");
    revalidatePath("/journal");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function updatePost(id: number, formData: FormData): Promise<{ error?: string; ok?: boolean }> {
  try {
    const payload = postPayload(formData);
    await mutate(`/api/admin/posts/${id}`, "PUT", payload);
    revalidatePath("/admin/posts");
    revalidatePath("/journal");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function deletePost(id: number): Promise<{ error?: string; ok?: boolean }> {
  try {
    await del(`/api/admin/posts/${id}`);
    revalidatePath("/admin/posts");
    revalidatePath("/journal");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

function postPayload(f: FormData) {
  return {
    slug: String(f.get("slug") || "").trim(),
    title: String(f.get("title") || "").trim(),
    body: String(f.get("body") || ""),
    postType: String(f.get("postType") || "journal"),
    published: f.get("published") === "on",
  };
}

export async function createPhoto(formData: FormData): Promise<{ error?: string; ok?: boolean }> {
  try {
    const payload = {
      url: String(formData.get("url") || "").trim(),
      caption: String(formData.get("caption") || "").trim(),
      location: String(formData.get("location") || "").trim(),
      takenAt: String(formData.get("takenAt") || "").trim() || null,
      published: formData.get("published") === "on",
    };
    await mutate("/api/admin/photos", "POST", payload);
    revalidatePath("/admin/photos");
    revalidatePath("/photos");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function deletePhoto(id: number): Promise<{ error?: string; ok?: boolean }> {
  try {
    await del(`/api/admin/photos/${id}`);
    revalidatePath("/admin/photos");
    revalidatePath("/photos");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function deleteComment(id: number): Promise<{ error?: string; ok?: boolean }> {
  try {
    await del(`/api/admin/comments/${id}`);
    revalidatePath("/admin/comments");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function deleteMessage(id: number): Promise<{ error?: string; ok?: boolean }> {
  try {
    await del(`/api/admin/messages/${id}`);
    revalidatePath("/admin/messages");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function setSiteSetting(key: string, value: unknown): Promise<{ error?: string; ok?: boolean }> {
  try {
    await mutate(`/api/admin/site-settings/${key}`, "PUT", { value });
    revalidatePath("/admin/settings");
    revalidatePath("/");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function uploadImage(formData: FormData): Promise<{ url?: string; error?: string }> {
  try {
    const file = formData.get("file") as File | null;
    if (!file) return { error: "No file provided" };
    const body = new FormData();
    body.append("file", file);
    const res = await fetch(`${API}/api/admin/uploads`, {
      method: "POST",
      headers: { Authorization: `Bearer ${await token()}` },
      body,
    });
    const json = await res.json().catch(() => null);
    if (!res.ok) {
      return { error: json?.error || `Upload failed (${res.status})` };
    }
    revalidatePath("/admin/photos");
    revalidatePath("/photos");
    return { url: json.url as string };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function uploadAndCreatePhoto(formData: FormData): Promise<{ error?: string; ok?: boolean }> {
  try {
    const file = formData.get("file") as File | null;
    if (!file) return { error: "No file selected" };
    const body = new FormData();
    body.append("file", file);
    const res = await fetch(`${API}/api/admin/uploads`, {
      method: "POST",
      headers: { Authorization: `Bearer ${await token()}` },
      body,
    });
    const json = await res.json().catch(() => null);
    if (!res.ok) {
      return { error: json?.error || `Upload failed (${res.status})` };
    }
    const payload = {
      url: json.url as string,
      caption: String(formData.get("caption") || "").trim(),
      location: String(formData.get("location") || "").trim(),
      takenAt: String(formData.get("takenAt") || "").trim() || null,
      published: formData.get("published") === "on",
    };
    await mutate("/api/admin/photos", "POST", payload);
    revalidatePath("/admin/photos");
    revalidatePath("/photos");
    return { ok: true };
  } catch (e: unknown) {
    return { error: e instanceof Error ? e.message : String(e) };
  }
}

export async function logout() {
  (await cookies()).delete("admin_token");
  redirect("/admin/login");
}
