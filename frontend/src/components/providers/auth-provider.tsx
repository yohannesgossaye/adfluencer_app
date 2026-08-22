"use client";

import * as React from "react";

import {
  AUTH_TOKEN_STORAGE_KEY,
  REFRESH_TOKEN_STORAGE_KEY,
  AUTH_USER_STORAGE_KEY,
} from "@/lib/constants";
import {
  authApi,
  type AuthUser,
  type AuthProfile,
  type AuthTokens,
} from "@/lib/api/auth";

type AuthContextValue = {
  accessToken: string | null;
  user: AuthUser | null;
  profile: AuthProfile | null;
  isLoading: boolean;
  setAccessToken: (token: string | null) => void;
  setSession: (
    tokens: AuthTokens | null,
    user?: AuthUser | null,
    profile?: AuthProfile | null
  ) => void;
  logout: () => Promise<void>;
  refreshMe: () => Promise<void>;
};

const AuthContext = React.createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [accessToken, setAccessTokenState] = React.useState<string | null>(null);
  const [user, setUser] = React.useState<AuthUser | null>(null);
  const [profile, setProfile] = React.useState<AuthProfile | null>(null);
  const [isLoading, setIsLoading] = React.useState<boolean>(true);

  const refreshMe = React.useCallback(async () => {
    const token = window.localStorage.getItem(AUTH_TOKEN_STORAGE_KEY);
    if (!token) {
      setUser(null);
      setProfile(null);
      setIsLoading(false);
      return;
    }

    try {
      const res = await authApi.getMe();
      if (res.data) {
        setUser(res.data.user);
        setProfile(res.data.profile);
        window.localStorage.setItem(
          AUTH_USER_STORAGE_KEY,
          JSON.stringify({ user: res.data.user, profile: res.data.profile })
        );
      }
    } catch {
      // 401 interceptor handles token refresh or cleanup
    } finally {
      setIsLoading(false);
    }
  }, []);

  React.useEffect(() => {
    const storedToken = window.localStorage.getItem(AUTH_TOKEN_STORAGE_KEY);
    setAccessTokenState(storedToken);

    const storedUser = window.localStorage.getItem(AUTH_USER_STORAGE_KEY);
    if (storedUser) {
      try {
        const parsed = JSON.parse(storedUser);
        if (parsed.user) setUser(parsed.user);
        if (parsed.profile) setProfile(parsed.profile);
      } catch {
        // Ignore JSON parse errors
      }
    }

    if (storedToken) {
      refreshMe();
    } else {
      setIsLoading(false);
    }
  }, [refreshMe]);

  const setAccessToken = React.useCallback((token: string | null) => {
    setAccessTokenState(token);
    if (token) {
      window.localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, token);
    } else {
      window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
    }
  }, []);

  const setSession = React.useCallback(
    (
      tokens: AuthTokens | null,
      newUser?: AuthUser | null,
      newProfile?: AuthProfile | null
    ) => {
      if (tokens) {
        setAccessTokenState(tokens.access_token);
        window.localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, tokens.access_token);
        if (tokens.refresh_token) {
          window.localStorage.setItem(
            REFRESH_TOKEN_STORAGE_KEY,
            tokens.refresh_token
          );
        }
      } else {
        setAccessTokenState(null);
        window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
        window.localStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
      }

      if (newUser !== undefined) {
        setUser(newUser);
      }
      if (newProfile !== undefined) {
        setProfile(newProfile);
      }

      if (newUser && newProfile) {
        window.localStorage.setItem(
          AUTH_USER_STORAGE_KEY,
          JSON.stringify({ user: newUser, profile: newProfile })
        );
      } else if (!tokens && !newUser) {
        window.localStorage.removeItem(AUTH_USER_STORAGE_KEY);
      }
    },
    []
  );

  const logout = React.useCallback(async () => {
    const refreshToken = window.localStorage.getItem(REFRESH_TOKEN_STORAGE_KEY) || undefined;
    try {
      await authApi.logout(refreshToken);
    } catch {
      // Continue client cleanup even if network fails
    } finally {
      setAccessTokenState(null);
      setUser(null);
      setProfile(null);
      window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
      window.localStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
      window.localStorage.removeItem(AUTH_USER_STORAGE_KEY);
    }
  }, []);

  const value = React.useMemo(
    () => ({
      accessToken,
      user,
      profile,
      isLoading,
      setAccessToken,
      setSession,
      logout,
      refreshMe,
    }),
    [accessToken, user, profile, isLoading, setAccessToken, setSession, logout, refreshMe]
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = React.useContext(AuthContext);
  if (!ctx) {
    throw new Error("useAuth must be used within AuthProvider");
  }
  return ctx;
}
