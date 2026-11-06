"use client";

import * as React from "react";
import { Loader2, ShieldAlert } from "lucide-react";

import { useAuth } from "@/components/providers/auth-provider";
import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { authApi } from "@/lib/api/auth";
import { renewSession } from "@/lib/api/api-client";
import { REFRESH_TOKEN_STORAGE_KEY, AUTH_TOKEN_STORAGE_KEY } from "@/lib/constants";
import {
  RENEW_EVERY_MS,
  WARNING_MS,
  clearStoredSession,
  isIdle,
  lastRenewal,
  loginUrl,
  markActivity,
  signOutAt,
  type SignOutReason,
} from "@/lib/session";

const ACTIVITY_EVENTS = ["pointerdown", "keydown", "wheel", "touchstart", "scroll", "mousemove"] as const;
const REASON_KEY = "adfluence_signed_out_reason";

/** Ends the session on the server, clears it here, and opens the login page with the reason. */
async function signOut(reason: SignOutReason) {
  const refreshToken = window.localStorage.getItem(REFRESH_TOKEN_STORAGE_KEY) || undefined;
  const url = loginUrl(reason);
  try {
    window.localStorage.setItem(REASON_KEY, reason);
  } catch {
    // ignore
  }
  // Don't keep the person waiting if the server is slow or down.
  await Promise.race([authApi.logout(refreshToken).catch(() => undefined), new Promise((r) => setTimeout(r, 2000))]);
  clearStoredSession();
  window.location.assign(url);
}

/**
 * Bank-style sign-out for a logged-in person: warns 60 seconds before the
 * idle limit, signs out at it, and renews the server session while they are
 * active (so a long form never times out mid-typing).
 */
function ActiveSessionGuard() {
  const [secondsLeft, setSecondsLeft] = React.useState<number | null>(null);
  const [staying, setStaying] = React.useState(false);
  const signingOut = React.useRef(false);

  const end = React.useCallback((reason: SignOutReason) => {
    if (signingOut.current) return;
    signingOut.current = true;
    void signOut(reason);
  }, []);

  // Activity: the person, in any tab. Renew the server session at most once a minute.
  React.useEffect(() => {
    let lastMark = 0;
    const onActivity = () => {
      const now = Date.now();
      if (signingOut.current || isIdle(now)) return;
      if (now - lastMark < 5_000) return;
      lastMark = now;
      markActivity(now);
      if (now - lastRenewal() > RENEW_EVERY_MS) void renewSession().catch(() => undefined);
    };
    for (const e of ACTIVITY_EVENTS) window.addEventListener(e, onActivity, { passive: true, capture: true });
    return () => {
      for (const e of ACTIVITY_EVENTS) window.removeEventListener(e, onActivity, { capture: true });
    };
  }, []);

  // Once a second: warn before the limit, sign out at it.
  React.useEffect(() => {
    const tick = () => {
      const { at, reason } = signOutAt();
      const left = at - Date.now();
      if (left <= 0) {
        end(reason);
        return;
      }
      setSecondsLeft(left <= WARNING_MS ? Math.ceil(left / 1000) : null);
    };
    tick();
    const timer = window.setInterval(tick, 1000);
    return () => window.clearInterval(timer);
  }, [end]);

  // Signed out in another tab: follow it to the login page.
  React.useEffect(() => {
    const onStorage = (e: StorageEvent) => {
      if (e.key !== AUTH_TOKEN_STORAGE_KEY || e.newValue) return;
      let reason: SignOutReason | undefined;
      try {
        reason = (window.localStorage.getItem(REASON_KEY) as SignOutReason | null) ?? undefined;
      } catch {
        reason = undefined;
      }
      signingOut.current = true;
      window.location.assign(loginUrl(reason));
    };
    window.addEventListener("storage", onStorage);
    return () => window.removeEventListener("storage", onStorage);
  }, []);

  const stay = async () => {
    setStaying(true);
    markActivity();
    try {
      await renewSession();
    } catch {
      // A network error keeps the session; a rejected session redirects to login.
    } finally {
      setStaying(false);
      setSecondsLeft(null);
    }
  };

  return (
    <Dialog open={secondsLeft !== null} onOpenChange={(open) => !open && void stay()}>
      <DialogContent>
        <DialogHeader className="pr-14">
          <DialogTitle className="flex items-center gap-2">
            <ShieldAlert className="size-5 text-amber-300" aria-hidden /> Are you still there?
          </DialogTitle>
          <DialogDescription>
            For your security, you&apos;ll be signed out in{" "}
            <span className="font-semibold tabular-nums text-white">{secondsLeft ?? 0} seconds</span> because there has been no activity.
          </DialogDescription>
        </DialogHeader>
        <DialogFooter className="mt-4">
          <Button type="button" variant="ghost" onClick={() => end("idle")} className="h-11 text-white hover:bg-white/10 sm:h-10">
            Sign out
          </Button>
          <Button
            type="button"
            disabled={staying}
            onClick={() => void stay()}
            className="h-11 rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)] sm:h-10"
          >
            {staying ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
            Stay signed in
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

/** Runs only while someone is logged in. */
export function SessionGuard() {
  const { user, accessToken } = useAuth();
  return user && accessToken ? <ActiveSessionGuard /> : null;
}
