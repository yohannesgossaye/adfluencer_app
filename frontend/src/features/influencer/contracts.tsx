"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";
import { influencerApi } from "@/lib/api/influencer";
import {
  influencerContractStats,
  initialInfluencerContracts,
  type InfluencerContractRecord,
  type InfluencerContractStatus,
} from "@/data/influencer-contracts-mock";
import { cn } from "@/lib/utils";

function formatEtb(n: number) {
  if (n >= 1000) return `ETB ${Math.round(n / 1000)}K`;
  return `ETB ${n.toLocaleString()}`;
}

function statusBadge(status: string) {
  const s = status.toLowerCase();
  switch (s) {
    case "pending_review":
    case "pending":
      return "border-amber-400/30 bg-amber-500/10 text-amber-200";
    case "accepted":
      return "border-violet-400/30 bg-violet-500/10 text-violet-200";
    case "active":
      return "border-emerald-400/30 bg-emerald-500/10 text-emerald-200";
    case "rejected":
      return "border-rose-400/30 bg-rose-500/10 text-rose-200";
    case "completed":
      return "border-sky-400/30 bg-sky-500/10 text-sky-200";
    default:
      return "border-white/10 bg-white/5 text-[#CBD5E1]";
  }
}

function statusLabel(status: string) {
  return status.replace("_", " ");
}

