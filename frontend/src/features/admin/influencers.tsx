"use client";

import Image from "next/image";
import { useCallback, useEffect, useMemo, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  adminApi,
  type InfluencerApplicationItem,
  type InfluencerApplicationDetail,
} from "@/lib/api/admin";
import { cn } from "@/lib/utils";

const tabs: Array<{ id: "all" | "pending" | "approved" | "rejected"; label: string }> = [
  { id: "all", label: "All" },
  { id: "pending", label: "Pending" },
  { id: "approved", label: "Approved" },
  { id: "rejected", label: "Rejected" },
];

function statusBadge(s: string) {
  const lower = s.toLowerCase();
  if (lower === "pending") return "border-amber-400/30 bg-amber-500/10 text-amber-200";
  if (lower === "approved" || lower === "active") return "border-emerald-400/30 bg-emerald-500/10 text-emerald-200";
  return "border-rose-400/30 bg-rose-500/10 text-rose-200";
}

export function AdminInfluencers() {
  const [rows, setRows] = useState<InfluencerApplicationItem[]>([]);
  const [tab, setTab] = useState<(typeof tabs)[number]["id"]>("all");
  const [search, setSearch] = useState("");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [selectedDetail, setSelectedDetail] = useState<InfluencerApplicationDetail | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [noteDraft, setNoteDraft] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  const loadList = async () => {
    try {
      const res = await adminApi.listInfluencerApplications({
        status: tab === "all" ? undefined : tab,
        search: search.trim() || undefined,
      });
      if (res.data?.items) {
        setRows(res.data.items);
        if (!selectedId && res.data.items[0]) {
          setSelectedId(res.data.items[0].user_id);
        }
      }
    } catch {
      // fallback
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadList();
  }, [tab, search]);

  useEffect(() => {
    if (!selectedId) return;
    async function loadItem() {
      try {
        const res = await adminApi.getInfluencerApplication(selectedId!);
        if (res.data) {
          setSelectedDetail(res.data);
          setNoteDraft(res.data.admin_note || "");
        }
      } catch {
        // fallback
      }
    }
    loadItem();
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
        languages: ["Amharic", "English"],
        portfolio_items: [],
      },
      completeness_percent: item.completeness_percent,
      admin_note: "",
      submitted_at: item.submitted_at,
    } as InfluencerApplicationDetail) : null;
  }, [selectedDetail, rows, selectedId]);

  const setStatus = async (userId: string, status: "approved" | "rejected") => {
    try {
      if (status === "approved") {
        await adminApi.approveInfluencerApplication(userId, { admin_note: noteDraft });
      } else {
        await adminApi.rejectInfluencerApplication(userId, { reason: "Application rejected by admin.", admin_note: noteDraft });
      }
      await loadList();
    } catch {
      // fallback
    }
  };

  const saveNotes = async () => {
    if (!selectedId) return;
    try {
      await adminApi.updateInfluencerApplicationNotes(selectedId, { admin_note: noteDraft });
      await loadList();
    } catch {
      // fallback
    }
  };

  const selectApplication = useCallback((userId: string) => {
    setSelectedId(userId);
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
                ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
                : "border-white/10 bg-[rgba(255,255,255,0.06)] text-[#CBD5E1] hover:text-white",
            )}
          >
            {t.label}
          </button>
        ))}
      </div>

      <Input
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Search by name, niche, or platform"
        className="h-11 max-w-md rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8]"
      />

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_380px]">
        <div className={cn("min-w-0 space-y-3", mobileShowDetail && "hidden xl:block")}>
          {rows.length === 0 && !isLoading ? (
            <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center">
              <p className="text-white">No influencer applications found.</p>
            </div>
          ) : (
            rows.map((r) => (
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
                <div className="relative h-16 w-16 shrink-0 overflow-hidden rounded-full border border-white/10">
                  <Image
                    src={r.profile_image_url || "/images/influencers/influencer-1.jpg"}
                    alt={r.creator_name}
                    fill
                    className="object-cover"
                    sizes="64px"
                  />
                </div>
                <div className="min-w-0 flex-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <p className="font-semibold text-white">{r.creator_name}</p>
                    <span
                      className={cn(
                        "rounded-full border px-2 py-0.5 text-[10px] font-medium capitalize",
                        statusBadge(r.status),
                      )}
                    >
                      {r.status}
                    </span>
                  </div>
                  <p className="text-xs text-[#CBD5E1]">
                    {r.primary_niche || "General"} · {r.primary_platform || "Multi"} · {r.total_followers.toLocaleString()} followers
                    {r.engagement_rate != null ? ` · ${r.engagement_rate}% engagement` : ""}
                  </p>
                  <p className="mt-1 text-xs text-[#94A3B8]">
                    Submitted {new Date(r.submitted_at).toLocaleDateString()}
                  </p>
                  <p className="mt-1 text-xs text-[#CBD5E1]">
                    Profile {r.completeness_percent}% complete
                  </p>
                  <div className="mt-2 flex flex-wrap gap-2">
                    <Button
                      type="button"
                      size="sm"
                      className="h-11 bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:h-8"
                      disabled={r.status !== "pending"}
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
                      disabled={r.status !== "pending"}
                      onClick={(e) => {
                        e.stopPropagation();
                        setStatus(r.user_id, "rejected");
                      }}
                    >
                      Reject
                    </Button>
                  </div>
                </div>
              </div>
            ))
          )}
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
                <span className="text-white">Pricing:</span> ETB {selected.profile.min_price_etb.toLocaleString()} - {selected.profile.max_price_etb?.toLocaleString()}
              </p>
            ) : null}

            {selected.profile?.portfolio_items && selected.profile.portfolio_items.length > 0 ? (
              <>
                <p className="mt-3 text-xs font-medium uppercase text-[#94A3B8]">Sample content</p>
                <div className="relative mt-2 aspect-video overflow-hidden rounded-xl border border-white/10">
                  <Image
                    src={selected.profile.portfolio_items[0].file_url}
                    alt={selected.profile.creator_name}
                    fill
                    className="object-cover"
                    sizes="360px"
                  />
                </div>
              </>
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
              onClick={saveNotes}
            >
              Save notes
            </Button>
          </aside>
        ) : null}
      </div>
    </div>
  );
}
