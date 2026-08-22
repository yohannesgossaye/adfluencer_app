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

func InitChatPersistence(db *sql.DB, log utils.Logger) storage.ChatPersistence {
	return &chatPersistence{
		db:  db,
		log: log,
	}
}

func (p *chatPersistence) GetOrCreateConversation(ctx context.Context, userOneID, userTwoID string, campaignID *string) (*dto.ConversationItem, error) {
	var existingID string
	var query string
	var args []interface{}

	if campaignID != nil && *campaignID != "" {
		query = findConversationWithCampaign
		args = []interface{}{userOneID, userTwoID, *campaignID}
	} else {
		query = findConversationWithoutCampaign
		args = []interface{}{userOneID, userTwoID}
	}

	err := p.db.QueryRowContext(ctx, query, args...).Scan(&existingID)
	if err == nil && existingID != "" {
		return p.GetConversationSummary(ctx, existingID, userOneID)
	}

	var userOneRole, userTwoRole string
	err = p.db.QueryRowContext(ctx, getUserRole, userOneID).Scan(&userOneRole)
	if err != nil {
		return nil, fmt.Errorf("user one not found: %w", err)
	}
	err = p.db.QueryRowContext(ctx, getUserRole, userTwoID).Scan(&userTwoRole)
	if err != nil {
		return nil, fmt.Errorf("user two not found: %w", err)
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var newConvID string
	err = tx.QueryRowContext(ctx, insertConversation, campaignID).Scan(&newConvID)
	if err != nil {
		return nil, fmt.Errorf("insert conversation: %w", err)
	}

	if _, err := tx.ExecContext(ctx, insertParticipant, newConvID, userOneID, userOneRole); err != nil {
		return nil, fmt.Errorf("insert participant 1: %w", err)
	}
	if _, err := tx.ExecContext(ctx, insertParticipant, newConvID, userTwoID, userTwoRole); err != nil {
		return nil, fmt.Errorf("insert participant 2: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return p.GetConversationSummary(ctx, newConvID, userOneID)
}

func (p *chatPersistence) GetConversationSummary(ctx context.Context, conversationID, currentUserID string) (*dto.ConversationItem, error) {
	var item dto.ConversationItem
	err := p.db.QueryRowContext(ctx, getConversationSummary, conversationID, currentUserID).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("conversation not found")
		}
		return nil, err
	}

	return &item, nil
}

func (p *chatPersistence) ListConversationsByUserID(ctx context.Context, userID string, limit, offset int) ([]dto.ConversationItem, error) {
	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}

	rows, err := p.db.QueryContext(ctx, listConversationsByUser, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]dto.ConversationItem, 0)
	for rows.Next() {
		var item dto.ConversationItem
		err := rows.Scan(
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
		items = append(items, item)
	}

	return items, rows.Err()
}

func (p *chatPersistence) IsConversationParticipant(ctx context.Context, conversationID, userID string) (bool, error) {
	var exists bool
	err := p.db.QueryRowContext(ctx, checkParticipantExists, conversationID, userID).Scan(&exists)
	return exists, err
}

func (p *chatPersistence) GetOtherParticipantID(ctx context.Context, conversationID, currentUserID string) (string, error) {
	var otherID string
	err := p.db.QueryRowContext(ctx, getOtherParticipant, conversationID, currentUserID).Scan(&otherID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("participant not found")
		}
		return "", err
	}
	return otherID, nil
}

func (p *chatPersistence) InsertMessage(ctx context.Context, req dto.SendMessageRequest, senderUserID string) (*dto.Message, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	mediaType := "text"
	if req.MediaType != "" {
		mediaType = req.MediaType
	}

	var msg dto.Message
	err = tx.QueryRowContext(ctx, insertMessage,
		req.ConversationID,
		senderUserID,
		req.ReceiverUserID,
		req.Body,
		req.MediaURL,
		mediaType,
	).Scan(
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
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("insert message: %w", err)
	}

	preview := req.Body
	if len(preview) > 100 {
		preview = preview[:97] + "..."
	}
	if req.MediaURL != nil && *req.MediaURL != "" && req.Body == "" {
		preview = "Attachment: " + mediaType
	}

	if _, err := tx.ExecContext(ctx, updateConversationLastMessage, preview, req.ConversationID); err != nil {
		return nil, fmt.Errorf("update conversation last message: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (p *chatPersistence) GetMessages(ctx context.Context, conversationID, currentUserID, cursor string, limit int) ([]dto.Message, *dto.CursorPaginationMeta, error) {
	if limit <= 0 || limit > 100 {
		limit = 50
	}

	rows, err := p.db.QueryContext(ctx, getMessagesByConversation, conversationID, currentUserID, cursor, limit)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	messages := make([]dto.Message, 0)
	for rows.Next() {
		var msg dto.Message
		err := rows.Scan(
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
			&msg.CreatedAt,
			&msg.UpdatedAt,
		)
		if err != nil {
			return nil, nil, err
		}
		messages = append(messages, msg)
	}

	if err := rows.Err(); err != nil {
		return nil, nil, err
	}

	var nextCursor string
	hasMore := len(messages) == limit
	if len(messages) > 0 && hasMore {
		nextCursor = messages[0].ID
	}

	meta := &dto.CursorPaginationMeta{
		NextCursor: nextCursor,
		HasMore:    hasMore,
		Limit:      limit,
	}

	return messages, meta, nil
}

func (p *chatPersistence) MarkMessagesAsRead(ctx context.Context, conversationID, currentUserID string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, markMessagesRead, conversationID, currentUserID); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, updateParticipantLastRead, conversationID, currentUserID); err != nil {
		return err
	}

	return tx.Commit()
}

func (p *chatPersistence) DeleteMessage(ctx context.Context, messageID, currentUserID, deleteType string) (*dto.Message, error) {
	msg, err := p.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if msg.SenderUserID != currentUserID && msg.ReceiverUserID != currentUserID {
		return nil, errors.New("unauthorized to delete this message")
	}

	var updateQuery string
	if deleteType == "everyone" {
		if msg.SenderUserID != currentUserID {
			return nil, errors.New("only sender can delete for everyone")
		}
		updateQuery = deleteMessageForEveryone
	} else {
		if msg.SenderUserID == currentUserID {
			updateQuery = deleteMessageForSender
		} else {
			updateQuery = deleteMessageForReceiver
		}
	}

	var updated dto.Message
	err = p.db.QueryRowContext(ctx, updateQuery, messageID).Scan(
		&updated.ID,
		&updated.ConversationID,
		&updated.SenderUserID,
		&updated.ReceiverUserID,
		&updated.Body,
		&updated.MediaURL,
		&updated.MediaType,
		&updated.IsRead,
		&updated.ReadAt,
		&updated.DeletedForSender,
		&updated.DeletedForReceiver,
		&updated.IsDeletedForEveryone,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &updated, nil
}

func (p *chatPersistence) GetMessageByID(ctx context.Context, messageID string) (*dto.Message, error) {
	var msg dto.Message
	err := p.db.QueryRowContext(ctx, getMessageByID, messageID).Scan(
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
		&msg.CreatedAt,
		&msg.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}

	return &msg, nil
}
