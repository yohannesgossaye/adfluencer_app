"use client";

import Image from "next/image";
import Link from "next/link";
import { useCallback, useMemo, useState } from "react";
import { ArrowLeft, Paperclip, Send } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import type { BusinessChatMessage, BusinessConversation } from "@/data/business-chat-mock";
import { initialBusinessConversations } from "@/data/business-chat-mock";
import { cn } from "@/lib/utils";

function formatTime(iso: string) {
  try {
    return new Intl.DateTimeFormat("en-GB", {
      day: "numeric",
      month: "short",
      hour: "2-digit",
      minute: "2-digit",
    }).format(new Date(iso));
  } catch {
    return iso;
  }
}

function formatListTime(iso: string) {
  try {
    const d = new Date(iso);
    const now = new Date();
    const sameDay =
      d.getDate() === now.getDate() &&
      d.getMonth() === now.getMonth() &&
      d.getFullYear() === now.getFullYear();
    if (sameDay) {
      return new Intl.DateTimeFormat("en-GB", {
        hour: "2-digit",
        minute: "2-digit",
      }).format(d);
    }
    return new Intl.DateTimeFormat("en-GB", {
      day: "numeric",
      month: "short",
    }).format(d);
  } catch {
    return "";
  }
}

const quickReplies = [
  "Can you share your rate?",
  "Let’s confirm deliverables",
  "Can we start this week?",
] as const;

