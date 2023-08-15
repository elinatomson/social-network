CREATE TABLE IF NOT EXISTS `messages` (
    `messageID`				INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
    `message`			    TEXT,
    `first_name_from`	    TEXT,
    `first_name_to`			TEXT,
    `date`       			DATETIME,
    `read`          		INTEGER DEFAULT 0 
);