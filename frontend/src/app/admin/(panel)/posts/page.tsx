import Link from "next/link";
import AdminNav from "@/components/AdminNav";
import { adminFetch } from "@/lib/admin";
import { deletePost } from "@/lib/adminActions";

type Post = {
  id: number;
  slug: string;
  title: string;
  postType: "journal" | "note";
  published: boolean;
  createdAt: string;
};

export default async function AdminPostsPage() {
  const posts = await adminFetch<Post[]>("/api/admin/posts");
  return (
    <div>
      <AdminNav active="posts" />
      <div className="flex items-end justify-between">
        <h1 className="font-serif text-3xl text-ink">Posts</h1>
        <Link
          href="/admin/posts/new"
          className="rounded-full bg-moss px-4 py-2 font-mono text-sm text-paper hover:bg-moss-dark"
        >
          + new post
        </Link>
      </div>
      <ul className="mt-6 space-y-3">
        {posts.map((p) => (
          <li
            key={p.id}
            className="flex flex-wrap items-center gap-3 rounded-2xl border border-line bg-paper-deep/40 px-5 py-4"
          >
            <div className="min-w-0 flex-1">
              <p className="truncate font-serif text-lg text-ink">{p.title}</p>
              <p className="truncate font-mono text-xs text-muted">
                {p.slug} · {p.postType} {p.published ? "· live" : "· draft"}
              </p>
            </div>
            <Link
              href={`/admin/posts/${p.id}`}
              className="rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-ink-soft hover:bg-moss hover:text-paper"
            >
              edit
            </Link>
            <form
              action={async () => {
                "use server";
                await deletePost(p.id);
              }}
            >
              <button
                type="submit"
                className="rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-muted hover:bg-ink hover:text-paper"
              >
                delete
              </button>
            </form>
          </li>
        ))}
        {posts.length === 0 ? <li className="font-mono text-sm text-muted">no posts yet</li> : null}
      </ul>
    </div>
  );
}
