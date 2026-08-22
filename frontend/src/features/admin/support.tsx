"use client";

import { useCallback, useEffect, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  adminApi,
  type AdminSupportTicketItem,
  type AdminSupportTicketDetail,
  type AdminSupportTicketStats,
} from "@/lib/api/admin";
import { cn } from "@/lib/utils";

type FilterKey = "all" | "business" | "influencer" | "open" | "in_review" | "resolved" | "high";

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
    default:
      return "Open";
  }
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
  { key: "resolved", label: "Resolved" },
  { key: "high", label: "High priority" },
];

export function AdminSupport() {
  const [tickets, setTickets] = useState<AdminSupportTicketItem[]>([]);
  const [stats, setStats] = useState<AdminSupportTicketStats>({
    open_count: 0,
    high_priority_count: 0,
    business_count: 0,
    influencer_count: 0,
    resolved_count: 0,
  });
  const [filter, setFilter] = useState<FilterKey>("all");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [selectedDetail, setSelectedDetail] = useState<AdminSupportTicketDetail | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [notesDraft, setNotesDraft] = useState("");
  const [replyDraft, setReplyDraft] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  const loadTickets = async () => {
    try {
      const roleParam = filter === "business" ? "business" : filter === "influencer" ? "influencer" : undefined;
      const statusParam = filter === "open" ? "open" : filter === "in_review" ? "in_review" : filter === "resolved" ? "resolved" : undefined;
      const priorityParam = filter === "high" ? "high" : undefined;

      const res = await adminApi.listSupportTickets({
        role: roleParam,
        status: statusParam,
        priority: priorityParam,
      });

      if (res.data?.items) {
        setTickets(res.data.items);
        if (!selectedId && res.data.items[0]) {
          setSelectedId(res.data.items[0].id);
        }
      }
      if (res.data?.stats) {
        setStats(res.data.stats);
      }
    } catch {
      // fallback
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadTickets();
  }, [filter]);

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
      } catch {
        // fallback
      }
    }
    loadDetail();
  }, [selectedId]);

  const setStatus = async (id: string, status: string) => {
    try {
      await adminApi.updateSupportTicketStatus(id, { status });
      await loadTickets();
      if (selectedId === id) {
        const res = await adminApi.getSupportTicket(id);
        if (res.data) setSelectedDetail(res.data);
      }
    } catch {
      // fallback
    }
  };

  const saveNotesForSelected = async () => {
    if (!selectedId) return;
    try {
      await adminApi.replySupportTicket(selectedId, {
        admin_reply: replyDraft.trim() || undefined,
        internal_note: notesDraft.trim() || undefined,
      });
      await loadTickets();
      const res = await adminApi.getSupportTicket(selectedId);
      if (res.data) setSelectedDetail(res.data);
    } catch {
      // fallback
    }
  };

  const selectTicket = useCallback((id: string) => {
    setSelectedId(id);
    if (typeof window !== "undefined" && window.matchMedia("(max-width: 1023px)").matches) {
      setMobileShowDetail(true);
    }
  }, []);

  const statCards = [
    { label: "Open tickets", value: String(stats.open_count) },
    { label: "High priority", value: String(stats.high_priority_count) },
    { label: "Business requests", value: String(stats.business_count) },
    { label: "Influencer requests", value: String(stats.influencer_count) },
    { label: "Resolved", value: String(stats.resolved_count) },
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
            selectedDetail && mobileShowDetail && "hidden lg:block",
            selectedDetail && "lg:max-w-[min(100%,520px)] xl:max-w-[560px]",
          )}
        >
          <h2 className="text-sm font-semibold uppercase tracking-wide text-[#94A3B8]">Tickets</h2>
          {tickets.length === 0 && !isLoading ? (
            <p className="text-sm text-[#94A3B8]">No tickets match this filter.</p>
          ) : (
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
                          disabled={t.status === "in_review"}
                          onClick={() => setStatus(t.id, "in_review")}
                        >
                          Mark in review
                        </Button>
                        <Button
                          type="button"
                          size="sm"
                          className="h-11 flex-1 bg-emerald-500/20 text-emerald-100 sm:h-8 sm:flex-none"
                          disabled={t.status === "resolved"}
                          onClick={() => setStatus(t.id, "resolved")}
                        >
                          Mark resolved
                        </Button>
                        <Button
                          type="button"
                          size="sm"
                          className="h-11 flex-1 bg-rose-500/20 text-rose-100 sm:h-8 sm:flex-none"
                          disabled={t.status === "escalated"}
                          onClick={() => setStatus(t.id, "escalated")}
                        >
                          Escalate
                        </Button>
                      </div>
                    </article>
                  </li>
                );
              })}
            </ul>
          )}
        </section>

        <section
          className={cn(
            "min-w-0 flex-1 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.06)] p-4 backdrop-blur-md md:p-5 lg:max-w-xl xl:max-w-none",
            !selectedDetail ? "hidden lg:block" : !mobileShowDetail && "hidden lg:block",
          )}
        >
          {!selectedDetail ? (
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
                <h2 className="text-lg font-semibold text-white">{selectedDetail.subject}</h2>
                <p className="mt-1 text-xs text-[#94A3B8]">
                  {selectedDetail.sender_name} · {selectedDetail.sender_email} · {new Date(selectedDetail.created_at).toLocaleDateString()}
                </p>
              </div>
              <div className="flex flex-wrap gap-2 text-xs">
                <span className="rounded-full border border-white/15 bg-white/5 px-2 py-1 capitalize text-[#CBD5E1]">
                  {selectedDetail.sender_role}
                </span>
                <span className="rounded-full border border-white/15 bg-white/5 px-2 py-1 text-[#CBD5E1]">
                  {selectedDetail.issue_type_label || selectedDetail.issue_type}
                </span>
                <span className={cn("rounded-full border px-2 py-1", priorityClass(selectedDetail.priority))}>
                  {selectedDetail.priority} priority
                </span>
                <span className="rounded-full border border-white/15 bg-white/5 px-2 py-1 text-[#CBD5E1]">
                  {statusLabel(selectedDetail.status)}
                </span>
              </div>
              <div>
                <p className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">Message</p>
                <p className="mt-2 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)]/60 p-3 text-sm leading-relaxed text-[#CBD5E1]">
                  {selectedDetail.message}
                </p>
              </div>
              <div>
                <p className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                  Related campaign / context
                </p>
                <p className="mt-1 text-sm text-white">{selectedDetail.related_campaign_title ?? "—"}</p>
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
                <Input
                  value={replyDraft}
                  onChange={(e) => setReplyDraft(e.target.value)}
                  placeholder="e.g. We’re looking into this and will update within 24h."
                  className={cn(fieldShell, "h-10 border-white/10 bg-[color:var(--adfluence-surface)]/60")}
                />
              </div>
              <Button
                type="button"
                className="h-11 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] lg:h-10 lg:w-auto"
                onClick={saveNotesForSelected}
              >
                Save note & reply
              </Button>
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
