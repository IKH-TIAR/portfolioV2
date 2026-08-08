import Link from "next/link";
import type { Project } from "@/lib/api";

export default function ProjectCard({ project }: { project: Project }) {
  return (
    <article className="flex flex-col rounded-2xl border border-line bg-paper-deep/50 p-6 transition-colors hover:border-moss/40">
      {project.imageUrl ? (
        // eslint-disable-next-line @next/next/no-img-element
        <img
          src={project.imageUrl}
          alt={project.title}
          className="mb-5 aspect-[16/9] w-full rounded-xl object-cover"
        />
      ) : null}
      <h3 className="font-serif text-xl text-ink">{project.title}</h3>
      {project.tagline ? (
        <p className="mt-1 text-muted">{project.tagline}</p>
      ) : null}
      <p className="mt-3 line-clamp-3 text-sm leading-relaxed text-ink-soft">
        {project.description}
      </p>
      {project.tech.length > 0 ? (
        <ul className="mt-4 flex flex-wrap gap-2">
          {project.tech.map((t) => (
            <li
              key={t}
              className="rounded-full border border-line bg-paper px-2.5 py-0.5 font-mono text-xs text-muted"
            >
              {t}
            </li>
          ))}
        </ul>
      ) : null}
      <div className="mt-5 flex gap-5 border-t border-line pt-4 font-mono text-sm">
        {project.repoUrl ? (
          <a
            href={project.repoUrl}
            target="_blank"
            rel="noreferrer"
            className="text-moss-dark underline underline-offset-3 hover:text-moss"
          >
            source
          </a>
        ) : null}
        {project.demoUrl ? (
          <a
            href={project.demoUrl}
            target="_blank"
            rel="noreferrer"
            className="text-moss-dark underline underline-offset-3 hover:text-moss"
          >
            live
          </a>
        ) : null}
      </div>
    </article>
  );
}

export function ProjectLink({ project }: { project: Project }) {
  return (
    <Link
      href={`/projects#${project.slug}`}
      className="group flex items-baseline justify-between gap-4 border-b border-line py-5 transition-colors hover:border-moss/40"
    >
      <div className="min-w-0">
        <h3 className="font-serif text-xl text-ink transition-colors group-hover:text-moss-dark sm:text-2xl">
          {project.title}
        </h3>
        {project.tagline ? <p className="mt-1 text-muted">{project.tagline}</p> : null}
      </div>
      <span className="shrink-0 font-mono text-sm text-moss">→</span>
    </Link>
  );
}
