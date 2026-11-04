"use client";

import Image from "next/image";
import { useCallback, useEffect, useMemo, useState } from "react";
import axios from "axios";
import { ArrowLeft } from "lucide-react";
import { toast } from "sonner";

import { ListStatus, LoadMore } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  adminApi,
  type InfluencerApplicationItem,
  type InfluencerApplicationDetail,
} from "@/lib/api/admin";
import { usePagedList } from "@/hooks/use-paged-list";
import { cn } from "@/lib/utils";
import { UserAvatar } from "@/components/common/user-avatar";
import { formatDate, formatEtbRange } from "@/lib/contracts";
import { StatusBadge } from "@/components/common/status-badge";
import { creatorStatus } from "@/lib/status";

const tabs: Array<{ id: "all" | "pending" | "approved" | "rejected"; label: string }> = [
  { id: "all", label: "All" },
  { id: "pending", label: "Pending" },
  { id: "approved", label: "Approved" },
  { id: "rejected", label: "Rejected" },
];


const canApprove = (status: string) => status === "pending" || status === "pending_approval";
// The backend rejects only submitted or approved creators.
const canReject = (status: string) => canApprove(status) || status === "active" || status === "approved";


export function AdminInfluencers() {
  const [tab, setTab] = useState<(typeof tabs)[number]["id"]>("pending");
  const [search, setSearch] = useState("");
  const [query, setQuery] = useState("");
  const [pickedId, setPickedId] = useState<string | null>(null);
  const [selectedDetail, setSelectedDetail] = useState<InfluencerApplicationDetail | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [noteDraft, setNoteDraft] = useState("");
  const [rejectFor, setRejectFor] = useState<string | null>(null);
  const [rejectReason, setRejectReason] = useState("");
  const [busyId, setBusyId] = useState<string | null>(null);
  const [savingNotes, setSavingNotes] = useState(false);

  useEffect(() => {
    const t = setTimeout(() => setQuery(search.trim()), 300);
    return () => clearTimeout(t);
  }, [search]);

  const errorText = (err: unknown, fallback: string) =>
    (axios.isAxiosError(err) ? (err.response?.data as { message?: string } | undefined)?.message : undefined) || fallback;

  const fetchPage = useCallback(
    async (page: number) =>
      (
        await adminApi.listInfluencerApplications({
          status: tab === "all" ? undefined : tab,
          search: query || undefined,
          page,
        })
      ).data,
    [tab, query],
  );
  const list = usePagedList(fetchPage);
  const rows: InfluencerApplicationItem[] = list.items;
  const loadList = list.reload;
  // The selection follows the list: if the picked creator isn't in it (other tab or search), show the first one.
  const selectedId = rows.some((r) => r.user_id === pickedId) ? pickedId : (rows[0]?.user_id ?? null);

  const loadDetail = async (userId: string) => {
    try {
      const res = await adminApi.getInfluencerApplication(userId);
      if (res.data) {
        setSelectedDetail(res.data);
        setNoteDraft(res.data.admin_note || "");
      }
    } catch (err) {
      toast.error(errorText(err, "This application could not be opened."));
    }
  };

  useEffect(() => {
    if (selectedId) void Promise.resolve(selectedId).then(loadDetail);
    // eslint-disable-next-line react-hooks/exhaustive-deps -- reload only when the selection changes
  }, [selectedId]);

  const selected = useMemo(() => {
    if (selectedDetail && selectedDetail.user_id === selectedId) {
      return selectedDetail;
    }
    const item = rows.find((r) => r.user_id === selectedId);
    return item ? ({
      user_id: item.user_id,
      status: item.status,
      email: "",
      profile: {
        creator_name: item.creator_name,
        profile_image_url: item.profile_image_url,
        niches: item.primary_niche ? [{ id: "1", name: item.primary_niche, slug: "" }] : [],
        platforms: item.primary_platform ? [item.primary_platform] : [],
        social_links: {},
        total_followers: item.total_followers,
        engagement_rate: item.engagement_rate,
        languages: [],
        portfolio_items: [],
      },
      completeness_percent: item.completeness_percent,
      admin_note: "",
      submitted_at: item.submitted_at,
    } as InfluencerApplicationDetail) : null;
  }, [selectedDetail, rows, selectedId]);

  const setStatus = async (userId: string, status: "approved" | "rejected") => {
    if (status === "rejected") {
      // Rejection needs a reason the creator will see: open the reason box.
      setPickedId(userId);
      setRejectFor(userId);
      setRejectReason("");
      if (window.matchMedia("(max-width: 1279px)").matches) setMobileShowDetail(true);
      return;
    }
    setBusyId(userId);
    try {
      // The notes box belongs to the selected creator; don't copy it onto another one.
      await adminApi.approveInfluencerApplication(userId, userId === selectedId ? { admin_note: noteDraft } : {});
      toast.success("Creator approved. They can now receive invites.");
      await loadList();
      if (userId === selectedId) await loadDetail(userId);
    } catch (err) {
      toast.error(errorText(err, "Could not approve this creator."));
    } finally {
      setBusyId(null);
    }
  };

  const sendRejection = async () => {
    if (!rejectFor || rejectReason.trim().length < 5) return;
    const userId = rejectFor;
    setBusyId(userId);
    try {
      await adminApi.rejectInfluencerApplication(userId, { reason: rejectReason.trim(), admin_note: noteDraft });
      toast.success("Rejection sent. The creator will see your reason.");
      setRejectFor(null);
      setRejectReason("");
      await loadList();
      await loadDetail(userId);
    } catch (err) {
      toast.error(errorText(err, "Could not reject this application."));
    } finally {
      setBusyId(null);
    }
  };

  const saveNotes = async () => {
    if (!selectedId) return;
    setSavingNotes(true);
    try {
      await adminApi.updateInfluencerApplicationNotes(selectedId, { admin_note: noteDraft });
      toast.success("Notes saved.");
      await loadList();
    } catch (err) {
      toast.error(errorText(err, "The notes could not be saved."));
    } finally {
      setSavingNotes(false);
    }
  };

  const selectApplication = useCallback((userId: string) => {
    setPickedId(userId);
    if (typeof window !== "undefined" && window.matchMedia("(max-width: 1279px)").matches) {
      setMobileShowDetail(true);
    }
  }, []);

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Influencer approvals</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Review creator applications before they appear on the marketplace.
        </p>
      </div>

      <div className="flex flex-wrap gap-2">
        {tabs.map((t) => (
          <button
            key={t.id}
            type="button"
            onClick={() => setTab(t.id)}
            className={cn(
              "inline-flex min-h-11 items-center rounded-full border px-4 py-2 text-sm font-medium transition-colors sm:min-h-0",
              tab === t.id
                ? "border-transparent bg-[color:var(--adfluence-action)] text-white"
                : "border-white/10 bg-[rgba(255,255,255,0.06)] text-[#CBD5E1] hover:text-white",
            )}
          >
            {t.label}
          </button>
        ))}
      </div>

      <Input
        aria-label="Search creators"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Search by name, niche, or platform"
        className="h-11 max-w-md rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8]"
      />

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <div className={cn("min-w-0 space-y-3", mobileShowDetail && selected && "hidden xl:block")}>
          <ListStatus
            state={list.state}
            isEmpty={rows.length === 0}
            emptyText="No creators match this tab or search."
            errorText="Applications could not be loaded."
            onRetry={list.retry}
          >
            {rows.map((r) => (
              <div
                key={r.user_id}
                role="button"
                tabIndex={0}
                onClick={() => selectApplication(r.user_id)}
                onKeyDown={(e) => {
                  if (e.key === "Enter" || e.key === " ") {
                    e.preventDefault();
                    selectApplication(r.user_id);
                  }
                }}
                className={cn(
                  "flex w-full min-w-0 cursor-pointer gap-4 rounded-2xl border p-4 text-left transition-colors",
                  selectedId === r.user_id
                    ? "border-[color:var(--adfluence-orange)]/35 bg-[color:var(--adfluence-orange)]/10"
                    : "border-white/10 bg-[rgba(255,255,255,0.08)] hover:border-white/20",
                )}
              >
                <UserAvatar
                  src={r.profile_image_url}
                  name={r.creator_name}
                  size="xl"
                  fallback="CR"
                />
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <p className="font-semibold text-white">{r.creator_name}</p>
                    <StatusBadge status={creatorStatus(r.status)} />
                  </div>
                  <p className="text-xs text-[#CBD5E1]">
                    {r.primary_niche || "General"} · {r.primary_platform || "Multi"} · {r.total_followers.toLocaleString()} followers
                    {r.engagement_rate != null ? ` · ${r.engagement_rate}% engagement` : ""}
                  </p>
                  <p className="mt-1 text-xs text-[#94A3B8]">
                    {r.status === "pending_profile" ? "Last updated" : "Submitted"} {formatDate(r.submitted_at)}
                  </p>
                  <p className="mt-1 text-xs text-[#CBD5E1]">
                    Profile {r.completeness_percent}% complete
                  </p>
                  <div className="mt-2 flex flex-wrap gap-2">
                    <Button
                      type="button"
                      size="sm"
                      className="h-11 bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)] sm:h-8"
                      disabled={!canApprove(r.status) || busyId === r.user_id}
                      onClick={(e) => {
                        e.stopPropagation();
                        setStatus(r.user_id, "approved");
                      }}
                    >
                      Approve
                    </Button>
                    <Button
                      type="button"
                      size="sm"
                      variant="outline"
                      className="h-11 border-white/15 text-white hover:bg-white/10 sm:h-8"
                      disabled={!canReject(r.status) || busyId === r.user_id}
                      onClick={(e) => {
                        e.stopPropagation();
                        setStatus(r.user_id, "rejected");
                      }}
                    >
                      {r.status === "approved" || r.status === "active" ? "Request changes" : "Reject"}
                    </Button>
                  </div>
                  {r.status === "pending_profile" ? (
                    <p className="mt-2 text-xs text-[#94A3B8]">
                      Waiting for the creator to complete and submit their profile. You can approve it once they submit.
                    </p>
                  ) : null}
                </div>
              </div>
            ))}
            <LoadMore
              shown={rows.length}
              total={list.total}
              hasMore={list.hasMore}
              loading={list.loadingMore}
              onLoadMore={list.loadMore}
            />
          </ListStatus>
        </div>

        {selected ? (
          <aside
            className={cn(
              "rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 xl:sticky xl:top-4 xl:self-start",
              !mobileShowDetail && "hidden xl:block",
            )}
          >
            <div className="mb-3 xl:hidden">
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="shrink-0 text-white hover:bg-white/10"
                aria-label="Back to applications"
                onClick={() => setMobileShowDetail(false)}
              >
                <ArrowLeft className="size-5" />
              </Button>
            </div>
            <h2 className="text-lg font-semibold text-white">Details — {selected.profile?.creator_name || "Creator"}</h2>
            <p className="mt-1">
              <StatusBadge status={creatorStatus(selected.status)} />
            </p>
            {canApprove(selected.status) || canReject(selected.status) ? (
              <div className="mt-3 flex flex-wrap gap-2">
                {canApprove(selected.status) ? (
                  <Button
                    type="button"
                    disabled={busyId === selected.user_id}
                    className="h-11 flex-1 bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)] sm:h-9"
                    onClick={() => void setStatus(selected.user_id, "approved")}
                  >
                    Approve
                  </Button>
                ) : null}
                {rejectFor !== selected.user_id ? (
                  <Button
                    type="button"
                    variant="outline"
                    disabled={busyId === selected.user_id}
                    className="h-11 flex-1 border-white/15 text-white hover:bg-white/10 sm:h-9"
                    onClick={() => void setStatus(selected.user_id, "rejected")}
                  >
                    {selected.status === "approved" || selected.status === "active" ? "Request changes" : "Reject"}
                  </Button>
                ) : null}
              </div>
            ) : null}
            <p className="mt-2 text-sm text-[#CBD5E1]">{selected.profile?.bio || "No bio provided."}</p>
            <div className="mt-3 space-y-1 text-xs text-[#94A3B8]">
              <p>TikTok: {selected.profile?.social_links?.tiktok || "—"}</p>
              <p>Instagram: {selected.profile?.social_links?.instagram || "—"}</p>
              <p>YouTube: {selected.profile?.social_links?.youtube || "—"}</p>
            </div>
            {selected.profile?.audience_description ? (
              <p className="mt-3 text-sm text-[#CBD5E1]">
                <span className="text-white">Audience:</span> {selected.profile.audience_description}
              </p>
            ) : null}
            {selected.profile?.min_price_etb != null ? (
              <p className="mt-2 text-sm text-[#CBD5E1]">
                <span className="text-white">Pricing:</span> {formatEtbRange(selected.profile.min_price_etb, selected.profile.max_price_etb)}
              </p>
            ) : null}

            {selected.profile?.portfolio_items && selected.profile.portfolio_items.length > 0 ? (
              <>
                <p className="mt-3 text-xs font-medium uppercase text-[#94A3B8]">Sample content</p>
                <div className="relative mt-2 aspect-video overflow-hidden rounded-xl border border-white/10 bg-black/40">
                  {selected.profile.portfolio_items[0].media_type === "video" ? (
                    <video
                      src={selected.profile.portfolio_items[0].file_url}
                      controls
                      className="h-full w-full object-cover"
                    />
                  ) : (
                    <Image
                      src={selected.profile.portfolio_items[0].file_url}
                      alt={selected.profile.creator_name}
                      fill
                      unoptimized
                      className="object-cover"
                      sizes="360px"
                    />
                  )}
                </div>
              </>
            ) : null}

            {rejectFor === selected.user_id ? (
              <div className="mt-4 space-y-2 rounded-xl border border-rose-400/30 bg-rose-500/10 p-3">
                <label className="block text-sm">
                  <span className="text-rose-100">Reason (shown to the creator)</span>
                  <textarea
                    id="reject-reason"
                    value={rejectReason}
                    autoFocus
                    maxLength={1000}
                    rows={3}
                    onChange={(e) => setRejectReason(e.target.value)}
                    placeholder="Tell them what to fix, e.g. add a link to your TikTok profile"
                    className="mt-1 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-white"
                  />
                  <span className="text-xs text-rose-100/80">{rejectReason.trim().length}/1000 · at least 5 characters</span>
                </label>
                <div className="flex flex-wrap gap-2">
                  <Button type="button" disabled={rejectReason.trim().length < 5 || busyId === rejectFor} onClick={() => void sendRejection()} className="h-11 bg-rose-600 text-white hover:bg-rose-700 sm:h-9">
                    Send rejection
                  </Button>
                  <Button type="button" variant="ghost" onClick={() => setRejectFor(null)} className="h-11 text-white hover:bg-white/10 sm:h-9">
                    Cancel
                  </Button>
                </div>
              </div>
            ) : null}

            <label className="mt-4 block text-sm">
              <span className="text-[#CBD5E1]">Admin notes</span>
              <textarea
                value={noteDraft}
                onChange={(e) => setNoteDraft(e.target.value)}
                rows={3}
                className="mt-1 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-2 text-sm text-white"
              />
            </label>
            <Button
              type="button"
              variant="outline"
              className="mt-2 h-11 w-full border-white/15 text-white hover:bg-white/10 sm:h-9"
              disabled={savingNotes}
              onClick={saveNotes}
            >
              {savingNotes ? "Saving…" : "Save notes"}
            </Button>
          </aside>
        ) : null}
      </div>
    </div>
  );
}
