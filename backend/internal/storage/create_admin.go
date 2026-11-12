package storage

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"context"
)

type CreateAdminPersistence interface {
	CreateAdmin(ctx context.Context, req adminDto.CreateAdminRequest) error
}
