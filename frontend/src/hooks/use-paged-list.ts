"use client";

import { useCallback, useMemo } from "react";
import { keepPreviousData, useInfiniteQuery, type QueryKey } from "@tanstack/react-query";
import { toast } from "sonner";

export type LoadState = "loading" | "ready" | "error";

// The API's list meta (common.ListMeta) reports the count as total_docs.
type Page = { items?: unknown[]; pagination?: { total_docs?: number } };
type ItemOf<R extends Page> = NonNullable<R["items"]>[number];

/**
 * A paged list with "Load more", cached by React Query.
 *
 * `key` names the list and its filters, e.g. ["admin", "users", { role, status }];
 * `fetchPage` loads one page for those same filters. Coming back to a list
 * shows the cached rows at once and refreshes them in the background; while
 * the filters change, the previous rows stay on screen until the new ones
 * arrive. `last` is the first page's raw response, for extras such as stats.
 */
export function usePagedList<R extends Page>(key: QueryKey, fetchPage: (page: number) => Promise<R | undefined>) {
  type T = ItemOf<R>;
  const query = useInfiniteQuery({
    queryKey: key,
    queryFn: async ({ pageParam }) => (await fetchPage(pageParam)) ?? ({} as R),
    initialPageParam: 1,
    getNextPageParam: (lastPage, pages) => {
      const total = lastPage.pagination?.total_docs;
      const loaded = pages.reduce((n, p) => n + (p.items?.length ?? 0), 0);
      return total != null && loaded < total && (lastPage.items?.length ?? 0) > 0 ? pages.length + 1 : undefined;
    },
    placeholderData: keepPreviousData,
  });

  const pages = query.data?.pages;
  const items = useMemo(() => (pages ?? []).flatMap((p) => (p.items ?? []) as T[]), [pages]);
  const total = pages?.length ? (pages[pages.length - 1].pagination?.total_docs ?? null) : null;

  // Rows already on screen stay there if a background refresh fails.
  const state: LoadState = query.data ? "ready" : query.isError && !query.isFetching ? "error" : "loading";

  const { refetch, fetchNextPage } = query;
  /** Reload after an action changed the data; the current rows stay on screen meanwhile. */
  const reload = useCallback(async () => {
    await refetch();
  }, [refetch]);
  /** Try again after an error. */
  const retry = reload;
  const loadMore = useCallback(async () => {
    const res = await fetchNextPage();
    if (res.isFetchNextPageError) toast.error("More results could not be loaded. Please try again.");
  }, [fetchNextPage]);

  return {
    items,
    last: pages?.[0],
    state,
    error: query.error,
    total,
    hasMore: query.hasNextPage,
    loadingMore: query.isFetchingNextPage,
    loadMore,
    reload,
    retry,
  };
}
