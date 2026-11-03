"use client";

import { useCallback, useEffect, useState } from "react";
import axios from "axios";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { adminApi, type EscrowQueueItem } from "@/lib/api/admin";
import { formatDate, formatEtb } from "@/lib/contracts";

function apiError(err: unknown, fallback: string) {
  if (axios.isAxiosError(err)) {
    const data = err.response?.data as { message?: string; error?: { type?: string } } | undefined;
    if (data?.error?.type === "invalid_contract_status") return "This contract is no longer waiting for payment.";
    return data?.message || (!err.response ? "You appear to be offline." : fallback);
  }
  return fallback;
}

/**
 * Until online payments exist, a business pays Adfluence directly (e.g. bank
 * transfer). An admin checks the money arrived and confirms it here, which
 * holds it in escrow and lets the creator start.
 */
export function AdminPayments() {
  const [items, setItems] = useState<EscrowQueueItem[]>([]);
  const [state, setState] = useState<"loading" | "ready" | "error">("loading");
  const [openId, setOpenId] = useState<string | null>(null);
  const [reference, setReference] = useState("");
  const [busy, setBusy] = useState(false);

  const load = useCallback(async () => {
    try {
      const res = await adminApi.listEscrowQueue();
      setItems(res.data?.items ?? []);
      setState("ready");
    } catch {
      setState("error");
    }
  }, []);

  useEffect(() => {
    void Promise.resolve().then(load);
  }, [load]);

  const confirm = async (item: EscrowQueueItem) => {
    setBusy(true);
    try {
      await adminApi.fundEscrow(item.id, reference.trim());
      toast.success(`Payment confirmed. "${item.title}" is now in progress.`);
      setOpenId(null);
      setReference("");
      await load();
    } catch (err) {
      toast.error(apiError(err, "The payment could not be confirmed."));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div className="space-y-6 text-white">
      <header>
        <h1 className="text-2xl font-semibold">Payments to confirm</h1>
        <p className="mt-1 max-w-2xl text-sm text-[#CBD5E1]">
          Contracts a creator accepted. Confirm each business&apos;s payment only after you have checked it arrived; the
          creator then starts the work.
        </p>
      </header>

      {state === "loading" ? (
        <p role="status" className="flex items-center gap-2 text-sm text-[#94A3B8]">
          <Loader2 className="size-4 animate-spin" aria-hidden /> Loading…
        </p>
      ) : state === "error" ? (
        <div role="alert" className="rounded-2xl border border-rose-400/30 bg-rose-500/10 p-6 text-sm text-rose-100">
          Payments could not be loaded.{" "}
          <button type="button" className="underline" onClick={() => { setState("loading"); void load(); }}>
            Try again
          </button>
        </div>
      ) : items.length === 0 ? (
        <div className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-8 text-center text-sm text-[#CBD5E1]">
          No payments waiting for confirmation.
        </div>
      ) : (
        <ul className="space-y-3">
          {items.map((it) => (
            <li key={it.id} className="rounded-2xl border border-white/10 bg-[rgba(255,255,255,0.08)] p-4">
              <div className="flex flex-wrap items-start justify-between gap-3">
                <div className="min-w-0">
                  <p className="font-semibold">{it.title}</p>
                  <p className="text-sm text-[#CBD5E1]">
                    {it.business_name} → {it.creator_name} · {it.campaign_title}
                  </p>
                  <p className="mt-1 text-xs text-[#94A3B8]">Accepted by the creator {formatDate(it.accepted_at)}</p>
                </div>
                <div className="text-right">
                  <p className="text-lg font-semibold tabular-nums">{formatEtb(it.total_business_payable_etb)}</p>
                  <p className="text-xs text-[#94A3B8]">to receive from the business (fee included)</p>
                </div>
              </div>

              {openId === it.id ? (
                <form
                  className="mt-3 space-y-2 rounded-xl border border-white/15 bg-black/20 p-3"
                  onSubmit={(e) => {
                    e.preventDefault();
                    void confirm(it);
                  }}
                >
                  <label className="block space-y-1">
                    <span className="text-sm text-[#CBD5E1]">Payment reference (e.g. bank transfer number)</span>
                    <Input
                      value={reference}
                      onChange={(e) => setReference(e.target.value)}
                      maxLength={200}
                      autoFocus
                      className="h-11 rounded-xl border-white/10 bg-[#06111F] text-white"
                    />
                  </label>
                  <p className="text-xs text-amber-100">
                    Only confirm after checking that {formatEtb(it.total_business_payable_etb)} arrived. This cannot be undone.
                  </p>
                  <div className="flex flex-wrap gap-2">
                    <Button
                      type="submit"
                      disabled={busy || reference.trim().length < 3}
                      className="h-11 rounded-xl bg-[color:var(--adfluence-action)] text-white hover:bg-[color:var(--adfluence-action-hover)] sm:h-10"
                    >
                      {busy ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
                      Confirm payment received
                    </Button>
                    <Button type="button" variant="ghost" disabled={busy} onClick={() => setOpenId(null)} className="h-11 text-white hover:bg-white/10 sm:h-10">
                      Cancel
                    </Button>
                  </div>
                </form>
              ) : (
                <Button
                  type="button"
                  onClick={() => { setOpenId(it.id); setReference(""); }}
                  className="mt-3 h-11 rounded-xl bg-white/10 text-white hover:bg-white/15 sm:h-10"
                >
                  Confirm payment…
                </Button>
              )}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
