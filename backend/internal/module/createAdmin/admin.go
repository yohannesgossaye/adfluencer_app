package createadmin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	"adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/utils"
	"context"
)

type createAdminModule struct {
	admin  storage.CreateAdminPersistence
	logger utils.Logger
}

func InitCreateAdminModule(admin storage.CreateAdminPersistence, log utils.Logger) storage.CreateAdminPersistence {
	return &createAdminModule{
		admin:  admin,
		logger: log,
	}
}

func (m *createAdminModule) CreateAdmin(ctx context.Context, req adminDto.CreateAdminRequest) error {
	return m.admin.CreateAdmin(ctx, req)
}
