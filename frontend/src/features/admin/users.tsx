"use client";

import { useCallback, useEffect, useState } from "react";

import { toast } from "sonner";

import { ConfirmActionDialog, type ConfirmActionRequest } from "@/components/common/confirm-action-dialog";
import { ListStatus, LoadMore } from "@/components/common/list-status";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  adminApi,
  type AdminUserListItem,
  type AdminUserStats,
} from "@/lib/api/admin";
import { usePagedList } from "@/hooks/use-paged-list";
import { adminErrorMessage } from "@/lib/api-errors";
import { cn } from "@/lib/utils";
import { formatDate } from "@/lib/contracts";
import { StatusBadge } from "@/components/common/status-badge";
import { accountStatus } from "@/lib/status";

/** Accounts that can't log in and can be restored. */
function isLocked(status: string) {
  return status === "blocked" || status === "removed" || status === "suspended";
}


export function AdminUsers() {
  const [roleFilter, setRoleFilter] = useState<"all" | "business" | "influencer" | "admin">("all");
  const [statusFilter, setStatusFilter] = useState<"all" | "active" | "pending" | "blocked" | "removed">("all");
  const [search, setSearch] = useState("");
  const [query, setQuery] = useState("");

  useEffect(() => {
    const t = setTimeout(() => setQuery(search.trim()), 300);
    return () => clearTimeout(t);
  }, [search]);

  const fetchPage = useCallback(
    async (page: number) =>
      (
        await adminApi.listUsers({
          role: roleFilter === "all" ? undefined : roleFilter,
          status: statusFilter === "all" ? undefined : statusFilter,
          search: query || undefined,
          page,
        })
      ).data,
    [roleFilter, statusFilter, query],
  );
  const list = usePagedList(fetchPage);
  const users = list.items;
  const stats: AdminUserStats = list.last?.stats ?? { total: 0, active: 0, blocked: 0, pending: 0 };
  const loadUsers = list.reload;

  const [confirm, setConfirm] = useState<ConfirmActionRequest | null>(null);

  const run = async (action: () => Promise<unknown>, success: string, failure: string) => {
    try {
      await action();
      toast.success(success);
      await loadUsers();
    } catch (err) {
      toast.error(adminErrorMessage(err, failure));
      throw err;
    }
  };

  const askBlock = (u: AdminUserListItem) =>
    setConfirm({
      title: `Block ${u.name}?`,
      description: "They are signed out everywhere and can't log in until you unblock them. Their data is kept.",
      confirmLabel: "Block account",
      tone: "danger",
      reasonLabel: "Reason (kept on the account)",
      reasonPlaceholder: "e.g. Repeated spam messages to creators",
      onConfirm: (reason) =>
        run(() => adminApi.blockUser(u.user_id, { reason }), `${u.name} is blocked.`, "The account could not be blocked."),
    });

  const askRemove = (u: AdminUserListItem) =>
    setConfirm({
      title: `Remove ${u.name}?`,
      description:
        "The account is closed and signed out everywhere. Nothing is deleted, so you can restore it later from the Removed filter.",
      confirmLabel: "Remove account",
      tone: "danger",
      reasonLabel: "Reason (kept on the account)",
      reasonPlaceholder: "e.g. Duplicate account",
      onConfirm: (reason) =>
        run(() => adminApi.removeUser(u.user_id, { reason }), `${u.name} is removed.`, "The account could not be removed."),
    });

  const askRestore = (u: AdminUserListItem) =>
    setConfirm({
      title: `Restore ${u.name}?`,
      description:
        "The account goes back to the status it had before (for example, a creator awaiting approval stays awaiting it) and can log in again.",
      confirmLabel: u.status === "removed" ? "Restore account" : "Unblock account",
      onConfirm: () =>
        run(() => adminApi.unblockUser(u.user_id), `${u.name} is restored.`, "The account could not be restored."),
    });

  return (
    <div className="space-y-6 text-white">
      <div>
        <h1 className="text-2xl font-semibold text-white">Users</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Manage businesses, influencers, and admin accounts. Blocked users cannot log in.
        </p>
      </div>

      <section className="grid gap-3 sm:grid-cols-2 lg:grid-cols-4">
        {[
          { label: "Total users", value: list.state === "ready" ? String(stats.total) : "—" },
          { label: "Active", value: list.state === "ready" ? String(stats.active) : "—" },
          { label: "Blocked", value: list.state === "ready" ? String(stats.blocked) : "—" },
          { label: "Pending", value: list.state === "ready" ? String(stats.pending) : "—" },
        ].map((s) => (
          <article
            key={s.label}
            className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
          >
            <p className="text-xs text-[#CBD5E1]">{s.label}</p>
            <p className="mt-2 text-xl font-bold text-[color:var(--adfluence-orange)]">
              {s.value}
            </p>
          </article>
        ))}
      </section>

      <div className="flex flex-wrap gap-2">
        {(["all", "business", "influencer", "admin"] as const).map((r) => (
          <button
            key={r}
            type="button"
            onClick={() => setRoleFilter(r)}
            className={cn(
              "inline-flex min-h-11 items-center rounded-full border px-4 py-2 text-sm font-medium capitalize transition-colors sm:min-h-0 sm:px-3 sm:py-1.5 sm:text-xs",
              roleFilter === r
                ? "border-transparent bg-[color:var(--adfluence-action)] text-white"
                : "border-white/10 bg-white/5 text-[#CBD5E1] hover:text-white",
            )}
          >
            {r}
          </button>
        ))}
        {(["all", "active", "pending", "blocked", "removed"] as const).map((s) => (
          <button
            key={s}
            type="button"
            onClick={() => setStatusFilter(s)}
            className={cn(
              "inline-flex min-h-11 items-center rounded-full border px-4 py-2 text-sm font-medium capitalize transition-colors sm:min-h-0 sm:px-3 sm:py-1.5 sm:text-xs",
              statusFilter === s
                ? "border-white/20 bg-white/15 text-white"
                : "border-white/10 bg-white/5 text-[#CBD5E1] hover:text-white",
            )}
          >
            {s}
          </button>
        ))}
      </div>

      <Input
        aria-label="Search users"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Search name or email"
        className="h-11 max-w-md rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white"
      />

      <div className="space-y-3 min-w-0">
        <ListStatus
          state={list.state}
          isEmpty={users.length === 0}
          emptyText="No users match these filters."
          errorText="Users could not be loaded."
          onRetry={list.retry}
        >
          {users.map((u) => (
            <article
              key={u.user_id}
              className="w-full min-w-0 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
            >
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                  <p className="font-semibold text-white">{u.name}</p>
                  <p className="text-sm text-[#CBD5E1]">{u.email}</p>
                  <p className="mt-1 text-xs capitalize text-[#94A3B8]">
                    {u.role} · Joined {formatDate(u.joined_at)}
                    {u.last_active_at ? ` · Last active ${formatDate(u.last_active_at)}` : ""}
                  </p>
                  <p className="mt-1 text-xs text-[#94A3B8]">
                    Campaigns {u.campaigns_count} · Contracts {u.contracts_count}
                  </p>
                </div>
                <StatusBadge status={accountStatus(u.status)} />
              </div>
              {u.status_reason && isLocked(u.status) ? (
                <p className="mt-2 text-xs text-[#CBD5E1]">Reason: {u.status_reason}</p>
              ) : null}
              {u.role === "admin" ? (
                <p className="mt-3 text-xs text-[#94A3B8]">Admin accounts can&apos;t be blocked or removed here.</p>
              ) : (
                <div className="mt-3 flex flex-wrap gap-2">
                  {isLocked(u.status) ? (
                    <Button
                      size="sm"
                      className="h-11 border border-white/15 bg-transparent text-white hover:bg-white/10 sm:h-8"
                      onClick={() => askRestore(u)}
                    >
                      {u.status === "removed" ? "Restore account" : "Unblock"}
                    </Button>
                  ) : (
                    <Button
                      size="sm"
                      className="h-11 border border-rose-400/50 bg-transparent text-rose-200 hover:bg-rose-500/10 sm:h-8"
                      onClick={() => askBlock(u)}
                    >
                      Block…
                    </Button>
                  )}
                  {u.status !== "removed" ? (
                    <Button
                      size="sm"
                      variant="outline"
                      className="h-11 border-white/15 text-white hover:bg-white/10 sm:h-8"
                      onClick={() => askRemove(u)}
                    >
                      Remove…
                    </Button>
                  ) : null}
                </div>
              )}
            </article>
          ))}
          <LoadMore
            shown={users.length}
            total={list.total}
            hasMore={list.hasMore}
            loading={list.loadingMore}
            onLoadMore={list.loadMore}
          />
        </ListStatus>
      </div>

      <ConfirmActionDialog request={confirm} onClose={() => setConfirm(null)} />
    </div>
  );
}
