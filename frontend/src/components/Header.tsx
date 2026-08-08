import Link from "next/link";

const nav = [
  { href: "/journal", label: "journal" },
  { href: "/projects", label: "projects" },
  { href: "/photos", label: "photos" },
  { href: "/#contact", label: "contact" },
];

export default function Header() {
  return (
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
        <nav className="flex items-center gap-5 sm:gap-8">
          {nav.map((item) => (
            <Link
              key={item.href}
              href={item.href}
              className="font-mono text-sm text-muted transition-colors hover:text-moss-dark"
            >
              {item.label}
            </Link>
          ))}
        </nav>
      </div>
    </header>
  );
}
