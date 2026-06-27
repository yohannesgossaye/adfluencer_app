package initiator

import (
	"adinfluence_service_app/internal/handler"
	auth "adinfluence_service_app/internal/handler/rest/http/auth"
	adminHTTP "adinfluence_service_app/internal/handler/rest/http/admin"
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

func InitHandler(mo modules, logger utils.Logger) handlers {
	return handlers{
		auth:       auth.InitAuthHandler(mo.auth, logger.Named("auth")),
		influencer: influencerHTTP.InitInfluencerHandler(mo.influencer, mo.auth, logger.Named("influencer")),
		business:   businessHTTP.InitBusinessHandler(mo.business, mo.auth, logger.Named("business")),
		admin:      adminHTTP.InitAdminHandler(mo.admin, logger.Named("admin")),
	}
}
