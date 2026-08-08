import type { Metadata } from "next";
import { getProjects } from "@/lib/api";
import { ProjectLink } from "@/components/ProjectCard";
import SectionHeading from "@/components/SectionHeading";

export const metadata: Metadata = {
  title: "Projects",
  description:
    "AI agents, automations, and backend projects: built with LangChain, LangGraph, FastAPI, n8n, and Go.",
};

export default async function ProjectsPage() {
  const projects = await getProjects();

  return (
    <div className="mx-auto w-full max-w-3xl px-5 py-16 sm:px-8">
      <SectionHeading path="~/projects" title="Projects" />
      <p className="-mt-4 mb-12 max-w-xl leading-relaxed text-muted">
        Things I&apos;ve built: agents that answer questions, bots that run
        small businesses, and the boring-but-necessary systems behind them.
      </p>
      {projects.length > 0 ? (
        <div>
          {projects.map((project) => (
            <ProjectLink key={project.id} project={project} />
          ))}
        </div>
      ) : (
        <p className="font-mono text-sm text-muted">
          no projects published yet. The workshop is open.
        </p>
      )}
    </div>
  );
}
