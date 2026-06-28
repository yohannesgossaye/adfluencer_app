package initiator

import (
	"adinfluence_service_app/internal/handler"
	adminHTTP "adinfluence_service_app/internal/handler/rest/http/admin"
	auth "adinfluence_service_app/internal/handler/rest/http/auth"
	businessHTTP "adinfluence_service_app/internal/handler/rest/http/business"
	influencerHTTP "adinfluence_service_app/internal/handler/rest/http/influencer"
	"adinfluence_service_app/platform/utils"
)

type handlers struct {
	auth        handler.AuthHandler
	influencer  handler.InfluencerHandler
	admin       handler.AdminHandler
	business    handler.BusinessHandler
	marketplace handler.MarketplaceHandler
}

func InitHandler(mo modules, log utils.Logger) handlers {
	return handlers{
		auth:       auth.InitAuthHandler(mo.auth, log),
		influencer: influencerHTTP.InitInfluencerHandler(mo.influencer, log),
		business:   businessHTTP.InitBusinessHandler(mo.business, log),
		admin:      adminHTTP.InitAdminHandler(mo.admin, log),
	}
}
