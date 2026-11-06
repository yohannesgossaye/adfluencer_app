"use client";

import { useCallback, useEffect, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { toast } from "sonner";

import { ListStatus, LoadMore } from "@/components/common/list-status";
import { FilterChips } from "@/components/common/filter-chips";
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
import { formatDate } from "@/lib/contracts";
import { StatusBadge } from "@/components/common/status-badge";
import { priorityLevel, ticketStatus } from "@/lib/status";
import { ReplyBox, SupportThread } from "@/components/support/support-thread";
import { isTicketDone as isDone } from "@/lib/support";

type RoleFilter = "all" | "business" | "influencer";
type StatusFilter = "all" | "open" | "in_review" | "escalated" | "resolved";
type PriorityFilter = "all" | "high";

const fieldShell =
  "w-full min-w-0 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)]/60 px-3 py-2 text-sm text-white outline-none placeholder:text-[#94A3B8] focus:border-[color:var(--adfluence-orange)]/40 focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25";




function messagePreview(text: string, max = 100) {
  const t = (text || "").trim();
  if (t.length <= max) return t;
  return `${t.slice(0, max)}…`;
}


export function AdminSupport() {
  // Three independent filters, e.g. "Business + Open + High priority".
  const [roleFilter, setRoleFilter] = useState<RoleFilter>("all");
  const [statusFilter, setStatusFilter] = useState<StatusFilter>("all");
  const [priorityFilter, setPriorityFilter] = useState<PriorityFilter>("all");
  const [pickedId, setPickedId] = useState<string | null>(null);
  const [selectedDetail, setSelectedDetail] = useState<AdminSupportTicketDetail | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [notesDraft, setNotesDraft] = useState("");

  const fetchPage = useCallback(
    async (page: number) => {
      const role = roleFilter === "all" ? undefined : roleFilter;
      const status = statusFilter === "all" ? undefined : statusFilter;
      const priority = priorityFilter === "all" ? undefined : priorityFilter;
      return (await adminApi.listSupportTickets({ role, status, priority, page })).data;
    },
    [roleFilter, statusFilter, priorityFilter],
  );
  const list = usePagedList(["admin", "support", { roleFilter, statusFilter, priorityFilter }], fetchPage);
  const tickets: AdminSupportTicketItem[] = list.items;
  const loadTickets = list.reload;
  const stats: AdminSupportTicketStats = list.last?.stats ?? {
    open_count: 0,
    high_priority_count: 0,
    business_count: 0,
    influencer_count: 0,
    resolved_count: 0,
    needs_reply_count: 0,
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
      toast.success(`Ticket marked ${ticketStatus(status).label.toLowerCase()}.`);
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

  const refreshSelected = async (id: string) => {
    const res = await adminApi.getSupportTicket(id);
    if (res.data) setSelectedDetail(res.data);
  };

  // A public reply: added to the conversation and emailed to the user.
  const sendReply = async (id: string, text: string) => {
    try {
      await adminApi.replySupportTicket(id, { admin_reply: text });
    } catch (err) {
      toast.error(adminErrorMessage(err, "The reply could not be sent."));
      throw err;
    }
    toast.success("Reply sent. The user has been emailed.");
    await Promise.all([loadTickets(), refreshSelected(id).catch(() => {})]);
  };

  const saveNote = async () => {
    if (!selectedId || !notesDraft.trim()) return;
    setSaving(true);
    try {
      await adminApi.replySupportTicket(selectedId, { internal_note: notesDraft.trim() });
      toast.success("Internal note saved.");
    } catch (err) {
      toast.error(adminErrorMessage(err, "The note could not be saved."));
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
    { label: "Waiting for reply", value: list.state !== "ready" ? "—" : String(stats.needs_reply_count) },
    { label: "Open tickets", value: list.state !== "ready" ? "—" : String(stats.open_count) },
    { label: "High priority", value: list.state !== "ready" ? "—" : String(stats.high_priority_count) },
    { label: "Business requests", value: list.state !== "ready" ? "—" : String(stats.business_count) },
    { label: "Creator requests", value: list.state !== "ready" ? "—" : String(stats.influencer_count) },
    { label: "Resolved", value: list.state !== "ready" ? "—" : String(stats.resolved_count) },
  ];

  return (
    <div className="min-w-0 space-y-6 overflow-x-hidden text-white">
      <header className="space-y-1">
        <h1 className="text-2xl font-semibold md:text-3xl">Support & Comments</h1>
        <p className="max-w-3xl text-sm text-[#CBD5E1] md:text-base">
          Support requests, feedback and reported issues from businesses and creators.
        </p>
      </header>

      <div className="grid gap-3 sm:grid-cols-3 xl:grid-cols-6">
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

      <div className="space-y-3">
        <FilterChips
          label="From"
          showLabel
          value={roleFilter}
          onChange={setRoleFilter}
          options={[
            { value: "all", label: "Everyone" },
            { value: "business", label: "Businesses" },
            { value: "influencer", label: "Creators" },
          ]}
        />
        <FilterChips
          label="Status"
          showLabel
          value={statusFilter}
          onChange={setStatusFilter}
          options={[
            { value: "all", label: "Any status" },
            { value: "open", label: "Open" },
            { value: "in_review", label: "In review" },
            { value: "escalated", label: "Escalated" },
            { value: "resolved", label: "Resolved" },
          ]}
        />
        <FilterChips
          label="Priority"
          showLabel
          value={priorityFilter}
          onChange={setPriorityFilter}
          options={[
            { value: "all", label: "Any priority" },
            { value: "high", label: "High priority" },
          ]}
        />
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
            <ul className="max-h-[70vh] space-y-3 overflow-y-auto pr-1 lg:max-h-[calc(100vh_-_220px)]">
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
                            {t.sender_name} · {t.sender_role === "influencer" ? "Creator" : "Business"} · {formatDate(t.last_activity_at || t.created_at)}
                          </p>
                        </div>
                        <div className="flex flex-col items-end gap-1">
                          {!isDone(t.status) && t.last_reply_by === "user" ? (
                            <StatusBadge status={{ label: "User replied", tone: "warning" }} />
                          ) : !isDone(t.status) && t.last_reply_by === "" ? (
                            <StatusBadge status={{ label: "New", tone: "info" }} />
                          ) : null}
                          <StatusBadge status={ticketStatus(t.status)} />
                          <StatusBadge status={{ ...priorityLevel(t.priority), label: `${priorityLevel(t.priority).label} priority` }} />
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
                          className="h-11 flex-1 bg-white/10 text-white hover:bg-white/15 sm:h-8 sm:flex-none"
                          disabled={busyId === t.id || t.status === "in_review"}
                          onClick={() => setStatus(t.id, "in_review")}
                        >
                          Mark in review
                        </Button>
                        <Button
                          type="button"
                          size="sm"
                          className="h-11 flex-1 bg-white/10 text-white hover:bg-white/15 sm:h-8 sm:flex-none"
                          disabled={busyId === t.id}
                          onClick={() => setStatus(t.id, isDone(t.status) ? "open" : "resolved")}
                        >
                          {isDone(t.status) ? "Reopen" : "Mark resolved"}
                        </Button>
                        <Button
                          type="button"
                          size="sm"
                          className="h-11 flex-1 border border-rose-400/50 bg-transparent text-rose-200 hover:bg-rose-500/10 sm:h-8 sm:flex-none"
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
              Select a ticket to see the conversation and reply.
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
                  {shownDetail.sender_name} · {shownDetail.sender_email} · {formatDate(shownDetail.created_at)}
                </p>
              </div>
              <div className="flex flex-wrap gap-2 text-xs">
                <span className="rounded-full border border-white/15 bg-white/5 px-2 py-0.5 text-[#CBD5E1]">
                  {shownDetail.sender_role === "influencer" ? "Creator" : "Business"}
                </span>
                <span className="rounded-full border border-white/15 bg-white/5 px-2 py-0.5 text-[#CBD5E1]">
                  {shownDetail.issue_type_label || shownDetail.issue_type}
                </span>
                <StatusBadge status={{ ...priorityLevel(shownDetail.priority), label: `${priorityLevel(shownDetail.priority).label} priority` }} />
                <StatusBadge status={ticketStatus(shownDetail.status)} />
              </div>
              {shownDetail.related_campaign_title ? (
                <p className="text-xs text-[#94A3B8]">
                  About: <span className="text-white">{shownDetail.related_campaign_title}</span>
                </p>
              ) : null}
              <SupportThread
                viewer="support"
                opening={{ name: shownDetail.sender_name, body: shownDetail.message, created_at: shownDetail.created_at }}
                messages={shownDetail.messages ?? []}
              />
              <ReplyBox
                key={shownDetail.id}
                label="Reply to the user"
                placeholder="e.g. We've confirmed your payment. The creator can start now."
                hint={`Added to the conversation and emailed to ${shownDetail.sender_email}.`}
                onSend={(text) => sendReply(shownDetail.id, text)}
              />
              <div className="space-y-1.5 border-t border-white/10 pt-4">
                <label htmlFor="support-internal-note" className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                  Internal note (only admins see this)
                </label>
                <textarea
                  id="support-internal-note"
                  className={cn(fieldShell, "min-h-[72px] resize-y")}
                  value={notesDraft}
                  onChange={(e) => setNotesDraft(e.target.value)}
                  placeholder="e.g. Waiting for the bank statement."
                />
                <Button
                  type="button"
                  variant="outline"
                  className="h-11 w-full rounded-xl border-white/15 text-white lg:h-9 lg:w-auto"
                  disabled={saving || !notesDraft.trim()}
                  onClick={saveNote}
                >
                  {saving ? "Saving…" : "Save note"}
                </Button>
              </div>
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
