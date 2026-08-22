"use client";

import { useEffect, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  adminApi,
  type AdminUserListItem,
  type AdminUserStats,
} from "@/lib/api/admin";
import { cn } from "@/lib/utils";

function statusStyle(s: string) {
  const lower = s.toLowerCase();
  switch (lower) {
    case "active":
      return "border-emerald-400/30 bg-emerald-500/10 text-emerald-200";
    case "pending":
      return "border-amber-400/30 bg-amber-500/10 text-amber-200";
    case "blocked":
    case "suspended":
    case "rejected":
      return "border-rose-400/30 bg-rose-500/10 text-rose-200";
    default:
      return "border-white/10 bg-white/5 text-[#94A3B8]";
  }
}

export function AdminUsers() {
  const [users, setUsers] = useState<AdminUserListItem[]>([]);
  const [stats, setStats] = useState<AdminUserStats>({
    total: 0,
    active: 0,
    blocked: 0,
    pending: 0,
  });
  const [roleFilter, setRoleFilter] = useState<"all" | "business" | "influencer" | "admin">("all");
  const [statusFilter, setStatusFilter] = useState<"all" | "active" | "pending" | "blocked" | "removed">("all");
  const [search, setSearch] = useState("");
  const [isLoading, setIsLoading] = useState(true);

  const loadUsers = async () => {
    try {
      const res = await adminApi.listUsers({
        role: roleFilter === "all" ? undefined : roleFilter,
        status: statusFilter === "all" ? undefined : statusFilter,
        search: search.trim() || undefined,
      });
      if (res.data?.items) {
        setUsers(res.data.items);
      }
      if (res.data?.stats) {
        setStats(res.data.stats);
      }
    } catch {
      // fallback
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadUsers();
  }, [roleFilter, statusFilter, search]);

  const block = async (u: AdminUserListItem) => {
    try {
      await adminApi.blockUser(u.user_id, { reason: "Blocked by administrator" });
      await loadUsers();
    } catch {
      // fallback
    }
  };

  const unblock = async (u: AdminUserListItem) => {
    try {
      await adminApi.unblockUser(u.user_id);
      await loadUsers();
    } catch {
      // fallback
    }
  };

  const remove = async (u: AdminUserListItem) => {
    try {
      await adminApi.removeUser(u.user_id, { reason: "Purged by administrator" });
      await loadUsers();
    } catch {
      // fallback
    }
  };

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
          { label: "Total users", value: String(stats.total) },
          { label: "Active", value: String(stats.active) },
          { label: "Blocked", value: String(stats.blocked) },
          { label: "Pending", value: String(stats.pending) },
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
                ? "border-transparent bg-[color:var(--adfluence-orange)] text-white"
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
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        placeholder="Search name or email"
        className="h-11 max-w-md rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white"
      />

      <div className="space-y-3 min-w-0">
        {users.length === 0 && !isLoading ? (
          <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center">
            <p className="text-white">No users found.</p>
          </div>
        ) : (
          users.map((u) => (
            <article
              key={u.user_id}
              className="w-full min-w-0 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
            >
              <div className="flex flex-wrap items-start justify-between gap-2">
                <div>
                  <p className="font-semibold text-white">{u.name}</p>
                  <p className="text-sm text-[#CBD5E1]">{u.email}</p>
                  <p className="mt-1 text-xs capitalize text-[#94A3B8]">
                    {u.role} · Joined {new Date(u.joined_at).toLocaleDateString()}
                    {u.last_active_at ? ` · Last active ${new Date(u.last_active_at).toLocaleDateString()}` : ""}
                  </p>
                  <p className="mt-1 text-xs text-[#94A3B8]">
                    Campaigns {u.campaigns_count} · Contracts {u.contracts_count}
                  </p>
                </div>
                <span
                  className={cn(
                    "rounded-full border px-2.5 py-0.5 text-xs capitalize",
                    statusStyle(u.status),
                  )}
                >
                  {u.status}
                </span>
              </div>
              {u.status === "blocked" || u.status === "suspended" ? (
                <p className="mt-2 text-xs text-rose-200">
                  Blocked users cannot log in.
                </p>
              ) : null}
              {u.status === "removed" ? (
                <p className="mt-2 text-xs text-amber-300">
                  Removed account. Click &quot;Restore &amp; Activate&quot; to reverse removal and allow login.
                </p>
              ) : null}
              <div className="mt-3 flex flex-wrap gap-2">
                {u.status === "removed" || u.status === "blocked" || u.status === "suspended" ? (
                  <Button
                    size="sm"
                    className="h-11 sm:h-8 bg-emerald-500/25 border border-emerald-500/40 text-emerald-100 hover:bg-emerald-500/35"
                    onClick={() => unblock(u)}
                  >
                    {u.status === "removed" ? "Restore & Activate" : "Unblock & Activate"}
                  </Button>
                ) : null}
                {u.status !== "blocked" ? (
                  <Button
                    size="sm"
                    className="h-11 sm:h-8 bg-rose-500/20 text-rose-100 hover:bg-rose-500/30"
                    disabled={u.status === "removed"}
                    onClick={() => block(u)}
                  >
                    Block
                  </Button>
                ) : null}
                {u.status !== "removed" ? (
                  <Button
                    size="sm"
                    variant="outline"
                    className="h-11 sm:h-8 border-white/15 text-white hover:bg-white/10"
                    onClick={() => remove(u)}
                  >
                    Remove
                  </Button>
                ) : null}
              </div>
            </article>
          ))
        )}
      </div>
    </div>
  );
}
