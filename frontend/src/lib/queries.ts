"use client";

/**
 * Cached reads shared by several screens (React Query). One key per kind of
 * data, so two screens asking for the same thing share one request, and an
 * action that changes it can refresh every screen with `invalidate…`.
 */
import { keepPreviousData, useQuery, useQueryClient, type UseQueryResult } from "@tanstack/react-query";

import type { LoadState } from "@/hooks/use-paged-list";
import { businessApi } from "@/lib/api/business";
import { chatApi } from "@/lib/api/chat";
import { influencerApi } from "@/lib/api/influencer";

export type Portal = "business" | "influencer";

export const queryKeys = {
  conversations: ["chat", "conversations"] as const,
  dashboard: (portal: Portal) => [portal, "dashboard"] as const,
  contracts: (portal: Portal, status?: string) => [portal, "contracts", { status: status ?? "" }] as const,
  contract: (portal: Portal, id: string) => [portal, "contract", id] as const,
};

/** The loading/ready/error state the list components expect. Data already shown stays on refresh errors. */
export function loadState(q: Pick<UseQueryResult, "data" | "isError" | "isFetching">): LoadState {
  if (q.data !== undefined) return "ready";
  return q.isError && !q.isFetching ? "error" : "loading";
}

export function useConversations() {
  return useQuery({
    queryKey: queryKeys.conversations,
    queryFn: async () => {
      const data = (await chatApi.listConversations()).data;
      return Array.isArray(data) ? data : [];
    },
  });
}

export function useBusinessDashboard() {
  return useQuery({
    queryKey: queryKeys.dashboard("business"),
    queryFn: async () => (await businessApi.getDashboardSummary()).data ?? null,
  });
}

export function useInfluencerDashboard() {
  return useQuery({
    queryKey: queryKeys.dashboard("influencer"),
    queryFn: async () => (await influencerApi.getDashboardSummary()).data ?? null,
  });
}

const contractApi = { business: businessApi, influencer: influencerApi };

export function useContracts(portal: Portal, status?: string) {
  return useQuery({
    queryKey: queryKeys.contracts(portal, status),
    queryFn: async () => (await contractApi[portal].listContracts({ status: status || undefined, limit: 50 })).data?.items ?? [],
    // Changing the filter keeps the current rows until the new ones arrive.
    placeholderData: keepPreviousData,
  });
}

export function useContract(portal: Portal, id: string | null) {
  return useQuery({
    queryKey: queryKeys.contract(portal, id ?? ""),
    queryFn: async () => (await contractApi[portal].getContract(id!)).data ?? null,
    enabled: Boolean(id),
  });
}

/** After a contract action: refresh the contract, every contract list and the dashboard. */
export function useInvalidateContracts(portal: Portal) {
  const client = useQueryClient();
  return (id?: string) =>
    Promise.all([
      client.invalidateQueries({ queryKey: [portal, "contracts"] }),
      id ? client.invalidateQueries({ queryKey: queryKeys.contract(portal, id) }) : null,
      client.invalidateQueries({ queryKey: queryKeys.dashboard(portal) }),
    ]);
}
