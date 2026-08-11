import type { Metadata } from "next";
import { Figtree, IBM_Plex_Mono, Instrument_Serif } from "next/font/google";
import "./globals.css";

const instrument = Instrument_Serif({
  weight: "400",
  subsets: ["latin"],
  variable: "--font-instrument",
  display: "swap",
});

const figtree = Figtree({
  subsets: ["latin"],
  variable: "--font-figtree",
  display: "swap",
});

const plex = IBM_Plex_Mono({
  weight: ["400", "500"],
  subsets: ["latin"],
  variable: "--font-plex",
  display: "swap",
});

export const metadata: Metadata = {
  title: {
    default: "Ikhtiar: AI developer & automation engineer",
    template: "%s · Ikhtiar",
  },
  description:
    "AI developer and automation engineer. Building agents with LangChain and LangGraph, automations with n8n, and learning Go. Journaling the journey.",
  metadataBase: new URL("https://portfolio-v2.vercel.app"),
  openGraph: {
    title: "Ikhtiar: AI developer & automation engineer",
    description:
      "Building agents, automations, and a public learning journal.",
    type: "website",
  },
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html
      lang="en"
      className={`${instrument.variable} ${figtree.variable} ${plex.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col">
        {children}
      </body>
    </html>
  );
}
