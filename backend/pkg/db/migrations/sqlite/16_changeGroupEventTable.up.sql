CREATE TEMP TABLE tmp_ge AS SELECT  * FROM 'group_event';
DROP TABLE 'group_event';
CREATE TABLE  'group_event' (
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		dateCreate TIMESTAMP NOT NULL,
		dateEvent TIMESTAMP NOT NULL,
		createByID INTEGER NOT NULL,
		UNIQUE (title, dateEvent),
		FOREIGN KEY (createByID) REFERENCES 'group_members'(id) 
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);

INSERT INTO 'group_event' (id, title, description, dateCreate, dateEvent, createByID) 
  SELECT tmp_ge.id, tmp_ge.title, tmp_ge.description, tmp_ge.dateCreate, tmp_ge.dateEvent, gm.id 
    FROM 'tmp_ge' LEFT JOIN 'groups' gr ON tmp_ge.groupID = gr.id
    LEFT JOIN 'group_members' gm ON tmp_ge.groupID = gm.groupID WHERE gm.userID = gr.creatorID;

DROP TABLE 'tmp_ge';
