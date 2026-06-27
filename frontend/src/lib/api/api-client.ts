import axios, { type AxiosError } from "axios";

import { AUTH_TOKEN_STORAGE_KEY } from "@/lib/constants";
import { getPublicEnv } from "@/lib/env";

const { apiBaseUrl } = getPublicEnv();

export const apiClient = axios.create({
  baseURL: apiBaseUrl || undefined,
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
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    if (error.response?.status === 401 && typeof window !== "undefined") {
      window.localStorage.removeItem(AUTH_TOKEN_STORAGE_KEY);
      const path = `${window.location.pathname}${window.location.search}`;
      if (!path.startsWith("/auth")) {
        window.location.assign(`/auth/login?next=${encodeURIComponent(path)}`);
      }
    }
    return Promise.reject(error);
  }
);
