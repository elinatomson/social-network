CREATE TABLE IF NOT EXISTS `posts` (
    `post_id`		INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    `user_id` 		INTEGER,
    `content`		TEXT NOT NULL, 
	`first_name` 	TEXT NOT NULL,
	`last_name` 	TEXT NOT NULL,
    `privacy` 	    TEXT NOT NULL,
    `image` 	    TEXT,
    `date` 	        DATETIME
);
