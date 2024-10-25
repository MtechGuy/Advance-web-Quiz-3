CREATE TABLE IF NOT EXISTS signup (
    id bigserial PRIMARY KEY,
    email text NOT NULL,
    full_name text NOT NULL,
    version integer NOT NULL DEFAULT 1
);