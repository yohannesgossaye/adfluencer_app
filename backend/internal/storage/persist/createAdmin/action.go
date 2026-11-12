package createadmin

import (
	adminDto "adinfluence_service_app/internal/constant/model/dto/admin"
	persistent "adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type CreateAdminPersistence struct {
	db     *sql.DB
	logger utils.Logger
}

func InitCreateAdminPersistence(db *sql.DB, log utils.Logger) persistent.CreateAdminPersistence {
	return &CreateAdminPersistence{
		db:     db,
		logger: log,
	}
}

func (p *CreateAdminPersistence) CreateAdmin(ctx context.Context, req adminDto.CreateAdminRequest) error {
	p.logger.Info(ctx, "creating admin user",
		zap.String("name", req.AdminName),
		zap.String("email", req.AdminEmail),
	)

	now := time.Now()
	adminID := uuid.New().String()
	adminProfileID := uuid.New().String()

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		p.logger.Error(ctx, "failed to begin transaction", zap.Error(err))
		return err
	}

	defer func() {
		if rollbackErr := tx.Rollback(); rollbackErr != nil && rollbackErr != sql.ErrTxDone {
			p.logger.Error(ctx, "failed to rollback transaction", zap.Error(rollbackErr))
		}
	}()

	result, err := tx.ExecContext(ctx, createAdmin, adminID, req.AdminEmail,
		nil, req.Password, "admin", "active", true, false, now, now)
	if err != nil {
		p.logger.Error(ctx, "failed to insert admin", zap.Error(err))
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("admin with email %s already exists", req.AdminEmail)
	}

	_, err = tx.ExecContext(ctx, insertAdminProfile, adminProfileID, adminID, req.AdminName, "super_admin", now, now)
	if err != nil {
		p.logger.Error(ctx, "failed to insert admin profile", zap.Error(err))
		return err
	}

	if err = tx.Commit(); err != nil {
		p.logger.Error(ctx, "failed to commit transaction", zap.Error(err))
		return err
	}

	p.logger.Info(ctx, "admin user created successfully",
		zap.String("name", req.AdminName),
		zap.String("email", req.AdminEmail),
	)

	return nil
}
