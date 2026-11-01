"use client";

import * as React from "react";

import { useAuth } from "@/components/providers/auth-provider";
import { authApi } from "@/lib/api/auth";
import {
  chatApi,
  type ChatConversationItem,
  type WSIncomingMessage,
  type WSOutgoingMessage,
} from "@/lib/api/chat";
import { AUTH_TOKEN_STORAGE_KEY } from "@/lib/constants";
import {
  ChatConnection,
  type ChatConnectionState,
} from "@/lib/chat/connection";

/** While the realtime connection is down, refresh this often (visible tab only). */
const FALLBACK_POLL_MS = 30_000;

type ChatEventListener = (
  event: WSOutgoingMessage | { type: "resync" },
) => void;

type ChatContextValue = {
  currentUserId: string;
  connectionState: ChatConnectionState;
  conversations: ChatConversationItem[];
  conversationsLoaded: boolean;
  conversationsError: string | null;
  totalUnread: number;
  refreshConversations: () => Promise<void>;
  /** Adds or replaces a conversation in the list (e.g. after opening one). */
  upsertConversation: (item: ChatConversationItem) => void;
  /** Marks a conversation's unread count locally (after mark-read). */
  setUnread: (conversationId: string, count: number) => void;
  /** The conversation the user is looking at; its messages are not counted as unread. */
  setFocusedConversation: (conversationId: string | null) => void;
  /** Sends a frame over the socket; false when not connected (use REST instead). */
  send: (frame: WSIncomingMessage) => boolean;
  subscribe: (listener: ChatEventListener) => () => void;
};

const ChatContext = React.createContext<ChatContextValue | null>(null);

export function useChat(): ChatContextValue {
  const ctx = React.useContext(ChatContext);
  if (!ctx) throw new Error("useChat must be used inside <ChatProvider>");
  return ctx;
}

/** Like useChat, but returns null outside a ChatProvider (e.g. the admin portal). */
export function useOptionalChat(): ChatContextValue | null {
  return React.useContext(ChatContext);
}

function sortByActivity(list: ChatConversationItem[]) {
  return [...list].sort((a, b) => {
    const ta = Date.parse(a.last_message_at ?? a.created_at);
    const tb = Date.parse(b.last_message_at ?? b.created_at);
    return tb - ta;
  });
}

