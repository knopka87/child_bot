package handler

import (
	"log"
	"net/http"

	"child-bot/api/internal/api/middleware"
	"child-bot/api/internal/api/response"
	"child-bot/api/internal/store"
)

// ReferralHandler обрабатывает запросы друзей и рефералов
type ReferralHandler struct {
	store  *store.Store
	appURL string
}

// NewReferralHandler создает новый ReferralHandler
func NewReferralHandler(store *store.Store, appURL string) *ReferralHandler {
	return &ReferralHandler{
		store:  store,
		appURL: appURL,
	}
}

// ReferralData данные реферальной программы
type ReferralData struct {
	ReferralCode     string            `json:"referral_code"`
	ReferralLink     string            `json:"referral_link"`
	TotalInvited     int               `json:"total_invited"`
	ActiveInvited    int               `json:"active_invited"`
	TotalRewards     int               `json:"total_rewards"` // монеты
	InvitedFriends   []InvitedFriend   `json:"invited_friends"`
	RewardMilestones []RewardMilestone `json:"reward_milestones"`
}

// InvitedFriend приглашенный друг
type InvitedFriend struct {
	ID           string `json:"id"`
	DisplayName  string `json:"display_name"`
	AvatarURL    string `json:"avatar_url"`
	JoinedAt     string `json:"joined_at"`
	IsActive     bool   `json:"is_active"`
	RewardEarned int    `json:"reward_earned"`
}

// RewardMilestone вознаграждение за количество приглашенных
type RewardMilestone struct {
	FriendsCount int    `json:"friends_count"`
	Reward       int    `json:"reward"` // монеты
	IsClaimed    bool   `json:"is_claimed"`
	Description  string `json:"description"`
}

// LeaderboardEntry запись в leaderboard
type LeaderboardEntry struct {
	Rank          int    `json:"rank"`
	ChildID       string `json:"child_id"`
	DisplayName   string `json:"display_name"`
	AvatarURL     string `json:"avatar_url"`
	TasksSolved   int    `json:"tasks_solved"`
	Level         int    `json:"level"`
	IsCurrentUser bool   `json:"is_current_user"`
}

// InviteRequest запрос на приглашение
type InviteRequest struct {
	Platform string `json:"platform"` // vk, telegram, link
}

// GetReferralData получает данные реферальной программы
// GET /friends/referrals
func (h *ReferralHandler) GetReferralData(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Получаем реферальный код
	refCode, err := h.store.GetReferralCode(r.Context(), childProfileID)
	if err != nil {
		log.Printf("GetReferralCode error: %v", err)
		response.InternalError(w, "Failed to get referral code")
		return
	}

	// Если кода нет, возвращаем пустые данные
	if refCode == nil {
		response.OK(w, ReferralData{
			ReferralCode:     "",
			ReferralLink:     "",
			TotalInvited:     0,
			ActiveInvited:    0,
			TotalRewards:     0,
			InvitedFriends:   []InvitedFriend{},
			RewardMilestones: []RewardMilestone{},
		})
		return
	}

	// Получаем статистику
	stats, err := h.store.GetReferralStats(r.Context(), childProfileID)
	if err != nil {
		log.Printf("GetReferralStats error: %v", err)
		response.InternalError(w, "Failed to get referral stats")
		return
	}

	// Получаем список друзей
	friendsDB, err := h.store.GetInvitedFriends(r.Context(), childProfileID)
	if err != nil {
		log.Printf("GetInvitedFriends error: %v", err)
		response.InternalError(w, "Failed to get invited friends")
		return
	}

	// Преобразуем друзей в формат API
	friends := make([]InvitedFriend, 0, len(friendsDB))
	for _, f := range friendsDB {
		friends = append(friends, InvitedFriend{
			ID:           f.ID,
			DisplayName:  f.DisplayName,
			AvatarURL:    "/assets/avatars/" + f.AvatarID + ".png",
			JoinedAt:     f.InvitedAt.Format("2006-01-02T15:04:05Z07:00"),
			IsActive:     f.IsActive,
			RewardEarned: f.RewardCoins,
		})
	}

	// Получаем milestone'ы
	milestonesDB, err := h.store.GetRewardMilestones(r.Context(), childProfileID)
	if err != nil {
		log.Printf("GetRewardMilestones error: %v", err)
		response.InternalError(w, "Failed to get reward milestones")
		return
	}

	// Преобразуем milestone'ы в формат API
	milestones := make([]RewardMilestone, 0, len(milestonesDB))
	for _, m := range milestonesDB {
		milestones = append(milestones, RewardMilestone{
			FriendsCount: m.FriendsCount,
			Reward:       m.RewardCoins,
			IsClaimed:    m.IsClaimed,
			Description:  m.Description,
		})
	}

	// Формируем реферальную ссылку
	referralLink := h.appURL + "?ref=" + refCode.Code

	data := ReferralData{
		ReferralCode:     refCode.Code,
		ReferralLink:     referralLink,
		TotalInvited:     stats.TotalInvited,
		ActiveInvited:    stats.ActiveInvited,
		TotalRewards:     stats.TotalRewards,
		InvitedFriends:   friends,
		RewardMilestones: milestones,
	}

	response.OK(w, data)
}

// Invite создает приглашение
// POST /friends/invite
func (h *ReferralHandler) Invite(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - создание приглашения через service layer
	// invite, err := h.service.CreateInvite(r.Context(), childProfileID, platform)

	// Placeholder
	invite := map[string]interface{}{
		"referral_link": "https://app.example.com/join?ref=ABCD1234",
		"share_text":    "Присоединяйся к Объяснятель! Помогу с домашкой 🦉",
		"platform":      "link",
	}

	response.OK(w, invite)
}

// GetLeaderboard получает leaderboard друзей
// GET /friends/leaderboard
func (h *ReferralHandler) GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// Query parameters
	// period := r.URL.Query().Get("period") // week, month, all

	// TODO: Phase 4 - получение leaderboard через service layer
	// leaderboard, err := h.service.GetFriendsLeaderboard(r.Context(), childProfileID, period)

	// Placeholder
	leaderboard := []LeaderboardEntry{
		{
			Rank:          1,
			ChildID:       "child_1",
			DisplayName:   "Ваня",
			AvatarURL:     "/assets/avatars/avatar_2.png",
			TasksSolved:   50,
			Level:         6,
			IsCurrentUser: false,
		},
		{
			Rank:          2,
			ChildID:       childProfileID,
			DisplayName:   "Ты",
			AvatarURL:     "/assets/avatars/avatar_1.png",
			TasksSolved:   42,
			Level:         5,
			IsCurrentUser: true,
		},
	}

	response.OK(w, leaderboard)
}

// ListFriends получает список друзей
// GET /friends
func (h *ReferralHandler) ListFriends(w http.ResponseWriter, r *http.Request) {
	childProfileID := middleware.GetChildProfileID(r.Context())
	if childProfileID == "" {
		response.Unauthorized(w, "Missing child_profile_id")
		return
	}

	// TODO: Phase 4 - получение списка друзей через service layer
	// friends, err := h.service.ListFriends(r.Context(), childProfileID)

	// Placeholder
	friends := []InvitedFriend{}

	response.OK(w, friends)
}
