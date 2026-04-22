package service

import (
	"context"
	"log"

	"child-bot/api/internal/store"
)

const XPForAchievement = 50

// AchievementService бизнес-логика для достижений
type AchievementService struct {
	store *store.Store
}

// NewAchievementService создает новый AchievementService
func NewAchievementService(store *store.Store) *AchievementService {
	return &AchievementService{store: store}
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
		for range unlockedIDs {
			_, _, _ = s.store.AddXP(ctx, childProfileID, XPForAchievement, store.DefaultXPConfig)
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
