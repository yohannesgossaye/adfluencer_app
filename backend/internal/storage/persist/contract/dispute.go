package contract

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/platform/utils"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	selectLatestDispute = `
	SELECT d.id, d.status::text, COALESCE(d.reason, ''), d.description,
		CASE WHEN d.opened_by_user_id = c.business_user_id THEN 'business' ELSE 'influencer' END,
		d.outcome, d.resolution_note, d.created_at, d.resolved_at
	FROM disputes d
	JOIN contracts c ON c.id = d.contract_id
	WHERE d.contract_id = $1
	ORDER BY d.created_at DESC
	LIMIT 1
	`

	lockContractForDispute = `
	SELECT business_user_id, influencer_user_id, status::text
	FROM contracts
	WHERE id = $1 AND (business_user_id = $2 OR influencer_user_id = $2)
	FOR UPDATE
	`

	insertDispute = `
	INSERT INTO disputes (id, contract_id, opened_by_user_id, against_user_id, issue_type, description, severity, status, reason)
	VALUES ($1, $2, $3, $4, 'contract_issue', $5, $6::dispute_severity, 'open', $7)
	`

	setContractStatus = `
	UPDATE contracts SET status = $2::contract_status, updated_at = NOW() WHERE id = $1
	`

	lockDisputeForResolve = `
	SELECT d.contract_id, d.status::text, c.status::text, c.influencer_user_id
	FROM disputes d
	JOIN contracts c ON c.id = d.contract_id
	WHERE d.id = $1
	FOR UPDATE OF d, c
	`

	resolveDispute = `
	UPDATE disputes
	SET status = 'resolved', outcome = $2, resolution_note = $3, resolved_by = $4,
		resolved_at = NOW(), updated_at = NOW()
	WHERE id = $1
	`

	completeDisputedContract = `
	UPDATE contracts SET status = 'completed', completed_at = NOW(), updated_at = NOW() WHERE id = $1
	`

	cancelDisputedContract = `
	UPDATE contracts SET status = 'cancelled', cancelled_at = NOW(), updated_at = NOW() WHERE id = $1
	`

	selectDisputeParties = `
	SELECT c.id, c.title, bp.business_name, ub.email, ip.creator_name, ui.email
	FROM contracts c
	JOIN business_profiles bp ON bp.user_id = c.business_user_id
	JOIN users ub ON ub.id = c.business_user_id
	JOIN influencer_profiles ip ON ip.user_id = c.influencer_user_id
	JOIN users ui ON ui.id = c.influencer_user_id
	WHERE c.id = $1
	`
)

func (p *contractRepo) latestDispute(ctx context.Context, contractID string) (*contractDto.DisputeSummary, error) {
	var d contractDto.DisputeSummary
	var outcome, note sql.NullString
	var resolvedAt sql.NullTime
	err := p.db.QueryRowContext(ctx, selectLatestDispute, contractID).Scan(
		&d.ID, &d.Status, &d.Reason, &d.Description, &d.OpenedByRole, &outcome, &note, &d.CreatedAt, &resolvedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("latest dispute: %w", err)
	}
	d.ReasonLabel = contractDto.DisputeReasonLabel(d.Reason)
	if outcome.Valid {
		d.Outcome = &outcome.String
	}
	if note.Valid {
		d.ResolutionNote = &note.String
	}
	if resolvedAt.Valid {
		t := resolvedAt.Time
		d.ResolvedAt = &t
	}
	return &d, nil
}

// OpenDispute reports a problem on an in-progress contract the user is part
// of. The contract moves to "disputed", so the business can no longer
// approve the work (and release the payout) until an admin decides.
func (p *contractRepo) OpenDispute(ctx context.Context, userID, contractID, reason, description string) (*contractDto.DisputeSummary, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("open dispute: begin: %w", err)
	}
	defer tx.Rollback()

	var businessID, creatorID, status string
	err = tx.QueryRowContext(ctx, lockContractForDispute, contractID, userID).Scan(&businessID, &creatorID, &status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, contractDto.ErrContractNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("open dispute: lock contract: %w", err)
	}
	if status == "disputed" {
		return nil, contractDto.ErrDisputeAlreadyOpen
	}
	if status != "active" {
		return nil, contractDto.ErrInvalidContractStatus
	}
	against := creatorID
	if userID == creatorID {
		against = businessID
	}

	disputeID := utils.GenerateUUID()
	if _, err := tx.ExecContext(ctx, insertDispute, disputeID, contractID, userID, against,
		strings.TrimSpace(description), contractDto.DisputeSeverity(reason), reason); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, contractDto.ErrDisputeAlreadyOpen
		}
		return nil, fmt.Errorf("open dispute: insert: %w", err)
	}
	if _, err := tx.ExecContext(ctx, setContractStatus, contractID, "disputed"); err != nil {
		return nil, fmt.Errorf("open dispute: contract status: %w", err)
	}
	if _, err := tx.ExecContext(ctx, insertContractStatusHistory, utils.GenerateUUID(), contractID, "active", "disputed", userID,
		"problem reported: "+contractDto.DisputeReasonLabel(reason)); err != nil {
		return nil, fmt.Errorf("open dispute: history: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("open dispute: commit: %w", err)
	}
	return p.latestDispute(ctx, contractID)
}

// ResolveDispute settles a dispute and moves its contract accordingly:
// continue → in progress, complete → completed with a payout queued,
// terminate → cancelled (the team refunds the business by hand for now).
func (p *contractRepo) ResolveDispute(ctx context.Context, adminID, disputeID, outcome, note string) (string, error) {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("resolve dispute: begin: %w", err)
	}
	defer tx.Rollback()

	var contractID, disputeStatus, contractStatus, creatorID string
	err = tx.QueryRowContext(ctx, lockDisputeForResolve, disputeID).Scan(&contractID, &disputeStatus, &contractStatus, &creatorID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", contractDto.ErrDisputeNotFound
	}
	if err != nil {
		return "", fmt.Errorf("resolve dispute: lock: %w", err)
	}
	if disputeStatus == "resolved" || disputeStatus == "closed" {
		return "", contractDto.ErrDisputeAlreadyResolved
	}
	note = strings.TrimSpace(note)

	newStatus := contractStatus
	if contractStatus == "disputed" {
		switch outcome {
		case contractDto.DisputeOutcomeContinue:
			newStatus = "active"
			_, err = tx.ExecContext(ctx, setContractStatus, contractID, newStatus)
		case contractDto.DisputeOutcomeComplete:
			newStatus = "completed"
			payment, perr := p.getPaymentSummary(ctx, contractID)
			if perr != nil {
				return "", perr
			}
			if payment == nil || payment.Status != "paid" {
				return "", contractDto.ErrPaymentRequired
			}
			var receivable float64
			if err := tx.QueryRowContext(ctx, `SELECT influencer_receivable_etb FROM contracts WHERE id = $1`, contractID).Scan(&receivable); err != nil {
				return "", fmt.Errorf("resolve dispute: amount: %w", err)
			}
			if _, err = tx.ExecContext(ctx, completeDisputedContract, contractID); err == nil {
				_, err = tx.ExecContext(ctx, insertPayout, utils.GenerateUUID(), payment.ID, creatorID, receivable)
			}
		case contractDto.DisputeOutcomeTerminate:
			newStatus = "cancelled"
			_, err = tx.ExecContext(ctx, cancelDisputedContract, contractID)
		}
		if err != nil {
			return "", fmt.Errorf("resolve dispute: contract: %w", err)
		}
		if _, err := tx.ExecContext(ctx, insertContractStatusHistory, utils.GenerateUUID(), contractID, "disputed", newStatus, adminID,
			"dispute resolved ("+outcome+"): "+note); err != nil {
			return "", fmt.Errorf("resolve dispute: history: %w", err)
		}
	}
	if _, err := tx.ExecContext(ctx, resolveDispute, disputeID, outcome, note, adminID); err != nil {
		return "", fmt.Errorf("resolve dispute: update: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("resolve dispute: commit: %w", err)
	}
	return contractID, nil
}

func (p *contractRepo) DisputeParties(ctx context.Context, contractID string) (*contractDto.DisputeParties, error) {
	var d contractDto.DisputeParties
	err := p.db.QueryRowContext(ctx, selectDisputeParties, contractID).Scan(
		&d.ContractID, &d.ContractTitle, &d.BusinessName, &d.BusinessEmail, &d.CreatorName, &d.CreatorEmail)
	if err != nil {
		return nil, fmt.Errorf("dispute parties: %w", err)
	}
	return &d, nil
}
