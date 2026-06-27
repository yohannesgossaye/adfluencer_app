package initiator

import (
	"adinfluence_service_app/internal/storage"
	auth "adinfluence_service_app/internal/storage/persist/auth"
	adminPersist "adinfluence_service_app/internal/storage/persist/admin"
	businessPersist "adinfluence_service_app/internal/storage/persist/business"
	contractPersist "adinfluence_service_app/internal/storage/persist/contract"
	influencerPersist "adinfluence_service_app/internal/storage/persist/influencer"
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
}

func InitPersistence(db *sql.DB, log utils.Logger) persistence {
	return persistence{
		auth:       auth.InitAuthPersistence(db, log),
		influencer: influencerPersist.InitInfluencerPersistence(db, log),
		business:   businessPersist.InitBusinessPersistence(db, log),
		contract:   contractPersist.InitContractPersistence(db, log),
		admin:      adminPersist.InitAdminPersistence(db, log),
		support:    supportPersist.InitSupportPersistence(db),
	}
}
