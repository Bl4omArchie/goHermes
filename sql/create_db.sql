-- SQLite

CREATE TABLE Document (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT,
    document_link TEXT NOT NULL,
    document_type TEXT NOT NULL
);