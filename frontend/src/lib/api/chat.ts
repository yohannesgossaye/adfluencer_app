import { apiClient } from "@/lib/api/api-client";
import type { ApiResponse } from "@/lib/api/auth";

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
  receiver_user_id: string;
  body: string;
  media_url?: string;
  media_type?: string;
};

export type MarkReadRequest = {
  conversation_id: string;
};

export type DeleteMessageRequest = {
  delete_type: "me" | "everyone";
};

export type WSIncomingMessage = {
  type: "send_message" | "mark_read" | "delete_message";
  conversation_id: string;
  receiver_user_id?: string;
  body?: string;
  media_url?: string;
  media_type?: string;
  message_id?: string;
  delete_type?: "me" | "everyone";
};

export type WSOutgoingMessage = {
  type: "message_created" | "messages_read" | "message_deleted" | "error";
  payload: any;
};

export type WSReadEvent = {
  conversation_id: string;
  read_by_user_id: string;
  read_at: string;
};

export type WSDeleteEvent = {
  message_id: string;
  conversation_id: string;
  delete_type: string;
  tombstone?: string;
};

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
  listConversations: async (): Promise<ApiResponse<{ items: ChatConversationItem[] }>> => {
    const res = await apiClient.get<ApiResponse<{ items: ChatConversationItem[] }>>(
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
  }): Promise<ApiResponse<{ items: ChatMessage[]; pagination: CursorPaginationMeta }>> => {
    const res = await apiClient.get<ApiResponse<{ items: ChatMessage[]; pagination: CursorPaginationMeta }>>(
      "/api/v1/chat/messages",
      { params }
    );
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
  ): Promise<ApiResponse<{ message_id: string; delete_type: string; tombstone?: string }>> => {
    const res = await apiClient.delete<ApiResponse<{ message_id: string; delete_type: string; tombstone?: string }>>(
      `/api/v1/chat/messages/${messageId}`,
      { data }
    );
    return res.data;
  },

  // 8. WebSocket connection helper
  getWebSocketUrl: (token: string): string => {
    const baseUrl = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";
    const wsProto = baseUrl.startsWith("https") ? "wss:" : "ws:";
    const host = baseUrl.replace(/^https?:\/\//, "");
    return `${wsProto}//${host}/api/v1/chat/ws?token=${encodeURIComponent(token)}`;
  },
};
