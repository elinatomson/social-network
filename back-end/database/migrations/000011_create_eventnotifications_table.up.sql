CREATE TABLE IF NOT EXISTS `eventnotifications`(
    `id`                INTEGER PRIMARY KEY AUTOINCREMENT,
	`event_id`		    INTEGER,
    `member_id`	        INTEGER,
    `group_id`	        INTEGER
);
