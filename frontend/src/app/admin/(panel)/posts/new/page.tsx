import Link from "next/link";
import AdminNav from "@/components/AdminNav";
import PostForm from "@/components/admin/PostForm";

export default function NewPostPage() {
  return (
    <div>
      <AdminNav active="posts" />
      <Link href="/admin/posts" className="font-mono text-sm text-moss hover:text-ink">
        ← back to posts
      </Link>
      <h1 className="mt-2 font-serif text-3xl text-ink">New post</h1>
      <PostForm />
    </div>
  );
}
