CREATE TABLE users(u_id SERIAL PRIMARY KEY,
username VARCHAR(255) UNIQUE NOT NULL,
password VARCHAR(255) NOT NULL
);

CREATE TYPE status AS ENUM('Not Started', 'In Progress', 'Completed');

CREATE TABLE todos(t_id SERIAL PRIMARY KEY,
title TEXT NOT NULL,
current_status status DEFAULT 'Not Started',
u_id INT NOT NULL REFERENCES users(u_id) ON DELETE CASCADE
);

