import Link from "next/link";

export default function NotFound() {
  return (
    <div className="mx-auto flex w-full max-w-3xl flex-col items-start px-5 py-32 sm:px-8">
      <p className="font-mono text-sm text-moss">404</p>
      <h1 className="mt-4 font-serif text-4xl text-ink sm:text-5xl">
        This page wandered off
      </h1>
      <p className="mt-4 max-w-md leading-relaxed text-muted">
        Whatever you were looking for isn&apos;t here. The journal is probably
        a better place to be anyway.
      </p>
      <Link
        href="/"
        className="mt-8 rounded-xl bg-moss px-6 py-3 font-mono text-sm font-medium text-paper transition-colors hover:bg-moss-dark"
      >
        back home
      </Link>
    </div>
  );
}
