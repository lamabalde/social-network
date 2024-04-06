INSERT INTO 'chat_members'(id, chatID, userID) VALUES (1, 'chat1', '1');
INSERT INTO 'chat_members'(id, chatID, userID) VALUES (2, 'chat1', '2');
INSERT INTO 'chat_members'(id, chatID, userID) VALUES (3, 'chat2', '1');
INSERT INTO 'chat_members'(id, chatID, userID) VALUES (4, 'chat2', '3');

INSERT INTO 'chat_messages' (id, content, chat_membersID, dateCreate) VALUES ('1', 'hello from 1 to 2', '1', '2023-11-22 10:41:23.656479916+00:00');
INSERT INTO 'chat_messages' (id, content, chat_membersID, dateCreate) VALUES ('2', 'hello from 2 to 1', '2', '2023-11-22 10:55:23.656479916+00:00');
INSERT INTO 'chat_messages' (id, content, chat_membersID, dateCreate) VALUES ('3', 'mes2 from 1 to 2', '1', '2023-11-22 10:56:23.656479916+00:00');
INSERT INTO 'chat_messages' (id, content, chat_membersID, dateCreate) VALUES ('4', 'mes3 from 1 to 2', '1', '2023-11-22 10:57:23.656479916+00:00');
INSERT INTO 'chat_messages' (id, content, chat_membersID, dateCreate) VALUES ('5', 'mes2 from 2 to 1', '2', '2023-11-22 10:58:23.656479916+00:00');
INSERT INTO 'chat_messages' (id, content, chat_membersID, dateCreate) VALUES ('6', 'mes4 from 1 to 2', '1', '2023-11-22 10:58:33.656479916+00:00');
INSERT INTO 'chat_messages' (id, content, chat_membersID, dateCreate) VALUES ('7', 'mes3 from 2 to 1', '2', '2023-11-22 10:59:33.656479916+00:00');
INSERT INTO 'chat_messages' (id, content, chat_membersID, dateCreate) VALUES ('8', 'hello from 1 to 3', '3', '2023-11-23 10:59:33.656479916+00:00');
INSERT INTO 'chat_messages' (id, content, chat_membersID, dateCreate) VALUES ('9', 'mes1 from 1 to 3', '3', '2023-11-24 10:59:33.656479916+00:00');
