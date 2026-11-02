"use client";

import { useCallback, useEffect, useState } from "react";
import axios from "axios";
import { ArrowLeft, Loader2 } from "lucide-react";
import { toast } from "sonner";

import { ContractDetailView } from "@/components/contracts/contract-detail";
import { Button } from "@/components/ui/button";
import { influencerApi, type ContractDetail, type ContractListItem } from "@/lib/api/influencer";
import { contractStatus, formatEtb } from "@/lib/contracts";
import { cn } from "@/lib/utils";

const FILTERS = [
  { id: "", label: "All" },
  { id: "pending_influencer", label: "Needs your answer" },
  { id: "accepted", label: "Waiting for payment" },
  { id: "active", label: "In progress" },
  { id: "completed", label: "Completed" },
] as const;

function apiError(err: unknown, fallback: string) {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as { message?: string; error?: { type?: string } } | undefined;
    if (data?.error?.type === "account_not_verified") return "Your profile must be approved before you can accept contracts.";
    if (data?.error?.type === "invalid_contract_status") return "This contract can no longer be changed.";
    return data?.message || (!err.response ? "You appear to be offline." : fallback);
  }
  return fallback;
}

export function InfluencerContracts() {
  const [filter, setFilter] = useState<string>("");
  const [rows, setRows] = useState<ContractListItem[]>([]);
  const [listState, setListState] = useState<"loading" | "ready" | "error">("loading");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [detail, setDetail] = useState<ContractDetail | null>(null);
  const [detailError, setDetailError] = useState<string | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [confirm, setConfirm] = useState<"accept" | "decline" | null>(null);
  const [busy, setBusy] = useState(false);

  const loadList = useCallback(async () => {
    try {
      const res = await influencerApi.listContracts({ status: filter || undefined, limit: 50 });
      setRows(res.data?.items ?? []);
      setListState("ready");
    } catch {
      setListState("error");
    }
  }, [filter]);

  const loadDetail = useCallback(async (id: string) => {
    try {
      const res = await influencerApi.getContract(id);
      setDetail(res.data ?? null);
      setDetailError(null);
    } catch (err) {
      setDetailError(apiError(err, "Could not load this contract."));
    }
  }, []);

  useEffect(() => {
    void Promise.resolve().then(loadList);
  }, [loadList]);

  const select = (id: string) => {
    setSelectedId(id);
    setDetail(null);
    setConfirm(null);
    void loadDetail(id);
    if (window.matchMedia("(max-width: 1279px)").matches) setMobileShowDetail(true);
  };

  const run = async () => {
    if (!detail || !confirm) return;
    setBusy(true);
    try {
      if (confirm === "accept") {
        await influencerApi.acceptContract(detail.id);
        toast.success("Contract accepted. The work starts once the payment is confirmed.");
      } else {
        await influencerApi.rejectContract(detail.id);
        toast.success("Contract declined.");
      }
      setConfirm(null);
      await Promise.all([loadList(), loadDetail(detail.id)]);
    } catch (err) {
      toast.error(apiError(err, "The action could not be completed."));
    } finally {
      setBusy(false);
    }
  };

  const actions =
    detail && detail.status === "pending_influencer" ? (
      confirm ? (
        <div className="w-full space-y-2 rounded-xl border border-white/15 bg-black/20 p-3">
          <p className="text-sm">
            {confirm === "accept"
              ? `Accept this contract? You'll receive ${formatEtb(detail.fees.influencer_receivable_etb)} when the business approves your work.`
              : "Decline this contract? The business will be told you declined."}
          </p>
          <div className="flex flex-wrap gap-2">
            <Button
              type="button"
              disabled={busy}
              onClick={() => void run()}
              className={cn("h-10 rounded-xl text-white", confirm === "accept" ? "bg-emerald-600 hover:bg-emerald-700" : "bg-rose-600 hover:bg-rose-700")}
            >
              {busy ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
              {confirm === "accept" ? "Accept contract" : "Decline contract"}
            </Button>
            <Button type="button" variant="ghost" disabled={busy} onClick={() => setConfirm(null)} className="h-10 rounded-xl text-white hover:bg-white/10">
              Not now
            </Button>
          </div>
        </div>
      ) : (
        <>
          <Button type="button" onClick={() => setConfirm("accept")} className="h-10 rounded-xl bg-emerald-600 text-white hover:bg-emerald-700">
            Accept
          </Button>
          <Button type="button" variant="outline" onClick={() => setConfirm("decline")} className="h-10 rounded-xl border-rose-400/40 bg-transparent text-rose-200 hover:bg-rose-500/10">
            Decline
          </Button>
        </>
      )
    ) : null;

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold">Contracts</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">Review terms from businesses, accept the ones you want, and follow your payouts.</p>
      </div>

      <div className="flex flex-wrap gap-2" role="tablist" aria-label="Contract status">
        {FILTERS.map((f) => (
          <button
            key={f.id || "all"}
            type="button"
            role="tab"
            aria-selected={filter === f.id}
            onClick={() => { setFilter(f.id); setListState("loading"); }}
            className={cn(
              "inline-flex min-h-11 items-center rounded-full border px-4 py-2 text-sm font-medium transition-colors sm:min-h-0",
              filter === f.id ? "border-transparent bg-[color:var(--adfluence-orange)] text-white" : "border-white/10 bg-[rgba(255,255,255,0.06)] text-[#CBD5E1] hover:text-white",
            )}
          >
            {f.label}
          </button>
        ))}
      </div>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_440px]">
        <div className={cn("min-w-0 space-y-3", mobileShowDetail && "hidden xl:block")}>
          {listState === "loading" ? (
            <p className="flex items-center gap-2 p-6 text-sm text-[#94A3B8]"><Loader2 className="size-4 animate-spin" aria-hidden /> Loading contracts…</p>
          ) : listState === "error" ? (
            <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-6 text-sm text-rose-100">
              Could not load contracts.
              <Button type="button" variant="link" className="ml-2 h-auto p-0 text-white" onClick={() => void loadList()}>Try again</Button>
            </div>
          ) : rows.length === 0 ? (
            <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center text-sm text-[#CBD5E1]">
              No contracts here yet. Businesses send you a contract after you accept their campaign invite.
            </div>
          ) : (
            rows.map((r) => {
              const st = contractStatus(r.status, "influencer");
              return (
                <button
                  key={r.id}
                  type="button"
                  onClick={() => select(r.id)}
                  aria-current={selectedId === r.id ? "true" : undefined}
                  className={cn(
                    "flex w-full min-w-0 items-start justify-between gap-4 rounded-2xl border p-4 text-left transition-colors focus-visible:outline-2 focus-visible:outline-[color:var(--adfluence-orange)]",
                    selectedId === r.id ? "border-[color:var(--adfluence-orange)]/35 bg-[color:var(--adfluence-orange)]/10" : "border-white/10 bg-[rgba(255,255,255,0.08)] hover:border-white/20",
                  )}
                >
                  <div className="min-w-0">
                    <p className="font-semibold text-white">{r.title}</p>
                    <p className="text-xs text-[#CBD5E1]">{r.counterparty_name} · {r.campaign_title}</p>
                    <span className={cn("mt-2 inline-block rounded-full border px-2 py-0.5 text-[10px] font-medium", st.badgeClass)}>{st.label}</span>
                  </div>
                  <span className="shrink-0 text-sm font-semibold tabular-nums">{formatEtb(r.influencer_receivable_etb || r.amount_etb)}</span>
                </button>
              );
            })
          )}
        </div>

        {selectedId ? (
          <aside aria-label="Contract details" className={cn("min-w-0 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 xl:sticky xl:top-4 xl:self-start", !mobileShowDetail && "hidden xl:block")}>
            <Button type="button" variant="ghost" size="sm" onClick={() => setMobileShowDetail(false)} className="mb-3 -ml-2 text-white hover:bg-white/10 xl:hidden">
              <ArrowLeft className="mr-1 size-4" aria-hidden /> Back to contracts
            </Button>
            {detailError ? (
              <p role="alert" className="text-sm text-rose-200">{detailError}</p>
            ) : !detail ? (
              <p className="flex items-center gap-2 text-sm text-[#94A3B8]"><Loader2 className="size-4 animate-spin" aria-hidden /> Loading contract…</p>
            ) : (
              <ContractDetailView contract={detail} role="influencer" actions={actions} />
            )}
          </aside>
        ) : (
          <aside className="hidden rounded-2xl border border-dashed border-white/10 p-8 text-center text-sm text-[#94A3B8] xl:block">Select a contract to see its details.</aside>
        )}
      </div>
    </div>
  );
}
