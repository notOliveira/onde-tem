CREATE EXTENSION IF NOT EXISTS postgis;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE establishments (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    types TEXT[] NOT NULL,
    email TEXT,
    website TEXT,
    phones JSONB,
    location GEOGRAPHY(Point, 4326) NOT NULL,
    address JSONB,
    timezone TEXT,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE INDEX idx_establishments_slug ON establishments(slug);
CREATE INDEX idx_establishments_types ON establishments USING GIN(types);
CREATE INDEX idx_establishments_location ON establishments USING GIST(location);