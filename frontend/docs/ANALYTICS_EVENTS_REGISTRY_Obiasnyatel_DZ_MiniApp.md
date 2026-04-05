# ANALYTICS_EVENTS_REGISTRY — Объяснятель ДЗ mini app

## 1. Назначение
Документ фиксирует реестр аналитических событий для mini app «Объяснятель ДЗ».
Цель — дать команде фронта, бэка и аналитики единый источник истины:
- какие события отправляются,
- когда именно они срабатывают,
- какие параметры обязательны,
- какие параметры опциональны,
- кто является отправителем события.

---

## 2. Общие правила

### 2.1. Базовые поля для большинства событий
Эти поля рекомендуется передавать почти во все события, если они доступны:

- `event_name`
- `event_time`
- `parent_user_id`
- `child_profile_id`
- `platform_type` (`vk` / `max` / `web`)
- `session_id`
- `app_version`
- `screen_name`
- `entry_point`
- `grade`
- `subscription_status`
- `trial_status`

Если событие связано с попыткой:
- `attempt_id`
- `mode` (`help` / `check`)
- `attempt_status`
- `scenario_type`

Если событие связано с изображением:
- `attempt_image_id`
- `image_role`
- `upload_source`

---

## 3. Реестр событий

| Event name | Когда срабатывает | Обязательные параметры | Необязательные параметры | Кто шлёт |
|---|---|---|---|---|
| `onboarding_opened` | Открыт первый экран онбординга | `platform_type`, `session_id` | `entry_point` | front |
| `registration_opened` | Открыт экран регистрации | `platform_type`, `session_id` |  | front |
| `consent_screen_opened` | Открыт экран согласий | `platform_type`, `session_id` |  | front |
| `grade_selected` | Пользователь выбрал класс | `grade` | `child_profile_id` | front |
| `avatar_selected` | Выбран аватар | `avatar_id` | `child_profile_id` | front |
| `display_name_entered` | Введено имя профиля | `name_length` | `child_profile_id` | front |
| `adult_consent_checked` | Поставлена/снята галочка согласия | `checked` |  | front |
| `privacy_policy_opened` | Открыта политика конфиденциальности | `platform_type`, `session_id` | `policy_version` | front |
| `privacy_policy_accepted` | Принята политика | `policy_version` |  | front |
| `terms_opened` | Открыто пользовательское соглашение | `platform_type`, `session_id` | `terms_version` | front |
| `terms_accepted` | Принято пользовательское соглашение | `terms_version` |  | front |
| `email_entered` | Введён e-mail | `email_domain` |  | front |
| `email_verification_sent` | Отправлено письмо подтверждения | `email_domain` |  | back |
| `email_verification_success` | E-mail подтверждён | `parent_user_id` | `email_domain` | back |
| `onboarding_completed` | Онбординг завершён | `parent_user_id`, `child_profile_id` | `grade` | front |

| `home_opened` | Открыта Главная | `child_profile_id` | `entry_point` | front |
| `level_bar_viewed` | Верхняя строка прогресса показана | `child_profile_id`, `level` | `level_progress_percent` | front |
| `coins_balance_viewed` | Показан баланс монеток | `child_profile_id`, `coins_balance` |  | front |
| `tasks_correct_count_viewed` | Показан счётчик правильно выполненных заданий | `child_profile_id`, `tasks_solved_correct_count` |  | front |
| `home_help_clicked` | Нажата кнопка «Помоги разобраться» | `child_profile_id` |  | front |
| `home_check_clicked` | Нажата кнопка «Проверка ДЗ» | `child_profile_id` |  | front |
| `unfinished_attempt_modal_shown` | Показан модал незавершённой попытки | `child_profile_id`, `attempt_id` | `mode` | front |
| `unfinished_attempt_continue_clicked` | Нажато «Продолжить» | `child_profile_id`, `attempt_id` | `mode` | front |
| `unfinished_attempt_new_task_clicked` | Нажато «Новое задание» | `child_profile_id`, `attempt_id` | `mode` | front |
| `mascot_clicked` | Клик по маскоту | `child_profile_id`, `mascot_id` | `mascot_state` | front |
| `villain_clicked` | Клик по злодею | `child_profile_id`, `villain_id` | `villain_state` | front |
| `mascot_stats_opened` | Открыта статистика маскота | `child_profile_id`, `mascot_id` |  | front |
| `villain_stats_opened` | Открыта статистика злодея | `child_profile_id`, `villain_id` |  | front |
| `mascot_joke_viewed` | Показана шутка маскота | `child_profile_id`, `mascot_id` | `joke_id` | front |
| `recent_attempt_clicked` | Клик по карточке последней попытки | `child_profile_id`, `attempt_id` | `history_status` | front |
| `recent_attempts_view_all_clicked` | Нажато «Смотреть все» в последних попытках | `child_profile_id` |  | front |

