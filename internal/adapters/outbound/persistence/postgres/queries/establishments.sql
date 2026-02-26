-- name: CreateEstablishment :exec
INSERT INTO establishments (
    id,
    name,
    slug,
    types,
    email,
    website,
    timezone,
    phones,
    address,
    location,
    created_at,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    ST_SetSRID(ST_MakePoint($10, $11), 4326),
    $12,
    $13
);

-- name: GetEstablishmentByID :one
SELECT * FROM establishments
WHERE id = $1;

-- name: GetEstablishmentBySlug :one
SELECT * FROM establishments
WHERE slug = $1;

-- name: DeleteEstablishment :exec
DELETE FROM establishments
WHERE id = $1;