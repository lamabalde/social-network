CREATE TEMP TABLE tmp_ge AS SELECT  * FROM 'group_event';
DROP TABLE 'group_event';
CREATE TABLE  'group_event' (
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		dateCreate TIMESTAMP NOT NULL,
		dateEvent TIMESTAMP NOT NULL,
		groupID INTEGER NOT NULL,
		UNIQUE (title, dateEvent),
		FOREIGN KEY (groupID) REFERENCES 'groups'(id) 
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);

INSERT INTO 'group_event' (id, title, description, dateCreate, dateEvent, groupID) 
  SELECT tmp_ge.id, tmp_ge.title, tmp_ge.description, tmp_ge.dateCreate, tmp_ge.dateEvent, gm.groupID 
    FROM 'tmp_ge' LEFT JOIN 'group_members' gm ON tmp_ge.createByID = gm.id;

DROP TABLE 'tmp_ge';
