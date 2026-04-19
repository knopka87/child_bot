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
	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForCorrectAnswer, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] Failed to award correct answer XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from correct answer! child=%s, level=%d", childProfileID, level)
	}

	return nil
}

// AwardFixErrors начисляет XP за исправление ошибок
func (s *ProfileService) AwardFixErrors(ctx context.Context, childProfileID string) error {
	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForFixErrors, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] Failed to award fix errors XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from fixing errors! child=%s, level=%d", childProfileID, level)
	}

	return nil
}

// AwardHintRequest начисляет XP за запрос подсказки
func (s *ProfileService) AwardHintRequest(ctx context.Context, childProfileID string) error {
	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForHintRequest, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] Failed to award hint request XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from hint request! child=%s, level=%d", childProfileID, level)
	}

	return nil
}

// AwardDailyLogin начисляет XP за ежедневный вход
func (s *ProfileService) AwardDailyLogin(ctx context.Context, childProfileID string) error {
	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForDailyLogin, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] Failed to award daily login XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from daily login! child=%s, level=%d", childProfileID, level)
	}

	return nil
}

// AwardVillainDefeat начисляет XP за победу над злодеем
func (s *ProfileService) AwardVillainDefeat(ctx context.Context, childProfileID string) error {
	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForVillainDefeat, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] Failed to award villain defeat XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from villain defeat! child=%s, level=%d", childProfileID, level)
	}

	return nil
}

// AwardAchievementUnlock начисляет XP за разблокировку достижения
func (s *ProfileService) AwardAchievementUnlock(ctx context.Context, childProfileID string) error {
	level, leveledUp, err := s.store.AddXP(ctx, childProfileID, XPForAchievement, store.DefaultXPConfig)
	if err != nil {
		log.Printf("[ProfileService] Failed to award achievement XP for %s: %v", childProfileID, err)
		return err
	}

	if leveledUp {
		log.Printf("[ProfileService] 🎉 Level up from achievement! child=%s, level=%d", childProfileID, level)
	}

	return nil
}
