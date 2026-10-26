"use client";

import { useActionState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { createPost, updatePost } from "@/lib/adminActions";

type Props = {
  post?: {
    id: number;
    slug: string;
    title: string;
    body: string;
    postType: "journal" | "note";
    tags: string[];
    published: boolean;
  };
};

type State = { error?: string; ok?: boolean };

export default function PostForm({ post }: Props) {
  const router = useRouter();
  const isNew = !post;
  const [state, formAction] = useActionState<State, FormData>(
    isNew ? (_, fd) => createPost(fd) : post ? (_, fd) => updatePost(post.id, fd) : (_, fd) => createPost(fd),
    {}
  );
  useEffect(() => {
    if (state.ok) {
      router.push("/admin/posts");
    }
  }, [state.ok, router]);
  return (
    <form action={formAction} className="mt-6 space-y-5">
      {state.error ? (
        <p className="rounded-xl border border-red-300 bg-red-50 px-4 py-2 font-mono text-sm text-red-700">
          {state.error}
        </p>
      ) : null}
      <div className="grid gap-5 sm:grid-cols-2">
        <div>
          <label className="font-mono text-xs uppercase tracking-wider text-muted">Title</label>
          <input
            name="title"
            defaultValue={post?.title}
            required
            className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink"
          />
        </div>
        <div>
          <label className="font-mono text-xs uppercase tracking-wider text-muted">Slug</label>
          <input
            name="slug"
            defaultValue={post?.slug}
            required
            className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink"
          />
        </div>
      </div>
      <div>
        <label className="font-mono text-xs uppercase tracking-wider text-muted">Type</label>
        <select
          name="postType"
          defaultValue={post?.postType ?? "journal"}
          className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink"
        >
          <option value="journal">journal</option>
          <option value="note">note</option>
        </select>
      </div>
      <div>
        <label className="font-mono text-xs uppercase tracking-wider text-muted">
          Body (markdown)
        </label>
        <textarea
          name="body"
          rows={18}
          defaultValue={post?.body}
          className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 font-mono text-sm text-ink"
        />
      </div>
      <div>
        <label className="font-mono text-xs uppercase tracking-wider text-muted">
          Tags (comma separated)
        </label>
        <input
          name="tags"
          defaultValue={post?.tags.join(", ") ?? ""}
          placeholder="ai, backend, go, fun fact"
          className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink"
        />
      </div>
      <label className="flex items-center gap-2 font-mono text-sm text-ink-soft">
        <input type="checkbox" name="published" defaultChecked={post?.published ?? false} className="h-4 w-4" />
        Published (unchecked = draft, hidden from the journal)
      </label>
      <button
        type="submit"
        className="rounded-full bg-moss px-5 py-2 font-mono text-sm text-paper hover:bg-moss-dark"
      >
        {isNew ? "Create post" : "Save changes"}
      </button>
    </form>
  );
}
