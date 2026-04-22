-- Откат тестовых данных
DELETE FROM child_referral_milestones WHERE child_profile_id = '462d5291-d5f3-4626-bcf9-3cd933d3e5be';
DELETE FROM referrals WHERE referrer_id = '462d5291-d5f3-4626-bcf9-3cd933d3e5be';
DELETE FROM referral_codes WHERE child_profile_id = '462d5291-d5f3-4626-bcf9-3cd933d3e5be';
DELETE FROM child_profiles WHERE id IN (
    '11111111-1111-1111-1111-111111111111',
    '22222222-2222-2222-2222-222222222222',
    '33333333-3333-3333-3333-333333333333'
);
