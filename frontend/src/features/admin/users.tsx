"use client";

import { useEffect, useMemo, useState } from "react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  loadAdminUsers,
  saveAdminUsers,
  type AdminUserRecord,
  type AdminUserRole,
  type AdminUserStatus,
} from "@/data/admin-users-mock";
import { blockUserEmail, unblockUserEmail } from "@/lib/blocked-users";
import { cn } from "@/lib/utils";

function statusStyle(s: AdminUserStatus) {
  switch (s) {
    case "active":
      return "border-emerald-400/30 bg-emerald-500/10 text-emerald-200";
    case "pending":
      return "border-amber-400/30 bg-amber-500/10 text-amber-200";
    case "blocked":
      return "border-rose-400/30 bg-rose-500/10 text-rose-200";
    default:
      return "border-white/10 bg-white/5 text-[#94A3B8]";
  }
}

export function AdminUsers() {
  const [users, setUsers] = useState<AdminUserRecord[]>([]);
  const [roleFilter, setRoleFilter] = useState<"all" | AdminUserRole>("all");
  const [statusFilter, setStatusFilter] = useState<"all" | AdminUserStatus>("all");
  const [search, setSearch] = useState("");

  useEffect(() => {
    setUsers(loadAdminUsers());
  }, []);

  const filtered = useMemo(() => {
    const q = search.trim().toLowerCase();
    return users.filter((u) => {
      const roleOk = roleFilter === "all" || u.role === roleFilter;
      const statusOk = statusFilter === "all" || u.status === statusFilter;
      const searchOk =
        !q ||
        u.name.toLowerCase().includes(q) ||
        u.email.toLowerCase().includes(q);
      return roleOk && statusOk && searchOk;
    });
  }, [users, roleFilter, statusFilter, search]);

  const persist = (next: AdminUserRecord[]) => {
    saveAdminUsers(next);
    setUsers(next);
  };

  const block = (u: AdminUserRecord) => {
    blockUserEmail(u.email);
    persist(users.map((x) => (x.id === u.id ? { ...x, status: "blocked" as const } : x)));
  };

  const unblock = (u: AdminUserRecord) => {
    unblockUserEmail(u.email);
    persist(users.map((x) => (x.id === u.id ? { ...x, status: "active" as const } : x)));
  };

  const remove = (u: AdminUserRecord) => {
    persist(users.map((x) => (x.id === u.id ? { ...x, status: "removed" as const } : x)));
  };

  const stats = useMemo(() => {
    return {
      total: users.length,
      active: users.filter((u) => u.status === "active").length,
      blocked: users.filter((u) => u.status === "blocked").length,
      pending: users.filter((u) => u.status === "pending").length,
    };
  }, [users]);

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
        {filtered.map((u) => (
          <article
            key={u.id}
            className="w-full min-w-0 rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4"
          >
            <div className="flex flex-wrap items-start justify-between gap-2">
              <div>
                <p className="font-semibold text-white">{u.name}</p>
                <p className="text-sm text-[#CBD5E1]">{u.email}</p>
                <p className="mt-1 text-xs capitalize text-[#94A3B8]">
                  {u.role} · Joined {u.joinedAt} · Last active {u.lastActive}
                </p>
                <p className="mt-1 text-xs text-[#94A3B8]">
                  Campaigns {u.campaignsCount} · Contracts {u.contractsCount}
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
            {u.status === "blocked" ? (
              <p className="mt-2 text-xs text-rose-200">
                Blocked users cannot log in.
              </p>
            ) : null}
            <div className="mt-3 flex flex-wrap gap-2">
              <Button size="sm" className="h-11 sm:h-8 bg-white/10 text-white hover:bg-white/15">
                View
              </Button>
              <Button
                size="sm"
                className="h-11 sm:h-8 bg-rose-500/20 text-rose-100 hover:bg-rose-500/30"
                disabled={u.status === "blocked" || u.status === "removed"}
                onClick={() => block(u)}
              >
                Block
              </Button>
              <Button
                size="sm"
                variant="outline"
                className="h-11 sm:h-8 border-white/15 text-white hover:bg-white/10"
                disabled={u.status !== "blocked"}
                onClick={() => unblock(u)}
              >
                Unblock
              </Button>
              <Button
                size="sm"
                variant="outline"
                className="h-11 sm:h-8 border-white/15 text-white hover:bg-white/10"
                disabled={u.status === "removed"}
                onClick={() => remove(u)}
              >
                Remove
              </Button>
            </div>
          </article>
        ))}
      </div>
    </div>
  );
}
