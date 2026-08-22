import axios, { type AxiosError } from "axios";

import { AUTH_TOKEN_STORAGE_KEY, REFRESH_TOKEN_STORAGE_KEY, AUTH_USER_STORAGE_KEY } from "@/lib/constants";
import { getPublicEnv } from "@/lib/env";

const { apiBaseUrl } = getPublicEnv();

export const apiClient = axios.create({
  baseURL: apiBaseUrl || "http://localhost:8080",
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

let isRefreshing = false;
let failedQueue: Array<{
  resolve: (value?: unknown) => void;
  reject: (reason?: unknown) => void;
}> = [];

const processQueue = (error: Error | null, token: string | null = null) => {
  failedQueue.forEach((prom) => {
    if (error) {
      prom.reject(error);
    } else {
      prom.resolve(token);
    }
  });
  failedQueue = [];
};

apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as typeof error.config & { _retry?: boolean };

    if (
      error.response?.status === 401 &&
      originalRequest &&
      !originalRequest._retry &&
      !originalRequest.url?.includes("/api/v1/auth/login") &&
      !originalRequest.url?.includes("/api/v1/auth/signup") &&
      !originalRequest.url?.includes("/api/v1/auth/refresh") &&
      typeof window !== "undefined"
    ) {
      const refreshToken = window.localStorage.getItem(REFRESH_TOKEN_STORAGE_KEY);
      if (!refreshToken) {
        window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
        window.localStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
        window.localStorage.removeItem(AUTH_USER_STORAGE_KEY);
        const path = `${window.location.pathname}${window.location.search}`;
        if (!path.startsWith("/auth")) {
          window.location.assign(`/auth/login?next=${encodeURIComponent(path)}`);
        }
        return Promise.reject(error);
      }

      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        })
          .then((token) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            return apiClient(originalRequest);
          })
          .catch((err) => Promise.reject(err));
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        const refreshRes = await axios.post<{ data: { access_token: string; refresh_token: string } }>(
          `${apiBaseUrl || "http://localhost:8080"}/api/v1/auth/refresh`,
          { refresh_token: refreshToken },
          { headers: { "Content-Type": "application/json" } }
        );

        const newAccessToken = refreshRes.data?.data?.access_token;
        const newRefreshToken = refreshRes.data?.data?.refresh_token;

        if (newAccessToken) {
          window.localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, newAccessToken);
          if (newRefreshToken) {
            window.localStorage.setItem(REFRESH_TOKEN_STORAGE_KEY, newRefreshToken);
          }
          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          processQueue(null, newAccessToken);
          return apiClient(originalRequest);
        }
      } catch (refreshErr) {
        processQueue(refreshErr as Error, null);
        window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
        window.localStorage.removeItem(REFRESH_TOKEN_STORAGE_KEY);
        window.localStorage.removeItem(AUTH_USER_STORAGE_KEY);
        const path = `${window.location.pathname}${window.location.search}`;
        if (!path.startsWith("/auth")) {
          window.location.assign(`/auth/login?next=${encodeURIComponent(path)}`);
        }
        return Promise.reject(refreshErr);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);

