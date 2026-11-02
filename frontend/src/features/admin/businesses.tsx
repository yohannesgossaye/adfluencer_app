"use client";

import { useCallback, useEffect, useRef, useState, type ReactNode } from "react";
import axios from "axios";
import { ArrowLeft, BadgeCheck, ExternalLink, Loader2, ShieldX } from "lucide-react";
import { toast } from "sonner";

import { UserAvatar } from "@/components/common/user-avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  adminApi,
  type BusinessVerificationDetail,
  type BusinessVerificationFilter,
  type BusinessVerificationItem,
} from "@/lib/api/admin";
import { cn } from "@/lib/utils";

const tabs: Array<{ id: BusinessVerificationFilter; label: string }> = [
  { id: "pending", label: "Awaiting review" },
  { id: "verified", label: "Verified" },
  { id: "rejected", label: "Not approved" },
  { id: "all", label: "All" },
];

const statusLabel: Record<string, string> = {
  pending: "Awaiting review",
  verified: "Verified",
  rejected: "Not approved",
  unsubmitted: "Not submitted",
};

function badge(status: string) {
  if (status === "pending") return "border-amber-400/30 bg-amber-500/10 text-amber-200";
  if (status === "verified") return "border-emerald-400/30 bg-emerald-500/10 text-emerald-200";
  return "border-rose-400/30 bg-rose-500/10 text-rose-200";
}

function since(iso?: string) {
  if (!iso) return "";
  const days = Math.floor((Date.now() - new Date(iso).getTime()) / 86_400_000);
  if (days <= 0) return "today";
  if (days === 1) return "1 day ago";
  return `${days} days ago`;
}

