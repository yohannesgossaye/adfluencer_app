"use client";

import { Suspense } from "react";
import { useSearchParams } from "next/navigation";
import { ShieldCheck } from "lucide-react";

import { SIGN_OUT_MESSAGES, type SignOutReason } from "@/lib/session";

function Notice() {
  const reason = useSearchParams().get("reason") as SignOutReason | null;
  if (!reason || !(reason in SIGN_OUT_MESSAGES)) return null;
  return (
    <p role="status" className="flex items-start gap-2 rounded-xl border border-sky-400/30 bg-sky-500/10 px-3 py-2 text-sm text-sky-50">
      <ShieldCheck className="mt-0.5 size-4 shrink-0" aria-hidden />
      {SIGN_OUT_MESSAGES[reason]}
    </p>
  );
}

/** "You were signed out after inactivity" on the login pages (?reason=idle|expired). */
export function SignedOutNotice() {
  return (
    <Suspense fallback={null}>
      <Notice />
    </Suspense>
  );
}
