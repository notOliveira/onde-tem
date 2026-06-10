-- name: CreateEstablishment :one
INSERT INTO establishments (
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
    @name,
    @slug,
    @types,
    @email,
    @website,
    @timezone,
    @phones,
    @address,
    ST_SetSRID(ST_MakePoint(@lon::double precision, @lat::double precision), 4326),
    @created_at,
    @updated_at
)
RETURNING id;

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
    CAST(ST_Y(location) AS double precision) AS lat,
    CAST(ST_X(location) AS double precision) AS lon,
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

-- name: UpdateEstablishment :exec
UPDATE establishments
SET
    name = @name,
    slug = @slug,
    types = @types,
    email = @email,
    website = @website,
    timezone = @timezone,
    phones = @phones,
    address = @address,
    location = ST_SetSRID(ST_MakePoint(@lon::double precision, @lat::double precision), 4326),
    updated_at = @updated_at
WHERE id = @id;

-- name: ListEstablishments :many
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
    CAST(ST_Y(location::geometry) AS double precision) AS lat,
    CAST(ST_X(location::geometry) AS double precision) AS lon,
    created_at,
    updated_at
FROM establishments
WHERE
    (
        @types::text[] IS NULL
        OR types && @types::text[]
    )
    AND (
        @search::text IS NULL
        OR @search::text = ''
        OR name ILIKE '%' || @search::text || '%'
        OR slug ILIKE '%' || @search::text || '%'
    )
ORDER BY created_at DESC
LIMIT $1
OFFSET $2;


-- name: DeleteEstablishment :exec
DELETE FROM establishments
WHERE id = $1;

