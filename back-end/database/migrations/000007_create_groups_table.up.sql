CREATE TABLE IF NOT EXISTS `groups` (
    `group_id`		INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    `title`		    TEXT NOT NULL, 
    `description`	TEXT NOT NULL, 
    `user_id` 		INTEGER,
	`first_name` 	TEXT NOT NULL,
	`last_name` 	TEXT NOT NULL,
    `selected_user_id` 	    TEXT
);
