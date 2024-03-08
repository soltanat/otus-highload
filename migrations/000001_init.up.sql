CREATE TABLE users
(
    id          uuid PRIMARY KEY,
    first_name  VARCHAR(255),
    second_name VARCHAR(255),
    birthdate   DATE,
    biography   TEXT,
    city        VARCHAR(255),
    password    VARCHAR(255) NOT NULL
);
