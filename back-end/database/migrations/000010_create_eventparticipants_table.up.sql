CREATE TABLE IF NOT EXISTS `eventparticipants`(
    `id`                INTEGER PRIMARY KEY AUTOINCREMENT,
	`event_id`		    INTEGER,
    `participant_id`	INTEGER,
    `first_name` 	    TEXT NOT NULL,
	`last_name` 	    TEXT NOT NULL,
    `going`             BOOLEAN NOT NULL
);
