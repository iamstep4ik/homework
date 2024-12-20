--+goose Up
-- +goose StatementBegin
create extension if not exists citext;

create table Users(
    id bigint primary key generated by default as identity,
    username citext not null unique,
    email citext not null unique,
    password_hash varchar(255),
    created_at timestamp(0) not null default(now() at time zone 'utc'),
    registred_events int[]

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table Users;
-- +goose StatementEnd 
