CREATE TABLE IF NOT EXISTS `sessions` (
    `user_id` 			INTEGER,
    `email` 			TEXT NOT NULL,
    `first_name` 		TEXT NOT NULL,
    `last_name` 		TEXT NOT NULL,
    `cookie`			TEXT NOT NULL
);