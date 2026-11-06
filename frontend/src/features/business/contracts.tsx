"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import axios from "axios";
import { ArrowLeft, Loader2, Plus } from "lucide-react";
import { toast } from "sonner";

import { ContractDetailView } from "@/components/contracts/contract-detail";
import { ReportProblemButton } from "@/components/contracts/report-problem";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Sheet, SheetContent, SheetDescription, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import {
  businessApi,
  type CampaignListItem,
  type InviteResponse,
} from "@/lib/api/business";
import { DELIVERABLE_TYPES, contractStatus, formatEtb, previewFees } from "@/lib/contracts";
import { cn } from "@/lib/utils";
import { StatusBadge } from "@/components/common/status-badge";
import { FilterChips } from "@/components/common/filter-chips";
import { loadState, useContract, useContracts, useInvalidateContracts } from "@/lib/queries";

const FILTERS = [
  { id: "", label: "All" },
  { id: "pending_influencer", label: "Waiting for creator" },
  { id: "accepted", label: "Awaiting payment" },
  { id: "active", label: "In progress" },
  { id: "disputed", label: "In dispute" },
  { id: "completed", label: "Completed" },
] as const;

type Confirm = { kind: "approve" | "cancel" } | null;

function apiError(err: unknown, fallback: string) {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as { message?: string; error?: { type?: string; details?: Array<{ description?: string }> } } | undefined;
    switch (data?.error?.type) {
      case "account_not_verified":
        return "Your business must be verified to create or manage contracts.";
      case "invite_not_accepted":
        return "The creator has not accepted this invite yet.";
      case "contract_already_exists":
        return "A contract already exists for this invite.";
    }
    return data?.error?.details?.[0]?.description || data?.message || (!err.response ? "You appear to be offline." : fallback);
  }
  return fallback;
}

const inputClass =
  "h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25";
const selectClass =
  "h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25";

