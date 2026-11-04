package storage

import (
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"context"
)

type ContractPersistence interface {
	// Disputes ("Report a problem" on a contract).
	OpenDispute(ctx context.Context, userID, contractID, reason, description string) (*contractDto.DisputeSummary, error)
	ResolveDispute(ctx context.Context, adminID, disputeID, outcome, note string) (contractID string, err error)
	DisputeParties(ctx context.Context, contractID string) (*contractDto.DisputeParties, error)
	GetAcceptedInvite(ctx context.Context, businessUserID, inviteID string) (*contractDto.AcceptedInviteRecord, error)
	ContractExistsForInvite(ctx context.Context, inviteID string) (bool, error)
	CreateContract(ctx context.Context, input contractDto.SaveContractInput) error

	GetContractForBusiness(ctx context.Context, businessUserID, contractID string) (*contractDto.Record, error)
	GetContractForInfluencer(ctx context.Context, influencerUserID, contractID string) (*contractDto.Record, error)
	ListContractsForBusiness(ctx context.Context, businessUserID string, query contractDto.ListQuery) ([]contractDto.ListItem, int, error)
	ListContractsForInfluencer(ctx context.Context, influencerUserID string, query contractDto.ListQuery) ([]contractDto.ListItem, int, error)

	AcceptContract(ctx context.Context, influencerUserID, contractID string) error
	RejectContract(ctx context.Context, influencerUserID, contractID string) error
	FundEscrow(ctx context.Context, contractID, paymentID, actorUserID, reference string) (*contractDto.PaymentSummary, error)
	ListContractsAwaitingFunding(ctx context.Context) ([]contractDto.EscrowQueueItem, error)
	ApproveWork(ctx context.Context, businessUserID, contractID, payoutID string) (*contractDto.PayoutSummary, error)
	CancelContract(ctx context.Context, businessUserID, contractID string) error
}
