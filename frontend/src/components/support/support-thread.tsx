"use client";

import { useEffect, useRef, useState } from "react";
import { Loader2, Send } from "lucide-react";

import { Button } from "@/components/ui/button";
import { SUPPORT_MESSAGE_MAX, type SupportMessage } from "@/lib/support";
import { cn } from "@/lib/utils";

function when(iso: string) {
  try {
    return new Intl.DateTimeFormat("en-GB", { day: "numeric", month: "short", hour: "2-digit", minute: "2-digit" }).format(new Date(iso));
  } catch {
    return iso;
  }
}

/**
 * A support conversation: the opening message, then every reply, oldest
 * first. `viewer` decides which side is "mine" (shown on the right).
 */
export function SupportThread({
  opening,
  messages,
  viewer,
}: {
  opening: { name: string; body: string; created_at: string };
  messages: SupportMessage[];
  viewer: "user" | "support";
}) {
  const all: SupportMessage[] = [{ id: "opening", from: "user", sender_name: opening.name, body: opening.body, created_at: opening.created_at }, ...messages];
  return (
    <ol aria-label="Conversation" className="space-y-3">
      {all.map((m) => {
        const mine = m.from === viewer;
        const name = mine ? "You" : m.from === "support" ? `${m.sender_name} · Adfluence support` : m.sender_name;
        return (
          <li key={m.id} className={cn("flex", mine ? "justify-end" : "justify-start")}>
            <div
              className={cn(
                "max-w-[85%] rounded-2xl border px-3.5 py-2.5 text-sm",
                mine
                  ? "rounded-br-md border-[color:var(--adfluence-orange)]/25 bg-[color:var(--adfluence-orange)]/10"
                  : "rounded-bl-md border-white/10 bg-[rgba(255,255,255,0.06)]",
              )}
            >
              <p className="text-xs text-[#94A3B8]">
                <span className="font-medium text-[#E2E8F0]">{name}</span> · {when(m.created_at)}
              </p>
              <p className="mt-1 whitespace-pre-wrap break-words text-[#F1F5F9]">{m.body}</p>
            </div>
          </li>
        );
      })}
    </ol>
  );
}

/** The box under a conversation. `onSend` throws to keep the draft. */
export function ReplyBox({
  onSend,
  label,
  placeholder,
  hint,
  sendLabel = "Send reply",
}: {
  onSend: (text: string) => Promise<void>;
  label: string;
  placeholder: string;
  hint?: string;
  sendLabel?: string;
}) {
  const [draft, setDraft] = useState("");
  const [busy, setBusy] = useState(false);
  const ref = useRef<HTMLTextAreaElement>(null);
  const text = draft.trim();

  useEffect(() => {
    const el = ref.current;
    if (!el) return;
    el.style.height = "auto";
    el.style.height = `${Math.min(el.scrollHeight, 240)}px`;
  }, [draft]);

  const send = async () => {
    if (!text || busy) return;
    setBusy(true);
    try {
      await onSend(text);
      setDraft("");
    } catch {
      // The caller showed the error; keep the draft so nothing is lost.
    } finally {
      setBusy(false);
    }
  };

  return (
    <form
      className="space-y-2"
      onSubmit={(e) => {
        e.preventDefault();
        void send();
      }}
    >
      <label className="block space-y-1">
        <span className="text-sm text-[#CBD5E1]">{label}</span>
        <textarea
          ref={ref}
          value={draft}
          rows={3}
          maxLength={SUPPORT_MESSAGE_MAX}
          onChange={(e) => setDraft(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === "Enter" && (e.ctrlKey || e.metaKey)) {
              e.preventDefault();
              void send();
            }
          }}
          placeholder={placeholder}
          className="w-full resize-none rounded-xl border border-white/10 bg-[#06111F] px-3 py-2 text-sm text-white outline-none placeholder:text-[#64748B] focus:border-[color:var(--adfluence-orange)] focus:ring-2 focus:ring-[color:var(--adfluence-orange)]/25"
        />
      </label>
      <div className="flex flex-wrap items-center justify-between gap-2">
        <p className="text-xs text-[#94A3B8]">{hint}</p>
        <Button
          type="submit"
          disabled={!text || busy}
          className="h-11 rounded-xl bg-[color:var(--adfluence-action)] px-5 text-white hover:bg-[color:var(--adfluence-action-hover)] sm:h-10"
        >
          {busy ? <Loader2 className="mr-2 size-4 animate-spin" aria-hidden /> : <Send className="mr-2 size-4" aria-hidden />}
          {sendLabel}
        </Button>
      </div>
    </form>
  );
}
