package storage

import (
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"context"
)

type ContractPersistence interface {
	GetAcceptedInvite(ctx context.Context, businessUserID, inviteID string) (*contractDto.AcceptedInviteRecord, error)
	ContractExistsForInvite(ctx context.Context, inviteID string) (bool, error)
	CreateContract(ctx context.Context, input contractDto.SaveContractInput) (*contractDto.Record, error)

	GetContractForBusiness(ctx context.Context, businessUserID, contractID string) (*contractDto.Record, error)
	GetContractForInfluencer(ctx context.Context, influencerUserID, contractID string) (*contractDto.Record, error)
	ListContractsForBusiness(ctx context.Context, businessUserID string, query contractDto.ListQuery) ([]contractDto.ListItem, int, error)
	ListContractsForInfluencer(ctx context.Context, influencerUserID string, query contractDto.ListQuery) ([]contractDto.ListItem, int, error)

	AcceptContract(ctx context.Context, influencerUserID, contractID string) error
	RejectContract(ctx context.Context, influencerUserID, contractID string) error
	StubPayContract(ctx context.Context, businessUserID, contractID, paymentID string) (*contractDto.PaymentSummary, error)
	ApproveWork(ctx context.Context, businessUserID, contractID, payoutID string) (*contractDto.PayoutSummary, error)
	CancelContract(ctx context.Context, businessUserID, contractID string) error
}
