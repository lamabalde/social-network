INSERT INTO 'groups' (id, title, description, creatorID, dateCreate) VALUES ('1', 'group1', 'group1 by user 2, id 1', '2', '2023-11-23 10:55:23.656479916+00:00');
INSERT INTO 'groups' (id, title, description, creatorID, dateCreate) VALUES ('2', 'group2', 'group2 by user 2, id 2', '2', '2023-11-24 10:55:23.656479916+00:00');
INSERT INTO 'groups' (id, title, description, creatorID, dateCreate) VALUES ('3', 'group3', 'group3 by user 1, id 3', '1', '2023-11-25 10:55:23.656479916+00:00');

INSERT INTO 'group_members' (groupID, userID, isMember) VALUES ('1', '2', true);
INSERT INTO 'group_members' (groupID, userID, isMember) VALUES ('1', '3', true);
INSERT INTO 'group_members' (groupID, userID, isMember) VALUES ('2', '2', true);
INSERT INTO 'group_members' (groupID, userID, isMember) VALUES ('2', '1', true);
INSERT INTO 'group_members' (groupID, userID, isMember) VALUES ('2', '3', true);
INSERT INTO 'group_members' (groupID, userID, isMember) VALUES ('3', '1', true);

INSERT INTO 'group_chat' (id, content, group_membersID, dateCreate) VALUES ('1', 'hello from 2 in 1', 1, '2023-11-22 10:05:23.656479916+00:00');
INSERT INTO 'group_chat' (id, content, group_membersID, dateCreate) VALUES ('2', 'hello from 3 in 1', 2, '2023-11-22 10:15:23.656479916+00:00');
INSERT INTO 'group_chat' (id, content, group_membersID, dateCreate) VALUES ('3', 'mess2 from 2 in 1', 1, '2023-11-22 10:25:23.656479916+00:00');
INSERT INTO 'group_chat' (id, content, group_membersID, dateCreate) VALUES ('4', 'mess3 from 2 in 1', 1, '2023-11-22 10:35:23.656479916+00:00');
INSERT INTO 'group_chat' (id, content, group_membersID, dateCreate) VALUES ('5', 'mess1 from 2 in 2', 3, '2023-11-22 10:45:23.656479916+00:00');
INSERT INTO 'group_chat' (id, content, group_membersID, dateCreate) VALUES ('6', 'mess1 from 1 in 2', 4, '2023-11-22 10:55:33.656479916+00:00');
INSERT INTO 'group_chat' (id, content, group_membersID, dateCreate) VALUES ('7', 'mess1 from 3 in 1', 2, '2023-11-22 10:59:33.656479916+00:00');

INSERT INTO 'group_event' (id, title, description, dateCreate, dateEvent, createByID) VALUES ('1', 'event1gr2', 'event1 group 2, id 1', '2023-11-23 10:55:23.656479916+00:00', '2023-12-23 10:55:23.656479916+00:00', 3);
INSERT INTO 'group_event' (id, title, description, dateCreate, dateEvent, createByID) VALUES ('2', 'event2gr2', 'event2 group 2, id 2', '2023-11-23 10:56:23.656479916+00:00', '2023-12-10 10:55:23.656479916+00:00', 3);
INSERT INTO 'group_event' (id, title, description, dateCreate, dateEvent, createByID) VALUES ('3', 'event1gr1', 'event1 group 1, id 3', '2023-11-23 10:57:23.656479916+00:00', '2023-12-20 10:55:23.656479916+00:00', 1);

INSERT INTO 'group_event_members' (group_eventID, group_membersID, mark) VALUES ('1', 5, 0);
INSERT INTO 'group_event_members' (group_eventID, group_membersID, mark) VALUES ('1', 3, 1);
INSERT INTO 'group_event_members' (group_eventID, group_membersID, mark) VALUES ('1', 4, 2);
INSERT INTO 'group_event_members' (group_eventID, group_membersID, mark) VALUES ('2', 4, 0);
INSERT INTO 'group_event_members' (group_eventID, group_membersID, mark) VALUES ('2', 5, 2);
INSERT INTO 'group_event_members' (group_eventID, group_membersID, mark) VALUES ('3', 2, 0);
