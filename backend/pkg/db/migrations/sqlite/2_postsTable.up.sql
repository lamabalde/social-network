CREATE TABLE posts (
    id TEXT PRIMARY KEY NOT NULL,
	theme TEXT NOT NULL DEFAULT ('(No theme)'),
	content TEXT NOT NULL,
	images TEXT, 
	category text,
	postType INT NOT NULL DEFAULT 0,
	userID TEXT NOT NULL,
	dateCreate TIMESTAMP NOT NULL,
	commentsQuantity INT NOT NULL,
	FOREIGN KEY (userID) REFERENCES users(id) 
    	ON UPDATE CASCADE
    	ON DELETE CASCADE
);
