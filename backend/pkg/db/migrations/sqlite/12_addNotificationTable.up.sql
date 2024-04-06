CREATE TABLE 'notifications' (
    target TEXT,
    source TEXT,
    type TEXT,
    body TEXT,
    read INT,
    FOREIGN KEY (target) REFERENCES 'users'(id) 
		ON DELETE CASCADE
		ON UPDATE CASCADE,
    FOREIGN KEY (source) REFERENCES 'users'(id) 
		ON DELETE CASCADE
		ON UPDATE CASCADE
);