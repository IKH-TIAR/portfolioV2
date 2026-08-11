import AdminNav from "@/components/AdminNav";
import SettingsForms from "@/components/admin/SettingsForms";
import { adminFetch } from "@/lib/admin";

type Settings = Record<string, { value: unknown }>;

export default async function AdminSettingsPage() {
  const settings = await adminFetch<Settings>("/api/admin/site-settings");
  const hero = (settings.hero?.value ?? {}) as { title?: string; tagline?: string };
  const about = (settings.about?.value ?? {}) as { body?: string };

  return (
    <div>
      <AdminNav active="settings" />
      <h1 className="font-serif text-3xl text-ink">Site settings</h1>
      <SettingsForms hero={hero} about={about} />
    </div>
  );
}
