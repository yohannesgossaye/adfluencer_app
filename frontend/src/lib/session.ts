/**
 * Bank-style sessions on the website: you are signed out after a period
 * without activity (the server enforces the same limit), and after a
 * maximum time however active you are.
 *
 * "Activity" means the person: typing, clicking, scrolling, touching, in any
 * Adfluence tab. Background traffic (chat updates, reconnects) never counts,
 * and never renews a session that has gone idle.
 */
import { AUTH_TOKEN_STORAGE_KEY, AUTH_USER_STORAGE_KEY, REFRESH_TOKEN_STORAGE_KEY } from "@/lib/constants";

const POLICY_KEY = "adfluence_session_policy";
const ACTIVITY_KEY = "adfluence_last_activity";
const RENEWED_KEY = "adfluence_session_renewed";

/** The site signs out this much before the server would, so "Stay signed in" always works. */
export const SAFETY_MARGIN_MS = 60_000;
/** The warning shows this long before signing out. */
export const WARNING_MS = 60_000;
/** While someone is active, the server session is renewed at most this often. */
export const RENEW_EVERY_MS = 60_000;

type Policy = { idleMs: number; endsAt: number | null };

/** Fallback when a login predates the server sending its limits. */
function defaultIdleMs(): number {
  try {
    const role = JSON.parse(window.localStorage.getItem(AUTH_USER_STORAGE_KEY) ?? "{}")?.user?.role;
    return (role === "admin" ? 10 : 15) * 60_000;
  } catch {
    return 15 * 60_000;
  }
}

function read(key: string): string | null {
  try {
    return window.localStorage.getItem(key);
  } catch {
    return null;
  }
}

function write(key: string, value: string) {
  try {
    window.localStorage.setItem(key, value);
  } catch {
    // Storage blocked: the server still enforces the limits.
  }
}

/** Saves the limits the server sent with a login or refresh. */
export function rememberSessionLimits(tokens: { idle_timeout_seconds?: number; session_ends_at?: string }) {
  const policy = {
    idleMs: tokens.idle_timeout_seconds ? tokens.idle_timeout_seconds * 1000 : undefined,
    endsAt: tokens.session_ends_at ? Date.parse(tokens.session_ends_at) : undefined,
  };
  write(POLICY_KEY, JSON.stringify(policy));
  write(RENEWED_KEY, String(Date.now()));
}

export function sessionPolicy(): Policy {
  try {
    const p = JSON.parse(read(POLICY_KEY) ?? "{}");
    return { idleMs: p.idleMs || defaultIdleMs(), endsAt: Number.isFinite(p.endsAt) ? p.endsAt : null };
  } catch {
    return { idleMs: defaultIdleMs(), endsAt: null };
  }
}

/** Records that the person did something (shared by all tabs). */
export function markActivity(now = Date.now()) {
  write(ACTIVITY_KEY, String(now));
}

export function lastActivity(): number {
  const v = Number(read(ACTIVITY_KEY));
  if (Number.isFinite(v) && v > 0) return v;
  // No activity recorded yet (e.g. logged in before this version): start counting now.
  const now = Date.now();
  markActivity(now);
  return now;
}

export function lastRenewal(): number {
  const v = Number(read(RENEWED_KEY));
  return Number.isFinite(v) && v > 0 ? v : 0;
}

/** When the site signs out: idle limit after the last activity, or the session's end, whichever is first. */
export function signOutAt(): { at: number; reason: SignOutReason } {
  const { idleMs, endsAt } = sessionPolicy();
  const idleAt = lastActivity() + idleMs - SAFETY_MARGIN_MS;
  if (endsAt !== null && endsAt - SAFETY_MARGIN_MS < idleAt) return { at: endsAt - SAFETY_MARGIN_MS, reason: "expired" };
  return { at: idleAt, reason: "idle" };
}

/** True when the person has been inactive past the limit: never renew then. */
export function isIdle(now = Date.now()): boolean {
  return now >= signOutAt().at;
}

export type SignOutReason = "idle" | "expired";

/** Clears the saved session in this browser (all tabs see it go). */
export function clearStoredSession() {
  for (const key of [AUTH_TOKEN_STORAGE_KEY, REFRESH_TOKEN_STORAGE_KEY, AUTH_USER_STORAGE_KEY, POLICY_KEY, ACTIVITY_KEY, RENEWED_KEY]) {
    try {
      window.localStorage.removeItem(key);
    } catch {
      // ignore
    }
  }
}

/** The login page for where the person is, remembering where to come back to. */
export function loginUrl(reason?: SignOutReason) {
  const path = `${window.location.pathname}${window.location.search}`;
  const login = path.startsWith("/admin") ? "/admin/login" : "/auth/login";
  const params = new URLSearchParams();
  if (!path.startsWith("/auth") && !path.startsWith("/admin/login") && path !== "/") params.set("next", path);
  if (reason) params.set("reason", reason);
  const q = params.toString();
  return q ? `${login}?${q}` : login;
}

export const SIGN_OUT_MESSAGES: Record<SignOutReason, string> = {
  idle: "For your security, you were signed out after a period of inactivity. Please log in again.",
  expired: "Your session reached its time limit. Please log in again.",
};
