"use client";

import Image from "next/image";
import { useState } from "react";
import { cn, getInitials } from "@/lib/utils";

export interface UserAvatarProps {
  src?: string | null;
  name?: string | null;
  size?: "sm" | "md" | "lg" | "xl" | "2xl" | "custom";
  className?: string;
  textClassName?: string;
  fallback?: string;
  sizes?: string;
}

const sizeClasses = {
  sm: "h-8 w-8 text-xs",
  md: "h-11 w-11 text-sm",
  lg: "h-14 w-14 text-base",
  xl: "h-16 w-16 text-lg",
  "2xl": "h-24 w-24 text-2xl font-bold",
  custom: "",
};

export function UserAvatar({
  src,
  name,
  size = "md",
  className,
  textClassName,
  fallback = "U",
  sizes = "64px",
}: UserAvatarProps) {
  const [hasError, setHasError] = useState(false);
  const initials = getInitials(name, fallback);

  return (
    <div
      className={cn(
        "relative flex shrink-0 items-center justify-center overflow-hidden rounded-full border border-white/10 bg-gradient-to-br from-[color:var(--adfluence-orange)]/35 via-white/10 to-white/5 select-none",
        size !== "custom" && sizeClasses[size],
        className
      )}
    >
      {src && !hasError ? (
        <Image
          src={src}
          alt={name || "Avatar"}
          fill
          unoptimized
          className="object-cover"
          sizes={sizes}
          onError={() => setHasError(true)}
        />
      ) : (
        <span className={cn("font-bold text-white uppercase tracking-wider", textClassName)}>
          {initials}
        </span>
      )}
    </div>
  );
}
