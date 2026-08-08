import { notFound } from "next/navigation";
import Link from "next/link";
import { getProject, formatDate, type Project } from "@/lib/api";

interface Params {
  slug: string;
}

export async function generateMetadata({ params }: { params: Promise<Params> }) {
  const { slug } = await params;
  const project = await getProject(slug).catch(() => null);
  if (!project) return {};
  return {
    title: project.title,
    description: project.tagline || project.description,
    openGraph: {
      title: project.title,
      description: project.tagline || project.description,
    },
  };
}

export default async function ProjectPage({ params }: { params: Promise<Params> }) {
  const { slug } = await params;
  let project: Project;
  try {
    project = await getProject(slug);
  } catch {
    notFound();
  }

  const links: { label: string; href: string | null }[] = [
    { label: "Repository", href: project.repoUrl },
    { label: "Live demo", href: project.demoUrl },
  ];

  return (
    <article className="border-b border-line">
      <div className="mx-auto w-full max-w-5xl px-5 py-20 sm:px-8">
        <p className="font-mono text-sm text-moss">~/projects/{project.slug}</p>
        <h1 className="mt-4 font-serif text-4xl text-ink sm:text-5xl">{project.title}</h1>
        {project.tagline ? (
          <p className="mt-4 max-w-2xl text-xl text-ink-soft">{project.tagline}</p>
        ) : null}

        {project.imageUrl ? (
          <img
            alt={`${project.title} preview`}
            src={project.imageUrl}
            className="mt-8 w-full rounded-lg border border-line object-cover"
            loading="lazy"
          />
        ) : null}

        <div
          className="mt-8 max-w-prose text-ink-soft"
          dangerouslySetInnerHTML={{ __html: project.description }}
        />

        {project.tech.length > 0 ? (
          <ul className="mt-8 flex flex-wrap gap-2">
            {project.tech.map((t) => (
              <li
                key={t}
                className="rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-ink-soft"
              >
                {t}
              </li>
            ))}
          </ul>
        ) : null}

        <div className="mt-10 flex flex-wrap gap-4">
          {links.map((l) =>
            l.href ? (
              <Link
                key={l.label}
                href={l.href}
                target="_blank"
                rel="noreferrer"
                className="rounded-full border border-line bg-paper-deep/60 px-4 py-2 font-mono text-sm text-ink-soft hover:bg-moss hover:text-paper"
              >
                {l.label}
              </Link>
            ) : null
          )}
          <Link
            href="/projects"
            className="rounded-full border border-line bg-paper-deep/60 px-4 py-2 font-mono text-sm text-ink-soft hover:bg-moss hover:text-paper"
          >
            All projects
          </Link>
        </div>

        <p className="mt-6 font-mono text-xs text-ink-soft/60">
          Last updated {formatDate(project.updatedAt)}
        </p>
      </div>
    </article>
  );
}
