import { redirect } from "next/navigation";
import { cookies } from "next/headers";

export default async function AdminIndexPage() {
  const token = (await cookies()).get("admin_token")?.value || "";
  redirect(token ? "/admin/projects" : "/admin/login");
}
