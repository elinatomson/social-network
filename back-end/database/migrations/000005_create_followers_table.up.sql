CREATE TABLE IF NOT EXISTS `followers`(
    `id`                INTEGER PRIMARY KEY AUTOINCREMENT,
	`follower_id`		INTEGER,
    `following_id`		INTEGER,
    `request_pending`   BOOLEAN NOT NULL
);
