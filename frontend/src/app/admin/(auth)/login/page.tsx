import { redirect } from "next/navigation";
import Link from "next/link";
import { cookies } from "next/headers";
import { login } from "@/lib/admin";

export const dynamic = "force-dynamic";

export default async function AdminLoginPage({
  searchParams,
}: {
  searchParams: Promise<{ callbackUrl?: string; error?: string }>;
}) {
  const sp = await searchParams;
  const token = (await cookies()).get("admin_token")?.value;
  if (token) redirect(sp.callbackUrl || "/admin/projects");

  return (
    <div className="mx-auto w-full max-w-sm px-4 py-16">
      <h1 className="font-serif text-2xl text-ink">Admin login</h1>
      {sp.error ? <p className="mt-4 text-sm text-moss">Invalid credentials.</p> : null}
      <form
        action={async (formData: FormData) => {
          "use server";
          const email = String(formData.get("email") || "");
          const password = String(formData.get("password") || "");
          const callbackUrl = String(formData.get("callbackUrl") || "/admin/projects");
          try {
            const token = await login(email, password);
            (await cookies()).set("admin_token", token, {
              httpOnly: true,
              sameSite: "lax",
              path: "/admin",
              maxAge: 60 * 60 * 24,
            });
            redirect(callbackUrl);
          } catch {
            redirect(`/admin/login?callbackUrl=${encodeURIComponent(callbackUrl)}&error=1`);
          }
        }}
        className="mt-6 space-y-4"
      >
        <input type="hidden" name="callbackUrl" value="/admin/projects" />
        <div>
          <label className="font-mono text-xs uppercase tracking-wider text-muted">Email</label>
          <input
            type="email"
            name="email"
            placeholder="you@domain.com"
            required
            className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink"
          />
        </div>
        <div>
          <label className="font-mono text-xs uppercase tracking-wider text-muted">Password</label>
          <input
            type="password"
            name="password"
            required
            className="mt-1 w-full rounded-xl border border-line bg-paper px-3 py-2 text-ink"
          />
        </div>
        <button
          type="submit"
          className="w-full rounded-full bg-moss px-4 py-2 font-mono text-sm text-paper hover:bg-moss-dark"
        >
          Sign in
        </button>
        <p className="pt-2 text-center text-xs text-muted">
          <Link href="/" className="underline">Back to site</Link>
        </p>
      </form>
    </div>
  );
}
