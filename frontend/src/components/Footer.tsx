export default function Footer() {
  return (
    <footer className="border-t border-line">
      <div className="mx-auto flex w-full max-w-5xl flex-col gap-3 px-5 py-8 font-mono text-sm text-muted sm:flex-row sm:items-center sm:justify-between sm:px-8">
        <p>© {new Date().getFullYear()} ikhtiar</p>
        <p>
          built with <span className="text-moss-dark">go</span> +{" "}
          <span className="text-moss-dark">next.js</span> · source on{" "}
          <a
            href="https://github.com/IKH-TIAR"
            target="_blank"
            rel="noreferrer"
            className="text-moss-dark underline underline-offset-3 hover:text-moss"
          >
            github
          </a>
        </p>
      </div>
    </footer>
  );
}
