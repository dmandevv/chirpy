-- name: CreateChirp :one
insert into chirps (id, created_at, updated_at, body, user_id) 
values (
    gen_random_uuid(), 
    now(),
    now(), 
    $1,
    $2
)
RETURNING *;

-- name: GetAllChirps :many
select *
from chirps
order by created_at asc;

-- name: GetAllChirpsOfUser :many
select *
from chirps
where user_id = $1
order by created_at asc;

-- name: GetChirpByID :one
select *
from chirps
where id = $1;

-- name: DeleteChirpByID :exec
delete from chirps
where id = $1;