"use client";

import { useState } from "react";
import axios from "axios";
import { Flag, Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { businessApi } from "@/lib/api/business";
import { influencerApi } from "@/lib/api/influencer";
import { DESCRIPTION_MAX, DESCRIPTION_MIN, DISPUTE_REASONS, type DisputeRole } from "@/lib/disputes";
import { cn } from "@/lib/utils";

function reportError(err: unknown) {
  if (!axios.isAxiosError(err)) return "The report could not be sent. Please try again.";
  if (!err.response) return "You appear to be offline. Please try again.";
  const data = err.response.data as { message?: string; error?: { type?: string; details?: Array<{ description?: string }> } } | undefined;
  switch (data?.error?.type) {
    case "dispute_already_open":
      return "A problem has already been reported on this contract. Our team is on it.";
    case "invalid_contract_status":
      return "Problems can only be reported while the work is in progress.";
  }
  return data?.error?.details?.[0]?.description || data?.message || "The report could not be sent. Please try again.";
}

/**
 * "Report a problem" on an in-progress contract. The contract is put on hold,
 * the payment stays held, and our team decides with both sides.
 */
export function ReportProblemButton({
  contractId,
  role,
  onReported,
}: {
  contractId: string;
  role: DisputeRole;
  onReported: () => void | Promise<void>;
}) {
  const [open, setOpen] = useState(false);
  return (
    <>
      <Button
        type="button"
        variant="ghost"
        onClick={() => setOpen(true)}
        className="h-10 rounded-xl text-rose-200 hover:bg-rose-500/10 hover:text-rose-100"
      >
        <Flag className="mr-2 size-4" aria-hidden /> Report a problem
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        {open ? (
          <ReportForm
            contractId={contractId}
            role={role}
            onDone={async () => {
              setOpen(false);
              await onReported();
            }}
            onCancel={() => setOpen(false)}
          />
        ) : null}
      </Dialog>
    </>
  );
}

function ReportForm({
  contractId,
  role,
  onDone,
  onCancel,
}: {
  contractId: string;
  role: DisputeRole;
  onDone: () => Promise<void>;
  onCancel: () => void;
}) {
  const [reason, setReason] = useState("");
  const [description, setDescription] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const length = description.trim().length;
  const valid = reason !== "" && length >= DESCRIPTION_MIN && length <= DESCRIPTION_MAX;
  const other = role === "business" ? "creator" : "business";

  const submit = async () => {
    if (!valid || busy) return;
    setBusy(true);
    setError(null);
    try {
      const api = role === "business" ? businessApi : influencerApi;
      await api.reportProblem(contractId, { reason, description: description.trim() });
      toast.success("Problem reported. We've emailed you and will be in touch.");
      await onDone();
    } catch (err) {
      setError(reportError(err));
    } finally {
      setBusy(false);
    }
  };

  return (
    <DialogContent onInteractOutside={(e) => busy && e.preventDefault()} onEscapeKeyDown={(e) => busy && e.preventDefault()}>
      <form
        onSubmit={(e) => {
          e.preventDefault();
          void submit();
        }}
      >
        <DialogHeader className="pr-14">
          <DialogTitle>Report a problem</DialogTitle>
          <DialogDescription>
            The contract goes on hold and the payment stays safely held while our team looks into it with you and the {other}. Both of
            you are told by email.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 px-6 pt-4">
          <fieldset>
            <legend className="mb-2 text-sm text-[#CBD5E1]">What went wrong?</legend>
            <div className="space-y-1.5">
              {DISPUTE_REASONS[role].map((r) => (
                <label
                  key={r.id}
                  className={cn(
                    "flex cursor-pointer items-center gap-3 rounded-xl border px-3 py-2.5 text-sm transition-colors",
                    reason === r.id
                      ? "border-[color:var(--adfluence-orange)]/50 bg-[color:var(--adfluence-orange)]/10 text-white"
                      : "border-white/10 text-[#E2E8F0] hover:border-white/20",
                  )}
                >
                  <input
                    type="radio"
                    name="reason"
                    value={r.id}
                    checked={reason === r.id}
                    onChange={() => setReason(r.id)}
                    className="size-4 accent-[color:var(--adfluence-orange)]"
                  />
                  {r.label}
                </label>
              ))}
            </div>
          </fieldset>

          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">Tell us what happened</span>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              maxLength={DESCRIPTION_MAX}
              rows={4}
              placeholder="Dates, what was agreed, what you've already tried. Our team reads this first."
              className="w-full rounded-xl border border-white/10 bg-[#06111F] px-3 py-2 text-sm text-white outline-none placeholder:text-[#64748B] focus:border-white/30"
            />
            <span className={cn("block text-xs", length > 0 && length < DESCRIPTION_MIN ? "text-amber-200" : "text-[#94A3B8]")}>
              {length < DESCRIPTION_MIN ? `At least ${DESCRIPTION_MIN} characters (${length} so far).` : `${length} / ${DESCRIPTION_MAX}`}
            </span>
          </label>

          {error ? (
            <p role="alert" className="rounded-xl border border-rose-400/30 bg-rose-500/10 px-3 py-2 text-sm text-rose-100">
              {error}
            </p>
          ) : null}
        </div>

        <DialogFooter className="mt-4">
          <Button type="button" variant="ghost" disabled={busy} onClick={onCancel} className="h-11 text-white hover:bg-white/10 sm:h-10">
            Cancel
          </Button>
          <Button type="submit" disabled={!valid || busy} className="h-11 rounded-xl bg-rose-600 text-white hover:bg-rose-500 sm:h-10">
            {busy ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
            Report problem
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}
