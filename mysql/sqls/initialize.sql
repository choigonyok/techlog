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
        `writetime` VARCHAR(100) NOT NULL DEFAULT "00-00-00",
        `imgpath` VARCHAR(255),
        `imgnum` INT);

CREATE TABLE `reply` (
        `id` INT AUTO_INCREMENT PRIMARY KEY,
        `admin` TINYINT(1) NOT NULL DEFAULT 0,
        `writerid` VARCHAR(255) NOT NULL DEFAULT "",
        `writerpw` VARCHAR(255) NOT NULL DEFAULT "",
        `text` TEXT NOT NULL DEFAULT "",
        `commentid` INT NOT NULL DEFAULT 0);