package service

import (
	"context"
	"log"

	"child-bot/api/internal/store"
)

const XPForAchievement = 50

// AchievementService бизнес-логика для достижений
type AchievementService struct {
	store          *store.Store
	profileService *ProfileService
}

// NewAchievementService создает новый AchievementService
func NewAchievementService(store *store.Store) *AchievementService {
	return &AchievementService{store: store}
}

// SetProfileService устанавливает ProfileService (для избежания циклических зависимостей)
func (s *AchievementService) SetProfileService(profileService *ProfileService) {
	s.profileService = profileService
}

// CheckStreakAchievements проверяет достижения за streak (дни подряд)
func (s *AchievementService) CheckStreakAchievements(ctx context.Context, childProfileID string) error {
	streakDays, err := s.store.GetCurrentStreakDays(ctx, childProfileID)
	if err != nil {
		log.Printf("[AchievementService] Failed to get streak days for %s: %v", childProfileID, err)
		return err
	}

	unlockedIDs, err := s.store.CheckAndUpdateAchievementsByType(ctx, childProfileID, "streak_days", streakDays)
	if err != nil {
		log.Printf("[AchievementService] Failed to check streak achievements for %s: %v", childProfileID, err)
		return err
	}

	if len(unlockedIDs) > 0 {
		log.Printf("[AchievementService] 🎉 Unlocked %d streak achievements for child %s: %v",
			len(unlockedIDs), childProfileID, unlockedIDs)

		// Начисляем XP за каждое разблокированное достижение
		if s.profileService != nil {
			for range unlockedIDs {
				if err := s.profileService.AwardAchievementUnlock(ctx, childProfileID); err != nil {
					log.Printf("[AchievementService] Failed to award achievement XP: %v", err)
				}
			}
		}

		// Начисляем награды (монеты)
		if err := s.AwardAchievementRewards(ctx, childProfileID, unlockedIDs); err != nil {
			log.Printf("[AchievementService] Failed to award rewards for streak achievements: %v", err)
		}

		// Стрик награждает стикерами, проверяем коллекционера
		s.checkCollectorAfterUnlock(ctx, childProfileID)
	}

	return nil
}

// CheckVillainAchievements проверяет достижения за побеждённых монстров
func (s *AchievementService) CheckVillainAchievements(ctx context.Context, childProfileID string) error {
	villainsCount, err := s.store.GetVillainsDefeatedCount(ctx, childProfileID)
	if err != nil {
		log.Printf("[AchievementService] Failed to get villains count for %s: %v", childProfileID, err)
		return err
	}

	unlockedIDs, err := s.store.CheckAndUpdateAchievementsByType(ctx, childProfileID, "villains_defeated", villainsCount)
	if err != nil {
		log.Printf("[AchievementService] Failed to check villain achievements for %s: %v", childProfileID, err)
		return err
	}

	if len(unlockedIDs) > 0 {
		log.Printf("[AchievementService] 🎉 Unlocked %d villain achievements for child %s: %v",
			len(unlockedIDs), childProfileID, unlockedIDs)

		// Начисляем награды (монеты)
		if err := s.AwardAchievementRewards(ctx, childProfileID, unlockedIDs); err != nil {
			log.Printf("[AchievementService] Failed to award rewards for villain achievements: %v", err)
		}

		// Злодеи награждают стикерами, проверяем коллекционера
		s.checkCollectorAfterUnlock(ctx, childProfileID)
	}

	return nil
}

// CheckTasksCorrectAchievements проверяет достижения за правильно решённые задачи
func (s *AchievementService) CheckTasksCorrectAchievements(ctx context.Context, childProfileID string) error {
	tasksCount, err := s.store.GetTasksCorrectCount(ctx, childProfileID)
	if err != nil {
		log.Printf("[AchievementService] Failed to get tasks correct count for %s: %v", childProfileID, err)
		return err
	}

	unlockedIDs, err := s.store.CheckAndUpdateAchievementsByType(ctx, childProfileID, "tasks_correct", tasksCount)
	if err != nil {
		log.Printf("[AchievementService] Failed to check tasks correct achievements for %s: %v", childProfileID, err)
		return err
	}

	if len(unlockedIDs) > 0 {
		log.Printf("[AchievementService] 🎉 Unlocked %d tasks correct achievements for child %s: %v",
			len(unlockedIDs), childProfileID, unlockedIDs)

		// Начисляем награды (монеты)
		if err := s.AwardAchievementRewards(ctx, childProfileID, unlockedIDs); err != nil {
			log.Printf("[AchievementService] Failed to award rewards for tasks correct achievements: %v", err)
		}
	}

	return nil
}

