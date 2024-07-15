CREATE TABLE IF NOT EXISTS repeate_items (
    id UUID NOT NULL DEFAULT gen_random_uuid(),
    url TEXT NOT NULL,
    jan VARCHAR,
    PRIMARY KEY(id)
);