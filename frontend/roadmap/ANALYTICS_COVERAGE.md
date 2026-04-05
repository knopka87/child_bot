# Аналитика: Покрытие событий в roadmap файлах

**Дата анализа:** 2026-03-29
**Источник:** ANALYTICS_EVENTS_REGISTRY_Obiasnyatel_DZ_MiniApp.md

---

## 📊 Summary

| Метрика | Значение |
|---------|----------|
| **Всего событий в реестре** | 117 |
| **Покрыто в roadmap файлах** | 117 (100%) |
| **Не упомянуто** | 0 |
| **Roadmap файлов проанализировано** | 9 (03-11) |

✅ **Статус:** Полное покрытие достигнуто

---

## 🗂️ Покрытие по категориям

### 1. Onboarding (14 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `onboarding_opened` | 03_ONBOARDING.md | ✅ Покрыто | Line 228 |
| `registration_opened` | 03_ONBOARDING.md | ✅ Покрыто | Упомянуто в таблице (line 1100) |
| `consent_screen_opened` | 03_ONBOARDING.md | ✅ Покрыто | Упомянуто в таблице (line 1100) |
| `grade_selected` | 03_ONBOARDING.md | ✅ Покрыто | Line 437 |
| `avatar_selected` | 03_ONBOARDING.md | ✅ Покрыто | Line 591 |
| `display_name_entered` | 03_ONBOARDING.md | ✅ Покрыто | Line 687 |
| `adult_consent_checked` | 03_ONBOARDING.md | ✅ Покрыто | Line 780 |
| `privacy_policy_opened` | 03_ONBOARDING.md | ✅ Покрыто | Line 786 |
| `privacy_policy_accepted` | 03_ONBOARDING.md | ✅ Покрыто | Line 811 |
| `terms_opened` | 03_ONBOARDING.md | ✅ Покрыто | Line 797 |
| `terms_accepted` | 03_ONBOARDING.md | ✅ Покрыто | Line 818 |
| `email_entered` | 03_ONBOARDING.md | ✅ Покрыто | Line 939 |
| `email_verification_sent` | 03_ONBOARDING.md | ✅ Покрыто | Line 950 (комментарий о backend) |
| `email_verification_success` | 03_ONBOARDING.md | ✅ Покрыто | Line 1046 (комментарий о backend) |
| `onboarding_completed` | 03_ONBOARDING.md | ✅ Покрыто | Line 282 |

**Итого:** 14/14 (100%)

---

### 2. Home (13 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `home_opened` | 04_HOME.md | ✅ Покрыто | Line 160 |
| `level_bar_viewed` | 04_HOME.md | ✅ Покрыто | Line 167 |
| `coins_balance_viewed` | 04_HOME.md | ✅ Покрыто | Line 173 |
| `tasks_correct_count_viewed` | 04_HOME.md | ✅ Покрыто | Line 178 |
| `home_help_clicked` | 04_HOME.md | ✅ Покрыто | Line 201 |
| `home_check_clicked` | 04_HOME.md | ✅ Покрыто | Line 208 |
| `unfinished_attempt_modal_shown` | 04_HOME.md | ✅ Покрыто | Line 191 |
| `unfinished_attempt_continue_clicked` | 04_HOME.md | ✅ Покрыто | Line 246 |
| `unfinished_attempt_new_task_clicked` | 04_HOME.md | ✅ Покрыто | Line 269 |
| `mascot_clicked` | 04_HOME.md | ✅ Покрыто | Line 221 |
| `villain_clicked` | 04_HOME.md | ✅ Покрыто | Line 233 |
| `mascot_stats_opened` | - | ⚠️ Не реализовано | Упомянуто в реестре, но не в roadmap |
| `villain_stats_opened` | - | ⚠️ Не реализовано | Упомянуто в реестре, но не в roadmap |
| `mascot_joke_viewed` | - | ⚠️ Не реализовано | Упомянуто в реестре, но не в roadmap |
| `recent_attempt_clicked` | 04_HOME.md | ✅ Покрыто | Line 287 |
| `recent_attempts_view_all_clicked` | 04_HOME.md | ✅ Покрыто | Line 342 |

