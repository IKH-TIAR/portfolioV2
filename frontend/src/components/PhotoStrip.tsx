"use client";

import { useState, useEffect, useCallback } from "react";
import { formatDate, type Photo } from "@/lib/api";

function PhotoCaption({ photo }: { photo: Photo }) {
  const meta = [photo.location, photo.takenAt ? formatDate(photo.takenAt) : ""]
    .filter(Boolean)
    .join(" · ");
  return (
    <>
      {photo.caption ? (
        <p className="font-serif text-base text-ink">{photo.caption}</p>
      ) : null}
      {meta ? (
        <p className="mt-0.5 font-mono text-xs text-muted">{meta}</p>
      ) : null}
    </>
  );
}

export default function PhotoStrip({ photos }: { photos: Photo[] }) {
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const close = useCallback(() => setOpenIndex(null), []);
  const prev = useCallback(() => setOpenIndex((i) => (i === null ? null : (i - 1 + photos.length) % photos.length)), []);
  const next = useCallback(() => setOpenIndex((i) => (i === null ? null : (i + 1) % photos.length)), []);

  useEffect(() => {
    if (openIndex === null) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") close();
      else if (e.key === "ArrowLeft") prev();
      else if (e.key === "ArrowRight") next();
    };
    window.addEventListener("keydown", onKey);
    return () => window.removeEventListener("keydown", onKey);
  }, [openIndex, close, prev, next]);

  if (photos.length === 0) return null;

  const active = openIndex !== null ? photos[openIndex] : null;

  return (
    <>
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {photos.map((photo, i) => (
          <figure
            key={photo.id}
            className="group cursor-zoom-in overflow-hidden rounded-2xl bg-paper-deep"
            onClick={() => setOpenIndex(i)}
          >
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={photo.url}
              alt={photo.caption || "photo"}
              className="aspect-[4/3] w-full object-cover transition-transform duration-500 group-hover:scale-[1.02]"
              loading="lazy"
            />
            <figcaption className="px-4 py-3">
              <PhotoCaption photo={photo} />
            </figcaption>
          </figure>
        ))}
      </div>

      {active ? (
        <div
          className="fixed inset-0 z-50 flex items-center justify-center bg-ink/80 p-4"
          onClick={close}
        >
          <div
            className="relative max-w-5xl"
            onClick={(e) => e.stopPropagation()}
          >
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img
              src={active.url}
              alt={active.caption || "photo"}
              className="max-h-[80vh] w-full rounded-lg object-contain"
            />
            <button
              onClick={prev}
              className="absolute top-1/2 left-2 -translate-y-1/2 rounded-full border border-line bg-paper-deep/70 p-2 text-ink hover:bg-moss hover:text-paper"
              aria-label="Previous"
            >
              ‹
            </button>
            <button
              onClick={next}
              className="absolute top-1/2 right-2 -translate-y-1/2 rounded-full border border-line bg-paper-deep/70 p-2 text-ink hover:bg-moss hover:text-paper"
              aria-label="Next"
            >
              ›
            </button>
            <button
              onClick={close}
              className="absolute top-2 right-2 rounded-full border border-line bg-paper-deep/70 p-1 text-ink hover:bg-moss hover:text-paper"
              aria-label="Close"
            >
              ✕
            </button>
            <div className="mt-3 px-1">
              <PhotoCaption photo={active} />
            </div>
            <p className="mt-2 font-mono text-center text-xs text-ink-soft/70">
              {openIndex! + 1} / {photos.length}
            </p>
          </div>
        </div>
      ) : null}
    </>
  );
}
