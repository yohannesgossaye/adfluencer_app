import type { StatusInfo, StatusTone } from "@/lib/status";
import { cn } from "@/lib/utils";

const TONE: Record<StatusTone, string> = {
  success: "border-emerald-400/40 bg-emerald-400/10 text-emerald-200",
  warning: "border-amber-300/40 bg-amber-300/10 text-amber-100",
  danger: "border-rose-400/50 bg-rose-400/10 text-rose-200",
  info: "border-sky-400/40 bg-sky-400/10 text-sky-200",
  neutral: "border-white/15 bg-white/5 text-[#CBD5E1]",
};

const DOT: Record<StatusTone, string> = {
  success: "bg-emerald-300",
  warning: "bg-amber-300",
  danger: "bg-rose-300",
  info: "bg-sky-300",
  neutral: "bg-[#94A3B8]",
};

/**
 * The one status pill used on every screen. Pass the result of a helper from
 * lib/status (e.g. campaignStatus(row.status)) rather than a raw status code.
 */
export function StatusBadge({ status, className }: { status: StatusInfo; className?: string }) {
  return (
    <span
      className={cn(
        "inline-flex w-fit shrink-0 items-center gap-1.5 rounded-full border px-2.5 py-0.5 text-xs font-medium",
        TONE[status.tone],
        className,
      )}
    >
      <span aria-hidden className={cn("size-1.5 rounded-full", DOT[status.tone])} />
      {status.label}
    </span>
  );
}
