package chat

import (
	dto "adinfluence_service_app/internal/constant/model/dto/marketplace"
	"adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type chatPersistence struct {
	db  *sql.DB
	log utils.Logger
}

const (
	defaultConversationLimit = 30
	maxConversationLimit     = 50
	defaultMessageLimit      = 50
	maxMessageLimit          = 100
)

func InitChatPersistence(db *sql.DB, log utils.Logger) storage.ChatPersistence {
	return &chatPersistence{
		db:  db,
		log: log,
	}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanMessage(row rowScanner) (*dto.Message, error) {
	var msg dto.Message
	err := row.Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.SenderUserID,
		&msg.ReceiverUserID,
		&msg.Body,
		&msg.MediaURL,
		&msg.MediaType,
		&msg.IsRead,
		&msg.ReadAt,
		&msg.DeletedForSender,
		&msg.DeletedForReceiver,
		&msg.IsDeletedForEveryone,
		&msg.ClientMessageID,
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func scanConversation(row rowScanner) (*dto.ConversationItem, error) {
	var item dto.ConversationItem
	err := row.Scan(
		&item.ID,
		&item.CampaignID,
		&item.CampaignTitle,
		&item.Status,
		&item.LastMessagePreview,
		&item.LastMessageAt,
		&item.CreatedAt,
		&item.OtherUserID,
		&item.OtherUserRole,
		&item.OtherUserName,
		&item.OtherUserAvatar,
		&item.UnreadCount,
	)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// GetOrCreateConversation returns the conversation between the two users (for
// the campaign, if given), creating it when needed. created reports whether a
// new conversation was made. Safe under concurrent calls: the unique index
// makes the second insert a no-op and both callers get the same row.
func (p *chatPersistence) GetOrCreateConversation(ctx context.Context, userOneID, userTwoID string, campaignID *string) (*dto.ConversationItem, bool, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, fmt.Errorf("begin conversation tx: %w", err)
	}
	defer tx.Rollback()

	created := true
	var convID string
	err = tx.QueryRowContext(ctx, insertConversation, campaignID, userOneID, userTwoID).Scan(&convID)
	if errors.Is(err, sql.ErrNoRows) {
		created = false
		err = tx.QueryRowContext(ctx, findConversationByPair, userOneID, userTwoID, campaignID).Scan(&convID)
	}
	if err != nil {
		return nil, false, fmt.Errorf("get or create conversation: %w", err)
	}

	if created {
		for _, userID := range []string{userOneID, userTwoID} {
			if _, err := tx.ExecContext(ctx, insertParticipant, convID, userID); err != nil {
				return nil, false, fmt.Errorf("insert participant: %w", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("commit conversation: %w", err)
	}

	item, err := p.GetConversationSummary(ctx, convID, userOneID)
	if err != nil {
		return nil, false, err
	}
	return item, created, nil
}

func (p *chatPersistence) IsBusinessVerified(ctx context.Context, userID string) (bool, error) {
	var verified bool
	if err := p.db.QueryRowContext(ctx, isBusinessVerified, userID).Scan(&verified); err != nil {
		return false, fmt.Errorf("check business verified: %w", err)
	}
	return verified, nil
}

func (p *chatPersistence) GetCampaignOwner(ctx context.Context, campaignID string) (string, error) {
	var ownerID string
	err := p.db.QueryRowContext(ctx, getCampaignOwner, campaignID).Scan(&ownerID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", dto.ErrCampaignNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get campaign owner: %w", err)
	}
	return ownerID, nil
}

func (p *chatPersistence) GetConversationSummary(ctx context.Context, conversationID, currentUserID string) (*dto.ConversationItem, error) {
	item, err := scanConversation(p.db.QueryRowContext(ctx, getConversationSummary, currentUserID, conversationID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dto.ErrConversationNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get conversation summary: %w", err)
	}
	return item, nil
}

func (p *chatPersistence) ListConversationsByUserID(ctx context.Context, userID string, limit, offset int) ([]dto.ConversationItem, error) {
	if limit <= 0 {
		limit = defaultConversationLimit
	}
	if limit > maxConversationLimit {
		limit = maxConversationLimit
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := p.db.QueryContext(ctx, listConversationsByUser, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("list conversations: %w", err)
	}
	defer rows.Close()

	items := make([]dto.ConversationItem, 0)
	for rows.Next() {
		item, err := scanConversation(rows)
		if err != nil {
			return nil, fmt.Errorf("scan conversation: %w", err)
		}
		items = append(items, *item)
	}

	return items, rows.Err()
}

func (p *chatPersistence) IsConversationParticipant(ctx context.Context, conversationID, userID string) (bool, error) {
	var exists bool
	if err := p.db.QueryRowContext(ctx, checkParticipantExists, conversationID, userID).Scan(&exists); err != nil {
		return false, fmt.Errorf("check participant: %w", err)
	}
	return exists, nil
}

func (p *chatPersistence) GetOtherParticipantID(ctx context.Context, conversationID, currentUserID string) (string, error) {
	var otherID string
	err := p.db.QueryRowContext(ctx, getOtherParticipant, conversationID, currentUserID).Scan(&otherID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", dto.ErrConversationNotFound
	}
	if err != nil {
		return "", fmt.Errorf("get other participant: %w", err)
	}
	return otherID, nil
}

// InsertMessage stores a message and updates the conversation preview.
// created is false when the sender already sent a message with the same
// client_message_id; the original message is returned unchanged.
func (p *chatPersistence) InsertMessage(ctx context.Context, req dto.SendMessageRequest, senderUserID string) (*dto.Message, bool, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, false, fmt.Errorf("begin message tx: %w", err)
	}
	defer tx.Rollback()

	mediaType := "text"
	if req.MediaType != "" {
		mediaType = req.MediaType
	}

	msg, err := scanMessage(tx.QueryRowContext(ctx, insertMessage,
		req.ConversationID,
		senderUserID,
		req.ReceiverUserID,
		req.Body,
		req.MediaURL,
		mediaType,
		req.ClientMessageID,
	))
	if errors.Is(err, sql.ErrNoRows) {
		existing, err := scanMessage(tx.QueryRowContext(ctx, getMessageByClientID, senderUserID, req.ClientMessageID))
		if err != nil {
			return nil, false, fmt.Errorf("load duplicate message: %w", err)
		}
		return existing, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("insert message: %w", err)
	}

	preview := messagePreview(req.Body)
	if req.MediaURL != nil && *req.MediaURL != "" && req.Body == "" {
		preview = "Attachment: " + mediaType
	}

	if _, err := tx.ExecContext(ctx, updateConversationLastMessage, preview, req.ConversationID); err != nil {
		return nil, false, fmt.Errorf("update conversation last message: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, false, fmt.Errorf("commit message: %w", err)
	}

	return msg, true, nil
}

func (p *chatPersistence) GetMessages(ctx context.Context, conversationID, currentUserID, cursor string, limit int) ([]dto.Message, *dto.CursorPaginationMeta, error) {
	if limit <= 0 {
		limit = defaultMessageLimit
	}
	if limit > maxMessageLimit {
		limit = maxMessageLimit
	}

	var cursorArg sql.NullString
	if cursor != "" {
		cursorArg = sql.NullString{String: cursor, Valid: true}
	}

	// Fetch one extra row to learn whether older messages exist.
	rows, err := p.db.QueryContext(ctx, getMessagesByConversation, conversationID, currentUserID, cursorArg, limit+1)
	if err != nil {
		return nil, nil, fmt.Errorf("get messages: %w", err)
	}
	defer rows.Close()

	messages := make([]dto.Message, 0, limit+1)
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, nil, fmt.Errorf("scan message: %w", err)
		}
		messages = append(messages, *msg)
	}
	if err := rows.Err(); err != nil {
		return nil, nil, fmt.Errorf("iterate messages: %w", err)
	}

	hasMore := len(messages) > limit
	if hasMore {
		// Rows are oldest→newest; the extra row is the oldest one.
		messages = messages[1:]
	}

	meta := &dto.CursorPaginationMeta{HasMore: hasMore, Limit: limit}
	if hasMore && len(messages) > 0 {
		meta.NextCursor = messages[0].ID
	}

	return messages, meta, nil
}

// MarkMessagesAsRead marks every unread message for the user in the
// conversation as read and returns how many changed.
func (p *chatPersistence) MarkMessagesAsRead(ctx context.Context, conversationID, currentUserID string) (int, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("begin mark read tx: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, markMessagesRead, conversationID, currentUserID)
	if err != nil {
		return 0, fmt.Errorf("mark messages read: %w", err)
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return 0, nil
	}

	if _, err := tx.ExecContext(ctx, updateParticipantLastRead, conversationID, currentUserID); err != nil {
		return 0, fmt.Errorf("update last read: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("commit mark read: %w", err)
	}
	return int(n), nil
}

func (p *chatPersistence) DeleteMessage(ctx context.Context, messageID, currentUserID, deleteType string) (*dto.Message, error) {
	msg, err := p.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if msg.SenderUserID != currentUserID && msg.ReceiverUserID != currentUserID {
		return nil, dto.ErrMessageNotFound
	}

	var updateQuery string
	switch {
	case deleteType == "everyone" && msg.SenderUserID != currentUserID:
		return nil, dto.ErrDeleteNotAllowed
	case deleteType == "everyone":
		updateQuery = deleteMessageForEveryone
	case msg.SenderUserID == currentUserID:
		updateQuery = deleteMessageForSender
	default:
		updateQuery = deleteMessageForReceiver
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin delete tx: %w", err)
	}
	defer tx.Rollback()

	updated, err := scanMessage(tx.QueryRowContext(ctx, updateQuery, messageID))
	if err != nil {
		return nil, fmt.Errorf("delete message: %w", err)
	}

	if deleteType == "everyone" {
		if _, err := tx.ExecContext(ctx, refreshConversationPreview, updated.ConversationID); err != nil {
			return nil, fmt.Errorf("refresh conversation preview: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit delete: %w", err)
	}
	return updated, nil
}

func (p *chatPersistence) GetMessageByID(ctx context.Context, messageID string) (*dto.Message, error) {
	msg, err := scanMessage(p.db.QueryRowContext(ctx, getMessageByID, messageID))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, dto.ErrMessageNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get message: %w", err)
	}
	return msg, nil
}
