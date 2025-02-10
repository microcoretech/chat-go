-- Copyright 2025 Mykhailo Bobrovskyi
--
-- Licensed under the Apache License, Version 2.0 (the "License");
-- you may not use this file except in compliance with the License.
-- You may obtain a copy of the License at
--
--     http://www.apache.org/licenses/LICENSE-2.0
--
-- Unless required by applicable law or agreed to in writing, software
-- distributed under the License is distributed on an "AS IS" BASIS,
-- WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
-- See the License for the specific language governing permissions and
-- limitations under the License.

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