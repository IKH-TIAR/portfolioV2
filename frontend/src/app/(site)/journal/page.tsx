import type { Metadata } from "next";
import Link from "next/link";
import { getPosts } from "@/lib/api";
import JournalList from "@/components/JournalList";
import SectionHeading from "@/components/SectionHeading";

export const metadata: Metadata = {
  title: "Journal",
  description:
    "Today I learned: a public learning journal of a junior AI developer.",
};

export default async function JournalPage({
  searchParams,
}: {
  searchParams: Promise<{ tag?: string | string[] }>;
}) {
  const { tag } = await searchParams;
  const activeTag =
    typeof tag === "string" ? tag.trim().toLowerCase() : "";
  const posts = await getPosts();

  const allTags = [
    ...new Set(posts.flatMap((p) => p.tags ?? [])),
  ].sort();

  const filtered = activeTag
    ? posts.filter((p) => (p.tags ?? []).includes(activeTag))
    : posts;

  return (
    <div className="mx-auto w-full max-w-3xl px-5 py-16 sm:px-8">
      <SectionHeading path="~/journal" title="The journal" />
      <p className="-mt-4 mb-12 max-w-xl leading-relaxed text-muted">
        A public learning log: what I built, what broke, what I figured out.
        Written for future me, readable by anyone.
      </p>

      {allTags.length > 0 ? (
        <div className="mb-8 flex flex-wrap items-center gap-2">
          <Link
            href="/journal"
            className={
              activeTag
                ? "rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-ink-soft transition-colors hover:bg-moss hover:text-paper"
                : "rounded-full bg-moss px-3 py-1 font-mono text-sm text-paper"
            }
          >
            all
          </Link>
          {allTags.map((t) => (
            <Link
              key={t}
              href={`/journal?tag=${encodeURIComponent(t)}`}
              className={
                activeTag === t
                  ? "rounded-full bg-moss px-3 py-1 font-mono text-sm text-paper"
                  : "rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-ink-soft transition-colors hover:bg-moss hover:text-paper"
              }
            >
              {t}
            </Link>
          ))}
        </div>
      ) : null}

      {activeTag ? (
        <p className="mb-8 font-mono text-sm text-muted">
          {filtered.length} {filtered.length === 1 ? "post" : "posts"} tagged{" "}
          “{activeTag}”
        </p>
      ) : null}

      <JournalList posts={filtered} />
    </div>
  );
}
