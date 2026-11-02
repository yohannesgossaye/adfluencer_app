"use client";

import { useCallback, useState } from "react";

import { toast } from "sonner";

import { ConfirmActionDialog, type ConfirmActionRequest } from "@/components/common/confirm-action-dialog";
import { ListStatus } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { adminApi, type AdminCategoryItem } from "@/lib/api/admin";
import { usePagedList } from "@/hooks/use-paged-list";
import { adminErrorMessage } from "@/lib/api-errors";

export function AdminCategories() {
  const [newName, setNewName] = useState("");

  // Categories are one short list; the hook only adds loading and error states.
  const fetchAll = useCallback(async () => (await adminApi.listCategories()).data, []);
  const list = usePagedList(fetchAll);
  const rows: AdminCategoryItem[] = list.items;
  const loadCategories = list.reload;

  const [confirm, setConfirm] = useState<ConfirmActionRequest | null>(null);
  const [adding, setAdding] = useState(false);

  const enable = async (c: AdminCategoryItem) => {
    try {
      await adminApi.enableCategory(c.id);
      toast.success(`"${c.name}" is enabled.`);
      await loadCategories();
    } catch (err) {
      toast.error(adminErrorMessage(err, "The category could not be enabled."));
    }
  };

  const askDisable = (c: AdminCategoryItem) =>
    setConfirm({
      title: `Disable "${c.name}"?`,
      description: `It disappears from the niche lists for new profiles and campaigns. The ${c.influencer_count} creator(s) and ${c.campaign_count} campaign(s) already using it keep it. You can enable it again any time.`,
      confirmLabel: "Disable category",
      tone: "danger",
      onConfirm: async () => {
        try {
          await adminApi.disableCategory(c.id);
          toast.success(`"${c.name}" is disabled.`);
          await loadCategories();
        } catch (err) {
          toast.error(adminErrorMessage(err, "The category could not be disabled."));
          throw err;
        }
      },
    });

  const rename = async (c: AdminCategoryItem, input: HTMLInputElement) => {
    const name = input.value.trim();
    if (!name || name === c.name) {
      input.value = c.name;
      return;
    }
    try {
      await adminApi.updateCategory(c.id, { name });
      toast.success(`Renamed to "${name}".`);
      await loadCategories();
    } catch (err) {
      input.value = c.name;
      toast.error(adminErrorMessage(err, "The category could not be renamed."));
    }
  };

  const add = async () => {
    const n = newName.trim();
    if (!n || adding) return;
    setAdding(true);
    try {
      await adminApi.createCategory({ name: n });
      setNewName("");
      toast.success(`"${n}" added.`);
      await loadCategories();
    } catch (err) {
      toast.error(adminErrorMessage(err, "The category could not be added."));
    } finally {
      setAdding(false);
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

      <form
        className="flex flex-wrap items-end gap-2 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
        onSubmit={(e) => {
          e.preventDefault();
          void add();
        }}
      >
        <Input
          aria-label="New category name"
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          placeholder="New category name"
          className="h-10 max-w-xs flex-1 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white"
        />
        <Button
          type="submit"
          disabled={adding || !newName.trim()}
          className="h-10 rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)]"
        >
          Add category
        </Button>
      </form>

      <div className="space-y-3 min-w-0">
        <ListStatus
          state={list.state}
          isEmpty={rows.length === 0}
          emptyText="No categories yet. Add the first one above."
          errorText="Categories could not be loaded."
          onRetry={list.retry}
        >
          {rows.map((c) => (
            <article
              key={c.id}
              className="flex w-full min-w-0 flex-col gap-3 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4 md:flex-row md:items-center md:justify-between"
            >
              <div className="flex flex-1 flex-col gap-2 sm:flex-row sm:items-center sm:gap-4">
                <Input
                  key={`${c.id}-${c.name}`}
                  aria-label={`Rename ${c.name}`}
                  defaultValue={c.name}
                  onBlur={(e) => void rename(c, e.currentTarget)}
                  onKeyDown={(e) => {
                    if (e.key === "Enter") e.currentTarget.blur();
                  }}
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
                  onClick={() => (c.is_active ? askDisable(c) : void enable(c))}
                >
                  {c.is_active ? "Disable" : "Enable"}
                </Button>
              </div>
            </article>
          ))}
        </ListStatus>
      </div>

      <ConfirmActionDialog request={confirm} onClose={() => setConfirm(null)} />
    </div>
  );
}
