CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    passport_serie INTEGER NOT NULL,
    passport_number INTEGER NOT NULL,
    surname VARCHAR(255),
    name VARCHAR(255),

    UNIQUE (passport_serie, passport_number)
);

CREATE TABLE tasks (
    task_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    task_name VARCHAR(255) NOT NULL,
    start_time TIMESTAMPTZ  NOT NULL DEFAULT CURRENT_TIMESTAMP,
    end_time TIMESTAMPTZ
);
