package contract

import (
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/platform/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

type contractRepo struct {
	db     *sql.DB
	logger utils.Logger
}

func InitContractPersistence(db *sql.DB, log utils.Logger) *contractRepo {
	return &contractRepo{db: db, logger: log}
}

func (p *contractRepo) GetAcceptedInvite(ctx context.Context, businessUserID, inviteID string) (*contractDto.AcceptedInviteRecord, error) {
	var rec contractDto.AcceptedInviteRecord
	err := p.db.QueryRowContext(ctx, getAcceptedInvite, inviteID, businessUserID).Scan(
		&rec.ID,
		&rec.CampaignID,
		&rec.BusinessUserID,
		&rec.InfluencerUserID,
		&rec.CampaignStatus,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, contractDto.ErrInviteNotAccepted
		}
		return nil, fmt.Errorf("get accepted invite: %w", err)
	}
	return &rec, nil
}

func (p *contractRepo) ContractExistsForInvite(ctx context.Context, inviteID string) (bool, error) {
	var exists bool
	if err := p.db.QueryRowContext(ctx, contractExistsForInvite, inviteID).Scan(&exists); err != nil {
		return false, fmt.Errorf("contract exists for invite: %w", err)
	}
	return exists, nil
}

func (p *contractRepo) CreateContract(ctx context.Context, input contractDto.SaveContractInput) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin create contract: %w", err)
	}
	defer tx.Rollback()

	fees := input.Fees
	if _, err := tx.ExecContext(ctx, insertContract,
		input.ID,
		input.CampaignID,
		input.BusinessUserID,
		input.InfluencerUserID,
		input.InviteID,
		input.Title,
		nullString(input.Description),
		input.AmountETB,
		fees.BusinessFeePercent,
		fees.InfluencerFeePercent,
		fees.BusinessFeeAmountETB,
		fees.InfluencerFeeAmountETB,
		fees.TotalBusinessPayableETB,
		fees.InfluencerReceivableETB,
		fees.PlatformTotalFeeETB,
		nullDate(input.StartDate),
		nullDate(input.EndDate),
	); err != nil {
		return fmt.Errorf("insert contract: %w", err)
	}

	for _, d := range input.Deliverables {
		if _, err := tx.ExecContext(ctx, insertContractDeliverable,
			utils.GenerateUUID(),
			input.ID,
			strings.ToLower(strings.TrimSpace(d.Type)),
			normalizeQty(d.Quantity),
			nullString(trimOptional(d.Description)),
		); err != nil {
			return fmt.Errorf("insert contract deliverable: %w", err)
		}
	}

	if _, err := tx.ExecContext(ctx, insertContractStatusHistory,
		utils.GenerateUUID(),
		input.ID,
		nil,
		"pending_influencer",
		input.BusinessUserID,
		"contract created",
	); err != nil {
		return fmt.Errorf("insert contract status history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit create contract: %w", err)
	}

	return nil
}

func (p *contractRepo) GetContractForBusiness(ctx context.Context, businessUserID, contractID string) (*contractDto.Record, error) {
	return p.loadContract(ctx, contractID, "business", businessUserID)
}

func (p *contractRepo) GetContractForInfluencer(ctx context.Context, influencerUserID, contractID string) (*contractDto.Record, error) {
	return p.loadContract(ctx, contractID, "influencer", influencerUserID)
}

