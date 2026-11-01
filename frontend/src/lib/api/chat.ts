import { apiClient } from "@/lib/api/api-client";
import type { ApiResponse } from "@/lib/api/auth";
import { getPublicEnv } from "@/lib/env";

/** Subprotocol offered with the access token when opening the chat socket. */
export const CHAT_WS_SUBPROTOCOL = "adfluence.v1";

export type ChatMessage = {
  id: string;
  conversation_id: string;
  sender_user_id: string;
  receiver_user_id: string;
  body: string;
  media_url?: string;
  media_type: string;
  is_read: boolean;
  read_at?: string;
  deleted_for_sender: boolean;
  deleted_for_receiver: boolean;
  is_deleted_for_everyone: boolean;
  client_message_id?: string;
  created_at: string;
  updated_at: string;
};

export type CursorPaginationMeta = {
  next_cursor?: string;
  has_more: boolean;
  limit: number;
};

export type ChatConversationItem = {
  id: string;
  campaign_id?: string;
  campaign_title?: string;
  status: string;
  last_message_preview?: string;
  last_message_at?: string;
  other_user_id: string;
  other_user_name: string;
  other_user_avatar?: string;
  other_user_role: string;
  unread_count: number;
  created_at: string;
};

export type CreateConversationRequest = {
  recipient_user_id: string;
  campaign_id?: string;
};

export type SendMessageRequest = {
  conversation_id: string;
  body: string;
  media_url?: string;
  media_type?: string;
  /** Client-generated UUID; resending it never creates a duplicate. */
  client_message_id?: string;
};

/** Response envelope of GET /chat/messages: cursor info lives in `meta`. */
export type MessagesPage = ApiResponse<ChatMessage[]> & { meta?: CursorPaginationMeta };

export type MarkReadRequest = {
  conversation_id: string;
};

export type DeleteMessageRequest = {
  delete_type: "me" | "everyone";
};

export type WSIncomingMessage = {
  type: "send_message" | "mark_read" | "delete_message";
  conversation_id?: string;
  client_message_id?: string;
  body?: string;
  media_url?: string;
  media_type?: string;
  message_id?: string;
  delete_type?: "me" | "everyone";
};

export type WSReadEvent = {
  conversation_id: string;
  read_by_user_id: string;
  read_count: number;
  read_at: string;
};

export type WSDeleteEvent = {
  message_id: string;
  conversation_id: string;
  delete_type: string;
  tombstone?: string;
};

export type WSErrorEvent = {
  code: string;
  message: string;
  client_message_id?: string;
};

export type WSOutgoingMessage =
  | { type: "message_created"; payload: ChatMessage }
  | { type: "messages_read"; payload: WSReadEvent }
  | { type: "message_deleted"; payload: WSDeleteEvent }
  | { type: "error"; payload: WSErrorEvent };

export const chatApi = {
  // 1. Start or get direct conversation
  getOrCreateConversation: async (
    data: CreateConversationRequest
  ): Promise<ApiResponse<ChatConversationItem>> => {
    const res = await apiClient.post<ApiResponse<ChatConversationItem>>(
      "/api/v1/chat/conversations",
      data
    );
    return res.data;
  },

  // 2. List my active conversations + previews + unread counts
  listConversations: async (): Promise<ApiResponse<ChatConversationItem[]>> => {
    const res = await apiClient.get<ApiResponse<ChatConversationItem[]>>(
      "/api/v1/chat/conversations"
    );
    return res.data;
  },

  // 3. Get single conversation details
  getConversation: async (
    id: string
  ): Promise<ApiResponse<ChatConversationItem>> => {
    const res = await apiClient.get<ApiResponse<ChatConversationItem>>(
      `/api/v1/chat/conversations/${id}`
    );
    return res.data;
  },

  // 4. Send message via REST
  sendMessage: async (
    data: SendMessageRequest
  ): Promise<ApiResponse<ChatMessage>> => {
    const res = await apiClient.post<ApiResponse<ChatMessage>>(
      "/api/v1/chat/messages",
      data
    );
    return res.data;
  },

  // 5. Cursor-based message history
  getMessages: async (params: {
    conversation_id: string;
    cursor?: string;
    limit?: number;
  }): Promise<MessagesPage> => {
    const res = await apiClient.get<MessagesPage>("/api/v1/chat/messages", { params });
    return res.data;
  },

  // 6. Mark conversation messages as read
  markMessagesAsRead: async (
    data: MarkReadRequest
  ): Promise<ApiResponse<{ conversation_id: string; read_count: number }>> => {
    const res = await apiClient.post<ApiResponse<{ conversation_id: string; read_count: number }>>(
      "/api/v1/chat/messages/read",
      data
    );
    return res.data;
  },

  // 7. Delete message
  deleteMessage: async (
    messageId: string,
    data: DeleteMessageRequest
  ): Promise<ApiResponse<ChatMessage>> => {
    const res = await apiClient.delete<ApiResponse<ChatMessage>>(
      `/api/v1/chat/messages/${messageId}`,
      { params: data }
    );
    return res.data;
  },

  // 8. WebSocket URL. The token is NOT put in the URL (it would end up in
  // server logs); it is sent as a subprotocol — see CHAT_WS_SUBPROTOCOL.
  getWebSocketUrl: (): string => {
    const baseUrl = getPublicEnv().apiBaseUrl || "http://localhost:8080";
    const url = new URL("/api/v1/chat/ws", baseUrl);
    url.protocol = url.protocol === "https:" ? "wss:" : "ws:";
    return url.toString();
  },
};
