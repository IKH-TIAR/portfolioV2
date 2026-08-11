import { notFound } from "next/navigation";
import Link from "next/link";
import AdminNav from "@/components/AdminNav";
import ProjectForm from "@/components/admin/ProjectForm";
import { adminFetch } from "@/lib/admin";

type Params = { id: string };

type Project = {
  id: number;
  slug: string;
  title: string;
  tagline: string;
  description: string;
  tech: string[];
  repoUrl: string;
  demoUrl: string;
  imageUrl: string;
  featured: boolean;
  sortOrder: number;
  published: boolean;
};

export default async function EditProjectPage({ params }: { params: Promise<Params> }) {
  const { id } = await params;
  let project: Project;
  try {
    project = await adminFetch<Project>(`/api/admin/projects/${id}`);
  } catch {
    notFound();
  }
  return (
    <div>
      <AdminNav active="projects" />
      <Link href="/admin/projects" className="font-mono text-sm text-moss hover:text-ink">
        ← back to projects
      </Link>
      <h1 className="mt-2 font-serif text-3xl text-ink">Edit: {project.title}</h1>
      <ProjectForm project={project} />
    </div>
  );
}