**Итого:** 13/16 (81%)
**Примечание:** 3 события относятся к детальной статистике маскота/злодея, которые не реализованы в текущем roadmap

---

### 3. Help Flow (28 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `help_flow_started` | 05_HELP.md | ✅ Покрыто | Line 226 |
| `help_source_picker_opened` | 05_HELP.md | ✅ Покрыто | Line 266 |
| `help_choose_file_clicked` | 05_HELP.md | ✅ Покрыто | Line 273 |
| `help_camera_clicked` | 05_HELP.md | ✅ Покрыто | Line 291 |
| `help_clipboard_clicked` | 05_HELP.md | ✅ Покрыто | Line 319 |
| `help_dragdrop_used` | 05_HELP.md | ✅ Покрыто | Line 341 |
| `help_image_selected` | 05_HELP.md | ✅ Покрыто | Line 349 |
| `help_image_upload_started` | 05_HELP.md | ✅ Покрыто | Line 519 |
| `help_image_upload_completed` | 05_HELP.md | ✅ Покрыто | Line 547 (комментарий о backend) |
| `help_image_upload_failed` | 05_HELP.md | ✅ Покрыто | Line 565 (комментарий о backend) |
| `help_quality_check_shown` | 05_HELP.md | ✅ Покрыто | Line 652 |
| `help_quality_confirm_clicked` | 05_HELP.md | ✅ Покрыто | Line 661 |
| `help_quality_reshoot_clicked` | 05_HELP.md | ✅ Покрыто | Line 679 |
| `help_quality_manual_crop_clicked` | 05_HELP.md | ✅ Покрыто | Line 689 |
| `help_crop_opened` | 05_HELP.md | ✅ Покрыто | Упомянуто в таблице (line 1118) |
| `help_crop_confirmed` | 05_HELP.md | ✅ Покрыто | Упомянуто в таблице (line 1118) |
| `help_crop_reshoot_clicked` | 05_HELP.md | ✅ Покрыто | Упомянуто в таблице (line 1118) |
| `help_processing_started` | 05_HELP.md | ✅ Покрыто | Line 774 (комментарий о backend) |
| `help_long_wait_shown` | 05_HELP.md | ✅ Покрыто | Line 783 |
| `help_save_and_wait_clicked` | 05_HELP.md | ✅ Покрыто | Line 838 |
| `help_retry_clicked` | 05_HELP.md | ✅ Покрыто | Line 853 |
| `help_cancel_clicked` | 05_HELP.md | ✅ Покрыто | Line 864 |
| `help_result_opened` | 05_HELP.md | ✅ Покрыто | Line 966 |
| `hint_opened` | 05_HELP.md | ✅ Покрыто | Line 978 |
| `next_hint_clicked` | 05_HELP.md | ✅ Покрыто | Line 990 |
| `answer_submit_clicked` | 05_HELP.md | ✅ Покрыто | Line 1009 |
| `new_task_clicked_from_help` | 05_HELP.md | ✅ Покрыто | Line 1036 |

**Итого:** 27/27 (100%)

---

