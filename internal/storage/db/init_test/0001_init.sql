CREATE USER test_avito_merch WITH ENCRYPTED PASSWORD 'password';

CREATE DATABASE test_avito_merch OWNER 'test_avito_merch';

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO test_avito_merch;
