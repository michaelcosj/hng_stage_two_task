-- Write your migrate up statements here
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(16)
);

CREATE TABLE organisations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL
);

CREATE TABLE user_organisations (
    user_id UUID REFERENCES users (id),
    org_id UUID REFERENCES organisations (id),

    PRIMARY KEY (user_id, org_id)
);

---- create above / drop below ----
DROP TABLE users;
DROP TABLE organisations;
DROP TABLE user_organisations;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
