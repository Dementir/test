create table region (
	region_id uuid not null primary key,
	name text not null
);
comment on table region is 'Справочник регионов';

create table channel (
	channel_id uuid not null primary key,
	name text not null
);
comment on table channel is 'Справочник каналов';

create table status (
	status_id uuid not null primary key,
	name text not null
);
comment on table status is 'Справочник статусов';

create table service (
	service_id uuid not null primary key,
	name text not null
);
comment on table service is 'Справочник сервисов';

create table log (
	stamp timestamp without time zone not null,
	channel uuid not null references channel(channel_id),
	region uuid not null references region(region_id),
	service uuid not null references service(service_id),
	status uuid not null references status(status_id),
	"user" bigint not null
);
comment on table log is 'Лог событий из лог-файла';

-- JOB2
create table "user" (
    user_id bigint not null primary key,
    name text
);
comment on table "user" is 'Список пользователей с идентификаторами';

create table question (
    question_id serial not null primary key,
    question_text text not null,
    right_answer text not null,
    wrong_answer text[] not null
);
comment on table question is 'Каталог вопросов с ответами';
comment on column question.question_id is 'Идентификатор вопроса';
comment on column question.question_text is 'Текст вопроса';
comment on column question.right_answer is 'Правильный вариант ответа';
comment on column question.wrong_answer is 'Массив неправильных вариантов ответов';

create table game (
    "user" bigint not null references "user"(user_id),
    question int not null references question(question_id),
    right_answer text,
    answered boolean not null default false
);
comment on table game is 'Информация об играх пользователей';
comment on column game.user is 'Идентификатор пользователя';
comment on column game.question is 'Идентификатор заданного вопроса';
comment on column game.right_answer is 'Идентификатор правильного ответа';
comment on column game.answered is 'Ответил ли пользователь на данный вопрос';

CREATE INDEX idx_user_game ON game("user");

create table rating (
    "user" bigint not null,
    right_answers int not null default 0,
    answer_time timestamp without time zone
);
comment on table rating is 'Результаты пользователей';
comment on column rating.user is 'Идентификатор пользователя';
comment on column rating.right_answers is 'Количество правильных ответов у пользователя';
comment on column rating.answer_time is 'Время последнего ответа пользователя';
