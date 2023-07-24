CREATE TABLE IF NOT EXISTS `users`(
	`user_id`		INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
	`email` 		TEXT UNIQUE NOT NULL,
	`password` 		TEXT NOT NULL,
	`first_name` 	TEXT NOT NULL,
	`last_name` 	TEXT NOT NULL,
	`date_of_birth` TEXT NOT NULL,
	`avatar` 		TEXT,
	`nickname` 		TEXT,
	`about_me` 		TEXT
);
