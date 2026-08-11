import type { Metadata } from "next";
import Link from "next/link";
import { notFound } from "next/navigation";
import ReactMarkdown from "react-markdown";
import { formatDate, getComments, getPost } from "@/lib/api";
import CommentSection from "@/components/CommentSection";

export async function generateMetadata({
  params,
}: {
  params: Promise<{ slug: string }>;
}): Promise<Metadata> {
  const { slug } = await params;
  try {
    const post = await getPost(slug);
    return { title: post.title, description: post.body.slice(0, 160) };
  } catch {
    return { title: "Journal" };
  }
}

export default async function JournalPost({
  params,
}: {
  params: Promise<{ slug: string }>;
}) {
  const { slug } = await params;

  let post;
  try {
    post = await getPost(slug);
  } catch {
    notFound();
  }

  const comments = await getComments(slug);

  return (
    <article className="mx-auto w-full max-w-3xl px-5 py-16 sm:px-8">
       <p className="font-mono text-xs uppercase tracking-wider text-muted">
         ~/journal/{post.slug}
       </p>
       <Link
         href="/journal"
         className="mt-4 inline-flex items-center gap-1 font-mono text-sm text-moss hover:text-ink"
       >
         ← Back to journal
       </Link>
       <h1 className="mt-4 font-serif text-4xl leading-tight text-ink sm:text-5xl">
        {post.title}
      </h1>
      <p className="mt-4 font-mono text-sm text-muted">
        {formatDate(post.createdAt)}
      </p>
      <div className="post-body mt-10">
        <ReactMarkdown>{post.body}</ReactMarkdown>
      </div>
      <CommentSection slug={post.slug} initial={comments} />
    </article>
  );
}
