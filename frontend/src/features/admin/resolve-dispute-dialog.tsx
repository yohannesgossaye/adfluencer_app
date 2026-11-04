"use client";

import { useState } from "react";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";

import { Button } from "@/components/ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "@/components/ui/dialog";
import { adminApi, type AdminDisputeDetail } from "@/lib/api/admin";
import { adminErrorMessage } from "@/lib/api-errors";
import { DECISION_NOTE_MIN, DISPUTE_OUTCOMES, type DisputeOutcome } from "@/lib/disputes";
import { cn } from "@/lib/utils";

/** The admin's decision on a dispute. Both sides are emailed the outcome and the note. */
export function ResolveDisputeDialog({
  dispute,
  onClose,
  onResolved,
}: {
  dispute: AdminDisputeDetail | null;
  onClose: () => void;
  onResolved: () => Promise<void>;
}) {
  return (
    <Dialog open={dispute !== null} onOpenChange={(open) => !open && onClose()}>
      {dispute ? <ResolveBody key={dispute.id} dispute={dispute} onClose={onClose} onResolved={onResolved} /> : null}
    </Dialog>
  );
}

function ResolveBody({
  dispute,
  onClose,
  onResolved,
}: {
  dispute: AdminDisputeDetail;
  onClose: () => void;
  onResolved: () => Promise<void>;
}) {
  const [outcome, setOutcome] = useState<DisputeOutcome | null>(null);
  const [note, setNote] = useState("");
  const [busy, setBusy] = useState(false);
  const noteLength = note.trim().length;
  const valid = outcome !== null && noteLength >= DECISION_NOTE_MIN;
  const onHold = dispute.contract_status === "disputed";

  const submit = async () => {
    if (!valid || busy) return;
    setBusy(true);
    try {
      await adminApi.resolveDispute(dispute.id, { outcome, note: note.trim() });
      toast.success(`Dispute resolved: ${DISPUTE_OUTCOMES[outcome].label.toLowerCase()}. Both sides have been emailed.`);
      onClose();
      await onResolved();
    } catch (err) {
      toast.error(adminErrorMessage(err, "The dispute could not be resolved."));
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
          <DialogTitle>Resolve this dispute</DialogTitle>
          <DialogDescription>
            {dispute.business_name} and {dispute.influencer_name} are both emailed your decision and note.
            {onHold ? "" : " This contract is no longer on hold, so the decision only closes the dispute."}
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4 px-6 pt-4">
          <fieldset>
            <legend className="mb-2 text-sm text-[#CBD5E1]">Decision</legend>
            <div className="space-y-1.5">
              {(Object.keys(DISPUTE_OUTCOMES) as DisputeOutcome[]).map((id) => {
                const o = DISPUTE_OUTCOMES[id];
                return (
                  <label
                    key={id}
                    className={cn(
                      "flex cursor-pointer items-start gap-3 rounded-xl border px-3 py-2.5 text-sm transition-colors",
                      outcome === id
                        ? o.tone === "danger"
                          ? "border-rose-400/50 bg-rose-500/10"
                          : "border-[color:var(--adfluence-orange)]/50 bg-[color:var(--adfluence-orange)]/10"
                        : "border-white/10 hover:border-white/20",
                    )}
                  >
                    <input
                      type="radio"
                      name="outcome"
                      value={id}
                      checked={outcome === id}
                      onChange={() => setOutcome(id)}
                      className="mt-0.5 size-4 accent-[color:var(--adfluence-orange)]"
                    />
                    <span>
                      <span className="block font-medium text-white">{o.label}</span>
                      <span className="block text-xs text-[#94A3B8]">{o.effect}</span>
                    </span>
                  </label>
                );
              })}
            </div>
          </fieldset>
          {outcome === "terminate" && onHold ? (
            <p className="rounded-xl border border-amber-400/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-50">
              Refunds are manual for now: after ending the contract, send the business their refund and note the reference.
            </p>
          ) : null}

          <label className="block space-y-1">
            <span className="text-sm text-[#CBD5E1]">Note to both sides</span>
            <textarea
              value={note}
              onChange={(e) => setNote(e.target.value)}
              maxLength={2000}
              rows={3}
              placeholder="e.g. The creator posted the video on 3 Nov; the business confirmed it. Payout released."
              className="w-full rounded-xl border border-white/10 bg-[#06111F] px-3 py-2 text-sm text-white outline-none placeholder:text-[#64748B] focus:border-white/30"
            />
            {noteLength > 0 && noteLength < DECISION_NOTE_MIN ? (
              <span className="text-xs text-amber-200">Please write at least {DECISION_NOTE_MIN} characters.</span>
            ) : null}
          </label>
        </div>

        <DialogFooter className="mt-4">
          <Button type="button" variant="ghost" disabled={busy} onClick={onClose} className="h-11 text-white hover:bg-white/10 sm:h-10">
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={!valid || busy}
            className={cn(
              "h-11 rounded-xl text-white sm:h-10",
              outcome === "terminate"
                ? "bg-rose-600 hover:bg-rose-500"
                : "bg-[color:var(--adfluence-action)] hover:bg-[color:var(--adfluence-action-hover)]",
            )}
          >
            {busy ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
            {outcome ? DISPUTE_OUTCOMES[outcome].label : "Choose a decision"}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}
