create table chat
(
    id         integer not null
        constraint chat_pk
            primary key,
    type       text,
    username   text,
    first_name text,
    last_name  text
);

