import Link from "next/link";
import AdminNav from "@/components/AdminNav";
import PostForm from "@/components/admin/PostForm";
import { adminFetch } from "@/lib/admin";
import { notFound } from "next/navigation";

type Params = { id: string };

type Post = {
  id: number;
  slug: string;
  title: string;
  body: string;
  postType: "journal" | "note";
  published: boolean;
};

export default async function EditPostPage({ params }: { params: Promise<Params> }) {
  const { id } = await params;
  let post: Post;
  try {
    post = await adminFetch<Post>(`/api/admin/posts/${id}`);
  } catch {
    notFound();
  }
  return (
    <div>
      <AdminNav active="posts" />
      <Link href="/admin/posts" className="font-mono text-sm text-moss hover:text-ink">
        ← back to posts
      </Link>
      <h1 className="mt-2 font-serif text-3xl text-ink">Edit: {post.title}</h1>
      <PostForm post={post} />
    </div>
  );
}