// CheckTasksNoHintsAchievements проверяет достижения за задачи без подсказок
func (s *AchievementService) CheckTasksNoHintsAchievements(ctx context.Context, childProfileID string) error {
	tasksCount, err := s.store.GetTasksNoHintsCount(ctx, childProfileID)
	if err != nil {
		log.Printf("[AchievementService] Failed to get tasks no hints count for %s: %v", childProfileID, err)
		return err
	}

	unlockedIDs, err := s.store.CheckAndUpdateAchievementsByType(ctx, childProfileID, "tasks_no_hints", tasksCount)
	if err != nil {
		log.Printf("[AchievementService] Failed to check tasks no hints achievements for %s: %v", childProfileID, err)
		return err
	}

	if len(unlockedIDs) > 0 {
		log.Printf("[AchievementService] 🎉 Unlocked %d tasks no hints achievements for child %s: %v",
			len(unlockedIDs), childProfileID, unlockedIDs)

		// Начисляем награды (монеты)
		if err := s.AwardAchievementRewards(ctx, childProfileID, unlockedIDs); err != nil {
			log.Printf("[AchievementService] Failed to award rewards for tasks no hints achievements: %v", err)
		}
	}

	return nil
}

// CheckFriendsInvitedAchievements проверяет достижения за приглашённых друзей
func (s *AchievementService) CheckFriendsInvitedAchievements(ctx context.Context, childProfileID string) error {
	friendsCount, err := s.store.GetFriendsInvitedCount(ctx, childProfileID)
	if err != nil {
		log.Printf("[AchievementService] Failed to get friends invited count for %s: %v", childProfileID, err)
		return err
	}

	unlockedIDs, err := s.store.CheckAndUpdateAchievementsByType(ctx, childProfileID, "friends_invited", friendsCount)
	if err != nil {
		log.Printf("[AchievementService] Failed to check friends invited achievements for %s: %v", childProfileID, err)
		return err
	}

	if len(unlockedIDs) > 0 {
		log.Printf("[AchievementService] 🎉 Unlocked %d friends invited achievements for child %s: %v",
			len(unlockedIDs), childProfileID, unlockedIDs)

		// Начисляем награды (монеты)
		if err := s.AwardAchievementRewards(ctx, childProfileID, unlockedIDs); err != nil {
			log.Printf("[AchievementService] Failed to award rewards for friends invited achievements: %v", err)
		}

		// Дружба награждает стикерами, проверяем коллекционера
		s.checkCollectorAfterUnlock(ctx, childProfileID)
	}

	return nil
}

// CheckHintsUsedAchievements проверяет достижения за использованные подсказки (Мудрая сова)
func (s *AchievementService) CheckHintsUsedAchievements(ctx context.Context, childProfileID string) error {
	hintsCount, err := s.store.GetHintsUsedCount(ctx, childProfileID)
	if err != nil {
		log.Printf("[AchievementService] Failed to get hints used count for %s: %v", childProfileID, err)
		return err
	}

	unlockedIDs, err := s.store.CheckAndUpdateAchievementsByType(ctx, childProfileID, "hints_used", hintsCount)
	if err != nil {
		log.Printf("[AchievementService] Failed to check hints used achievements for %s: %v", childProfileID, err)
		return err
	}

	if len(unlockedIDs) > 0 {
		log.Printf("[AchievementService] 🦉 Unlocked hints used achievement for child %s: %v",
			childProfileID, unlockedIDs)

		// Начисляем награды (монеты)
		if err := s.AwardAchievementRewards(ctx, childProfileID, unlockedIDs); err != nil {
			log.Printf("[AchievementService] Failed to award rewards for hints used achievements: %v", err)
		}
	}

	return nil
}

