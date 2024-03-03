CREATE TABLE IF NOT EXISTS user_credentials
(
    user_id  BIGINT REFERENCES users ("id") ON DELETE CASCADE,
    password VARCHAR NOT NULL
);

INSERT INTO user_credentials
    (user_id, password)
VALUES
    --password: 'mohnIeih4Ju9zHYE1VPWL0mHyzBjyFPl'
    ((SELECT id FROM users WHERE email = 'admin@gmail.com'), '$2a$12$0x5RuE//aFNXsaID8WF8B.8HPWLeRHLJgQYtBa7hLMz2o4o4nBkJi'),
    --password: 'mohnIeih4Ju9zHYE1VPWL0mHyzBjyFPl'
    ((SELECT id FROM users WHERE email = 'mikhail.bobrovsky@gmail.com'), '$2a$12$0x5RuE//aFNXsaID8WF8B.8HPWLeRHLJgQYtBa7hLMz2o4o4nBkJi');