function formatDate(iso?: string) {
  if (!iso) return "—";
  try {
    return new Intl.DateTimeFormat("en-GB", { day: "numeric", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" }).format(new Date(iso));
  } catch {
    return iso;
  }
}

function apiError(err: unknown, fallback: string) {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as { message?: string; error?: { details?: Array<{ description?: string }> } } | undefined;
    return data?.error?.details?.[0]?.description || data?.message || (!err.response ? "You appear to be offline." : fallback);
  }
  return fallback;
}

function Row({ label, children, highlight }: { label: string; children: ReactNode; highlight?: boolean }) {
  return (
    <div className="grid grid-cols-[130px_minmax(0,1fr)] gap-2 py-1.5 text-sm">
      <dt className="text-[#94A3B8]">{label}</dt>
      <dd className={cn("min-w-0 break-words text-white", highlight && "font-mono")}>{children || <span className="text-[#64748B]">Not provided</span>}</dd>
    </div>
  );
}

function ExtLink({ href }: { href?: string }) {
  if (!href) return null;
  return (
    <a href={href} target="_blank" rel="noopener noreferrer" className="inline-flex items-center gap-1 text-[color:var(--adfluence-orange)] hover:underline">
      {href.replace(/^https?:\/\//, "")} <ExternalLink className="size-3" aria-hidden />
    </a>
  );
}

type Decision = "verify" | "reject" | "revoke";

export function AdminBusinesses() {
  const [tab, setTab] = useState<BusinessVerificationFilter>("pending");
  const [search, setSearch] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [rows, setRows] = useState<BusinessVerificationItem[]>([]);
  const [total, setTotal] = useState(0);
  const [listState, setListState] = useState<"loading" | "ready" | "error">("loading");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [detail, setDetail] = useState<BusinessVerificationDetail | null>(null);
  const [detailError, setDetailError] = useState<string | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [decision, setDecision] = useState<Decision | null>(null);
  const [text, setText] = useState("");
  const [busy, setBusy] = useState(false);
  const listRequest = useRef(0);

  useEffect(() => {
    const t = setTimeout(() => setDebouncedSearch(search.trim()), 300);
    return () => clearTimeout(t);
  }, [search]);

  const loadList = useCallback(async () => {
    const id = ++listRequest.current;
    try {
      const res = await adminApi.listBusinessVerifications({ status: tab, search: debouncedSearch || undefined, limit: 50 });
      if (id !== listRequest.current) return; // a newer search replaced this one
      setRows(res.data?.items ?? []);
      setTotal(res.data?.pagination?.total_docs ?? res.data?.items.length ?? 0);
      setListState("ready");
    } catch {
      if (id === listRequest.current) setListState("error");
    }
  }, [tab, debouncedSearch]);

  useEffect(() => {
    void Promise.resolve().then(loadList);
  }, [loadList]);

  const loadDetail = useCallback(async (userId: string) => {
    try {
      const res = await adminApi.getBusinessVerification(userId);
      setDetail(res.data ?? null);
      setDetailError(null);
    } catch (err) {
      setDetailError(apiError(err, "Could not load this business."));
    }
  }, []);

  const select = (userId: string) => {
    setSelectedId(userId);
    setDetail(null);
    setDecision(null);
    setText("");
    void loadDetail(userId);
    if (window.matchMedia("(max-width: 1279px)").matches) setMobileShowDetail(true);
  };

  const startDecision = (d: Decision) => {
    setDecision(d);
    setText("");
  };

  const reasonTooShort = decision !== "verify" && text.trim().length < 5;

  const confirm = async () => {
    if (!detail || !decision || reasonTooShort) return;
    setBusy(true);
    try {
      if (decision === "verify") {
        await adminApi.verifyBusiness(detail.user_id, { admin_note: text.trim() || undefined });
        toast.success(`${detail.business_name} is now verified.`);
      } else if (decision === "reject") {
        await adminApi.rejectBusinessVerification(detail.user_id, { reason: text.trim() });
        toast.success(`${detail.business_name} was not approved. They will see your reason.`);
      } else {
        await adminApi.revokeBusinessVerification(detail.user_id, { reason: text.trim() });
        toast.success(`Verification revoked for ${detail.business_name}.`);
      }
      setDecision(null);
      setText("");
      await Promise.all([loadList(), loadDetail(detail.user_id)]);
    } catch (err) {
      toast.error(apiError(err, "The action could not be completed."));
    } finally {
      setBusy(false);
    }
  };

  const decisionCopy: Record<Decision, { title: string; label: string; placeholder: string; button: string }> = {
    verify: {
      title: "Verify this business",
      label: "Note (optional, visible to the business)",
      placeholder: "e.g. TIN checked against the trade license",
      button: "Verify business",
    },
    reject: {
      title: "Do not approve",
      label: "Reason (required, shown to the business)",
      placeholder: "Tell them exactly what to fix, e.g. the TIN does not match the legal name",
      button: "Send rejection",
    },
    revoke: {
      title: "Revoke verification",
      label: "Reason (required, shown to the business)",
      placeholder: "e.g. Reported for misleading campaigns",
      button: "Revoke verification",
    },
  };

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold">Business verification</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Businesses can only publish campaigns, invite creators and create contracts once verified. Oldest submissions are listed first.
        </p>
      </div>

      <div className="flex flex-wrap gap-2" role="tablist" aria-label="Verification status">
        {tabs.map((t) => (
          <button
            key={t.id}
            type="button"
            role="tab"
            aria-selected={tab === t.id}
            onClick={() => {
              setTab(t.id);
              setListState("loading");
            }}
            className={cn(
              "inline-flex min-h-11 items-center rounded-full border px-4 py-2 text-sm font-medium transition-colors sm:min-h-0",
              tab === t.id
                ? "border-transparent bg-[color:var(--adfluence-action)] text-white"
                : "border-white/10 bg-[rgba(255,255,255,0.06)] text-[#CBD5E1] hover:text-white",
            )}
          >
            {t.label}
            {tab === t.id && listState === "ready" ? <span className="ml-2 rounded-full bg-white/20 px-1.5 text-xs">{total}</span> : null}
          </button>
        ))}
      </div>

      <Input
        id="business-search"
        aria-label="Search businesses"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Search by business, legal name, owner, TIN or email"
        className="h-11 max-w-md rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8]"
      />

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_420px]">
        <div className={cn("min-w-0 space-y-3", mobileShowDetail && "hidden xl:block")}>
          {listState === "loading" ? (
            <p className="flex items-center gap-2 p-6 text-sm text-[#94A3B8]">
              <Loader2 className="size-4 animate-spin" aria-hidden /> Loading businesses…
            </p>
          ) : listState === "error" ? (
            <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-6 text-sm text-rose-100">
              Could not load businesses.
              <Button type="button" variant="link" className="ml-2 h-auto p-0 text-white" onClick={() => void loadList()}>
                Try again
              </Button>
            </div>
          ) : rows.length === 0 ? (
            <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center text-sm text-[#CBD5E1]">
              {debouncedSearch ? `No businesses match “${debouncedSearch}”.` : tab === "pending" ? "No businesses are waiting for review." : "Nothing here yet."}
            </div>
          ) : (
            rows.map((r) => (
              <button
                key={r.user_id}
                type="button"
                onClick={() => select(r.user_id)}
                aria-current={selectedId === r.user_id ? "true" : undefined}
                className={cn(
                  "flex w-full min-w-0 gap-4 rounded-2xl border p-4 text-left transition-colors focus-visible:outline-2 focus-visible:outline-[color:var(--adfluence-orange)]",
                  selectedId === r.user_id
                    ? "border-[color:var(--adfluence-orange)]/35 bg-[color:var(--adfluence-orange)]/10"
                    : "border-white/10 bg-[rgba(255,255,255,0.08)] hover:border-white/20",
                )}
              >
                <UserAvatar src={r.logo_url} name={r.business_name} size="lg" />
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <p className="font-semibold text-white">{r.business_name}</p>
                    <span className={cn("rounded-full border px-2 py-0.5 text-[10px] font-medium", badge(r.verification_status))}>
                      {statusLabel[r.verification_status] ?? r.verification_status}
                    </span>
                  </div>
                  <p className="text-xs text-[#CBD5E1]">
                    {r.owner_full_name}
                    {r.industry ? ` · ${r.industry}` : ""}
                    {r.city ? ` · ${r.city}` : ""}
                  </p>
                  <p className="mt-1 text-xs text-[#94A3B8]">
                    {r.verification_status === "pending" && r.submitted_at ? `Submitted ${since(r.submitted_at)}` : r.email}
                  </p>
                </div>
              </button>
            ))
          )}
        </div>

        {selectedId ? (
          <aside
            aria-label="Business details"
            className={cn(
              "min-w-0 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 xl:sticky xl:top-4 xl:self-start",
              !mobileShowDetail && "hidden xl:block",
            )}
          >
            <Button
              type="button"
              variant="ghost"
              size="sm"
              onClick={() => setMobileShowDetail(false)}
              className="mb-3 -ml-2 text-white hover:bg-white/10 xl:hidden"
            >
              <ArrowLeft className="mr-1 size-4" aria-hidden /> Back to list
            </Button>
            {detailError ? (
              <p role="alert" className="text-sm text-rose-200">{detailError}</p>
            ) : !detail ? (
              <p className="flex items-center gap-2 text-sm text-[#94A3B8]">
                <Loader2 className="size-4 animate-spin" aria-hidden /> Loading details…
              </p>
            ) : (
              <div className="space-y-5">
                <div className="flex items-center gap-3">
                  <UserAvatar src={detail.logo_url} name={detail.business_name} size="xl" />
                  <div className="min-w-0">
                    <p className="text-lg font-semibold">{detail.business_name}</p>
                    <span className={cn("rounded-full border px-2 py-0.5 text-[10px] font-medium", badge(detail.verification_status))}>
                      {statusLabel[detail.verification_status] ?? detail.verification_status}
                    </span>
                  </div>
                </div>

                {detail.verification_note ? (
                  <p className="rounded-xl border border-white/10 bg-black/20 px-3 py-2 text-sm">
                    <span className="text-[#94A3B8]">Last note to the business: </span>
                    {detail.verification_note}
                  </p>
                ) : null}

                <div>
                  <h2 className="mb-1 text-xs font-semibold tracking-wide text-[#94A3B8] uppercase">Identity</h2>
                  <dl className="divide-y divide-white/5">
                    <Row label="Legal name">{detail.legal_business_name}</Row>
                    <Row label="TIN" highlight>{detail.tin_number}</Row>
                    <Row label="Type">{[detail.business_type, detail.business_size].filter(Boolean).join(" · ")}</Row>
                    <Row label="Owner">{detail.owner_full_name}</Row>
                  </dl>
                </div>

                <div>
                  <h2 className="mb-1 text-xs font-semibold tracking-wide text-[#94A3B8] uppercase">Contact</h2>
                  <dl className="divide-y divide-white/5">
                    <Row label="Email">
                      {detail.email} {detail.email_verified ? <span className="text-xs text-emerald-300">(confirmed)</span> : <span className="text-xs text-amber-300">(not confirmed)</span>}
                    </Row>
                    <Row label="Phone">
                      {detail.phone_number}{" "}
                      {detail.phone_number ? (
                        detail.phone_verified ? <span className="text-xs text-emerald-300">(confirmed)</span> : <span className="text-xs text-amber-300">(not confirmed)</span>
                      ) : null}
                    </Row>
                    <Row label="Location">{[detail.address, detail.city, detail.region].filter(Boolean).join(", ")}</Row>
                  </dl>
                </div>

                <div>
                  <h2 className="mb-1 text-xs font-semibold tracking-wide text-[#94A3B8] uppercase">Business</h2>
                  <dl className="divide-y divide-white/5">
                    <Row label="Industry">{detail.industry}</Row>
                    <Row label="Website"><ExtLink href={detail.website_url} /></Row>
                    <Row label="Social">
                      {[detail.social_links.instagram, detail.social_links.tiktok, detail.social_links.facebook, detail.social_links.linkedin].some(Boolean) ? (
                        <span className="flex flex-col gap-1">
                          <ExtLink href={detail.social_links.instagram} />
                          <ExtLink href={detail.social_links.tiktok} />
                          <ExtLink href={detail.social_links.facebook} />
                          <ExtLink href={detail.social_links.linkedin} />
                        </span>
                      ) : null}
                    </Row>
                  </dl>
                  {detail.description ? <p className="mt-2 text-sm whitespace-pre-wrap text-[#CBD5E1]">{detail.description}</p> : null}
                </div>

                <dl className="divide-y divide-white/5 text-xs">
                  <Row label="Registered">{formatDate(detail.registered_at)}</Row>
                  <Row label="Submitted">{formatDate(detail.submitted_at)}</Row>
                  {detail.verified_at ? <Row label="Verified">{formatDate(detail.verified_at)}</Row> : null}
                </dl>

                {decision ? (
                  <div className="space-y-3 rounded-xl border border-white/15 bg-black/20 p-3">
                    <p className="font-medium">{decisionCopy[decision].title}</p>
                    <label className="block space-y-1">
                      <span className="text-xs text-[#CBD5E1]">{decisionCopy[decision].label}</span>
                      <textarea
                        id="decision-text"
                        value={text}
                        maxLength={1000}
                        rows={3}
                        autoFocus
                        onChange={(e) => setText(e.target.value)}
                        placeholder={decisionCopy[decision].placeholder}
                        className="w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-white outline-none placeholder:text-[#94A3B8] focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
                      />
                      <span className="text-xs text-[#94A3B8]">
                        {text.trim().length}/1000{decision !== "verify" ? " · at least 5 characters" : ""}
                      </span>
                    </label>
                    <div className="flex flex-wrap gap-2">
                      <Button
                        type="button"
                        disabled={busy || reasonTooShort}
                        onClick={() => void confirm()}
                        className={cn(
                          "h-10 rounded-xl text-white",
                          decision === "verify" ? "bg-emerald-600 hover:bg-emerald-700" : "bg-rose-600 hover:bg-rose-700",
                        )}
                      >
                        {busy ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
                        {decisionCopy[decision].button}
                      </Button>
                      <Button type="button" variant="ghost" disabled={busy} onClick={() => setDecision(null)} className="h-10 rounded-xl text-white hover:bg-white/10">
                        Cancel
                      </Button>
                    </div>
                  </div>
                ) : (
                  <div className="flex flex-wrap gap-2">
                    {detail.verification_status === "rejected" ? (
                      <p className="w-full text-xs text-[#94A3B8]">
                        Not approved. The business can update its profile and resubmit. If you have changed your mind, you can
                        verify it now.
                      </p>
                    ) : null}
                    {detail.verification_status === "pending" || detail.verification_status === "rejected" ? (
                      <Button type="button" onClick={() => startDecision("verify")} className="h-10 rounded-xl bg-emerald-600 text-white hover:bg-emerald-700">
                        <BadgeCheck className="mr-1 size-4" aria-hidden /> {detail.verification_status === "rejected" ? "Verify anyway" : "Verify"}
                      </Button>
                    ) : null}
                    {detail.verification_status === "pending" ? (
                      <Button type="button" variant="outline" onClick={() => startDecision("reject")} className="h-10 rounded-xl border-rose-400/40 bg-transparent text-rose-200 hover:bg-rose-500/10">
                        <ShieldX className="mr-1 size-4" aria-hidden /> Do not approve
                      </Button>
                    ) : null}
                    {detail.verification_status === "verified" ? (
                      <Button type="button" variant="outline" onClick={() => startDecision("revoke")} className="h-10 rounded-xl border-rose-400/40 bg-transparent text-rose-200 hover:bg-rose-500/10">
                        <ShieldX className="mr-1 size-4" aria-hidden /> Revoke verification
                      </Button>
                    ) : null}
                  </div>
                )}
              </div>
            )}
          </aside>
        ) : (
          <aside className="hidden rounded-2xl border border-dashed border-white/10 p-8 text-center text-sm text-[#94A3B8] xl:block">
            Select a business to review its details.
          </aside>
        )}
      </div>
    </div>
  );
}
