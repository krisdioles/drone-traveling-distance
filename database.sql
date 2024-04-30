-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- This is test table. Remove this table and replace with your own tables. 
CREATE TABLE test (
	id serial PRIMARY KEY,
	name VARCHAR ( 50 ) UNIQUE NOT NULL,
);

CREATE TABLE estates (
	id uuid DEFAULT gen_random_uuid(),
	length INT NOT NULL,
	width INT NOT NULL
);

CREATE TABLE trees (
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	x INT NOT NULL,
	y INT NOT NULL,
	height INT NOT NULL,
	estate_id uuid NOT NULL
);
