package router

import (
	"net/http"

	"child-bot/api/internal/api/handler"
	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/config"
	"child-bot/api/internal/llm"
	"child-bot/api/internal/service"
	"child-bot/api/internal/store"
)

// Dependencies содержит зависимости для handlers
type Dependencies struct {
	Store      *store.Store
	LLMClient  *llm.Client
	Config     *config.Config
	DefaultLLM string
}

// New создает новый router с middleware
func New(deps *Dependencies) http.Handler {
	mux := http.NewServeMux()

	// Health check (public)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		response.OK(w, map[string]string{"status": "ok"})
	})

	// Инициализируем services
	attemptService := service.NewAttemptService(deps.Store, deps.LLMClient, deps.DefaultLLM)
	profileService := service.NewProfileService(deps.Store)
	villainService := service.NewVillainService(deps.Store)
	achievementService := service.NewAchievementService(deps.Store)

	// Устанавливаем зависимости между сервисами (для избежания циклических зависимостей)
	attemptService.SetProfileService(profileService)
	attemptService.SetVillainService(villainService)
	attemptService.SetAchievementService(achievementService)
	profileService.SetAchievementService(achievementService)
	villainService.SetAchievementService(achievementService)

	homeService := service.NewHomeService(deps.Store, attemptService, profileService, villainService)

	// Инициализируем handlers с сервисами
	attemptHandler := handler.NewAttemptHandlerWithService(attemptService, profileService)
	homeHandler := handler.NewHomeHandler(homeService)
	profileHandler := handler.NewProfileHandler(profileService)
	achievementHandler := handler.NewAchievementHandler(deps.Store)
	villainHandler := handler.NewVillainHandler(villainService)
	subscriptionHandler := handler.NewSubscriptionHandler(deps.Store)
	referralHandler := handler.NewReferralHandler(deps.Store, deps.Config.AppURL)
	avatarHandler := handler.NewAvatarHandler()
	consentHandler := handler.NewConsentHandler(deps.Store)
	analyticsHandler := handler.NewAnalyticsHandler()
	legalHandler := handler.NewLegalHandler(deps.Store)
	emailHandler := handler.NewEmailHandler(deps.Store)

	// Регистрация routes
	registerAttemptRoutes(mux, attemptHandler)
	registerHomeRoutes(mux, homeHandler)
	registerProfileRoutes(mux, profileHandler)
	registerAchievementRoutes(mux, achievementHandler)
	registerVillainRoutes(mux, villainHandler)
	registerSubscriptionRoutes(mux, subscriptionHandler)
	registerReferralRoutes(mux, referralHandler)
	registerAvatarRoutes(mux, avatarHandler)
	registerConsentRoutes(mux, consentHandler)
	registerAnalyticsRoutes(mux, analyticsHandler)
	registerLegalRoutes(mux, legalHandler)
	registerEmailRoutes(mux, emailHandler)

	// Применяем middleware в правильном порядке:
	// Recovery -> Logging -> CORS -> Auth
	return middleware.Chain(
		middleware.Recovery,
		middleware.Logging,
		middleware.CORS,
		middleware.Auth,
	)(mux)
}

// registerAttemptRoutes регистрирует routes для attempts
func registerAttemptRoutes(mux *http.ServeMux, h *handler.AttemptHandler) {
	mux.HandleFunc("POST /attempts", h.Create)
	mux.HandleFunc("GET /attempts/unfinished", h.GetUnfinished)
	mux.HandleFunc("GET /attempts/recent", h.GetRecent)
	mux.HandleFunc("POST /attempts/{id}/images", h.UploadImage)
	mux.HandleFunc("POST /attempts/{id}/process", h.Process)
	mux.HandleFunc("GET /attempts/{id}/result", h.GetResult)
	mux.HandleFunc("POST /attempts/{id}/next-hint", h.NextHint)
	mux.HandleFunc("DELETE /attempts/{id}", h.Delete)
}

// registerHomeRoutes регистрирует routes для home
func registerHomeRoutes(mux *http.ServeMux, h *handler.HomeHandler) {
	mux.HandleFunc("GET /home/{childProfileId}", h.GetHomeData)
}

