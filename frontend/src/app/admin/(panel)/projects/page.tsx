import Link from "next/link";
import AdminNav from "@/components/AdminNav";
import { adminFetch } from "@/lib/admin";
import { deleteProject } from "@/lib/adminActions";

type Project = {
  id: number;
  slug: string;
  title: string;
  tagline: string;
  featured: boolean;
  published: boolean;
  sortOrder: number;
};

export default async function AdminProjectsPage() {
  const projects = await adminFetch<Project[]>("/api/admin/projects");
  return (
    <div>
      <AdminNav active="projects" />
      <div className="flex items-end justify-between">
        <h1 className="font-serif text-3xl text-ink">Projects</h1>
        <Link
          href="/admin/projects/new"
          className="rounded-full bg-moss px-4 py-2 font-mono text-sm text-paper hover:bg-moss-dark"
        >
          + new project
        </Link>
      </div>
      <ul className="mt-6 space-y-3">
        {projects.map((p) => (
          <li
            key={p.id}
            className="flex flex-wrap items-center gap-3 rounded-2xl border border-line bg-paper-deep/40 px-5 py-4"
          >
            <div className="min-w-0 flex-1">
              <p className="truncate font-serif text-lg text-ink">
                {p.title}
                {p.featured ? <span className="ml-2 font-mono text-xs text-moss">★ featured</span> : null}
              </p>
              <p className="truncate font-mono text-xs text-muted">
                {p.slug} {p.published ? "· live" : "· draft"} · order {p.sortOrder}
              </p>
            </div>
            <Link
              href={`/admin/projects/${p.id}`}
              className="rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-ink-soft hover:bg-moss hover:text-paper"
            >
              edit
            </Link>
            <form
              action={async () => {
                "use server";
                await deleteProject(p.id);
              }}
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
        {projects.length === 0 ? <li className="font-mono text-sm text-muted">no projects yet</li> : null}
      </ul>
    </div>
  );
}
