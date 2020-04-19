CREATE TABLE local_user (
    id serial PRIMARY KEY,
    email VARCHAR(200) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL,
    lid UUID DEFAULT uuid_generate_v4()
);
