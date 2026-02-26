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
    @id,
    @name,
    @slug,
    @types,
    @email,
    @website,
    @timezone,
    @phones,
    @address,
    ST_SetSRID(ST_MakePoint(@lon, @lat), 4326),
    @created_at,
    @updated_at
);

-- name: GetEstablishmentByID :one
SELECT
    id,
    name,
    slug,
    types,
    email,
    website,
    timezone,
    phones,
    address,
    ST_Y(location) AS lat,
    ST_X(location) AS lon,
    created_at,
    updated_at
FROM establishments
WHERE id = $1;

-- name: GetEstablishmentBySlug :one
SELECT
    id,
    name,
    slug,
    types,
    email,
    website,
    timezone,
    phones,
    address,
    CAST(ST_Y(location) AS double precision) AS lat,
    CAST(ST_X(location) AS double precision) AS lon,
    created_at,
    updated_at
FROM establishments
WHERE slug = $1;

-- name: DeleteEstablishment :exec
DELETE FROM establishments
WHERE id = $1;