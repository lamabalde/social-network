	CREATE TABLE 'groups'(
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		creatorID TEXT NOT NULL,
		dateCreate TIMESTAMP NOT NULL,
		UNIQUE (title)
	);
	CREATE TABLE  'group_members' (
		id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		groupID TEXT NOT NULL,
		userID TEXT NOT NULL,
		UNIQUE (groupID, userID),
		FOREIGN KEY (groupID) REFERENCES 'groups'(id) 
			ON DELETE CASCADE
			ON UPDATE CASCADE,
		FOREIGN KEY (userID) REFERENCES 'users'(id) 
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
	CREATE TABLE  'group_chat' (
		id TEXT PRIMARY KEY NOT NULL,
		content TEXT NOT NULL,  
		images TEXT, 
		group_membersID INTEGER NOT NULL,
		dateCreate TIMESTAMP NOT NULL,
		FOREIGN KEY (group_membersID) REFERENCES 'group_members'(id) 
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
	CREATE TABLE  'group_event' (
		id TEXT PRIMARY KEY NOT NULL,
		title TEXT NOT NULL,
		description TEXT NOT NULL,
		dateCreate TIMESTAMP NOT NULL,
		dateEvent TIMESTAMP NOT NULL,
		groupID TEXT NOT NULL,
		UNIQUE (title, dateEvent),
		FOREIGN KEY (groupID) REFERENCES 'groups'(id) 
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
	CREATE TABLE  'group_event_members' (
		id INTEGER PRIMARY KEY NOT NULL,
		group_eventID TEXT NOT NULL,
		group_membersID INTEGER NOT NULL,
		mark INTEGER NOT NULL, 
		UNIQUE (group_eventID, group_membersID),
		FOREIGN KEY (group_eventID) REFERENCES 'group_event'(id) 
			ON DELETE CASCADE
			ON UPDATE CASCADE,
		FOREIGN KEY (group_membersID) REFERENCES 'group_members'(id) 
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
			