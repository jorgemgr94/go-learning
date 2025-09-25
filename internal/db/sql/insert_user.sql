INSERT INTO users (id, name, email, created_at, updated_at)
VALUES ($1, $2, $3, NOW(), NOW())
RETURNING id;
