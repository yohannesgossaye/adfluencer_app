"use client";

import Image from "next/image";
import Link from "next/link";
import { useCallback, useEffect, useMemo, useState } from "react";
import { ArrowLeft, FileDown, Plus, CheckCircle, ShieldCheck, AlertCircle } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Sheet, SheetContent, SheetHeader, SheetTitle, SheetTrigger } from "@/components/ui/sheet";
import {
  businessApi,
  type ContractDetail,
  type CampaignListItem,
  type InviteResponse,
} from "@/lib/api/business";
import { cn } from "@/lib/utils";
import { UserAvatar } from "@/components/common/user-avatar";

function statusBadgeClass(status: string) {
  switch (status.toLowerCase()) {
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
    case "cancelled":
      return "bg-rose-500/15 text-rose-200 border-rose-400/30";
    default:
      return "bg-white/10 text-[#CBD5E1] border-white/15";
  }
}

function formatEtb(n: number) {
  if (n >= 1000) return `ETB ${Math.round(n / 1000)}K`;
  return `ETB ${(n || 0).toLocaleString()}`;
}

export function BusinessContracts() {
  const [contracts, setContracts] = useState<ContractDetail[]>([]);
  const [selectedId, setSelectedId] = useState<string>("");
  const [isLoading, setIsLoading] = useState(true);
  const [actionLoading, setActionLoading] = useState(false);
  const [statusMessage, setStatusMessage] = useState("");
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [sheetOpen, setSheetOpen] = useState(false);

  // New contract form dependencies
  const [campaigns, setCampaigns] = useState<CampaignListItem[]>([]);
  const [invites, setInvites] = useState<InviteResponse[]>([]);
  const [form, setForm] = useState({
    campaignId: "",
    inviteId: "",
    title: "",
    deliverables: "1 TikTok Video, 2 Instagram Stories",
    amountEtb: "25000",
    startDate: new Date().toISOString().slice(0, 10),
    endDate: new Date(Date.now() + 14 * 86400000).toISOString().slice(0, 10),
    notes: "",
  });

  const loadContracts = useCallback(async () => {
    try {
      const res = await businessApi.listContracts();
      if (res.data?.items) {
        setContracts(res.data.items);
        if (res.data.items.length > 0 && !selectedId) {
          setSelectedId(res.data.items[0].id);
        }
      }
    } catch {
      // Fallback
    } finally {
      setIsLoading(false);
    }
  }, [selectedId]);

  useEffect(() => {
    loadContracts();
  }, [loadContracts]);

  // Load campaigns when opening new contract sheet
  useEffect(() => {
    if (!sheetOpen) return;
    async function loadFormContext() {
      try {
        const cRes = await businessApi.listCampaigns({ page: 1, limit: 50 });
        const items = cRes.data?.items;
        if (items && items.length > 0) {
          setCampaigns(items);
          if (!form.campaignId) {
            setForm((prev) => ({
              ...prev,
              campaignId: items[0].id,
              title: `${items[0].title} Collaboration`,
            }));
          }
        }
      } catch {
        // Fallback
      }
    }
    loadFormContext();
  }, [sheetOpen, form.campaignId]);

  // Load invites for the selected campaign in form
  useEffect(() => {
    if (!form.campaignId) return;
    async function loadCampaignInvites() {
      try {
        const invRes = await businessApi.listInvites(form.campaignId);
        const items = invRes.data?.items;
        if (items) {
          setInvites(items);
          if (items.length > 0 && !form.inviteId) {
            setForm((prev) => ({ ...prev, inviteId: items[0].id }));
          }
        }
      } catch {
        setInvites([]);
      }
    }
    loadCampaignInvites();
  }, [form.campaignId, form.inviteId]);

  const selected = useMemo(
    () => contracts.find((c) => c.id === selectedId) ?? contracts[0],
    [contracts, selectedId],
  );

  const stats = useMemo(() => {
    let pending = 0;
    let active = 0;
    let completed = 0;
    let totalValue = 0;

    for (const c of contracts) {
      const s = c.status.toLowerCase();
      if (s === "pending" || s === "draft") pending++;
      else if (s === "active") active++;
      else if (s === "completed") completed++;

      totalValue += c.fees?.amount_etb || c.fees?.total_business_payable_etb || 0;
    }

    return { pending, active, completed, totalValue };
  }, [contracts]);

  const selectContract = useCallback((id: string) => {
    setSelectedId(id);
    if (typeof window !== "undefined" && window.matchMedia("(max-width: 1279px)").matches) {
      setMobileShowDetail(true);
    }
  }, []);

  const handleCreateContract = async () => {
    if (!form.campaignId || !form.inviteId) {
      setStatusMessage("Please select a campaign and an invited creator.");
      setTimeout(() => setStatusMessage(""), 3000);
      return;
    }

    const amount = Number(form.amountEtb.replace(/\D/g, "")) || 20000;
    const deliverablesList = form.deliverables
      .split(",")
      .map((d) => d.trim())
      .filter(Boolean)
      .map((d) => ({
        type: d,
        quantity: 1,
        description: d,
      }));

    setActionLoading(true);
    try {
      const res = await businessApi.createContract({
        invite_id: form.inviteId,
        title: form.title || "Influencer Collaboration Agreement",
        description: form.notes || undefined,
        amount_etb: amount,
        deliverables:
          deliverablesList.length > 0
            ? deliverablesList
            : [{ type: "Content Post", quantity: 1 }],
        start_date: form.startDate,
        end_date: form.endDate,
      });

      const newContract = res.data?.contract;
      if (newContract) {
        setContracts((prev) => [newContract, ...prev]);
        setSelectedId(newContract.id);
        setSheetOpen(false);
        setStatusMessage("Contract created successfully!");
        setTimeout(() => setStatusMessage(""), 3000);
      }
    } catch {
      setStatusMessage("Failed to create contract. Ensure invited creator is valid.");
      setTimeout(() => setStatusMessage(""), 3500);
    } finally {
      setActionLoading(false);
    }
  };

  const handlePayEscrow = async (contractId: string) => {
    setActionLoading(true);
    try {
      const res = await businessApi.payContract(contractId, "telebirr");
      const updatedContract = res.data?.contract;
      if (updatedContract) {
        setContracts((prev) => prev.map((c) => (c.id === contractId ? updatedContract : c)));
        setStatusMessage("Escrow payment funded successfully! Contract is now Active.");
        setTimeout(() => setStatusMessage(""), 3500);
      }
    } catch {
      setStatusMessage("Payment simulation failed.");
      setTimeout(() => setStatusMessage(""), 3500);
    } finally {
      setActionLoading(false);
    }
  };

  const handleApproveWork = async (contractId: string) => {
    setActionLoading(true);
    try {
      const res = await businessApi.approveWork(contractId);
      const updatedContract = res.data?.contract;
      if (updatedContract) {
        setContracts((prev) => prev.map((c) => (c.id === contractId ? updatedContract : c)));
        setStatusMessage("Work approved! Escrow funds released to creator.");
        setTimeout(() => setStatusMessage(""), 3500);
      }
    } catch {
      setStatusMessage("Failed to approve work.");
      setTimeout(() => setStatusMessage(""), 3500);
    } finally {
      setActionLoading(false);
    }
  };

  const handleDownloadPDF = () => {
    window.print();
  };

  return (
    <div className="space-y-6 text-white">
      <section className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <h1 className="text-2xl font-semibold text-white">Contracts</h1>
          <p className="mt-1 max-w-xl text-sm text-[#CBD5E1]">
            Create, review, and manage influencer collaboration agreements and escrow payments.
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
              <SheetTitle className="text-white">New Collaboration Contract</SheetTitle>
            </SheetHeader>
            <div className="mt-4 flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto pr-1">
              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Campaign</span>
                {campaigns.length > 0 ? (
                  <select
                    value={form.campaignId}
                    onChange={(e) => {
                      const cId = e.target.value;
                      const cObj = campaigns.find((c) => c.id === cId);
                      setForm((f) => ({
                        ...f,
                        campaignId: cId,
                        title: cObj ? `${cObj.title} Agreement` : f.title,
                      }));
                    }}
                    className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-bg)] px-3 text-white outline-none focus:border-[color:var(--adfluence-orange)]"
                  >
                    {campaigns.map((c) => (
                      <option key={c.id} value={c.id} className="bg-[color:var(--adfluence-bg)]">
                        {c.title} ({c.status})
                      </option>
                    ))}
                  </select>
                ) : (
                  <p className="text-xs text-amber-200">No campaigns available. Create a campaign first.</p>
                )}
              </label>

              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Invited Creator</span>
                {invites.length > 0 ? (
                  <select
                    value={form.inviteId}
                    onChange={(e) => setForm((f) => ({ ...f, inviteId: e.target.value }))}
                    className="h-11 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-bg)] px-3 text-white outline-none focus:border-[color:var(--adfluence-orange)]"
                  >
                    {invites.map((inv) => (
                      <option
                        key={inv.id}
                        value={inv.id}
                        className="bg-[color:var(--adfluence-bg)]"
                      >
                        {inv.creator_name || "Creator"} ({inv.status})
                      </option>
                    ))}
                  </select>
                ) : (
                  <p className="text-xs text-amber-200">
                    No invited creators for this campaign. Invite creators from Recommendations first.
                  </p>
                )}
              </label>

              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Contract Title</span>
                <Input
                  value={form.title}
                  onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
                  placeholder="e.g. Betu Cosmetics x Hana Launch"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-bg)] text-white placeholder:text-[#94A3B8]"
                />
              </label>

              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Deliverables (comma separated)</span>
                <Input
                  value={form.deliverables}
                  onChange={(e) => setForm((f) => ({ ...f, deliverables: e.target.value }))}
                  placeholder="e.g. 1 TikTok Video, 2 Instagram Stories"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-bg)] text-white placeholder:text-[#94A3B8]"
                />
              </label>

              <div className="grid grid-cols-2 gap-3">
                <label className="space-y-2 text-sm">
                  <span className="text-[#CBD5E1]">Start Date</span>
                  <Input
                    type="date"
                    value={form.startDate}
                    onChange={(e) => setForm((f) => ({ ...f, startDate: e.target.value }))}
                    className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-bg)] text-white"
                  />
                </label>
                <label className="space-y-2 text-sm">
                  <span className="text-[#CBD5E1]">End Date</span>
                  <Input
                    type="date"
                    value={form.endDate}
                    onChange={(e) => setForm((f) => ({ ...f, endDate: e.target.value }))}
                    className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-bg)] text-white"
                  />
                </label>
              </div>

              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Payment Amount (ETB)</span>
                <Input
                  value={form.amountEtb}
                  onChange={(e) => setForm((f) => ({ ...f, amountEtb: e.target.value }))}
                  placeholder="e.g. 25000"
                  className="h-11 rounded-xl border-white/10 bg-[color:var(--adfluence-bg)] text-white placeholder:text-[#94A3B8]"
                />
              </label>

              <label className="space-y-2 text-sm">
                <span className="text-[#CBD5E1]">Notes / Instructions</span>
                <textarea
                  value={form.notes}
                  onChange={(e) => setForm((f) => ({ ...f, notes: e.target.value }))}
                  rows={3}
                  className="w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-bg)] px-3 py-2 text-white placeholder:text-[#94A3B8]"
                  placeholder="Brand guidelines or specific requirements"
                />
              </label>
            </div>

            <Button
              type="button"
              disabled={actionLoading || !form.campaignId || !form.inviteId}
              className="mt-4 h-11 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
              onClick={handleCreateContract}
            >
              {actionLoading ? "Creating Contract..." : "Create & Issue Contract"}
            </Button>
          </SheetContent>
        </Sheet>
      </section>

      {statusMessage ? (
        <p className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 px-4 py-2.5 text-sm text-emerald-200">
          {statusMessage}
        </p>
      ) : null}

      {/* Contract KPIs */}
      <section className="grid gap-3 sm:grid-cols-2 xl:grid-cols-4">
        {[
          { label: "Pending Review", value: String(stats.pending) },
          { label: "Active Contracts", value: String(stats.active) },
          { label: "Completed", value: String(stats.completed) },
          { label: "Total Value", value: formatEtb(stats.totalValue) },
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
        <section className={cn("space-y-3", mobileShowDetail && "hidden xl:block")}>
          <h2 className="text-sm font-medium uppercase tracking-wide text-[#CBD5E1]">
            All contracts
          </h2>

          {contracts.length > 0 ? (
            <div className="space-y-3">
              {contracts.map((c) => {
                const rowActive = c.id === selectedId;
                const influencerName = c.influencer?.creator_name || "Creator";
                const campaignTitle = c.campaign?.title || "Campaign Collaboration";
                const deliverablesSummary = c.deliverables?.length
                  ? c.deliverables.map((d) => `${d.quantity} ${d.type}`).join(", ")
                  : "Campaign deliverables";

                return (
                  <button
                    key={c.id}
                    type="button"
                    onClick={() => selectContract(c.id)}
                    className={cn(
                      "flex w-full min-w-0 flex-col gap-3 rounded-2xl border p-4 text-left transition-colors sm:flex-row sm:items-center",
                      rowActive
                        ? "border-[color:var(--adfluence-orange)]/35 bg-[color:var(--adfluence-orange)]/10"
                        : "border-white/10 bg-[rgba(255,255,255,0.08)] hover:border-white/20",
                    )}
                  >
                    <UserAvatar
                      src={c.influencer?.profile_image_url}
                      name={influencerName}
                      size="lg"
                      fallback="IN"
                    />
                    <div className="min-w-0 flex-1">
                      <p className="font-semibold text-white truncate">{c.title}</p>
                      <p className="text-xs text-[#CBD5E1] truncate">{campaignTitle} · {influencerName}</p>
                      <p className="mt-1 text-xs text-[#94A3B8] truncate">{deliverablesSummary}</p>
                      <div className="mt-2 flex flex-wrap gap-2 text-xs text-[#CBD5E1]">
                        <span>{formatEtb(c.fees?.amount_etb || 0)}</span>
                        <span>·</span>
                        <span>
                          {c.start_date
                            ? `${new Date(c.start_date).toLocaleDateString()} - ${new Date(c.end_date || "").toLocaleDateString()}`
                            : "Flexible"}
                        </span>
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
                        {c.influencer?.user_id ? (
                          <Button
                            asChild
                            size="sm"
                            className="h-11 rounded-lg bg-white/10 px-3 text-xs text-white hover:bg-white/15 sm:h-8"
                          >
                            <Link href={`/business/influencers/${c.influencer.user_id}`}>Profile</Link>
                          </Button>
                        ) : null}
                        <Button
                          asChild
                          size="sm"
                          className="h-11 rounded-lg bg-white/10 px-3 text-xs text-white hover:bg-white/15 sm:h-8"
                        >
                          <Link href="/business/messages">Message</Link>
                        </Button>
                      </div>
                    </div>
                  </button>
                );
              })}
            </div>
          ) : (
            <div className="rounded-2xl border border-white/10 bg-white/5 p-12 text-center">
              <p className="font-medium text-white">No contracts created yet</p>
              <p className="mt-1 text-sm text-[#CBD5E1]">
                Invite creators from your campaigns and click &quot;New Contract&quot; to issue collaboration agreements.
              </p>
            </div>
          )}
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

              <h3 className="text-lg font-semibold text-white">{selected.title}</h3>
              <p className="mt-1 text-sm text-[#CBD5E1]">
                {selected.campaign?.title} × {selected.influencer?.creator_name || "Creator"}
              </p>
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
                  <span className="text-white">Contract Amount:</span>{" "}
                  {formatEtb(selected.fees?.amount_etb || 0)}
                </p>
                <p>
                  <span className="text-white">Platform Fee (5%):</span>{" "}
                  {formatEtb(selected.fees?.business_fee_amount_etb || 0)}
                </p>
                <p>
                  <span className="text-white font-medium">Total Escrow Payable:</span>{" "}
                  <span className="font-semibold text-white">
                    {formatEtb(selected.fees?.total_business_payable_etb || selected.fees?.amount_etb || 0)}
                  </span>
                </p>
                <p>
                  <span className="text-white">Timeline:</span>{" "}
                  {selected.start_date
                    ? `${new Date(selected.start_date).toLocaleDateString()} - ${new Date(selected.end_date || "").toLocaleDateString()}`
                    : "Flexible"}
                </p>
                <p>
                  <span className="text-white">Notes:</span> {selected.description || "—"}
                </p>
              </div>

              {selected.deliverables && selected.deliverables.length > 0 ? (
                <div className="mt-4">
                  <p className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                    Deliverables
                  </p>
                  <ul className="mt-2 space-y-2">
                    {selected.deliverables.map((item, idx) => (
                      <li key={item.id || idx} className="flex items-center gap-2 text-sm text-[#E2E8F0]">
                        <CheckCircle className="size-4 text-emerald-400" />
                        <span>{item.quantity}x {item.type} {item.description ? `(${item.description})` : ""}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              ) : null}

              {/* Action Buttons based on status */}
              <div className="mt-5 flex flex-col gap-2">
                {selected.status === "pending" || selected.status === "draft" ? (
                  <Button
                    type="button"
                    disabled={actionLoading}
                    className="h-10 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
                    onClick={() => handlePayEscrow(selected.id)}
                  >
                    <ShieldCheck className="mr-2 size-4" />
                    Fund Escrow ({formatEtb(selected.fees?.total_business_payable_etb || selected.fees?.amount_etb || 0)})
                  </Button>
                ) : null}

                {selected.status === "active" ? (
                  <Button
                    type="button"
                    disabled={actionLoading}
                    className="h-10 w-full rounded-xl bg-emerald-600 text-white hover:bg-emerald-700"
                    onClick={() => handleApproveWork(selected.id)}
                  >
                    <CheckCircle className="mr-2 size-4" />
                    Approve Work & Release Payout
                  </Button>
                ) : null}

                <Button
                  type="button"
                  variant="outline"
                  className="h-10 w-full rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10"
                  onClick={handleDownloadPDF}
                >
                  <FileDown className="mr-2 size-4" />
                  Print / Save Contract (PDF)
                </Button>
              </div>
            </div>
          </aside>
        ) : null}
      </div>
    </div>
  );
}
