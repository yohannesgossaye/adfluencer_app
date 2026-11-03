"use client";

import { useState } from "react";
import { Loader2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { cn } from "@/lib/utils";

export type ConfirmActionRequest = {
  title: string;
  description: string;
  confirmLabel: string;
  /** "danger" for actions that lock someone out or hide their content. */
  tone?: "danger" | "default";
  /** When set, the admin must give a reason (stored and shown to the user where the backend supports it). */
  reasonLabel?: string;
  reasonPlaceholder?: string;
  /** Runs the action. Throw to keep the dialog open; the caller shows the error. */
  onConfirm: (reason: string) => Promise<void>;
};

const MIN_REASON = 3;

/**
 * One confirmation step for admin actions. Render it once per screen and
 * pass the pending request (or null when closed).
 */
export function ConfirmActionDialog({
  request,
  onClose,
}: {
  request: ConfirmActionRequest | null;
  onClose: () => void;
}) {
  return (
    <Dialog open={request !== null} onOpenChange={(open) => !open && onClose()}>
      {request ? <ConfirmBody key={request.title} request={request} onClose={onClose} /> : null}
    </Dialog>
  );
}

function ConfirmBody({ request, onClose }: { request: ConfirmActionRequest; onClose: () => void }) {
  const [reason, setReason] = useState("");
  const [busy, setBusy] = useState(false);
  const needsReason = Boolean(request.reasonLabel);
  const reasonOk = !needsReason || reason.trim().length >= MIN_REASON;

  const submit = async () => {
    if (!reasonOk || busy) return;
    setBusy(true);
    try {
      await request.onConfirm(reason.trim());
      onClose();
    } catch {
      // The caller has already shown the error; stay open so the admin can retry.
    } finally {
      setBusy(false);
    }
  };

  return (
    <DialogContent
      onInteractOutside={(e) => busy && e.preventDefault()}
      onEscapeKeyDown={(e) => busy && e.preventDefault()}
    >
      <form
        onSubmit={(e) => {
          e.preventDefault();
          void submit();
        }}
      >
        <DialogHeader className="pr-14">
          <DialogTitle>{request.title}</DialogTitle>
          <DialogDescription>{request.description}</DialogDescription>
        </DialogHeader>
        {needsReason ? (
          <label className="block space-y-1 px-6 pt-4">
            <span className="text-sm text-[#CBD5E1]">{request.reasonLabel}</span>
            <textarea
              value={reason}
              onChange={(e) => setReason(e.target.value)}
              placeholder={request.reasonPlaceholder}
              maxLength={500}
              rows={3}
              autoFocus
              className="w-full rounded-xl border border-white/10 bg-[#06111F] px-3 py-2 text-sm text-white outline-none placeholder:text-[#64748B] focus:border-white/30"
            />
            {reason.length > 0 && !reasonOk ? (
              <span className="text-xs text-amber-200">Please write at least {MIN_REASON} characters.</span>
            ) : null}
          </label>
        ) : null}
        <DialogFooter className="mt-4">
          <Button
            type="button"
            variant="ghost"
            disabled={busy}
            onClick={onClose}
            className="h-11 text-white hover:bg-white/10 sm:h-10"
          >
            Cancel
          </Button>
          <Button
            type="submit"
            disabled={busy || !reasonOk}
            className={cn(
              "h-11 rounded-xl text-white sm:h-10",
              request.tone === "danger"
                ? "bg-rose-600 hover:bg-rose-500"
                : "bg-[color:var(--adfluence-action)] hover:bg-[color:var(--adfluence-action-hover)]",
            )}
          >
            {busy ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : null}
            {request.confirmLabel}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  );
}
