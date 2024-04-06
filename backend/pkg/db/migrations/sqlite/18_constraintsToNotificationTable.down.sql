CREATE TEMP TABLE tmp_no AS SELECT  * FROM 'notifications';
DROP TABLE 'notifications';
CREATE TABLE 'notifications' (
    id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    userID TEXT NOT NULL,
    fromUserID TEXT,
    type TEXT NOT NULL,
    body TEXT NOT NULL,
    groupID TEXT,
    postID TEXT,
    read INT NOT NULL,
    dateCreate TIMESTAMP NOT NULL,
    FOREIGN KEY (userID) REFERENCES 'users'(id) 
		ON DELETE CASCADE
		ON UPDATE CASCADE,
    FOREIGN KEY (fromUserID) REFERENCES 'users'(id) 
		ON DELETE CASCADE
		ON UPDATE CASCADE,
    FOREIGN KEY (groupID) REFERENCES 'groups'(id) 
		ON DELETE CASCADE
		ON UPDATE CASCADE,
    FOREIGN KEY (postID) REFERENCES 'posts'(id) 
		ON DELETE CASCADE
		ON UPDATE CASCADE
);

INSERT INTO 'notifications'  
  SELECT *
  FROM 'tmp_no'; 

DROP TABLE 'tmp_no';