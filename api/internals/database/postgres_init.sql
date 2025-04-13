CREATE DATABASE moviesdb;

\c moviesdb;

CREATE TABLE IF NOT EXISTS USERS (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(100) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS MOVIES (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT NOT NULL,
    images TEXT[],
    release_date DATE,
    director VARCHAR(100),
    rating_average NUMERIC(3, 2) DEFAULT 0.0,
    duration INT -- duration in minutes
);

CREATE TABLE IF NOT EXISTS RATINGS (
    id SERIAL PRIMARY KEY,
    user_id INT,
    movie_id INT NOT NULL,
    rating NUMERIC(3,1) NOT NULL,
    review TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE SET NULL,
    FOREIGN KEY (movie_id) REFERENCES MOVIES(id) ON DELETE CASCADE
);
