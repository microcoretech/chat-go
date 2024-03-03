CREATE TABLE IF NOT EXISTS users
(
    id         BIGSERIAL PRIMARY KEY,
    email      VARCHAR   NOT NULL,
    username   VARCHAR   NOT NULL,
    role       SMALLINT  NOT NULL DEFAULT 1,
    first_name VARCHAR   NOT NULL,
    last_name  VARCHAR   NOT NULL,
    about_me   VARCHAR   NOT NULL DEFAULT '',
    image_url  VARCHAR   NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

INSERT INTO users
    (email, username, role, first_name, last_name)
VALUES
    ('admin@gmail.com', 'admin', 2, 'Admin', 'Admin'),
    ('mikhail.bobrovsky@gmail.com', 'mikhailbobrovsky', 1, 'Mikhail', 'Bobrovsky');