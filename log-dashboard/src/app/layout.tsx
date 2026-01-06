import type { Metadata } from "next";
import { Fira_Code, Geist, Geist_Mono } from "next/font/google";
import "./globals.css";

import { ToastProvider } from "@/providers/ToastProvider";
import CookieBanner from "@/components/layout/CookieBanner";

const geistSans = Geist({
  variable: "--font-geist-sans",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

const firaCode = Fira_Code({
  variable: "--font-fira-code",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "Sijil",
  description: "High-throughput log ingestion",
  icons: {
    icon: "../../public/logo.png",
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${geistSans.variable} ${geistMono.variable} ${firaCode.variable} `}
      >
        <ToastProvider>{children}</ToastProvider>
        <CookieBanner />
      </body>
    </html>
  );
}
