DROP DATABASE IF EXISTS blogdb;

CREATE DATABASE blogdb;

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

USE blogdb;

DROP TABLE IF EXISTS `cookie`;
CREATE TABLE `cookie` (
    `value` VARCHAR(255) NOT NULL DEFAULT ''
);

DROP TABLE IF EXISTS `visitor`;
CREATE TABLE `visitor` (
        `today` INT NOT NULL DEFAULT 0,
        `total` INT NOT NULL DEFAULT 0,
        `date` DATE DEFAULT NULL
);
INSERT INTO `visitor` VALUES (1,199,'2023-10-29');

DROP TABLE IF EXISTS `post`;
CREATE TABLE `post` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `tags` VARCHAR(40) NOT NULL DEFAULT '',
    `title` VARCHAR(40) NOT NULL DEFAULT '',
    `text` TEXT,
    `writeTime` DATE NOT NULL DEFAULT '2000-01-01',
    PRIMARY KEY (`id`)
);
INSERT INTO `post` (`tags`, `title`, `text`, `writeTime`) VALUES ("MSA DEV OPS","title1", "text1", '2023-10-29');
INSERT INTO `post` (`tags`, `title`, `text`, `writeTime`) VALUES ("DEV GOLANG","title2", "text2", '2023-10-30');

DROP TABLE IF EXISTS `image`;
CREATE TABLE `image` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `postID` INT NOT NULL DEFAULT 0,
    `imageName` VARCHAR(20) NOT NULL DEFAULT '',
    `thumbnail` TINYINT(1),
    PRIMARY KEY (`id`),
    FOREIGN KEY (`postID`) REFERENCES `post` (`id`) ON DELETE CASCADE
);
INSERT INTO `image` (`postID`, `imageName`, `thumbnail`) VALUES (1, 'image1.JPEG', 1);
INSERT INTO `image` (`postID`, `imageName`, `thumbnail`) VALUES (2, 'image2.jpg', 1);

DROP TABLE IF EXISTS `comment`;
CREATE TABLE `comment` (
    `id` INT NOT NULL AUTO_INCREMENT,
    `text` TEXT,
    `writerID` VARCHAR(13) NOT NULL DEFAULT '',
    `writerPW` VARCHAR(8) NOT NULL DEFAULT '',
    `admin` TINYINT(1) NOT NULL DEFAULT 0,
    `postID`INT NOT NULL DEFAULT 0,
    PRIMARY KEY (`id`),
    FOREIGN KEY (`postID`) REFERENCES `post` (`id`) ON DELETE CASCADE
);
INSERT INTO `comment` (`text`, `writerID`, `writerPW`, `admin`, `postID`) VALUES ('comment1', 'comid1', 
'1234', 1, 1);
INSERT INTO `comment` (`text`, `writerID`, `writerPW`, `admin`, `postID`) VALUES ('comment2', 'comid2', '2345', 0, 1);

DROP TABLE IF EXISTS `reply`;
CREATE TABLE `reply` (
        `id` INT NOT NULL AUTO_INCREMENT,
        `admin` TINYINT(1) NOT NULL DEFAULT 0,
        `writerID` VARCHAR(13) NOT NULL DEFAULT '',
        `writerPW` VARCHAR(8) NOT NULL DEFAULT '',
        `text` TEXT,
        `commentID` INT NOT NULL DEFAULT 0,
        PRIMARY KEY (`id`),
        FOREIGN KEY (`commentID`) REFERENCES `comment` (`id`) ON DELETE CASCADE
);
INSERT INTO `reply` (`admin`, `writerID`, `writerPW`, `text`, `commentID`) VALUES (1, 'id1', '1234', 'reply1', 1);
INSERT INTO `reply` (`admin`, `writerID`, `writerPW`, `text`, `commentID`) VALUES (2, 'id2', '2345', 'reply2', 2);