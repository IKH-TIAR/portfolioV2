import type { ReactNode } from "react";
import { requireAuth } from "@/lib/admin";

export default async function AdminPanelLayout({ children }: { children: ReactNode }) {
  await requireAuth();
  return <>{children}</>;
}
