-- name: CreateRefreshToken :one
insert into refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
values (
    $1,
    now(),
    now(),
    $2,
    $3,
    null
)
returning *;

-- name: GetRefreshToken :one
select *
from refresh_tokens
where token = $1;

-- name: GetUserFromRefreshToken :one
select u.id, u.created_at, u.updated_at, u.email, u.hashed_password
from users u
inner join refresh_tokens rt on rt.user_id = u.id
where rt.token = $1 and rt.revoked_at is null and rt.expires_at > now();

-- name: RevokeRefreshToken :exec
update refresh_tokens
set revoked_at = now(), updated_at = now()
where token = $1;
