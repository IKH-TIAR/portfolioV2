import Link from "next/link";
import AdminNav from "@/components/AdminNav";
import ProjectForm from "@/components/admin/ProjectForm";

export default function NewProjectPage() {
  return (
    <div>
      <AdminNav active="projects" />
      <Link href="/admin/projects" className="font-mono text-sm text-moss hover:text-ink">
        ← back to projects
      </Link>
      <h1 className="mt-2 font-serif text-3xl text-ink">New project</h1>
      <ProjectForm />
    </div>
  );
}
