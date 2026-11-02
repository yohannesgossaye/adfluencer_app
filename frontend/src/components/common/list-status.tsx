"use client";

import type { ReactNode } from "react";
import { Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import type { LoadState } from "@/hooks/use-paged-list";

/**
 * Shows loading, a load failure (with Try again) or the empty message, and
 * otherwise the children. A failed load never reads as "nothing found".
 */
export function ListStatus({
  state,
  isEmpty,
  emptyText,
  errorText,
  onRetry,
  children,
}: {
  state: LoadState;
  isEmpty: boolean;
  emptyText: ReactNode;
  errorText: string;
  onRetry: () => void;
  children: ReactNode;
}) {
  if (state === "loading") {
    return (
      <p role="status" className="flex items-center gap-2 py-6 text-sm text-[#94A3B8]">
        <Loader2 className="size-4 animate-spin" aria-hidden /> Loading…
      </p>
    );
  }
  if (state === "error") {
    return (
      <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-6 text-sm text-rose-100">
        {errorText} Check your connection and{" "}
        <button type="button" className="font-medium underline underline-offset-2" onClick={onRetry}>
          try again
        </button>
        .
      </div>
    );
  }
  if (isEmpty) {
    return (
      <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center text-sm text-[#CBD5E1]">
        {emptyText}
      </div>
    );
  }
  return <>{children}</>;
}

/** "Load more" under a paged list, with how many are shown. */
export function LoadMore({
  shown,
  total,
  hasMore,
  loading,
  onLoadMore,
}: {
  shown: number;
  total: number | null;
  hasMore: boolean;
  loading: boolean;
  onLoadMore: () => void;
}) {
  if (total === null || shown === 0) return null;
  return (
    <div className="flex flex-col items-center gap-2 pt-2">
      <p className="text-xs text-[#94A3B8]">
        Showing {shown} of {total}
      </p>
      {hasMore ? (
        <Button
          type="button"
          variant="outline"
          disabled={loading}
          onClick={onLoadMore}
          className="h-11 rounded-xl border-white/15 bg-transparent text-white hover:bg-white/10 sm:h-9"
        >
          {loading ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
          Load more
        </Button>
      ) : null}
    </div>
  );
}
