"use client";

import { useState } from "react";
import { postContact } from "@/lib/api";

const inputClasses =
  "w-full rounded-xl border border-line bg-paper px-4 py-3 text-ink placeholder:text-muted/70 focus:border-moss focus:outline-none";

export default function ContactForm() {
  const [name, setName] = useState("");
  const [email, setEmail] = useState("");
  const [body, setBody] = useState("");
  const [status, setStatus] = useState<"idle" | "sending" | "sent" | "error">("idle");
  const [error, setError] = useState("");

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setStatus("sending");
    setError("");
    try {
      await postContact({ name, email, body });
      setStatus("sent");
      setName("");
      setEmail("");
      setBody("");
    } catch (err) {
      setStatus("error");
      setError(err instanceof Error ? err.message : "Something went wrong. Try again.");
    }
  }

  if (status === "sent") {
    return (
      <div className="rounded-2xl border border-moss/30 bg-moss-light p-8 text-center">
        <p className="font-serif text-2xl text-moss-dark">Message sent</p>
        <p className="mt-2 text-muted">Thanks for reaching out. I&apos;ll get back to you soon.</p>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="space-y-5">
      <div className="grid gap-5 sm:grid-cols-2">
        <div>
          <label htmlFor="contact-name" className="mb-1.5 block font-mono text-sm text-muted">
            your name
          </label>
          <input
            id="contact-name"
            type="text"
            required
            maxLength={120}
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Ada Lovelace"
            className={inputClasses}
          />
        </div>
        <div>
          <label htmlFor="contact-email" className="mb-1.5 block font-mono text-sm text-muted">
            your email
          </label>
          <input
            id="contact-email"
            type="email"
            required
            maxLength={254}
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="ada@example.com"
            className={inputClasses}
          />
        </div>
      </div>
      <div>
        <label htmlFor="contact-body" className="mb-1.5 block font-mono text-sm text-muted">
          message
        </label>
        <textarea
          id="contact-body"
          required
          maxLength={5000}
          rows={5}
          value={body}
          onChange={(e) => setBody(e.target.value)}
          placeholder="What are you building?"
          className={inputClasses}
        />
      </div>
      {status === "error" ? (
        <p role="alert" className="font-mono text-sm text-red-700">
          {error}
        </p>
      ) : null}
      <button
        type="submit"
        disabled={status === "sending"}
        className="rounded-xl bg-moss px-6 py-3 font-mono text-sm font-medium text-paper transition-colors hover:bg-moss-dark disabled:cursor-not-allowed disabled:opacity-60"
      >
        {status === "sending" ? "sending…" : "send message"}
      </button>
    </form>
  );
}
