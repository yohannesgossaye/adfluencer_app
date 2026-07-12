"use client";

import {
  createContext,
  useCallback,
  useContext,
  useMemo,
  type ReactNode,
} from "react";
import { usePathname, useRouter, useSearchParams } from "next/navigation";

import {
  buildMarketingAuthHref,
  parseMarketingAuthRole,
  type MarketingAuthMode,
  type MarketingAuthRole,
} from "@/lib/marketing-auth-url";

type MarketingAuthContextValue = {
  mode: MarketingAuthMode | null;
  signupRole: MarketingAuthRole;
  isOpen: boolean;
  openLogin: () => void;
  openSignup: (role?: MarketingAuthRole) => void;
  closeAuth: () => void;
  switchToLogin: () => void;
  switchToSignup: (role?: MarketingAuthRole) => void;
};

const MarketingAuthContext = createContext<MarketingAuthContextValue | null>(null);

export function MarketingAuthProvider({ children }: { children: ReactNode }) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();

  const authParam = searchParams.get("auth");
  const mode: MarketingAuthMode | null =
    authParam === "login" || authParam === "signup" ? authParam : null;
  const signupRole = parseMarketingAuthRole(searchParams.get("role")) ?? "business";
  const isOpen = mode !== null;

  const pushAuth = useCallback(
    (nextMode: MarketingAuthMode, role?: MarketingAuthRole) => {
      router.push(buildMarketingAuthHref(pathname, nextMode, role), { scroll: false });
    },
    [pathname, router],
  );

  const openLogin = useCallback(() => {
    pushAuth("login");
  }, [pushAuth]);

  const openSignup = useCallback(
    (role: MarketingAuthRole = "business") => {
      pushAuth("signup", role);
    },
    [pushAuth],
  );

  const closeAuth = useCallback(() => {
    router.replace(pathname, { scroll: false });
  }, [pathname, router]);

  const switchToLogin = useCallback(() => {
    pushAuth("login");
  }, [pushAuth]);

  const switchToSignup = useCallback(
    (role: MarketingAuthRole = signupRole) => {
      pushAuth("signup", role);
    },
    [pushAuth, signupRole],
  );

  const value = useMemo(
    () => ({
      mode,
      signupRole,
      isOpen,
      openLogin,
      openSignup,
      closeAuth,
      switchToLogin,
      switchToSignup,
    }),
    [
      mode,
      signupRole,
      isOpen,
      openLogin,
      openSignup,
      closeAuth,
      switchToLogin,
      switchToSignup,
    ],
  );

  return (
    <MarketingAuthContext.Provider value={value}>
      {children}
    </MarketingAuthContext.Provider>
  );
}

export function useMarketingAuth() {
  const context = useContext(MarketingAuthContext);
  if (!context) {
    throw new Error("useMarketingAuth must be used within MarketingAuthProvider");
  }
  return context;
}
