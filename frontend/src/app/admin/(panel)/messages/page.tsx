import AdminNav from "@/components/AdminNav";
import { adminFetch } from "@/lib/admin";
import { deleteMessage } from "@/lib/adminActions";
import { formatDate } from "@/lib/api";

type Message = {
  id: number;
  name: string;
  email: string;
  body: string;
  read: boolean;
  createdAt: string;
};

export default async function AdminMessagesPage() {
  const messages = await adminFetch<Message[]>("/api/admin/messages");
  return (
    <div>
      <AdminNav active="messages" />
      <h1 className="font-serif text-3xl text-ink">Messages</h1>
      <ul className="mt-6 space-y-3">
        {messages.map((m) => (
          <li
            key={m.id}
            className="rounded-2xl border border-line bg-paper-deep/40 px-5 py-4"
          >
            <div className="flex flex-wrap items-center gap-3">
              <p className="font-serif text-lg text-ink">{m.name}</p>
              <a href={`mailto:${m.email}`} className="font-mono text-sm text-moss underline underline-offset-2">
                {m.email}
              </a>
              <span className="ml-auto font-mono text-xs text-muted">{formatDate(m.createdAt)}</span>
            </div>
            <p className="mt-2 whitespace-pre-wrap text-ink-soft">{m.body}</p>
            <form
              action={async () => {
                "use server";
                await deleteMessage(m.id);
              }}
              className="mt-3"
            >
              <button
                type="submit"
                className="rounded-full border border-line bg-paper-deep/60 px-3 py-1 font-mono text-sm text-muted hover:bg-ink hover:text-paper"
              >
                delete
              </button>
            </form>
          </li>
        ))}
        {messages.length === 0 ? <li className="font-mono text-sm text-muted">inbox empty</li> : null}
      </ul>
    </div>
  );
}
