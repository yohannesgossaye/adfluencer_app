"use client";

import Image from "next/image";
import Link from "next/link";
import { useMemo, useState } from "react";
import { FileDown, Plus } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  businessContractStats,
  initialBusinessContracts,
  type BusinessContractRecord,
  type BusinessContractRecordStatus,
} from "@/data/business-contracts-mock";
import { influencers } from "@/data/influencers";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet";
import { cn } from "@/lib/utils";

function statusBadgeClass(status: BusinessContractRecordStatus) {
  switch (status) {
    case "draft":
      return "bg-white/10 text-[#CBD5E1] border-white/15";
    case "pending":
      return "bg-[color:var(--adfluence-orange)]/15 text-amber-200 border-[color:var(--adfluence-orange)]/30";
    case "accepted":
      return "bg-[color:var(--adfluence-purple-mid)]/15 text-violet-200 border-[color:var(--adfluence-purple-accent)]/35";
    case "active":
      return "bg-emerald-500/15 text-emerald-200 border-emerald-400/30";
    case "completed":
      return "bg-sky-500/15 text-sky-200 border-sky-400/30";
    default:
      return "bg-white/10 text-[#CBD5E1]";
  }
}

function formatEtb(n: number) {
  if (n >= 1000) return `ETB ${Math.round(n / 1000)}K`;
  return `ETB ${n.toLocaleString()}`;
}

