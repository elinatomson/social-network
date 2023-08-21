CREATE TABLE IF NOT EXISTS `groupmembers`(
    `id`                INTEGER PRIMARY KEY AUTOINCREMENT,
	`group_id`		INTEGER,
    `group_title`		TEXT,
    `group_creator_id`		INTEGER,
    `requester_id`		INTEGER,
    `request_pending`   BOOLEAN NOT NULL,
    `invitation_pending`   BOOLEAN NOT NULL
);
