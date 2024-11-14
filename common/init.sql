CREATE DATABASE visitor_db;
CREATE DATABASE post_db;

\c visitor_db;

CREATE TABLE visitor (
    date VARCHAR(100) PRIMARY KEY,
    count SMALLINT NOT NULL DEFAULT 0
);

\c post_db;

CREATE TABLE posts (
    id VARCHAR(100) PRIMARY KEY,
    title VARCHAR(40) NOT NULL,
    subtitle VARCHAR(40) NOT NULL,
    text TEXT NOT NULL,
    write_time VARCHAR(40) NOT NULL
);

CREATE TABLE images (
    id VARCHAR(100) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    post_id VARCHAR(100) NOT NULL,
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE TABLE thumbnails (
    image_id VARCHAR(100) NOT NULL,
    post_id VARCHAR(100) NOT NULL,
    PRIMARY KEY (image_id, post_id),
    FOREIGN KEY (image_id) REFERENCES images(id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);

CREATE TABLE tags (
    name VARCHAR(100) NOT NULL,
    post_id VARCHAR(100) NOT NULL,
    PRIMARY KEY (name, post_id),
    FOREIGN KEY (post_id) REFERENCES posts(id)
);