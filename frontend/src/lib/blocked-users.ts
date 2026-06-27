const BLOCKED_KEY = "adfluence_blocked_user_emails";

function parseList(): string[] {
  if (typeof window === "undefined") return [];
  try {
    const raw = localStorage.getItem(BLOCKED_KEY);
    if (!raw) return [];
    const arr = JSON.parse(raw) as unknown;
    return Array.isArray(arr) ? arr.map(String) : [];
  } catch {
    return [];
  }
}

export function getBlockedEmails(): string[] {
  return parseList();
}

export function isEmailBlocked(email: string): boolean {
  const e = email.trim().toLowerCase();
  return getBlockedEmails().some((x) => x.toLowerCase() === e);
}

export function blockUserEmail(email: string) {
  if (typeof window === "undefined") return;
  const e = email.trim().toLowerCase();
  const next = [...new Set([...parseList(), e])];
  localStorage.setItem(BLOCKED_KEY, JSON.stringify(next));
}

export function unblockUserEmail(email: string) {
  if (typeof window === "undefined") return;
  const e = email.trim().toLowerCase();
  const next = parseList().filter((x) => x.toLowerCase() !== e);
  localStorage.setItem(BLOCKED_KEY, JSON.stringify(next));
}
