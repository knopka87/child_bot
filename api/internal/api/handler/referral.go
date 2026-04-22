package handler

import (
	"context"
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
	ReferralCode   string          `json:"referral_code"`
	ReferralLink   string          `json:"referral_link"`
	TotalInvited   int             `json:"total_invited"`
	ActiveInvited  int             `json:"active_invited"`
	TotalRewards   int             `json:"total_rewards"` // монеты (устаревшее, всегда 0)
	InvitedFriends []InvitedFriend `json:"invited_friends"`
	// Данные о текущей цели достижения "Дружба"
	CurrentAchievement *CurrentFriendshipAchievement `json:"current_achievement,omitempty"`
}

// CurrentFriendshipAchievement информация о текущем достижении "Дружба"
type CurrentFriendshipAchievement struct {
	AchievementID string `json:"achievement_id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	Icon          string `json:"icon"`
	TargetCount   int    `json:"target_count"`         // Абсолютная цель (5, 10, 15, 20, ...)
	CurrentCount  int    `json:"current_count"`        // Сколько уже приглашено (абсолютное)
	PreviousLevel int    `json:"previous_level"`       // Предыдущий полученный уровень (0, 5, 10, 15, ...)
	IsUnlocked    bool   `json:"is_unlocked"`          // Получено ли достижение
	NextLevel     *int   `json:"next_level,omitempty"` // Следующий уровень или null если это последний
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
			ReferralCode:   "",
			ReferralLink:   "",
			TotalInvited:   0,
			ActiveInvited:  0,
			TotalRewards:   0,
			InvitedFriends: []InvitedFriend{},
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

	// Формируем реферальную ссылку
	referralLink := h.appURL + "?ref=" + refCode.Code

	// Получаем информацию о текущем достижении "Дружба"
	currentAchievement, err := h.getCurrentFriendshipAchievement(r.Context(), childProfileID, stats.ActiveInvited)
	if err != nil {
		log.Printf("getCurrentFriendshipAchievement error: %v", err)
		// Не критично, продолжаем без achievement
	}

	data := ReferralData{
		ReferralCode:       refCode.Code,
		ReferralLink:       referralLink,
		TotalInvited:       stats.TotalInvited,
		ActiveInvited:      stats.ActiveInvited,
		TotalRewards:       stats.TotalRewards,
		InvitedFriends:     friends,
		CurrentAchievement: currentAchievement,
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

// getCurrentFriendshipAchievement получает информацию о текущем достижении "Дружба"
// Для страницы Friends нужна особая логика:
// - Если текущее достижение уже разблокировано И есть следующий уровень → показываем следующий уровень как цель
// - Иначе показываем текущее
func (h *ReferralHandler) getCurrentFriendshipAchievement(ctx context.Context, childProfileID string, currentFriendsCount int) (*CurrentFriendshipAchievement, error) {
	// Получаем ВСЕ достижения "Дружба" напрямую из базы
	allFriendshipAchievements, err := h.getAllFriendshipAchievements(ctx, childProfileID)
	if err != nil {
		return nil, err
	}

	if len(allFriendshipAchievements) == 0 {
		return nil, nil // Нет достижений "Дружба"
	}

	// Ищем первое неразблокированное достижение (это и есть текущая цель)
	var targetAchievement *store.CombinedAchievement
	var nextLevel *int
	previousLevel := 0 // Предыдущий полученный уровень

	for i, ach := range allFriendshipAchievements {
		if !ach.IsUnlocked {
			// Нашли первую неполученную цель
			targetAchievement = &allFriendshipAchievements[i]

			// Предыдущий уровень - это requirement_value предыдущего достижения или 0
			if i > 0 {
				previousLevel = allFriendshipAchievements[i-1].RequirementValue
			}

			// Проверяем есть ли ещё уровень после этого
			if i+1 < len(allFriendshipAchievements) {
				nextReq := allFriendshipAchievements[i+1].RequirementValue
				nextLevel = &nextReq
			}
			break
		}
	}

	// Если все достижения разблокированы, показываем последнее
	if targetAchievement == nil {
		lastIdx := len(allFriendshipAchievements) - 1
		targetAchievement = &allFriendshipAchievements[lastIdx]
		// Предыдущий уровень - это предпоследнее достижение или само последнее если оно единственное
		if lastIdx > 0 {
			previousLevel = allFriendshipAchievements[lastIdx-1].RequirementValue
		} else {
			previousLevel = targetAchievement.RequirementValue
		}
	}

	result := &CurrentFriendshipAchievement{
		AchievementID: targetAchievement.ID,
		Title:         targetAchievement.Title,
		Description:   targetAchievement.Description,
		Icon:          targetAchievement.Icon,
		TargetCount:   targetAchievement.RequirementValue,
		CurrentCount:  targetAchievement.CurrentProgress,
		PreviousLevel: previousLevel,
		IsUnlocked:    targetAchievement.IsUnlocked,
		NextLevel:     nextLevel,
	}

	return result, nil
}

// getAllFriendshipAchievements получает ВСЕ достижения "Дружба" с прогрессом, отсортированные по requirement_value
func (h *ReferralHandler) getAllFriendshipAchievements(ctx context.Context, childProfileID string) ([]store.CombinedAchievement, error) {
	query := `
		SELECT
			a.id, a.type, a.title, a.description, a.icon,
			a.requirement_type, a.requirement_value,
			a.reward_type, a.reward_id, a.reward_name, a.reward_amount,
			a.priority,
			COALESCE(ca.current_progress, 0) as current_progress,
			COALESCE(ca.is_unlocked, FALSE) as is_unlocked,
			ca.unlocked_at
		FROM achievements a
		LEFT JOIN child_achievements ca
			ON a.id = ca.achievement_id
			AND ca.child_profile_id = $1
		WHERE a.reward_type = 'sticker'
		  AND a.reward_name = 'Дружба'
		ORDER BY a.requirement_value ASC
	`

	rows, err := h.store.DB.QueryContext(ctx, query, childProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []store.CombinedAchievement
	for rows.Next() {
		var ach store.CombinedAchievement
		err := rows.Scan(
			&ach.ID, &ach.Type, &ach.Title, &ach.Description, &ach.Icon,
			&ach.RequirementType, &ach.RequirementValue,
			&ach.RewardType, &ach.RewardID, &ach.RewardName, &ach.RewardAmount,
			&ach.Priority,
			&ach.CurrentProgress, &ach.IsUnlocked, &ach.UnlockedAt,
		)
		if err != nil {
			return nil, err
		}
		achievements = append(achievements, ach)
	}

	return achievements, rows.Err()
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
