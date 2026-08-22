package handler

import "net/http"

type ChatHandler interface {
	GetOrCreateConversation(w http.ResponseWriter, r *http.Request)
	ListConversations(w http.ResponseWriter, r *http.Request)
	GetConversation(w http.ResponseWriter, r *http.Request)
	SendMessage(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
	MarkMessagesAsRead(w http.ResponseWriter, r *http.Request)
	DeleteMessage(w http.ResponseWriter, r *http.Request)
	WebSocket(w http.ResponseWriter, r *http.Request)
}
