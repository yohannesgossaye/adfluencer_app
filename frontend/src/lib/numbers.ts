/**
 * Parses an amount typed by a person: "12000", "12,000", "12k", "1.2m", " 15 000 ".
 * Returns undefined for an empty string and NaN for anything unreadable, so
 * callers can tell "left blank" from "typed something invalid".
 */
export function parseAmount(raw: string): number | undefined {
  const s = raw.trim().toLowerCase().replace(/[,\s]/g, "");
  if (s === "") return undefined;
  const m = /^(\d+(?:\.\d+)?)([km])?$/.exec(s);
  if (!m) return NaN;
  const n = Number(m[1]);
  return m[2] === "k" ? n * 1_000 : m[2] === "m" ? n * 1_000_000 : n;
}
