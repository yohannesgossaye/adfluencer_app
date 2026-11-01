"use client";

import { useEffect, type ReactNode } from "react";
import { usePathname, useRouter } from "next/navigation";
import { Loader2 } from "lucide-react";

import { useAuth } from "@/components/providers/auth-provider";
import { homeFor } from "@/lib/role-home";

/**
 * Shows a portal only to users with the matching role. Signed-out visitors go
 * to login (and come back afterwards); users of another role go to their own
 * portal. The API still enforces access; this keeps people out of pages that
 * would only show errors.
 */
export function RoleGuard({ role, children }: { role: "business" | "influencer" | "admin"; children: ReactNode }) {
  const { user, accessToken, isLoading } = useAuth();
  const router = useRouter();
  const pathname = usePathname();

  const signedOut = !isLoading && !accessToken;
  const wrongRole = !isLoading && !!user && user.role !== role;

  useEffect(() => {
    if (signedOut) {
      const login = role === "admin" ? "/admin/login" : "/auth/login";
      router.replace(`${login}?next=${encodeURIComponent(pathname)}`);
    } else if (wrongRole) {
      router.replace(homeFor(user));
    }
  }, [signedOut, wrongRole, role, router, pathname, user]);

  if (isLoading || signedOut || wrongRole || !user) {
    return (
      <div className="flex min-h-svh items-center justify-center bg-[color:var(--adfluence-bg)] text-sm text-[#94A3B8]">
        <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> Loading…
      </div>
    );
  }
  return <>{children}</>;
}