function NewContractSheet({
  open,
  onOpenChange,
  onCreated,
  preset,
}: {
  open: boolean;
  onOpenChange: (v: boolean) => void;
  onCreated: (id: string) => void;
  /** Campaign and accepted invite to preselect, e.g. from "Send contract". */
  preset?: { campaignId: string; inviteId: string };
}) {
  const [campaigns, setCampaigns] = useState<CampaignListItem[]>([]);
  const [invites, setInvites] = useState<InviteResponse[]>([]);
  const [campaignId, setCampaignId] = useState("");
  const [inviteId, setInviteId] = useState("");
  const [title, setTitle] = useState("");
  const [amount, setAmount] = useState("");
  const [quantities, setQuantities] = useState<Record<string, number>>({});
  const [startDate, setStartDate] = useState("");
  const [endDate, setEndDate] = useState("");
  const [notes, setNotes] = useState("");
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    if (!open) return;
    businessApi
      .listCampaigns({ status: "active", page: 1, limit: 50 })
      .then((res) => setCampaigns(res.data?.items ?? []))
      .catch(() => setCampaigns([]));
    if (preset) {
      void Promise.resolve().then(() => {
        setCampaignId(preset.campaignId);
        setInviteId(preset.inviteId);
      });
    }
  }, [open, preset]);

  useEffect(() => {
    if (!campaignId) return;
    businessApi
      .listInvites(campaignId)
      .then((res) => setInvites((res.data?.items ?? []).filter((i) => i.status === "accepted")))
      .catch(() => setInvites([]));
  }, [campaignId]);

  const amountNumber = Number(amount);
  const fees = previewFees(Number.isFinite(amountNumber) ? amountNumber : 0);
  const deliverables = Object.entries(quantities)
    .filter(([, q]) => q > 0)
    .map(([type, quantity]) => ({ type, quantity }));

  const problem = !campaignId
    ? "Choose a live campaign."
    : !inviteId
      ? "Choose a creator who accepted your invite."
      : title.trim().length < 3
        ? "Give the contract a title (at least 3 characters)."
        : !(amountNumber >= 1)
          ? "Enter the amount you agreed with the creator."
          : deliverables.length === 0
            ? "Add at least one deliverable."
            : startDate && endDate && endDate < startDate
              ? "The end date must be on or after the start date."
              : null;

  const create = async () => {
    if (problem) return;
    setSaving(true);
    try {
      const res = await businessApi.createContract({
        invite_id: inviteId,
        title: title.trim(),
        description: notes.trim() || undefined,
        amount_etb: amountNumber,
        deliverables,
        start_date: startDate || undefined,
        end_date: endDate || undefined,
      });
      if (res.data?.id) {
        toast.success("Contract sent to the creator.");
        onCreated(res.data.id);
        onOpenChange(false);
        setInviteId("");
        setTitle("");
        setAmount("");
        setQuantities({});
        setNotes("");
      }
    } catch (err) {
      toast.error(apiError(err, "Could not create the contract."));
    } finally {
      setSaving(false);
    }
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-full overflow-y-auto border-white/10 bg-[color:var(--adfluence-bg)] text-white sm:max-w-lg">
        <SheetHeader>
          <SheetTitle className="text-white">New contract</SheetTitle>
          <SheetDescription className="text-[#94A3B8]">
            Send terms to a creator who accepted your invite. Once they accept, Adfluence confirms your payment and the work starts.
          </SheetDescription>
        </SheetHeader>
        <div className="space-y-4 px-4 pb-6">
          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">Campaign *</span>
            <select id="contract-campaign" value={campaignId} onChange={(e) => { setCampaignId(e.target.value); setInviteId(""); }} className={selectClass}>
              <option value="">{campaigns.length ? "Select a live campaign" : "No live campaigns"}</option>
              {campaigns.map((c) => (
                <option key={c.id} value={c.id}>{c.title}</option>
              ))}
            </select>
          </label>
          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">Creator *</span>
            <select id="contract-invite" value={inviteId} disabled={!campaignId} onChange={(e) => {
              setInviteId(e.target.value);
              const inv = invites.find((i) => i.id === e.target.value);
              const camp = campaigns.find((c) => c.id === campaignId);
              if (inv && !title) setTitle(`${camp?.title ?? "Campaign"} × ${inv.creator_name ?? "creator"}`);
            }} className={selectClass}>
              <option value="">{campaignId && invites.length === 0 ? "No creator has accepted an invite yet" : "Select a creator"}</option>
              {invites.map((i) => (
                <option key={i.id} value={i.id}>{i.creator_name || "Creator"}</option>
              ))}
            </select>
          </label>
          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">Title *</span>
            <Input id="contract-title" value={title} maxLength={255} onChange={(e) => setTitle(e.target.value)} className={inputClass} />
          </label>
          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">Agreed amount (ETB) *</span>
            <Input id="contract-amount" value={amount} inputMode="decimal" onChange={(e) => setAmount(e.target.value)} className={inputClass} placeholder="e.g. 20000" />
          </label>
          {amountNumber >= 1 ? (
            <div className="rounded-xl border border-white/10 bg-black/20 p-3 text-sm">
              <div className="flex justify-between"><span className="text-[#94A3B8]">You pay (amount + 5%)</span><span className="tabular-nums">{formatEtb(fees.totalBusinessPays)}</span></div>
              <div className="flex justify-between"><span className="text-[#94A3B8]">Creator receives (amount − 5%)</span><span className="tabular-nums">{formatEtb(fees.influencerReceives)}</span></div>
            </div>
          ) : null}
          <fieldset className="space-y-2">
            <legend className="text-sm text-[#CBD5E1]">Deliverables *</legend>
            {DELIVERABLE_TYPES.map((d) => (
              <div key={d.id} className="flex items-center justify-between gap-3 rounded-xl border border-white/10 px-3 py-2">
                <span className="text-sm">{d.label}</span>
                <Input
                  id={`qty-${d.id}`}
                  aria-label={`${d.label} quantity`}
                  type="number"
                  min={0}
                  max={20}
                  value={quantities[d.id] ?? 0}
                  onChange={(e) => setQuantities((q) => ({ ...q, [d.id]: Math.max(0, Math.min(20, Number(e.target.value) || 0)) }))}
                  className="h-9 w-20 rounded-lg border-white/10 bg-[color:var(--adfluence-surface)] text-white"
                />
              </div>
            ))}
          </fieldset>
          <div className="grid grid-cols-2 gap-3">
            <label className="space-y-1">
              <span className="text-sm text-[#CBD5E1]">Start date</span>
              <Input id="contract-start" type="date" value={startDate} onChange={(e) => setStartDate(e.target.value)} className={inputClass} />
            </label>
            <label className="space-y-1">
              <span className="text-sm text-[#CBD5E1]">End date</span>
              <Input id="contract-end" type="date" value={endDate} onChange={(e) => setEndDate(e.target.value)} className={inputClass} />
            </label>
          </div>
          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">Notes for the creator</span>
            <textarea id="contract-notes" value={notes} maxLength={2000} rows={3} onChange={(e) => setNotes(e.target.value)} className="w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-white outline-none focus:border-[color:var(--adfluence-orange)]" />
          </label>
          {problem ? <p className="text-xs text-[#94A3B8]">{problem}</p> : null}
          <Button type="button" disabled={!!problem || saving} onClick={() => void create()} className="h-11 w-full rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)]">
            {saving ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
            Send contract to creator
          </Button>
        </div>
      </SheetContent>
    </Sheet>
  );
}

