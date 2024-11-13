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

insert into posts (id, title, subtitle, text, write_time) values ('asdgadsgfhdc-ewahgdfa', 'Test Title', 'Test SubTitle', 'test \n## Contents', '2024/11/13');
insert into posts (id, title, subtitle, text, write_time) values ('esagdasgdsageasg', 'Test Title2', 'Test SubTitle2', 'test \n## Contents123', '2024/11/12');
insert into posts (id, title, subtitle, text, write_time) values ('kgvajndsvklnka', 'Test Title3', 'Test SubTitle3', 'test \n## Contents566', '2024/11/11');

insert into tags (name, post_id) values ('tag1', 'asdgadsgfhdc-ewahgdfa');
insert into tags (name, post_id) values ('tag2', 'asdgadsgfhdc-ewahgdfa');
insert into tags (name, post_id) values ('tag3', 'asdgadsgfhdc-ewahgdfa');
insert into tags (name, post_id) values ('Kubernetes', 'esagdasgdsageasg');
insert into tags (name, post_id) values ('tag2', 'esagdasgdsageasg');
insert into tags (name, post_id) values ('tag3', 'kgvajndsvklnka');
-- docker run -p 5432:5432 --name postgresql -e POSTGRES_PASSWORD=tester1234 -e TZ=Asia/Seoul -v ./hack/init.sql:/docker-entrypoint-initdb.d/init.sql -d postgres:latest