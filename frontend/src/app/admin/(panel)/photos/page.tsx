import AdminNav from "@/components/AdminNav";
import PhotoManager from "@/components/admin/PhotoManager";
import { adminFetch } from "@/lib/admin";

type Photo = {
  id: number;
  url: string;
  caption: string;
  location: string;
  takenAt: string | null;
  published: boolean;
  createdAt: string;
};

export default async function AdminPhotosPage() {
  const photos = await adminFetch<Photo[]>("/api/admin/photos");
  return (
    <div>
      <AdminNav active="photos" />
      <h1 className="font-serif text-3xl text-ink">Photos</h1>
      <PhotoManager photos={photos} />
    </div>
  );
}
