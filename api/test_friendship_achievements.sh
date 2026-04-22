#!/bin/bash

# Тестовый скрипт для проверки системы достижений "Стикер Дружба"

DB_URL="postgres://child_bot:dev_secret@localhost:5432/child_bot"

echo "=== Тест системы достижений Стикер Дружба ==="
echo ""

# 1. Создаем основного пользователя (реферера)
echo "1. Создаем основного пользователя..."
REFERRER_ID=$(psql $DB_URL -t -A -c "
    INSERT INTO child_profiles (display_name, avatar_id, grade, platform_id, platform_user_id)
    VALUES ('Главный', 'cat', 3, 'test', 'referrer_001')
    ON CONFLICT (platform_id, platform_user_id)
    DO UPDATE SET updated_at = NOW()
    RETURNING id;
" | head -1 | tr -d ' \n')
echo "   Referrer ID: $REFERRER_ID"

# 2. Создаем реферальный код
echo "2. Создаем реферальный код..."
REFERRAL_CODE=$(psql $DB_URL -t -A -c "
    INSERT INTO referral_codes (child_profile_id, code)
    VALUES ('$REFERRER_ID', 'TEST' || substr(md5(random()::text), 1, 4))
    ON CONFLICT (child_profile_id)
    DO UPDATE SET uses_count = referral_codes.uses_count
    RETURNING code;
" | head -1 | tr -d ' \n')
echo "   Referral code: $REFERRAL_CODE"

# 3. Создаем 7 друзей и активируем их
echo "3. Создаем и активируем друзей..."
for i in {1..7}; do
    FRIEND_ID=$(psql $DB_URL -t -A -c "
        INSERT INTO child_profiles (display_name, avatar_id, grade, platform_id, platform_user_id)
        VALUES ('Друг $i', 'dog', 3, 'test', 'friend_00$i')
        ON CONFLICT (platform_id, platform_user_id)
        DO UPDATE SET updated_at = NOW()
        RETURNING id;
    " | head -1 | tr -d ' \n')

    psql $DB_URL -q -c "
        INSERT INTO referrals (referrer_id, referred_id, is_active, reward_coins)
        VALUES ('$REFERRER_ID', '$FRIEND_ID', true, 50)
        ON CONFLICT (referrer_id, referred_id) DO NOTHING;
    " >/dev/null 2>&1
    echo "   Друг $i создан и активирован (ID: $FRIEND_ID)"
done

# 4. Проверяем количество активных друзей
echo ""
echo "4. Проверяем количество активных друзей..."
ACTIVE_COUNT=$(psql $DB_URL -t -A -c "
    SELECT COUNT(*)
    FROM referrals
    WHERE referrer_id = '$REFERRER_ID' AND is_active = true;
" | head -1 | tr -d ' \n')
echo "   Активных друзей: $ACTIVE_COUNT"

# 5. Проверяем разблокированные достижения
echo ""
echo "5. Проверяем разблокированные достижения..."
psql $DB_URL -c "
    SELECT
        a.title,
        a.requirement_value as требуется,
        ca.current_progress as прогресс,
        ca.is_unlocked as разблокировано,
        ca.unlocked_at
    FROM achievements a
    LEFT JOIN child_achievements ca ON a.id = ca.achievement_id AND ca.child_profile_id = '$REFERRER_ID'
    WHERE a.requirement_type = 'friends_invited'
    ORDER BY a.priority;
"

# 6. Запускаем проверку достижений через API (если сервер запущен)
echo ""
echo "6. Тестируем API проверки достижений..."
curl -s -X POST "http://localhost:8080/test/check-achievements/$REFERRER_ID" 2>/dev/null || echo "   API недоступен (это нормально, если нет тестового эндпоинта)"

echo ""
echo "=== Тест завершен ==="
echo "Referrer ID для ручного тестирования: $REFERRER_ID"
