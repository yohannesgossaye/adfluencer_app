"use client";

import { useCallback, useEffect, useState } from "react";
import axios from "axios";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { businessApi, type SupportTicket } from "@/lib/api/business";
import { influencerApi } from "@/lib/api/influencer";
import { cn } from "@/lib/utils";

export type SupportRole = "business" | "influencer";

/** Issue types accepted by the backend for each role (support/types.go). */
const ISSUE_TYPES: Record<SupportRole, Array<{ value: string; label: string }>> = {
  business: [
    { value: "campaign_issue", label: "Campaign issue" },
    { value: "contract_issue", label: "Contract issue" },
    { value: "payment_question", label: "Payment question" },
    { value: "creator_communication", label: "Creator communication" },
    { value: "technical_issue", label: "Technical issue" },
    { value: "general_feedback", label: "General feedback" },
  ],
  influencer: [
    { value: "profile_approval", label: "Profile approval" },
    { value: "campaign_issue", label: "Campaign invite" },
    { value: "contract_issue", label: "Contract issue" },
    { value: "payment_question", label: "Payment or payout question" },
    { value: "complaint", label: "Problem with a business" },
    { value: "technical_issue", label: "Technical issue" },
    { value: "general_feedback", label: "General feedback" },
  ],
};

const PRIORITIES = [
  { value: "low", label: "Low" },
  { value: "normal", label: "Normal" },
  { value: "high", label: "High" },
  { value: "urgent", label: "Urgent" },
];

const STATUS_LABEL: Record<string, string> = {
  open: "Open",
  in_review: "In review",
  resolved: "Resolved",
  escalated: "Escalated",
  closed: "Closed",
};

const field =
  "w-full min-w-0 rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-white outline-none placeholder:text-[#94A3B8] focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25";

function apiError(err: unknown, fallback: string) {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as { message?: string } | undefined;
    return data?.message || (!err.response ? "You appear to be offline. Your ticket was not sent." : fallback);
  }
  return fallback;
}

