package initiator

import (
	"adinfluence_service_app/internal/module"
	createadmin "adinfluence_service_app/internal/module/createAdmin"
	adminModule "adinfluence_service_app/internal/module/admin"
	auth "adinfluence_service_app/internal/module/auth"
	businessModule "adinfluence_service_app/internal/module/business"
	disputeModule "adinfluence_service_app/internal/module/dispute"
	influencer "adinfluence_service_app/internal/module/influencer"
	chatModule "adinfluence_service_app/internal/module/marketplace/chat"
	supportModule "adinfluence_service_app/internal/module/support"
	"adinfluence_service_app/platform/jwt"
	"adinfluence_service_app/platform/mail"
	"adinfluence_service_app/platform/notify"
	"adinfluence_service_app/platform/ratelimit"
	"adinfluence_service_app/platform/sms"
	platformStorage "adinfluence_service_app/platform/storage"
	"adinfluence_service_app/platform/utils"
)

type modules struct {
	auth        module.AuthService
	influencer  module.InfluencerService
	business    module.BusinessService
	admin       module.AdminService
	chat        module.ChatService
	createAdmin module.CreateAdminService
}

func InitModules(persistence persistence, tokens *jwt.TokenManager, email *mail.Client, smsClient *sms.Client, dispatcher *notify.Dispatcher, limits ratelimit.Store, files platformStorage.Service, frontendURL, contactInbox string, logSecrets bool, log utils.Logger) modules {
	// One dispute service (and email sender) shared by the business, creator and admin modules.
	emailer := notify.NewEmailer(email, dispatcher, frontendURL)
	disputes := disputeModule.New(persistence.contract, emailer, log)
	threads := supportModule.New(persistence.support, emailer, log)
	return modules{
		auth:        auth.InitAuthModule(persistence.auth, tokens, email, smsClient, dispatcher, limits, frontendURL, contactInbox, logSecrets, log),
		influencer:  influencer.InitInfluencerModule(persistence.influencer, persistence.contract, persistence.support, persistence.auth, files, disputes, threads, log),
		business:    businessModule.InitBusinessModule(persistence.business, persistence.influencer, persistence.contract, persistence.support, persistence.auth, files, disputes, threads, log),
		admin:       adminModule.InitAdminModule(persistence.admin, persistence.influencer, persistence.support, persistence.contract, disputes, threads, log),
		chat:        chatModule.InitChatService(persistence.chat, persistence.auth, log),
		createAdmin: createadmin.InitCreateAdminModule(persistence.createAdmin, log),
	}
}
