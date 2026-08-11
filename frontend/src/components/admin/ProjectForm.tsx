"use client";

import { useActionState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { createProject, updateProject } from "@/lib/adminActions";

type Props = {
  project?: {
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
    published: boolean;
    sortOrder: number;
  };
};

type State = { error?: string; ok?: boolean };

export default function ProjectForm({ project }: Props) {
  const router = useRouter();
  const isNew = !project;
  const [state, formAction] = useActionState<State, FormData>(
    isNew ? (_, fd) => createProject(fd) : project ? (_, fd) => updateProject(project.id, fd) : (_, fd) => createProject(fd),
    {}
  );
  useEffect(() => {
    if (state.ok) {
      router.push("/admin/projects");
    }
  }, [state.ok, router]);
  return (
    <form action={formAction} className="mt-6 space-y-5">
      {state.error ? (
        <p className="rounded-xl border border-red-300 bg-red-50 px-4 py-2 font-mono text-sm text-red-700">
          {state.error}
        </p>
      ) : null}
      <input type="hidden" name="sortOrder" value={project?.sortOrder ?? 0} />
      <div className="grid gap-5 sm:grid-cols-2">
        <Field label="Title" name="title" value={project?.title} required />
        <Field label="Slug" name="slug" value={project?.slug} required />
      </div>
      <Field label="Tagline" name="tagline" value={project?.tagline} />
      <div>
        <Label>Description (HTML or markdown)</Label>
        <textarea
          name="description"
          rows={6}
          defaultValue={project?.description}
          className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 font-mono text-sm text-ink"
        />
      </div>
      <Field label="Tech (comma separated)" name="tech" value={project?.tech.join(", ")} />
      <div className="grid gap-5 sm:grid-cols-2">
        <Field label="Repo URL" name="repoUrl" value={project?.repoUrl} />
        <Field label="Demo URL" name="demoUrl" value={project?.demoUrl} />
      </div>
      <Field label="Image URL" name="imageUrl" value={project?.imageUrl} />
      <div className="flex flex-wrap gap-6">
        <Checkbox label="Featured" name="featured" checked={project?.featured ?? false} />
        <Checkbox label="Published" name="published" checked={project?.published ?? false} />
      </div>
      <button
        type="submit"
        className="rounded-full bg-moss px-5 py-2 font-mono text-sm text-paper hover:bg-moss-dark"
      >
        {isNew ? "Create project" : "Save changes"}
      </button>
    </form>
  );
}

function Field({
  label,
  name,
  value,
  required,
}: {
  label: string;
  name: string;
  value?: string;
  required?: boolean;
}) {
  return (
    <div>
      <Label>{label}</Label>
      <input
        name={name}
        defaultValue={value}
        required={required}
        className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink"
      />
    </div>
  );
}

function Label({ children }: { children: React.ReactNode }) {
  return (
    <label className="font-mono text-xs uppercase tracking-wider text-muted">{children}</label>
  );
}

function Checkbox({ label, name, checked }: { label: string; name: string; checked: boolean }) {
  return (
    <label className="flex items-center gap-2 font-mono text-sm text-ink-soft">
      <input type="checkbox" name={name} defaultChecked={checked} className="h-4 w-4" />
      {label}
    </label>
  );
}
