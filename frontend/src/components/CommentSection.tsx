"use client";

import { useState } from "react";
import { formatDate, postComment, type Comment } from "@/lib/api";

export default function CommentSection({ slug, initial }: { slug: string; initial: Comment[] }) {
  const [comments, setComments] = useState<Comment[]>(initial);
  const [author, setAuthor] = useState("");
  const [body, setBody] = useState("");
  const [status, setStatus] = useState<"idle" | "sending" | "error">("idle");
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setStatus("sending");
    setError("");
    try {
      const created = await postComment(slug, { author, body });
      setComments((prev) => [...prev, created]);
      setAuthor("");
      setBody("");
      setStatus("idle");
    } catch (err) {
      setStatus("error");
      setError(err instanceof Error ? err.message : "Could not post the comment. Try again.");
    }
  }

  return (
    <section className="mt-16">
      <h2 className="font-serif text-2xl text-ink">comments</h2>

      {comments.length > 0 ? (
        <ul className="mt-6 space-y-6">
          {comments.map((c) => (
            <li key={c.id} className="border-l-2 border-moss/40 pl-4">
              <p className="font-mono text-sm text-muted">
                {c.author} · {formatDate(c.createdAt)}
              </p>
              <p className="mt-1 leading-relaxed text-ink-soft">{c.body}</p>
            </li>
          ))}
        </ul>
      ) : (
        <p className="mt-6 font-mono text-sm text-muted">
          no comments yet. Be the first.
        </p>
      )}

      <form onSubmit={handleSubmit} className="mt-8 space-y-4 rounded-2xl border border-line bg-paper-deep/40 p-6">
        <div>
          <label htmlFor="comment-author" className="mb-1.5 block font-mono text-sm text-muted">
            name
          </label>
          <input
            id="comment-author"
            type="text"
            required
            maxLength={80}
            value={author}
            onChange={(e) => setAuthor(e.target.value)}
            placeholder="how should we call you?"
            className="w-full rounded-xl border border-line bg-paper px-4 py-2.5 text-ink placeholder:text-muted/70 focus:border-moss focus:outline-none"
          />
        </div>
        <div>
          <label htmlFor="comment-body" className="mb-1.5 block font-mono text-sm text-muted">
            comment
          </label>
          <textarea
            id="comment-body"
            required
            maxLength={2000}
            rows={3}
            value={body}
            onChange={(e) => setBody(e.target.value)}
            placeholder="say something kind or curious"
            className="w-full rounded-xl border border-line bg-paper px-4 py-2.5 text-ink placeholder:text-muted/70 focus:border-moss focus:outline-none"
          />
        </div>
        {status === "error" ? (
          <p role="alert" className="font-mono text-sm text-red-700">
            {error}
          </p>
        ) : null}
        <button
          type="submit"
          disabled={status === "sending"}
          className="rounded-xl bg-moss px-5 py-2.5 font-mono text-sm font-medium text-paper transition-colors hover:bg-moss-dark disabled:cursor-not-allowed disabled:opacity-60"
        >
          {status === "sending" ? "posting…" : "post comment"}
        </button>
      </form>
    </section>
  );
}
