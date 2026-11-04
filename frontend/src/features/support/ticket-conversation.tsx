"use client";

import { useCallback, useEffect, useState } from "react";
import axios from "axios";
import { ArrowLeft, Loader2 } from "lucide-react";
import { toast } from "sonner";

import { StatusBadge } from "@/components/common/status-badge";
import { ReplyBox, SupportThread } from "@/components/support/support-thread";
import { Button } from "@/components/ui/button";
import { businessApi, type SupportTicket } from "@/lib/api/business";
import { influencerApi } from "@/lib/api/influencer";
import { formatDate } from "@/lib/contracts";
import { ticketStatus } from "@/lib/status";
import { isTicketDone } from "@/lib/support";

import type { SupportRole } from "./user-support-page";

function replyError(err: unknown) {
  if (!axios.isAxiosError(err)) return "Your reply could not be sent. Please try again.";
  if (!err.response) return "You appear to be offline. Your reply was not sent.";
  const data = err.response.data as { message?: string; error?: { details?: Array<{ description?: string }> } } | undefined;
  return data?.error?.details?.[0]?.description || data?.message || "Your reply could not be sent. Please try again.";
}

/** One support request and its conversation, for the person who opened it. */
export function TicketConversation({
  role,
  ticketId,
  onBack,
  onChanged,
}: {
  role: SupportRole;
  ticketId: string;
  onBack: () => void;
  onChanged: () => void;
}) {
  const [ticket, setTicket] = useState<SupportTicket | null>(null);
  const [state, setState] = useState<"loading" | "ready" | "error">("loading");
  const api = role === "business" ? businessApi : influencerApi;

  const load = useCallback(async () => {
    try {
      const res = await api.getSupportTicket(ticketId);
      setTicket(res.data ?? null);
      setState("ready");
    } catch {
      setState("error");
    }
  }, [api, ticketId]);

  useEffect(() => {
    void Promise.resolve().then(load);
  }, [load]);

  const send = async (text: string) => {
    try {
      const res = await api.replySupportTicket(ticketId, text);
      toast.success(res.data?.reopened ? "Reply sent. We've reopened your request." : "Reply sent. We'll answer here and by email.");
    } catch (err) {
      toast.error(replyError(err));
      throw err;
    }
    await load();
    onChanged();
  };

  return (
    <div className="space-y-4">
      <Button type="button" variant="ghost" size="sm" onClick={onBack} className="-ml-2 text-white hover:bg-white/10">
        <ArrowLeft className="mr-1 size-4" aria-hidden /> My requests
      </Button>
      {state === "loading" ? (
        <p className="flex items-center gap-2 text-sm text-[#94A3B8]">
          <Loader2 className="size-4 animate-spin" aria-hidden /> Loading conversation…
        </p>
      ) : state === "error" || !ticket ? (
        <p role="alert" className="text-sm text-rose-200">
          This request could not be loaded.{" "}
          <button type="button" className="underline" onClick={() => void load()}>
            Try again
          </button>
        </p>
      ) : (
        <>
          <div>
            <div className="flex flex-wrap items-center gap-2">
              <h2 className="text-lg font-semibold">{ticket.subject}</h2>
              <StatusBadge status={ticketStatus(ticket.status)} />
            </div>
            <p className="mt-0.5 text-xs text-[#94A3B8]">
              {ticket.issue_type_label || ticket.issue_type} · opened {formatDate(ticket.created_at)}
              {ticket.related_campaign_title ? ` · ${ticket.related_campaign_title}` : ""}
            </p>
          </div>
          <SupportThread
            viewer="user"
            opening={{ name: "You", body: ticket.message ?? "", created_at: ticket.created_at }}
            messages={ticket.messages ?? []}
          />
          {ticket.messages?.length === 0 ? (
            <p className="text-xs text-[#94A3B8]">Our team will reply here, and you&apos;ll get an email when they do.</p>
          ) : null}
          <ReplyBox
            label={isTicketDone(ticket.status) ? "Still need help?" : "Add a message"}
            placeholder="Write your message…"
            hint={isTicketDone(ticket.status) ? "Sending a message reopens this request." : "Ctrl + Enter to send"}
            sendLabel="Send"
            onSend={send}
          />
        </>
      )}
    </div>
  );
}
