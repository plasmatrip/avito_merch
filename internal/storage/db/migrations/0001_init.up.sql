BEGIN;

-- create tables and indexes for users
CREATE TABLE IF NOT EXISTS users (
	id uuid NOT NULL UNIQUE PRIMARY KEY,
    date timestamp with time zone NOT NULL, -- registration date
	login varchar(64) NOT NULL UNIQUE, -- login
	password varchar(64) NOT NULL -- password hash
);
CREATE INDEX IF NOT EXISTS users_id ON users USING GIN (id);
CREATE INDEX IF NOT EXISTS users_login ON users (login);

-- create tables and indexes for transactions
CREATE TABLE IF NOT EXISTS transactions (
    id uuid NOT NULL UNIQUE PRIMARY KEY ,
    date timestamp with time zone NOT NULL, -- transaction date
    from_user_id uuid NOT NULL REFERENCES users (id), -- sender
    to_user_id uuid NOT NULL REFERENCES users (id), -- recipient
    amount money NOT NULL DEFAULT '0' -- amount of coins
);
CREATE INDEX IF NOT EXISTS transactions_id ON transactions USING GIN (id);
CREATE INDEX IF NOT EXISTS transactions_from_user_id ON transactions USING GIN (from_user_id);
CREATE INDEX IF NOT EXISTS transactions_to_user_id ON transactions USING GIN (to_user_id);

-- create tables and indexes for merch
CREATE TABLE IF NOT EXISTS merch (
    id uuid NOT NULL UNIQUE PRIMARY KEY ,
    name varchar(64) NOT NULL UNIQUE, -- merch name
    price money NOT NULL DEFAULT '0' -- price
);
CREATE INDEX IF NOT EXISTS merch_id ON merch USING GIN (id);
CREATE INDEX IF NOT EXISTS merch_name ON merch (name);

-- insert some merch
INSERT INTO merch (id, name, price) VALUES
(gen_random_uuid (), 't-shirt', 80),
(gen_random_uuid (), 'cup', 20),
(gen_random_uuid (), 'book', 50),
(gen_random_uuid (), 'pen', 10),
(gen_random_uuid (), 'powerbank', 200),
(gen_random_uuid (), 'hoody', 300),
(gen_random_uuid (), 'umbrella', 200),
(gen_random_uuid (), 'socks', 10),
(gen_random_uuid (), 'wallet', 50),
(gen_random_uuid (), 'pink-hoody', 500);

-- create tables and indexes for purchases
CREATE TABLE IF NOT EXISTS purchases (
    id uuid NOT NULL UNIQUE PRIMARY KEY,
    date timestamp with time zone NOT NULL, -- purchase date
    user_id uuid NOT NULL REFERENCES users (id), -- user id
    merch_id uuid NOT NULL REFERENCES merch (id) -- merch id
);
CREATE INDEX IF NOT EXISTS purchases_id ON purchases USING GIN (id);
CREATE INDEX IF NOT EXISTS purchases_user_id ON purchases USING GIN (user_id);
CREATE INDEX IF NOT EXISTS purchases_merch_id ON purchases USING GIN (merch_id);

-- create tables and indexes for accounts
CREATE TABLE IF NOT EXISTS accounts (
	id uuid NOT NULL UNIQUE PRIMARY KEY,
	user_id uuid NOT NULL UNIQUE REFERENCES users (id), -- user id
	amount money NOT NULL DEFAULT '0' -- amount of coins
);
CREATE INDEX IF NOT EXISTS accounts_id ON accounts USING GIN (id);
CREATE INDEX IF NOT EXISTS accounts_user_id ON accounts USING GIN (user_id);

COMMIT;