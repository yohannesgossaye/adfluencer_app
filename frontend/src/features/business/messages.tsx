"use client";

import Image from "next/image";
import Link from "next/link";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";
import { ArrowLeft, Paperclip, Send, Trash2 } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  chatApi,
  type ChatConversationItem,
  type ChatMessage,
} from "@/lib/api/chat";
import { useAuth } from "@/components/providers/auth-provider";
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

function formatListTime(iso?: string) {
  if (!iso) return "";
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
  const { user: authUser } = useAuth();
  const [conversations, setConversations] = useState<ChatConversationItem[]>([]);
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [selectedId, setSelectedId] = useState<string>("");
  const [mobileShowChat, setMobileShowChat] = useState(false);
  const [search, setSearch] = useState("");
  const [draft, setDraft] = useState("");
  const [isLoading, setIsLoading] = useState(true);
  const messagesEndRef = useRef<HTMLDivElement | null>(null);
  const wsRef = useRef<WebSocket | null>(null);

  const currentUserId = authUser?.id || "";

  const loadConversations = async () => {
    try {
      const res = await chatApi.listConversations();
      if (res.data?.items) {
        setConversations(res.data.items);
        if (!selectedId && res.data.items.length > 0) {
          setSelectedId(res.data.items[0].id);
        }
      }
    } catch {
      // fallback
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    loadConversations();
  }, []);

  const selected = useMemo(
    () => conversations.find((c) => c.id === selectedId) ?? conversations[0] ?? null,
    [conversations, selectedId],
  );

  const loadMessages = async (convId: string) => {
    if (!convId) return;
    try {
      const res = await chatApi.getMessages({ conversation_id: convId, limit: 50 });
      if (res.data?.items) {
        setMessages(res.data.items);
        setTimeout(() => {
          messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
        }, 100);
      }
      await chatApi.markMessagesAsRead({ conversation_id: convId });
      setConversations((prev) =>
        prev.map((c) => (c.id === convId ? { ...c, unread_count: 0 } : c)),
      );
    } catch {
      // fallback
    }
  };

  useEffect(() => {
    if (selected?.id) {
      loadMessages(selected.id);
    }
  }, [selected?.id]);

  // Real-time WebSocket connection
  useEffect(() => {
    const token = typeof window !== "undefined" ? localStorage.getItem("auth_access_token") : null;
    if (!token) return;

    const wsUrl = chatApi.getWebSocketUrl(token);
    const ws = new WebSocket(wsUrl);
    wsRef.current = ws;

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "message_created") {
          const newMsg: ChatMessage = data.payload;
          if (newMsg.conversation_id === selectedId) {
            setMessages((prev) => [...prev, newMsg]);
            setTimeout(() => {
              messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
            }, 100);
            chatApi.markMessagesAsRead({ conversation_id: selectedId });
          }
          loadConversations();
        } else if (data.type === "messages_read") {
          const { conversation_id, read_by_user_id } = data.payload;
          if (conversation_id === selectedId && read_by_user_id !== currentUserId) {
            setMessages((prev) =>
              prev.map((m) => (m.sender_user_id === currentUserId ? { ...m, is_read: true } : m)),
            );
          }
        } else if (data.type === "message_deleted") {
          const { message_id } = data.payload;
          setMessages((prev) => prev.filter((m) => m.id !== message_id));
        }
      } catch {
        // ignore malformed ws messages
      }
    };

    return () => {
      ws.close();
    };
  }, [selectedId, currentUserId]);

  const filteredList = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return conversations;
    return conversations.filter(
      (c) =>
        c.other_user_name.toLowerCase().includes(q) ||
        (c.campaign_title && c.campaign_title.toLowerCase().includes(q)) ||
        (c.last_message_preview && c.last_message_preview.toLowerCase().includes(q)),
    );
  }, [conversations, search]);

  const selectConversation = useCallback((id: string) => {
    setSelectedId(id);
    if (typeof window !== "undefined" && window.matchMedia("(max-width: 767px)").matches) {
      setMobileShowChat(true);
    }
  }, []);

  const sendMessage = useCallback(
    async (text: string) => {
      const trimmed = text.trim();
      if (!trimmed || !selected) return;

      try {
        const res = await chatApi.sendMessage({
          conversation_id: selected.id,
          receiver_user_id: selected.other_user_id,
          body: trimmed,
        });
        const created = res.data;
        if (created) {
          setMessages((prev) => [...prev, created]);
          setTimeout(() => {
            messagesEndRef.current?.scrollIntoView({ behavior: "smooth" });
          }, 100);
        }
        setDraft("");
        loadConversations();
      } catch {
        // fallback
      }
    },
    [selected],
  );

  const deleteMessage = async (msgId: string) => {
    try {
      await chatApi.deleteMessage(msgId, { delete_type: "everyone" });
      setMessages((prev) => prev.filter((m) => m.id !== msgId));
    } catch {
      // fallback
    }
  };

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
            {filteredList.length === 0 && !isLoading ? (
              <p className="p-4 text-center text-xs text-[#94A3B8]">No conversations yet.</p>
            ) : (
              filteredList.map((c) => {
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
                        src={c.other_user_avatar || "/images/influencers/influencer-1.jpg"}
                        alt={c.other_user_name}
                        fill
                        className="object-cover"
                        sizes="48px"
                      />
                    </div>
                    <div className="min-w-0 flex-1">
                      <div className="flex items-start justify-between gap-2">
                        <p className="truncate text-sm font-medium text-white">{c.other_user_name}</p>
                        <span className="shrink-0 text-[10px] text-[#CBD5E1]">
                          {formatListTime(c.last_message_at)}
                        </span>
                      </div>
                      <p className="truncate text-xs text-[#94A3B8]">{c.campaign_title || "Direct Message"}</p>
                      <p className="mt-0.5 truncate text-xs text-[#CBD5E1]">{c.last_message_preview || "Start conversation"}</p>
                      {c.unread_count > 0 ? (
                        <span className="mt-1 inline-flex min-w-[1.25rem] items-center justify-center rounded-full bg-[color:var(--adfluence-orange)] px-1.5 py-0.5 text-[10px] font-semibold text-white">
                          {c.unread_count}
                        </span>
                      ) : null}
                    </div>
                  </button>
                );
              })
            )}
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
                    src={selected.other_user_avatar || "/images/influencers/influencer-1.jpg"}
                    alt={selected.other_user_name}
                    fill
                    className="object-cover"
                    sizes="48px"
                  />
                </div>
                <div className="min-w-0 flex-1">
                  <p className="font-semibold text-white">{selected.other_user_name}</p>
                  <p className="text-xs text-[#CBD5E1]">
                    {selected.campaign_title || "Direct Message"} · <span className="capitalize">{selected.other_user_role}</span>
                  </p>
                </div>
                <div className="flex w-full shrink-0 flex-wrap gap-2 sm:w-auto sm:justify-end">
                  <Button
                    asChild
                    size="sm"
                    className="h-9 flex-1 rounded-xl bg-white/10 text-white hover:bg-white/15 sm:flex-none"
                  >
                    <Link href={`/business/influencers/${selected.other_user_id}`}>
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
                {messages.length === 0 ? (
                  <p className="py-8 text-center text-xs text-[#94A3B8]">
                    No messages yet. Send a message below to start collaborating!
                  </p>
                ) : (
                  messages.map((m) => {
                    const mine = m.sender_user_id === currentUserId || (!currentUserId && m.receiver_user_id !== selected.other_user_id);
                    return (
                      <div key={m.id} className={cn("group flex", mine ? "justify-end" : "justify-start")}>
                        <div
                          className={cn(
                            "relative max-w-[85%] rounded-2xl border px-3 py-2 text-sm md:max-w-[70%]",
                            mine
                              ? "border-[color:var(--adfluence-orange)]/25 bg-[color:var(--adfluence-orange)]/15 text-white"
                              : "border-white/10 bg-[rgba(255,255,255,0.08)] text-[#E2E8F0]",
                          )}
                        >
                          <p>{m.body}</p>
                          {m.media_url ? (
                            <a
                              href={m.media_url}
                              target="_blank"
                              rel="noreferrer"
                              className="mt-2 inline-flex items-center gap-1 rounded-lg border border-white/10 bg-[color:var(--adfluence-surface)] px-2 py-1 text-xs text-[#CBD5E1] hover:underline"
                            >
                              📎 Attachment ({m.media_type})
                            </a>
                          ) : null}
                          <div className="mt-1 flex items-center justify-between gap-2 text-[10px] text-[#94A3B8]">
                            <span>{formatTime(m.created_at)}</span>
                            {mine ? <span>{m.is_read ? "✓✓ Read" : "✓ Sent"}</span> : null}
                          </div>
                          {mine ? (
                            <button
                              type="button"
                              onClick={() => deleteMessage(m.id)}
                              className="absolute -right-2 -top-2 hidden rounded-full bg-rose-500/80 p-1 text-white hover:bg-rose-600 group-hover:block"
                              title="Delete message"
                            >
                              <Trash2 className="size-3" />
                            </button>
                          ) : null}
                        </div>
                      </div>
                    );
                  })
                )}
                <div ref={messagesEndRef} />
              </div>

              <div className="shrink-0 border-t border-white/10 bg-[rgba(255,255,255,0.04)] px-3 py-3 md:px-4">
                <div className="mb-2 flex flex-wrap gap-2">
                  {quickReplies.map((q) => (
                    <button
                      key={q}
                      type="button"
                      onClick={() => sendMessage(q)}
                      className="inline-flex min-h-11 items-center rounded-full border border-white/10 bg-[color:var(--adfluence-surface)] px-4 py-2.5 text-sm text-[#CBD5E1] transition-colors hover:border-[color:var(--adfluence-orange)]/40 hover:text-white md:min-h-0 md:px-3 md:py-1 md:text-xs"
                    >
                      {q}
                    </button>
                  ))}
                </div>
                <div className="flex gap-2">
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
          ) : (
            <div className="flex h-full min-h-[300px] flex-col items-center justify-center p-6 text-center text-sm text-[#94A3B8]">
              Select a conversation to start chatting with creators in real time.
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
