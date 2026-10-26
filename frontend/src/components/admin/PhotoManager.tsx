"use client";

import { useActionState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { createPhoto, deletePhoto, uploadAndCreatePhoto } from "@/lib/adminActions";

type Photo = {
  id: number;
  url: string;
  caption: string;
  location: string;
  takenAt: string | null;
  published: boolean;
  createdAt: string;
};

type State = { error?: string; ok?: boolean };

export default function PhotoManager({ photos }: { photos: Photo[] }) {
  const router = useRouter();
  const [uploadState, uploadAction, uploadPending] = useActionState<State, FormData>(
    (_, fd) => uploadAndCreatePhoto(fd),
    {}
  );
  const [urlState, urlAction] = useActionState<State, FormData>((_, fd) => createPhoto(fd), {});
  useEffect(() => {
    if (uploadState.ok || urlState.ok) {
      router.refresh();
    }
  }, [uploadState.ok, urlState.ok, router]);
  return (
    <div className="mt-6 space-y-8">
      <form action={uploadAction} className="space-y-4 rounded-2xl border border-line bg-paper-deep/40 p-6">
        {uploadState.error ? <ErrorBox message={uploadState.error} /> : null}
        <h2 className="font-serif text-xl text-ink">Upload photos</h2>
        <input
          type="file"
          name="file"
          multiple
          required
          accept="image/jpeg,image/png,image/gif,image/webp,image/avif"
          className="mt-1 w-full text-sm text-ink-soft file:mr-3 file:rounded-full file:border-0 file:bg-moss file:px-4 file:py-2 file:text-sm file:text-paper"
        />
        <p className="font-mono text-xs text-muted">
          select one or more images (up to 20, 10MB each) — they share the
          caption, location and date below
        </p>
        <div className="grid gap-4 sm:grid-cols-3">
          <div>
            <label className="font-mono text-xs uppercase tracking-wider text-muted">Caption</label>
            <input name="caption" className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink" />
          </div>
          <div>
            <label className="font-mono text-xs uppercase tracking-wider text-muted">Location</label>
            <input name="location" className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink" />
          </div>
          <div>
            <label className="font-mono text-xs uppercase tracking-wider text-muted">Taken</label>
            <input name="takenAt" type="date" className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink" />
          </div>
        </div>
        <div className="flex flex-wrap items-center gap-6">
          <label className="flex items-center gap-2 font-mono text-sm text-ink-soft">
            <input type="checkbox" name="published" defaultChecked className="h-4 w-4" />
            Published
          </label>
          <button
            type="submit"
            disabled={uploadPending}
            className="rounded-full bg-moss px-5 py-2 font-mono text-sm text-paper hover:bg-moss-dark disabled:opacity-50"
          >
            {uploadPending ? "Uploading…" : "Upload photos"}
          </button>
        </div>
      </form>

      <form action={urlAction} className="space-y-4 rounded-2xl border border-dashed border-line p-6">
        {urlState.error ? <ErrorBox message={urlState.error} /> : null}
        <h2 className="font-serif text-xl text-ink">Add by URL</h2>
        <div className="grid gap-4 sm:grid-cols-3">
          <div>
            <label className="font-mono text-xs uppercase tracking-wider text-muted">Image URL</label>
            <input name="url" required placeholder="https://..." className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink" />
          </div>
          <div>
            <label className="font-mono text-xs uppercase tracking-wider text-muted">Caption</label>
            <input name="caption" className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink" />
          </div>
          <div>
            <label className="font-mono text-xs uppercase tracking-wider text-muted">Location</label>
            <input name="location" className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink" />
          </div>
        </div>
        <button
          type="submit"
          className="rounded-full border border-line bg-paper-deep/60 px-5 py-2 font-mono text-sm text-ink-soft hover:bg-moss hover:text-paper"
        >
          Add photo
        </button>
      </form>

      {photos.length > 0 ? (
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {photos.map((p) => (
            <div key={p.id} className="overflow-hidden rounded-2xl border border-line bg-paper-deep">
              {/* eslint-disable-next-line @next/next/no-img-element */}
              <img src={p.url} alt={p.caption || "photo"} className="aspect-[4/3] w-full object-cover" />
              <div className="px-4 py-3">
                <p className="truncate font-serif text-sm text-ink">{p.caption || "(no caption)"}</p>
                <p className="mt-0.5 truncate font-mono text-xs text-muted">
                  {p.location || "no location"} {p.published ? "· live" : "· draft"}
                </p>
                <form
                  action={() => {
                    void deletePhoto(p.id);
                  }}
                  className="mt-2"
                >
                  <button
                    type="submit"
                    className="font-mono text-xs text-muted underline underline-offset-2 hover:text-ink"
                  >
                    delete
                  </button>
                </form>
              </div>
            </div>
          ))}
        </div>
      ) : (
        <p className="font-mono text-sm text-muted">no photos yet</p>
      )}
    </div>
  );
}

function ErrorBox({ message }: { message: string }) {
  return (
    <p className="rounded-xl border border-red-300 bg-red-50 px-4 py-2 font-mono text-sm text-red-700">
      {message}
    </p>
  );
}
