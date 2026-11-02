"use client";

import Link from "next/link";
import { ArrowLeft } from "lucide-react";

import { cn } from "@/lib/utils";

type AuthPageCardProps = {
  children: React.ReactNode;
  size?: "md" | "lg";
  className?: string;
  closeHref?: string;
};

export function AuthPageCard({
  children,
  size = "md",
  className,
  closeHref = "/",
}: AuthPageCardProps) {
  return (
    <div
      className={cn(
        "relative mx-auto w-full rounded-2xl border border-white/10 bg-[color:var(--adfluence-bg)] text-white shadow-[0_24px_80px_rgba(6,17,31,0.55)]",
        size === "lg" ? "max-w-xl" : "max-w-md",
        className,
      )}
    >
      {/* A page, not a dialog: a plain way back instead of a close (X) button. */}
      <Link
        href={closeHref}
        className="absolute top-4 right-4 z-10 inline-flex min-h-11 items-center gap-1.5 rounded-lg px-2 text-sm text-[#CBD5E1] transition-colors hover:text-white sm:min-h-0 sm:py-1"
      >
        <ArrowLeft className="size-4" aria-hidden />
        Home
      </Link>
      <div className="space-y-6 p-6 sm:p-8">{children}</div>
    </div>
  );
}
