package influencer

import (
	"adinfluence_service_app/internal/constant/model"
	authDto "adinfluence_service_app/internal/constant/model/dto/auth"
	"context"

	"go.uber.org/zap"
)

func (m *influencerModule) writeAuditAsync(input authDto.CreateAuditLogInput) {
	go func(input authDto.CreateAuditLogInput) {
		defer func() {
			if r := recover(); r != nil {
				m.logger.Error(context.Background(), "audit log panic recovered",
					zap.Any("panic", r),
					zap.String("action", input.Action),
				)
			}
		}()

		auditCtx, cancel := context.WithTimeout(context.Background(), model.AuditAsyncTimeout)
		defer cancel()

		if err := m.influencer.CreateAuditLog(auditCtx, input); err != nil {
			m.logger.Warn(auditCtx, "audit log failed", zap.String("reason", err.Error()))
		}
	}(input)
}
