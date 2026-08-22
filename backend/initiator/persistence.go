package initiator

import (
	"adinfluence_service_app/internal/storage"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	auth "adinfluence_service_app/internal/storage/persist/auth"
	businessPersist "adinfluence_service_app/internal/storage/persist/business"
	contractPersist "adinfluence_service_app/internal/storage/persist/contract"
	influencerPersist "adinfluence_service_app/internal/storage/persist/influencer"
	chatPersist "adinfluence_service_app/internal/storage/persist/marketplace/chat"
	supportPersist "adinfluence_service_app/internal/storage/persist/support"
	"adinfluence_service_app/platform/utils"
	"database/sql"
)

type persistence struct {
	auth       storage.AuthPersisitence
	influencer storage.InfluencerPersistence
	business   storage.BusinessPersistence
	contract   storage.ContractPersistence
	admin      storage.AdminPersistence
	support    storage.SupportPersistence
	chat       storage.ChatPersistence
}

func InitPersistence(db *sql.DB, log utils.Logger) persistence {
	return persistence{
		auth:       auth.InitAuthPersistence(db, log),
		influencer: influencerPersist.InitInfluencerPersistence(db, log),
		business:   businessPersist.InitBusinessPersistence(db, log),
		contract:   contractPersist.InitContractPersistence(db, log),
		admin:      adminPersist.InitAdminPersistence(db, log),
		support:    supportPersist.InitSupportPersistence(db),
		chat:       chatPersist.InitChatPersistence(db, log),
	}
}