export function BusinessMessages() {
  const [conversations, setConversations] = useState<BusinessConversation[]>(() =>
    structuredClone(initialBusinessConversations),
  );
  const [selectedId, setSelectedId] = useState<string>("conv-1");
  const [mobileShowChat, setMobileShowChat] = useState(false);
  const [search, setSearch] = useState("");
  const [draft, setDraft] = useState("");

  const selected = useMemo(
    () => conversations.find((c) => c.id === selectedId) ?? conversations[0],
    [conversations, selectedId],
  );

  const filteredList = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return conversations;
    return conversations.filter(
      (c) =>
        c.influencerName.toLowerCase().includes(q) ||
        c.campaignName.toLowerCase().includes(q) ||
        c.lastMessage.toLowerCase().includes(q),
    );
  }, [conversations, search]);

  const selectConversation = useCallback((id: string) => {
    setSelectedId(id);
    setConversations((prev) =>
      prev.map((c) => (c.id === id ? { ...c, unreadCount: 0 } : c)),
    );
    if (typeof window !== "undefined" && window.matchMedia("(max-width: 767px)").matches) {
      setMobileShowChat(true);
    }
  }, []);

  const sendMessage = useCallback(
    (text: string) => {
      const trimmed = text.trim();
      if (!trimmed || !selected) return;
      const newMsg: BusinessChatMessage = {
        id: `local-${Date.now()}`,
        role: "business",
        body: trimmed,
        sentAt: new Date().toISOString(),
      };
      setConversations((prev) =>
        prev.map((c) => {
          if (c.id !== selected.id) return c;
          return {
            ...c,
            messages: [...c.messages, newMsg],
            lastMessage: trimmed,
            lastTimeIso: newMsg.sentAt,
          };
        }),
      );
      setDraft("");
    },
    [selected],
  );

  const attachPlaceholder = useCallback(() => {
    if (!selected) return;
    const newMsg: BusinessChatMessage = {
      id: `local-${Date.now()}`,
      role: "business",
      body: "Shared an attachment",
      sentAt: new Date().toISOString(),
      attachmentLabel: "attachment.pdf",
    };
    setConversations((prev) =>
      prev.map((c) =>
        c.id === selected.id
          ? {
              ...c,
              messages: [...c.messages, newMsg],
              lastMessage: newMsg.body,
              lastTimeIso: newMsg.sentAt,
            }
          : c,
      ),
    );
  }, [selected]);

  return (
    <div className="flex min-h-0 flex-col gap-4 text-white">
      <div className="flex min-h-0 flex-1 flex-col gap-4 md:flex-row md:gap-0 md:overflow-hidden md:rounded-2xl md:border md:border-white/10">
        {/* Conversation list */}
        <aside
          className={cn(
            "flex w-full shrink-0 flex-col border-white/10 bg-[rgba(255,255,255,0.08)] md:w-[300px] md:border-r md:py-4",
            mobileShowChat && "hidden md:flex",
          )}
        >
          <div className="border-b border-white/10 px-3 pb-3 md:px-4">
            <h1 className="text-lg font-semibold text-white">Messages</h1>
            <p className="mt-0.5 text-xs text-[#CBD5E1]">Campaign conversations</p>
            <Input
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search conversations"
              className="mt-3 h-10 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-sm text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
            />
          </div>
          <div className="min-h-0 flex-1 space-y-1 overflow-y-auto px-2 py-2 md:px-3">
            {filteredList.map((c) => {
              const active = c.id === selectedId;
              return (
                <button
                  key={c.id}
                  type="button"
                  onClick={() => selectConversation(c.id)}
                  className={cn(
                    "flex w-full gap-3 rounded-xl border border-transparent p-3 text-left transition-colors",
                    active
                      ? "border-white/15 bg-[color:var(--adfluence-orange)]/12"
                      : "hover:border-white/10 hover:bg-white/5",
                  )}
                >
                  <div className="relative h-12 w-12 shrink-0 overflow-hidden rounded-full border border-white/10">
                    <Image
                      src={c.imageUrl}
                      alt={c.influencerName}
                      fill
                      className="object-cover"
                      sizes="48px"
                    />
                    {c.online ? (
                      <span className="absolute bottom-0 right-0 h-3 w-3 rounded-full border-2 border-[color:var(--adfluence-surface)] bg-emerald-400" />
                    ) : null}
                  </div>
                  <div className="min-w-0 flex-1">
                    <div className="flex items-start justify-between gap-2">
                      <p className="truncate text-sm font-medium text-white">{c.influencerName}</p>
                      <span className="shrink-0 text-[10px] text-[#CBD5E1]">
                        {formatListTime(c.lastTimeIso)}
                      </span>
                    </div>
                    <p className="truncate text-xs text-[#94A3B8]">{c.campaignName}</p>
                    <p className="mt-0.5 truncate text-xs text-[#CBD5E1]">{c.lastMessage}</p>
                    {c.unreadCount > 0 ? (
                      <span className="mt-1 inline-flex min-w-[1.25rem] items-center justify-center rounded-full bg-[color:var(--adfluence-orange)] px-1.5 py-0.5 text-[10px] font-semibold text-white">
                        {c.unreadCount}
                      </span>
                    ) : null}
                  </div>
                </button>
              );
            })}
          </div>
        </aside>

        {/* Chat */}
        <section
          className={cn(
            "flex min-h-[70vh] flex-1 flex-col bg-[#0B1A2B] md:min-h-0",
            !mobileShowChat && "hidden md:flex",
          )}
        >
          {selected ? (
            <>
              <header className="flex shrink-0 flex-wrap items-start gap-3 border-b border-white/10 bg-[rgba(255,255,255,0.06)] px-3 py-3 md:px-4">
                <Button
                  type="button"
                  variant="ghost"
                  size="icon"
                  className="shrink-0 text-white hover:bg-white/10 md:hidden"
                  aria-label="Back to conversations"
                  onClick={() => setMobileShowChat(false)}
                >
                  <ArrowLeft className="size-5" />
                </Button>
                <div className="relative h-12 w-12 shrink-0 overflow-hidden rounded-full border border-white/10">
                  <Image
                    src={selected.imageUrl}
                    alt={selected.influencerName}
                    fill
                    className="object-cover"
                    sizes="48px"
                  />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="font-semibold text-white">{selected.influencerName}</p>
                  <p className="text-xs text-[#CBD5E1]">
                    {selected.niche} · {selected.platform}
                  </p>
                  <p className="text-xs text-emerald-300/90">
                    {selected.online ? "Online" : selected.replyHint}
                  </p>
                </div>
                <div className="flex w-full shrink-0 flex-wrap gap-2 sm:w-auto sm:justify-end">
                  <Button
                    asChild
                    size="sm"
                    className="h-9 flex-1 rounded-xl bg-white/10 text-white hover:bg-white/15 sm:flex-none"
                  >
                    <Link href={`/business/influencers/${selected.influencerId}`}>
                      View Profile
                    </Link>
                  </Button>
                  <Button
                    asChild
                    size="sm"
                    className="h-9 flex-1 rounded-xl bg-[color:var(--adfluence-orange)] text-white hover:bg-[color:var(--adfluence-orange-hover)] sm:flex-none"
                  >
                    <Link href="/business/contracts">Create Contract</Link>
                  </Button>
                </div>
              </header>

              <div className="min-h-0 flex-1 space-y-3 overflow-y-auto px-3 py-4 md:px-4">
                {selected.messages.map((m) => {
                  if (m.role === "system") {
                    return (
                      <div key={m.id} className="flex justify-center">
                        <span className="rounded-full border border-white/10 bg-white/5 px-3 py-1 text-xs text-[#CBD5E1]">
                          {m.body}
                        </span>
                      </div>
                    );
                  }
                  const mine = m.role === "business";
                  return (
                    <div key={m.id} className={cn("flex", mine ? "justify-end" : "justify-start")}>
                      <div
                        className={cn(
                          "max-w-[85%] rounded-2xl border px-3 py-2 text-sm md:max-w-[70%]",
                          mine
                            ? "border-[color:var(--adfluence-orange)]/25 bg-[color:var(--adfluence-orange)]/15 text-white"
                            : "border-white/10 bg-[rgba(255,255,255,0.08)] text-[#E2E8F0]",
                        )}
                      >
                        <p>{m.body}</p>
                        {m.attachmentLabel ? (
                          <span className="mt-2 inline-flex items-center gap-1 rounded-lg border border-white/10 bg-[color:var(--adfluence-surface)] px-2 py-1 text-xs text-[#CBD5E1]">
                            📎 {m.attachmentLabel}
                          </span>
                        ) : null}
                        <p className="mt-1 text-[10px] text-[#94A3B8]">{formatTime(m.sentAt)}</p>
                      </div>
                    </div>
                  );
                })}
              </div>

              <div className="shrink-0 border-t border-white/10 bg-[rgba(255,255,255,0.04)] px-3 py-3 md:px-4">
                <div className="mb-2 flex flex-wrap gap-2">
                  {quickReplies.map((q) => (
                    <button
                      key={q}
                      type="button"
                      onClick={() => sendMessage(q)}
                      className="rounded-full border border-white/10 bg-[color:var(--adfluence-surface)] px-3 py-1 text-xs text-[#CBD5E1] transition-colors hover:border-[color:var(--adfluence-orange)]/40 hover:text-white"
                    >
                      {q}
                    </button>
                  ))}
                </div>
                <div className="flex gap-2">
                  <Button
                    type="button"
                    size="icon"
                    variant="outline"
                    className="h-11 shrink-0 rounded-xl border-white/15 bg-[color:var(--adfluence-surface)] text-white hover:bg-white/10"
                    aria-label="Attach file"
                    onClick={attachPlaceholder}
                  >
                    <Paperclip className="size-4" />
                  </Button>
                  <Input
                    value={draft}
                    onChange={(e) => setDraft(e.target.value)}
                    onKeyDown={(e) => {
                      if (e.key === "Enter" && !e.shiftKey) {
                        e.preventDefault();
                        sendMessage(draft);
                      }
                    }}
                    placeholder="Write a message…"
                    className="h-11 min-w-0 flex-1 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                  />
                  <Button
                    type="button"
                    className="h-11 shrink-0 rounded-xl bg-[color:var(--adfluence-orange)] px-4 text-white hover:bg-[color:var(--adfluence-orange-hover)]"
                    onClick={() => sendMessage(draft)}
                  >
                    <Send className="size-4 md:mr-1" />
                    <span className="hidden md:inline">Send</span>
                  </Button>
                </div>
              </div>
            </>
          ) : null}
        </section>
      </div>
    </div>
  );
}
