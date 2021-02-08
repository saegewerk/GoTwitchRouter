-- name: InsertUser :one
INSERT INTO twitch_user
(nick,name, accesslevel,created_at,updated_at)
VALUES ($1,$2,$3,$4,$5) RETURNING id;

-- name: GetUserUUIDByNick :one
SELECT id FROM twitch_user WHERE nick=$1;