export function UserSupportPage({ role, subtitle }: { role: SupportRole; subtitle: string }) {
  const issueTypes = ISSUE_TYPES[role];
  const [issueType, setIssueType] = useState(issueTypes[0].value);
  const [subject, setSubject] = useState("");
  const [message, setMessage] = useState("");
  const [priority, setPriority] = useState("normal");
  const [related, setRelated] = useState("");
  const [relatedOptions, setRelatedOptions] = useState<Array<{ id: string; label: string }>>([]);
  const [tickets, setTickets] = useState<SupportTicket[]>([]);
  const [ticketsState, setTicketsState] = useState<"loading" | "ready" | "error">("loading");
  const [submitting, setSubmitting] = useState(false);

  const loadTickets = useCallback(async () => {
    try {
      const res = role === "business" ? await businessApi.listSupportTickets() : await influencerApi.listSupportTickets();
      setTickets(res.data?.items ?? []);
      setTicketsState("ready");
    } catch {
      setTicketsState("error");
    }
  }, [role]);

  useEffect(() => {
    void Promise.resolve().then(loadTickets);
    const options =
      role === "business"
        ? businessApi.listSupportCampaignOptions().then((r) => (r.data?.items ?? []).map((c) => ({ id: c.id, label: c.title })))
        : influencerApi
            .listSupportInviteOptions()
            .then((r) => (r.data?.items ?? []).map((i) => ({ id: i.id, label: `${i.campaign_title} · ${i.business_name}` })));
    options.then(setRelatedOptions).catch(() => setRelatedOptions([]));
  }, [role, loadTickets]);

  const canSubmit = subject.trim().length > 0 && message.trim().length > 0 && !submitting;

  const submit = async () => {
    if (!canSubmit) return;
    setSubmitting(true);
    try {
      const base = { issue_type: issueType, subject: subject.trim(), message: message.trim(), priority };
      if (role === "business") {
        await businessApi.createSupportTicket({ ...base, related_campaign_id: related || undefined });
      } else {
        await influencerApi.createSupportTicket({ ...base, related_invite_id: related || undefined });
      }
      toast.success("Your request was sent. Our team will reply here.");
      setSubject("");
      setMessage("");
      setPriority("normal");
      setRelated("");
      await loadTickets();
    } catch (err) {
      toast.error(apiError(err, "Your request could not be sent. Please try again."));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="min-w-0 space-y-6 text-white">
      <header className="space-y-1">
        <h1 className="text-2xl font-semibold md:text-3xl">Support & feedback</h1>
        <p className="max-w-2xl text-sm text-[#CBD5E1] md:text-base">{subtitle}</p>
      </header>

      <div className="grid gap-6 xl:grid-cols-[minmax(0,1fr)_minmax(0,1fr)]">
        <form
          className="space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5"
          onSubmit={(e) => {
            e.preventDefault();
            void submit();
          }}
        >
          <h2 className="text-lg font-semibold">New request</h2>
          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">What do you need help with?</span>
            <select id="support-issue" value={issueType} onChange={(e) => setIssueType(e.target.value)} className={cn(field, "h-11")}>
              {issueTypes.map((t) => (
                <option key={t.value} value={t.value}>{t.label}</option>
              ))}
            </select>
          </label>
          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">{role === "business" ? "Related campaign" : "Related invite"}</span>
            <select id="support-related" value={related} onChange={(e) => setRelated(e.target.value)} className={cn(field, "h-11")}>
              <option value="">{role === "business" ? "Not about a specific campaign" : "Not about a specific invite"}</option>
              {relatedOptions.map((o) => (
                <option key={o.id} value={o.id}>{o.label}</option>
              ))}
            </select>
          </label>
          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">Subject *</span>
            <Input id="support-subject" value={subject} maxLength={200} onChange={(e) => setSubject(e.target.value)} className={cn(field, "h-11")} placeholder="Short summary" />
          </label>
          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">Message *</span>
            <textarea id="support-message" value={message} rows={5} maxLength={5000} onChange={(e) => setMessage(e.target.value)} className={field} placeholder="Tell us what happened and what you expected." />
          </label>
          <fieldset className="space-y-1">
            <legend className="text-sm text-[#CBD5E1]">Priority</legend>
            <div className="flex flex-wrap gap-2">
              {PRIORITIES.map((p) => (
                <button
                  key={p.value}
                  type="button"
                  aria-pressed={priority === p.value}
                  onClick={() => setPriority(p.value)}
                  className={cn(
                    "min-h-11 rounded-full border px-4 text-sm md:min-h-9",
                    priority === p.value ? "border-transparent bg-[color:var(--adfluence-orange)] text-white" : "border-white/10 text-[#CBD5E1] hover:text-white",
                  )}
                >
                  {p.label}
                </button>
              ))}
            </div>
          </fieldset>
          <Button type="submit" disabled={!canSubmit} className="h-11 rounded-xl bg-[color:var(--adfluence-orange)] px-6 text-white hover:bg-[color:var(--adfluence-orange-hover)]">
            {submitting ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
            Send request
          </Button>
        </form>

        <section className="space-y-3 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5">
          <h2 className="text-lg font-semibold">My requests</h2>
          {ticketsState === "loading" ? (
            <p className="flex items-center gap-2 text-sm text-[#94A3B8]"><Loader2 className="size-4 animate-spin" aria-hidden /> Loading…</p>
          ) : ticketsState === "error" ? (
            <p role="alert" className="text-sm text-rose-200">
              Could not load your requests.{" "}
              <button type="button" className="underline" onClick={() => void loadTickets()}>Try again</button>
            </p>
          ) : tickets.length === 0 ? (
            <p className="text-sm text-[#CBD5E1]">You haven’t contacted support yet.</p>
          ) : (
            <ul className="space-y-2">
              {tickets.map((t) => (
                <li key={t.id} className="rounded-xl border border-white/10 bg-[rgba(255,255,255,0.04)] p-3">
                  <div className="flex flex-wrap items-center justify-between gap-2">
                    <p className="font-medium">{t.subject}</p>
                    <span className="rounded-full border border-white/15 px-2 py-0.5 text-[10px] text-[#CBD5E1]">{STATUS_LABEL[t.status] ?? t.status}</span>
                  </div>
                  <p className="mt-0.5 text-xs text-[#94A3B8]">
                    {t.issue_type_label || t.issue_type} · {new Date(t.created_at).toLocaleDateString()}
                    {t.related_campaign_title ? ` · ${t.related_campaign_title}` : ""}
                  </p>
                  {t.admin_reply ? (
                    <p className="mt-2 rounded-lg border border-emerald-400/20 bg-emerald-500/10 px-3 py-2 text-sm text-emerald-50">
                      <span className="font-medium">Adfluence: </span>
                      {t.admin_reply}
                    </p>
                  ) : null}
                </li>
              ))}
            </ul>
          )}
        </section>
      </div>
    </div>
  );
}