// registerProfileRoutes регистрирует routes для profile
func registerProfileRoutes(mux *http.ServeMux, h *handler.ProfileHandler) {
	mux.HandleFunc("POST /profiles/child", h.CreateChild)
	mux.HandleFunc("GET /profile", h.Get)
	mux.HandleFunc("PUT /profile", h.Update)
	mux.HandleFunc("GET /profile/history", h.GetHistory)
	mux.HandleFunc("GET /profile/stats", h.GetStats)
}

// registerAchievementRoutes регистрирует routes для achievements
func registerAchievementRoutes(mux *http.ServeMux, h *handler.AchievementHandler) {
	mux.HandleFunc("GET /achievements", h.List)
	mux.HandleFunc("GET /achievements/unlocked", h.GetUnlocked)
	mux.HandleFunc("GET /achievements/stats", h.GetStats)
	mux.HandleFunc("GET /achievements/has-new", h.HasNew)
	mux.HandleFunc("POST /achievements/mark-viewed", h.MarkViewed)
	mux.HandleFunc("GET /achievements/{id}", h.GetByID)
	mux.HandleFunc("POST /achievements/{id}/claim", h.Claim)
}

// registerVillainRoutes регистрирует routes для villains
func registerVillainRoutes(mux *http.ServeMux, h *handler.VillainHandler) {
	mux.HandleFunc("GET /villains", h.List)
	mux.HandleFunc("GET /villains/active", h.GetActive)
	mux.HandleFunc("GET /villains/{id}", h.GetByID)
	mux.HandleFunc("GET /villains/{id}/battle", h.GetBattle)
	mux.HandleFunc("GET /villains/{id}/victory", h.GetVictory)
	mux.HandleFunc("POST /villains/{id}/damage", h.DealDamage)
}

// registerSubscriptionRoutes регистрирует routes для subscription
func registerSubscriptionRoutes(mux *http.ServeMux, h *handler.SubscriptionHandler) {
	mux.HandleFunc("GET /subscription/status", h.GetStatus)
	mux.HandleFunc("GET /subscription/plans", h.GetPlans)
	mux.HandleFunc("POST /subscription/subscribe", h.Subscribe)
	mux.HandleFunc("DELETE /subscription/cancel", h.Cancel)
	mux.HandleFunc("POST /subscription/resume", h.Resume)
}

// registerReferralRoutes регистрирует routes для friends/referrals
func registerReferralRoutes(mux *http.ServeMux, h *handler.ReferralHandler) {
	mux.HandleFunc("GET /friends", h.ListFriends)
	mux.HandleFunc("POST /friends/invite", h.Invite)
	mux.HandleFunc("GET /friends/referrals", h.GetReferralData)
	mux.HandleFunc("GET /friends/leaderboard", h.GetLeaderboard)
}

// registerAvatarRoutes регистрирует routes для avatars
func registerAvatarRoutes(mux *http.ServeMux, h *handler.AvatarHandler) {
	mux.HandleFunc("GET /avatars", h.GetAll)
}

// registerConsentRoutes регистрирует routes для consent
func registerConsentRoutes(mux *http.ServeMux, h *handler.ConsentHandler) {
	mux.HandleFunc("POST /consent", h.SaveConsent)
	mux.HandleFunc("GET /consent", h.GetConsent)
	mux.HandleFunc("GET /consent/check", h.CheckConsent)
}

// registerAnalyticsRoutes регистрирует routes для analytics
func registerAnalyticsRoutes(mux *http.ServeMux, h *handler.AnalyticsHandler) {
	mux.HandleFunc("POST /analytics/events", h.SendEvents)
}

// registerLegalRoutes регистрирует routes для legal documents
func registerLegalRoutes(mux *http.ServeMux, h *handler.LegalHandler) {
	mux.HandleFunc("GET /legal/privacy", h.GetPrivacyPolicy)
	mux.HandleFunc("GET /legal/terms", h.GetTermsOfService)
}

// registerEmailRoutes регистрирует routes для email verification
func registerEmailRoutes(mux *http.ServeMux, h *handler.EmailHandler) {
	mux.HandleFunc("POST /email/verify/send", h.SendVerification)
	mux.HandleFunc("POST /email/verify/check", h.VerifyCode)
	mux.HandleFunc("GET /email/verify/status", h.CheckVerification)
}
