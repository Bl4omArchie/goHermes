-- database: ./eprint.db
DROP TABLE Document;
DROP TABLE Author;

CREATE TABLE Document (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    link TEXT NOT NULL UNIQUE,
    release TEXT NOT NULL,
    filepath TEXT UNIQUE,
    hash TEXT UNIQUE,
    license TEXT
);

CREATE TABLE Author (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT
);

CREATE TABLE AuthorDocument (
    author_id INTEGER,
    doc_id INTEGER,
    FOREIGN KEY (author_id) REFERENCES Author(id),
    FOREIGN KEY (doc_id) REFERENCES Document(id),
    PRIMARY KEY (author_id, doc_id)
);