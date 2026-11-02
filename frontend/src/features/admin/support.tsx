"use client";

import { useCallback, useEffect, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { toast } from "sonner";

import { ListStatus, LoadMore } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import {
  adminApi,
  type AdminSupportTicketItem,
  type AdminSupportTicketDetail,
  type AdminSupportTicketStats,
} from "@/lib/api/admin";
import { usePagedList } from "@/hooks/use-paged-list";
import { adminErrorMessage } from "@/lib/api-errors";
import { cn } from "@/lib/utils";

type FilterKey = "all" | "business" | "influencer" | "open" | "in_review" | "escalated" | "resolved" | "high";

const fieldShell =
  "w-full min-w-0 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)]/60 px-3 py-2 text-sm text-white outline-none placeholder:text-[#94A3B8] focus:border-[color:var(--adfluence-orange)]/40 focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25";

function priorityClass(p: string) {
  const lower = p.toLowerCase();
  if (lower === "high") return "text-rose-200";
  if (lower === "medium") return "text-amber-200";
  return "text-[#CBD5E1]";
}

function statusLabel(s: string) {
  const lower = s.toLowerCase();
  switch (lower) {
    case "in_review":
      return "In review";
    case "resolved":
      return "Resolved";
    case "escalated":
      return "Escalated";
    case "closed":
      return "Closed";
    default:
      return "Open";
  }
}

function isDone(status: string) {
  return status === "resolved" || status === "closed";
}

function messagePreview(text: string, max = 100) {
  const t = (text || "").trim();
  if (t.length <= max) return t;
  return `${t.slice(0, max)}…`;
}

const filterButtons: { key: FilterKey; label: string }[] = [
  { key: "all", label: "All" },
  { key: "business", label: "Business" },
  { key: "influencer", label: "Influencer" },
  { key: "open", label: "Open" },
  { key: "in_review", label: "In review" },
  { key: "escalated", label: "Escalated" },
  { key: "resolved", label: "Resolved" },
  { key: "high", label: "High priority" },
];

export function AdminSupport() {
  const [filter, setFilter] = useState<FilterKey>("all");
  const [pickedId, setPickedId] = useState<string | null>(null);
  const [selectedDetail, setSelectedDetail] = useState<AdminSupportTicketDetail | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [notesDraft, setNotesDraft] = useState("");
  const [replyDraft, setReplyDraft] = useState("");

  const fetchPage = useCallback(
    async (page: number) => {
      const role = filter === "business" || filter === "influencer" ? filter : undefined;
      const status =
        filter === "open" || filter === "in_review" || filter === "escalated" || filter === "resolved" ? filter : undefined;
      const priority = filter === "high" ? "high" : undefined;
      return (await adminApi.listSupportTickets({ role, status, priority, page })).data;
    },
    [filter],
  );
  const list = usePagedList(fetchPage);
  const tickets: AdminSupportTicketItem[] = list.items;
  const loadTickets = list.reload;
  const stats: AdminSupportTicketStats = list.last?.stats ?? {
    open_count: 0,
    high_priority_count: 0,
    business_count: 0,
    influencer_count: 0,
    resolved_count: 0,
  };
  // The selection follows the list (filters change it).
  const selectedId = tickets.some((t) => t.id === pickedId) ? pickedId : (tickets[0]?.id ?? null);

  useEffect(() => {
    if (!selectedId) return;
    async function loadDetail() {
      try {
        const res = await adminApi.getSupportTicket(selectedId!);
        if (res.data) {
          setSelectedDetail(res.data);
          setNotesDraft(res.data.internal_note || "");
          setReplyDraft(res.data.admin_reply || "");
        }
      } catch (err) {
        toast.error(adminErrorMessage(err, "The ticket could not be opened."));
      }
    }
    loadDetail();
  }, [selectedId]);

  const shownDetail = selectedDetail?.id === selectedId ? selectedDetail : null;

  const [busyId, setBusyId] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  const setStatus = async (id: string, status: string) => {
    setBusyId(id);
    try {
      await adminApi.updateSupportTicketStatus(id, { status });
      toast.success(`Ticket marked ${statusLabel(status).toLowerCase()}.`);
      await loadTickets();
      if (selectedId === id) {
        const res = await adminApi.getSupportTicket(id);
        if (res.data) setSelectedDetail(res.data);
      }
    } catch (err) {
      toast.error(adminErrorMessage(err, "The ticket could not be updated."));
    } finally {
      setBusyId(null);
    }
  };

  const saveNotesForSelected = async () => {
    if (!selectedId) return;
    setSaving(true);
    try {
      await adminApi.replySupportTicket(selectedId, {
        admin_reply: replyDraft.trim() || undefined,
        internal_note: notesDraft.trim() || undefined,
      });
      toast.success(replyDraft.trim() ? "Reply and notes saved." : "Notes saved.");
      await loadTickets();
      const res = await adminApi.getSupportTicket(selectedId);
      if (res.data) setSelectedDetail(res.data);
    } catch (err) {
      toast.error(adminErrorMessage(err, "The reply could not be saved."));
    } finally {
      setSaving(false);
    }
  };

  const selectTicket = useCallback((id: string) => {
    setPickedId(id);
    if (typeof window !== "undefined" && window.matchMedia("(max-width: 1023px)").matches) {
      setMobileShowDetail(true);
    }
  }, []);

  const statCards = [
    { label: "Open tickets", value: list.state !== "ready" ? "—" : String(stats.open_count) },
    { label: "High priority", value: list.state !== "ready" ? "—" : String(stats.high_priority_count) },
    { label: "Business requests", value: list.state !== "ready" ? "—" : String(stats.business_count) },
    { label: "Influencer requests", value: list.state !== "ready" ? "—" : String(stats.influencer_count) },
    { label: "Resolved", value: list.state !== "ready" ? "—" : String(stats.resolved_count) },
  ];

  return (
    <div className="min-w-0 space-y-6 overflow-x-hidden text-white">
      <header className="space-y-1">
        <h1 className="text-2xl font-semibold md:text-3xl">Support & Comments</h1>
        <p className="max-w-3xl text-sm text-[#CBD5E1] md:text-base">
          Review support requests, feedback, and reported issues from businesses and influencers.
        </p>
      </header>

      <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-5">
        {statCards.map((s) => (
          <article
            key={s.label}
            className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.06)] p-4 shadow-[0_12px_40px_rgba(0,0,0,0.28)] backdrop-blur-md"
          >
            <p className="text-[11px] font-medium uppercase tracking-wide text-[#94A3B8]">{s.label}</p>
            <p className="mt-2 text-xl font-bold text-[color:var(--adfluence-orange)] md:text-2xl">
              {s.value}
            </p>
          </article>
        ))}
      </div>

      <div className="flex min-w-0 gap-2 overflow-x-auto pb-1 [-ms-overflow-style:none] [scrollbar-width:none] [&::-webkit-scrollbar]:hidden">
        {filterButtons.map((b) => (
          <button
            key={b.key}
            type="button"
            onClick={() => setFilter(b.key)}
            className={cn(
              "inline-flex min-h-11 shrink-0 items-center rounded-full border px-4 py-2 text-sm font-medium transition-colors sm:min-h-0 sm:px-3 sm:py-1.5 sm:text-xs",
              filter === b.key
                ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
                : "border-white/15 bg-white/5 text-[#CBD5E1] hover:bg-white/10",
            )}
          >
            {b.label}
          </button>
        ))}
      </div>

      <div className="flex min-w-0 flex-col gap-6 lg:flex-row lg:items-stretch">
        <section
          className={cn(
            "min-h-0 min-w-0 flex-1 space-y-3 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.06)] p-4 backdrop-blur-md md:p-5",
            shownDetail && mobileShowDetail && "hidden lg:block",
            shownDetail && "lg:max-w-[min(100%,520px)] xl:max-w-[560px]",
          )}
        >
          <h2 className="text-sm font-semibold uppercase tracking-wide text-[#94A3B8]">Tickets</h2>
          <ListStatus
            state={list.state}
            isEmpty={tickets.length === 0}
            emptyText="No tickets match this filter."
            errorText="Tickets could not be loaded."
            onRetry={list.retry}
          >
            <ul className="max-h-[70vh] space-y-3 overflow-y-auto pr-1 lg:max-h-[calc(100vh-220px)]">
              {tickets.map((t) => {
                const active = t.id === selectedId;
                return (
                  <li key={t.id}>
                    <article
                      className={cn(
                        "rounded-xl border p-3 transition-colors md:p-4",
                        active
                          ? "border-[color:var(--adfluence-orange)]/35 bg-[color:var(--adfluence-orange)]/10"
                          : "border-white/10 bg-[color:var(--adfluence-surface)]/40",
                      )}
                    >
                      <div className="flex flex-wrap items-start justify-between gap-2">
                        <div className="min-w-0">
                          <p className="truncate font-semibold text-white">{t.subject}</p>
                          <p className="text-xs text-[#94A3B8]">
                            {t.sender_name} · <span className="capitalize">{t.sender_role}</span> · {new Date(t.created_at).toLocaleDateString()}
                          </p>
                        </div>
                        <div className="text-right text-xs">
                          <p className={cn("font-medium", priorityClass(t.priority))}>{t.priority} priority</p>
                          <p className="mt-0.5 text-[#CBD5E1]">{statusLabel(t.status)}</p>
                        </div>
                      </div>
                      <p className="mt-2 line-clamp-2 text-xs text-[#CBD5E1]">{messagePreview(t.message_preview)}</p>
                      <p className="mt-1 text-[11px] text-[#94A3B8]">{t.issue_type_label || t.issue_type}</p>
                      <div className="mt-3 flex flex-wrap gap-2">
                        <Button
                          type="button"
                          size="sm"
                          className="h-11 flex-1 bg-white/10 text-white sm:h-8 sm:flex-none"
                          onClick={() => selectTicket(t.id)}
                        >
                          View details
                        </Button>
                        <Button
                          type="button"
                          size="sm"
                          className="h-11 flex-1 border border-amber-400/30 bg-amber-500/15 text-amber-100 sm:h-8 sm:flex-none"
                          disabled={busyId === t.id || t.status === "in_review"}
                          onClick={() => setStatus(t.id, "in_review")}
                        >
                          Mark in review
                        </Button>
                        <Button
                          type="button"
                          size="sm"
                          className="h-11 flex-1 bg-emerald-500/20 text-emerald-100 sm:h-8 sm:flex-none"
                          disabled={busyId === t.id}
                          onClick={() => setStatus(t.id, isDone(t.status) ? "open" : "resolved")}
                        >
                          {isDone(t.status) ? "Reopen" : "Mark resolved"}
                        </Button>
                        <Button
                          type="button"
                          size="sm"
                          className="h-11 flex-1 bg-rose-500/20 text-rose-100 sm:h-8 sm:flex-none"
                          disabled={busyId === t.id || t.status === "escalated" || isDone(t.status)}
                          onClick={() => setStatus(t.id, "escalated")}
                        >
                          Escalate
                        </Button>
                      </div>
                    </article>
                  </li>
                );
              })}
              <li>
                <LoadMore
                  shown={tickets.length}
                  total={list.total}
                  hasMore={list.hasMore}
                  loading={list.loadingMore}
                  onLoadMore={list.loadMore}
                />
              </li>
            </ul>
          </ListStatus>
        </section>

        <section
          className={cn(
            "min-w-0 flex-1 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.06)] p-4 backdrop-blur-md md:p-5 lg:max-w-xl xl:max-w-none",
            !shownDetail ? "hidden lg:block" : !mobileShowDetail && "hidden lg:block",
          )}
        >
          {!shownDetail ? (
            <div className="flex h-full min-h-[200px] flex-col items-center justify-center text-center text-sm text-[#94A3B8]">
              Select a ticket to view details, internal notes, and post an official reply.
            </div>
          ) : (
            <div className="space-y-4">
              <div className="lg:hidden">
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="shrink-0 text-white hover:bg-white/10"
                  aria-label="Back to tickets"
                  onClick={() => setMobileShowDetail(false)}
                >
                  <ArrowLeft className="size-5" />
                </Button>
              </div>
              <div>
                <h2 className="text-lg font-semibold text-white">{shownDetail.subject}</h2>
                <p className="mt-1 text-xs text-[#94A3B8]">
                  {shownDetail.sender_name} · {shownDetail.sender_email} · {new Date(shownDetail.created_at).toLocaleDateString()}
                </p>
              </div>
              <div className="flex flex-wrap gap-2 text-xs">
                <span className="rounded-full border border-white/15 bg-white/5 px-2 py-1 capitalize text-[#CBD5E1]">
                  {shownDetail.sender_role}
                </span>
                <span className="rounded-full border border-white/15 bg-white/5 px-2 py-1 text-[#CBD5E1]">
                  {shownDetail.issue_type_label || shownDetail.issue_type}
                </span>
                <span className={cn("rounded-full border px-2 py-1", priorityClass(shownDetail.priority))}>
                  {shownDetail.priority} priority
                </span>
                <span className="rounded-full border border-white/15 bg-white/5 px-2 py-1 text-[#CBD5E1]">
                  {statusLabel(shownDetail.status)}
                </span>
              </div>
              <div>
                <p className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">Message</p>
                <p className="mt-2 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)]/60 p-3 text-sm leading-relaxed text-[#CBD5E1]">
                  {shownDetail.message}
                </p>
              </div>
              <div>
                <p className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                  Related campaign / context
                </p>
                <p className="mt-1 text-sm text-white">{shownDetail.related_campaign_title ?? "—"}</p>
              </div>
              <div className="space-y-1.5">
                <label className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                  Admin internal notes
                </label>
                <textarea
                  className={cn(fieldShell, "min-h-[88px] resize-y")}
                  value={notesDraft}
                  onChange={(e) => setNotesDraft(e.target.value)}
                  placeholder="Internal notes visible only to admins"
                />
              </div>
              <div className="space-y-1.5">
                <label className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                  Official reply to user
                </label>
                <textarea
                  value={replyDraft}
                  onChange={(e) => setReplyDraft(e.target.value)}
                  placeholder="e.g. We’re looking into this and will update within 24h."
                  className={cn(fieldShell, "min-h-[110px] resize-y")}
                />
                <p className="text-xs text-[#94A3B8]">The user sees this reply on their support page.</p>
              </div>
              <Button
                type="button"
                className="h-11 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] lg:h-10 lg:w-auto"
                disabled={saving}
                onClick={saveNotesForSelected}
              >
                {saving ? "Saving…" : "Save note & reply"}
              </Button>
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
