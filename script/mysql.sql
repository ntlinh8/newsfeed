create database newsfeed;

use newsfeed;

CREATE TABLE db_users (
    id int auto_increment primary key,
    username varchar(255),
    hash_password varchar(255),
    last_name varchar(255),
    first_name varchar(255),
    dob int,
    email varchar(255)
);