-- Création de la base de données
CREATE DATABASE ePrint;

-- Connexion à la base de données nouvellement créée
\c eprint;

-- Création de la table des utilisateurs
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Création de la table des commandes
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT REFERENCES users(user_id) ON DELETE CASCADE,
    product_name VARCHAR(100) NOT NULL,
    amount DECIMAL(10, 2) NOT NULL,
    order_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- Table Papers
CREATE TABLE Papers
(
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    link TEXT NOT NULL,
    publication_year INTEGER,
    category VARCHAR(255),
    file_data bytea
);