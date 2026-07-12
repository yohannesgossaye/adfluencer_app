"use client";

import Link from "next/link";
import { X } from "lucide-react";

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
        "relative w-full rounded-2xl border border-white/10 bg-[color:var(--adfluence-bg)] text-white shadow-[0_24px_80px_rgba(6,17,31,0.55)]",
        size === "lg" ? "max-w-xl" : "max-w-md",
        className,
      )}
    >
      <Link
        href={closeHref}
        aria-label="Close and return home"
        className="absolute top-4 right-4 z-10 inline-flex h-11 w-11 items-center justify-center rounded-xl border border-white/10 bg-white/5 text-white/80 transition-colors hover:bg-white/10 hover:text-white sm:h-auto sm:w-auto sm:rounded-lg sm:p-1.5"
      >
        <X className="size-5 sm:size-4" />
        <span className="sr-only">Close</span>
      </Link>
      <div className="space-y-6 p-6 pr-16 sm:pr-14">{children}</div>
    </div>
  );
}
