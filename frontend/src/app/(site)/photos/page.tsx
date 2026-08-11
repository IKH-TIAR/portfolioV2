import type { Metadata } from "next";
import { getPhotos } from "@/lib/api";
import PhotoStrip from "@/components/PhotoStrip";
import SectionHeading from "@/components/SectionHeading";

export const metadata: Metadata = {
  title: "Photos",
  description: "Field notes: photographs from around the way.",
};

export default async function PhotosPage() {
  const photos = await getPhotos();

  return (
    <div className="mx-auto w-full max-w-5xl px-5 py-16 sm:px-8">
      <SectionHeading path="~/photos" title="Field notes" />
      <p className="-mt-4 mb-12 max-w-xl leading-relaxed text-muted">
        Nature, whenever it&apos;s around. No post-processing, no thesis,
        just where I was.
      </p>
      {photos.length > 0 ? (
        <PhotoStrip photos={photos} />
      ) : (
        <p className="font-mono text-sm text-muted">
          no photos yet. The camera is on the shelf.
        </p>
      )}
    </div>
  );
}
