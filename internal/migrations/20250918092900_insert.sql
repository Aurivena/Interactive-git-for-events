-- +goose Up
-- +goose StatementBegin
INSERT INTO place (id, title, address, description, lon, lat, tier, kind, tags) VALUES
-- ===== КИНО =====
('11111111-1111-1111-1111-111111111111','Россия','Курган, ул. Володарского, 75',NULL,65.339361,55.439371,'standard','cinema',
 jsonb_build_object(
         'website','https://россия45.рф',
         'phone','+7 (3522) 60-52-50',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday',   'start','09:00','end','21:00','spans_midnight',false),
                 jsonb_build_object('week','tuesday',  'start','09:00','end','21:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','09:00','end','21:00','spans_midnight',false),
                 jsonb_build_object('week','thursday', 'start','09:00','end','21:00','spans_midnight',false),
                 jsonb_build_object('week','friday',   'start','09:00','end','21:00','spans_midnight',false),
                 jsonb_build_object('week','saturday', 'start','09:00','end','21:00','spans_midnight',false),
                 jsonb_build_object('week','sunday',   'start','09:00','end','21:00','spans_midnight',false)
                     )
 )
),
('11111111-1111-1111-1111-111111111112','Pushka','Курган, ул. Пушкина, 25, ТРЦ «Пушкинский», 3 этаж',NULL,65.318954,55.432190,'standard','cinema',
 jsonb_build_object(
         'website','https://cinema.pushka.club/kurgan/pushka',
         'phone','+7 (3522) 60-70-55',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday','start','09:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','tuesday','start','09:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','09:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','thursday','start','09:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','friday','start','09:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','saturday','start','09:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','sunday','start','09:00','end','23:00','spans_midnight',false)
                     )
 )
),
('11111111-1111-1111-1111-111111111113','Ultra Cinema Курган','Курган, ул. Коли Мяготина, 8, ТРЦ «Hyper City», 2 этаж',NULL,65.280027,55.426618,'standard','cinema',
 jsonb_build_object(
         'website','https://kurgan.ultra-cinema.ru',
         'phone','+7 (3522) 22-89-87',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday',   'start','10:00','end','00:00','spans_midnight',true),
                 jsonb_build_object('week','tuesday',  'start','10:00','end','00:00','spans_midnight',true),
                 jsonb_build_object('week','wednesday','start','10:00','end','00:00','spans_midnight',true),
                 jsonb_build_object('week','thursday', 'start','10:00','end','00:00','spans_midnight',true),
                 jsonb_build_object('week','friday',   'start','10:00','end','00:00','spans_midnight',true),
                 jsonb_build_object('week','saturday', 'start','10:00','end','00:00','spans_midnight',true),
                 jsonb_build_object('week','sunday',   'start','10:00','end','00:00','spans_midnight',true)
                     )
 )
),
('11111111-1111-1111-1111-111111111114','Клумба Синема','Курган, 2-й микрорайон, 17, ТРЦ «Стрекоза», 2 этаж',NULL,65.264725,55.464850,'value','cinema',
 jsonb_build_object(
         'website','https://klumba-cinema.ru',
         'phone','+7 (963) 869-80-49',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday',   'start','09:00','end','02:00','spans_midnight',true),
                 jsonb_build_object('week','tuesday',  'start','09:00','end','02:00','spans_midnight',true),
                 jsonb_build_object('week','wednesday','start','09:00','end','02:00','spans_midnight',true),
                 jsonb_build_object('week','thursday', 'start','09:00','end','02:00','spans_midnight',true),
                 jsonb_build_object('week','friday',   'start','09:00','end','02:00','spans_midnight',true),
                 jsonb_build_object('week','saturday', 'start','09:00','end','02:00','spans_midnight',true),
                 jsonb_build_object('week','sunday',   'start','09:00','end','02:00','spans_midnight',true)
                     )
 )
),

