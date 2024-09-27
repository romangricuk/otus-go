-- Убедимся, что расширение uuid-ossp установлено
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Вставка тестовых данных в таблицы events и notifications
DO
$$
    DECLARE
        user_id              UUID;
        event_id             UUID;
        event_date           DATE;
        event_title          TEXT;
        event_description    TEXT;
        notification_time    TIME;
        notification_message TEXT;
    BEGIN
        -- Цикл по пользователям
        FOR user_counter IN 1..10
            LOOP
                user_id := uuid_generate_v4();

                -- Цикл по месяцам
                FOR month_counter IN 1..12
                    LOOP
                        -- Устанавливаем дату первого дня месяца
                        event_date := make_date(2024, month_counter, 1);

                        -- Цикл по дням в месяце
                        FOR day_counter IN 1..31
                            LOOP
                                -- Проверка на существование дня в месяце
                                BEGIN
                                    event_date := make_date(2024, month_counter, day_counter);
                                EXCEPTION
                                    WHEN OTHERS THEN
                                        -- Пропускаем несуществующие дни (например, 30 февраля)
                                        CONTINUE;
                                END;

                                event_id := uuid_generate_v4();
                                event_title := 'Событие ' || user_counter || ' - ' || to_char(event_date, 'YYYY-MM-DD');
                                event_description := 'Описание для ' || event_title;
                                notification_time := '08:00:00';
                                notification_message := 'Напоминание для ' || event_title;

                                -- Вставка события
                                INSERT INTO events (id, title, description, start_time, end_time, user_id)
                                VALUES (event_id,
                                        event_title,
                                        event_description,
                                        event_date + time '10:00:00',
                                        event_date + time '12:00:00',
                                        user_id);

                                -- Вставка уведомления
                                INSERT INTO notifications (id, event_id, user_id, time, message, sent)
                                VALUES (uuid_generate_v4(),
                                        event_id,
                                        user_id,
                                        event_date + notification_time,
                                        notification_message,
                                        'wait');
                            END LOOP;
                    END LOOP;
            END LOOP;
    END
$$;
