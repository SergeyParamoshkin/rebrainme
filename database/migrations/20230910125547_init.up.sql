CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
   id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
   username VARCHAR (50) UNIQUE NOT NULL,
   password VARCHAR (50) NOT NULL,
   email VARCHAR (300) UNIQUE NOT NULL
);

INSERT INTO users (username, password, email) VALUES ('admin', '', 'serg3091@gmail.com');