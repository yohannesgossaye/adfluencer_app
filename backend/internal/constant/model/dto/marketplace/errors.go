package marketplace

import "errors"

// Chat errors. Handlers map them to HTTP statuses; their text is safe to show.
var (
	ErrConversationNotFound = errors.New("conversation not found")
	ErrMessageNotFound      = errors.New("message not found")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrCampaignNotFound     = errors.New("campaign not found")
	ErrNotParticipant       = errors.New("you are not a participant in this conversation")
	ErrChatNotAllowed       = errors.New("chat is only available between businesses and approved creators")
	ErrCampaignNotOwned     = errors.New("the campaign does not belong to this conversation's business")
	ErrAccountCannotChat    = errors.New("your account cannot use chat right now")
	ErrRecipientUnavailable = errors.New("this user is not available for chat")
	ErrDeleteNotAllowed     = errors.New("you can only delete your own messages for everyone")
	ErrSendRateLimited      = errors.New("you are sending messages too quickly, please slow down")
)

// ValidationError describes invalid chat input.
type ValidationError struct{ Msg string }

func (e *ValidationError) Error() string { return e.Msg }

func Invalid(msg string) error { return &ValidationError{Msg: msg} }
