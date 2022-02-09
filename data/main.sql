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
    `user_id` binary(16) not null,
    `is_blacklist` bit,
	`enabled` bit not null,
    foreign key(`user_id`) references `user`(`id`),
    constraint `one_blacklist_per_user` unique (`user_id`, `is_blacklist`)
);

create table if not exists `playlist_song` (
    `id` binary(16) primary key,
    `song_id` binary(16) not null,
    `playlist_id` binary(16) not null,
    foreign key(`song_id`) references `song`(`id`),
    foreign key(`playlist_id`) references `playlist`(`id`),
    constraint `no_duplicate_song_per_playlist` unique (`song_id`, `playlist_id`)
);