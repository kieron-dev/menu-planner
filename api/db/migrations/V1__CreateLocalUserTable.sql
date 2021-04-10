CREATE TABLE local_user (
    id serial PRIMARY KEY,
    email VARCHAR(200) UNIQUE NOT NULL,
    name VARCHAR(200) NOT NULL
);

CREATE UNIQUE INDEX local_user__email
    ON local_user (email);