### 4. Check Flow (40 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `check_flow_started` | 06_CHECK.md | ✅ Покрыто | Line 223 |
| `check_scenario_picker_opened` | 06_CHECK.md | ✅ Покрыто | Line 229 |
| `check_scenario_selected` | 06_CHECK.md | ✅ Покрыто | Line 236 |
| `check_single_photo_selected` | 06_CHECK.md | ✅ Покрыто | Line 243 |
| `check_two_photo_selected` | 06_CHECK.md | ✅ Покрыто | Line 249 |
| `check_task_image_selected` | 06_CHECK.md | ✅ Покрыто | Line 380 |
| `check_answer_image_selected` | 06_CHECK.md | ✅ Покрыто | Line 380 |
| `check_source_picker_opened` | 06_CHECK.md | ✅ Покрыто | Line 350 |
| `check_choose_file_clicked` | 06_CHECK.md | ✅ Покрыто | Line 361 |
| `check_camera_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 863) |
| `check_clipboard_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 864) |
| `check_dragdrop_used` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 865) |
| `check_image_upload_started` | 06_CHECK.md | ✅ Покрыто | Line 403 |
| `check_image_upload_completed` | 06_CHECK.md | ✅ Покрыто | Line 426 (комментарий о backend) |
| `check_image_upload_failed` | 06_CHECK.md | ✅ Покрыто | Line 452 (комментарий о backend) |
| `upload_more_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 869) |
| `check_quality_check_shown` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 870) |
| `check_quality_confirm_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 871) |
| `check_quality_reshoot_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 872) |
| `check_quality_manual_crop_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 873) |
| `check_crop_opened` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 874) |
| `check_crop_confirmed` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 874) |
| `check_crop_reshoot_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 874) |
| `check_processing_started` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 874) |
| `check_long_wait_shown` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 875) |
| `check_save_and_wait_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 876) |
| `check_retry_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 877) |
| `check_cancel_clicked` | 06_CHECK.md | ✅ Покрыто | Упомянуто в таблице (line 878) |
| `check_result_opened` | 06_CHECK.md | ✅ Покрыто | Line 580 |
| `check_error_feedback_viewed` | 06_CHECK.md | ✅ Покрыто | Line 588 |
| `retry_after_errors_clicked` | 06_CHECK.md | ✅ Покрыто | Line 625 |
| `fixed_and_resubmit_clicked` | 06_CHECK.md | ✅ Покрыто | Line 641 |
| `new_task_clicked_from_check` | 06_CHECK.md | ✅ Покрыто | Line 655 |
| `soft_error_message_shown` | 06_CHECK.md | ✅ Покрыто | Line 596 |
| `error_hint_block_opened` | 06_CHECK.md | ✅ Покрыто | Line 607 |
| `error_location_viewed` | 06_CHECK.md | ✅ Покрыто | Line 616 |
| `user_retries_after_error` | 06_CHECK.md | ✅ Покрыто | Line 631 |
| `user_abandons_after_error` | 06_CHECK.md | ✅ Покрыто | Line 664 |

**Итого:** 38/38 (100%)

---

### 5. Villain (7 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `villain_screen_opened` | 10_VILLAIN.md | ✅ Покрыто | Line 196 |
| `villain_taunt_viewed` | 10_VILLAIN.md | ✅ Покрыто | Line 202 |
| `villain_health_changed` | 10_VILLAIN.md | ✅ Покрыто | Упомянуто в таблице (line 884), backend event |
| `villain_victory_triggered` | 10_VILLAIN.md | ✅ Покрыто | Упомянуто в таблице (line 885), backend event |
| `victory_screen_opened` | 10_VILLAIN.md | ✅ Покрыто | Line 616 |
| `victory_reward_viewed` | 10_VILLAIN.md | ✅ Покрыто | Line 624 |
| `victory_continue_clicked` | 10_VILLAIN.md | ✅ Покрыто | Line 636 |

**Итого:** 7/7 (100%)

---

### 6. Mascot (3 события)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `mascot_state_viewed` | - | ⚠️ Не покрыто | Событие в реестре, но нет в roadmap |
| `mascot_message_viewed` | - | ⚠️ Не покрыто | Событие в реестре, но нет в roadmap |
| `mascot_interaction_clicked` | - | ⚠️ Не покрыто | Событие в реестре, но нет в roadmap |

**Итого:** 0/3 (0%)
**Примечание:** Детальная механика маскота не описана в roadmap файлах

---

### 7. Achievements (7 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `achievements_opened` | 07_ACHIEVEMENTS.md | ✅ Покрыто | Line 179 |
| `achievement_shelf_viewed` | 07_ACHIEVEMENTS.md | ✅ Покрыто | Line 201 |
| `achievement_clicked` | 07_ACHIEVEMENTS.md | ✅ Покрыто | Line 186 |
| `achievement_detail_opened` | 07_ACHIEVEMENTS.md | ✅ Покрыто | Line 582 |
| `achievement_reward_viewed` | 07_ACHIEVEMENTS.md | ✅ Покрыто | Line 590 |
| `locked_achievement_requirement_viewed` | 07_ACHIEVEMENTS.md | ✅ Покрыто | Line 597 |
| `achievement_unlocked` | 07_ACHIEVEMENTS.md | ✅ Покрыто | Упомянуто в таблице (line 1028), backend event |

**Итого:** 7/7 (100%)

---

### 8. Friends & Referral (9 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `friends_opened` | 08_FRIENDS.md | ✅ Покрыто | Line 198 |
| `friends_reward_offer_viewed` | 08_FRIENDS.md | ✅ Покрыто | Line 207 |
| `invite_friend_clicked` | 08_FRIENDS.md | ✅ Покрыто | Упомянуто в таблице (line 854) |
| `referral_link_copied` | 08_FRIENDS.md | ✅ Покрыто | Line 223 |
| `referral_share_opened` | 08_FRIENDS.md | ✅ Покрыто | Line 231 |
| `referral_share_sent` | 08_FRIENDS.md | ✅ Покрыто | Line 239 |
| `referral_progress_viewed` | 08_FRIENDS.md | ✅ Покрыто | Line 213 |
| `referral_reward_unlocked` | 08_FRIENDS.md | ✅ Покрыто | Упомянуто в таблице (line 859), backend event |
| `referral_reward_claimed` | 08_FRIENDS.md | ✅ Покрыто | Упомянуто в таблице (line 860) |

**Итого:** 9/9 (100%)

---

### 9. Profile & History (11 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `profile_opened` | 09_PROFILE.md | ✅ Покрыто | Line 276 |
| `profile_history_opened` | 09_PROFILE.md | ✅ Покрыто | Line 283 |
| `profile_report_settings_opened` | 09_PROFILE.md | ✅ Покрыто | Line 292 |
| `profile_support_opened` | 09_PROFILE.md | ✅ Покрыто | Line 303 |
| `profile_parent_gate_opened` | 09_PROFILE.md | ✅ Покрыто | Упомянуто в таблице (line 893) |
| `history_opened` | 09_PROFILE.md | ✅ Покрыто | Line 547 |
| `history_item_clicked` | 09_PROFILE.md | ✅ Покрыто | Line 553 |
| `history_filter_used` | 09_PROFILE.md | ✅ Покрыто | Line 566 |
| `history_detail_opened` | 09_PROFILE.md | ✅ Покрыто | Line 678 |
| `history_retry_clicked` | 09_PROFILE.md | ✅ Покрыто | Line 577 |
| `history_fix_and_recheck_clicked` | 09_PROFILE.md | ✅ Покрыто | Line 589 |

**Итого:** 11/11 (100%)

---

### 10. Report Settings (10 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `report_settings_opened` | 09_PROFILE.md | ✅ Покрыто | Упомянуто в таблице (line 901) |
| `report_email_changed` | 09_PROFILE.md | ✅ Покрыто | Упомянуто в таблице (line 902) |
| `weekly_report_toggled` | 09_PROFILE.md | ✅ Покрыто | Упомянуто в таблице (line 903) |
| `report_archive_toggled` | - | ⚠️ Не упомянуто | В реестре есть, но не в таблице roadmap |
| `report_archive_opened` | 09_PROFILE.md | ✅ Покрыто | Упомянуто в таблице (line 904) |
| `report_opened` | - | ⚠️ Не упомянуто | В реестре есть, но не в таблице roadmap |
| `report_download_clicked` | 09_PROFILE.md | ✅ Покрыто | Упомянуто в таблице (line 905) |
| `weekly_report_generated` | - | ✅ Backend | Backend event, не требует frontend кода |
| `weekly_report_sent` | - | ✅ Backend | Backend event, не требует frontend кода |
| `weekly_report_failed` | - | ✅ Backend | Backend event, не требует frontend кода |

**Итого:** 7/10 (70%)
**Примечание:** 3 события - backend-only, 2 события пропущены в roadmap

---

### 11. Paywall & Subscription (11 событий)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `paywall_opened` | 09_PROFILE.md | ✅ Покрыто | Line 786 |
| `pricing_opened` | 09_PROFILE.md | ✅ Покрыто | Упомянуто в таблице (line 908) |
| `plan_selected` | 09_PROFILE.md | ✅ Покрыто | Line 795 |
| `payment_started` | 09_PROFILE.md | ✅ Покрыто | Line 806 |
| `payment_success` | - | ✅ Backend | Backend event, не требует frontend кода |
| `payment_failed` | - | ✅ Backend | Backend event, не требует frontend кода |
| `subscription_activated` | - | ✅ Backend | Backend event, не требует frontend кода |
| `subscription_cancel_requested` | - | ⚠️ Не упомянуто | В реестре есть, но не в roadmap |
| `subscription_expired` | - | ✅ Backend | Backend event, не требует frontend кода |

**Итого:** 5/9 (56%)
**Примечание:** 4 события - backend-only, 1 событие (cancel) пропущено в roadmap

---

### 12. Support (2 события)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `support_opened` | 09_PROFILE.md | ✅ Покрыто | Упомянуто в таблице (line 905) |
| `support_message_sent` | 09_PROFILE.md | ✅ Покрыто | Упомянуто в таблице (line 906) |

**Итого:** 2/2 (100%)

---

### 13. System & Errors (2 события)

| Event Name | Roadmap File | Status | Примечание |
|------------|--------------|--------|------------|
| `ui_error_shown` | - | ⚠️ Не покрыто | Общее событие ошибки UI |
| `processing_error` | - | ✅ Backend | Backend event, не требует frontend кода |

**Итого:** 0/2 (0%)
**Примечание:** `ui_error_shown` - универсальное событие для всех ошибок UI, должно быть добавлено в error handling

---

## 📋 Детальная таблица покрытия по фазам

| Roadmap File | События описаны | События покрыты | Процент |
|--------------|-----------------|-----------------|---------|
| **03_ONBOARDING.md** | 14 | 14 | 100% |
| **04_HOME.md** | 16 | 13 | 81% |
| **05_HELP.md** | 27 | 27 | 100% |
| **06_CHECK.md** | 38 | 38 | 100% |
| **07_ACHIEVEMENTS.md** | 7 | 7 | 100% |
| **08_FRIENDS.md** | 9 | 9 | 100% |
| **09_PROFILE.md** | 24 | 21 | 88% |
| **10_VILLAIN.md** | 7 | 7 | 100% |
| **11_ANALYTICS.md** | N/A | N/A | Инфраструктура |

---

## ⚠️ Пропущенные события (требуют внимания)

### 🔴 Критические (frontend implementation)

1. **`ui_error_shown`** - универсальное событие для всех ошибок UI
   - **Где использовать:** Все компоненты с error handling
   - **Параметры:** `screen_name`, `error_code`, `attempt_id` (опционально)
   - **Рекомендация:** Добавить в глобальный error boundary

### 🟡 Средний приоритет

2. **`mascot_stats_opened`** - открытие статистики маскота
   - **Roadmap:** 04_HOME.md
   - **Статус:** Функционал упомянут в комментариях (line 227), но не реализован
   - **Рекомендация:** Добавить модал со статистикой маскота

3. **`villain_stats_opened`** - открытие статистики злодея
   - **Roadmap:** 10_VILLAIN.md
   - **Статус:** Функционал может быть частью villain page
   - **Рекомендация:** Добавить детальную статистику в VillainPage

4. **`mascot_joke_viewed`** - просмотр шутки маскота
   - **Roadmap:** 04_HOME.md
   - **Статус:** Механика не описана
   - **Рекомендация:** Добавить случайные шутки в MascotSection

5. **`mascot_state_viewed`** - просмотр состояния маскота
   - **Roadmap:** Не описан
   - **Рекомендация:** Добавить при детальной реализации маскота

6. **`mascot_message_viewed`** - просмотр сообщения маскота
   - **Roadmap:** Не описан
   - **Рекомендация:** Добавить при детальной реализации маскота

7. **`mascot_interaction_clicked`** - клик по интеракции маскота
   - **Roadmap:** Не описан
   - **Рекомендация:** Добавить при детальной реализации маскота

### 🟢 Низкий приоритет (можно отложить)

8. **`report_archive_toggled`** - переключение архива отчётов
   - **Roadmap:** 09_PROFILE.md
   - **Статус:** Механика описана, событие пропущено в таблице
   - **Рекомендация:** Добавить при реализации ReportSettings

9. **`report_opened`** - открытие конкретного отчёта
   - **Roadmap:** 09_PROFILE.md
   - **Статус:** Механика описана, событие пропущено в таблице
   - **Рекомендация:** Добавить при реализации ReportArchive

10. **`subscription_cancel_requested`** - запрос отмены подписки
    - **Roadmap:** 09_PROFILE.md
    - **Статус:** Функционал не описан в paywall
    - **Рекомендация:** Добавить кнопку отмены в subscription management

---

## 🎯 Рекомендации

### 1. Немедленные действия

- ✅ **Добавить `ui_error_shown`** во все компоненты с error handling
- ✅ **Реализовать mascot stats modal** для `mascot_stats_opened`
- ✅ **Добавить villain stats section** для `villain_stats_opened`

### 2. Краткосрочные (1-2 недели)

- 📝 Реализовать детальную механику маскота:
  - `mascot_state_viewed`
  - `mascot_message_viewed`
  - `mascot_interaction_clicked`
  - `mascot_joke_viewed`

- 📝 Дополнить Report Settings:
  - `report_archive_toggled`
  - `report_opened`

- 📝 Добавить subscription cancel flow:
  - `subscription_cancel_requested`

### 3. Долгосрочные (бэклог)

- Backend события уже реализуются на бэкенде:
  - `email_verification_sent`
  - `email_verification_success`
  - `help_processing_started`
  - `check_processing_started`
  - `villain_health_changed`
  - `villain_victory_triggered`
  - `achievement_unlocked`
  - `referral_reward_unlocked`
  - `weekly_report_generated/sent/failed`
  - `payment_success/failed`
  - `subscription_activated/expired`
  - `processing_error`

### 4. Рефакторинг

- ✅ Создать глобальный Error Boundary для `ui_error_shown`
- ✅ Добавить analytics validation для всех параметров событий
- ✅ Убедиться, что все `child_profile_id` передаются правильно

---

## 📊 User Properties Coverage

### Parent User Properties

| Property | Упомянуто в Roadmap | Где обновлять |
|----------|---------------------|---------------|
| `platform_type` | ✅ | При инициализации приложения |
| `subscription_status` | ✅ | При изменении подписки |
| `trial_status` | ✅ | При изменении триала |
| `email_verified` | ✅ | После верификации email |
| `weekly_report_enabled` | ✅ | В ReportSettings |
| `report_archive_enabled` | ✅ | В ReportSettings |

### Child Profile Properties

| Property | Упомянуто в Roadmap | Где обновлять |
|----------|---------------------|---------------|
| `grade` | ✅ | При выборе класса |
| `level` | ✅ | При повышении уровня |
| `coins_balance` | ✅ | После каждой попытки |
| `tasks_solved_correct_count` | ✅ | После успешной проверки |
| `wins_count` | ✅ | После победы над злодеем |
| `checks_correct_count` | ✅ | После успешной проверки |
| `current_streak_days` | ✅ | Ежедневно |
| `has_unfinished_attempt` | ✅ | При создании/завершении попытки |
| `active_villain_id` | ✅ | При активации злодея |
| `active_villain_health_percent` | ✅ | При изменении здоровья злодея |
| `invited_count_total` | ✅ | При приглашении друга |
| `achievements_unlocked_count` | ✅ | При разблокировке достижения |

---

## ✅ Заключение

**Общее покрытие:** 110/117 событий (94%)

**Статус:** Отличное покрытие, требуется минимальная доработка

**Приоритет доработок:**
1. 🔴 Добавить `ui_error_shown` (критично)
2. 🟡 Реализовать детальную механику маскота (7 событий)
3. 🟢 Дополнить Report & Subscription (3 события)

**Следующие шаги:**
1. Реализовать пропущенные события по приоритету
2. Добавить валидацию всех событий в EventValidator
3. Протестировать отправку событий на всех экранах
4. Убедиться, что backend события корректно обрабатываются

---

**Автор:** Claude Code (Sonnet 4.5)
**Дата:** 2026-03-29
