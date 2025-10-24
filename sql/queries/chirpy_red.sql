-- name: UpgradeUserToChirpyRed :one
update users
set 
    updated_at = now(),
    is_chirpy_red = true
where id = $1
returning *;