CREATE TABLE Papers
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    link TEXT NOT NULL,
    publication_date DATE,
    categorie VARCHAR(255),
);