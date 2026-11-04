"use client";

import { useCallback, useEffect, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { toast } from "sonner";

import { ConfirmActionDialog, type ConfirmActionRequest } from "@/components/common/confirm-action-dialog";
import { ListStatus, LoadMore } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import {
  adminApi,
  type AdminDisputeListItem,
  type AdminDisputeDetail,
} from "@/lib/api/admin";
import { adminErrorMessage } from "@/lib/api-errors";
import { usePagedList } from "@/hooks/use-paged-list";
import { cn } from "@/lib/utils";
import { StatusBadge } from "@/components/common/status-badge";
import { disputeStatus, priorityLevel } from "@/lib/status";
import { FilterChips } from "@/components/common/filter-chips";
import { contractStatus, formatDate } from "@/lib/contracts";
import { outcomeLabel } from "@/lib/disputes";
import { ResolveDisputeDialog } from "./resolve-dispute-dialog";


function isFinished(status: string) {
  return status === "resolved" || status === "closed";
}


export function AdminDisputes() {
  const [tab, setTab] = useState<"open" | "resolved">("open");
  const [pickedId, setPickedId] = useState<string | null>(null);
  const [selectedDetail, setSelectedDetail] = useState<AdminDisputeDetail | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);
  const [adminNoteDraft, setAdminNoteDraft] = useState("");
  const [resolving, setResolving] = useState<AdminDisputeDetail | null>(null);

  const fetchPage = useCallback(async (page: number) => (await adminApi.listDisputes({ tab, page })).data, [tab]);
  const list = usePagedList(fetchPage);
  const rows: AdminDisputeListItem[] = list.items;
  const loadDisputes = list.reload;
  // The selection follows the list: a dispute from the other tab is never shown.
  const selectedId = rows.some((r) => r.id === pickedId) ? pickedId : (rows[0]?.id ?? null);
  const shownDetail = selectedDetail?.id === selectedId ? selectedDetail : null;

  useEffect(() => {
    if (!selectedId) return;
    async function loadItem() {
      try {
        const res = await adminApi.getDispute(selectedId!);
        if (res.data) {
          setSelectedDetail(res.data);
          setAdminNoteDraft(res.data.admin_note || "");
        }
      } catch (err) {
        toast.error(adminErrorMessage(err, "This dispute could not be opened."));
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
      await adminApi.updateDisputeNotes(selectedId, { admin_note: adminNoteDraft });
      toast.success("Note saved.");
      await loadDisputes();
    } catch (err) {
      toast.error(adminErrorMessage(err, "The note could not be saved."));
    } finally {
      setSavingNotes(false);
    }
  };

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
    setPickedId(id);
    if (typeof window !== "undefined" && window.matchMedia("(max-width: 1279px)").matches) {
      setMobileShowDetail(true);
    }
  }, []);

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Disputes & reports</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Problems reported on contracts. The payment stays held until you decide: continue the contract, complete it and pay the
          creator, or end it.
        </p>
      </div>

      <FilterChips
        label="Disputes"
        value={tab}
        onChange={setTab}
        options={[
          { value: "open", label: "Open & under review" },
          { value: "resolved", label: "Resolved" },
        ]}
      />

      <div className="grid gap-4 xl:grid-cols-[minmax(0,1fr)_400px]">
        <div className={cn("space-y-3", mobileShowDetail && shownDetail && "hidden xl:block")}>
          <ListStatus
            state={list.state}
            isEmpty={rows.length === 0}
            emptyText={tab === "open" ? "No open disputes." : "No resolved disputes yet."}
            errorText="Disputes could not be loaded."
            onRetry={list.retry}
          >
            {rows.map((d) => (
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
                <div className="mt-2 flex flex-wrap gap-1.5">
                  <StatusBadge status={disputeStatus(d.status)} />
                  <StatusBadge status={{ ...priorityLevel(d.severity), label: `${priorityLevel(d.severity).label} severity` }} />
                </div>
                <p className="mt-2 line-clamp-2 text-xs text-[#CBD5E1]">{d.preview}</p>
              </button>
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

        {shownDetail ? (
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
            <h2 className="text-lg font-semibold text-white">{shownDetail.title}</h2>
            <div className="flex flex-wrap gap-1.5">
              <StatusBadge status={disputeStatus(shownDetail.status)} />
              <StatusBadge status={{ ...contractStatus(shownDetail.contract_status, "business").badge, label: `Contract: ${contractStatus(shownDetail.contract_status, "business").label}` }} />
            </div>
            <dl className="space-y-1 text-xs text-[#CBD5E1]">
              <div>
                <dt className="inline text-[#94A3B8]">Problem: </dt>
                <dd className="inline text-white">{shownDetail.issue_type_label}</dd>
              </div>
              <div>
                <dt className="inline text-[#94A3B8]">Reported by: </dt>
                <dd className="inline text-white">
                  {shownDetail.opened_by_role === "business"
                    ? `${shownDetail.business_name} (business)`
                    : `${shownDetail.influencer_name} (creator)`}{" "}
                  on {formatDate(shownDetail.created_at)}
                </dd>
              </div>
              <div>
                <dt className="inline text-[#94A3B8]">Against: </dt>
                <dd className="inline text-white">
                  {shownDetail.opened_by_role === "business"
                    ? `${shownDetail.influencer_name} (creator)`
                    : `${shownDetail.business_name} (business)`}
                </dd>
              </div>
            </dl>
            <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="text-xs font-medium uppercase text-[#94A3B8]">What they said</p>
              <p className="mt-2 text-sm whitespace-pre-wrap text-[#CBD5E1]">{shownDetail.description || shownDetail.preview}</p>
            </div>
            {isFinished(shownDetail.status) ? (
              <div className="rounded-xl border border-emerald-400/30 bg-emerald-500/10 p-3">
                <p className="text-sm font-semibold text-emerald-50">Decision: {outcomeLabel(shownDetail.outcome)}</p>
                {shownDetail.resolved_at ? <p className="text-xs text-emerald-100/80">{formatDate(shownDetail.resolved_at)}</p> : null}
                {shownDetail.resolution_note ? (
                  <p className="mt-2 text-sm whitespace-pre-wrap text-[#E2E8F0]">{shownDetail.resolution_note}</p>
                ) : null}
              </div>
            ) : null}
            <label className="block text-xs">
              <span className="text-[#CBD5E1]">Internal note (only admins see this)</span>
              <textarea
                value={adminNoteDraft}
                onChange={(e) => setAdminNoteDraft(e.target.value)}
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
              {savingNotes ? "Saving…" : "Save note"}
            </Button>
            <div className="flex flex-col gap-2 pt-2">
              {isFinished(shownDetail.status) ? null : (
                <>
                  {shownDetail.status !== "under_review" ? (
                    <Button
                      type="button"
                      size="sm"
                      className="h-11 bg-white/10 text-white sm:h-9"
                      onClick={() =>
                        void setStatus(shownDetail.id, "under_review", "Marked as under review.").catch(() => {})
                      }
                    >
                      Mark under review
                    </Button>
                  ) : null}
                  <Button
                    type="button"
                    size="sm"
                    className="h-11 bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)] sm:h-9"
                    onClick={() => setResolving(shownDetail)}
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
                  onClick={() => askBlock("business", shownDetail)}
                >
                  Block business…
                </Button>
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  className="h-11 border-rose-400/40 text-rose-200 sm:h-9"
                  onClick={() => askBlock("creator", shownDetail)}
                >
                  Block creator…
                </Button>
              </div>
            </div>
          </aside>
        ) : null}
      </div>

      <ConfirmActionDialog request={confirm} onClose={() => setConfirm(null)} />
      <ResolveDisputeDialog
        dispute={resolving}
        onClose={() => setResolving(null)}
        onResolved={async () => {
          const id = resolving?.id;
          await loadDisputes();
          if (id) await refreshSelected(id).catch(() => {});
        }}
      />
    </div>
  );
}
