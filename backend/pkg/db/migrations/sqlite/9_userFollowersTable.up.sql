CREATE TABLE followers (
    followerID TEXT NOT NULL REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    followingID TEXT NOT NULL CHECK(followerID!=followingID) REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    followStatus TEXT NOT NULL,
    UNIQUE (followerID, followingID)
);

ALTER TABLE users
DROP COLUMN followers;
