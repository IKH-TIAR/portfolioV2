import Link from "next/link";
import { logout } from "@/lib/adminActions";

const links = [
  { href: "/admin/projects", label: "projects" },
  { href: "/admin/posts", label: "posts" },
  { href: "/admin/photos", label: "photos" },
  { href: "/admin/messages", label: "messages" },
  { href: "/admin/comments", label: "comments" },
  { href: "/admin/settings", label: "settings" },
];

export default function AdminNav({ active }: { active: string }) {
  return (
    <nav className="mb-8 flex flex-wrap items-center gap-1 border-b border-line pb-4">
      {links.map((l) => (
        <Link
          key={l.href}
          href={l.href}
          className={
            active === l.href
              ? "rounded-full bg-moss px-3 py-1 font-mono text-sm text-paper"
              : "rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-ink-soft hover:bg-moss hover:text-paper"
          }
        >
          {l.label}
        </Link>
      ))}
      <form action={logout} className="ml-auto">
        <button
          type="submit"
          className="rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-muted hover:bg-ink hover:text-paper"
        >
          sign out
        </button>
      </form>
    </nav>
  );
}
