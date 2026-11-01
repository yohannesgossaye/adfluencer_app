import type { ReactNode } from "react";
import { BadgeCheck, Clock, FileWarning, ShieldAlert } from "lucide-react";

import type { VerificationStatus } from "@/lib/api/business";
import { cn } from "@/lib/utils";

function formatDate(iso?: string) {
  if (!iso) return "";
  try {
    return new Intl.DateTimeFormat("en-GB", { day: "numeric", month: "short", year: "numeric" }).format(new Date(iso));
  } catch {
    return "";
  }
}

const copy: Record<
  VerificationStatus,
  { title: string; body: string; tone: string; icon: typeof Clock }
> = {
  unsubmitted: {
    title: "Get verified to work with creators",
    body: "Complete your business profile and submit it for verification. Until then you can prepare campaign drafts, but you can’t publish them, see recommended creators, invite creators or create contracts.",
    tone: "border-sky-400/30 bg-sky-500/10 text-sky-100",
    icon: FileWarning,
  },
  pending: {
    title: "Your business is being verified",
    body: "Our team is checking your business details. You can keep preparing campaign drafts; publishing and inviting creators unlock once you’re verified.",
    tone: "border-amber-400/30 bg-amber-500/10 text-amber-100",
    icon: Clock,
  },
  verified: {
    title: "Your business is verified",
    body: "You can publish campaigns, invite creators and create contracts.",
    tone: "border-emerald-400/30 bg-emerald-500/10 text-emerald-100",
    icon: BadgeCheck,
  },
  rejected: {
    title: "Verification was not approved",
    body: "Update your profile as described below and submit it again.",
    tone: "border-rose-400/30 bg-rose-500/10 text-rose-100",
    icon: ShieldAlert,
  },
};

/**
 * Shows where the business stands in verification, and why when rejected.
 */
export function VerificationStatusCard({
  status,
  note,
  submittedAt,
  verifiedAt,
  action,
  className,
}: {
  status: VerificationStatus;
  note?: string;
  submittedAt?: string;
  verifiedAt?: string;
  action?: ReactNode;
  className?: string;
}) {
  const c = copy[status] ?? copy.unsubmitted;
  const Icon = c.icon;
  const date =
    status === "verified" && verifiedAt
      ? `Verified on ${formatDate(verifiedAt)}`
      : status === "pending" && submittedAt
        ? `Submitted on ${formatDate(submittedAt)}`
        : "";

  return (
    <section
      role={status === "rejected" ? "alert" : "status"}
      className={cn("flex flex-col gap-3 rounded-2xl border p-4 sm:flex-row sm:items-start", c.tone, className)}
    >
      <Icon className="size-6 shrink-0" aria-hidden />
      <div className="min-w-0 flex-1 space-y-1">
        <p className="font-semibold text-white">{c.title}</p>
        <p className="text-sm opacity-90">{c.body}</p>
        {note && status !== "verified" ? (
          <p className="mt-2 rounded-xl border border-white/10 bg-black/20 px-3 py-2 text-sm text-white">
            <span className="font-medium">{status === "rejected" ? "Reason from Adfluence: " : "Note: "}</span>
            {note}
          </p>
        ) : null}
        {date ? <p className="text-xs opacity-75">{date}</p> : null}
      </div>
      {action ? <div className="shrink-0">{action}</div> : null}
    </section>
  );
}
