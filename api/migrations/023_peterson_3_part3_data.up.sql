-- Данные учебника Петерсон 3 класс, часть 3
-- PostgreSQL compatible

INSERT INTO textbooks (subject, grade, authors, title, part, year, publisher, source_url)
VALUES ('math', 3, 'Петерсон Л.Г.', 'Математика 3 класс', 3, 2022, 'Просвещение', 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik')
ON CONFLICT (subject, grade, authors, part) DO UPDATE SET year = EXCLUDED.year;

-- Задачи
DO $$
DECLARE
    v_textbook_id BIGINT;
    v_task_id BIGINT;
BEGIN
    SELECT id INTO v_textbook_id FROM textbooks WHERE subject = 'math' AND grade = 3 AND authors = 'Петерсон Л.Г.' AND part = 3;

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 3, '1', 0, 'Миша прошёл на лыжах расстояние 80 м за 20 с, а Игорь – 45 м за 15 с. Кто из них прошёл большее расстояние, а кто – меньшее? Кто шёл больше времени, а кто – меньше? Кто шёл быстрее, а кто – медленнее? Какие величины характеризуют движение объектов?', '</p> \n<p class="text">Миша прошёл на лыжах расстояние 80 м за 20 с, а Игорь – 45 м за 15 с. Кто из них прошёл большее расстояние, а кто – меньшее? Кто шёл больше времени, а кто – меньше? Кто шёл быстрее, а кто – медленнее? Какие величины характеризуют движение объектов?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica3-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 3, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 3, номер 1, год 2022."/>\n</div>\n</div>', 'Из них прошёл большее расстояние Миша 80 м, а Игорь 45 м – меньшее. Миша шёл больше времени 20 с, а Игорь 15 с – меньше. Миша 80 : 20 = 4 (м/с) шёл быстрее, а Игорь 45 : 15 = 3 (м/с) – медленнее. Скорость, расстояние и время – величины характеризуют движение объектов.', '<p>\nИз них прошёл большее расстояние Миша 80 м, а Игорь 45 м – меньшее. Миша шёл больше времени 20 с, а Игорь 15 с – меньше. Миша 80 : 20 = 4 (м/с) шёл быстрее, а Игорь 45 : 15 = 3 (м/с) – медленнее. Скорость, расстояние и время – величины характеризуют движение объектов.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-3/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica3-nomer1.jpg', 'peterson/3/part3/page3/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd3d288c60de357eecaab84537fa7790c7be6a268c16e8df29be1a26f377565c8', '15,20,45,80', '["больше","меньше"]'::jsonb, 'миша прошёл на лыжах расстояние 80 м за 20 с, а игорь-45 м за 15 с. кто из них прошёл большее расстояние, а кто-меньшее? кто шёл больше времени, а кто-меньше? кто шёл быстрее, а кто-медленнее? какие величины характеризуют движение объектов');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 3, '2', 1, 'Объясни смысл предложений: а) Самолёт летит со скоростью 800 км/ч. б) Скорость теплохода 45 км/ч. в) Человек идёт со скоростью 4 км/ч. г) Земля движется по орбите со скоростью 30 км/с. д) Черепаха ползёт со скоростью 4 м/мин.', '</p> \n<p class="text">Объясни смысл предложений:</p> \n\n<p class="description-text"> \nа) Самолёт летит со скоростью 800 км/ч.<br/>\nб) Скорость теплохода 45 км/ч.<br/>\nв) Человек идёт со скоростью 4 км/ч.<br/>\nг) Земля движется по орбите со скоростью 30 км/с.<br/>\nд) Черепаха ползёт со скоростью 4 м/мин.\n\n</p>\n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica3-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 3, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 3, номер 2, год 2022."/>\n</div>\n</div>', 'а) Самолёт пролетел расстояние 800 км за 1 ч. б) Теплоход проплыл 45 км за 1 ч. в) Человек проходит 4 км за 1 ч. г) Земля движется по орбите преодолевая расстояние 30 км в 1 с. д) Черепаха ползёт расстояние 4 м за 1 мин.', '<p>\nа) Самолёт пролетел расстояние 800 км за 1 ч.<br/>\nб) Теплоход проплыл 45 км за 1 ч.<br/>\nв) Человек проходит 4 км за 1 ч.<br/>\nг) Земля движется по орбите преодолевая расстояние 30 км в 1 с.<br/>\nд) Черепаха ползёт расстояние 4 м за 1 мин.\n</p>', 'Скорость. Время. Расстояние. При движении автомобиля, автобуса, поезда нас интересует пройденное ими расстояние, время в пути и скорость (быстрее или медленнее они едут). Под расстоянием, пройденным движущимся объектом, мы будем понимать длину дороги, соединяющей начало и конец пути. Скоростью мы будем называть расстояние, пройденное в единицу времени. Скорость является величиной. В качестве единиц измерения скорости используют такие единицы, как метр в секунду (м/c), метр в минуту (м/мин), километр в час (км/ч) и т. д.', '<div class="recomended-block">\n<span class="title">Скорость. Время. Расстояние.</span>\n<p>\n\nПри движении автомобиля, автобуса, поезда нас интересует пройденное ими расстояние,  время в пути и скорость (быстрее или медленнее они едут).<br/>\nПод расстоянием, пройденным движущимся объектом, мы будем понимать длину дороги, соединяющей начало и конец пути.<br/>\nСкоростью мы будем называть расстояние, пройденное в единицу времени.<br/>\nСкорость является величиной. В качестве единиц измерения скорости используют такие единицы, как метр в секунду (м/c), метр в минуту (м/мин), километр в час (км/ч) и т. д.\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica3-spravka.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 3, справка, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 3, справка, год 2022."/>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-3/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica3-nomer2.jpg', 'peterson/3/part3/page3/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '55ba92b6e982ae79aa3aac41fb8d623c1f9b47885ba112d7dd66f071544ebb90', '4,30,45,800', NULL, 'объясни смысл предложений:а) самолёт летит со скоростью 800 км/ч. б) скорость теплохода 45 км/ч. в) человек идёт со скоростью 4 км/ч. г) земля движется по орбите со скоростью 30 км/с. д) черепаха ползёт со скоростью 4 м/мин');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 4, '3', 0, 'Найди: а) Скорость космического корабля, если он пролетел 56 км за 8 с. б) Скорость плота на реке, если он за 4 ч проплыл 16 км. в) Скорость автобуса, если он прошёл 120 км за 3 ч. г) Скорость велосипедиста, если он проехал 36 км за 2 ч.', '</p> \n<p class="text">Найди:<br/>\nа) Скорость космического корабля, если он пролетел 56 км за 8 с.<br/>\nб) Скорость плота на реке, если он за 4 ч проплыл 16 км.<br/>\nв) Скорость автобуса, если он прошёл 120 км за 3 ч.<br/>\nг) Скорость велосипедиста, если он проехал 36 км за 2 ч.\n</p>', 'а) 56 : 8 = 7 (км/с) б) 16 : 4 = 4 (км/ч) в) 120 : 3 = 40 (км/ч) г) 36 : 2 = 18 (км/ч)', '<p>\nа) 56 : 8 = 7 (км/с)<br/>\nб) 16 : 4 = 4 (км/ч)<br/>\nв) 120 : 3 = 40 (км/ч)<br/>\nг) 36 : 2 = 18 (км/ч)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-4/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8fdb4fc5ee6c389f2845b97ce8dc22d386157550c331994f639b798cf42af5dc', '2,3,4,8,16,36,56,120', '["найди"]'::jsonb, 'найди:а) скорость космического корабля, если он пролетел 56 км за 8 с. б) скорость плота на реке, если он за 4 ч проплыл 16 км. в) скорость автобуса, если он прошёл 120 км за 3 ч. г) скорость велосипедиста, если он проехал 36 км за 2 ч');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 4, '4', 1, 'Найди карточки, на которых указана скорость: а) самолёта; б) поезда; в) автомобиля; г) пешехода; д) велосипедиста; е) ракеты. Сделай по желанию рисунок и подпиши значение скорости.', '</p> \n<p class="text">Найди карточки, на которых указана скорость: а) самолёта; б) поезда; в) автомобиля; г) пешехода; д) велосипедиста; е) ракеты. Сделай по желанию рисунок и подпиши значение скорости.</p> \n\n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica4-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 4, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 4, номер 4, год 2022."/>\n</div>\n\n</div>', 'а) 900 км/ч б) 60 км/ч в) 90 км/ч г) 5 км/ч д) 20 км/ч е) 6 км/с', '<p>\nа) 900 км/ч<br/>\nб) 60 км/ч<br/>\nв) 90 км/ч<br/>\nг) 5 км/ч<br/>\nд) 20 км/ч<br/>\nе) 6 км/с\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-4/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica4-nomer4.jpg', 'peterson/3/part3/page4/task4_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '54169a6611b45eaab7913155145ed513a49ecc71755de13e28e400c520dbf421', NULL, '["найди"]'::jsonb, 'найди карточки, на которых указана скорость:а) самолёта; б) поезда; в) автомобиля; г) пешехода; д) велосипедиста; е) ракеты. сделай по желанию рисунок и подпиши значение скорости');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 4, '5', 2, 'а) Поезд прошёл 224 км за 4 часа. Его скорость в 3 раза меньше скорости вертолёта. Чему равна скорость вертолёта? б) Плот проплыл 27 км за 9 ч, а моторная лодка – 24 км за 2 ч. Чья скорость больше и на сколько?', '</p> \n<p class="text">а) Поезд прошёл 224 км за 4 часа. Его скорость в 3 раза меньше скорости вертолёта. Чему равна скорость вертолёта?<br/>\nб) Плот проплыл 27 км за 9 ч, а моторная лодка – 24 км за 2 ч. Чья скорость больше и на сколько?\n</p>', 'а) 224 : 4 · 3 = 56 · 3 = 168 (км/ч) Ответ: 168 км/ч скорость вертолёта. б) 24 : 2 - 27 : 9 = 12 - 3 = 9 (км/ч) Ответ: на 9 км/ч больше скорость моторной лодки чем скорость плота.', '<p>\nа) 224 : 4 · 3 = 56 · 3 = 168 (км/ч)<br/>\n<b>Ответ:</b> 168 км/ч скорость вертолёта. <br/><br/>\nб) 24 : 2 - 27 : 9 = 12 - 3 = 9 (км/ч)<br/>\n<b>Ответ:</b> на 9 км/ч больше скорость моторной лодки чем скорость плота.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Поезд прошёл 224 км за 4 часа. Его скорость в 3 раза меньше скорости вертолёта. Чему равна скорость вертолёта?","solution":"224 : 4 · 3 = 56 · 3 = 168 (км/ч) Ответ: 168 км/ч скорость вертолёта."},{"letter":"б","condition":"Плот проплыл 27 км за 9 ч, а моторная лодка – 24 км за 2 ч. Чья скорость больше и на сколько?","solution":"24 : 2 - 27 : 9 = 12 - 3 = 9 (км/ч) Ответ: на 9 км/ч больше скорость моторной лодки чем скорость плота."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-4/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '44a2e5ea82c22c92f81be24b6c02824874c7bc9a5e21d8803f10a27a5024fe06', '2,3,4,9,24,27,224', '["больше","меньше","раз","раза"]'::jsonb, 'а) поезд прошёл 224 км за 4 часа. его скорость в 3 раза меньше скорости вертолёта. чему равна скорость вертолёта? б) плот проплыл 27 км за 9 ч, а моторная лодка-24 км за 2 ч. чья скорость больше и на сколько');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 4, '6', 3, 'а) Грузовая машина за 8 ч прошла 280 км, а легковая машина это же расстояние – за 4 ч. Во сколько раз скорость грузовой машины меньше скорости легковой? б) Велосипедист за 3 ч проехал 57 км, а мотоциклист за 2 ч проехал на 71 км больше. На сколько километров в час скорость велосипедиста меньше скорости мотоциклиста?', '</p> \n<p class="text">а) Грузовая машина за 8 ч прошла 280 км, а легковая машина это же расстояние – за 4 ч. Во сколько раз скорость грузовой машины меньше скорости легковой? <br/>\nб) Велосипедист за 3 ч проехал 57 км, а мотоциклист за 2 ч проехал на 71 км больше. На сколько километров в час скорость велосипедиста меньше скорости мотоциклиста?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica4-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 4, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 4, номер 6, год 2022."/>\n</div>\n</div>', 'а) 280 : 4 : 280 : 8 = 70 : 35 = 2 (раза) Ответ: в 2 раза скорость грузовой машины меньше скорости легковой. б) (57 + 71) : 2 – 57 : 3 = 128 : 2 - 19 = 64 – 19 = 45 (км/ч) Ответ: на 45 километров в час скорость велосипедиста меньше скорости мотоциклиста.', '<p>\nа) 280 : 4 : 280 : 8 = 70 : 35 = 2 (раза)<br/>\n<b>Ответ:</b> в 2 раза скорость грузовой машины меньше скорости легковой.<br/><br/>\nб) (57 + 71) : 2 – 57 : 3 = 128 : 2 - 19 = 64 – 19 = 45 (км/ч)<br/>\n<b>Ответ:</b> на 45 километров в час скорость велосипедиста меньше скорости мотоциклиста.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Грузовая машина за 8 ч прошла 280 км, а легковая машина это же расстояние – за 4 ч. Во сколько раз скорость грузовой машины меньше скорости легковой?","solution":"280 : 4 : 280 : 8 = 70 : 35 = 2 (раза) Ответ: в 2 раза скорость грузовой машины меньше скорости легковой."},{"letter":"б","condition":"Велосипедист за 3 ч проехал 57 км, а мотоциклист за 2 ч проехал на 71 км больше. На сколько километров в час скорость велосипедиста меньше скорости мотоциклиста?","solution":"(57 + 71) : 2 – 57 : 3 = 128 : 2 - 19 = 64 – 19 = 45 (км/ч) Ответ: на 45 километров в час скорость велосипедиста меньше скорости мотоциклиста."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-4/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica4-nomer6.jpg', 'peterson/3/part3/page4/task6_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c176591c08740667015d2c8edc7dff4e195af93730b1efc76c26fcca6553abb4', '2,3,4,8,57,71,280', '["больше","меньше","раз"]'::jsonb, 'а) грузовая машина за 8 ч прошла 280 км, а легковая машина это же расстояние-за 4 ч. во сколько раз скорость грузовой машины меньше скорости легковой? б) велосипедист за 3 ч проехал 57 км, а мотоциклист за 2 ч проехал на 71 км больше. на сколько километров в час скорость велосипедиста меньше скорости мотоциклиста');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 4, '7', 4, 'Реши уравнения с комментированием и сделай проверку: а) (40 · x) : 10 = 28    б) y : 9 - 28 = 32    в) 39 + 490 : k = 46', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (40 · x) : 10 = 28    б) y : 9 - 28 = 32    в) 39 + 490 : k = 46\n</p>', 'а) (40 · x) : 10 = 28 чтобы найти сомножитель делимого х надо частное разделить на множитель делимого и умножить на делитель х = 28 : 40 · 10 х = 7 Проверка: (40 · 7) : 10 = 28 б) y : 9 - 28 = 32 Чтобы найти делимое уменьшаемого у надо найти произведение суммы разности с вычитаемым на делитель уменьшаемого у = (32 + 28) · 9 у = 60 · 9 у = 540 Проверка: 540 : 9 – 28 = 32 в) 39 + 490 : k = 46 Чтобы найти делитель слагаемого k надо найти частное делимого слагаемого и разности суммы на слагаемое k = 490 : (46 - 39) k = 490 : 7 k = 70 Проверка: 39 + 490 : 70 = 46', '<p>\nа) (40 · x) : 10 = 28 <br/> \nчтобы найти сомножитель делимого х надо частное разделить на множитель делимого и умножить на делитель<br/>\nх = 28 : 40 · 10<br/>\nх = 7<br/>   \n<b>Проверка:</b> (40 · 7) : 10 = 28<br/><br/>\nб) y : 9 - 28 = 32 <br/>       \nЧтобы найти делимое уменьшаемого у надо найти произведение суммы разности с вычитаемым на делитель уменьшаемого<br/>  \nу = (32 + 28) · 9<br/>\nу = 60 · 9<br/>\nу = 540<br/>\n<b>Проверка:</b> 540 : 9 – 28 = 32<br/><br/>\nв) 39 + 490 : k = 46<br/>\nЧтобы найти делитель слагаемого k надо найти частное делимого слагаемого и разности суммы на слагаемое<br/>\nk = 490 : (46 - 39)<br/>\nk = 490 : 7<br/>\nk = 70<br/>\n<b>Проверка:</b> 39 + 490 : 70 = 46\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-4/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3b43bf0dd17560dca89617430a60e8e2d8a3c4d256b305eff2a6c3cb2142f2d3', '9,10,28,32,39,40,46,490', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (40*x):10=28    б) y:9-28=32    в) 39+490:k=46');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 4, '8', 5, 'Выполни действия и сделай проверку: а) 547923 + 83699221        в) 90560 · 200 б) 4758036 - 50854            г) 3027600 : 6', '</p> \n<p class="text">Выполни действия и сделай проверку:</p> \n\n<p class="description-text"> \nа) 547923 + 83699221        в) 90560 · 200<br/>\nб) 4758036 - 50854            г) 3027600 : 6\n\n</p>', 'а) 547923 + 83699221 = 84247144 22 + 3 · 2 + 4 · 4 = 22 + 6 + 16 = 44 (ноги) Ответ: 44 ног гуляет теперь по двору.', '<p>\nа) 547923 + 83699221 = 84247144\n</p>\n\n<div class="img-wrapper-460">\n<img width="220" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica4-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 4, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 4, номер 8, год 2022."/>\n\n\n<p>\n22 + 3 · 2 + 4 · 4 = 22 + 6 + 16 = 44 (ноги)<br/>\n<b>Ответ:</b> 44 ног гуляет теперь по двору.\n\n</p>', 'Замечание: В задачах на движение будем считать, что скорость в течение всего времени движения не изменяется, а движение происходит по прямой дороге. Такое движение называют равномерным прямолинейным.', '<div class="recomended-block">\n<span class="title">Замечание:</span>\n<p>\nВ задачах на движение будем считать, что скорость в течение всего времени движения не изменяется, а движение происходит по прямой дороге. Такое движение называют равномерным прямолинейным.\n</p>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-4/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica4-nomer8.jpg', 'peterson/3/part3/page4/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '7c3307189bf941e92e7c2edb6fb03f0beeb55a5de577d3e453c651285b523c5f', '6,200,50854,90560,547923,3027600,4758036,83699221', NULL, 'выполни действия и сделай проверку:а) 547923+83699221        в) 90560*200 б) 4758036-50854            г) 3027600:6');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 5, '1', 0, 'Прочитай формулы. Что они означают? Какие ещё формулы ты знаешь? S = a · b        P = (a + b) · 2 V = a · b · c    a = b · c + r Зачем нужны формулы и как их устанавливают?', '</p> \n<p class="text">Прочитай формулы. Что они означают? Какие ещё формулы ты знаешь?<br/><br/>\nS = a · b        P = (a + b) · 2<br/>\nV = a · b · c    a = b · c + r<br/><br/>\nЗачем нужны формулы и как их устанавливают?\n</p>', 'Площадь равна произведению длины на ширину, Периметр равен удвоенной сумме длины и ширины, Объём равен произведению длины, ширины и высоты, А равно сумме произведения b на с и r. Формулы нужны для нахождения параметров фигур.', '<p>\nПлощадь равна произведению длины на ширину,<br/>\nПериметр равен удвоенной сумме длины и ширины,<br/>\nОбъём равен произведению длины, ширины и высоты,<br/>\nА равно сумме произведения b на с и r.<br/>\nФормулы нужны для нахождения параметров фигур.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-5/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e6f43d32bc7c71b417824ac1a4b38edd1acaba2501c00308296150270ccd25a9', '2', NULL, 'прочитай формулы. что они означают? какие ещё формулы ты знаешь? s=a*b        p=(a+b)*2 v=a*b*c    a=b*c+r зачем нужны формулы и как их устанавливают');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 5, '2', 1, 'Аэросани едут со скоростью v = 45 км/ч. Построй в тетради числовой луч и покажи на нём движение саней*. Какое расстояние преодолеют аэросани за 1 ч, 2 ч, 3 ч, 4 ч, t ч? Составь и заполни в тетради таблицу. Напиши формулу, выражающую зависимость пройденного расстояния s от времени t.', '</p> \n<p class="text">Аэросани едут со скоростью v = 45 км/ч. Построй в тетради числовой луч и покажи на нём движение саней*.</p> \n\n<div class="description-text">  \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica5-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 5, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 5, номер 2, год 2022."/>\n</div>\n</div>\n\n<p class="text">Какое расстояние преодолеют аэросани за 1 ч, 2 ч, 3 ч, 4 ч, t ч? Составь и заполни в тетради таблицу. Напиши формулу, выражающую зависимость пройденного расстояния s от времени t.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica5-nomer2-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 5, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 5, номер 2, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica5-nomer2-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 5, номер 2, ответ, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 5, ответ, номер 2, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-5/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica5-nomer2.jpg', 'peterson/3/part3/page5/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica5-nomer2-1.jpg', 'peterson/3/part3/page5/task2_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica5-nomer2-2.jpg', 'peterson/3/part3/page5/task2_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '74e8d10e50ee960b6a86dcb8136f6cae40b21b0209d832b12164f1bd15f0f84e', '1,2,3,4,45', '["заполни","числовой луч"]'::jsonb, 'аэросани едут со скоростью v=45 км/ч. построй в тетради числовой луч и покажи на нём движение саней*. какое расстояние преодолеют аэросани за 1 ч, 2 ч, 3 ч, 4 ч, t ч? составь и заполни в тетради таблицу. напиши формулу, выражающую зависимость пройденного расстояния s от времени t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 5, '3', 2, 'Проанализируй решение предыдущей задачи. Установи, как найти расстояние s, пройденное объектом, если он двигался со скоростью v в течение времени t.', '</p> \n<p class="text">Проанализируй решение предыдущей задачи. Установи, как найти расстояние  s, пройденное объектом, если он двигался со скоростью v в  течение времени t.</p>', 'Найдём расстояние s, пройденное объектом, если он двигался со скоростью v в течение времени t: s = v · t', '<p>\nНайдём расстояние s, пройденное объектом, если он двигался со скоростью v в течение времени t: s = v · t\n</p>', 'Формула пути Пусть v – скорость движения некоторого объекта, t – время и s – расстояние, пройденное за время t. Зависимость между этими величинами устанавливает формула пути: s = v · t (Для записи формулы пути используются строчные буквы s, v и t, чтобы не путать их с обозначением площади – S и объёма – V.) Формула пути означает, что расстояние равно скорости, умноженной на время движения. Из формулы пути по правилу нахождения неизвестного множителя следует, что: v = s : t      t = s : v • Скорость равна расстоянию, делённому на время движения. • Время движения равно расстоянию, делённому на скорость.', '<div class="recomended-block">\n<span class="title">Формула пути</span>\n<p>\nПусть  v – скорость движения некоторого объекта, t – время и s – расстояние, пройденное за время t. Зависимость между этими величинами устанавливает формула пути:<br/> \ns = v · t<br/>\n(Для записи формулы пути используются строчные буквы s, v и t, чтобы не путать их с обозначением площади – S и объёма – V.) <br/>\nФормула пути означает, что расстояние равно скорости,  умноженной на время движения. <br/>\nИз формулы пути по правилу нахождения неизвестного множителя следует, что:<br/>\nv = s  :  t      t  =  s  :  v<br/>\n• Скорость равна расстоянию, делённому на время движения.<br/>\n• Время движения равно расстоянию, делённому на скорость.\n\n</p>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-5/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1c99c5121f2fd10eaf0de71fd2c49d13c66bf06bc9902d34a2b7c3525a87aaef', NULL, NULL, 'проанализируй решение предыдущей задачи. установи, как найти расстояние s, пройденное объектом, если он двигался со скоростью v в течение времени t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 6, '4', 0, 'Найди неизвестные значения величин по формуле пути s = v · t:', '</p> \n<p class="text">Найди неизвестные значения величин по формуле пути s = v · t: </p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica6-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 6, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 6, номер 4, год 2022."/>\n</div>\n</div>', 'а) s = 5 · 9 = 45 (м/с) 48 = v · 6, v = 48 : 6 = 8 (км/ч) 21 = 7 · t, t = 21 : 7 = 3 (м/мин) б) 320 = v · 80, v = 320 : 80 = 4 (км/ч) 810 = 9 · t, t = 810 : 9 = 90 (мин) s = 60 · 50 = 3000 (м)', '<p>\nа) s = 5 · 9 = 45 (м/с)<br/>\n48 = v · 6, v = 48 : 6 = 8 (км/ч)<br/>\n21 = 7 · t, t = 21 : 7 = 3 (м/мин)<br/><br/>\nб) 320 = v · 80, v = 320 : 80 = 4 (км/ч)<br/>\n810 = 9 · t, t = 810 : 9 = 90 (мин)<br/>\ns = 60 · 50 = 3000 (м)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-6/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica6-nomer4.jpg', 'peterson/3/part3/page6/task4_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '393928ea37622759cc553539c3b6e10ab53d8c58571d30e78977ae0c03b9b450', NULL, '["найди"]'::jsonb, 'найди неизвестные значения величин по формуле пути s=v*t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 6, '5', 1, 'Реши задачи по формуле пути s = v · t: а) Всадник едет со скоростью 8 км/ч. Какое расстояние он проедет за 4 часа? б) Чему равна скорость почтового голубя, если за 2 ч он пролетает 120 км? в) Пчела летит со скоростью 6 м/с. За какое время она долетит до улья, если находится на расстоянии 360 м от него?', '</p> \n<p class="text">Реши задачи по формуле пути s = v · t:<br/>\nа) Всадник едет со скоростью 8 км/ч. Какое расстояние он проедет за 4 часа?<br/>\nб) Чему равна скорость почтового голубя, если за 2 ч он пролетает 120 км?<br/>\nв) Пчела летит со скоростью 6 м/с. За какое время она долетит до улья, если находится на расстоянии 360 м от него?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica6-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 6, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 6, номер 5, год 2022."/>\n</div>\n</div>', 'а) 8 · 4 = 32 (км) Ответ: 32 километра расстояние он проедет за 4 часа. б) 120 : 2 = 60 (км/ч) Ответ: 60 км/ч скорость почтового голубя, если за 2 ч он пролетает 120 км. в) 360 : 6 = 60 (с) Ответ: за 60 секунд она долетит до улья, если находится на расстоянии 360 м от него.', '<p>\nа) 8 · 4 = 32 (км)<br/>\n<b>Ответ:</b> 32 километра расстояние он проедет за 4 часа.<br/><br/>\nб) 120 : 2 = 60 (км/ч)<br/>\n<b>Ответ:</b> 60 км/ч скорость почтового голубя, если за 2 ч он пролетает 120 км.<br/><br/>\nв) 360 : 6 = 60 (с)<br/>\n<b>Ответ:</b> за 60 секунд она долетит до улья, если находится на расстоянии 360 м от него.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-6/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica6-nomer5.jpg', 'peterson/3/part3/page6/task5_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'dc5c3f2104cac2e6007e74ae01873457bf64fcd0b597ac29e663c6c5d10ba046', '2,4,6,8,120,360', '["реши"]'::jsonb, 'реши задачи по формуле пути s=v*t:а) всадник едет со скоростью 8 км/ч. какое расстояние он проедет за 4 часа? б) чему равна скорость почтового голубя, если за 2 ч он пролетает 120 км? в) пчела летит со скоростью 6 м/с. за какое время она долетит до улья, если находится на расстоянии 360 м от него');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 6, '6', 2, 'Между городом и деревней 250 км. Машина выехала из города в 10 часов утра и прибыла в деревню в 3 часа дня. С какой скоростью она ехала?', '</p> \n<p class="text">Между городом и деревней 250 км. Машина выехала из города в 10 часов утра и прибыла в деревню в 3 часа дня. С какой скоростью она ехала?</p>', '3 часа дня это 15 часов – 10 часов = 5 часов; 250 : 5 = 50 (км/ч). Ответ: 50 км/ч скоростью она ехала.', '<p>\n3 часа дня это 15 часов – 10 часов = 5 часов;<br/>\n250 : 5 = 50 (км/ч).<br/>\n<b>Ответ:</b> 50 км/ч скоростью она ехала.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-6/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5337dfe1672b7398b3e81a3ec7fa261444724c1791d713538195504183dff840', '3,10,250', NULL, 'между городом и деревней 250 км. машина выехала из города в 10 часов утра и прибыла в деревню в 3 часа дня. с какой скоростью она ехала');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 6, '7', 3, 'Аквариум имеет форму прямоугольного параллелепипеда. Длина аквариума – 50 см, ширина – 35 см, а высота – 40 см. Его боковые стенки стеклянные. Определи площадь поверхности стекла и объём аквариума.', '</p> \n<p class="text">Аквариум имеет форму прямоугольного параллелепипеда. Длина аквариума – 50 см, ширина – 35 см, а высота – 40 см. Его боковые стенки стеклянные. Определи площадь поверхности стекла и объём аквариума.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica6-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 6, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 6, номер 7, год 2022."/>\n</div>\n</div>', '50 · 35 · 40 = 70000 (см 3 ) объем аквариума. (50 · 40 + 35 · 40) · 2 = 6800 (см 2 ) площадь стеклянных стенок аквариума. Ответ: 68 дм 2 площадь стенок и 70000 см 3 объем аквариума.', '<p>\n50 · 35 · 40 = 70000 (см<sup>3</sup>) объем аквариума. <br/>\n(50 · 40 + 35 · 40) · 2 = 6800 (см<sup>2</sup>) площадь стеклянных стенок аквариума.<br/> \n<b>Ответ:</b> 68 дм<sup>2</sup> площадь стенок и 70000 см<sup>3</sup> объем аквариума.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-6/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica6-nomer7.jpg', 'peterson/3/part3/page6/task7_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c090917bf2794f703edb8988c7df78156835c3b981181c702701f2b54d218010', '35,40,50', '["площадь"]'::jsonb, 'аквариум имеет форму прямоугольного параллелепипеда. длина аквариума-50 см, ширина-35 см, а высота-40 см. его боковые стенки стеклянные. определи площадь поверхности стекла и объём аквариума');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 6, '8', 4, 'Реши уравнения с комментированием и сделай проверку: а) (25 - a) · 7 = 63 б) 400 : b - 32 = 48 в) 250 + 9 · c = 520', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (25 - a) · 7 = 63<br/>       б) 400 : b - 32 = 48<br/>       в) 250 + 9 · c = 520\n</p>', 'а) (25 - a) · 7 = 63 Чтобы найти вычитаемое разности произведения а надо вычесть из уменьшаемого разности произведения частное произведения и множителя а = 25 - 63 : 7 а = 16 Проверка: (25 - 16) · 7 = 63 б) 400 : b - 32 = 48 b = (48 + 32) · 400 b = 2000 Проверка: 400 : 2000 - 32 = 48 в) 250 + 9 · c = 520 Чтобы найти множитель слагаемого с надо отнять от суммы слагаемое и разделить на множитель слагаемого с = (520 - 250) : 9 с = 30 Проверка: 250 + 9 · 30 = 520', '<p>\nа) (25 - a) · 7 = 63<br/>   \nЧтобы найти вычитаемое разности произведения а надо вычесть из уменьшаемого разности произведения частное произведения и множителя<br/>\nа = 25 - 63 : 7<br/>\nа = 16<br/>\n<b>Проверка:</b> (25 - 16) · 7 = 63<br/><br/>    \nб) 400 : b - 32 = 48<br/>\nb = (48 + 32) · 400<br/>\nb = 2000<br/>\n<b>Проверка:</b> 400 : 2000 - 32 = 48 <br/><br/>      \nв) 250 + 9 · c = 520<br/>\nЧтобы найти множитель слагаемого с надо отнять от суммы слагаемое и разделить на множитель слагаемого<br/>\nс = (520 - 250) : 9<br/>\nс = 30<br/>\n<b>Проверка:</b> 250 + 9 · 30 = 520\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-6/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e8d545861732943f1717cb851df64681507db2f4a231347d093b793343d51f6e', '7,9,25,32,48,63,250,400,520', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (25-a)*7=63 б) 400:b-32=48 в) 250+9*c=520');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 6, '9', 5, 'Запиши множество делителей и множество кратных числа 11 *. * В заданиях учебника делители и кратные – натуральные числа.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 11 *.</p> \n\n<p class="description-text"> \n* В заданиях учебника делители и кратные – натуральные числа.\n</p>', '{1; 11}', '<p>\n{1; 11}\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-6/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8bc60dde6c078f828202a34ae7c78e9ad558552cbc2144bedf68c23e19bfc131', '11', NULL, 'запиши множество делителей и множество кратных числа 11*.*в заданиях учебника делители и кратные-натуральные числа');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 6, '10', 6, 'Составь программу действий и вычисли: а) (63200856 - 4916321) : 1 + 8006512 · (36 - 36) б) 1 · 7007503 - 29867 · (387915 : 387915)', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) (63200856 - 4916321) : 1 + 8006512 · (36 - 36)<br/>\nб) 1 · 7007503 - 29867 · (387915 : 387915)\n\n</p>', 'а) (63200856 - 4916321) : 1 + 8006512 · (36 - 36) = 58284535 63200856 - 4916321 = 58284535 а) В одном году 52 недели и 1 день - если год обычный и 52 недели и 2 дня - если год високосный. Чтобы посчитать количество недель в году нужно количество дней в году разделить на количество дней в неделе. В обычном году - 365 дней, в високосном - 366. Обычный год: 365 : 7 = 52 недели и 1 день Високосный год: 366 : 7 = 52 недели и 2 дня б) В году 365 дней и 53 вторника какой день недели был 1 января этого года. Этот год состоял из 52 недель и еще одного дня. В 52 неделях было 52 вторника, значит последний день года – вторник. Каждая из 52 предшествующих этому вторнику недель оканчивалась понедельником и, следовательно, начиналась с воскресенья. Следовательно, этот год начался с воскресенья. Воскресенье был 1 января этого года.', '<p>\nа) (63200856 - 4916321) : 1 + 8006512 · (36 - 36) = 58284535<br/>\n63200856 - 4916321 = 58284535\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="200" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica6-nomer10.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 6, номер 10, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 6, номер 10, год 2022."/>\n\n\n<p>\nа) В одном году 52 недели и 1 день - если год обычный и 52 недели и 2 дня - если год високосный. Чтобы посчитать количество недель в году нужно количество дней в году разделить на количество дней в неделе. В обычном году - 365 дней, в високосном - 366.<br/>\nОбычный год: 365 : 7 = 52 недели и 1 день<br/>\nВисокосный год: 366 : 7 = 52 недели и 2 дня<br/><br/>\nб) В году 365 дней и 53 вторника какой день недели был 1 января этого года. Этот год состоял из 52 недель и еще одного дня. В 52 неделях было 52 вторника, значит последний день года – вторник. Каждая из 52 предшествующих этому вторнику недель оканчивалась понедельником и, следовательно, начиналась с воскресенья. Следовательно, этот год начался с воскресенья. Воскресенье был 1 января этого года.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-6/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica6-nomer10.jpg', 'peterson/3/part3/page6/task10_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3f2373e41ec975a7d2c3c932afb860a283ea9d57d2f7f1f8cbb6ee3424cd57a2', '1,36,29867,387915,4916321,7007503,8006512,63200856', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) (63200856-4916321):1+8006512*(36-36) б) 1*7007503-29867*(387915:387915)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 7, '1', 0, 'Назови величины, характеризующие движение объектов. Объясни смысл предложений: а) Скорость воробья примерно 40 км/ч. б) Самая быстрая в мире птица сапсан способна развивать скорость до 200 км/ч. в) Африканский страус не может летать, зато разгоняется до 72 км/ч. г) Меч-рыба плывёт со скоростью 100 км/ч.', '</p> \n<p class="text">Назови величины, характеризующие движение объектов. Объясни смысл предложений: <br/>\nа) Скорость воробья примерно 40 км/ч.<br/>\nб) Самая быстрая в мире птица сапсан способна развивать скорость до 200 км/ч.<br/>\nв) Африканский страус не может летать, зато разгоняется до 72 км/ч.<br/>\nг) Меч-рыба плывёт со скоростью 100 км/ч.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 7, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 7, номер 1, год 2022."/>\n</div>\n</div>', 'а) Воробей примерно 40 км пролетает за 1 ч. б) Самая быстрая в мире птица сапсан способна преодолевать до 200 км за 1 ч. в) Африканский страус не может летать, зато преодолевает до 72 км за 1 ч. г) Меч - рыба проплывает 100 км за 1 ч.', '<p>\nа) Воробей примерно 40 км пролетает за 1 ч.<br/>\nб) Самая быстрая в мире птица сапсан способна преодолевать до 200 км за 1 ч.<br/>\nв) Африканский страус не может летать, зато преодолевает до 72 км за 1 ч.<br/>\nг) Меч - рыба проплывает 100 км за 1 ч.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-7/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer1.jpg', 'peterson/3/part3/page7/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '4726f073e7425c85f69fe2093a4b35a9cde16bc98b03908067ef80f370e89c45', '40,72,100,200', '["раз"]'::jsonb, 'назови величины, характеризующие движение объектов. объясни смысл предложений:а) скорость воробья примерно 40 км/ч. б) самая быстрая в мире птица сапсан способна развивать скорость до 200 км/ч. в) африканский страус не может летать, зато разгоняется до 72 км/ч. г) меч-рыба плывёт со скоростью 100 км/ч');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 7, '2', 1, 'Используя формулу пути s = v · t, найди неизвестные значения величин:', '</p> \n<p class="text">Используя формулу пути s = v · t, найди неизвестные значения величин:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 7, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 7, номер 2, год 2022."/>\n</div>\n</div>', '60 = v · 3, v = 60 : 3 = 2 (км/ч) s = 9 · 40 = 360 (м) 75 = 3 · t, t = 75 : 3 = 25 (с) 48 = 2 · t, t = 48 : 2 = 24 (мин) 540 = v · 18, v = 540 : 18 =30 (дм/с) s = 64 · 4 = 256 (км)', '<p>\n60 = v · 3, v = 60 : 3 = 2 (км/ч)	<br/>	\ns = 9 · 40 = 360 (м)	<br/>			\n75 = 3 · t, t = 75 : 3 = 25 (с)	<br/><br/>	\n\n48 = 2 · t, t = 48 : 2 = 24 (мин)<br/>\n540 = v · 18, v = 540 : 18 =30 (дм/с)<br/>\ns = 64 · 4 = 256 (км)\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-7/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer2.jpg', 'peterson/3/part3/page7/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c9d405093478e9cc584fa97f7fd13b6d2388b72e1daf5ed6f13ab982098779a9', NULL, '["найди"]'::jsonb, 'используя формулу пути s=v*t, найди неизвестные значения величин');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 7, '3', 2, 'По реке плывёт плот со скоростью v = 2 км/ч. Построй в тетради числовой луч и покажи на нём движение плота. Какое расстояние пройдёт плот за 1 ч, 3 ч, 5 ч, 7 ч, t ч? Составь и заполни в тетради таблицу. Напиши формулу, выражающую зависимость пройденного расстояния s от времени t.', '</p> \n<p class="text">По реке плывёт плот со скоростью v = 2 км/ч. Построй в тетради числовой луч и покажи на нём движение плота.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 7, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 7, номер 3, год 2022."/>\n</div>\n</div>\n\n<p class="text">Какое расстояние пройдёт плот за 1 ч, 3 ч, 5 ч, 7 ч, t ч? Составь и заполни в тетради таблицу. Напиши формулу, выражающую зависимость пройденного расстояния s от времени t.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer3-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 7, номер 3-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 7, номер 3-1, год 2022."/>\n</div>\n</div>', 's = 2 · 1 = 2 (км) s = 2 · 3 = 6 (км) s = 2 · 5 = 10 (км) s = 2 · 7 = 14 (км) s = 2 · t', '<p>\ns = 2 · 1 = 2 (км)<br/>\ns = 2 · 3 = 6 (км)<br/>\ns = 2 · 5 = 10 (км)<br/>\ns = 2 · 7 = 14 (км)\ns = 2 · t\n\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer3-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 7, номер 3-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 7, номер 3-2, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-7/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer3.jpg', 'peterson/3/part3/page7/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer3-1.jpg', 'peterson/3/part3/page7/task3_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer3-2.jpg', 'peterson/3/part3/page7/task3_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '997ecc106f1ffec6d3559cdd56a17becc18ce77df52917ad51b27abe39fe71b2', '1,2,3,5,7', '["заполни","числовой луч"]'::jsonb, 'по реке плывёт плот со скоростью v=2 км/ч. построй в тетради числовой луч и покажи на нём движение плота. какое расстояние пройдёт плот за 1 ч, 3 ч, 5 ч, 7 ч, t ч? составь и заполни в тетради таблицу. напиши формулу, выражающую зависимость пройденного расстояния s от времени t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 7, '4', 3, 'а) Космический корабль летит со скоростью 9 км/с. За какое время он пролетит 441 км? б) Сколько метров проплывёт окунь за 8 мин, если будет плыть со скоростью 80 м/мин? в) Подводная лодка проплыла 228 км за 6 ч. Чему равна её скорость? г) Улитка ползёт со скоростью 5 м/ч. За какое время она проползёт 35 м?', '</p> \n<p class="text">а) Космический корабль летит со скоростью 9 км/с. За какое время он пролетит 441 км? б) Сколько метров проплывёт окунь за 8 мин, если будет плыть со скоростью 80 м/мин? в) Подводная лодка проплыла 228 км за 6 ч. Чему равна её скорость? г) Улитка ползёт со скоростью 5 м/ч. За какое время она проползёт 35 м?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 7, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 7, номер 4, год 2022."/>\n</div>\n</div>', 'а) t = 441 : 9 = 49 (с) Ответ: за 49 секунд он пролетит 441 км. б) s = 80 · 8 = 640 (м) Ответ: 640 метров проплывёт окунь за 8 мин, если будет плыть со скоростью 80 м/мин. в) v = 228 : 6 = 38 (км/ч) Ответ: 38 км/ч скорость подводной лодки. г) t = 35 : 5 = 7 (ч) Ответ: за 7 часов она проползёт 35 м.', '<p>\nа) t = 441 : 9 = 49 (с)<br/>\n<b>Ответ:</b> за 49 секунд он пролетит 441 км.<br/><br/>\nб) s = 80 · 8 = 640 (м) <br/>\n<b>Ответ:</b> 640 метров проплывёт окунь за 8 мин, если будет плыть со скоростью 80 м/мин.<br/><br/>\nв) v = 228 : 6 = 38 (км/ч)<br/>\n<b>Ответ:</b> 38 км/ч скорость подводной лодки.<br/><br/>\nг) t = 35 : 5 = 7 (ч)<br/>\n<b>Ответ:</b> за 7 часов она проползёт 35 м.\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Космический корабль летит со скоростью 9 км/с. За какое время он пролетит 441 км?","solution":"t = 441 : 9 = 49 (с) Ответ: за 49 секунд он пролетит 441 км."},{"letter":"б","condition":"Сколько метров проплывёт окунь за 8 мин, если будет плыть со скоростью 80 м/мин?","solution":"s = 80 · 8 = 640 (м) Ответ: 640 метров проплывёт окунь за 8 мин, если будет плыть со скоростью 80 м/мин."},{"letter":"в","condition":"Подводная лодка проплыла 228 км за 6 ч. Чему равна её скорость?","solution":"v = 228 : 6 = 38 (км/ч) Ответ: 38 км/ч скорость подводной лодки."},{"letter":"г","condition":"Улитка ползёт со скоростью 5 м/ч. За какое время она проползёт 35 м?","solution":"t = 35 : 5 = 7 (ч) Ответ: за 7 часов она проползёт 35 м."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-7/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica7-nomer4.jpg', 'peterson/3/part3/page7/task4_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '562a4e66b9d5d2f10d8879e160a5acf9593c2363875a0e4516cd0956da36dd8a', '5,6,8,9,35,80,228,441', NULL, 'а) космический корабль летит со скоростью 9 км/с. за какое время он пролетит 441 км? б) сколько метров проплывёт окунь за 8 мин, если будет плыть со скоростью 80 м/мин? в) подводная лодка проплыла 228 км за 6 ч. чему равна её скорость? г) улитка ползёт со скоростью 5 м/ч. за какое время она проползёт 35 м');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 8, '5', 0, 'Составь программу действий и вычисли: а) 50 - (600 · 3) : (4 · 25) - 5 · (40 - 7 · 5) б) (80 · 8 + 420 : 7) : 100 + (140 : 20 + 38 : 19) · 3', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) 50 - (600 · 3) : (4 · 25) - 5 · (40 - 7 · 5)<br/>\nб) (80 · 8 + 420 : 7) : 100 + (140 : 20 + 38 : 19) · 3\n\n</p>', 'а) 50 - (600 · 3) : (4 · 25) - 5 · (40 - 7 · 5) = 7 600 · 3 = 1800 4 · 25 = 100 1800 : 100 = 18 7 · 5 = 35 40 - 35 = 5 5 · 5 = 25 50 - 18 = 32 32 - 25 = 7 б) (80 · 8 + 420 : 7) : 100 + (140 : 20 + 38 : 19) · 3 = 34 80 · 8 = 640 420 : 7 = 60 640 + 60 = 700 700 : 100 = 7 140 : 20 = 7 38 : 19 = 2 7 + 2 = 9 9 · 3 = 27 7 + 27 = 34', '<p>\nа) 50 - (600 · 3) : (4 · 25) - 5 · (40 - 7 · 5) = 7<br/>\n600 · 3 = 1800<br/>\n4 · 25 = 100<br/>\n1800 : 100 = 18<br/>\n7 · 5 = 35<br/>\n40 - 35 = 5<br/>\n5 · 5 = 25<br/>\n50 - 18 = 32<br/>\n32 -  25 = 7<br/><br/>\nб) (80 · 8 + 420 : 7) : 100 + (140 : 20 + 38 : 19) · 3 = 34<br/>\n80 · 8 = 640<br/>\n420 : 7 = 60<br/>\n640 + 60 = 700<br/>\n700 : 100 = 7<br/>\n140 : 20 = 7<br/>\n38 : 19 = 2<br/>\n7 + 2 = 9<br/>\n9 · 3 = 27<br/>\n7 + 27 = 34\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-8/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3603fb9387e5a03cb743d29489366e2a9e9d12e5010c964c625f20eae9be5350', '3,4,5,7,8,19,20,25,38,40', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) 50-(600*3):(4*25)-5*(40-7*5) б) (80*8+420:7):100+(140:20+38:19)*3');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 8, '6', 1, 'Определи по спидометру скорость движения каждой машины:', '</p> \n<p class="text">Определи по спидометру скорость движения каждой машины:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica8-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 8, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 8, номер 6, год 2022."/>\n</div>\n</div>', 'Легковой автомобиль 100 км/ч Автобус 90 км/ч Грузовик 50 км/ч', '<p>\nЛегковой автомобиль 100 км/ч<br/>\nАвтобус 90 км/ч<br/>\nГрузовик 50 км/ч\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-8/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica8-nomer6.jpg', 'peterson/3/part3/page8/task6_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b3f229f219e7c189e6f89c01bf4b40589c6ab660a4ad573b6013e74b09f0e2e4', NULL, NULL, 'определи по спидометру скорость движения каждой машины');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 8, '7', 2, 'Придумай задачу, в которой надо найти скорость по известному расстоянию и времени, и реши её.', '</p> \n<p class="text">Придумай задачу, в которой надо найти скорость по известному расстоянию и времени, и реши её.</p>', 'Грузовик проехал 50 км за 1 час. С какой скоростью двигался грузовик? 50 : 1 = 50 (км/ч) Ответ: 50 км/ч двигался грузовик.', '<p>\nГрузовик проехал 50 км за 1 час. С какой скоростью двигался грузовик?<br/>\n50 : 1 = 50 (км/ч)<br/>\n<b>Ответ:</b> 50 км/ч двигался грузовик.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-8/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '4ef504a8d1692e16c88fe03a4b21aba3cb04c2caf482a1aaf1b4d83052d5ec00', NULL, '["реши"]'::jsonb, 'придумай задачу, в которой надо найти скорость по известному расстоянию и времени, и реши её');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 8, '8', 3, 'Сравни: 5 ч 6 мин □ 56 мин      9 мин 20 с □ 560 с      1 сут. 15 ч □ 115 ч 108 мин □ 1 ч 8 мин      734 с □ 7 мин 34 с      206 ч □ 2 сут. 6 ч', '</p> \n<p class="text">Сравни: </p> \n\n<p class="description-text"> \n5 ч 6 мин □ 56 мин      9 мин 20 с □ 560 с      1 сут. 15 ч □ 115 ч<br/>\n108 мин □ 1 ч 8 мин      734 с □ 7 мин 34 с      206 ч □ 2 сут.  6 ч\n\n</p>', '5 ч 6 мин ˃ 56 мин      9 мин 20 с = 560 с      1 сут. 15 ч ˂ 115 ч 108 мин ˃ 1 ч 8 мин      734 с ˃ 7 мин 34 с      206 ч ˃ 2 сут. 6 ч', '<p>\n5 ч 6 мин ˃ 56 мин      9 мин 20 с = 560 с      1 сут. 15 ч ˂ 115 ч<br/>\n108 мин ˃ 1 ч 8 мин      734 с ˃ 7 мин 34 с      206 ч ˃ 2 сут.  6 ч\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-8/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '430e4e44e11ebf88f9cee8f66021ddf19ff16a2a9cd04276d65c5802fa9c739b', '1,2,5,6,7,8,9,15,20,34', '["сравни"]'::jsonb, 'сравни:5 ч 6 мин □ 56 мин      9 мин 20 с □ 560 с      1 сут. 15 ч □ 115 ч 108 мин □ 1 ч 8 мин      734 с □ 7 мин 34 с      206 ч □ 2 сут. 6 ч');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 8, '9', 4, 'Реши уравнения с комментированием и сделай проверку: а) (780 - m · 60) : 6 = 70 б) 640 : (x · 9 + 8) = 8', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (780 - m · 60) : 6 = 70<br/>\nб) 640 : (x · 9 + 8) = 8\n</p>', 'а) (780 - m · 60) : 6 = 70 m = (780 - 70 · 6) : 60 m = (780 - 420) : 60 m = 360 : 60 m = 6 Проверка: (780 - 6 · 60) : 6 = 70 б) 640 : (x · 9 + 8) = 8 x = (640 : 8 – 8) : 9 х = (80 - 8) : 9 х = 72 : 9 х = 8 Проверка: 640 : (8 · 9 + 8) = 8', '<p>\nа) (780 - m · 60) : 6 = 70 <br/> \nm = (780 - 70 · 6) : 60<br/>\nm = (780 - 420) : 60<br/>\nm = 360 : 60<br/>\nm = 6    <br/>\nПроверка: (780 - 6 · 60) : 6 = 70 <br/><br/>   \nб) 640 : (x · 9 + 8) = 8<br/>\nx = (640 : 8 – 8) : 9<br/>\nх = (80 - 8) : 9 <br/>\nх = 72 : 9<br/>\nх = 8<br/>\nПроверка: 640 : (8 · 9 + 8) = 8\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-8/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a47137e79395ac396e53b2ff32297f2dc52693eb66f8f492a5f897689ff35b2c', '6,8,9,60,70,640,780', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (780-m*60):6=70 б) 640:(x*9+8)=8');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 8, '10', 5, 'Запиши множество делителей и множество кратных числа 12.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 12.</p>', 'Делители: 12, 24, 36, 48, 60, 72, 84, кратные: 1, 2, 3,4, 6, 12. Делители: 1,12,3,4,2,6. Кратные: 12, 24,36.', '<p>\nДелители: 12, 24, 36, 48, 60, 72, 84, кратные: 1, 2, 3,4, 6, 12. Делители: 1,12,3,4,2,6. Кратные: 12, 24,36.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-8/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '69392f0b1357aea3d5448bef7ab15460dd40b9584ff2d0111adb9550058ed383', '12', NULL, 'запиши множество делителей и множество кратных числа 12');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 8, '11', 6, 'По диаграмме Эйлера – Венна определи, из каких элементов состоят множества A и B. Запиши эти множества с помощью фигурных скобок. Найди их пересечение и объединение.', '</p> \n<p class="text">По диаграмме Эйлера – Венна определи, из каких элементов состоят множества A и B. Запиши эти множества с помощью фигурных скобок. Найди их пересечение и объединение. </p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="250" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica8-nomer11.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 8, номер 11, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 8, номер 11, год 2022."/>\n</div>\n</div>', 'A = {5; d; ; }, B = {d; ; n; 3; 7}, А В = {d; } и обведи А В = {3; 5; 7; n; d; ; } 2x + 1 + 3 = 20 х = (20 - 4) : 2 х = 16 : 2 х = 8 Ответ: 8 было гусей.', '<p>\nA = {5; d;  ;  }, B = {d;  ; n; 3; 7}, А  В = {d;  } и обведи А   В = {3; 5; 7; n; d;  ;  }\n</p>\n\n\n<p>\n2x + 1 + 3 = 20<br/>\nх = (20 - 4) : 2<br/>\nх = 16 : 2 <br/>\nх = 8<br/>\n<b>Ответ:</b> 8 было гусей.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-8/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica8-nomer11.jpg', 'peterson/3/part3/page8/task11_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8a9b4b21af6c6206fd7a36eea13dc74a74b94949ee02758d8389be8ca215eece', NULL, '["найди"]'::jsonb, 'по диаграмме эйлера-венна определи, из каких элементов состоят множества a и b. запиши эти множества с помощью фигурных скобок. найди их пересечение и объединение');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 9, '1', 0, 'а) Из Москвы в Селижарово выехал автобус со скоростью 60 км/ч. Построй числовой луч и покажи на нём движение автобуса. Какое расстояние прошёл автобус за 1 ч, 2 ч, 3 ч, 4 ч, 5 ч, 6 ч, t ч? Через какое время он приедет в Селижарово? Составь и заполни таблицу. Напиши формулу зависимости пройденного расстояния s от времени t. б) Проанализируй решение предыдущей задачи. Объясни, как можно построить формулу зависимости одной величины от другой.', '</p> \n<p class="text">а) Из Москвы в Селижарово выехал автобус со скоростью 60 км/ч. Построй числовой луч и покажи на нём движение автобуса.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica9-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 9, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 9, номер 1, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">Какое расстояние прошёл автобус за 1 ч, 2 ч, 3 ч, 4 ч, 5 ч, 6 ч, t ч? Через какое время он приедет в Селижарово? Составь и заполни таблицу. Напиши формулу зависимости пройденного расстояния s от времени t.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica9-nomer1-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 9, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 9, номер 1, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Проанализируй решение предыдущей задачи. Объясни, как можно построить формулу зависимости одной величины от другой.</p>', 'а)', '<p>\nа) \n</p>\n\n<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica9-nomer1-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 9, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 9, номер 1, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Из Москвы в Селижарово выехал автобус со скоростью 60 км/ч. Построй числовой луч и покажи на нём движение автобуса. Какое расстояние прошёл автобус за 1 ч, 2 ч, 3 ч, 4 ч, 5 ч, 6 ч, t ч? Через какое время он приедет в Селижарово? Составь и заполни таблицу. Напиши формулу зависимости пройденного расстояния s от времени t.","solution":""},{"letter":"б","condition":"Проанализируй решение предыдущей задачи. Объясни, как можно построить формулу зависимости одной величины от другой.","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-9/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica9-nomer1.jpg', 'peterson/3/part3/page9/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica9-nomer1-1.jpg', 'peterson/3/part3/page9/task1_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica9-nomer1-2.jpg', 'peterson/3/part3/page9/task1_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5eafa8922c95feb525c6ef6af804c30a44c4daed214d30962d4a7af87a85832e', '1,2,3,4,5,6,60', '["заполни","числовой луч"]'::jsonb, 'а) из москвы в селижарово выехал автобус со скоростью 60 км/ч. построй числовой луч и покажи на нём движение автобуса. какое расстояние прошёл автобус за 1 ч, 2 ч, 3 ч, 4 ч, 5 ч, 6 ч, t ч? через какое время он приедет в селижарово? составь и заполни таблицу. напиши формулу зависимости пройденного расстояния s от времени t. б) проанализируй решение предыдущей задачи. объясни, как можно построить формулу зависимости одной величины от другой');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 9, '2', 1, 'Расстояние между двумя городами 180 км. С какой скоростью надо ехать, чтобы преодолеть это расстояние за 1 ч, 2 ч, 3 ч, 4 ч, t ч? Построй формулу зависимости скорости движения v от времени t.', '</p> \n<p class="text">Расстояние между двумя городами 180 км. С какой скоростью надо ехать, чтобы преодолеть это расстояние за 1 ч, 2 ч, 3 ч, 4 ч, t ч? Построй формулу зависимости скорости движения v от времени t.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="120" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica9-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 9, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 9, номер 2, год 2022."/>\n</div>\n</div>', '180 : 1 = 180 (км/ч) 180 : 2 = 90 (км/ч) 180 : 3 = 60 (км/ч) 180 : 4 = 45 (км/ч) v = s : t', '<p>\n180 : 1 = 180 (км/ч)<br/>\n180 : 2 = 90 (км/ч)<br/>\n180 : 3 = 60 (км/ч)<br/>\n180 : 4 = 45 (км/ч)<br/>\nv = s : t\n</p>', 'Алгоритм построения формул зависимостей между величинами. 1. Составить таблицу соответствующих значений величин (если нужно, использовать схему). 2. Понаблюдать, как изменяются значения одной величины при изменении другой. 3. Выявить закономерность и записать её в виде формулы. Пример:', '<div class="recomended-block">\n<span class="title">Алгоритм построения формул зависимостей между величинами.</span>\n<p>\n1. Составить таблицу соответствующих значений величин (если нужно, использовать схему).<br/>\n2. Понаблюдать, как изменяются значения одной величины при изменении другой.<br/>\n3. Выявить закономерность и записать её в виде формулы. <br/><br/>\nПример:\n\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica9-spravka.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 9, справка, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 9, справка, год 2022."/>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-9/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica9-nomer2.jpg', 'peterson/3/part3/page9/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '87530333ee7e103d00f3e0c161ca3e6ca790c6605c00de00897ce2d4465b299c', '1,2,3,4,180', NULL, 'расстояние между двумя городами 180 км. с какой скоростью надо ехать, чтобы преодолеть это расстояние за 1 ч, 2 ч, 3 ч, 4 ч, t ч? построй формулу зависимости скорости движения v от времени t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 10, '3', 0, 'Какое расстояние пройдёт поезд за 5 ч, если движется со скоростью 70 км/ч, 82 км/ч, 90 км/ч, 100 км/ч, v км/ч? Составь формулу зависимости расстояния s от скорости v.', '</p> \n<p class="text">Какое расстояние пройдёт поезд за 5 ч, если движется со скоростью 70 км/ч, 82 км/ч, 90 км/ч, 100 км/ч, v км/ч? Составь формулу зависимости расстояния s от скорости v.</p>', 's = v · t 70 · 5 = 350 (км) 82 · 5 = 410 (км) 90 · 5 = 450 (км) 100 · 5 = 500 (км)', '<p>\ns = v · t<br/>\n70 · 5 = 350 (км)<br/>\n82 · 5 = 410 (км)<br/>\n90 · 5 = 450 (км)<br/>\n100 · 5 = 500 (км)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-10/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '73a216458826a29ce629d9305954e74ce0431bb4e7c3f7c71e6ab4114a7f44c6', '5,70,82,90,100', NULL, 'какое расстояние пройдёт поезд за 5 ч, если движется со скоростью 70 км/ч, 82 км/ч, 90 км/ч, 100 км/ч, v км/ч? составь формулу зависимости расстояния s от скорости v');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 10, '4', 1, 'Сколько времени потребуется велосипедисту, чтобы проехать 60 км, если скорость его движения 10 км/ч, 12 км/ч, 15 км/ч, 20 км/ч, v км/ч? Составь формулу зависимости времени t от скорости v.', '</p> \n<p class="text">Сколько времени потребуется велосипедисту, чтобы проехать 60 км, если скорость его движения 10 км/ч, 12 км/ч, 15 км/ч, 20 км/ч, v км/ч? Составь формулу зависимости времени t от скорости v. </p> \n\n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 10, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 10, номер 4, год 2022."/>\n</div>\n</div>', '60 : 10 = 6 (ч) 60 : 12 = 5 (ч) 60 : 15 = 4 (ч) 60 : 20 = 3 (ч) t = s : v', '<p>\n60 : 10 = 6 (ч)<br/>\n60 : 12 = 5 (ч)<br/>\n60 : 15 = 4 (ч)<br/>\n60 : 20 = 3 (ч)<br/>\nt = s : v\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-10/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer4.jpg', 'peterson/3/part3/page10/task4_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a317bd08523abe08a6d08a8b0bc9544771ab1443055ab123c11a021fab3b7fb7', '10,12,15,20,60', NULL, 'сколько времени потребуется велосипедисту, чтобы проехать 60 км, если скорость его движения 10 км/ч, 12 км/ч, 15 км/ч, 20 км/ч, v км/ч? составь формулу зависимости времени t от скорости v');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 10, '5', 2, 'а) Расстояние между двумя пристанями 160 км. Может ли катер пройти это расстояние за 9 ч, если будет идти со скоростью 18 км/ч? б) От Сашиного дома до школы 1 км. Успеет ли он прийти в школу за 15 мин, если будет идти со скоростью 80 м/мин?', '</p> \n<p class="text">а) Расстояние между двумя пристанями 160 км. Может ли катер пройти это расстояние за 9 ч, если будет идти со скоростью 18 км/ч?<br/>\nб) От Сашиного дома до школы 1 км. Успеет ли он прийти в школу за 15 мин, если будет идти со скоростью 80 м/мин?\n</p>', 'а) 9 · 18 = 162 (км)', '<p>\nа) 9 · 18 = 162 (км)\n</p>\n\n<div class="img-wrapper-460">\n<img width="80" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 10, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 10, номер 5, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Расстояние между двумя пристанями 160 км. Может ли катер пройти это расстояние за 9 ч, если будет идти со скоростью 18 км/ч?","solution":"9 · 18 = 162 (км)"},{"letter":"б","condition":"От Сашиного дома до школы 1 км. Успеет ли он прийти в школу за 15 мин, если будет идти со скоростью 80 м/мин?","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-10/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer5.jpg', 'peterson/3/part3/page10/task5_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c195f9dc305f30e64ec8d3401b515840c1e117fdf786c7193eb2642dc2ef6e7d', '1,9,15,18,80,160', NULL, 'а) расстояние между двумя пристанями 160 км. может ли катер пройти это расстояние за 9 ч, если будет идти со скоростью 18 км/ч? б) от сашиного дома до школы 1 км. успеет ли он прийти в школу за 15 мин, если будет идти со скоростью 80 м/мин');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 10, '6', 3, 'Реши уравнения с комментированием и сделай проверку: а) 14 - 360 : x = 8 б) (450 : y + 50) : 70 = 2 в) (3 · z + 160) : 7 = 40', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) 14 - 360 : x = 8<br/>       \nб) (450 : y + 50) : 70 = 2<br/>      \nв) (3 · z + 160) : 7 = 40\n</p>', 'а) 14 - 360 : x = 8 х = 360 : (14 – 8) х = 360 : 6 х = 60 Проверка: 14 - 360 : 60 = 8 б) (450 : y + 50) : 70 = 2 у = 450 : (2 · 70 - 50) у = 450 : 90 у = 5 Проверка: (450 : 5 + 50) : 70 = 2 в) (3 · z + 160) : 7 = 40 z = (40 · 7 - 160) : 3 z = 120 : 3 z = 40 Проверка: (3 · 40 + 160) : 7 = 40', '<p>\nа) 14 - 360 : x = 8 <br/>   \nх = 360 : (14 – 8) <br/>\nх = 360 : 6<br/>\nх = 60<br/>\n<b>Проверка:</b> 14 - 360 : 60 = 8 <br/><br/> \nб) (450 : y + 50) : 70 = 2<br/>\nу = 450 : (2 · 70 - 50)<br/>\nу = 450 : 90<br/>\nу = 5<br/>\n<b>Проверка:</b> (450 : 5 + 50) : 70 = 2 <br/><br/>    \nв) (3 · z + 160) : 7 = 40<br/>\nz = (40 · 7 - 160) : 3<br/>\nz = 120 : 3<br/>\nz = 40<br/>\n<b>Проверка:</b> (3 · 40 + 160) : 7 = 40\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-10/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0dc95ff1e51569555aadec89e52eef55b3899763313a398f139697ec68134eef', '2,3,7,8,14,40,50,70,160,360', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) 14-360:x=8 б) (450:y+50):70=2 в) (3*z+160):7=40');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 10, '7', 4, 'Запиши множество делителей и множество кратных числа 13.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 13.</p>', 'Число 1 и само число 13. кратных числу 13. те чисел. делящихся нацело на 13. можно найти бесконечное множество 26, 39, 52, 520, 104, 130, 2600.', '<p>\nЧисло 1 и само число 13. кратных числу 13. те чисел. делящихся нацело на 13. можно найти бесконечное множество 26, 39, 52, 520, 104, 130, 2600.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-10/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b948c2204e80aefb0c100d6cc208c26c15745247cb09b698621f5a9e38adffbc', '13', NULL, 'запиши множество делителей и множество кратных числа 13');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 10, '8', 5, 'Найди пропущенные цифры. Проверь с помощью калькулятора. Сделай проверку, выполнив обратные действия.', '</p> \n<p class="text">Найди пропущенные цифры. Проверь с помощью калькулятора.<br/>\nСделай проверку, выполнив обратные действия.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 10, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 10, номер 8, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-460">\n<img width="100" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer8-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 10, номер 8-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 10, номер 8-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-10/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer8.jpg', 'peterson/3/part3/page10/task8_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer8-1.jpg', 'peterson/3/part3/page10/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b6ae258854b61f2121890ed9b31fe3f41ba25ef672f66c4201f9f6c80d725eaa', NULL, '["найди"]'::jsonb, 'найди пропущенные цифры. проверь с помощью калькулятора. сделай проверку, выполнив обратные действия');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 10, '9', 6, 'Кате надо было разделить число 48236 на 8. У неё получилось частное 629 и остаток 2. Проверь её вычисления с помощью формулы деления с остатком. Найди ошибку и выполни деление правильно.', '</p> \n<p class="text">Кате надо было разделить число 48236 на 8. У неё получилось частное 629 и остаток 2. Проверь её вычисления с помощью формулы деления с остатком. Найди ошибку и выполни деление правильно.</p>', '', '<div class="img-wrapper-460">\n<img width="90" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 10, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 10, номер 9, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-10/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer9.jpg', 'peterson/3/part3/page10/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '52381c0f1933cd270e3c1c5b434987ad4e12a14969388df3ea091813bba868de', '2,8,629,48236', '["раздели","найди","частное","раз","остаток"]'::jsonb, 'кате надо было разделить число 48236 на 8. у неё получилось частное 629 и остаток 2. проверь её вычисления с помощью формулы деления с остатком. найди ошибку и выполни деление правильно');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 10, '10', 7, 'Найди частное и остаток при делении: а) числа 14 на число 5; б) числа 6 на число 3; в) числа 2 на число 3. Обоснуй свой ответ, пользуясь формулой деления с остатком.', '</p> \n<p class="text">Найди частное и остаток при делении: а) числа 14 на число 5; б) числа 6 на число 3; в) числа 2 на число 3. Обоснуй свой ответ, пользуясь формулой деления с остатком.</p>', 'а) 14 : 5 = 2 + 4 – делится на 2 с остатком 4; б) 6 : 3 = 2 + 0 – делится ровно на 2 без остатка; в) 2 : 3 = 0 + 2 – нет целых в остатке 2.', '<p>\nа) 14 : 5 = 2 + 4 – делится на 2 с остатком 4;<br/> \nб) 6 : 3 = 2 + 0 – делится ровно на 2 без остатка;<br/> \nв) 2 : 3 = 0 + 2 – нет целых в остатке 2.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-10/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '358e012c78ee7bd1c614615ed3d60f016611c36a6ebdd278b6c4ff5055ff5428', '2,3,5,6,14', '["найди","частное","остаток"]'::jsonb, 'найди частное и остаток при делении:а) числа 14 на число 5; б) числа 6 на число 3; в) числа 2 на число 3. обоснуй свой ответ, пользуясь формулой деления с остатком');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 10, '11', 8, 'Выполни действия и сделай проверку: а) 483567823 + 998430           в) 37090 · 6000 б) 2666990000 - 89607787      г) 210040000 : 500', '</p> \n<p class="text">Выполни действия и сделай проверку:</p> \n\n<p class="description-text"> \nа) 483567823 + 998430           в) 37090 · 6000<br/> \nб) 2666990000 - 89607787      г) 210040000 : 500\n\n</p>', 'а) 483567823 + 998430 = 484566253 а) 1 сутки дольше длится, чем 1000 минут; б) 1000 часов дольше длится, чем 1 месяц; в) 1 год = 12 · 30 · 24 · 60 · 60 = 31104000 дольше длится, чем 1000000 секунд.', '<p>\nа) 483567823 + 998430 = 484566253\n</p>\n\n<div class="img-wrapper-460">\n<img width="260" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer11.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 10, номер 11, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 10, номер 11, год 2022."/>\n\n\n<p>\nа) 1 сутки дольше длится, чем 1000 минут; <br/>\nб) 1000 часов дольше длится, чем 1 месяц; <br/>\nв) 1 год = 12 · 30 · 24 · 60 · 60 = 31104000 дольше длится, чем 1000000 секунд.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-10/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica10-nomer11.jpg', 'peterson/3/part3/page10/task11_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c3b3dc465d7f05f12394ccc0fff0200b09714b704c4a966ee720a085f69dc40d', '500,6000,37090,998430,89607787,210040000,483567823,2666990000', NULL, 'выполни действия и сделай проверку:а) 483567823+998430           в) 37090*6000 б) 2666990000-89607787      г) 210040000:500');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 11, '1', 0, 'а) На числовом луче показано движение велосипедиста. Объясни, откуда он выехал. Куда и с какой скоростью он едет? В какой точке числового луча он был через 1 ч, 2 ч, 3 ч, 4 ч? Сколько времени он затратил на весь путь? б) Пусть s – путь, который проехал велосипедист, d – его расстояние до Ромашково и D – расстояние до Горок. Как изменяются d и D в зависимости от времени t – уменьшаются или увеличиваются? в) Заполни таблицу. Запиши формулу зависимости каждой из величин s, d, D от времени движения t.', '</p> \n<p class="text">а) На числовом луче показано движение велосипедиста. Объясни, откуда он выехал. Куда и с какой скоростью он едет? В какой точке числового луча он был через 1 ч, 2 ч, 3 ч, 4 ч? Сколько времени он затратил на весь путь?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 11, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 11, номер 1, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Пусть s – путь, который проехал велосипедист, d – его расстояние до Ромашково и D – расстояние до Горок. Как изменяются d и D в зависимости от времени t – уменьшаются или увеличиваются?<br/><br/>\nв) Заполни таблицу. Запиши формулу зависимости каждой из величин s,  d,  D от времени движения t.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer1-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 11, номер 1-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 11, номер 1-1, год 2022."/>\n</div>\n</div>', 'а) велосипедист выехал из Петушки. Поехал в Ромашково со скоростью 15 км/ч. Он был через 1 ч - 60, 2 ч - 45, 3 ч - 30, 4 ч – 15. На весть путь он затратил 5 ч. б) d уменьшается в зависимости от времени t, D в зависимости от времени t – увеличиваются. в) s = 15 · t, d = 75 - 15 · t, D = 75 + 15 · t.', '<p>\nа) велосипедист выехал из Петушки. Поехал в Ромашково со скоростью 15 км/ч. Он был через 1 ч - 60, 2 ч - 45, 3 ч - 30, 4 ч – 15. На весть путь он затратил 5 ч.<br/><br/>\nб) d уменьшается в зависимости от времени t, D в зависимости от времени t – увеличиваются. <br/><br/>\nв) s = 15 · t, d = 75 - 15 · t, D = 75 + 15 · t.\n\n</p>\n\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer1-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 11, номер 1-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 11, номер 1-2, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"На числовом луче показано движение велосипедиста. Объясни, откуда он выехал. Куда и с какой скоростью он едет? В какой точке числового луча он был через 1 ч, 2 ч, 3 ч, 4 ч? Сколько времени он затратил на весь путь?","solution":"велосипедист выехал из Петушки. Поехал в Ромашково со скоростью 15 км/ч. Он был через 1 ч - 60, 2 ч - 45, 3 ч - 30, 4 ч – 15. На весть путь он затратил 5 ч."},{"letter":"б","condition":"Пусть s – путь, который проехал велосипедист, d – его расстояние до Ромашково и D – расстояние до Горок. Как изменяются d и D в зависимости от времени t – уменьшаются или увеличиваются?","solution":"d уменьшается в зависимости от времени t, D в зависимости от времени t – увеличиваются."},{"letter":"в","condition":"Заполни таблицу. Запиши формулу зависимости каждой из величин s, d, D от времени движения t.","solution":"s = 15 · t, d = 75 - 15 · t, D = 75 + 15 · t."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-11/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer1.jpg', 'peterson/3/part3/page11/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer1-1.jpg', 'peterson/3/part3/page11/task1_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer1-2.jpg', 'peterson/3/part3/page11/task1_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'eadf69cba42ce6bc9b549f089d33da1efd97d24995bd0b630df4fb42afb3f7c1', '1,2,3,4', '["заполни"]'::jsonb, 'а) на числовом луче показано движение велосипедиста. объясни, откуда он выехал. куда и с какой скоростью он едет? в какой точке числового луча он был через 1 ч, 2 ч, 3 ч, 4 ч? сколько времени он затратил на весь путь? б) пусть s-путь, который проехал велосипедист, d-его расстояние до ромашково и d-расстояние до горок. как изменяются d и d в зависимости от времени t-уменьшаются или увеличиваются? в) заполни таблицу. запиши формулу зависимости каждой из величин s, d, d от времени движения t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 11, '2', 1, 'а) Определи по рисунку, откуда вышел турист, куда и с какой скоростью он идёт. Построй в тетради числовой луч и покажи на нём движение туриста. б) Пусть s км – путь, пройденный туристом, d км – расстояние между туристом и Москвой, D км – расстояние до Икши. Заполни таблицу. Запиши формулу зависимости каждой из величин s, d, D от времени движения t.', '</p> \n<p class="text">а) Определи по рисунку, откуда вышел турист, куда и с какой скоростью он идёт. Построй в тетради числовой луч и покажи на нём движение туриста.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 11, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 11, номер 2, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Пусть s км – путь, пройденный туристом, d км – расстояние между туристом и Москвой, D км – расстояние до Икши. Заполни таблицу. Запиши формулу зависимости каждой из величин s,  d,  D от времени движения  t.</p> \n\n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer2-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 11, номер 2-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 11, номер 2-1, год 2022."/>\n</div>\n</div>', 'a) турист вышел с Турбазы и пошёл к Икша со скоростью 3 км/ч. б) s = 3 · t, d = 12 + 3 · t, D = 18 - 3 · t.', '<p>\na) турист вышел с Турбазы и пошёл к Икша со скоростью 3 км/ч.<br/><br/> \nб) s = 3 · t, d = 12 + 3 · t, D = 18 - 3 · t.\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer2-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 11, номер 2-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 11, номер 2-2, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Определи по рисунку, откуда вышел турист, куда и с какой скоростью он идёт. Построй в тетради числовой луч и покажи на нём движение туриста.","solution":""},{"letter":"б","condition":"Пусть s км – путь, пройденный туристом, d км – расстояние между туристом и Москвой, D км – расстояние до Икши. Заполни таблицу. Запиши формулу зависимости каждой из величин s, d, D от времени движения t.","solution":"s = 3 · t, d = 12 + 3 · t, D = 18 - 3 · t."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-11/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer2.jpg', 'peterson/3/part3/page11/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer2-1.jpg', 'peterson/3/part3/page11/task2_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica11-nomer2-2.jpg', 'peterson/3/part3/page11/task2_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '725fef3574ea037a1f664a274a2a83f106f0be2f5fe9bbbc9cf27ef7efec58bc', NULL, '["заполни","числовой луч"]'::jsonb, 'а) определи по рисунку, откуда вышел турист, куда и с какой скоростью он идёт. построй в тетради числовой луч и покажи на нём движение туриста. б) пусть s км-путь, пройденный туристом, d км-расстояние между туристом и москвой, d км-расстояние до икши. заполни таблицу. запиши формулу зависимости каждой из величин s, d, d от времени движения t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 12, '3', 0, 'Расстояние от деревни до станции 40 км. Всадник едет из деревни на станцию со скоростью 14 км/ч. Успеет ли он доскакать до станции за 3 часа?', '</p> \n<p class="text">Расстояние от деревни до станции 40 км. Всадник едет из деревни на станцию со скоростью 14 км/ч. Успеет ли он доскакать до станции за 3 часа?</p>', '14 · 3 = 42 (км) Ответ: не успеет он доскакать до станции за 3 часа.', '<p>\n14 · 3 = 42 (км)<br/>\n<b>Ответ:</b> не успеет он доскакать до станции за 3 часа.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-12/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'fbef280b8a92f13b81c4d4c314d06d32c4122d2f06642b8bcec76663ca59b339', '3,14,40', NULL, 'расстояние от деревни до станции 40 км. всадник едет из деревни на станцию со скоростью 14 км/ч. успеет ли он доскакать до станции за 3 часа');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 12, '4', 1, 'Туристы решили пройти за день 30 км. Они уже прошли 3 ч со скоростью 6 км/ч. Какое расстояние им осталось пройти? За какое время они пройдут это расстояние, двигаясь с прежней скоростью?', '</p> \n<p class="text">Туристы решили пройти за день 30 км. Они уже прошли 3 ч со скоростью  6  км/ч.  Какое расстояние им осталось пройти? За какое время они пройдут это расстояние, двигаясь с прежней скоростью?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 12, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 12, номер 4, год 2022."/>\n</div>\n</div>', '30 - 6 · 3 = 30 - 18 = 12 (км) 12 : 6 = 2 (ч) Ответ: 12 километров им осталось пройти. За 2 часа они пройдут это расстояние, двигаясь с прежней скоростью.', '<p>\n30 - 6 · 3 = 30 - 18 = 12 (км)<br/>\n12 : 6 = 2 (ч)<br/>\n<b>Ответ:</b> 12 километров им осталось пройти. За 2 часа они пройдут это расстояние, двигаясь с прежней скоростью.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-12/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer4.jpg', 'peterson/3/part3/page12/task4_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'cf6a1ef68814949e47b5f889ca4f2c77d1fef2eef9f5ad54b8876713930eac2d', '3,6,30', '["реши"]'::jsonb, 'туристы решили пройти за день 30 км. они уже прошли 3 ч со скоростью 6 км/ч. какое расстояние им осталось пройти? за какое время они пройдут это расстояние, двигаясь с прежней скоростью');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 12, '5', 2, 'БЛИЦтурнир* а) Маша прошла n км. Чему равна её скорость, если она затратила на путь k часов? б) Лена шла a ч со скоростью b км/ч. Какое расстояние она прошла за это время? в) Витя пробежал x метров за 5 мин, а Саша – за 6 мин. У кого из них скорость больше и на сколько?', '</p> \n<p class="text"><b>БЛИЦтурнир*</b><br/>\nа) Маша прошла n км. Чему равна её скорость, если она затратила на путь k часов?<br/>\nб) Лена шла a ч со скоростью b км/ч. Какое расстояние она прошла за это время?<br/>\nв) Витя пробежал x метров за 5 мин, а Саша – за 6 мин. У кого из них скорость больше и на сколько?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 12, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 12, номер 5, год 2022."/>\n</div>\n</div>', 'а) n : k (км/ч) Ответ: n : k км/ч её скорость, если она затратила на путь k часов. б) a · b (км) Ответ: a · b километров она прошла за это время. в) Витя - x : 5 (м/мин), а Саша – х : 6 (м/мин), x : 6 - х : 5 = х : (6 - 5) = х : 1 = х (м/мин) Ответ: у Саши из них скорость больше и на х м/мин.', '<p>\nа) n : k (км/ч)<br/>\n<b>Ответ:</b> n : k км/ч её скорость, если она затратила на путь k часов.<br/><br/>\nб) a · b (км) <br/>\n<b>Ответ:</b> a · b километров она прошла за это время.<br/><br/>\nв) Витя - x : 5 (м/мин), а Саша – х : 6 (м/мин), x : 6 - х : 5 = х : (6 - 5) = х : 1 = х (м/мин)<br/>\n<b>Ответ:</b> у Саши из них скорость больше и на х м/мин.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-12/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer5.jpg', 'peterson/3/part3/page12/task5_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5204df787c180afcd3f51125c9108325dc7ddf11d49a9753eb3ddd2c991e92b6', '5,6', '["больше"]'::jsonb, 'блицтурнир*а) маша прошла n км. чему равна её скорость, если она затратила на путь k часов? б) лена шла a ч со скоростью b км/ч. какое расстояние она прошла за это время? в) витя пробежал x метров за 5 мин, а саша-за 6 мин. у кого из них скорость больше и на сколько');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 12, '6', 3, 'Какие свойства сложения и вычитания выражают данные равенства? Объясни их смысл, используя графические модели. 1) a - (b + c) = (a - b) - c = (a - c) - b 2) (a + b) - c = (a - c) + b = a + (b - c)', '</p> \n<p class="text">Какие свойства сложения и вычитания выражают данные равенства? Объясни их смысл, используя графические модели.</p> \n\n<p class="description-text"> \n1)  a - (b + c) = (a - b) - c = (a - c) - b<br/>\n2) (a + b) - c = (a - c) + b = a + (b - c)\n\n</p>\n\n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 12, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 12, номер 6, год 2022."/>\n</div>\n</div>', 'a - (b + c) = (a - b) - c = (a - c) - b − свойство вычитания числа из суммы: чтобы из числа вычесть сумму чисел, можно сначала вычесть одно из слагаемых, а затем второе. (a + b) - c = (a - c) + b = a + (b - c) − свойство вычитания числа из суммы: чтобы из суммы вычесть число, можно из первого слагаемого вычесть число и прибавить второе слагаемое, или к первому слагаемому прибавить разность второго слагаемого и числа.', '<p>\na - (b + c) = (a - b) - c = (a - c) - b − свойство вычитания числа из суммы: чтобы из числа вычесть сумму чисел, можно сначала вычесть одно из слагаемых, а затем второе.<br/>\n(a + b) - c = (a - c) + b = a + (b - c) − свойство вычитания числа из суммы: чтобы из суммы вычесть число, можно из первого слагаемого вычесть число и прибавить второе слагаемое, или к первому слагаемому прибавить разность второго слагаемого и числа.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-12/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer6.jpg', 'peterson/3/part3/page12/task6_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '058cb712acf06a18970bf29e78dfed60d3f7e7e92fd9bc9f76592216119e2acb', '1,2', NULL, 'какие свойства сложения и вычитания выражают данные равенства? объясни их смысл, используя графические модели. 1) a-(b+c)=(a-b)-c=(a-c)-b 2) (a+b)-c=(a-c)+b=a+(b-c)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 12, '7', 4, 'Вычисли наиболее удобным способом: а) 894 - (294 + 80)      в) (586 + 245) - 486      д) 232 - (95 + 132) б) 715 - 99 - 101        г) (324 + 498) - 298      е) (629 + 56) - 629', '</p> \n<p class="text">Вычисли наиболее удобным способом:</p> \n\n<p class="description-text"> \nа) 894 - (294 + 80)      в) (586 + 245) - 486      д) 232 - (95 + 132)<br/>\nб) 715 - 99 - 101        г) (324 + 498) - 298      е) (629 + 56) - 629\n</p>', 'а) 894 - (294 + 80) = 894 - 294 + 80 = 680; б) 715 - 99 - 101 = 715 - (99 + 101) = 515; в) (586 + 245) - 486 = 586 - 486 + 245 = 345; г) (324 + 498) - 298 = 324 + 498 - 298 = 524; д) 232 - (95 + 132) = 232 - 132 - 95 = 5; е) (629 + 56) - 629 = 629 - 629 + 56 = 56.', '<p>\nа) 894 - (294 + 80) = 894 - 294 + 80 = 680;<br/> 		\nб) 715 - 99 - 101 = 715 - (99 + 101) = 515;<br/>\nв) (586 + 245) - 486 = 586 - 486 + 245 = 345;<br/>\nг) (324 + 498) - 298 = 324 + 498 - 298 = 524;<br/>\nд) 232 - (95 + 132) = 232 - 132 - 95 = 5;<br/>\nе) (629 + 56) - 629 = 629 - 629 + 56 = 56.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-12/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '878e0f47d48eecc762646c338bf1616f2c9f3340aa820dd7c6b43e9a4279e99b', '56,80,95,99,101,132,232,245,294,298', '["вычисли"]'::jsonb, 'вычисли наиболее удобным способом:а) 894-(294+80)      в) (586+245)-486      д) 232-(95+132) б) 715-99-101        г) (324+498)-298      е) (629+56)-629');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 12, '8', 5, 'Реши уравнения с комментированием и сделай проверку: а) (a · 80) : 4 = 120 б) 9 · (560 : b - 5) = 27 в) (14 - c) · 4 - 9 = 19', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (a · 80) : 4 = 120<br/>  б) 9 · (560 : b - 5) = 27<br/>  в) (14 - c) · 4 - 9 = 19\n</p>', 'а) (a · 80) : 4 = 120 а = (120 · 4) : 80 а = 6 Проверка: (6 · 80) : 4 = 120 б) 9 · (560 : b - 5) = 27 b = 560 : (27 : 9 + 5) b = 560 : 8 b = 70 Проверка: 9 · (560 : 70 - 5) = 27 в) (14 - c) · 4 - 9 = 19 с = 14 - (19 + 9) : 4 с = 7 Проверка: (14 - 7) · 4 - 9 = 19', '<p>\nа) (a · 80) : 4 = 120<br/>\nа = (120 · 4) : 80<br/>\nа = 6<br/>\n<b>Проверка:</b> (6 · 80) : 4 = 120<br/><br/>\nб) 9 · (560 : b - 5) = 27 <br/>\nb = 560 : (27 : 9 + 5)<br/>\nb = 560 : 8<br/>\nb = 70<br/>\n<b>Проверка:</b> 9 · (560 : 70 - 5) = 27<br/><br/> \nв) (14 - c) · 4 - 9 = 19<br/>\nс = 14 - (19 + 9) : 4<br/>\nс = 7<br/>\n<b>Проверка:</b> (14 - 7) · 4 - 9 = 19\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-12/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '20fbdea249b54c7098abf3f12bee91481de350259b27a0d979c6d6b847a737c8', '4,5,9,14,19,27,80,120,560', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (a*80):4=120 б) 9*(560:b-5)=27 в) (14-c)*4-9=19');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 12, '9', 6, 'Составь программу действий и вычисли: а) (6543508 + 34592) : 9 - 700900 · 70 : 100 б) 81650204 - (54867 + 295 · 60) : 9 + 2989685', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) (6543508 + 34592) : 9 - 700900 · 70 : 100<br/>\nб) 81650204 - (54867 + 295 · 60) : 9 + 2989685\n</p>', 'а) (6543508 + 34592) : 9 - 700900 · 70 : 100 = 240270 6543508 + 34592 = 6578100', '<p>\nа) (6543508 + 34592) : 9 - 700900 · 70 : 100 = 240270<br/>\n6543508 + 34592 = 6578100\n</p>\n\n<div class="img-wrapper-460">\n<img width="180" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 12, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 12, номер 9, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-12/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer9.jpg', 'peterson/3/part3/page12/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2019a292303bec4ef9a00e0aba0a7973aa882fb062938fd23f72bce64b4654e1', '9,60,70,100,295,34592,54867,700900,2989685,6543508', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) (6543508+34592):9-700900*70:100 б) 81650204-(54867+295*60):9+2989685');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 12, '10', 7, '1 января 2018 года было понедельником. Каким днём недели будет 1 января 2019 года, 1 января 2020 года, 1 января 2021 года?', '</p> \n<p class="text">1 января 2018 года было понедельником. Каким днём недели будет 1 января 2019 года, 1 января 2020 года, 1 января 2021 года? </p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer10.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 12, номер 10, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 12, номер 10, год 2022."/>\n</div>\n\n</div>', '1 января 2019 года – это вторник, 1 января 2020 года – среда, 1 января 2021 года – пятница. Дни идут по порядку, кроме 1 января 2021. 2020 является високосным годом - в нём 366 дней - добавляется 29 февраля).', '<p>\n1 января 2019 года – это вторник, 1 января 2020 года – среда, 1 января 2021 года – пятница. Дни идут по порядку, кроме 1 января 2021. 2020 является високосным годом - в нём 366 дней - добавляется 29 февраля).\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-12/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica12-nomer10.jpg', 'peterson/3/part3/page12/task10_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0b2259b31040eb77840bfadedee12c4ff79f5557e036005f29b1118ea533d193', '1,2018,2019,2020,2021', NULL, '1 января 2018 года было понедельником. каким днём недели будет 1 января 2019 года, 1 января 2020 года, 1 января 2021 года');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 13, '1', 0, 'а) Определи по рисунку, из какого города вышел поезд. Куда и с какой скоростью он идёт? Построй числовой луч и покажи на нём движение поезда. б) Пусть s – путь, который прошёл поезд, d – его расстояние до Вологды, D – расстояние до Калуги. Заполни таблицу. Запиши формулы зависимостей каждой из величин s, d, D от времени движения t.', '</p> \n<p class="text">а) Определи по рисунку, из какого города вышел поезд. Куда и с какой скоростью он идёт? Построй числовой луч и покажи на нём движение поезда.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 13, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 13, номер 1, год 2022."/>\n</div>\n\n</div>\n\n\n\n<p class="text">б) Пусть s – путь, который прошёл поезд, d – его расстояние до Вологды, D – расстояние до Калуги. Заполни таблицу. Запиши формулы зависимостей каждой из величин s,  d,  D от времени движения t.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer1-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 13, номер 1-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 13, номер 1-1, год 2022."/>\n</div>\n</div>', 'а) Поезд вышел из города Москва в Вологду со скоростью 80 км/ч. б) s = 17 · t, d = n · t, D = 102 - 17 · t', '<p>\nа) Поезд вышел из города Москва в Вологду со скоростью 80 км/ч.<br/>\nб) s = 17 · t, d = n · t, D = 102 - 17 · t\n\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer1-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 13, номер 1-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 13, номер 1-2, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Определи по рисунку, из какого города вышел поезд. Куда и с какой скоростью он идёт? Построй числовой луч и покажи на нём движение поезда.","solution":"Поезд вышел из города Москва в Вологду со скоростью 80 км/ч."},{"letter":"б","condition":"Пусть s – путь, который прошёл поезд, d – его расстояние до Вологды, D – расстояние до Калуги. Заполни таблицу. Запиши формулы зависимостей каждой из величин s, d, D от времени движения t.","solution":"s = 17 · t, d = n · t, D = 102 - 17 · t"}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-13/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer1.jpg', 'peterson/3/part3/page13/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer1-1.jpg', 'peterson/3/part3/page13/task1_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer1-2.jpg', 'peterson/3/part3/page13/task1_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '32f64c9ea7d5380357489036cd34e1f44ba142e0082b55eed725b0af051860be', NULL, '["заполни","числовой луч"]'::jsonb, 'а) определи по рисунку, из какого города вышел поезд. куда и с какой скоростью он идёт? построй числовой луч и покажи на нём движение поезда. б) пусть s-путь, который прошёл поезд, d-его расстояние до вологды, d-расстояние до калуги. заполни таблицу. запиши формулы зависимостей каждой из величин s, d, d от времени движения t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 13, '2', 1, 'а) Определи по рисунку, из какого города выехал автобус. Куда и с какой скоростью он едет? Построй числовой луч и покажи на нём движение автобуса. б) Пусть s – путь, который прошёл автобус, d – его расстояние до Брянска, D – расстояние до Воронежа. Заполни таблицу. Запиши формулы зависимостей каждой из величин s, d, D от времени движения t.', '</p> \n<p class="text">а) Определи по рисунку, из какого города выехал автобус. Куда и с какой скоростью он едет? Построй числовой луч и покажи на нём движение автобуса.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 13, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 13, номер 2, год 2022."/>\n</div>\n\n</div>\n\n\n<p class="text">б) Пусть s – путь, который прошёл автобус, d – его расстояние до Брянска, D – расстояние до Воронежа. Заполни таблицу. Запиши формулы зависимостей каждой из величин s, d, D от времени движения t.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer2-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 13, номер 2-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 13, номер 2-1, год 2022."/>\n</div>\n\n</div>', 'а) Автобус выехал из города Москва в Брянск со скоростью 60 км/ч. б) s = 60 · t, d = 480 - 60 · t, D = 160 + 60 · t', '<p>\nа) Автобус выехал из города Москва в Брянск со скоростью 60 км/ч.<br/>\nб) s = 60 · t, d = 480 - 60 · t, D = 160 + 60 · t\n\n</p>\n\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer2-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 13, номер 2-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 13, номер 2-2, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Определи по рисунку, из какого города выехал автобус. Куда и с какой скоростью он едет? Построй числовой луч и покажи на нём движение автобуса.","solution":"Автобус выехал из города Москва в Брянск со скоростью 60 км/ч."},{"letter":"б","condition":"Пусть s – путь, который прошёл автобус, d – его расстояние до Брянска, D – расстояние до Воронежа. Заполни таблицу. Запиши формулы зависимостей каждой из величин s, d, D от времени движения t.","solution":"s = 60 · t, d = 480 - 60 · t, D = 160 + 60 · t"}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-13/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer2.jpg', 'peterson/3/part3/page13/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer2-1.jpg', 'peterson/3/part3/page13/task2_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer2-2.jpg', 'peterson/3/part3/page13/task2_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '31df72f588e53beba4829cd04346433357197eecd5ba6d0dd284b9454a9d7162', NULL, '["заполни","числовой луч"]'::jsonb, 'а) определи по рисунку, из какого города выехал автобус. куда и с какой скоростью он едет? построй числовой луч и покажи на нём движение автобуса. б) пусть s-путь, который прошёл автобус, d-его расстояние до брянска, d-расстояние до воронежа. заполни таблицу. запиши формулы зависимостей каждой из величин s, d, d от времени движения t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 13, '3', 2, 'Выполни действия. Проверь результаты с помощью калькулятора. а) 49237 + 181048 в) 700 · 209530 б) 6080010 – 5550481 г) 60002400 : 80', '</p> \n<p class="text">\nВыполни действия. Проверь результаты с помощью калькулятора.\n</p> \n\n<p class="description-text"> \nа) 49237 + 181048  в) 700 · 209530 <br/>\nб) 6080010 – 5550481  г) 60002400 : 80\n</p>', 'а) 49237 + 181048 = 230285', '<p>\nа) 49237 + 181048 = 230285\n</p>\n\n<div class="img-wrapper-460">\n<img width="180" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 13, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 13, номер 3, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-13/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica13-nomer3.jpg', 'peterson/3/part3/page13/task3_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ba6381c26de8c5ff35ba40ba1911422bf34d3008e6847c0f2c2a79c106a5170b', '80,700,49237,181048,209530,5550481,6080010,60002400', NULL, 'выполни действия. проверь результаты с помощью калькулятора. а) 49237+181048 в) 700*209530 б) 6080010-5550481 г) 60002400:80');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 14, '4', 0, 'Расстояние от посёлка Солнечное до Тучково 18 км, а от Тучково до Маросейкино – в 4 раза больше. Автобус едет из Солнечного в Маросейкино через Тучково со скоростью 45 км/ч. За какое время он проедет весь этот путь?', '</p> \n<p class="text">Расстояние от посёлка Солнечное до Тучково 18 км, а от Тучково до Маросейкино – в 4 раза больше. Автобус едет из Солнечного в Маросейкино через Тучково со скоростью 45 км/ч. За какое время он проедет весь этот путь?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 14, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 14, номер 4, год 2022."/>\n</div>\n</div>', '(18 + 18 · 4) : 45 = (18 + 72) : 45 = 90 : 45 = 2 (ч) Ответ: за 2 часа он проедет весть этот путь.', '<p>\n(18 + 18 · 4) : 45 = (18 + 72) : 45 = 90 : 45 = 2 (ч)<br/>\n<b>Ответ:</b> за 2 часа он проедет весть этот путь.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-14/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer4.jpg', 'peterson/3/part3/page14/task4_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0816e511e1d099d4fb6902855ef91cbac16eefc1c0664ead9dc67b32bccabec1', '4,18,45', '["больше","раз","раза"]'::jsonb, 'расстояние от посёлка солнечное до тучково 18 км, а от тучково до маросейкино-в 4 раза больше. автобус едет из солнечного в маросейкино через тучково со скоростью 45 км/ч. за какое время он проедет весь этот путь');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 14, '5', 1, 'Стоянка геологов находится на расстоянии 250 км от города. Что бы добраться до стоянки, геологи сначала ехали из города 3 ч на машине со скоростью 72 км/ч, затем 2 ч ехали на лошадях со скоростью 9 км/ч, а потом 4 ч шли пешком. С какой скоростью они шли пешком?', '</p> \n<p class="text">Стоянка геологов находится на расстоянии 250 км от города. Что бы добраться до стоянки, геологи сначала ехали из города 3 ч на машине со скоростью 72 км/ч, затем 2 ч ехали на лошадях со скоростью 9 км/ч, а потом 4 ч шли пешком. С какой скоростью они шли пешком?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 14, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 14, номер 5, год 2022."/>\n</div>\n</div>', '(250 - 72 · 3 - 9 · 2) : 4 = (250 - 216 - 18) : 4 = (34 - 18) : 4 = 16 : 4 = 4 (км/ч) Ответ: 4 км/ч они шли пешком.', '<p>\n(250 - 72 · 3 - 9 · 2) : 4 = (250 - 216 - 18) : 4 = (34 - 18) : 4 = 16 : 4 = 4 (км/ч)<br/>\n<b>Ответ:</b> 4 км/ч они шли пешком.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-14/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer5.jpg', 'peterson/3/part3/page14/task5_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6a0d14cfe45beab688885c85c3f45ff2cff6abee446703de88ec9c40608508d9', '2,3,4,9,72,250', NULL, 'стоянка геологов находится на расстоянии 250 км от города. что бы добраться до стоянки, геологи сначала ехали из города 3 ч на машине со скоростью 72 км/ч, затем 2 ч ехали на лошадях со скоростью 9 км/ч, а потом 4 ч шли пешком. с какой скоростью они шли пешком');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 14, '6', 2, 'Реши уравнения с комментированием: а) 540 : (17 - x) = 60 б) (8 · y - 30) : 9 = 50', '</p> \n<p class="text">Реши уравнения с комментированием:</p> \n\n<p class="description-text"> \nа) 540 : (17 - x) = 60<br/>   \nб) (8 · y - 30) : 9 = 50\n</p>', 'а) 540 : (17 - x) = 60 чтобы найти неизвестное делимое 8 · y - 30, нужно частное умножить на делитель х = 540 : 60 + 17 х = 9 + 17 х = 26 б) (8 · y - 30) : 9 = 50 чтобы найти неизвестное уменьшаемое 8 · y, нужно к разности прибавить вычитаемое у = (50 · 9 + 30) : 8 у = (450 + 30) : 8 у = 480 : 8 у = 60', '<p>\nа) 540 : (17 - x) = 60  <br/> \nчтобы найти неизвестное делимое 8 · y - 30, нужно частное умножить на делитель<br/>   \nх = 540 : 60 + 17<br/> \nх = 9 + 17<br/> \nх = 26  <br/> <br/> \nб) (8 · y - 30) : 9 = 50<br/> \nчтобы найти неизвестное уменьшаемое 8 · y, нужно к разности прибавить вычитаемое<br/> \nу = (50 · 9 + 30) : 8<br/> \nу = (450 + 30) : 8<br/> \nу = 480 : 8<br/> \nу = 60\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-14/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0108d0216fe71d2ba06edc250e7e7662b78465186ef95c2b17977fd29c168139', '8,9,17,30,50,60,540', '["реши"]'::jsonb, 'реши уравнения с комментированием:а) 540:(17-x)=60 б) (8*y-30):9=50');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 14, '7', 3, 'Выполни действия. Расположи ответы примеров в порядке возрастания и расшифруй имя героя книги. Кто это?', '</p> \n<p class="text">Выполни действия. Расположи ответы примеров в порядке возрастания и расшифруй имя героя книги. Кто это?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 14, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 14, номер 7, год 2022."/>\n</div>\n</div>', 'Л - 48756 + 192317 + 392 = 241073 + 392 = 241465', '<p>\nЛ - 48756 + 192317 + 392 = 241073 + 392 = 241465\n</p>\n\n<div class="img-wrapper-460">\n<img width="170" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer7-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 14, номер 7-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 14, номер 7-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-14/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer7.jpg', 'peterson/3/part3/page14/task7_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer7-1.jpg', 'peterson/3/part3/page14/task7_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e781c126d7db745f6f441ed1d57e1f80fd2eda7b58077b19f0d84d296a685214', NULL, NULL, 'выполни действия. расположи ответы примеров в порядке возрастания и расшифруй имя героя книги. кто это');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 14, '8', 4, 'Запиши множество делителей и множество кратных числа 14.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 14.</p>', 'Множество делителей числа 14: 1, 2, 7, 14. Множество кратных числа 14: 14, 28, 42, 56, 70, 84, 98, 102, 116, 130, 144.', '<p>\nМножество делителей числа 14: 1, 2, 7, 14.<br/>  \nМножество кратных числа 14: 14, 28, 42, 56, 70, 84, 98, 102, 116, 130, 144.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-14/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f422d3bc5b3ec9d9a2a0e54e3639d9f7a90643c0a6f51f3b579d56edcbc4e049', '14', NULL, 'запиши множество делителей и множество кратных числа 14');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 14, '9', 5, 'A – множество остатков, которые могут получиться при делении на 5, а B – множество остатков, возможных при делении на 7. а) Задай множества A и B перечислением и запиши элементы с помощью фигурных скобок. б) Построй диаграмму Эйлера – Венна множеств A и B. Какое из множеств является подмножеством другого? в) Найди A ⋂ B и A ⋃ B.', '</p> \n<p class="text">A – множество остатков, которые могут получиться при делении на 5, а B – множество остатков, возможных при делении на 7. <br/> \nа) Задай множества A и B перечислением и запиши элементы с помощью фигурных скобок.<br/> \nб) Построй диаграмму Эйлера – Венна множеств A и B. Какое из множеств является подмножеством другого?<br/> \nв) Найди A ⋂ B и A ⋃ B.\n</p>', 'а) А = {0; 1; 2; 3; 4}, B = {0; 1; 2; 3; 4; 5; 6} б)', '<p>\nа) А = {0; 1; 2; 3; 4}, B = {0; 1; 2; 3; 4; 5; 6}<br/> <br/> \nб) \n\n</p>\n\n<div class="img-wrapper-460">\n<img width="200" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 14, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 14, номер 9, год 2022."/>\n\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer10-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 14, номер 10-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 14, номер 10-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-14/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer9.jpg', 'peterson/3/part3/page14/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica14-nomer10-1.jpg', 'peterson/3/part3/page14/task9_solution_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '4e1260b4395e8c047a68509d9a24bd347975e1cead660e0b2003d82f7c2ab23d', '5,7', '["найди"]'::jsonb, 'a-множество остатков, которые могут получиться при делении на 5, а b-множество остатков, возможных при делении на 7. а) задай множества a и b перечислением и запиши элементы с помощью фигурных скобок. б) построй диаграмму эйлера-венна множеств a и b. какое из множеств является подмножеством другого? в) найди a ⋂ b и a ⋃ b');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 15, '1', 0, 'Пчела Майя стала соединять формулы с их названиями. Все линии перепутались. Определи, правильно ли пчела Майя выполнила задание. Какие ещё формулы ты знаешь?', '</p> \n<p class="text">Пчела Майя стала соединять формулы с их названиями. Все линии перепутались. Определи, правильно ли пчела Майя выполнила задание.<br/>\nКакие ещё формулы ты знаешь?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica15-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 15, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 15, номер 1, год 2022."/>\n</div>\n</div>', 'Правильно пчела Майя выполнила задание.', '<p>\nПравильно пчела Майя выполнила задание.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-15/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica15-nomer1.jpg', 'peterson/3/part3/page15/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e6ee4c5a275b5601439228d1d51448392e18491268f1c7a9f011b323feaddc95', NULL, NULL, 'пчела майя стала соединять формулы с их названиями. все линии перепутались. определи, правильно ли пчела майя выполнила задание. какие ещё формулы ты знаешь');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 15, '2', 1, 'Прочитай задачу и объясни, как составлена таблица: «Заяц сначала бежал 2 ч со скоростью 24 км/ч, затем 3 ч ехал на велосипеде, а после этого 5 ч ехал на поезде со скоростью 48 км/ч. Всего заяц пробежал и проехал 357 км. С какой скоростью он ехал на велосипеде?» Используя таблицу, ответь на вопросы: а) Какой путь пробежал заяц за первые 2 ч? б) Какой путь он проехал на поезде за последние 5 ч? в) Какой путь проехал заяц на велосипеде за 3 ч? г) С какой скоростью он ехал на велосипеде? Составь план решения задачи и запиши решение в тетради. Сделай вывод: как можно решить задачу с помощью таблицы?', '</p> \n<p class="text">Прочитай задачу и объясни, как составлена таблица:<br/>\n«Заяц сначала бежал 2 ч со скоростью 24 км/ч, затем 3 ч ехал на велосипеде, а после этого 5 ч ехал на поезде со скоростью 48 км/ч. Всего заяц пробежал и проехал 357 км. С какой скоростью он ехал на велосипеде?»<br/><br/>\n\nИспользуя таблицу, ответь на вопросы:<br/>\nа) Какой путь пробежал заяц за первые 2 ч?<br/>\nб) Какой путь он проехал на поезде за последние 5 ч?<br/>\nв) Какой путь проехал заяц на велосипеде за 3 ч?<br/>\nг) С какой скоростью он ехал на велосипеде?<br/>\nСоставь план решения задачи и запиши решение в тетради. Сделай вывод: как можно решить задачу с помощью таблицы?\n\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica15-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 15, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 15, номер 2, год 2022."/>\n</div>\n</div>', 'а) 2 · 24 = 48 (км) б) 5 · 48 = 240 (км) в) 357 - (240 + 48) = 69 (км) г) 69 : 3 = 23 (км) План 1) Надо узнать какой путь он проехал на велосипеде, для этого из всего пути отнимем путь, который он бежал и ехал на поезде. 2) Поделим путь на время (3 ч) 3) Узнаем скорость передвижения на велосипеде. Решение: 1) 2 · 24 = 48 (км) – путь который пробежал. 2) 5 · 48 = 240 (км) – путь который проехал на поезде. 3) 357 - (240 + 48) = 69 (км) – путь, который проехал на велосипеде. 4) 69 : 3 = 23 (км/ч) – скорость зайца на велосипеде. Ответ: 23 км/ч.', '<p>\nа) 2 · 24 = 48 (км)<br/> \nб) 5 · 48 = 240 (км) <br/>\nв) 357 - (240 + 48) = 69 (км)<br/>\nг) 69 : 3 = 23 (км)<br/><br/>\nПлан<br/>\n1) Надо узнать какой путь он проехал на велосипеде, для этого из всего пути отнимем путь, который он бежал и ехал на поезде.<br/>\n2) Поделим путь на время (3 ч)<br/>\n3) Узнаем скорость передвижения на велосипеде.<br/><br/>\nРешение:<br/>\n1) 2 · 24 = 48 (км) – путь который пробежал.<br/>\n2) 5 · 48 = 240 (км) – путь который проехал на поезде. <br/>\n3) 357 - (240 + 48) = 69 (км) – путь, который проехал на велосипеде.<br/>\n4) 69 : 3 = 23 (км/ч) – скорость зайца на велосипеде.<br/>\n<b>Ответ:</b> 23 км/ч.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-15/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica15-nomer2.jpg', 'peterson/3/part3/page15/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0e4d0cc61ad362b166b2c109872ca2afea7a30f58e4aef586c4b68e672daacf6', '2,3,5,24,48,357', '["реши"]'::jsonb, 'прочитай задачу и объясни, как составлена таблица:"заяц сначала бежал 2 ч со скоростью 24 км/ч, затем 3 ч ехал на велосипеде, а после этого 5 ч ехал на поезде со скоростью 48 км/ч. всего заяц пробежал и проехал 357 км. с какой скоростью он ехал на велосипеде?" используя таблицу, ответь на вопросы:а) какой путь пробежал заяц за первые 2 ч? б) какой путь он проехал на поезде за последние 5 ч? в) какой путь проехал заяц на велосипеде за 3 ч? г) с какой скоростью он ехал на велосипеде? составь план решения задачи и запиши решение в тетради. сделай вывод:как можно решить задачу с помощью таблицы');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 16, '3', 0, 'Составь в тетради таблицы и реши задачи: а) Вертолёт пролетает 840 км за 3 ч, а автомобиль проходит это же расстояние за 7 ч. Чья скорость больше и на сколько? б) Поезд проходит 320 км за 5 ч. Какое расстояние он пройдёт за 8 ч, двигаясь с этой же скоростью? в) Караван верблюдов шёл в первый день 8 ч со скоростью 9 км/ч, во второй день – 6 ч со скоростью 8 км/ч, а в третий день – 9 ч со скоростью 7 км/ч. Какое расстояние прошёл караван за 3 дня?', '</p> \n<p class="text">Составь в тетради таблицы и реши задачи:<br/>\nа) Вертолёт пролетает 840 км за 3 ч, а автомобиль проходит это же расстояние за 7 ч. Чья скорость больше и на сколько?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 16, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 16, номер 3, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Поезд проходит 320 км за 5 ч. Какое расстояние он пройдёт за 8 ч, двигаясь с этой же скоростью?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer3-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 16, номер 3-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 16, номер 3-1, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">в) Караван верблюдов шёл в первый день 8 ч со скоростью 9 км/ч, во второй день – 6 ч со скоростью 8 км/ч, а в третий день – 9 ч со скоростью 7 км/ч. Какое расстояние прошёл караван за 3 дня?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer3-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 16, номер 3-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 16, номер 3-2, год 2022."/>\n</div>\n</div>', 'а) 840 : 3 - 840 : 7 = 280 - 120 = 160 (км/ч)', '<p>\nа) 840 : 3 - 840 : 7 = 280 - 120 = 160 (км/ч)\n</p>\n\n<div class="img-wrapper-460">\n<img width="300" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer3-3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 16, номер 3-3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 16, номер 3-3, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-16/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer3.jpg', 'peterson/3/part3/page16/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer3-1.jpg', 'peterson/3/part3/page16/task3_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 2, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer3-2.jpg', 'peterson/3/part3/page16/task3_condition_2.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer3-3.jpg', 'peterson/3/part3/page16/task3_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd785e691514785e870cc70d49e1a04047f23f3e09317bcf64b887bfd1c1aba60', '3,5,6,7,8,9,320,840', '["реши","больше"]'::jsonb, 'составь в тетради таблицы и реши задачи:а) вертолёт пролетает 840 км за 3 ч, а автомобиль проходит это же расстояние за 7 ч. чья скорость больше и на сколько? б) поезд проходит 320 км за 5 ч. какое расстояние он пройдёт за 8 ч, двигаясь с этой же скоростью? в) караван верблюдов шёл в первый день 8 ч со скоростью 9 км/ч, во второй день-6 ч со скоростью 8 км/ч, а в третий день-9 ч со скоростью 7 км/ч. какое расстояние прошёл караван за 3 дня');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 16, '4', 1, 'Реши уравнения с комментированием и сделай проверку: а) x · 7 - 80 = 340 б) (900 - y) : 9 = 80 в) (350 : y + 10) · 7 = 560', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) x · 7 - 80 = 340<br/>        \nб) (900 - y) : 9 = 80 <br/>       \nв) (350 : y + 10) · 7 = 560\n</p>', 'а) x · 7 - 80 = 340 Чтобы найти неизвестное уменьшаемое х · 7, нужно к разности прибавить вычитаемое х · 7 = 340 + 80 х · 7 = 420 Чтобы найти неизвестный множитель х, нужно произведение разделить на известный множитель х = 420 : 7 х = 60 Проверка: 60 · 7 − 80 = 340 б) (900 - y) : 9 = 80 Чтобы найти делимое (900 - y) надо делитель умножить на частное 900 - y = 9 · 80 900 - у = 720 Чтобы найти вычитаемое у надо вычесть из уменьшаемого разность у = 900 - 720 у = 180 Проверка: (900 - 180) : 9 = 80 в) (350 : y + 10) · 7 = 560 Чтобы найти неизвестный множитель 350 : у + 10, нужно произведение разделить на известный множитель 350 : у + 10 = 560 : 7 350 : у + 10 = 80 Чтобы найти неизвестное слагаемое, нужно из суммы вычесть известное 350 : у = 80 - 10 350 : у = 70 Чтобы найти делитель у надо делимое разделить на частное у = 350 : 70 у = 5 Проверка: (350 : 5 + 10) · 7 = 560', '<p>\nа) x · 7 - 80 = 340   <br/>  \nЧтобы найти неизвестное уменьшаемое х · 7, нужно к разности прибавить вычитаемое<br/> \nх · 7 = 340 + 80 <br/> \nх · 7 = 420 <br/> \nЧтобы найти неизвестный множитель х, нужно произведение разделить на известный множитель<br/> \nх = 420 : 7 <br/> \nх = 60<br/> \n<b>Проверка:</b> 60 · 7 − 80 = 340 <br/> <br/>     \nб) (900 - y) : 9 = 80 <br/>  \nЧтобы найти делимое (900 - y)  надо делитель умножить на частное<br/> \n900 - y = 9 · 80<br/> \n900 - у = 720<br/>\nЧтобы найти вычитаемое у надо вычесть из уменьшаемого разность<br/>\nу = 900 - 720<br/>\nу = 180  <br/> \n<b>Проверка:</b> (900 - 180) : 9 = 80  <br/><br/> \nв) (350 : y + 10) · 7 = 560<br/>\nЧтобы найти неизвестный множитель 350 : у + 10, нужно произведение разделить на известный множитель<br/>\n350 : у + 10 = 560 : 7 <br/>\n350 : у + 10 = 80 <br/>\nЧтобы найти неизвестное слагаемое, нужно из суммы вычесть известное<br/>\n350 : у = 80 - 10<br/>\n350 : у = 70<br/>\nЧтобы найти делитель у надо делимое разделить на частное<br/>\nу = 350 : 70<br/>\nу = 5<br/>\n<b>Проверка:</b> (350 : 5 + 10) · 7 = 560 \n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-16/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1637794306c3c6707b7b4d9e0c121428c4bddeb0ff2e0d03886bbc874c37525d', '7,9,10,80,340,350,560,900', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) x*7-80=340 б) (900-y):9=80 в) (350:y+10)*7=560');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 16, '5', 2, 'Прочитай числа и расположи их в порядке возрастания: 94517, 3896, 3002650, 302650, 32650 Найди разность наибольшего и наименьшего из этих чисел.', '</p> \n<p class="text">Прочитай числа и расположи их в порядке возрастания:<br/>\n94517, 3896, 3002650, 302650, 32650<br/>\nНайди разность наибольшего и наименьшего из этих чисел.\n</p>', 'Девяносто четыре тысячи пятьсот семнадцать, три тысячи восемьсот девяносто шесть, три миллиона две тысячи шестьсот пятьдесят, триста две тысячи шестьсот пятьдесят, тридцать две тысячи шестьсот пятьдесят 3896, 32650, 94517, 302650, 3002650 3002650 - 3896 = 2998754', '<p>\nДевяносто четыре тысячи пятьсот семнадцать, три тысячи восемьсот девяносто шесть, три миллиона две тысячи шестьсот пятьдесят, триста две тысячи шестьсот пятьдесят, тридцать две тысячи шестьсот пятьдесят<br/>\n3896, 32650, 94517, 302650, 3002650<br/>\n3002650 - 3896 = 2998754\n</p>\n\n<div class="img-wrapper-460">\n<img width="190" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 16, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 16, номер 5, год 2022."/>', 'Алгоритм решения задач с помощью таблиц 1. Внимательно прочитать условие задачи. 2. Отметить в таблице известные и неизвестные величины. 3. Составить план решения задачи: какие неизвестные величины и в каком порядке нужно найти. 4. Решить задачу по плану.', '<div class="recomended-block">\n<span class="title">Алгоритм решения задач с помощью таблиц</span>\n<p>\n1. Внимательно прочитать условие задачи.<br/>\n2. Отметить в таблице известные и неизвестные величины.<br/>\n3. Составить план решения задачи: какие неизвестные величины и в каком порядке нужно найти.<br/>\n4. Решить задачу по плану.\n\n</p>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-16/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica16-nomer5.jpg', 'peterson/3/part3/page16/task5_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b1e0c684251eae2d6f6434f26af1a7fd917fce6a03679570572feba236888cbb', '3896,32650,94517,302650,3002650', '["найди","разность","больше","меньше","раз"]'::jsonb, 'прочитай числа и расположи их в порядке возрастания:94517, 3896, 3002650, 302650, 32650 найди разность наибольшего и наименьшего из этих чисел');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 17, '6', 0, '1) Запиши число 40560 в виде суммы разрядных слагаемых. 2) Сколько единиц в разряде сотен числа 40560? Сколько всего сотен в этом числе? Вырази его: а) в сотнях и единицах; б) в тысячах и единицах. 3) Вырази величины в указанных единицах измерения: 40560 м = … км … м        40560 кг = … ц … кг 40560 кг = … т … кг        40560 мм = … дм … мм 40560 мм = … м … мм     40560 г = … кг … г', '</p> \n<p class="text">1) Запиши число 40560 в виде суммы разрядных слагаемых.<br/>\n2) Сколько единиц в разряде сотен числа 40560? Сколько всего сотен в этом числе? Вырази его: а) в сотнях и единицах; б) в тысячах и единицах.<br/>\n3) Вырази величины в указанных единицах измерения:<br/>\n40560 м = … км … м        40560 кг = … ц … кг<br/>\n40560 кг = … т … кг        40560 мм = … дм … мм<br/> \n40560 мм = … м … мм     40560 г = … кг … г\n</p>', '1) 40000 + 500 + 60 2) 5 единиц в разряде сотен, всего сотен 405. а) 405 сотен 60 единиц б) 40 тысяч 560 единиц 3) 40560 м = 40 км 560 м     40560 кг = 405 ц 60 кг 40560 кг = 40 т 560 кг          40560 мм = 405 дм 60 мм 40560 мм = 40 м 560 мм      40560 г = 40 кг 560 г', '<p>\n1) 40000 + 500 + 60<br/>\n2) 5 единиц в разряде сотен, всего сотен 405.<br/>\nа) 405 сотен 60 единиц б) 40 тысяч 560 единиц<br/> \n3) 40560 м = 40 км 560 м     40560 кг = 405 ц 60 кг<br/>\n40560 кг = 40 т 560 кг          40560 мм = 405 дм 60 мм<br/> \n40560 мм = 40 м 560 мм      40560 г = 40 кг 560 г\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-17/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8b68570bb69e509e175ef1d9b8481254eaf4acc9ec87d8d78f29de8a911bf687', '1,2,3,40560', '["раз"]'::jsonb, '1) запиши число 40560 в виде суммы разрядных слагаемых. 2) сколько единиц в разряде сотен числа 40560? сколько всего сотен в этом числе? вырази его:а) в сотнях и единицах; б) в тысячах и единицах. 3) вырази величины в указанных единицах измерения:40560 м=... км ... м        40560 кг=... ц ... кг 40560 кг=... т ... кг        40560 мм=... дм ... мм 40560 мм=... м ... мм     40560 г=... кг ... г');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 17, '7', 1, 'Найди пропущенные цифры при делении с остатком углом. Сделай проверку по формуле деления с остатком: a = b · c + r, r < b.', '</p> \n<p class="text">Найди пропущенные цифры при делении с остатком углом. Сделай проверку по формуле деления с остатком: a = b · c + r, r &lt; b.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica17-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 17, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 17, номер 7, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica17-nomer7-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 17, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 17, номер 7, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-17/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica17-nomer7.jpg', 'peterson/3/part3/page17/task7_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica17-nomer7-1.jpg', 'peterson/3/part3/page17/task7_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '30c62a56d56d5ff57533c38d62e91db0a827a590d77120e4e2d9dd78bc27934b', NULL, '["найди"]'::jsonb, 'найди пропущенные цифры при делении с остатком углом. сделай проверку по формуле деления с остатком:a=b*c+r, r<b');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 17, '8', 2, 'Запиши множество делителей и множество кратных числа 15', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 15</p>', 'Множество делителей числа 15 = {1, 3, 5, 15}. Множество кратных числа 15 = {15, 30, 45, 60, 75}.', '<p>\nМножество делителей числа 15 = {1, 3, 5, 15}.<br/> \nМножество кратных числа 15 = {15, 30, 45, 60, 75}.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-17/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '20d59eb3f3c3de14bcdf20fb7029fe317085e8a6f89032872e7b807270666fac', '15', NULL, 'запиши множество делителей и множество кратных числа 15');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 17, '9', 3, 'Расшифруй имя славного защитника Руси. Что ты о нём знаешь?', '</p> \n<p class="text">Расшифруй имя славного защитника Руси. Что ты о нём знаешь?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica17-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 17, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 17, номер 9, год 2022."/>\n</div>\n</div>', 'Р - 839 - 625 = 214 У - 247 + 53 = 300 О - 400 - 265 = 135 Ь - 218 + 26 = 244 Л - 325 - 43 = 282 Я - 350 : 7 · 8 = 400 Ц - 9 · 4 + 82 = 118 К - 172 - 72 : 4 = 154 Е - 567 - 60 · 4 = 327 И - (320 : 40) · 8 = 64 Т - 900 : (25 · 6) = 6 М - 90 · 2 : 30 · 70 = 420 31 + 30 + 31 = 61 + 31 = 92 (дня) − в четвертом квартале. Високосный год: меньше всего дней в I и II кварталах - 91 день, больше всего − в III и IV кварталах − 92 дня. Обычный год: меньше всего дней в I квартале - 90 дней, больше всего − в III и IV кварталах − 92 дня.', '<p>\nР - 839 - 625 = 214<br/>\nУ - 247 + 53 = 300<br/>\nО - 400 - 265 = 135<br/>\nЬ - 218 + 26 = 244<br/>\nЛ - 325 - 43 = 282<br/>\nЯ - 350 : 7 · 8 = 400<br/>\nЦ - 9 · 4 + 82 = 118<br/>\nК - 172 - 72 : 4 = 154<br/>\nЕ - 567 - 60 · 4 = 327<br/>\nИ - (320 : 40) · 8 = 64<br/>\nТ - 900 : (25 · 6) = 6<br/>\nМ - 90 · 2 : 30 · 70 = 420\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica17-nomer9-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 17, номер 9-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 17, номер 9-1, год 2022."/>\n\n\n<p>\n31 + 30 + 31 = 61 + 31 = 92 (дня) − в четвертом квартале. <br/>\nВисокосный год: меньше всего дней в I и II кварталах - 91 день, больше всего − в III и IV кварталах − 92 дня. <br/>\nОбычный год: меньше всего дней в I квартале - 90 дней, больше всего − в III и IV кварталах − 92 дня.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-17/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica17-nomer9.jpg', 'peterson/3/part3/page17/task9_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica17-nomer9-1.jpg', 'peterson/3/part3/page17/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8de288e4db140a442ed514a9138a24f74ce0e569c5414f9209c18a31d2af30ae', NULL, NULL, 'расшифруй имя славного защитника руси. что ты о нём знаешь');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 18, '1', 0, 'а) Ира прошла 320 м за 5 мин, а Петя – 225 м за 3 мин. У кого из ребят скорость больше и на сколько? б) Орёл за 9 с пролетел 270 м, а сокол за это же время пролетел 189 м. На сколько метров в секунду скорость сокола меньше скорости орла? в) Первый лыжник за 3 ч пробежал 51 км, а второй лыжник за это же время пробежал на 6 км больше. На сколько километров в час скорость второго лыжника больше скорости первого?', '</p> \n<p class="text">а) Ира прошла 320 м за 5 мин, а Петя – 225 м за 3 мин. У кого из ребят скорость больше и на сколько?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica18-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 18, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 18, номер 1, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Орёл за 9 с пролетел 270 м, а сокол за это же время пролетел 189 м. На сколько метров в секунду скорость сокола меньше скорости орла?<br/><br/>\nв) Первый лыжник за 3 ч пробежал 51 км, а второй лыжник за это же время пробежал на 6 км больше. На сколько километров в час скорость второго лыжника больше скорости первого?\n</p>', 'а) Ира 320 : 5 = 64 (м/мин), а Петя – 225 : 3 = 75 (м/мин). 75 - 64 = 11 (м/мин) Ответ: у Пети скорость больше на 11 м/мин. б) Орёл - 270 : 9 = 30 (м/с), а сокол – 189 : 9 = 21 (м/с). 30 - 21 = 9 (м/с) Ответ: на 9 метров в секунду скорость сокола меньше скорости орла. в) Первый лыжник – 51 : 3 = 17 (км/ч), а второй лыжник – (51 + 6) : 3 = 57 : 3 = 19 (км/ч). 19 - 17 = 2 (км/ч) Ответ: на 2 километров в час скорость второго лыжника больше скорости первого.', '<p>\nа) Ира 320 : 5 = 64 (м/мин), а Петя – 225 : 3 = 75 (м/мин). 75 - 64 = 11 (м/мин)<br/>\n<b>Ответ:</b> у Пети скорость больше на 11 м/мин.<br/><br/>\nб) Орёл - 270 : 9 = 30 (м/с), а сокол – 189 : 9 = 21 (м/с). 30 - 21 = 9 (м/с)<br/>\n<b>Ответ:</b> на 9 метров в секунду скорость сокола меньше скорости орла.<br/><br/>\nв) Первый лыжник – 51 : 3 = 17 (км/ч), а второй лыжник – (51 + 6) : 3 = 57 : 3 = 19 (км/ч). 19 - 17 = 2 (км/ч)<br/>\n<b>Ответ:</b> на 2 километров в час скорость второго лыжника больше скорости первого.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Ира прошла 320 м за 5 мин, а Петя – 225 м за 3 мин. У кого из ребят скорость больше и на сколько?","solution":"Ира 320 : 5 = 64 (м/мин), а Петя – 225 : 3 = 75 (м/мин). 75 - 64 = 11 (м/мин) Ответ: у Пети скорость больше на 11 м/мин."},{"letter":"б","condition":"Орёл за 9 с пролетел 270 м, а сокол за это же время пролетел 189 м. На сколько метров в секунду скорость сокола меньше скорости орла?","solution":"Орёл - 270 : 9 = 30 (м/с), а сокол – 189 : 9 = 21 (м/с). 30 - 21 = 9 (м/с) Ответ: на 9 метров в секунду скорость сокола меньше скорости орла."},{"letter":"в","condition":"Первый лыжник за 3 ч пробежал 51 км, а второй лыжник за это же время пробежал на 6 км больше. На сколько километров в час скорость второго лыжника больше скорости первого?","solution":"Первый лыжник – 51 : 3 = 17 (км/ч), а второй лыжник – (51 + 6) : 3 = 57 : 3 = 19 (км/ч). 19 - 17 = 2 (км/ч) Ответ: на 2 километров в час скорость второго лыжника больше скорости первого."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-18/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica18-nomer1.jpg', 'peterson/3/part3/page18/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '77c5f52063898c7cdf6da524869a07d50250a4fabf9dace41b4fce1cdc82af34', '3,5,6,9,51,189,225,270,320', '["больше","меньше"]'::jsonb, 'а) ира прошла 320 м за 5 мин, а петя-225 м за 3 мин. у кого из ребят скорость больше и на сколько? б) орёл за 9 с пролетел 270 м, а сокол за это же время пролетел 189 м. на сколько метров в секунду скорость сокола меньше скорости орла? в) первый лыжник за 3 ч пробежал 51 км, а второй лыжник за это же время пробежал на 6 км больше. на сколько километров в час скорость второго лыжника больше скорости первого');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 18, '2', 1, 'а) От деревни до станции 4 км. Ваня идёт из деревни на станцию со скоростью 80 м/мин. Какое расстояние ему останется пройти через полчаса после выхода? Сколько времени ему потребуется, чтобы пройти оставшееся расстояние? б) Автомобиль за 6 ч проехал 480 км. Какое расстояние мог бы проехать автомобиль за это же время, если бы увеличил скорость на 12 км/ч?', '</p> \n<p class="text">а) От деревни до станции 4 км. Ваня идёт из деревни на станцию со скоростью 80 м/мин. Какое расстояние ему останется пройти через полчаса после выхода? Сколько времени ему потребуется, чтобы пройти оставшееся расстояние?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica18-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 18, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 18, номер 2, год 2022."/>\n</div>\n</div>\n\n<p class="text">б) Автомобиль за 6 ч проехал 480 км. Какое расстояние мог бы проехать автомобиль за это же время, если бы увеличил скорость на 12 км/ч?</p>', 'а) 4000 - 80 · 30 = 4000 - 2400 = 1600 (м) , 1600 : 80 = 20 (мин) Ответ: 1600 метров ему останется пройти через полчаса после выхода. 20 минут ему потребуется, чтобы пройти оставшееся расстояние. б) (480 : 6 + 12) · 6 = (80 + 12) · 6 = 92 · 6 = 552 (км) Ответ: 552 километра мог бы проехать автомобиль за это же время, если бы увеличил скорость на 12 км/ч.', '<p>\nа) 4000 - 80 · 30 = 4000 - 2400 = 1600 (м) , 1600 : 80 = 20 (мин)<br/>\n<b>Ответ:</b> 1600 метров ему останется пройти через полчаса после выхода. 20 минут ему потребуется, чтобы пройти оставшееся расстояние.<br/><br/>\nб) (480 : 6 + 12) · 6 = (80 + 12) · 6 = 92 · 6 = 552 (км)<br/>\n<b>Ответ:</b> 552 километра мог бы проехать автомобиль за это же время, если бы увеличил скорость на 12 км/ч.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"От деревни до станции 4 км. Ваня идёт из деревни на станцию со скоростью 80 м/мин. Какое расстояние ему останется пройти через полчаса после выхода? Сколько времени ему потребуется, чтобы пройти оставшееся расстояние?","solution":"4000 - 80 · 30 = 4000 - 2400 = 1600 (м) , 1600 : 80 = 20 (мин) Ответ: 1600 метров ему останется пройти через полчаса после выхода. 20 минут ему потребуется, чтобы пройти оставшееся расстояние."},{"letter":"б","condition":"Автомобиль за 6 ч проехал 480 км. Какое расстояние мог бы проехать автомобиль за это же время, если бы увеличил скорость на 12 км/ч?","solution":"(480 : 6 + 12) · 6 = (80 + 12) · 6 = 92 · 6 = 552 (км) Ответ: 552 километра мог бы проехать автомобиль за это же время, если бы увеличил скорость на 12 км/ч."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-18/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica18-nomer2.jpg', 'peterson/3/part3/page18/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2cdaaac9c7727e40ec1074a650e075980466b3b2201b6ddf3e36a39913ac6c39', '4,6,12,80,480', NULL, 'а) от деревни до станции 4 км. ваня идёт из деревни на станцию со скоростью 80 м/мин. какое расстояние ему останется пройти через полчаса после выхода? сколько времени ему потребуется, чтобы пройти оставшееся расстояние? б) автомобиль за 6 ч проехал 480 км. какое расстояние мог бы проехать автомобиль за это же время, если бы увеличил скорость на 12 км/ч');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 18, '3', 2, 'БЛИЦтурнир а) Таня шла сначала по шоссе a км, а потом по просёлку b км. С какой скоростью шла Таня, если весь путь занял t часов? б) Костя шёл лесом a км, а полем на b км больше. Весь путь занял t часов. С какой скоростью шёл Костя? в) Расстояние от села Горшково до деревни Светлая a км, а от деревни Светлая до города в b раз меньше. Грузовик проехал от Горшково до города через деревню Светлая со скоростью v км/ч. Сколько времени ехал грузовик?', '</p> \n<p class="text">БЛИЦтурнир<br/>\nа) Таня шла сначала по шоссе a км, а потом по просёлку  b км. С какой скоростью шла Таня, если весь путь занял t часов?<br/>\nб) Костя шёл лесом a км, а полем на b км больше. Весь путь занял t часов. С какой скоростью шёл Костя?<br/>\nв) Расстояние от села Горшково до деревни Светлая a км, а от деревни Светлая до города в b раз меньше. Грузовик проехал от Горшково до города через деревню Светлая со скоростью v км/ч. Сколько времени ехал грузовик?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica18-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 18, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 18, номер 3, год 2022."/>\n</div>\n</div>', 'а) (a + b) : t (км/ч). Ответ: с (a + b) : t км/ч шла Таня, если весь путь занял t часов. б) (а + а + b) : t = (2а + b) : t (км/ч). Ответ: с (2а + b) : t км/ч шёл Костя. в) (а + а : b) : v (ч). Ответ: (а + а : b) : v часов ехал грузовик.', '<p>\nа) (a + b) : t (км/ч). <br/>\n<b>Ответ:</b> с (a + b) : t км/ч шла Таня, если весь путь занял t часов.<br/><br/>\nб) (а + а + b) : t = (2а + b) : t (км/ч). <br/>\n<b>Ответ:</b> с (2а + b) : t км/ч шёл Костя.<br/><br/>\nв) (а + а : b) : v (ч). <br/>\n<b>Ответ:</b> (а + а : b) : v часов ехал грузовик.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-18/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica18-nomer3.jpg', 'peterson/3/part3/page18/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '853f2e62bd066503911244bd07e342f53dd11ec7869cc05c08850d277272f5c0', NULL, '["больше","меньше","раз"]'::jsonb, 'блицтурнир а) таня шла сначала по шоссе a км, а потом по просёлку b км. с какой скоростью шла таня, если весь путь занял t часов? б) костя шёл лесом a км, а полем на b км больше. весь путь занял t часов. с какой скоростью шёл костя? в) расстояние от села горшково до деревни светлая a км, а от деревни светлая до города в b раз меньше. грузовик проехал от горшково до города через деревню светлая со скоростью v км/ч. сколько времени ехал грузовик');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 19, '4', 0, 'Составь выражение и найди его значение при данных значениях букв: а) Лодка проплывает a км вниз по реке со скоростью b км/ч, а возвращается со скоростью c км/ч. Какое время затратит лодка на весь путь туда и обратно? (a = 30, b = 10, c = 6) б) Валя прошла за k часов x км, а Серёжа за то же время прошёл y км. На сколько скорость Серёжи больше скорости Вали? (x = 12, y = 15, k = 3) в) Машина проехала за n часов d км. Какое расстояние она проедет за m часов, если будет ехать с той же скоростью? (d = 240, n = 4, m = 7)', '</p> \n<p class="text">Составь выражение и найди его значение при данных значениях букв:<br/><br/>\nа) Лодка проплывает a км вниз по реке со скоростью b км/ч, а возвращается со скоростью c км/ч. Какое время затратит лодка на весь путь туда и обратно? (a = 30, b = 10, c = 6)<br/>\nб) Валя прошла за k часов x км, а Серёжа за то же время прошёл y км. На сколько скорость Серёжи больше скорости Вали? (x = 12, y = 15, k = 3)<br/>\nв) Машина проехала за n часов d км. Какое расстояние она проедет за m часов, если будет ехать с той же скоростью? (d = 240, n = 4, m = 7)\n</p>', 'а) а : b + a : c = 30 : 10 + 30 : 6 = 3 + 5 = 8 (ч) Ответ: 8 часов затратит лодка на весь путь туда и обратно. б) у : k - x : k = 15 : 3 - 12 : 3 = 5 - 4 = 1 (км/ч) Ответ: на 1 км/ч скорость Серёжи больше скорости Вали. в) d : n · m = 240 : 4 · 7 = 60 · 7 = 420 (км) Ответ: 420 километра проедет она за m часов, если будет ехать с той же скоростью.', '<p>\nа) а : b + a : c = 30 : 10 + 30 : 6 = 3 + 5 = 8 (ч)<br/>\n<b>Ответ:</b> 8 часов затратит лодка на весь путь туда и обратно. <br/><br/>\nб) у : k - x : k = 15 : 3 - 12 : 3 = 5 - 4 = 1 (км/ч)<br/>\n<b>Ответ:</b> на 1 км/ч скорость Серёжи больше скорости Вали.<br/><br/>\nв) d : n · m = 240 : 4 · 7 = 60 · 7 = 420 (км)<br/>\n<b>Ответ:</b> 420 километра проедет она за m часов, если будет ехать с той же скоростью.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-19/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '05eebfa08a358f2b137a2f755b4e3f626982c67d7836317b1f192141e8dd0641', '3,4,6,7,10,12,15,30,240', '["найди","больше"]'::jsonb, 'составь выражение и найди его значение при данных значениях букв:а) лодка проплывает a км вниз по реке со скоростью b км/ч, а возвращается со скоростью c км/ч. какое время затратит лодка на весь путь туда и обратно? (a=30, b=10, c=6) б) валя прошла за k часов x км, а серёжа за то же время прошёл y км. на сколько скорость серёжи больше скорости вали? (x=12, y=15, k=3) в) машина проехала за n часов d км. какое расстояние она проедет за m часов, если будет ехать с той же скоростью? (d=240, n=4, m=7)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 19, '5', 1, 'Пусть a – длина прямоугольника, а b – его ширина. Объясни смысл выражений: a + b      a · 2 + b · 2      a · b a - b      (a + b) · 2        a : b', '</p> \n<p class="text">Пусть a – длина прямоугольника, а b – его ширина. Объясни смысл выражений:</p> \n\n<p class="description-text"> \na + b      a · 2 + b · 2      a · b<br/>\na - b      (a + b) · 2        a : b\n\n</p>\n\n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica19-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 19, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 19, номер 5, год 2022."/>\n</div>\n</div>', 'a + b – сумма ширины и длины прямоугольника a · 2 + b · 2 – сумма удвоенной ширины и удвоенной длины прямоугольника a · b – произведение длины и ширины прямоугольника a - b – разность длины и ширины прямоугольника (a + b) · 2 – удвоенная сумма длины и ширины прямоугольника a : b – частное длины и ширины прямоугольника', '<p>\na + b – сумма ширины и длины прямоугольника	<br/>\na · 2 + b · 2 – сумма удвоенной ширины и удвоенной длины прямоугольника<br/> 	\na · b – произведение длины и ширины прямоугольника<br/>\na - b – разность длины и ширины прямоугольника<br/>\n(a + b) · 2 – удвоенная сумма длины и ширины прямоугольника<br/>  		\na : b – частное длины и ширины прямоугольника\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-19/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica19-nomer5.jpg', 'peterson/3/part3/page19/task5_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd254f102308e35223d99e08153c313a7d6bfed6164d3e000a985c82f63d04f64', '2', NULL, 'пусть a-длина прямоугольника, а b-его ширина. объясни смысл выражений:a+b      a*2+b*2      a*b a-b      (a+b)*2        a:b');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 19, '6', 2, 'Найди площадь закрашенных фигур:', '</p> \n<p class="text">Найди площадь закрашенных фигур:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica19-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 19, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 19, номер 6, год 2022."/>\n</div>\n</div>', 'а) 12 · 6 + (6 - 4) · 4 = 72 + 8 = 80 (м 2 ) б) 80 · 96 - 40 · 28 = 7680 - 1120 = 6560 (см 2 )', '<p>\nа) 12 · 6 + (6 - 4) · 4 = 72 + 8 = 80 (м<sup>2</sup>)<br/>\nб) 80 · 96 - 40 · 28 = 7680 - 1120 = 6560 (см<sup>2</sup>)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-19/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica19-nomer6.jpg', 'peterson/3/part3/page19/task6_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '23cb115a3a23caabe7fe6eb77e104b4ff822f48872e53ca6e0020be2bb03d31a', NULL, '["найди","площадь"]'::jsonb, 'найди площадь закрашенных фигур');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 19, '7', 3, 'Реши уравнения с комментированием и проверкой: а) (150 : x + 6) : 7 = 8 б) 800 - (y · 8 - 20) = 100', '</p> \n<p class="text">Реши уравнения с комментированием и проверкой:</p> \n\n<p class="description-text"> \nа) (150 : x + 6) : 7 = 8<br/>        б) 800 - (y · 8 - 20) = 100\n</p>', 'а) (150 : x + 6) : 7 = 8 Чтобы найти делимое (150 : x + 6) надо делитель умножить на частное 150 : x + 6 = 8 · 7 150 : x + 6 = 56 Чтобы найти слагаемое 150 : х надо из суммы вычесть известное слагаемое 150 : х = 56 - 6 150 : х = 50 Чтобы найти делитель х надо делимое разделить на частное х = 150 : 50 х = 3 Проверка: (150 : 3 + 6) : 7 = 8 б) 800 - (y · 8 - 20) = 100 Чтобы найти вычитаемое (y · 8 - 20) надо из уменьшаемого вычесть разность y · 8 - 20 = 800 - 100 y · 8 - 20 = 700 Чтобы найти уменьшаемое y · 8 надо сложить вычитаемое и разность y · 8 = 20 + 700 y · 8 = 720 Чтобы найти множитель у надо произведение разделить на известный множитель y = 720 : 8 у = 90 Проверка: 800 - (90 · 8 - 20) = 100', '<p>\nа) (150 : x + 6) : 7 = 8 <br/>     \nЧтобы найти делимое (150 : x + 6) надо делитель умножить на частное<br/>\n150 : x + 6 = 8 · 7<br/>\n150 : x + 6 = 56<br/>\nЧтобы найти слагаемое 150 : х надо из суммы вычесть известное слагаемое<br/>\n150 : х = 56 - 6<br/>\n150 : х = 50<br/>\nЧтобы найти делитель х надо делимое разделить на частное<br/>\nх = 150 : 50<br/>\nх = 3<br/>\n<b>Проверка:</b> (150 : 3 + 6) : 7 = 8<br/><br/>\n\nб) 800 - (y · 8 - 20) = 100<br/>\nЧтобы найти вычитаемое (y · 8 - 20) надо из уменьшаемого вычесть разность<br/>\ny · 8 - 20 = 800 - 100<br/>\ny · 8 - 20 = 700<br/>\nЧтобы найти уменьшаемое y · 8 надо сложить вычитаемое и разность<br/>\ny · 8 = 20 + 700<br/>\ny · 8 = 720<br/>\nЧтобы найти множитель у надо произведение разделить на известный множитель<br/>\ny = 720 : 8<br/>\nу = 90<br/>\n<b>Проверка:</b> 800 - (90 · 8 - 20) = 100 \n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-19/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1975cbea4d14c04fca518f9a08175a3b6f0fce1709d67cf3cf5b612a7e0fba33', '6,7,8,20,100,150,800', '["реши"]'::jsonb, 'реши уравнения с комментированием и проверкой:а) (150:x+6):7=8 б) 800-(y*8-20)=100');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 19, '8', 4, 'Составь программу действий и вычисли: а) 0 · 19 + (45 : 1 - 0) · 1 - 18 · (12 : 12) б) 1 · 0 + (3 · 8 - 6 · 4) · 5 - 0 : (945 - 732)', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) 0 · 19 + (45 : 1 - 0) · 1 - 18 · (12 : 12)<br/>        \nб) 1 · 0 + (3 · 8 - 6 · 4) · 5 - 0 : (945 - 732)\n</p>', 'а) 0 · 19 + (45 : 1 - 0) · 1 – 18 · (12 : 12) = 27 45 : 1 = 45 45 - 0 = 45 45 · 1 = 45 12 : 12 = 1 18 · 1 = 18 0 · 19 = 0 0 + 45 = 45 45 - 18 = 27 б) 1 · 0 + (3 · 8 - 6 · 4) · 5 - 0 : (945 - 732) = 0 3 · 8 = 24 6 · 4 = 24 24 - 24 = 0 0 · 5 = 0 945 - 732 = 213 0 : 213 = 0 1 · 0 = 0 0 + 0 = 0 0 - 0 = 0', '<p>\nа) 0 · 19 + (45 : 1 - 0) · 1 – 18 · (12 : 12) = 27<br/>\n45 : 1 = 45<br/>\n45 - 0 = 45<br/>\n45 · 1 = 45<br/>\n12 : 12 = 1<br/>\n18 · 1 = 18<br/>\n0 · 19 = 0<br/>\n0 + 45 = 45<br/>\n45 - 18 = 27<br/><br/>\nб) 1 · 0 + (3 · 8 - 6 · 4) · 5 - 0 : (945 - 732) = 0<br/>\n3 · 8 = 24<br/>\n6 · 4 = 24<br/>\n24 - 24 = 0<br/>\n0 · 5 = 0<br/>\n945 - 732 = 213<br/>\n0 : 213 =  0<br/>\n1 · 0 = 0<br/>\n0 + 0 = 0<br/>\n0 - 0 = 0\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-19/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1f140447aec22a07b0253a17193f9b909750ffedbb8b49a0bfe163497945555c', '0,1,3,4,5,6,8,12,18,19', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) 0*19+(45:1-0)*1-18*(12:12) б) 1*0+(3*8-6*4)*5-0:(945-732)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 19, '9', 5, 'Запиши множество делителей и множество кратных числа 16.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 16.</p>', '1, 2, 4, 8, 16 16, 32, 48, 64, 80', '<p>\n1, 2, 4, 8, 16<br/>\n16, 32, 48, 64, 80\n</p>\n\n\n<p>\n</p><div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica19-nomer10-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 19, номер 10-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 19, номер 10-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-19/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica19-nomer10-1.jpg', 'peterson/3/part3/page19/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '96800d5e28a75ff8a38dfb44736afa7531d2544a544b9f1427144dec24fa13c5', '16', NULL, 'запиши множество делителей и множество кратных числа 16');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 20, '1', 0, 'а) После вспышки молнии Марина услышала гром через 5 с. На каком расстоянии от неё ударила молния? (Скорость распространения звука в воздухе равна 330 м/с.) б) Скорость распространения света 300000 км/с. На Солнце произошла вспышка. Через какое время её увидят на Земле, если расстояние от Земли до Солнца равно 150000000 км?', '</p> \n<p class="text">а) После вспышки молнии Марина услышала гром через 5 с. На каком расстоянии от неё ударила молния? (Скорость распространения звука в воздухе равна 330 м/с.)<br/>\nб) Скорость распространения света 300000 км/с. На Солнце произошла вспышка. Через какое время её увидят на Земле, если расстояние от Земли до Солнца равно 150000000 км?\n</p>', 'а) 330 · 5 = 1650 (м) Ответ: на расстоянии 1650 метров от неё ударила молния. б) 150000000 : 300000 = 500 (с) Ответ: через 500 секунд её увидят на Земле, если расстояние от Земли до Солнца равно 150000000 км.', '<p>\nа) 330 · 5 = 1650 (м)<br/>\n<b>Ответ:</b> на расстоянии 1650 метров от неё ударила молния.<br/><br/>\nб) 150000000 : 300000 = 500 (с)<br/>\n<b>Ответ:</b> через 500 секунд её увидят на Земле, если расстояние от Земли до Солнца равно 150000000 км.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"После вспышки молнии Марина услышала гром через 5 с. На каком расстоянии от неё ударила молния? (Скорость распространения звука в воздухе равна 330 м/с.)","solution":"330 · 5 = 1650 (м) Ответ: на расстоянии 1650 метров от неё ударила молния."},{"letter":"б","condition":"Скорость распространения света 300000 км/с. На Солнце произошла вспышка. Через какое время её увидят на Земле, если расстояние от Земли до Солнца равно 150000000 км?","solution":"150000000 : 300000 = 500 (с) Ответ: через 500 секунд её увидят на Земле, если расстояние от Земли до Солнца равно 150000000 км."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-20/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8414b0c18ca24fab763d37718673d2541e933dec178cd3f8bb44d7f535c975b1', '5,330,300000,150000000', '["равно"]'::jsonb, 'а) после вспышки молнии марина услышала гром через 5 с. на каком расстоянии от неё ударила молния? (скорость распространения звука в воздухе равна 330 м/с.) б) скорость распространения света 300000 км/с. на солнце произошла вспышка. через какое время её увидят на земле, если расстояние от земли до солнца равно 150000000 км');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 20, '2', 1, 'а) Грузовик проехал расстояние из города A в город B со скоростью 45 км/ч за 4 часа. Обратно из B в A он возвращался по той же дороге со скоростью на 15 км/ч больше. На сколько меньше времени затратил грузовик на обратный путь? б) Автобус проехал 432 км за 8 часов. На сколько километров в час он должен был увеличить скорость, чтобы проехать это расстояние на 2 часа быстрее?', '</p> \n<p class="text">а) Грузовик проехал расстояние из города A в город B со скоростью 45 км/ч за 4 часа. Обратно из B в A он возвращался по той же дороге со скоростью на 15 км/ч больше. На сколько меньше времени затратил грузовик на обратный путь?</p> \n\n<div class="description-text"> \n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica20-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 20, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 20, номер 2, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Автобус проехал 432 км за 8 часов. На сколько километров в час он должен был увеличить скорость, чтобы проехать это расстояние на 2 часа быстрее?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica20-nomer2-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 20, номер 2-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 20, номер 2-1, год 2022."/>\n</div>\n</div>', 'а) 45 · 4 = 180 (км) 180 : (45 + 15) = 180 : 60 = 3 (ч) 4 - 3 = 1 (ч) Ответ: на 1 час меньше времени затратил грузовик на обратный путь. б) 432 : 8 = 54 (км/ч) 8 - 2 = 6 (ч) 432 : 6 = 72 (км/ч) 72 - 54 = 18 (км/ч) Ответ: на 18 километров в час он должен был увеличить скорость, чтобы проехать это расстояние на 2 часа быстрее.', '<p>\nа) 45 · 4 = 180 (км) <br/>\n180 : (45 + 15) = 180 : 60 = 3 (ч)<br/>\n4 - 3 = 1 (ч)<br/>\n<b>Ответ:</b> на 1 час меньше времени затратил грузовик на обратный путь.<br/><br/>\nб) 432 : 8 = 54 (км/ч)<br/>\n8 - 2 = 6 (ч)<br/>\n432 : 6 = 72 (км/ч)<br/>\n72 - 54 = 18 (км/ч)<br/>\n<b>Ответ:</b> на 18 километров в час он должен был увеличить скорость, чтобы проехать это расстояние на 2 часа быстрее.\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Грузовик проехал расстояние из города A в город B со скоростью 45 км/ч за 4 часа. Обратно из B в A он возвращался по той же дороге со скоростью на 15 км/ч больше. На сколько меньше времени затратил грузовик на обратный путь?","solution":"45 · 4 = 180 (км) 180 : (45 + 15) = 180 : 60 = 3 (ч) 4 - 3 = 1 (ч) Ответ: на 1 час меньше времени затратил грузовик на обратный путь."},{"letter":"б","condition":"Автобус проехал 432 км за 8 часов. На сколько километров в час он должен был увеличить скорость, чтобы проехать это расстояние на 2 часа быстрее?","solution":"432 : 8 = 54 (км/ч) 8 - 2 = 6 (ч) 432 : 6 = 72 (км/ч) 72 - 54 = 18 (км/ч) Ответ: на 18 километров в час он должен был увеличить скорость, чтобы проехать это расстояние на 2 часа быстрее."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-20/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica20-nomer2.jpg', 'peterson/3/part3/page20/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica20-nomer2-1.jpg', 'peterson/3/part3/page20/task2_condition_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '66ebd5c1f82ca9e57195fddafc727497638561fa86ed9e744aa521de63518315', '2,4,8,15,45,432', '["больше","меньше"]'::jsonb, 'а) грузовик проехал расстояние из города a в город b со скоростью 45 км/ч за 4 часа. обратно из b в a он возвращался по той же дороге со скоростью на 15 км/ч больше. на сколько меньше времени затратил грузовик на обратный путь? б) автобус проехал 432 км за 8 часов. на сколько километров в час он должен был увеличить скорость, чтобы проехать это расстояние на 2 часа быстрее');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 20, '3', 2, 'Вычисли устно:', '</p> \n<p class="text">Вычисли устно:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica20-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 20, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 20, номер 3, год 2022."/>\n</div>\n</div>', 'S = 5 · 5 = 25 м 2           Р = 28 см Р = 4 · 5 = 20 м            S = 7 · 7 = 49 см 2 S = 12 дм 2                     P = 18 см P = 2(6 + 2) = 16 дм    S = 4 · 5 = 20', '<p>\nS = 5 · 5 = 25 м<sup>2</sup>          Р = 28 см<br/> 	 \nР = 4 · 5 = 20 м            S = 7 · 7 = 49 см<sup>2</sup><br/>	\n\nS = 12 дм<sup>2</sup>                    P = 18 см<br/>\nP = 2(6 + 2) = 16 дм    S = 4 · 5 = 20 \n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-20/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica20-nomer3.jpg', 'peterson/3/part3/page20/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '055b590c9ac65960536b9b017c90a950022df9ee500bfd3479e5868eea7f6d05', NULL, '["вычисли"]'::jsonb, 'вычисли устно');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 20, '4', 3, 'Автомобиль должен за 7 часов проехать 630 км. Первые 2 ч он ехал со скоростью 70 км/ч, а в следующие 3 ч увеличил скорость на 20 км/ч. С какой скоростью автомобиль должен ехать оставшийся путь, чтобы прибыть в пункт назначения вовремя?', '</p> \n<p class="text">Автомобиль должен за 7 часов проехать 630 км. Первые 2 ч он ехал со скоростью 70 км/ч, а в следующие 3 ч увеличил скорость на 20 км/ч. С какой скоростью автомобиль должен ехать оставшийся путь, чтобы прибыть в пункт назначения вовремя?</p>', '7 - 2 - 3 = 2 (ч) – время на оставшийся путь 630 - 70 · 2 - (70 + 20) · 3 = 630 - 140 - 90 · 3 = 490 - 270 = 220 (км) – оставшийся путь 220 : 2 = 110 (км/ч) Ответ: 110 километров в час должен автомобиль ехать оставшийся путь, чтобы прибыть в пункт назначения вовремя.', '<p>\n7 - 2 - 3 = 2 (ч) – время на оставшийся путь<br/>\n630 - 70 · 2 - (70 + 20) · 3 = 630 - 140 - 90 · 3 = 490 - 270 = 220 (км) – оставшийся путь<br/>\n220 : 2 = 110 (км/ч)<br/>\n<b>Ответ:</b> 110 километров в час должен автомобиль ехать оставшийся путь, чтобы прибыть в пункт назначения вовремя.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-20/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '139c28a094032677e2107106a387797941068440a4145b53661c53e5be918b5c', '2,3,7,20,70,630', NULL, 'автомобиль должен за 7 часов проехать 630 км. первые 2 ч он ехал со скоростью 70 км/ч, а в следующие 3 ч увеличил скорость на 20 км/ч. с какой скоростью автомобиль должен ехать оставшийся путь, чтобы прибыть в пункт назначения вовремя');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 21, '5', 0, 'Иван Иванович отправился из дома на озеро Медвежье ловить рыбу. Три часа он ехал на поезде со скоростью 75 км/ч, а потом 2 часа шёл по лесу со скоростью 4 км/ч. Какой путь проделал Иван Иванович от дома до озера?', '</p> \n<p class="text">Иван Иванович отправился из дома на озеро Медвежье ловить рыбу. Три часа он ехал на поезде со скоростью 75 км/ч, а потом 2 часа шёл по лесу со скоростью 4 км/ч. Какой путь проделал Иван Иванович от дома до озера?</p>', '75 · 3 + 4 · 2 = 225 + 8 = 233 (км) Ответ: 223 километра проделал Иван Иванович от дома до озера.', '<p>\n75 · 3 + 4 · 2 = 225 + 8 = 233 (км)<br/>\n<b>Ответ:</b> 223 километра проделал Иван Иванович от дома до озера.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-21/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6960956eea4f8df5095fb3daf7aba348bc22ce9ac65d994fb15cb38e567377bd', '2,4,75', NULL, 'иван иванович отправился из дома на озеро медвежье ловить рыбу. три часа он ехал на поезде со скоростью 75 км/ч, а потом 2 часа шёл по лесу со скоростью 4 км/ч. какой путь проделал иван иванович от дома до озера');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 21, '6', 1, 'Составь выражение и найди его значение, если a = 90: «Велосипедист проехал расстояние, равное a км, за 5 ч, а автобус – за 2 ч. На сколько километров в час скорость автобуса больше скорости велосипедиста?»', '</p> \n<p class="text">Составь выражение и найди его значение, если a = 90:<br/>\n«Велосипедист проехал расстояние, равное a км, за 5 ч, а автобус – за 2 ч. На сколько километров в час скорость автобуса больше скорости велосипедиста?»\n</p>', 'а : 2 - а : 5 = 90 : 2 - 90 : 5 = 45 - 19 = 26 (км/ч) Ответ: на 26 километров час скорость автобуса больше скорости велосипедиста.', '<p>\nа : 2 - а : 5 = 90 : 2 - 90 : 5 = 45 - 19 = 26 (км/ч)<br/>\n<b>Ответ:</b> на 26 километров час скорость автобуса больше скорости велосипедиста.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-21/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '996e7bb0e37a51f753daa56f514836550352b88e2760968364392f2c5dc59329', '2,5,90', '["найди","больше","равно"]'::jsonb, 'составь выражение и найди его значение, если a=90:"велосипедист проехал расстояние, равное a км, за 5 ч, а автобус-за 2 ч. на сколько километров в час скорость автобуса больше скорости велосипедиста?"');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 21, '7', 2, 'Прочитай выражения и найди их значения: а) 800 · n, если n = 70540 б) 278100 : c, если c = 90', '</p> \n<p class="text">Прочитай выражения и найди их значения:</p> \n\n<p class="description-text"> \nа) 800 · n, если n = 70540<br/> 	\nб) 278100 : c, если c  = 90\n</p>', 'а) 800 · n = 800 · 70540 = 56432000, если n = 70540 восемьсот умножить на n б) 278100 : c = 278100 : 90 = 390, если c = 90 двести семьдесят восемь тысяч сто разделить на с', '<p>\nа) 800 · n = 800 · 70540 = 56432000, если n = 70540  <br/> 	\nвосемьсот умножить на n<br/> <br/> \nб) 278100 : c = 278100 : 90 = 390, если c = 90<br/> \nдвести семьдесят восемь тысяч сто разделить на с\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-21/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e9125c37dfec6e3d85d89c35b9f0c2190f4429bd872f2fe62323b91d2ca9b341', '90,800,70540,278100', '["найди"]'::jsonb, 'прочитай выражения и найди их значения:а) 800*n, если n=70540 б) 278100:c, если c=90');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 21, '8', 3, 'Реши уравнения с комментированием и сделай проверку: а) (200 + 20 · a) : 6 = 60 б) 320 : (b · 8 - 40) = 10', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (200 + 20 · a) : 6 = 60 <br/>  		\nб) 320 : (b · 8 - 40) = 10\n</p>', 'а) (200 + 20 · a) : 6 = 60 Чтобы найти делимое 200 + 20 · а надо делитель умножить на частное 200 + 20 · а = 6 · 60 200 + 20 · а = 360 Чтобы найти слагаемое 20 · а надо из суммы отнять известное слагаемое 20 · а = 360 – 200 20 · а = 160 Чтобы найти множитель а надо произведение разделить на известный множитель а = 160 : 20 а = 8 Проверка: (200 + 20 · 8) : 6 = 60 б) 320 : (b · 8 - 40) = 10 Чтобы найти делитель b · 8 – 40 надо делимое разделить на частное b · 8 - 40 = 320 : 10 b · 8 - 40 = 32 Чтобы найти уменьшаемое b · 8 надо вычитаемое сложить с разностью b · 8 = 40 + 32 b · 8 = 72 Чтобы найти множитель b надо произведение разделить на известный множитель b = 72 : 8 b = 9 Проверка: 320 : (9 · 8 - 40) = 10', '<p>\nа) (200 + 20 · a) : 6 = 60  <br/>	\nЧтобы найти делимое 200 + 20 · а надо делитель умножить на частное<br/>\n200 + 20 · а = 6 · 60<br/>\n200 + 20 · а = 360<br/>\nЧтобы найти слагаемое 20 · а надо из суммы отнять известное слагаемое<br/>\n20 · а = 360 – 200<br/>\n20 · а = 160<br/>\nЧтобы найти множитель а надо произведение разделить на известный множитель<br/>\nа = 160 : 20<br/>\nа = 8  	<br/>\nПроверка: (200 + 20 · 8) : 6 = 60<br/><br/>\nб) 320 : (b · 8 - 40) = 10<br/>\nЧтобы найти делитель b · 8 – 40 надо делимое разделить на частное<br/>\nb · 8 - 40 = 320 : 10<br/>\nb · 8 - 40 = 32<br/>\nЧтобы найти уменьшаемое b · 8 надо вычитаемое сложить с разностью<br/>\nb · 8 = 40 + 32<br/>\nb · 8 = 72<br/>\nЧтобы найти множитель b надо произведение разделить на известный множитель<br/>\nb = 72 : 8<br/>\nb = 9<br/>\nПроверка: 320 : (9 · 8 - 40) = 10  \n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-21/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c721f6184ff212676eeeb0531c4c9617109d65c0c5cf49672e4fc1c7400f0395', '6,8,10,20,40,60,200,320', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (200+20*a):6=60 б) 320:(b*8-40)=10');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 21, '9', 4, 'а) Туристы вышли из посёлка Дачное. В каком направлении и с какой скоростью они идут? Построй числовой луч и покажи на нём движение туристов. б) Пусть s км – путь, пройденный туристами, d км – расстояние от туристов до Грибцова, а D км – до Земляничной Поляны. Заполни таблицу. Запиши формулу зависимости каждой из величин s, d, D от времени движения t.', '</p> \n<p class="text">а) Туристы вышли из посёлка Дачное. В каком направлении и с какой скоростью они идут? Построй числовой луч и покажи на нём движение туристов.</p> \n\n<div class="description-text"> \n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica21-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 21, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 21, номер 9, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Пусть s км – путь, пройденный туристами, d км – расстояние от туристов до Грибцова, а D км – до Земляничной Поляны. Заполни таблицу. Запиши формулу зависимости каждой из величин s, d,  D от времени движения t.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica21-nomer9-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 21, номер 9-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 21, номер 9-1, год 2022."/>\n</div>\n</div>', 'а) Туристы вышли из посёлка Дачное в направлении Грибцово с 6 км/ч они идут. б)', '<p>\nа) Туристы вышли из посёлка Дачное в направлении Грибцово с 6 км/ч они идут.<br/>\nб) \n\n</p>\n\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica21-nomer9-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 21, номер 9-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 21, номер 9-2, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Туристы вышли из посёлка Дачное. В каком направлении и с какой скоростью они идут? Построй числовой луч и покажи на нём движение туристов.","solution":"Туристы вышли из посёлка Дачное в направлении Грибцово с 6 км/ч они идут."},{"letter":"б","condition":"Пусть s км – путь, пройденный туристами, d км – расстояние от туристов до Грибцова, а D км – до Земляничной Поляны. Заполни таблицу. Запиши формулу зависимости каждой из величин s, d, D от времени движения t.","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-21/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica21-nomer9.jpg', 'peterson/3/part3/page21/task9_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica21-nomer9-1.jpg', 'peterson/3/part3/page21/task9_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica21-nomer9-2.jpg', 'peterson/3/part3/page21/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '797c8b9ab2c38972f5351c0c69498a284ab51e8865314acd812509f4321d47a4', NULL, '["заполни","числовой луч"]'::jsonb, 'а) туристы вышли из посёлка дачное. в каком направлении и с какой скоростью они идут? построй числовой луч и покажи на нём движение туристов. б) пусть s км-путь, пройденный туристами, d км-расстояние от туристов до грибцова, а d км-до земляничной поляны. заполни таблицу. запиши формулу зависимости каждой из величин s, d, d от времени движения t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 21, '10', 5, 'Запиши множество делителей и множество кратных числа 17.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 17.</p>', 'Делители числа 17 = 1, 17. Кратные числа 17 = 17, 34, 51, 68, 85 99999999999999', '<p>\nДелители числа 17 = 1, 17.<br/>\nКратные числа 17 = 17, 34, 51, 68, 85\n</p>\n\n\n<p>\n99999999999999\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-21/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '90d3665f50cd14ee5b091c359094ee8efa2f2cd2cc74578796eb5b4d669b15be', '17', NULL, 'запиши множество делителей и множество кратных числа 17');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 22, '1', 0, 'Повтори таблицу мер длины, массы и времени. Вырази данные значения величин в указанных единицах измерения: а) 8 км 16 м = … м      б) 2 сут. 9 ч 25 мин = … мин 5 м 9 мм = … мм          3 ч 12 мин 46 с = … с 2 т 3 ц 6 кг = … кг       870 мин = … ч … мин 4 ц 7 кг 8 г = … г         3520 с = … мин … с', '</p> \n<p class="text">Повтори таблицу мер длины, массы и времени. Вырази данные значения величин в указанных единицах измерения:</p> \n\n<p class="description-text"> \nа)  8 км 16 м = … м      б) 2 сут. 9 ч 25 мин = … мин<br/>  \n5 м 9 мм = … мм          3 ч 12 мин 46 с = … с<br/>  \n2 т 3 ц 6 кг = … кг       870 мин = … ч … мин<br/>  \n4 ц 7 кг 8 г = … г         3520 с = … мин … с\n</p>', 'а) 8 км 16 м = 8016 м 5 м 9 мм = 5009 мм 2 т 3 ц 6 кг = 2036 кг 4 ц 7 кг 8 г = 407008 г б) 2 сут. 9 ч 25 мин = 3445 мин 3 ч 12 мин 46 с = 12566 с 870 мин = 14 ч 30 мин 3520 с = 58 мин 40 с', '<p>\nа) 8 км 16 м = 8016 м<br/>  	  \n5 м 9 мм = 5009 мм <br/>         \n2 т 3 ц 6 кг = 2036 кг<br/>      	 \n4 ц 7 кг 8 г = 407008 г<br/><br/>     \n\nб) 2 сут. 9 ч 25 мин = 3445 мин<br/>\n3 ч 12 мин 46 с = 12566 с<br/> \n870 мин = 14 ч 30 мин<br/> \n3520 с = 58 мин 40 с\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-22/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ed56d86380ad86b7363663610f47736f4fdcff8a327691e00a6f8363fa0e4210', '2,3,4,5,6,7,8,9,12,16', '["раз"]'::jsonb, 'повтори таблицу мер длины, массы и времени. вырази данные значения величин в указанных единицах измерения:а) 8 км 16 м=... м      б) 2 сут. 9 ч 25 мин=... мин 5 м 9 мм=... мм          3 ч 12 мин 46 с=... с 2 т 3 ц 6 кг=... кг       870 мин=... ч ... мин 4 ц 7 кг 8 г=... г         3520 с=... мин ... с');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 22, '2', 1, 'а) Лыжник прошёл 18 км за 2 часа. Какое расстояние он пройдёт за такое же время, если увеличит скорость на 3 км/ч? б) Моторная лодка прошла по течению реки 5 ч со скоростью 24 км/ч. На обратный путь она затратила на 1 час больше времени. Чему равна скорость моторной лодки против течения реки?', '</p> \n<p class="text">а) Лыжник прошёл 18 км за 2 часа. Какое расстояние он пройдёт за такое же время, если увеличит скорость на 3 км/ч?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica22-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 22, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 22, номер 2, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Моторная лодка прошла по течению реки 5 ч со скоростью 24 км/ч. На обратный путь она затратила на 1 час больше времени. Чему равна скорость моторной лодки против течения реки?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica22-nomer2-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 22, номер 2-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 22, номер 2-1, год 2022."/>\n</div>\n</div>', 'а) (18 : 2 + 3) · 2 = (9 + 3) · 2 = 12 · 2 = 24 (км) Ответ: 24 километра он пройдёт за такое же время, если увеличит скорость на 3 км/ч. б) 5 · 24 : (5 + 1) = 120 : 6 = 20 (км/ч) Ответ: 20 километров в час скорость моторной лодки против течения реки.', '<p>\nа) (18 : 2 + 3) · 2 = (9 + 3) · 2 = 12 · 2 = 24 (км)<br/>\n<b>Ответ:</b> 24 километра он пройдёт за такое же время, если увеличит скорость на 3 км/ч.<br/><br/>\nб) 5 · 24 : (5 + 1) = 120 : 6 = 20 (км/ч)<br/>\n<b>Ответ:</b> 20 километров в час скорость моторной лодки против течения реки.\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Лыжник прошёл 18 км за 2 часа. Какое расстояние он пройдёт за такое же время, если увеличит скорость на 3 км/ч?","solution":"(18 : 2 + 3) · 2 = (9 + 3) · 2 = 12 · 2 = 24 (км) Ответ: 24 километра он пройдёт за такое же время, если увеличит скорость на 3 км/ч."},{"letter":"б","condition":"Моторная лодка прошла по течению реки 5 ч со скоростью 24 км/ч. На обратный путь она затратила на 1 час больше времени. Чему равна скорость моторной лодки против течения реки?","solution":"5 · 24 : (5 + 1) = 120 : 6 = 20 (км/ч) Ответ: 20 километров в час скорость моторной лодки против течения реки."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-22/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica22-nomer2.jpg', 'peterson/3/part3/page22/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica22-nomer2-1.jpg', 'peterson/3/part3/page22/task2_condition_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e8be2b27355bd6dbdc31c1c1a75aab8d36df3fee3270083d15ab2722698e1e1a', '1,2,3,5,18,24', '["больше"]'::jsonb, 'а) лыжник прошёл 18 км за 2 часа. какое расстояние он пройдёт за такое же время, если увеличит скорость на 3 км/ч? б) моторная лодка прошла по течению реки 5 ч со скоростью 24 км/ч. на обратный путь она затратила на 1 час больше времени. чему равна скорость моторной лодки против течения реки');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 22, '3', 2, 'Вычисли устно:', '</p> \n<p class="text">Вычисли устно:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica22-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 22, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 22, номер 3, год 2022."/>\n</div>\n</div>', 'V = 9 · 3 · 2 = 9 · 6 = 54 (м 3 ) S = 12 · 5 - 6 · 2 = 60 – 12 = 48 (дм 2 )', '<p>\nV = 9 · 3 · 2 = 9 · 6 = 54 (м<sup>3</sup>)<br/>\nS = 12 · 5 - 6 · 2 = 60 – 12 = 48 (дм<sup>2</sup>)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-22/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica22-nomer3.jpg', 'peterson/3/part3/page22/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '055b590c9ac65960536b9b017c90a950022df9ee500bfd3479e5868eea7f6d05', NULL, '["вычисли"]'::jsonb, 'вычисли устно');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 22, '4', 3, 'Катер прошёл путь между двумя пристанями со скоростью 30 км/ч, а обратный путь – со скоростью на 10 км/ч большей. Расстояние между этими пристанями равно 240 км. Какое время затратил катер на путь туда и обратно?', '</p> \n<p class="text">Катер прошёл путь между двумя пристанями со скоростью 30 км/ч, а обратный путь – со скоростью на 10 км/ч большей. Расстояние между этими пристанями равно 240 км. Какое время затратил катер на путь туда и обратно?</p>', '240 : 30 + 240 : (30 + 10) = 8 + 6 = 14 (ч) Ответ: 14 часов затратил катер на путь туда и обратно.', '<p>\n240 : 30 + 240 : (30 + 10) = 8 + 6 = 14 (ч)<br/>\n<b>Ответ:</b> 14 часов затратил катер на путь туда и обратно.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-22/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c4c5e3d068bdbc7a2e079818331f3640e826f4d84a619f3c5307bad38ab865a9', '10,30,240', '["больше","равно"]'::jsonb, 'катер прошёл путь между двумя пристанями со скоростью 30 км/ч, а обратный путь-со скоростью на 10 км/ч большей. расстояние между этими пристанями равно 240 км. какое время затратил катер на путь туда и обратно');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 22, '5', 4, 'Прочитай выражения и найди их значения: а) 10000 - x : 70, если x = 644560 б) (y · 6004) : 500, если y = 4000', '</p> \n<p class="text">Прочитай выражения и найди их значения:</p> \n\n<p class="description-text"> \nа) 10000 - x : 70, если x = 644560<br/>  		\nб) (y · 6004) : 500, если y = 4000\n</p>', 'а) 10000 - x : 70 = 10000 - 644560 : 70 = 10000 - 9208 = 792, если x = 644560 разность десяти тысяч и частного х на семьдесят', '<p>\nа) 10000 - x : 70 = 10000 - 644560 : 70 = 10000 - 9208 = 792, если x = 644560<br/>  	\nразность десяти тысяч и частного х на семьдесят\n</p>\n\n<div class="img-wrapper-460">\n<img width="100" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica22-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 22, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 22, номер 5, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-22/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica22-nomer5.jpg', 'peterson/3/part3/page22/task5_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '04e0e029dd935a45828eb26f98312a78a24a1b516aad07be3d07518e849d0877', '70,500,4000,6004,10000,644560', '["найди"]'::jsonb, 'прочитай выражения и найди их значения:а) 10000-x:70, если x=644560 б) (y*6004):500, если y=4000');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 22, '6', 5, 'Реши уравнения с комментированием и сделай проверку: а) (n : 4 - 35) · 6 = 150 б) 90 · (m - 8) + 60 = 510', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (n : 4 - 35) · 6 = 150<br/>           \nб) 90 · (m - 8) + 60 = 510\n</p>', 'а) (n : 4 - 35) · 6 = 150 Чтобы найти множитель n : 4 – 35 надо произведение разделить на известный множитель n : 4 - 35 = 150 : 6 n : 4 - 35 = 25 Чтобы найти уменьшаемое n : 4 надо вычитаемое сложить с разностью n : 4 = 35 + 25 n : 4 = 60 Чтобы найти делимое n надо делитель умножить на частное n = 60 · 4 n = 15 Проверка: (15 : 4 - 35) · 6 = 150 б) 90 · (m - 8) + 60 = 510 Чтобы найти слагаемое 90 · (m - 8) надо из суммы вычесть звестное слагаемое 90 · (m - 8) = 510 - 60 90 · (m - 8) = 450 Чтобы найти множитель m - 8 надо произведение разделить на известный множитель m - 8 = 450 : 90 m - 8 = 5 Чтобы найти уменьшаемое m надо вычитаемое умножить на разность m = 8 · 5 m = 40 Проверка: 90 · (m - 8) = 450', '<p>\nа) (n : 4 - 35) · 6 = 150 <br/>           \nЧтобы найти множитель n : 4 – 35 надо произведение разделить на известный множитель<br/> \nn : 4 - 35 = 150 : 6<br/> \nn : 4 - 35 = 25<br/> \nЧтобы найти уменьшаемое n : 4 надо вычитаемое сложить с разностью<br/> \nn : 4 = 35 + 25<br/> \nn : 4 = 60<br/> \nЧтобы найти делимое n надо делитель умножить на частное<br/> \nn = 60 · 4<br/> \nn = 15<br/> \n<b>Проверка:</b> (15 : 4 - 35) · 6 = 150<br/> <br/> \nб) 90 · (m - 8) + 60 = 510<br/> \nЧтобы найти слагаемое 90 · (m - 8) надо из суммы вычесть звестное слагаемое<br/> \n90 · (m - 8) = 510 - 60 <br/> \n90 · (m - 8) = 450<br/> \nЧтобы найти множитель m - 8 надо произведение разделить на известный множитель<br/>  \nm - 8 = 450 : 90<br/> \nm - 8 = 5<br/> \nЧтобы найти уменьшаемое m надо вычитаемое умножить на разность<br/> \nm = 8 · 5 <br/> \nm = 40<br/> \n<b>Проверка:</b> 90 · (m - 8) = 450 \n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-22/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6db71e74124ced80a19bdfe8da687c8ece7a52230de7d055692c4a19f849709e', '4,6,8,35,60,90,150,510', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (n:4-35)*6=150 б) 90*(m-8)+60=510');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 23, '7', 0, 'Составь и реши уравнения: а) Миша задумал число, умножил его на 5 и полученное произведение вычел из 41. В результате у него получилось 16. Какое число задумал Миша? б) Галя задумала число, вычла его из 50, результат разделила на 7. У неё получилось 7. Какое число задумала Галя? в) Тимоша задумал число, затем разделил 54 на задуманное число, прибавил к результату 26 и полученную сумму разделил на 8. В ответе у него получилось 4. Какое число задумал Тимоша?', '</p> \n<p class="text">Составь и реши уравнения:</p> \n\n<p class="description-text"> \nа) Миша задумал число, умножил его на 5 и полученное произведение вычел из 41. В результате у него получилось 16. Какое число задумал Миша?<br/>\nб) Галя задумала число, вычла его из 50, результат разделила на 7. У неё получилось 7. Какое число задумала Галя?<br/>\nв) Тимоша задумал число, затем разделил 54 на задуманное число, прибавил к результату 26 и полученную сумму разделил на 8. В ответе у него получилось 4. Какое число задумал Тимоша?\n</p>', 'Та) Пусть Миша задумал число х. Миша умножил его на 5, то получил 5х. Полученное произведение Миша вычел из 41, то получил (41 - 5х) = 16 5x = 41 - 16 5x = 25 x = 25 : 5 x = 5 Ответ. Миша задумал число 5. б) x − задуманное число (50 - x) : 7 = 7 50 - x = 7 · 7 50 - x = 49 x = 50 − 49 x = 1 Ответ: 1 − задуманное число. в) (54 : х + 26) : 8 = 4 54 : х + 26 = 4 · 8 54 : х + 26 = 32 54 : х = 32 - 26 54 : х = 6 х = 9 Ответ: Тимоша задумал 4.', '<p>\nТа) Пусть Миша задумал число х. Миша умножил его на 5, то получил 5х. Полученное произведение Миша вычел из 41, то получил (41 - 5х) = 16<br/>\n5x = 41 - 16<br/>\n5x = 25<br/>\nx = 25 : 5<br/>\nx = 5 <br/>\n<b>Ответ.</b> Миша задумал число 5.<br/><br/>\nб) x − задуманное число<br/>\n(50 - x) : 7 = 7 <br/>\n50 - x = 7 · 7 <br/>\n50 - x = 49 <br/>\nx = 50 − 49 <br/>\nx = 1<br/>\n<b>Ответ:</b> 1 − задуманное число.<br/><br/>\nв) (54 : х + 26) : 8 = 4 <br/>\n54 : х + 26 = 4 · 8 <br/>\n54 : х + 26 = 32<br/> \n54 : х = 32 - 26<br/> \n54 : х = 6 <br/>\nх = 9<br/>\n<b>Ответ:</b> Тимоша задумал 4.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-23/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1f2ca67ba69b52dd7547bd0ea0e53c81016470a5043540c631874406dcacc0bd', '4,5,7,8,16,26,41,50,54', '["раздели","реши","произведение","раз"]'::jsonb, 'составь и реши уравнения:а) миша задумал число, умножил его на 5 и полученное произведение вычел из 41. в результате у него получилось 16. какое число задумал миша? б) галя задумала число, вычла его из 50, результат разделила на 7. у неё получилось 7. какое число задумала галя? в) тимоша задумал число, затем разделил 54 на задуманное число, прибавил к результату 26 и полученную сумму разделил на 8. в ответе у него получилось 4. какое число задумал тимоша');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 23, '8', 1, 'Составь программу действий и вычисли: а) (1800 : 2 : 30 + 18) : 6 + (70 · 7 - 140 : 2) : 60 б) (60 - 16 : 4) : 8 · 40 - (80 · 8 - 20 · 5) : 6', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) (1800 : 2 : 30 + 18) : 6 + (70 · 7 - 140 : 2) : 60<br/>\nб) (60 - 16 : 4) : 8 · 40 - (80 · 8 - 20 · 5) : 6\n</p>', 'а) (1800 : 2 : 30 + 18) : 6 + (70 · 7 - 140 : 2) : 60 = 15 1800 : 2 = 900 900 : 30 = 30 30 + 18 = 48 48 : 6 = 8 70 · 7 = 490 140 : 2 = 70 490 - 70 = 420 420 : 60 = 7 8 + 7 = 15 б) (60 - 16 : 4) : 8 · 40 - (80 · 8 - 20 · 5) : 6 = 190 16 : 4 = 4 60 - 4 = 56 56 : 8 = 7 7 · 40 = 280 80 · 8 = 640 20 · 5 = 100 640 - 100 = 540 540 : 6 = 90 280 - 90 = 190', '<p>\nа) (1800 : 2 : 30 + 18) : 6 + (70 · 7 - 140 : 2) : 60 = 15<br/>\n1800 : 2 = 900<br/>\n900 : 30 = 30<br/>\n30 + 18 = 48<br/>\n48 : 6 = 8<br/>\n70 · 7 = 490<br/>\n140 : 2 = 70<br/>\n490 - 70 = 420<br/>\n420 : 60 = 7<br/>\n8 + 7 = 15<br/><br/>\nб) (60 - 16 : 4) : 8 · 40 - (80 · 8 - 20 · 5) : 6 = 190<br/>\n16 : 4 = 4<br/>\n60 - 4 = 56<br/>\n56 : 8 = 7<br/>\n7 · 40 = 280<br/>\n80 · 8 = 640<br/>\n20 · 5 = 100<br/>\n640 - 100 = 540<br/>\n540 : 6 = 90<br/>\n280 - 90 = 190\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-23/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'cd3d2c11157223b6df07df6ebf86515e83affd23ebe720ec3a6a563c6024e2a5', '2,4,5,6,7,8,16,18,20,30', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) (1800:2:30+18):6+(70*7-140:2):60 б) (60-16:4):8*40-(80*8-20*5):6');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 23, '9', 2, 'БЛИЦтурнир а) Стрекоза пролетает a км за 2 ч. Какое расстояние она пролетит за 5 ч, если будет лететь с той же скоростью? б) Заяц пробежал b км за 3 ч, а волк пробежал то же расстояние за 4 ч. У кого из них скорость больше и на сколько? в) Крокодил Гена проехал 3 ч на поезде со скоростью n км/ч и 2 ч на автобусе со скоростью m км/ч. Сколько всего километров он проехал? г) Черепаха Тортила 5 ч ползла со скоростью c км/ч. Всего ей надо проползти d км. Какое расстояние ей ещё осталось проползти?', '</p> \n<p class="text">БЛИЦтурнир<br/>\nа) Стрекоза пролетает a км за 2 ч. Какое расстояние она пролетит за 5 ч, если будет лететь с той же скоростью?<br/>\nб) Заяц пробежал b км за 3 ч, а волк пробежал то же расстояние за 4 ч. У кого из них скорость больше и на сколько?<br/>\nв) Крокодил Гена проехал 3 ч на поезде со скоростью n км/ч и 2 ч на автобусе со скоростью m км/ч. Сколько всего километров он проехал?<br/>\nг) Черепаха Тортила 5 ч ползла со скоростью c км/ч. Всего ей надо проползти d км. Какое расстояние ей ещё осталось проползти?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica23-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 23, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 23, номер 9, год 2022."/>\n</div>\n</div>', 'а) a : 2 · 5 (км) Ответ: 5 километров она пролетит за 5 ч, если будет лететь с той же скоростью. б) b : 3 - b : 4 = b(4 - 3) : 12 = b : 12 (км/ч) Ответ: у зайца из них скорость больше на b : 12 км/ч. в) 3n + 2m (км) Ответ: 3n + 2m километров всего километров он проехал. г) d - 5c (км) Ответ: d - 5c километров ей ещё осталось проползти.', '<p>\nа) a : 2 · 5 (км)<br/>\n<b>Ответ:</b> 5 километров она пролетит за 5 ч, если будет лететь с той же скоростью.<br/><br/>\nб) b : 3 - b : 4 = b(4 - 3) : 12 = b : 12 (км/ч)<br/>\n<b>Ответ:</b> у зайца из них скорость больше на b : 12 км/ч. <br/><br/>\nв) 3n + 2m (км) <br/>\n<b>Ответ:</b> 3n + 2m километров всего километров он проехал.<br/><br/>\nг) d - 5c (км) <br/>\n<b>Ответ:</b> d - 5c километров ей ещё осталось проползти.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-23/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica23-nomer9.jpg', 'peterson/3/part3/page23/task9_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '34a046b273b74ba8a14ddf0cb85d405c20951b132c4768fcf1054a03c1dea98c', '2,3,4,5', '["больше"]'::jsonb, 'блицтурнир а) стрекоза пролетает a км за 2 ч. какое расстояние она пролетит за 5 ч, если будет лететь с той же скоростью? б) заяц пробежал b км за 3 ч, а волк пробежал то же расстояние за 4 ч. у кого из них скорость больше и на сколько? в) крокодил гена проехал 3 ч на поезде со скоростью n км/ч и 2 ч на автобусе со скоростью m км/ч. сколько всего километров он проехал? г) черепаха тортила 5 ч ползла со скоростью c км/ч. всего ей надо проползти d км. какое расстояние ей ещё осталось проползти');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 23, '10', 3, 'Запиши множество делителей и множество кратных числа 18.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 18.</p>', 'Множество делителей числа 1, 2, 3, 6, 9, 18 Множество кратных числа 18, 36, 54, 72, 90', '<p>\nМножество делителей числа 1, 2, 3, 6, 9, 18<br/>\nМножество кратных числа 18, 36, 54, 72, 90\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-23/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'fa09f6026dca587ac76aeb57f36c68df162f1c06202c6f21bc35439787278160', '18', NULL, 'запиши множество делителей и множество кратных числа 18');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 23, '11', 4, 'Найди пропущенные цифры. Проверь с помощью калькулятора.', '</p> \n<p class="text">\nНайди пропущенные цифры. Проверь с помощью калькулятора.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica23-nomer11.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 23, номер 11, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 23, номер 11, год 2022."/>\n</div>\n</div>', 'Илья проехал дальше всех. Костя проехал меньше всех.', '<p>\n</p><div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica23-nomer11-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 23, номер 11-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 23, номер 11-1, год 2022."/>\n\n\n<p>\nИлья проехал дальше всех. Костя проехал меньше всех.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-23/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica23-nomer11.jpg', 'peterson/3/part3/page23/task11_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica23-nomer11-1.jpg', 'peterson/3/part3/page23/task11_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'bd2ab484242a48c0cc57106fe83abde922c2af27e83f78d7655f55c6dcc67566', NULL, '["найди"]'::jsonb, 'найди пропущенные цифры. проверь с помощью калькулятора');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 24, '1', 0, 'а) Вертолёт пролетел 840 км за 4 часа, а автобус проехал расстояние в 2 раза меньшее, затратив на 2 часа больше. Во сколько раз скорость автобуса меньше скорости вертолёта? б) Лыжник пробежал 36 км за 2 ч, а пешеход прошёл половину этого расстояния за время в 3 раза большее. На сколько километров в час скорость пешехода меньше скорости лыжника?', '</p> \n<p class="text">а) Вертолёт пролетел 840 км за 4 часа, а автобус проехал расстояние в 2 раза меньшее, затратив на 2 часа больше. Во сколько раз скорость автобуса меньше скорости вертолёта?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica24-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 24, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 24, номер 1, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Лыжник пробежал 36 км за 2 ч, а пешеход прошёл половину этого расстояния за время в 3 раза большее. На сколько километров в час скорость пешехода меньше скорости лыжника?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica24-nomer1-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 24, номер 1-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 24, номер 1-1, год 2022."/>\n</div>\n</div>', 'а) (840 : 4) : (840 : 2 : (4 + 2)) = 210 : (210 : 6) = 210 : 35 = 6 (раз) Ответ: в 6 раз скорость автобуса меньше скорости вертолёта. б) 36 : 2 - 36 : 2 : 2 · 3 = 18 - 27 = 9 (км/ч) Ответ: на 9 километров в час скорость пешехода меньше скорости лыжника.', '<p>\nа) (840 : 4) : (840 : 2 : (4 + 2)) = 210 : (210 : 6) = 210 : 35 = 6 (раз)<br/>  \n<b>Ответ:</b> в 6 раз скорость автобуса меньше скорости вертолёта.<br/><br/>\nб) 36 : 2 - 36 : 2 : 2 · 3 = 18 - 27 = 9 (км/ч)<br/>\n<b>Ответ:</b> на 9 километров в час скорость пешехода меньше скорости лыжника.\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Вертолёт пролетел 840 км за 4 часа, а автобус проехал расстояние в 2 раза меньшее, затратив на 2 часа больше. Во сколько раз скорость автобуса меньше скорости вертолёта?","solution":"(840 : 4) : (840 : 2 : (4 + 2)) = 210 : (210 : 6) = 210 : 35 = 6 (раз) Ответ: в 6 раз скорость автобуса меньше скорости вертолёта."},{"letter":"б","condition":"Лыжник пробежал 36 км за 2 ч, а пешеход прошёл половину этого расстояния за время в 3 раза большее. На сколько километров в час скорость пешехода меньше скорости лыжника?","solution":"36 : 2 - 36 : 2 : 2 · 3 = 18 - 27 = 9 (км/ч) Ответ: на 9 километров в час скорость пешехода меньше скорости лыжника."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-24/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica24-nomer1.jpg', 'peterson/3/part3/page24/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica24-nomer1-1.jpg', 'peterson/3/part3/page24/task1_condition_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2c7af91896e919564e4b543c2af17a60800d472ea953249b4b2977ec0d07dd20', '2,3,4,36,840', '["больше","меньше","раз","раза"]'::jsonb, 'а) вертолёт пролетел 840 км за 4 часа, а автобус проехал расстояние в 2 раза меньшее, затратив на 2 часа больше. во сколько раз скорость автобуса меньше скорости вертолёта? б) лыжник пробежал 36 км за 2 ч, а пешеход прошёл половину этого расстояния за время в 3 раза большее. на сколько километров в час скорость пешехода меньше скорости лыжника');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 24, '2', 1, 'БЛИЦтурнир а) После того как поезд проехал 4 часа со скоростью n км/ч, ему ещё осталось проехать b км. Чему равен весь путь поезда? б) Спортсмен бежал 2 часа со скоростью v км/ч. Длина всей дистанции равна m км. Сколько километров ему ещё осталось пробежать? в) Самолёт пролетел s км за 3 часа, а в обратную сторону – за 2 часа. На сколько километров в час больше была его скорость на обратном пути?', '</p> \n<p class="text">БЛИЦтурнир<br/>\nа) После того как поезд проехал 4 часа со скоростью n км/ч, ему ещё осталось проехать b км. Чему равен весь путь поезда?<br/>\nб) Спортсмен бежал 2 часа со скоростью v км/ч. Длина всей дистанции равна m км. Сколько километров ему ещё осталось пробежать?<br/>\nв) Самолёт пролетел s км за 3 часа, а в обратную сторону – за 2 часа. На сколько километров в час больше была его скорость на обратном пути?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica24-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 24, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 24, номер 2, год 2022."/>\n</div>\n</div>', 'а) 4n + b (км) Ответ: 4n + b километров равен весь путь поезда. б) m - 2v (км) Ответ: m - 2v километров ему ещё осталось пробежать. в) s : 2 - s : 3 = (3s - 2s) : 6 = s : 6 (км/ч) Ответ: на s : 6 километров в час больше была его скорость на обратном пути.', '<p>\nа) 4n + b (км)<br/>\n<b>Ответ:</b> 4n + b километров равен весь путь поезда.<br/><br/>\nб) m - 2v (км)<br/>\n<b>Ответ:</b> m - 2v километров ему ещё осталось пробежать.<br/><br/>\nв) s : 2 - s : 3 = (3s - 2s) : 6 = s : 6 (км/ч)<br/>\n<b>Ответ:</b> на s : 6 километров в час больше была его скорость на обратном пути.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-24/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica24-nomer2.jpg', 'peterson/3/part3/page24/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'afb30618564bdee5b6da8a70ebf7cfe7e176ea425fccbadd722ddd586119633c', '2,3,4', '["больше"]'::jsonb, 'блицтурнир а) после того как поезд проехал 4 часа со скоростью n км/ч, ему ещё осталось проехать b км. чему равен весь путь поезда? б) спортсмен бежал 2 часа со скоростью v км/ч. длина всей дистанции равна m км. сколько километров ему ещё осталось пробежать? в) самолёт пролетел s км за 3 часа, а в обратную сторону-за 2 часа. на сколько километров в час больше была его скорость на обратном пути');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 24, '3', 2, 'Выполни действия: а) 5 ч 12 мин - 3 ч 48 мин б) 16 мин 39 с + 4 мин 56 с в) 42 ц 94 кг + 2 т 6 кг г) 12 т 50 кг - 52 ц 90 кг', '</p> \n<p class="text">Выполни действия: </p> \n\n<p class="description-text"> \nа) 5 ч 12 мин - 3 ч 48 мин <br/>   	\nб) 16 мин 39 с + 4 мин 56 с <br/> 	\nв) 42 ц 94 кг + 2 т 6 кг<br/>\nг) 12 т 50 кг - 52 ц 90 кг\n</p>', 'а) 5 ч 12 мин - 3 ч 48 мин = 1 ч 24 мин б) 16 мин 39 с + 4 мин 56 с = 21 мин 35 с в) 42 ц 94 кг + 2 т 6 кг = 6 т 3 ц г) 12 т 50 кг - 52 ц 90 кг = 6 т 7 ц 60 кг', '<p>\nа) 5 ч 12 мин - 3 ч 48 мин = 1 ч 24 мин <br/>   	 \nб) 16 мин 39 с + 4 мин 56 с = 21 мин 35 с <br/> 	\nв) 42 ц 94 кг + 2 т 6 кг = 6 т 3 ц <br/> \nг) 12 т 50 кг - 52 ц 90 кг = 6 т 7 ц 60 кг\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-24/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'dfea15bf6674064b4cb0b08ca3fafe041b6143a344c0fd4365bc59bcddaca7ad', '2,3,4,5,6,12,16,39,42,48', NULL, 'выполни действия:а) 5 ч 12 мин-3 ч 48 мин б) 16 мин 39 с+4 мин 56 с в) 42 ц 94 кг+2 т 6 кг г) 12 т 50 кг-52 ц 90 кг');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 24, '4', 3, 'Составь программу действий и вычисли: а) 80 : (16 · 4 + 320 : 20) + 74 · 0 - (18 - 18) : 30 б) 0 : 48 + 50 · (10000 - 9999) - 40 : (27 · 3 - 320 : 4)', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) 80 : (16 · 4 + 320 : 20) + 74 · 0 - (18 - 18) : 30 <br/>\nб) 0 : 48 + 50 · (10000 - 9999) -  40 : (27 · 3 - 320 : 4)\n</p>', 'а) 80 : (16 · 4 + 320 : 20) + 74 · 0 - (18 - 18) : 30 = 1 16 · 4 = 64 320 : 20 = 16 64 + 16 = 80 80 : 80 = 1 18 - 18 = 0 0 : 30 = 0 74 · 0 = 0 1 + 0 = 1 1 - 0 = 1 б) 0 : 48 + 50 · (10000 - 9999) - 40 : (27 · 3 - 320 : 4) = 10 10000 - 9999 = 1 50 · 1 = 50 27 · 3 = 81 320 : 4 = 80 81 - 80 = 1 40 : 1 = 40 0 : 48 = 0 0 + 50 = 50 50 - 40 = 10', '<p>\nа) 80 : (16 · 4 + 320 : 20) + 74 · 0 - (18 - 18) : 30 = 1<br/>\n16 · 4 = 64<br/>\n320 : 20 = 16<br/>\n64 + 16 = 80<br/>\n80 : 80 = 1<br/>\n18 - 18 = 0<br/>\n0 : 30 = 0<br/>\n74 · 0 = 0<br/>\n1 + 0 = 1<br/>\n1 - 0 = 1<br/><br/>\nб) 0 : 48 + 50 · (10000 - 9999) -  40 : (27 · 3 - 320 : 4) = 10<br/>\n10000 - 9999 = 1<br/>\n50 · 1 = 50<br/>\n27 · 3 = 81<br/>\n320 : 4 = 80<br/>\n81 - 80 = 1<br/>\n40 : 1 = 40<br/>\n0 : 48 = 0<br/>\n0 + 50 = 50<br/>\n50 - 40 = 10\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-24/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a05b9ec55d95ffcddd5777a20355124af552885a5bc35ad22b94a15afa8c563f', '0,3,4,16,18,20,27,30,40,48', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) 80:(16*4+320:20)+74*0-(18-18):30 б) 0:48+50*(10000-9999)-40:(27*3-320:4)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 24, '5', 4, 'Реши уравнения с комментированием и сделай проверку: а) (90 · b + 60) : 3 = 80 б) 1400 : (35 - y) - 29 = 41', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (90 · b + 60) : 3 = 80<br/>             \nб) 1400 : (35 - y) - 29 = 41\n</p>', 'а) (90 · b + 60) : 3 = 80 Чтобы найти делимое 90 · b + 60 надо делитель умножить на частное 90 · b + 60 = 3 · 80 90 · b + 60 = 240 Чтобы найти слагаемое 90 · b надо из суммы вычесть известное слагаемое 90 · b = 240 - 60 90 · b = 80 Чтобы найти множитель b надо произведение разделить на известный множитель b = 80 : 90 b = 8 : 9 Проверка: (90 · 8 : 9 + 60) : 3 = 80 б) 1400 : (35 - y) - 29 = 41 Чтобы найти уменьшаемое 1400 : (35 - y) надо сложить вычитаемое с разностью 1400 : (35 - y) = 41 + 29 1400 : (35 - y) = 70 Чтобы найти делитель 35 – y надо делимое разделить на частное 35 - y = 1400 : 70 35 - y = 20 Чтобы найти вычитаемое у надо из уменьшаемого отнять разность y = 35 - 20 y = 15 Проверка: 1400 : (35 - 15) = 41', '<p>\nа) (90 · b + 60) : 3 = 80<br/>          \nЧтобы найти делимое 90 · b + 60 надо делитель умножить на частное<br/>\n90 · b + 60 = 3 · 80<br/>\n90 · b + 60 = 240<br/>\nЧтобы найти слагаемое 90 · b надо из суммы вычесть известное слагаемое<br/>\n90 · b = 240 - 60<br/>\n90 · b = 80<br/>\nЧтобы найти множитель b надо произведение разделить на известный множитель<br/>\nb = 80 : 90<br/>\nb = 8 : 9<br/>\n<b>Проверка:</b> (90 · 8 : 9 + 60) : 3 = 80<br/><br/>\nб) 1400 : (35 - y) - 29 = 41<br/>\nЧтобы найти уменьшаемое 1400 : (35 - y) надо сложить вычитаемое с разностью<br/>\n1400 : (35 - y) = 41 + 29<br/>\n1400 : (35 - y) = 70<br/>\nЧтобы найти делитель 35 – y надо делимое разделить на частное<br/>\n35 - y = 1400 : 70<br/>\n35 - y = 20<br/>\nЧтобы найти вычитаемое у надо из уменьшаемого отнять разность<br/>\ny = 35 - 20<br/>\ny = 15<br/>\n<b>Проверка:</b> 1400 : (35 - 15) = 41\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-24/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f786f084b53899a3fdd2bdaaf5d120cf15a0c5b40be40f8bbd494386e2d6f79d', '3,29,35,41,60,80,90,1400', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (90*b+60):3=80 б) 1400:(35-y)-29=41');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 24, '6', 5, 'Периметр прямоугольника равен 50 см, а его длина – 15 см. Чему равна площадь этого прямоугольника?', '</p> \n<p class="text">Периметр прямоугольника равен 50 см, а его длина – 15 см. Чему равна площадь этого прямоугольника?</p>', '50 = 2(a + 15), а = 50 : 2 - 15 S = а · 15, S = 15 · (50 : 2 - 15) S = 15 · 10 S = 150 (см 2 )', '<p>\n50 = 2(a + 15), а = 50 : 2 - 15<br/>\nS = а · 15,<br/> \nS = 15 · (50 : 2 - 15)<br/>\nS = 15 · 10<br/>\nS = 150 (см<sup>2</sup>)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-24/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'bca29cb62f1379d05a6b1dc308cd29d0006aacb2ade265b8472863f4fde590bf', '15,50', '["периметр","площадь"]'::jsonb, 'периметр прямоугольника равен 50 см, а его длина-15 см. чему равна площадь этого прямоугольника');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 25, '7', 0, 'Длина коробки, имеющей форму прямоугольного параллелепипеда, равна 30 см, а ширина – 20 см. 1) Чему равна высота коробки, если её объём равен 7200 см 3 ? 2) Какую площадь и какой периметр имеет дно коробки? 3) Коробку надо перевязать лентой, как показано на рисунке. Какой длины должна быть эта лента, если на узел и бант надо дополнительно предусмотреть 26 см?', '</p> \n<p class="text">Длина коробки, имеющей форму прямоугольного параллелепипеда, равна 30 см, а ширина – 20 см. </p> \n\n<p class="description-text"> \n1) Чему равна высота коробки, если её объём равен 7200 см<sup>3</sup>? <br/>\n2) Какую площадь и какой периметр имеет дно коробки?<br/>\n3) Коробку надо перевязать лентой, как показано на рисунке. Какой длины должна быть эта лента, если на узел и бант надо дополнительно предусмотреть 26 см?\n\n</p>\n\n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 25, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 25, номер 7, год 2022."/>\n</div>\n</div>', '1) 7200 : (30 · 20) = 7200 : 600 = 12 (см) 2) Р = 2(20 + 30) Р = 2 · 50 Р = 100 (см) S = 20 · 30 S = 600 (см 2 ) 3) 26 + (20 + 30) · 2 + 12 · 4 = 26 + 100 + 48 = 174 (см)', '<p>\n1) 7200 : (30 · 20) = 7200 : 600 = 12 (см)<br/>\n2) Р = 2(20 + 30)<br/>\nР = 2 · 50<br/>\nР = 100 (см)<br/>\nS = 20 · 30<br/>\nS = 600 (см<sup>2</sup>)<br/>\n3) 26 + (20 + 30) · 2 + 12 · 4 = 26 + 100 + 48 = 174 (см)\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-25/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer7.jpg', 'peterson/3/part3/page25/task7_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '961b99b44eaf6cb1cb954e40ee11a4f51e4904a900fe78ebea2e7e2b1ece760b', '1,2,3,20,26,30,7200', '["периметр","площадь"]'::jsonb, 'длина коробки, имеющей форму прямоугольного параллелепипеда, равна 30 см, а ширина-20 см. 1) чему равна высота коробки, если её объём равен 7200 см 3 ? 2) какую площадь и какой периметр имеет дно коробки? 3) коробку надо перевязать лентой, как показано на рисунке. какой длины должна быть эта лента, если на узел и бант надо дополнительно предусмотреть 26 см');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 25, '8', 1, 'Сравни выражения, не выполняя вычислений. Обоснуй свой ответ. 3974 + 815     815 + 3794      786 · 29     786 + 29 76012 - 32     76012 - 23       3420 : 6     3420 · 2 9083 - 96     9100 - 96            2158 : 26     2158 : 83', '</p> \n<p class="text">Сравни выражения, не выполняя вычислений. Обоснуй свой ответ.</p> \n\n<p class="description-text"> \n3974 + 815 <span class="okon">   </span> 815 + 3794      786 · 29 <span class="okon">   </span> 786 + 29<br/>\n76012 - 32 <span class="okon">   </span> 76012 - 23       3420 : 6 <span class="okon">   </span> 3420 · 2<br/>\n9083 - 96 <span class="okon">   </span> 9100 - 96            2158 : 26 <span class="okon">   </span> 2158 : 83\n</p>', '3974 + 815 ˃ 815 + 3794      786 · 29 ˃ 786 + 29 76012 - 32 < 76012 - 23       3420 : 6 < 3420 · 2 9083 - 96 < 9100 - 96            2158 : 26 ˃ 2158 : 83', '<p>\n3974 + 815 ˃ 815 + 3794      786 · 29 ˃ 786 + 29<br/>\n76012 - 32 &lt; 76012 - 23       3420 : 6 &lt; 3420 · 2<br/>\n9083 - 96 &lt; 9100 - 96            2158 : 26 ˃ 2158 : 83\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-25/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'cd41f20ed52b7188cb6541c97f5212710ea2b07ea7b3888e9e7417c256b9c8da', '2,6,23,26,29,32,83,96,786,815', '["сравни"]'::jsonb, 'сравни выражения, не выполняя вычислений. обоснуй свой ответ. 3974+815     815+3794      786*29     786+29 76012-32     76012-23       3420:6     3420*2 9083-96     9100-96            2158:26     2158:83');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 25, '9', 2, 'Запиши множество делителей и множество кратных числа 19.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 19.</p>', '19, 38, 57, 76, 95, 114, 133 и т. д.', '<p>\n19, 38, 57, 76, 95, 114, 133 и т. д.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-25/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '536645c9c23dcb255ef75b9dd1e756075f8b37e88b1afe93492dc9d29e223d98', '19', NULL, 'запиши множество делителей и множество кратных числа 19');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 25, '10', 3, 'Найди площадь прямоугольного участка по указанным размерам. Сколько различных способов решения имеет эта задача? Что ты замечаешь?', '</p> \n<p class="text">Найди площадь прямоугольного участка по указанным размерам.  Сколько различных способов решения имеет эта задача? Что ты замечаешь?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="200" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer10.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 25, номер 10, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 25, номер 10, год 2022."/>\n</div>\n</div>', 'S = 50 · 38 + 6 · 38 = 1900 + 228 = 2128 (м 2 ) S = (50 + 6) · 38 = 56 · 38 = 2128 (м 2 )', '<p>\nS = 50 · 38 + 6 · 38 = 1900 + 228 = 2128 (м<sup>2</sup>)<br/>\nS = (50 + 6) · 38 = 56 · 38 = 2128 (м<sup>2</sup>)\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="100" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer10-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 25, номер 10-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 25, номер 10-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-25/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer10.jpg', 'peterson/3/part3/page25/task10_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer10-1.jpg', 'peterson/3/part3/page25/task10_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '068078ba35ecac3c2abcb93a060f4a364475f1e078192a46d605e13d4c0f3075', NULL, '["найди","площадь","раз"]'::jsonb, 'найди площадь прямоугольного участка по указанным размерам. сколько различных способов решения имеет эта задача? что ты замечаешь');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 25, '11', 4, 'Найди площадь прямоугольника, разбивая его на части удобным способом:', '</p> \n<p class="text">Найди площадь прямоугольника, разбивая его на части удобным способом:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer11.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 25, номер 11, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 25, номер 11, год 2022."/>\n</div>\n</div>', 'S = 90 · 70 + 7 · 5 = 6300 + 35 = 6335 (дм 2 )', '<p>\nS = 90 · 70 + 7 · 5 = 6300 + 35 = 6335 (дм<sup>2</sup>)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-25/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer11.jpg', 'peterson/3/part3/page25/task11_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c9c5f668d51644812e427764b7f63c6f7482f470356f84d62389c625a47e57c9', NULL, '["найди","площадь","раз"]'::jsonb, 'найди площадь прямоугольника, разбивая его на части удобным способом');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 25, '12', 5, 'Вычисли. Расположи ответы в порядке убывания и расшифруй имя сказочного героя. Из какой он сказки?', '</p> \n<p class="text">Вычисли. Расположи ответы в порядке убывания и расшифруй имя сказочного героя. Из какой он сказки?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer12.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 25, номер 12, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 25, номер 12, год 2022."/>\n</div>\n</div>', 'Д – 5632084 - 5294352 = 337732 У – 19050 · 50 = 952500 Н – 313920 : 4 = 78630 В – 94203 + 186902 + 56618 = 94203 + 243520 = 337723 Г – 3052 · 600 = 1831200 И – 647040 : 8 = 80880 ГУДВИН Один из героев сказочного цикла о Волшебной стране и Изумрудном городе.', '<p>\nД – 5632084 - 5294352 = 337732<br/>  \nУ – 19050 · 50 = 952500<br/> \nН – 313920 : 4 = 78630<br/>\nВ – 94203 + 186902 + 56618 = 94203 + 243520 = 337723<br/>\nГ – 3052 · 600 = 1831200<br/>\nИ – 647040 : 8 = 80880<br/>\nГУДВИН<br/>\nОдин из героев сказочного цикла о Волшебной стране и Изумрудном городе.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-25/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica25-nomer12.jpg', 'peterson/3/part3/page25/task12_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3f1619d62a8ec99c5e891b3923c0163fd634647a2570566feaae74d6ec570084', NULL, '["вычисли"]'::jsonb, 'вычисли. расположи ответы в порядке убывания и расшифруй имя сказочного героя. из какой он сказки');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 25, '13', 6, 'D – множество девочек класса, M – множество мальчиков этого же класса. Что представляют собой множества D ⋂ M и D ⋃ M?', '</p> \n<p class="text">D – множество девочек класса, M – множество мальчиков этого же класса. Что представляют собой множества D ⋂ M и D ⋃ M?</p>', 'D ⋂ M нету и D ⋃ M этот класс 20 : 4=5 (слив) 20 - 5 = 15 (слив) 15 : 3 = 5 (слив) 20 : 4 = 5 (слив) 20 - 5 = 15 (слив) 15 : 3 = 5 (слив) 5 + 5 = 10 – слив всего взяла Наташа', '<p>\nD ⋂ M нету и D ⋃ M этот класс\n</p>\n\n\n<p>\n20 : 4=5 (слив)<br/> \n20 - 5 = 15 (слив)<br/> \n15 : 3 = 5 (слив)<br/> \n20 : 4 = 5 (слив)<br/> \n20 - 5 = 15 (слив)<br/> \n15 : 3 = 5 (слив) <br/>\n5 + 5 = 10 – слив всего взяла Наташа\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-25/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '428fe9c36cc440ee4c5a790b468e700680d626ccd8a76155373e06d2a0fdaa6d', NULL, NULL, 'd-множество девочек класса, m-множество мальчиков этого же класса. что представляют собой множества d ⋂ m и d ⋃ m');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 26, '1', 0, 'а) Объясни по рисунку, как умножить число на сумму, и выполни умножение: a · (b + c) = a · b + a · c 21 · 56 = 21 · (50 + 6) = ... б) Используя рисунок, объясни способ записи умножения на двузначное число в столбик:', '</p> \n<p class="text">а) Объясни по рисунку, как умножить число на сумму, и выполни умножение:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="180" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica26-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 26, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 26, номер 1, год 2022."/>\n</div>\n</div>\n\n<p class="description-text"> \na · (b + c) = a · b + a · c<br/>\n21 · 56 = 21 · (50 + 6) = ...\n</p>\n\n<p class="text">б) Используя рисунок, объясни способ записи умножения на двузначное число в столбик:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica26-nomer1-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 26, номер 1-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 26, номер 1-1, год 2022."/>\n</div>\n</div>', 'Умножение многозначного числа на двузначное. Чтобы умножить любое число на двузначное, можно умножить это число сначала на единицы, а потом на десятки и полученные произведения сложить. В записи суммы число десятков сдвигают на 1 разряд влево.', '<p>\nУмножение многозначного числа на двузначное. <br/>\nЧтобы умножить любое число на двузначное, можно умножить это число сначала на единицы, а потом на десятки и полученные произведения сложить.<br/>\nВ записи суммы число десятков сдвигают на 1 разряд влево.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Объясни по рисунку, как умножить число на сумму, и выполни умножение: a · (b + c) = a · b + a · c 21 · 56 = 21 · (50 + 6) = ...","solution":""},{"letter":"б","condition":"Используя рисунок, объясни способ записи умножения на двузначное число в столбик:","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-26/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica26-nomer1.jpg', 'peterson/3/part3/page26/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica26-nomer1-1.jpg', 'peterson/3/part3/page26/task1_condition_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5d351d78b898ad02ad2aa9c9eed9c21ed0e7e2681b979b8ae971f6a3fd1a220b', '6,21,50,56', '["столбик"]'::jsonb, 'а) объясни по рисунку, как умножить число на сумму, и выполни умножение:a*(b+c)=a*b+a*c 21*56=21*(50+6)=... б) используя рисунок, объясни способ записи умножения на двузначное число в столбик');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 26, '2', 1, 'В кинотеатре 18 рядов по 32 места в каждом ряду. Сколько всего мест в кинотеатре? Найди в данной записи ответы на вопросы: Сколько мест в 8 рядах? Сколько мест в 10 рядах? Сколько всего мест в кинотеатре?', '</p> \n<p class="text">В кинотеатре 18 рядов по 32 места в каждом ряду. Сколько всего мест в кинотеатре?<br/>\nНайди в данной записи ответы на вопросы:<br/>\nСколько мест в 8 рядах?<br/>\nСколько мест в 10 рядах?<br/>\nСколько всего мест в кинотеатре?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="220" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica26-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 26, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 26, номер 2, год 2022."/>\n</div>\n</div>', '18 · 32 = 576 мест всего 8 · 32 = 256 мест в 8 рядах 10 · 32 = 320 мест в 10 рядах.', '<p>\n18 · 32 = 576 мест всего <br/>\n8 · 32 = 256 мест в 8 рядах <br/>\n10 · 32 = 320 мест в 10 рядах.\n\n</p>', 'Умножение многозначного числа на двузначное Чтобы умножить любое число на двузначное, можно умножить это число сначала на единицы, а потом на десятки и полученные произведения сложить. В записи суммы число десятков сдвигают на 1 разряд влево. Пример:', '<div class="recomended-block">\n<span class="title">Умножение многозначного числа на двузначное </span>\n<p>\nЧтобы умножить любое число на двузначное, можно умножить это число сначала на единицы, а потом на десятки и полученные произведения сложить.<br/>\nВ записи суммы число десятков сдвигают на 1 разряд влево.<br/>\nПример:\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica26-spravka.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 26, справка, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 26, справка, год 2022."/>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-26/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica26-nomer2.jpg', 'peterson/3/part3/page26/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '85adc4c01a41c070c73ba07d27d3c9f7f451b5f15299a4eec8406e483d6d64d6', '8,10,18,32', '["найди"]'::jsonb, 'в кинотеатре 18 рядов по 32 места в каждом ряду. сколько всего мест в кинотеатре? найди в данной записи ответы на вопросы:сколько мест в 8 рядах? сколько мест в 10 рядах? сколько всего мест в кинотеатре');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 27, '3', 0, 'Правильно ли Максим решил и прокомментировал пример?', '</p> \n<p class="text">Правильно ли Максим решил и прокомментировал пример?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica27-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 27, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 27, номер 3, год 2022."/>\n</div>\n</div>', 'Верно.', '<p>\nВерно.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-27/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica27-nomer3.jpg', 'peterson/3/part3/page27/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'be9f36b3057a90709a37f3041a7ff0ac9d6cba31b71b21702d2c74183e38475a', NULL, '["реши"]'::jsonb, 'правильно ли максим решил и прокомментировал пример');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 27, '4', 1, 'Реши примеры с комментированием: а) 92 · 89            в) 138 · 56 б) 57 · 95            г) 296 · 23 д) 906 · 15          ж) 2384 · 47 е) 709 · 84          з) 9051 · 72', '</p> \n<p class="text">Реши примеры с комментированием:</p> \n\n<p class="description-text"> \nа) 92 · 89            в) 138 · 56<br/>  	\nб) 57 · 95            г) 296 · 23<br/>  	\nд) 906 · 15          ж) 2384 · 47<br/>\nе) 709 · 84          з) 9051 · 72\n\n</p>', 'а) 92 · 89 = 8188', '<p>\nа) 92 · 89 = 8188\n</p>\n\n<div class="img-wrapper-460">\n<img width="110" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica27-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 27, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 27, номер 4, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-27/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica27-nomer4.jpg', 'peterson/3/part3/page27/task4_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '013fa1f7ea2e795b1f558c20932f5fd19dba8bc2c20a5e77a1a49a43af696fb1', '15,23,47,56,57,72,84,89,92,95', '["реши"]'::jsonb, 'реши примеры с комментированием:а) 92*89            в) 138*56 б) 57*95            г) 296*23 д) 906*15          ж) 2384*47 е) 709*84          з) 9051*72');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 27, '5', 2, 'а) Лыжники были в походе 7 дней. Каждый день они шли по 6 ч со скоростью 9 км/ч. Сколько километров прошли лыжники? б) Миша пробежал 8 кругов со скоростью 200 м/мин. Сколько времени он бежал, если длина одного круга 400 м?', '</p> \n<p class="text">а) Лыжники были в походе 7 дней. Каждый день они шли по 6 ч со скоростью 9 км/ч. Сколько километров прошли лыжники?<br/>\nб) Миша пробежал 8 кругов со скоростью 200 м/мин. Сколько времени он бежал, если длина одного круга 400 м?\n</p>', 'а) 7 · 6 · 9 = 42 · 9 = 378 (км) Ответ: 378 километров прошли лыжники. б) 8 · 400 : 200 = 3200 : 200 = 16 (мин) Ответ: 16 минут он бежал.', '<p>\nа) 7 · 6 · 9 = 42 · 9 = 378 (км)<br/>\n<b>Ответ:</b> 378 километров прошли лыжники.<br/><br/>\nб) 8 · 400 : 200 = 3200 : 200 = 16 (мин)<br/>\n<b>Ответ:</b> 16 минут он бежал.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Лыжники были в походе 7 дней. Каждый день они шли по 6 ч со скоростью 9 км/ч. Сколько километров прошли лыжники?","solution":"7 · 6 · 9 = 42 · 9 = 378 (км) Ответ: 378 километров прошли лыжники."},{"letter":"б","condition":"Миша пробежал 8 кругов со скоростью 200 м/мин. Сколько времени он бежал, если длина одного круга 400 м?","solution":"8 · 400 : 200 = 3200 : 200 = 16 (мин) Ответ: 16 минут он бежал."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-27/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '4689cdebcb460783796765cb76ddedc337247ad80e1376dbe58ed90c5f5e7437', '6,7,8,9,200,400', NULL, 'а) лыжники были в походе 7 дней. каждый день они шли по 6 ч со скоростью 9 км/ч. сколько километров прошли лыжники? б) миша пробежал 8 кругов со скоростью 200 м/мин. сколько времени он бежал, если длина одного круга 400 м');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 27, '6', 3, 'Расстояние от Москвы до Новосибирска 3320 км. Поезд проходит его за 40 ч, а самолёт пролетает – в 10 раз быстрее. На сколько часов меньше лететь до Новосибирска самолётом, чем ехать поездом? Во сколько раз скорость поезда меньше скорости самолёта?', '</p> \n<p class="text">Расстояние от Москвы до Новосибирска 3320 км. Поезд проходит его за 40 ч, а самолёт пролетает – в 10 раз быстрее. На сколько часов меньше лететь до Новосибирска самолётом, чем ехать поездом? Во сколько раз скорость поезда меньше скорости самолёта?</p>', '40 - 40 : 10 = 40 - 4 = 36 (ч) 3320 : 40 = 83 (км/ч) – скорость поезда 83 · 10 = 830 (км/ч) – скорость самолёта 830 : 83 = 10 (раз) Ответ: на 36 часов меньше лететь до Новосибирска самолётом, чем ехать поездом. В 10 раз скорость поезда меньше скорости самолёта.', '<p>\n40 - 40 : 10 = 40 - 4 = 36 (ч)<br/>\n3320 : 40 = 83 (км/ч) – скорость поезда<br/>\n83 · 10 = 830 (км/ч) – скорость самолёта<br/>\n830 : 83 = 10 (раз)<br/>\n<b>Ответ:</b> на 36 часов меньше лететь до Новосибирска самолётом, чем ехать поездом. В 10 раз скорость поезда меньше скорости самолёта. \n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-27/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '20c128f00bac2ca62ceb050835be6578d95507f5fb00885c4a67b5cb5b523b79', '10,40,3320', '["меньше","раз"]'::jsonb, 'расстояние от москвы до новосибирска 3320 км. поезд проходит его за 40 ч, а самолёт пролетает-в 10 раз быстрее. на сколько часов меньше лететь до новосибирска самолётом, чем ехать поездом? во сколько раз скорость поезда меньше скорости самолёта');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 27, '7', 4, 'Составь и реши уравнения: а) На сколько надо умножить число 60, чтобы получить 4320? б) Какое число надо разделить на 700, чтобы получить 506? в) На сколько надо разделить 8500, чтобы получить 500?', '</p> \n<p class="text">Составь и реши уравнения:<br/>\nа) На сколько надо умножить число 60, чтобы получить 4320?<br/>\nб) Какое число надо разделить на 700, чтобы получить 506?<br/>\nв) На сколько надо разделить 8500, чтобы получить 500?\n</p>', 'а) 4320 : 60 = 72 60 · 72 = 4320 б) 700 · 506 = 354200 354200 : 700 = 506 в) 8500 : 500 = 17 8500 : 17 = 500', '<p>\nа) 4320 : 60 = 72<br/>\n60 · 72 = 4320<br/><br/>\nб) 700 · 506 = 354200<br/>\n354200 : 700 = 506<br/><br/>\nв) 8500 : 500 = 17<br/>\n8500 : 17 = 500\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-27/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8fe9a85d9cd06fa6c0ed0da214eefde01b8ec271ca21b0f996ece17ed6557fa3', '60,500,506,700,4320,8500', '["раздели","реши","раз"]'::jsonb, 'составь и реши уравнения:а) на сколько надо умножить число 60, чтобы получить 4320? б) какое число надо разделить на 700, чтобы получить 506? в) на сколько надо разделить 8500, чтобы получить 500');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 27, '8', 5, 'Запиши множество делителей и множество кратных числа 20.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 20.</p>', '0,40,60,80,100,120,140,160,180 и так далее', '<p>\n0,40,60,80,100,120,140,160,180 и так далее\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-27/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'de7aaaad30ce9eec9ded21003ee5d83ffdc487e0833b299c8b466d064cca9c60', '20', NULL, 'запиши множество делителей и множество кратных числа 20');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 27, '9', 6, 'Выполни действия: а) 4 ч 58 мин + 2 ч 17 мин - 3 ч 29 мин б) 18 мин 9 с - 7 мин 46 с + 48 мин 35 с в) 4 мин 52 с · 5 г) 7 ч 3 мин : 9', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 4 ч 58 мин + 2 ч 17 мин - 3 ч 29 мин <br/> 	\nб) 18 мин 9 с - 7 мин 46 с + 48 мин 35 с <br/> 	\nв) 4 мин 52 с · 5<br/>\nг) 7 ч 3 мин : 9\n</p>', 'а) 4 ч 58 мин + 2 ч 17 мин - 3 ч 29 мин = 3 ч 46 мин б) 18 мин 9 с - 7 мин 46 с + 48 мин 35 с = 10 мин 23 с + 48 мин 35 с = 58 мин 58 с в) 4 мин 52 с · 5 = 292 с · 5 = 1460 с = 24 мин 20 с г) 7 ч 3 мин : 9 = 45 мин : 9 = 5 мин', '<p>\nа) 4 ч 58 мин + 2 ч 17 мин - 3 ч 29 мин = 3 ч 46 мин <br/> 	\nб) 18 мин 9 с - 7 мин 46 с + 48 мин 35 с = 10 мин 23 с + 48 мин 35 с = 58 мин 58 с<br/>\nв) 4 мин 52 с · 5 = 292 с · 5 = 1460 с = 24 мин 20 с<br/>\nг) 7 ч 3 мин : 9 = 45 мин : 9 = 5 мин\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-27/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1649ba1896279052fdc1f56f0002ad9bc2dcc3b0a081d7cf8b3c7d3310d9cfb8', '2,3,4,5,7,9,17,18,29,35', NULL, 'выполни действия:а) 4 ч 58 мин+2 ч 17 мин-3 ч 29 мин б) 18 мин 9 с-7 мин 46 с+48 мин 35 с в) 4 мин 52 с*5 г) 7 ч 3 мин:9');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 27, '10', 7, 'В вазе лежат персик, ананас и банан. Сколькими различными способами из неё можно взять один, два или три фрукта?', '</p> \n<p class="text">В вазе лежат персик, ананас и банан. Сколькими различными способами из неё можно взять один, два или три фрукта?</p>', 'Один фрукт: банан, ананас, персик – тремя способами. Два фрукта: банан и персик, персик и ананас, ананас и банан – тремя способами. Три фрукта: ананас и банан, и персик – одним способом.', '<p>\nОдин фрукт: банан, ананас, персик – тремя способами.<br/>\nДва фрукта: банан и персик, персик и ананас, ананас и банан – тремя способами.<br/>\nТри фрукта: ананас и банан, и персик – одним способом.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-27/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '07a00c3d5915b82876ae60786c585818052ab4f2fad2e1746da5a4016bd82eab', NULL, '["раз"]'::jsonb, 'в вазе лежат персик, ананас и банан. сколькими различными способами из неё можно взять один, два или три фрукта');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 28, '1', 0, 'Составь выражение и найди его значение: а) Одна ручка стоит 17 р. Сколько надо заплатить за 5 таких ручек? б) Метр ткани стоит 120 р. Сколько стоят 3 м этой ткани? в) Литр сока стоит a р. Сколько стоят n л этого сока? Что общего во всех этих задачах? О каких величинах в них идёт речь? Как найти стоимость товара, зная его цену и количество?', '</p> \n<p class="text">Составь выражение и найди его значение:<br/> \nа) Одна ручка стоит 17 р. Сколько надо заплатить за 5 таких ручек?<br/>\nб) Метр ткани стоит 120 р. Сколько стоят 3 м этой ткани?<br/>\nв) Литр сока стоит a р. Сколько стоят n л этого сока?<br/>\nЧто общего во всех этих задачах? О каких величинах в них идёт речь? Как найти стоимость товара, зная его цену и количество?\n</p>', 'а) 17 · 5 = 85 (р) б) 120 · 3 = 360 (р) в) a · n (р) Нужно найти стоимость различного товара по разной цене. Рубль, единицу длины – метр. Взаимосвязь между ценой, количеством и стоимостью.', '<p>\nа) 17 · 5 = 85 (р)<br/>\nб) 120 · 3 = 360 (р)<br/>\nв) a · n (р)<br/>\nНужно найти стоимость различного товара по разной цене. Рубль, единицу длины – метр. Взаимосвязь между ценой, количеством и стоимостью.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-28/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '7dd6e5d01c1e19655310872d8a892d2a645b9c3f6f959b59742c1ee471d2a31b', '3,5,17,120', '["найди"]'::jsonb, 'составь выражение и найди его значение:а) одна ручка стоит 17 р. сколько надо заплатить за 5 таких ручек? б) метр ткани стоит 120 р. сколько стоят 3 м этой ткани? в) литр сока стоит a р. сколько стоят n л этого сока? что общего во всех этих задачах? о каких величинах в них идёт речь? как найти стоимость товара, зная его цену и количество');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 28, '2', 1, 'Найди неизвестные значения величин по формуле стоимости C = a · n:', '</p> \n<p class="text">Найди неизвестные значения величин по формуле стоимости C = a · n:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica28-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 27, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 27, номер 2, год 2022."/>\n</div>\n</div>', 'а) 360 : 60 = 6 (кг), 40 · 4 = 160 (р.), 950 : 5 = 190 (р./м) б) 840 : 4 = 210 (р./шт), 56 : 8 = 7 (л), 70 · 5 = 350 (р.)', '<p>\nа) 360 : 60 = 6 (кг), 40 · 4 = 160 (р.), 950 : 5 = 190 (р./м)<br/>  \nб) 840 : 4 = 210 (р./шт), 56 : 8 = 7 (л), 70 · 5 = 350 (р.)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-28/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica28-nomer2.jpg', 'peterson/3/part3/page28/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b42972078882a6714420d09970ec4372fffa96d3d04a0b45287512e07d99a6a3', NULL, '["найди"]'::jsonb, 'найди неизвестные значения величин по формуле стоимости c=a*n');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 28, '3', 2, 'Цена книги 45 р. Чему равна стоимость 2 книг, 4 книг, 6 книг, n книг? Заполни в тетради таблицу. Запиши формулу зависимости стоимости C купленных книг от их количества n.', '</p> \n<p class="text">Цена книги 45 р. Чему равна стоимость 2 книг, 4 книг, 6 книг, n книг? Заполни в тетради таблицу. Запиши формулу зависимости стоимости C купленных книг от их количества n.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica28-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 27, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 27, номер 3, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica28-nomer3-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 27, номер 3-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 27, номер 3-1, год 2022."/>', 'Формула стоимости Пусть C – стоимость товара, a – его цена (то есть стоимость единицы товара – 1 штуки, 1 метра, 1 килограмма, 1 литра и т. д.), а n – количество товара в выбранных единицах. Тогда: C = a · n Полученное равенство называется формулой стоимости. Оно означает, что стоимость равна цене, умноженной на количество товара. Из формулы стоимости по правилу нахождения неизвестного множителя легко выразить величины a и n: a = C : n      n = C : a • Цена равна стоимости, делённой на количество товара. • Количество товара равно стоимости, делённой на цену.', '<div class="recomended-block">\n<span class="title">Формула стоимости</span>\n<p>\nПусть C – стоимость товара, a – его цена (то есть стоимость единицы товара – 1 штуки, 1 метра, 1 килограмма, 1 литра и т. д.), а n – количество товара в выбранных единицах. Тогда:<br/>\nC = a · n<br/>\nПолученное равенство называется формулой стоимости. Оно означает, что стоимость равна цене, умноженной на количество товара.<br/>\nИз формулы стоимости по правилу нахождения неизвестного множителя легко выразить величины a и n:<br/>\na = C : n      n = C : a<br/>\n• Цена равна стоимости, делённой на количество товара.<br/>\n• Количество товара равно стоимости, делённой на цену.\n\n</p>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-28/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica28-nomer3.jpg', 'peterson/3/part3/page28/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica28-nomer3-1.jpg', 'peterson/3/part3/page28/task3_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1a4355ed3f2b40c2f19ae6641f1cf079a99fe33dfc59a7aaa7fd7fd839f22e95', '2,4,6,45', '["заполни"]'::jsonb, 'цена книги 45 р. чему равна стоимость 2 книг, 4 книг, 6 книг, n книг? заполни в тетради таблицу. запиши формулу зависимости стоимости c купленных книг от их количества n');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 29, '4', 0, 'У Игоря 240 р. Сколько тетрадей он сможет купить, если их цена 10 р., 12 р., 15 р., 20 р., a р.? Заполни таблицу. Запиши формулу зависимости количества купленных тетрадей n от их цены a.', '</p> \n<p class="text">У Игоря 240 р. Сколько тетрадей он сможет купить, если их цена 10 р., 12 р., 15 р., 20 р., a р.? Заполни таблицу. Запиши формулу зависимости количества купленных тетрадей n от их цены a.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 29, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 29, номер 4, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer4-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 29, номер 4-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 29, номер 4-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-29/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer4.jpg', 'peterson/3/part3/page29/task4_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer4-1.jpg', 'peterson/3/part3/page29/task4_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '67b9c343dc77da09b5230436c5097642207c2279aafc431cc8987737a6e62a0f', '10,12,15,20,240', '["заполни"]'::jsonb, 'у игоря 240 р. сколько тетрадей он сможет купить, если их цена 10 р., 12 р., 15 р., 20 р., a р.? заполни таблицу. запиши формулу зависимости количества купленных тетрадей n от их цены a');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 29, '5', 1, 'Реши примеры с комментированием. Найди сумму и разность наибольшего и наименьшего из получившихся чисел: а) 85 · 54        б) 279 · 68 в) 406 · 49      г) 9032 · 97', '</p> \n<p class="text">Реши примеры с комментированием. Найди сумму и разность наибольшего и наименьшего из получившихся чисел:</p> \n\n<p class="description-text"> \nа) 85 · 54        б) 279 · 68<br/>\nв) 406 · 49      г) 9032 · 97\n</p>', 'а) 85 · 54 = 4590', '<p>\nа) 85 · 54 = 4590\n</p>\n\n<div class="img-wrapper-460">\n<img width="110" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 29, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 29, номер 5, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-29/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer5.jpg', 'peterson/3/part3/page29/task5_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ef065f07f70f19c4cd53e345016745f0ac5b1d895f78113e688ae32f65856132', '49,54,68,85,97,279,406,9032', '["найди","реши","разность","больше","меньше","раз"]'::jsonb, 'реши примеры с комментированием. найди сумму и разность наибольшего и наименьшего из получившихся чисел:а) 85*54        б) 279*68 в) 406*49      г) 9032*97');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 29, '6', 2, 'Выполни действия: а) 415 · 36          б) 709 · 79 в) 3705 · 68        г) 20507 · 94', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 415 · 36          б) 709 · 79<br/>  \nв) 3705 · 68        г) 20507 · 94\n</p>', 'а) 415 · 36 = 14940', '<p>\nа) 415 · 36 = 14940\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 29, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 29, номер 6, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-29/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer6.jpg', 'peterson/3/part3/page29/task6_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '59d7362d9992cf9c5fe4cdf7f5a0d23eb64c27a46efafa73254dc96135f3eea4', '36,68,79,94,415,709,3705,20507', NULL, 'выполни действия:а) 415*36          б) 709*79 в) 3705*68        г) 20507*94');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 29, '7', 3, 'Мотоциклист выехал из Москвы в Клин со скоростью 45 км/ч. В дороге он сделал две остановки: одну на – 25 мин, а вторую – на 35 мин. В Клин мотоциклист прибыл в 13 ч 20 мин. В котором часу он выехал из Москвы, если расстояние от Москвы до Клина равно 90 км?', '</p> \n<p class="text">Мотоциклист выехал из Москвы в Клин со скоростью 45 км/ч. В дороге он сделал две остановки: одну на – 25 мин, а вторую – на 35 мин.  В Клин мотоциклист прибыл в 13 ч 20 мин. В котором часу он выехал из Москвы, если расстояние от Москвы до Клина равно 90 км?</p>', '13 ч 20 мин - 90 км : 45 км/ч - 25 мин - 35 мин = 13 ч 20 мин - 2 ч - 25 мин - 35 мин = 10 ч 20 мин Ответ: в 10 ч 20 минут он выехал из Москвы, если расстояние от Москвы до Клина равно 90 км.', '<p>\n13 ч 20 мин - 90 км : 45 км/ч - 25 мин - 35 мин = 13 ч 20 мин - 2 ч - 25 мин - 35 мин = 10 ч 20 мин<br/>\n<b>Ответ:</b> в 10 ч 20 минут он выехал из Москвы, если расстояние от Москвы до Клина равно 90 км.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-29/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '4bd483080b11aaf7410f43dc6f69320e067e37c723f0d09f7dad751bc5e74d88', '13,20,25,35,45,90', '["равно"]'::jsonb, 'мотоциклист выехал из москвы в клин со скоростью 45 км/ч. в дороге он сделал две остановки:одну на-25 мин, а вторую-на 35 мин. в клин мотоциклист прибыл в 13 ч 20 мин. в котором часу он выехал из москвы, если расстояние от москвы до клина равно 90 км');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 29, '8', 4, 'Повтори таблицу мер длины. Используя её, вырази данные величины в указанных единицах измерения: а) 3 см 5 мм = … мм       б) 3 км 5 м = … м 3 дм 5 см = … см           3 км 5 м = … дм 3 дм 5 мм = … мм          3 км 5 м = … см 3 дм 5 см = … мм           3 км 5 м = … мм 3 м 5 дм = … см             3 км 5 см = … мм', '</p> \n<p class="text">Повтори таблицу мер длины. Используя её, вырази данные величины в указанных единицах измерения:</p> \n\n<p class="description-text"> \nа) 3 см 5 мм = … мм       б)  3 км 5 м = … м<br/>  		\n3 дм 5 см = … см           3 км 5 м = … дм<br/>  		\n3 дм 5 мм = … мм          3 км 5 м = … см<br/>	\n3 дм 5 см = … мм           3 км 5 м = … мм<br/> 			\n3 м 5 дм = … см             3 км 5 см = … мм\n\n</p>\n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 29, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 29, номер 8, год 2022."/>\n</div>\n\n</div>', 'а) 3 см 5 мм = 35 мм 3 дм 5 см = 35 см 3 дм 5 мм = 305 мм 3 дм 5 см = 350 мм 3 м 5 дм = 350 см б) 3 км 5 м = 3005 м 3 км 5 м = 30050 дм 3 км 5 м = 300500 см 3 км 5 м = 3005000 мм 3 км 5 см = 3000500 мм', '<p>\nа) 3 см 5 мм = 35 мм <br/>    		\n3 дм 5 см = 35 см<br/>          			\n3 дм 5 мм = 305 мм<br/>  	 	\n3 дм 5 см = 350 мм<br/>         			\n3 м 5 дм = 350 см<br/><br/>         	\n\nб) 3 км 5 м = 3005 м<br/>  	\n3 км 5 м = 30050 дм <br/>\n3 км 5 м = 300500 см<br/> \n3 км 5 м = 3005000 мм<br/> \n3 км 5 см = 3000500 мм\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-29/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer8.jpg', 'peterson/3/part3/page29/task8_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5e97a6b59d7c55482b1ad7d4fb962ddb150f0827050d47c0166c6174013f0c11', '3,5', '["раз"]'::jsonb, 'повтори таблицу мер длины. используя её, вырази данные величины в указанных единицах измерения:а) 3 см 5 мм=... мм       б) 3 км 5 м=... м 3 дм 5 см=... см           3 км 5 м=... дм 3 дм 5 мм=... мм          3 км 5 м=... см 3 дм 5 см=... мм           3 км 5 м=... мм 3 м 5 дм=... см             3 км 5 см=... мм');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 29, '9', 5, 'Выполни действия: а) (30 км - 5 км 964 м) : 6 б) 40 км 20 м - 78 м 28 мм · 500', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) (30 км - 5 км 964 м) : 6 <br/>     \nб) 40 км 20 м - 78 м 28 мм · 500\n</p>', 'а) (30 км - 5 км 964 м) : 6 = 24 км 6 м : 6 = 24006 м : 6 = 4001 м = 4 км 1 м б) 40 км 20 м - 78 м 28 мм · 500 = 40020000 мм - 78028 мм · 500 = 40020000 мм - 39014000 мм = 1006000 = 1 км 6 м', '<p>\nа) (30 км - 5 км 964 м) : 6 = 24 км 6 м : 6 = 24006 м : 6 = 4001 м = 4 км 1 м<br/>\nб) 40 км 20 м - 78 м 28 мм · 500 = 40020000 мм - 78028 мм · 500 = 40020000 мм - 39014000 мм = 1006000 = 1 км 6 м\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="210" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 29, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 29, номер 9, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-29/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer9.jpg', 'peterson/3/part3/page29/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'abebf7f7b742212aa18d49fe86172febd27fde2f8cc61bcaff51273bc1e904c4', '5,6,20,28,30,40,78,500,964', NULL, 'выполни действия:а) (30 км-5 км 964 м):6 б) 40 км 20 м-78 м 28 мм*500');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 29, '10', 6, 'Запиши множество делителей и множество кратных числа 21.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 21.</p>', '21,42,63,84,105,126 и т. д.', '<p>\n21,42,63,84,105,126 и т. д.\n</p>\n\n\n<div class="img-wrapper-460">\n<img width="160" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer11-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 29, номер 11-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 29, номер 11-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-29/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica29-nomer11-1.jpg', 'peterson/3/part3/page29/task10_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e8e85398445de48da68e450e47339a46ab89194404e5f1b289b3048b0da8246b', '21', NULL, 'запиши множество делителей и множество кратных числа 21');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 30, '1', 0, 'Выполни умножение с комментированием: а) 36 · 79         в) 635 · 46 б) 17 · 54         г) 281 · 38 д) 508 · 75       ж) 4205 · 97 е) 902 · 23       з) 9003 · 61', '</p> \n<p class="text">Выполни умножение с комментированием:</p> \n\n<p class="description-text"> \nа) 36 · 79         в) 635 · 46<br/>  	\nб) 17 · 54         г) 281 · 38<br/>  	\nд) 508 · 75       ж) 4205 · 97<br/>\nе) 902 · 23       з) 9003 · 61\n\n</p>', 'а) 36 · 79 = 2844', '<p>\nа) 36 · 79 = 2844  \n</p>\n\n<div class="img-wrapper-460">\n<img width="110" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 30, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 30, номер 1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-30/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer1.jpg', 'peterson/3/part3/page30/task1_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '7fecbbca182069aacce4a0c6b7a0de939b28d871ff04dcfa2cf2f1d9d9be1e1c', '17,23,36,38,46,54,61,75,79,97', NULL, 'выполни умножение с комментированием:а) 36*79         в) 635*46 б) 17*54         г) 281*38 д) 508*75       ж) 4205*97 е) 902*23       з) 9003*61');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 30, '2', 1, 'Найди неизвестные значения величин по формуле стоимости С = а · n:', '</p> \n<p class="text">Найди неизвестные значения величин по формуле стоимости С = а · n:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 30, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 30, номер 2, год 2022."/>\n</div>\n\n</div>', 'а) 58 · 3 = 174', '<p>\nа) 58 · 3 = 174\n</p>\n\n<div class="img-wrapper-460">\n<img width="80" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer2-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 30, номер 2-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 30, номер 2-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-30/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer2.jpg', 'peterson/3/part3/page30/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer2-1.jpg', 'peterson/3/part3/page30/task2_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ee2d666c4c62f34df2cc3a855f5d1532e4c6e8f21aa3bc30fa5dd30be4553336', NULL, '["найди"]'::jsonb, 'найди неизвестные значения величин по формуле стоимости с=а*n');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 30, '3', 2, 'Вырази в указанных единицах измерения: а) 2 км 8 м = … м       б) 2 ч 8 мин = … мин 2 м 8 см = … см         2 сут. 8 ч = … мин 2 дм 8 мм = … мм       2 т 8 кг = … кг 2 м 8 мм = … мм         2 т 8 ц = … кг', '</p> \n<p class="text">Вырази в указанных единицах измерения: </p> \n\n<p class="description-text"> \nа) 2 км 8 м = … м       б) 2 ч 8 мин = … мин<br/>  \n2 м 8 см = … см         2 сут. 8 ч = … мин<br/>  \n2 дм 8 мм = … мм       2 т 8 кг = … кг<br/>\n2 м 8 мм = … мм         2 т 8 ц = … кг\n\n</p>', 'а) 2 км 8 м = 2008 м 2 м 8 см = 208 см 2 дм 8 мм = 208 мм 2 м 8 мм = 2008 мм б) 2 ч 8 мин = 128 мин 2 сут. 8 ч = 2360 мин 2 т 8 кг = 2008 кг 2 т 8 ц = 2800 кг', '<p>\nа) 2 км 8 м = 2008 м<br/>  		\n2 м 8 см = 208 см  <br/>  		  \n2 дм 8 мм = 208 мм <br/>   		 \n2 м 8 мм = 2008 мм<br/><br/>	    	\n\nб) 2 ч 8 мин = 128 мин<br/>\n2 сут. 8 ч = 2360 мин<br/>  \n2 т 8 кг = 2008 кг<br/> \n2 т 8 ц = 2800 кг\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-30/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'debbfc735c45e7fd258cd386dbc92d535d082dac06214915c55195ca88da25f2', '2,8', '["раз"]'::jsonb, 'вырази в указанных единицах измерения:а) 2 км 8 м=... м       б) 2 ч 8 мин=... мин 2 м 8 см=... см         2 сут. 8 ч=... мин 2 дм 8 мм=... мм       2 т 8 кг=... кг 2 м 8 мм=... мм         2 т 8 ц=... кг');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 30, '4', 3, 'а) Килограмм клубники стоит 90 р. Сколько рублей надо заплатить за 2 кг такой клубники? б) За 5 одинаковых электронных дисков с играми заплатили 800 р. Сколько рублей стоит один диск?', '</p> \n<p class="text">а) Килограмм клубники стоит 90 р. Сколько рублей надо заплатить за 2 кг такой клубники?<br/> \nб) За 5 одинаковых электронных дисков с играми заплатили 800 р. Сколько рублей стоит один диск?\n</p>', 'а) 90 · 2 = 180 (р.) Ответ: 180 рублей надо заплатить за 2 кг такой клубники. б) 800 : 5 = 160 (р.) Ответ: 160 рублей стоит один диск.', '<p>\nа) 90 · 2 = 180 (р.)<br/>\n<b>Ответ:</b> 180 рублей надо заплатить за 2 кг такой клубники.<br/><br/>\nб) 800 : 5 = 160 (р.)<br/>\n<b>Ответ:</b> 160 рублей стоит один диск.\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Килограмм клубники стоит 90 р. Сколько рублей надо заплатить за 2 кг такой клубники?","solution":"90 · 2 = 180 (р.) Ответ: 180 рублей надо заплатить за 2 кг такой клубники."},{"letter":"б","condition":"За 5 одинаковых электронных дисков с играми заплатили 800 р. Сколько рублей стоит один диск?","solution":"800 : 5 = 160 (р.) Ответ: 160 рублей стоит один диск."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-30/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8acbcb98ec044c1105f007d9838e090e2ff70e3110cc828fe67c1f5184e3baf0', '2,5,90,800', NULL, 'а) килограмм клубники стоит 90 р. сколько рублей надо заплатить за 2 кг такой клубники? б) за 5 одинаковых электронных дисков с играми заплатили 800 р. сколько рублей стоит один диск');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 30, '5', 4, 'Цена одного билета на поезд равна 800 р. Сколько таких билетов можно купить на 2000 р.? Сколько рублей ещё останется?', '</p> \n<p class="text">Цена одного билета на поезд равна 800 р. Сколько таких билетов можно купить на 2000 р.? Сколько рублей ещё останется?</p>', '2000 : 800 = 2 + 400 Ответ: 2 таких билетов можно купить на 2000 р., 400 рублей ещё останется.', '<p>\n2000 : 800 = 2 + 400 <br/>\n<b>Ответ:</b> 2 таких билетов можно купить на 2000 р., 400 рублей ещё останется.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-30/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6b3cc20994c9a66f6bef2b2a656a9f8e1acccc3dc13e0eb89ea1a0969d58feae', '800,2000', NULL, 'цена одного билета на поезд равна 800 р. сколько таких билетов можно купить на 2000 р.? сколько рублей ещё останется');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 30, '6', 5, 'Реши уравнения с комментированием и сделай проверку: а) (900 - x : 6) · 5 = 4200 б) 325 + (90 - n) : 17 = 330', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (900 - x : 6) · 5 = 4200<br/>  	\nб) 325 + (90 - n) : 17 = 330\n</p>', 'а) (900 - x : 6) · 5 = 4200 Чтобы множитель (900 - х : 6) надо произведение разделить на известный множитель 900 - х : 6 = 4200 : 5 900 - х : 6 = 840 Чтобы найти вычитаемое х : 6 надо из уменьшаемое вычесть разность х : 6 = 900 - 840 х : 6 = 60 Чтобы найти делимое х надо делитель умножить на частное х = 6 · 60 х = 360 Проверка: (900 - 360 : 6) · 5 = 4200 б) 325 + (90 - n) : 17 = 330 Чтобы найти слагаемое (90 - n) : 17 надо вычесть из суммы известное слагаемое (90 - n) : 17 = 330 - 325 (90 - n) : 17 = 5 Чтобы найти делимое (90 - n) надо делитель умножить на частное 90 - n = 17 · 5 90 - n = 85 Чтобы найти вычитаемое n надо из уменьшаемого вычесть разность n = 90 - 85 n = 5 Проверка: 325 + (90 - 5) : 17 = 330', '<p>\nа) (900 - x : 6) · 5 = 4200  <br/>	\nЧтобы множитель (900 - х : 6) надо произведение разделить на известный множитель<br/>\n900 - х : 6 = 4200 : 5<br/>\n900 - х : 6 = 840<br/>\nЧтобы найти вычитаемое х : 6 надо из уменьшаемое вычесть разность<br/>\nх : 6 = 900 - 840<br/>\nх : 6 = 60<br/>\nЧтобы найти делимое х надо делитель умножить на частное<br/>\nх = 6 · 60<br/>\nх = 360 <br/>\n<b>Проверка:</b> (900 - 360 : 6) · 5 = 4200<br/><br/>\nб) 325 + (90 - n) : 17 = 330<br/>\nЧтобы найти слагаемое (90 - n) : 17 надо вычесть из суммы известное слагаемое<br/>\n(90 - n) : 17 = 330 - 325<br/>\n(90 - n) : 17 = 5<br/>\nЧтобы найти делимое (90 - n) надо делитель умножить на частное<br/>\n90 - n = 17 · 5<br/>\n90 - n = 85<br/>\nЧтобы найти вычитаемое n надо из уменьшаемого вычесть разность<br/>\nn = 90 - 85<br/>\nn = 5<br/>\n<b>Проверка:</b> 325 + (90 - 5) : 17 = 330  \n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-30/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '05a8a5f9c3cf32a4ebb930ee8742288c4079bbf509db89c9d35562add28699aa', '5,6,17,90,325,330,900,4200', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (900-x:6)*5=4200 б) 325+(90-n):17=330');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 30, '7', 6, 'В строке 56 печатных знаков, а на странице – 36 строк. Сколько печатных знаков уместится на 64 страницах?', '</p> \n<p class="text">В строке 56 печатных знаков, а на странице – 36 строк. Сколько печатных знаков уместится на 64 страницах?</p>', '56 · 36 · 64 = 2016 · 64 = 129024', '<p>\n56 · 36 · 64 = 2016 · 64 = 129024\n</p>\n\n<div class="img-wrapper-460">\n<img width="260" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 30, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 30, номер 7, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-30/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer7.jpg', 'peterson/3/part3/page30/task7_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '9b80d5604b36961e36d7c71ac8c63bab9fdf3c1251383309e77b733f3e3af759', '36,56,64', NULL, 'в строке 56 печатных знаков, а на странице-36 строк. сколько печатных знаков уместится на 64 страницах');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 30, '8', 7, 'а) Поезд шёл 18 ч со скоростью 76 км/ч и 16 ч со скоростью 72 км/ч. Какое расстояние прошёл поезд за всё это время? б) Почтальон проехал на велосипеде 36 км за 2 ч. Затем он уменьшил скорость на 2 км/ч и ехал ещё 3 ч. Сколько всего километров проехал на велосипеде почтальон?', '</p> \n<p class="text">а) Поезд шёл 18 ч со скоростью 76 км/ч и 16 ч со скоростью 72 км/ч. Какое расстояние прошёл поезд за всё это время?<br/>\nб) Почтальон проехал на велосипеде 36 км за 2 ч. Затем он уменьшил скорость на 2 км/ч и ехал ещё 3 ч. Сколько всего километров проехал на велосипеде почтальон?\n</p>', 'а) 76 · 18 + 16 · 72 = 1368 + 1152 = 1440 (км)', '<p>\nа) 76 · 18 + 16 · 72 = 1368 + 1152 = 1440 (км)\n</p>\n\n<div class="img-wrapper-460">\n<img width="140" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 30, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 30, номер 8, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Поезд шёл 18 ч со скоростью 76 км/ч и 16 ч со скоростью 72 км/ч. Какое расстояние прошёл поезд за всё это время?","solution":"76 · 18 + 16 · 72 = 1368 + 1152 = 1440 (км)"},{"letter":"б","condition":"Почтальон проехал на велосипеде 36 км за 2 ч. Затем он уменьшил скорость на 2 км/ч и ехал ещё 3 ч. Сколько всего километров проехал на велосипеде почтальон?","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-30/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica30-nomer8.jpg', 'peterson/3/part3/page30/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '834d6a1824806b9f372be274b8802ce8955930bc6794ee3de1ca7a34bf2715e7', '2,3,16,18,36,72,76', NULL, 'а) поезд шёл 18 ч со скоростью 76 км/ч и 16 ч со скоростью 72 км/ч. какое расстояние прошёл поезд за всё это время? б) почтальон проехал на велосипеде 36 км за 2 ч. затем он уменьшил скорость на 2 км/ч и ехал ещё 3 ч. сколько всего километров проехал на велосипеде почтальон');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 30, '9', 8, 'Почтовый голубь должен доставить донесение на расстояние 130 км. Скорость голубя 50 км/ч. Успеет ли он доставить это донесение: а) за 2 часа? б) за 3 часа?', '</p> \n<p class="text">Почтовый голубь должен доставить донесение на расстояние 130 км. Скорость голубя 50 км/ч. Успеет ли он доставить это донесение: а) за 2 часа? б) за 3 часа?</p>', 'а) 50 · 2 = 100 (км) Ответ: голубь успеет доставить это донесение за 2 часа. б) 50 · 3 = 150 (км) Ответ: голубь не успеет доставить это донесение за 3 часа.', '<p>\nа) 50 · 2 = 100 (км)<br/>\n<b>Ответ:</b> голубь успеет доставить это донесение за 2 часа.<br/><br/>\nб) 50 · 3 = 150 (км)<br/>\n<b>Ответ:</b> голубь не успеет доставить это донесение за 3 часа.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-30/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '319c6a515dca43f9bc1ae994bc2354693e426d3905dff91b04e07e3d6e546ab0', '2,3,50,130', NULL, 'почтовый голубь должен доставить донесение на расстояние 130 км. скорость голубя 50 км/ч. успеет ли он доставить это донесение:а) за 2 часа? б) за 3 часа');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 31, '10', 0, 'Запиши предложение в виде равенства: а) n на 17 = меньше, чем m б) x в 8 раз меньше, чем = y в) a на 92 меньше, чем b г) k в 5 раз больше, чем d', '</p> \n<p class="text">Запиши предложение в виде равенства:</p> \n\n<p class="description-text"> \nа) n на 17 = меньше, чем m<br/> 	\nб) x в 8 раз меньше, чем = y<br/>	\nв) a на 92 меньше, чем b<br/>\nг) k в 5 раз больше, чем d\n</p>', 'а) n = m - 17       в) a = b - 92 б) x · 8 = y          г) k = 5 · d', '<p>\nа) n = m - 17       в) a = b - 92<br/>\nб) x · 8 = y          г) k = 5 · d\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-31/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '70b0bc4370b14ca9d8d2f19b1975ab5e912584cae197e76dfa9a7a7faaa469bc', '5,8,17,92', '["больше","меньше","раз"]'::jsonb, 'запиши предложение в виде равенства:а) n на 17=меньше, чем m б) x в 8 раз меньше, чем=y в) a на 92 меньше, чем b г) k в 5 раз больше, чем d');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 31, '11', 1, 'Продолжительность дня равна t ч. Чему равна продолжительность ночи? Составь выражение и найди его значение, если t = 8, 10, 12. Какие значения может принимать переменная t?', '</p> \n<p class="text">Продолжительность дня равна t ч. Чему равна продолжительность ночи? Составь выражение и найди его значение, если t = 8, 10, 12. Какие значения может принимать переменная t?</p>', '24 - t 24 - 8 = 16 (ч) 24 - 10 = 14 (ч) 24 - 12 = 12 (ч) Переменная t может принимать значения от 0 до 24 часов.', '<p>\n24 - t <br/>\n24 - 8 = 16 (ч)<br/>\n24 - 10 = 14 (ч)<br/>\n24 - 12 = 12 (ч) <br/>\nПеременная t может принимать значения от 0 до 24 часов.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-31/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '12bba50a6869d52412005fdaec9915afa9724b8beb4f41f7c390feb7f7043bc3', '8,10,12', '["найди"]'::jsonb, 'продолжительность дня равна t ч. чему равна продолжительность ночи? составь выражение и найди его значение, если t=8, 10, 12. какие значения может принимать переменная t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 31, '12', 2, 'Запиши программу действий в виде выражения со скобками. Найди значение полученного выражения.', '</p> \n<p class="text">Запиши программу действий в виде выражения со скобками. Найди значение полученного выражения.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica31-nomer12.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 31, номер 12, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 31, номер 12, год 2022."/>\n</div>\n</div>', '(2488 + 4512) · 593 - (485830 - 37598) : 8 = 7000 · 593 - 448232 : 8 = 4151000 - 56029 = 4094971', '<p>\n(2488 + 4512) · 593 - (485830 - 37598) : 8 = 7000 · 593 - 448232 : 8 = 4151000 - 56029 = 4094971\n</p>\n\n<div class="img-wrapper-460">\n<img width="200" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica31-nomer12-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 31, номер 12-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 31, номер 12-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-31/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica31-nomer12.jpg', 'peterson/3/part3/page31/task12_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica31-nomer12-1.jpg', 'peterson/3/part3/page31/task12_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '48756d77a43d1f728cf175e80dc6b9d05bdd1dce4de7558ce8f1cea41eb9ed05', NULL, '["найди"]'::jsonb, 'запиши программу действий в виде выражения со скобками. найди значение полученного выражения');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 31, '13', 3, 'Реши уравнения устно. Расположи ответы в порядке возрастания. Расшифруй имя сказочного героя. Узнай название книги и имя её автора.', '</p> \n<p class="text">Реши уравнения устно. Расположи ответы в порядке возрастания. Расшифруй имя сказочного героя. Узнай название книги и имя её автора.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica31-nomer13.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 31, номер 13, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 31, номер 13, год 2022."/>\n</div>\n</div>', 'М - 9 + b = 12 b = 12 - 9 b = 3 И - 8 · m = 480 m = 480 : 8 m = 60 С - 40 - c = 12 c = 40 - 12 c = 28 О - 90 : d = 5 d = 90 : 5 d = 18 Т - a · 50 = 250 a = 250 : 50 a = 5 К - n – 27 = 8 n = 27 + 8 n = 35 А - 52 : t = 13 t = 52 : 13 t = 4 Н - k : 19 = 4 k = 19 · 4 k = 76 Р - 34 - x = 17 x = 34 - 17 x = 17 МАТРОСКИН Дядя Фёдор, пёс и кот – Эдуард Успенский', '<p>\nМ - 9 + b = 12<br/>\nb = 12 - 9<br/>\nb = 3<br/>\nИ - 8 · m = 480	<br/>\nm = 480 : 8<br/>\nm = 60<br/>\nС - 40 - c = 12<br/>\nc = 40 - 12<br/>\nc = 28<br/><br/>\n\nО - 90 : d = 5<br/>\nd = 90 : 5<br/>\nd = 18	<br/>\nТ - a · 50 = 250<br/>\na = 250 : 50<br/>\na = 5	<br/>\nК - n – 27 = 8<br/>\nn = 27 + 8<br/>\nn = 35	<br/><br/>\n\nА - 52 : t = 13<br/>\nt = 52 : 13<br/>\nt = 4<br/>\nН - k : 19 = 4<br/>\nk = 19 · 4<br/>\nk = 76<br/>\nР - 34 - x = 17<br/>\nx = 34 - 17<br/>\nx = 17<br/><br/>\n\nМАТРОСКИН<br/>\nДядя Фёдор, пёс и кот – Эдуард Успенский\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-31/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica31-nomer13.jpg', 'peterson/3/part3/page31/task13_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'aa19a9dc8a775d8602ae5c4a4a404b3c6ea4db9d1afc30471be443905b839378', NULL, '["реши"]'::jsonb, 'реши уравнения устно. расположи ответы в порядке возрастания. расшифруй имя сказочного героя. узнай название книги и имя её автора');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 31, '14', 4, 'Запиши множество делителей и множество кратных числа 22.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 22.</p>', 'Делители: 1, 22, 2, 11. Кратные: 22, 44, 66, 88, 11, 132. 25 - (16 + 16 - 10) = 25 - 22 = 3 (ученика) Ответ: 3 ученика не записались ни в один из этих кружков.', '<p>\nДелители: 1, 22, 2, 11. Кратные: 22, 44, 66, 88, 11, 132.\n</p>\n\n\n<p>\n25 - (16 + 16 - 10) = 25 - 22 = 3 (ученика)<br/>\n<b>Ответ:</b> 3 ученика не записались ни в один из этих кружков.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-31/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd269b807426efa35f54ce8ea4c7639c1145894f61b79980bc75b27cb63cd4c59', '22', NULL, 'запиши множество делителей и множество кратных числа 22');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 32, '1', 0, 'а) Что общего в выражениях? Вспомни правило умножения круглых чисел и вычисли: 400 · 70        160 · 300 9 · 80000      250 · 4000 б) Как записывают умножение круглых чисел в столбик? Почему? Приведи свой пример.', '</p> \n<p class="text">а) Что общего в выражениях? Вспомни правило умножения круглых чисел и вычисли:<br/>\n400 · 70        160 · 300<br/>               \n9 · 80000      250 · 4000<br/><br/> \nб) Как записывают умножение круглых чисел в столбик? Почему? Приведи свой пример.\n</p>', 'а) 400 · 70 = 28000 160 · 300 = 48000 9 · 80000 = 720000 250 · 4000 = 1000000 б) сначала будем умножать на однозначное число, поэтому это число записываем под единицами, а 0 в стороне: Сначала умножаем на однозначное число, а затем, чтобы умножить на 100, к результату приписываем 00.', '<p>\nа) 400 · 70 = 28000<br/>	\n160 · 300 = 48000<br/>		\n9 · 80000 = 720000<br/>	\n250 · 4000 = 1000000<br/><br/>\nб) сначала будем умножать на однозначное число, поэтому это число записываем под единицами, а 0 в стороне: Сначала умножаем на однозначное число, а затем, чтобы умножить на 100, к результату приписываем 00.\n</p>\n\n<div class="img-wrapper-460">\n<img width="140" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 32, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 32, номер 1, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Что общего в выражениях? Вспомни правило умножения круглых чисел и вычисли: 400 · 70        160 · 300 9 · 80000      250 · 4000","solution":"400 · 70 = 28000 160 · 300 = 48000 9 · 80000 = 720000 250 · 4000 = 1000000"},{"letter":"б","condition":"Как записывают умножение круглых чисел в столбик? Почему? Приведи свой пример.","solution":"сначала будем умножать на однозначное число, поэтому это число записываем под единицами, а 0 в стороне: Сначала умножаем на однозначное число, а затем, чтобы умножить на 100, к результату приписываем 00."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-32/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer1.jpg', 'peterson/3/part3/page32/task1_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '57269c49be1fc0a2c1aa2b33edba9ce25c885d85d0bba07b9de4241672bcfd6e', '9,70,160,250,300,400,4000,80000', '["вычисли","столбик"]'::jsonb, 'а) что общего в выражениях? вспомни правило умножения круглых чисел и вычисли:400*70        160*300 9*80000      250*4000 б) как записывают умножение круглых чисел в столбик? почему? приведи свой пример');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 32, '2', 1, 'Выполни действия: а) 360 · 7500      б) 2800 · 940 в) 50900 · 62      г) 73050 · 8600', '</p> \n<p class="text">Выполни действия:<br/>	\nа) 360 · 7500      б) 2800 · 940<br/>	 	\nв) 50900 · 62      г) 73050 · 8600\n</p>', 'а) 360 · 7500 = 2700000', '<p>\nа) 360 · 7500 = 2700000  \n</p>\n\n<div class="img-wrapper-460">\n<img width="180" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 32, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 32, номер 2, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-32/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer2.jpg', 'peterson/3/part3/page32/task2_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6a3c9ef8ecda1072df8ad8880f6d5996b09ba8d0a220c374da057a77d1618293', '62,360,940,2800,7500,8600,50900,73050', NULL, 'выполни действия:а) 360*7500      б) 2800*940 в) 50900*62      г) 73050*8600');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 32, '3', 2, 'Вычисли. Расположи ответы в порядке убывания. Расшифруй имя короля сказочного государства, который избавил детей от скучных занятий в школе. Узнай название этой книги и имя её автора.', '</p> \n<p class="text">Вычисли. Расположи ответы в порядке убывания. Расшифруй имя короля сказочного государства, который избавил детей от скучных занятий в школе. Узнай название этой книги и имя её автора.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="300" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 32, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 32, номер 3, год 2022."/>\n</div>\n</div>', 'Ш - 5400 · 62 = 334800', '<p>\nШ - 5400 · 62 = 334800   \n</p>\n\n<div class="img-wrapper-460">\n<img width="160" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer3-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 32, номер 3-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 32, номер 3-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-32/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer3.jpg', 'peterson/3/part3/page32/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer3-1.jpg', 'peterson/3/part3/page32/task3_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '75a3c70313afa27267852b051962404975ba1995fa134e4b2524f20b095d3b6b', NULL, '["вычисли"]'::jsonb, 'вычисли. расположи ответы в порядке убывания. расшифруй имя короля сказочного государства, который избавил детей от скучных занятий в школе. узнай название этой книги и имя её автора');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 32, '4', 3, 'Составь программу действий и вычисли: а) 860 · 900 - 6750 : 5 · (24 + 44) б) (64 + 137) · 28 · 910 - 560 772 : 9', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) 860 · 900 - 6750 : 5 · (24 + 44)<br/> \nб) (64 + 137) · 28 · 910 - 560 772 : 9\n</p>', 'а) 860 · 900 - 6750 : 5 · (24 + 44) = 773691 24 + 44 = 68 5 · 68 = 340 6750 : 340 = 19 + 290', '<p>\nа) 860 · 900 - 6750 : 5 · (24 + 44) = 773691<br/>\n24 + 44 = 68<br/>\n5 · 68 = 340<br/>\n6750 : 340 = 19 + 290\n</p>\n\n<div class="img-wrapper-460">\n<img width="70" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 32, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 32, номер 4, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-32/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-nomer4.jpg', 'peterson/3/part3/page32/task4_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ec97291d31d1b01b38bba72bed1e6de86a4cfaace3f7ad59c8a74989a5bb45ee', '5,9,24,28,44,64,137,560,772,860', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) 860*900-6750:5*(24+44) б) (64+137)*28*910-560 772:9');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 32, '5', 4, 'Сравни в каждом равенстве числа, обозначенные буквами. Какое из них больше, а какое меньше? На сколько? a = b + 18      k - t = 5 x = y - 9         n - 4 = m', '</p> \n<p class="text">Сравни в каждом равенстве числа, обозначенные буквами. Какое из них больше, а какое меньше? На сколько?<br/>\na = b + 18      k - t = 5 <br/> 	\nx = y - 9         n - 4 = m\n</p>', 'а больше b на 18 k больше t на 5 х меньше у на 9 n больше m на 4', '<p>\nа больше b на 18<br/>\nk больше t на 5<br/>\nх меньше у на 9<br/>\nn больше m на 4\n\n</p>', 'Алгоритм умножения круглых многозначных чисел 1. Записать множители в столбик, не глядя на нули. 2. Выполнить умножение многозначных чисел, не глядя на нули. 3. Записать в произведении справа столько нулей, сколько в обоих множителях вместе.', '<div class="recomended-block">\n<span class="title">Алгоритм умножения круглых многозначных чисел</span>\n<p>\n1. Записать множители в столбик, не глядя на нули.<br/>\n2. Выполнить умножение многозначных чисел, не глядя на нули.<br/>\n3. Записать в произведении справа столько нулей, сколько в обоих множителях вместе.\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica32-spravka.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 32, справка, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 32, справка, год 2022."/>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-32/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'adcdd02eae79624feccb07a1b9c8bc77202bc3fa882796d319893dd95cc5571f', '4,5,9,18', '["сравни","больше","меньше"]'::jsonb, 'сравни в каждом равенстве числа, обозначенные буквами. какое из них больше, а какое меньше? на сколько? a=b+18      k-t=5 x=y-9         n-4=m');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 33, '6', 0, 'Реши уравнения с комментированием и сделай проверку: а) (k : 16) · 13 + 11 = 50 б) 14 - 72 : (d - 3) = 8', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (k : 16) · 13 + 11 = 50<br/> \nб) 14 - 72 : (d - 3) = 8\n</p>', 'а) (k : 16) · 13 + 11 = 50 Чтобы найти слагаемое (k : 16) · 13 надо из суммы вычесть известное слагаемое (k : 16) · 13 = 50 - 11 (k : 16) · 13 = 39 Чтобы найти множитель (k : 16) надо произведение разделить на известный множитель (k : 16) = 39 : 13 k : 16 = 3 Чтобы найти делимое надо частное умножить на делитель k = 3 · 16 k = 48 Проверка: (48 : 16) · 13 + 11 = 50 б) 14 - 72 : (d - 3) = 8 Чтоб найти вычитаемое 72 : (d - 3) надо из уменьшаемого вычесть разность 72 : (d - 3) = 14 - 8 72 : (d - 3) = 6 Чтобы найти делитель (d - 3) надо делимое разделить на частное d - 3 = 72 : 6 d - 3 = 12 Чтобы найти уменьшаемое надо к частному прибавить вычитаемое d = 12 + 3 d = 15 Проверка: 14 – 72 : (15 - 3) = 8', '<p>\nа) (k : 16) · 13 + 11 = 50 <br/>  \nЧтобы найти слагаемое (k : 16) · 13 надо из суммы вычесть известное слагаемое<br/> \n(k : 16) · 13 = 50 - 11<br/> \n(k : 16) · 13 = 39<br/> \nЧтобы найти множитель (k : 16) надо произведение разделить на известный множитель<br/> \n(k : 16) = 39 : 13<br/> \nk : 16 = 3 <br/> \nЧтобы найти делимое надо частное умножить на делитель<br/> \nk = 3 · 16<br/> \nk = 48<br/> \n<b>Проверка:</b> (48 : 16) · 13 + 11 = 50<br/><br/> \nб) 14 - 72 : (d - 3) = 8<br/>\nЧтоб найти вычитаемое 72 : (d - 3) надо из уменьшаемого вычесть разность<br/>\n72 : (d - 3) = 14 - 8 <br/>\n72 : (d - 3) = 6<br/>\nЧтобы найти делитель (d - 3) надо делимое разделить на частное<br/>\nd - 3 = 72 : 6<br/>\nd - 3 = 12<br/>\nЧтобы найти уменьшаемое надо к частному прибавить вычитаемое<br/>\nd = 12 + 3<br/>\nd = 15<br/>\n<b>Проверка:</b> 14 – 72 : (15 - 3) = 8\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-33/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e7ed1e8c15a54954f4e2b4265d3339418b1ad012d193968ab5c878941a7beb0c', '3,8,11,13,14,16,50,72', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (k:16)*13+11=50 б) 14-72:(d-3)=8');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 33, '7', 1, 'а) Одна роза стоит 40 р. Сколько надо заплатить за букет из 7 роз? б) Одна конфета стоит 12 р. Сколько таких конфет можно купить на 60 рублей? в) За 20 календарей заплатили 1800 р. Сколько рублей стоит один календарь?', '</p> \n<p class="text">а) Одна роза стоит 40 р. Сколько надо заплатить за букет из 7 роз?<br/>\nб) Одна конфета стоит 12 р. Сколько таких конфет можно купить на 60 рублей?<br/> \nв) За 20 календарей заплатили 1800 р. Сколько рублей стоит один календарь?\n</p>', 'а) 40 · 7 = 280 (р.) Ответ: 28 рублей надо заплатить за букет из 7 роз. б) 60 : 12 = 5 (конфет) Ответ: 5 таких конфет можно купить на 60 рублей. в) 1800 : 20 = 90 (р.) Ответ: 90 рублей стоит один календарь.', '<p>\nа) 40 · 7 = 280 (р.) <br/>\n<b>Ответ:</b> 28 рублей надо заплатить за букет из 7 роз.<br/><br/>\nб) 60 : 12 = 5 (конфет)<br/>\n<b>Ответ:</b> 5 таких конфет можно купить на 60 рублей.<br/><br/>\nв) 1800 : 20 = 90 (р.)<br/>\n<b>Ответ:</b> 90 рублей стоит один календарь.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Одна роза стоит 40 р. Сколько надо заплатить за букет из 7 роз?","solution":"40 · 7 = 280 (р.) Ответ: 28 рублей надо заплатить за букет из 7 роз."},{"letter":"б","condition":"Одна конфета стоит 12 р. Сколько таких конфет можно купить на 60 рублей?","solution":"60 : 12 = 5 (конфет) Ответ: 5 таких конфет можно купить на 60 рублей."},{"letter":"в","condition":"За 20 календарей заплатили 1800 р. Сколько рублей стоит один календарь?","solution":"1800 : 20 = 90 (р.) Ответ: 90 рублей стоит один календарь."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-33/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6e45c859d9af1ee1f245de92e05dd62b5fa431fd30f8ab41b376a3e8cde3597c', '7,12,20,40,60,1800', NULL, 'а) одна роза стоит 40 р. сколько надо заплатить за букет из 7 роз? б) одна конфета стоит 12 р. сколько таких конфет можно купить на 60 рублей? в) за 20 календарей заплатили 1800 р. сколько рублей стоит один календарь');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 33, '8', 2, 'За футболку и 4 пары носков заплатили 200 рублей. Футболка стоит 80 р. Сколько рублей стоит одна пара носков?', '</p> \n<p class="text">За футболку и 4 пары носков заплатили 200 рублей. Футболка стоит 80 р. Сколько рублей стоит одна пара носков?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica33-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 33, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 33, номер 8, год 2022."/>\n</div>\n</div>', '(200 - 80) : 4 = 120 : 4 = 30 (р.) Ответ: 30 рублей стоит одна пара носков.', '<p>\n(200 - 80) : 4 = 120 : 4 = 30 (р.)<br/>\n<b>Ответ:</b> 30 рублей стоит одна пара носков.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-33/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica33-nomer8.jpg', 'peterson/3/part3/page33/task8_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2f35f510ac0d9a308b8d31bf7a0f034a285226130677382da45b7d0eb05229a0', '4,80,200', NULL, 'за футболку и 4 пары носков заплатили 200 рублей. футболка стоит 80 р. сколько рублей стоит одна пара носков');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 33, '9', 3, 'У Оли было 200 р. Она купила 3 тетради по цене 15 р., 2 ручки по 37 р. и 6 карандашей по 8 р. Сколько денег у неё осталось? Сможет ли она купить на них шоколадку за 32 р.?', '</p> \n<p class="text">У Оли было 200 р. Она купила 3 тетради по цене 15 р., 2 ручки по 37 р. и 6 карандашей по 8 р. Сколько денег у неё осталось? Сможет ли она купить на них шоколадку за 32 р.?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica33-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 33, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 33, номер 9, год 2022."/>\n</div>\n</div>', '200 - (3 · 15 + 2 · 27 + 6 · 8) = 200 - (45 + 54 + 48) = 200 - 147 = 53 (р.) 53 - 32 = 21 (р.) Ответ: 53 рубля у неё осталось, она сможет купить на них шоколадку за 32 р..', '<p>\n200 - (3 · 15 + 2 · 27 + 6 · 8) = 200 - (45 + 54 + 48) = 200 - 147 = 53 (р.)<br/>\n53 - 32 = 21 (р.)<br/>\n<b>Ответ:</b> 53 рубля у неё осталось, она сможет купить на них шоколадку за 32 р..\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-33/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica33-nomer9.jpg', 'peterson/3/part3/page33/task9_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8bc8d21b87685fd251201438e01ec26a57de5cd4a0135262c293745f2d7a6687', '2,3,6,8,15,32,37,200', NULL, 'у оли было 200 р. она купила 3 тетради по цене 15 р., 2 ручки по 37 р. и 6 карандашей по 8 р. сколько денег у неё осталось? сможет ли она купить на них шоколадку за 32 р');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 33, '10', 4, 'Запиши множество делителей и множество кратных числа 23.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 23.</p>', 'Делителей два это 1, 23. Кратных числа 23 множество.', '<p>\nДелителей два это 1, 23. Кратных числа 23 множество.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-33/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '411c4f27b90e88693a0f07d39216af08933039711fccaee6d4aaf5d1511960a0', '23', NULL, 'запиши множество делителей и множество кратных числа 23');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 33, '11', 5, 'Масса первого арбуза равна a кг. Масса второго арбуза – на 3 кг меньше массы первого. А масса третьего арбуза – в 2 раза больше массы второго. Чему равна масса трёх арбузов вместе? Составь выражение и найди его значение при a = 8.', '</p> \n<p class="text">Масса первого арбуза равна a кг. Масса второго арбуза – на 3 кг меньше массы первого. А масса третьего арбуза – в 2 раза больше массы второго. Чему равна масса трёх арбузов вместе?<br/>\nСоставь выражение и найди его значение при a = 8.\n</p>', 'а + а - 3 + (а - 3) · 2 8 + (8 - 3) + (8 - 3) · 2 = 8 + 5 + 10 = 23 (кг) Ответ: 23 килограмма равна масса трёх арбузов.', '<p>\nа + а - 3 + (а - 3) · 2 <br/>\n8 + (8 - 3) + (8 - 3) · 2 = 8 + 5 + 10 = 23 (кг)<br/>\n<b>Ответ:</b> 23 килограмма равна масса трёх арбузов.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-33/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '4e533ee7b82f2ec64b1e82b6bae3de2e28e1d809f92d02932e87e846669477b8', '2,3,8', '["найди","больше","меньше","раз","раза"]'::jsonb, 'масса первого арбуза равна a кг. масса второго арбуза-на 3 кг меньше массы первого. а масса третьего арбуза-в 2 раза больше массы второго. чему равна масса трёх арбузов вместе? составь выражение и найди его значение при a=8');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 33, '12', 6, 'Тарас бежит со скоростью 150 м/мин, а Юра – со скоростью 12 км/ч. Кто из них бежит быстрее?', '</p> \n<p class="text">Тарас бежит со скоростью 150 м/мин, а Юра – со скоростью 12 км/ч. Кто из них бежит быстрее?</p>', '12 км/ч = 12 · 1000 : 60 = 12000 : 60 = 200 (м/мин) 200 м/мин ˃ 150 м/мин Ответ: Юра бежит быстрее Тараса. 7 - 7 + 7 - 7 + 7 - 7 + 7 = 7 Наибольшее значение: 7 · 7 · 7 · 7 · 7 · 7 · 7 = 49 · 49 · 49 · 7 = 2401 · 343 = 823543', '<p>\n12 км/ч = 12 · 1000 : 60 = 12000 : 60 = 200 (м/мин)<br/>\n200 м/мин ˃ 150 м/мин<br/>\n<b>Ответ:</b> Юра бежит быстрее Тараса.\n</p>\n\n\n<p>\n7 - 7 + 7 - 7 + 7 - 7 + 7 = 7<br/>\nНаибольшее значение: 7 · 7 · 7 · 7 · 7 · 7 · 7 = 49 · 49 · 49 · 7 = 2401 · 343 = 823543\n</p>\n\n<div class="img-wrapper-460">\n<img width="200" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica33-nomer13-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 33, номер 13-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 33, номер 13-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-33/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica33-nomer13-1.jpg', 'peterson/3/part3/page33/task12_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c32f4f19dd3d9a316936459e84d6864c7cf68367dc7f2afcca6d93c8c70e5f21', '12,150', NULL, 'тарас бежит со скоростью 150 м/мин, а юра-со скоростью 12 км/ч. кто из них бежит быстрее');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 34, '1', 0, 'Прочитай задачу и объясни, как составлена таблица. Составь план решения задачи и найди ответ. «Месяц назад 2 одинаковые порции мороженого стоили 36 р. Сейчас его цена увеличилась на 2 р. Сколько теперь надо заплатить за 5 таких порций мороженого?»', '</p> \n<p class="text">Прочитай задачу и объясни, как составлена таблица. Составь план решения задачи и найди ответ.<br/>\n«Месяц назад 2 одинаковые порции мороженого стоили 36 р. Сейчас его цена увеличилась на 2 р. Сколько теперь надо заплатить за 5 таких порций мороженого?»\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica34-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 34, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 34, номер 1, год 2022."/>\n</div>\n</div>', '36 : 2 = 18 (р.) – стоило мороженое. 18 + 2 = 20 (р.) – стоит мороженое сейчас. 20 · 5 = 100 (р.) – стоит 5 мороженных. Ответ: 100 рублей теперь надо заплатить за 5 таких порций мороженого.', '<p>\n36 : 2 = 18 (р.) – стоило мороженое.<br/>\n18 + 2 = 20 (р.) – стоит мороженое сейчас.<br/>\n20 · 5 = 100 (р.) – стоит 5 мороженных.<br/>\n<b>Ответ:</b> 100 рублей теперь надо заплатить за 5 таких порций мороженого.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-34/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica34-nomer1.jpg', 'peterson/3/part3/page34/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'dd9e58d6ad3144f0ab357b98f498c9a1d2a1810b3213c07fd786ba65dfd4fdf8', '2,5,36', '["найди"]'::jsonb, 'прочитай задачу и объясни, как составлена таблица. составь план решения задачи и найди ответ. "месяц назад 2 одинаковые порции мороженого стоили 36 р. сейчас его цена увеличилась на 2 р. сколько теперь надо заплатить за 5 таких порций мороженого?"');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 34, '2', 1, 'Реши задачи с помощью таблиц: а) Маша купила 7 заколок, а Вера – на 2 заколки меньше. Цена всех заколок одинаковая. Маша заплатила на 140 р. больше Веры. Сколько стоит одна заколка? Сколько рублей заплатила за заколки каждая из девочек? б) Саша и Дима купили вместе 20 солдатиков по одинаковой цене. Саша заплатил 720 р., а Дима – на 240 р. меньше. Сколько солдатиков купил каждый из них?', '</p> \n<p class="text">Реши задачи с помощью таблиц:<br/>\nа) Маша купила 7 заколок, а Вера – на 2 заколки меньше. Цена всех заколок одинаковая. Маша заплатила на 140 р. больше Веры. Сколько стоит одна заколка? Сколько рублей заплатила за заколки каждая из девочек?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="370" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica34-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 34, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 34, номер 2, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Саша и Дима купили вместе 20 солдатиков по одинаковой цене. Саша заплатил 720 р., а Дима – на 240 р. меньше. Сколько солдатиков купил каждый из них?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="370" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica34-nomer2-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 34, номер 2-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 34, номер 2-1, год 2022."/>\n</div>\n</div>', 'а) 140 : 2 = 70 (р.) – одна заколка 70 · 7 = 490 (р.) – Маша, 70 · (7 – 2) = 70 · 5 = 350 (р.) - Вера Ответ: 70 рублей стоит одна заколка. 490 рублей заплатила Маша и 350 рублей заплатила Вера за заколки. б) (720 + (720 - 240)) : 20 = (720 + 480) : 20 = 1200 : 20 = 60 (р.) – один солдатик 720 : 60 = 12 (солдатиков) – Саша (720 - 240) : 60 = 480 : 60 = 8 (солдатиков) – Дима Ответ: 12 солдатиков купил Саша и 8 солдатиков купил Дима.', '<p>\nа) 140 : 2 = 70 (р.) – одна заколка<br/>\n70 · 7 = 490 (р.) – Маша, <br/>\n70 · (7 – 2) = 70 · 5 = 350 (р.) - Вера<br/>\n<b>Ответ:</b> 70 рублей стоит одна заколка. 490 рублей заплатила Маша и 350 рублей заплатила Вера за заколки.<br/><br/>\nб) (720 + (720 - 240)) : 20 = (720 + 480) : 20 = 1200 : 20 = 60 (р.) – один солдатик<br/>\n720 : 60 = 12 (солдатиков) – Саша<br/>\n(720 - 240) : 60 = 480 : 60 = 8 (солдатиков) – Дима<br/>\n<b>Ответ:</b> 12 солдатиков купил Саша и 8 солдатиков купил Дима.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-34/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica34-nomer2.jpg', 'peterson/3/part3/page34/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica34-nomer2-1.jpg', 'peterson/3/part3/page34/task2_condition_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '11153af858444cde086049cfc2aeab22c66447e1df6fd9b2d5504940ab82cfac', '2,7,20,140,240,720', '["реши","больше","меньше"]'::jsonb, 'реши задачи с помощью таблиц:а) маша купила 7 заколок, а вера-на 2 заколки меньше. цена всех заколок одинаковая. маша заплатила на 140 р. больше веры. сколько стоит одна заколка? сколько рублей заплатила за заколки каждая из девочек? б) саша и дима купили вместе 20 солдатиков по одинаковой цене. саша заплатил 720 р., а дима-на 240 р. меньше. сколько солдатиков купил каждый из них');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 34, '3', 2, 'а) 9 пирожных, имеющих одну цену, стоят 234 р. Сколько рублей надо заплатить за 7 таких пирожных? б) Мама сначала купила 3 кг яблок по цене 40 р. за килограмм, а потом ещё 2 кг таких же яблок. Сколько денег она заплатила?', '</p> \n<p class="text">а) 9 пирожных, имеющих одну цену, стоят 234 р. Сколько рублей надо заплатить за 7 таких пирожных?<br/>\nб) Мама сначала купила 3 кг яблок по цене 40 р. за килограмм, а потом ещё 2 кг таких же яблок. Сколько денег она заплатила?\n</p>', 'а) 234 : 9 · 7 = 26 · 7 = 182 (р.) Ответ: 182 рубля надо заплатить за 7 таких пирожных. б) 3 · 40 + 2 · 40 = 120 + 80 = 200 (р.) Ответ: 200 рублей заплатила она.', '<p>\nа) 234 : 9 · 7 = 26 · 7 = 182 (р.)<br/>\n<b>Ответ:</b> 182 рубля надо заплатить за 7 таких пирожных. <br/><br/>\nб) 3 · 40 + 2 · 40 = 120 + 80 = 200 (р.)<br/>\n<b>Ответ:</b> 200 рублей заплатила она.\n</p>', '', '', TRUE, '[{"letter":"а","condition":"9 пирожных, имеющих одну цену, стоят 234 р. Сколько рублей надо заплатить за 7 таких пирожных?","solution":"234 : 9 · 7 = 26 · 7 = 182 (р.) Ответ: 182 рубля надо заплатить за 7 таких пирожных."},{"letter":"б","condition":"Мама сначала купила 3 кг яблок по цене 40 р. за килограмм, а потом ещё 2 кг таких же яблок. Сколько денег она заплатила?","solution":"3 · 40 + 2 · 40 = 120 + 80 = 200 (р.) Ответ: 200 рублей заплатила она."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-34/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3f2b2e884e3146fe645666b5e33ee1fe6a932753b1a5b6f593c798446f0d304a', '2,3,7,9,40,234', NULL, 'а) 9 пирожных, имеющих одну цену, стоят 234 р. сколько рублей надо заплатить за 7 таких пирожных? б) мама сначала купила 3 кг яблок по цене 40 р. за килограмм, а потом ещё 2 кг таких же яблок. сколько денег она заплатила');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 34, '4', 3, 'Для осенних посадок купили 60 пакетов луковиц тюльпанов по цене 15 р. за пакет, а нарциссов – на 25 пакетов меньше. Цена пакета нарциссов на 3 р. меньше, чем цена пакета тюльпанов. Сколько рублей надо заплатить за всю эту покупку?', '</p> \n<p class="text">Для осенних посадок купили 60 пакетов луковиц тюльпанов по цене 15 р. за пакет, а нарциссов – на 25 пакетов меньше. Цена пакета нарциссов на 3 р. меньше, чем цена пакета тюльпанов. Сколько рублей надо заплатить за всю эту покупку?</p>', '60 · 15 + (60 - 25) · 60 : (15 - 3) = 900 + 35 · 12 = 900 + 420 = 1320 (р.) Ответ: 1320 рублей надо заплатить за всю эту покупку.', '<p>\n60 · 15 + (60 - 25) · 60 : (15 - 3) = 900 + 35 · 12 = 900 + 420 = 1320 (р.)<br/>\n<b>Ответ:</b> 1320 рублей надо заплатить за всю эту покупку.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-34/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3180e5d8b7ec9c0df3053dc2ba83288afc9f1a1d0c9a6688113866f2df80b87c', '3,15,25,60', '["меньше"]'::jsonb, 'для осенних посадок купили 60 пакетов луковиц тюльпанов по цене 15 р. за пакет, а нарциссов-на 25 пакетов меньше. цена пакета нарциссов на 3 р. меньше, чем цена пакета тюльпанов. сколько рублей надо заплатить за всю эту покупку');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 35, '5', 0, 'Вычисли устно наиболее удобным способом: а) 126 + 99            д) 997 · 452 + 3 · 452 б) 532 - 98            е) 284 + 98 + 116 + 2 в) 20 · 142 · 5        ж) (939 + 56) - 239 г) 73 · 25 · 4          з) 721 - 96 - 621', '</p> \n<p class="text">Вычисли устно наиболее удобным способом:</p> \n\n<p class="description-text"> \nа) 126 + 99            д) 997 · 452 + 3 · 452<br/>\nб) 532 - 98            е) 284 + 98 + 116 + 2<br/>\nв) 20 · 142 · 5        ж) (939 + 56) - 239<br/>\nг) 73 · 25 · 4          з) 721 - 96 - 621\n</p>', 'а) 126 + 99 = 101 + 99 + 25 = 225 б) 532 - 98 = 432 - 100 - 98 = 432 - 2 = 430 в) 20 · 142 · 5 = 100 · 142 = 14200 г) 73 · 25 · 4 = 73 · 100 = 7300 д) 997 · 452 + 3 · 452 = 452 · (997 + 3) = 452 · 1000 = 452000 е) 284 + 98 + 116 + 2 = 100 + 400 = 500 ж) (939 + 56) - 239 = 939 - 239 + 56 = 700 + 56 = 756 з) 721 - 96 - 621 = 100 - 96 = 4', '<p>\nа) 126 + 99 = 101 + 99 + 25 = 225<br/>  		\nб) 532 - 98 = 432 - 100 - 98 = 432 - 2 = 430<br/>	\nв) 20 · 142 · 5 = 100 · 142 = 14200<br/>	\nг) 73 · 25 · 4 = 73 · 100 = 7300<br/>\nд) 997 · 452 + 3 · 452 = 452 · (997 + 3) = 452 · 1000 = 452000<br/> \nе) 284 + 98 + 116 + 2 = 100 + 400 = 500<br/>\nж) (939 + 56) - 239 = 939 - 239 + 56 = 700 + 56 = 756<br/>\nз) 721 - 96 - 621 = 100 - 96 = 4\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-35/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '81367ebf92f654a965a6fe912b40205cf4defb552a75b374df8d7c46c10697eb', '2,3,4,5,20,25,56,73,96,98', '["вычисли"]'::jsonb, 'вычисли устно наиболее удобным способом:а) 126+99            д) 997*452+3*452 б) 532-98            е) 284+98+116+2 в) 20*142*5        ж) (939+56)-239 г) 73*25*4          з) 721-96-621');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 35, '6', 1, 'Выполни умножение: а) 450 · 7600         б) 58000 · 4700 в) 20560 · 950       г) 69 · 300800', '</p> \n<p class="text">Выполни умножение:</p> \n\n<p class="description-text"> \nа) 450 · 7600         б) 58000 · 4700<br/> \nв) 20560 · 950       г) 69 · 300800\n</p>', 'а) 450 · 7600 = 3420000', '<p>\nа) 450 · 7600 = 3420000\n</p>\n\n<div class="img-wrapper-460">\n<img width="180" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica35-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 35, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 35, номер 6, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-35/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica35-nomer6.jpg', 'peterson/3/part3/page35/task6_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1f7aefb487ed06de853619a27177af9720d8b9383595c70c0222b54a7233146c', '69,450,950,4700,7600,20560,58000,300800', NULL, 'выполни умножение:а) 450*7600         б) 58000*4700 в) 20560*950       г) 69*300800');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 35, '7', 2, 'Реши уравнения с комментированием и сделай проверку: а) (980 : n) · 18 - 84 = 276 б) 96 + (80 - x) : 14 = 100', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (980 : n) · 18 - 84 = 276<br/>  \nб) 96 + (80 - x) : 14 = 100\n</p>', 'а) (980 : n) · 18 - 84 = 276 Чтобы найти уменьшаемое (980 : n) · 18 надо к разности прибавить вычитаемое (980 : n) · 18 = 276 + 84 (980 : n) · 18 = 360 Чтобы найти множитель (980 : n) надо произведение разделить на известный множитель (980 : n) = 360 : 18 (980 : n) = 20 Чтобы найти делитель n надо делимое разделить на частное n = 980 : 20 n = 49 Проверка: (980 : 49) · 18 - 84 = 276 б) 96 + (80 - x) : 14 = 100 Чтобы найти слагаемое (80 - x) : 14 надо из суммы отнять известное слагаемое (80 - x) : 14 = 100 - 96 (80 - x) : 14 = 4 Чтобы найти делимое (80 - x) надо делитель умножить на частное (80 - x) = 14 · 4 80 - х = 56 Чтобы найти вычитаемое надо отнять от уменьшаемого разность х = 80 - 56 х = 24 Проверка: 96 + (80 - 24) : 14 = 100', '<p>\nа) (980 : n) · 18 - 84 = 276<br/>  \nЧтобы найти уменьшаемое (980 : n) · 18 надо к разности прибавить вычитаемое<br/>\n(980 : n) · 18 = 276 + 84<br/>\n(980 : n) · 18 = 360<br/>\nЧтобы найти множитель (980 : n) надо произведение разделить на известный множитель<br/>\n(980 : n) = 360 : 18<br/>\n(980 : n) = 20<br/>\nЧтобы найти делитель n надо делимое разделить на частное<br/>\nn = 980 : 20<br/>\nn = 49<br/>\n<b>Проверка:</b> (980 : 49) · 18 - 84 = 276 <br/><br/>\n\nб) 96 + (80 - x) : 14 = 100<br/>\nЧтобы найти слагаемое (80 - x) : 14 надо из суммы отнять известное слагаемое<br/>\n(80 - x) : 14 = 100 - 96<br/>\n(80 - x) : 14 = 4<br/>\nЧтобы найти делимое (80 - x) надо делитель умножить на частное<br/>\n(80 - x) = 14 · 4<br/>\n80 - х = 56<br/>\nЧтобы найти вычитаемое надо отнять от уменьшаемого разность<br/>\nх = 80 - 56<br/>\nх = 24<br/>\n<b>Проверка:</b> 96 + (80 - 24) : 14 = 100\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-35/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b27b5a95ee9b729b5c8db5d0a161ad86b941c224e390273f4d13abed03d3476b', '14,18,80,84,96,100,276,980', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (980:n)*18-84=276 б) 96+(80-x):14=100');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 35, '8', 3, 'По рисунку найди делимое, делитель, частное и остаток. Запиши соотношение между ними с помощью формулы a = b · c + r, r < b. Проверь записанное равенство с помощью вычислений.', '</p> \n<p class="text">По рисунку найди делимое, делитель, частное и остаток. Запиши соотношение между ними с помощью формулы a = b · c + r, r &lt; b. Проверь записанное равенство с помощью вычислений.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica35-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 35, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 35, номер 8, год 2022."/>\n</div>\n</div>', 'а) 76 = 17 · 4 + 8 б) 81 = 26 · 3 + 3', '<p>\nа) 76 = 17 · 4 + 8<br/>\nб) 81 = 26 · 3 + 3\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-35/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica35-nomer8.jpg', 'peterson/3/part3/page35/task8_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'fcdff53221fa094372ac0c8702a3f1144d851a0e6c22e36a482574516f2c9149', NULL, '["найди","частное","делитель","делимое","остаток"]'::jsonb, 'по рисунку найди делимое, делитель, частное и остаток. запиши соотношение между ними с помощью формулы a=b*c+r, r<b. проверь записанное равенство с помощью вычислений');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 35, '9', 4, 'Запиши множество делителей и множество кратных числа 24.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 24.</p>', 'Множество делителей числа 24: 1, 2, 3, 4, 6, 8, 12, 24. Множество кратных числа 24: 24, 48, 72, 96, 120, ... .', '<p>\nМножество делителей числа 24: 1, 2, 3, 4, 6, 8, 12, 24.<br/>\nМножество кратных числа 24: 24, 48, 72, 96, 120, ... .\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-35/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '49d4214d9944c5101a33ce0982a51166a0ac441b3a22448faacac7f390919fac', '24', NULL, 'запиши множество делителей и множество кратных числа 24');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 35, '10', 5, 'Сравни в каждом равенстве числа, обозначенные буквами. Какое из них больше, а какое меньше? На сколько? n = m · 3      c · 10 = d      k : t = 2 a : b = 6      p : 5 = r        y = x : 8', '</p> \n<p class="text">Сравни в каждом равенстве числа, обозначенные буквами. Какое из них больше, а какое меньше? На сколько?</p> \n\n<p class="description-text"> \nn = m · 3      c · 10 = d      k : t = 2<br/>\na : b = 6      p : 5 = r        y = x : 8\n</p>', 'n = m · 3, n больше m в 3 раза c · 10 = d, c меньше d в 10 раз k : t = 2, k больше t в 2 раза a : b = 6, a больше b в 6 раз p : 5 = r, p больше r в 5 раз y = x : 8, y меньше x в 8 раз', '<p>\nn = m · 3, n больше m в 3 раза <br/> 	\nc · 10 = d, c меньше d в 10 раз<br/>	\nk : t = 2, k больше t в 2 раза<br/>\na : b = 6, a больше b в 6 раз <br/>	\np : 5 = r, p больше r в 5 раз<br/>	\ny = x : 8, y меньше x в 8 раз\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-35/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '11eeb152403ee5aa9b817650b82c63e52b4e583f0b8391fcea389221bdfdd1b1', '2,3,5,6,8,10', '["сравни","больше","меньше"]'::jsonb, 'сравни в каждом равенстве числа, обозначенные буквами. какое из них больше, а какое меньше? на сколько? n=m*3      c*10=d      k:t=2 a:b=6      p:5=r        y=x:8');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 35, '11', 6, 'Длина класса, имеющего форму прямоугольного параллелепипеда, равна 12 м, ширина – 8 м, а высота – 4 м. Найди объём этого класса, площадь его пола, потолка, стен.', '</p> \n<p class="text">Длина класса, имеющего форму прямоугольного параллелепипеда, равна 12 м, ширина – 8 м, а высота – 4 м. Найди объём этого класса, площадь его пола, потолка, стен.</p>', 'Длина класса, имеющего форму прямоугольного параллелепипеда, равна 12 м, ширина – 8 м, а высота – 4 м. V = a · b · h V = 12 · 8 · 4 = 384 (м 3 ) - объём этого класса, S = а · b S = 12 · 8 = 96 (м 2 ) - площадь его пола, S = 12 · 8 = 96 (м 2 ) - площадь его потолка, S = (12 · 4) · 2 + (8 · 4) · 2 = 48 · 2 + 36 · 2 = (48 + 36) · 2 = 84 · 2 = 168 (м 2 ) - площадь его стен. А, D и E', '<p>\nДлина класса, имеющего форму прямоугольного параллелепипеда, равна 12 м, ширина – 8 м, а высота – 4 м. <br/>\nV = a · b · h <br/>\nV = 12 · 8 · 4 = 384 (м<sup>3</sup>) - объём этого класса, <br/>\nS = а · b <br/>\nS = 12 · 8 = 96 (м<sup>2</sup>) - площадь его пола, <br/>\nS = 12 · 8 = 96 (м<sup>2</sup>) - площадь его потолка, <br/>\nS = (12 · 4) · 2 + (8 · 4) · 2 = 48 · 2 + 36 · 2 = (48 + 36) · 2 = 84 · 2 = 168 (м<sup>2</sup>) - площадь его стен.\n\n</p>\n\n\n<p>\nА, D и E\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-35/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '290b64df392b1218a516601b5815ba1cd6287523131c384f14cdcf0037012779', '4,8,12', '["найди","площадь"]'::jsonb, 'длина класса, имеющего форму прямоугольного параллелепипеда, равна 12 м, ширина-8 м, а высота-4 м. найди объём этого класса, площадь его пола, потолка, стен');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 36, '1', 0, 'Выполни действия: а) 916 · 73              б) 850 · 3800 в) 20900 · 9400     г) 60080 · 460', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 916 · 73              б) 850 · 3800 <br/> \nв) 20900 · 9400     г) 60080 · 460\n</p>', 'а) 916 · 73 = 66868', '<p>\nа) 916 · 73 = 66868\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica36-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 36, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 36, номер 1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-36/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica36-nomer1.jpg', 'peterson/3/part3/page36/task1_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f6de1db3fe75c544a8b5dbcc61c2e9542bac2e3648f3d580d24ac1ba0a5e27e2', '73,460,850,916,3800,9400,20900,60080', NULL, 'выполни действия:а) 916*73              б) 850*3800 в) 20900*9400     г) 60080*460');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 36, '2', 1, 'Придумай задачи по таблицам и реши их с помощью формулы стоимости:', '</p> \n<p class="text">Придумай задачи по таблицам и реши их с помощью формулы стоимости:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica36-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 36, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 36, номер 2, год 2022."/>\n</div>\n</div>', 'а) в первой корзине 5 яблок по 46 р. и во второй корзине 8 груш по 27 р. Сколько стоит первая, вторая корзины и на сколько отличается их стоимость? 46 · 5 - 27 · 8 = 230 - 216 = 14 (р.) Ответ: яблоки в первой корзине стоят 230 рублей, груши во второй корзине стоят 216 рублей и яблоки дороже груш на 14 рублей. б) в первом магазине было потрачено на покупку 6 альбомов 192 р., а во втором магазине за 4 альбома 384 р. Сколько стоит альбом в первом, во втором магазинах? Во сколько раз отличается стоимость альбома в магазинах? (384 : 4) : (192 : 6) = 96 : 32 = 3 (раза) Ответ: 96 рублей стоит альбом в первом магазине, 32 рубля стоит альбом во втором магазине, их стоимость различается в 3 раза.', '<p>\nа) в первой корзине 5 яблок по 46 р. и во второй корзине 8 груш по 27 р. Сколько стоит первая, вторая корзины и на сколько отличается их стоимость?<br/>\n46 · 5 - 27 · 8 = 230 - 216 = 14 (р.)<br/>\n<b>Ответ:</b> яблоки в первой корзине стоят 230 рублей, груши во второй корзине стоят 216 рублей и яблоки дороже груш на 14 рублей.<br/><br/>\n\nб) в первом магазине было потрачено на покупку 6 альбомов 192 р., а во втором магазине за 4 альбома 384 р. Сколько стоит альбом в первом, во втором магазинах? Во сколько раз отличается стоимость альбома в магазинах?<br/>\n(384 : 4) : (192 : 6) = 96 : 32 = 3 (раза)<br/>\n<b>Ответ:</b> 96 рублей стоит альбом в первом магазине, 32 рубля стоит альбом во втором магазине, их стоимость различается в 3 раза.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-36/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica36-nomer2.jpg', 'peterson/3/part3/page36/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '264d2cd8453095e92bbceddef9edfd603ec7c9ecadd9b5df6a9ea932b793e8d8', NULL, '["реши"]'::jsonb, 'придумай задачи по таблицам и реши их с помощью формулы стоимости');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 36, '3', 2, 'БЛИЦтурнир а) Мама купила 3 м шёлка по a р. за метр и 5 м ситца по b р. за метр. Сколько рублей она заплатила за всю покупку? б) Цена конфеты n р. Вадим купил 6 таких конфет, и у него ещё осталось t р. Сколько денег у него было вначале? в) Саше надо купить 7 бубликов по k р. за штуку. В кассу он отдал y р. Сколько сдачи он должен получить? г) Цена арбуза a р. за килограмм, а дыни – b р. за килограмм. На сколько рублей дыня массой в 5 кг дороже арбуза массой 6 кг?', '</p> \n<p class="text">БЛИЦтурнир<br/>\nа) Мама купила 3 м шёлка по a р. за метр и 5 м ситца по b р. за метр. Сколько рублей она заплатила за всю покупку?<br/>\nб) Цена конфеты n р. Вадим купил 6 таких конфет, и у него ещё осталось t р. Сколько денег у него было вначале?<br/>\nв) Саше надо купить 7 бубликов по k р. за штуку. В кассу он отдал y р. Сколько сдачи он должен получить?<br/>\nг) Цена арбуза a р. за килограмм, а дыни – b р. за килограмм. На сколько рублей дыня массой в 5 кг дороже арбуза массой 6 кг?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica36-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 36, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 36, номер 3, год 2022."/>\n</div>\n</div>', 'а) 3a + 5b (р.) Ответ: 3a + 5b рублей она заплатила за всю покупку. б) 6n + t (р.) Ответ: 6n + t рублей у него было вначале. в) y - 7k (р.) Ответ: y - 7k рублей сдачи он должен получить. г) 5b - 6а (р.) Ответ: на 5b - 6а рублей дыня массой в 5 кг дороже арбуза массой 6 кг.', '<p>\nа) 3a + 5b (р.)<br/>\n<b>Ответ:</b> 3a + 5b рублей она заплатила за всю покупку.<br/><br/>\nб) 6n + t (р.)<br/>\n<b>Ответ:</b> 6n + t рублей у него было вначале.<br/><br/>\nв) y - 7k (р.)<br/>\n<b>Ответ:</b> y - 7k рублей сдачи он должен получить.<br/><br/>\nг) 5b - 6а (р.) <br/>\n<b>Ответ:</b> на 5b - 6а рублей дыня массой в 5 кг дороже арбуза массой 6 кг.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-36/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica36-nomer3.jpg', 'peterson/3/part3/page36/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f70c1a19b2d632e8cb3b6053d4a6eabf08606773abd7aad56fefd20f89835801', '3,5,6,7', NULL, 'блицтурнир а) мама купила 3 м шёлка по a р. за метр и 5 м ситца по b р. за метр. сколько рублей она заплатила за всю покупку? б) цена конфеты n р. вадим купил 6 таких конфет, и у него ещё осталось t р. сколько денег у него было вначале? в) саше надо купить 7 бубликов по k р. за штуку. в кассу он отдал y р. сколько сдачи он должен получить? г) цена арбуза a р. за килограмм, а дыни-b р. за килограмм. на сколько рублей дыня массой в 5 кг дороже арбуза массой 6 кг');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 36, '4', 3, 'У Алёши в кошельке 6 монет по 5 р., две монеты по 10 р. и одна купюра 50 р. Он купил 3 тетради по цене 18 р., ластик за 12 р. и линейку за 19 р. На оставшиеся деньги он решил купить ластики. Сколько ластиков он сможет купить, если их цена 5 р. за штуку?', '</p> \n<p class="text">У Алёши в кошельке 6 монет по 5 р., две монеты по 10 р. и одна купюра 50 р. Он купил 3 тетради по цене 18 р., ластик за 12 р. и линейку за 19 р. На оставшиеся деньги он решил купить ластики. Сколько ластиков он сможет купить, если их цена 5 р. за штуку?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica36-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 36, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 36, номер 4, год 2022."/>\n</div>\n</div>', '((6 · 5 + 2 · 10 + 50) - (3 · 18 + 12 + 19)) : 5 = ((30 + 20 + 50) - (54 + 31)) : 5 = (100 - 85) : 5 = 15 : 5 = 3 (шт.) Ответ: 3 ластика он сможет купить, если их цена 5 р. за штуку.', '<p>\n((6 · 5 + 2 · 10 + 50) - (3 · 18 + 12 + 19)) : 5 = ((30 + 20 + 50) - (54 + 31)) : 5 = (100 - 85) : 5 = 15 : 5 = 3 (шт.)<br/>\n<b>Ответ:</b> 3 ластика он сможет купить, если их цена 5 р. за штуку.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-36/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica36-nomer4.jpg', 'peterson/3/part3/page36/task4_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '72a6a0a36cff6733895acc8c7f043b88503597a01d1176a1241786798c0eae91', '3,5,6,10,12,18,19,50', '["реши"]'::jsonb, 'у алёши в кошельке 6 монет по 5 р., две монеты по 10 р. и одна купюра 50 р. он купил 3 тетради по цене 18 р., ластик за 12 р. и линейку за 19 р. на оставшиеся деньги он решил купить ластики. сколько ластиков он сможет купить, если их цена 5 р. за штуку');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 36, '5', 4, 'Реши уравнения с комментированием и проверкой: а) (24 - 360 : x) · 6 = 90 б) 4 + (у - 14) : 3 = 20', '</p> \n<p class="text">Реши уравнения с комментированием и проверкой:</p>\n\n<p class="description-text"> \nа) (24 - 360 : x) · 6 = 90 <br/>      \nб) 4 + (у - 14) : 3 = 20\n</p>', 'а) (24 - 360 : x) · 6 = 90 Чтобы найти множитель (24 - 360 : x) надо произведение разделить на известный множитель (24 - 360 : x) = 90 : 6 (24 - 360 : x) = 15 Чтобы найти вычитаемое 360 : x надо из уменьшаемого отнять разность 360 : x = 24 - 15 360 : x = 9 Чтобы найти делитель надо делимое разделить на частное х = 360 : 9 х = 40 Проверка: (24 - 360 : 40) · 6 = 90 б) 4 + (у - 14) : 3 = 20 Чтобы найти слагаемое (у - 14) : 3 надо из суммы вычесть известное слагаемое (у - 14) : 3 = 20 - 4 (у - 14) : 3 = 16 Чтобы найти делимое (у - 14) надо делитель умножить на частное (у - 14) = 3 · 16 (у - 14) = 48 Чтобы найти уменьшаемое надо вычитаемое прибавить к разности у = 14 + 48 у = 62 Проверка: 4 + (62 - 14) : 3 = 20', '<p>\nа) (24 - 360 : x) · 6 = 90  <br/>   \nЧтобы найти множитель (24 - 360 : x) надо произведение разделить на известный множитель<br/>\n(24 - 360 : x) = 90 : 6<br/>\n(24 - 360 : x) = 15<br/>\nЧтобы найти вычитаемое 360 : x надо из уменьшаемого отнять разность<br/>\n360 : x = 24 - 15<br/>\n360 : x = 9<br/>\nЧтобы найти делитель надо делимое разделить на частное<br/>\nх = 360 : 9<br/>\nх = 40 <br/>\n<b>Проверка:</b> (24 - 360 : 40) · 6 = 90    <br/><br/>\nб) 4 + (у - 14) : 3 = 20<br/>\nЧтобы найти слагаемое (у - 14) : 3 надо из суммы вычесть известное слагаемое<br/>\n(у - 14) : 3 = 20 - 4<br/>\n(у - 14) : 3 = 16<br/>\nЧтобы найти делимое (у - 14) надо делитель умножить на частное<br/>\n(у - 14) = 3 · 16<br/>\n(у - 14) = 48<br/>\nЧтобы найти уменьшаемое надо вычитаемое прибавить к разности<br/>\nу = 14 + 48<br/>\nу = 62<br/>\n<b>Проверка:</b> 4 + (62 - 14) : 3 = 20\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-36/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '63a04e873fe99893a2f02c69f14ab06f8f2b885c9c20018c3c8e46c813e833c2', '3,4,6,14,20,24,90,360', '["реши"]'::jsonb, 'реши уравнения с комментированием и проверкой:а) (24-360:x)*6=90 б) 4+(у-14):3=20');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 37, '6', 0, 'Сравни в каждом равенстве числа, обозначенные буквами: n - 8 = d      a - k = 2      x · 5 = y p = t + 9      c : b = 8      r = m : 7', '</p> \n<p class="text">Сравни в каждом равенстве числа, обозначенные буквами:</p> \n\n<p class="description-text"> \nn - 8 = d      a - k = 2      x · 5 = y<br/>\np = t + 9      c : b = 8      r = m : 7\n</p>', 'n больше на 8, чем d     a больше, чем k на 2 p больше, чем t на 9      c больше, чем b в 8 раз x меньше в 5 раз, чем y r меньше, чем m в 7 раз', '<p>\nn больше на 8, чем d     a больше, чем k на 2<br/>		\np больше, чем t на 9      c больше, чем b в 8 раз<br/>		\nx меньше в 5 раз, чем y<br/>\nr меньше, чем m в 7 раз\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-37/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '540aea4f210ef48a6358bbf4c1d0a523dd5f1336b91ba92045e1bc1940bc0843', '2,5,7,8,9', '["сравни"]'::jsonb, 'сравни в каждом равенстве числа, обозначенные буквами:n-8=d      a-k=2      x*5=y p=t+9      c:b=8      r=m:7');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 37, '7', 1, 'Вырази величины в указанных единицах измерения:я а) 4 км 25 м = … м 4 м 25 см = … см 4 м 25 мм = … мм 4 м 2 25 дм 2 = … дм 2 4 дм 3 25 см 3 = … см 3 б) 4 ц 25 кг = … кг 4 т 25 кг = … кг 4 кг 25 г = … г 4 ч 25 мин = … мин 4 мин 25 с = … с', '</p> \n<p class="text">Вырази величины в указанных единицах измерения:я</p> \n\n<p class="description-text"> \nа) 4 км 25 м = … м<br/>  	\n4 м 25 см = … см <br/>    	\n4 м 25 мм = … мм <br/>   	\n4 м<sup>2</sup> 25 дм<sup>2</sup> = … дм<sup>2</sup>    <br/>	\n4 дм<sup>3</sup> 25 см<sup>3</sup> = … см<sup>3</sup> <br/><br/>   	\n\nб) 4 ц 25 кг = … кг<br/>\n4 т 25 кг = … кг<br/>\n4 кг 25 г = … г<br/>\n4 ч 25 мин = … мин<br/>\n4 мин 25 с = … с\n</p>', 'а) 4 км 25 м = 4025 м 4 м 25 см = 425 см 4 м 25 мм = 10025 мм 4 м 2 25 дм 2 = 425 дм 2 4 дм 3 25 см 3 = 4025 см 3 б) 4 ц 25 кг = 425 кг 4 т 25 кг = 4025 кг 4 кг 25 г = 4025 г 4 ч 25 мин = 265 мин 4 мин 25 с = 265 с', '<p>\nа) 4 км 25 м = 4025 м <br/> 		\n4 м 25 см = 425 см <br/>    		\n4 м 25 мм = 10025 мм <br/>   	\n4 м<sup>2</sup> 25 дм<sup>2</sup> = 425 дм<sup>2</sup> <br/>   		\n4 дм<sup>3</sup> 25 см<sup>3</sup> = 4025 см<sup>3</sup>   <br/><br/> 	\n\nб) 4 ц 25 кг = 425 кг<br/>\n4 т 25 кг = 4025 кг<br/>\n4 кг 25 г = 4025 г<br/>\n4 ч 25 мин = 265 мин<br/>\n4 мин 25 с = 265 с\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-37/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c718dd68eef20324bb6431a09d6e16c39621d3c10cd32cb3ff377fb1b810e789', '2,3,4,25', '["раз"]'::jsonb, 'вырази величины в указанных единицах измерения:я а) 4 км 25 м=... м 4 м 25 см=... см 4 м 25 мм=... мм 4 м 2 25 дм 2=... дм 2 4 дм 3 25 см 3=... см 3 б) 4 ц 25 кг=... кг 4 т 25 кг=... кг 4 кг 25 г=... г 4 ч 25 мин=... мин 4 мин 25 с=... с');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 37, '8', 2, 'Вычисли. Расположи ответы в порядке убывания и расшифруй слово. Найди в словаре, что оно означает. Припомни, а с тобой это случалось?', '</p> \n<p class="text">Вычисли. Расположи ответы в порядке убывания и расшифруй слово. Найди в словаре, что оно означает. Припомни, а с тобой это случалось?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica37-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 37, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 37, номер 8, год 2022."/>\n</div>\n</div>', 'Ё - 892 · 53 = 47276', '<p>\nЁ - 892 · 53 = 47276\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica37-nomer8-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 37, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 37, номер 8, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-37/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica37-nomer8.jpg', 'peterson/3/part3/page37/task8_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica37-nomer8-1.jpg', 'peterson/3/part3/page37/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '038b59bf8d493b8741f86eee9166c0c41eca6cd1d0a907eff26c7160af3e809a', NULL, '["вычисли","найди"]'::jsonb, 'вычисли. расположи ответы в порядке убывания и расшифруй слово. найди в словаре, что оно означает. припомни, а с тобой это случалось');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 37, '9', 3, 'Запиши множество делителей и множество кратных числа 25.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 25.</p>', 'Множество делителей числа 25: 1, 5, 25. Множество кратных числа 25: 25, 50, 75, 100, 125, 150, ....', '<p>\nМножество делителей числа 25: 1, 5, 25. Множество кратных числа 25: 25, 50, 75, 100, 125, 150, ....\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-37/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0caa32f9da227d89bc488c046a6b53d30904b23b8ad6b427e2c92d8be1c2634c', '25', NULL, 'запиши множество делителей и множество кратных числа 25');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 37, '10', 4, 'Набери указанную сумму денег наименьшим возможным числом монет и купюр. Составь и заполни таблицу в тетради.', '</p> \n<p class="text">Набери указанную сумму денег наименьшим возможным числом монет и купюр. Составь и заполни таблицу в тетради.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica37-nomer10.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 37, номер 10, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 37, номер 10, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica37-nomer10-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 37, номер 10-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 37, номер 10-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-37/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica37-nomer10.jpg', 'peterson/3/part3/page37/task10_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica37-nomer10-1.jpg', 'peterson/3/part3/page37/task10_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '42d9a0edc4bd731bb92e554e69f4b7a4b794bea5e87b13b1a27de7f06181bd3d', NULL, '["заполни"]'::jsonb, 'набери указанную сумму денег наименьшим возможным числом монет и купюр. составь и заполни таблицу в тетради');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 37, '11', 5, 'Подбери корни уравнений. Обоснуй свой ответ. а) x + x + x + x = 4 · 752 б) (y + 7) · 5 = 8 · 5 + 7 · 5', '</p> \n<p class="text">Подбери корни уравнений. Обоснуй свой ответ.</p> \n\n<p class="description-text"> \nа)  x + x + x + x = 4 · 752 <br/>             \nб) (y + 7) · 5 = 8 · 5 + 7 · 5\n</p>', 'а) 752 + 752 + 752 + 752 = 4 · 752 корень уравнения равен 752, только в этом случае равенство верно. б) (8 + 7) · 5 = 8 · 5 + 7 · 5 корень уравнения равен 8, используем правило умножения суммы на число. 12х · у : 3 12 : 3(х · у) 4(х · у) – произведение этих чисел увеличилось в 4 раза.', '<p>\nа) 752 + 752 + 752 + 752 = 4 · 752 корень уравнения равен 752, только в этом случае равенство верно.<br/>\nб) (8 + 7) · 5 = 8 · 5 + 7 · 5 корень уравнения равен 8, используем правило умножения суммы на число.\n</p>\n\n\n<p>\n12х · у : 3 <br/>\n12 : 3(х · у)<br/>\n4(х · у) – произведение этих чисел увеличилось в 4 раза.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-37/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'db5bec161722a980163fe25d9f586658d6c2df275979b68aa8f6697f9d518aae', '4,5,7,8,752', NULL, 'подбери корни уравнений. обоснуй свой ответ. а) x+x+x+x=4*752 б) (y+7)*5=8*5+7*5');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 38, '1', 0, 'а) Объясни по рисунку, как умножить число на сумму, и выполни умножение: б) Используя рисунок, объясни способ записи умножения на трёхзначное число в столбик:', '</p> \n<p class="text">а) Объясни по рисунку, как умножить число на сумму, и выполни умножение:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica38-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 38, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 38, номер 1, год 2022."/>\n</div>\n</div>\n\n<p class="text">б) Используя рисунок, объясни способ записи умножения на трёхзначное число в столбик:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica38-nomer1-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 38, номер 1-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 38, номер 1-1, год 2022."/>\n</div>\n</div>', 'а) а · (b + c + d) = a · b + a · c + a · d – чтобы умножить число на сумму надо это число умножить на каждое слагаемое суммы и сложить полученные значения 156 · 324 = 156 · (300 + 20 + 4) = 156 · 300 + 156 · 20 + 156 · 4 = 15600 + 3120 + 624 = 18720 + 624 = 19344 б) чтобы умножить число на трёхзначное число надо умножить число на единицы, потом на десятки, потом на сотни и сложить полученные числа. В записи суммы число десятков сдвигается на 1 разряд влево, а число сотен – на 2 разряда влево.', '<p>\nа) а · (b + c + d) = a · b + a · c + a · d – чтобы умножить число на сумму надо это число умножить на каждое слагаемое суммы и сложить полученные значения<br/>\n156 · 324 = 156 · (300 + 20 + 4) = 156 · 300 + 156 · 20 + 156 · 4 = 15600 + 3120 + 624 = 18720 + 624 = 19344<br/><br/>\nб) чтобы умножить число на трёхзначное число надо умножить число на единицы, потом на десятки, потом на сотни и сложить полученные числа. В записи суммы число десятков сдвигается на 1 разряд влево, а число сотен – на 2 разряда влево.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Объясни по рисунку, как умножить число на сумму, и выполни умножение:","solution":"а · (b + c + d) = a · b + a · c + a · d – чтобы умножить число на сумму надо это число умножить на каждое слагаемое суммы и сложить полученные значения 156 · 324 = 156 · (300 + 20 + 4) = 156 · 300 + 156 · 20 + 156 · 4 = 15600 + 3120 + 624 = 18720 + 624 = 19344"},{"letter":"б","condition":"Используя рисунок, объясни способ записи умножения на трёхзначное число в столбик:","solution":"чтобы умножить число на трёхзначное число надо умножить число на единицы, потом на десятки, потом на сотни и сложить полученные числа. В записи суммы число десятков сдвигается на 1 разряд влево, а число сотен – на 2 разряда влево."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-38/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica38-nomer1.jpg', 'peterson/3/part3/page38/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica38-nomer1-1.jpg', 'peterson/3/part3/page38/task1_condition_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'bcd1b0901a6704cb363976ffacc8528d21b3ee3a247a83b65a7bd79dd3338c2a', NULL, '["столбик"]'::jsonb, 'а) объясни по рисунку, как умножить число на сумму, и выполни умножение:б) используя рисунок, объясни способ записи умножения на трёхзначное число в столбик');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 38, '2', 1, 'В одной упаковке 248 ластиков. Сколько ластиков в 536 упаковках? Найди ответ в данной записи примера. Можно ли по этой записи определить, сколько ластиков в 6 упаковках, в 30 упаковках, в 500 упаковках, в 5360 упаковках?', '</p> \n<p class="text">В одной упаковке 248 ластиков. Сколько ластиков в 536 упаковках? Найди ответ в данной записи примера. Можно ли по этой записи определить, сколько ластиков в 6 упаковках, в 30 упаковках, в 500 упаковках, в 5360 упаковках?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica38-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 38, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 38, номер 2, год 2022."/>\n</div>\n</div>', '248 · 6 = 1488 (ластиков) 248 · 30 = 7440 (ластиков) 248 · 500 = 124000 (ластиков) 1488 + 7440 + 124000 = 8928 + 124000 = 132928 (ластиков) В 5360 упаковках: 132928 · 10 = 1329280 (ластиков) Ответ: 132928 ластиков в 536 упаковках. 1488 ластиков в 6 упаковка, 7440 ластиков в 30 упаковках, 124000 ластиков в 500 упаковках, 1329280 ластиков в 5360 упаковках.', '<p>\n248 · 6 = 1488 (ластиков)<br/>\n248 · 30 = 7440 (ластиков)<br/>\n248 · 500 = 124000 (ластиков)<br/>\n1488 + 7440 + 124000 = 8928 + 124000 = 132928 (ластиков)<br/>\nВ 5360 упаковках: 132928 · 10 = 1329280 (ластиков)<br/>\n<b>Ответ:</b> 132928 ластиков в 536 упаковках. 1488 ластиков в 6 упаковка, 7440 ластиков в 30 упаковках, 124000 ластиков в 500 упаковках, 1329280 ластиков в 5360 упаковках.\n\n</p>', 'Умножение многозначного числа на трёхзначное Чтобы умножить любое число на трёхзначное, можно умножить это число последовательно на единицы, десятки и сотни трёхзначного числа, а затем полученные произведения сложить. В записи суммы число десятков сдвигается на 1 разряд влево, а число сотен – на 2 разряда влево. Пример:', '<div class="recomended-block">\n<span class="title">Умножение многозначного числа на трёхзначное </span>\n<p>\nЧтобы умножить любое число на трёхзначное, можно умножить это число последовательно на единицы, десятки и сотни трёхзначного числа, а затем полученные произведения сложить.<br/>\nВ записи суммы число десятков сдвигается на 1 разряд влево, а число сотен – на 2 разряда влево.<br/><br/>\nПример:\n\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica38-spravka.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 38, справка, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 38, справка, год 2022."/>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-38/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica38-nomer2.jpg', 'peterson/3/part3/page38/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '17920400b66826d83e4d37333bb1c41c1ef3190fa8fb5e562d59fc014d237045', '6,30,248,500,536,5360', '["найди"]'::jsonb, 'в одной упаковке 248 ластиков. сколько ластиков в 536 упаковках? найди ответ в данной записи примера. можно ли по этой записи определить, сколько ластиков в 6 упаковках, в 30 упаковках, в 500 упаковках, в 5360 упаковках');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 39, '3', 0, 'Завод за один день выпускает 485 автомобилей. На какие вопросы можно ответить по данной записи примеров? Можно ли, не вычисляя, сказать, на сколько второе произведение больше первого?', '</p> \n<p class="text">Завод за один день выпускает 485 автомобилей. На какие вопросы можно ответить по данной записи примеров?<br/>\nМожно ли, не вычисляя, сказать, на сколько второе произведение больше первого?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 39, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 39, номер 3, год 2022."/>\n</div>\n</div>', 'В году 365 или 366 дней и поэтому можно сказать, что второе произведение больше первого на 485 автомобилей, так как оно больше на 1 день в котором завод выпускает 485 автомобилей.', '<p>\nВ году 365 или 366 дней и поэтому можно сказать, что второе произведение больше первого на 485 автомобилей, так как оно больше на 1 день в котором завод выпускает 485 автомобилей.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-39/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer3.jpg', 'peterson/3/part3/page39/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '57b8376a7562dc92fa324f11c7bd4fbde4ee6eea797205586a07a4df1e984a17', '485', '["произведение","больше"]'::jsonb, 'завод за один день выпускает 485 автомобилей. на какие вопросы можно ответить по данной записи примеров? можно ли, не вычисляя, сказать, на сколько второе произведение больше первого');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 39, '4', 1, 'Найди значения выражений: а) 752 · 128        в) 405 · 527 б) 246 · 496        г) 906 · 358 д) 1029 · 374      ж) 5007 · 716 е) 8503 · 982      з) 30209 · 245', '</p> \n<p class="text">Найди значения выражений:</p> \n\n<p class="description-text"> \nа) 752 · 128        в) 405 · 527<br/>  \nб) 246 · 496        г) 906 · 358<br/>  \nд) 1029 · 374      ж) 5007 · 716<br/>\nе) 8503 · 982      з) 30209 · 245\n</p>', 'а) 752 · 128 = 96256', '<p>\nа) 752 · 128 = 96256\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 39, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 39, номер 4, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-39/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer4.jpg', 'peterson/3/part3/page39/task4_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '017c820c8dcb516ee35860fe8fe180745bf6e69f1900b5683a70af26dfc6fb4b', '128,245,246,358,374,405,496,527,716,752', '["найди"]'::jsonb, 'найди значения выражений:а) 752*128        в) 405*527 б) 246*496        г) 906*358 д) 1029*374      ж) 5007*716 е) 8503*982      з) 30209*245');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 39, '5', 2, 'Вычисли. Расшифруй слово, расположив ответы примеров в порядке возрастания. Кто это? Найди информацию о нём в Интернете или энциклопедии.', '</p> \n<p class="text">Вычисли. Расшифруй слово, расположив ответы примеров в порядке возрастания. Кто это? Найди информацию о нём в Интернете или энциклопедии.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="380" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 39, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 39, номер 5, год 2022."/>\n</div>\n</div>', 'Р - 706 · 329 = 232274', '<p>\nР - 706 · 329 = 232274\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer5-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 39, номер 5-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 39, номер 5-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-39/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer5.jpg', 'peterson/3/part3/page39/task5_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer5-1.jpg', 'peterson/3/part3/page39/task5_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8e66b10e377cd891fdd7e66012fa621b24b1d62f05e768557d0f3c44c8b5fe99', NULL, '["вычисли","найди"]'::jsonb, 'вычисли. расшифруй слово, расположив ответы примеров в порядке возрастания. кто это? найди информацию о нём в интернете или энциклопедии');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 39, '6', 3, 'Реши уравнения и сделай проверку: а) 62 - (116 + x) : 5 = 34 б) 540 : (y · 3 - 60) = 6', '</p> \n<p class="text">Реши уравнения и сделай проверку:</p> \n\n<p class="description-text"> \nа) 62 - (116 + x) : 5 = 34<br/>     \nб) 540 : (y · 3 - 60) = 6\n</p>', 'а) 62 - (116 + x) : 5 = 34 (116 + x) : 5 = 62 - 34 (116 + x) : 5 = 28 116 + x = 28 · 5 116 + х = 140 х = 140 - 116 х = 24 Проверка: 62 - (116 + 24) : 5 = 34 б) 540 : (y · 3 - 60) = 6 y · 3 - 60 = 540 : 6 y · 3 - 60 = 90 y · 3 = 90 + 60 y · 3 = 150 у = 150 : 3 у = 50 Проверка: 540 : (50 · 3 - 60) = 6', '<p>\nа) 62 - (116 + x) : 5 = 34 <br/>   \n (116 + x) : 5 = 62 - 34<br/>\n(116 + x) : 5 = 28<br/>\n116 + x = 28 · 5<br/>\n116 + х = 140<br/>\nх = 140 - 116<br/>\nх = 24<br/>\n<b>Проверка:</b> 62 - (116 + 24) : 5 = 34  <br/><br/>  \nб) 540 : (y · 3 - 60) = 6<br/>\ny · 3 - 60 = 540 : 6<br/>\ny · 3 - 60 = 90<br/>\ny · 3 = 90 + 60<br/>\ny · 3 = 150<br/>\nу = 150 : 3<br/>\nу = 50<br/>\n<b>Проверка:</b> 540 : (50 · 3 - 60) = 6\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-39/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a8938de61313afcdc889bf3a453fafa759b78cf3af75acd730e44e0bb46bda69', '3,5,6,34,60,62,116,540', '["реши"]'::jsonb, 'реши уравнения и сделай проверку:а) 62-(116+x):5=34 б) 540:(y*3-60)=6');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 39, '7', 4, 'БЛИЦтурнир а) Олег съел n пирожков, а Саша – на 3 пирожка меньше. Во сколько раз меньше пирожков съел Саша, чем Олег? б) У Маши b марок, а у Гены в 5 раз меньше. Сколько марок у них вместе? в) Аня шла 2 ч со скоростью x км/ч, а Полина – 4 ч со скоростью y км/ч. На сколько километров больше прошла Полина, чем Аня? г) В вазе было c груш. Из неё взяли 6 раз по d груш. Сколько груш осталось в вазе? д) Три одинаковые конфеты стоят k р. Сколько рублей надо заплатить за 8 таких конфет?', '</p> \n<p class="text">БЛИЦтурнир<br/>\nа) Олег съел n пирожков, а Саша – на 3 пирожка меньше. Во сколько раз меньше пирожков съел Саша, чем Олег?<br/>\nб) У Маши b марок, а у Гены в 5 раз меньше. Сколько марок у них вместе?<br/>\nв) Аня шла 2 ч со скоростью x км/ч, а Полина – 4 ч со скоростью y км/ч. На сколько километров больше прошла Полина, чем Аня?<br/>\nг) В вазе было c груш. Из неё взяли 6 раз по d груш. Сколько груш осталось в вазе?<br/>\nд) Три одинаковые конфеты стоят k р. Сколько рублей надо заплатить за 8 таких конфет?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 39, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 39, номер 7, год 2022."/>\n</div>\n</div>', 'а) n : (n - 3) (раз) Ответ: во n : (n - 3) раз меньше пирожков съел Саша, чем Олег. б) b + b : 5 (марок) Ответ: b + b : 5 марок у них вместе. в) 4y - 2х (км) Ответ: на 4y - 2х километров больше прошла Полина, чем Аня. г) c - 6d (груш) Ответ: c - 6d груш осталось в вазе. д) k : 3 · 8 (р) Ответ: k : 3 · 8 рублей надо заплатить за 8 таких конфет.', '<p>\nа) n : (n - 3) (раз)<br/>\n<b>Ответ:</b> во n : (n - 3) раз меньше пирожков съел Саша, чем Олег.<br/><br/>\nб) b + b : 5 (марок)<br/>\n<b>Ответ:</b> b + b : 5 марок у них вместе.<br/><br/>\nв) 4y - 2х (км)<br/>\n<b>Ответ:</b> на 4y - 2х километров больше прошла Полина, чем Аня.<br/><br/>\nг) c - 6d (груш)<br/>\n<b>Ответ:</b> c - 6d груш осталось в вазе.<br/><br/>\nд) k : 3 · 8 (р)<br/>\n<b>Ответ:</b> k : 3 · 8 рублей надо заплатить за 8 таких конфет.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-39/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica39-nomer7.jpg', 'peterson/3/part3/page39/task7_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e5972de9d1ae79e7f0062bc8bc981f79ba3ebdd8521fe8c4afd6ae967c07a5f9', '2,3,4,5,6,8', '["больше","меньше","раз"]'::jsonb, 'блицтурнир а) олег съел n пирожков, а саша-на 3 пирожка меньше. во сколько раз меньше пирожков съел саша, чем олег? б) у маши b марок, а у гены в 5 раз меньше. сколько марок у них вместе? в) аня шла 2 ч со скоростью x км/ч, а полина-4 ч со скоростью y км/ч. на сколько километров больше прошла полина, чем аня? г) в вазе было c груш. из неё взяли 6 раз по d груш. сколько груш осталось в вазе? д) три одинаковые конфеты стоят k р. сколько рублей надо заплатить за 8 таких конфет');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 40, '8', 0, 'а) Поезд прошёл расстояние 560 км со скоростью 70 км/ч, а расстояние 240 км – со скоростью 60 км/ч. Сколько времени он был в пути? б) Для спортивного зала купили на 560 р. резиновые мячи по цене 70 р. за штуку и на 240 р. теннисные мячи по цене 60 р. за штуку. Сколько всего мячей купили? Что ты замечаешь? Придумай свою задачу с другими величинами, которая решается так же.', '</p> \n<p class="text">а) Поезд прошёл расстояние 560 км со скоростью 70 км/ч, а расстояние 240 км – со скоростью 60 км/ч. Сколько времени он был в пути?<br/>\nб) Для спортивного зала купили на 560 р. резиновые мячи по цене 70 р. за штуку и на 240 р. теннисные мячи по цене 60 р. за штуку. Сколько всего мячей купили?<br/>\nЧто ты замечаешь? Придумай свою задачу с другими величинами, которая решается так же.\n</p>', 'а) 560 : 70 + 240 : 60 = 8 + 4 = 12 (ч) Ответ: 12 часов он был в пути. б) 560 : 70 + 240 : 60 = 8 + 4 = 12 (штук) Ответ: 12 всего мячей купили. Заметно, что число результата одно и тоже. Своя задача с другими величинами, которая решается так же: На столе корзина с яблоками общей стоимостью 560 р. по цене 70 р. за яблоко и корзина с грушами общей стоимостью 240 р. по цене 60 р. за грушу. Сколько всего яблок и груш на столе? 560 : 70 + 240 : 60 = 8 + 4 = 12 (штук) Ответ: 12 всего яблок и груш на столе.', '<p>\nа) 560 : 70 + 240 : 60 = 8 + 4 = 12 (ч) <br/>\n<b>Ответ:</b> 12 часов он был в пути.<br/><br/>\nб) 560 : 70 + 240 : 60 = 8 + 4 = 12 (штук) <br/>\n<b>Ответ:</b> 12 всего мячей купили.<br/><br/>\nЗаметно, что число результата одно и тоже. Своя задача с другими величинами, которая решается так же: На столе корзина с яблоками общей стоимостью 560 р. по цене 70 р. за яблоко и корзина с грушами общей стоимостью 240 р. по цене 60 р. за грушу. Сколько всего яблок и груш на столе?<br/>\n560 : 70 + 240 : 60 = 8 + 4 = 12 (штук) <br/>\n<b>Ответ:</b> 12 всего яблок и груш на столе.\n\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Поезд прошёл расстояние 560 км со скоростью 70 км/ч, а расстояние 240 км – со скоростью 60 км/ч. Сколько времени он был в пути?","solution":"560 : 70 + 240 : 60 = 8 + 4 = 12 (ч) Ответ: 12 часов он был в пути."},{"letter":"б","condition":"Для спортивного зала купили на 560 р. резиновые мячи по цене 70 р. за штуку и на 240 р. теннисные мячи по цене 60 р. за штуку. Сколько всего мячей купили? Что ты замечаешь? Придумай свою задачу с другими величинами, которая решается так же.","solution":"560 : 70 + 240 : 60 = 8 + 4 = 12 (штук) Ответ: 12 всего мячей купили. Заметно, что число результата одно и тоже. Своя задача с другими величинами, которая решается так же: На столе корзина с яблоками общей стоимостью 560 р. по цене 70 р. за яблоко и корзина с грушами общей стоимостью 240 р. по цене 60 р. за грушу. Сколько всего яблок и груш на столе? 560 : 70 + 240 : 60 = 8 + 4 = 12 (штук) Ответ: 12 всего яблок и груш на столе."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-40/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3a28276a7f04ffb78c90d8749e854d48c0cc80d0b1235220379ee20a19180b00', '60,70,240,560', NULL, 'а) поезд прошёл расстояние 560 км со скоростью 70 км/ч, а расстояние 240 км-со скоростью 60 км/ч. сколько времени он был в пути? б) для спортивного зала купили на 560 р. резиновые мячи по цене 70 р. за штуку и на 240 р. теннисные мячи по цене 60 р. за штуку. сколько всего мячей купили? что ты замечаешь? придумай свою задачу с другими величинами, которая решается так же');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 40, '9', 1, 'Выполни вычисления по алгоритму, заданному блок-схемой. Составь и заполни таблицу в тетради.', '</p> \n<p class="text">Выполни вычисления по алгоритму, заданному блок-схемой. Составь и заполни таблицу в тетради.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica40-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 40, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 40, номер 9, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica40-nomer9-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 40, номер 9-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 40, номер 9-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-40/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica40-nomer9.jpg', 'peterson/3/part3/page40/task9_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica40-nomer9-1.jpg', 'peterson/3/part3/page40/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0a1e228842ee20dc80f198635f570b1460bbfe0696979b6b143b90d99811c9b9', NULL, '["заполни"]'::jsonb, 'выполни вычисления по алгоритму, заданному блок-схемой. составь и заполни таблицу в тетради');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 40, '10', 2, 'Сравни, не вычисляя: 352 · 218     218 · 352 920 · 614     614 · 920 516 · 724     724 · 521 306 · 825     294 · 438 368 : 8     368 : 23 504 : 56     672 : 56', '</p> \n<p class="text">Сравни, не вычисляя:</p> \n\n<p class="description-text"> \n352 · 218 <span class="okon">   </span> 218 · 352<br/>  	\n920 · 614 <span class="okon">   </span> 614 · 920  <br/>	   \n516 · 724 <span class="okon">   </span> 724 · 521  <br/><br/>  	\n\n\n306 · 825 <span class="okon">   </span> 294 · 438<br/>\n368 : 8 <span class="okon">   </span> 368 : 23<br/>\n504 : 56 <span class="okon">   </span> 672 : 56\n</p>', '352 · 218 = 218 · 352 920 · 614 = 614 · 920 516 · 724 ˂ 724 · 521 306 · 825 ˃ 294 · 438 368 : 8 ˂ 368 : 23 504 : 56 ˂ 672 : 56', '<p>\n352 · 218 =  218 · 352<br/>  	\n920 · 614 = 614 · 920  <br/>	   \n516 · 724 ˂ 724 · 521  <br/><br/>  	\n\n\n306 · 825 ˃ 294 · 438<br/>\n368 : 8 ˂ 368 : 23<br/>\n504 : 56 ˂ 672 : 56\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-40/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '036bd4a3f55328cfdf38ccc4135e03565b9358c9a5b2eeef5e3a4276a25c21b9', '8,23,56,218,294,306,352,368,438,504', '["сравни"]'::jsonb, 'сравни, не вычисляя:352*218     218*352 920*614     614*920 516*724     724*521 306*825     294*438 368:8     368:23 504:56     672:56');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 40, '11', 3, 'Повтори римскую нумерацию (ч. 1, с. 62). а) Запиши арабскими цифрами числа: VII, IX, XXIV, XLVI, CCCIV, DCCXII, MLVI. б) Запиши римскими цифрами числа: 4, 11, 36, 59, 93, 125, 408, 2002.', '</p> \n<p class="text">Повтори римскую нумерацию (ч. 1, с. 62).<br/>\nа) Запиши арабскими цифрами числа:<br/>\nVII, IX, XXIV, XLVI, CCCIV, DCCXII, MLVI.<br/>\nб) Запиши римскими цифрами числа: <br/>\n4, 11, 36, 59, 93, 125, 408, 2002.\n</p>', 'а) арабскими цифрами числа: VII = 7, IX = 9, XXIV = 24, XLVI = 46, CCCIV = 304, DCCXII = 712, MLVI = 1056. б) римскими цифрами числа: 4 = IV, 11 = XI, 36 = XXXVI, 59 = LIX, 93 = XCIII, 125 = CXXV, 408 = CDVIII, 2002 = MMII. 1749 год издания книги.', '<p>\nа) арабскими цифрами числа:<br/>\nVII = 7, IX = 9, XXIV = 24, XLVI = 46, CCCIV = 304, DCCXII = 712, MLVI = 1056.<br/><br/>\nб) римскими цифрами числа: <br/>\n4 = IV, 11 = XI, 36 = XXXVI, 59 = LIX, 93 = XCIII, 125 = CXXV, 408 = CDVIII, 2002 = MMII.\n</p>\n\n\n<p>\n1749 год издания книги.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-40/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'eb1fb757f155daf70ba6c1216a58bc634dea2c1ff2585de309e42013974d92d8', '1,4,11,36,59,62,93,125,408,2002', '["римск"]'::jsonb, 'повтори римскую нумерацию (ч. 1, с. 62). а) запиши арабскими цифрами числа:vii, ix, xxiv, xlvi, ccciv, dccxii, mlvi. б) запиши римскими цифрами числа:4, 11, 36, 59, 93, 125, 408, 2002');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 41, '1', 0, 'Рассмотри два способа умножения на трёхзначное число, в разряде десятков которого стоит 0. Чем отличаются эти способы? Почему в практике вычислений обычно используется второй способ?', '</p> \n<p class="text">Рассмотри два способа умножения на трёхзначное число, в разряде десятков которого стоит 0. Чем отличаются эти способы?<br/>\nПочему в практике вычислений обычно используется второй способ?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 41, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 41, номер 1, год 2022."/>\n</div>\n</div>', 'В первом способе записывается строчка с нулями и так точно не забыть почему сдвинулось следующее произведение на единицу. Второй способ удобнее. В нем не прописываются нули и число сразу сдвигается влево.', '<p>\nВ первом способе записывается строчка с нулями и так точно не забыть почему сдвинулось следующее произведение на единицу. <br/>\nВторой способ удобнее. В нем не прописываются нули и число сразу сдвигается влево.\n \n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-41/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer1.jpg', 'peterson/3/part3/page41/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '32f12e46327adf2adb03e0c74606b323213075ecc83dab547c3c0c1fb58f4c49', '0', '["раз"]'::jsonb, 'рассмотри два способа умножения на трёхзначное число, в разряде десятков которого стоит 0. чем отличаются эти способы? почему в практике вычислений обычно используется второй способ');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 41, '2', 1, 'Найди значения произведений: а) 963 · 407        в) 529 · 104 б) 216 · 809        г) 745 · 902 д) 807 · 307        ж) 402 · 609 е) 201 · 508        з) 905 · 106', '</p> \n<p class="text">Найди значения произведений:</p> \n\n<p class="description-text"> \nа) 963 · 407        в) 529 · 104<br/>  	\nб) 216 · 809        г) 745 · 902<br/>  	\nд) 807 · 307        ж) 402 · 609<br/>\nе) 201 · 508        з) 905 · 106\n\n</p>', 'а) 963 · 407 = 391941', '<p>\nа) 963 · 407 = 391941\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 41, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 41, номер 2, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-41/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer2.jpg', 'peterson/3/part3/page41/task2_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c4e1db5c7bbc2e807b94909c49fe2dc516c6c3c3f56a6ba59ccdfe339faaea79', '104,106,201,216,307,402,407,508,529,609', '["найди"]'::jsonb, 'найди значения произведений:а) 963*407        в) 529*104 б) 216*809        г) 745*902 д) 807*307        ж) 402*609 е) 201*508        з) 905*106');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 41, '3', 2, 'Вычисли. Расшифруй название старинной единицы объёма сыпучих тел во Франции. Узнай, скольким литрам она примерно равна.', '</p> \n<p class="text">Вычисли. Расшифруй название старинной единицы объёма сыпучих тел во Франции. Узнай, скольким литрам она примерно равна.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="250" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 41, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 41, номер 3, год 2022."/>\n</div>\n</div>', 'Д - 864 · 508 = 438912', '<p>\nД - 864 · 508 = 438912\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer3-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 41, номер 3-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 41, номер 3-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-41/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer3.jpg', 'peterson/3/part3/page41/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer3-1.jpg', 'peterson/3/part3/page41/task3_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3b13700a2e2761a0949f6c6af84cc0b5e1b8eaab4bed1d60fc6a46736ea766f1', NULL, '["вычисли"]'::jsonb, 'вычисли. расшифруй название старинной единицы объёма сыпучих тел во франции. узнай, скольким литрам она примерно равна');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 41, '4', 3, 'Найди значение выражения 527 · a, если a = 48, 250, 673, 901.', '</p> \n<p class="text">Найди значение выражения 527 · a, если a = 48, 250, 673, 901.</p>', '527 · 48 = 25296', '<p>\n527 · 48 = 25296\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 41, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 41, номер 4, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-41/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica41-nomer4.jpg', 'peterson/3/part3/page41/task4_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c12b05f938a98deb7df228efc09287413b9c74bfad5c9dcf513c7cf631dc6f1d', '48,250,527,673,901', '["найди"]'::jsonb, 'найди значение выражения 527*a, если a=48, 250, 673, 901');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 41, '5', 4, 'Реши задачи и сравни их. Что ты замечаешь? а) Слава бежал 3 мин со скоростью 200 м/мин. Затем он увеличил скорость на 40 м/мин и бежал ещё 2 мин. После этого ему осталось пробежать 120 м. Сколько всего метров надо пробежать Славе? б) Блузка стоит 200 р., а юбка – на 40 р. дороже. Нина купила 3 блузки и 2 юбки, и после этого у неё осталось 120 р. Сколько денег было у Нины вначале? Придумай ещё какую-нибудь задачу, которая решается так же.', '</p> \n<p class="text">Реши задачи и сравни их. Что ты замечаешь?<br/>\nа) Слава бежал 3 мин со скоростью 200 м/мин. Затем он увеличил скорость на 40 м/мин и бежал ещё 2 мин. После этого ему осталось пробежать 120 м. Сколько всего метров надо пробежать Славе?<br/>\nб) Блузка стоит 200 р., а юбка – на 40 р. дороже. Нина купила 3 блузки и 2 юбки, и после этого у неё осталось 120 р. Сколько денег было у Нины вначале? Придумай ещё какую-нибудь задачу, которая решается так же.\n</p>', 'а) 3 ·200 + (200 + 40) · 2 + 120 = 600 + 240 · 2 + 120 = 720 + 480 = 1200 (м) Ответ: 1200 всего метров надо пробежать Славе. б) 200 · 3 + (200 + 40) · 2 + 120 = 600 + 240 · 2 + 120 = 720 + 480 = 1200 (р) Ответ: 1200 рублей было у Нины вначале. Придумана ещё задача, которая решается так же: яблоки стоят 200 р. за килограмм, а груша – на 40 р. дороже. Мама купила 3 кг и 2 кг груш, и после этого у неё осталось 120 р. Сколько денег было у мамы вначале? 200 · 3 + (200 + 40) · 2 + 120 = 600 + 240 · 2 + 120 = 720 + 480 = 1200 (р) Ответ: 1200 рублей было у мамы вначале.', '<p>\nа) 3 ·200 + (200 + 40) · 2 + 120 = 600 + 240 · 2 + 120 = 720 + 480 = 1200 (м)<br/>\n<b>Ответ:</b> 1200 всего метров надо пробежать Славе.<br/><br/>\nб) 200 · 3 + (200 + 40) · 2 + 120 = 600 + 240 · 2 + 120 = 720 + 480 = 1200 (р)<br/>\n<b>Ответ:</b> 1200 рублей было у Нины вначале.<br/><br/>\nПридумана ещё задача, которая решается так же: яблоки стоят 200 р. за килограмм, а груша – на 40 р. дороже. Мама купила 3 кг и 2 кг груш, и после этого у неё осталось 120 р. Сколько денег было у мамы вначале?<br/>\n200 · 3 + (200 + 40) · 2 + 120 = 600 + 240 · 2 + 120 = 720 + 480 = 1200 (р)<br/>\n<b>Ответ:</b> 1200 рублей было у мамы вначале.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-41/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1a44dcfcbe7de548371e66a88213eafa4ad9130f899ebdca2d7241a451a7bcdd', '2,3,40,120,200', '["реши","сравни"]'::jsonb, 'реши задачи и сравни их. что ты замечаешь? а) слава бежал 3 мин со скоростью 200 м/мин. затем он увеличил скорость на 40 м/мин и бежал ещё 2 мин. после этого ему осталось пробежать 120 м. сколько всего метров надо пробежать славе? б) блузка стоит 200 р., а юбка-на 40 р. дороже. нина купила 3 блузки и 2 юбки, и после этого у неё осталось 120 р. сколько денег было у нины вначале? придумай ещё какую-нибудь задачу, которая решается так же');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 41, '6', 5, 'Запиши множество делителей и множество кратных числа 26.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 26.</p>', 'Множество делителей числа 26: 1, 2, 13, 26. Множество кратных числа 26: 26, 52, 78, 104, ....', '<p>\nМножество делителей числа 26: 1, 2, 13, 26.<br/>\nМножество кратных числа 26: 26, 52, 78, 104, ....\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-41/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5a723a6cffecec74a0b4128603af1a18ad384e8e4e942365584669ddb53d7ca3', '26', NULL, 'запиши множество делителей и множество кратных числа 26');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 42, '7', 0, 'Выполни действия: а) 4 дм 5 см + 3 м 7 см                 г) 8 т 96 кг - 429 кг б) 5 км 32 м + 4 км 756 м            д) 6 ч 32 мин + 19 ч 58 мин в) 7 дм 2 6 см 2 + 18 дм 2 68 см 2     е) 40 мин 2 с - 34 мин 25 с', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 4 дм 5 см + 3 м 7 см                 г) 8 т 96 кг - 429 кг<br/>\nб) 5 км 32 м + 4 км 756 м            д) 6 ч 32 мин + 19 ч 58 мин<br/>\nв) 7 дм<sup>2</sup> 6 см<sup>2</sup> + 18 дм<sup>2</sup> 68 см<sup>2</sup>    е) 40 мин 2 с - 34 мин 25 с\n\n</p>', 'а) 4 дм 5 см + 3 м 7 см = 45 см + 307 см = 352 см = 3 м 5 дм 2 см б) 5 км 32 м + 4 км 756 м = 50032 м + 40756 м = 90788 м = 9 км 788 м в) 7 дм 2 6 см 2 + 18 дм 2 68 см 2 = 706 см 2 + 1868 см 2 = 2574 см 2 = 25 дм 2 74 см 2 г) 8 т 96 кг - 429 кг = 8098 кг - 429 кг = 7669 кг = 7 т 669 кг д) 6 ч 32 мин + 19 ч 58 мин = (6 ч + 19 ч) + (32 мин + 58 мин) = 25 ч + 1 ч 30 мин = 26 ч 30 мин е) 40 мин 2 с - 34 мин 25 с = (39 мин - 34 мин) + (62 с - 25 с) = 5 мин 37 с', '<p>\nа) 4 дм 5 см + 3 м 7 см = 45 см + 307 см = 352 см = 3 м 5 дм 2 см<br/>			\nб) 5 км 32 м + 4 км 756 м = 50032 м + 40756 м = 90788 м = 9 км 788 м<br/>		\nв) 7 дм<sup>2</sup> 6 см<sup>2</sup> + 18 дм<sup>2</sup> 68 см<sup>2</sup> = 706 см<sup>2</sup> + 1868 см<sup>2</sup> = 2574 см<sup>2</sup> = 25 дм<sup>2</sup> 74 см<sup>2</sup>	<br/>\nг) 8 т 96 кг - 429 кг = 8098 кг - 429 кг = 7669 кг = 7 т 669 кг<br/>\nд) 6 ч 32 мин + 19 ч 58 мин = (6 ч + 19 ч) + (32 мин + 58 мин) = 25 ч + 1 ч 30 мин = 26 ч 30 мин<br/>\nе) 40 мин 2 с - 34 мин 25 с = (39 мин - 34 мин) + (62 с - 25 с) = 5 мин 37 с\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-42/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '27db3cc2116fbdfa57a8cd6b6d9ac6af2dc53946250c4849da539ea765052b7a', '2,3,4,5,6,7,8,18,19,25', NULL, 'выполни действия:а) 4 дм 5 см+3 м 7 см                 г) 8 т 96 кг-429 кг б) 5 км 32 м+4 км 756 м            д) 6 ч 32 мин+19 ч 58 мин в) 7 дм 2 6 см 2+18 дм 2 68 см 2     е) 40 мин 2 с-34 мин 25 с');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 42, '8', 1, 'Игра «Кто какое число задумал?» а) Зайка – попрыгайка задумал число, прибавил его к числу 26, сумму умножил на 5 и из полученного произведения вычел 42. В результате у него получилось 138. Какое число задумал Зайка – попрыгайка? б) Мышка – норушка вычла задуманное число из 31, разность разделила на 9 и к полученному результату прибавила 8. В ответе у неё получилось 11. Какое число задумала Мышка – норушка? в) Лягушка – квакушка разделила 250 на задуманное число, вычла из частного 24 и разность умножила на 2. В результате у неё получилось 52. Какое число задумала Лягушка – квакушка?', '</p> \n<p class="text">Игра «Кто какое число задумал?»<br/>\nа) Зайка – попрыгайка задумал число, прибавил его к числу 26, сумму умножил на 5 и из полученного произведения вычел 42. В результате у него получилось 138. Какое число задумал Зайка – попрыгайка?<br/>\nб) Мышка – норушка вычла задуманное число из 31, разность разделила на 9 и к полученному результату прибавила 8. В ответе у неё получилось 11. Какое число задумала Мышка – норушка? <br/>\nв) Лягушка – квакушка разделила 250 на задуманное число, вычла из частного 24 и разность умножила на 2. В результате у неё получилось 52. Какое число задумала Лягушка – квакушка?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 42, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 42, номер 8, год 2022."/>\n</div>\n</div>', 'а) (х + 26) · 5 - 42 = 138 (х + 26) · 5 = 138 + 42 (х + 26) · 5 = 180 х + 26 = 180 : 5 х + 26 = 36 х = 36 - 26 х = 10 Ответ: число 10 задумал Зайка – попрыгайка. б) (31 - х) : 9 + 8 = 11 (31 - х) : 9 = 11 - 8 (31 - х) : 9 = 3 31 - х = 3 · 9 31 - х = 27 х = 31 - 27 х = 4 Ответ: число 4 задумала Мышка – норушка. в) (250 : х - 24) · 2 = 52 250 : х - 24 = 52 : 2 250 : х - 24 = 26 250 : х = 26 + 24 250 : х = 50 х = 250 : 50 х = 5 Ответ: число 5 задумала Лягушка – квакушка.', '<p>\nа) (х + 26) · 5 - 42 = 138<br/>\n(х + 26) · 5 = 138 + 42<br/>\n(х + 26) · 5 = 180<br/>\nх + 26 = 180 : 5<br/>\nх + 26 = 36<br/>\nх = 36 - 26<br/>\nх = 10<br/>\n<b>Ответ:</b> число 10 задумал Зайка – попрыгайка.<br/><br/>\nб) (31 - х) : 9 + 8 = 11<br/>\n(31 - х) : 9 = 11 - 8<br/>\n(31 - х) : 9 = 3<br/>\n31 - х = 3 · 9<br/>\n31 - х = 27<br/>\nх = 31 - 27<br/>\nх = 4<br/>\n<b>Ответ:</b> число 4 задумала Мышка – норушка.<br/><br/>\nв) (250 : х - 24) · 2 = 52<br/>\n250 : х - 24 = 52 : 2<br/>\n250 : х - 24 = 26<br/>\n250 : х = 26 + 24<br/>\n250 : х = 50<br/>\nх = 250 : 50<br/>\nх = 5<br/>\n<b>Ответ:</b> число 5 задумала Лягушка – квакушка.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-42/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer8.jpg', 'peterson/3/part3/page42/task8_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2aefce4919b957d410097d6f9c89ef08c5cc224f89f69bda9e51a9d53924fe63', '2,5,8,9,11,24,26,31,42,52', '["раздели","разность","раз"]'::jsonb, 'игра "кто какое число задумал?" а) зайка-попрыгайка задумал число, прибавил его к числу 26, сумму умножил на 5 и из полученного произведения вычел 42. в результате у него получилось 138. какое число задумал зайка-попрыгайка? б) мышка-норушка вычла задуманное число из 31, разность разделила на 9 и к полученному результату прибавила 8. в ответе у неё получилось 11. какое число задумала мышка-норушка? в) лягушка-квакушка разделила 250 на задуманное число, вычла из частного 24 и разность умножила на 2. в результате у неё получилось 52. какое число задумала лягушка-квакушка');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 42, '9', 2, 'Составь программу действий и вычисли: 72 · 480 + 789 · 295 - (34188 + 392012) : 100', '</p> \n<p class="text">Составь программу действий и вычисли: </p> \n\n<p class="description-text"> \n72 · 480 + 789 · 295 - (34188 + 392012) : 100\n</p>', '72 · 480 + 789 · 295 – (34188 + 392012) : 100 = 263053 34188 + 392012 = 426200', '<p>\n72 · 480 + 789 · 295 – (34188 + 392012) : 100 = 263053<br/>\n34188 + 392012 = 426200\n</p>\n\n<div class="img-wrapper-460">\n<img width="160" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 42, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 42, номер 9, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-42/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer9.jpg', 'peterson/3/part3/page42/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '44ee74985df3c95c725c7d6f1a4ac8d0ca87fd3031bf89343c4d9419ce34366e', '72,100,295,480,789,34188,392012', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:72*480+789*295-(34188+392012):100');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 42, '10', 3, 'Начерти пятиугольник ABCDE и проведи прямую l так, чтобы она разбила пятиугольник: а) на треугольник и шестиугольник; б) на треугольник и пятиугольник; в) на четырёхугольник и пятиугольник; г) на два четырёхугольника.', '</p> \n<p class="text">Начерти пятиугольник ABCDE и проведи прямую l так, чтобы  она разбила пятиугольник: а) на треугольник и шестиугольник; б) на треугольник и пятиугольник; в) на четырёхугольник и пятиугольник; г) на два четырёхугольника.</p>', '', '<div class="img-wrapper-460">\n<img width="300" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer10.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 42, номер 10, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 42, номер 10, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-42/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer10.jpg', 'peterson/3/part3/page42/task10_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '07102924893b998d47f3e9285f5c4a0d8e7563e8622448397c274803ad6d447e', NULL, '["раз"]'::jsonb, 'начерти пятиугольник abcde и проведи прямую l так, чтобы она разбила пятиугольник:а) на треугольник и шестиугольник; б) на треугольник и пятиугольник; в) на четырёхугольник и пятиугольник; г) на два четырёхугольника');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 42, '11', 4, 'Какой из прямоугольных параллелепипедов, изображённых на рисунке, вместительнее?', '</p> \n<p class="text">Какой из прямоугольных параллелепипедов, изображённых на рисунке, вместительнее?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="300" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer11.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 42, номер 11, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 42, номер 11, год 2022."/>\n</div>\n</div>', '87 · 56 · 43 = 4872 · 43 = 209496 (cм 3 ) 64 · 64 · 64 = 4096 · 64 = 262144 Равно потому что половина трети равно 1/6, а треть половины тоже 1/6. Ответ: равны. 20 - 3 = 17 – осталось 7 кусочков 17 - 3 = 14 – осталось 6 кусочков 14 - 3 = 11 – осталось 5 кусочков 11 - 3 = 8 – осталось 4 кусочка 8 - 3 = 5 – осталось 3 кусочка 5 - 3 = 2 – осталось 2 кусочка Ответ: 6 кусочков разрезал Дима.', '<p>\n87 · 56 · 43 = 4872 · 43 = 209496 (cм<sup>3</sup>)\n</p>\n\n<div class="img-wrapper-460">\n<img width="270" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer11-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 42, номер 11-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 42, номер 11-1, год 2022."/>\n\n\n<p>\n64 · 64 · 64 = 4096 · 64 = 262144\n</p>\n\n<div class="img-wrapper-460">\n<img width="270" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer12.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 42, номер 12, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 42, номер 12, год 2022."/>\n\n\n<p>\nРавно потому что половина трети равно 1/6, а треть половины тоже 1/6.<br/>\n<b>Ответ:</b> равны.\n</p>\n\n\n<p>\n20 - 3 = 17 – осталось 7 кусочков<br/>\n17 - 3 = 14 – осталось 6 кусочков<br/>\n14 - 3 = 11 – осталось 5 кусочков<br/>\n11 - 3 = 8 – осталось 4 кусочка<br/>\n8 - 3 = 5 – осталось 3 кусочка<br/>\n5 - 3 = 2 – осталось 2 кусочка<br/>\n<b>Ответ:</b> 6 кусочков разрезал Дима.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-42/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer11.jpg', 'peterson/3/part3/page42/task11_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer11-1.jpg', 'peterson/3/part3/page42/task11_solution_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica42-nomer12.jpg', 'peterson/3/part3/page42/task11_solution_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '9aabd6515f30a21d4c8979ffb57452a19e16dcad7ba72f4f4a08e7f497afcdb5', NULL, NULL, 'какой из прямоугольных параллелепипедов, изображённых на рисунке, вместительнее');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 43, '1', 0, 'Найди ошибки в записи и решении примеров: Запиши и реши их в тетради правильно.', '</p> \n<p class="text">Найди ошибки в записи и решении примеров:<br/>\nЗапиши и реши их в тетради правильно.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 43, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 43, номер 1, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-460">\n<img width="190" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer1-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 43, номер 1-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 43, номер 1-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-43/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer1.jpg', 'peterson/3/part3/page43/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer1-1.jpg', 'peterson/3/part3/page43/task1_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8731b0fd547fa7b6209d2d90e16c90fd7914464d95b7b6b805cf8cd106773f4a', NULL, '["найди","реши"]'::jsonb, 'найди ошибки в записи и решении примеров:запиши и реши их в тетради правильно');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 43, '2', 1, 'Выполни действия: а) 318 · 956        б) 729 · 304 в) 407 · 501         г) 60080 · 264', '</p> \n<p class="text">\nВыполни действия:\n</p> \n\n<p class="description-text"> \nа) 318 · 956        б) 729 · 304 <br/>  \nв) 407 · 501         г) 60080 · 264\n</p>', 'а) 318 · 956 = 304008', '<p>\nа) 318 · 956 = 304008\n</p>\n\n<div class="img-wrapper-460">\n<img width="160" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 43, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 43, номер 2, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-43/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer2.jpg', 'peterson/3/part3/page43/task2_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'fa5712fa812b1edbb33783c480220b5ed013d17e12295fd41d7f12db1604dcae', '264,304,318,407,501,729,956,60080', NULL, 'выполни действия:а) 318*956        б) 729*304 в) 407*501         г) 60080*264');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 43, '3', 2, 'Запиши формулу стоимости. Используя её, заполни таблицу:', '</p> \n<p class="text">\nЗапиши формулу стоимости. Используя её, заполни таблицу:\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 43, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 43, номер 3, год 2022."/>\n</div>\n</div>', 'C = a · n C = 92 · 6 C = 552', '<p>\nC = a · n<br/>\nC = 92 · 6<br/>\nC = 552\n</p>\n\n<div class="img-wrapper-460">\n<img width="70" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer3-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 43, номер 3-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 43, номер 3-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-43/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer3.jpg', 'peterson/3/part3/page43/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer3-1.jpg', 'peterson/3/part3/page43/task3_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd08c009a5d0293a6207fa7c7d6dc2c51e8159ad459901f9e5700fd2b26b9405d', NULL, '["заполни"]'::jsonb, 'запиши формулу стоимости. используя её, заполни таблицу');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 43, '4', 3, 'Для букета купили 5 роз и 6 гербер. Каждая роза стоит 45 р., а гербера – в 3 раза дешевле. Сколько рублей стоит весь букет?', '</p> \n<p class="text">Для букета купили 5 роз и 6 гербер. Каждая роза стоит 45 р., а гербера – в 3 раза дешевле. Сколько рублей стоит весь букет?</p>', '5 · 45 + 6 · 45 : 3 = 225 + 6 · 15 = 225 + 90 = 315 (р.) Ответ: 315 рублей стоит весь букет.', '<p>\n5 · 45 + 6 · 45 : 3 =  225 + 6 · 15 = 225 + 90 = 315 (р.)<br/>\n<b>Ответ:</b> 315 рублей стоит весь букет.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-43/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c665ccca4240e58815bbcbd6a5884f3ebe8fdcbbc51970c657ef2e491c639f19', '3,5,6,45', '["раз","раза"]'::jsonb, 'для букета купили 5 роз и 6 гербер. каждая роза стоит 45 р., а гербера-в 3 раза дешевле. сколько рублей стоит весь букет');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 43, '5', 4, 'Стол и 4 одинаковых стула стоят 2800 р. За стол заплатили 1200 р. Сколько рублей стоит один стул?', '</p> \n<p class="text">Стол и 4 одинаковых стула стоят 2800 р. За стол заплатили 1200 р. Сколько рублей стоит один стул?</p>', '(2800 - 1200) : 4 = 1600 : 4 = 400 (р.) Ответ: 400 рублей стоит один стул.', '<p>\n(2800 - 1200) : 4 = 1600 : 4 = 400 (р.)<br/>\n<b>Ответ:</b> 400 рублей стоит один стул.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-43/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a2d376d2444d26135dc348e708a4a301ab67cd92302db875e5d7bff3f5cb1bd8', '4,1200,2800', NULL, 'стол и 4 одинаковых стула стоят 2800 р. за стол заплатили 1200 р. сколько рублей стоит один стул');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 43, '6', 5, 'а) Составь выражение к задаче: «Расстояние от Бреста до Киева примерно 600 км. Поезд ехал из Бреста в Киев сначала 2 ч со скоростью 60 км/ч, а потом t ч со скоростью 80 км/ч. Сколько километров ему осталось проехать?» Найди значение выражения при t = 1, 2, 3, 4, 5, 6. Может ли t принять значение, равное 10? б) Пусть d км – расстояние, оставшееся до Киева. Заполни таблицу и составь формулу зависимости d от t:', '</p> \n<p class="text">а) Составь выражение к задаче: «Расстояние от Бреста до Киева примерно 600 км. Поезд ехал из Бреста в Киев сначала 2 ч со скоростью 60 км/ч, а потом t ч со скоростью 80 км/ч. Сколько километров ему осталось проехать?» Найди значение выражения при t = 1, 2, 3, 4, 5, 6. Может ли t принять значение, равное 10?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 43, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 43, номер 6, год 2022."/>\n</div>\n</div>\n\n<p class="text">б) Пусть d км – расстояние, оставшееся до Киева. Заполни таблицу и составь формулу зависимости d от t:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer6-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 43, номер 6-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 43, номер 6-1, год 2022."/>\n</div>\n</div>', 'а) 600 - (2 · 60 + t · 80) при t = 1 600 - (2 · 60 + 1 · 80) = 600 - (120 + 80) = 600 - 200 = 400 (км) Ответ: 400 километров ему осталось проехать. при t = 2 600 - (2 · 60 + 2 · 80) = 600 - (120 + 160) = 600 - 280 = 320 (км) Ответ: 320 километров ему осталось проехать. при t = 3 600 - (2 · 60 + 3 · 80) = 600 - (120 + 240) = 600 - 360 = 240 (км) Ответ: 240 километров ему осталось проехать. при t = 4 600 - (2 · 60 + 4 · 80) = 600 - (120 + 320) = 600 - 440 = 160 (км) Ответ: 160 километров ему осталось проехать. при t = 5 600 - (2 · 60 + 5 · 80) = 600 - (120 + 400) = 600 - 520 = 80 (км) Ответ: 80 километров ему осталось проехать. при t = 6 600 - (2 · 60 + 6 · 80) = 600 - (120 + 480) = 600 - 600 = 0 (км) Ответ: 0 километров ему осталось проехать. При t = 10 600 - (2 · 60 + 10 · 80) = 600 - (120 + 800) = 600 - 920 Ответ: не может t принять значение, равное 10. Б) 600 - (2 · 60 + 0 · 80) = 600 - (120 + 0) = 600 - 120 = 480 (км)', '<p>\nа) 600 - (2 · 60 + t · 80)<br/>\nпри t = 1<br/>\n600 - (2 · 60 + 1 · 80) = 600 - (120 + 80) = 600 - 200 = 400 (км)<br/>\n<b>Ответ:</b> 400 километров ему осталось проехать. <br/><br/>\nпри t = 2<br/>\n600 - (2 · 60 + 2 · 80) = 600 - (120 + 160) = 600 - 280 = 320 (км)<br/>\n<b>Ответ:</b> 320 километров ему осталось проехать. <br/><br/>\nпри t = 3<br/>\n600 - (2 · 60 + 3 · 80) = 600 - (120 + 240) = 600 - 360 = 240 (км)<br/>\n<b>Ответ:</b> 240 километров ему осталось проехать. <br/><br/>\nпри t = 4<br/>\n600 - (2 · 60 + 4 · 80) = 600 - (120 + 320) = 600 - 440 = 160 (км)<br/>\n<b>Ответ:</b> 160 километров ему осталось проехать. <br/><br/>\nпри t = 5<br/>\n600 - (2 · 60 + 5 · 80) = 600 - (120 + 400) = 600 - 520 = 80 (км)<br/>\n<b>Ответ:</b> 80 километров ему осталось проехать. <br/><br/>\nпри t = 6<br/>\n600 - (2 · 60 + 6 · 80) = 600 - (120 + 480) = 600 - 600 = 0 (км)<br/>\n<b>Ответ:</b> 0 километров ему осталось проехать. <br/><br/>\nПри t = 10<br/>\n600 - (2 · 60 + 10 · 80) = 600 - (120 + 800) = 600 - 920<br/>\n<b>Ответ:</b> не может t принять значение, равное 10.<br/><br/>\n\n\nБ)<br/>\n600 - (2 · 60 + 0 · 80) = 600 - (120 + 0) = 600 - 120 = 480 (км)\n\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer6-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 43, номер 6-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 43, номер 6-2, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Составь выражение к задаче: «Расстояние от Бреста до Киева примерно 600 км. Поезд ехал из Бреста в Киев сначала 2 ч со скоростью 60 км/ч, а потом t ч со скоростью 80 км/ч. Сколько километров ему осталось проехать?» Найди значение выражения при t = 1, 2, 3, 4, 5, 6. Может ли t принять значение, равное 10?","solution":"600 - (2 · 60 + t · 80) при t = 1 600 - (2 · 60 + 1 · 80) = 600 - (120 + 80) = 600 - 200 = 400 (км) Ответ: 400 километров ему осталось проехать. при t = 2 600 - (2 · 60 + 2 · 80) = 600 - (120 + 160) = 600 - 280 = 320 (км) Ответ: 320 километров ему осталось проехать. при t = 3 600 - (2 · 60 + 3 · 80) = 600 - (120 + 240) = 600 - 360 = 240 (км) Ответ: 240 километров ему осталось проехать. при t = 4 600 - (2 · 60 + 4 · 80) = 600 - (120 + 320) = 600 - 440 = 160 (км) Ответ: 160 километров ему осталось проехать. при t = 5 600 - (2 · 60 + 5 · 80) = 600 - (120 + 400) = 600 - 520 = 80 (км) Ответ: 80 километров ему осталось проехать. при t = 6 600 - (2 · 60 + 6 · 80) = 600 - (120 + 480) = 600 - 600 = 0 (км) Ответ: 0 километров ему осталось проехать. При t = 10 600 - (2 · 60 + 10 · 80) = 600 - (120 + 800) = 600 - 920 Ответ: не может t принять значение, равное 10. Б) 600 - (2 · 60 + 0 · 80) = 600 - (120 + 0) = 600 - 120 = 480 (км)"},{"letter":"б","condition":"Пусть d км – расстояние, оставшееся до Киева. Заполни таблицу и составь формулу зависимости d от t:","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-43/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer6.jpg', 'peterson/3/part3/page43/task6_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer6-1.jpg', 'peterson/3/part3/page43/task6_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica43-nomer6-2.jpg', 'peterson/3/part3/page43/task6_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e141beecfa9afa19860afd6ebee768fc7060989f8f162ed118353f3edd2f41dc', '1,2,3,4,5,6,10,60,80,600', '["найди","заполни","равно"]'::jsonb, 'а) составь выражение к задаче:"расстояние от бреста до киева примерно 600 км. поезд ехал из бреста в киев сначала 2 ч со скоростью 60 км/ч, а потом t ч со скоростью 80 км/ч. сколько километров ему осталось проехать?" найди значение выражения при t=1, 2, 3, 4, 5, 6. может ли t принять значение, равное 10? б) пусть d км-расстояние, оставшееся до киева. заполни таблицу и составь формулу зависимости d от t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 44, '7', 0, 'Автобус проехал 180 км за 4 часа, а обратный путь – на 1 час быстрее. На сколько километров в час увеличилась скорость автобуса на обратном пути?', '</p> \n<p class="text">Автобус проехал 180 км за 4 часа, а обратный путь – на 1 час быстрее. На сколько километров в час увеличилась скорость автобуса на обратном пути?</p>', '180 : (4 - 1) - 180 : 4 = 180 : 3 - 45 = 60 - 45 = 15 (км/ч) Ответ: на 15 километров в час увеличилась скорость автобуса на обратном пути.', '<p>\n180 : (4 - 1) - 180 : 4 = 180 : 3 - 45 = 60 - 45 = 15 (км/ч)<br/> \n<b>Ответ:</b> на 15 километров в час увеличилась скорость автобуса на обратном пути.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-44/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '9a47778cce92725a679299c971a0444236d59898752a3edb6cab79e8a06e0f27', '1,4,180', NULL, 'автобус проехал 180 км за 4 часа, а обратный путь-на 1 час быстрее. на сколько километров в час увеличилась скорость автобуса на обратном пути');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 44, '8', 1, '(Устно.) Подбери корни уравнений или объясни, почему их нет. Сделай проверку: 7 + x = 7      n - 0 = 7      a - a = 7 7 - y = 0      t - 7 = 0      b - b = 0', '</p> \n<p class="text">(Устно.) Подбери корни уравнений или объясни, почему их нет. Сделай проверку:</p> \n\n<p class="description-text"> \n7 + x = 7      n - 0 = 7      a - a = 7<br/> \n7 - y = 0      t - 7 = 0      b - b = 0\n\n</p>', '7 + x = 7 х = 7 - 7 х = 0 Проверка: 7 + 0 = 7 7 - y = 0 у = 7 - 0 у = 7 Проверка: 7 - 7 = 0 n - 0 = 7 n = 7 + 0 n = 7 Проверка: 7 - 0 = 7 t - 7 = 0 t = 0 + 7 t = 7 Проверка: 7 - 7 = 0 a - a = 7 разность одинаковых чисел равна 0, поэтому нет корней b - b = 0 b имеет множество корней.', '<p>\n7 + x = 7<br/> \nх = 7 - 7<br/> \nх = 0<br/> \n<b>Проверка:</b> 7 + 0 = 7<br/> <br/> \n7 - y = 0  <br/> \nу = 7 - 0<br/> \nу = 7<br/> \n<b>Проверка:</b> 7 - 7 = 0<br/> <br/> \nn - 0 = 7 <br/>  \nn = 7 + 0<br/> \nn = 7<br/> \n<b>Проверка:</b> 7 - 0 = 7<br/> <br/> \nt - 7 = 0 <br/>  \nt = 0 + 7<br/> \nt = 7<br/> \n<b>Проверка:</b> 7 - 7 = 0<br/> <br/> \na - a = 7<br/> \nразность одинаковых чисел равна 0, поэтому нет корней<br/> \nb - b = 0<br/> \nb имеет множество корней.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-44/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2afa36291dea8d39e41edf8d3095e0b446bd947b312399a693ba840d82eabe05', '0,7', NULL, '(устно.) подбери корни уравнений или объясни, почему их нет. сделай проверку:7+x=7      n-0=7      a-a=7 7-y=0      t-7=0      b-b=0');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 44, '9', 2, 'БЛИЦтурнир а) За 5 банок краски заплатили k р., а за 9 банок лака – n р. На сколько рублей банка краски дороже банки лака? б) Три рюкзака стоят а р., а две палатки – на b р. дороже. На сколько рублей рюкзак дешевле палатки? в) Пешеход прошёл d км за 4 часа. Скорость велосипедиста – на m км/ч больше. С какой скоростью ехал велосипедист? г) Лодка проплыла s км за 5 ч, а катер это же расстояние – за 2 ч. На сколько километров в час скорость катера больше скорости лодки?', '</p> \n<p class="text">БЛИЦтурнир<br/>\nа) За 5 банок краски заплатили k р., а за 9 банок лака – n р. На сколько рублей банка краски дороже банки лака?<br/>\nб) Три рюкзака стоят а р., а две палатки – на b р. дороже. На сколько рублей рюкзак дешевле палатки? <br/>\nв) Пешеход прошёл d км за 4 часа. Скорость велосипедиста – на m км/ч больше. С какой скоростью ехал велосипедист?<br/>\nг) Лодка проплыла s км за 5 ч, а катер это же расстояние – за 2 ч. На сколько километров в час скорость катера больше скорости лодки?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica44-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 44, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 44, номер 9, год 2022."/>\n</div>\n</div>', 'а) k : 5 - n : 9 (р.) Ответ: на k : 5 - n : 9 рублей банка краски дороже банки лака. б) (а + b) : 2 – а : 3 (р.) Ответ: на (а + b) : 2 - а : 3 рублей рюкзак дешевле палатки. в) d : 4 + m (км/ч) Ответ: с d : 4 + m километров в час ехал велосипедист. г) s : 2 - s : 5 (км/ч) Ответ: на s : 2 - s : 5 километров в час скорость катера больше скорости лодки.', '<p>\nа) k : 5 - n : 9 (р.)<br/> \n<b>Ответ:</b> на k : 5 - n : 9 рублей банка краски дороже банки лака.<br/> <br/> \nб) (а + b) : 2 – а : 3 (р.)<br/> \n<b>Ответ:</b> на (а + b) : 2 - а : 3 рублей рюкзак дешевле палатки.<br/> <br/>  \nв) d : 4 + m (км/ч)<br/> \n<b>Ответ:</b> с d : 4 + m километров в час ехал велосипедист.<br/> <br/> \nг) s : 2 - s : 5 (км/ч)<br/> \n<b>Ответ:</b> на s : 2 - s : 5 километров в час скорость катера больше скорости лодки.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-44/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica44-nomer9.jpg', 'peterson/3/part3/page44/task9_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5ac6822f36e1417ef958a24dd80c2cc85fafc9ec4f3285c14db4bfe4a35d2538', '2,4,5,9', '["больше"]'::jsonb, 'блицтурнир а) за 5 банок краски заплатили k р., а за 9 банок лака-n р. на сколько рублей банка краски дороже банки лака? б) три рюкзака стоят а р., а две палатки-на b р. дороже. на сколько рублей рюкзак дешевле палатки? в) пешеход прошёл d км за 4 часа. скорость велосипедиста-на m км/ч больше. с какой скоростью ехал велосипедист? г) лодка проплыла s км за 5 ч, а катер это же расстояние-за 2 ч. на сколько километров в час скорость катера больше скорости лодки');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 44, '10', 3, 'Выполни действия: а) 985468 + 45032    в) 8000 · 8090      д) 121212 · 350 б) 507000 - 92944    г) 4905600 : 70    е) 795 · 270', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 985468 + 45032    в) 8000 · 8090      д) 121212 · 350<br/>\nб) 507000 - 92944    г) 4905600 : 70    е) 795 · 270\n\n</p>', 'а) 985468 + 45032 = 1030500', '<p>\nа) 985468 + 45032 = 1030500\n</p>\n\n<div class="img-wrapper-460">\n<img width="170" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica44-nomer10.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 44, номер 10, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 44, номер 10, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-44/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica44-nomer10.jpg', 'peterson/3/part3/page44/task10_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2cfef9786d4350cc0fd6eb4c79012c9e4d6ddad15b914c72ecbc2ec43bc9a1de', '70,270,350,795,8000,8090,45032,92944,121212,507000', NULL, 'выполни действия:а) 985468+45032    в) 8000*8090      д) 121212*350 б) 507000-92944    г) 4905600:70    е) 795*270');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 44, '11', 4, 'Какие точки на рисунке принадлежат прямой l, а какие – не принадлежат? Запиши в тетради, используя знаки ∈ и ∉.', '</p> \n<p class="text">Какие точки на рисунке принадлежат прямой l, а какие – не принадлежат? Запиши в тетради, используя знаки ∈ и ∉.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica44-nomer11.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 44, номер 11, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 44, номер 11, год 2022."/>\n</div>\n</div>', 'A ∉ l          D ∈ l B ∈ l          E ∉ l C ∉ l          K ∈ l A = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18} B = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23,24, 25, 26, 27} A ∩ B = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18} наибольший общий делитель чисел 18 и 27 : 18.', '<p>\nA ∉ l          D ∈ l<br/>\nB ∈ l          E ∉ l<br/>\nC ∉ l          K ∈ l\n</p>\n\n\n<p>\nA = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}<br/>\nB = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23,24, 25, 26, 27}<br/>\nA ∩ B = {1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18}<br/>\nнаибольший общий делитель чисел 18 и 27 : 18.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-44/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica44-nomer11.jpg', 'peterson/3/part3/page44/task11_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6acf40966f65621fafb3b37dd801cea974b04a49ddd13c7a0512b8fe5c49b62d', NULL, NULL, 'какие точки на рисунке принадлежат прямой l, а какие-не принадлежат? запиши в тетради, используя знаки ∈ и ∉');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 45, '1', 0, 'Вале и Гале было поручено сделать флажки для ёлки. Валя сделала за 2 часа 40 флажков, а Галя за 3 часа – 45 флажков. Кто из девочек сделал больше флажков, а кто – меньше? Кто работал больше времени, а кто – меньше? Кто работал быстрее, а кто – медленнее? Какие величины характеризуют работу? Как они связаны между собой?', '</p> \n<p class="text">Вале и Гале было поручено сделать флажки для ёлки. Валя сделала за 2 часа 40 флажков, а Галя за 3 часа – 45 флажков.<br/>\nКто из девочек сделал больше флажков, а кто – меньше? Кто работал больше времени, а кто – меньше? Кто работал быстрее, а кто – медленнее?<br/>\nКакие величины характеризуют работу? Как они связаны между собой?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica45-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 45, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 45, номер 1, год 2022."/>\n</div>\n\n</div>', 'Валя сделала меньше флажков. Галя работал больше времени, а Валя меньше. 40 : 2 = 20 (флажков/ч) – Валя 45 : 3 = 15 (флажков/ч) – Галя Валя работала быстрее, а Галя медленнее. Работа – это скорость работы. При выполнении работы нас интересует объём работы (работа) – сколько всего сделано, время работы и производительность – быстрее или медленнее она выполнялась. A = w · t Это равенство называют формулой работы. Оно означает, что работа равна производительности, умноженной на время работы.', '<p>\nВаля сделала меньше флажков. Галя работал больше времени, а Валя меньше. <br/>\n40 : 2 = 20 (флажков/ч) – Валя<br/>\n45 : 3 = 15 (флажков/ч) – Галя<br/>\nВаля работала быстрее, а Галя медленнее.<br/>\nРабота – это скорость работы. При выполнении работы нас интересует объём работы (работа) – сколько всего сделано, время работы и производительность – быстрее или медленнее она выполнялась.<br/>\nA = w · t<br/>\nЭто равенство называют формулой работы. Оно означает, что работа равна производительности, умноженной на время работы.\n\n</p>', 'Формула работы При выполнении работы нас интересует объём работы (работа) – сколько всего сделано, время работы и производительность – быстрее или медленнее она выполнялась. Производительность – это работа, выполненная за единицу времени. Или, другими словами, это «скорость работы». При решении задач на работу мы будем считать, что производительность не меняется. Задача: Строитель уложил 120 кирпичей за 10 мин. С какой производительностью он работал? Решение: 120 : 10 = 12 кирпичей в минуту. Производительность является величиной. В качестве единиц её измерения используют такие единицы, как штуки в минуту (шт./мин), тонны в час (т/ч), литры в секунду (л/с) и т. д. Если обозначить всю выполненную работу буквой A, производительность – буквой w, а время работы – буквой t, то можно записать равенство:', '<div class="recomended-block">\n<span class="title">Формула работы</span>\n<p>\nПри выполнении работы нас интересует <b class="black">объём работы (работа)</b> – сколько всего сделано, <b class="black">время</b> работы и <b class="black">производительность</b> – быстрее или медленнее она выполнялась.<br/>\n<b class="black">Производительность</b> – это работа, выполненная за единицу времени. Или, другими словами, это «скорость работы».<br/>\nПри решении задач на работу мы будем считать, что производительность не меняется.\n<br/>\n</p>\n\n\n<span class="title">Задача:</span>\n\n<p>\nСтроитель уложил 120 кирпичей за 10 мин. С какой производительностью он работал?<br/>\nРешение: <br/>\n120 : 10 = 12 кирпичей в минуту.<br/>\nПроизводительность является величиной. В качестве единиц её измерения используют такие единицы, как штуки в минуту (шт./мин), тонны в час (т/ч), литры в секунду (л/с) и т. д.<br/>\nЕсли обозначить всю выполненную работу буквой A, производительность – буквой w, а время работы – буквой t, то можно записать равенство:\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica45-spravka.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 45, справка, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 45, справка, год 2022."/>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-45/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica45-nomer1.jpg', 'peterson/3/part3/page45/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'eccce810675255c8770146f455461f3da5796ebcd265d12311f84ab978be0564', '2,3,40,45', '["больше","меньше"]'::jsonb, 'вале и гале было поручено сделать флажки для ёлки. валя сделала за 2 часа 40 флажков, а галя за 3 часа-45 флажков. кто из девочек сделал больше флажков, а кто-меньше? кто работал больше времени, а кто-меньше? кто работал быстрее, а кто-медленнее? какие величины характеризуют работу? как они связаны между собой');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 46, '2', 0, 'Объясни смысл предложений: а) Оля лепит пельмени с производительностью 2 штуки в минуту. б) Денис делает табуретки с производительностью 4 табуретки в день. в) Гена копает картошку с производительностью 3 ведра в час. г) Ира печатает текст с производительностью 120 знаков в минуту. Как из формулы работы найти производительность? Как найти время работы?', '</p> \n<p class="text">Объясни смысл предложений:<br/>\nа) Оля лепит пельмени с производительностью 2 штуки в минуту.<br/>\nб) Денис делает табуретки с производительностью 4 табуретки в день.<br/>\nв) Гена копает картошку с производительностью 3 ведра в час.<br/>\nг) Ира печатает текст с производительностью 120 знаков в минуту.<br/>\nКак из формулы работы найти производительность? Как найти время работы?\n</p>', 'а) Оля успевает слепить 2 пельменя за минуту. б) Денис успевает сделать 4 табуретки за день. в) Гена успевает выкопать 3 ведра картошки в час. г) Ира успевает напечатать 120 знаков текста в минуту. A = w · t Это равенство называют формулой работы. Оно означает, что работа равна производительности, умноженной на время работы. Поэтому из формулы работы можно найти производительность: w = A : t. Поэтому можно найти время работы: t = A : w.', '<p>\nа) Оля успевает слепить 2 пельменя за минуту.<br/>\nб) Денис успевает сделать 4 табуретки за день.<br/>\nв) Гена успевает выкопать 3 ведра картошки в час.<br/>\nг) Ира успевает напечатать 120 знаков текста в минуту.<br/><br/>\nA = w · t<br/>\nЭто равенство называют формулой работы. Оно означает, что работа равна производительности, умноженной на время работы. <br/>\nПоэтому из формулы работы можно найти производительность: w = A : t.<br/>\nПоэтому можно найти время работы: t = A : w.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-46/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a153fef4af8acabe32150f97431b75889e5c41b6f3462c5be344859e7d3465f0', '2,3,4,120', NULL, 'объясни смысл предложений:а) оля лепит пельмени с производительностью 2 штуки в минуту. б) денис делает табуретки с производительностью 4 табуретки в день. в) гена копает картошку с производительностью 3 ведра в час. г) ира печатает текст с производительностью 120 знаков в минуту. как из формулы работы найти производительность? как найти время работы');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 46, '3', 1, 'Найди неизвестные значения величин по формуле работы А = w · t:', '</p> \n<p class="text">Найди неизвестные значения величин по формуле работы А = w · t:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica46-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 46, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 46, номер 3, год 2022."/>\n</div>\n</div>', 'а) t = 60 : 4 t = 15(ч) A = 8 · 20 А = 160 (л) w = 450 : 15 w = 30 (шт/с) б) w = 240 : 8 w = 30 (зн/мин) A = 12 · 4 А = 48 (шт.) t = 480 : 80 t = 6 (ч)', '<p>\nа) t = 60 : 4<br/>\nt = 15(ч) <br/>\nA = 8 · 20<br/>\nА = 160 (л)<br/>\nw = 450 : 15<br/>\nw = 30 (шт/с)<br/><br/>\nб) w = 240 : 8<br/>\nw = 30 (зн/мин)<br/>\nA = 12 · 4<br/>\nА = 48 (шт.)<br/>\nt = 480 : 80<br/>\nt = 6 (ч)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-46/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica46-nomer3.jpg', 'peterson/3/part3/page46/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0d6bb74ef000f6ed1b6b3982f54843ae090c305a4f355c746c6b705ef123c3ef', NULL, '["найди"]'::jsonb, 'найди неизвестные значения величин по формуле работы а=w*t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 46, '4', 2, 'а) Завод выпускает 208 автомобилей в день. Сколько автомобилей выпустит завод в год? (Считать, что в году 256 рабочих дней.) б) Автомат закрыл 10800 банок за 6 ч. С какой производительностью он работает?', '</p> \n<p class="text">а) Завод выпускает 208 автомобилей в день. Сколько автомобилей выпустит завод в год? (Считать, что в году 256 рабочих дней.)<br/>\nб) Автомат закрыл 10800 банок за 6 ч. С какой производительностью он работает?\n</p>', 'а) 208 · 256 = 53248 (автомобилей)', '<p>\nа) 208 · 256 = 53248 (автомобилей)\n</p>\n\n<div class="img-wrapper-460">\n<img width="120" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica46-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 46, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 46, номер 4, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Завод выпускает 208 автомобилей в день. Сколько автомобилей выпустит завод в год? (Считать, что в году 256 рабочих дней.)","solution":"208 · 256 = 53248 (автомобилей)"},{"letter":"б","condition":"Автомат закрыл 10800 банок за 6 ч. С какой производительностью он работает?","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-46/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica46-nomer4.jpg', 'peterson/3/part3/page46/task4_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1a93b591fce32dff80a55da0825aa67727197ffebb1caf8097ca01e1b2c36f18', '6,208,256,10800', NULL, 'а) завод выпускает 208 автомобилей в день. сколько автомобилей выпустит завод в год? (считать, что в году 256 рабочих дней.) б) автомат закрыл 10800 банок за 6 ч. с какой производительностью он работает');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 46, '5', 3, 'Мастер вытачивает 8 деталей в час. Сколько деталей он сделает за 2 ч, 4 ч, 6 ч, 7 ч, 9 ч, t ч? Заполни в тетради таблицу. Запиши формулу зависимости работы A, выполненной мастером, от времени работы t.', '</p> \n<p class="text">Мастер вытачивает 8 деталей в час. Сколько деталей он сделает за 2 ч, 4 ч, 6 ч, 7 ч, 9 ч, t ч? Заполни в тетради таблицу. Запиши формулу зависимости работы A, выполненной мастером, от времени работы t. </p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica46-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 46, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 46, номер 5, год 2022."/>\n</div>\n</div>', 'A = w · t 8 · 2 = 16 (дет.) 8 · 4 = 32 (дет.) 8 · 6 = 48 (дет.) 8 · 7 = 56 (дет.) 8 · 9 = 72 (дет.)', '<p>\nA = w · t<br/>\n8 · 2 = 16 (дет.)<br/>\n8 · 4 = 32 (дет.)<br/>\n8 · 6 = 48 (дет.)<br/>\n8 · 7 = 56 (дет.)<br/>\n8 · 9 = 72 (дет.)\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica46-nomer5-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 46, номер 5-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 46, номер 5-1, год 2022."/>', 'Справка: Из формулы работы можно найти величины w и t по правилу нахождения неизвестного множителя: w = A : t      t = A : w • Производительность равна работе, делённой на время работы. • Время равно работе, делённой на производительность.', '<div class="recomended-block">\n<span class="title">Справка:</span>\n<p>\nИз формулы работы можно найти величины w и t по правилу нахождения неизвестного множителя:<br/>\nw = A : t      t = A : w<br/>\n• Производительность равна работе, делённой на время работы.<br/>\n• Время равно работе, делённой на производительность.\n\n</p>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-46/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica46-nomer5.jpg', 'peterson/3/part3/page46/task5_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica46-nomer5-1.jpg', 'peterson/3/part3/page46/task5_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '808f3ec100c87bf724a8395192144ec09675dbfca0da2c5de5f46c00bab9eab2', '2,4,6,7,8,9', '["заполни"]'::jsonb, 'мастер вытачивает 8 деталей в час. сколько деталей он сделает за 2 ч, 4 ч, 6 ч, 7 ч, 9 ч, t ч? заполни в тетради таблицу. запиши формулу зависимости работы a, выполненной мастером, от времени работы t');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 47, '6', 0, 'Тане надо вымыть 36 тарелок. Сколько времени она затратит на эту работу, если будет мыть в минуту 2 тарелки, 3 тарелки, 4 тарелки, 6 тарелок, 9 тарелок, w тарелок? Заполни таблицу. Запиши формулу зависимости времени работы t от производительности w.', '</p> \n<p class="text">Тане надо вымыть 36 тарелок. Сколько времени она затратит на эту работу, если будет мыть в минуту 2 тарелки, 3 тарелки, 4 тарелки, 6 тарелок, 9 тарелок, w тарелок? Заполни таблицу. Запиши формулу зависимости времени работы t от производительности w.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 47, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 47, номер 6, год 2022."/>\n</div>\n</div>', 't = A : w 36 : 2 = 18 (мин) 36 : 3 = 12 (мин) 36 : 4 =9 (мин) 36 : 6 = 6 (мин) 36 : 9 = 4 (мин)', '<p>\nt = A : w<br/>\n36 : 2 = 18 (мин)<br/>\n36 : 3 = 12 (мин)<br/>\n36 : 4 =9 (мин)<br/>\n36 : 6 = 6 (мин)<br/>\n36 : 9 = 4 (мин)\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer6-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 47, номер 6-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 47, номер 6-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-47/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer6.jpg', 'peterson/3/part3/page47/task6_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer6-1.jpg', 'peterson/3/part3/page47/task6_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '82e40be4be912634adaf936ea9c1cfcc41989f194df47b921bcb2128818d0c1b', '2,3,4,6,9,36', '["заполни"]'::jsonb, 'тане надо вымыть 36 тарелок. сколько времени она затратит на эту работу, если будет мыть в минуту 2 тарелки, 3 тарелки, 4 тарелки, 6 тарелок, 9 тарелок, w тарелок? заполни таблицу. запиши формулу зависимости времени работы t от производительности w');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 47, '7', 1, 'Выполни действия: а) 152 · 387        б) 492 · 604 в) 999 · 555        г) 333 · 707', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 152 · 387        б) 492 · 604<br/>   \nв) 999 · 555        г) 333 · 707\n</p>', 'а) 152 · 387 = 58824', '<p>\nа) 152 · 387 = 58824\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 47, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 47, номер 7, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-47/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer7.jpg', 'peterson/3/part3/page47/task7_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd427e9cf2395b949bf9e14d6ab87b2686c8018bc01d116c03152dd7cdd0f72ac', '152,333,387,492,555,604,707,999', NULL, 'выполни действия:а) 152*387        б) 492*604 в) 999*555        г) 333*707');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 47, '8', 2, 'Вычисли. Расположи ответы примеров в порядке убывания и расшифруй название цветка. Узнай, почему он так называется.', '</p> \n<p class="text">Вычисли. Расположи ответы примеров в порядке убывания и расшифруй название цветка. Узнай, почему он так называется.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="270" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 47, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 47, номер 8, год 2022."/>\n</div>\n</div>', 'Я - 960 · 24 = 23040', '<p>\nЯ - 960 · 24 = 23040\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer8-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 47, номер 8-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 47, номер 8-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-47/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer8.jpg', 'peterson/3/part3/page47/task8_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer8-1.jpg', 'peterson/3/part3/page47/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '13a2c5b209574cd58353c5b26ae8752b42fceb50362935e2f0a48483bd6f993c', NULL, '["вычисли"]'::jsonb, 'вычисли. расположи ответы примеров в порядке убывания и расшифруй название цветка. узнай, почему он так называется');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 47, '9', 3, 'Составь программу действий и вычисли: а) 234240 : 6 · 9 - (20030 - 7358) : 4 б) 834024 + 7900 · 25 - (483 · 504) : 8 · 10', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) 234240 : 6 · 9 - (20030 - 7358) : 4<br/>\nб) 834024 + 7900 · 25 - (483 · 504) : 8 · 10\n</p>', 'а) 234240 : 6 · 9 - (20030 - 7358) : 4 = 348192 20030 - 7358 = 12672', '<p>\nа) 234240 : 6 · 9 - (20030 - 7358) : 4 = 348192<br/>\n20030 - 7358 = 12672\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 47, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 47, номер 9, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-47/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer9.jpg', 'peterson/3/part3/page47/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0329e53863b4b969eb34e1b578b4416e4881fda117c034e6fec59621beac5a88', '4,6,8,9,10,25,483,504,7358,7900', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) 234240:6*9-(20030-7358):4 б) 834024+7900*25-(483*504):8*10');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 47, '10', 4, '7 дм 5 мм     75 мм            6 т 8 ц     6800 кг 9 м 2 дм     920 дм            6 кг 8 г     6800 г 2 км 32 м     203200 см      6 ч 8 мин     68 мин', '</p> \n<p class="text">7 дм 5 мм <span class="okon">   </span> 75 мм            6 т 8 ц <span class="okon">   </span> 6800 кг<br/> \n9 м 2 дм <span class="okon">   </span> 920 дм            6 кг 8 г <span class="okon">   </span> 6800 г<br/>\n2 км 32 м <span class="okon">   </span> 203200 см      6 ч 8 мин <span class="okon">   </span> 68 мин \n</p>', '7 дм 5 мм > 75 мм            6 т 8 ц = 6800 кг 9 м 2 дм < 920 дм            6 кг 8 г < 6800 г 2 км 32 м = 203200 см      6 ч 8 мин > 68 мин', '<p>\n7 дм 5 мм &gt; 75 мм            6 т 8 ц = 6800 кг<br/> \n9 м 2 дм &lt; 920 дм            6 кг 8 г &lt; 6800 г<br/>\n2 км 32 м = 203200 см      6 ч 8 мин &gt; 68 мин \n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-47/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'addc070515e2205bd8a862bdf6d8e122d81424f1085f655d40378b0bd6340f8e', '2,5,6,7,8,9,32,68,75,920', NULL, '7 дм 5 мм     75 мм            6 т 8 ц     6800 кг 9 м 2 дм     920 дм            6 кг 8 г     6800 г 2 км 32 м     203200 см      6 ч 8 мин     68 мин');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 47, '11', 5, 'Реши уравнения с комментированием и сделай проверку: а) (700 : x + 20) : 4 = 40 б) 2 · (500 - y : 3) = 820', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (700 : x + 20) : 4 = 40<br/>  		\nб) 2 · (500 - y : 3) = 820\n</p>', 'а) (700 : x + 20) : 4 = 40 Чтобы найти делимое надо делитель (700 : x + 20) умножить на частное (700 : x + 20) = 4 · 40 (700 : x + 20) = 160 Что бы найти слагаемое 700 : x надо из суммы вычесть известное слагаемое 700 : x = 160 - 20 700 : x = 140 Чтобы найти делитель надо делимое разделить на частное х = 700 : 140 х = 5 Проверка: (700 : 5 + 20) : 4 = 40 б) 2 · (500 - y : 3) = 820 Чтобы найти сомножитель (500 - y : 3) надо произведение разделить на известный сомножитель (500 - y : 3) = 820 : 2 (500 - y : 3) = 410 Чтобы найти вычитаемое y : 3 надо из уменьшаемого вычесть разность y : 3 = 500 - 410 у : 3 = 90 Чтобы найти делимое надо делитель умножить на частное у = 90 · 3 у = 270 Проверка: 2 · (500 - 270 : 3) = 820', '<p>\nа) (700 : x + 20) : 4 = 40  <br/>		\nЧтобы найти делимое надо делитель (700 : x + 20) умножить на частное<br/>\n(700 : x + 20) = 4 · 40<br/>\n(700 : x + 20) = 160<br/>\nЧто бы найти слагаемое 700 : x надо из суммы вычесть известное слагаемое<br/>\n700 : x = 160 - 20<br/>\n700 : x = 140<br/>\nЧтобы найти делитель надо делимое разделить на частное<br/>\nх = 700 : 140<br/>\nх = 5<br/>\n<b>Проверка:</b> (700 : 5 + 20) : 4 = 40  <br/><br/>\nб) 2 · (500 - y : 3) = 820<br/>\nЧтобы найти сомножитель (500 - y : 3) надо произведение разделить на известный сомножитель<br/>\n(500 - y : 3) = 820 : 2<br/>\n(500 - y : 3) = 410<br/>\nЧтобы найти вычитаемое y : 3 надо из уменьшаемого вычесть разность<br/>\ny : 3 = 500 - 410<br/>\nу : 3 = 90<br/>\nЧтобы найти делимое надо делитель умножить на частное<br/>\nу = 90 · 3<br/>\nу = 270 <br/>\n<b>Проверка:</b> 2 · (500 - 270 : 3) = 820\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-47/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '93d21b2e9dece76d468827a6cca7e859f1e7ab000f67a1d5d3af88dbe021ec92', '2,3,4,20,40,500,700,820', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (700:x+20):4=40 б) 2*(500-y:3)=820');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 47, '12', 6, 'Запиши множество делителей и множество кратных числа 27.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 27.</p>', 'Делители числа 27: 1, 3, 9, 27. Кратные числа 27 = 27, 54, 81, 108.', '<p>\nДелители числа 27: 1, 3, 9, 27. Кратные числа 27 = 27, 54, 81, 108.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-47/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f7b1a725aa5a1e2a4da08319114a4c1a4acd65419eeaa51082cf4653a04b669b', '27', NULL, 'запиши множество делителей и множество кратных числа 27');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 47, '13', 7, 'Проведи прямую а. Отметь на рисунке точки K, L, M и N, такие что K ∈ a, L ∉ a, M ∉ a, N ∈ a.', '</p> \n<p class="text">Проведи прямую а. Отметь на рисунке точки K,  L,  M и N, такие что K ∈ a, L ∉ a, M ∉ a, N ∈ a.</p>', 'А = {1, 2, 3, 4} В = {3, 4, 5, 6} А⋃В = {1,2, 3, 4, 5, 6} А⋂В = {3, 4}', '<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer13.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 47, номер 13, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 47, номер 13, год 2022."/>\n\n\n<p>\nА = {1, 2, 3, 4}<br/>\nВ = {3, 4, 5, 6}<br/>\nА⋃В = {1,2, 3, 4, 5, 6}<br/>\nА⋂В = {3, 4}\n</p>\n\n<div class="img-wrapper-460">\n<img width="250" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer14.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 47, номер 14, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 47, номер 14, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-47/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer13.jpg', 'peterson/3/part3/page47/task13_solution_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica47-nomer14.jpg', 'peterson/3/part3/page47/task13_solution_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f29896fde032c71d45b3034903283ad04782cb040aecdce52d0be65ae9ab9c6a', NULL, NULL, 'проведи прямую а. отметь на рисунке точки k, l, m и n, такие что k ∈ a, l ∉ a, m ∉ a, n ∈ a');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 48, '1', 0, 'Прочитай задачу и объясни, как составлена таблица. Составь план решения задачи и найди ответ. «Один оператор набрал на компьютере за 5 часов 90 страниц рукописи, а другой за 7 часов – 98 страниц. У кого из них производительность больше и на сколько?»', '</p> \n<p class="text">Прочитай задачу и объясни, как составлена таблица. Составь план решения задачи и найди ответ.<br/>\n«Один оператор набрал на компьютере за 5 часов 90 страниц рукописи, а другой за 7 часов – 98 страниц. У кого из них производительность больше и на сколько?»\n</p> \n\n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 48, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 48, номер 1, год 2022."/>\n</div>\n</div>', 'В первой строчке, что известно о первом операторе и во второй строчке таблицы, что известно о втором операторе. Необходимо найти разность их производительности. w = A : t 90 : 5 - 98 : 7 = 18 - 14 = 4 (стр./ч) Ответ: на 4 страницы в час больше производительность у первого оператора чем у второго.', '<p>\nВ первой строчке, что известно о первом операторе и во второй строчке таблицы, что известно о втором операторе. Необходимо найти разность их производительности.<br/>\nw = A : t <br/>   \n90 : 5 - 98 : 7 = 18 - 14 = 4 (стр./ч)<br/>\n<b>Ответ:</b> на 4 страницы в час больше производительность у первого оператора чем у второго.  \n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-48/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer1.jpg', 'peterson/3/part3/page48/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1d6ddde33606640822b8eb22659a0118cb16f794c199b807a1ad49431a009437', '5,7,90,98', '["найди","больше"]'::jsonb, 'прочитай задачу и объясни, как составлена таблица. составь план решения задачи и найди ответ. "один оператор набрал на компьютере за 5 часов 90 страниц рукописи, а другой за 7 часов-98 страниц. у кого из них производительность больше и на сколько?"');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 48, '2', 1, 'Реши задачи с помощью таблиц: а) За 6 дней на фабрике сшили 1926 костюмов. Сколько костюмов сошьют на этой фабрике за год (256 рабочих дней), если будут работать с той же производительностью? б) Экскаватор за 1 час копает 18 м канавы. Одну канаву он выкопал за 7 ч, а другую – за 19 ч. Сколько метров канавы выкопал экскаватор за всё это время?', '</p> \n<p class="text">Реши задачи с помощью таблиц:<br/> \nа) За 6 дней на фабрике сшили 1926 костюмов. Сколько костюмов сошьют на этой фабрике за год (256 рабочих дней), если будут работать с той же производительностью?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 48, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 48, номер 2, год 2022."/>\n</div>\n</div>\n\n<p class="text">б) Экскаватор за 1 час копает 18 м канавы. Одну канаву он выкопал за 7 ч, а другую – за 19 ч. Сколько метров канавы выкопал экскаватор за всё это время?</p> \n\n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer2-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 48, номер 2-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 48, номер 2-1, год 2022."/>\n</div>\n</div>', 'а) w = A : t 1926 : 6 = 321 (к./дн.) А = w · t 321 · 256 = 82176 (к.)', '<p>\nа) w = A : t<br/>\n1926 : 6 = 321 (к./дн.)<br/>\nА = w · t<br/>\n321 · 256 = 82176 (к.)\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer2-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 48, номер 2-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 48, номер 2-2, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-48/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer2.jpg', 'peterson/3/part3/page48/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer2-1.jpg', 'peterson/3/part3/page48/task2_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer2-2.jpg', 'peterson/3/part3/page48/task2_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '83231d63925700e54333a446dc3a9bc5bcc15e72b6f003ccb12d467b46bec289', '1,6,7,18,19,256,1926', '["реши"]'::jsonb, 'реши задачи с помощью таблиц:а) за 6 дней на фабрике сшили 1926 костюмов. сколько костюмов сошьют на этой фабрике за год (256 рабочих дней), если будут работать с той же производительностью? б) экскаватор за 1 час копает 18 м канавы. одну канаву он выкопал за 7 ч, а другую-за 19 ч. сколько метров канавы выкопал экскаватор за всё это время');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 48, '3', 2, 'а) Два друга взяли в библиотеке одинаковые книги. Первый читает 3 страницы в день, а второй – 9 страниц в день. Кто из них прочитает эту книгу раньше и на сколько дней, если в книге 360 страниц? б) Мастер сделал на станке 72 детали за 3 часа. Сколько деталей он сделает за 8 часов, если будет работать с той же производительностью?', '</p> \n<p class="text">а) Два друга взяли в библиотеке одинаковые книги. Первый читает 3 страницы в день, а второй – 9 страниц в день. Кто из них прочитает эту книгу раньше и на сколько дней, если в книге 360 страниц? <br/> \nб) Мастер сделал на станке 72 детали за 3 часа. Сколько деталей он сделает за 8 часов, если будет работать с той же производительностью?\n</p>', 'а) t = А : w 360 : 3 - 360 : 9 = 120 - 40 = 80 (д.) Ответ: второй читает 9 страниц в день из них он прочитает эту книгу раньше на 80 дней, если в книге 360 страниц. б) w = А : t 72 : 3 = 24 (д./ч) А = w · t 24 · 8 = 192 (д.)', '<p>\nа) t = А : w<br/>\n360 : 3 - 360 : 9 = 120 - 40 = 80 (д.)<br/>\n<b>Ответ:</b> второй читает 9 страниц в день из них он прочитает эту книгу раньше на 80 дней, если в книге 360 страниц.<br/><br/>  \nб) w = А : t <br/>\n72 : 3 = 24 (д./ч)<br/>\nА = w · t<br/>\n24 · 8 = 192 (д.)\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="80" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 48, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 48, номер 3, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Два друга взяли в библиотеке одинаковые книги. Первый читает 3 страницы в день, а второй – 9 страниц в день. Кто из них прочитает эту книгу раньше и на сколько дней, если в книге 360 страниц?","solution":"t = А : w 360 : 3 - 360 : 9 = 120 - 40 = 80 (д.) Ответ: второй читает 9 страниц в день из них он прочитает эту книгу раньше на 80 дней, если в книге 360 страниц."},{"letter":"б","condition":"Мастер сделал на станке 72 детали за 3 часа. Сколько деталей он сделает за 8 часов, если будет работать с той же производительностью?","solution":"w = А : t 72 : 3 = 24 (д./ч) А = w · t 24 · 8 = 192 (д.)"}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-48/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer3.jpg', 'peterson/3/part3/page48/task3_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0ea5de13daa77f0d8034127743a6b199a95db56af3e5a454d150b5cd819c9356', '3,8,9,72,360', NULL, 'а) два друга взяли в библиотеке одинаковые книги. первый читает 3 страницы в день, а второй-9 страниц в день. кто из них прочитает эту книгу раньше и на сколько дней, если в книге 360 страниц? б) мастер сделал на станке 72 детали за 3 часа. сколько деталей он сделает за 8 часов, если будет работать с той же производительностью');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 48, '4', 3, 'Маляр должен покрасить заводской забор длиной 243 м. Первые 3 дня он красил по 18 м забора в день. За сколько времени он выполнил всю работу, если в оставшиеся дни он увеличил производительность на 3 метра в день?', '</p> \n<p class="text">Маляр должен покрасить заводской забор длиной 243 м. Первые 3 дня он красил по 18 м забора в день. За сколько времени он выполнил всю работу, если в оставшиеся дни он увеличил производительность на 3 метра в день?</p>', 't = А : w 243 - 18 - 18 - 18 = 225 - 18 - 18 = 207 - 18 = 189 (м) – отсталость после первых 3 дней', '<p>\nt = А : w<br/>\n243 - 18 - 18 - 18 = 225 - 18 - 18 = 207 - 18 = 189 (м) – отсталость после первых 3 дней\n</p>\n\n<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 48, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 48, номер 4, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-48/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica48-nomer4.jpg', 'peterson/3/part3/page48/task4_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd41c76052671c9bb44db178d25d4d6d209687c9a820defb4fa7b010d09ec6a16', '3,18,243', NULL, 'маляр должен покрасить заводской забор длиной 243 м. первые 3 дня он красил по 18 м забора в день. за сколько времени он выполнил всю работу, если в оставшиеся дни он увеличил производительность на 3 метра в день');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 49, '5', 0, 'Найди ошибки в записи в решении примеров: Запиши и реши их в тетради правильно.', '</p> \n<p class="text">Найди ошибки в записи в решении примеров:<br/>\nЗапиши и реши их в тетради правильно.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 49, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 49, номер 5, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-460">\n<img width="170" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer5-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 49, номер 5-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 49, номер 5-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-49/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer5.jpg', 'peterson/3/part3/page49/task5_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer5-1.jpg', 'peterson/3/part3/page49/task5_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8cafe8375fe4bbb8185b194511a183f2dbe55c8ddb34181683c50721fd84dd2a', NULL, '["найди","реши"]'::jsonb, 'найди ошибки в записи в решении примеров:запиши и реши их в тетради правильно');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 49, '6', 1, 'Выполни действия. Проверь результаты с помощью калькулятора. а) 254 · 966        б) 809 · 421 в) 358 · 604        г) 705 · 108', '</p> \n<p class="text">Выполни действия. Проверь результаты с помощью калькулятора. </p> \n\n<p class="description-text"> \nа) 254 · 966        б) 809  · 421<br/>  \nв) 358 · 604        г) 705 · 108\n</p>', 'а) 254 · 966 = 245364', '<p>\nа) 254 · 966 = 245364\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 49, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 49, номер 6, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-49/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer6.jpg', 'peterson/3/part3/page49/task6_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ca93090091f8a9ca02d05d626bffa991564994e80fc4c6196c758cb304de1ac1', '108,254,358,421,604,705,809,966', NULL, 'выполни действия. проверь результаты с помощью калькулятора. а) 254*966        б) 809*421 в) 358*604        г) 705*108');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 49, '7', 2, 'а) По формуле a = b · c + r, r < b найди делимое, если делитель равен 8, частное 25, а остаток 5. б) Выполни деление с остатком и сделай проверку: 976326 : 7    702514 : 5    183600 : 70', '</p> \n<p class="text">а) По формуле a = b · c + r, r &lt; b найди делимое, если делитель равен 8, частное 25, а остаток 5.<br/><br/>\nб) Выполни деление с остатком и сделай проверку:<br/>\n976326 : 7    702514 : 5    183600 : 70\n</p>', 'а) а = 8 · 25 + 5 а = 200 + 5 а = 205', '<p>\nа) а = 8 · 25 + 5 <br/>\nа = 200 + 5<br/>\nа = 205\n</p>\n\n<div class="img-wrapper-460">\n<img width="80" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 49, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 49, номер 7, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"По формуле a = b · c + r, r \u003c b найди делимое, если делитель равен 8, частное 25, а остаток 5.","solution":"а = 8 · 25 + 5 а = 200 + 5 а = 205"},{"letter":"б","condition":"Выполни деление с остатком и сделай проверку: 976326 : 7    702514 : 5    183600 : 70","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-49/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer7.jpg', 'peterson/3/part3/page49/task7_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd5a70a9afa6be002d8126c91b9ba3612442f3296081c199b6bcdf6dd8f3d6a8e', '5,7,8,25,70,183600,702514,976326', '["найди","частное","делитель","делимое","остаток"]'::jsonb, 'а) по формуле a=b*c+r, r<b найди делимое, если делитель равен 8, частное 25, а остаток 5. б) выполни деление с остатком и сделай проверку:976326:7    702514:5    183600:70');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 49, '8', 3, 'Найди значение выражения: (720 - 99) · 324 - (728 + 50 · 90)', '</p> \n<p class="text">Найди значение выражения:</p> \n\n<p class="description-text"> \n(720 - 99) · 324 - (728 + 50 · 90)\n</p>', '(720 - 99) · 324 - (728 + 50 · 90) = 621 · 324 - (778 · 90) = 201204 - 70020 = 131184', '<p>\n(720 - 99) · 324 - (728 + 50 · 90) = 621 · 324 - (778 · 90) = 201204 - 70020 = 131184\n</p>\n\n<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 49, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 49, номер 8, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-49/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer8.jpg', 'peterson/3/part3/page49/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1bd75bda3b9cbeea50b899356b6526405e4a4764a892333ebfef3c94c33aeccc', '50,90,99,324,720,728', '["найди"]'::jsonb, 'найди значение выражения:(720-99)*324-(728+50*90)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 49, '9', 4, 'Реши уравнения с комментированием:я а) (720 - t · 6) : 9 = 60 б) 4 · (250 : a + 12) = 68', '</p> \n<p class="text">Реши уравнения с комментированием:я</p> \n\n<p class="description-text"> \nа) (720 - t · 6) : 9 = 60 <br/>      \nб) 4 · (250 : a + 12) = 68\n</p>', 'а) (720 - t · 6) : 9 = 60 Чтобы найти делимое (720 - t · 6) надо делитель умножить на частное (720 - t · 6) = 9 · 60 720 - t · 6 = 540', '<p>\nа) (720 - t · 6) : 9 = 60  <br/> \nЧтобы найти делимое (720 - t · 6) надо делитель умножить на частное<br/> \n(720 - t · 6) = 9 · 60<br/> \n720 - t · 6 = 540\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="70" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 49, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 49, номер 9, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-49/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer9.jpg', 'peterson/3/part3/page49/task9_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '245821a58e7790d1b82cd297854bfd90934d7d2ba0171d346fda693d035ecd6b', '4,6,9,12,60,68,250,720', '["реши"]'::jsonb, 'реши уравнения с комментированием:я а) (720-t*6):9=60 б) 4*(250:a+12)=68');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 49, '10', 5, 'Вычисли. Расположи ответы в порядке убывания. Расшифруй, как называли в Древнем Риме богинь красоты? Узнай, сколько их было? Какие у них имена?', '</p> \n<p class="text">Вычисли. Расположи ответы в порядке убывания. Расшифруй, как называли в Древнем Риме богинь красоты? Узнай, сколько их было? Какие у них имена?</p> \n\n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="370" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer10.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 49, номер 10, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 49, номер 10, год 2022."/>\n</div></div>', 'Я - 140 + 60 - 280 : 7 · 5 = 200 - 40 · 5 = 200 - 200 = 0 Г - 90 · 3 + 20 – 140 : 5 = 270 + 20 - 28 = 262 А - (400 - 25 · 3 · 2) : 10 = (400 - 150) : 10 = 250 : 10 = 25 И - (17 + 7 · 9 + 5 · 8) : 20 = (17 + 63 + 40) : 20 = 120 : 20 = 6 Р - 130 · 2 – 360 : 30 = 260 - 12 = 248 Ц - (270 - 240 : 4 · 3) : 9 = (270 - 180) : 9 = 90 : 9 = 10 ГРАЦИЯ – Три Грации. Дочери Зевса, называемые Гесиодом Аглая (сияющая), Евфросина (благомыслящая) и Талия (цветущая), олицетворяют доброе, радостное, вечно юное начало жизни. Грации часто сопровождают богиню любви Афродиту.', '<p>\nЯ - 140 + 60 - 280 : 7 · 5 = 200 - 40 · 5 = 200 - 200 = 0<br/>\nГ - 90 · 3 + 20 – 140 : 5 = 270 + 20 - 28 = 262<br/>\nА - (400 - 25 · 3 · 2) : 10 = (400 - 150) : 10 = 250 : 10 = 25<br/>\nИ - (17 + 7 · 9 + 5 · 8) : 20 = (17 + 63 + 40) : 20 = 120 : 20 = 6<br/>\nР - 130 · 2 – 360 : 30 = 260 - 12 = 248<br/>\nЦ - (270 - 240 : 4 · 3) : 9 = (270 - 180) : 9 = 90 : 9 = 10<br/>\nГРАЦИЯ – Три Грации. Дочери Зевса, называемые Гесиодом Аглая (сияющая), Евфросина (благомыслящая) и Талия (цветущая), олицетворяют доброе, радостное, вечно юное начало жизни. Грации часто сопровождают богиню любви Афродиту.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-49/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica49-nomer10.jpg', 'peterson/3/part3/page49/task10_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '28b2874e6047503a5b67a822d8a616a2bebefcc974063a4a833e5e9fc882381c', NULL, '["вычисли"]'::jsonb, 'вычисли. расположи ответы в порядке убывания. расшифруй, как называли в древнем риме богинь красоты? узнай, сколько их было? какие у них имена');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 49, '11', 6, 'Запиши множество делителей и множество кратных числа 28.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 28.</p>', 'Множество делителей числа 28: 1, 2, 4, 7, 14, 28.', '<p>\nМножество делителей числа 28: 1, 2, 4, 7, 14, 28.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-49/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'bfe1eb94a2a95ca8cc8cddee06bb44b024ce62c8c70d49a0864c9166b55dd309', '28', NULL, 'запиши множество делителей и множество кратных числа 28');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 49, '12', 7, 'Ребро куба равно 11 см. Найди площадь поверхности куба и сумму длин всех его рёбер. Чему равен объём этого куба?', '</p> \n<p class="text">Ребро куба равно 11 см. Найди площадь поверхности куба и сумму длин всех его рёбер. Чему равен объём этого куба?</p>', 'S = 6a 2 S = 6 · 11 2 S = 6 · 121 S = 726 (см 2 ) 11 · 4 + 11 · 4 + 11 · 4 = 44 + 44 + 44 = 88 + 44 = 132 (см) V= a 3 V= 11 3 V= 1331 (см 3 ) Ответ: площадь поверхности куба равна 726 см 2 , сумма длин всех его рёбер равна 132 см. 1331 см 3 равен объём этого куба. 100000000000000000000 за самым большим 20 - значным числом следует самое маленькое 21 - значное число. 14 + 6 = 20 квадратов', '<p>\nS = 6a<sup>2</sup><br/>\nS = 6 · 11<sup>2</sup><br/>\nS = 6 · 121<br/>\nS = 726 (см<sup>2</sup>)<br/>\n11 · 4 + 11 · 4 + 11 · 4 = 44 + 44 + 44 = 88 + 44 = 132 (см)<br/>\nV= a<sup>3</sup><br/>\nV= 11<sup>3</sup><br/>\nV= 1331 (см<sup>3</sup>)<br/>\n<b>Ответ:</b> площадь поверхности куба равна 726 см<sup>2</sup>, сумма длин всех его рёбер равна 132 см. 1331 см<sup>3</sup> равен объём этого куба.\n\n</p>\n\n\n<p>\n100000000000000000000 за самым большим 20 - значным числом следует самое маленькое 21 - значное число.\n</p>\n\n\n<p>\n14 + 6 = 20 квадратов\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-49/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'bb4b0b1f54d832ebb27081e0380a6e4332144dfd7f9d805676e68b43e7b99c4a', '11', '["найди","площадь","равно"]'::jsonb, 'ребро куба равно 11 см. найди площадь поверхности куба и сумму длин всех его рёбер. чему равен объём этого куба');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 50, '1', 0, 'Выбери примеры на умножение круглых чисел и вычисли: а) 725 · 8200        б) 349 · 506 в) 8070 · 3680      г) 40300 · 9040', '</p> \n<p class="text">Выбери примеры на умножение круглых чисел и вычисли:</p> \n\n<p class="description-text"> \nа) 725 · 8200        б) 349 · 506 <br/> \nв) 8070 · 3680      г) 40300 · 9040\n</p>', 'а) 725 · 8200 = 5945000', '<p>а) 725 · 8200 = 5945000\n</p>\n\n<div class="img-wrapper-460">\n<img width="180" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica50-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 50, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 50, номер 1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-50/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica50-nomer1.jpg', 'peterson/3/part3/page50/task1_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'dad0feae2c5f5a59795e983bb825c05ca818e617aad7cb7be6f705bc0275b7a6', '349,506,725,3680,8070,8200,9040,40300', '["вычисли"]'::jsonb, 'выбери примеры на умножение круглых чисел и вычисли:а) 725*8200        б) 349*506 в) 8070*3680      г) 40300*9040');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 50, '2', 1, 'Найди производительность, если: а) бабушка связала 36 рядов за 3 часа; б) ученик прочитал 270 слов за 6 минут; в) садовник посадил 54 цветка за 2 часа; г) завод выпустил 480 машин за 4 дня.', '</p> \n<p class="text">Найди производительность, если:<br/>\nа) бабушка связала 36 рядов за 3 часа;<br/>\nб) ученик прочитал 270 слов за 6 минут; <br/>\nв) садовник посадил 54 цветка за 2 часа;<br/>\nг) завод выпустил 480 машин за 4 дня.\n</p>', 'а) 36 : 3 = 12 (рядов/ч) б) 270 : 6 = 45 (слов/минуту) в) 54 : 2 = 27 (цветка/час) г) 480 : 4 = 120 (машин/день)', '<p>\nа) 36 : 3 = 12 (рядов/ч)<br/>\nб) 270 : 6 = 45 (слов/минуту) <br/>\nв) 54 : 2 = 27 (цветка/час)<br/>\nг) 480 : 4 = 120 (машин/день)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-50/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e3d448ab263b6ddddeb45671e8ace95871d816b98f390756663a8b950cde332d', '2,3,4,6,36,54,270,480', '["найди"]'::jsonb, 'найди производительность, если:а) бабушка связала 36 рядов за 3 часа; б) ученик прочитал 270 слов за 6 минут; в) садовник посадил 54 цветка за 2 часа; г) завод выпустил 480 машин за 4 дня');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 50, '3', 2, 'Найди пропущенные значения величин:', '</p> \n<p class="text">Найди пропущенные значения величин:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica50-nomer3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 50, номер 3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 50, номер 3, год 2022."/>\n</div>\n</div>', 'w = 72 : 6            А = 8 · 7 w = 12 (шт/ч)      А = 56 (шт.) А = 50 · 3             t = 900 : 150 А = 150 (т)           t = 6 (дней) t = 400 : 80         w = 420 : 14 t = 5 (год)            w = 30 (шт./ч)', '<p>\nw = 72 : 6            А = 8 · 7<br/>\nw = 12 (шт/ч)      А = 56 (шт.)<br/>\nА = 50 · 3             t = 900 : 150<br/>\nА = 150 (т)           t = 6 (дней)<br/>\nt = 400 : 80         w = 420 : 14<br/>\nt = 5 (год)            w = 30 (шт./ч)\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-50/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica50-nomer3.jpg', 'peterson/3/part3/page50/task3_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e65a4c11d69eee7f78cf09a0b0bf2bac39f2b16b64e89825a0bf169f9ca1646a', NULL, '["найди"]'::jsonb, 'найди пропущенные значения величин');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 50, '4', 3, 'Мастер получил заказ на изготовление 600 деталей. Первые 4 часа он делал по 70 деталей в час. Затем он увеличил производительность на 10 деталей в час. За сколько часов он выполнил весь заказ?', '</p> \n<p class="text">Мастер получил заказ на изготовление 600 деталей. Первые 4 часа он делал по 70 деталей в час. Затем он увеличил производительность на 10 деталей в час. За сколько часов он выполнил весь заказ?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica50-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 50, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 50, номер 4, год 2022."/>\n</div>\n</div>', '70 · 4 = 280 (д.) 600 - 280 = 320 (д.) – деталей потом 320 : (70 + 10) = 4 (ч) 4 + 4 = 8 (ч) Ответ: за 8 часов он выполнил весь заказ.', '<p>\n70 · 4 = 280 (д.)<br/>\n600 - 280 = 320 (д.) – деталей потом<br/>\n320 : (70 + 10) = 4 (ч)<br/>\n4 + 4 = 8 (ч) <br/>\n<b>Ответ:</b> за 8 часов он выполнил весь заказ.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-50/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica50-nomer4.jpg', 'peterson/3/part3/page50/task4_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'dfcf71b01368bd7c4784ffc139dc2ed4dd966adfb5ad7cd1aabb285620eeb328', '4,10,70,600', NULL, 'мастер получил заказ на изготовление 600 деталей. первые 4 часа он делал по 70 деталей в час. затем он увеличил производительность на 10 деталей в час. за сколько часов он выполнил весь заказ');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 50, '5', 4, 'Реши задачи и сравни их. Что ты замечаешь? а) Токарь вытачивает 240 деталей за 3 дня, а его ученик – за 4 дня. На сколько производительность токаря выше производительности ученика? б) У Димы в копилке 240 р. Он может купить на них 3 книги по одной цене или 4 одинаковых альбома. На сколько альбом дешевле книги? в) Расстояние между Москвой и Ярославлем равно 240 км. Автобус проходит это расстояние за 4 ч, а поезд – за 3 ч. На сколько километров в час скорость поезда больше скорости автобуса? г) Бассейн, объём которого 240 м 3 , наполняется первой трубой за 3 ч, а второй трубой – за 4 ч. На сколько скорость наполнения бассейна первой трубой больше скорости наполнения второй трубой? Придумай задачу с другими величинами, которая решается так же.', '</p> \n<p class="text">Реши задачи и сравни их. Что ты замечаешь?<br/>\nа) Токарь вытачивает 240 деталей за 3 дня, а его ученик – за 4 дня. На сколько производительность токаря выше производительности ученика? <br/>\nб) У Димы в копилке 240 р. Он может купить на них 3 книги по одной цене или 4 одинаковых альбома. На сколько альбом дешевле книги?<br/>\nв) Расстояние между Москвой и Ярославлем равно 240 км. Автобус проходит это расстояние за 4 ч, а поезд – за 3 ч. На сколько километров в час скорость поезда больше скорости автобуса?<br/>\nг) Бассейн, объём которого 240 м<sup>3</sup>, наполняется первой трубой за 3 ч, а второй трубой – за 4 ч. На сколько скорость наполнения бассейна первой трубой больше скорости наполнения второй трубой?<br/>\nПридумай задачу с другими величинами, которая решается так же.\n</p>', 'а) 240 : 3 - 240 : 4 = 80 - 60 = 20 (д./дня) Ответ: на 20 деталей в день производительность токаря выше производительности ученика. б) 240 : 3 - 240 : 4 = 80 - 60 = 20 (р.) Ответ: на 20 р. альбом дешевле книги. в) 240 : 3 - 240 : 4 = 80 - 60 = 20 (км/ч) Ответ: на 20 километров в час скорость поезда больше скорости автобуса. г) 240 : 3 - 240 : 4 = 80 - 60 = 20 (м 3 /ч) Ответ: на 20 скорость наполнения бассейна первой трубой больше скорости наполнения второй трубой. Задача: У мамы в копилке 240 р. Она может купить на них 3 кг яблок по одной цене или 4 кг груш. На сколько груши дешевле яблок. 240 : 3 - 240 : 4 = 80 - 60 = 20 (р) Ответ: на 20 р. груши дешевле яблок.', '<p>\nа) 240 : 3 - 240 : 4 = 80 - 60 = 20 (д./дня)<br/>\n<b>Ответ:</b> на 20 деталей в день производительность токаря выше производительности ученика.<br/>\nб) 240 : 3 - 240 : 4 = 80 - 60 = 20 (р.)<br/>\n<b>Ответ:</b> на 20 р. альбом дешевле книги.<br/>\nв) 240 : 3 - 240 : 4 = 80 - 60 = 20 (км/ч)<br/>\n<b>Ответ:</b> на 20 километров в час скорость поезда больше скорости автобуса.<br/>\nг) 240 : 3 - 240 : 4 = 80 - 60 = 20 (м<sup>3</sup>/ч)<br/>\n<b>Ответ:</b> на 20 скорость наполнения бассейна первой трубой больше скорости наполнения второй трубой.<br/><br/>\n<b>Задача:</b> <br/>\nУ мамы в копилке 240 р. Она может купить на них 3 кг яблок по одной цене или 4 кг груш. На сколько груши дешевле яблок. <br/>\n240 : 3 - 240 : 4 = 80 - 60 = 20 (р)<br/>\n<b>Ответ:</b> на 20 р. груши дешевле яблок.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-50/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b380f1e16dab63cd2f9e3ced594f69716ce3f60bbc91d47adc021df71e046da0', '3,4,240', '["реши","сравни","больше","равно"]'::jsonb, 'реши задачи и сравни их. что ты замечаешь? а) токарь вытачивает 240 деталей за 3 дня, а его ученик-за 4 дня. на сколько производительность токаря выше производительности ученика? б) у димы в копилке 240 р. он может купить на них 3 книги по одной цене или 4 одинаковых альбома. на сколько альбом дешевле книги? в) расстояние между москвой и ярославлем равно 240 км. автобус проходит это расстояние за 4 ч, а поезд-за 3 ч. на сколько километров в час скорость поезда больше скорости автобуса? г) бассейн, объём которого 240 м 3 , наполняется первой трубой за 3 ч, а второй трубой-за 4 ч. на сколько скорость наполнения бассейна первой трубой больше скорости наполнения второй трубой? придумай задачу с другими величинами, которая решается так же');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 51, '6', 0, 'Вычисли. Расположи ответы примеров в порядке убывания. Кто это? Найди информацию о нём в Интернете или энциклопедии.', '</p> \n<p class="text">Вычисли. Расположи ответы примеров в порядке убывания. Кто это? Найди информацию о нём в Интернете или энциклопедии.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 51, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 51, номер 6, год 2022."/>\n</div>\n</div>', 'И - 340 · 750 = 255000', '<p>\nИ - 340 · 750 = 255000\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer6-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 51, номер 6-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 51, номер 6-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-51/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer6.jpg', 'peterson/3/part3/page51/task6_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer6-1.jpg', 'peterson/3/part3/page51/task6_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a5b2c5ea1e330e18a0e4294653f006a2557c2f529de3e354c6f876b4beb8c09d', NULL, '["вычисли","найди"]'::jsonb, 'вычисли. расположи ответы примеров в порядке убывания. кто это? найди информацию о нём в интернете или энциклопедии');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 51, '7', 1, 'Вычисли и сравни значения выражений. Что ты замечаешь? 3524120 - 398705 : 5 · 40 (3524120 - 398705) : 5 · 40 (3524120 - 398705 : 5) · 40', '</p> \n<p class="text">Вычисли и сравни значения выражений. Что ты замечаешь?</p> \n\n<p class="description-text"> \n3524120 - 398705 : 5 · 40<br/>\n(3524120 - 398705) : 5 · 40<br/>\n(3524120 - 398705 : 5) · 40\n</p>', '3524120 - 398705 : 5 · 40 = 3524120 - 79741 · 40 = 3524120 - 3189640 = 334480', '<p>\n3524120 - 398705 : 5 · 40 = 3524120 - 79741 · 40 = 3524120 - 3189640 = 334480\n</p>\n\n<div class="img-wrapper-460">\n<img width="220" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 51, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 51, номер 7, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-51/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer7.jpg', 'peterson/3/part3/page51/task7_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0a9de2348d6ae28c0397bb7e16e3782fcef5886d4678db025dd8c9f91f7c0366', '5,40,398705,3524120', '["вычисли","сравни"]'::jsonb, 'вычисли и сравни значения выражений. что ты замечаешь? 3524120-398705:5*40 (3524120-398705):5*40 (3524120-398705:5)*40');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 51, '8', 2, 'Выполни действия: а) 7 м 85 см · 412        в) 6 дм 3 94 см 3 · 904 б) 4 см 2 6 мм 2 · 503        г) 3 кг 68 г · 706 д) 8 мин 24 с · 375 е) 1 ч 15 мин · 576', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 7 м 85 см · 412        в) 6 дм<sup>3</sup> 94 см<sup>3</sup> · 904<br/>  	\nб) 4 см<sup>2</sup> 6 мм<sup>2</sup> · 503        г) 3 кг 68 г · 706  <br/><br/>     	\n\n\nд) 8 мин 24 с · 375<br/>\nе) 1 ч 15 мин · 576\n\n</p>', 'а) 7 м 85 см · 412 = 323420 см = 3234 м 20 см', '<p>\nа) 7 м 85 см · 412 = 323420 см = 3234 м 20 см\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 51, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 51, номер 8, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-51/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer8.jpg', 'peterson/3/part3/page51/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c4af974c2d3261c9ba8ffadfc2624f5aa1c5115fc291e88da61513c4947d0d5b', '1,2,3,4,6,7,8,15,24,68', NULL, 'выполни действия:а) 7 м 85 см*412        в) 6 дм 3 94 см 3*904 б) 4 см 2 6 мм 2*503        г) 3 кг 68 г*706 д) 8 мин 24 с*375 е) 1 ч 15 мин*576');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 51, '9', 3, 'Площадь нижней грани прямоугольного параллелепипеда равна 800 см2. Определи высоту этого параллелепипеда, если его объём равен 24000 см 3 .', '</p> \n<p class="text">Площадь нижней грани прямоугольного параллелепипеда равна 800 см2. Определи высоту этого параллелепипеда, если его объём равен 24000 см<sup>3</sup>.</p>', '24000 см 3 : 800 см 2 = 30 (см) Ответ: 30 сантиметров высота этого параллелепипеда.', '<p>\n24000 см<sup>3</sup> : 800 см<sup>2</sup> = 30 (см)<br/>\n<b>Ответ:</b> 30 сантиметров высота этого параллелепипеда.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-51/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '19af738fca94312f9c4b0fe20258ec26cd49319fc3f4eeff36fb4f5532ff5c57', '2,3,800,24000', '["площадь"]'::jsonb, 'площадь нижней грани прямоугольного параллелепипеда равна 800 см2. определи высоту этого параллелепипеда, если его объём равен 24000 см 3');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 51, '10', 4, 'Напиши формулу объёма прямоугольного параллелепипеда, если у него: а) длина равна 8, ширина 4, высота c; б) площадь основания 45, а высота h; в) площадь основания S, а высота h.', '</p> \n<p class="text">Напиши формулу объёма прямоугольного параллелепипеда, если у него: <br/>\nа) длина равна 8, ширина 4, высота c; <br/>\nб) площадь основания 45, а высота h; <br/>\nв) площадь основания S, а высота h.\n</p>', 'а) V = l · w · h, где V - объем, l - длина, w - ширина и h - высота параллелепипеда. V = 8 · 4 · с; б) V = 45 · h; в) V = S · h.', '<p>\nа) V = l · w · h, где V - объем, l - длина, w - ширина и h - высота параллелепипеда. <br/>\nV = 8 · 4 · с; <br/>\nб) V = 45 · h; <br/>\nв) V = S · h.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-51/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '87e8c8978c3547db77985441b06a7f9dcc869fd4b7074589e3d15bb55c6ecd72', '4,8,45', '["площадь"]'::jsonb, 'напиши формулу объёма прямоугольного параллелепипеда, если у него:а) длина равна 8, ширина 4, высота c; б) площадь основания 45, а высота h; в) площадь основания s, а высота h');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 51, '11', 5, 'Рассмотри таблицы. Как связаны между собой переменные x и у ? Составь формулу, выражающую y через x.', '</p> \n<p class="text">Рассмотри таблицы. Как связаны между собой переменные x и у ? Составь формулу, выражающую y через x.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer11.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 51, номер 11, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 51, номер 11, год 2022."/>\n</div>\n</div>', 'а) у = 3х б) у = х + 4', '<p>\nа) у = 3х<br/>\nб) у = х + 4\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-51/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer11.jpg', 'peterson/3/part3/page51/task11_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0708bdca7ccf5a48f3fb9b082016c210faf13338f42bd784e13d1ea3c5f93c1d', NULL, NULL, 'рассмотри таблицы. как связаны между собой переменные x и у ? составь формулу, выражающую y через x');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 51, '12', 6, 'На рисунке все фигуры, кроме одной, имеют общее свойство. Какая фигура «лишняя»?', '</p> \n<p class="text">На рисунке все фигуры, кроме одной, имеют общее свойство. Какая фигура «лишняя»?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer12.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 51, номер 12, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 51, номер 12, год 2022."/>\n</div>\n</div>', 'Все фигуры имеют ось симметрии горизонтальную и вертикальную кроме Е. У фигуры Е только одна ось симметрии.', '<p>\nВсе фигуры имеют ось симметрии горизонтальную и вертикальную кроме Е. У фигуры Е только одна ось симметрии.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-51/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica51-nomer12.jpg', 'peterson/3/part3/page51/task12_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '07f9bdfd3943975056c64f44cd2f3b0d7a2d78b382fd335283665be58720c157', NULL, NULL, 'на рисунке все фигуры, кроме одной, имеют общее свойство. какая фигура "лишняя"');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 52, '1', 0, 'Проанализируй, как связаны между собой величины каждой строки. Запиши формулу зависимости между ними. Что общего у всех записанных формул? Замени все формулы одной общей формулой.', '</p> \n<p class="text">\nПроанализируй, как связаны между собой величины каждой строки. Запиши формулу зависимости между ними.\n<br/>\nЧто общего у всех записанных формул? Замени все формулы одной общей формулой.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica52-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 52, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 52, номер 1, год 2022."/>\n</div>\n</div>', '1 – s = v · t 2 – C = a · n 3 – F = w · t 4 – S = a · b 5 – V = a · t 6 – K = k · n 7 – T = t · n 8 – M = m · n 9 – P = h · n Формула зависимости всех выражается через произведение. А = b · c.', '<p>\n1 – s = v · t<br/>\n2 – C = a · n<br/>\n3 – F = w · t<br/>\n4 – S = a · b<br/>\n5 – V = a · t<br/>\n6 – K = k · n<br/>\n7 – T = t · n<br/>\n8 – M = m · n<br/>\n9 – P = h · n<br/>\nФормула зависимости всех выражается через произведение. А = b · c.\n\n</p>', 'Формула произведения Формулы зависимостей между величинами – такие, как формула пути (s = v · t), формула стоимости (C = a · n), формула работы (A = w · t) и др., – можно записать одной общей формулой: a = b · c Эту общую формулу мы будем называть формулой произведения. Величины b и c в формуле произведения можно найти по общему правилу нахождения неизвестного множителя: b = a : c      c = a : b', '<div class="recomended-block">\n<span class="title">Формула произведения</span>\n<p>\nФормулы зависимостей между величинами – такие, как формула пути (s = v · t), формула стоимости (C = a · n), формула работы (A = w · t) и др., – можно записать одной общей формулой:<br/>\n<b class="black">a = b · c</b><br/>\nЭту общую формулу мы будем называть <b class="black">формулой произведения.</b><br/>\nВеличины b и c в формуле произведения можно найти по общему правилу нахождения неизвестного множителя:<br/>\nb = a : c      c = a : b\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica52-spravka.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 52, справка, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 52, справка, год 2022."/>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-52/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica52-nomer1.jpg', 'peterson/3/part3/page52/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e534403306fa1f57c79f3447155819ae9649ef398d8266ab0972e8ec6b17263c', NULL, NULL, 'проанализируй, как связаны между собой величины каждой строки. запиши формулу зависимости между ними. что общего у всех записанных формул? замени все формулы одной общей формулой');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 53, '2', 0, 'Реши задачи. Сравни их условия и решения. а) Турист прошёл в первый день 32 км, а во второй – 24 км. Всего он шёл в эти 2 дня 14 часов. Сколько времени шёл турист в каждый из этих дней, если его скорость не изменялась? б) Первый мастер сделал 32 игрушки, а второй – 24 игрушки. На всю эту работу в сумме они затратили 14 часов. Сколько времени работал каждый мастер, если их производительность одинаковая? в) Две подружки из Цветограда купили вместе 14 одинаковых воздушных шариков. Первая уплатила за свою покупку 32 монеты, а вторая – 24 монеты. Всего они купили 14 шариков. Сколько шариков купила каждая из подруг? г) Из двух отрезов шёлка сшили 14 одинаковых юбок. В первом отрезе было 32 м, а во втором – 24 м. Сколько юбок сшили из каждого отреза? Что ты замечаешь? Как это можно объяснить?', '</p> \n<p class="text">Реши задачи. Сравни их условия и решения.<br/> \nа) Турист прошёл в первый день 32 км, а во второй – 24 км. Всего он шёл в эти 2 дня 14 часов. Сколько времени шёл турист в каждый из этих дней, если его скорость не изменялась?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica53-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 53, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 53, номер 2, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Первый мастер сделал 32 игрушки, а второй – 24 игрушки. На всю эту работу в сумме они затратили 14 часов. Сколько времени работал каждый мастер, если их производительность одинаковая?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica53-nomer2-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 53, номер 2-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 53, номер 2-1, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">в) Две подружки из Цветограда купили вместе 14 одинаковых воздушных шариков. Первая уплатила за свою покупку 32 монеты, а вторая – 24 монеты. Всего они купили 14 шариков. Сколько шариков купила каждая из подруг?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica53-nomer2-2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 53, номер 2-2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 53, номер 2-2, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">г) Из двух отрезов шёлка сшили 14 одинаковых юбок. В первом отрезе было 32 м, а во втором – 24 м. Сколько юбок сшили из каждого отреза?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica53-nomer2-3.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 53, номер 2-3, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 53, номер 2-3, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">Что ты замечаешь? Как это можно объяснить?</p>', 'а) (32 + 24) : 14 = 56 : 14 = 4 (км/ч) 32 : 4 = 8 (ч) – первый день 24 : 4 = 6 (ч) – второй день Ответ: 8 часов в первый и 6 часов во второй день шёл турист, если его скорость не изменялась. б) (32 + 24) : 14 = 56 : 14 = 4 (км/ч) 32 : 4 = 8 (ч) – первый день 24 : 4 = 6 (ч) – второй день Ответ: 8 часов первый и 6 часов второй мастер работал, если их производительность одинаковая. в) (32 + 24) : 14 = 56 : 14 = 4 (монет) 32 : 4 = 8 (ш.) – первый день 24 : 4 = 6 (ш.) – второй день Ответ: 8 шариков первая и 6 шариков вторая подруга купила. г) Из двух отрезов шёлка сшили 14 одинаковых юбок. В первом отрезе было 32 м, а во втором – 24 м. (32 + 24) : 14 = 56 : 14 = 4 (м) 32 : 4 = 8 (ю.) – первый день 24 : 4 = 6 (ю.) – второй день Ответ: 8 юбок из первого и 6 юбок из второго отреза сшили. Заметно что действия одинаковые и числа результата одинаковы. Это можно объяснить тем что числа одинаковые в условиях и условия поиска основаны на поиске по общему правилу нахождения неизвестного множителя.', '<p>\nа) (32 + 24) : 14 = 56 : 14 = 4 (км/ч)<br/>\n32 : 4 = 8 (ч) – первый день<br/>\n24 : 4 = 6 (ч) – второй день<br/>\n<b>Ответ:</b> 8 часов в первый и 6 часов во второй день шёл турист, если его скорость не изменялась.<br/><br/>\n\nб) (32 + 24) : 14 = 56 : 14 = 4 (км/ч)<br/>\n32 : 4 = 8 (ч) – первый день<br/>\n24 : 4 = 6 (ч) – второй день<br/>\n<b>Ответ:</b> 8 часов первый и 6 часов второй мастер работал, если их производительность одинаковая.<br/><br/>\n\nв) (32 + 24) : 14 = 56 : 14 = 4 (монет)<br/>\n32 : 4 = 8 (ш.) – первый день<br/>\n24 : 4 = 6 (ш.) – второй день<br/>\n<b>Ответ:</b> 8 шариков первая и 6 шариков вторая подруга купила.<br/><br/>\n\nг) Из двух отрезов шёлка сшили 14 одинаковых юбок. В первом отрезе было 32 м, а во втором – 24 м. (32 + 24) : 14 = 56 : 14 = 4 (м)<br/>\n32 : 4 = 8 (ю.) – первый день<br/>\n24 : 4 = 6 (ю.) – второй день<br/>\n<b>Ответ:</b> 8 юбок из первого и 6 юбок из второго отреза сшили.<br/><br/>\n\nЗаметно что действия одинаковые и числа результата одинаковы. Это можно объяснить тем что числа одинаковые в условиях и условия поиска основаны на поиске по общему правилу нахождения неизвестного множителя.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-53/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica53-nomer2.jpg', 'peterson/3/part3/page53/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica53-nomer2-1.jpg', 'peterson/3/part3/page53/task2_condition_1.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 2, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica53-nomer2-2.jpg', 'peterson/3/part3/page53/task2_condition_2.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 3, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica53-nomer2-3.jpg', 'peterson/3/part3/page53/task2_condition_3.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'abb10deb8e9772f9b64b2950c55dca9b5551eee71f235579d539b8dd9c0cdec7', '2,14,24,32', '["реши","сравни"]'::jsonb, 'реши задачи. сравни их условия и решения. а) турист прошёл в первый день 32 км, а во второй-24 км. всего он шёл в эти 2 дня 14 часов. сколько времени шёл турист в каждый из этих дней, если его скорость не изменялась? б) первый мастер сделал 32 игрушки, а второй-24 игрушки. на всю эту работу в сумме они затратили 14 часов. сколько времени работал каждый мастер, если их производительность одинаковая? в) две подружки из цветограда купили вместе 14 одинаковых воздушных шариков. первая уплатила за свою покупку 32 монеты, а вторая-24 монеты. всего они купили 14 шариков. сколько шариков купила каждая из подруг? г) из двух отрезов шёлка сшили 14 одинаковых юбок. в первом отрезе было 32 м, а во втором-24 м. сколько юбок сшили из каждого отреза? что ты замечаешь? как это можно объяснить');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 54, '3', 0, 'Реши задачи. Для каждой из них придумай задачу с другими величинами, которая решается так же. а) Алёша купил 3 календарика по 8 р. и 7 открыток по 12 р. за штуку. Сколько всего денег заплатил Алёша? б) Фрегат проплыл сначала 2 ч, а потом ещё 4 ч с той же скоростью. Всего он проплыл 216 км. С какой скоростью он плыл? в) Дима почистил 12 картофелин за 6 мин, а Ира – 15 картофелин за 5 мин. Кто из них чистит картошку быстрее и на сколько?', '</p> \n<p class="text">Реши задачи. Для каждой из них придумай задачу с другими величинами, которая решается так же.</p> \n\n<p class="description-text"> \nа) Алёша купил 3 календарика по 8 р. и 7 открыток по 12 р. за штуку. Сколько всего денег заплатил Алёша?<br/>\nб) Фрегат проплыл сначала 2 ч, а потом ещё 4 ч с той же скоростью. Всего он проплыл 216 км. С какой скоростью он плыл?<br/>\nв) Дима почистил 12 картофелин за 6 мин, а Ира – 15 картофелин за 5 мин. Кто из них чистит картошку быстрее и на сколько?\n\n</p>', 'а) 3 · 8 + 7 · 12 = 24 + 84 = 108 (р.) Ответ: 108 рублей всего денег заплатил Алёша. б) 216 : (2 + 4) = 216 : 6 = 36 (км/ч) Ответ: 36 км/ч он плыл. в) 15 : 5 - 12 : 6 = 75 - 72 = 3 (к./мин) Ответ: Ира из них чистит картошку быстрее и на 3 картошки в минуту.', '<p>\nа) 3 · 8 + 7 · 12 = 24 + 84 = 108 (р.)<br/>\n<b>Ответ:</b> 108 рублей всего денег заплатил Алёша.<br/><br/>\nб) 216 : (2 + 4) = 216 : 6 = 36 (км/ч)<br/>\n<b>Ответ:</b> 36 км/ч он плыл.<br/><br/>\nв) 15 : 5 - 12 : 6 = 75 - 72 = 3 (к./мин) <br/>\n<b>Ответ:</b> Ира из них чистит картошку быстрее и на 3 картошки в минуту.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-54/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b89ad31c07a04f7115997b998fba4e081b04d9ddfa6984b0fe2ea7e0644e22be', '2,3,4,5,6,7,8,12,15,216', '["реши"]'::jsonb, 'реши задачи. для каждой из них придумай задачу с другими величинами, которая решается так же. а) алёша купил 3 календарика по 8 р. и 7 открыток по 12 р. за штуку. сколько всего денег заплатил алёша? б) фрегат проплыл сначала 2 ч, а потом ещё 4 ч с той же скоростью. всего он проплыл 216 км. с какой скоростью он плыл? в) дима почистил 12 картофелин за 6 мин, а ира-15 картофелин за 5 мин. кто из них чистит картошку быстрее и на сколько');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 54, '4', 1, 'Составь программу действий и вычисли: а) (154800 : 10 : 9 - 47 · 6) · (97840 : 80 + 77) б) 76000 · 90 : 1000 - 96 : (48 : 8) · 109 - 5400 : 600', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) (154800 : 10 : 9 - 47 · 6) · (97840 : 80 + 77)<br/>\nб) 76000 · 90 : 1000 - 96 : (48 : 8) · 109 - 5400 : 600\n</p>', 'а) (154800 : 10 : 9 - 47 · 6) · (97840 : 80 + 77) = 1869400 154800 : 10 = 15480 15480 : 9 = 1720 47 · 6 = 282 1720 - 282 = 1438 97840 : 80 = 1223 1223 + 77 = 1300 1438 · 1300 = 1869400', '<p>\nа) (154800 : 10 : 9 - 47 · 6) · (97840 : 80 + 77) = 1869400<br/>\n154800 : 10 = 15480<br/>\n15480 : 9 = 1720<br/>\n47 · 6 = 282<br/>\n1720 - 282 = 1438<br/>\n97840 : 80 = 1223<br/>\n1223 + 77 = 1300<br/>\n1438 · 1300 = 1869400\n</p>\n\n<div class="img-wrapper-460">\n<img width="160" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica54-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 54, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 54, номер 4, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-54/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica54-nomer4.jpg', 'peterson/3/part3/page54/task4_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f9c8a9f26c570657ea61dc7e54f161df49ff0dfd8bb670dc0c18561ae47fdfc4', '6,8,9,10,47,48,77,80,90,96', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) (154800:10:9-47*6)*(97840:80+77) б) 76000*90:1000-96:(48:8)*109-5400:600');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 54, '5', 2, 'Запиши множества делителей чисел 7 и 31. Что общего у этих двух множеств? Придумай своё число, множество делителей которого обладает тем же свойством.', '</p> \n<p class="text">Запиши множества делителей чисел 7 и 31. Что общего у этих двух множеств? Придумай своё число, множество делителей которого обладает тем же свойством.</p>', 'Делители 7: 1; 7. Делители 31: 1; 31. У чисел 7 и 31 есть общий делитель - число 1. Делители 3: 1; 3.', '<p>\nДелители 7: 1; 7. Делители 31: 1; 31. У чисел 7 и 31 есть общий делитель - число 1. Делители 3: 1; 3.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-54/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5dd380e5b8ea03025789b578ef657ad78dccd606a9e97ff8f454116fb9c2e274', '7,31', NULL, 'запиши множества делителей чисел 7 и 31. что общего у этих двух множеств? придумай своё число, множество делителей которого обладает тем же свойством');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 54, '6', 3, 'Реши уравнения с комментированием и сделай проверку: а) (3 · m - 20) : 5 = 50 б) 480 : (13 - t) + 20 = 100', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) (3 · m - 20) : 5 = 50<br/>     \nб) 480 : (13 - t) + 20 = 100\n</p>', 'а) (3 · m – 20) : 5 = 50 Чтобы найти делимое (3 · m - 20) надо делитель умножить частное (3 · m - 20) = 50 · 5 (3 · m - 20) = 250 Чтобы найти уменьшаемое 3 · m надо вычитаемое прибавить к разности 3 · m = 250 + 20 3 · m = 270 Чтобы найти множитель надо произведение разделить на известный множитель m = 270 : 3 m = 90 Проверка: (3 · 90 – 20) : 5 = 50 б) 480 : (13 - t) + 20 = 100 Чтобы найти слагаемое 480 : (13 - t) надо из суммы вычесть известное слагаемое 480 : (13 - t) = 100 - 20 480 : (13 - t) = 80 Чтобы найти делитель 13 - t надо делимое разделить на частное 13 - t = 480 : 80 13 - t = 6 Что бы найти вычитаемое надо из уменьшаемого вычесть разность t = 13 - 6 t = 7 Проверка: 480 : (13 - 7) + 20 = 100', '<p>\nа) (3 · m – 20) : 5 = 50<br/>   \nЧтобы найти делимое (3 · m - 20) надо делитель умножить частное<br/>\n(3 · m - 20) = 50 · 5<br/>\n(3 · m - 20) = 250<br/>\nЧтобы найти уменьшаемое 3 · m надо вычитаемое прибавить к разности<br/>\n3 · m = 250 + 20<br/>\n3 · m = 270<br/>\nЧтобы найти множитель надо произведение разделить на известный множитель<br/>\nm = 270 : 3<br/>\nm = 90<br/>\n<b>Проверка:</b> (3 · 90 – 20) : 5 = 50<br/><br/>\nб) 480 : (13 - t) + 20 = 100<br/>\nЧтобы найти слагаемое 480 : (13 - t) надо из суммы вычесть известное слагаемое<br/>\n480 : (13 - t) = 100 - 20<br/>\n480 : (13 - t) = 80<br/>\nЧтобы найти делитель 13 - t надо делимое разделить на частное<br/>\n13 - t = 480 : 80 <br/>\n13 - t = 6<br/>\nЧто бы найти вычитаемое надо из уменьшаемого вычесть разность<br/>\nt = 13 - 6<br/>\nt = 7<br/>\n<b>Проверка:</b> 480 : (13 - 7) + 20 = 100\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-54/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '06bcb54c84571dd4096c60088616f53eba73736a7ad461c3316a2952d6fbbe22', '3,5,13,20,50,100,480', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) (3*m-20):5=50 б) 480:(13-t)+20=100');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 54, '7', 4, 'Выполни действия. Расположи ответы примеров в порядке убывания и расшифруй название игры. Узнай, как играют в эту игру.', '</p> \n<p class="text">Выполни действия. Расположи ответы примеров в порядке убывания и расшифруй название игры. Узнай, как играют в эту игру.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="250" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica54-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 54, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 54, номер 7, год 2022."/>\n</div>\n</div>', 'Б - 4700 · 750 = 3525000 171, 252, 333 Петр – 4, 4 : 2 = 2 Иван – 2, Михаил – 2 и Герасим – 2, 2 : 2 = 1 Яков 1 4 + 2 + 2 + 2 + 1 = 11 Ответ: Пасти овец должен Петр 4 дня, Иван 2 дня, Михаил 2 дня, Герасим 2 дня и Яков 1 день.', '<p>\nБ - 4700 · 750 = 3525000\n</p>\n\n<div class="img-wrapper-460">\n<img width="160" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica54-nomer7-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 54, номер 7-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 54, номер 7-1, год 2022."/>\n\n\n<p>\n171, 252, 333\n</p>\n\n\n<p>\nПетр – 4, 4 : 2 = 2<br/>\nИван – 2, Михаил – 2 и Герасим – 2, 2 : 2 = 1 <br/>\nЯков 1 <br/>\n4 + 2 + 2 + 2 + 1 = 11<br/>\n<b>Ответ:</b> Пасти овец должен Петр 4 дня, Иван 2 дня, Михаил 2 дня, Герасим 2 дня и Яков 1 день.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-54/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica54-nomer7.jpg', 'peterson/3/part3/page54/task7_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica54-nomer7-1.jpg', 'peterson/3/part3/page54/task7_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '20c47702901bd634b343e0a63f1ea5cd74baf16e9a765a51753c50d929fe13d4', NULL, NULL, 'выполни действия. расположи ответы примеров в порядке убывания и расшифруй название игры. узнай, как играют в эту игру');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 55, '1', 0, 'Реши уравнения и сделай проверку: а) 3600 : (18 - x) - 120 = 280 б) (y : 8 + 18) · 9 = 540', '</p> \n<p class="text">Реши уравнения и сделай проверку:</p> \n\n<p class="description-text"> \nа) 3600 : (18 - x) - 120 = 280 <br/>\nб) (y : 8 + 18) · 9 = 540\n</p>', 'а) 3600 : (18 - x) - 120 = 280 3600 : (18 - x) = 280 + 120 3600 : (18 - x) = 400 18 - x = 3600 : 400 18 - х = 9 х = 18 - 9 х = 9 Проверка: 3600 : (18 - 9) - 120 = 280 б) (y : 8 + 18) · 9 = 540 (y : 8 + 18) = 540 : 9 (y : 8 + 18) = 60 у : 8 = 60 - 18 у : 8 = 42 у = 42 · 8 у = 336 Проверка: (336 : 8 + 18) · 9 = 540', '<p>\nа) 3600 : (18 - x) - 120 = 280 <br/>\n3600 : (18 - x) = 280 + 120<br/>\n3600 : (18 - x) = 400<br/>\n18 - x = 3600 : 400<br/>\n18 - х = 9<br/>\nх = 18 - 9<br/>\nх = 9<br/>\n<b>Проверка:</b> 3600 : (18 - 9) - 120 = 280<br/><br/>\nб) (y : 8 + 18) · 9 = 540<br/>\n(y : 8 + 18) = 540 : 9<br/>\n(y : 8 + 18) = 60<br/>\nу : 8 = 60 - 18<br/>\nу : 8 = 42<br/>\nу = 42 · 8<br/>\nу = 336<br/>\n<b>Проверка:</b> (336 : 8 + 18) · 9 = 540\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-55/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e9cf87c333df0501f15ecf7379d7ff5672eb288338d08f56290b2ce83749f39b', '8,9,18,120,280,540,3600', '["реши"]'::jsonb, 'реши уравнения и сделай проверку:а) 3600:(18-x)-120=280 б) (y:8+18)*9=540');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 55, '2', 1, 'БЛИЦтурнир а) Строитель уложил m кирпичей за 4 ч. За сколько времени он уложит d кирпичей, если будет работать с той же производительностью? б) Самолёт пролетел s км за 2 ч, а вертолёт пролетел это же расстояние за 3 ч. На сколько скорость самолёта больше скорости вертолёта? в) За 6 м льняной ткани заплатили k р. А один метр шёлка на n р. дороже метра льняной ткани. Чему равна цена метра шёлка? г) Мастеру надо было изготовить a деталей. Он уже сделал b деталей. Чему должна быть равна его производительность, чтобы он успел сделать оставшиеся детали за t часов? Для одной из данных задач придумай задачу с другими величинами, которые решаются так же.', '</p> \n<p class="text">БЛИЦтурнир<br/>\nа) Строитель уложил m кирпичей за 4 ч. За сколько времени он уложит d кирпичей, если будет работать с той же производительностью?<br/>\nб) Самолёт пролетел s км за 2 ч, а вертолёт пролетел это же расстояние за 3 ч. На сколько скорость самолёта больше скорости вертолёта?<br/>\nв) За 6 м льняной ткани заплатили k р. А один метр шёлка на n р. дороже метра льняной ткани. Чему равна цена метра шёлка?<br/>\nг) Мастеру надо было изготовить a деталей. Он уже сделал b деталей. Чему должна быть равна его производительность, чтобы он успел сделать оставшиеся детали за t часов? Для одной из данных задач придумай задачу с другими величинами, которые решаются так же.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica55-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 55, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 55, номер 2, год 2022."/>\n</div>\n</div>', 'а) m : 4 · d (кирпичей/ч) Ответ: за m : 4 · d времени он уложит d кирпичей, если будет работать с той же производительностью. б) s : 2 - s : 3 (км/ч) Ответ: на s : 2 - s : 3 км/ч скорость самолёта больше скорости вертолёта. в) 6 : k + n (р.) Ответ: 6 : k + n р. равна цена метра шёлка. г) (a – b) : t (деталей/ч) Ответ: (a – b) : t деталей в час должна быть равна его производительность, чтобы он успел сделать оставшиеся детали за t часов. Придуманная задача: Велосипедисту надо было проехать a км. Он уже проехал b км. Чему должна быть равна его скорость, чтобы он успел проехать оставшиеся км за t часов? (a – b) : t (км/ч) Ответ: (a – b) : t километров в час должна быть равна его скорость, чтобы он успел проехать оставшиеся км за t часов.', '<p>\nа) m : 4 · d (кирпичей/ч) <br/>\n<b>Ответ:</b> за m : 4 · d времени он уложит d кирпичей, если будет работать с той же производительностью.<br/><br/>\nб) s : 2 - s : 3 (км/ч)<br/>\n<b>Ответ:</b> на s : 2 - s : 3 км/ч скорость самолёта больше скорости вертолёта.<br/><br/>\nв) 6 : k + n (р.)<br/>\n<b>Ответ:</b> 6 : k + n р. равна цена метра шёлка.<br/><br/>\nг) (a – b) : t (деталей/ч) <br/>\n<b>Ответ:</b> (a – b) : t деталей в час должна быть равна его производительность, чтобы он успел сделать оставшиеся детали за t часов.<br/><br/>\nПридуманная задача:<br/>\nВелосипедисту надо было проехать a км. Он уже проехал b км. Чему должна быть равна его скорость, чтобы он успел проехать оставшиеся км за t часов? <br/>\n(a – b) : t (км/ч) <br/>\n<b>Ответ:</b> (a – b) : t километров в час должна быть равна его скорость, чтобы он успел проехать оставшиеся км за t часов.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-55/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica55-nomer2.jpg', 'peterson/3/part3/page55/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1b69343345074ec78e36abe5b3ee7c8fed2036872f7dd0bc0dc2f15381d7ac2b', '2,3,4,6', '["больше"]'::jsonb, 'блицтурнир а) строитель уложил m кирпичей за 4 ч. за сколько времени он уложит d кирпичей, если будет работать с той же производительностью? б) самолёт пролетел s км за 2 ч, а вертолёт пролетел это же расстояние за 3 ч. на сколько скорость самолёта больше скорости вертолёта? в) за 6 м льняной ткани заплатили k р. а один метр шёлка на n р. дороже метра льняной ткани. чему равна цена метра шёлка? г) мастеру надо было изготовить a деталей. он уже сделал b деталей. чему должна быть равна его производительность, чтобы он успел сделать оставшиеся детали за t часов? для одной из данных задач придумай задачу с другими величинами, которые решаются так же');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 55, '3', 2, 'Вырази в указанных единицах измерения: а) 5 дм 6 мм = … мм      б) 5 ц 6 кг = … кг 5 м 6 см = … мм             5 кг 6 г = … г 5 км 6 м = … м               5 сут. 6 ч = … ч 5 км 6 м = … дм             5 мин 6 с = … с 5 км 6 м = … мм             5 ч 6 мин = … мин', '</p> \n<p class="text">Вырази в указанных единицах измерения:</p> \n\n<p class="description-text"> \nа)  5 дм 6 мм = … мм      б)  5 ц 6 кг = … кг <br/> \n5 м 6 см = … мм             5 кг 6 г = … г  <br/>\n5 км 6 м = … м               5 сут. 6 ч = … ч <br/>  \n5 км 6 м = … дм             5 мин 6 с = … с  <br/> \n5 км 6 м = … мм             5 ч 6 мин = … мин\n\n</p>', 'а) 5 дм 6 мм = 506 мм        б) 5 ц 6 кг = 506 кг 5 м 6 см = 5010 мм             5 кг 6 г = 5006 г 5 км 6 м = 5006 м               5 сут. 6 ч = 126 ч 5 км 6 м = 50060 дм           5 мин 6 с = 306 с 5 км 6 м = 5006000 мм      5 ч 6 мин = 306 мин', '<p>\nа)  5 дм 6 мм = 506 мм        б)  5 ц 6 кг = 506 кг<br/>  \n5 м 6 см = 5010 мм             5 кг 6 г = 5006 г  <br/>\n5 км 6 м = 5006 м               5 сут. 6 ч = 126 ч  <br/> \n5 км 6 м = 50060 дм           5 мин 6 с = 306 с  <br/> \n5 км 6 м = 5006000 мм      5 ч 6 мин = 306 мин\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-55/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '096da3d7a2e1e605149657f2298082aaef540223577ed36b694894b00d5b04b9', '5,6', '["раз"]'::jsonb, 'вырази в указанных единицах измерения:а) 5 дм 6 мм=... мм      б) 5 ц 6 кг=... кг 5 м 6 см=... мм             5 кг 6 г=... г 5 км 6 м=... м               5 сут. 6 ч=... ч 5 км 6 м=... дм             5 мин 6 с=... с 5 км 6 м=... мм             5 ч 6 мин=... мин');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 55, '4', 3, 'Выполни действия: а) 3 т 2 ц 6 кг - 29 ц 48 кг        в) 9 мин 15 с · 8 б) 5 км 19 м + 1 км 981 м         г) 6 м 1 дм 2 мм : 3', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 3 т 2 ц 6 кг - 29 ц 48 кг        в) 9 мин 15 с · 8<br/>\nб) 5 км 19 м + 1 км 981 м         г) 6 м 1 дм 2 мм : 3\n\n</p>', 'а) 3 т 2 ц 6 кг - 29 ц 48 кг = 3206 кг - 2948 кг = 258 кг = 2 ц 58 кг', '<p>\nа) 3 т 2 ц 6 кг - 29 ц 48 кг = 3206 кг - 2948 кг = 258 кг = 2 ц 58 кг   \n</p>\n\n<div class="img-wrapper-460">\n<img width="120" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica55-nomer4.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 55, номер 4, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 55, номер 4, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-55/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica55-nomer4.jpg', 'peterson/3/part3/page55/task4_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b3c39973cabc9e2abe402d469400d546a76459d35579d0e03e8597a374bd6371', '1,2,3,5,6,8,9,15,19,29', NULL, 'выполни действия:а) 3 т 2 ц 6 кг-29 ц 48 кг        в) 9 мин 15 с*8 б) 5 км 19 м+1 км 981 м         г) 6 м 1 дм 2 мм:3');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 55, '5', 4, 'Для сада купили в питомнике 14 кустов красной и чёрной смородины по одинаковой цене. За красную смородину заплатили 250 р., а за чёрную – 450 р. Каких кустов купили больше и на сколько?', '</p> \n<p class="text">Для сада купили в питомнике 14 кустов красной и чёрной смородины по одинаковой цене. За красную смородину заплатили 250 р., а за чёрную – 450 р. Каких кустов купили больше и на сколько?</p>', '(250 + 450) : 14 = 700 : 14 = 50 (р.) 250 : 50 = 5 (кустов) – красная смородина 450 : 50 = 9 (кустов) – чёрная смородина 9 - 5 = 4 (куста) Ответ: чёрной смородины кустов купили больше на 4 куста.', '<p>\n(250 + 450) : 14 = 700 : 14 = 50 (р.)\n</p>\n\n\n\n<p>\n250 : 50 = 5 (кустов) – красная смородина<br/>\n450 : 50 = 9 (кустов) – чёрная смородина<br/>\n9 - 5 = 4 (куста)<br/>\n<b>Ответ:</b> чёрной смородины кустов купили больше на 4 куста.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-55/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '9d4dcb3ab17f45e105682cbb82b2d6b4860ff676fb6ad09d0c7ef589d25215bb', '14,250,450', '["больше"]'::jsonb, 'для сада купили в питомнике 14 кустов красной и чёрной смородины по одинаковой цене. за красную смородину заплатили 250 р., а за чёрную-450 р. каких кустов купили больше и на сколько');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 56, '6', 0, 'В пошивочной мастерской в первый день сшили 24 одинаковых комплекта белья, а во второй – на два таких комплекта больше. На все комплекты было израсходовано за два дня 800 м ткани. Сколько метров ткани израсходовали в каждый из этих дней?', '</p> \n<p class="text">В пошивочной мастерской в первый день сшили 24 одинаковых комплекта белья, а во второй – на два таких комплекта больше. На все комплекты было израсходовано за два дня 800 м ткани. Сколько метров ткани израсходовали в каждый из этих дней?</p>', '24 + 2 = 26 (комплекта) – 2 день 800 : (24 + 26) = 800 : 50 = 16 (м) – 1 комплект 24 · 16 = 364 (м) – 1 день', '<p>\n24 + 2 = 26 (комплекта) – 2 день<br/>\n800 : (24 + 26) = 800 : 50 = 16 (м) – 1 комплект<br/>\n24 · 16 = 364 (м) – 1 день\n</p>\n\n<div class="img-wrapper-460">\n<img width="80" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica56-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 56, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 56, номер 6, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-56/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica56-nomer6.jpg', 'peterson/3/part3/page56/task6_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '87fafe60151d8b9641f4cdc4eb46a0ecf4356c528ed991cb6ab38204f7e9a00c', '24,800', '["больше"]'::jsonb, 'в пошивочной мастерской в первый день сшили 24 одинаковых комплекта белья, а во второй-на два таких комплекта больше. на все комплекты было израсходовано за два дня 800 м ткани. сколько метров ткани израсходовали в каждый из этих дней');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 56, '7', 1, 'Найди значения выражений: а) 270 : 9 · 7 - 360 : (16 : 4) + (42 : 7 · 6 + 14) б) 125 · 0 : (45 · 4) + (120 · 10 : 100 - 8) · (15 · 1000 : 5)', '</p> \n<p class="text">Найди значения выражений:</p> \n\n<p class="description-text"> \nа) 270 : 9 · 7 - 360 : (16 : 4) + (42 : 7 · 6 + 14)<br/>\nб) 125 · 0 : (45 · 4) + (120 · 10 : 100 - 8) · (15 · 1000 : 5)\n</p>', 'а) 270 : 9 · 7 - 360 : (16 : 4) + (42 : 7 · 6 + 14) = 30 · 7 - 360 : 4 + (6 · 6 + 14) = 210 - 90 + (36 + 14) = 120 + 50 = 170 б) 125 · 0 : (45 · 4) + (120 · 10 : 100 – 8) · (15 · 1000 : 5) = 0 : 180 + (1200 : 100 - 8) · (15000 : 5) = 0 + (12 - 8) · 3000 = 4 · 3000 = 12000', '<p>\nа) 270 : 9 · 7 - 360 : (16 : 4) + (42 : 7 · 6 + 14) = 30 · 7 - 360 : 4 + (6 · 6 + 14) = 210 - 90 + (36 + 14) = 120 + 50 = 170<br/>\nб) 125 · 0 : (45 · 4) + (120 · 10 : 100 – 8) · (15 · 1000 : 5) = 0 : 180 + (1200 : 100 - 8) · (15000 : 5) = 0 + (12 - 8) · 3000 = 4 · 3000 = 12000 \n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-56/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '235ed9af6e973ca2c96963cd44ab1a5f0daa6d82334a6da2c82f4bae2d928702', '0,4,5,6,7,8,9,10,14,15', '["найди"]'::jsonb, 'найди значения выражений:а) 270:9*7-360:(16:4)+(42:7*6+14) б) 125*0:(45*4)+(120*10:100-8)*(15*1000:5)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 56, '8', 2, 'Выполни умножение. Найди сумму и разность самого большого и самого маленького из получившихся чисел: 2590 · 763      9450 · 4560      49300 · 807', '</p> \n<p class="text">Выполни умножение. Найди сумму и разность самого большого и самого маленького из получившихся чисел: </p> \n\n<p class="description-text"> \n2590 · 763      9450 · 4560      49300 · 807\n</p>', '2590 · 763 = 1976170', '<p>\n2590 · 763 = 1976170\n</p>\n\n<div class="img-wrapper-460">\n<img width="190" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica56-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 56, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 56, номер 8, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-56/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica56-nomer8.jpg', 'peterson/3/part3/page56/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '61a1a0ea4ad2265e0fd917469b22fb2dad1f06aadf793b32bef2daf1c415e56c', '763,807,2590,4560,9450,49300', '["найди","разность","раз"]'::jsonb, 'выполни умножение. найди сумму и разность самого большого и самого маленького из получившихся чисел:2590*763      9450*4560      49300*807');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 56, '9', 3, 'Пересекаются ли: а) прямая l и луч AB; б) прямая l и луч TS; в) прямая l и отрезок MK; г) прямая l и отрезок CD; д) лучи AB и TS; е) отрезки MK и CD; ж) луч TS и отрезок MK; з) луч TS и отрезок EF?', '</p> \n<p class="text">\nПересекаются ли: а) прямая l и луч AB; б) прямая l и луч TS; в) прямая l и отрезок MK; г) прямая l и отрезок CD; д) лучи AB и TS; е) отрезки MK и CD; ж) луч TS и отрезок MK; з) луч TS и отрезок EF?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica56-nomer9.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 56, номер 9, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 56, номер 9, год 2022."/>\n</div>\n</div>', 'а) прямая l и луч AB пересекаются; б) прямая l и луч TS не пересекаются; в) прямая l и отрезок MK не пересекаются; г) прямая l и отрезок CD не пересекаются; д) лучи AB и TS пересекаются; е) отрезки MK и CD не пересекаются; ж) луч TS и отрезок MK не пересекаются; з) луч TS и отрезок EF не пересекаются. а) 16 = 4 · 4 4 + 4 = 8 16 = 2 · 8 2 + 8 = 10 16 = 1 · 16 1 + 16 = 17 В случае 16 = 4 · 4 получилась наименьшая сумма. б) 36 = 6 · 6 6 + 6 = 12 36 = 12 · 3 12 + 3 = 15 36 = 36 · 1 36 + 1 = 37 В случае 36 = 6 · 6 получилась наименьшая сумма, 64 = 8 · 8 8 + 8 = 16 64 = 2 · 32 2 + 32 = 34 64 = 64 · 1 64 + 1 = 64 В случае 64 = 8 · 8 получилась наименьшая сумма. Можно высказать предположение (гипотезу), что в случае табличного произведения получается наименьшая сумма. Так как сторона зелёного кубика в два раза больше стороны красного, то Роману построившему большой куб из 64 красных кубиков потребуется в два раза меньше зелёных кубиков, чтобы построить точно такой же куб. Ответ: в два раза меньше нужно зелёных кубиков, чтобы построить точно такой же куб.', '<p>\nа) прямая l и луч AB пересекаются; б) прямая l и луч TS не пересекаются; в) прямая l и отрезок MK не пересекаются; г) прямая l и отрезок CD не пересекаются; д) лучи AB и TS пересекаются; е) отрезки MK и CD не пересекаются; ж) луч TS и отрезок MK не пересекаются; з) луч TS и отрезок EF не пересекаются. \n</p>\n\n\n<p>\nа) 16 = 4 · 4<br/> \n4 + 4 = 8<br/>\n16 = 2 · 8<br/>\n2 + 8 = 10<br/>\n16 = 1 · 16<br/>\n1 + 16 = 17<br/>\nВ случае 16 = 4 · 4 получилась наименьшая сумма. <br/><br/>\nб) 36 = 6 · 6<br/> \n6 + 6 = 12<br/>\n36 = 12 · 3<br/>\n12 + 3 = 15<br/>\n36 = 36 · 1 <br/>\n36 + 1 = 37<br/>\nВ случае 36 = 6 · 6 получилась наименьшая сумма,<br/>\n64 = 8 · 8<br/> \n8 + 8 = 16<br/>\n64 = 2 · 32 <br/>\n2 + 32 = 34<br/>\n64 = 64 · 1 <br/>\n64 + 1 = 64<br/>\nВ случае 64 = 8 · 8 получилась наименьшая сумма. <br/>\nМожно высказать предположение (гипотезу), что в случае табличного произведения получается наименьшая сумма.\n</p>\n\n\n<p>\nТак как сторона зелёного кубика в два раза больше стороны красного, то Роману построившему большой куб из 64 красных кубиков потребуется в два раза меньше зелёных кубиков, чтобы построить точно такой же куб.\n<b>Ответ:</b> в два раза меньше нужно зелёных кубиков, чтобы построить точно такой же куб.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-56/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica56-nomer9.jpg', 'peterson/3/part3/page56/task9_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'df14d8b67feb04c135b83c0f60268eab68a271c4a7582d50d665b27927a7bbca', NULL, NULL, 'пересекаются ли:а) прямая l и луч ab; б) прямая l и луч ts; в) прямая l и отрезок mk; г) прямая l и отрезок cd; д) лучи ab и ts; е) отрезки mk и cd; ж) луч ts и отрезок mk; з) луч ts и отрезок ef');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 59, '1', 0, 'Определи тип простой задачи и реши её: а) Миша нашёл 24 гриба, а Витя – в 3 раза меньше. Сколько грибов нашёл Витя? б) Таня испекла 15 пирожков. Из них 8 пирожков съели за ужином. Сколько пирожков осталось? в) Дима прошёл за 20 минут 1 км 600 м. С какой скоростью он шёл? г) Лариса посадила в своём цветнике 36 тюльпанов и 42 нарцисса. Каких цветов она посадила больше и на сколько?', '</p> \n<p class="text">Определи тип простой задачи и реши её:<br/>\nа) Миша нашёл 24 гриба, а Витя – в 3 раза меньше. Сколько грибов нашёл Витя? <br/>\nб) Таня испекла 15 пирожков. Из них 8 пирожков съели за ужином. Сколько пирожков осталось?<br/>\nв) Дима прошёл за 20 минут 1 км 600 м. С какой скоростью он шёл?<br/>\nг) Лариса посадила в своём цветнике 36 тюльпанов и 42 нарцисса. Каких цветов она посадила больше и на сколько?\n</p>', 'а) вид зависимости: a = b · c 24 : 3 = 6 Ответ: 6 грибов нашёл Витя. б) вид зависимости: a = b + c 15 - 8 = 7 (пирожков) Ответ: 7 пирожков осталось. в) вид зависимости: a = b · c 1 км 600 м : 20 минут = 1600 м : 20 минут = 80 (м/мин) Ответ: 80 км/мин он шёл. г) вид зависимости: a = b + c 42 - 36 = 6 (цветка) Ответ: нарциссов она посадила больше на 6.', '<p>\nа) вид зависимости: a = b · c<br/>\n24 : 3 = 6<br/>\n<b>Ответ:</b> 6 грибов нашёл Витя. <br/><br/>\nб) вид зависимости: a = b + c <br/>\n15 - 8 = 7 (пирожков) <br/>\n<b>Ответ:</b> 7 пирожков осталось.<br/><br/>\nв) вид зависимости: a = b · c <br/>\n1 км 600 м : 20 минут = 1600 м : 20 минут = 80 (м/мин)<br/>\n<b>Ответ:</b> 80 км/мин он шёл.<br/><br/>\nг) вид зависимости: a = b + c <br/>\n42 - 36 = 6 (цветка)<br/>\n<b>Ответ:</b> нарциссов она посадила больше на 6.\n</p>', 'Способы решения составных задач Встречаются также простые задачи, в которых величины сравниваются: на сколько или во сколько раз одна величина больше (меньше) другой. Правила их решения нам хорошо известны. Таким образом, компас решения простых задач (в одно действие) можно представить так:', '<div class="recomended-block">\n<span class="title">Способы решения составных задач</span>\n<p>\nВстречаются также простые задачи, в которых величины сравниваются: на сколько или во сколько раз одна величина больше (меньше) другой. Правила их решения нам хорошо известны. <br/>\nТаким образом, компас решения простых задач (в одно действие) можно представить так:\n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica59-spravka.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 59, справка, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 59, справка, год 2022."/>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-59/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'bec1aa78389a4278d0ff66485064b18b46baad3b9ec52a83c99891ee466573ad', '1,3,8,15,20,24,36,42,600', '["реши","больше","меньше","раз","раза"]'::jsonb, 'определи тип простой задачи и реши её:а) миша нашёл 24 гриба, а витя-в 3 раза меньше. сколько грибов нашёл витя? б) таня испекла 15 пирожков. из них 8 пирожков съели за ужином. сколько пирожков осталось? в) дима прошёл за 20 минут 1 км 600 м. с какой скоростью он шёл? г) лариса посадила в своём цветнике 36 тюльпанов и 42 нарцисса. каких цветов она посадила больше и на сколько');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 60, '2', 0, 'В магазин привезли 120 кг яблок, груш – в 2 раза меньше, чем яблок, а персиков – на 12 кг больше, чем груш. Сколько всего килограммов яблок, груш и персиков привезли в магазин?', '</p> \n<p class="text">В магазин привезли 120 кг яблок, груш – в 2 раза меньше, чем яблок, а персиков – на 12 кг больше, чем груш. Сколько всего килограммов яблок, груш и персиков привезли в магазин?</p>', '120 + 120 : 2 + (120 : 2 + 12) = 120 + 60 + (60 + 12) = 180 + 72 = 252 (кг) Ответ: 252 всего килограммов яблок, груш и персиков привезли в магазин.', '<p>\n120 + 120 : 2 + (120 : 2 + 12) = 120 + 60 + (60 + 12) = 180 + 72 = 252 (кг)<br/>\n<b>Ответ:</b> 252 всего килограммов яблок, груш и персиков привезли в магазин.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-60/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c776718525e4b1868c633c4008e49b1a197e8cf34ffab5cbae2df3dbde958483', '2,12,120', '["больше","меньше","раз","раза"]'::jsonb, 'в магазин привезли 120 кг яблок, груш-в 2 раза меньше, чем яблок, а персиков-на 12 кг больше, чем груш. сколько всего килограммов яблок, груш и персиков привезли в магазин');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 60, '3', 1, 'На шоссе стоят четыре восьмиэтажных жилых дома. На каждом этаже каждого из этих домов по 9 квартир. Из всех квартир 128 однокомнатных, 96 двухкомнатных, а остальные – трёхкомнатные. Сколько всего трёхкомнатных квартир в этих домах?', '</p> \n<p class="text">На шоссе стоят четыре восьмиэтажных жилых дома. На каждом этаже каждого из этих домов по 9 квартир. Из всех квартир 128 однокомнатных, 96 двухкомнатных, а остальные – трёхкомнатные. Сколько всего трёхкомнатных квартир в этих домах?</p>', '8 · 9 = 72 (квартир) – в одном доме 72 · 4 = 288 (квартир) – во всех домах 288 - 128 - 96 = 160 - 96 = 64 (квартиры) Ответ: 64 всего трёхкомнатных квартир в этих домах.', '<p>\n8 · 9 = 72 (квартир) – в одном доме<br/>\n72 · 4 = 288 (квартир) – во всех домах<br/>\n288 - 128 - 96 = 160 - 96 = 64 (квартиры)<br/>\n<b>Ответ:</b> 64 всего трёхкомнатных квартир в этих домах.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-60/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '31c35e617c56dfc16c13c740d4da8744062ca14cd99ec8a18beb3a60a34594ad', '9,96,128', NULL, 'на шоссе стоят четыре восьмиэтажных жилых дома. на каждом этаже каждого из этих домов по 9 квартир. из всех квартир 128 однокомнатных, 96 двухкомнатных, а остальные-трёхкомнатные. сколько всего трёхкомнатных квартир в этих домах');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 60, '4', 2, 'Расстояние между Москвой и Минском 720 км. Автомобиль ехал из Москвы в Минск со скоростью 80 км/ч, а на обратном пути – увеличил скорость на 10 км/ч. Сколько времени затратил автомобиль на весь путь из Москвы в Минск и обратно?', '</p> \n<p class="text">Расстояние между Москвой и Минском 720 км. Автомобиль ехал из Москвы в Минск со скоростью 80 км/ч, а на обратном пути – увеличил скорость на 10 км/ч. Сколько времени затратил автомобиль на весь путь из Москвы в Минск и обратно?</p>', '(720 км : 80 км/ч) + (720 км : (80км/ч + 10км/ч)) = 9 + (720 км/ч : 90 км/ч) = 9 ч + 8 ч = 17 (ч) Ответ: 17 часов затратил автомобиль на весь путь из Москвы в Минск и обратно.', '<p>\n(720 км : 80 км/ч) + (720 км : (80км/ч + 10км/ч)) = 9 + (720 км/ч : 90 км/ч) = 9 ч + 8 ч = 17 (ч)<br/>\n<b>Ответ:</b> 17 часов затратил автомобиль на весь путь из Москвы в Минск и обратно.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-60/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a3c4b824e423215d7de4d65415c14e87206166cf55c4734c50ee527f3e9c6307', '10,80,720', NULL, 'расстояние между москвой и минском 720 км. автомобиль ехал из москвы в минск со скоростью 80 км/ч, а на обратном пути-увеличил скорость на 10 км/ч. сколько времени затратил автомобиль на весь путь из москвы в минск и обратно');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 60, '5', 3, 'В первом куске 12 м ткани, а во втором – 8 м такой же ткани. Первый кусок дороже второго на 320 р. Сколько рублей стоит каждый из этих кусков ткани?', '</p> \n<p class="text">В первом куске 12 м ткани, а во втором – 8 м такой же ткани. Первый кусок дороже второго на 320 р. Сколько рублей стоит каждый из этих кусков ткани?</p>', '320 : (12 - 8) = 320 : 4 = 80 (р) – стоимость 1 м ткани 12 · 80 = 960 (р) – первый 8 · 80 = 640 (р) – второй Ответ: 960 рублей первый и 640 рублей второй кусок ткани стоит.', '<p>\n320 : (12 - 8) = 320 : 4 = 80 (р) – стоимость 1 м ткани<br/>\n12 · 80 = 960 (р) – первый<br/>\n8 · 80 = 640 (р) – второй<br/>\n<b>Ответ:</b> 960 рублей первый и 640 рублей второй кусок ткани стоит.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-60/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '7768202d7187ecee26798b1eed577dc293dab4c954a20b61e1719913d9511626', '8,12,320', NULL, 'в первом куске 12 м ткани, а во втором-8 м такой же ткани. первый кусок дороже второго на 320 р. сколько рублей стоит каждый из этих кусков ткани');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 60, '6', 4, 'Сумма площадей двух прямоугольников, имеющих одинаковую длину, равна 220 дм 2 . Ширина первого прямоугольника 4 дм, а ширина второго – на 3 дм больше, чем первого. Чему равна длина этих прямоугольников?', '</p> \n<p class="text">Сумма площадей двух прямоугольников, имеющих одинаковую длину, равна 220 дм<sup>2</sup>. Ширина первого прямоугольника 4 дм, а ширина второго – на 3 дм больше, чем первого. Чему равна длина этих прямоугольников?</p>', '220 : (4 + 4 + 3) = 220 : 11 = 20 (дм) Ответ: 20 дм длина этих прямоугольников.', '<p>\n220 : (4 + 4 + 3) = 220 : 11 = 20 (дм)<br/>\n<b>Ответ:</b> 20 дм длина этих прямоугольников.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-60/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5215568704599875f5e381dae4c256cf8088c2048d9a496d34cfdc5d95166483', '2,3,4,220', '["сумма","больше"]'::jsonb, 'сумма площадей двух прямоугольников, имеющих одинаковую длину, равна 220 дм 2 . ширина первого прямоугольника 4 дм, а ширина второго-на 3 дм больше, чем первого. чему равна длина этих прямоугольников');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 60, '7', 5, 'Выполни действия: а) 374 · 75        в) 850 · 39800      д) 7263000 : 90 б) 908 · 132      г) 4620 · 5040      е) 24040000 : 800', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 374 · 75        в) 850 · 39800      д) 7263000 : 90<br/>\nб) 908 · 132      г)  4620 · 5040      е) 24040000 : 800\n</p>', 'а) 374 · 75 = 28050', '<p>\nа) 374 · 75 = 28050\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica60-nomer7.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 60, номер 7, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 60, номер 7, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-60/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica60-nomer7.jpg', 'peterson/3/part3/page60/task7_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3d3eb20697067a1cc30b317d1a87b26eee2fa0497b9173eb9c3a39b3f543c8b9', '75,90,132,374,800,850,908,4620,5040,39800', NULL, 'выполни действия:а) 374*75        в) 850*39800      д) 7263000:90 б) 908*132      г) 4620*5040      е) 24040000:800');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 60, '8', 6, 'Составь программу действий и вычисли: а) (18560 - 17915) · (4235 : 5 + 9535) б) (600300 - 728 · 604) : 4 · (1700 · 390)', '</p> \n<p class="text">Составь программу действий и вычисли: </p> \n\n<p class="description-text"> \nа) (18560 - 17915) · (4235 : 5 + 9535)<br/>\nб) (600300 - 728 · 604) : 4 · (1700 · 390)\n</p>', 'а) (18560 - 17915) · (4235 : 5 + 9535) = 6696390 18560 - 17915 = 645', '<p>\nа) (18560 - 17915) · (4235 : 5 + 9535) = 6696390<br/>\n18560 - 17915 = 645\n</p>\n\n<div class="img-wrapper-460">\n<img width="140" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica60-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 60, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 60, номер 8, год 2022."/>\n\n\n<div class="img-wrapper-460">\n<img width="180" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica60-nomer9-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 60, номер 9-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 60, номер 9-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-60/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica60-nomer8.jpg', 'peterson/3/part3/page60/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica60-nomer9-1.jpg', 'peterson/3/part3/page60/task8_solution_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '97fa7b1c542810b815782e1ec9fd78af771670aaca5506b742678bdb33e4bd91', '4,5,390,604,728,1700,4235,9535,17915,18560', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) (18560-17915)*(4235:5+9535) б) (600300-728*604):4*(1700*390)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 61, '1', 0, 'Составь задачи по таблицам и реши их. Что ты замечаешь? Придумай и реши аналогичные задачи на движение и стоимость.', '</p> \n<p class="text">Составь задачи по таблицам и реши их. Что ты замечаешь?  Придумай и реши аналогичные задачи на движение и стоимость.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica61-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 61, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 61, номер 1, год 2022."/>\n</div>\n</div>', 'а) 96 штук первого товара выпущено за 8 дней и 60 штук второго товара за оставшиеся дни. За сколько дней выпущен второй товар? С какой одинаковой производительностью выпускаются эти товары? 96 : 8 = 12 (шт./день) 60 : 12 = 5 (дней) Ответ: за 5 дней выпущен второй товар, 12 шт/день производительность выпуска этих товаров. б) Площадь первого прямоугольника 96 см 2 и ширина 8 см. Площадь второго прямоугольника 60 см 2 . Длины этих прямоугольников равны. Какова длина прямоугольников? Какова ширина второго прямоугольника? 96 : 8 = 12 (см) 60 : 12 = 5 (см) Ответ: 12 см длина этих прямоугольников, 5 см ширина второго прямоугольника.', '<p>\nа) 96 штук первого товара выпущено за 8 дней и 60 штук второго товара за оставшиеся дни. За сколько дней выпущен второй товар? С какой одинаковой производительностью выпускаются эти товары?<br/> \n96 : 8 = 12 (шт./день)<br/> \n60 : 12 = 5 (дней)<br/> \n<b>Ответ:</b> за 5 дней выпущен второй товар, 12 шт/день производительность выпуска этих товаров.<br/> <br/> \nб) Площадь первого прямоугольника 96 см<sup>2</sup> и ширина 8 см. Площадь второго прямоугольника 60 см<sup>2</sup>. Длины этих прямоугольников равны. Какова длина прямоугольников? Какова ширина второго прямоугольника?<br/> \n96 : 8 = 12 (см)<br/> \n60 : 12 = 5 (см)<br/> \n<b>Ответ:</b> 12 см длина этих прямоугольников, 5 см ширина второго прямоугольника.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-61/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica61-nomer1.jpg', 'peterson/3/part3/page61/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '99ae91912fc793379378cfc01885f7dc35ac5d1adc6d69fbefaa63ab4cad3e5f', NULL, '["реши"]'::jsonb, 'составь задачи по таблицам и реши их. что ты замечаешь? придумай и реши аналогичные задачи на движение и стоимость');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 61, '2', 1, 'Портниха за 4 дня сшила на 42 комплекта белья меньше, чем за 7 дней. С какой производительностью она работала? Сколько комплектов белья сошьёт эта портниха за 20 дней, если будет работать с той же производительностью?', '</p> \n<p class="text">Портниха за 4 дня сшила на 42 комплекта белья меньше, чем за 7 дней. С какой производительностью она работала? Сколько комплектов белья сошьёт эта портниха за 20 дней, если будет работать с той же производительностью?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica61-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 61, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 61, номер 2, год 2022."/>\n</div>\n</div>', '42 : (7 - 4) = 42 : 3 = 14 (к./день) 14 · 20 = 180 (к.) Ответ: 14 к./день равна производительность портнихи, 180 комплектов белья сошьёт эта портниха за 20 дней, если будет работать с той же производительностью.', '<p>\n42 : (7 - 4) = 42 : 3 = 14 (к./день)<br/> \n14 · 20 = 180 (к.)<br/> \n<b>Ответ:</b> 14 к./день равна производительность портнихи, 180 комплектов белья сошьёт эта портниха за 20 дней, если будет работать с той же производительностью.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-61/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica61-nomer2.jpg', 'peterson/3/part3/page61/task2_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '792176bc5cbf4db668b89d4c0db418e1da459eabf2e735d6900a1a0f6c1c96ff', '4,7,20,42', '["меньше"]'::jsonb, 'портниха за 4 дня сшила на 42 комплекта белья меньше, чем за 7 дней. с какой производительностью она работала? сколько комплектов белья сошьёт эта портниха за 20 дней, если будет работать с той же производительностью');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 61, '3', 2, 'Составь и реши уравнения: а) Задумано число. К нему прибавили 19, сумму умножили на 5 и из полученного произведения вычли 16. Получилось 139. Какое число задумано? б) Задумано число. Его вычли из 480, разность разделили на 6 и полученное частное увеличили на 89. В результате получилось 165. Какое число задумано?', '</p> \n<p class="text">Составь и реши уравнения:<br/>\nа) Задумано число. К нему прибавили 19, сумму умножили на 5 и из полученного произведения вычли 16. Получилось 139. Какое число задумано?<br/>\nб) Задумано число. Его вычли из 480, разность разделили на 6 и полученное частное увеличили на 89. В результате получилось 165. Какое число задумано?\n</p>', 'а) (х + 19) · 5 - 16 = 139 (х + 19) · 5 = 139 + 16 х + 19 = 155 : 5 х + 19 = 31 х = 31 19 х = 12 Ответ: число 12 задумано. б) (480 - х) : 6 + 89 = 165 (480 - х) : 6 = 165 89 (480 - х) : 6 = 76 480 - х = 76 · 6 480 - х = 456 х = 480 - 456 х = 24 Ответ: число 24 задумано.', '<p>\nа) (х + 19) · 5 - 16 = 139<br/>\n(х + 19) · 5 = 139 + 16<br/>\nх + 19 = 155 : 5<br/>\nх + 19 = 31<br/>\nх = 31  19<br/>\nх = 12<br/>\n<b>Ответ:</b> число 12 задумано.<br/><br/>\nб) (480 - х) : 6 + 89 = 165<br/>\n(480 - х) : 6 = 165  89<br/>\n(480 - х) : 6 = 76<br/>\n480 - х = 76 · 6<br/>\n480 - х = 456<br/>\nх = 480 - 456<br/>\nх = 24<br/>\n<b>Ответ:</b> число 24 задумано.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-61/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '91b33903df0e04284377d9e6c6f1f7f90ee7230a5d23dc0d5aa454227c086ee0', '5,6,16,19,89,139,165,480', '["раздели","реши","разность","частное","раз"]'::jsonb, 'составь и реши уравнения:а) задумано число. к нему прибавили 19, сумму умножили на 5 и из полученного произведения вычли 16. получилось 139. какое число задумано? б) задумано число. его вычли из 480, разность разделили на 6 и полученное частное увеличили на 89. в результате получилось 165. какое число задумано');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 61, '4', 3, 'Составь программу действий и вычисли: а) 560 : (720 : 90) - 900 : 50 · 3 + (6 · 8 : 4 + 28) : 5 б) 7 · (45 : 9 · 6 - 23) + 84 : (320 : 80) - 13 · (51 : 17)', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) 560 : (720 : 90) - 900 : 50 · 3 + (6 · 8 : 4 + 28) : 5<br/>\nб) 7 · (45 : 9 · 6 - 23) + 84 : (320 : 80) - 13 · (51 : 17)\n</p>', 'а) 560 : (720 : 90) - 900 : 50 · 3 + (6 · 8 : 4 + 28) : 5 = 24 720 : 90 = 8 560 : 8 = 70 900 : 50 = 18 18 · 3 = 54 6 · 8 = 48 48 : 4 = 12 12 + 28 = 40 40 : 5 = 8 70 - 54 = 16 16 + 8 = 24 б) 7 · (45 : 9 · 6 - 23) + 84 : (320 : 80) - 13 · (51 : 17) = 31 45 : 9 = 5 5 · 6 = 30 30 - 23 = 7 7 · 7 = 49 320 : 80 = 4 84 : 4 = 21 51 : 17 = 3 13 · 3 = 39 49 + 21 = 70 70 - 39 = 31', '<p>\nа) 560 : (720 : 90) - 900 : 50 · 3 + (6 · 8 : 4 + 28) : 5 = 24<br/>\n720 : 90 = 8<br/>\n560 : 8 = 70<br/>\n900 : 50 = 18<br/>\n18 · 3 = 54<br/>\n6 · 8 = 48<br/>\n48 : 4 = 12<br/>\n12 + 28 = 40<br/>\n40 : 5 = 8<br/>\n70 - 54 = 16<br/>\n16 + 8 = 24<br/><br/>\nб) 7 · (45 : 9 · 6 - 23) + 84 : (320 : 80) - 13 · (51 : 17) = 31<br/>\n45 : 9 = 5<br/>\n5 · 6 = 30<br/>\n30 - 23 = 7<br/>\n7 · 7 = 49<br/>\n320 : 80 = 4<br/>\n84 : 4 = 21<br/>\n51 : 17 = 3<br/>\n13 · 3 = 39<br/>\n49 + 21 = 70<br/>\n70 - 39 = 31\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-61/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '22ed3b8c8273eff8f38c7e92b9b226421073a17346000cca4be94c70f9615935', '3,4,5,6,7,8,9,13,17,23', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) 560:(720:90)-900:50*3+(6*8:4+28):5 б) 7*(45:9*6-23)+84:(320:80)-13*(51:17)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 61, '5', 4, 'Найди значения произведений: а) 3015 · 24      в) 81030 · 2600      д) 8170 · 706 б) 527 · 609      г) 12800 · 3560      е) 9030 · 9040', '</p> \n<p class="text">Найди значения произведений: </p> \n\n<p class="description-text"> \nа) 3015 · 24      в) 81030 · 2600      д) 8170 · 706 <br/>\nб) 527 · 609      г) 12800 · 3560      е) 9030 · 9040\n\n</p>', 'а) 3015 · 24 = 72360', '<p>\nа) 3015 · 24 = 72360\n</p>\n\n<div class="img-wrapper-460">\n<img width="120" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica61-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 61, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 61, номер 5, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-61/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica61-nomer5.jpg', 'peterson/3/part3/page61/task5_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b008b309c9ae5f8cbb15f97f115be8862237371a4553d8d5a5d630e18bbaa67a', '24,527,609,706,2600,3015,3560,8170,9030,9040', '["найди"]'::jsonb, 'найди значения произведений:а) 3015*24      в) 81030*2600      д) 8170*706 б) 527*609      г) 12800*3560      е) 9030*9040');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 62, '6', 0, 'Выполни действия: а) 5 ч 18 мин + 4 ч 56 мин      в) 2 т 30 кг - 12 ц 80 кг б) 7 мин 2 с - 1 мин 35 с          г) 1 кг 236 г + 6 кг 764 г', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 5 ч 18 мин + 4 ч 56 мин      в) 2 т 30 кг - 12 ц 80 кг<br/>\nб) 7 мин 2 с - 1 мин 35 с          г) 1 кг 236 г + 6 кг 764 г\n</p>', 'а) 5 ч 18 мин + 4 ч 56 мин = (5 ч + 4 ч) + (18 мин + 56 мин) = 9 ч 74 мин = 10 ч 14 мин б) 7 мин 2 с - 1 мин 35 с = (6 мин - 1 мин) + (60 с + 2 с - 35 с) = 5 мин 27 с в) 2 т 30 кг - 12 ц 80 кг = (19 ц - 12 ц) + (100 кг + 30 кг - 80 кг) = 5 ц 50 кг г) 1 кг 236 г + 6 кг 764 г = (1 кг + 6 кг) + (236 г + 764 г) = 5 кг + 1000 г = 6 кг', '<p>\nа) 5 ч 18 мин + 4 ч 56 мин = (5 ч + 4 ч) + (18 мин + 56 мин) = 9 ч 74 мин = 10 ч 14 мин<br/>\nб) 7 мин 2 с - 1 мин 35 с = (6 мин - 1 мин) + (60 с + 2 с - 35 с) = 5 мин 27 с<br/>\nв) 2 т 30 кг - 12 ц 80 кг = (19 ц - 12 ц) + (100 кг + 30 кг - 80 кг) = 5 ц 50 кг<br/>\nг) 1 кг 236 г + 6 кг 764 г = (1 кг + 6 кг) + (236 г + 764 г) = 5 кг + 1000 г = 6 кг\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-62/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '9b7d87adcc4618275c6855ea023853f4141c6e3012ea3e1e59a838e03069cfb8', '1,2,4,5,6,7,12,18,30,35', NULL, 'выполни действия:а) 5 ч 18 мин+4 ч 56 мин      в) 2 т 30 кг-12 ц 80 кг б) 7 мин 2 с-1 мин 35 с          г) 1 кг 236 г+6 кг 764 г');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 62, '7', 1, 'В автопробеге Париж – Дакар участвовало 420 машин. Экипаж каждой машины состоял из 3 человек. До финиша не дошли 248 машин. Сколько спортсменов прибыли к финишу?', '</p> \n<p class="text">В автопробеге Париж – Дакар участвовало 420 машин. Экипаж каждой машины состоял из 3 человек. До финиша не дошли 248 машин. Сколько спортсменов прибыли к финишу?</p>', '(420 - 248) · 3 = 172 · 3 = 516 (спортсменов) Ответ: 516 спортсменов прибыли к финишу.', '<p>\n(420 - 248) · 3 = 172 · 3 = 516 (спортсменов)<br/>\n<b>Ответ:</b> 516 спортсменов прибыли к финишу.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-62/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '522230c1f821dd46d01d0f03fa64dffe5d4866a9eed14a37e033267b8a85af0f', '3,248,420', NULL, 'в автопробеге париж-дакар участвовало 420 машин. экипаж каждой машины состоял из 3 человек. до финиша не дошли 248 машин. сколько спортсменов прибыли к финишу');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 62, '8', 2, 'Несколько мальчиков ловили рыбу. Всего они поймали 75 рыб. Сколько было мальчиков, если двое поймали по 10 рыб, а остальные – по 11?', '</p> \n<p class="text">Несколько мальчиков ловили рыбу. Всего они поймали 75 рыб. Сколько было мальчиков, если двое поймали по 10 рыб, а остальные – по 11?</p>', '(75 - (2 · 10)) : 11 = (75 - 20) : 11 = 55 : 11 = 5 (мальчиков) – по 11 рыб 5 + 2 = 7 (мальчиков) Ответ: 7 было мальчиков, если двое поймали по 10 рыб, а остальные – по 11.', '<p>\n(75 - (2 · 10)) : 11 = (75 - 20) : 11 = 55 : 11 = 5 (мальчиков) – по 11 рыб<br/>\n5 + 2 = 7 (мальчиков)<br/>\n<b>Ответ:</b> 7 было мальчиков, если двое поймали по 10 рыб, а остальные – по 11. \n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-62/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5aeaaf669cc416df9018bbab2397df1468d5291f15685a2a50527e1a192034fb', '10,11,75', NULL, 'несколько мальчиков ловили рыбу. всего они поймали 75 рыб. сколько было мальчиков, если двое поймали по 10 рыб, а остальные-по 11');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 62, '9', 3, 'Реши задачи и сравни их решения. Что ты замечаешь? 1) Магазин продал за день 16 одинаковых банок вишнёвого варенья и 20 таких же банок малинового. Малинового варенья было продано на 8 кг больше, чем вишнёвого. Сколько килограммов варенья каждого сорта было продано за этот день? 2) Магазин продал за день 32 кг вишнёвого варенья и 40 кг малинового. Всё варенье было разложено в одинаковые банки, причём банок с вишнёвым вареньем было на 4 меньше, чем с малиновым. Сколько банок варенья каждого сорта было продано?', '</p> \n<p class="text">Реши задачи и сравни их решения. Что ты замечаешь?<br/>\n1) Магазин продал за день 16 одинаковых банок вишнёвого варенья и 20 таких же банок малинового. Малинового варенья было продано на 8 кг больше, чем вишнёвого. Сколько килограммов варенья каждого сорта было продано за этот день?<br/>\n2) Магазин продал за день 32 кг вишнёвого варенья и 40 кг малинового. Всё варенье было разложено в одинаковые банки, причём банок с вишнёвым вареньем было на 4 меньше, чем с малиновым. Сколько банок варенья каждого сорта было продано?\n</p>', '1) 8 : (20 - 16) = 8 : 4 = 2 (кг) – в 1 банке 16 · 2 = 32 (кг) – вишнёвое 20 · 2 = 40 (кг) – малиновое Ответ: 32 кг вишнёвого и 40 кг малинового варенья было продано за этот день. 2) (40 - 32) : 4 = 8 : 4 = 2 (кг) – в 1 банке 40 : 2 = 20 (банок) – малинового 32 : 2 = 16 (банок) – вишнёвого Ответ: 20 банок вишнёвого и 16 банок малинового варенья было продано.', '<p>\n1) 8 : (20 - 16) = 8 : 4 = 2 (кг) – в 1 банке<br/>\n16 · 2 = 32 (кг) – вишнёвое<br/>\n20 · 2 = 40 (кг) – малиновое<br/>\n<b>Ответ:</b> 32 кг вишнёвого и 40 кг малинового варенья было продано за этот день.<br/><br/>\n2) (40 - 32) : 4 = 8 : 4 = 2 (кг) – в 1 банке<br/>\n40 : 2 = 20 (банок) – малинового<br/>\n32 : 2 = 16 (банок) – вишнёвого<br/>\n<b>Ответ:</b> 20 банок вишнёвого и 16 банок малинового варенья было продано.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-62/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd17ea2744e8b5efcabcfa9db8f3a25295e9409be2548967e4dae331fd47979b5', '1,2,4,8,16,20,32,40', '["реши","сравни","больше","меньше","раз"]'::jsonb, 'реши задачи и сравни их решения. что ты замечаешь? 1) магазин продал за день 16 одинаковых банок вишнёвого варенья и 20 таких же банок малинового. малинового варенья было продано на 8 кг больше, чем вишнёвого. сколько килограммов варенья каждого сорта было продано за этот день? 2) магазин продал за день 32 кг вишнёвого варенья и 40 кг малинового. всё варенье было разложено в одинаковые банки, причём банок с вишнёвым вареньем было на 4 меньше, чем с малиновым. сколько банок варенья каждого сорта было продано');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 62, '10', 4, 'Сравни выражения*: 118 + n     n + 45       k : 4     k : 6 29 - b     40 - b          14 · d     21 · d x - 35     x - 45          50 : m     15 : m а · b - c     b · а + c m · (n + k)     m · n + k 4 · х + 8 · х     (х · 6) · 2 * Во всех заданиях на сравнение значения букв – натуральные числа и все действия выполнимы.', '</p> \n<p class="text">Сравни выражения*:</p> \n\n<p class="description-text"> \n118 + n <span class="okon">   </span> n + 45       k : 4 <span class="okon">   </span> k : 6<br/> 		\n29 - b <span class="okon">   </span> 40 - b          14 · d <span class="okon">   </span> 21 · d<br/> 		\nx - 35 <span class="okon">   </span>  x - 45          50 : m <span class="okon">   </span> 15 : m<br/><br/>		\n\nа · b - c <span class="okon">   </span> b · а + c<br/>\nm · (n + k) <span class="okon">   </span> m · n + k<br/>\n4 · х + 8 · х <span class="okon">   </span> (х · 6) · 2<br/><br/>\n* Во всех заданиях на сравнение значения букв – натуральные числа и все действия выполнимы.\n</p>', '118 + n > n + 45       k : 4 > k : 6 29 - b < 40 - b          14 · d < 21 · d x - 35 > x - 45          50 : m > 15 : m а · b - c < b · а + c m · (n + k) < m · n + k 4 · х + 8 · х = (х · 6) · 2 20 : 5 = 4 (раза) – остановки на заправку 4 · 2 = 8 (ч) – на заправку 20 + 8 = 28 (ч) Ответ: через 28 часов он прибудет в город В, если дорога идёт вдоль реки. а) 0, 15, 30, 45, 60, 60 + 15 = 75, 75 + 15 = 90, 90 + 15 = 105 б) 1, 4, 9, 16, 25, 25 + 11 = 36, 36 + 13 = 49, 49 + 15 = 64', '<p>\n118 + n &gt; n + 45       k : 4 &gt;  k : 6<br/> 		\n29 - b &lt; 40 - b          14 · d &lt; 21 · d<br/> 		\nx - 35 &gt;  x - 45          50 : m &gt; 15 : m<br/><br/>		\n\nа · b - c &lt; b · а + c<br/>\nm · (n + k) &lt; m · n + k<br/>\n4 · х + 8 · х = (х · 6) · 2\n</p>\n\n\n<p>\n20 : 5 = 4 (раза) – остановки на заправку<br/>\n4 · 2 = 8 (ч) – на заправку<br/>\n20 + 8 = 28 (ч)<br/>\n<b>Ответ:</b> через 28 часов он прибудет в город В, если дорога идёт вдоль реки.\n\n</p>\n\n\n<p>\nа) 0, 15, 30, 45, 60, 60 + 15 = 75, 75 + 15 = 90, 90 + 15 = 105 <br/>      \nб) 1, 4, 9, 16, 25, 25 + 11 = 36, 36 + 13 = 49, 49 + 15 = 64\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-62/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f47df5dc14f790a9399fa61214e82fe375240fb3fd5e600a3bb9d68b32a14c99', '2,4,6,8,14,15,21,29,35,40', '["сравни"]'::jsonb, 'сравни выражения*:118+n     n+45       k:4     k:6 29-b     40-b          14*d     21*d x-35     x-45          50:m     15:m а*b-c     b*а+c m*(n+k)     m*n+k 4*х+8*х     (х*6)*2*во всех заданиях на сравнение значения букв-натуральные числа и все действия выполнимы');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 63, '1', 0, 'Умножение натуральных чисел на четырёхзначное, пятизначное, шестизначное и т. д. число выполняется аналогично тому, как выполняется умножение на трёхзначное число, например: Объясни, как произведены вычисления.', '</p> \n<p class="text">Умножение натуральных чисел на четырёхзначное, пятизначное, шестизначное и т. д. число выполняется аналогично тому, как выполняется умножение на трёхзначное число, например:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica63-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 63, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 63, номер 1, год 2022."/>\n</div>\n</div>\n\n<p class="text">Объясни, как произведены вычисления.</p>', 'Пишем 2318 (начинаем с того, у которого больше разрядов). Под ним записываем 1011 (с новой строки). При этом важно, чтобы одинаковые разряды обоих чисел были расположены строго друг под другом (десятки под десятками, сотни под сотнями и т.д.) Под сомножителями чертим горизонтальную линию, которая будет отделять их от результата. Начинаем выполнять умножение: Крайнюю правую цифру второго множителя (разряд – единицы) поочередно умножаем на каждую цифру первого числа (справа налево) 2318. При этом если ответ оказался двузначным, в текущем разряде оставляем последнюю цифру, а первую переносим в следующий, сложив со значением, полученным в результате умножения. Иногда в результате такого переноса в ответе появляется новый разряд. Затем переходим к следующей цифре второго множителя (десятки) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 2318. Затем переходим к следующей цифре второго множителя (сотни) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 2318. Получившиеся числа складываем и получаем ответ. 2343498 Пишем 704500 (начинаем с того, у которого больше разрядов). Под ним записываем 1001 (с новой строки). При этом важно, чтобы одинаковые разряды обоих чисел были расположены строго друг под другом (десятки под десятками, сотни под сотнями и т.д.), не беря нули во внимание. Их добавляем к готовому ответу. Под сомножителями чертим горизонтальную линию, которая будет отделять их от результата. Начинаем выполнять умножение: Крайнюю правую цифру второго множителя (разряд – единицы) поочередно умножаем на каждую цифру первого числа (справа налево) 7045. При этом если ответ оказался двузначным, в текущем разряде оставляем последнюю цифру, а первую переносим в следующий, сложив со значением, полученным в результате умножения. Иногда в результате такого переноса в ответе появляется новый разряд. Умножение на ноль даст 0000. Их не пишем и следующий ряд сдвигаем влево на одну цифру. И следующий ряд умножение на ноль даст 0000. Их не пишем и следующий ряд сдвигаем влево на одну цифру. Затем переходим к следующей цифре второго множителя (тысячи) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 7045. Получившиеся числа складываем, добавляем 00 и получаем ответ. 705204500 Пишем 44440 (начинаем с того, у которого больше разрядов). Под ним записываем 222200 (с новой строки). При этом важно, чтобы одинаковые разряды обоих чисел были расположены строго друг под другом (десятки под десятками, сотни под сотнями и т.д.), не беря нули во внимание. Их добавляем к готовому ответу. Под сомножителями чертим горизонтальную линию, которая будет отделять их от результата. Начинаем выполнять умножение: Крайнюю правую цифру второго множителя (разряд – единицы) поочередно умножаем на каждую цифру первого числа (справа налево) 8888. При этом если ответ оказался двузначным, в текущем разряде оставляем последнюю цифру, а первую переносим в следующий, сложив со значением, полученным в результате умножения. Иногда в результате такого переноса в ответе появляется новый разряд. Затем переходим к следующей цифре второго множителя (десятки) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 8888. Затем переходим к следующей цифре второго множителя (сотни) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 8888. Затем переходим к следующей цифре второго множителя (тысячи) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 8888. Получившиеся числа складываем, добавляем 00 от первого и от второго множителя 0 и получаем ответ. 9874568000', '<p>\nПишем 2318 (начинаем с того, у которого больше разрядов).<br/>\nПод ним записываем 1011 (с новой строки). При этом важно, чтобы одинаковые разряды обоих чисел были расположены строго друг под другом (десятки под десятками, сотни под сотнями и т.д.)<br/>\nПод сомножителями чертим горизонтальную линию, которая будет отделять их от результата.<br/>\nНачинаем выполнять умножение:<br/>\nКрайнюю правую цифру второго множителя (разряд – единицы) поочередно умножаем на каждую цифру первого числа (справа налево) 2318. При этом если ответ оказался двузначным, в текущем разряде оставляем последнюю цифру, а первую переносим в следующий, сложив со значением, полученным в результате умножения. Иногда в результате такого переноса в ответе появляется новый разряд.<br/>\nЗатем переходим к следующей цифре второго множителя (десятки) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 2318.<br/>\nЗатем переходим к следующей цифре второго множителя (сотни) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 2318.<br/>\nПолучившиеся числа складываем и получаем ответ. 2343498<br/><br/>\n\nПишем 704500 (начинаем с того, у которого больше разрядов).<br/>\nПод ним записываем 1001 (с новой строки). При этом важно, чтобы одинаковые разряды обоих чисел были расположены строго друг под другом (десятки под десятками, сотни под сотнями и т.д.), не беря нули во внимание. Их добавляем к готовому ответу.<br/>\nПод сомножителями чертим горизонтальную линию, которая будет отделять их от результата.<br/>\nНачинаем выполнять умножение:<br/>\nКрайнюю правую цифру второго множителя (разряд – единицы) поочередно умножаем на каждую цифру первого числа (справа налево) 7045. При этом если ответ оказался двузначным, в текущем разряде оставляем последнюю цифру, а первую переносим в следующий, сложив со значением, полученным в результате умножения. Иногда в результате такого переноса в ответе появляется новый разряд.<br/>\nУмножение на ноль даст 0000. Их не пишем и следующий ряд сдвигаем влево на одну цифру. <br/>\nИ следующий ряд умножение на ноль даст 0000. Их не пишем и следующий ряд сдвигаем влево на одну цифру. <br/>\nЗатем переходим к следующей цифре второго множителя (тысячи) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 7045.<br/>\nПолучившиеся числа складываем, добавляем 00 и получаем ответ. 705204500<br/><br/>\n\nПишем 44440 (начинаем с того, у которого больше разрядов).<br/>\nПод ним записываем 222200 (с новой строки). При этом важно, чтобы одинаковые разряды обоих чисел были расположены строго друг под другом (десятки под десятками, сотни под сотнями и т.д.), не беря нули во внимание. Их добавляем к готовому ответу.<br/>\nПод сомножителями чертим горизонтальную линию, которая будет отделять их от результата.<br/>\nНачинаем выполнять умножение:<br/>\nКрайнюю правую цифру второго множителя (разряд – единицы) поочередно умножаем на каждую цифру первого числа (справа налево) 8888. При этом если ответ оказался двузначным, в текущем разряде оставляем последнюю цифру, а первую переносим в следующий, сложив со значением, полученным в результате умножения. Иногда в результате такого переноса в ответе появляется новый разряд.<br/>\nЗатем переходим к следующей цифре второго множителя (десятки) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 8888.<br/>\nЗатем переходим к следующей цифре второго множителя (сотни) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 8888.<br/>\nЗатем переходим к следующей цифре второго множителя (тысячи) и выполняем аналогичные действия, записывая результат со сдвигом на один разряд влево 8888.<br/>\nПолучившиеся числа складываем, добавляем 00 от первого и от второго множителя 0 и получаем ответ. 9874568000\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-63/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica63-nomer1.jpg', 'peterson/3/part3/page63/task1_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f4142663d49a82a0fba1cc83eb66dae1d6abf770d93d0e5ff38d34b7646b4d0b', NULL, NULL, 'умножение натуральных чисел на четырёхзначное, пятизначное, шестизначное и т. д. число выполняется аналогично тому, как выполняется умножение на трёхзначное число, например:объясни, как произведены вычисления');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 63, '2', 1, 'Выполни действия: а) 7032 · 2102 б) 80800 · 7777 в) 12340 · 5609', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) 7032 · 2102  б) 80800 · 7777  в) 12340 · 5609\n</p>', 'а) 7032 · 2102 = 14781264', '<p>\nа) 7032 · 2102 = 14781264\n</p>\n\n<div class="img-wrapper-460">\n<img width="170" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica63-nomer2.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 63, номер 2, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 63, номер 2, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-63/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica63-nomer2.jpg', 'peterson/3/part3/page63/task2_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '82ec7a7f9e596796599adcd7f0bdfee7a8acab5661f1398dd39098296bc6e160', '2102,5609,7032,7777,12340,80800', NULL, 'выполни действия:а) 7032*2102 б) 80800*7777 в) 12340*5609');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 63, '3', 2, 'Практическая работа № 1 Сколько прошло дней, часов, минут, секунд с момента твоего рождения до сегодняшнего дня? (Для простоты вычислений считай день своего рождения и сегодняшний день полностью прожитыми днями.)', '</p> \n<p class="text">Практическая работа № 1<br/>\nСколько прошло дней, часов, минут, секунд с момента твоего рождения до сегодняшнего дня? (Для простоты вычислений считай день своего рождения и сегодняшний день полностью прожитыми днями.)\n</p>', 'Я родился 24 марта 2014 года. Сегодня 20 октября 2023 года. Решение: До 24 марта 2023 года я прожил полных 9 лет, причём 2018 год и 2022 год были високосными: 365 · 9 + 2 = 3287 (дней) С 24 марта 2015 года до 24 сентября этого же года прошло полных 6 месяцев (4 месяца по 31 дню и 2 месяца по 30 дней), и до 20 октября ещё 27 дней: 31 · 4 + 30 · 2 + 27 = 211 (дней) Итак, всего: 3287 + 211 = 3498 дней 24 · 3498 = 83952 часов 60 · 83952 = 5037120 минут 60 · 5037120 = 302227200 секунд', '<p>\nЯ родился 24 марта 2014 года. Сегодня 20 октября 2023 года.<br/>\nРешение: <br/>\nДо 24 марта 2023 года я прожил полных 9 лет, причём 2018 год и 2022 год были високосными:<br/>\n365 · 9 + 2 = 3287 (дней)<br/>\nС 24 марта 2015 года до 24 сентября этого же года прошло полных 6 месяцев (4 месяца по 31 дню и 2 месяца по 30 дней), и до 20 октября ещё 27 дней: <br/>\n31 · 4 + 30 · 2 + 27 = 211 (дней)<br/>\nИтак, всего: <br/>\n3287 + 211 = 3498 дней<br/>\n24 · 3498 = 83952 часов <br/>\n60 · 83952 = 5037120 минут<br/>\n60 · 5037120 = 302227200 секунд\n</p>', 'Образец: Коля Васечкин родился 24 марта 2006 года. Сколько времени он прожил (в днях, минутах, секундах) до 20 октября 2015 года? Решение: До 24 марта 2015 года Коля Васечкин прожил полных 9 лет, причём 2008 год и 2012 год были високосными: 365 · 9 + 2 = 3287 (дней) С 24 марта 2015 года до 24 сентября этого же года прошло полных 6 месяцев (4 месяца по 31 дню и 2 месяца по 30 дней), и до 20 октября ещё 27 дней: 31 · 4 + 30 · 2 + 27 = 211 (дней) Итак, Коля Васечкин прожил всего: 3287 + 211 = 3498 дней 24 · 3498 = 83952 часов 60 · 83952 = 5037120 минут 60 · 5037120 = 302227200 секунд', '<div class="recomended-block">\n<span class="title">Образец:</span>\n<p>\nКоля Васечкин родился 24 марта 2006 года. Сколько времени он прожил (в днях, минутах, секундах) до 20 октября 2015 года?<br/>\nРешение: <br/>\nДо 24 марта 2015 года Коля Васечкин прожил полных 9 лет, причём 2008 год и 2012 год были високосными:<br/>\n365 · 9 + 2 = 3287 (дней)<br/>\nС 24 марта 2015 года до 24 сентября этого же года прошло полных 6 месяцев (4 месяца по 31 дню и 2 месяца по 30 дней), и до 20 октября ещё 27 дней: <br/>\n31 · 4 + 30 · 2 + 27 = 211 (дней)<br/>\nИтак, Коля Васечкин прожил всего: <br/>\n3287 + 211 = 3498 дней<br/>\n24 · 3498 = 83952 часов <br/>\n60 · 83952 = 5037120 минут<br/>\n60 · 5037120 = 302227200 секунд\n</p>\n</div>', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-63/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'cf6cb528a3f9186115b99b5bbb1a62ead1ff28933128e37ca8b7511f28dc3003', '1', NULL, 'практическая работа № 1 сколько прошло дней, часов, минут, секунд с момента твоего рождения до сегодняшнего дня? (для простоты вычислений считай день своего рождения и сегодняшний день полностью прожитыми днями.)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 64, '4', 0, 'Практическая работа № 2 Узнай, сколько дней, часов, минут, секунд прожил кто-либо из твоих родных или друзей (по твоему выбору) с момента рождения до сегодняшнего дня.', '</p> \n<p class="text"><b>Практическая работа № 2</b><br/>\nУзнай, сколько дней, часов, минут, секунд прожил кто-либо из твоих родных или друзей (по твоему выбору) с момента рождения до сегодняшнего дня.\n</p>', 'Моя сестра - близнец родилась 24 марта 2014 года. Сегодня 20 октября 2023 года. Решение: До 24 марта 2023 года моя сестра прожила полных 9 лет, причём 2018 год и 2022 год были високосными: 365 · 9 + 2 = 3287 (дней) С 24 марта 2015 года до 24 сентября этого же года прошло полных 6 месяцев (4 месяца по 31 дню и 2 месяца по 30 дней), и до 20 октября ещё 27 дней: 31 · 4 + 30 · 2 + 27 = 211 (дней) Итак, всего: 3287 + 211 = 3498 дней 24 · 3498 = 83952 часов 60 · 83952 = 5037120 минут 60 · 5037120 = 302227200 секунд', '<p>\nМоя сестра - близнец родилась 24 марта 2014 года. Сегодня 20 октября 2023 года.<br/>\nРешение: <br/>\nДо 24 марта 2023 года моя сестра прожила полных 9 лет, причём 2018 год и 2022 год были високосными:<br/>\n365 · 9 + 2 = 3287 (дней)<br/>\nС 24 марта 2015 года до 24 сентября этого же года прошло полных 6 месяцев (4 месяца по 31 дню и 2 месяца по 30 дней), и до 20 октября ещё 27 дней: <br/>\n31 · 4 + 30 · 2 + 27 = 211 (дней)<br/>\nИтак, всего: <br/>\n3287 + 211 = 3498 дней<br/>\n24 · 3498 = 83952 часов <br/>\n60 · 83952 = 5037120 минут<br/>\n60 · 5037120 = 302227200 секунд\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-64/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '9ee2c4b825ca027f8629e79240c2a36f0e2231ec6a14bbf17251ee27856bf69f', '2', NULL, 'практическая работа № 2 узнай, сколько дней, часов, минут, секунд прожил кто-либо из твоих родных или друзей (по твоему выбору) с момента рождения до сегодняшнего дня');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 64, '5', 1, 'В библиотеке три хранилища. В первом хранилище 15789 книг, во втором на 2634 книги меньше, чем в первом, а в третьем в 6 раз меньше, чем в первых двух хранилищах вместе. Сколько всего книг в библиотеке?', '</p> \n<p class="text">В библиотеке три хранилища. В первом хранилище 15789 книг, во втором на 2634 книги меньше, чем в первом, а в третьем в 6 раз меньше, чем в первых двух хранилищах вместе. Сколько всего книг в библиотеке?</p>', '15789 + (15789 - 2634) + (15789 + (15789 – 2634)) : 6 = 15789 + 13155 + (15789 + 13155) : 6 = 28944 + 28944 : 6 = 28944 + 4824 = 33768 (книг)', '<p>\n15789 + (15789 - 2634) + (15789 + (15789 – 2634)) : 6 = 15789 + 13155 + (15789 + 13155) : 6 = 28944 + 28944 : 6 = 28944 + 4824 = 33768 (книг)\n</p>\n\n<div class="img-wrapper-460">\n<img width="310" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica64-nomer5.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 64, номер 5, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 64, номер 5, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-64/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica64-nomer5.jpg', 'peterson/3/part3/page64/task5_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '71045ef9d88a29c422c27c34fee1870f4d61c2eacc212a4765c2e3a021909e24', '6,2634,15789', '["меньше","раз"]'::jsonb, 'в библиотеке три хранилища. в первом хранилище 15789 книг, во втором на 2634 книги меньше, чем в первом, а в третьем в 6 раз меньше, чем в первых двух хранилищах вместе. сколько всего книг в библиотеке');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 64, '6', 2, 'Купили три отреза одинаковой ткани. В первом отрезе 7 м ткани, во втором – в 2 раза больше, чем в первом, а в третьем – на 5 м меньше, чем во втором. За все три отреза заплатили 43200 р. Сколько стоит каждый отрез?', '</p> \n<p class="text">Купили три отреза одинаковой ткани. В первом отрезе 7 м ткани, во втором – в 2 раза больше, чем в первом, а в третьем – на 5 м меньше, чем во втором. За все три отреза заплатили 43200 р. Сколько стоит каждый отрез?</p>', '43200 : (7 + (7 · 2) + (7 · 2) - 5) = 43200 : (7 + 14 + 9) = 43200 : 30 = 1440 (р.) - 1 метр 7 · 1440 = 10080 (р.) – первый 7 · 2 · 1440 = 14 · 1440 = 20160 (р.) – второй', '<p>\n43200 : (7 + (7 · 2) + (7 · 2) - 5) = 43200 : (7 + 14 + 9) = 43200 : 30 = 1440 (р.) - 1 метр<br/>\n\n7 · 1440 = 10080 (р.) – первый<br/>\n\n7 · 2 · 1440 = 14 · 1440 = 20160 (р.) – второй\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="120" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica64-nomer6.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 64, номер 6, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 64, номер 6, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-64/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica64-nomer6.jpg', 'peterson/3/part3/page64/task6_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2cd1596cf2b49371bbd24ec182e1477d5fc4edef69cfdb35c3a92a123768b833', '2,5,7,43200', '["больше","меньше","раз","раза"]'::jsonb, 'купили три отреза одинаковой ткани. в первом отрезе 7 м ткани, во втором-в 2 раза больше, чем в первом, а в третьем-на 5 м меньше, чем во втором. за все три отреза заплатили 43200 р. сколько стоит каждый отрез');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 64, '7', 3, 'Междугородний автобус должен проехать расстояние между двумя городами, равное 350 км, за 7 часов. Но первые два часа из-за сильного дождя он ехал со скоростью на 5 км/ч меньше, чем предполагалось. С какой скоростью автобус должен проехать оставшийся путь, чтобы прийти в пункт назначения без опоздания?', '</p> \n<p class="text">Междугородний автобус должен проехать расстояние между двумя городами, равное 350 км, за 7 часов. Но первые два часа из-за сильного дождя он ехал со скоростью на 5 км/ч меньше, чем предполагалось. С какой скоростью автобус должен проехать оставшийся путь, чтобы прийти в пункт назначения без опоздания?</p>', '350 : 7 = 50 (км/ч) 50 - 5 = 45 (км/ч) – скорость в первые 2 часа 45 · 2 = 90 (км) – расстояние за первые 2 часа 350 - 90 = 260 (км) – оставшееся расстояние 7 - 2 = 5 (ч) – оставшееся время 260 : 5 = 52 (км/ч) Ответ: 52 км/ч автобус должен проехать оставшийся путь, чтобы прийти в пункт назначения без опоздания.', '<p>\n350 : 7 = 50 (км/ч)<br/>\n50 - 5 = 45 (км/ч) – скорость в первые 2 часа<br/>\n45 · 2 = 90 (км) – расстояние за первые 2 часа<br/>\n350 - 90 = 260 (км) – оставшееся расстояние<br/>\n7 - 2 = 5 (ч) – оставшееся время <br/>\n260 : 5 = 52 (км/ч)<br/>\n<b>Ответ:</b> 52 км/ч автобус должен проехать оставшийся путь, чтобы прийти в пункт назначения без опоздания.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-64/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1b800aa175bf80d92cb2cdf1f8d94f6036294650fdc2b1083e985ff5b4ab6b46', '5,7,350', '["меньше","равно"]'::jsonb, 'междугородний автобус должен проехать расстояние между двумя городами, равное 350 км, за 7 часов. но первые два часа из-за сильного дождя он ехал со скоростью на 5 км/ч меньше, чем предполагалось. с какой скоростью автобус должен проехать оставшийся путь, чтобы прийти в пункт назначения без опоздания');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 64, '8', 4, 'Олег пробежал 1 км за 5 мин. На сколько быстрее он пробежит это расстояние, если увеличит скорость на 50 м/мин?', '</p> \n<p class="text">Олег пробежал 1 км за 5 мин. На сколько быстрее он пробежит это расстояние, если увеличит скорость на 50 м/мин?</p>', '1000 м : 5 мин + 50 м/мин = 200 м/мин + 50 м/мин = 250 м/мин 1000 : 250 = 4 (мин) 5 - 4 = 1 (мин) Ответ: на 1 минуту быстрее он пробежит это расстояние, если увеличит скорость на 50 м/мин.', '<p>\n1000 м : 5 мин + 50 м/мин = 200 м/мин + 50 м/мин = 250 м/мин<br/>\n1000 : 250 = 4 (мин)<br/>\n5 - 4 = 1 (мин)<br/>\n<b>Ответ:</b> на 1 минуту быстрее он пробежит это расстояние, если увеличит скорость на 50 м/мин.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-64/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '60c39279c0e46de8db4b407f889051c7f3b11f220abdf446d19c3b67fc417f01', '1,5,50', NULL, 'олег пробежал 1 км за 5 мин. на сколько быстрее он пробежит это расстояние, если увеличит скорость на 50 м/мин');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 64, '9', 5, 'Найди значение выражения 450 - 9 · x, если x = 0, 1, 6, 8, 9, 40. Какое наибольшее значение может принимать x?', '</p> \n<p class="text">Найди значение выражения 450 - 9 · x, если x = 0, 1, 6, 8, 9, 40. Какое наибольшее значение может принимать x?</p>', '450 - 9 · 0 = 450 450 - 9 · 1 = 441 450 - 9 · 6 = 450 - 54 = 406 450 - 9 · 8 = 450 - 72 = 378 450 - 9 · 9 = 450 - 81 = 369 450 - 9 · 40 = 450 - 360 = 90 х = 50 450 - 9 · 50 = 450 - 450 = 0 Ответ: 50 наибольшее значение может принимать x.', '<p>\n450 - 9 · 0 = 450<br/>\n450 - 9 · 1 = 441<br/>\n450 - 9 · 6 = 450 - 54 = 406<br/>\n450 - 9 · 8 = 450 - 72 = 378<br/>\n450 - 9 · 9 = 450 - 81 = 369<br/>\n450 - 9 · 40 = 450 - 360 = 90<br/>\nх = 50<br/>\n450 - 9 · 50 = 450 - 450 = 0<br/>\n<b>Ответ:</b> 50 наибольшее значение может принимать x.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-64/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '626bbdde3d5b265371ccc15265aa1cd055f1c689c743fb4bb69bf66b593d2f36', '0,1,6,8,9,40,450', '["найди","больше"]'::jsonb, 'найди значение выражения 450-9*x, если x=0, 1, 6, 8, 9, 40. какое наибольшее значение может принимать x');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 64, '10', 6, 'Запиши множество делителей и множество кратных числа 32.', '</p> \n<p class="text">Запиши множество делителей и множество кратных числа 32.</p>', 'Делители 32: 1, 2, 4, 8, 16, 32. Кратные 32: 32, 64, 96, 128, 160, 192 и так далее с каждым разом прибавляется число 32 от последнего числа.', '<p>\nДелители 32: 1, 2, 4, 8, 16, 32. <br/>\nКратные 32: 32, 64, 96, 128, 160, 192 и так далее с каждым разом прибавляется число 32 от последнего числа.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-64/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '75b2cf8557a61adeb642b642db24c9ea3728826e6d068220ac2a1a5a8199965d', '32', NULL, 'запиши множество делителей и множество кратных числа 32');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 64, '11', 7, 'Во сколько раз число A больше, чем число B: A - (35302 - 28394) · 1500 : 400 + 479145 B - 57912 - 180 · (119486 + 3964) : 3000', '</p> \n<p class="text">Во сколько раз число A больше, чем число B:<br/>\nA - (35302 - 28394) · 1500 : 400 + 479145<br/>\nB - 57912 - 180 · (119486 + 3964) : 3000\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica64-nomer11.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 64, номер 11, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 64, номер 11, год 2022."/>\n</div>\n</div>', 'A - (35302 - 28394) · 1500 : 400 + 479145 = 6908 · 1500 : 400 + 479145 =', '<p>\nA - (35302 - 28394) · 1500 : 400 + 479145 = 6908 · 1500 : 400 + 479145 = \n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica64-nomer11-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 64, номер 11-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 64, номер 11-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-64/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica64-nomer11.jpg', 'peterson/3/part3/page64/task11_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica64-nomer11-1.jpg', 'peterson/3/part3/page64/task11_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'db4239a9d5f51280ff0abe541e5d331d24d0e94f2c46d3f6e12ce1a9671e3e52', '180,400,1500,3000,3964,28394,35302,57912,119486,479145', '["больше","раз"]'::jsonb, 'во сколько раз число a больше, чем число b:a-(35302-28394)*1500:400+479145 b-57912-180*(119486+3964):3000');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 64, '12', 8, 'Определи, какого данного не хватает для ответа на вопрос задачи. Подбери возможное значение и реши задачу. а) Мама купила 4 кг гречки. Сколько денег она заплатила? б) Турист прошел за день 24 км. Сколько времени он был в пути? в) Рабочий делает 3 детали в час. Сколько всего деталей он сделал? г) Самолёт вылетел из Москвы в 9 ч 25 мин утра. В котором часу он приземлился в Новосибирске?', '</p> \n<p class="text">Определи, какого данного не хватает для ответа на вопрос задачи. Подбери возможное значение и реши задачу.<br/> \nа) Мама купила 4 кг гречки. Сколько денег она заплатила?<br/>\nб) Турист прошел за день 24 км. Сколько времени он был в пути?<br/>\nв) Рабочий делает 3 детали в час. Сколько всего деталей он сделал? <br/>\nг) Самолёт вылетел из Москвы в 9 ч 25 мин утра. В котором часу он приземлился в Новосибирске?\n</p>', 'а) не хватает стоимости 1 кг = 60 р. Тогда 60 · 4 = 240 (р.) Ответ: 240 рублей она заплатила. б) не хватает скорости туриста 3 км/ч Тогда 24 : 3 = 8 (ч) Ответ: 8 часов он был в пути. в) не хватает количество часов 4 часа Тогда 4 · 3 = 12 (детали) Ответ: 12 всего деталей он сделал. г) не хватает времени полёта 4 часа и разница часовых поясов 4 часа Тогда 9 ч 25 мин + 4 ч + 4 ч = 17 ч 25 мин Ответ: в 17 ч 25 мин он приземлился в Новосибирске.', '<p>\nа) не хватает стоимости 1 кг = 60 р. <br/>\nТогда 60 · 4 = 240 (р.)<br/>\n<b>Ответ:</b> 240 рублей она заплатила.<br/><br/>\n\nб) не хватает скорости туриста 3 км/ч<br/>\nТогда 24 : 3 = 8 (ч)<br/>\n<b>Ответ:</b> 8 часов он был в пути.<br/><br/>\n\nв) не хватает количество часов 4 часа<br/>\nТогда 4 · 3 = 12 (детали)<br/>\n<b>Ответ:</b> 12 всего деталей он сделал.<br/><br/>\n\nг) не хватает времени полёта 4 часа и разница часовых поясов 4 часа<br/>\nТогда 9 ч 25 мин + 4 ч + 4 ч = 17 ч 25 мин<br/>\n<b>Ответ:</b> в 17 ч 25 мин он приземлился в Новосибирске.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-64/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8c833e14377c682feb3071aaa65375405be814196ea70efbb8cd78babfbcb8cb', '3,4,9,24,25', '["реши"]'::jsonb, 'определи, какого данного не хватает для ответа на вопрос задачи. подбери возможное значение и реши задачу. а) мама купила 4 кг гречки. сколько денег она заплатила? б) турист прошел за день 24 км. сколько времени он был в пути? в) рабочий делает 3 детали в час. сколько всего деталей он сделал? г) самолёт вылетел из москвы в 9 ч 25 мин утра. в котором часу он приземлился в новосибирске');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 65, '13', 0, '1) Верно ли, что 1 л равен 1 м 3 ? 2) Верно ли, что масса арбуза может быть равна 5 кг? 3) Верно ли, что скорость пешехода равна примерно 30 км/ч? 4) Верно ли, что 7 км2 24 м 2 больше по площади, чем 80000 дм 2 ?', '</p> \n<p class="text">1) Верно ли, что 1 л равен 1 м<sup>3</sup>? <br/>\n2) Верно ли, что масса арбуза может быть равна 5 кг?<br/>\n3) Верно ли, что скорость пешехода равна примерно 30 км/ч?<br/>\n4) Верно ли, что 7 км2 24 м<sup>2</sup>  больше по площади, чем 80000 дм<sup>2</sup>?\n</p>', '1) верно 2) верно 3) не верно 4) верно 700000000 дм 2 24 м 2 ˃ 80000 дм 2 1 · 9 + 2 = 11 12 · 9 + 3 = 111 123 · 9 + 4 = 1111 1234 · 9 + 5 = 11111 12345 · 9 + 6 = 111111 123456 · 9 + 7 = 1111111 1234567 · 9 + 8 = 11111111 12345678 · 9 + 9 = 111111111 123456789 · 9 + 10 = 1111111101 + 10 = 1111111111 Ответ: сохранилась данная закономерность для следующей строки.', '<p>\n1) верно<br/>\n2) верно<br/>\n3) не верно<br/>\n4) верно 700000000 дм<sup>2</sup> 24 м<sup>2</sup> ˃ 80000 дм<sup>2</sup> \n</p>\n\n\n<p>\n1 · 9 + 2 = 11<br/>\n12 · 9 + 3 = 111<br/>\n123 · 9 + 4 = 1111<br/>\n1234 · 9 + 5 = 11111<br/>\n12345 · 9 + 6 = 111111<br/>\n123456 · 9 + 7 = 1111111<br/>\n1234567 · 9 + 8 = 11111111<br/>\n12345678 · 9 + 9 = 111111111<br/>\n123456789 · 9 + 10 = 1111111101 + 10 = 1111111111<br/>\n<b>Ответ:</b> сохранилась данная закономерность для следующей строки.\n</p>\n\n\n<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica65-nomer15-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 65, номер 15-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 65, номер 15-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-65/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica65-nomer15-1.jpg', 'peterson/3/part3/page65/task13_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c54544b9aa5def081aaf5a39e2b36d6b77fa187750de98908b1e725e29f4831b', '1,2,3,4,5,7,24,30,80000', '["больше"]'::jsonb, '1) верно ли, что 1 л равен 1 м 3 ? 2) верно ли, что масса арбуза может быть равна 5 кг? 3) верно ли, что скорость пешехода равна примерно 30 км/ч? 4) верно ли, что 7 км2 24 м 2 больше по площади, чем 80000 дм 2');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 66, '1', 0, 'Продолжи ряд на два числа, сохраняя закономерность: а) 0, 19, 38, 57 ...                     г) 1, 9, 25, 49, 81, 121 ... б) 318, 422, 526 ...                  д) 0, 2, 6, 12, 20, 30 ... в) 72574, 72561, 72548 ...    е) 2, 3, 5, 8, 12, 17 ...', '</p> \n<p class="text">Продолжи ряд на два числа, сохраняя закономерность:</p> \n\n<p class="description-text"> \nа) 0, 19, 38, 57 ...                     г) 1, 9, 25, 49, 81, 121 ...<br/>\nб) 318, 422, 526 ...                  д) 0, 2, 6, 12, 20, 30 ...<br/>\nв) 72574, 72561, 72548 ...    е) 2, 3, 5, 8, 12, 17 ...\n</p>', 'а) 0, 19, 38, 57, 76, 95, 114', '<p>\nа) 0, 19, 38, 57, 76, 95, 114\n</p>\n\n<div class="img-wrapper-460">\n<img width="240" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica66-nomer1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 66, номер 1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 66, номер 1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-66/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica66-nomer1.jpg', 'peterson/3/part3/page66/task1_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3b5ec03add58eab174b359c2022286498dee9ca26d9a5a9966a8170436158065', '0,1,2,3,5,6,8,9,12,17', NULL, 'продолжи ряд на два числа, сохраняя закономерность:а) 0, 19, 38, 57 ...                     г) 1, 9, 25, 49, 81, 121 ... б) 318, 422, 526 ...                  д) 0, 2, 6, 12, 20, 30 ... в) 72574, 72561, 72548 ...    е) 2, 3, 5, 8, 12, 17');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 66, '2', 1, 'Что общего в примерах каждого столбика? Объясни приёмы вычислений. 36 + 9        50 - 23        24 · 3 27 + 48      71 - 15        4 · 19 75 : 5        68 : 17 84 : 6        92 : 46', '</p> \n<p class="text">Что общего в примерах каждого столбика? Объясни приёмы вычислений.</p> \n\n<p class="description-text"> \n36 + 9        50 - 23        24 · 3 <br/> 	\n27 + 48      71 - 15        4 · 19  <br/><br/>	\n\n75 : 5        68 : 17<br/>\n84 : 6        92 : 46\n\n</p>', '36 + 9 = 45 27 + 48 = 75 Сложить единицы, если двузначное число, то десятки сложить с десятками 50 - 23 = 27 71 - 15 = 56 Вычитаем единиц и десятки из десяток 24 · 3 = 72 4 · 19 = 76 Умножаем единицы на единицы и на десятки двузначного числа 75 : 5 = 15 (50 + 25) : 5 = 10 + 5 = 15 84 : 6 = 14 (60 + 24) : 6 = 10 + 4 = 14 Представим двузначное число в виде суммы разрядных или удобных слагаемых, Разделим каждое слагаемое на это число, Сложим полученные результаты. 68 : 17 = 4 28 : 7 = 4, поэтому по 4 92 : 46 = 2 12 : 6 = 2, поэтому по 2 Подбор по 2, по 3, по 4', '<p>\n36 + 9 = 45<br/>  							\n27 + 48 = 75<br/>							\nСложить единицы, если двузначное число, то десятки сложить с десятками<br/>\n50 - 23 = 27<br/>\n71 - 15 = 56<br/>\nВычитаем единиц и десятки из десяток<br/>\n24 · 3 = 72<br/>\n4 · 19 = 76<br/>\nУмножаем единицы на единицы и на десятки двузначного числа<br/>\n75 : 5 = 15<br/>\n(50 + 25) : 5 = 10 + 5 = 15<br/>\n84 : 6 = 14<br/>\n(60 + 24) : 6 = 10 + 4 = 14<br/>\nПредставим двузначное число в виде суммы разрядных или удобных слагаемых, <br/>\nРазделим каждое слагаемое на это число,<br/>\nСложим полученные результаты.<br/>\n68 : 17 = 4<br/>\n28 : 7 = 4, поэтому по 4<br/>\n92 : 46 = 2<br/>\n12 : 6 = 2, поэтому по 2<br/>\nПодбор по 2, по 3, по 4\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-66/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '480aefb514a092f9643ade527a44d0f6fdbcc92a46702c2edeb55ce0dd7e6580', '3,4,5,6,9,15,17,19,23,24', '["столбик"]'::jsonb, 'что общего в примерах каждого столбика? объясни приёмы вычислений. 36+9        50-23        24*3 27+48      71-15        4*19 75:5        68:17 84:6        92:46');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 66, '3', 2, 'Запиши на математическом языке: а) переместительное свойство сложения и умножения; б) сочетательное свойство сложения и умножения; в) распределительное свойство умножения; г) правило деления суммы на число; д) правило вычитания числа из суммы; е) правило вычитания суммы из числа. Объясни их смысл.', '</p> \n<p class="text">Запиши на математическом языке: а) переместительное свойство сложения и умножения; б) сочетательное свойство сложения и умножения; в) распределительное свойство умножения; г) правило деления суммы на число; д) правило вычитания числа из суммы; е) правило вычитания суммы из числа. Объясни их смысл.</p>', 'а) a · b = b · a, a + b = b + a От перестановки слагаемых местами их сумма не изменится. От перестановки мест множителей произведение не меняется. б) a + (b + c) = (a + b) + c, a · (b · c) = (a · b) · c Чтобы умножить число на произведение двух чисел, можно сначала умножить его на первый множитель, а потом полученное произведение умножить на второй множитель. Чтобы к сумме двух чисел прибавить третье число, можно к первому прибавить сумму второго и третьего чисел. в) a · (b · c) = a · b + a · c Чтобы умножить сумму на число, нужно умножить на это число каждое слагаемое и сложить полученные результаты. г) (a + b) : c = a : c + b : c Чтобы разделить сумму на число, можно разделить на это число каждое слагаемое и полученные результаты сложить. д) a - (b + c) = a - b - c Чтобы из суммы вычесть число, можно вычесть его из одного слагаемого, а к полученной разности прибавить другое слагаемое. е) a - (b + c) = a - b - c Чтобы вычесть сумму из числа, можно сначала вычесть из этого числа первое слагаемое, а потом из полученной разности – второе слагаемое.', '<p>\nа) a · b = b · a, a + b = b + a<br/>\nОт перестановки слагаемых местами их сумма не изменится. <br/>\nОт перестановки мест множителей произведение не меняется.<br/><br/>\nб) a + (b + c) = (a + b) + c, a · (b · c) = (a · b) · c<br/>\nЧтобы умножить число на произведение двух чисел, можно сначала умножить его на первый множитель, а потом полученное произведение умножить на второй множитель.<br/>\nЧтобы к сумме двух чисел прибавить третье число, можно к первому прибавить сумму второго и третьего чисел.<br/><br/>\nв) a · (b · c) = a · b + a · c<br/>\nЧтобы умножить сумму на число, нужно умножить на это число каждое слагаемое и сложить полученные результаты.<br/><br/>\nг) (a + b) : c = a : c + b : c<br/>\nЧтобы разделить сумму на число, можно разделить на это число каждое слагаемое и полученные результаты сложить.<br/><br/>\nд) a - (b + c) = a - b - c<br/>\nЧтобы из суммы вычесть число, можно вычесть его из одного слагаемого, а к полученной разности прибавить другое слагаемое. <br/><br/>\nе) a - (b + c) = a - b - c<br/>\nЧтобы вычесть сумму из числа, можно сначала вычесть из этого числа первое слагаемое, а потом из полученной разности – второе слагаемое.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-66/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ad56b2850c517c6193ac841261db91f8022201e334e819d40d074475c870c1a9', NULL, '["делитель"]'::jsonb, 'запиши на математическом языке:а) переместительное свойство сложения и умножения; б) сочетательное свойство сложения и умножения; в) распределительное свойство умножения; г) правило деления суммы на число; д) правило вычитания числа из суммы; е) правило вычитания суммы из числа. объясни их смысл');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 66, '4', 3, 'Пользуясь свойствами арифметических действий, упрости выражения: 99 + 1 + а        34 - (27 + с) 16 + b + 9        (d + 46) - 45 8 · m · 3          5 · х - 2 · х n · 25 · 4        9 · у + у', '</p> \n<p class="text">Пользуясь свойствами арифметических действий, упрости выражения:</p> \n\n<p class="description-text"> \n99 + 1 + а        34 - (27 + с) <br/>			\n16 + b + 9        (d + 46) - 45  <br/><br/>		\n\n8 · m · 3          5 · х - 2 · х<br/>\nn · 25 · 4        9 · у + у\n</p>', '99 + 1 + а = 100 + a 16 + b + 9 = 25 + b 34 - (27 + с) = 34 - 27 - c = 7 - c (d + 46) - 45 = d + 46 - 45 = d + 1 8 · m · 3 = 24 · m n · 25 · 4 = n · 100 5 · х - 2 · х = (5 - 2) · x = 3 · x 9 · у + у = (9 + 1) · y = 10 · y', '<p>\n99 + 1 + а = 100 + a<br/>\n16 + b + 9 = 25 + b<br/>\n34 - (27 + с) = 34 - 27 - c = 7 - c<br/>\n(d + 46) - 45 = d + 46 - 45 = d + 1<br/>\n8 · m · 3 = 24 · m<br/>\nn · 25 · 4 = n · 100<br/>\n5 · х - 2 · х = (5 - 2) · x = 3 · x<br/>\n9 · у + у = (9 + 1) · y = 10 · y\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-66/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '18f0f26e467fd6df3a7e4e93ed412dfd06145c2aa2a3452cb23fe35f96f7ae6b', '1,2,3,4,5,8,9,16,25,27', NULL, 'пользуясь свойствами арифметических действий, упрости выражения:99+1+а        34-(27+с) 16+b+9        (d+46)-45 8*m*3          5*х-2*х n*25*4        9*у+у');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 66, '5', 4, 'Вычисли наиболее удобным способом: а) 32 + 34 + 36 + 38      г) (786 + 195) - 586 б) 5 · 19 · 5 · 3 · 2 · 2      д) 903 - 672 - 28 в) 47 · 15 + 53 · 15        е) 245 · 64 - 245 · 54', '</p> \n<p class="text">Вычисли наиболее удобным способом:</p> \n\n<p class="description-text"> \nа) 32 + 34 + 36 + 38      г) (786 + 195) - 586<br/>\nб) 5 · 19 · 5 · 3 · 2 · 2      д) 903 - 672 - 28<br/>\nв) 47 · 15 + 53 · 15        е) 245 · 64 - 245 · 54\n</p>', 'а) 32 + 34 + 36 + 38 = 70 + 70 = 140 б) 5 · 19 · 5 · 3 · 2 · 2 = 10 · 10 · 57 = 5700 в) 47 · 15 + 53 · 15 = (47 + 53) · 15 = 100 · 15 = 1500 г) (786 + 195) - 586 = 195 + (786 - 586) = 195 + 200 = 395 д) 903 - 672 - 28 = 903 - (672 + 28) = 903 - 700 = 203 е) 245 · 64 - 245 · 54 = 245 · (64 - 54) = 245 · 10 = 2450', '<p>\nа) 32 + 34 + 36 + 38 = 70 + 70 = 140  <br/>		\nб) 5 · 19 · 5 · 3 · 2 · 2 = 10 · 10 · 57 = 5700  	<br/>	\nв) 47 · 15 + 53 · 15 = (47 + 53) · 15 = 100 · 15 = 1500  	<br/>	\nг) (786 + 195) - 586 = 195 + (786 - 586) = 195 + 200 = 395<br/>\nд) 903 - 672 - 28 = 903 - (672 + 28) = 903 - 700 = 203<br/>\nе) 245 · 64 - 245 · 54 = 245 · (64 - 54) = 245 · 10 = 2450\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-66/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0f680506dfef7b084c74367fe772fd379d6e313c94b219b3426c461535a1eae2', '2,3,5,15,19,28,32,34,36,38', '["вычисли"]'::jsonb, 'вычисли наиболее удобным способом:а) 32+34+36+38      г) (786+195)-586 б) 5*19*5*3*2*2      д) 903-672-28 в) 47*15+53*15        е) 245*64-245*54');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 66, '6', 5, 'БЛИЦтурнир а) У Ани а марок, а у Тани на с марок меньше. Сколько марок у Ани и Тани вместе? б) Купили n слив. За обедом съели х слив, а за ужином – k слив. Сколько слив осталось? в) Было d красных шариков и k синих. Их разделили поровну на 3 человек. Сколько шариков досталось каждому? г) Артём поймал а рыбок, а Юра – в 4 раза больше. На сколько рыбок меньше поймал Артём, чем Юра? д) После того как в саду посадили 4 ряда вишен по t вишен в ряду, осталось посадить ещё m вишен. Сколько всего вишен должны посадить в саду?', '</p> \n<p class="text">БЛИЦтурнир <br/>\nа) У Ани а марок, а у Тани на с марок меньше. Сколько марок у Ани и Тани вместе? <br/>\nб) Купили n слив. За обедом съели х слив, а за ужином – k слив. Сколько слив осталось? <br/>\nв) Было d красных шариков и k синих. Их разделили поровну на 3 человек.  Сколько шариков досталось каждому? <br/>\nг) Артём поймал а рыбок, а Юра – в 4 раза больше. На сколько рыбок меньше поймал Артём, чем Юра? <br/>\nд) После того как в саду посадили 4 ряда вишен по t вишен в ряду, осталось посадить ещё m вишен. Сколько всего вишен должны посадить в саду?\n</p>', 'а) a + (a - с) (марок) Ответ: a + (a - с) марок у Ани и Тани вместе. б) n - х - k (слив) Ответ: n - х - k слив осталось. в) (d + k) : 3 (шариков) Ответ: (d + k) : 3 шариков досталось каждому? г) 4 · а - а = а · (4 - 1) = 3 · а (рыбок) Ответ: на 3 · а рыбок меньше поймал Артём, чем Юра. д) 4 · t + m (вишен) Ответ: 4 · t + m всего вишен должны посадить в саду.', '<p>\nа) a + (a - с)  (марок)<br/>\n<b>Ответ:</b> a + (a - с) марок у Ани и Тани вместе.<br/><br/>\nб) n - х - k (слив)<br/>\n<b>Ответ:</b> n - х - k слив осталось.<br/><br/>\nв) (d + k) : 3 (шариков) <br/>\n<b>Ответ:</b> (d + k) : 3 шариков досталось каждому?<br/><br/>\nг) 4 · а - а = а · (4 - 1) = 3 · а (рыбок)<br/>\n<b>Ответ:</b> на 3 · а рыбок меньше поймал Артём, чем Юра.<br/><br/>\nд) 4 · t + m (вишен) <br/>\n<b>Ответ:</b> 4 · t + m всего вишен должны посадить в саду.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-66/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3537c6781756e805f62752211ff56df9b661ff9890a7bfa065b7bf81222539c5', '3,4', '["раздели","больше","меньше","раз","раза"]'::jsonb, 'блицтурнир а) у ани а марок, а у тани на с марок меньше. сколько марок у ани и тани вместе? б) купили n слив. за обедом съели х слив, а за ужином-k слив. сколько слив осталось? в) было d красных шариков и k синих. их разделили поровну на 3 человек. сколько шариков досталось каждому? г) артём поймал а рыбок, а юра-в 4 раза больше. на сколько рыбок меньше поймал артём, чем юра? д) после того как в саду посадили 4 ряда вишен по t вишен в ряду, осталось посадить ещё m вишен. сколько всего вишен должны посадить в саду');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 67, '7', 0, 'Найди значения выражений: а) 6 · x, если x = 17 б) 90 - у : 8, если у = 64 в) (75 + а) - (94 + b), если а = 25, b = 3', '</p> \n<p class="text">Найди значения выражений:</p> \n\n<p class="description-text"> \nа) 6 · x, если x = 17<br/>\nб) 90 - у : 8, если у = 64<br/>\nв) (75 + а) - (94 + b), если а = 25, b = 3\n</p>', 'а) 6 · 17 = 102 б) 90 - 64 : 8 = 90 - 8 = 82 в) (75 + 25) - (94 + 3) = 100 - 97 = 3', '<p>\nа) 6 · 17 = 102<br/>\nб) 90 - 64 : 8 = 90 - 8 = 82<br/>\nв) (75 + 25) - (94 + 3) = 100 - 97 = 3\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-67/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a4743e5e351537cbad098943e1ee9ca9d7f0a8dfcee84312fdce00b63daa4151', '3,6,8,17,25,64,75,90,94', '["найди"]'::jsonb, 'найди значения выражений:а) 6*x, если x=17 б) 90-у:8, если у=64 в) (75+а)-(94+b), если а=25, b=3');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 67, '8', 1, 'Викторина «В мире музыки» Вычисли. Расшифруй фамилии известных композиторов. Узнай, в какое время и в какой стране они жили. Слушаешь ли ты их музыку?', '</p> \n<p class="text"><b>Викторина «В мире музыки»</b><br/>\nВычисли. <br/>\nРасшифруй фамилии известных композиторов. Узнай, в какое время и в какой стране они жили. Слушаешь ли ты их музыку?\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica67-nomer8.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 67, номер 8, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 67, номер 8, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica67-nomer8-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 67, номер 8-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 67, номер 8-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-67/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica67-nomer8.jpg', 'peterson/3/part3/page67/task8_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica67-nomer8-1.jpg', 'peterson/3/part3/page67/task8_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '92138266812b9e3f802baf5682df8ca450a4bc76b1d5967d43b27fd08c84667f', NULL, '["вычисли"]'::jsonb, 'викторина "в мире музыки" вычисли. расшифруй фамилии известных композиторов. узнай, в какое время и в какой стране они жили. слушаешь ли ты их музыку');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 67, '9', 2, 'Разбей на классы и прочитай числа: 3609, 92820, 720053, 9113004, 50886999, 45012870, 5380024597, 12345678910, 376000000200.', '</p> \n<p class="text">Разбей на классы и прочитай числа:</p> \n\n<p class="description-text"> \n3609, 92820, 720053, 9113004, 50886999,<br/>\n45012870, 5380024597, 12345678910, 376000000200.\n</p>', '3 609 три тысячи шестьсот девять, 92 820 девяносто две тысячи восемьсот двадцать, 720 053 семьсот двадцать тысяч пятьдесят три, 9 113 004 девять миллионов сто тринадцать тысяч четыре, 50 886 999 пятьдесят миллионов восемьсот восемьдесят шесть тысяч девятьсот девяносто девять, 45 012 870 сорок пять миллионов двенадцать тысяч восемьсот семьдесят, 5 380 024 597 пять миллиарда триста восемьдесят миллионов двадцать четыре тысячи пятьсот девяносто семь, 12 345 678 910 двенадцать миллиарда триста сорок пять миллиона шестьсот семьдесят восемь тысячи девятьсот десять, 376 000 000 200 триста семьдесят шесть миллиарда двести.', '<p>\n3 609 три тысячи шестьсот девять,   <br/>\n92 820 девяносто две тысячи восемьсот двадцать, <br/>  \n720 053 семьсот двадцать тысяч пятьдесят три,  <br/> \n9 113 004 девять миллионов сто тринадцать тысяч четыре,<br/>   \n50 886 999 пятьдесят миллионов восемьсот восемьдесят шесть тысяч девятьсот девяносто девять,<br/>\n45 012 870 сорок пять миллионов двенадцать тысяч восемьсот семьдесят,   	<br/>\n5 380 024 597 пять миллиарда триста восемьдесят миллионов двадцать четыре тысячи пятьсот девяносто семь,  <br/> \n12 345 678 910 двенадцать миллиарда триста сорок пять миллиона шестьсот семьдесят восемь тысячи девятьсот десять,   	\n376 000 000 200 триста семьдесят шесть миллиарда двести.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-67/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6cf6ee594325f34ff7315d83dedafea3a35be8089590eb620e51480c0e1d3533', '3609,92820,720053,9113004,45012870,50886999,5380024597,12345678910,376000000200', '["раз"]'::jsonb, 'разбей на классы и прочитай числа:3609, 92820, 720053, 9113004, 50886999, 45012870, 5380024597, 12345678910, 376000000200');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 67, '10', 3, 'а) Какое число идёт при счёте за числом 82355, 739999? б) Какое число предшествует в натуральном ряду числу 3480, 26000?', '</p> \n<p class="text">а) Какое число идёт при счёте за числом 82355, 739999?<br/>\nб) Какое число предшествует в натуральном ряду числу 3480, 26000?\n</p>', 'а) 82356, 740000. б) 3479, 25999.', '<p>\nа) 82356, 740000.<br/>\nб) 3479, 25999.\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Какое число идёт при счёте за числом 82355, 739999?","solution":"82356, 740000."},{"letter":"б","condition":"Какое число предшествует в натуральном ряду числу 3480, 26000?","solution":"3479, 25999."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-67/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '53c046f744e75a8a1e3d599a4fc93096ea9df58cf65dabbf8b408652dc27c674', '3480,26000,82355,739999', NULL, 'а) какое число идёт при счёте за числом 82355, 739999? б) какое число предшествует в натуральном ряду числу 3480, 26000');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 67, '11', 4, 'Запиши в виде суммы разрядных слагаемых числа 817, 3029, 53082, 706480.', '</p> \n<p class="text">Запиши в виде суммы разрядных слагаемых числа 817, 3029, 53082, 706480.</p>', '8 · 100 + 1 · 10 + 7 · 1, 3 · 1000 + 2 · 10 + 9 · 1, 5 · 10000 + 3 · 1000 + 8 · 10 + 2 · 1, 7 · 100000 + 6 · 1000 + 4 · 100 + 8 · 10.', '<p>\n8 · 100 + 1 · 10 + 7 · 1,<br/>\n3 · 1000 + 2 · 10 + 9 · 1,<br/> \n5 · 10000 + 3 · 1000 + 8 · 10 + 2 · 1,<br/> \n7 · 100000 + 6 · 1000 + 4 · 100 + 8 · 10.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-67/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ab11c21d2ff501e79a97cfcdfce1b3692b7d2b42c500a4b9405cb6fe06ecf247', '817,3029,53082,706480', '["раз"]'::jsonb, 'запиши в виде суммы разрядных слагаемых числа 817, 3029, 53082, 706480');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 68, '12', 0, 'Запиши цифрами числа: а) 4 тыс. 549 ед.        д) 439 млн 972 тыс. 508 ед. б) 8 тыс. 20 ед.          е) 5 млн 2 тыс. 16 ед. в) 76 тыс. 9 ед.          ж) 29 млн 396 ед. г) 318 тыс. 690 ед.    з) 4 млн 7 тыс.', '</p> \n<p class="text">Запиши цифрами числа:</p> \n\n<p class="description-text"> \nа) 4 тыс. 549 ед.        д) 439 млн 972 тыс. 508 ед.<br/>\nб) 8 тыс. 20 ед.          е) 5 млн 2 тыс. 16 ед.<br/>\nв) 76 тыс. 9 ед.          ж) 29 млн 396 ед.<br/>\nг) 318 тыс. 690 ед.    з) 4 млн 7 тыс.\n</p>', 'а) 4549            д) 439972508 б) 8020            е) 5002016 в) 76009          ж) 29000396 г) 318690        з) 4007000', '<p>\nа) 4549            д) 439972508<br/>\nб) 8020            е) 5002016<br/>\nв) 76009          ж) 29000396<br/>\nг) 318690        з) 4007000\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-68/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '9878c9b62f29067ed13b41f84d987ae636a86e20a89ff820a605dfeb80960dd6', '2,4,5,7,8,9,16,20,29,76', NULL, 'запиши цифрами числа:а) 4 тыс. 549 ед.        д) 439 млн 972 тыс. 508 ед. б) 8 тыс. 20 ед.          е) 5 млн 2 тыс. 16 ед. в) 76 тыс. 9 ед.          ж) 29 млн 396 ед. г) 318 тыс. 690 ед.    з) 4 млн 7 тыс');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 68, '13', 1, 'Найди в таблице: а) наибольшее четырёхзначное число; б) наименьшее четырёхзначное число; в) наименьшее трёхзначное число с цифрой 8 в разряде единиц; г) наибольшее четырёхзначное число с цифрой 5 в разряде десятков; д) наибольшее пятизначное число с цифрой 7 в разряде сотен; е) наибольшее четырёхзначное число с разными цифрами; ж) наименьшее четырёхзначное число с разными цифрами.', '</p> \n<p class="text">Найди в таблице:<br/>\nа) наибольшее четырёхзначное число;<br/>\nб) наименьшее четырёхзначное число;<br/>\nв) наименьшее трёхзначное число с цифрой 8 в разряде единиц;<br/>\nг) наибольшее четырёхзначное число с цифрой 5 в разряде десятков;<br/>\nд) наибольшее пятизначное число с цифрой 7 в разряде сотен;<br/>\nе) наибольшее четырёхзначное число с разными цифрами;<br/>\nж) наименьшее четырёхзначное число с разными цифрами.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica68-nomer13.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 68, номер 13, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 68, номер 13, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica68-nomer13-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 68, номер 13-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 68, номер 13-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-68/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica68-nomer13.jpg', 'peterson/3/part3/page68/task13_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica68-nomer13-1.jpg', 'peterson/3/part3/page68/task13_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '531ede41aa3f00320b70d632b0deb93ba0d97594c423975859335d3d0c70d4cb', '5,7,8', '["найди","больше","меньше","раз"]'::jsonb, 'найди в таблице:а) наибольшее четырёхзначное число; б) наименьшее четырёхзначное число; в) наименьшее трёхзначное число с цифрой 8 в разряде единиц; г) наибольшее четырёхзначное число с цифрой 5 в разряде десятков; д) наибольшее пятизначное число с цифрой 7 в разряде сотен; е) наибольшее четырёхзначное число с разными цифрами; ж) наименьшее четырёхзначное число с разными цифрами');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 68, '14', 2, 'Прочитай число 28057000094. Какая цифра стоит в разряде единиц миллионов этого числа? Сколько в нём всего миллионов?', '</p> \n<p class="text">Прочитай число 28057000094. Какая цифра стоит в разряде единиц миллионов этого числа? Сколько в нём всего миллионов? </p>', 'Двадцать восемь миллиарда пятьдесят семь миллионов девяносто четыре 7 стоит в разряде единиц миллионов этого числа. 28057 в нём всего миллионов.', '<p>\nДвадцать восемь миллиарда пятьдесят семь миллионов девяносто четыре<br/>\n7 стоит в разряде единиц миллионов этого числа. 28057 в нём всего миллионов.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-68/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b7abd0431a4a5326e379006a62ff6ece7cada257f1b1ad1433414d253543331c', '28057000094', '["раз"]'::jsonb, 'прочитай число 28057000094. какая цифра стоит в разряде единиц миллионов этого числа? сколько в нём всего миллионов');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 68, '15', 3, 'Сравни с помощью знаков >, <, = : 352     235        4003     999 98     3060        5300     5299 7425     74000 82016     82106', '</p> \n<p class="text">Сравни с помощью знаков &gt;, &lt;, = : </p> \n\n<p class="description-text"> \n352 <span class="okon">   </span> 235        4003 <span class="okon">   </span> 999<br/>\n98 <span class="okon">   </span> 3060        5300 <span class="okon">   </span> 5299  <br/><br/>	\n\n7425 <span class="okon">   </span> 74000<br/>\n82016 <span class="okon">   </span> 82106\n</p>', '352 > 235        4003 > 999 98 < 3060        5300 > 5299 7425 < 74000 82016 < 82106', '<p>\n352 &gt; 235        4003 &gt; 999 <br/> 		\n98 &lt; 3060        5300 &gt; 5299 <br/><br/> 	\n\n7425 &lt; 74000 <br/>\n82016 &lt; 82106\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-68/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '20eeef3993837e5a7a8231fa867ad0e0a458acc773bd1fc0d442a9c89eeb0f19', '98,235,352,999,3060,4003,5299,5300,7425,74000', '["сравни"]'::jsonb, 'сравни с помощью знаков>,<,=:352     235        4003     999 98     3060        5300     5299 7425     74000 82016     82106');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 68, '16', 4, 'Выполни действия. Проверь с помощью калькулятора. а) 305246 - 21237           д) 23715926 + 3276315 б) 524032 + 78369           е) 944502483 - 25360157 в) 4061497 + 938708      ж) 726524996 + 873475104 г) 80000425 - 536842     з) 120036705 - 92759318', '</p> \n<p class="text">Выполни действия. Проверь с помощью калькулятора.</p> \n\n<p class="description-text"> \nа) 305246 - 21237           д) 23715926 + 3276315<br/>\nб) 524032 + 78369           е) 944502483 - 25360157<br/>\nв) 4061497 + 938708      ж) 726524996 + 873475104<br/>\nг) 80000425 - 536842     з) 120036705 - 92759318\n</p>', 'а) 305246 - 21237 = 284009', '<p>\nа) 305246 - 21237 = 284009\n</p>\n\n<div class="img-wrapper-460">\n<img width="180" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica68-nomer16.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 68, номер 16, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 68, номер 16, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-68/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica68-nomer16.jpg', 'peterson/3/part3/page68/task16_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '948c1acae8bfd0ace1596a8ea6d76babad0ecf9a8e68a24312fba23dd2322b0f', '21237,78369,305246,524032,536842,938708,3276315,4061497,23715926,25360157', NULL, 'выполни действия. проверь с помощью калькулятора. а) 305246-21237           д) 23715926+3276315 б) 524032+78369           е) 944502483-25360157 в) 4061497+938708      ж) 726524996+873475104 г) 80000425-536842     з) 120036705-92759318');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 68, '17', 5, 'Составь программу действий и вычисли: (9452 + 13808) - (55400 - 39326) + 1227', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \n(9452 + 13808) - (55400 - 39326) + 1227\n</p>', '(9452 + 13808) - (55400 - 39326) + 1227 = 8413 9452 + 13808 = 23260', '<p>\n(9452 + 13808) - (55400 - 39326) + 1227 = 8413<br/>\n9452 + 13808 = 23260\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica68-nomer17.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 68, номер 17, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 68, номер 17, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-68/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica68-nomer17.jpg', 'peterson/3/part3/page68/task17_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '88424666736e1174177d9f57253055aac92205d783c33aa2b0c37cc33aa18769', '1227,9452,13808,39326,55400', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:(9452+13808)-(55400-39326)+1227');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 69, '18', 0, 'а) На сколько число 32856 меньше числа 40912? б) На сколько число 51045 больше числа 6387?', '</p> \n<p class="text">а) На сколько число 32856 меньше числа 40912?<br/> \nб) На сколько число 51045 больше числа 6387?\n</p>', 'а) 40912 - 32856 = 8056', '<p>\nа) 40912 - 32856 = 8056\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer18.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 69, номер 18, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 69, номер 18, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"На сколько число 32856 меньше числа 40912?","solution":"40912 - 32856 = 8056"},{"letter":"б","condition":"На сколько число 51045 больше числа 6387?","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-69/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer18.jpg', 'peterson/3/part3/page69/task18_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e04a4b31f9c84a835437594d1cbf163d9cb43f2569b6ea1191c94a88d0965151', '6387,32856,40912,51045', '["больше","меньше"]'::jsonb, 'а) на сколько число 32856 меньше числа 40912? б) на сколько число 51045 больше числа 6387');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 69, '19', 1, 'Найди неизвестные числа:', '</p> \n<p class="text">Найди неизвестные числа:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer19.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 69, номер 19, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 69, номер 19, год 2022."/>\n</div>\n</div>', '98002 - 72 = 97930', '<p>\n98002 - 72 = 97930\n</p>\n\n<div class="img-wrapper-460">\n<img width="140" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer19-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 69, номер 19-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 69, номер 19-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-69/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer19.jpg', 'peterson/3/part3/page69/task19_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer19-1.jpg', 'peterson/3/part3/page69/task19_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '12b19e3ab6d17dc350097cbf06f9287fc0ef68eecdc42673d8eb110691c46ba5', NULL, '["найди"]'::jsonb, 'найди неизвестные числа');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 69, '20', 2, 'Вычисли длину неизвестного отрезка, используя взаимосвязь между частью и целым.', '</p> \n<p class="text">Вычисли длину неизвестного отрезка, используя взаимосвязь между частью и целым.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer20.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 69, номер 20, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 69, номер 20, год 2022."/>\n</div>\n</div>', 'а) 18 - 6 - 9 = 9 - 6 = 3 (см) б) 36 - 23 = 13 (мм), 36 - 29 = 7 (мм), 36 - 13 - 7 = 23 - 7 = 16 (мм) в) 7 + 14 + 16 = 37 (дм) г) 34 - 6 = 28 (м), 48 - 28 = 20 (м)', '<p>\nа) 18 - 6 - 9 = 9 - 6 = 3 (см)<br/>\nб) 36 - 23 = 13 (мм), 36 - 29 = 7 (мм), 36 - 13 - 7 = 23 - 7 = 16 (мм)<br/>\nв) 7 + 14 + 16 = 37 (дм)<br/>\nг) 34 - 6 = 28 (м), 48 - 28 = 20 (м)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-69/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer20.jpg', 'peterson/3/part3/page69/task20_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd8a7e1aaf61da002063b54e808f0fc92ba77a102c33f6e4277217ed26b2f836a', NULL, '["вычисли"]'::jsonb, 'вычисли длину неизвестного отрезка, используя взаимосвязь между частью и целым');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 69, '21', 3, 'На отрезке MK длиной 26 см отметили точку A так, что AM = 19 см, и точку B так, что BK = 12 см. Найди длину отрезка AB.', '</p> \n<p class="text">На отрезке MK длиной 26 см отметили точку A так, что AM = 19 см, и точку B так, что BK = 12 см. Найди длину отрезка AB.</p>', '26 - 19 = 7 (см) – АК 26 - 12 = 14 (см) – ВМ 26 - 14 - 7 = 12 - 7 = 5 (см) Ответ: АВ равен 5 сантиметров.', '<p>\n26 - 19 = 7 (см) – АК<br/>\n26 - 12 = 14 (см) – ВМ<br/>\n26 - 14 - 7 = 12 - 7 = 5 (см)<br/>\n<b>Ответ:</b> АВ равен 5 сантиметров.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-69/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6d1631d7180c3064ba1e51fb7984a348b44bb71311ca8f4dc87cf1bdf3e9c767', '12,19,26', '["найди"]'::jsonb, 'на отрезке mk длиной 26 см отметили точку a так, что am=19 см, и точку b так, что bk=12 см. найди длину отрезка ab');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 69, '22', 4, 'Составь все возможные равенства из чисел 3409, 596, 4005. Как найти целое? Как найти часть?', '</p> \n<p class="text">Составь все возможные равенства из чисел 3409, 596, 4005. Как найти целое? Как найти часть?</p>', '3409 + 596 = 4005 4005 - 596 = 3409 596 + 3409 = 4005 4005 − 3409 = 596. Находим целое, складываем части. Находим часть, из целого вычитаем другу часть.', '<p>\n3409 + 596 = 4005<br/>\n4005 - 596 = 3409<br/>\n596 + 3409 = 4005<br/>\n4005 − 3409 = 596.<br/>\nНаходим целое, складываем части.<br/>\nНаходим часть, из целого вычитаем другу часть.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-69/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1218f086f6d42d554e59d65c4dd1f999d286ac8bd24149d16cf5f2f2765a9391', '596,3409,4005', NULL, 'составь все возможные равенства из чисел 3409, 596, 4005. как найти целое? как найти часть');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 69, '23', 5, 'Реши уравнения с комментированием и сделай проверку: а) x - 18910 = 3459 б) 6207 + y = 50000 в) 45180 - z = 7652', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа)  x - 18910 = 3459  <br/>\nб) 6207 + y = 50000  <br/>\nв) 45180 - z = 7652\n</p>', 'а) x - 18910 = 3459 Что бы найти уменьшаемое надо вычитаемое прибавить к разности х = 18910 + 3459', '<p>\nа)  x - 18910 = 3459<br/>\nЧто бы найти уменьшаемое надо вычитаемое прибавить к разности<br/>\nх = 18910 + 3459\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer23.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 69, номер 23, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 69, номер 23, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-69/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer23.jpg', 'peterson/3/part3/page69/task23_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '37fd188db0fd48a39c4f787cfd2e04854d0dd14c9b87fba7f1e6f061c4149e56', '3459,6207,7652,18910,45180,50000', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) x-18910=3459 б) 6207+y=50000 в) 45180-z=7652');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 69, '24', 6, 'В автомобильных гонках участвовало три команды. У первой команды было 24 автомобиля, что на 3 автомобиля больше, чем у второй команды. Сколько автомобилей было у третьей команды, если всего в гонках участвовало 80 автомобилей?', '</p> \n<p class="text">В автомобильных гонках участвовало три команды. У первой команды было 24 автомобиля, что на 3 автомобиля больше, чем у второй команды. Сколько автомобилей было у третьей команды, если всего в гонках участвовало 80 автомобилей?</p>', '80 - 24 - (24 - 3) = 56 - 21 = 35 (автомобилей) Ответ: 35 автомобилей было у третьей команды, если всего в гонках участвовало 80 автомобилей.', '<p>\n80 - 24 - (24 - 3) = 56 - 21 = 35 (автомобилей)<br/>\n<b>Ответ:</b> 35 автомобилей было у третьей команды, если всего в гонках участвовало 80 автомобилей.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-69/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a8e6e9837c02ae70a6778cfc6acc503321eea0bc1e4de1048d2f0d5311f24882', '3,24,80', '["больше"]'::jsonb, 'в автомобильных гонках участвовало три команды. у первой команды было 24 автомобиля, что на 3 автомобиля больше, чем у второй команды. сколько автомобилей было у третьей команды, если всего в гонках участвовало 80 автомобилей');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 69, '25', 7, 'Стакан чая стоит 5 р., что на 12 р. дешевле, чем булочка. А чай вместе с булочкой стоят столько же, сколько апельсин. Сколько рублей надо заплатить за стакан чая, булочку и апельсин?', '</p> \n<p class="text">Стакан чая стоит 5 р., что на 12 р. дешевле, чем булочка. А чай вместе с булочкой стоят столько же, сколько апельсин. Сколько рублей надо заплатить за стакан чая, булочку и апельсин?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="300" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer25.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 69, номер 25, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 69, номер 25, год 2022."/>\n</div>\n</div>', '(5 + 12 + 5) + (5 + 12 + 5) = 22 + 22 = 44 (р.) Ответ: 44 рублей надо заплатить за стакан чая, булочку и апельсин.', '<p>\n(5 + 12 + 5) + (5 + 12 + 5) = 22 + 22 = 44 (р.) <br/>\n<b>Ответ:</b> 44 рублей надо заплатить за стакан чая, булочку и апельсин.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-69/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica69-nomer25.jpg', 'peterson/3/part3/page69/task25_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '82158ddff0961ea40ad6a664312c718fe8e36343c0118a52650acd4f9eee7f12', '5,12', NULL, 'стакан чая стоит 5 р., что на 12 р. дешевле, чем булочка. а чай вместе с булочкой стоят столько же, сколько апельсин. сколько рублей надо заплатить за стакан чая, булочку и апельсин');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 70, '26', 0, 'Придумай задачу по схеме и реши её:', '</p> \n<p class="text">Придумай задачу по схеме и реши её:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="290" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica70-nomer26.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 70, номер 26, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 70, номер 26, год 2022."/>\n</div>\n</div>', 'Карамелек 18 штук, батончиков на 2 штуки больше карамелек. Ирисок в 3 раза меньше карамелек. Шоколадных столько, сколько вместе карамельки, батончики и ириски. На сколько карамелек больше ирисок? Сколько конфет всего? 18 - 18 : 3 = 18 – 6 = 12 (штук) – ирисок (18 + (18 + 2) + 12) + (18 + (18 + 2) + 12) = (30 + 20) + (30 + 20) = 50 + 50 = 100 (штук) Ответ: 100 конфет всего.', '<p>\nКарамелек 18 штук, батончиков на 2 штуки больше карамелек. Ирисок в 3 раза меньше карамелек. Шоколадных столько, сколько вместе карамельки, батончики и ириски. На сколько карамелек больше ирисок? Сколько конфет всего?<br/>\n18 - 18 : 3 = 18 – 6 = 12 (штук) – ирисок<br/>\n(18 + (18 + 2) + 12) + (18 + (18 + 2) + 12) = (30 + 20) + (30 + 20) = 50 + 50 = 100 (штук)<br/>\n<b>Ответ:</b> 100 конфет всего.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-70/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica70-nomer26.jpg', 'peterson/3/part3/page70/task26_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '91eb9943371a5d3fb392d4c938245bbb03f4db7e19266bec873675699a7b3f39', NULL, '["реши"]'::jsonb, 'придумай задачу по схеме и реши её');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 70, '27', 1, 'а) Первая сторона треугольника равна 14 дм, а вторая сторона в 2 раза больше первой. Найди третью сторону треугольника, если его периметр равен 64 дм. б) Длина первой стороны треугольника 24 см. Это в 2 раза больше длины второй стороны и на 5 см меньше длины третьей. Найди периметр этого треугольника.', '</p> \n<p class="text">а) Первая сторона треугольника равна 14 дм, а вторая сторона в 2 раза больше первой. Найди третью сторону треугольника, если его периметр равен 64 дм.<br/>\nб) Длина первой стороны треугольника 24 см. Это в 2 раза больше длины второй стороны и на 5 см меньше длины третьей. Найди периметр этого треугольника.\n</p>', 'а) 64 - 14 - 14 · 2 = 50 - 7 = 43 (дм) б) 24 + (24 : 2) + (24 : 2 + 5) = 24 + 12 + 17 = 53 (см)', '<p>\nа) 64 - 14 - 14 · 2 = 50 - 7 = 43 (дм)<br/>\nб) 24 + (24 : 2) + (24 : 2 + 5) = 24 + 12 + 17 = 53 (см)\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Первая сторона треугольника равна 14 дм, а вторая сторона в 2 раза больше первой. Найди третью сторону треугольника, если его периметр равен 64 дм.","solution":"64 - 14 - 14 · 2 = 50 - 7 = 43 (дм)"},{"letter":"б","condition":"Длина первой стороны треугольника 24 см. Это в 2 раза больше длины второй стороны и на 5 см меньше длины третьей. Найди периметр этого треугольника.","solution":"24 + (24 : 2) + (24 : 2 + 5) = 24 + 12 + 17 = 53 (см)"}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-70/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '69f60ff580651bc9cb1a3e2994d87f3c70dec0e039e9af81c735b2e2c4942f5d', '2,5,14,24,64', '["найди","периметр","сторона","больше","меньше","раз","раза"]'::jsonb, 'а) первая сторона треугольника равна 14 дм, а вторая сторона в 2 раза больше первой. найди третью сторону треугольника, если его периметр равен 64 дм. б) длина первой стороны треугольника 24 см. это в 2 раза больше длины второй стороны и на 5 см меньше длины третьей. найди периметр этого треугольника');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 70, '28', 2, 'а) Ширина прямоугольника равна 84 м, что на 6 м меньше его длины. Найди периметр и площадь этого прямоугольника. б) Площадь прямоугольника равна 750 м 2 , а длина – 30 м. На сколько метров ширина этого прямоугольника меньше длины?', '</p> \n<p class="text">а) Ширина прямоугольника равна 84 м, что на 6 м меньше его длины. Найди периметр и площадь этого прямоугольника.<br/>\nб) Площадь прямоугольника равна 750 м<sup>2</sup>, а длина – 30 м. На сколько метров ширина этого прямоугольника меньше длины?\n</p>', 'а) (84 + 84 + 6) · 2 = 174 · 2 = 348 (м) – периметр 84 · (84 + 6) = 84 · 90 = 7560 (м 2 ) – площадь Ответ: 348 м - периметр и 7560 м 2 – площадь этого прямоугольника. б) 30 - 750 : 30 = 30 - 25 = 5 (м) Ответ: на 5 метров ширина этого прямоугольника меньше длины.', '<p>\nа) (84 + 84 + 6) · 2 = 174 · 2 = 348 (м) – периметр<br/>\n84 · (84 + 6) = 84 · 90 = 7560 (м<sup>2</sup>) – площадь<br/>\n<b>Ответ:</b> 348 м - периметр и 7560 м<sup>2</sup> – площадь этого прямоугольника.<br/><br/>\nб) 30 - 750 : 30 = 30 - 25 = 5 (м)<br/>\n<b>Ответ:</b> на 5 метров ширина этого прямоугольника меньше длины.\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Ширина прямоугольника равна 84 м, что на 6 м меньше его длины. Найди периметр и площадь этого прямоугольника.","solution":"(84 + 84 + 6) · 2 = 174 · 2 = 348 (м) – периметр 84 · (84 + 6) = 84 · 90 = 7560 (м 2 ) – площадь Ответ: 348 м - периметр и 7560 м 2 – площадь этого прямоугольника."},{"letter":"б","condition":"Площадь прямоугольника равна 750 м 2 , а длина – 30 м. На сколько метров ширина этого прямоугольника меньше длины?","solution":"30 - 750 : 30 = 30 - 25 = 5 (м) Ответ: на 5 метров ширина этого прямоугольника меньше длины."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-70/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2328a327dd58197a20e61afc74ef651d9c21f0faf1c575d06d1f06576460e057', '2,6,30,84,750', '["найди","периметр","площадь","меньше"]'::jsonb, 'а) ширина прямоугольника равна 84 м, что на 6 м меньше его длины. найди периметр и площадь этого прямоугольника. б) площадь прямоугольника равна 750 м 2 , а длина-30 м. на сколько метров ширина этого прямоугольника меньше длины');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 70, '29', 3, 'Вычисли площади фигур:', '</p> \n<p class="text">Вычисли площади фигур:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica70-nomer29.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 70, номер 29, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 70, номер 29, год 2022."/>\n</div>\n</div>', 'а) 8 · 8 + 5 · 3 = 64 + 15 = 79 (м 2 ) б) 56 · 40 - 20 · 14 = 2240 - 280 = 1960 (см 2 )', '<p>\nа) 8 · 8 + 5 · 3 = 64 + 15 = 79 (м<sup>2</sup>)<br/>\nб) 56 · 40 - 20 · 14 = 2240 - 280 = 1960 (см<sup>2</sup>)\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-70/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica70-nomer29.jpg', 'peterson/3/part3/page70/task29_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '72ab4b5b1eac9ac1763724a9f87652fc5a901c015f6f7a864559ac7af28511ff', NULL, '["вычисли"]'::jsonb, 'вычисли площади фигур');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 70, '30', 4, 'Построй квадрат со стороной 4 см. Затем построй прямоугольник, ширина которого на 2 см меньше, а длина – на 2 см больше стороны квадрата.', '</p> \n<p class="text">Построй квадрат со стороной 4 см. Затем построй прямоугольник, ширина которого на 2 см меньше, а длина – на 2 см больше стороны квадрата.</p>', '', '<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica70-nomer30.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 70, номер 30, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 70, номер 30, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-70/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica70-nomer30.jpg', 'peterson/3/part3/page70/task30_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'af3c10d098ad7dd4fe9e9109863ac7309ff291420c8f1cc66cc625842e68a7f8', '2,4', '["больше","меньше"]'::jsonb, 'построй квадрат со стороной 4 см. затем построй прямоугольник, ширина которого на 2 см меньше, а длина-на 2 см больше стороны квадрата');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 70, '31', 5, 'Длина классной комнаты 12 м, ширина 10 м, а высота 4 м. Найди её объём.', '</p> \n<p class="text">Длина классной комнаты 12 м, ширина 10 м, а высота 4 м. Найди её объём.</p>', '12 · 10 · 4 = 480 (м 3 ) Ответ: объем классной комнаты равен 480 м 3 .', '<p>\n12 · 10 · 4 = 480 (м<sup>3</sup>)<br/>\n<b>Ответ:</b> объем классной комнаты равен 480 м<sup>3</sup>.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-70/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3b0f4605d113b415d564e89a744dadc8d32c4a6b55f95cfcbd7a6f72a6c3a882', '4,10,12', '["найди"]'::jsonb, 'длина классной комнаты 12 м, ширина 10 м, а высота 4 м. найди её объём');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 70, '32', 6, 'а) Вырази число 15340 в десятках; в сотнях и единицах; в тысячах и единицах. б) Вырази 15340 см в дециметрах; в метрах и дециметрах. в) Вырази 15340 м в километрах и метрах. г) Вырази 15340 г в килограммах и граммах. д) Вырази 15340 кг в центнерах и килограммах; в тоннах и килограммах.', '</p> \n<p class="text">а) Вырази число 15340 в десятках; в сотнях и единицах; в тысячах и единицах.<br/>\nб) Вырази 15340 см в дециметрах; в метрах и дециметрах.<br/>\nв) Вырази 15340 м в километрах и метрах.<br/>\nг) Вырази 15340 г в килограммах и граммах.<br/>\nд) Вырази 15340 кг в центнерах и килограммах; в тоннах и килограммах.\n</p>', 'а) 1534 десятков, 153 десяток + 40 единиц, 15 тысяч + 340 единиц. б) 1534 дм, 153 м 4 дм, в) 15 км 340 м г) 15кг 340 г д) 153 ц 40 кг, 15 т 340 кг.', '<p>\nа) 1534 десятков, 153 десяток + 40 единиц, 15 тысяч + 340 единиц.<br/>\nб) 1534 дм, 153 м 4 дм, <br/>\nв) 15 км 340 м<br/>\nг) 15кг 340 г <br/>\nд) 153 ц 40 кг, 15 т 340 кг.\n</p>', '', '', TRUE, '[{"letter":"а","condition":"Вырази число 15340 в десятках; в сотнях и единицах; в тысячах и единицах.","solution":"1534 десятков, 153 десяток + 40 единиц, 15 тысяч + 340 единиц."},{"letter":"б","condition":"Вырази 15340 см в дециметрах; в метрах и дециметрах.","solution":"1534 дм, 153 м 4 дм,"},{"letter":"в","condition":"Вырази 15340 м в километрах и метрах.","solution":"15 км 340 м"},{"letter":"г","condition":"Вырази 15340 г в килограммах и граммах.","solution":"15кг 340 г"},{"letter":"д","condition":"Вырази 15340 кг в центнерах и килограммах; в тоннах и килограммах.","solution":"153 ц 40 кг, 15 т 340 кг."}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-70/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '249333f9920802d7f6bf3c41310e1341ed2dac10f01fec676e3916a2b18183d2', '15340', '["раз"]'::jsonb, 'а) вырази число 15340 в десятках; в сотнях и единицах; в тысячах и единицах. б) вырази 15340 см в дециметрах; в метрах и дециметрах. в) вырази 15340 м в километрах и метрах. г) вырази 15340 г в килограммах и граммах. д) вырази 15340 кг в центнерах и килограммах; в тоннах и килограммах');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 71, '33', 0, 'Вспомни таблицы мер длины, площади, объёма, массы. Выполни действия: а) 5 м 96 см + 32 дм 4 см б) 6 дм 3 см 2 мм - 48 см в) 4 км 788 м + 6 км 20 м г) 12 км 52 м - 8 км 258 м д) 9 кг 200 г - 5 кг 540 г е) 17 ц 69 кг + 3 т 831 ж) 15 м 2 2 см 2 - 9 м 2 5 дм 2 27 см 2 з) 12 дм 3 - 3 дм 3 4 см 3', '</p> \n<p class="text">Вспомни таблицы мер длины, площади, объёма, массы. Выполни действия:</p> \n\n<p class="description-text"> \nа) 5 м 96 см + 32 дм 4 см <br/> 	 \nб) 6 дм 3 см 2 мм - 48 см <br/> 	\nв) 4 км 788 м + 6 км 20 м  <br/>	\nг) 12 км 52 м - 8 км 258 м <br/><br/> 	\n\nд) 9 кг 200 г - 5 кг 540 г<br/>\nе) 17 ц 69 кг + 3 т 831 <br/>\nж) 15 м<sup>2</sup> 2 см<sup>2</sup> - 9 м<sup>2</sup> 5 дм<sup>2</sup> 27 см<sup>2</sup> <br/>\nз) 12 дм<sup>3</sup> - 3 дм<sup>3</sup> 4 см<sup>3</sup>\n\n</p>', 'а) 5 м 96 см + 32 дм 4 см = 50 дм + 32 дм + 96 см + 4 см = 82 дм 100 см = 9 м 2 дм б) 6 дм 3 см 2 мм - 48 см = 1 дм 53 см - 48 см + 2 мм = 1 дм 5 см 2 мм в) 4 км 788 м + 6 км 20 м = 10 км 808 м г) 12 км 52 м - 8 км 258 м = 3 км 1052 м - 258 м = 3 км 794 м д) 9 кг 200 г - 5 кг 540 г = 3 кг 1200 г - 540 г = 3 кг 660 г е) 17 ц 69 кг + 3 т 831 кг = 17 ц + 30 ц 900 кг = 47 ц 900 г = 4 т 7 ц 900 г ж) 15 м 2 2 см 2 - 9 м 2 5 дм 2 27 см 2 = (14 м 2 - 9 м 2 ) + (10002 см 2 - 527 см 2 ) = 5 м 2 9475 см 2 = 5 м 2 94 дм 2 75 см 2 з) 12 дм 3 - 3 дм 3 4 см 3 = 9 дм 3 4 см 3 .', '<p>\nа) 5 м 96 см + 32 дм 4 см = 50 дм + 32 дм + 96 см + 4 см = 82 дм 100 см = 9 м 2 дм<br/>\nб) 6 дм 3 см 2 мм - 48 см = 1 дм 53 см - 48 см + 2 мм = 1 дм 5 см 2 мм  <br/>\nв) 4 км 788 м + 6 км 20 м = 10 км 808 м<br/>\nг) 12 км 52 м - 8 км 258 м = 3 км 1052 м - 258 м = 3 км 794 м <br/>\nд) 9 кг 200 г - 5 кг 540 г = 3 кг 1200 г - 540 г = 3 кг 660 г<br/>\nе) 17 ц 69 кг + 3 т 831 кг = 17 ц + 30 ц 900 кг = 47 ц 900 г = 4 т 7 ц 900 г<br/>\nж) 15 м<sup>2</sup> 2 см<sup>2</sup> - 9 м<sup>2</sup> 5 дм<sup>2</sup> 27 см<sup>2</sup> = (14 м<sup>2</sup> - 9 м<sup>2</sup>) + (10002 см<sup>2</sup> - 527 см<sup>2</sup>) = 5 м<sup>2</sup> 9475 см<sup>2</sup> = 5 м<sup>2</sup> 94 дм<sup>2</sup> 75 см<sup>2</sup> <br/>\nз) 12 дм<sup>3</sup> - 3 дм<sup>3</sup> 4 см<sup>3</sup> = 9 дм<sup>3</sup> 4 см<sup>3</sup>.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-71/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b7c548735146100ea6b58abfe8a5e8f6f9e056961a028b29a3eb6a27fa46dd85', '2,3,4,5,6,8,9,12,15,17', NULL, 'вспомни таблицы мер длины, площади, объёма, массы. выполни действия:а) 5 м 96 см+32 дм 4 см б) 6 дм 3 см 2 мм-48 см в) 4 км 788 м+6 км 20 м г) 12 км 52 м-8 км 258 м д) 9 кг 200 г-5 кг 540 г е) 17 ц 69 кг+3 т 831 ж) 15 м 2 2 см 2-9 м 2 5 дм 2 27 см 2 з) 12 дм 3-3 дм 3 4 см 3');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 71, '34', 1, 'Составь программу действий. Что ты замечаешь? (a + b) · c - d : (k + m) · n (a + b · c) - (d : k + m) · n', '</p> \n<p class="text">Составь программу действий. Что ты замечаешь?</p> \n\n<p class="description-text"> \n(a + b) · c - d : (k + m) · n <br/>      \n (a + b · c) - (d : k + m) · n\n</p>', '(a + b) · c - d : (k + m) · n a + b (a + b) · c k + m d : (k + m) d : (k + m) · n (a + b) · c - d : (k + m) · n (a + b · c) - (d : k + m) · n b · c a + b · c d : k d : k + m (d : k + m) · n (a + b · c) - (d : k + m) · n Порядок действий изменился от расположения скобок.', '<p>\n(a + b) · c - d : (k + m) · n<br/>    \na + b<br/>\n(a + b) · c<br/>\nk + m<br/>\nd : (k + m)<br/>\nd : (k + m) · n  <br/>  \n(a + b) · c - d : (k + m) · n <br/><br/>   \n\n(a + b · c) - (d : k + m) · n<br/>\nb · c<br/>\na + b · c<br/>\nd : k<br/>\nd : k + m<br/>\n(d : k + m) · n<br/>\n(a + b · c) - (d : k + m) · n<br/>\nПорядок действий изменился от расположения скобок.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-71/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ecdc2fc542b594102d2af97cc124c97fc8afc2f28683d88607c9c657362ebbca', NULL, NULL, 'составь программу действий. что ты замечаешь? (a+b)*c-d:(k+m)*n (a+b*c)-(d:k+m)*n');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 71, '35', 2, 'Какие знаки арифметических действий можно поставить вместо звёздочек? Возможны ли другие варианты? a * 0 = a            1 * a = a a * a = 0            a * 1 = a a * a = 1            a * 0 = 0 0 * a = 0            0 * a = a', '</p> \n<p class="text">Какие знаки арифметических действий можно поставить вместо звёздочек? Возможны ли другие варианты?</p> \n\n<p class="description-text"> \na * 0 = a            1 * a = a <br/>  	\na * a = 0            a * 1 = a<br/>  	    \na * a = 1            a * 0 = 0<br/>\n0 * a = 0            0 * a = a\n\n</p>', 'a + 0 = a или а - 0 = а      1 · a = a a - a = 0             a · 1 = a или а : 1 = а a : a = 1             a · 0 = 0 или а : 0 = 0 0 · a = 0 или 0 : а = 0      0 + a = a или 0 – а = а', '<p>\na + 0 = a или а - 0 = а      1 · a = a  <br/>						\na - a = 0             a · 1 = a или а : 1 = а <br/><br/> 	\n\na : a = 1             a · 0 = 0 или а : 0 = 0<br/>\n0 · a = 0 или 0 : а = 0      0 + a = a или 0 – а = а\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-71/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '91b8bfb359bac646868b35f691be6b756e196a727e19f5ad6817c24085d9f646', '0,1', NULL, 'какие знаки арифметических действий можно поставить вместо звёздочек? возможны ли другие варианты? a*0=a            1*a=a a*a=0            a*1=a a*a=1            a*0=0 0*a=0            0*a=a');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 71, '36', 3, 'Составь программу действий и вычисли: а) 24 : 1 - (4 · 5 - 14) · 4 + 8 : 8 б) 0 · (15 - 6) : 3 + (7 · 8 + 4) : 60 - 1 · 0', '</p> \n<p class="text">Составь программу действий и вычисли:</p> \n\n<p class="description-text"> \nа) 24 : 1 - (4 · 5 - 14) · 4 + 8 : 8 <br/>\nб) 0 · (15 - 6) : 3 + (7 · 8 + 4) : 60 - 1 · 0\n</p>', 'а) 24 : 1 - (4 · 5 - 14) · 4 + 8 : 8 = 1 4 · 5 = 20 20 - 14 = 6 6 · 4 = 24 24 : 1 = 24 8 : 8 = 1 24 - 24 = 0 0 + 1 = 1 б) 0 · (15 - 6) : 3 + (7 · 8 + 4) : 60 - 1 · 0 = 1 15 - 6 = 9 0 · 9 = 0 0 : 3 = 0 7 · 8 = 56 56 + 4 = 60 60 : 60 = 1 1 · 0 = 0 0 + 1 = 1 1 - 0 = 1', '<p>\nа) 24 : 1 - (4 · 5 - 14) · 4 + 8 : 8 = 1<br/>\n4 · 5 = 20<br/>\n20 - 14 = 6<br/> \n6 · 4 = 24<br/>\n24 : 1 = 24<br/>\n8 : 8 = 1 <br/>\n24 - 24 = 0<br/>\n0 + 1 = 1<br/><br/>\nб) 0 · (15 - 6) : 3 + (7 · 8 + 4) : 60 - 1 · 0 = 1<br/>\n15 - 6 = 9<br/>\n0 · 9 = 0<br/>\n0 : 3 = 0<br/>\n7 · 8 = 56<br/>\n56 + 4 = 60<br/>\n60 : 60 = 1<br/>\n1 · 0 = 0<br/>\n0 + 1 = 1<br/>\n1 - 0 = 1\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-71/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'b8f4aaa77dd36b32e5b82d4039aed9bc257f5dd2aa3ce17770cd7f9c1c4d1d50', '0,1,3,4,5,6,7,8,14,15', '["вычисли"]'::jsonb, 'составь программу действий и вычисли:а) 24:1-(4*5-14)*4+8:8 б) 0*(15-6):3+(7*8+4):60-1*0');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 71, '37', 4, 'Составь 4 равенства из чисел 12, 5, 60. Прочитай эти равенства разными способами и построй графическую модель.', '</p> \n<p class="text">Составь 4 равенства из чисел 12, 5, 60. Прочитай эти равенства разными способами и построй графическую модель.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica71-nomer37.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 71, номер 37, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 71, номер 37, год 2022."/>\n</div>\n</div>', '12 · 5 = 60, 60 : 12 = 5, 60 : 5 = 12, 5 · 12 = 60', '<p>\n12 · 5 = 60, 60 : 12 = 5, 60 : 5 = 12, 5 · 12 = 60\n</p>\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica71-nomer37-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 71, номер 37-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 71, номер 37-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-71/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica71-nomer37.jpg', 'peterson/3/part3/page71/task37_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica71-nomer37-1.jpg', 'peterson/3/part3/page71/task37_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5873c107c58e41408bfd53b0c2536dca822759f4832dde7d1c03e70f17a0157f', '4,5,12,60', '["раз"]'::jsonb, 'составь 4 равенства из чисел 12, 5, 60. прочитай эти равенства разными способами и построй графическую модель');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 71, '38', 5, 'Реши уравнения c комментированием и сделай проверку: а) x : 9 = 4056 б) 8 · x = 24016 в) 351900 : x = 5', '</p> \n<p class="text">Реши уравнения c комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) x : 9 = 4056  б) 8 · x = 24016  в) 351900 : x = 5\n</p>', 'а) x : 9 = 4056 Что бы найти делимое надо делитель умножить на частное х = 4056 · 9', '<p>\nа) x : 9 = 4056  <br/>\nЧто бы найти делимое надо делитель умножить на частное<br/>\nх = 4056 · 9\n\n</p>\n\n<div class="img-wrapper-460">\n<img width="130" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica71-nomer38.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 71, номер 38, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 71, номер 38, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-71/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica71-nomer38.jpg', 'peterson/3/part3/page71/task38_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'dc46683ed0f8edded2856a201bbf6a15a3628064b5abe5c98fd78631e7c51d61', '5,8,9,4056,24016,351900', '["реши"]'::jsonb, 'реши уравнения c комментированием и сделай проверку:а) x:9=4056 б) 8*x=24016 в) 351900:x=5');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 71, '39', 6, 'Как умножить и как разделить круглые числа? Вычисли: а) 86700 · 6          в) 34500 · 80 б) 200 · 709          г) 5010 · 3000 д) 42800 : 40        ж) 21063000 : 700 е) 260400 : 50      з) 50402700 : 900', '</p> \n<p class="text">Как умножить и как разделить круглые числа? Вычисли:</p> \n\n<p class="description-text"> \nа) 86700 · 6          в) 34500 · 80 <br/> \nб) 200 · 709          г) 5010 · 3000 <br/><br/> \n\nд) 42800 : 40        ж) 21063000 : 700<br/>\nе) 260400 : 50      з) 50402700 : 900\n</p>', 'а) 86700 · 6 = 520200 б) 200 · 709 = 141800 в) 34500 · 80 = 2770000 г) 5010 · 3000 = 15030000 д) 42800 : 40 = 1712000 е) 260400 : 50 = 13020000 ж) 21063000 : 700 = 14744100000 з) 50402700 : 900 = 45362430000', '<p>\nа) 86700 · 6 = 520200  <br/>			\nб) 200 · 709 = 141800	<br/>		\nв) 34500 · 80 = 2770000  <br/>\nг) 5010 · 3000 = 15030000<br/>\nд) 42800 : 40 = 1712000<br/>\nе) 260400 : 50 = 13020000<br/>\nж) 21063000 : 700 = 14744100000<br/>\nз) 50402700 : 900 = 45362430000\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-71/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3d41e4b46a55bcb9b610612d54cb973019ac1a6da038021eaeae334cd858221b', '6,40,50,80,200,700,709,900,3000,5010', '["раздели","вычисли","раз"]'::jsonb, 'как умножить и как разделить круглые числа? вычисли:а) 86700*6          в) 34500*80 б) 200*709          г) 5010*3000 д) 42800:40        ж) 21063000:700 е) 260400:50      з) 50402700:900');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 71, '40', 7, 'Вычисли устно. Сделай проверку, используя формулу деления с остатком: 56 : 9        83 : 5        35 : 17 47 : 6        92 : 8        70 : 12 52 : 15      81 : 23 93 : 14      64 : 49', '</p> \n<p class="text">Вычисли устно. Сделай проверку, используя формулу деления с остатком:</p> \n\n<p class="description-text"> \n56 : 9        83 : 5        35 : 17  <br/>	\n47 : 6        92 : 8        70 : 12 <br/> 	\n52 : 15      81 : 23<br/>\n93 : 14      64 : 49\n</p>', '56 : 9 = 6 + 2 Проверка: 56 = 9 · 6 + 2 47 : 6 = 7 + 5 Проверка: 47 = 6 · 7 + 5 83 : 5 = 16 + 3 Проверка: 83 = 5 · 16 + 3 92 : 8 = 11 + 4 Проверка: 92 = 8 · 11 + 4 35 : 17 = 2 + 1 Проверка: 35 = 17 · 2 + 1 70 : 12 = 5 + 10 Проверка: 70 = 12 · 5 + 10 52 : 15 = 3 + 7 Проверка: 52 = 15 · 3 + 7 93 : 14 = 6 + 9 Проверка: 93 = 14 · 6 + 9 81 : 23 = 3 + 12 Проверка: 81 = 23 · 3 + 12 64 : 49 = 1 + 15 Проверка: 64 = 49 · 1 + 15', '<p>\n56 : 9 = 6 + 2<br/>\n<b>Проверка:</b> 56 = 9 · 6 + 2<br/><br/>\n47 : 6 = 7 + 5<br/>\n<b>Проверка:</b> 47 = 6 · 7 + 5<br/><br/>\n83 : 5 =  16 + 3<br/>\n<b>Проверка:</b> 83 = 5 · 16 + 3<br/><br/>\n92 : 8 = 11 + 4<br/>\n<b>Проверка:</b> 92 = 8 · 11 + 4<br/><br/>\n35 : 17 = 2 + 1<br/>\n<b>Проверка:</b> 35 = 17 · 2 + 1<br/><br/>\n70 : 12 = 5 + 10<br/>\n<b>Проверка:</b> 70 = 12 · 5 + 10<br/><br/>\n52 : 15 = 3 + 7<br/>\n<b>Проверка:</b> 52 = 15 · 3 + 7<br/><br/>\n93 : 14 = 6 + 9<br/>\n<b>Проверка:</b> 93 = 14 · 6 + 9<br/><br/>\n81 : 23 = 3 + 12<br/>\n<b>Проверка:</b> 81 = 23 · 3 + 12<br/><br/>\n64 : 49 = 1 + 15<br/>\n<b>Проверка:</b> 64 = 49 · 1 + 15\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-71/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '9cc30b352457e0f7ff1da697e32c0508d9ba49df1f335020e0502631d5c8a368', '5,6,8,9,12,14,15,17,23,35', '["вычисли"]'::jsonb, 'вычисли устно. сделай проверку, используя формулу деления с остатком:56:9        83:5        35:17 47:6        92:8        70:12 52:15      81:23 93:14      64:49');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 71, '41', 8, 'Выполни деление с остатком и сделай проверку: а) 5108 : 7      в) 40153 : 5        д) 840260 : 80 б) 3275 : 3      г) 603240 : 9      е) 360450 : 60', '</p> \n<p class="text">Выполни деление с остатком и сделай проверку:</p> \n\n<p class="description-text"> \nа) 5108 : 7      в) 40153 : 5        д) 840260 : 80<br/>\nб) 3275 : 3      г) 603240 : 9      е) 360450 : 60\n</p>', 'а) 5108 : 7 = 729 + 5', '<p>\nа) 5108 : 7 = 729 + 5\n</p>\n\n<div class="img-wrapper-460">\n<img width="70" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica71-nomer41.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 71, номер 41, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 71, номер 41, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-71/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica71-nomer41.jpg', 'peterson/3/part3/page71/task41_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '15c135a5cd2483cc3074941692fd6112cf4425f3e59db280224d5a49320c8eb4', '3,5,7,9,60,80,3275,5108,40153,360450', NULL, 'выполни деление с остатком и сделай проверку:а) 5108:7      в) 40153:5        д) 840260:80 б) 3275:3      г) 603240:9      е) 360450:60');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 72, '42', 0, 'Сравни выражения, где буквы обозначают натуральные числа: m + 48     80 + m       36 : x     24 : x 60 - n     25 - n          b : 5     b : 3 k - 18     k - 53          (9 + c) · 4     9 + c · 4 a + a + a     2 · a         d · 6 - d     d · 5', '</p> \n<p class="text">Сравни выражения, где буквы обозначают натуральные числа: </p> \n\n<p class="description-text"> \nm + 48 <span class="okon">   </span> 80 + m       36 : x <span class="okon">   </span> 24 : x<br/> \n60 - n <span class="okon">   </span> 25 - n          b : 5 <span class="okon">   </span> b : 3<br/> \nk - 18 <span class="okon">   </span> k - 53          (9 + c) · 4 <span class="okon">   </span> 9 + c · 4<br/> \na + a + a <span class="okon">   </span> 2 · a         d · 6 - d <span class="okon">   </span> d · 5\n</p>', 'm + 48 ˂ 80 + m    36 : x ˃ 24 : x 60 - n ˃ 25 - n        b : 5 ˂ b : 3 k - 18 ˃ k - 53        (9 + c) · 4 ˂ 9 + c · 4 a + a + a ˃ 2 · a       d · 6 - d = d · 5', '<p>\nm + 48 ˂ 80 + m    36 : x ˃ 24 : x<br/> \n60 - n ˃ 25 - n        b : 5 ˂ b : 3<br/> \nk - 18 ˃ k - 53        (9 + c) · 4 ˂ 9 + c · 4<br/> \na + a + a ˃ 2 · a       d · 6 - d = d · 5\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-72/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5fe87030c00b38a31c3f13ca230faa8655c72deac0031d2cd2b285ce180d74c2', '2,3,4,5,6,9,18,24,25,36', '["сравни"]'::jsonb, 'сравни выражения, где буквы обозначают натуральные числа:m+48     80+m       36:x     24:x 60-n     25-n          b:5     b:3 k-18     k-53          (9+c)*4     9+c*4 a+a+a     2*a         d*6-d     d*5');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 72, '43', 1, 'Ворон живёт 60 лет, а овца – в 5 раз меньше ворона. Лошадь живёт на 4 года больше овцы, а хомяк – в 8 раз меньше лошади. Сколько лет живёт хомяк?', '</p> \n<p class="text">Ворон живёт 60 лет, а овца – в 5 раз меньше ворона. Лошадь живёт на 4 года больше овцы, а хомяк – в 8 раз меньше лошади. Сколько лет живёт хомяк?</p>', '60 : 5 + 4 - 8 = 12 + 4 - 8 = 8 - 8 = 0 Ответ: хомяк живёт меньше года.', '<p>\n60 : 5 + 4 - 8 = 12 + 4 - 8 = 8 - 8 = 0<br/>\n<b>Ответ:</b> хомяк живёт меньше года.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-72/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd57687febede7284e7b3033ebeb53d341096ceb9d5e28c3f177b16932fcf2a9b', '4,5,8,60', '["больше","меньше","раз"]'::jsonb, 'ворон живёт 60 лет, а овца-в 5 раз меньше ворона. лошадь живёт на 4 года больше овцы, а хомяк-в 8 раз меньше лошади. сколько лет живёт хомяк');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 72, '44', 2, 'Четыре зайчишки - братишки пошли в поле за морковками. Каждый из них принёс домой по 45 морковок. За ужином съели 36 морковок, а остальные разложили поровну в 3 пакета. Сколько морковок в каждом пакете?', '</p> \n<p class="text">Четыре зайчишки - братишки пошли в поле за морковками. Каждый из них принёс домой по 45 морковок. За ужином съели 36 морковок, а остальные разложили поровну в 3 пакета. Сколько морковок в каждом пакете?</p>', '(4 · 45 - 36) : 3 = (180 - 36) : 3 = 144 : 3 = 4 : 3 = 8 (морковок) Ответ: 8 морковок в каждом пакете.', '<p>\n(4 · 45 - 36) : 3 = (180 - 36) : 3 = 144 : 3 = 4 : 3 = 8 (морковок)<br/>\n<b>Ответ:</b> 8 морковок в каждом пакете.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-72/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f302404290e62df851d37bbc73aaaa9947f851d1b9b39877e32bc4d6289a509a', '3,36,45', '["раз"]'::jsonb, 'четыре зайчишки-братишки пошли в поле за морковками. каждый из них принёс домой по 45 морковок. за ужином съели 36 морковок, а остальные разложили поровну в 3 пакета. сколько морковок в каждом пакете');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 72, '45', 3, 'В саду у Хрюши росла яблонька. Осенью он собрал урожай – 50 яблок. По 2 яблока Хрюша подарил 5 белочкам и по 3 яблока дал 3 ёжикам. Сколько яблок у него ещё осталось?', '</p> \n<p class="text">В саду у Хрюши росла яблонька. Осенью он собрал урожай – 50 яблок. По 2 яблока Хрюша подарил 5 белочкам и по 3 яблока дал 3 ёжикам. Сколько яблок у него ещё осталось?</p>', '50 - (2 · 5 + 3 · 3) = 50 - (10 + 9) = 50 - 19 = 31 (яблок) Ответ: 31 яблок у него ещё осталось.', '<p>\n50 - (2 · 5 + 3 · 3) = 50 - (10 + 9) = 50 - 19 = 31 (яблок)<br/>\n<b>Ответ:</b> 31 яблок у него ещё осталось.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-72/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'a3b2108618e6c8b78769739b24b727df84a9f3644dbb38dd78261cbdc23de611', '2,3,5,50', NULL, 'в саду у хрюши росла яблонька. осенью он собрал урожай-50 яблок. по 2 яблока хрюша подарил 5 белочкам и по 3 яблока дал 3 ёжикам. сколько яблок у него ещё осталось');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 72, '46', 4, 'Бежала Мышка по полю и нашла 6 колосков по 40 зёрен в каждом. Чтобы испечь пирог, ей нужно 30 зёрен. Сколько пирогов сможет испечь Мышка из найденных колосков?', '</p> \n<p class="text">Бежала Мышка по полю и нашла 6 колосков по 40 зёрен в каждом. Чтобы испечь пирог, ей нужно 30 зёрен. Сколько пирогов сможет испечь Мышка из найденных колосков?</p>', '6 · 40 : 30 = 240 : 30 = 8 (пирогов) Ответ: 8 пирогов сможет испечь Мышка из найденных колосков.', '<p>\n6 · 40 : 30 = 240 : 30 = 8 (пирогов)<br/>\n<b>Ответ:</b> 8 пирогов сможет испечь Мышка из найденных колосков.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-72/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'f4d2a51a67a863ae267e5bbc313d1c3003ca83360e0bbb4236afacb869a4d6cc', '6,30,40', NULL, 'бежала мышка по полю и нашла 6 колосков по 40 зёрен в каждом. чтобы испечь пирог, ей нужно 30 зёрен. сколько пирогов сможет испечь мышка из найденных колосков');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 72, '47', 5, 'Белочка заготавливает грибы на зиму – каждый день одинаковое число грибов. За 5 дней она успела заготовить 40 грибов. Сколько грибов она сможет заготовить за неделю (7 дней)? За сколько дней она заготовит 200 грибов?', '</p> \n<p class="text">Белочка заготавливает грибы на зиму – каждый день одинаковое число грибов. За 5 дней она успела заготовить 40 грибов. Сколько грибов она сможет заготовить за неделю (7 дней)? За сколько дней она заготовит 200 грибов?</p>', '40 : 5 = 8 (грибов) – за 1 день 7 · 8 = 56 (грибов) – за неделю 200 : 8 = 25 (дней) Ответ: 56 грибов она сможет заготовить за неделю (7 дней). За 25 дней она заготовит 200 грибов.', '<p>\n40 : 5 = 8 (грибов) – за 1 день<br/>\n7 · 8 = 56 (грибов) – за неделю<br/>\n200 : 8 = 25 (дней)<br/>\n<b>Ответ:</b> 56 грибов она сможет заготовить за неделю (7 дней). За 25 дней она заготовит 200 грибов.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-72/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0752cb4ffd09753c2b87d481b303c9f9903cc5db23d6a36efbbc49b85c61bc79', '5,7,40,200', NULL, 'белочка заготавливает грибы на зиму-каждый день одинаковое число грибов. за 5 дней она успела заготовить 40 грибов. сколько грибов она сможет заготовить за неделю (7 дней)? за сколько дней она заготовит 200 грибов');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 72, '48', 6, 'Корней и Матвей подтянулись вместе 36 раз. Корней подтянулся на 14 раз меньше Матвея. Сколько раз подтянулся каждый из них?', '</p> \n<p class="text">Корней и Матвей подтянулись вместе 36 раз. Корней подтянулся на 14 раз меньше Матвея. Сколько раз подтянулся каждый из них?</p>', '(х - 14) + х = 36 х - 14 + х = 36 2х = 36 + 14 х = 50 : 2 х = 25 (раз) – Матвей 25 - 14 = 11 (раз) – Корней Ответ: Матвей – 25, Корней – 11 раз подтянулись.', '<p>\n(х - 14) + х = 36<br/>\nх - 14 + х = 36<br/>\n2х = 36 + 14<br/>\nх = 50 : 2<br/>\nх = 25 (раз) – Матвей<br/>\n25 - 14 = 11 (раз) – Корней<br/>\n<b>Ответ:</b> Матвей –  25, Корней – 11 раз подтянулись.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-72/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2359fb02a5397c34c54572c238f22db6fbbd0ff089a117fa185e91bc90ea2cba', '14,36', '["меньше","раз"]'::jsonb, 'корней и матвей подтянулись вместе 36 раз. корней подтянулся на 14 раз меньше матвея. сколько раз подтянулся каждый из них');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 72, '49', 7, 'Прошлым летом Пантелей собрал урожай абрикосов с трёх деревьев. С первого дерева он собрал 312 абрикосов, со второго – в 2 раза меньше, чем с первого, а с третьего – на 28 абрикосов больше, чем со второго. Из них 652 абрикоса съела коза Бориска, а остальные достались Пантелею. Сколько абрикосов ему досталось?', '</p> \n<p class="text">Прошлым летом Пантелей собрал урожай абрикосов с трёх деревьев. С первого дерева он собрал 312 абрикосов, со второго – в 2 раза меньше, чем с первого, а с третьего – на 28 абрикосов больше, чем со второго. Из них 652 абрикоса съела коза Бориска, а остальные достались Пантелею. Сколько абрикосов ему досталось?</p>', '312 + 312 : 2 + (312 : 2 + 28) - 652 = 312 + 156 + (156 + 28) - 652 = 468 + 184 - 652 = 652 - 652 = 0 (абрикосов) Ответ: ноль абрикосов ему досталось.', '<p>\n312 + 312 : 2 + (312 : 2 + 28) - 652 = 312 + 156 + (156 + 28) - 652 = 468 + 184 - 652 = 652 - 652 = 0 (абрикосов)<br/>\n<b>Ответ:</b> ноль абрикосов ему досталось.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-72/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3aa9ac9d0c60740b13a0c5dd02f069f70780536603113b2381c7415a0aa54691', '2,28,312,652', '["больше","меньше","раз","раза"]'::jsonb, 'прошлым летом пантелей собрал урожай абрикосов с трёх деревьев. с первого дерева он собрал 312 абрикосов, со второго-в 2 раза меньше, чем с первого, а с третьего-на 28 абрикосов больше, чем со второго. из них 652 абрикоса съела коза бориска, а остальные достались пантелею. сколько абрикосов ему досталось');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 73, '50', 0, 'Найди х: 1) x + m = n       3) x - c = d 2) a - x = b         4) k + x = p 5) x : a = c          7) x · m = r 6) t · x = k          8) b : x = d', '</p> \n<p class="text">Найди  х:</p> \n\n<p class="description-text"> \n1) x + m = n       3) x - c = d<br/>	\n2) a - x = b         4) k + x = p <br/><br/> 		\n\n5) x : a = c          7) x · m = r<br/>\n6) t · x = k          8) b : x = d\n\n</p>', '1) x + m = n х = n - m 2) a - x = b x = a - b 3) x - c = d x = d + c 4) k + x = p x = p - k 5) x : a = c x = a · c 6) t · x = k x = k : t 7) x · m = r x = r : m 8) b : x = d x = b : d', '<p>\n1) x + m = n <br/> \nх = n - m<br/>\n2) a - x = b<br/>  \nx = a - b<br/>\n3) x - c = d<br/>  \nx = d + c<br/>\n4) k + x = p<br/>  \nx = p - k<br/>\n5) x : a = c<br/>  \nx = a · c<br/>\n6) t · x = k<br/>  \nx = k : t<br/>\n7) x · m = r<br/>\nx = r : m<br/>\n8) b : x = d<br/>\nx = b : d \n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-73/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1d696e0a1ced82d5c6f89cb9346319add413c99aa6e0dcf95a284b2e3d50b64a', '1,2,3,4,5,6,7,8', '["найди"]'::jsonb, 'найди х:1) x+m=n       3) x-c=d 2) a-x=b         4) k+x=p 5) x:a=c          7) x*m=r 6) t*x=k          8) b:x=d');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 73, '51', 1, 'Текст из условия задания a + 3 · b          x : 2 - y (c + d) · (m - n) (8 · k) : (p + 4)', '</p> \n<p class="text">Текст из условия задания</p> \n\n<p class="description-text"> \na + 3 · b          x : 2 - y<br/> 	\n(c + d) · (m - n) 		<br/>\n(8 · k) : (p + 4)\n</p>', 'a + 3 · b сумма слагаемого а и произведения 3 на b, к числу а прибавим произведение чисел 3 и b x : 2 - y разность частного делимого х и делителя 2 с числом у из х разделенного на двое вычесть у (c + d) · (m - n) Произведение суммы чисел c и d с разностью чисел m и n Сумма чисел с и d умножить на разность уменьшаемого m и вычитаемого n (8 · k) : (p + 4) Частное произведения чисел 8 и k с суммой чисел р и 4 Произведение чисел 8 и k разделим на сумму чисел р и 4.', '<p>\na + 3 · b <br/>\nсумма слагаемого а и произведения 3 на b,<br/>\n к числу а прибавим произведение чисел 3 и b	<br/>\nx : 2 - y 	<br/>\nразность частного делимого х и делителя 2 с числом у<br/>\nиз х разделенного на двое вычесть у <br/>\n(c + d) · (m - n) <br/>\nПроизведение суммы чисел c и d с разностью чисел m и n	<br/>\nСумма чисел с и d умножить на разность уменьшаемого m и вычитаемого n	<br/>\n(8 · k) : (p + 4)<br/>\nЧастное произведения чисел 8 и k с суммой чисел р и 4<br/>\nПроизведение чисел 8 и k разделим на сумму чисел р и 4.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-73/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '02888ec6c9a7a912d75383f3286a40339b092d0afae6006c36e1250a7df02eee', '2,3,4,8', NULL, 'текст из условия задания a+3*b          x:2-y (c+d)*(m-n) (8*k):(p+4)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 73, '52', 2, 'Реши уравнения с комментированием и сделай проверку: а) 64 + 36 : (x · 3 - 15) = 70 б) 124 - 24 · (480 : x - 56) = 28', '</p> \n<p class="text">Реши уравнения с комментированием и сделай проверку:</p> \n\n<p class="description-text"> \nа) 64 + 36 : (x · 3 - 15) = 70  <br/>    \nб) 124 - 24 · (480 : x - 56) = 28\n</p>', 'а) 64 + 36 : (x · 3 - 15) = 70 Чтобы найти слагаемое 36 : (x · 3 - 15) надо из суммы вычесть известное слагаемое 36 : (x · 3 - 15) = 70 - 64 36 : (x · 3 - 15) = 6 Чтобы найти делитель (x · 3 - 15) надо делимое разделить на частное (x · 3 - 15) = 36 : 6 (x · 3 - 15) = 6 Чтобы найти уменьшаемое x · 3 надо вычитаемое сложить с разностью x · 3 = 15 + 6 x · 3 = 21 Чтобы найти множитель надо произведение разделить на известный множитель х = 21 : 3 х = 7 Проверка: 64 + 36 : (7 · 3 - 15) = 70 б) 124 - 24 · (480 : x - 56) = 28 Чтобы найти вычитаемое 24 · (480 : x - 56) надо из уменьшаемого вычесть разность 24 · (480 : x - 56) = 124 - 28 24 · (480 : x - 56) = 96 Чтобы найти множитель (480 : x - 56) надо произведение разделит на известный множитель (480 : x - 56) = 96 : 24 (480 : x - 56) = 4 Чтобы найти уменьшаемое 480 : x надо вычитаемое сложить с разностью 480 : x = 56 + 4 480 : x = 60 Чтобы найти делитель надо делимое разделить на частное х = 480 : 60 х = 8 Проверка: 124 - 24 · (480 : 8 - 56) = 28', '<p>\nа) 64 + 36 : (x · 3 - 15) = 70   <br/>   \nЧтобы найти слагаемое 36 : (x · 3 - 15) надо из суммы вычесть известное слагаемое<br/>\n36 : (x · 3 - 15) = 70 - 64<br/>\n36 : (x · 3 - 15) = 6<br/>\nЧтобы найти делитель (x · 3 - 15) надо делимое разделить на частное<br/>\n(x · 3 - 15) = 36 : 6 <br/>\n(x · 3 - 15) = 6<br/>\nЧтобы найти уменьшаемое x · 3 надо вычитаемое сложить с разностью<br/>\nx · 3 = 15 + 6<br/>\nx · 3 = 21<br/>\nЧтобы найти множитель надо произведение разделить на известный множитель<br/>\nх = 21 : 3<br/>\nх = 7<br/>\n<b>Проверка:</b> 64 + 36 : (7 · 3 - 15) = 70 <br/><br/>\nб) 124 - 24 · (480 : x - 56) = 28<br/>\nЧтобы найти вычитаемое 24 · (480 : x - 56) надо из уменьшаемого вычесть разность<br/>\n24 · (480 : x - 56) = 124 - 28<br/>\n24 · (480 : x - 56) = 96<br/>\nЧтобы найти множитель (480 : x - 56) надо произведение разделит на известный множитель<br/>\n(480 : x - 56) = 96 : 24 <br/>\n(480 : x - 56) = 4<br/>\nЧтобы найти уменьшаемое 480 : x надо вычитаемое сложить с разностью<br/>\n480 : x = 56 + 4<br/>\n480 : x = 60<br/>\nЧтобы найти делитель надо делимое разделить на частное<br/>\nх = 480 : 60<br/>\nх = 8<br/>\n<b>Проверка:</b> 124 - 24 · (480 : 8 - 56) = 28\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-73/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e72abcbfe3740396b497b5760ad32fae276a5e7ae7993bb526afe4fda1d661ec', '3,15,24,28,36,56,64,70,124,480', '["реши"]'::jsonb, 'реши уравнения с комментированием и сделай проверку:а) 64+36:(x*3-15)=70 б) 124-24*(480:x-56)=28');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 73, '53', 3, 'Ваня задумал число, увеличил его в 7 раз, вычел 9, разделил на 6, к результату прибавил 15, разделил на 3 и получил 8. Какое число задумал Ваня?', '</p> \n<p class="text">Ваня задумал число, увеличил его в 7 раз, вычел 9, разделил на 6, к результату прибавил 15, разделил на 3 и получил 8. Какое число задумал Ваня?</p>', '((х · 7 - 9) : 6 + 15) : 3 = 8 (х · 7 - 9) : 6 + 15 = 8 · 3 (х · 7 - 9) : 6 + 15 = 24 (х · 7 - 9) : 6 = 24 - 15 (х · 7 - 9) : 6 = 9 (х · 7 - 9) = 9 · 6 (х · 7 - 9) = 54 х · 7 = 54 + 9 х · 7 = 63 х = 63 : 7 х = 9 Ответ: число 9 задумал Ваня.', '<p>\n((х · 7 - 9) : 6 + 15) : 3 = 8<br/>\n(х · 7 - 9) : 6 + 15 = 8 · 3<br/>\n(х · 7 - 9) : 6 + 15 = 24<br/>\n(х · 7 - 9) : 6 = 24 - 15<br/>\n(х · 7 - 9) : 6 = 9<br/>\n(х · 7 - 9) = 9 · 6<br/>\n(х · 7 - 9) = 54<br/>\nх · 7 = 54 + 9<br/>\nх · 7 = 63<br/>\nх = 63 : 7<br/>\nх = 9<br/>\n<b>Ответ:</b> число 9 задумал Ваня.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-73/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '374b91aa46f85e042481392dcad8520fac6f7c694392d7316ade21842fa47265', '3,6,7,8,9,15', '["раздели","раз"]'::jsonb, 'ваня задумал число, увеличил его в 7 раз, вычел 9, разделил на 6, к результату прибавил 15, разделил на 3 и получил 8. какое число задумал ваня');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 73, '54', 4, 'Найди произведения: а) 35 · 18          в) 74 · 953 б) 279 · 42        г) 506 · 125 д) 817 · 304      ж) 123450 · 7800 е) 608 · 207      з) 69080 · 10500', '</p> \n<p class="text">Найди произведения:</p> \n\n<p class="description-text"> \nа) 35 · 18          в) 74 · 953<br/>  		\nб) 279 · 42        г) 506 · 125 <br/><br/> 	\n\nд) 817 · 304      ж) 123450 · 7800<br/>\nе) 608 · 207      з) 69080 · 10500\n\n</p>', 'а) 35 · 18 = 630', '<p>\nа) 35 · 18 = 630\n</p>\n\n<div class="img-wrapper-460">\n<img width="70" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica73-nomer54.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 73, номер 54, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 73, номер 54, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-73/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica73-nomer54.jpg', 'peterson/3/part3/page73/task54_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '501636027ba095adde186227990b40394dbe939decd7816c8dfeaec48ccc3959', '18,35,42,74,125,207,279,304,506,608', '["найди"]'::jsonb, 'найди произведения:а) 35*18          в) 74*953 б) 279*42        г) 506*125 д) 817*304      ж) 123450*7800 е) 608*207      з) 69080*10500');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 73, '55', 5, 'Найди значения выражений: а) (729 · 8 + 729 · 492) : 90 · (520800 : 400 - 498) б) 405 · (803 - 597) : 6 + 876000 : (3104 - 72 · 38 + 432)', '</p> \n<p class="text">Найди значения выражений:</p> \n\n<p class="description-text"> \nа) (729 · 8 + 729 · 492) : 90 · (520800 : 400 - 498)<br/>\nб) 405 · (803 - 597) : 6 + 876000 : (3104 - 72 · 38 + 432)\n</p>', 'а) (729 · 8 + 729 · 492) : 90 · (520800 : 400 - 498) = 729 · (8 + 492) : 90 · (1302 - 498) = 729 · 500 : 90 · 804 = 364500 : 90 · 804 = 4050 · 804 = 3256200', '<p>\nа) (729 · 8 + 729 · 492) : 90 · (520800 : 400 - 498) = 729 · (8 + 492) : 90 · (1302 - 498) = 729 · 500 : 90 · 804 = 364500 : 90 · 804 = 4050 · 804 = 3256200 \n</p>\n\n\n<div class="img-wrapper-460">\n<img width="250" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica73-nomer55.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 73, номер 55, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 73, номер 55, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-73/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica73-nomer55.jpg', 'peterson/3/part3/page73/task55_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'd8b024104e94855b1d8dbd9f93bdf8df320f3042f186090bc6f8a8da8460e652', '6,8,38,72,90,400,405,432,492,498', '["найди"]'::jsonb, 'найди значения выражений:а) (729*8+729*492):90*(520800:400-498) б) 405*(803-597):6+876000:(3104-72*38+432)');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 73, '56', 6, 'Составь выражения и найди их значения: а) Белоснежка приготовила m порций мороженого. Из них n порций она отдала своему другу Медвежонку, а остальные разделила поровну между 7 гномами. Сколько порций мороженого получил каждый гном? (m = 17, n = 3) б) Чтобы добраться до замка Принцессы, Кот в сапогах преодолел m км. Первые n км он ехал на повозке. Остальной путь он шёл пешком в течение недели, проходя каждый день поровну. Сколько километров проходил Кот в сапогах за один день? (m = 500, n = 150) Что общего и что различного в этих задачах? Придумай свою задачу про сказочных героев, имеющую такое же решение.', '</p> \n<p class="text">Составь выражения и найди их значения:<br/>\nа) Белоснежка приготовила m порций мороженого. Из них n порций она отдала своему другу Медвежонку, а остальные разделила поровну между 7 гномами. Сколько порций мороженого получил каждый гном? (m = 17, n = 3)<br/>\nб) Чтобы добраться до замка Принцессы, Кот в сапогах преодолел m км. Первые n км он ехал на повозке. Остальной путь он шёл пешком в течение недели, проходя каждый день поровну. Сколько километров проходил Кот в сапогах за один день? (m = 500, n = 150)<br/>\nЧто общего и что различного в этих задачах? Придумай свою задачу про сказочных героев, имеющую такое же решение.\n</p>', 'а) (m - n) : 7 (17 - 3) : 7 = 14 : 7 = 2 (порций) Ответ: 2 порций мороженого получил каждый гном. б) (m - n) : 7 (500 - 150) : 7 = 350 : 7 = 50 (км) Ответ: 50 километров проходил Кот в сапогах за один день. Общее выражение для поиска искомого и различное в этих задачах результат, часть данных. Своя задача про сказочных героев, имеющая такое же решение. Колобок преодолел m км. Первые n км его провезли на повозке. Остальной путь он катился сам в течение недели, преодолевая каждый день поровну. Сколько километров преодолел Колобок за один день? (m = 500, n = 150) (m - n) : 7 (500 - 150) : 7 = 350 : 7 = 50 (км) Ответ: 50 километров преодолел Колобок за один день.', '<p>\nа) (m - n) : 7 <br/>\n(17 - 3) : 7 = 14 : 7 = 2 (порций)<br/>\n<b>Ответ:</b> 2 порций мороженого получил каждый гном.<br/><br/>\nб) (m - n) : 7<br/>\n(500 - 150) : 7 = 350 : 7 = 50 (км)<br/>\n<b>Ответ:</b> 50 километров проходил Кот в сапогах за один день.<br/><br/>\nОбщее выражение для поиска искомого и различное в этих задачах результат, часть данных.<br/>\nСвоя задача про сказочных героев, имеющая такое же решение.<br/>\nКолобок преодолел m км. Первые n км его провезли на повозке. Остальной путь он катился сам в течение недели, преодолевая каждый день поровну. Сколько километров преодолел Колобок за один день? (m = 500, n = 150)<br/>\n(m - n) : 7<br/>\n(500 - 150) : 7 = 350 : 7 = 50 (км)<br/>\n<b>Ответ:</b> 50 километров преодолел Колобок за один день.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-73/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c7d628cf6840afa792b330845490a2563d56624b16b6ab0c8d289f136d17b9cb', '3,7,17,150,500', '["раздели","найди","раз"]'::jsonb, 'составь выражения и найди их значения:а) белоснежка приготовила m порций мороженого. из них n порций она отдала своему другу медвежонку, а остальные разделила поровну между 7 гномами. сколько порций мороженого получил каждый гном? (m=17, n=3) б) чтобы добраться до замка принцессы, кот в сапогах преодолел m км. первые n км он ехал на повозке. остальной путь он шёл пешком в течение недели, проходя каждый день поровну. сколько километров проходил кот в сапогах за один день? (m=500, n=150) что общего и что различного в этих задачах? придумай свою задачу про сказочных героев, имеющую такое же решение');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 74, '57', 0, 'Придумай задачи по таблицам: Что ты замечаешь? Придумай задачи с другими величинами, которые решаются так же.', '</p> \n<p class="text">Придумай задачи по таблицам:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica74-nomer57.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 74, номер 57, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 74, номер 57, год 2022."/>\n</div>\n</div>\n\n<p class="text">Что ты замечаешь? Придумай задачи с другими величинами, которые решаются так же.</p>', 'а) Первую часть Колобок преодолел в течении 2 ч со скоростью 60 км/ч, а вторую часть пути в течении 3 ч со скоростью 50 км/ч. Сколько составил км каждый путь и всего? 60 · 2 + 50 · 3 = 120 + 150 = 270 (км) Ответ: первая часть пути – 120 км, вторая часть пути – 150 км и 270 км - всего. б) В первые 2 ч рабочие имели производительность 60 шт./ч, в следующие 3 ч – 50 шт./ч. Сколько произведено в первые 2 ч и в следующие 3 ч? Сколько всего произведено? 60 · 2 + 50 · 3 = 120 + 150 = 270 (шт.) Ответ: первая 2 ч – 120 шт., вторые 3 ч – 150 шт. и 270 шт. - всего.', '<p>\nа) Первую часть Колобок преодолел в течении 2 ч со скоростью 60 км/ч, а вторую часть пути в течении 3 ч со скоростью 50 км/ч. Сколько составил км каждый путь и всего?<br/>\n60 · 2 + 50 · 3 = 120 + 150 = 270 (км)<br/>\n<b>Ответ:</b> первая часть пути – 120 км, вторая часть пути – 150 км и 270 км - всего. <br/><br/>  \nб) В первые 2 ч рабочие имели производительность 60 шт./ч, в следующие 3 ч – 50 шт./ч. Сколько произведено в первые 2 ч и в следующие 3 ч? Сколько всего произведено?<br/>\n 60 · 2 + 50 · 3 = 120 + 150 = 270 (шт.)<br/>\n <b>Ответ:</b> первая 2 ч – 120 шт., вторые 3 ч – 150 шт. и 270 шт. - всего.   \n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-74/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica74-nomer57.jpg', 'peterson/3/part3/page74/task57_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c4343ea31ef06dd24bc28bdbff88b8933be35ecc79ceab6a59041d3c844d1baf', NULL, NULL, 'придумай задачи по таблицам:что ты замечаешь? придумай задачи с другими величинами, которые решаются так же');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 74, '58', 1, 'Велосипедист проехал расстояние 32 км за 2 ч. За сколько времени он проедет 80 км, если его скорость не изменится?', '</p> \n<p class="text">Велосипедист проехал расстояние 32 км за 2 ч. За сколько времени он проедет 80 км, если его скорость не изменится?</p>', '80 : (32 : 2) = 80 : 16 = 5 (ч) Ответ: за 5 ч он проедет 80 км, если его скорость не изменится.', '<p>\n80 : (32 : 2) = 80 : 16 = 5 (ч)<br/>\n<b>Ответ:</b> за 5 ч он проедет 80 км, если его скорость не изменится.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-74/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '07edf423efdf06c6396afbf69bbaf8f7d795ef8eca1a63cb2dea559b2369d95a', '2,32,80', NULL, 'велосипедист проехал расстояние 32 км за 2 ч. за сколько времени он проедет 80 км, если его скорость не изменится');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 74, '59', 2, 'Катер проплыл расстояние 84 км за 3 ч, после чего ему осталось проплыть 140 км. За сколько времени он проплывёт оставшееся расстояние, если увеличит скорость на 7 км/ч?', '</p> \n<p class="text">Катер проплыл расстояние 84 км за 3 ч, после чего ему осталось проплыть 140 км. За сколько времени он проплывёт оставшееся расстояние, если увеличит скорость на 7 км/ч?</p>', '140 : (84 : 3 + 7) = 140 : (28 + 7) = 140 : 35 = 4 (ч) Ответ: за 4 ч он проплывёт оставшееся расстояние, если увеличит скорость на 7 км/ч.', '<p>\n140 : (84 : 3 + 7) = 140 : (28 + 7) = 140 : 35 = 4 (ч)<br/>\n<b>Ответ:</b> за 4 ч он проплывёт оставшееся расстояние, если увеличит скорость на 7 км/ч.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-74/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'e1306fbefb26d79cac4f0eb284e2b7b3ba14e17c7f98ecd16e6600ce37fe7df8', '3,7,84,140', NULL, 'катер проплыл расстояние 84 км за 3 ч, после чего ему осталось проплыть 140 км. за сколько времени он проплывёт оставшееся расстояние, если увеличит скорость на 7 км/ч');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 74, '60', 3, 'Мастер должен был изготовить 90 деталей за 6 ч. Однако он успевал сделать в час на 3 детали больше, чем предполагал. На сколько часов быстрее он сделал эту работу?', '</p> \n<p class="text">Мастер должен был изготовить 90 деталей за 6 ч. Однако он успевал сделать в час на 3 детали больше, чем предполагал. На сколько часов быстрее он сделал эту работу?</p>', '6 - 90 : (90 : 6 + 3) = 6 - 90 : (15 + 3) = 6 - 90 : 18 = 6 - 5 = 1 (ч) Ответ: на 1 ч быстрее он сделал эту работу.', '<p>\n6 - 90 : (90 : 6 + 3) = 6 - 90 : (15 + 3) = 6 - 90 : 18 = 6 - 5 = 1 (ч)<br/>\n<b>Ответ:</b> на 1 ч быстрее он сделал эту работу.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-74/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '5a711a1ecf5d53b7be21388278960b946723393f3a0d80520ecec79b56cefa48', '3,6,90', '["больше"]'::jsonb, 'мастер должен был изготовить 90 деталей за 6 ч. однако он успевал сделать в час на 3 детали больше, чем предполагал. на сколько часов быстрее он сделал эту работу');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 74, '61', 4, 'Лида и Оля купили тесьму на 48 р. каждая: Лида – по цене 8 р., а Оля – 12 р. за метр. Кто из них купил больше тесьмы и на сколько?', '</p> \n<p class="text">Лида и Оля купили тесьму на 48 р. каждая: Лида – по цене 8 р., а Оля – 12 р. за метр. Кто из них купил больше тесьмы и на сколько?</p>', '48 : 8 - 48 : 12 = 6 - 4 = 2 (м) Ответ: Лида купила больше тесьмы на 2 м больше.', '<p>\n48 : 8 - 48 : 12 = 6 - 4 = 2 (м)<br/>\n<b>Ответ:</b> Лида купила больше тесьмы на 2 м больше. \n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-74/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2f516e04e775915d0e8f6386ab8335c9e7930b06b150f7817ecb4f8141fe8d3a', '8,12,48', '["больше"]'::jsonb, 'лида и оля купили тесьму на 48 р. каждая:лида-по цене 8 р., а оля-12 р. за метр. кто из них купил больше тесьмы и на сколько');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 74, '62', 5, 'В летнем лагере «Орлёнок» отдыхало на 120 детей больше, чем в лагере «Следопыт». По окончании смены для отправки детей в город лагерю «Орлёнок» потребовалось 19 автобусов, а лагерю «Следопыт» – 14 таких же автобусов. Сколько детей отдыхало в этих лагерях, если в каждом автобусе ехало одинаковое количество детей?', '</p> \n<p class="text">В летнем лагере «Орлёнок» отдыхало на 120 детей больше, чем в лагере «Следопыт». По окончании смены для отправки детей в город лагерю «Орлёнок» потребовалось 19 автобусов, а лагерю «Следопыт» – 14 таких же автобусов. Сколько детей отдыхало в этих лагерях, если в каждом автобусе ехало одинаковое количество детей?</p>', '120 : (19 - 14) = 120 : 5 = 24 (детей) – 1 автобус 19 · 24 = 456 (детей) – «Орлёнок»', '<p>\n120 : (19 - 14) = 120 : 5 = 24 (детей) – 1 автобус<br/>\n19 · 24 = 456 (детей) – «Орлёнок»\n</p>\n\n<div class="img-wrapper-460">\n<img width="80" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica74-nomer62.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 74, номер 62, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 74, номер 62, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-74/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica74-nomer62.jpg', 'peterson/3/part3/page74/task62_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8c25a9c068c54a5ca2e884b6e6393dd248014f36a533cf4e705fd01194a379df', '14,19,120', '["больше"]'::jsonb, 'в летнем лагере "орлёнок" отдыхало на 120 детей больше, чем в лагере "следопыт". по окончании смены для отправки детей в город лагерю "орлёнок" потребовалось 19 автобусов, а лагерю "следопыт"-14 таких же автобусов. сколько детей отдыхало в этих лагерях, если в каждом автобусе ехало одинаковое количество детей');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 74, '63', 6, 'Автомобиль проехал с одинаковой скоростью в первый день 960 км, а во второй – 720 км. В первый день он был в пути на 3 ч больше, чем во второй день. Какое расстояние он проедет за 7 ч, двигаясь с той же скоростью?', '</p> \n<p class="text">Автомобиль проехал с одинаковой скоростью в первый день 960 км, а во второй – 720 км. В первый день он был в пути на 3 ч больше, чем во второй день. Какое расстояние он проедет за 7 ч, двигаясь с той же скоростью?</p>', '(960 - 720) : 3 · 7 = 240 : 3 · 7 = 80 · 7 = 560 (км) Ответ: 560 км он проедет за 7 ч, двигаясь с той же скоростью.', '<p>\n(960 - 720) : 3 · 7 = 240 : 3 · 7 = 80 · 7 = 560 (км)<br/>\n<b>Ответ:</b> 560 км он проедет за 7 ч, двигаясь с той же скоростью.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-74/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'c544c076209bcc8cd5c40511c54880773604dee34b0b440f8cfb45bf5b43b196', '3,7,720,960', '["больше"]'::jsonb, 'автомобиль проехал с одинаковой скоростью в первый день 960 км, а во второй-720 км. в первый день он был в пути на 3 ч больше, чем во второй день. какое расстояние он проедет за 7 ч, двигаясь с той же скоростью');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 74, '64', 7, 'Реши задачи и сравни их решения. Как называют такие задачи? а) Для двух классов купили 8 одинаковых пачек учебников. Один класс получил 45 учебников, а другой – 75. Сколько пачек учебников получил каждый класс? б) Для двух классов купили 120 учебников в одинаковых пачках. Один класс получил 3 пачки, а другой – 5 пачек. Сколько учебников получил каждый класс?', '</p> \n<p class="text">Реши задачи и сравни их решения. Как называют такие задачи?<br/>\nа) Для двух классов купили 8 одинаковых пачек учебников. Один класс получил 45 учебников, а другой –  75. Сколько пачек учебников получил каждый класс?<br/>\nб) Для двух классов купили 120 учебников в одинаковых пачках. Один класс получил 3 пачки, а другой – 5 пачек. Сколько учебников получил каждый класс?\n</p>', 'а) (75 + 45) : 8 = 120 : 8 = 15 (уч./п.) 45 : 15 = 3 (п.) – один класс 75 : 15 = 5 (п.) – другой класс Ответ: 3 в один класс и 5 в другой класс пачек учебников получил каждый класс. б) 120 : (3 + 5) = 120 : 8 = 15 (уч./п.) 15 · 3 = 45 (уч.) – один класс 15 ·5 = 75 (уч.) – другой класс Ответ: 45 один класс и 75 другой класс учебников получили.', '<p>\nа) (75 + 45) : 8 = 120 : 8 = 15 (уч./п.)<br/>\n45 : 15 = 3 (п.) – один класс<br/>\n75 : 15 = 5 (п.) – другой класс<br/>\n<b>Ответ:</b> 3 в один класс и 5 в другой класс пачек учебников получил каждый класс.<br/><br/>\nб) 120 : (3 + 5) = 120 : 8 = 15 (уч./п.)<br/>\n15 · 3 = 45 (уч.) – один класс<br/>\n15 ·5 = 75 (уч.) – другой класс<br/>\n<b>Ответ:</b> 45 один класс и 75 другой класс учебников получили.\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-74/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0f6a6153e552ed4fdad4ff15ec85e94b67afac1c7d6361688a98b4310c5b7c75', '3,5,8,45,75,120', '["реши","сравни"]'::jsonb, 'реши задачи и сравни их решения. как называют такие задачи? а) для двух классов купили 8 одинаковых пачек учебников. один класс получил 45 учебников, а другой-75. сколько пачек учебников получил каждый класс? б) для двух классов купили 120 учебников в одинаковых пачках. один класс получил 3 пачки, а другой-5 пачек. сколько учебников получил каждый класс');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 75, '65', 0, 'Вадим купил для себя 18 одинаковых тетрадей, а для соседа – 12 таких же тетрадей. За всю покупку он заплатил 450 р. Сосед принёс ему купюру в 500 р. Сколько сдачи Вадим должен ему вернуть?', '</p> \n<p class="text">Вадим купил для себя 18 одинаковых тетрадей, а для соседа – 12 таких же тетрадей. За всю покупку он заплатил 450 р. Сосед принёс ему купюру в 500 р. Сколько сдачи Вадим должен ему вернуть?</p>', '500 - 450 = 50 (р.) Ответ: 50 р. сдачи Вадим должен ему вернуть.', '<p>\n500 - 450 = 50 (р.)<br/>\n<b>Ответ:</b> 50 р. сдачи Вадим должен ему вернуть.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-75/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '8c6809cda48e1f4262d4bb1567f239995fb11c20e6aaa6beb0e18be742599caf', '12,18,450,500', NULL, 'вадим купил для себя 18 одинаковых тетрадей, а для соседа-12 таких же тетрадей. за всю покупку он заплатил 450 р. сосед принёс ему купюру в 500 р. сколько сдачи вадим должен ему вернуть');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 75, '66', 1, 'Первый маляр за 3 ч покрасил потолок в комнате площадью 27 м 2 . Второй маляр, выполняя такую же работу, потратил на 2 ч больше времени. Но площадь его комнаты была на 13 м 2 больше, чем у первого. У кого из них производительность больше и на сколько?', '</p> \n<p class="text">Первый маляр за 3 ч покрасил потолок в комнате площадью 27 м<sup>2</sup>. Второй маляр, выполняя такую же работу, потратил на 2 ч больше времени. Но площадь его комнаты была на 13 м<sup>2</sup> больше, чем у первого. У кого из них производительность больше и на сколько?</p>', '27 : 3 - (27 + 13) : (3 + 2) = 9 - 40 : 5 = 9 - 8 = 1 (м 2 /ч) Ответ: у первого из них производительность больше на 1 м 2 /ч.', '<p>\n27 : 3 - (27 + 13) : (3 + 2) = 9 - 40 : 5 = 9 - 8 = 1 (м<sup>2</sup>/ч)<br/>\n<b>Ответ:</b> у первого из них производительность больше на 1 м<sup>2</sup>/ч.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-75/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '3cee04f494ceb1e30ae1809b84593be02fb0e0cdb8c939f87f522f215b672ac2', '2,3,13,27', '["площадь","больше"]'::jsonb, 'первый маляр за 3 ч покрасил потолок в комнате площадью 27 м 2 . второй маляр, выполняя такую же работу, потратил на 2 ч больше времени. но площадь его комнаты была на 13 м 2 больше, чем у первого. у кого из них производительность больше и на сколько');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 75, '67', 2, 'Саше надо отметить точку М, нарисовать луч АК, отрезок ВС и прямую EF. На рисунке показан его чертёж. Какие ошибки он допустил? Нарисуй в тетради указанные фигуры правильно.', '</p> \n<p class="text">Саше надо отметить точку М, нарисовать луч АК, отрезок ВС и прямую EF. На рисунке показан его чертёж. Какие ошибки он допустил? Нарисуй в тетради указанные фигуры правильно.</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer67.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 75, номер 67, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 75, номер 67, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer67-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 75, номер 67-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 75, номер 67-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-75/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer67.jpg', 'peterson/3/part3/page75/task67_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer67-1.jpg', 'peterson/3/part3/page75/task67_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2748c0a23e3c18dde53f5755eae1cef3bf40c4045c8f2055573ae3b363b6abb6', NULL, NULL, 'саше надо отметить точку м, нарисовать луч ак, отрезок вс и прямую ef. на рисунке показан его чертёж. какие ошибки он допустил? нарисуй в тетради указанные фигуры правильно');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 75, '68', 3, 'Построй: а) прямую АМ; б) отрезок АМ; в) луч АМ; г) луч МА.', '</p> \n<p class="text">Построй: а) прямую АМ; б) отрезок АМ; в) луч АМ; г) луч МА.</p>', '', '<div class="img-wrapper-460">\n<img width="230" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer68.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 75, номер 68, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 75, номер 68, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-75/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer68.jpg', 'peterson/3/part3/page75/task68_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '419e289e2231bef08ae4aba3230efd80661fa229a475cd3dc6910840e1bf4710', NULL, NULL, 'построй:а) прямую ам; б) отрезок ам; в) луч ам; г) луч ма');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 75, '69', 4, 'а) Отметь две точки А и В, проведи через них прямую. Начерти луч ОМ, пересекающий прямую АВ, и луч КС, её не пересекающий. б) Отметь точки М и D и проведи луч DM. Начерти прямую EK, которая пересекает луч DM, и прямую АС, которая его не пересекает.', '</p> \n<p class="text">а) Отметь две точки А и В, проведи через них прямую. Начерти луч ОМ, пересекающий прямую АВ, и луч КС, её не пересекающий.<br/>\nб) Отметь точки М и D и проведи луч DM. Начерти прямую EK, которая пересекает луч DM, и прямую АС, которая его не пересекает.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer69.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 75, номер 69, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 75, номер 69, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-460">\n<img width="250" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer69-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 75, номер 69-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 75, номер 69-1, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Отметь две точки А и В, проведи через них прямую. Начерти луч ОМ, пересекающий прямую АВ, и луч КС, её не пересекающий.","solution":""},{"letter":"б","condition":"Отметь точки М и D и проведи луч DM. Начерти прямую EK, которая пересекает луч DM, и прямую АС, которая его не пересекает.","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-75/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer69.jpg', 'peterson/3/part3/page75/task69_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer69-1.jpg', 'peterson/3/part3/page75/task69_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '7b1a10c7869286e8cc1777324c3ee48975a98ddc149714948a5652dad209b49d', NULL, NULL, 'а) отметь две точки а и в, проведи через них прямую. начерти луч ом, пересекающий прямую ав, и луч кс, её не пересекающий. б) отметь точки м и d и проведи луч dm. начерти прямую ek, которая пересекает луч dm, и прямую ас, которая его не пересекает');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 75, '70', 5, 'Построй отрезок AB = 5 см 4 мм и отметь на нём точки C и D так, чтобы точка C лежала между точками B и D. Измерь отрезок BC.', '</p> \n<p class="text">Построй отрезок AB = 5 см 4 мм и отметь на нём точки C и D так, чтобы точка C лежала между точками B и D. Измерь отрезок BC.</p>', '', '<div class="img-wrapper-460">\n<img width="350" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer70.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 75, номер 70, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 75, номер 70, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-75/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer70.jpg', 'peterson/3/part3/page75/task70_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '1e6a1790b10db31bf8ce403c5b8c9cc36b2c208ba8c67e2c9630676b9f58e15c', '4,5', NULL, 'построй отрезок ab=5 см 4 мм и отметь на нём точки c и d так, чтобы точка c лежала между точками b и d. измерь отрезок bc');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 75, '71', 6, 'Измерь с помощью линейки стороны многоугольника и найди его периметр. Сколько у него острых углов, прямых, тупых?', '</p> \n<p class="text">Измерь с помощью линейки стороны многоугольника и найди его периметр. Сколько у него острых углов, прямых, тупых?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer71.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 75, номер 71, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 75, номер 71, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer71-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 75, номер 71-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 75, номер 71-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-75/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer71.jpg', 'peterson/3/part3/page75/task71_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer71-1.jpg', 'peterson/3/part3/page75/task71_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '6ed7235aba98ee363c933e96bf43386767d62587be387cb1eb5b9d565b905e30', NULL, '["найди","периметр"]'::jsonb, 'измерь с помощью линейки стороны многоугольника и найди его периметр. сколько у него острых углов, прямых, тупых');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 75, '72', 7, 'Найди в окружающей обстановке предметы, которые могут служить моделями отрезков. Рассмотри с помощью этих моделей возможные случаи взаимного расположения двух отрезков. Опиши их словами и изобрази на чертеже.', '</p> \n<p class="text">Найди в окружающей обстановке предметы, которые могут служить моделями отрезков. Рассмотри с помощью этих моделей возможные случаи взаимного расположения двух отрезков. Опиши их словами и изобрази на чертеже.</p>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer72.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 75, номер 72, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 75, номер 72, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-75/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica75-nomer72.jpg', 'peterson/3/part3/page75/task72_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'dd4c5f05486c6d21e59be6877549100f3fdd9a1c0ba281b0d2c6aaee3890c34d', NULL, '["найди","раз"]'::jsonb, 'найди в окружающей обстановке предметы, которые могут служить моделями отрезков. рассмотри с помощью этих моделей возможные случаи взаимного расположения двух отрезков. опиши их словами и изобрази на чертеже');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 76, '73', 0, 'а) Построй треугольник АВС. Построй треугольник, симметричный треугольнику АВС относительно стороны ВС. Перенеси полученный треугольник вправо на 8 клеточек. Опиши обратное преобразование. б) Построй квадрат АВСD со стороной 3 см. Построй квадрат, симметричный ему относительно стороны СD.', '</p> \n<p class="text">а) Построй треугольник АВС. Построй треугольник, симметричный треугольнику АВС относительно стороны ВС. Перенеси полученный треугольник вправо на 8 клеточек. Опиши обратное преобразование.<br/>\nб) Построй квадрат АВСD со стороной 3 см. Построй квадрат, симметричный ему относительно стороны СD.\n</p> \n\n<div class="description-text"> \n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer73.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 76, номер 73, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 76, номер 73, год 2022."/>\n</div>\n</div>', 'а)', '<p>\nа) \n</p>\n\n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer73-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 76, номер 73-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 76, номер 73-1, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Построй треугольник АВС. Построй треугольник, симметричный треугольнику АВС относительно стороны ВС. Перенеси полученный треугольник вправо на 8 клеточек. Опиши обратное преобразование.","solution":""},{"letter":"б","condition":"Построй квадрат АВСD со стороной 3 см. Построй квадрат, симметричный ему относительно стороны СD.","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-76/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer73.jpg', 'peterson/3/part3/page76/task73_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer73-1.jpg', 'peterson/3/part3/page76/task73_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '20ed7bf33ff70d87a27a60e223c2cef985a818fac2f066631365eb2000da92aa', '3,8', '["раз"]'::jsonb, 'а) построй треугольник авс. построй треугольник, симметричный треугольнику авс относительно стороны вс. перенеси полученный треугольник вправо на 8 клеточек. опиши обратное преобразование. б) построй квадрат авсd со стороной 3 см. построй квадрат, симметричный ему относительно стороны сd');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 76, '74', 1, 'а) Найди симметричные фигуры и укажи оси симметрии. Какими способами можно проверить правильность ответа? б) Сколько осей симметрии имеют прямоугольник, квадрат, круг? Построй их и укажи оси симметрии.', '</p> \n<p class="text">а) Найди симметричные фигуры и укажи оси симметрии. Какими способами можно проверить правильность ответа?</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer74.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 76, номер 74, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 76, номер 74, год 2022."/>\n</div>\n</div>\n\n\n<p class="text">б) Сколько осей симметрии имеют прямоугольник, квадрат, круг? Построй их и укажи оси симметрии.</p>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer74-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 76, номер 74-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 76, номер 74-1, год 2022."/>', '', '', TRUE, '[{"letter":"а","condition":"Найди симметричные фигуры и укажи оси симметрии. Какими способами можно проверить правильность ответа?","solution":""},{"letter":"б","condition":"Сколько осей симметрии имеют прямоугольник, квадрат, круг? Построй их и укажи оси симметрии.","solution":""}]'::jsonb, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-76/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer74.jpg', 'peterson/3/part3/page76/task74_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer74-1.jpg', 'peterson/3/part3/page76/task74_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'fa6d4007388c4cf4f749919587335191549a4f570affb0ed270314f10d82430d', NULL, '["найди"]'::jsonb, 'а) найди симметричные фигуры и укажи оси симметрии. какими способами можно проверить правильность ответа? б) сколько осей симметрии имеют прямоугольник, квадрат, круг? построй их и укажи оси симметрии');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 76, '75', 2, 'Выполни действия: а) (3 мин 48 с + 16 мин 36 с - 6 мин 54 с) · 120 б) (4 сут. 6 ч 15 мин - 18 ч 29 мин + 5 сут. 12 ч 14 мин) : 9', '</p> \n<p class="text">Выполни действия:</p> \n\n<p class="description-text"> \nа) (3 мин 48 с + 16 мин 36 с - 6 мин 54 с) · 120<br/>\nб) (4 сут. 6 ч 15 мин - 18 ч 29 мин + 5 сут. 12 ч 14 мин) : 9\n</p>', 'а) (3 мин 48 с + 16 мин 36 с - 6 мин 54 с) · 120 = 13 мин 30 с · 120 = 780 с · 120 = 93600 с = 1560 мин = 26 ч = 1 сут. 2 ч', '<p>\nа) (3 мин 48 с + 16 мин 36 с - 6 мин 54 с) · 120 = 13 мин 30 с · 120 = 780 с · 120 = 93600 с = 1560 мин = 26 ч = 1 сут. 2 ч\n</p>\n\n<div class="img-wrapper-460">\n<img width="220" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer75.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 76, номер 75, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 76, номер 75, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-76/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer75.jpg', 'peterson/3/part3/page76/task75_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '556c2cf8254b1097b8a07cfa4e30ed3daa5a09bc873dd580bd5949329c5719dd', '3,4,5,6,9,12,14,15,16,18', NULL, 'выполни действия:а) (3 мин 48 с+16 мин 36 с-6 мин 54 с)*120 б) (4 сут. 6 ч 15 мин-18 ч 29 мин+5 сут. 12 ч 14 мин):9');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 76, '76', 3, 'По таблице построй формулу зависимости y от x:', '</p> \n<p class="text">По таблице построй формулу зависимости y от x:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer76.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 76, номер 76, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 76, номер 76, год 2022."/>\n</div>\n</div>', 'а) у = х + 9 б) у = 9 · х', '<p>\nа) у = х + 9 <br/>\nб) у = 9 · х\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-76/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica76-nomer76.jpg', 'peterson/3/part3/page76/task76_condition_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '2c16ea651a3facb746bb0e38d62702b98d757d8ddeef15ae24774d573c78c3ab', NULL, NULL, 'по таблице построй формулу зависимости y от x');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 76, '77', 4, 'Подбери корни уравнений и сделай проверку: а) х · х + 4 = 29 б) (х - 2) · (х + 5) = 0', '</p> \n<p class="text">Подбери корни уравнений и сделай проверку: </p> \n\n<p class="description-text"> \nа) х · х + 4 = 29<br/>        б) (х - 2) · (х + 5) = 0\n</p>', 'а) х · х + 4 = 29 х · х = 29 - 4 х · х = 25 х 2 = 5 2 х = 5 Проверка: 5 · 5 + 4 = 29 б) (х - 2) · (х + 5) = 0 х - 2 = 0 х = 2 х + 5 = 0 х = -5 Проверка: (2 - 2) · (5 - 5) = 0', '<p>\nа) х · х + 4 = 29<br/>\nх · х = 29 - 4<br/>\nх · х = 25<br/>\nх<sup>2</sup> = 5<sup>2</sup><br/>\nх = 5    <br/>\n<b>Проверка:</b> 5 · 5 + 4 = 29<br/><br/>\nб) (х - 2) · (х + 5) = 0<br/>\nх - 2 = 0<br/>\nх = 2<br/>\nх + 5 = 0<br/>\nх = -5<br/>\n<b>Проверка:</b> (2 - 2) · (5 - 5) = 0\n\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-76/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '96f94ddbf1d8da9181513887eec15ab65710132c0395ac645f720e6bfebb5b96', '0,2,4,5,29', NULL, 'подбери корни уравнений и сделай проверку:а) х*х+4=29 б) (х-2)*(х+5)=0');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 76, '78', 5, 'Как называется множество: а) людей, обслуживающих самолёт в полёте; б) фруктовых деревьев на пришкольном участке; в) машин, движущихся по дороге; г) верблюдов, идущих друг за другом по пустыне?', '</p> \n<p class="text">Как называется множество:<br/>\nа) людей, обслуживающих самолёт в полёте;<br/>\nб) фруктовых деревьев на пришкольном участке;<br/>\nв) машин, движущихся по дороге;<br/>\nг) верблюдов, идущих друг за другом по пустыне?\n</p>', 'а) экипаж; б) сад; в) автомобильный поток; г) караван.', '<p>\nа) экипаж;<br/>\nб) сад;<br/>\nв) автомобильный поток;<br/>\nг) караван.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-76/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '44cdd1cacdcc178bc85263eebce44d747f148dafcf985a43b4b0803cd4af18eb', NULL, NULL, 'как называется множество:а) людей, обслуживающих самолёт в полёте; б) фруктовых деревьев на пришкольном участке; в) машин, движущихся по дороге; г) верблюдов, идущих друг за другом по пустыне');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 76, '79', 6, 'К – множество планет Солнечной системы. Принадлежит ли этому множеству Марс, Земля, Луна, Полярная звезда?', '</p> \n<p class="text">К – множество планет Солнечной системы. Принадлежит ли этому множеству Марс, Земля, Луна, Полярная звезда?</p>', 'Марс принадлежит к множеству К, Земля принадлежит к множеству К, Луна не принадлежит к множеству К, Полярная звезда не принадлежит к множеству К.', '<p>\nМарс принадлежит к множеству К, Земля принадлежит к множеству К, Луна не принадлежит к множеству К, Полярная звезда не принадлежит к множеству К.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-76/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, 'ff1435b3804b9fc26df8b53551d539586ab28622ea15c327f4b0d3364836da99', NULL, NULL, 'к-множество планет солнечной системы. принадлежит ли этому множеству марс, земля, луна, полярная звезда');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 77, '80', 0, 'А – множество трёхзначных чисел, В – множество чисел, оканчивающихся цифрой 2. Принадлежат ли этим множествам числа: 724, 42, 531, 1022, 738, 63? Сделай записи, используя знаки ∈ и ∉ .', '</p> \n<p class="text">\nА – множество трёхзначных чисел, В – множество чисел, оканчивающихся цифрой 2. Принадлежат ли этим множествам числа: 724, 42, 531, 1022, 738, 63? Сделай записи, используя знаки ∈ и ∉ .\n</p>', '724 ∈ А и 724 ∉ В, 42 ∈ В и 42 ∉ А, 531 ∈ А и 531 ∉ В, 1022 ∈ В и 1022 ∉ А, 738 ∈ А и 738 ∉ В, 63 ∉ А и 63 ∉ В.', '<p>\n724 ∈ А и 724 ∉ В, 42 ∈ В и 42 ∉ А, 531 ∈ А и 531 ∉ В, 1022 ∈ В и 1022 ∉ А, \n738 ∈ А и 738 ∉ В, 63 ∉ А и 63 ∉ В.\n</p>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-77/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '87318a4cefeffca824b3384166f218ee2f717022f5a2937af890fc9a2b18c1c7', '2,42,63,531,724,738,1022', NULL, 'а-множество трёхзначных чисел, в-множество чисел, оканчивающихся цифрой 2. принадлежат ли этим множествам числа:724, 42, 531, 1022, 738, 63? сделай записи, используя знаки ∈ и ∉');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 77, '81', 1, 'По диаграмме Эйлера–Венна определи, из каких элементов состоят множества А и В. Составь множества А ⋂ В и А ⋃ В. Сделай записи, используя знаки ⊂ и ⊄ :{4, ∆} … А {∆, 2} … В {∆} … А ⋂ В', '</p> \n<p class="text">По диаграмме Эйлера–Венна определи, из каких элементов состоят множества А и В.<br/>\nСоставь множества А ⋂ В и А ⋃ В.<br/>\nСделай записи, используя знаки ⊂ и ⊄ :{4, ∆} … А      {∆, 2} … В      {∆} … А ⋂ В\n</p> \n\n<div class="img-wrapper-460">\n<img width="200" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica77-nomer81.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 77, номер 81, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 77, номер 81, год 2022."/>\n</div>', 'А = {Т, 4, ∆} B = {∆, 3, n} А ⋂ В = {∆} А ⋃ В = {∆, n, T, 4, 3} {4, ∆} ⊂ А    {∆, 2} ⊄ В    {∆} ⊂ А ⋂ В A = {м, о, р, е}, D = {д, о, м}, E = {д, ы, м}, A ⋂ D = {о, м}, D ⋂ м E = {д, м}, (A ⋂ D) ⋂ E = {м}, A ⋂ (D ⋂ E) = {м}. M ⋃ K = {1; 3; 5; 7; 9; 10}, K ⋃ T = {3; 5; 6; 9; 10}, (M ⋃ K) ⋃ T = {1; 3; 5; 6; 7; 9; 10}, M ⋃ (K ⋃ T) = {1; 3; 5; 6; 7; 9; 10}. а) красный и синий, синий и жёлтый, красный и жёлтый – 3 способа б) красный и синий, синий и жёлтый, красный и жёлтый, красный и красный, синий и синий, жёлтый и жёлтый – 6 способов. Так как она выбрала только 2 стихотворения М. Ю. Лермонтова, 2 стихотворения А. Блока, а необходимы стихотворения разных авторов, то только 2 программы своего выступления сможет составить Аня из этих стихов, если порядок их чтения не имеет значения.', '<p>\nА = {Т, 4, ∆}<br/>\nB = {∆, 3, n}<br/>\nА ⋂ В = {∆}<br/>\nА ⋃ В = {∆, n, T, 4, 3}<br/>\n{4, ∆} ⊂ А    {∆, 2} ⊄ В    {∆} ⊂ А ⋂ В\n</p>\n\n\n<p>\nA = {м, о, р, е}, D = {д, о, м}, E = {д, ы, м}, A ⋂ D = {о, м}, D ⋂ м E = {д, м}, <br/> \n(A ⋂ D) ⋂ E = {м}, A ⋂ (D ⋂ E) = {м}.\n</p>\n\n\n<p>\nM ⋃ K = {1; 3; 5; 7; 9; 10}, K ⋃ T = {3; 5; 6; 9; 10}, (M ⋃ K) ⋃ T = {1; 3; 5; 6; 7; 9; 10}, M ⋃ (K ⋃ T) = {1; 3; 5; 6; 7; 9; 10}.\n</p>\n\n\n<p>\nа) красный и синий, синий и жёлтый, красный и жёлтый – 3 способа<br/>\nб) красный и синий, синий и жёлтый, красный и жёлтый, красный и красный, синий и синий, жёлтый и жёлтый – 6 способов.\n</p>\n\n\n<p>\nТак как она выбрала только 2 стихотворения М. Ю. Лермонтова, 2 стихотворения А. Блока, а необходимы стихотворения разных авторов, то только 2 программы своего выступления сможет составить Аня из этих стихов, если порядок их чтения не имеет значения.\n</p>\n\n\n<div class="img-wrapper-460">\n<img width="200" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica77-nomer86-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 77, номер 86-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 77, номер 86-1, год 2022."/>\n\n\n<div class="img-wrapper-460">\n<img width="150" src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica77-nomer87-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 77, номер 87-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 77, номер 87-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-77/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica77-nomer81.jpg', 'peterson/3/part3/page77/task81_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica77-nomer86-1.jpg', 'peterson/3/part3/page77/task81_solution_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 1, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica77-nomer87-1.jpg', 'peterson/3/part3/page77/task81_solution_1.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '0409846135703cf6ec011acbb58d0a215edb4675ad2198994a1003d4c69fc5d3', '2,4', NULL, 'по диаграмме эйлера-венна определи, из каких элементов состоят множества а и в. составь множества а ⋂ в и а ⋃ в. сделай записи, используя знаки ⊂ и ⊄:{4, ∆} ... а {∆, 2} ... в {∆} ... а ⋂ в');

    INSERT INTO textbook_tasks
    (textbook_id, page_number, task_number, task_order, condition_text, condition_html, solution_text, solution_html, hints_text, hints_html, has_sub_items, sub_items_json, source_url)
    VALUES (v_textbook_id, 79, '89', 0, 'Пользуясь заданным алгоритмом, найди значения х и сопоставь их соответствующим буквам. Расшифруй слово, расположив ответы примеров в порядке убывания:', '</p> \n<p class="text">Пользуясь заданным алгоритмом, найди значения х и сопоставь их соответствующим буквам. Расшифруй слово, расположив ответы примеров в порядке убывания:</p> \n\n<div class="description-text"> \n<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica79-nomer89.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 79, номер 89, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 79, номер 89, год 2022."/>\n</div>\n</div>', '', '<div class="img-wrapper-450">\n<img src="/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica79-nomer89-1.jpg" title="Учебник по математике 3 класс Петерсон - Часть 3, страница 79, номер 89-1, год 2022." alt="Учебник по математике 3 класс Петерсон, часть 3, страница 79, номер 89-1, год 2022."/>', '', '', FALSE, NULL, 'https://gdz-raketa.ru/matematika/3-klass/peterson-uchebnik/3-chast-stranitsa-79/')
    RETURNING id INTO v_task_id;

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'condition', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica79-nomer89.jpg', 'peterson/3/part3/page79/task89_condition_0.jpg', '');

    INSERT INTO textbook_task_images
    (task_id, image_type, image_order, sub_item_letter, original_url, local_path, alt_text)
    VALUES (v_task_id, 'solution', 0, NULL, 'https://gdz-raketa.ru/images/gdz/matematika/3klass/uchebnik-peterson/chast3/stranica79-nomer89-1.jpg', 'peterson/3/part3/page79/task89_solution_0.jpg', '');

    INSERT INTO textbook_task_index
    (task_id, textbook_id, grade, normalized_hash, numbers_signature, keywords, normalized_text)
    VALUES (v_task_id, v_textbook_id, 3, '36229a96e97e0627bf0625c25ec35a0739f7ffd1095610acc85b3c022af82f77', NULL, '["найди"]'::jsonb, 'пользуясь заданным алгоритмом, найди значения х и сопоставь их соответствующим буквам. расшифруй слово, расположив ответы примеров в порядке убывания');

END $$;
