"use client";

import { createContext, useCallback, useContext, useMemo, type ReactNode } from "react";
import { usePathname, useRouter } from "next/navigation";

import { buildMarketingAuthHref, type MarketingAuthMode, type MarketingAuthRole } from "@/lib/marketing-auth-url";

/**
 * Opens and closes the login/signup dialog by changing the URL (?auth=…).
 * It deliberately does not read the URL: useSearchParams here would force
 * every marketing page to render only in the browser (empty HTML for
 * Google). MarketingAuthModal reads the URL inside its own Suspense.
 */
type MarketingAuthContextValue = {
  openLogin: () => void;
  openSignup: (role?: MarketingAuthRole) => void;
  closeAuth: () => void;
};

const MarketingAuthContext = createContext<MarketingAuthContextValue | null>(null);

export function MarketingAuthProvider({ children }: { children: ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();

  const pushAuth = useCallback(
    (mode: MarketingAuthMode, role?: MarketingAuthRole) => {
      router.push(buildMarketingAuthHref(pathname, mode, role), { scroll: false });
    },
    [pathname, router],
  );
  const openLogin = useCallback(() => pushAuth("login"), [pushAuth]);
  const openSignup = useCallback((role: MarketingAuthRole = "business") => pushAuth("signup", role), [pushAuth]);
  const closeAuth = useCallback(() => router.replace(pathname, { scroll: false }), [pathname, router]);

  const value = useMemo(() => ({ openLogin, openSignup, closeAuth }), [openLogin, openSignup, closeAuth]);
  return <MarketingAuthContext.Provider value={value}>{children}</MarketingAuthContext.Provider>;
}

export function useMarketingAuth() {
  const context = useContext(MarketingAuthContext);
  if (!context) {
    throw new Error("useMarketingAuth must be used within MarketingAuthProvider");
  }
  return context;
}
