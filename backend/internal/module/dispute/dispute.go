// Package dispute is "Report a problem" on a contract, shared by the
// business and creator portals, and the admin's decision.
package dispute

import (
	"context"
	"errors"
	"fmt"
	"strings"

	customErr "adinfluence_service_app/internal/constant/errors"
	contractDto "adinfluence_service_app/internal/constant/model/dto/contract"
	"adinfluence_service_app/internal/storage"
	"adinfluence_service_app/platform/notify"
	"adinfluence_service_app/platform/utils"

	"go.uber.org/zap"
)

type Service struct {
	contracts storage.ContractPersistence
	email     *notify.Emailer
	log       utils.Logger
}

func New(contracts storage.ContractPersistence, email *notify.Emailer, log utils.Logger) *Service {
	return &Service{contracts: contracts, email: email, log: log}
}

// Open reports a problem. role is "business" or "influencer".
func (s *Service) Open(ctx context.Context, userID, role, contractID string, req contractDto.OpenDisputeRequest) (*contractDto.DisputeSummary, error) {
	req.Reason = strings.TrimSpace(req.Reason)
	if err := req.Validate(role); err != nil {
		return nil, err
	}
	d, err := s.contracts.OpenDispute(ctx, userID, contractID, req.Reason, req.Description)
	if err != nil {
		return nil, mapError(err)
	}
	s.log.Info(ctx, "dispute opened", zap.String("contract_id", contractID), zap.String("reason", req.Reason), zap.String("by", role))

	if p, err := s.contracts.DisputeParties(ctx, contractID); err == nil {
		opener, other, otherEmail := p.BusinessName, p.CreatorName, p.CreatorEmail
		openerEmail := p.BusinessEmail
		if role == "influencer" {
			opener, other, otherEmail, openerEmail = p.CreatorName, p.BusinessName, p.BusinessEmail, p.CreatorEmail
		}
		link := s.contractLink(role, contractID)
		s.email.Send(openerEmail, "We received your report: "+p.ContractTitle,
			fmt.Sprintf("Hello %s,\n\nWe received your report on the contract \"%s\" (%s).\n\nThe payment stays held by Adfluence while our team looks into it. We'll contact both sides and email you our decision.\n\n%s\n\nThe Adfluence team\n",
				opener, p.ContractTitle, d.ReasonLabel, link), "dispute_opened")
		s.email.Send(otherEmail, "A problem was reported on: "+p.ContractTitle,
			fmt.Sprintf("Hello %s,\n\n%s reported a problem on the contract \"%s\": %s.\n\nThe contract is paused and the payment stays held while our team reviews it. We may contact you for your side. You'll get an email with our decision.\n\n%s\n\nThe Adfluence team\n",
				other, opener, p.ContractTitle, d.ReasonLabel, s.contractLink(otherRole(role), contractID)), "dispute_opened")
	}
	return d, nil
}

var outcomeText = map[string]string{
	contractDto.DisputeOutcomeContinue:  "The contract continues and is back in progress.",
	contractDto.DisputeOutcomeComplete:  "The work is accepted and the payment is released to the creator.",
	contractDto.DisputeOutcomeTerminate: "The contract is ended. Adfluence will return the payment to the business.",
}

// Resolve records the admin's decision and moves the contract.
func (s *Service) Resolve(ctx context.Context, adminID, disputeID string, req contractDto.ResolveDisputeRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}
	contractID, err := s.contracts.ResolveDispute(ctx, adminID, disputeID, req.Outcome, req.Note)
	if err != nil {
		return mapError(err)
	}
	s.log.Info(ctx, "dispute resolved", zap.String("dispute_id", disputeID), zap.String("outcome", req.Outcome))

	if p, err := s.contracts.DisputeParties(ctx, contractID); err == nil {
		body := func(name, role string) string {
			return fmt.Sprintf("Hello %s,\n\nOur team reviewed the problem reported on \"%s\".\n\nDecision: %s\n\nNote from Adfluence: %s\n\n%s\n\nThe Adfluence team\n",
				name, p.ContractTitle, outcomeText[req.Outcome], strings.TrimSpace(req.Note), s.contractLink(role, contractID))
		}
		s.email.Send(p.BusinessEmail, "Decision on your contract: "+p.ContractTitle, body(p.BusinessName, "business"), "dispute_resolved")
		s.email.Send(p.CreatorEmail, "Decision on your contract: "+p.ContractTitle, body(p.CreatorName, "influencer"), "dispute_resolved")
	}
	return nil
}

func (s *Service) contractLink(role, contractID string) string {
	if s.email == nil || s.email.FrontendURL == "" {
		return ""
	}
	base := strings.TrimRight(s.email.FrontendURL, "/")
	if role == "influencer" {
		return "View the contract: " + base + "/influencer/contracts?contract_id=" + contractID
	}
	return "View the contract: " + base + "/business/contracts?contract_id=" + contractID
}

func otherRole(role string) string {
	if role == "influencer" {
		return "business"
	}
	return "influencer"
}

func mapError(err error) error {
	switch {
	case errors.Is(err, contractDto.ErrContractNotFound):
		return customErr.ErrContractNotFound
	case errors.Is(err, contractDto.ErrInvalidContractStatus):
		return customErr.ErrInvalidContractStatus
	case errors.Is(err, contractDto.ErrDisputeAlreadyOpen):
		return customErr.ErrDisputeAlreadyOpen
	case errors.Is(err, contractDto.ErrDisputeNotFound):
		return customErr.ErrDisputeNotFound
	case errors.Is(err, contractDto.ErrDisputeAlreadyResolved):
		return customErr.ErrDisputeAlreadyResolved
	case errors.Is(err, contractDto.ErrPaymentRequired):
		return customErr.ErrPaymentRequired
	}
	return err
}
