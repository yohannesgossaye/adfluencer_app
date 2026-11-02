"use client";

import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState, type ReactNode } from "react";
import axios from "axios";
import { ArrowDown, ArrowLeft, Loader2, RotateCcw, Send, Trash2 } from "lucide-react";
import { toast } from "sonner";

import { useChat } from "@/components/chat/chat-provider";
import { UserAvatar } from "@/components/common/user-avatar";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { chatApi, type ChatConversationItem, type ChatMessage } from "@/lib/api/chat";
import { cn } from "@/lib/utils";

export type MessagesPageConfig = {
  subtitle: string;
  searchPlaceholder: string;
  emptyHint: string;
  noSelection: string;
  quickReplies: readonly string[];
  headerActions?: (conversation: ChatConversationItem) => ReactNode;
};

type DeliveryStatus = "sending" | "failed";
type UiMessage = ChatMessage & { _status?: DeliveryStatus; _error?: string };
type Thread = {
  messages: UiMessage[];
  loaded: boolean;
  hasMore: boolean;
  nextCursor?: string;
  loadingOlder: boolean;
  error?: string;
};

const PAGE_SIZE = 30;
const WS_ACK_TIMEOUT_MS = 8_000;
const NEAR_BOTTOM_PX = 120;
const DELETED_TEXT = "This message was deleted";

const emptyThread: Thread = { messages: [], loaded: false, hasMore: false, loadingOlder: false };

function errorMessage(err: unknown, fallback: string) {
  if (axios.isAxiosError(err)) {
    const msg = (err.response?.data as { message?: string } | undefined)?.message;
    if (msg) return msg;
    if (!err.response) return "You appear to be offline. Check your connection.";
  }
  return fallback;
}

function formatTime(iso: string) {
  try {
    return new Intl.DateTimeFormat("en-GB", { day: "numeric", month: "short", hour: "2-digit", minute: "2-digit" }).format(new Date(iso));
  } catch {
    return iso;
  }
}

function formatListTime(iso?: string) {
  if (!iso) return "";
  try {
    const d = new Date(iso);
    const now = new Date();
    const sameDay = d.toDateString() === now.toDateString();
    return new Intl.DateTimeFormat("en-GB", sameDay ? { hour: "2-digit", minute: "2-digit" } : { day: "numeric", month: "short" }).format(d);
  } catch {
    return "";
  }
}

/** Merges messages by id (a pending message is replaced by its server copy) and sorts oldest→newest. */
function mergeMessages(existing: UiMessage[], incoming: UiMessage[]): UiMessage[] {
  const byKey = new Map<string, UiMessage>();
  for (const m of existing) byKey.set(m.id, m);
  for (const m of incoming) {
    if (m.client_message_id && byKey.has(m.client_message_id) && m.id !== m.client_message_id) {
      byKey.delete(m.client_message_id);
    }
    byKey.set(m.id, m);
  }
  return [...byKey.values()].sort((a, b) => a.created_at.localeCompare(b.created_at) || a.id.localeCompare(b.id));
}

function isMobileViewport() {
  return typeof window !== "undefined" && window.matchMedia("(max-width: 767px)").matches;
}

