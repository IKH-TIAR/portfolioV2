"use client";

import { useActionState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { setSiteSetting } from "@/lib/adminActions";

type State = { error?: string; ok?: boolean };

type Props = {
  hero?: { title?: string; tagline?: string };
  about?: { body?: string };
};

export default function SettingsForms({ hero, about }: Props) {
  const router = useRouter();
  const [heroState, heroAction] = useActionState<State, FormData>(
    (_, fd) =>
      setSiteSetting("hero", {
        title: fd.get("heroTitle") as string,
        tagline: fd.get("heroTagline") as string,
      }),
    {}
  );
  const [aboutState, aboutAction] = useActionState<State, FormData>(
    (_, fd) => setSiteSetting("about", { body: fd.get("aboutBody") as string }),
    {}
  );
  useEffect(() => {
    if (heroState.ok || aboutState.ok) {
      router.refresh();
    }
  }, [heroState.ok, aboutState.ok, router]);

  return (
    <>
      <form action={heroAction} className="mt-6 space-y-4 rounded-2xl border border-line bg-paper-deep/40 p-6">
        {heroState.error ? <ErrorBox message={heroState.error} /> : null}
        <h2 className="font-serif text-xl text-ink">Hero</h2>
        <div>
          <label className="font-mono text-xs uppercase tracking-wider text-muted">Title</label>
          <input
            name="heroTitle"
            defaultValue={hero?.title ?? ""}
            className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink"
          />
        </div>
        <div>
          <label className="font-mono text-xs uppercase tracking-wider text-muted">Tagline</label>
          <textarea
            name="heroTagline"
            rows={3}
            defaultValue={hero?.tagline ?? ""}
            className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink"
          />
        </div>
        <button
          type="submit"
          className="rounded-full bg-moss px-5 py-2 font-mono text-sm text-paper hover:bg-moss-dark"
        >
          Save hero
        </button>
      </form>

      <form action={aboutAction} className="mt-6 space-y-4 rounded-2xl border border-line bg-paper-deep/40 p-6">
        {aboutState.error ? <ErrorBox message={aboutState.error} /> : null}
        <h2 className="font-serif text-xl text-ink">About</h2>
        <div>
          <label className="font-mono text-xs uppercase tracking-wider text-muted">
            About body (markdown)
          </label>
          <textarea
            name="aboutBody"
            rows={8}
            defaultValue={about?.body ?? ""}
            className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 font-mono text-sm text-ink"
          />
        </div>
        <button
          type="submit"
          className="rounded-full bg-moss px-5 py-2 font-mono text-sm text-paper hover:bg-moss-dark"
        >
          Save about
        </button>
      </form>
    </>
  );
}

function ErrorBox({ message }: { message: string }) {
  return (
    <p className="rounded-xl border border-red-300 bg-red-50 px-4 py-2 font-mono text-sm text-red-700">
      {message}
    </p>
  );
}
