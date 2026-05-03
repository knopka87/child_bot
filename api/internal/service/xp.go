package service

import (
	"context"
	"log"

	"child-bot/api/internal/store"
)

// XP Rewards константы
const (
	XPForCorrectAnswer = 50  // Правильное решение задачи
	XPForFixErrors     = 20  // Исправление ошибок
	XPForHintRequest   = 10  // Запрос подсказки
	XPForDailyLogin    = 30  // Ежедневный вход
	XPForVillainDefeat = 100 // Победа над злодеем
	// XPForAchievement declared in achievement.go to avoid circular dependency
)

// AwardCorrectAnswer начисляет XP за правильное решение
func (s *ProfileService) AwardCorrectAnswer(ctx context.Context, childProfileID string) error {
	log.Printf("[ProfileService.AwardCorrectAnswer] 🎯 Starting to award %d XP for correct answer to child %s", XPForCorrectAnswer, childProfileID)

	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForCorrectAnswer, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] ❌ Failed to award correct answer XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from correct answer! child=%s, level=%d", childProfileID, level)
	} else {
		log.Printf("[ProfileService.AwardCorrectAnswer] ✅ Successfully awarded %d XP for correct answer to child %s (no level up)", XPForCorrectAnswer, childProfileID)
	}

	return nil
}

// AwardFixErrors начисляет XP за исправление ошибок
func (s *ProfileService) AwardFixErrors(ctx context.Context, childProfileID string) error {
	log.Printf("[ProfileService.AwardFixErrors] 🎯 Starting to award %d XP for fixing errors to child %s", XPForFixErrors, childProfileID)

	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForFixErrors, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] ❌ Failed to award fix errors XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from fixing errors! child=%s, level=%d", childProfileID, level)
	} else {
		log.Printf("[ProfileService.AwardFixErrors] ✅ Successfully awarded %d XP for fixing errors to child %s (no level up)", XPForFixErrors, childProfileID)
	}

	return nil
}

// AwardHintRequest начисляет XP за запрос подсказки
func (s *ProfileService) AwardHintRequest(ctx context.Context, childProfileID string) error {
	log.Printf("[ProfileService.AwardHintRequest] 🎯 Starting to award %d XP for hint request to child %s", XPForHintRequest, childProfileID)

	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForHintRequest, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] ❌ Failed to award hint request XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from hint request! child=%s, level=%d", childProfileID, level)
	} else {
		log.Printf("[ProfileService.AwardHintRequest] ✅ Successfully awarded %d XP for hint request to child %s (no level up)", XPForHintRequest, childProfileID)
	}

	return nil
}

// AwardDailyLogin начисляет XP за ежедневный вход
func (s *ProfileService) AwardDailyLogin(ctx context.Context, childProfileID string) error {
	log.Printf("[ProfileService.AwardDailyLogin] 🎯 Starting to award %d XP for daily login to child %s", XPForDailyLogin, childProfileID)

	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForDailyLogin, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] ❌ Failed to award daily login XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from daily login! child=%s, level=%d", childProfileID, level)
	} else {
		log.Printf("[ProfileService.AwardDailyLogin] ✅ Successfully awarded %d XP for daily login to child %s (no level up)", XPForDailyLogin, childProfileID)
	}

	return nil
}

// AwardVillainDefeat начисляет XP за победу над злодеем
func (s *ProfileService) AwardVillainDefeat(ctx context.Context, childProfileID string) error {
	log.Printf("[ProfileService.AwardVillainDefeat] 🎯 Starting to award %d XP for villain defeat to child %s", XPForVillainDefeat, childProfileID)

	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForVillainDefeat, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] ❌ Failed to award villain defeat XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from villain defeat! child=%s, level=%d", childProfileID, level)
	} else {
		log.Printf("[ProfileService.AwardVillainDefeat] ✅ Successfully awarded %d XP for villain defeat to child %s (no level up)", XPForVillainDefeat, childProfileID)
	}

	return nil
}

// AwardAchievementUnlock начисляет XP за разблокировку достижения
func (s *ProfileService) AwardAchievementUnlock(ctx context.Context, childProfileID string) error {
	log.Printf("[ProfileService.AwardAchievementUnlock] 🎯 Starting to award %d XP for achievement unlock to child %s", XPForAchievement, childProfileID)

	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForAchievement, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] ❌ Failed to award achievement XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from achievement! child=%s, level=%d", childProfileID, level)
	} else {
		log.Printf("[ProfileService.AwardAchievementUnlock] ✅ Successfully awarded %d XP for achievement unlock to child %s (no level up)", XPForAchievement, childProfileID)
	}

	return nil
}
