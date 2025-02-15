#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$postgres" --dbname "$password" <<-EOSQL
    CREATE USER avito_merch WITH ENCRYPTED PASSWORD 'password';
    CREATE DATABASE avito_merch OWNER 'avito_merch';
    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO avito_merch;
EOSQL