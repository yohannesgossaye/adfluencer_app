"use client";

import { useCallback, useEffect, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { toast } from "sonner";

import { ConfirmActionDialog, type ConfirmActionRequest } from "@/components/admin/confirm-action-dialog";
import { Button } from "@/components/ui/button";
import {
  adminApi,
  type AdminDisputeListItem,
  type AdminDisputeDetail,
} from "@/lib/api/admin";
import { adminErrorMessage } from "@/lib/api-errors";
import { cn } from "@/lib/utils";

const STATUS_LABEL: Record<string, string> = {
  open: "Open",
  under_review: "Under review",
  escalated: "Escalated",
  resolved: "Resolved",
  closed: "Closed",
};

function isFinished(status: string) {
  return status === "resolved" || status === "closed";
}

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

  const [confirm, setConfirm] = useState<ConfirmActionRequest | null>(null);
  const [savingNotes, setSavingNotes] = useState(false);

  const refreshSelected = async (id: string) => {
    const res = await adminApi.getDispute(id);
    if (res.data) setSelectedDetail(res.data);
  };

  const setStatus = async (id: string, status: string, success: string) => {
    try {
      await adminApi.updateDisputeStatus(id, { status });
      toast.success(success);
      await loadDisputes();
      if (selectedId === id) await refreshSelected(id);
    } catch (err) {
      toast.error(adminErrorMessage(err, "The dispute could not be updated."));
      throw err;
    }
  };

  const saveNotes = async () => {
    if (!selectedId) return;
    setSavingNotes(true);
    try {
      await adminApi.updateDisputeNotes(selectedId, {
        admin_note: adminNoteDraft,
        resolution_note: resolutionNoteDraft,
      });
      toast.success("Notes saved.");
      await loadDisputes();
    } catch (err) {
      toast.error(adminErrorMessage(err, "The notes could not be saved."));
    } finally {
      setSavingNotes(false);
    }
  };

  const askResolve = (d: AdminDisputeDetail) =>
    setConfirm({
      title: "Resolve this dispute?",
      description: "Record what was decided. The dispute moves to Resolved; you can reopen it later if needed.",
      confirmLabel: "Resolve dispute",
      reasonLabel: "Resolution note",
      reasonPlaceholder: "e.g. Creator re-posted the video; business confirmed. No refund.",
      onConfirm: async (note) => {
        try {
          await adminApi.updateDisputeNotes(d.id, { admin_note: adminNoteDraft, resolution_note: note });
        } catch (err) {
          toast.error(adminErrorMessage(err, "The resolution note could not be saved."));
          throw err;
        }
        setResolutionNoteDraft(note);
        await setStatus(d.id, "resolved", "Dispute resolved.");
      },
    });

  const askBlock = (party: "business" | "creator", d: AdminDisputeDetail) => {
    const userId = party === "business" ? d.business_user_id : d.influencer_user_id;
    const name = party === "business" ? d.business_name : d.influencer_name;
    setConfirm({
      title: `Block ${name}?`,
      description: `The ${party} is signed out everywhere and can't log in until unblocked from Users. The dispute stays as it is.`,
      confirmLabel: `Block ${party}`,
      tone: "danger",
      reasonLabel: "Reason (kept on the account)",
      reasonPlaceholder: "e.g. Took payment without delivering (dispute)",
      onConfirm: async (reason) => {
        try {
          await adminApi.blockUser(userId, { reason });
          toast.success(`${name} is blocked.`);
        } catch (err) {
          toast.error(adminErrorMessage(err, "The account could not be blocked."));
          throw err;
        }
      },
    });
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
                <p className={cn("text-xs font-medium capitalize", sev(d.severity))}>
                  {d.severity} severity · {STATUS_LABEL[d.status] ?? d.status}
                </p>
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
              Status: <span className="text-white">{STATUS_LABEL[selectedDetail.status] ?? selectedDetail.status}</span>
            </p>
            <p className="text-xs text-[#CBD5E1]">
              Between <span className="text-white">{selectedDetail.business_name}</span> (business) and{" "}
              <span className="text-white">{selectedDetail.influencer_name}</span> (creator)
              {selectedDetail.against_user_id ? (
                <>
                  {" "}· reported:{" "}
                  <span className="text-white">
                    {selectedDetail.against_user_id === selectedDetail.business_user_id
                      ? selectedDetail.business_name
                      : selectedDetail.influencer_name}
                  </span>
                </>
              ) : null}
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
              disabled={savingNotes}
              onClick={saveNotes}
            >
              {savingNotes ? "Saving…" : "Save notes"}
            </Button>
            <div className="flex flex-col gap-2 pt-2">
              {isFinished(selectedDetail.status) ? (
                <Button
                  type="button"
                  size="sm"
                  className="h-11 bg-white/10 text-white sm:h-9"
                  onClick={() =>
                    setConfirm({
                      title: "Reopen this dispute?",
                      description: "It moves back to Open & under review.",
                      confirmLabel: "Reopen",
                      onConfirm: () => setStatus(selectedDetail.id, "under_review", "Dispute reopened."),
                    })
                  }
                >
                  Reopen dispute
                </Button>
              ) : (
                <>
                  {selectedDetail.status !== "under_review" ? (
                    <Button
                      type="button"
                      size="sm"
                      className="h-11 bg-white/10 text-white sm:h-9"
                      onClick={() =>
                        void setStatus(selectedDetail.id, "under_review", "Marked as under review.").catch(() => {})
                      }
                    >
                      Mark under review
                    </Button>
                  ) : null}
                  <Button
                    type="button"
                    size="sm"
                    className="h-11 bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:h-9"
                    onClick={() => askResolve(selectedDetail)}
                  >
                    Resolve dispute…
                  </Button>
                </>
              )}
              <div className="grid gap-2 sm:grid-cols-2">
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  className="h-11 border-rose-400/40 text-rose-200 sm:h-9"
                  onClick={() => askBlock("business", selectedDetail)}
                >
                  Block business…
                </Button>
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  className="h-11 border-rose-400/40 text-rose-200 sm:h-9"
                  onClick={() => askBlock("creator", selectedDetail)}
                >
                  Block creator…
                </Button>
              </div>
            </div>
          </aside>
        ) : null}
      </div>

      <ConfirmActionDialog request={confirm} onClose={() => setConfirm(null)} />
    </div>
  );
}
