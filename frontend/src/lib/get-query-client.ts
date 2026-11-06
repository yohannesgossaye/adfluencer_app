import { isServer, QueryClient } from "@tanstack/react-query";
import axios from "axios";
import { cache } from "react";

/** Client errors (bad request, not allowed, not found…) don't get better by retrying. */
function shouldRetry(failureCount: number, error: unknown) {
  const status = axios.isAxiosError(error) ? error.response?.status : undefined;
  if (status && status >= 400 && status < 500) return false;
  return failureCount < 2;
}

function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        // Revisiting a page shows the saved answer at once and refreshes it in
        // the background, so it is instant but never out of date after the
        // user's own actions on another page.
        staleTime: 0,
        // Switching browser tabs doesn't re-ask the server; opening a page does.
        refetchOnWindowFocus: false,
        // Pages not on screen keep their data for 5 minutes.
        gcTime: 5 * 60 * 1000,
        retry: shouldRetry,
      },
    },
  });
}

const getServerQueryClient = cache(makeQueryClient);

let browserQueryClient: QueryClient | undefined;

export function getQueryClient() {
  if (isServer) {
    return getServerQueryClient();
  }
  if (!browserQueryClient) {
    browserQueryClient = makeQueryClient();
  }
  return browserQueryClient;
}
