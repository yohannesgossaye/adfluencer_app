"use client";

import { useCallback, useEffect, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  adminApi,
  type AdminDisputeListItem,
  type AdminDisputeDetail,
} from "@/lib/api/admin";
import { cn } from "@/lib/utils";

function sev(s: string) {
  const lower = s.toLowerCase();
  if (lower === "high") return "text-rose-300";
  if (lower === "medium") return "text-amber-200";
  return "text-[#CBD5E1]";
}

export function AdminDisputes() {
  const [rows, setRows] = useState<AdminDisputeListItem[]>([]);
  const [tab, setTab] = useState<"open" | "resolved">("open");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [selectedDetail, setSelectedDetail] = useState<AdminDisputeDetail | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [adminNoteDraft, setAdminNoteDraft] = useState("");
  const [resolutionNoteDraft, setResolutionNoteDraft] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  const loadDisputes = async () => {
    try {
      const res = await adminApi.listDisputes({ tab });
      if (res.data?.items) {
        setRows(res.data.items);
        if (!selectedId && res.data.items[0]) {
          setSelectedId(res.data.items[0].id);
        }
      }
    } catch {
      // fallback
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadDisputes();
  }, [tab]);

  useEffect(() => {
    if (!selectedId) return;
    async function loadItem() {
      try {
        const res = await adminApi.getDispute(selectedId!);
        if (res.data) {
          setSelectedDetail(res.data);
          setAdminNoteDraft(res.data.admin_note || "");
          setResolutionNoteDraft(res.data.resolution_note || "");
        }
      } catch {
        // fallback
      }
    }
    loadItem();
  }, [selectedId]);

  const setStatus = async (id: string, status: string) => {
    try {
      await adminApi.updateDisputeStatus(id, { status });
      await loadDisputes();
      if (selectedId === id) {
        const res = await adminApi.getDispute(id);
        if (res.data) setSelectedDetail(res.data);
      }
    } catch {
      // fallback
    }
  };

  const saveNotes = async () => {
    if (!selectedId) return;
    try {
      await adminApi.updateDisputeNotes(selectedId, {
        admin_note: adminNoteDraft,
        resolution_note: resolutionNoteDraft,
      });
      await loadDisputes();
    } catch {
      // fallback
    }
  };

  const blockUser = async () => {
    if (!selectedDetail?.influencer_user_id) return;
    try {
      await adminApi.blockUser(selectedDetail.influencer_user_id, { reason: "Blocked during dispute resolution" });
      window.alert("User account has been suspended.");
    } catch {
      // fallback
    }
  };

  const selectDispute = useCallback((id: string) => {
    setSelectedId(id);
    if (typeof window !== "undefined" && window.matchMedia("(max-width: 1279px)").matches) {
      setMobileShowDetail(true);
    }
  }, []);

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Disputes & reports</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Review reported conversations and disputes. Admin supervises escalated escrow claims.
        </p>
      </div>

      <div className="flex gap-2">
        {(["open", "resolved"] as const).map((t) => (
          <button
            key={t}
            type="button"
            onClick={() => setTab(t)}
            className={cn(
              "inline-flex min-h-11 items-center rounded-full border px-4 py-2 text-sm capitalize transition-colors sm:min-h-0",
              tab === t
                ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
                : "border-white/10 bg-white/5 text-[#CBD5E1]",
            )}
          >
            {t === "open" ? "Open & under review" : "Resolved"}
          </button>
        ))}
      </div>

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_400px]">
        <div className={cn("space-y-3", mobileShowDetail && "hidden xl:block")}>
          {rows.length === 0 && !isLoading ? (
            <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center">
              <p className="text-white">No disputes found.</p>
            </div>
          ) : (
            rows.map((d) => (
              <button
                key={d.id}
                type="button"
                onClick={() => selectDispute(d.id)}
                className={cn(
                  "w-full rounded-2xl border p-4 text-left transition-colors",
                  selectedId === d.id
                    ? "border-[color:var(--adfluence-orange)]/35 bg-[color:var(--adfluence-orange)]/10"
                    : "border-white/10 bg-[rgba(255,255,255,0.08)] hover:border-white/20",
                )}
              >
                <p className="font-semibold text-white">{d.title}</p>
                <p className="text-xs text-[#CBD5E1]">
                  {d.business_name} ↔ {d.influencer_name} · {d.campaign_title}
                </p>
                <p className="mt-1 text-xs text-[#94A3B8]">{d.issue_type_label || d.issue_type}</p>
                <p className={cn("text-xs font-medium", sev(d.severity))}>{d.severity} severity</p>
                <p className="mt-2 line-clamp-2 text-xs text-[#CBD5E1]">{d.preview}</p>
              </button>
            ))
          )}
        </div>

        {selectedDetail ? (
          <aside
            className={cn(
              "space-y-4 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 xl:sticky xl:top-4",
              !mobileShowDetail && "hidden xl:block",
            )}
          >
            <div className="xl:hidden">
              <Button
                type="button"
                variant="ghost"
                size="icon"
                className="shrink-0 text-white hover:bg-white/10"
                aria-label="Back to disputes"
                onClick={() => setMobileShowDetail(false)}
              >
                <ArrowLeft className="size-5" />
              </Button>
            </div>
            <h2 className="text-lg font-semibold text-white">{selectedDetail.title}</h2>
            <p className="text-xs text-[#CBD5E1]">
              Status: <span className="capitalize text-white">{selectedDetail.status.replace("_", " ")}</span>
            </p>
            <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="text-xs font-medium uppercase text-[#94A3B8]">Description & preview</p>
              <p className="mt-2 text-sm text-[#CBD5E1]">{selectedDetail.description || selectedDetail.preview}</p>
            </div>
            <label className="block text-xs">
              <span className="text-[#CBD5E1]">Admin internal note</span>
              <textarea
                value={adminNoteDraft}
                onChange={(e) => setAdminNoteDraft(e.target.value)}
                rows={2}
                className="mt-1 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-2 py-1 text-sm text-white"
              />
            </label>
            <label className="block text-xs">
              <span className="text-[#CBD5E1]">Resolution note</span>
              <textarea
                value={resolutionNoteDraft}
                onChange={(e) => setResolutionNoteDraft(e.target.value)}
                rows={2}
                className="mt-1 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-2 py-1 text-sm text-white"
              />
            </label>
            <Button
              type="button"
              variant="outline"
              size="sm"
              className="h-10 w-full border-white/15 text-white"
              onClick={saveNotes}
            >
              Save notes
            </Button>
            <div className="flex flex-col gap-2 pt-2">
              <Button
                type="button"
                size="sm"
                className="h-11 bg-white/10 text-white sm:h-9"
                onClick={() => setStatus(selectedDetail.id, "under_review")}
              >
                Mark under review
              </Button>
              <Button
                type="button"
                size="sm"
                className="h-11 bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:h-9"
                onClick={() => setStatus(selectedDetail.id, "resolved")}
              >
                Resolve dispute
              </Button>
              {selectedDetail.influencer_user_id ? (
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  className="h-11 border-rose-400/40 text-rose-200 sm:h-9"
                  onClick={blockUser}
                >
                  Block user
                </Button>
              ) : null}
            </div>
          </aside>
        ) : null}
      </div>
    </div>
  );
}
