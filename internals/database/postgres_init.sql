CREATE DATABASE moviesdb;

\c moviesdb;

CREATE EXTENSION IF NOT EXISTS vector;

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
    description_vec vector(3) NOT NULL
);

CREATE TABLE IF NOT EXISTS RATINGS (
    id SERIAL PRIMARY KEY,
    user_id INT,
    movie_id INT NOT NULL,
    rating NUMERIC(3,1) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE SET NULL,
    FOREIGN KEY (movie_id) REFERENCES MOVIES(id) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS CHATS (
    id SERIAL PRIMARY KEY,
    movie_id INT NOT NULL,
    user_id INT,
    text_content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (movie_id) REFERENCES MOVIES(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE SET NULL
);


CREATE TABLE IF NOT EXISTS WATCHLIST (
    id SERIAL PRIMARY KEY,
    movie_id INT NOT NULL,
    user_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (movie_id) REFERENCES MOVIES(id) ON DELETE SET NULL,
    FOREIGN KEY (user_id) REFERENCES USERS(id) ON DELETE CASCADE
);



-- performance optimization:

-- for fast ILIKE(%moviename%) search operation 
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE INDEX IF NOT EXISTS idx_movies_name_trgm ON movies USING GIN (name gin_trgm_ops);

-- for fast query on chats based on movie_id
CREATE INDEX IF NOT EXISTS idx_chats_movie_id ON chats(movie_id); 

-- faster query of ratings based on movie_id
CREATE INDEX IF NOT EXISTS idx_ratings_movie_id ON ratings(movie_id);

-- faster query of user ratings on named movies
CREATE INDEX IF NOT EXISTS idx_ratings_user_movie ON ratings(user_id, movie_id);


-- faster query of watchlist based on user_id
CREATE INDEX IF NOT EXISTS idx_watchlist_user_id ON watchlist(user_id);
