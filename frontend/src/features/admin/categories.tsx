"use client";

import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  initialAdminCategories,
  loadCategories,
  saveCategories,
  type AdminCategory,
} from "@/data/admin-categories-mock";

export function AdminCategories() {
  const [rows, setRows] = useState<AdminCategory[]>(initialAdminCategories);
  const [newName, setNewName] = useState("");

  useEffect(() => {
    setRows(loadCategories());
  }, []);

  const persist = (next: AdminCategory[]) => {
    saveCategories(next);
    setRows(next);
  };

  const toggle = (id: string) => {
    persist(rows.map((r) => (r.id === id ? { ...r, enabled: !r.enabled } : r)));
  };

  const rename = (id: string, name: string) => {
    persist(rows.map((r) => (r.id === id ? { ...r, name } : r)));
  };

  const add = () => {
    const n = newName.trim();
    if (!n) return;
    persist([
      ...rows,
      {
        id: `cat-${Date.now()}`,
        name: n,
        enabled: true,
        influencerCount: 0,
        campaignCount: 0,
      },
    ]);
    setNewName("");
  };

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Categories & niches</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Manage creator niches and campaign categories for the marketplace.
        </p>
      </div>

      <div className="flex flex-wrap items-end gap-2 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
        <Input
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          placeholder="New category name"
          className="h-10 max-w-xs flex-1 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white"
        />
        <Button
          type="button"
          className="h-10 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)]"
          onClick={add}
        >
          Add category
        </Button>
      </div>

      <div className="space-y-3 min-w-0">
        {rows.map((c) => (
          <article
            key={c.id}
            className="flex w-full min-w-0 flex-col gap-3 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:flex-row md:items-center md:justify-between"
          >
            <div className="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center sm:gap-4">
              <Input
                key={`${c.id}-${c.name}`}
                defaultValue={c.name}
                onBlur={(e) => rename(c.id, e.target.value.trim() || c.name)}
                className="h-10 max-w-[200px] rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white"
              />
              <p className="text-xs text-[#CBD5E1]">
                Influencers: {c.influencerCount} · Campaigns: {c.campaignCount}
              </p>
            </div>
            <div className="flex items-center gap-2">
              <span className="text-xs text-[#94A3B8]">{c.enabled ? "Enabled" : "Disabled"}</span>
              <Button
                type="button"
                size="sm"
                variant="outline"
                className="h-11 border-white/15 text-white sm:h-8"
                onClick={() => toggle(c.id)}
              >
                {c.enabled ? "Disable" : "Enable"}
              </Button>
            </div>
          </article>
        ))}
      </div>
    </div>
  );
}