export function BusinessContracts() {
  const [filter, setFilter] = useState<string>("");
  const searchParams = useSearchParams();
  // Emails link here with ?contract_id= to open that contract.
  const linkedId = searchParams.get("contract_id");
  const [selectedId, setSelectedId] = useState<string | null>(linkedId);
  const [mobileShowDetail, setMobileShowDetail] = useState(Boolean(linkedId));
  // "Send contract" on an accepted invite links here with both ids.
  const preset = useMemo(() => {
    const campaignId = searchParams.get("campaign_id");
    const inviteId = searchParams.get("invite_id");
    return campaignId && inviteId ? { campaignId, inviteId } : undefined;
  }, [searchParams]);
  const [sheetOpen, setSheetOpen] = useState(Boolean(preset));
  const [confirm, setConfirm] = useState<Confirm>(null);
  const [busy, setBusy] = useState(false);

  // Cached: switching filters or contracts shows what was already loaded at once.
  const listQuery = useContracts("business", filter);
  const rows = listQuery.data ?? [];
  const listState = loadState(listQuery);
  const detailQuery = useContract("business", selectedId);
  const detail = detailQuery.data ?? null;
  const detailError = detailQuery.isError && !detail ? apiError(detailQuery.error, "Could not load this contract.") : null;
  const refresh = useInvalidateContracts("business");

  const select = (id: string) => {
    setSelectedId(id);
    setConfirm(null);
    if (window.matchMedia("(max-width: 1279px)").matches) setMobileShowDetail(true);
  };

  const run = async () => {
    if (!detail || !confirm) return;
    setBusy(true);
    try {
      if (confirm.kind === "approve") {
        await businessApi.approveWork(detail.id);
        toast.success("Work approved. The creator's payout has been released.");
      } else {
        await businessApi.cancelContract(detail.id);
        toast.success("Contract cancelled.");
      }
      setConfirm(null);
      await refresh(detail.id);
    } catch (err) {
      toast.error(apiError(err, "The action could not be completed."));
    } finally {
      setBusy(false);
    }
  };

  const confirmCopy = useMemo(() => {
    if (!detail || !confirm) return null;
    if (confirm.kind === "approve")
      return {
        text: `Approve the delivered work? This releases ${formatEtb(detail.fees.influencer_receivable_etb)} to the creator and completes the contract.`,
        button: "Approve and release payout",
        tone: "bg-[color:var(--adfluence-action)] hover:bg-[color:var(--adfluence-action-hover)]",
      };
    return { text: "Cancel this contract? The creator will see it as cancelled.", button: "Cancel contract", tone: "bg-rose-600 hover:bg-rose-700" };
  }, [detail, confirm]);

  const actions = detail ? (
    confirm && confirmCopy ? (
      <div className="w-full space-y-2 rounded-xl border border-white/15 bg-black/20 p-3">
        <p className="text-sm">{confirmCopy.text}</p>
        <div className="flex flex-wrap gap-2">
          <Button type="button" disabled={busy} onClick={() => void run()} className={cn("h-10 rounded-xl text-white", confirmCopy.tone)}>
            {busy ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
            {confirmCopy.button}
          </Button>
          <Button type="button" variant="ghost" disabled={busy} onClick={() => setConfirm(null)} className="h-10 rounded-xl text-white hover:bg-white/10">
            Keep as is
          </Button>
        </div>
      </div>
    ) : (
      <>
        {detail.status === "accepted" ? (
          <p className="w-full rounded-xl border border-amber-400/30 bg-amber-500/10 p-3 text-sm text-amber-50">
            Waiting for Adfluence to confirm your payment of {formatEtb(detail.fees.total_business_payable_etb)}. The creator
            starts once it is confirmed. Already paid?{" "}
            <Link href="/business/support" className="underline">Send us the payment reference</Link>.
          </p>
        ) : null}
        {detail.status === "active" ? (
          <Button type="button" onClick={() => setConfirm({ kind: "approve" })} className="h-10 rounded-xl bg-[color:var(--adfluence-action)] hover:bg-[color:var(--adfluence-action-hover)] text-white">
            Approve work and release payout
          </Button>
        ) : null}
        {detail.status === "active" ? (
          <ReportProblemButton contractId={detail.id} role="business" onReported={() => refresh(detail.id).then(() => {})} />
        ) : null}
        {detail.status === "pending_influencer" || detail.status === "accepted" ? (
          <Button type="button" variant="outline" onClick={() => setConfirm({ kind: "cancel" })} className="h-10 rounded-xl border-rose-400/40 bg-transparent text-rose-200 hover:bg-rose-500/10">
            Cancel contract
          </Button>
        ) : null}
        <Button asChild variant="ghost" className="h-10 rounded-xl text-white hover:bg-white/10">
          <Link href={`/business/influencers/${detail.influencer.user_id}`}>View creator</Link>
        </Button>
      </>
    )
  ) : null;

  return (
    <div className="space-y-6 text-white">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold">Contracts</h1>
          <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">Send terms to creators, follow payment and release payouts when the work is done.</p>
        </div>
        <Button type="button" onClick={() => setSheetOpen(true)} className="h-11 rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)]">
          <Plus className="mr-1 size-4" aria-hidden /> New contract
        </Button>
      </div>

      <FilterChips
        label="Contract status"
        value={filter}
        onChange={(v) => {
          setFilter(v);
        }}
        options={FILTERS.map((f) => ({ value: f.id, label: f.label }))}
      />

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_440px]">
        <div className={cn("min-w-0 space-y-3", mobileShowDetail && "hidden xl:block")}>
          {listState === "loading" ? (
            <p className="flex items-center gap-2 p-6 text-sm text-[#94A3B8]"><Loader2 className="size-4 animate-spin" aria-hidden /> Loading contracts…</p>
          ) : listState === "error" ? (
            <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-6 text-sm text-rose-100">
              Could not load contracts.
              <Button type="button" variant="link" className="ml-2 h-auto p-0 text-white" onClick={() => void listQuery.refetch()}>Try again</Button>
            </div>
          ) : rows.length === 0 ? (
            <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center text-sm text-[#CBD5E1]">
              No contracts here yet. When a creator accepts your invite, use <span className="text-white">New contract</span> to send them terms.
            </div>
          ) : (
            rows.map((r) => {
              const st = contractStatus(r.status, "business");
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
                    <StatusBadge status={st.badge} className="mt-2" />
                  </div>
                  <span className="shrink-0 text-sm font-semibold tabular-nums">{formatEtb(r.total_business_payable_etb || r.amount_etb)}</span>
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
              <ContractDetailView contract={detail} role="business" actions={actions} />
            )}
          </aside>
        ) : (
          <aside className="hidden rounded-2xl border border-dashed border-white/10 p-8 text-center text-sm text-[#94A3B8] xl:block">Select a contract to see its details.</aside>
        )}
      </div>

      <NewContractSheet
        open={sheetOpen}
        preset={preset}
        onOpenChange={setSheetOpen}
        onCreated={(id) => {
          setFilter("");
          void refresh().then(() => select(id));
        }}
      />
    </div>
  );
}
