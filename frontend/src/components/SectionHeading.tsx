export default function SectionHeading({
  path,
  title,
  children,
}: {
  path: string;
  title: string;
  children?: React.ReactNode;
}) {
  return (
    <div className="mb-10 flex flex-wrap items-end justify-between gap-4">
      <div>
        <p className="font-mono text-sm text-moss">{path}</p>
        <h2 className="mt-2 font-serif text-3xl text-ink sm:text-4xl">{title}</h2>
      </div>
      {children}
    </div>
  );
}
