DROP SCHEMA IF EXISTS vk CASCADE;
DROP TABLE IF EXISTS vk.actors;
DROP TABLE IF EXISTS vk.movies;
DROP TABLE IF EXISTS vk.actor_movie;

CREATE SCHEMA IF NOT EXISTS vk;

CREATE TABLE vk.actors (
    actor_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    gender CHAR(1),
    date_of_birth DATE,
    CHECK (gender IN ('M', 'F'))
);

CREATE TABLE vk.movies (
    movie_id SERIAL PRIMARY KEY,
    title VARCHAR(150) NOT NULL CHECK (LENGTH(title) >= 1 AND LENGTH(title) <= 150),
    description VARCHAR(1000),
    release_date DATE,
    rating DECIMAL(2, 1) CHECK (rating >= 1 AND rating <= 10)
);

CREATE TABLE vk.actor_movie (
    actor_id INT NOT NULL,
    movie_id INT NOT NULL,
    PRIMARY KEY (actor_id, movie_id),
    FOREIGN KEY (actor_id) REFERENCES vk.actors (actor_id) ON DELETE CASCADE,
    FOREIGN KEY (movie_id) REFERENCES vk.movies (movie_id) ON DELETE CASCADE
);