| `help_flow_started` | Начат сценарий помощи | `child_profile_id`, `mode` | `entry_point` | front |
| `help_source_picker_opened` | Открыт выбор источника изображения для помощи | `child_profile_id` |  | front |
| `help_choose_file_clicked` | Выбран файл для помощи | `child_profile_id` |  | front |
| `help_camera_clicked` | Выбрана камера для помощи | `child_profile_id` |  | front |
| `help_clipboard_clicked` | Выбрана вставка из буфера | `child_profile_id` |  | front |
| `help_dragdrop_used` | Использован drag&drop | `child_profile_id` |  | front |
| `help_image_selected` | Выбрано изображение для помощи | `child_profile_id`, `upload_source` | `file_size_bucket`, `mime_type` | front |
| `help_image_upload_started` | Старт загрузки изображения помощи | `child_profile_id`, `upload_source` | `attempt_id` | front |
| `help_image_upload_completed` | Изображение помощи загружено | `child_profile_id`, `attempt_id`, `attempt_image_id` | `file_size_bucket` | back |
| `help_image_upload_failed` | Ошибка загрузки изображения помощи | `child_profile_id`, `error_code` | `attempt_id` | back |
| `help_quality_check_shown` | Показан экран проверки качества фото помощи | `child_profile_id`, `attempt_id` | `attempt_image_id` | front |
| `help_quality_confirm_clicked` | Подтверждено качество фото помощи | `child_profile_id`, `attempt_id` | `attempt_image_id` | front |
| `help_quality_reshoot_clicked` | Выбран пересъём в помощи | `child_profile_id`, `attempt_id` | `attempt_image_id` | front |
| `help_quality_manual_crop_clicked` | Выбран ручной crop в помощи | `child_profile_id`, `attempt_id` | `attempt_image_id` | front |
| `help_crop_opened` | Открыт экран crop в помощи | `child_profile_id`, `attempt_id` | `attempt_image_id` | front |
| `help_crop_confirmed` | Подтвержден crop в помощи | `child_profile_id`, `attempt_id` | `attempt_image_id` | front |
| `help_crop_reshoot_clicked` | После crop выбран пересъём | `child_profile_id`, `attempt_id` | `attempt_image_id` | front |
| `help_processing_started` | Старт обработки в помощи | `child_profile_id`, `attempt_id` |  | back |
| `help_long_wait_shown` | Показан long-wait в помощи | `child_profile_id`, `attempt_id` | `duration_seconds` | front |
| `help_save_and_wait_clicked` | Нажато «Сохранить и подождать» в помощи | `child_profile_id`, `attempt_id` |  | front |
| `help_retry_clicked` | Нажато «Повторить» в помощи | `child_profile_id`, `attempt_id` |  | front |
| `help_cancel_clicked` | Нажато «Отменить» в помощи | `child_profile_id`, `attempt_id` |  | front |
| `help_result_opened` | Открыт экран результата помощи | `child_profile_id`, `attempt_id` | `used_hints_count` | front |
| `hint_opened` | Открыта подсказка | `child_profile_id`, `attempt_id`, `hint_level` |  | front |
| `next_hint_clicked` | Нажата «Следующая подсказка» | `child_profile_id`, `attempt_id`, `from_hint_level` |  | front |
| `answer_submit_clicked` | Нажата кнопка отправки ответа из помощи | `child_profile_id`, `attempt_id` | `used_hints_count` | front |
| `new_task_clicked_from_help` | Нажато «Новое задание» из помощи | `child_profile_id`, `attempt_id` |  | front |

