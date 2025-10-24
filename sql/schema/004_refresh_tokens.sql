-- +goose Up
create table refresh_tokens (
    token text primary key,
    created_at timestamp not null default now(),
    updated_at timestamp not null default now(),
    user_id uuid references users(id) on delete cascade not null,
    expires_at timestamp not null,
    revoked_at timestamp default null
);

-- +goose Down
drop table if exists refresh_tokens;