export function BusinessContracts() {
  const [contracts, setContracts] = useState<BusinessContractRecord[]>(() =>
    structuredClone(initialBusinessContracts),
  );
  const [selectedId, setSelectedId] = useState<string>(initialBusinessContracts[0]?.id ?? "");
  const [sheetOpen, setSheetOpen] = useState(false);
  const [form, setForm] = useState({
    influencerId: influencers[0]?.id ?? "",
    campaign: "",
    deliverables: "",
    timeline: "",
    amount: "",
    notes: "",
  });

  const selected = useMemo(
    () => contracts.find((c) => c.id === selectedId) ?? contracts[0],
    [contracts, selectedId],
  );

  const addContract = () => {
    const inf = influencers.find((i) => i.id === form.influencerId) ?? influencers[0];
    if (!inf) return;
    const amount = Number(form.amount.replace(/\D/g, "")) || 100000;
    const newRow: BusinessContractRecord = {
      id: `bctr-${Date.now()}`,
      title: `${form.campaign || "New campaign"} x ${inf.name}`,
      influencerName: inf.name,
      influencerImage: inf.imageUrl ?? "/images/influencers/influencer-1.jpg",
      influencerId: inf.id,
      campaign: form.campaign || "Untitled campaign",
      deliverables: form.deliverables || "TBD",
      amountEtb: amount,
      timeline: form.timeline || "TBD",
      status: "draft",
      brief: "Draft collaboration created from the business portal.",
      deliverableChecklist: [
        { id: "n1", label: "Deliverable 1", done: false },
        { id: "n2", label: "Deliverable 2", done: false },
      ],
      notes: form.notes || "",
    };
    setContracts((prev) => [newRow, ...prev]);
    setSelectedId(newRow.id);
    setSheetOpen(false);
    setForm({
      influencerId: influencers[0]?.id ?? "",
      campaign: "",
      deliverables: "",
      timeline: "",
      amount: "",
      notes: "",
    });
  };

  const updateStatus = (id: string, status: BusinessContractRecordStatus) => {
    setContracts((prev) => prev.map((c) => (c.id === id ? { ...c, status } : c)));
  };

  const toggleChecklist = (contractId: string, itemId: string) => {
    setContracts((prev) =>
      prev.map((c) => {
        if (c.id !== contractId) return c;
        return {
          ...c,
          deliverableChecklist: c.deliverableChecklist.map((row) =>
            row.id === itemId ? { ...row, done: !row.done } : row,
          ),
        };
      }),
    );
  };

  const mockDownload = () => {
    window.alert("Download started (demo) — PDF export will connect to backend later.");
  };

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold text-white">Contracts</h1>
          <p className="mt-1 max-w-xl text-sm text-[#CBD5E1]">
            Create, review, and manage influencer collaboration agreements.
          </p>
        </div>
        <Sheet open={sheetOpen} onOpenChange={setSheetOpen}>
          <SheetTrigger asChild>
            <Button className="h-11 rounded-xl bg-[color:var(--adfluence-orange)] px-5 text-white hover:bg-[color:var(--adfluence-orange-hover)]">
              <Plus className="mr-2 size-4" />
              New Contract
            </Button>
          </SheetTrigger>
          <SheetContent
            side="right"
            className="flex w-full max-w-md flex-col border-white/10 bg-[color:var(--adfluence-surface)] text-white"
          >
            <SheetHeader>
              <SheetTitle className="text-white">New contract</SheetTitle>
            </SheetHeader>
            <div className="mt-4 flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto pr-1">
              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Influencer</span>
                <select
                  value={form.influencerId}
                  onChange={(e) => setForm((f) => ({ ...f, influencerId: e.target.value }))}
                  className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-bg)] px-3 text-white outline-none focus:border-[color:var(--adfluence-orange)]"
                >
                  {influencers.map((i) => (
                    <option key={i.id} value={i.id} className="bg-[color:var(--adfluence-bg)]">
                      {i.name}
                    </option>
                  ))}
                </select>
              </label>
              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Campaign</span>
                <Input
                  value={form.campaign}
                  onChange={(e) => setForm((f) => ({ ...f, campaign: e.target.value }))}
                  placeholder="Campaign name"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-bg)] text-white placeholder:text-[#94A3B8]"
                />
              </label>
              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Deliverables</span>
                <Input
                  value={form.deliverables}
                  onChange={(e) => setForm((f) => ({ ...f, deliverables: e.target.value }))}
                  placeholder="e.g. 2 Reels, 4 stories"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-bg)] text-white placeholder:text-[#94A3B8]"
                />
              </label>
              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Timeline</span>
                <Input
                  value={form.timeline}
                  onChange={(e) => setForm((f) => ({ ...f, timeline: e.target.value }))}
                  placeholder="e.g. Jun 1 – Jun 20, 2026"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-bg)] text-white placeholder:text-[#94A3B8]"
                />
              </label>
              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Payment amount (ETB)</span>
                <Input
                  value={form.amount}
                  onChange={(e) => setForm((f) => ({ ...f, amount: e.target.value }))}
                  placeholder="e.g. 150000"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-bg)] text-white placeholder:text-[#94A3B8]"
                />
              </label>
              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Notes</span>
                <textarea
                  value={form.notes}
                  onChange={(e) => setForm((f) => ({ ...f, notes: e.target.value }))}
                  rows={3}
                  className="w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-bg)] px-3 py-2 text-white placeholder:text-[#94A3B8]"
                  placeholder="Internal notes"
                />
              </label>
            </div>
            <Button
              type="button"
              className="mt-4 h-11 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
              onClick={addContract}
            >
              Save draft contract
            </Button>
          </SheetContent>
        </Sheet>
      </section>

      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {[
          { label: "Pending Review", value: String(businessContractStats.pendingReview) },
          { label: "Active Contracts", value: String(businessContractStats.active) },
          { label: "Completed", value: String(businessContractStats.completed) },
          { label: "Total Value", value: formatEtb(businessContractStats.totalValueEtb) },
        ].map((card) => (
          <article
            key={card.label}
            className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
          >
            <p className="text-xs text-[#CBD5E1]">{card.label}</p>
            <p className="mt-2 text-2xl font-bold text-[color:var(--adfluence-orange)]">
              {card.value}
            </p>
          </article>
        ))}
      </section>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <section className="space-y-3">
          <h2 className="text-sm font-medium uppercase tracking-wide text-[#CBD5E1]">
            All contracts
          </h2>
          <div className="space-y-3 overflow-x-auto">
            {contracts.map((c) => {
              const rowActive = c.id === selectedId;
              return (
                <button
                  key={c.id}
                  type="button"
                  onClick={() => setSelectedId(c.id)}
                  className={cn(
                    "flex w-full min-w-0 flex-col gap-3 rounded-2xl border p-4 text-left transition-colors sm:flex-row sm:items-center",
                    rowActive
                      ? "border-[color:var(--adfluence-orange)]/35 bg-[color:var(--adfluence-orange)]/10"
                      : "border-white/10 bg-[rgba(255,255,255,0.08)] hover:border-white/20",
                  )}
                >
                  <div className="relative h-14 w-14 shrink-0 overflow-hidden rounded-full border border-white/10">
                    <Image
                      src={c.influencerImage}
                      alt={c.influencerName}
                      fill
                      className="object-cover"
                      sizes="56px"
                    />
                  </div>
                  <div className="min-w-0 flex-1">
                    <p className="font-semibold text-white">{c.title}</p>
                    <p className="text-xs text-[#CBD5E1]">{c.campaign}</p>
                    <p className="mt-1 text-xs text-[#94A3B8]">{c.deliverables}</p>
                    <div className="mt-2 flex flex-wrap gap-2 text-xs text-[#CBD5E1]">
                      <span>{formatEtb(c.amountEtb)}</span>
                      <span>·</span>
                      <span>{c.timeline}</span>
                    </div>
                  </div>
                  <div className="flex shrink-0 flex-col items-start gap-2 sm:items-end">
                    <span
                      className={cn(
                        "rounded-full border px-2.5 py-0.5 text-xs font-medium capitalize",
                        statusBadgeClass(c.status),
                      )}
                    >
                      {c.status}
                    </span>
                    <div className="flex flex-wrap gap-2">
                      <Button
                        asChild
                        size="sm"
                        className="h-8 rounded-lg bg-white/10 px-3 text-xs text-white hover:bg-white/15"
                      >
                        <Link href={`/business/influencers/${c.influencerId}`}>Profile</Link>
                      </Button>
                      <Button
                        asChild
                        size="sm"
                        className="h-8 rounded-lg bg-white/10 px-3 text-xs text-white hover:bg-white/15"
                      >
                        <Link href="/business/messages">Message</Link>
                      </Button>
                    </div>
                  </div>
                </button>
              );
            })}
          </div>
        </section>

        {selected ? (
          <aside className="xl:sticky xl:top-4 xl:self-start">
            <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
              <h3 className="text-lg font-semibold text-white">Contract detail</h3>
              <p className="mt-1 text-sm text-[#CBD5E1]">{selected.title}</p>
              <span
                className={cn(
                  "mt-3 inline-flex rounded-full border px-2.5 py-0.5 text-xs font-medium capitalize",
                  statusBadgeClass(selected.status),
                )}
              >
                {selected.status}
              </span>

              <div className="mt-4 space-y-2 text-sm text-[#CBD5E1]">
                <p>
                  <span className="text-white">Campaign brief:</span> {selected.brief}
                </p>
                <p>
                  <span className="text-white">Amount:</span> {formatEtb(selected.amountEtb)}
                </p>
                <p>
                  <span className="text-white">Timeline:</span> {selected.timeline}
                </p>
                <p>
                  <span className="text-white">Notes:</span> {selected.notes || "—"}
                </p>
              </div>

              <div className="mt-4">
                <p className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                  Deliverables
                </p>
                <ul className="mt-2 space-y-2">
                  {selected.deliverableChecklist.map((item) => (
                    <li key={item.id} className="flex items-center gap-2 text-sm text-[#E2E8F0]">
                      <input
                        type="checkbox"
                        checked={item.done}
                        onChange={() => toggleChecklist(selected.id, item.id)}
                        className="h-4 w-4 rounded border-white/20 bg-[color:var(--adfluence-surface)] accent-[color:var(--adfluence-orange)]"
                      />
                      {item.label}
                    </li>
                  ))}
                </ul>
              </div>

              <div className="mt-5 flex flex-col gap-2">
                <Button
                  type="button"
                  className="h-10 w-full rounded-xl bg-white/10 text-white hover:bg-white/15"
                  onClick={() => updateStatus(selected.id, "pending")}
                  disabled={selected.status !== "draft"}
                >
                  Send Contract
                </Button>
                <Button
                  type="button"
                  className="h-10 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] disabled:opacity-40"
                  onClick={() => updateStatus(selected.id, "active")}
                  disabled={
                    selected.status === "active" ||
                    selected.status === "completed" ||
                    selected.status === "draft"
                  }
                >
                  Mark Active
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  className="h-10 w-full rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
                  onClick={mockDownload}
                >
                  <FileDown className="mr-2 size-4" />
                  Download PDF
                </Button>
              </div>
            </div>
          </aside>
        ) : null}
      </div>
    </div>
  );
}
