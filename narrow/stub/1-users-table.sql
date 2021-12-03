CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    name varchar NOT NULL,
    email varchar NOT NULL UNIQUE
);

-- stubbed users
INSERT INTO users
VALUES
       ('a6acab82-2b2e-484c-8e2b-f3f7736a26ed', 'foo', 'foo@bar.com'),
       ('20105afa-9a44-4f10-b50d-69884fbdd32c', 'bar', 'bar@bar.com'),
       ('1908007f-d7bc-4c39-845d-f11f436b579c', 'fiz', 'fiz@bar.com'),
       ('6f6b2530-bcf5-4208-b94a-3430533ec46e', 'buz', 'buz@bar.com');