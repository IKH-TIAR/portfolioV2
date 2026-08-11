import AdminNav from "@/components/AdminNav";
import { adminFetch } from "@/lib/admin";
import { deleteComment } from "@/lib/adminActions";
import { formatDate } from "@/lib/api";

type Post = {
  id: number;
  slug: string;
  title: string;
};

type Comment = {
  id: number;
  postId: number;
  author: string;
  body: string;
  createdAt: string;
};

export default async function AdminCommentsPage() {
  const posts = await adminFetch<Post[]>("/api/admin/posts");
  const commentsByPost: { post: Post; comments: Comment[] }[] = await Promise.all(
    posts.map(async (post) => {
      const comments = await adminFetch<Comment[]>(`/api/posts/${post.slug}/comments`);
      return { post, comments };
    })
  );
  const total = commentsByPost.reduce((n, p) => n + p.comments.length, 0);

  return (
    <div>
      <AdminNav active="comments" />
      <h1 className="font-serif text-3xl text-ink">Comments</h1>
      {total === 0 ? (
        <p className="mt-6 font-mono text-sm text-muted">no comments yet</p>
      ) : (
        commentsByPost
          .filter((p) => p.comments.length > 0)
          .map(({ post, comments }) => (
            <section key={post.id} className="mt-8">
              <h2 className="font-serif text-xl text-ink">
                {post.title}{" "}
                <span className="font-mono text-xs text-muted">/{post.slug}</span>
              </h2>
              <ul className="mt-4 space-y-3">
                {comments.map((c) => (
                  <li
                    key={c.id}
                    className="rounded-2xl border border-line bg-paper-deep/40 px-5 py-4"
                  >
                    <div className="flex flex-wrap items-center gap-3">
                      <p className="font-mono text-sm text-ink">{c.author}</p>
                      <span className="ml-auto font-mono text-xs text-muted">
                        {formatDate(c.createdAt)}
                      </span>
                    </div>
                    <p className="mt-2 whitespace-pre-wrap text-ink-soft">{c.body}</p>
                    <form
                      action={async () => {
                        "use server";
                        await deleteComment(c.id);
                      }}
                      className="mt-3"
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
              </ul>
            </section>
          ))
      )}
    </div>
  );
}
