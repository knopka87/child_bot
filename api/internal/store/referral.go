package store

import (
	"context"
	"database/sql"
	"time"
)

// ReferralCode данные реферального кода
type ReferralCode struct {
	ChildProfileID string
	Code           string
	UsesCount      int
	CreatedAt      time.Time
}

// ReferralStats статистика рефералов
type ReferralStats struct {
	TotalInvited  int
	ActiveInvited int
	TotalRewards  int
}

// InvitedFriendDB приглашенный друг из БД
type InvitedFriendDB struct {
	ID           string
	DisplayName  string
	AvatarID     string
	InvitedAt    time.Time
	ActivatedAt  sql.NullTime
	IsActive     bool
	RewardCoins  int
	RewardClaimed bool
}

// RewardMilestoneDB milestone из БД
type RewardMilestoneDB struct {
	ID           string
	FriendsCount int
	RewardCoins  int
	Description  string
	IsClaimed    bool
}

// GetReferralCode получает реферальный код пользователя
func (s *Store) GetReferralCode(ctx context.Context, childProfileID string) (*ReferralCode, error) {
	query := `
		SELECT child_profile_id, code, uses_count, created_at
		FROM referral_codes
		WHERE child_profile_id = $1
	`

	var rc ReferralCode
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(
		&rc.ChildProfileID,
		&rc.Code,
		&rc.UsesCount,
		&rc.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &rc, nil
}

// GetReferralStats получает статистику рефералов
func (s *Store) GetReferralStats(ctx context.Context, childProfileID string) (*ReferralStats, error) {
	query := `
		SELECT
			COUNT(*) as total_invited,
			COUNT(*) FILTER (WHERE is_active = true) as active_invited,
			COALESCE(SUM(reward_coins) FILTER (WHERE reward_claimed = true), 0) as total_rewards
		FROM referrals
		WHERE referrer_id = $1
	`

	var stats ReferralStats
	err := s.DB.QueryRowContext(ctx, query, childProfileID).Scan(
		&stats.TotalInvited,
		&stats.ActiveInvited,
		&stats.TotalRewards,
	)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}

// GetInvitedFriends получает список приглашенных друзей
func (s *Store) GetInvitedFriends(ctx context.Context, childProfileID string) ([]InvitedFriendDB, error) {
	query := `
		SELECT
			cp.id,
			cp.display_name,
			cp.avatar_id,
			r.invited_at,
			r.activated_at,
			r.is_active,
			r.reward_coins,
			r.reward_claimed
		FROM referrals r
		JOIN child_profiles cp ON r.referred_id = cp.id
		WHERE r.referrer_id = $1
		ORDER BY r.invited_at DESC
	`

	rows, err := s.DB.QueryContext(ctx, query, childProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []InvitedFriendDB
	for rows.Next() {
		var f InvitedFriendDB
		err := rows.Scan(
			&f.ID,
			&f.DisplayName,
			&f.AvatarID,
			&f.InvitedAt,
			&f.ActivatedAt,
			&f.IsActive,
			&f.RewardCoins,
			&f.RewardClaimed,
		)
		if err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return friends, nil
}

// GetRewardMilestones получает milestone'ы с их статусом
func (s *Store) GetRewardMilestones(ctx context.Context, childProfileID string) ([]RewardMilestoneDB, error) {
	query := `
		SELECT
			rm.id,
			rm.friends_count,
			rm.reward_coins,
			rm.description,
			COALESCE(crm.is_claimed, false) as is_claimed
		FROM referral_milestones rm
		LEFT JOIN child_referral_milestones crm
			ON rm.id = crm.milestone_id AND crm.child_profile_id = $1
		ORDER BY rm.friends_count ASC
	`

	rows, err := s.DB.QueryContext(ctx, query, childProfileID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var milestones []RewardMilestoneDB
	for rows.Next() {
		var m RewardMilestoneDB
		err := rows.Scan(
			&m.ID,
			&m.FriendsCount,
			&m.RewardCoins,
			&m.Description,
			&m.IsClaimed,
		)
		if err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return milestones, nil
}
