import type { ReactNode } from "react";
import Link from "next/link";

export default function AdminLayout({ children }: { children: ReactNode }) {
  return (
    <>
      <header className="border-b border-line">
        <div className="mx-auto flex w-full max-w-5xl items-center justify-between px-5 py-5 sm:px-8">
          <Link
            href="/"
            className="group flex items-baseline gap-2 font-mono text-sm text-ink"
          >
            <span className="text-moss transition-colors group-hover:text-moss-dark">
              ~/
            </span>
            <span className="font-medium">ikhtiar</span>
          </Link>
          <span className="font-mono text-xs uppercase tracking-wider text-muted">
            admin
          </span>
        </div>
      </header>
      <main className="mx-auto w-full max-w-5xl px-5 py-10 sm:px-8">{children}</main>
    </>
  );
}
