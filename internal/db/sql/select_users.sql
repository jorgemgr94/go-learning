SELECT 
    id,
    name,
    email,
    created_at,
    updated_at
FROM users 
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
