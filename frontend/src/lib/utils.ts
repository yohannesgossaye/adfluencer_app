import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}


export function getInitials(name?: string | null, fallback = "U"): string {
  if (!name || !name.trim()) return fallback.toUpperCase();

  const clean = name.trim();
  const text = clean.includes("@") ? clean.split("@")[0] : clean;
  const normalized = text.replace(/[._\-+]/g, " ").trim();
  const parts = normalized.split(/\s+/).filter(Boolean);

  if (parts.length >= 2) {
    return (parts[0][0] + parts[1][0]).toUpperCase();
  }
  if (parts.length === 1) {
    const single = parts[0];
    return (single.length >= 2 ? single.slice(0, 2) : single).toUpperCase();
  }
  return fallback.toUpperCase();
}
