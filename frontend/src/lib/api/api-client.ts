import axios, { type AxiosError, type InternalAxiosRequestConfig } from "axios";

import { AUTH_TOKEN_STORAGE_KEY, REFRESH_TOKEN_STORAGE_KEY, AUTH_USER_STORAGE_KEY } from "@/lib/constants";
import { getPublicEnv } from "@/lib/env";

const { apiBaseUrl } = getPublicEnv();
const baseURL = apiBaseUrl || "http://localhost:8080";

export const apiClient = axios.create({
  baseURL,
  headers: {
    Accept: "application/json",
    "Content-Type": "application/json",
  },
  paramsSerializer: {
    indexes: null,
  },
});

apiClient.interceptors.request.use((config) => {
  if (typeof window === "undefined") return config;
  const token = window.localStorage.getItem(AUTH_TOKEN_STORAGE_KEY);
  if (token && !config.headers.Authorization) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

type RetriableRequest = InternalAxiosRequestConfig & { _retry?: boolean };

const AUTH_ENDPOINTS = ["/api/v1/auth/login", "/api/v1/auth/signup", "/api/v1/auth/refresh"];
const REFRESH_LOCK_NAME = "adfluence-token-refresh";
const REFRESH_TIMEOUT_MS = 15_000;

/** The server rejected the refresh token: the session is really over. */
class SessionExpiredError extends Error {}

function isAuthEndpoint(url?: string) {
  return !!url && AUTH_ENDPOINTS.some((path) => url.includes(path));
}

function bearerToken(config: InternalAxiosRequestConfig): string | null {
  const header = config.headers?.Authorization;
  return typeof header === "string" && header.startsWith("Bearer ") ? header.slice(7) : null;
}

function clearSessionAndRedirect() {
  window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
  window.localStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
  window.localStorage.removeItem(AUTH_USER_STORAGE_KEY);
  const path = `${window.location.pathname}${window.location.search}`;
  if (path.startsWith("/auth") || path.startsWith("/admin/login")) return;
  // Admins sign in on their own page; everyone else on the main login.
  const login = path.startsWith("/admin") ? "/admin/login" : "/auth/login";
  window.location.assign(`${login}?next=${encodeURIComponent(path)}`);
}

async function requestNewTokens(): Promise<string> {
  const refreshToken = window.localStorage.getItem(REFRESH_TOKEN_STORAGE_KEY);
  if (!refreshToken) throw new SessionExpiredError("no refresh token");

  try {
    const res = await axios.post<{
      data?: { tokens?: { access_token: string; refresh_token?: string } };
    }>(
      `${baseURL}/api/v1/auth/refresh`,
      { refresh_token: refreshToken },
      { headers: { "Content-Type": "application/json" }, timeout: REFRESH_TIMEOUT_MS }
    );
    const tokens = res.data?.data?.tokens;
    if (!tokens?.access_token) throw new Error("refresh response did not include an access token");

    window.localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, tokens.access_token);
    if (tokens.refresh_token) {
      window.localStorage.setItem(REFRESH_TOKEN_STORAGE_KEY, tokens.refresh_token);
    }
    return tokens.access_token;
  } catch (err) {
    // Only a 401 means the refresh token is dead. Network errors, timeouts and
    // 5xx (server restarting, cold start) are temporary: keep the session.
    if (axios.isAxiosError(err) && err.response?.status === 401) {
      throw new SessionExpiredError("refresh token rejected");
    }
    throw err;
  }
}

let refreshInFlight: Promise<string> | null = null;

/**
 * Returns a fresh access token, refreshing at most once at a time across all
 * open tabs. Tabs share localStorage, so a tab that waited for the lock first
 * checks whether another tab already replaced the token it failed with.
 */
function getFreshAccessToken(failedToken: string | null): Promise<string> {
  if (refreshInFlight) return refreshInFlight;

  const run = async () => {
    const current = window.localStorage.getItem(AUTH_TOKEN_STORAGE_KEY);
    if (current && current !== failedToken) return current;
    return requestNewTokens();
  };

  const locked: Promise<string> =
    typeof navigator !== "undefined" && navigator.locks?.request
      ? navigator.locks.request(REFRESH_LOCK_NAME, run).then((token) => token)
      : run();

  const pending = locked.finally(() => {
    refreshInFlight = null;
  });
  refreshInFlight = pending;
  return pending;
}

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const original = error.config as RetriableRequest | undefined;

    if (
      error.response?.status !== 401 ||
      !original ||
      typeof window === "undefined" ||
      isAuthEndpoint(original.url)
    ) {
      return Promise.reject(error);
    }

    // Already retried with a fresh token and still unauthorized: the account
    // is no longer allowed in (e.g. blocked), so end the session.
    if (original._retry) {
      clearSessionAndRedirect();
      return Promise.reject(error);
    }
    original._retry = true;

    try {
      const token = await getFreshAccessToken(bearerToken(original));
      original.headers.Authorization = `Bearer ${token}`;
      return apiClient(original);
    } catch (refreshErr) {
      if (refreshErr instanceof SessionExpiredError) {
        clearSessionAndRedirect();
        return Promise.reject(error);
      }
      return Promise.reject(refreshErr);
    }
  }
);
