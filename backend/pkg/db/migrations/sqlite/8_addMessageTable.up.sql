CREATE TABLE messages (
messageID TEXT PRIMARY KEY,
messageBody TEXT,
sender TEXT REFERENCES 'users'(id),
receiver TEXT REFERENCES 'users'(id),
groupID TEXT REFERENCES 'groups'(id),
timeStamp TEXT);
