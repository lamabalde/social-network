package queries

import (
	"database/sql"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

const (
	FOLLOWERS = iota
	FOLLOWINGS
)

func (dbm *DBModel) GetFollowers(userID string) ([]models.UserBase, error) {
	q := `SELECT u.id, u.userName FROM users u
			RIGHT JOIN followers f ON u.id =followerID
			WHERE f.followingID = ? AND followStatus = ?`
	rows, err := dbm.DB.Query(q, userID, models.FOLLOW_STATUS_FOLLOWING)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserBase
	for rows.Next() {
		user := models.UserBase{}
		err := rows.Scan(&user.ID, &user.UserName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (dbm *DBModel) GetFollowing(userID string) ([]models.UserBase, error) {
	q := `SELECT u.id, u.userName FROM users u
	RIGHT JOIN followers f ON u.id =followingID
	WHERE f.followerID = ? AND followStatus = ?`
	rows, err := dbm.DB.Query(q, userID, models.FOLLOW_STATUS_FOLLOWING)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserBase
	for rows.Next() {
		user := models.UserBase{}
		err := rows.Scan(&user.ID, &user.UserName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (dbm *DBModel) GetFollowStatus(followerID, followingID string) (string, error) {
	q := `SELECT followStatus FROM followers
	WHERE followerID = ? AND followingID = ?`
	var followStatus string
	row := dbm.DB.QueryRow(q, followerID, followingID)

	if err := row.Scan(&followStatus); err != nil {
		if err == sql.ErrNoRows {
			return "not following", nil
		}
		return "not following", err
	}
	return followStatus, nil
}

func (dbm *DBModel) GetFollows(userID string) ([]models.UserBase, error) {
	q := `SELECT u.id, u.userName FROM users u
		RIGHT JOIN (SELECT followingID as fl FROM followers
		WHERE followerID = ? AND  followStatus = ? UNION SELECT followerID as fl FROM followers
		WHERE followingID = ? AND  followStatus = ?) on u.id=fl`
	rows, err := dbm.DB.Query(q, userID, models.FOLLOW_STATUS_FOLLOWING, userID, models.FOLLOW_STATUS_FOLLOWING)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.UserBase
	for rows.Next() {
		user := models.UserBase{}
		err := rows.Scan(&user.ID, &user.UserName)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

/*
returns list of users who filtered by userIDs.CheckID().
The list is ordered in descending order by date of chat messages sent from the user in the list to the user 'forUserID',
if there is no message from a user, shuch users will sort by their names.
*/
func (dbm *DBModel) GetFilteredFollowsOrderedByMessagesToGivenUser(userIDs models.IdChecker, forUserID string) ([]*models.User, error) {
	// q := `SELECT u.id, u.userName , max(ms.dateCreate) FROM users u
	// LEFT JOIN (SELECT  mb.id as mbID,  mb.userID as UserID FROM chat_members mb
	// 	WHERE mb.userID!=? AND mb.chatID IN (SELECT chatID FROM chat_members WHERE userID=?))
	// 	userMb ON u.id=userMb.UserID
	// LEFT JOIN chat_messages ms ON userMb.mbID=ms.chat_membersID
	// WHERE u.id!=?
	// GROUP BY u.id  ORDER BY ms.dateCreate desc, lower(userName)`
	q := `SELECT u.id, u.userName , max(ms.dateCreate) 
		FROM users u
		RIGHT JOIN 
			(SELECT followingID as flID FROM followers
			WHERE followerID = ? AND  followStatus = ? UNION SELECT followerID as flID FROM followers
			WHERE followingID = ? AND  followStatus = ?) on u.id=flID

		LEFT JOIN 
			chat_members mb
			ON u.id=mb.userID AND mb.chatID IN (SELECT chatID FROM chat_members WHERE userID=?)
		LEFT JOIN 
			chat_messages ms 
			ON mb.id=ms.chat_membersID
		WHERE u.id!=?
		GROUP BY u.id  ORDER BY ms.dateCreate desc, lower(userName)`

	rows, err := dbm.DB.Query(q, forUserID, models.FOLLOW_STATUS_FOLLOWING, forUserID, models.FOLLOW_STATUS_FOLLOWING, forUserID, forUserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		var dateCreate sql.NullString
		err := rows.Scan(&user.ID, &user.UserName, &dateCreate)
		if err != nil {
			return nil, err
		}
		if userIDs.CheckID(user.ID) {
			user.LastMessageDate = dateCreate.String
			users = append(users, user)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (dbm *DBModel) AddFollowing(followerID, followingID, followStatus string) error {
	q := `INSERT INTO followers (followerID, followingID, followStatus) VALUES (?,?,?)`
	_, err := dbm.DB.Exec(q, followerID, followingID, followStatus)
	if err != nil {
		return err
	}
	return nil
}

func (dbm *DBModel) SetFollowStatus(followerID, followingID, followStatus string) error {
	q := `UPDATE followers SET followStatus=? WHERE followerID = ? AND followingID = ?`

	res, err := dbm.DB.Exec(q, followStatus, followerID, followingID)
	if err != nil {
		return err
	}
	return dbm.checkUnique(res)
}

func (dbm *DBModel) DeleteFollowing(followerID, followingID string) error {
	q := `DELETE FROM followers WHERE followerID=? AND followingID=?`
	res, err := dbm.DB.Exec(q, followerID, followingID)
	if err != nil {
		return err
	}

	return dbm.checkUnique(res)
}
