package queries

import (
	"database/sql"
	"errors"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

const (
	POSTS_LIKES    = "posts_likes"
	COMMENTS_LIKES = "comments_likes"
)

/****
the group of function for getting likes
****/

/* returns quantity of likes/dislikes from the given table (posts or comments) for the given id of a message*/
func (dbm *DBModel) GetLikes(tableName string, messageID string, userIDForReaction string) ([]int, int8, error) {
	likes := []int{0, 0}
	q := `SELECT  count(CASE WHEN like THEN TRUE END) AS likes, count(CASE WHEN NOT like THEN TRUE END)  AS dislikes,   
	(SELECT like FROM posts_likes WHERE userID = ? AND messageID = ?) AS user_like 
	FROM ` + tableName + ` WHERE messageID=? `
	row := dbm.DB.QueryRow(q, userIDForReaction, messageID, messageID)

	var userLike sql.NullBool
	err := row.Scan(&likes[models.LIKE], &likes[models.DISLIKE], &userLike)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, -1, models.ErrNoRecords
		}
		return nil, -1, err
	}

	var usersReaction int8
	if userLike.Valid {
		if userLike.Bool {
			usersReaction = int8(models.LIKE)
		} else {
			usersReaction = int8(models.DISLIKE)
		}
	} else {
		usersReaction = -1
	}

	return likes, usersReaction, nil
}

func (dbm *DBModel) GetPostLikes(messageID string, userIDForReaction string) ([]int, int8, error) {
	return dbm.GetLikes(POSTS_LIKES, messageID, userIDForReaction)
}

func (dbm *DBModel) GetCommentLikes(messageID string, userIDForReaction string) ([]int, int8, error) {
	return dbm.GetLikes(COMMENTS_LIKES, messageID, userIDForReaction)
}

/* returns quantity of likes/dislikes from the given table (posts or comments) for the given user and message*/
func (dbm *DBModel) getUsersLike(tableName string, userID, messageID string) (int, bool, error) {
	var id int
	var like bool
	q := `SELECT id,like FROM ` + tableName + ` WHERE userID=? AND messageID=?`
	row := dbm.DB.QueryRow(q, userID, messageID)

	err := row.Scan(&id, &like)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, false, models.ErrNoRecords
		}
		return 0, false, err
	}

	return id, like, nil
}

func (dbm *DBModel) GetUsersPostLike(userID, messageID string) (int, bool, error) {
	return dbm.getUsersLike(POSTS_LIKES, userID, messageID)
}

func (dbm *DBModel) GetUsersCommentLike(userID, messageID string) (int, bool, error) {
	return dbm.getUsersLike(COMMENTS_LIKES, userID, messageID)
}

/****
the group of function for changing likes (inser, update, delete)
****/
/*inserts a like/dislike to the given table.*/
func (dbm *DBModel) insertLike(tableName string, userID, messageID string, like bool) (int, error) {
	q := `INSERT INTO ` + tableName + ` (userID, messageID, like) VALUES (?,?,?)`
	res, err := dbm.DB.Exec(q, userID, messageID, like)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

/*sets a new value of like/dislike in the given table.*/
func (dbm *DBModel) updateLike(tableName string, id int, like bool) error {
	q := `UPDATE ` + tableName + ` SET like=? WHERE id=?`
	res, err := dbm.DB.Exec(q, like, id)
	if err != nil {
		return err
	}

	return dbm.checkUnique(res)
}

/*deletes a row from the given table.*/
func (dbm *DBModel) deleteLike(tableName string, id int) error {
	q := `DELETE FROM ` + tableName + ` WHERE id=?`
	res, err := dbm.DB.Exec(q, id)
	if err != nil {
		return err
	}

	return dbm.checkUnique(res)
}

func (dbm *DBModel) InsertPostLike(userID, messageID string, like bool) (int, error) {
	return dbm.insertLike(POSTS_LIKES, userID, messageID, like)
}

func (dbm *DBModel) UpdatePostLike(id int, like bool) error {
	return dbm.updateLike(POSTS_LIKES, id, like)
}

func (dbm *DBModel) DeletePostLike(id int) error {
	return dbm.deleteLike(POSTS_LIKES, id)
}

func (dbm *DBModel) InsertCommentLike(userID, messageID string, like bool) (int, error) {
	return dbm.insertLike(COMMENTS_LIKES, userID, messageID, like)
}

func (dbm *DBModel) UpdateCommentLike(id int, like bool) error {
	return dbm.updateLike(COMMENTS_LIKES, id, like)
}

func (dbm *DBModel) DeleteCommentLike(id int) error {
	return dbm.deleteLike(COMMENTS_LIKES, id)
}

/*deletes a row from the given table by message ID.*/
func (dbm *DBModel) deleteLikeByMessageID(tableName string, messageID string) error {
	q := `DELETE FROM ` + tableName + ` WHERE messageID=?`
	_, err := dbm.DB.Exec(q, messageID)
	if err != nil {
		return err
	}

	return nil
}

func (dbm *DBModel) DeletePostLikeByMessageID(messageID string) error {
	return dbm.deleteLikeByMessageID(POSTS_LIKES, messageID)
}

func (dbm *DBModel) DeleteCommentLikeByMessageID(messageID string) error {
	return dbm.deleteLikeByMessageID(COMMENTS_LIKES, messageID)
}
