import Link from "next/link";
import {
  getPhotos,
  getPosts,
  getProjects,
  getSite,
} from "@/lib/api";
import JournalList from "@/components/JournalList";
import PhotoStrip from "@/components/PhotoStrip";
import ProjectCard from "@/components/ProjectCard";
import SectionHeading from "@/components/SectionHeading";
import ContactForm from "@/components/ContactForm";

const stack = [
  "python",
  "langchain",
  "langgraph",
  "n8n",
  "fastapi",
  "go",
  "postgres",
  "next.js",
];

export default async function Home() {
  const [projects, posts, photos, site] = await Promise.all([
    getProjects(),
    getPosts(3),
    getPhotos(),
    getSite(),
  ]);
  const featured = projects.filter((p) => p.featured).slice(0, 3);
  const heroTitle = site.hero?.title ?? "Ikhtiar";
  const heroTagline =
    site.hero?.tagline ??
    "What I cannot build, I do not understand.";

  return (
    <>
      <section className="border-b border-line">
        <div className="mx-auto w-full max-w-5xl px-5 pb-20 pt-16 sm:px-8 sm:pb-28 sm:pt-24">
          <p className="font-mono text-sm text-moss">~/whoami</p>
          <h1 className="mt-4 max-w-3xl font-serif text-5xl leading-[1.05] text-ink sm:text-7xl">
            {heroTitle}
          </h1>
          <p className="mt-6 max-w-xl text-lg leading-relaxed text-muted">
            {heroTagline}
          </p>
          <p className="mt-8 inline-flex items-center gap-2 rounded-full border border-line bg-paper-deep/60 px-4 py-2 font-mono text-sm text-ink-soft">
            <span aria-hidden="true" className="h-2 w-2 rounded-full bg-moss" />
            currently: jr ai developer · learning go
          </p>
        </div>
      </section>

      <section className="border-b border-line">
        <div className="mx-auto w-full max-w-5xl px-5 py-20 sm:px-8">
          <SectionHeading path="~/projects" title="Featured work">
            <Link
              href="/projects"
              className="font-mono text-sm text-moss-dark underline underline-offset-3 hover:text-moss"
            >
              all projects →
            </Link>
          </SectionHeading>
          {featured.length > 0 ? (
            <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
              {featured.map((project) => (
                <ProjectCard key={project.id} project={project} />
              ))}
            </div>
          ) : (
            <p className="font-mono text-sm text-muted">
              no projects published yet. The workshop is open.
            </p>
          )}
        </div>
      </section>

      <section className="border-b border-line">
        <div className="mx-auto w-full max-w-5xl px-5 py-20 sm:px-8">
          <SectionHeading path="~/journal" title="Today I learned">
            <Link
              href="/journal"
              className="font-mono text-sm text-moss-dark underline underline-offset-3 hover:text-moss"
            >
              read the journal →
            </Link>
          </SectionHeading>
          <JournalList posts={posts} showExcerpt={false} />
        </div>
      </section>

      {photos.length > 0 ? (
        <section className="border-b border-line">
          <div className="mx-auto w-full max-w-5xl px-5 py-20 sm:px-8">
            <SectionHeading path="~/photos" title="Field notes">
              <Link
                href="/photos"
                className="font-mono text-sm text-moss-dark underline underline-offset-3 hover:text-moss"
              >
                all photos →
              </Link>
            </SectionHeading>
            <PhotoStrip photos={photos.slice(0, 6)} />
          </div>
        </section>
      ) : null}

      <section className="border-b border-line">
        <div className="mx-auto grid w-full max-w-5xl gap-12 px-5 py-20 sm:px-8 lg:grid-cols-2">
          <div>
            <SectionHeading path="~/about" title="About me" />
            <div className="space-y-4 leading-relaxed text-ink-soft">
              <p>
                I&apos;m a junior AI developer with a stubborn habit of
                automating anything that happens twice. At work I build
                AI-powered systems: agents, RAG pipelines, and the glue
                between them. After hours I tinker, journal what I learn,
                and occasionally go outside with a camera.
              </p>
              <p>
                Right now I&apos;m going deeper into Go, system design, and
                microservices, because the bigger the automation, the more
                it needs to be built like a real system.
              </p>
            </div>
            <ul className="mt-8 flex flex-wrap gap-2">
              {stack.map((t) => (
                <li
                  key={t}
                  className="rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-ink-soft"
                >
                  {t}
                </li>
              ))}
            </ul>
          </div>
          <div>
            <SectionHeading path="~/now" title="What I'm up to" />
            <ul className="space-y-5">
              <li className="flex gap-4">
                <span className="mt-1.5 h-2 w-2 shrink-0 rounded-full bg-moss" />
                <p className="text-ink-soft">
                  <span className="font-medium text-ink">Learning Go.</span>{" "}
                  The portfolio you&apos;re looking at is built with it.
                </p>
              </li>
              <li className="flex gap-4">
                <span className="mt-1.5 h-2 w-2 shrink-0 rounded-full bg-moss" />
                <p className="text-ink-soft">
                  <span className="font-medium text-ink">Building agents.</span>{" "}
                  SQL agents, RAG chatbots, and n8n automations for small
                  businesses.
                </p>
              </li>
              <li className="flex gap-4">
                <span className="mt-1.5 h-2 w-2 shrink-0 rounded-full bg-moss" />
                <p className="text-ink-soft">
                  <span className="font-medium text-ink">Taking photos.</span>{" "}
                  Whatever nature happens to be around.
                </p>
              </li>
            </ul>
          </div>
        </div>
      </section>

      <section id="contact">
        <div className="mx-auto w-full max-w-5xl px-5 py-20 sm:px-8">
          <SectionHeading path="~/contact" title="Say hello" />
          <div className="grid gap-10 lg:grid-cols-2">
            <div>
              <p className="leading-relaxed text-muted">
                Have a project, a question, or just want to talk shop about
                agents and automations? My inbox is open.
              </p>
              <p className="mt-6 font-mono text-sm">
                <a
                  href="https://github.com/IKH-TIAR"
                  target="_blank"
                  rel="noreferrer"
                  className="text-moss-dark underline underline-offset-3 hover:text-moss"
                >
                  github.com/IKH-TIAR
                </a>
              </p>
            </div>
            <ContactForm />
          </div>
        </div>
      </section>
    </>
  );
}
