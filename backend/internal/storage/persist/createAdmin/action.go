package createadmin

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	customErr "adinfluence_service_app/internal/constant/errors"
	"adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

type createAdminPersistence struct {
	db     *sql.DB
	logger utils.Logger
}

func InitCreateAdminPersistence(db *sql.DB, log utils.Logger) storage.CreateAdminPersistence {
	return &createAdminPersistence{db: db, logger: log}
}

// CreateAdmin adds the user, their admin profile and an audit-log entry in
// one transaction: all three or nothing.
func (p *createAdminPersistence) CreateAdmin(ctx context.Context, in storage.CreateAdminInput) (string, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin create admin: %w", err)
	}
	defer tx.Rollback()

	id := utils.GenerateUUID()
	if _, err := tx.ExecContext(ctx, insertAdminUser, id, in.Email, in.PasswordHash); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return "", customErr.ErrEmailAlreadyExists
		}
		return "", fmt.Errorf("insert admin user: %w", err)
	}
	if _, err := tx.ExecContext(ctx, insertAdminProfile, id, in.Name, in.Level); err != nil {
		return "", fmt.Errorf("insert admin profile: %w", err)
	}
	if _, err := tx.ExecContext(ctx, insertAdminCreatedAudit, utils.GenerateUUID(), id, in.IPAddress, in.UserAgent); err != nil {
		return "", fmt.Errorf("insert admin audit log: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit create admin: %w", err)
	}
	return id, nil
}
