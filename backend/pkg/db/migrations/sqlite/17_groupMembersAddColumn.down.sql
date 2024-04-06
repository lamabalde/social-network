DELETE FROM group_members WHERE isMember = false;

ALTER TABLE group_members
DROP COLUMN isMember;