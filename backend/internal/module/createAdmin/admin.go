package createadmin

import (
	"context"

	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"adinfluence_service_app/internal/module"
	"adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

type createAdminModule struct {
	admin  storage.CreateAdminPersistence
	logger utils.Logger
}

func InitCreateAdminModule(admin storage.CreateAdminPersistence, log utils.Logger) module.CreateAdminService {
	return &createAdminModule{admin: admin, logger: log}
}

// CreateAdmin stores the admin with a hashed password (the same hashing as
// sign-up, so login works) and records it in the audit log. The request must
// already be normalized and validated.
func (m *createAdminModule) CreateAdmin(ctx context.Context, req adminDto.CreateAdminRequest, meta authDto.SessionMeta) (string, error) {
	hash, err := jwt.HashPassword(req.Password)
	if err != nil {
		return "", err
	}
	id, err := m.admin.CreateAdmin(ctx, storage.CreateAdminInput{
		Name:         req.AdminName,
		Email:        req.AdminEmail,
		PasswordHash: hash,
		Level:        req.AdminRole,
		IPAddress:    meta.IPAddress,
		UserAgent:    meta.UserAgent,
	})
	if err != nil {
		return "", err
	}
	m.logger.Info(ctx, "admin created", zap.String("user_id", id), zap.String("email", req.AdminEmail))
	return id, nil
}
