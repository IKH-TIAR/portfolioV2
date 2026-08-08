import type { Metadata } from "next";
import { getPosts } from "@/lib/api";
import JournalList from "@/components/JournalList";
import SectionHeading from "@/components/SectionHeading";

export const metadata: Metadata = {
  title: "Journal",
  description:
    "Today I learned: a public learning journal of a junior AI developer.",
};

export default async function JournalPage() {
  const posts = await getPosts();

  return (
    <div className="mx-auto w-full max-w-3xl px-5 py-16 sm:px-8">
      <SectionHeading path="~/journal" title="The journal" />
      <p className="-mt-4 mb-12 max-w-xl leading-relaxed text-muted">
        A public learning log: what I built, what broke, what I figured out.
        Written for future me, readable by anyone.
      </p>
      <JournalList posts={posts} />
    </div>
  );
}
