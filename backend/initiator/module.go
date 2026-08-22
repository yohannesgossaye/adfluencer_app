package initiator

import (
	"adinfluence_service_app/internal/module"
	adminModule "adinfluence_service_app/internal/module/admin"
	auth "adinfluence_service_app/internal/module/auth"
	businessModule "adinfluence_service_app/internal/module/business"
	influencer "adinfluence_service_app/internal/module/influencer"
	chatModule "adinfluence_service_app/internal/module/marketplace/chat"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/mail"
	"adinfluence_service_app/platform/sms"
	platformStorage "adinfluence_service_app/platform/storage"
	"adinfluence_service_app/platform/utils"
)

type modules struct {
	auth       module.AuthService
	influencer module.InfluencerService
	business   module.BusinessService
	admin      module.AdminService
	chat       module.ChatService
}

func InitModules(persistence persistence, tokens *jwt.TokenManager, email *mail.Client, smsClient *sms.Client, files platformStorage.Service, frontendURL string, log utils.Logger) modules {
	return modules{
		auth:       auth.InitAuthModule(persistence.auth, tokens, email, smsClient, frontendURL, log),
		influencer: influencer.InitInfluencerModule(persistence.influencer, persistence.contract, persistence.support, persistence.auth, files, log),
		business:   businessModule.InitBusinessModule(persistence.business, persistence.influencer, persistence.contract, persistence.support, persistence.auth, files, log),
		admin:      adminModule.InitAdminModule(persistence.admin, persistence.influencer, persistence.support, log),
		chat:       chatModule.InitChatService(persistence.chat, persistence.auth, log),
	}
}

