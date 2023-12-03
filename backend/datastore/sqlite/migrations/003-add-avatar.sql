CREATE TABLE avatars (
    username TEXT PRIMARY KEY,
    avatar BLOB NOT NULL,
    width INTEGER NOT NULL,
    last_modified TEXT NOT NULL
);
