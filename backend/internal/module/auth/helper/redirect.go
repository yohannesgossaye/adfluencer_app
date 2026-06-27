package helper

import "strings"

// SafeInternalPath returns next only when it is a safe relative in-app path.
func SafeInternalPath(next string) string {
	next = strings.TrimSpace(next)
	if next == "" {
		return ""
	}
	if !strings.HasPrefix(next, "/") || strings.HasPrefix(next, "//") {
		return ""
	}
	if strings.Contains(next, "://") {
		return ""
	}
	return next
}
