CREATE TABLE IF NOT EXISTS `events`(
    `event_id`          INTEGER PRIMARY KEY AUTOINCREMENT,
    `title`		        TEXT,
    `description`		TEXT,
    `user_id` 		    INTEGER,
	`first_name` 	    TEXT NOT NULL,
	`last_name` 	    TEXT NOT NULL,
    `time`		        TEXT NOT NULL,
    `group_id`		    INTEGER
);