// CheckStickersCollectedAchievements проверяет достижения за собранные стикеры (Коллекционер)
func (s *AchievementService) CheckStickersCollectedAchievements(ctx context.Context, childProfileID string) error {
	stickersCount, err := s.store.GetStickersCollectedCount(ctx, childProfileID)
	if err != nil {
		log.Printf("[AchievementService] Failed to get stickers collected count for %s: %v", childProfileID, err)
		return err
	}

	unlockedIDs, err := s.store.CheckAndUpdateAchievementsByType(ctx, childProfileID, "stickers_collected", stickersCount)
	if err != nil {
		log.Printf("[AchievementService] Failed to check stickers collected achievements for %s: %v", childProfileID, err)
		return err
	}

	if len(unlockedIDs) > 0 {
		log.Printf("[AchievementService] 🎨 Unlocked stickers collected achievement for child %s: %v",
			childProfileID, unlockedIDs)

		// Начисляем награды (монеты)
		if err := s.AwardAchievementRewards(ctx, childProfileID, unlockedIDs); err != nil {
			log.Printf("[AchievementService] Failed to award rewards for stickers collected achievements: %v", err)
		}
	}

	return nil
}

// checkCollectorAfterUnlock вспомогательный метод для проверки достижения "Коллекционер"
// после разблокировки стикеров. Вызывается автоматически после streak, villains, friends.
func (s *AchievementService) checkCollectorAfterUnlock(ctx context.Context, childProfileID string) {
	err := s.CheckStickersCollectedAchievements(ctx, childProfileID)
	if err != nil {
		log.Printf("[AchievementService] Failed to check collector achievement for %s: %v", childProfileID, err)
	}
}

// CheckErrorsFoundAchievements проверяет достижения за найденные ошибки
func (s *AchievementService) CheckErrorsFoundAchievements(ctx context.Context, childProfileID string) error {
	errorsCount, err := s.store.GetErrorsFoundCount(ctx, childProfileID)
	if err != nil {
		log.Printf("[AchievementService] Failed to get errors found count for %s: %v", childProfileID, err)
		return err
	}

	unlockedIDs, err := s.store.CheckAndUpdateAchievementsByType(ctx, childProfileID, "errors_found", errorsCount)
	if err != nil {
		log.Printf("[AchievementService] Failed to check errors found achievements for %s: %v", childProfileID, err)
		return err
	}

	if len(unlockedIDs) > 0 {
		log.Printf("[AchievementService] 📝 Unlocked errors found achievements for child %s: %v",
			childProfileID, unlockedIDs)

		// Исправленные ошибки награждают стикерами, проверяем коллекционера
		s.checkCollectorAfterUnlock(ctx, childProfileID)
	}

	return nil
}

// AwardAchievementRewards начисляет награды за разблокированные достижения (монеты, стикеры)
func (s *AchievementService) AwardAchievementRewards(ctx context.Context, childProfileID string, achievementIDs []string) error {
	if len(achievementIDs) == 0 {
		return nil
	}

	// Получаем информацию о наградах для разблокированных достижений
	query := `
		SELECT id, reward_type, reward_amount
		FROM achievements
		WHERE id = ANY($1) AND reward_type = 'coins'
	`

	rows, err := s.store.DB.QueryContext(ctx, query, achievementIDs)
	if err != nil {
		log.Printf("[AwardAchievementRewards] Failed to get rewards for achievements: %v", err)
		return err
	}
	defer rows.Close()

	totalCoins := 0
	var achievementsWithCoins []string

	for rows.Next() {
		var achievementID, rewardType string
		var rewardAmount int

		if err := rows.Scan(&achievementID, &rewardType, &rewardAmount); err != nil {
			log.Printf("[AwardAchievementRewards] Failed to scan reward: %v", err)
			continue
		}

		if rewardType == "coins" && rewardAmount > 0 {
			totalCoins += rewardAmount
			achievementsWithCoins = append(achievementsWithCoins, achievementID)
		}
	}

	// Начисляем монеты если есть
	if totalCoins > 0 && s.profileService != nil {
		err := s.profileService.AddCoins(ctx, childProfileID, totalCoins)
		if err != nil {
			log.Printf("[AwardAchievementRewards] Failed to add %d coins for child %s: %v", totalCoins, childProfileID, err)
			return err
		}

		log.Printf("[AwardAchievementRewards] 💰 Awarded %d coins for %d achievements to child %s: %v",
			totalCoins, len(achievementsWithCoins), childProfileID, achievementsWithCoins)
	}

	return nil
}
