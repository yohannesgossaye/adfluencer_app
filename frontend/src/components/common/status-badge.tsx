import { cn } from "@/lib/utils";
import type { CampaignStatus, ContractStatus, InfluencerStatus } from "@/types";

type Status = CampaignStatus | ContractStatus | InfluencerStatus;

const statusClassMap: Record<Status, string> = {
  draft: "bg-slate-100 text-slate-700",
  pending: "bg-amber-100 text-amber-700",
  active: "bg-emerald-100 text-emerald-700",
  completed: "bg-blue-100 text-blue-700",
  cancelled: "bg-rose-100 text-rose-700",
  approved: "bg-emerald-100 text-emerald-700",
  rejected: "bg-rose-100 text-rose-700",
  accepted: "bg-indigo-100 text-indigo-700",
};

export function StatusBadge({ status }: { status: Status }) {
  return (
    <span
      className={cn(
        "inline-flex rounded-full px-2.5 py-1 text-xs font-medium capitalize",
        statusClassMap[status],
      )}
    >
      {status}
    </span>
  );
}
