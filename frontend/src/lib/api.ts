export type Project = {
  id: number;
  slug: string;
  title: string;
  tagline: string;
  description: string;
  tech: string[];
  repoUrl: string;
  demoUrl: string;
  imageUrl: string;
  featured: boolean;
  sortOrder: number;
  published: boolean;
  createdAt: string;
  updatedAt: string;
};

export type Post = {
  id: number;
  slug: string;
  title: string;
  body: string;
  postType: "journal" | "note";
  published: boolean;
  createdAt: string;
  updatedAt: string;
};

export type Comment = {
  id: number;
  postId: number;
  author: string;
  body: string;
  createdAt: string;
};

export type Photo = {
  id: number;
  url: string;
  caption: string;
  location: string;
  takenAt: string | null;
  published: boolean;
  createdAt: string;
};

const API = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8080";

export class ApiError extends Error {
  status: number;

  constructor(status: number, message: string) {
    super(message);
    this.status = status;
  }
}

async function getJSON<T>(path: string, revalidate = 60): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    next: { revalidate },
    headers: { Accept: "application/json" },
  });
  if (!res.ok) {
    throw new ApiError(res.status, `API request failed: ${res.status}`);
  }
  return res.json() as Promise<T>;
}

export async function getProjects(): Promise<Project[]> {
  return getJSON("/api/projects");
}

export async function getProject(slug: string): Promise<Project> {
  return getJSON(`/api/projects/${slug}`);
}

export async function getPosts(limit?: number): Promise<Post[]> {
  const qs = limit ? `?limit=${limit}` : "";
  return getJSON(`/api/posts${qs}`);
}

export async function getPost(slug: string): Promise<Post> {
  return getJSON(`/api/posts/${slug}`);
}

export async function getComments(slug: string): Promise<Comment[]> {
  return getJSON(`/api/posts/${slug}/comments`);
}

export async function getPhotos(): Promise<Photo[]> {
  return getJSON("/api/photos");
}

export type SiteSettings = {
  hero?: { title?: string; tagline?: string };
  about?: { body?: string };
};

export async function getSite(): Promise<SiteSettings> {
  return getJSON("/api/site");
}

export async function postComment(slug: string, data: { author: string; body: string }): Promise<Comment> {
  const res = await fetch(`${API}/api/posts/${slug}/comments`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  const json = await res.json().catch(() => null);
  if (!res.ok) {
    const msg = json && json.error ? json.error : `Request failed (${res.status})`;
    throw new ApiError(res.status, msg);
  }
  return json as Comment;
}

export async function postContact(data: {
  name: string;
  email: string;
  body: string;
}): Promise<void> {
  const res = await fetch(`${API}/api/contact`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
  });
  const json = await res.json().catch(() => null);
  if (!res.ok) {
    const msg = json && json.error ? json.error : `Request failed (${res.status})`;
    throw new ApiError(res.status, msg);
  }
}

export function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

export function excerpt(markdown: string, max = 180): string {
  const plain = markdown
    .replace(/```[\s\S]*?```/g, " ")
    .replace(/`([^`]*)`/g, "$1")
    .replace(/!\[[^\]]*\]\([^)]*\)/g, " ")
    .replace(/\[([^\]]*)\]\([^)]*\)/g, "$1")
    .replace(/[#>*_~|-]/g, " ")
    .replace(/\s+/g, " ")
    .trim();
  return plain.length > max ? `${plain.slice(0, max).trimEnd()}…` : plain;
}
