DROP DATABASE IF EXISTS blog;

CREATE DATABASE blog;

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";

USE blog;

CREATE TABLE `comment` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `text` TEXT, 
        `writerid` VARCHAR(255) NOT NULL DEFAULT "",
        `writerpw` VARCHAR(255) NOT NULL DEFAULT "",
        `admin` TINYINT(1) NOT NULL DEFAULT 0,
        `postid` INT NOT NULL DEFAULT 0);

CREATE TABLE `cookie` (
        `value` VARCHAR(255) NOT NULL DEFAULT "");

CREATE TABLE `visitor` (
        `today` INT NOT NULL DEFAULT 0,
        `total` INT NOT NULL DEFAULT 0);

CREATE TABLE `post` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `tag` VARCHAR(255) NOT NULL DEFAULT "",
        `title` VARCHAR(255) NOT NULL DEFAULT "",
        `text` TEXT NOT NULL DEFAULT "",
        `writetime` DATETIME NOT NULL DEFAULT "00-00-00",
        `imgpath` VARCHAR(255),
        `imgnum` INT);

CREATE TABLE `reply` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `admin` TINYINT(1) NOT NULL DEFAULT 0,
        `writerid` VARCHAR(255) NOT NULL DEFAULT "",
        `writerpw` VARCHAR(255) NOT NULL DEFAULT "",
        `text` TEXT NOT NULL DEFAULT "",
        `commentid` INT NOT NULL DEFAULT 0);

CREATE TABLE `beabouttodelete` (
        `delete_date` DATETIME DEFAULT '0000-00-00 00:00:00',
        `connection_id` INT,
        FOREIGN KEY (`connection_id`) REFERENCES `connection`(`connection_id`) ON DELETE CASCADE);

CREATE TABLE `question` (
        `question_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `target_word` VARCHAR(255) NOT NULL,
        `question_contents` VARCHAR(255) NOT NULL);

CREATE TABLE `answer` (
        `answer_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `connection_id` INT NOT NULL,
        `first_answer` VARCHAR(255) DEFAULT 'not-written',
        `second_answer` VARCHAR(255) DEFAULT 'not-written',
        `answer_date` VARCHAR(255) NOT NULL,
        `question_id` INT,
        FOREIGN KEY (`question_id`) REFERENCES `question`(`question_id`) ON UPDATE CASCADE ON DELETE CASCADE);

CREATE TABLE `exceptionword` (
        `exception_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `connection_id` INT NOT NULL,
        `except_word` TEXT NOT NULL);

CREATE TABLE `anniversary` (
        `anniversary_id` INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `connection_id` INT NOT NULL,
        `year` INT NOT NULL,
        `month` INT NOT NULL,
        `date` INT NOT NULL,
        `contents` VARCHAR(255) NOT NULL,
        `d_day` TINYINT(1) NOT NULL);