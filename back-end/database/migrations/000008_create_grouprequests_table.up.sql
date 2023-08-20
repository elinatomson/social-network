CREATE TABLE IF NOT EXISTS `grouprequests`(
    `id`                INTEGER PRIMARY KEY AUTOINCREMENT,
	`group_id`		INTEGER,
    `group_title`		TEXT,
    `group_creator_id`		TEXT,
    `requester_id`		INTEGER,
    `request_pending`   BOOLEAN NOT NULL
);
