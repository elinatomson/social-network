CREATE TABLE IF NOT EXISTS `comments` (
    `comment_id`	INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    `post_id` 		INTEGER,
    `user_id` 		INTEGER,
    `comment`		TEXT NOT NULL, 
	`first_name` 	TEXT NOT NULL,
	`last_name` 	TEXT NOT NULL,
    `image` 	    TEXT,
    `date` 	        DATETIME
);
