UPDATE users
SET display_name = $1, email = $2, password = $3, update_time = $4
WHERE id = $5
RETURNING id, organization_id, display_name, email, create_time, update_time;