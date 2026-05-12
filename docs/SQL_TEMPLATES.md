# SQL Query Templates

Готовые SQL запросы для типичных задач.

## User Profile Management

### Get full profile info
```sql
SELECT
  id, display_name, platform_id, platform_user_id,
  level, experience_points, coins_balance,
  tasks_solved_total, tasks_solved_correct,
  hints_used_total, streak_days,
  created_at, last_activity_at
FROM child_profiles
WHERE id = '$PROFILE_ID';
```

### Find profile by VK user ID
```sql
SELECT id, display_name, created_at
FROM child_profiles
WHERE platform_id = 'vk' AND platform_user_id = '$VK_USER_ID';
```

### Find all profiles for user (across platforms)
```sql
SELECT platform_id, id, display_name, created_at, last_activity_at
FROM child_profiles
WHERE platform_user_id LIKE '%$USER_ID%'
ORDER BY created_at;
```

### Get profile with related data counts
```sql
SELECT
  cp.*,
  (SELECT COUNT(*) FROM attempts WHERE child_profile_id = cp.id) as attempts,
  (SELECT COUNT(*) FROM child_achievements WHERE child_profile_id = cp.id AND unlocked = true) as achievements,
  (SELECT COUNT(*) FROM referrals WHERE referrer_id = cp.id) as referrals_count,
  (SELECT COUNT(*) FROM villain_battles WHERE child_profile_id = cp.id) as villain_battles
FROM child_profiles cp
WHERE cp.id = '$PROFILE_ID';
```

## Duplicate Detection

### Find duplicate profiles (same user, different platforms)
```sql
-- By display_name (может быть неточно)
SELECT display_name, COUNT(*), array_agg(id), array_agg(platform_id)
FROM child_profiles
GROUP BY display_name
HAVING COUNT(*) > 1;

-- By similar activity time (within 1 hour)
SELECT
  cp1.id as id1, cp1.platform_id as platform1, cp1.display_name as name1,
  cp2.id as id2, cp2.platform_id as platform2, cp2.display_name as name2,
  cp1.created_at, cp2.created_at
FROM child_profiles cp1
JOIN child_profiles cp2 ON cp1.id < cp2.id
WHERE ABS(EXTRACT(EPOCH FROM (cp1.created_at - cp2.created_at))) < 3600
  AND cp1.platform_id != cp2.platform_id
ORDER BY cp1.created_at DESC;
```

### Find web profiles (potential duplicates)
```sql
SELECT
  id, display_name, platform_user_id,
  level, experience_points, coins_balance,
  created_at, last_activity_at,
  (SELECT COUNT(*) FROM attempts WHERE child_profile_id = child_profiles.id) as attempts_count
FROM child_profiles
WHERE platform_id = 'web'
ORDER BY created_at DESC
LIMIT 20;
```

## Profile Cleanup

### Check profile before deletion
```sql
SELECT
  'attempts' as table_name, COUNT(*) FROM attempts WHERE child_profile_id = '$PROFILE_ID'
UNION ALL
SELECT 'achievements', COUNT(*) FROM child_achievements WHERE child_profile_id = '$PROFILE_ID'
UNION ALL
SELECT 'referral_codes', COUNT(*) FROM referral_codes WHERE child_profile_id = '$PROFILE_ID'
UNION ALL
SELECT 'referrals', COUNT(*) FROM referrals WHERE referrer_id = '$PROFILE_ID' OR referred_id = '$PROFILE_ID'
UNION ALL
SELECT 'subscriptions', COUNT(*) FROM subscriptions WHERE child_profile_id = '$PROFILE_ID'
UNION ALL
SELECT 'villain_battles', COUNT(*) FROM villain_battles WHERE child_profile_id = '$PROFILE_ID';
```

### Delete profile (CASCADE will delete all related data)
```sql
-- ОСТОРОЖНО! Это удалит ВСЕ данные пользователя!
DELETE FROM child_profiles
WHERE id = '$PROFILE_ID'
RETURNING id, display_name, platform_id, platform_user_id;
```

### Safe delete (only web profiles without attempts)
```sql
DELETE FROM child_profiles
WHERE id = '$PROFILE_ID'
  AND platform_id = 'web'
  AND NOT EXISTS (SELECT 1 FROM attempts WHERE child_profile_id = '$PROFILE_ID')
RETURNING id, display_name;
```

## Activity Analysis

### Recent activity
```sql
SELECT id, display_name, platform_id, last_activity_at
FROM child_profiles
WHERE last_activity_at IS NOT NULL
ORDER BY last_activity_at DESC
LIMIT 20;
```

### Inactive users (no activity in 30 days)
```sql
SELECT id, display_name, platform_id,
       last_activity_at,
       EXTRACT(DAY FROM (NOW() - last_activity_at)) as days_inactive
FROM child_profiles
WHERE last_activity_at < NOW() - INTERVAL '30 days'
ORDER BY last_activity_at;
```

### Active users today
```sql
SELECT COUNT(*), platform_id
FROM child_profiles
WHERE last_activity_at >= CURRENT_DATE
GROUP BY platform_id;
```

## Statistics

### Platform distribution
```sql
SELECT
  platform_id,
  COUNT(*) as users,
  COUNT(*) FILTER (WHERE last_activity_at >= NOW() - INTERVAL '7 days') as active_7d,
  COUNT(*) FILTER (WHERE last_activity_at >= NOW() - INTERVAL '30 days') as active_30d
FROM child_profiles
GROUP BY platform_id
ORDER BY users DESC;
```

### Top users by XP
```sql
SELECT display_name, platform_id, level, experience_points, tasks_solved_correct
FROM child_profiles
ORDER BY experience_points DESC
LIMIT 10;
```

### Recent registrations
```sql
SELECT id, display_name, platform_id, created_at
FROM child_profiles
WHERE created_at >= CURRENT_DATE - INTERVAL '7 days'
ORDER BY created_at DESC;
```

## Achievements

### Get user achievements
```sql
SELECT
  a.id, a.name, a.description, a.icon,
  ca.unlocked, ca.progress, ca.max_progress,
  ca.unlocked_at, ca.viewed
FROM child_achievements ca
JOIN achievements a ON ca.achievement_id = a.id
WHERE ca.child_profile_id = '$PROFILE_ID'
ORDER BY ca.unlocked DESC, ca.unlocked_at DESC;
```

### Achievement stats
```sql
SELECT
  cp.id, cp.display_name,
  COUNT(ca.id) FILTER (WHERE ca.unlocked = true) as unlocked,
  COUNT(a.id) as total
FROM child_profiles cp
CROSS JOIN achievements a
LEFT JOIN child_achievements ca ON ca.child_profile_id = cp.id AND ca.achievement_id = a.id
WHERE cp.id = '$PROFILE_ID'
GROUP BY cp.id, cp.display_name;
```

## Quick Lookups

### Get profile ID by VK user
```sql
SELECT id FROM child_profiles
WHERE platform_id = 'vk' AND platform_user_id = '$VK_USER_ID';
```

### Check if profile exists
```sql
SELECT EXISTS(
  SELECT 1 FROM child_profiles
  WHERE platform_id = '$PLATFORM' AND platform_user_id = '$USER_ID'
) as exists;
```