create table "user"
(
    chat_id integer not null
        constraint user_pk
            primary key,
    grade   integer
);