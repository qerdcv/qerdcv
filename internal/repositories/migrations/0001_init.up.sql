CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(35) NOT NULL UNIQUE,
    password_hash VARCHAR(60) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- TODO: add user_id idx
CREATE TABLE IF NOT EXISTS user_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    expires_at TIMESTAMP NOT NULL
);


-- TODO: add user_id idx
CREATE TABLE IF NOT EXISTS budget_categories (
    id UUID PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    name varchar(35) NOT NULL,
    UNIQUE (user_id, name)
);

CREATE TABLE IF NOT EXISTS budget_transactions (
    id UUID PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    category_id UUID REFERENCES budget_categories(id),
    amount BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
