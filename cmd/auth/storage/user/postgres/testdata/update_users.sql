UPDATE users
SET username = $1, email = $2, password = $3, update_time = $4
WHERE id = $5
RETURNING id, username, email, create_time, update_time;