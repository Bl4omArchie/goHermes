CREATE TABLE Papers
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    link TEXT NOT NULL,
    publication_year INTEGER,
    category VARCHAR(255),
    file_data bytea,
);