CREATE TEMP TABLE tmp_no AS SELECT  * FROM 'notifications';
DROP TABLE 'notifications';
ALTER TABLE 'notifications' RENAME TO 'tmp';
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

INSERT INTO 'notifications' (target, source, type, body, read) 
  SELECT userID, fromUserID, type, body, read FROM 'tmp_no'; 

DROP TABLE 'tmp_no';