export function ChatProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const currentUserId = user?.id ?? "";

  const [connectionState, setConnectionState] =
    React.useState<ChatConnectionState>("idle");
  const [conversations, setConversations] = React.useState<
    ChatConversationItem[]
  >([]);
  const [conversationsLoaded, setConversationsLoaded] = React.useState(false);
  const [conversationsError, setConversationsError] = React.useState<
    string | null
  >(null);

  const listenersRef = React.useRef(new Set<ChatEventListener>());
  const connectionRef = React.useRef<ChatConnection | null>(null);
  const focusedRef = React.useRef<string | null>(null);
  const currentUserRef = React.useRef(currentUserId);
  const conversationsRef = React.useRef(conversations);

  React.useEffect(() => {
    currentUserRef.current = currentUserId;
  }, [currentUserId]);
  React.useEffect(() => {
    conversationsRef.current = conversations;
  }, [conversations]);

  const emit = React.useCallback(
    (event: WSOutgoingMessage | { type: "resync" }) => {
      listenersRef.current.forEach((listener) => listener(event));
    },
    [],
  );

  const refreshConversations = React.useCallback(async () => {
    try {
      const res = await chatApi.listConversations();
      setConversations(sortByActivity(Array.isArray(res.data) ? res.data : []));
      setConversationsError(null);
    } catch {
      setConversationsError(
        "Could not load conversations. Check your connection and try again.",
      );
    } finally {
      setConversationsLoaded(true);
    }
  }, []);

  const upsertConversation = React.useCallback((item: ChatConversationItem) => {
    setConversations((prev) =>
      sortByActivity([item, ...prev.filter((c) => c.id !== item.id)]),
    );
  }, []);

  const setUnread = React.useCallback(
    (conversationId: string, count: number) => {
      setConversations((prev) =>
        prev.map((c) =>
          c.id === conversationId ? { ...c, unread_count: count } : c,
        ),
      );
    },
    [],
  );

  const setFocusedConversation = React.useCallback(
    (conversationId: string | null) => {
      focusedRef.current = conversationId;
    },
    [],
  );

  const handleEvent = React.useCallback(
    (event: WSOutgoingMessage) => {
      const me = currentUserRef.current;
      if (event.type === "message_created") {
        const msg = event.payload;
        if (
          !conversationsRef.current.some((c) => c.id === msg.conversation_id)
        ) {
          // A conversation we have not loaded yet (e.g. someone just wrote first).
          void refreshConversations();
          emit(event);
          return;
        }
        const viewing =
          focusedRef.current === msg.conversation_id &&
          document.visibilityState === "visible";
        setConversations((prev) =>
          sortByActivity(
            prev.map((c) =>
              c.id === msg.conversation_id
                ? {
                    ...c,
                    last_message_preview:
                      msg.body ||
                      (msg.media_url ? `Attachment: ${msg.media_type}` : ""),
                    last_message_at: msg.created_at,
                    unread_count:
                      msg.receiver_user_id === me && !viewing
                        ? c.unread_count + 1
                        : c.unread_count,
                  }
                : c,
            ),
          ),
        );
      } else if (
        event.type === "messages_read" &&
        event.payload.read_by_user_id === me
      ) {
        // Read in another tab or device: clear the badge here too.
        setUnread(event.payload.conversation_id, 0);
      } else if (
        event.type === "message_deleted" &&
        event.payload.delete_type === "everyone"
      ) {
        void refreshConversations();
      }
      emit(event);
    },
    [emit, refreshConversations, setUnread],
  );

  // One connection per signed-in user.
  React.useEffect(() => {
    if (!currentUserId) return;

    const connection = new ChatConnection({
      url: chatApi.getWebSocketUrl(),
      getToken: async (fresh) => {
        if (fresh) {
          // Any authenticated call refreshes an expired access token.
          await authApi.getMe().catch(() => undefined);
        }
        return window.localStorage.getItem(AUTH_TOKEN_STORAGE_KEY);
      },
      onEvent: handleEvent,
      onStateChange: setConnectionState,
      onResync: () => {
        void refreshConversations();
        emit({ type: "resync" });
      },
    });
    connectionRef.current = connection;
    connection.start();
    void Promise.resolve().then(refreshConversations);

    const wake = () => {
      if (document.visibilityState === "visible") connection.nudge();
    };
    document.addEventListener("visibilitychange", wake);
    window.addEventListener("online", wake);

    return () => {
      document.removeEventListener("visibilitychange", wake);
      window.removeEventListener("online", wake);
      connection.stop();
      connectionRef.current = null;
    };
  }, [currentUserId, handleEvent, refreshConversations, emit]);

  // Fallback only while disconnected: no polling when the socket is healthy.
  React.useEffect(() => {
    if (
      !currentUserId ||
      connectionState === "open" ||
      connectionState === "idle"
    )
      return;
    const timer = setInterval(() => {
      if (document.visibilityState !== "visible") return;
      void refreshConversations();
      emit({ type: "resync" });
    }, FALLBACK_POLL_MS);
    return () => clearInterval(timer);
  }, [currentUserId, connectionState, refreshConversations, emit]);

  const send = React.useCallback(
    (frame: WSIncomingMessage) => connectionRef.current?.send(frame) ?? false,
    [],
  );

  const subscribe = React.useCallback((listener: ChatEventListener) => {
    listenersRef.current.add(listener);
    return () => {
      listenersRef.current.delete(listener);
    };
  }, []);

  const totalUnread = React.useMemo(
    () => conversations.reduce((sum, c) => sum + (c.unread_count || 0), 0),
    [conversations],
  );

  const value = React.useMemo<ChatContextValue>(
    () => ({
      currentUserId,
      connectionState,
      conversations,
      conversationsLoaded,
      conversationsError,
      totalUnread,
      refreshConversations,
      upsertConversation,
      setUnread,
      setFocusedConversation,
      send,
      subscribe,
    }),
    [
      currentUserId,
      connectionState,
      conversations,
      conversationsLoaded,
      conversationsError,
      totalUnread,
      refreshConversations,
      upsertConversation,
      setUnread,
      setFocusedConversation,
      send,
      subscribe,
    ],
  );

  return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>;
}