| `check_flow_started` | Начат сценарий проверки ДЗ | `child_profile_id`, `mode` | `entry_point` | front |
| `check_scenario_picker_opened` | Открыт экран выбора сценария проверки | `child_profile_id` |  | front |
| `check_scenario_selected` | Выбран сценарий проверки | `child_profile_id`, `check_scenario` |  | front |
| `check_single_photo_selected` | Выбран сценарий одного фото | `child_profile_id`, `check_scenario` |  | front |
| `check_two_photo_selected` | Выбран сценарий двух фото | `child_profile_id`, `check_scenario` |  | front |
| `check_task_image_selected` | Выбрано фото задания в сценарии двух фото | `child_profile_id`, `attempt_id`, `image_role` | `upload_source` | front |
| `check_answer_image_selected` | Выбрано фото ответа в сценарии двух фото | `child_profile_id`, `attempt_id`, `image_role` | `upload_source` | front |
| `check_source_picker_opened` | Открыт выбор источника изображения для проверки | `child_profile_id` | `check_scenario` | front |
| `check_choose_file_clicked` | Выбран файл для проверки | `child_profile_id` | `check_scenario` | front |
| `check_camera_clicked` | Выбрана камера для проверки | `child_profile_id` | `check_scenario` | front |
| `check_clipboard_clicked` | Выбрана вставка из буфера для проверки | `child_profile_id` | `check_scenario` | front |
| `check_dragdrop_used` | Использован drag&drop в проверке | `child_profile_id` | `check_scenario` | front |
| `check_image_upload_started` | Старт загрузки изображения проверки | `child_profile_id`, `attempt_id`, `image_role` | `upload_source` | front |
| `check_image_upload_completed` | Изображение проверки загружено | `child_profile_id`, `attempt_id`, `attempt_image_id`, `image_role` | `file_size_bucket` | back |
| `check_image_upload_failed` | Ошибка загрузки изображения проверки | `child_profile_id`, `attempt_id`, `image_role`, `error_code` |  | back |
| `upload_more_clicked` | Нажата кнопка «Загрузить ещё» | `child_profile_id`, `mode`, `image_role_expected` | `attempt_id` | front |
| `check_quality_check_shown` | Показан экран проверки качества фото в проверке | `child_profile_id`, `attempt_id` | `attempt_image_id`, `image_role` | front |
| `check_quality_confirm_clicked` | Подтверждено качество фото в проверке | `child_profile_id`, `attempt_id` | `attempt_image_id`, `image_role` | front |
| `check_quality_reshoot_clicked` | Выбран пересъём в проверке | `child_profile_id`, `attempt_id` | `attempt_image_id`, `image_role` | front |
| `check_quality_manual_crop_clicked` | Выбран ручной crop в проверке | `child_profile_id`, `attempt_id` | `attempt_image_id`, `image_role` | front |
| `check_crop_opened` | Открыт экран crop в проверке | `child_profile_id`, `attempt_id` | `attempt_image_id`, `image_role` | front |
| `check_crop_confirmed` | Подтвержден crop в проверке | `child_profile_id`, `attempt_id` | `attempt_image_id`, `image_role` | front |
| `check_crop_reshoot_clicked` | После crop выбран пересъём в проверке | `child_profile_id`, `attempt_id` | `attempt_image_id`, `image_role` | front |
| `check_processing_started` | Старт обработки проверки | `child_profile_id`, `attempt_id`, `check_scenario` |  | back |
| `check_long_wait_shown` | Показан long-wait в проверке | `child_profile_id`, `attempt_id` | `duration_seconds` | front |
| `check_save_and_wait_clicked` | Нажато «Сохранить и подождать» в проверке | `child_profile_id`, `attempt_id` |  | front |
| `check_retry_clicked` | Нажато «Повторить» в проверке | `child_profile_id`, `attempt_id` |  | front |
| `check_cancel_clicked` | Нажато «Отменить» в проверке | `child_profile_id`, `attempt_id` |  | front |
| `check_result_opened` | Открыт результат проверки | `child_profile_id`, `attempt_id`, `result_status` | `error_count` | front |
| `check_error_feedback_viewed` | Пользователь увидел блок ошибки | `child_profile_id`, `attempt_id`, `error_count` |  | front |
| `retry_after_errors_clicked` | Нажато повторить после ошибки | `child_profile_id`, `attempt_id` |  | front |
| `fixed_and_resubmit_clicked` | Нажато «Исправил» | `child_profile_id`, `attempt_id` |  | front |
| `new_task_clicked_from_check` | Нажато «Новое задание» из проверки | `child_profile_id`, `attempt_id` |  | front |
| `soft_error_message_shown` | Показана мягкая формулировка о наличии ошибки | `child_profile_id`, `attempt_id`, `error_count` |  | front |
| `error_hint_block_opened` | Открыт блок с указанием ошибки | `child_profile_id`, `attempt_id`, `error_block_id` | `step_number`, `line_reference` | front |
| `error_location_viewed` | Просмотрена локализация ошибки | `child_profile_id`, `attempt_id`, `location_type` |  | front |
| `user_retries_after_error` | Пользователь пытается ещё раз после ошибки | `child_profile_id`, `attempt_id` |  | front |
| `user_abandons_after_error` | Пользователь бросает сценарий после ошибки | `child_profile_id`, `attempt_id` |  | front |

