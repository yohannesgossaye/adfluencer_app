import type { ReactNode } from "react";

import type { ContractDetail } from "@/lib/api/business";
import { contractStatus, deliverableLabel, formatDate, formatEtb } from "@/lib/contracts";
import { cn } from "@/lib/utils";
import { StatusBadge } from "@/components/common/status-badge";

function Line({ label, value, strong }: { label: string; value: ReactNode; strong?: boolean }) {
  return (
    <div className={cn("flex items-baseline justify-between gap-4 py-1.5 text-sm", strong && "border-t border-white/10 pt-2 font-semibold")}>
      <span className={strong ? "text-white" : "text-[#94A3B8]"}>{label}</span>
      <span className="text-right text-white tabular-nums">{value}</span>
    </div>
  );
}

/**
 * Contract details shown to either side. Money is shown from the viewer's
 * point of view: what the business pays, or what the creator receives.
 */
export function ContractDetailView({
  contract,
  role,
  actions,
}: {
  contract: ContractDetail;
  role: "business" | "influencer";
  actions?: ReactNode;
}) {
  const status = contractStatus(contract.status, role);
  const fees = contract.fees;
  const counterparty =
    role === "business" ? contract.influencer?.creator_name || "Creator" : contract.business?.business_name || "Business";

  return (
    <div className="space-y-5">
      <div>
        <div className="flex flex-wrap items-center gap-2">
          <h2 className="text-lg font-semibold text-white">{contract.title}</h2>
          <StatusBadge status={status.badge} />
        </div>
        <p className="mt-1 text-sm text-[#CBD5E1]">
          {role === "business" ? "With " : "From "}
          <span className="text-white">{counterparty}</span> · {contract.campaign?.title}
        </p>
        {status.hint ? <p className="mt-2 rounded-xl border border-white/10 bg-black/20 px-3 py-2 text-sm text-[#E2E8F0]">{status.hint}</p> : null}
      </div>

      {actions ? <div className="flex flex-wrap gap-2">{actions}</div> : null}

      <section>
        <h3 className="mb-1 text-xs font-semibold tracking-wide text-[#94A3B8] uppercase">Payment</h3>
        {role === "business" ? (
          <>
            <Line label="Agreed amount" value={formatEtb(fees.amount_etb)} />
            <Line label={`Platform fee (${fees.business_fee_percent}%)`} value={formatEtb(fees.business_fee_amount_etb)} />
            <Line label="You pay" value={formatEtb(fees.total_business_payable_etb)} strong />
            <p className="mt-1 text-xs text-[#94A3B8]">
              The creator receives {formatEtb(fees.influencer_receivable_etb)} after their {fees.influencer_fee_percent}% platform fee.
            </p>
          </>
        ) : (
          <>
            <Line label="Agreed amount" value={formatEtb(fees.amount_etb)} />
            <Line label={`Platform fee (${fees.influencer_fee_percent}%)`} value={`− ${formatEtb(fees.influencer_fee_amount_etb)}`} />
            <Line label="You receive" value={formatEtb(fees.influencer_receivable_etb)} strong />
          </>
        )}
        {contract.payment ? (
          <p className="mt-2 text-xs text-[#CBD5E1]">
            Escrow: <span className="capitalize">{contract.payment.status}</span>
            {contract.payment.paid_at ? ` on ${formatDate(contract.payment.paid_at)}` : ""}
          </p>
        ) : null}
        {contract.payout ? (
          <p className="text-xs text-[#CBD5E1]">
            Payout to creator: <span className="capitalize">{contract.payout.status}</span> ({formatEtb(contract.payout.amount_etb)})
          </p>
        ) : null}
      </section>

      <section>
        <h3 className="mb-1 text-xs font-semibold tracking-wide text-[#94A3B8] uppercase">Deliverables</h3>
        <ul className="space-y-1.5">
          {contract.deliverables.map((d) => (
            <li key={d.id} className="flex items-start justify-between gap-3 rounded-xl border border-white/10 bg-[rgba(255,255,255,0.04)] px-3 py-2 text-sm">
              <span className="text-white">
                {d.quantity} × {deliverableLabel(d.type)}
                {d.description ? <span className="block text-xs text-[#94A3B8]">{d.description}</span> : null}
              </span>
            </li>
          ))}
        </ul>
      </section>

      <section className="text-sm">
        <h3 className="mb-1 text-xs font-semibold tracking-wide text-[#94A3B8] uppercase">Timeline</h3>
        <Line label="Start" value={formatDate(contract.start_date)} />
        <Line label="End" value={formatDate(contract.end_date)} />
        <Line label="Created" value={formatDate(contract.created_at)} />
        {contract.accepted_by_influencer_at ? <Line label="Accepted by creator" value={formatDate(contract.accepted_by_influencer_at)} /> : null}
        {contract.completed_at ? <Line label="Completed" value={formatDate(contract.completed_at)} /> : null}
      </section>

      {contract.description ? (
        <section>
          <h3 className="mb-1 text-xs font-semibold tracking-wide text-[#94A3B8] uppercase">Notes</h3>
          <p className="text-sm whitespace-pre-wrap text-[#CBD5E1]">{contract.description}</p>
        </section>
      ) : null}
    </div>
  );
}
