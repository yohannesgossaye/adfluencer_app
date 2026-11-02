"use client";

import { useCallback, useEffect, useRef, useState } from "react";
import { toast } from "sonner";

export type LoadState = "loading" | "ready" | "error";

type Page = { items?: unknown[]; pagination?: { total?: number } };
type ItemOf<R extends Page> = NonNullable<R["items"]>[number];

/**
 * Loads a paged list and appends further pages on demand ("Load more").
 *
 * `fetchPage` must be memoised (useCallback) with the filters it depends on;
 * when it changes, the list starts again from page 1. Late responses from an
 * older filter are ignored. `last` is the latest raw response, for extras such
 * as stats that come with page 1.
 */
export function usePagedList<R extends Page>(fetchPage: (page: number) => Promise<R | undefined>) {
  type T = ItemOf<R>;
  const [items, setItems] = useState<T[]>([]);
  const [last, setLast] = useState<R | undefined>();
  const [state, setState] = useState<LoadState>("loading");
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState<number | null>(null);
  const [loadingMore, setLoadingMore] = useState(false);
  const [error, setError] = useState<unknown>(null);
  const request = useRef(0);

  const load = useCallback(
    async (target: number) => {
      const id = ++request.current;
      if (target > 1) setLoadingMore(true);
      try {
        const res = await fetchPage(target);
        if (id !== request.current) return;
        const got = (res?.items ?? []) as T[];
        setItems((prev) => (target === 1 ? got : [...prev, ...got]));
        setTotal(res?.pagination?.total ?? null);
        setPage(target);
        if (target === 1) setLast(res);
        setError(null);
        setState("ready");
      } catch (err) {
        if (id !== request.current) return;
        // A failed "Load more" keeps what is shown (loadMore shows a toast).
        if (target === 1) {
          setError(err);
          setState("error");
        }
        else throw err;
      } finally {
        if (id === request.current) setLoadingMore(false);
      }
    },
    [fetchPage],
  );

  useEffect(() => {
    void Promise.resolve(1).then(load);
  }, [load]);

  /** Reload from page 1, e.g. after an action changed the data. Keeps the current rows on screen meanwhile. */
  const reload = useCallback(() => load(1), [load]);

  /** Full retry after an error: shows the loading state again. */
  const retry = useCallback(() => {
    setState("loading");
    return load(1);
  }, [load]);

  const loadMore = useCallback(
    () => load(page + 1).catch(() => toast.error("More results could not be loaded. Please try again.")),
    [load, page],
  );

  const hasMore = total !== null && items.length < total;

  return { items, setItems, last, state, error, total, hasMore, loadingMore, loadMore, reload, retry };
}
