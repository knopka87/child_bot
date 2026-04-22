-- Rollback auto-claim rewards

DROP TRIGGER IF EXISTS trigger_auto_claim_achievement ON child_achievements;
DROP TRIGGER IF EXISTS trigger_auto_claim_achievement_insert ON child_achievements;
DROP TRIGGER IF EXISTS trigger_auto_claim_referral ON referrals;
DROP TRIGGER IF EXISTS trigger_auto_claim_milestone ON child_referral_milestones;

DROP FUNCTION IF EXISTS auto_claim_achievement_reward();
DROP FUNCTION IF EXISTS auto_claim_achievement_on_insert();
DROP FUNCTION IF EXISTS auto_claim_referral_reward();
DROP FUNCTION IF EXISTS auto_claim_milestone_reward();
