"use client";

import * as React from "react";

import { AUTH_TOKEN_STORAGE_KEY } from "@/lib/constants";

type AuthContextValue = {
  accessToken: string | null;
  setAccessToken: (token: string | null) => void;
};

const AuthContext = React.createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [accessToken, setAccessTokenState] = React.useState<string | null>(null);

  React.useEffect(() => {
    const stored = window.localStorage.getItem(AUTH_TOKEN_STORAGE_KEY);
    setAccessTokenState(stored); // eslint-disable-line react-hooks/set-state-in-effect -- client-only localStorage sync
  }, []);

  const setAccessToken = React.useCallback((token: string | null) => {
    setAccessTokenState(token);
    if (token) {
      window.localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, token);
    } else {
      window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
    }
  }, []);

  const value = React.useMemo(
    () => ({ accessToken, setAccessToken }),
    [accessToken, setAccessToken]
  );

  return (
    <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = React.useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
}