| `villain_screen_opened` | Открыт экран злодея | `child_profile_id`, `villain_id` |  | front |
| `villain_taunt_viewed` | Показана реплика злодея | `child_profile_id`, `villain_id` |  | front |
| `villain_health_changed` | Здоровье злодея изменилось | `child_profile_id`, `villain_id`, `health_before`, `health_after`, `damage_amount` | `reason` | back |
| `villain_victory_triggered` | Сработала победа над злодеем | `child_profile_id`, `villain_id`, `attempt_id` |  | back |
| `victory_screen_opened` | Открыт экран победы | `child_profile_id`, `villain_id`, `attempt_id` |  | front |
| `victory_reward_viewed` | Пользователь увидел награду за победу | `child_profile_id`, `villain_id`, `reward_type` | `reward_id` | front |
| `victory_continue_clicked` | Нажато продолжение после победы | `child_profile_id`, `villain_id` |  | front |

| `mascot_state_viewed` | Показано состояние маскота | `child_profile_id`, `mascot_id`, `mascot_state` |  | front |
| `mascot_message_viewed` | Показано сообщение маскота | `child_profile_id`, `mascot_id` | `message_id` | front |
| `mascot_interaction_clicked` | Клик по интеракции маскота | `child_profile_id`, `mascot_id`, `interaction_type` |  | front |

| `achievements_opened` | Открыт экран достижений | `child_profile_id` |  | front |
| `achievement_shelf_viewed` | Просмотрена полка достижений | `child_profile_id`, `shelf_order` |  | front |
| `achievement_clicked` | Клик по достижению | `child_profile_id`, `achievement_id`, `is_unlocked` |  | front |
| `achievement_detail_opened` | Открыта карточка достижения | `child_profile_id`, `achievement_id` | `is_unlocked` | front |
| `achievement_reward_viewed` | Просмотрена награда открытого достижения | `child_profile_id`, `achievement_id` |  | front |
| `locked_achievement_requirement_viewed` | Просмотрено условие закрытого достижения | `child_profile_id`, `achievement_id` |  | front |
| `achievement_unlocked` | Достижение открыто | `child_profile_id`, `achievement_id`, `unlock_reason` |  | back |

| `friends_opened` | Открыт экран друзей | `child_profile_id` |  | front |
| `friends_reward_offer_viewed` | Просмотрен оффер награды за друзей | `child_profile_id`, `target_count`, `current_count` |  | front |
| `invite_friend_clicked` | Нажата кнопка приглашения друга | `child_profile_id` | `referral_code` | front |
| `referral_link_copied` | Скопирована реферальная ссылка | `child_profile_id`, `referral_code` |  | front |
| `referral_share_opened` | Открыт share sheet | `child_profile_id`, `referral_code` |  | front |
| `referral_share_sent` | Подтверждена отправка приглашения | `child_profile_id`, `referral_code` | `channel_type` | front |
| `referral_progress_viewed` | Просмотрен прогресс приглашений | `child_profile_id`, `invited_count_total`, `target_count` |  | front |
| `referral_reward_unlocked` | Разблокирована награда за друзей | `child_profile_id`, `reward_type` | `reward_id` | back |
| `referral_reward_claimed` | Получена награда за друзей | `child_profile_id`, `reward_type` | `reward_id` | front |

| `profile_opened` | Открыт Профиль | `child_profile_id` |  | front |
| `profile_history_opened` | Открыт блок Истории из Профиля | `child_profile_id` |  | front |
| `profile_report_settings_opened` | Открыт блок отчёта родителю | `parent_user_id` |  | front |
| `profile_support_opened` | Открыт блок поддержки | `parent_user_id` |  | front |
| `profile_parent_gate_opened` | Открыт взрослый контур / parent gate | `parent_user_id` |  | front |