export function MessagesPage({ config }: { config: MessagesPageConfig }) {
  const chat = useChat();
  const { currentUserId, conversations, conversationsLoaded, conversationsError, connectionState } = chat;
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const linkedConversationId = searchParams.get("conversation_id");

  const [selectedId, setSelectedId] = useState<string>(linkedConversationId ?? "");
  const [mobileShowChat, setMobileShowChat] = useState(false);
  const [search, setSearch] = useState("");
  const [draft, setDraft] = useState("");
  const [threads, setThreads] = useState<Record<string, Thread>>({});
  const [deleteMenuFor, setDeleteMenuFor] = useState<string | null>(null);
  const [showNewBelow, setShowNewBelow] = useState(false);

  const threadsRef = useRef(threads);
  const selectedRef = useRef(selectedId);
  const listRef = useRef<HTMLDivElement | null>(null);
  const pendingScroll = useRef<{ kind: "bottom" | "smooth" } | { kind: "preserve"; height: number; top: number } | null>(null);
  const deepLinkHandled = useRef(false);
  const ackTimers = useRef(new Map<string, ReturnType<typeof setTimeout>>());

  useEffect(() => {
    threadsRef.current = threads;
  }, [threads]);
  useEffect(() => {
    selectedRef.current = selectedId;
  }, [selectedId]);

  const selected = useMemo(() => conversations.find((c) => c.id === selectedId) ?? null, [conversations, selectedId]);
  const thread = (selectedId && threads[selectedId]) || emptyThread;

  const patchThread = useCallback((id: string, patch: (t: Thread) => Thread) => {
    setThreads((prev) => ({ ...prev, [id]: patch(prev[id] ?? emptyThread) }));
  }, []);

  const isNearBottom = () => {
    const el = listRef.current;
    return !el || el.scrollHeight - el.scrollTop - el.clientHeight < NEAR_BOTTOM_PX;
  };

  // ---- loading -----------------------------------------------------------

  const loadLatest = useCallback(
    async (conversationId: string, { merge }: { merge: boolean }) => {
      try {
        const res = await chatApi.getMessages({ conversation_id: conversationId, limit: PAGE_SIZE });
        const page = Array.isArray(res.data) ? res.data : [];
        if (selectedRef.current === conversationId && !merge) pendingScroll.current = { kind: "bottom" };
        patchThread(conversationId, (t) => ({
          ...t,
          messages: merge ? mergeMessages(t.messages, page) : mergeMessages(t.messages.filter((m) => m._status), page),
          loaded: true,
          hasMore: merge && t.loaded ? t.hasMore : Boolean(res.meta?.has_more),
          nextCursor: merge && t.loaded ? t.nextCursor : res.meta?.next_cursor,
          error: undefined,
        }));
      } catch (err) {
        patchThread(conversationId, (t) => ({ ...t, loaded: true, error: errorMessage(err, "Could not load messages.") }));
      }
    },
    [patchThread],
  );

  const loadOlder = useCallback(async () => {
    const id = selectedRef.current;
    const t = threadsRef.current[id];
    if (!id || !t?.hasMore || t.loadingOlder || !t.nextCursor) return;
    patchThread(id, (x) => ({ ...x, loadingOlder: true }));
    try {
      const res = await chatApi.getMessages({ conversation_id: id, cursor: t.nextCursor, limit: PAGE_SIZE });
      const el = listRef.current;
      if (el && selectedRef.current === id) pendingScroll.current = { kind: "preserve", height: el.scrollHeight, top: el.scrollTop };
      patchThread(id, (x) => ({
        ...x,
        messages: mergeMessages(x.messages, Array.isArray(res.data) ? res.data : []),
        hasMore: Boolean(res.meta?.has_more),
        nextCursor: res.meta?.next_cursor,
        loadingOlder: false,
      }));
    } catch (err) {
      patchThread(id, (x) => ({ ...x, loadingOlder: false }));
      toast.error(errorMessage(err, "Could not load older messages."));
    }
  }, [patchThread]);

  // ---- selection ---------------------------------------------------------

  const selectConversation = useCallback(
    (id: string) => {
      setSelectedId(id);
      setDeleteMenuFor(null);
      setShowNewBelow(false);
      if (isMobileViewport()) setMobileShowChat(true);
      router.replace(`${pathname}?conversation_id=${encodeURIComponent(id)}`, { scroll: false });
    },
    [pathname, router],
  );

  // Open the conversation from the link once (e.g. "Message creator"), then let the user move freely.
  useEffect(() => {
    if (deepLinkHandled.current || !conversationsLoaded) return;
    deepLinkHandled.current = true;
    if (!linkedConversationId) return;
    if (isMobileViewport()) setMobileShowChat(true);
    if (!conversations.some((c) => c.id === linkedConversationId)) {
      chatApi
        .getConversation(linkedConversationId)
        .then((res) => res.data && chat.upsertConversation(res.data))
        .catch((err) => toast.error(errorMessage(err, "That conversation could not be opened.")));
    }
  }, [conversationsLoaded, linkedConversationId, conversations, chat]);

  // Desktop: open the most recent conversation when nothing is selected.
  useEffect(() => {
    if (!selectedId && conversations.length > 0 && !isMobileViewport()) {
      setSelectedId(conversations[0].id);
    }
  }, [selectedId, conversations]);

  useEffect(() => {
    chat.setFocusedConversation(selectedId || null);
    return () => chat.setFocusedConversation(null);
  }, [selectedId, chat]);

  useEffect(() => {
    if (selectedId && !threadsRef.current[selectedId]?.loaded) void loadLatest(selectedId, { merge: false });
    else if (selectedId) pendingScroll.current = { kind: "bottom" };
  }, [selectedId, loadLatest]);

  // ---- read receipts -----------------------------------------------------

  const markReadInFlight = useRef(false);
  const markRead = useCallback(
    async (conversationId: string) => {
      if (markReadInFlight.current || document.visibilityState !== "visible") return;
      markReadInFlight.current = true;
      try {
        await chatApi.markMessagesAsRead({ conversation_id: conversationId });
        chat.setUnread(conversationId, 0);
      } catch {
        // Not critical: the badge stays until the next successful read.
      } finally {
        markReadInFlight.current = false;
      }
    },
    [chat],
  );

  const selectedUnread = selected?.unread_count ?? 0;
  useEffect(() => {
    if (selectedId && selectedUnread > 0 && thread.loaded) void markRead(selectedId);
  }, [selectedId, selectedUnread, thread.loaded, markRead]);

  useEffect(() => {
    const onVisible = () => {
      const id = selectedRef.current;
      if (document.visibilityState === "visible" && id) void markRead(id);
    };
    document.addEventListener("visibilitychange", onVisible);
    return () => document.removeEventListener("visibilitychange", onVisible);
  }, [markRead]);

  // ---- realtime events ---------------------------------------------------

  useEffect(() => {
    return chat.subscribe((event) => {
      if (event.type === "resync") {
        const id = selectedRef.current;
        if (id) void loadLatest(id, { merge: true });
        return;
      }
      if (event.type === "message_created") {
        const msg = event.payload;
        if (!threadsRef.current[msg.conversation_id]?.loaded) return;
        const timer = msg.client_message_id && ackTimers.current.get(msg.client_message_id);
        if (timer) {
          clearTimeout(timer);
          ackTimers.current.delete(msg.client_message_id!);
        }
        const viewing = selectedRef.current === msg.conversation_id;
        if (viewing) {
          const mine = msg.sender_user_id === currentUserId;
          if (mine || isNearBottom()) pendingScroll.current = { kind: "smooth" };
          else setShowNewBelow(true);
          if (!mine) void markRead(msg.conversation_id);
        }
        patchThread(msg.conversation_id, (t) => ({ ...t, messages: mergeMessages(t.messages, [msg]) }));
      } else if (event.type === "messages_read" && event.payload.read_by_user_id !== currentUserId) {
        patchThread(event.payload.conversation_id, (t) => ({
          ...t,
          messages: t.messages.map((m) => (m.sender_user_id === currentUserId ? { ...m, is_read: true } : m)),
        }));
      } else if (event.type === "message_deleted") {
        const { conversation_id, message_id, delete_type } = event.payload;
        patchThread(conversation_id, (t) => ({
          ...t,
          messages:
            delete_type === "me"
              ? t.messages.filter((m) => m.id !== message_id)
              : t.messages.map((m) =>
                  m.id === message_id ? { ...m, is_deleted_for_everyone: true, body: DELETED_TEXT, media_url: undefined } : m,
                ),
        }));
      } else if (event.type === "error" && event.payload.client_message_id) {
        const clientId = event.payload.client_message_id;
        const timer = ackTimers.current.get(clientId);
        if (timer) clearTimeout(timer);
        ackTimers.current.delete(clientId);
        markFailedRef.current(clientId, event.payload.message);
      }
    });
  }, [chat, currentUserId, loadLatest, markRead, patchThread]);

  // ---- sending -----------------------------------------------------------

  const markFailed = useCallback(
    (clientId: string, message: string) => {
      setThreads((prev) => {
        const next = { ...prev };
        for (const [id, t] of Object.entries(prev)) {
          if (t.messages.some((m) => m.id === clientId)) {
            next[id] = { ...t, messages: t.messages.map((m) => (m.id === clientId ? { ...m, _status: "failed", _error: message } : m)) };
          }
        }
        return next;
      });
    },
    [],
  );
  const markFailedRef = useRef(markFailed);
  useEffect(() => {
    markFailedRef.current = markFailed;
  }, [markFailed]);

  const sendViaRest = useCallback(
    async (conversationId: string, clientId: string, body: string) => {
      try {
        const res = await chatApi.sendMessage({ conversation_id: conversationId, body, client_message_id: clientId });
        if (res.data) patchThread(conversationId, (t) => ({ ...t, messages: mergeMessages(t.messages, [res.data!]) }));
      } catch (err) {
        markFailed(clientId, errorMessage(err, "Message not sent."));
      }
    },
    [markFailed, patchThread],
  );

  // Prefer the socket; if no acknowledgement arrives in time, the same
  // client_message_id is sent over REST — the server never stores it twice.
  const deliver = useCallback(
    (conversationId: string, clientId: string, body: string) => {
      const sent = chat.send({ type: "send_message", conversation_id: conversationId, body, client_message_id: clientId });
      if (!sent) {
        void sendViaRest(conversationId, clientId, body);
        return;
      }
      ackTimers.current.set(
        clientId,
        setTimeout(() => {
          ackTimers.current.delete(clientId);
          const pending = threadsRef.current[conversationId]?.messages.some((m) => m.id === clientId && m._status === "sending");
          if (pending) void sendViaRest(conversationId, clientId, body);
        }, WS_ACK_TIMEOUT_MS),
      );
    },
    [chat, sendViaRest],
  );

  useEffect(() => {
    const timers = ackTimers.current;
    return () => timers.forEach((t) => clearTimeout(t));
  }, []);

  const sendMessage = useCallback(
    (text: string) => {
      const body = text.trim();
      if (!body || !selected) return;
      const clientId = crypto.randomUUID();
      const now = new Date().toISOString();
      const optimistic: UiMessage = {
        id: clientId,
        client_message_id: clientId,
        conversation_id: selected.id,
        sender_user_id: currentUserId,
        receiver_user_id: selected.other_user_id,
        body,
        media_type: "text",
        is_read: false,
        deleted_for_sender: false,
        deleted_for_receiver: false,
        is_deleted_for_everyone: false,
        created_at: now,
        updated_at: now,
        _status: "sending",
      };
      pendingScroll.current = { kind: "smooth" };
      patchThread(selected.id, (t) => ({ ...t, messages: [...t.messages, optimistic] }));
      setDraft((current) => (current.trim() === body ? "" : current));
      deliver(selected.id, clientId, body);
    },
    [selected, currentUserId, patchThread, deliver],
  );

  const retry = useCallback(
    (m: UiMessage) => {
      patchThread(m.conversation_id, (t) => ({
        ...t,
        messages: t.messages.map((x) => (x.id === m.id ? { ...x, _status: "sending", _error: undefined } : x)),
      }));
      deliver(m.conversation_id, m.id, m.body);
    },
    [patchThread, deliver],
  );

  const deleteMessage = useCallback(
    async (m: UiMessage, deleteType: "me" | "everyone") => {
      setDeleteMenuFor(null);
      try {
        await chatApi.deleteMessage(m.id, { delete_type: deleteType });
        patchThread(m.conversation_id, (t) => ({
          ...t,
          messages:
            deleteType === "me"
              ? t.messages.filter((x) => x.id !== m.id)
              : t.messages.map((x) => (x.id === m.id ? { ...x, is_deleted_for_everyone: true, body: DELETED_TEXT, media_url: undefined } : x)),
        }));
      } catch (err) {
        toast.error(errorMessage(err, "Could not delete the message."));
      }
    },
    [patchThread],
  );

  // ---- scrolling ---------------------------------------------------------

  useLayoutEffect(() => {
    const el = listRef.current;
    const intent = pendingScroll.current;
    if (!el || !intent) return;
    pendingScroll.current = null;
    if (intent.kind === "preserve") {
      el.scrollTop = el.scrollHeight - intent.height + intent.top;
    } else {
      el.scrollTo({ top: el.scrollHeight, behavior: intent.kind === "smooth" ? "smooth" : "auto" });
      setShowNewBelow(false);
    }
  }, [thread.messages, selectedId]);

  const onScroll = () => {
    if (showNewBelow && isNearBottom()) setShowNewBelow(false);
    if (listRef.current && listRef.current.scrollTop < 40) void loadOlder();
  };

  // ---- render ------------------------------------------------------------

  const filteredList = useMemo(() => {
    const q = search.trim().toLowerCase();
    if (!q) return conversations;
    return conversations.filter(
      (c) =>
        c.other_user_name.toLowerCase().includes(q) ||
        c.campaign_title?.toLowerCase().includes(q) ||
        c.last_message_preview?.toLowerCase().includes(q),
    );
  }, [conversations, search]);

  const connectionNotice =
    connectionState === "stopped"
      ? "Chat is not available for your account right now."
      : connectionState === "reconnecting"
        ? "Reconnecting… new messages may be delayed."
        : null;

  return (
    <div className="flex min-h-0 flex-col gap-4 text-white">
      {connectionNotice ? (
        <p role="status" className="rounded-xl border border-amber-400/30 bg-amber-500/10 px-3 py-2 text-xs text-amber-100">
          {connectionNotice}
        </p>
      ) : null}
      <div className="flex min-h-0 flex-1 flex-col gap-4 md:flex-row md:gap-0 md:overflow-hidden md:rounded-2xl md:border md:border-white/10">
        <aside
          className={cn(
            "flex w-full shrink-0 flex-col border-white/10 bg-[rgba(255,255,255,0.08)] md:w-[300px] md:border-r md:py-4",
            mobileShowChat && "hidden md:flex",
          )}
        >
          <div className="border-b border-white/10 px-3 pb-3 md:px-4">
            <h1 className="text-lg font-semibold text-white">Messages</h1>
            <p className="mt-0.5 text-xs text-[#CBD5E1]">{config.subtitle}</p>
            <Input
              id="chat-search"
              aria-label="Search conversations"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder={config.searchPlaceholder}
              className="mt-3 h-10 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-sm text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
            />
          </div>
          <div className="min-h-0 flex-1 space-y-1 overflow-y-auto px-2 py-2 md:px-3">
            {conversationsError ? (
              <div role="alert" className="m-2 rounded-xl border border-rose-400/30 bg-rose-500/10 p-3 text-xs text-rose-100">
                <p>{conversationsError}</p>
                <Button type="button" variant="link" size="sm" className="mt-1 h-auto p-0 text-xs text-white" onClick={() => void chat.refreshConversations()}>
                  Try again
                </Button>
              </div>
            ) : null}
            {!conversationsLoaded ? (
              <p className="flex items-center justify-center gap-2 p-4 text-xs text-[#94A3B8]">
                <Loader2 className="size-3 animate-spin" aria-hidden /> Loading conversations…
              </p>
            ) : filteredList.length === 0 ? (
              search.trim() ? (
                <div className="p-4 text-center">
                  <p className="text-xs text-[#94A3B8]">No conversations match &ldquo;{search}&rdquo;</p>
                  <Button type="button" variant="link" size="sm" className="mt-1 h-auto p-0 text-xs text-[color:var(--adfluence-orange)]" onClick={() => setSearch("")}>
                    Clear search
                  </Button>
                </div>
              ) : (
                <div className="p-4 text-center">
                  <p className="text-xs text-[#94A3B8]">No conversations yet.</p>
                  <p className="mt-1 text-[11px] text-[#64748B]">{config.emptyHint}</p>
                </div>
              )
            ) : (
              filteredList.map((c) => {
                const active = c.id === selectedId;
                return (
                  <button
                    key={c.id}
                    type="button"
                    onClick={() => selectConversation(c.id)}
                    aria-current={active ? "true" : undefined}
                    className={cn(
                      "flex w-full gap-3 rounded-xl border border-transparent p-3 text-left transition-colors focus-visible:outline-2 focus-visible:outline-[color:var(--adfluence-orange)]",
                      active ? "border-white/15 bg-[color:var(--adfluence-orange)]/12" : "hover:border-white/10 hover:bg-white/5",
                    )}
                  >
                    <UserAvatar src={c.other_user_avatar} name={c.other_user_name} size="md" />
                    <div className="min-w-0 flex-1">
                      <div className="flex items-start justify-between gap-2">
                        <p className="truncate text-sm font-medium text-white">{c.other_user_name}</p>
                        <span className="shrink-0 text-[10px] text-[#CBD5E1]">{formatListTime(c.last_message_at)}</span>
                      </div>
                      <p className="truncate text-xs text-[#94A3B8]">{c.campaign_title || "Direct message"}</p>
                      <div className="mt-0.5 flex items-center gap-2">
                        <p className="min-w-0 flex-1 truncate text-xs text-[#CBD5E1]">{c.last_message_preview || "Start the conversation"}</p>
                        {c.unread_count > 0 ? (
                          <span
                            aria-label={`${c.unread_count} unread`}
                            className="inline-flex min-w-[1.25rem] items-center justify-center rounded-full bg-[color:var(--adfluence-orange)] px-1.5 py-0.5 text-[10px] font-semibold text-white"
                          >
                            {c.unread_count}
                          </span>
                        ) : null}
                      </div>
                    </div>
                  </button>
                );
              })
            )}
          </div>
        </aside>

        <section className={cn("relative flex min-h-[70vh] flex-1 flex-col bg-[#0B1A2B] md:min-h-0", !mobileShowChat && "hidden md:flex")}>
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
                <UserAvatar src={selected.other_user_avatar} name={selected.other_user_name} size="md" />
                <div className="min-w-0 flex-1">
                  <p className="font-semibold text-white">{selected.other_user_name}</p>
                  <p className="text-xs text-[#CBD5E1]">
                    {selected.campaign_title || "Direct message"} · <span className="capitalize">{selected.other_user_role}</span>
                  </p>
                </div>
                {config.headerActions ? (
                  <div className="flex w-full shrink-0 flex-wrap gap-2 sm:w-auto sm:justify-end">{config.headerActions(selected)}</div>
                ) : null}
              </header>

              <div
                ref={listRef}
                onScroll={onScroll}
                role="log"
                aria-live="polite"
                aria-label={`Conversation with ${selected.other_user_name}`}
                className="min-h-0 flex-1 space-y-3 overflow-y-auto px-3 py-4 md:px-4"
              >
                {thread.hasMore ? (
                  <div className="flex justify-center">
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      disabled={thread.loadingOlder}
                      onClick={() => void loadOlder()}
                      className="h-8 rounded-full text-xs text-[#CBD5E1] hover:bg-white/10"
                    >
                      {thread.loadingOlder ? <Loader2 className="mr-1 size-3 animate-spin" aria-hidden /> : null}
                      Load older messages
                    </Button>
                  </div>
                ) : null}
                {thread.error ? (
                  <div role="alert" className="mx-auto max-w-sm rounded-xl border border-rose-400/30 bg-rose-500/10 p-3 text-center text-xs text-rose-100">
                    <p>{thread.error}</p>
                    <Button type="button" variant="link" size="sm" className="mt-1 h-auto p-0 text-xs text-white" onClick={() => void loadLatest(selected.id, { merge: false })}>
                      Try again
                    </Button>
                  </div>
                ) : !thread.loaded ? (
                  <p className="flex items-center justify-center gap-2 py-8 text-xs text-[#94A3B8]">
                    <Loader2 className="size-3 animate-spin" aria-hidden /> Loading messages…
                  </p>
                ) : thread.messages.length === 0 ? (
                  <p className="py-8 text-center text-xs text-[#94A3B8]">No messages yet. Send a message below to start collaborating.</p>
                ) : (
                  thread.messages.map((m) => {
                    const mine = m.sender_user_id === currentUserId;
                    const deleted = m.is_deleted_for_everyone;
                    const canAct = !deleted && !m._status;
                    return (
                      <div key={m.id} className={cn("group flex", mine ? "justify-end" : "justify-start")}>
                        <div
                          className={cn(
                            "relative max-w-[85%] rounded-2xl border px-3 py-2 text-sm md:max-w-[70%]",
                            deleted
                              ? "border-white/10 bg-transparent text-[#94A3B8] italic"
                              : mine
                                ? "border-[color:var(--adfluence-orange)]/25 bg-[color:var(--adfluence-orange)]/15 text-white"
                                : "border-white/10 bg-[rgba(255,255,255,0.08)] text-[#E2E8F0]",
                            m._status === "sending" && "opacity-70",
                          )}
                        >
                          <p className="whitespace-pre-wrap break-words">{deleted ? DELETED_TEXT : m.body}</p>
                          {m.media_url && !deleted ? (
                            <a
                              href={m.media_url}
                              target="_blank"
                              rel="noopener noreferrer"
                              className="mt-2 inline-flex items-center gap-1 rounded-lg border border-white/10 bg-[color:var(--adfluence-surface)] px-2 py-1 text-xs text-[#CBD5E1] hover:underline"
                            >
                              Attachment ({m.media_type})
                            </a>
                          ) : null}
                          <div className="mt-1 flex items-center justify-between gap-2 text-[10px] text-[#94A3B8]">
                            <span>{formatTime(m.created_at)}</span>
                            {mine && !deleted ? (
                              m._status === "sending" ? (
                                <span>Sending…</span>
                              ) : m._status === "failed" ? (
                                <button type="button" onClick={() => retry(m)} className="inline-flex items-center gap-1 text-rose-300 hover:underline">
                                  <RotateCcw className="size-3" aria-hidden /> Not sent · Retry
                                </button>
                              ) : (
                                <span>{m.is_read ? "Read" : "Sent"}</span>
                              )
                            ) : null}
                          </div>
                          {m._status === "failed" && m._error ? <p className="mt-1 text-[10px] text-rose-300">{m._error}</p> : null}
                          {canAct ? (
                            <button
                              type="button"
                              aria-label="Delete message"
                              aria-expanded={deleteMenuFor === m.id}
                              onClick={() => setDeleteMenuFor(deleteMenuFor === m.id ? null : m.id)}
                              className="absolute -top-2 -right-2 rounded-full bg-[color:var(--adfluence-surface)] p-1 text-[#CBD5E1] opacity-100 ring-1 ring-white/15 transition-opacity hover:text-white focus-visible:opacity-100 md:opacity-0 md:group-hover:opacity-100 md:group-focus-within:opacity-100"
                            >
                              <Trash2 className="size-3" aria-hidden />
                            </button>
                          ) : null}
                          {deleteMenuFor === m.id ? (
                            <div className="mt-2 flex flex-wrap gap-2 border-t border-white/10 pt-2 text-xs not-italic">
                              <button type="button" className="rounded-lg bg-white/10 px-2 py-1 hover:bg-white/15" onClick={() => void deleteMessage(m, "me")}>
                                Delete for me
                              </button>
                              {mine ? (
                                <button type="button" className="rounded-lg bg-rose-500/80 px-2 py-1 text-white hover:bg-rose-600" onClick={() => void deleteMessage(m, "everyone")}>
                                  Delete for everyone
                                </button>
                              ) : null}
                              <button type="button" className="rounded-lg px-2 py-1 text-[#CBD5E1] hover:bg-white/10" onClick={() => setDeleteMenuFor(null)}>
                                Cancel
                              </button>
                            </div>
                          ) : null}
                        </div>
                      </div>
                    );
                  })
                )}
              </div>

              {showNewBelow ? (
                <button
                  type="button"
                  onClick={() => {
                    pendingScroll.current = { kind: "smooth" };
                    setShowNewBelow(false);
                    listRef.current?.scrollTo({ top: listRef.current.scrollHeight, behavior: "smooth" });
                  }}
                  className="absolute bottom-28 left-1/2 inline-flex -translate-x-1/2 items-center gap-1 rounded-full bg-[color:var(--adfluence-orange)] px-3 py-1.5 text-xs font-medium text-white shadow-lg"
                >
                  <ArrowDown className="size-3" aria-hidden /> New messages
                </button>
              ) : null}

              <div className="shrink-0 border-t border-white/10 bg-[rgba(255,255,255,0.04)] px-3 py-3 md:px-4">
                <div className="mb-2 flex flex-wrap gap-2">
                  {config.quickReplies.map((q) => (
                    <button
                      key={q}
                      type="button"
                      // Fill the box instead of sending, so the text can be checked or edited first.
                      onClick={() => {
                        setDraft(q);
                        document.getElementById("chat-message")?.focus();
                      }}
                      className="inline-flex min-h-11 items-center rounded-full border border-white/10 bg-[color:var(--adfluence-surface)] px-4 py-2.5 text-sm text-[#CBD5E1] transition-colors hover:border-[color:var(--adfluence-orange)]/40 hover:text-white md:min-h-0 md:px-3 md:py-1 md:text-xs"
                    >
                      {q}
                    </button>
                  ))}
                </div>
                <form
                  className="flex gap-2"
                  onSubmit={(e) => {
                    e.preventDefault();
                    sendMessage(draft);
                  }}
                >
                  <Input
                    id="chat-message"
                    aria-label="Message"
                    value={draft}
                    maxLength={4000}
                    onChange={(e) => setDraft(e.target.value)}
                    placeholder="Write a message…"
                    className="h-11 min-w-0 flex-1 rounded-xl border-white/10 bg-[color:var(--adfluence-surface)] text-white placeholder:text-[#94A3B8] focus-visible:border-[color:var(--adfluence-orange)] focus-visible:ring-[color:var(--adfluence-orange)]/25"
                  />
                  <Button
                    type="submit"
                    disabled={!draft.trim()}
                    className="h-11 shrink-0 rounded-xl bg-[color:var(--adfluence-orange)] px-4 text-white hover:bg-[color:var(--adfluence-orange-hover)]"
                  >
                    <Send className="size-4 md:mr-1" aria-hidden />
                    <span className="hidden md:inline">Send</span>
                    <span className="sr-only md:hidden">Send</span>
                  </Button>
                </form>
              </div>
            </>
          ) : (
            <div className="flex h-full min-h-[300px] flex-col items-center justify-center p-6 text-center text-sm text-[#94A3B8]">{config.noSelection}</div>
          )}
        </section>
      </div>
    </div>
  );
}
