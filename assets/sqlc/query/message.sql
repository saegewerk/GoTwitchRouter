-- name: InsertMessage :one
INSERT INTO msg
( msg,msg_id, event,fk_user,created_at)
VALUES ($1,$2,$3,$4,$5) RETURNING id;