package admin

import (
	commonDto "adinfluence_service_app/internal/constant/model/dto/common"
	"strings"
	"time"
)

type ListUsersQuery struct {
	Role   string
	Status string
	Search string
	Page   int
	Limit  int
}

func (q ListUsersQuery) Normalize() ListUsersQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}
	switch q.Role {
	case "", "business", "influencer", "admin":
	default:
		q.Role = ""
	}
	switch q.Status {
	case "", "active", "pending", "blocked", "removed", "suspended", "rejected":
	default:
		q.Status = ""
	}
	q.Search = strings.TrimSpace(q.Search)
	return q
}

type UserStats struct {
	Total   int `json:"total"`
	Active  int `json:"active"`
	Blocked int `json:"blocked"`
	Pending int `json:"pending"`
}

type UserListItem struct {
	ID              string     `json:"user_id"`
	Name            string     `json:"name"`
	Email           string     `json:"email"`
	Role            string     `json:"role"`
	Status          string     `json:"status"`
	JoinedAt        time.Time  `json:"joined_at"`
	LastActiveAt    *time.Time `json:"last_active_at,omitempty"`
	CampaignsCount  int        `json:"campaigns_count"`
	ContractsCount  int        `json:"contracts_count"`
}

type ListUsersResponse struct {
	Items      []UserListItem     `json:"items"`
	Stats      UserStats          `json:"stats"`
	Pagination commonDto.ListMeta `json:"pagination"`
}

type UserDetailResponse struct {
	UserListItem
	StatusReason *string `json:"status_reason,omitempty"`
}

type BlockUserRequest struct {
	Reason string `json:"reason,omitempty"`
}

type BlockUserResponse struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type UnblockUserResponse struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type RemoveUserRequest struct {
	Reason string `json:"reason,omitempty"`
}

type RemoveUserResponse struct {
	UserID string `json:"user_id"`
	Status string `json:"status"`
}

type UserRow struct {
	ID             string
	Email          string
	Role           string
	Status         string
	StatusReason   *string
	DisplayName    string
	CreatedAt      time.Time
	LastLoginAt    *time.Time
	CampaignsCount int
	ContractsCount int
}
