"use client";

import { useCallback, useEffect, useMemo, useState } from "react";
import { ArrowLeft } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  loadDisputes,
  saveDisputes,
  type AdminDispute,
  type DisputeStatus,
} from "@/data/admin-disputes-mock";
import { blockUserEmail } from "@/lib/blocked-users";
import { cn } from "@/lib/utils";

function sev(s: AdminDispute["severity"]) {
  if (s === "high") return "text-rose-300";
  if (s === "medium") return "text-amber-200";
  return "text-[#CBD5E1]";
}

export function AdminDisputes() {
  const [rows, setRows] = useState<AdminDispute[]>([]);
  const [tab, setTab] = useState<"open" | "resolved">("open");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [mobileShowDetail, setMobileShowDetail] = useState(false);

  useEffect(() => {
    const list = loadDisputes();
    setRows(list);
    if (list[0]) setSelectedId((p) => p ?? list[0].id);
  }, []);

  const selected = useMemo(() => rows.find((r) => r.id === selectedId) ?? null, [rows, selectedId]);

  const persist = (next: AdminDispute[]) => {
    saveDisputes(next);
    setRows(next);
  };

  const patchRow = (id: string, patch: Partial<AdminDispute>) => {
    setRows((prev) => {
      const next = prev.map((r) => (r.id === id ? { ...r, ...patch } : r));
      saveDisputes(next);
      return next;
    });
  };

  const setStatus = (id: string, status: DisputeStatus) => patchRow(id, { status });

  const list = useMemo(() => {
    if (tab === "resolved") return rows.filter((r) => r.status === "resolved");
    return rows.filter((r) => r.status === "open" || r.status === "under_review");
  }, [rows, tab]);

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
          Review reported conversations and disputes. Admin does not monitor all chats — only
          escalated items.
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
          {list.map((d) => (
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
                {d.business} ↔ {d.influencer} · {d.campaign}
              </p>
              <p className="mt-1 text-xs text-[#94A3B8]">{d.issueType}</p>
              <p className={cn("text-xs font-medium", sev(d.severity))}>{d.severity} severity</p>
              <p className="mt-2 line-clamp-2 text-xs text-[#CBD5E1]">{d.preview}</p>
            </button>
          ))}
        </div>

        {selected ? (
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
            <h2 className="text-lg font-semibold text-white">{selected.title}</h2>
            <p className="text-xs text-[#CBD5E1]">
              Status: <span className="capitalize text-white">{selected.status.replace("_", " ")}</span>
            </p>
            <div className="rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] p-3">
              <p className="text-xs font-medium uppercase text-[#94A3B8]">Conversation preview</p>
              <ul className="mt-2 space-y-2 text-sm text-[#CBD5E1]">
                {selected.messages.map((m) => (
                  <li key={m.id}>
                    <span className="text-white">{m.from}:</span> {m.body}
                  </li>
                ))}
              </ul>
            </div>
            <label className="block text-xs">
              <span className="text-[#CBD5E1]">Admin note</span>
              <textarea
                value={selected.adminNote}
                onChange={(e) => patchRow(selected.id, { adminNote: e.target.value })}
                rows={2}
                className="mt-1 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-2 py-1 text-sm text-white"
              />
            </label>
            <label className="block text-xs">
              <span className="text-[#CBD5E1]">Resolution note</span>
              <textarea
                value={selected.resolutionNote}
                onChange={(e) => patchRow(selected.id, { resolutionNote: e.target.value })}
                rows={2}
                className="mt-1 w-full rounded-xl border border-white/10 bg-[color:var(--adfluence-surface)] px-2 py-1 text-sm text-white"
              />
            </label>
            <div className="flex flex-col gap-2">
              <Button
                type="button"
                size="sm"
                className="h-11 bg-white/10 text-white sm:h-9"
                onClick={() => window.alert("Conversation viewer (demo)")}
              >
                View conversation
              </Button>
              <Button
                type="button"
                size="sm"
                className="h-11 bg-white/10 text-white sm:h-9"
                onClick={() => setStatus(selected.id, "under_review")}
              >
                Mark under review
              </Button>
              <Button
                type="button"
                size="sm"
                className="h-11 bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:h-9"
                onClick={() => setStatus(selected.id, "resolved")}
              >
                Resolve
              </Button>
              <Button
                type="button"
                size="sm"
                variant="outline"
                className="h-11 border-rose-400/40 text-rose-200 sm:h-9"
                onClick={() => {
                  if (selected.influencerEmail) blockUserEmail(selected.influencerEmail);
                  window.alert(`Blocked ${selected.influencerEmail || "user"} (demo)`);
                }}
              >
                Block user
              </Button>
            </div>
          </aside>
        ) : null}
      </div>
    </div>
  );
}
