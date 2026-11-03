"use client";

import Link from "next/link";
import { createContext, useCallback, useContext, useEffect, useState, type ReactNode } from "react";
import { Clock, FileEdit, Lock, ShieldAlert } from "lucide-react";

import { Button } from "@/components/ui/button";
import { influencerApi } from "@/lib/api/influencer";
import { REQUIRED_FOR_SUBMIT } from "@/lib/influencer-profile";

type Account = {
  state: "loading" | "ready" | "error";
  status?: string;
  statusReason?: string;
  missingRequired: string[];
  refresh: () => Promise<void>;
};

const AccountContext = createContext<Account | null>(null);

/**
 * Loads the creator's own account status once for the whole portal, so every
 * page can explain what is available instead of failing on a 403.
 */
export function InfluencerAccountProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<Account["state"]>("loading");
  const [status, setStatus] = useState<string>();
  const [statusReason, setStatusReason] = useState<string>();
  const [missingRequired, setMissingRequired] = useState<string[]>([]);

  const refresh = useCallback(async () => {
    try {
      const res = await influencerApi.getProfile();
      setStatus(res.data?.user.status);
      setStatusReason(res.data?.user.status_reason);
      setMissingRequired((res.data?.completeness.missing ?? []).filter((f) => f in REQUIRED_FOR_SUBMIT));
      setState("ready");
    } catch {
      setState("error");
    }
  }, []);

  useEffect(() => {
    void Promise.resolve().then(refresh);
  }, [refresh]);

  return (
    <AccountContext.Provider value={{ state, status, statusReason, missingRequired, refresh }}>
      {children}
    </AccountContext.Provider>
  );
}

export function useInfluencerAccount(): Account {
  const ctx = useContext(AccountContext);
  if (!ctx) throw new Error("useInfluencerAccount must be used inside InfluencerAccountProvider");
  return ctx;
}

/**
 * Shows `children` only to approved creators; everyone else gets a clear
 * explanation of where their application stands and what to do next.
 * If the status can't be loaded the page is shown and the server decides.
 */
export function ApprovalGate({ feature, children }: { feature: string; children: ReactNode }) {
  const { state, status } = useInfluencerAccount();
  if (state === "loading") {
    return <p role="status" className="p-6 text-sm text-[#94A3B8]">Loading…</p>;
  }
  if (state === "error" || status === "active") return <>{children}</>;
  return <AccountStatusPanel feature={feature} />;
}

export function AccountStatusPanel({ feature }: { feature: string }) {
  const { status, statusReason, missingRequired } = useInfluencerAccount();

  const view = (() => {
    switch (status) {
      case "pending_profile":
        return {
          Icon: FileEdit,
          title: "Finish your profile to get started",
          body: `${feature} open once an admin approves your profile. Complete your profile and submit it for review.`,
          cta: { href: "/influencer/profile", label: "Complete my profile" },
        };
      case "pending_approval":
        return {
          Icon: Clock,
          title: "Your profile is under review",
          body: `An admin is checking your profile. ${feature} open as soon as you're approved.`,
          cta: { href: "/influencer/profile", label: "View my profile" },
        };
      case "rejected":
        return {
          Icon: ShieldAlert,
          title: "Changes requested",
          body: `Update your profile and re-apply. ${feature} open once you're approved.`,
          cta: { href: "/influencer/profile", label: "Update my profile" },
        };
      default:
        return {
          Icon: Lock,
          title: "Your account isn't active",
          body: "Contact support if you think this is a mistake.",
          cta: { href: "/influencer/support", label: "Contact support" },
        };
    }
  })();

  return (
    <section className="mx-auto max-w-xl rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.06)] p-6 text-center text-white">
      <view.Icon className="mx-auto size-10 text-[color:var(--adfluence-orange)]" aria-hidden />
      <h1 className="mt-3 text-xl font-semibold">{view.title}</h1>
      <p className="mt-2 text-sm text-[#CBD5E1]">{view.body}</p>
      {status === "rejected" && statusReason ? (
        <p className="mt-3 rounded-xl border border-rose-400/30 bg-rose-500/10 p-3 text-left text-sm text-rose-100">
          <span className="font-medium">Admin note: </span>
          {statusReason}
        </p>
      ) : null}
      {status === "pending_profile" && missingRequired.length > 0 ? (
        <ul className="mt-4 space-y-1 text-left text-sm text-amber-100">
          {missingRequired.map((f) => (
            <li key={f}>• {REQUIRED_FOR_SUBMIT[f]}</li>
          ))}
        </ul>
      ) : null}
      <div className="mt-5 flex flex-wrap justify-center gap-2">
        <Button asChild className="h-11 rounded-xl bg-[color:var(--adfluence-action)] px-5 text-white hover:bg-[color:var(--adfluence-action-hover)]">
          <Link href={view.cta.href}>{view.cta.label}</Link>
        </Button>
        {view.cta.href !== "/influencer/support" ? (
          <Button asChild variant="outline" className="h-11 rounded-xl border-white/15 bg-transparent px-5 text-white hover:bg-white/10">
            <Link href="/influencer/support">Contact support</Link>
          </Button>
        ) : null}
      </div>
    </section>
  );
}
