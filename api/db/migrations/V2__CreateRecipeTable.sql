CREATE TABLE recipe (
    id serial PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    user_id INT,

    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES local_user(id)
);

CREATE INDEX recipe__user_id
    ON recipe (user_id);
