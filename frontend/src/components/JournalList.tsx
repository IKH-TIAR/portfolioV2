import Link from "next/link";
import { excerpt, formatDate, type Post } from "@/lib/api";

export default function JournalList({
  posts,
  showExcerpt = true,
}: {
  posts: Post[];
  showExcerpt?: boolean;
}) {
  if (posts.length === 0) {
    return (
      <div className="rounded-2xl border border-dashed border-line p-10 text-center font-mono text-sm text-muted">
        the journal is quiet for now. Check back soon.
      </div>
    );
  }

  return (
    <ol className="relative space-y-10 before:absolute before:bottom-2 before:left-[5px] before:top-2 before:w-px before:bg-line">
      {posts.map((post) => (
        <li key={post.id} className="relative pl-10">
          <span
            aria-hidden="true"
            className="absolute left-0 top-1.5 h-[11px] w-[11px] rounded-full border-2 border-moss bg-paper"
          />
          <p className="font-mono text-xs uppercase tracking-wider text-muted">
            {formatDate(post.createdAt)}
          </p>
          <Link
            href={`/journal/${post.slug}`}
            className="group mt-1 block font-serif text-2xl leading-snug text-ink transition-colors hover:text-moss-dark"
          >
            {post.title}
          </Link>
          {showExcerpt && excerpt(post.body) ? (
            <p className="mt-2 max-w-xl text-sm leading-relaxed text-muted">
              {excerpt(post.body, 160)}
            </p>
          ) : null}
        </li>
      ))}
    </ol>
  );
}
