create database if not exists `tunebot`;

use `tunebot`;

create table if not exists `user` (
    `id` binary(16) primary key,
    `username` varchar(32) not null unique,
    `password` binary(60) not null
);

create table if not exists `song` (
    `id` binary(16) primary key,
	`url` varchar(256) not null unique
);

create table if not exists `playlist` (
    `id` binary(16) primary key,
    `name` varchar(64) not null,
    `user_id` binary(16) not null,
	`enabled` bit not null,
    foreign key(`user_id`) references `user`(`id`),
    constraint `unique_playlist_name_per_user` unique (`user_id`, `name`)
);

create table if not exists `playlist_song` (
    `id` binary(16) primary key,
    `song_id` binary(16) not null,
    `playlist_id` binary(16) not null,
    foreign key(`song_id`) references `song`(`id`),
    foreign key(`playlist_id`) references `playlist`(`id`),
    constraint `no_duplicate_song_per_playlist` unique (`song_id`, `playlist_id`)
);

create user 'tunebot'@'localhost' identified by '6N9.h+Q.H*ah.zPZ';
grant all privileges on `tunebot`.* to 'tunebot'@'localhost';
flush privileges;