export function InfluencerContracts() {
  const [rows, setRows] = useState<InfluencerContractRecord[]>(() =>
    structuredClone(initialInfluencerContracts),
  );
  const [selectedId, setSelectedId] = useState(initialInfluencerContracts[0]?.id ?? "");
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  const loadContracts = async () => {
    try {
      const res = await influencerApi.listContracts();
      if (res.data?.items && res.data.items.length > 0) {
        const mapped: InfluencerContractRecord[] = res.data.items.map((c) => ({
          id: c.id,
          brand: c.business?.business_name || "Partner Brand",
          campaign: c.campaign?.title || "Campaign",
          deliverables: c.deliverables?.map((d) => d.type).join(", ") || "Post",
          timeline: c.start_date && c.end_date
            ? `${new Date(c.start_date).toLocaleDateString()} - ${new Date(c.end_date).toLocaleDateString()}`
            : "Ongoing",
          amountEtb: c.fees?.influencer_receivable_etb || c.fees?.amount_etb || 0,
          status: (c.status.toLowerCase() as InfluencerContractStatus) || "pending",
          terms: c.description || "Platform escrow standard terms (5% fee applied).",
          deadline: c.end_date ? new Date(c.end_date).toLocaleDateString() : "TBD",
          notes: c.payment?.status ? `Payment status: ${c.payment.status}` : "Standard escrow contract",
          checklist: c.deliverables?.map((d, i) => ({
            id: d.id || `deliv-${i}`,
            label: `${d.quantity}x ${d.type}`,
            done: d.status === "completed",
          })) || [],
        }));
        setRows(mapped);
        if (!selectedId || !mapped.some((m) => m.id === selectedId)) {
          setSelectedId(mapped[0].id);
        }
      }
    } catch {
      // fallback to mock
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadContracts();
  }, []);

  const selected = useMemo(
    () => rows.find((r) => r.id === selectedId) ?? rows[0],
    [rows, selectedId],
  );

  const selectContract = useCallback((id: string) => {
    setSelectedId(id);
    if (typeof window !== "undefined" && window.matchMedia("(max-width: 1279px)").matches) {
      setMobileShowDetail(true);
    }
  }, []);

  const acceptContract = async (id: string) => {
    try {
      await influencerApi.acceptContract(id);
      await loadContracts();
    } catch {
      setRows((prev) =>
        prev.map((r) => {
          if (r.id !== id) return r;
          if (r.status === "pending_review" || r.status === "pending") return { ...r, status: "active" };
          return r;
        }),
      );
    }
  };

  const rejectContract = async (id: string) => {
    try {
      await influencerApi.rejectContract(id);
      await loadContracts();
    } catch {
      setRows((prev) =>
        prev.map((r) => {
          if (r.id !== id) return r;
          if (r.status === "completed" || r.status === "rejected") return r;
          return { ...r, status: "rejected" };
        }),
      );
    }
  };

  const markComplete = (id: string) => {
    setRows((prev) =>
      prev.map((r) => (r.id === id && r.status === "active" ? { ...r, status: "completed" } : r)),
    );
  };

  const toggleCheck = (contractId: string, itemId: string) => {
    setRows((prev) =>
      prev.map((r) => {
        if (r.id !== contractId) return r;
        return {
          ...r,
          checklist: r.checklist.map((c) =>
            c.id === itemId ? { ...c, done: !c.done } : c,
          ),
        };
      }),
    );
  };

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Contracts</h1>
        <p className="mt-1 max-w-xl text-sm text-[#CBD5E1]">
          Review agreements, deliverables, and payouts from brand collaborations.
        </p>
      </div>

      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {[
          { label: "Pending contracts", value: String(influencerContractStats.pending) },
          { label: "Active contracts", value: String(influencerContractStats.active) },
          { label: "Completed", value: String(influencerContractStats.completed) },
          { label: "Total earnings", value: formatEtb(influencerContractStats.totalEarningsEtb) },
        ].map((c) => (
          <article
            key={c.label}
            className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
          >
            <p className="text-xs text-[#CBD5E1]">{c.label}</p>
            <p className="mt-2 text-2xl font-bold text-[color:var(--adfluence-orange)]">
              {c.value}
            </p>
          </article>
        ))}
      </section>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_400px]">
        <section className={cn("space-y-3", mobileShowDetail && "hidden xl:block")}>
          <h2 className="text-sm font-medium uppercase tracking-wide text-[#94A3B8]">Your contracts</h2>
          <div className="space-y-3">
            {rows.map((r) => {
              const rowOn = r.id === selectedId;
              return (
                <button
                  key={r.id}
                  type="button"
                  onClick={() => selectContract(r.id)}
                  className={cn(
                    "w-full rounded-2xl border p-4 text-left transition-colors",
                    rowOn
                      ? "border-[color:var(--adfluence-orange)]/35 bg-[color:var(--adfluence-orange)]/10"
                      : "border-white/10 bg-[rgba(255,255,255,0.08)] hover:border-white/20",
                  )}
                >
                  <div className="flex flex-wrap items-start justify-between gap-2">
                    <div>
                      <p className="font-semibold text-white">{r.brand}</p>
                      <p className="text-sm text-[#CBD5E1]">{r.campaign}</p>
                    </div>
                    <span
                      className={cn(
                        "rounded-full border px-2.5 py-0.5 text-xs capitalize",
                        statusBadge(r.status),
                      )}
                    >
                      {statusLabel(r.status)}
                    </span>
                  </div>
                  <p className="mt-2 text-xs text-[#94A3B8]">{r.deliverables}</p>
                  <div className="mt-2 flex flex-wrap gap-3 text-sm text-[#CBD5E1]">
                    <span>{formatEtb(r.amountEtb)}</span>
                    <span>·</span>
                    <span>{r.timeline}</span>
                  </div>
                  <div className="mt-3 flex flex-wrap gap-2">
                    <Button
                      type="button"
                      size="sm"
                      variant="outline"
                      className="h-11 rounded-lg border-white/15 bg-transparent text-xs text-white hover:bg-white/10 sm:h-8"
                      onClick={(e) => {
                        e.stopPropagation();
                        selectContract(r.id);
                      }}
                    >
                      Review
                    </Button>
                    <Button
                      type="button"
                      size="sm"
                      className="h-11 rounded-lg bg-[color:var(--adfluence-orange)] text-xs text-white hover:bg-[color:var(--adfluence-orange-hover)] disabled:opacity-40 sm:h-8"
                      disabled={r.status === "active" || r.status === "completed" || r.status === "rejected"}
                      onClick={(e) => {
                        e.stopPropagation();
                        acceptContract(r.id);
                      }}
                    >
                      Accept
                    </Button>
                    <Button
                      type="button"
                      size="sm"
                      variant="outline"
                      className="h-11 rounded-lg border-white/15 text-xs text-rose-200 hover:bg-rose-500/10 sm:h-8"
                      disabled={r.status === "completed" || r.status === "rejected"}
                      onClick={(e) => {
                        e.stopPropagation();
                        rejectContract(r.id);
                      }}
                    >
                      Reject
                    </Button>
                    <Button
                      type="button"
                      size="sm"
                      className="h-11 rounded-lg bg-white/10 text-xs text-white hover:bg-white/15 disabled:opacity-40 sm:h-8"
                      disabled={r.status !== "active"}
                      onClick={(e) => {
                        e.stopPropagation();
                        markComplete(r.id);
                      }}
                    >
                      Mark complete
                    </Button>
                  </div>
                </button>
              );
            })}
          </div>
        </section>

        {selected ? (
          <aside
            className={cn(
              "xl:sticky xl:top-4 xl:self-start",
              !mobileShowDetail && "hidden xl:block",
            )}
          >
            <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
              <div className="mb-3 xl:hidden">
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="shrink-0 text-white hover:bg-white/10"
                  aria-label="Back to contracts"
                  onClick={() => setMobileShowDetail(false)}
                >
                  <ArrowLeft className="size-5" />
                </Button>
              </div>
              <h3 className="text-lg font-semibold text-white">Contract detail</h3>
              <p className="mt-1 text-sm text-[#CBD5E1]">{selected.brand}</p>
              <span
                className={cn(
                  "mt-2 inline-flex rounded-full border px-2.5 py-0.5 text-xs capitalize",
                  statusBadge(selected.status),
                )}
              >
                {statusLabel(selected.status)}
              </span>

              <div className="mt-4 space-y-2 text-sm text-[#CBD5E1]">
                <p>
                  <span className="text-white">Campaign:</span> {selected.campaign}
                </p>
                <p>
                  <span className="text-white">Payment:</span> {formatEtb(selected.amountEtb)}
                </p>
                <p>
                  <span className="text-white">Timeline:</span> {selected.timeline}
                </p>
                <p>
                  <span className="text-white">Deadline:</span> {selected.deadline}
                </p>
                <p>
                  <span className="text-white">Terms:</span> {selected.terms}
                </p>
                <p>
                  <span className="text-white">Notes:</span> {selected.notes}
                </p>
              </div>

              <div className="mt-4">
                <p className="text-xs font-medium uppercase text-[#94A3B8]">Deliverables</p>
                <ul className="mt-2 space-y-2">
                  {selected.checklist.map((item) => (
                    <li key={item.id} className="flex items-center gap-2 text-sm text-[#E2E8F0]">
                      <input
                        type="checkbox"
                        checked={item.done}
                        onChange={() => toggleCheck(selected.id, item.id)}
                        className="h-4 w-4 rounded border-white/20 bg-[color:var(--adfluence-surface)] accent-[color:var(--adfluence-orange)]"
                      />
                      {item.label}
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </aside>
        ) : null}
      </div>
    </div>
  );
}