| `history_opened` | Открыта История | `child_profile_id` |  | front |
| `history_item_clicked` | Клик по карточке истории | `child_profile_id`, `attempt_id`, `history_status` |  | front |
| `history_filter_used` | Использован фильтр истории | `child_profile_id`, `filter_type`, `filter_value` |  | front |
| `history_detail_opened` | Открыта детальная карточка попытки | `child_profile_id`, `attempt_id` |  | front |
| `history_retry_clicked` | Нажато повторить из истории | `child_profile_id`, `attempt_id` |  | front |
| `history_fix_and_recheck_clicked` | Нажато «Исправить и проверить» из истории | `child_profile_id`, `attempt_id` |  | front |

| `report_settings_opened` | Открыты настройки отчёта | `parent_user_id` |  | front |
| `report_email_changed` | Изменён e-mail отчёта | `parent_user_id`, `email_domain` |  | front |
| `weekly_report_toggled` | Включён/выключен weekly report | `parent_user_id`, `enabled` |  | front |
| `report_archive_toggled` | Включён/выключен архив отчётов | `parent_user_id`, `enabled` |  | front |
| `report_archive_opened` | Открыт архив отчётов | `parent_user_id` |  | front |
| `report_opened` | Открыт конкретный отчёт | `parent_user_id`, `report_id` |  | front |
| `report_download_clicked` | Нажата загрузка отчёта | `parent_user_id`, `report_id` |  | front |
| `weekly_report_generated` | Система сформировала weekly report | `parent_user_id`, `report_id` | `period_start`, `period_end` | system |
| `weekly_report_sent` | Система отправила weekly report | `parent_user_id`, `report_id` | `email_domain` | system |
| `weekly_report_failed` | Ошибка генерации или отправки отчёта | `parent_user_id`, `report_id`, `error_code` |  | system |

| `paywall_opened` | Открыт paywall | `parent_user_id`, `entry_point`, `blocked_feature` |  | front |
| `pricing_opened` | Открыт экран тарифов | `parent_user_id` |  | front |
| `plan_selected` | Выбран тариф | `parent_user_id`, `billing_plan_id`, `billing_period`, `price_amount` |  | front |
| `payment_started` | Старт оплаты | `parent_user_id`, `billing_plan_id` | `amount`, `currency` | front |
| `payment_success` | Успешная оплата | `parent_user_id`, `billing_plan_id`, `amount` | `currency`, `provider` | back |
| `payment_failed` | Ошибка оплаты | `parent_user_id`, `billing_plan_id`, `error_code` | `amount`, `currency`, `provider` | back |
| `subscription_activated` | Подписка активирована | `parent_user_id`, `billing_plan_id`, `subscription_status` | `expires_at` | back |
| `subscription_cancel_requested` | Запрошена отмена подписки | `parent_user_id`, `subscription_id` |  | front |
| `subscription_expired` | Подписка истекла | `parent_user_id`, `subscription_id` | `expired_at` | back |

| `support_opened` | Открыта поддержка | `parent_user_id` | `screen_name` | front |
| `support_message_sent` | Отправлено сообщение в поддержку | `parent_user_id` | `message_length` | front |
| `ui_error_shown` | На экране показана UI-ошибка | `screen_name`, `error_code` | `attempt_id` | front |
| `processing_error` | Ошибка обработки попытки | `attempt_id`, `error_code`, `stage` | `mode`, `scenario_type` | back |

---

## 4. Пользовательские свойства (user properties)

### Для `parent_user`
- `platform_type`
- `subscription_status`
- `trial_status`
- `email_verified`
- `weekly_report_enabled`
- `report_archive_enabled`

### Для `child_profile`
- `grade`
- `level`
- `coins_balance`
- `tasks_solved_correct_count`
- `wins_count`
- `checks_correct_count`
- `current_streak_days`
- `has_unfinished_attempt`
- `active_villain_id`
- `active_villain_health_percent`
- `invited_count_total`
- `achievements_unlocked_count`

---

## 5. Принципы ownership

### Front отправляет:
- экранные события
- клики
- выбор сценариев
- открытие блоков
- локальные пользовательские действия

### Back / System отправляет:
- статусы обработки
- успешность/ошибки обработки
- изменение здоровья злодея
- выдачу наград
- создание и отправку отчётов
- успешность/ошибки оплаты
- активацию/истечение подписки

---

## 6. Следующий слой (можно делать отдельно)
На основе этого реестра потом стоит собрать:
1. `user_properties_registry`
2. `metrics_definition_sheet`
3. список витрин / dashboard metrics
