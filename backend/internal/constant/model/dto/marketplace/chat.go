package marketplace

import "time"

type Message struct {
	ID                   string     `json:"id"`
	ConversationID       string     `json:"conversation_id"`
	SenderUserID         string     `json:"sender_user_id"`
	ReceiverUserID       string     `json:"receiver_user_id"`
	Body                 string     `json:"body"`
	MediaURL             *string    `json:"media_url,omitempty"`
	MediaType            string     `json:"media_type"`
	IsRead               bool       `json:"is_read"`
	ReadAt               *time.Time `json:"read_at,omitempty"`
	DeletedForSender     bool       `json:"deleted_for_sender"`
	DeletedForReceiver   bool       `json:"deleted_for_receiver"`
	IsDeletedForEveryone bool       `json:"is_deleted_for_everyone"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

type CursorPaginationMeta struct {
	NextCursor string `json:"next_cursor,omitempty"`
	HasMore    bool   `json:"has_more"`
	Limit      int    `json:"limit"`
}

type ConversationItem struct {
	ID                 string     `json:"id"`
	CampaignID         *string    `json:"campaign_id,omitempty"`
	CampaignTitle      *string    `json:"campaign_title,omitempty"`
	Status             string     `json:"status"`
	LastMessagePreview *string    `json:"last_message_preview,omitempty"`
	LastMessageAt      *time.Time `json:"last_message_at,omitempty"`
	OtherUserID        string     `json:"other_user_id"`
	OtherUserName      string     `json:"other_user_name"`
	OtherUserAvatar    *string    `json:"other_user_avatar,omitempty"`
	OtherUserRole      string     `json:"other_user_role"`
	UnreadCount        int        `json:"unread_count"`
	CreatedAt          time.Time  `json:"created_at"`
}

type CreateConversationRequest struct {
	RecipientUserID string  `json:"recipient_user_id"`
	CampaignID      *string `json:"campaign_id,omitempty"`
}

type SendMessageRequest struct {
	ConversationID string  `json:"conversation_id"`
	ReceiverUserID string  `json:"receiver_user_id"`
	Body           string  `json:"body"`
	MediaURL       *string `json:"media_url,omitempty"`
	MediaType      string  `json:"media_type,omitempty"`
}

type MarkReadRequest struct {
	ConversationID string `json:"conversation_id"`
}

type DeleteMessageRequest struct {
	DeleteType string `json:"delete_type"` // "me" or "everyone"
}

type WSIncomingMessage struct {
	Type           string  `json:"type"` // "send_message", "mark_read", "delete_message"
	ConversationID string  `json:"conversation_id"`
	ReceiverUserID string  `json:"receiver_user_id,omitempty"`
	Body           string  `json:"body,omitempty"`
	MediaURL       *string `json:"media_url,omitempty"`
	MediaType      string  `json:"media_type,omitempty"`
	MessageID      string  `json:"message_id,omitempty"`
	DeleteType     string  `json:"delete_type,omitempty"`
}

type WSOutgoingMessage struct {
	Type    string      `json:"type"` // "message_created", "messages_read", "message_deleted", "error"
	Payload interface{} `json:"payload"`
}

type WSReadEvent struct {
	ConversationID string    `json:"conversation_id"`
	ReadByUserID   string    `json:"read_by_user_id"`
	ReadAt         time.Time `json:"read_at"`
}

type WSDeleteEvent struct {
	MessageID      string `json:"message_id"`
	ConversationID string `json:"conversation_id"`
	DeleteType     string `json:"delete_type"`
	Tombstone      string `json:"tombstone,omitempty"`
}
