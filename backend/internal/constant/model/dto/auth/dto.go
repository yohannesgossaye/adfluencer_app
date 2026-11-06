package auth

import "time"

type SignUpRequest struct {
	Name            string `json:"name"`
	Phone           string `json:"phone_number"`
	Email           string `json:"email"`
	Role            string `json:"role"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type SignupResponse struct {
	User       UserResponse    `json:"user"`
	Profile    ProfileResponse `json:"profile"`
	Tokens     TokenResponse   `json:"tokens"`
	RedirectTo string          `json:"redirect_to"`
}

type UserResponse struct {
	ID             string     `json:"id"`
	Email          string     `json:"email"`
	PhoneNumber    string     `json:"phone_number,omitempty"`
	Role           string     `json:"role"`
	Status         string     `json:"status"`
	EmailVerified  bool       `json:"email_verified"`
	PhoneVerified  bool       `json:"phone_verified"`
	StatusReason   *string    `json:"status_reason,omitempty"`
	SuspendedUntil *time.Time `json:"suspended_until,omitempty"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	// RememberMe is ignored: every session ends after inactivity (kept so
	// older clients that still send it are accepted).
	RememberMe bool   `json:"remember_me"`
	Portal     string `json:"portal,omitempty"` // "admin" when logging in from /admin/login
	Next       string `json:"next,omitempty"`   // safe internal redirect after login
}

type LoginResponse struct {
	User       UserResponse    `json:"user"`
	Profile    ProfileResponse `json:"profile"`
	Tokens     TokenResponse   `json:"tokens"`
	RedirectTo string          `json:"redirect_to"`
}

type MeResponse struct {
	User    UserResponse    `json:"user"`
	Profile ProfileResponse `json:"profile"`
}

type ProfileResponse struct {
	ID              string `json:"id"`
	BusinessName    string `json:"business_name,omitempty"`
	OwnerFullName   string `json:"owner_full_name,omitempty"`
	FullName        string `json:"full_name,omitempty"`
	CreatorName     string `json:"creator_name,omitempty"`
	PermissionLevel string `json:"permission_level,omitempty"`
	IsVerified      bool   `json:"is_verified"`
}

// UserCredential is used internally for login — password hash must never be returned to clients.
type UserCredential struct {
	UserResponse
	PasswordHash string
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	// The session's limits, so the site can warn before it signs the user out:
	// idle seconds allowed between activity, and when the session ends anyway.
	IdleTimeoutSeconds int64     `json:"idle_timeout_seconds"`
	SessionEndsAt      time.Time `json:"session_ends_at"`
}
