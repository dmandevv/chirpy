-- +goose Up
create table chirps (
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    body text not null,
    user_id UUID not null references users(id) on delete cascade
);
-- +goose Down
drop table chirps;