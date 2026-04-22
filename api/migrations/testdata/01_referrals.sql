-- Миграция с тестовыми данными для разработки
-- Применяется только при запуске с флагом make migrate-test

-- Реферальный код для тестового пользователя
INSERT INTO referral_codes (child_profile_id, code, uses_count)
VALUES ('462d5291-d5f3-4626-bcf9-3cd933d3e5be', 'PETYA123', 2)
ON CONFLICT (child_profile_id) DO UPDATE SET code = 'PETYA123', uses_count = 2;

-- Тестовые друзья
INSERT INTO child_profiles (id, display_name, avatar_id, grade, platform_id, platform_user_id, created_at)
VALUES
    ('11111111-1111-1111-1111-111111111111', 'Маша', 'avatar_girl_1', 5, 'web', 'friend_1', NOW() - INTERVAL '10 days'),
    ('22222222-2222-2222-2222-222222222222', 'Саша', 'avatar_boy_2', 6, 'web', 'friend_2', NOW() - INTERVAL '5 days'),
    ('33333333-3333-3333-3333-333333333333', 'Катя', 'avatar_girl_2', 5, 'web', 'friend_3', NOW() - INTERVAL '2 days')
ON CONFLICT (platform_id, platform_user_id) DO NOTHING;

-- Реферальные связи (2 активных друга + 1 неактивный)
INSERT INTO referrals (referrer_id, referred_id, is_active, reward_coins, reward_claimed, invited_at, activated_at, reward_claimed_at)
VALUES
    -- Маша - активна, награда получена
    ('462d5291-d5f3-4626-bcf9-3cd933d3e5be', '11111111-1111-1111-1111-111111111111', TRUE, 50, TRUE,
     NOW() - INTERVAL '10 days', NOW() - INTERVAL '9 days', NOW() - INTERVAL '9 days'),

    -- Саша - активен, награда получена
    ('462d5291-d5f3-4626-bcf9-3cd933d3e5be', '22222222-2222-2222-2222-222222222222', TRUE, 50, TRUE,
     NOW() - INTERVAL '5 days', NOW() - INTERVAL '4 days', NOW() - INTERVAL '4 days'),

    -- Катя - неактивна (только зарегистрировалась)
    ('462d5291-d5f3-4626-bcf9-3cd933d3e5be', '33333333-3333-3333-3333-333333333333', FALSE, 50, FALSE,
     NOW() - INTERVAL '2 days', NULL, NULL)
ON CONFLICT (referrer_id, referred_id) DO NOTHING;

-- Milestone для 1 и 3 друзей получены, следующий - 5 друзей
INSERT INTO child_referral_milestones (child_profile_id, milestone_id, is_claimed, claimed_at)
SELECT '462d5291-d5f3-4626-bcf9-3cd933d3e5be', id, TRUE, NOW() - INTERVAL '9 days'
FROM referral_milestones
WHERE friends_count = 1
ON CONFLICT (child_profile_id, milestone_id) DO NOTHING;

INSERT INTO child_referral_milestones (child_profile_id, milestone_id, is_claimed, claimed_at)
SELECT '462d5291-d5f3-4626-bcf9-3cd933d3e5be', id, TRUE, NOW() - INTERVAL '4 days'
FROM referral_milestones
WHERE friends_count = 3
ON CONFLICT (child_profile_id, milestone_id) DO NOTHING;
