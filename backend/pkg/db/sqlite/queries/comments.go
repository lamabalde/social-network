package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

const (
	SelectCommentsForPostIDQuery = `-- select comments.
		SELECT c.id, c.content, c.images, c.userID, u.userName, c.dateCreate, c.postID,
			count(CASE WHEN cl.like THEN TRUE END), count(CASE WHEN NOT cl.like THEN TRUE END), 
			(CASE WHEN c.id IN (SELECT messageID FROM comments_likes cl  WHERE cl.userID = ? AND cl.like=true)  THEN 1
				      WHEN c.id IN (SELECT messageID FROM comments_likes cl  WHERE cl.userID = ? AND cl.like=false) THEN 0
					  ELSE -1 END)
	    FROM comments c
		LEFT JOIN users u ON u.id=c.userID
	    LEFT JOIN comments_likes cl ON cl.messageID=c.id 
		WHERE c.postID = ?		 
		GROUP BY c.id
		ORDER BY c.dateCreate DESC
		LIMIT ? OFFSET ?;
`

	SelectCommentByIDQuery = `SELECT c.id, c.content, c.images, c.userID, u.userName, c.dateCreate, c.postID, 
			count(CASE WHEN cl.like THEN TRUE END), count(CASE WHEN NOT cl.like THEN TRUE END),
			(CASE WHEN c.id IN (SELECT messageID FROM comments_likes cl  WHERE cl.userID = ? AND cl.like=true)  THEN 1
				      WHEN c.id IN (SELECT messageID FROM comments_likes cl  WHERE cl.userID = ? AND cl.like=false) THEN 0
					  ELSE -1 END) 
	    FROM comments c
		LEFT JOIN users u ON u.id=c.userID
	    LEFT JOIN comments_likes cl ON cl.messageID=c.id 
		WHERE c.id = ?		 
		GROUP BY c.id;
		`
)

func (dbm *DBModel) GetCommentsByPostID(postID, userIDForReaction string, limit, offset int) ([]*models.Comment, error) {
	rows, err := dbm.DB.Query(SelectCommentsForPostIDQuery, userIDForReaction, userIDForReaction, postID, limit, offset)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		comment, err := scanRowToComment(rows)
		if err != nil {
			return nil, err
		}
		if len(comment.Content.Images) > 0 {
			comment.Content.Image = "img/comment/" + comment.Content.Images[0]
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

/*
scan and prefilles an item of models.Comment for getCommentsByPostID and getCommentsByID
*/
func scanRowToComment(rows *sql.Rows) (*models.Comment, error) {
	comment := &models.Comment{}
	comment.Content.Likes = make([]int, models.N_REACTIONS)

	var commentID, content, images, userID, userName, postID sql.NullString
	var commentCreate sql.NullTime

	// parse the row with fields:
	// c.id, c.content, c.images, c.userID, u.userName, c.dateCreate, c.postID,
	// count(CASE WHEN cl.like THEN TRUE END), count(CASE WHEN NOT cl.like THEN TRUE END)
	// (CASE WHEN p.id IN (SELECT messageID FROM comments_likes cl  WHERE cl.userID = ? AND cl.like=true)  THEN 1
	// 		      WHEN p.id IN (SELECT messageID FROM comments_likes cl  WHERE cl.userID = ? AND cl.like=false) THEN 0
	// 			  ELSE -1 END)
	err := rows.Scan(&commentID, &content, &images,
		&userID, &userName,
		&commentCreate, &postID,
		&comment.Content.Likes[models.LIKE], &comment.Content.Likes[models.DISLIKE],
		&comment.Content.UserReaction,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}
		return nil, err
	}

	if !commentID.Valid {
		return nil, models.ErrNoRecords
	}

	comment = assembleComment(comment, commentID, content, images, userID, userName, postID, commentCreate)

	return comment, err
}

func assembleComment(comment *models.Comment, commentID, content, images, userID, userName, postID sql.NullString, commentCreate sql.NullTime) *models.Comment {
	comment.ID = commentID.String
	comment.Content.Text = content.String
	comment.Content.Images = helpers.SplitNullString(images)
	comment.Content.UserID = userID.String
	comment.Content.UserName = userName.String
	comment.PostID = postID.String
	comment.Content.DateCreate = commentCreate.Time

	return comment
}

/*
search in the DB a comment by the given ID returns comment and its postID
*/
func (dbm *DBModel) GetCommentByID(id, userIDForReaction string) (*models.Comment, error) {
	rows, err := dbm.DB.Query(SelectCommentByIDQuery, userIDForReaction, userIDForReaction, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}
		return nil, err
	}
	defer rows.Close()

	var comment *models.Comment
	if rows.Next() {
		comment, err = scanRowToComment(rows)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, models.ErrNoRecords
	}

	return comment, nil
}

/*
inserts a new comment into DB, returns an ID for the comment
*/
func (dbm *DBModel) AddComment(postID string, text string, images []string, authorID string) (*models.Comment, error) {
	strOfImages := helpers.JoinToNullString(images)

	commentID, err := helpers.GenerateNewUUID()
	if err != nil {
		return nil, err
	}

	// TODO may be get date from frontend - it could be more precise
	dateCreate := time.Now()

	q := `INSERT INTO comments (id, content, images, userID, dateCreate, postID) VALUES (?,?,?,?,?,?)`
	_, err = dbm.DB.Exec(q, commentID, text, strOfImages, authorID, dateCreate, postID)
	if err != nil {
		return nil, err
	}

	err = dbm.increaseCommentsQuantityForPost(postID)
	if err != nil {
		return nil, err
	}

	content := models.Content{
		UserID:     authorID,
		Text:       text,
		Images:     images,
		DateCreate: dateCreate,
	}
	comment := &models.Comment{
		ID:      commentID,
		Content: content,
		PostID:  postID,
	}
	return comment, nil
}

/*
modify a comment with the given id
*/
func (dbm *DBModel) ModifyComment(id string, content string, images []string) error {
	fields := ""
	fieldsValues := []any{}
	if content != "" {
		fields += "content=?, "
		fieldsValues = append(fieldsValues, content)
	}
	if len(images) != 0 {
		fields += "images=?, "
		fieldsValues = append(fieldsValues, helpers.JoinToNullString(images))
	}
	fields, ok := strings.CutSuffix(fields, ", ")
	if !ok {
		panic("cant cut the , after fields list in func modufyPost")
	}
	fieldsValues = append(fieldsValues, id)

	q := fmt.Sprintf("UPDATE comments SET %s WHERE id=?", fields)
	_, err := dbm.DB.Exec(q, fieldsValues...)
	if err != nil {
		return err
	}

	return nil
}

/*
Delete a comment from the database (this is "hard" delete, we should use "soft" delete instead)
*/
func (dbm *DBModel) DeleteComment(id string) error {
	row := dbm.DB.QueryRow(`SELECT c.postID FROM comments c WHERE c.id = ? `, id)
	var postID string
	err := row.Scan(&postID)
	if err != nil {
		return fmt.Errorf("cant get postID for comment to delete: %w", err)
	}

	_, err = dbm.DB.Exec("DELETE FROM comments WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("comment delete failed: %w", err)
	}

	err = dbm.decreaseCommentsQuantityForPost(postID)
	if err != nil {
		return err
	}

	return nil
}