-- ===== РЕСТОРАНЫ =====
('22222222-2222-2222-2222-222222222221','Палермо','Курган, ул. Карла Маркса, 58',NULL,65.346068,55.438404,'premium','restaurant',
 jsonb_build_object(
         'website','https://palermo-express.ru',
         'phone','+7 (3522) 46-06-66',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','tuesday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','thursday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','friday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','saturday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','sunday','start','10:00','end','22:00','spans_midnight',false)
                     )
 )
),
('22222222-2222-2222-2222-222222222222','Виновники','Курган, ул. Гоголя, 83',NULL,65.343660,55.441780,'upscale','restaurant',
 jsonb_build_object(
         'website','https://vinovniki.ru',
         'phone','+7 (3522) 55-77-88',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','tuesday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','thursday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','friday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','saturday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','sunday','start','10:00','end','23:00','spans_midnight',false)
                     )
 )
),
('22222222-2222-2222-2222-222222222223','Мачете','Курган, ул. Гоголя, 61',NULL,65.342753,55.441351,'standard','restaurant',
 jsonb_build_object(
         'website','https://machete.pro',
         'phone','+7 (912) 835-21-11',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','tuesday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','thursday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','friday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','saturday','start','10:00','end','23:00','spans_midnight',false),
                 jsonb_build_object('week','sunday','start','10:00','end','23:00','spans_midnight',false)
                     )
 )
),
('22222222-2222-2222-2222-222222222224','Mamma Mia','Курган, ул. Рихарда Зорге, 35',NULL,65.315965,55.437327,'value','restaurant',
 jsonb_build_object(
         'website','https://mammamia45.ru',
         'phone','+7 (3522) 44-55-66',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','tuesday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','thursday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','friday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','saturday','start','10:00','end','22:00','spans_midnight',false),
                 jsonb_build_object('week','sunday','start','10:00','end','22:00','spans_midnight',false)
                     )
 )
),

-- ===== МУЗЕИ =====
('33333333-3333-3333-3333-333333333331','Художественный музей им. Травникова','Курган, ул. Максима Горького, 127/4',NULL,65.352805,55.440488,'standard','museum',
 jsonb_build_object(
         'website','https://culture.ru/institutes/44437',
         'phone','+7 (3522) 46-59-46',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday',   'closed',true),
                 jsonb_build_object('week','tuesday',  'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','thursday', 'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','friday',   'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','saturday', 'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','sunday',   'start','10:00','end','18:00','spans_midnight',false)
                     )
 )
),
('33333333-3333-3333-3333-333333333332','Музей истории города Кургана','Курган, ул. Куйбышева, 59',NULL,65.348682,55.434737,'value','museum',
 jsonb_build_object(
         'website','https://culture.ru/institutes/44442',
         'phone','+7 (3522) 41-24-94',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday',   'closed',true),
                 jsonb_build_object('week','tuesday',  'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','thursday', 'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','friday',   'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','saturday', 'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','sunday',   'start','10:00','end','18:00','spans_midnight',false)
                     )
 )
),
('33333333-3333-3333-3333-333333333333','Краеведческий музей','Курган, ул. Пушкина, 137',NULL,65.337381,55.440514,'standard','museum',
 jsonb_build_object(
         'website','https://culture.ru/institutes/44441',
         'phone','+7 (3522) 46-09-38',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday',   'closed',true),
                 jsonb_build_object('week','tuesday',  'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','thursday', 'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','friday',   'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','saturday', 'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','sunday',   'start','10:00','end','18:00','spans_midnight',false)
                     )
 )
),
('33333333-3333-3333-3333-333333333334','Дом-музей декабристов','Курган, ул. Климова, 80А',NULL,65.353470,55.434630,'value','museum',
 jsonb_build_object(
         'website','https://culture.ru/institutes/44443',
         'phone','+7 (3522) 46-67-65',
         'schedule', jsonb_build_array(
                 jsonb_build_object('week','monday',   'closed',true),
                 jsonb_build_object('week','tuesday',  'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','wednesday','start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','thursday', 'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','friday',   'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','saturday', 'start','10:00','end','18:00','spans_midnight',false),
                 jsonb_build_object('week','sunday',   'start','10:00','end','18:00','spans_midnight',false)
                     )
 )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
