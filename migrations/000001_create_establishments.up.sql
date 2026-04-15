CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE establishments (
    id uuid PRIMARY KEY DEFAULT uuidv7(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    types TEXT[] NOT NULL,

    email TEXT NOT NULL DEFAULT '',
    website TEXT NOT NULL DEFAULT '',
    timezone TEXT NOT NULL DEFAULT '',

    phones JSONB NOT NULL DEFAULT '[]'::jsonb,
    address JSONB NOT NULL DEFAULT '{}'::jsonb,

    location GEOGRAPHY(Point, 4326) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_establishments_slug ON establishments(slug);
CREATE INDEX idx_establishments_types ON establishments USING GIN(types);
CREATE INDEX idx_establishments_location ON establishments USING GIST(location);