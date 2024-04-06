ALTER TABLE group_members
ADD isMember BOOL NOT NULL DEFAULT false;

UPDATE group_members SET isMember = true;