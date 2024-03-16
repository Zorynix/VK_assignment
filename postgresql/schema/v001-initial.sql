DROP SCHEMA IF EXISTS vk CASCADE;
CREATE SCHEMA IF NOT EXISTS vk;
CREATE TABLE vk.actors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    gender CHAR(1),
    date_of_birth DATE,
    CHECK (gender IN ('M', 'F'))
);
CREATE TABLE vk.movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL CHECK (
        LENGTH(title) >= 1
        AND LENGTH(title) <= 150
    ),
    description VARCHAR(1000),
    release_date DATE,
    rating DECIMAL(2, 1) CHECK (
        rating >= 1
        AND rating <= 10
    )
);
CREATE TABLE vk.actor_movies (
    actor_name VARCHAR(255) NOT NULL,
    movie_title VARCHAR(150) NOT NULL,
    PRIMARY KEY (actor_name, movie_title),
    CONSTRAINT fk_actor_name FOREIGN KEY (actor_name) REFERENCES vk.actors(name) ON DELETE CASCADE,
    CONSTRAINT fk_movie_title FOREIGN KEY (movie_title) REFERENCES vk.movies(title) ON DELETE CASCADE
);