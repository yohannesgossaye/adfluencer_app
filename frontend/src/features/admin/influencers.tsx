"use client";

import Image from "next/image";
import { useEffect, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  loadApplications,
  saveApplications,
  type AdminApplicationStatus,
  type AdminInfluencerApplication,
} from "@/data/admin-influencer-applications";
import { cn } from "@/lib/utils";

const tabs: Array<{ id: "all" | AdminApplicationStatus; label: string }> = [
  { id: "all", label: "All" },
  { id: "pending", label: "Pending" },
  { id: "approved", label: "Approved" },
  { id: "rejected", label: "Rejected" },
];

function statusBadge(s: AdminApplicationStatus) {
  if (s === "pending") return "border-amber-400/30 bg-amber-500/10 text-amber-200";
  if (s === "approved") return "border-emerald-400/30 bg-emerald-500/10 text-emerald-200";
  return "border-rose-400/30 bg-rose-500/10 text-rose-200";
}

export function AdminInfluencers() {
  const [rows, setRows] = useState<AdminInfluencerApplication[]>([]);
  const [tab, setTab] = useState<(typeof tabs)[number]["id"]>("all");
  const [search, setSearch] = useState("");
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [noteDraft, setNoteDraft] = useState("");

  useEffect(() => {
    const list = loadApplications();
    setRows(list);
    if (list[0]) setSelectedId((prev) => prev ?? list[0].id);
  }, []);

  const selected = useMemo(
    () => rows.find((r) => r.id === selectedId) ?? null,
    [rows, selectedId],
  );

  useEffect(() => {
    const s = rows.find((r) => r.id === selectedId);
    if (s) setNoteDraft(s.adminNotes);
  }, [selectedId, rows]);

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    return rows.filter((r) => {
      const tabOk = tab === "all" || r.status === tab;
      const searchOk =
        !q ||
        r.name.toLowerCase().includes(q) ||
        r.niche.toLowerCase().includes(q) ||
        r.platform.toLowerCase().includes(q);
      return tabOk && searchOk;
    });
  }, [rows, tab, search]);

  const persist = (next: AdminInfluencerApplication[]) => {
    saveApplications(next);
    setRows(next);
  };

  const setStatus = (id: string, status: AdminApplicationStatus) => {
    persist(rows.map((r) => (r.id === id ? { ...r, status } : r)));
  };

  const saveNotes = () => {
    if (!selected) return;
    persist(
      rows.map((r) =>
        r.id === selected.id ? { ...r, adminNotes: noteDraft } : r,
      ),
    );
  };

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
              "rounded-full border px-4 py-2 text-sm font-medium transition-colors",
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
        <div className="space-y-3 overflow-x-auto">
          {filtered.map((r) => (
            <button
              key={r.id}
              type="button"
              onClick={() => setSelectedId(r.id)}
              className={cn(
                "flex w-full min-w-[280px] gap-4 rounded-2xl border p-4 text-left transition-colors",
                selectedId === r.id
                  ? "border-[color:var(--adfluence-orange)]/35 bg-[color:var(--adfluence-orange)]/10"
                  : "border-white/10 bg-[rgba(255,255,255,0.08)] hover:border-white/20",
              )}
            >
              <div className="relative h-16 w-16 shrink-0 overflow-hidden rounded-full border border-white/10">
                <Image src={r.imageUrl} alt={r.name} fill className="object-cover" sizes="64px" />
              </div>
              <div className="min-w-0 flex-1">
                <div className="flex flex-wrap items-center gap-2">
                  <p className="font-semibold text-white">{r.name}</p>
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
                  {r.niche} · {r.platform} · {r.followers.toLocaleString()} followers ·{" "}
                  {r.engagementRate}% engagement
                </p>
                <p className="mt-1 text-xs text-[#94A3B8]">Submitted {r.submittedAt}</p>
                <p className="mt-1 text-xs text-[#CBD5E1]">{r.completeness}</p>
                <div className="mt-2 flex flex-wrap gap-2">
                  <Button
                    type="button"
                    size="sm"
                    className="h-8 bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
                    disabled={r.status !== "pending"}
                    onClick={(e) => {
                      e.stopPropagation();
                      setStatus(r.id, "approved");
                    }}
                  >
                    Approve
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    className="h-8 border-white/15 text-white hover:bg-white/10"
                    disabled={r.status !== "pending"}
                    onClick={(e) => {
                      e.stopPropagation();
                      setStatus(r.id, "rejected");
                    }}
                  >
                    Reject
                  </Button>
                </div>
              </div>
            </button>
          ))}
        </div>

        {selected ? (
          <aside className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-5 xl:sticky xl:top-4 xl:self-start">
            <h2 className="text-lg font-semibold text-white">Details — {selected.name}</h2>
            <p className="mt-2 text-sm text-[#CBD5E1]">{selected.bio}</p>
            <div className="mt-3 space-y-1 text-xs text-[#94A3B8]">
              <p>TikTok: {selected.tiktok || "—"}</p>
              <p>Instagram: {selected.instagram || "—"}</p>
              <p>YouTube: {selected.youtube || "—"}</p>
            </div>
            <p className="mt-3 text-sm text-[#CBD5E1]">
              <span className="text-white">Audience:</span> {selected.audienceNotes}
            </p>
            <p className="mt-2 text-sm text-[#CBD5E1]">
              <span className="text-white">Pricing:</span> {selected.pricing}
            </p>
            <p className="mt-3 text-xs font-medium uppercase text-[#94A3B8]">Sample content</p>
            <div className="relative mt-2 aspect-video overflow-hidden rounded-xl border border-white/10">
              <Image src={selected.imageUrl} alt={selected.name} fill className="object-cover" sizes="360px" />
            </div>
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
              className="mt-2 h-9 w-full border-white/15 text-white hover:bg-white/10"
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
