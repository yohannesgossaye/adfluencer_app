"use client";

import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { adminApi, type AdminCategoryItem } from "@/lib/api/admin";

export function AdminCategories() {
  const [rows, setRows] = useState<AdminCategoryItem[]>([]);
  const [newName, setNewName] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  const loadCategories = async () => {
    try {
      const res = await adminApi.listCategories();
      if (res.data?.items) {
        setRows(res.data.items);
      }
    } catch {
      // fallback
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadCategories();
  }, []);

  const toggle = async (id: string, currentlyActive: boolean) => {
    try {
      if (currentlyActive) {
        await adminApi.disableCategory(id);
      } else {
        await adminApi.enableCategory(id);
      }
      await loadCategories();
    } catch {
      // fallback
    }
  };

  const rename = async (id: string, name: string) => {
    try {
      await adminApi.updateCategory(id, { name });
      await loadCategories();
    } catch {
      // fallback
    }
  };

  const add = async () => {
    const n = newName.trim();
    if (!n) return;
    try {
      await adminApi.createCategory({ name: n });
      setNewName("");
      await loadCategories();
    } catch {
      // fallback
    }
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
        {rows.length === 0 && !isLoading ? (
          <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center">
            <p className="text-white">No categories found.</p>
          </div>
        ) : (
          rows.map((c) => (
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
                  Influencers: {c.influencer_count} · Campaigns: {c.campaign_count}
                </p>
              </div>
              <div className="flex items-center gap-2">
                <span className="text-xs text-[#94A3B8]">{c.is_active ? "Enabled" : "Disabled"}</span>
                <Button
                  type="button"
                  size="sm"
                  variant="outline"
                  className="h-11 border-white/15 text-white sm:h-8"
                  onClick={() => toggle(c.id, c.is_active)}
                >
                  {c.is_active ? "Disable" : "Enable"}
                </Button>
              </div>
            </article>
          ))
        )}
      </div>
    </div>
  );
}
