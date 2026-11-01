package chat

const (
	previewMaxRunes  = 100
	previewKeepRunes = 97
)

// messagePreview shortens a message body for the conversation list.
// It counts characters (runes), not bytes: Amharic (Ge'ez) letters take
// 3 bytes and emoji 4, so cutting by bytes splits a character and produces
// invalid UTF-8, which Postgres rejects — failing the whole message send.
func messagePreview(body string) string {
	runes := []rune(body)
	if len(runes) <= previewMaxRunes {
		return body
	}
	return string(runes[:previewKeepRunes]) + "..."
}
