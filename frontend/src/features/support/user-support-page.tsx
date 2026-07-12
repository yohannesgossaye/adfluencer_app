"use client";

import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  appendSupportTicket,
  filterTicketsForPortalUser,
  loadSupportTickets,
  type SupportPriority,
  type SupportRole,
  type SupportTicket,
} from "@/data/support";
import { getMockUser } from "@/lib/mock-auth";
import { cn } from "@/lib/utils";

const fieldShell =
  "w-full min-w-0 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)]/60 px-3 py-2 text-sm text-white outline-none transition-[border,box-shadow] placeholder:text-[#94A3B8] focus:border-[color:var(--adfluence-orange)]/40 focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25";

function priorityBadgeClass(p: SupportPriority) {
  if (p === "high") return "border-rose-400/35 bg-rose-500/15 text-rose-100";
  if (p === "medium") return "border-amber-400/35 bg-amber-500/15 text-amber-100";
  return "border-white/15 bg-white/10 text-[#CBD5E1]";
}

function statusLabel(status: SupportTicket["status"]) {
  switch (status) {
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

type UserSupportPageProps = {
  role: SupportRole;
  subtitle: string;
  issueTypeOptions: readonly string[];
  relatedOptions: readonly string[];
};

export function UserSupportPage({
  role,
  subtitle,
  issueTypeOptions,
  relatedOptions,
}: UserSupportPageProps) {
  const [issueType, setIssueType] = useState<string>(issueTypeOptions[0] ?? "");
  const [subject, setSubject] = useState("");
  const [message, setMessage] = useState("");
  const [priority, setPriority] = useState<SupportPriority>("medium");
  const [related, setRelated] = useState<string>(relatedOptions[0] ?? "");
  const [success, setSuccess] = useState(false);
  const [myTickets, setMyTickets] = useState<SupportTicket[]>([]);

  const refreshMine = () => {
    const user = getMockUser();
    const all = loadSupportTickets();
    setMyTickets(filterTicketsForPortalUser(all, role, user?.email ?? null));
  };

  useEffect(() => {
    refreshMine();
  }, [role]);

  const sessionUser = getMockUser();

  const submit = () => {
    const trimmedSubject = subject.trim();
    const trimmedMessage = message.trim();
    if (!trimmedSubject || !trimmedMessage) return;
    const session = getMockUser();
    const email = session?.email ?? "";
    const senderName = session?.name ?? "Guest";
    const row: SupportTicket = {
      id: `usr-${Date.now()}`,
      senderName,
      role,
      email,
      subject: trimmedSubject,
      message: trimmedMessage,
      type: issueType,
      priority,
      status: "open",
      date: new Date().toISOString().slice(0, 10),
      relatedCampaign: related,
      source: "user",
    };
    appendSupportTicket(row);
    setSuccess(true);
    setSubject("");
    setMessage("");
    setPriority("medium");
    setIssueType(issueTypeOptions[0] ?? "");
    setRelated(relatedOptions[0] ?? "");
    refreshMine();
    window.setTimeout(() => setSuccess(false), 5000);
  };

  return (
    <div className="min-w-0 space-y-6 overflow-x-hidden text-white">
      <header className="space-y-1">
        <h1 className="text-2xl font-semibold text-white md:text-3xl">Support & Feedback</h1>
        <p className="max-w-2xl text-sm text-[#CBD5E1] md:text-base">{subtitle}</p>
      </header>

      <div className="grid min-w-0 gap-6 lg:grid-cols-2 lg:items-start">
        <section
          className={cn(
            "min-w-0 space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.06)] p-4 shadow-[0_18px_60px_rgba(0,0,0,0.35)] backdrop-blur-md md:p-6",
          )}
        >
          <h2 className="text-lg font-semibold text-white">Submit a request</h2>
          <div className="space-y-4">
            <div className="space-y-1.5">
              <label className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                Issue type
              </label>
              <select
                className={cn(fieldShell, "h-10 appearance-none bg-[#06111F]/80")}
                value={issueType}
                onChange={(e) => setIssueType(e.target.value)}
              >
                {issueTypeOptions.map((opt) => (
                  <option key={opt} value={opt} className="bg-[#06111F] text-white">
                    {opt}
                  </option>
                ))}
              </select>
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                Subject
              </label>
              <Input
                value={subject}
                onChange={(e) => setSubject(e.target.value)}
                placeholder="Short summary"
                className={cn(fieldShell, "h-10 border-white/10 bg-[#06111F]/60 text-white")}
              />
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                Message
              </label>
              <textarea
                className={cn(fieldShell, "min-h-[120px] resize-y")}
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder="Describe your question or issue"
              />
            </div>
            <div className="space-y-1.5">
              <span className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                Priority
              </span>
              <div className="flex flex-wrap gap-2">
                {(["low", "medium", "high"] as const).map((p) => (
                  <button
                    key={p}
                    type="button"
                    onClick={() => setPriority(p)}
                    className={cn(
                      "rounded-full border px-3 py-1.5 text-xs font-medium capitalize transition-colors",
                      priority === p
                        ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
                        : "border-white/15 bg-white/5 text-[#CBD5E1] hover:bg-white/10",
                    )}
                  >
                    {p}
                  </button>
                ))}
              </div>
            </div>
            <div className="space-y-1.5">
              <label className="text-xs font-medium uppercase tracking-wide text-[#94A3B8]">
                {role === "business" ? "Related campaign" : "Related campaign / invite"}
              </label>
              <select
                className={cn(fieldShell, "h-10 appearance-none bg-[#06111F]/80")}
                value={related}
                onChange={(e) => setRelated(e.target.value)}
              >
                {relatedOptions.map((opt) => (
                  <option key={opt} value={opt} className="bg-[#06111F] text-white">
                    {opt}
                  </option>
                ))}
              </select>
            </div>
            {success ? (
              <p className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-100">
                Support request submitted. Admin will review it.
              </p>
            ) : null}
            <Button
              type="button"
              className="h-11 w-full rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] md:h-10"
              onClick={submit}
            >
              Submit Support Request
            </Button>
            {!sessionUser ? (
              <p className="text-xs text-[#94A3B8]">
                Sign in with a demo account to attach your email to new tickets and see them under My
                Requests.
              </p>
            ) : null}
          </div>
        </section>

        <section
          className={cn(
            "min-w-0 space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.06)] p-4 shadow-[0_18px_60px_rgba(0,0,0,0.35)] backdrop-blur-md md:p-6",
          )}
        >
          <h2 className="text-lg font-semibold text-white">My Requests</h2>
          {myTickets.length === 0 ? (
            <p className="text-sm text-[#94A3B8]">
              No requests yet for this account. Submit a ticket to see it here.
            </p>
          ) : (
            <ul className="space-y-3">
              {myTickets
                .slice()
                .sort((a, b) => (a.date < b.date ? 1 : -1))
                .map((t) => (
                  <li
                    key={t.id}
                    className="rounded-xl border border-white/10 bg-[#06111F]/50 p-3 md:p-4"
                  >
                    <div className="flex flex-wrap items-start justify-between gap-2">
                      <p className="font-medium text-white">{t.subject}</p>
                      <span
                        className={cn(
                          "shrink-0 rounded-full border px-2 py-0.5 text-[10px] font-semibold uppercase tracking-wide",
                          priorityBadgeClass(t.priority),
                        )}
                      >
                        {t.priority}
                      </span>
                    </div>
                    <p className="mt-1 text-xs text-[#94A3B8]">
                      {t.type} · {statusLabel(t.status)} · {t.date}
                    </p>
                    <p className="mt-2 text-xs text-[#CBD5E1]">
                      <span className="text-[#94A3B8]">Admin reply: </span>
                      {t.adminReply?.trim() ? t.adminReply : "—"}
                    </p>
                  </li>
                ))}
            </ul>
          )}
        </section>
      </div>
    </div>
  );
}
