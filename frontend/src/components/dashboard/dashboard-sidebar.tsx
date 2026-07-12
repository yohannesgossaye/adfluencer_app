"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import {
  IconAlertTriangle,
  IconBriefcase,
  IconBuildingStore,
  IconHeadset,
  IconLayoutDashboard,
  IconLogout,
  IconMessageCircle,
  IconReceipt2,
  IconRosetteDiscountCheck,
  IconSparkles,
  IconTags,
  IconUserCircle,
  IconUsers,
} from "@tabler/icons-react";

import { BrandLogo } from "@/components/common/brand-logo";
import { useAuth } from "@/components/providers/auth-provider";
import type { NavItem } from "@/lib/navigation";
import { signOut } from "@/lib/sign-out";
import { cn } from "@/lib/utils";

const iconMap = {
  dashboard: IconLayoutDashboard,
  briefcase: IconBriefcase,
  sparkles: IconSparkles,
  recommendations: IconRosetteDiscountCheck,
  messages: IconMessageCircle,
  contracts: IconReceipt2,
  profile: IconUserCircle,
  users: IconUsers,
  building: IconBuildingStore,
  disputes: IconAlertTriangle,
  categories: IconTags,
  support: IconHeadset,
} as const;

export function DashboardSidebar({
  title,
  items,
  variant = "desktop",
}: {
  title: string;
  items: NavItem[];
  variant?: "desktop" | "drawer";
}) {
  const pathname = usePathname();
  const router = useRouter();

  const { setAccessToken } = useAuth();

  const homeHref =
    items.find((item) => item.icon === "dashboard")?.href ?? items[0]?.href ?? "/";

  const onLogout = () => {
    signOut();
    setAccessToken(null);
    router.push("/");
  };

  return (
    <aside
      className={cn(
        "flex h-full min-h-0 w-full flex-col border-r border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-5 text-white",
        variant === "drawer" && "min-h-[100dvh]",
      )}
    >
      <BrandLogo href={homeHref} tagline={title} className="mb-6 px-1" />

      <nav className="min-h-0 flex-1 space-y-1 overflow-y-auto pr-0.5">
        {items.map((item) => {
          const ItemIcon = iconMap[item.icon] ?? IconLayoutDashboard;
          const active =
            pathname === item.href || pathname.startsWith(`${item.href}/`);
          return (
            <Link
              key={item.href}
              href={item.href}
              className={cn(
                "group relative flex items-center gap-2 rounded-xl px-3 py-2.5 text-sm transition-colors",
                active
                  ? "bg-white/10 text-white"
                  : "text-[#CBD5E1] hover:bg-white/10 hover:text-white",
              )}
            >
              <span
                className={cn(
                  "absolute left-0 top-1/2 h-6 w-1 -translate-y-1/2 rounded-r-full opacity-0 transition-opacity",
                  active
                    ? "bg-[color:var(--adfluence-orange)] opacity-100"
                    : "group-hover:bg-[color:var(--adfluence-orange)]/50 group-hover:opacity-60",
                )}
              />
              <ItemIcon className="size-4 shrink-0" />
              <span className="truncate">{item.label}</span>
            </Link>
          );
        })}
      </nav>

      <div className="mt-4 shrink-0 border-t border-white/10 pt-4">
        <button
          type="button"
          onClick={onLogout}
          className="flex w-full items-center gap-2 rounded-xl px-3 py-2.5 text-left text-sm text-[#CBD5E1] transition-colors hover:bg-white/10 hover:text-white"
        >
          <IconLogout
            className="size-4 shrink-0 text-[color:var(--adfluence-orange)]/90"
            aria-hidden
          />
          <span className="truncate">Logout</span>
        </button>
      </div>
    </aside>
  );
}
