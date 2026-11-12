package initiator

import (
	"adinfluence_service_app/internal/handler"
	adminHTTP "adinfluence_service_app/internal/handler/rest/http/admin"
	auth "adinfluence_service_app/internal/handler/rest/http/auth"
	businessHTTP "adinfluence_service_app/internal/handler/rest/http/business"
	createAdminHTTP "adinfluence_service_app/internal/handler/rest/http/createdAdmin"
	influencerHTTP "adinfluence_service_app/internal/handler/rest/http/influencer"
	chatHTTP "adinfluence_service_app/internal/handler/rest/http/marketplace/chat"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/utils"
)

type handlers struct {
	auth        handler.AuthHandler
	influencer  handler.InfluencerHandler
	admin       handler.AdminHandler
	business    handler.BusinessHandler
	chat        handler.ChatHandler
	marketplace handler.MarketplaceHandler
	createAdmin handler.CreateAdminHandler
}

func InitHandler(mo modules, tokens *jwt.TokenManager, adminCreateApiKey string, allowedOrigins []string, log utils.Logger) handlers {
	return handlers{
		auth:        auth.InitAuthHandler(mo.auth, log),
		influencer:  influencerHTTP.InitInfluencerHandler(mo.influencer, log),
		business:    businessHTTP.InitBusinessHandler(mo.business, log),
		admin:       adminHTTP.InitAdminHandler(mo.admin, log),
		chat:        chatHTTP.InitChatHandler(mo.chat, tokens, allowedOrigins, log),
		createAdmin: createAdminHTTP.InitCreateAdminHandler(mo.createAdmin, adminCreateApiKey, log),
	}
}