func (p *contractRepo) loadContract(ctx context.Context, contractID, role, userID string) (*contractDto.Record, error) {
	query := getContractByID
	args := []any{contractID}
	switch role {
	case "business":
		query += " AND c.business_user_id = $2"
		args = append(args, userID)
	case "influencer":
		query += " AND c.influencer_user_id = $2"
		args = append(args, userID)
	}

	var rec contractDto.Record
	var description sql.NullString
	var startDate, endDate sql.NullTime
	var acceptedAt, completedAt sql.NullTime

	err := p.db.QueryRowContext(ctx, query, args...).Scan(
		&rec.ID,
		&rec.InviteID,
		&rec.Campaign.ID,
		&rec.BusinessUserID,
		&rec.InfluencerUserID,
		&rec.Title,
		&description,
		&rec.Fees.AmountETB,
		&rec.Fees.BusinessFeePercent,
		&rec.Fees.InfluencerFeePercent,
		&rec.Fees.BusinessFeeAmountETB,
		&rec.Fees.InfluencerFeeAmountETB,
		&rec.Fees.TotalBusinessPayableETB,
		&rec.Fees.InfluencerReceivableETB,
		&rec.Fees.PlatformTotalFeeETB,
		&startDate,
		&endDate,
		&rec.Status,
		&acceptedAt,
		&completedAt,
		&rec.CreatedAt,
		&rec.UpdatedAt,
		&rec.Campaign.Title,
		&rec.Business.BusinessName,
		&rec.Influencer.CreatorName,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, contractDto.ErrContractNotFound
		}
		return nil, fmt.Errorf("load contract: %w", err)
	}

	rec.Business.UserID = rec.BusinessUserID
	rec.Influencer.UserID = rec.InfluencerUserID
	if description.Valid {
		rec.Description = &description.String
	}
	if startDate.Valid {
		t := startDate.Time
		rec.StartDate = &t
	}
	if endDate.Valid {
		t := endDate.Time
		rec.EndDate = &t
	}
	if acceptedAt.Valid {
		t := acceptedAt.Time
		rec.AcceptedByInfluencerAt = &t
	}
	if completedAt.Valid {
		t := completedAt.Time
		rec.CompletedAt = &t
	}

	var (
		deliverables []contractDto.DeliverableResponse
		payment      *contractDto.PaymentSummary
		payout       *contractDto.PayoutSummary
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		deliverables, err = p.getDeliverables(ctx, rec.ID)
		return err
	})

	g.Go(func() error {
		var err error
		payment, err = p.getPaymentSummary(ctx, rec.ID)
		return err
	})

	g.Go(func() error {
		var err error
		payout, err = p.getPayoutSummary(ctx, rec.ID)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	rec.Deliverables = deliverables
	rec.Payment = payment
	rec.Payout = payout

	return &rec, nil
}

func (p *contractRepo) ListContractsForBusiness(ctx context.Context, businessUserID string, query contractDto.ListQuery) ([]contractDto.ListItem, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit

	rows, err := p.db.QueryContext(ctx, listBusinessContracts, businessUserID, query.Status, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list business contracts: %w", err)
	}
	defer rows.Close()

	var total int
	items, err := scanBusinessContractList(rows, &total)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (p *contractRepo) ListContractsForInfluencer(ctx context.Context, influencerUserID string, query contractDto.ListQuery) ([]contractDto.ListItem, int, error) {
	query = query.Normalize()
	offset := (query.Page - 1) * query.Limit

	rows, err := p.db.QueryContext(ctx, listInfluencerContracts, influencerUserID, query.Status, query.Limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list influencer contracts: %w", err)
	}
	defer rows.Close()

	var total int
	items, err := scanInfluencerContractList(rows, &total)
	if err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (p *contractRepo) AcceptContract(ctx context.Context, influencerUserID, contractID string) error {
	return p.transitionContract(ctx, acceptContract, contractID, influencerUserID, "pending_influencer", "accepted", influencerUserID, "influencer accepted contract")
}

func (p *contractRepo) RejectContract(ctx context.Context, influencerUserID, contractID string) error {
	return p.transitionContract(ctx, rejectContract, contractID, influencerUserID, "pending_influencer", "rejected", influencerUserID, "influencer rejected contract")
}

func (p *contractRepo) CancelContract(ctx context.Context, businessUserID, contractID string) error {
	return p.transitionContract(ctx, cancelBusinessContract, contractID, businessUserID, "", "cancelled", businessUserID, "business cancelled contract")
}

func (p *contractRepo) transitionContract(ctx context.Context, updateSQL, contractID, ownerID, oldStatus, newStatus, changedBy, reason string) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin contract transition: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, updateSQL, contractID, ownerID)
	if err != nil {
		return fmt.Errorf("update contract status: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("update contract status rows: %w", err)
	}
	if n == 0 {
		return contractDto.ErrInvalidContractStatus
	}

	var oldStatusArg any
	if oldStatus == "" {
		oldStatusArg = nil
	} else {
		oldStatusArg = oldStatus
	}

	if _, err := tx.ExecContext(ctx, insertContractStatusHistory,
		utils.GenerateUUID(),
		contractID,
		oldStatusArg,
		newStatus,
		changedBy,
		reason,
	); err != nil {
		return fmt.Errorf("insert contract status history: %w", err)
	}

	return tx.Commit()
}

func (p *contractRepo) StubPayContract(ctx context.Context, businessUserID, contractID, paymentID string) (*contractDto.PaymentSummary, error) {
	rec, err := p.GetContractForBusiness(ctx, businessUserID, contractID)
	if err != nil {
		return nil, err
	}
	if rec.Status != "accepted" {
		return nil, contractDto.ErrInvalidContractStatus
	}

	var paidExists bool
	if err := p.db.QueryRowContext(ctx, paidPaymentExists, contractID).Scan(&paidExists); err != nil {
		return nil, fmt.Errorf("check paid payment: %w", err)
	}
	if paidExists {
		return nil, contractDto.ErrPaymentAlreadyExists
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin stub pay: %w", err)
	}
	defer tx.Rollback()

	fees := rec.Fees
	if _, err := tx.ExecContext(ctx, insertPayment,
		paymentID,
		contractID,
		rec.BusinessUserID,
		rec.InfluencerUserID,
		fees.AmountETB,
		fees.BusinessFeeAmountETB,
		fees.InfluencerFeeAmountETB,
		fees.PlatformTotalFeeETB,
		fees.TotalBusinessPayableETB,
		fees.InfluencerReceivableETB,
		contractDto.StubPaymentMethod,
	); err != nil {
		return nil, fmt.Errorf("insert payment: %w", err)
	}

	res, err := tx.ExecContext(ctx, activateContractAfterPay, contractID, businessUserID)
	if err != nil {
		return nil, fmt.Errorf("activate contract: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("activate contract rows: %w", err)
	}
	if n == 0 {
		return nil, contractDto.ErrInvalidContractStatus
	}

	if _, err := tx.ExecContext(ctx, insertContractStatusHistory,
		utils.GenerateUUID(),
		contractID,
		"accepted",
		"active",
		businessUserID,
		"business paid (escrow held)",
	); err != nil {
		return nil, fmt.Errorf("insert contract status history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit stub pay: %w", err)
	}

	return p.getPaymentSummary(ctx, contractID)
}

func (p *contractRepo) ApproveWork(ctx context.Context, businessUserID, contractID, payoutID string) (*contractDto.PayoutSummary, error) {
	rec, err := p.GetContractForBusiness(ctx, businessUserID, contractID)
	if err != nil {
		return nil, err
	}
	if rec.Status != "active" {
		return nil, contractDto.ErrInvalidContractStatus
	}

	payment, err := p.getPaymentSummary(ctx, contractID)
	if err != nil {
		return nil, err
	}
	if payment == nil || payment.Status != "paid" {
		return nil, contractDto.ErrPaymentRequired
	}

	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin approve work: %w", err)
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, completeContractAfterApprove, contractID, businessUserID)
	if err != nil {
		return nil, fmt.Errorf("complete contract: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("complete contract rows: %w", err)
	}
	if n == 0 {
		return nil, contractDto.ErrInvalidContractStatus
	}

	if _, err := tx.ExecContext(ctx, insertPayout,
		payoutID,
		payment.ID,
		rec.InfluencerUserID,
		rec.Fees.InfluencerReceivableETB,
	); err != nil {
		return nil, fmt.Errorf("insert payout: %w", err)
	}

	if _, err := tx.ExecContext(ctx, insertContractStatusHistory,
		utils.GenerateUUID(),
		contractID,
		"active",
		"completed",
		businessUserID,
		"business approved work — payout queued",
	); err != nil {
		return nil, fmt.Errorf("insert contract status history: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit approve work: %w", err)
	}

	return p.getPayoutSummary(ctx, contractID)
}

func (p *contractRepo) getDeliverables(ctx context.Context, contractID string) ([]contractDto.DeliverableResponse, error) {
	rows, err := p.db.QueryContext(ctx, getContractDeliverables, contractID)
	if err != nil {
		return nil, fmt.Errorf("get contract deliverables: %w", err)
	}
	defer rows.Close()

	items := make([]contractDto.DeliverableResponse, 0)
	for rows.Next() {
		var item contractDto.DeliverableResponse
		var description sql.NullString
		if err := rows.Scan(&item.ID, &item.Type, &item.Quantity, &description, &item.Status); err != nil {
			return nil, fmt.Errorf("scan contract deliverable: %w", err)
		}
		if description.Valid {
			item.Description = &description.String
		}
		items = append(items, item)
	}
	if items == nil {
		items = []contractDto.DeliverableResponse{}
	}
	return items, rows.Err()
}

func (p *contractRepo) getPaymentSummary(ctx context.Context, contractID string) (*contractDto.PaymentSummary, error) {
	var payment contractDto.PaymentSummary
	var paidAt sql.NullTime
	err := p.db.QueryRowContext(ctx, getContractPayment, contractID).Scan(
		&payment.ID,
		&payment.Status,
		&payment.TotalPaidByBusiness,
		&payment.PaymentMethod,
		&paidAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get contract payment: %w", err)
	}
	if paidAt.Valid {
		t := paidAt.Time
		payment.PaidAt = &t
	}
	return &payment, nil
}

func (p *contractRepo) getPayoutSummary(ctx context.Context, contractID string) (*contractDto.PayoutSummary, error) {
	var payout contractDto.PayoutSummary
	var paidAt sql.NullTime
	err := p.db.QueryRowContext(ctx, getContractPayout, contractID).Scan(
		&payout.ID,
		&payout.Status,
		&payout.AmountETB,
		&payout.CreatedAt,
		&paidAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get contract payout: %w", err)
	}
	if paidAt.Valid {
		t := paidAt.Time
		payout.PaidAt = &t
	}
	return &payout, nil
}

func scanBusinessContractList(rows *sql.Rows, total *int) ([]contractDto.ListItem, error) {
	items := make([]contractDto.ListItem, 0)
	for rows.Next() {
		var item contractDto.ListItem
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Status,
			&item.AmountETB,
			&item.CampaignTitle,
			&item.CounterpartyName,
			&item.TotalBusinessPayable,
			&item.InfluencerReceivable,
			&item.CreatedAt,
			total,
		); err != nil {
			return nil, fmt.Errorf("scan business contract list: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []contractDto.ListItem{}
	}
	return items, rows.Err()
}

func scanInfluencerContractList(rows *sql.Rows, total *int) ([]contractDto.ListItem, error) {
	items := make([]contractDto.ListItem, 0)
	for rows.Next() {
		var item contractDto.ListItem
		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Status,
			&item.AmountETB,
			&item.CampaignTitle,
			&item.CounterpartyName,
			&item.TotalBusinessPayable,
			&item.InfluencerReceivable,
			&item.CreatedAt,
			total,
		); err != nil {
			return nil, fmt.Errorf("scan influencer contract list: %w", err)
		}
		items = append(items, item)
	}
	if items == nil {
		items = []contractDto.ListItem{}
	}
	return items, rows.Err()
}

func nullString(v *string) sql.NullString {
	if v == nil || strings.TrimSpace(*v) == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: strings.TrimSpace(*v), Valid: true}
}

func nullDate(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

func normalizeQty(qty int) int {
	if qty <= 0 {
		return 1
	}
	return qty
}

func trimOptional(v string) *string {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return &v
}
