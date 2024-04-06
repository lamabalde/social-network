package queries

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

const (
	SelectPostByIDQuery = `SELECT p.id, p.theme, p.content, p.images, p.category, p.userID, u.userName,  
	p.dateCreate, p.commentsQuantity, 
	p.groupID, p.postType,
count(CASE WHEN pl.like THEN TRUE END), count(CASE WHEN NOT pl.like THEN TRUE END), 
(CASE WHEN p.id IN (SELECT messageID FROM posts_likes pl  WHERE pl.userID = ? AND pl.like=true)  THEN 1
	 WHEN p.id IN (SELECT messageID FROM posts_likes pl  WHERE pl.userID = ? AND pl.like=false) THEN 0
	 ELSE -1 END)
FROM posts p
LEFT JOIN users u ON u.id=p.userID
LEFT JOIN posts_likes pl ON pl.messageID=p.id 
WHERE p.id = ?;
`

	// 	SelectPostsByConditionQuery = `SELECT p.id, p.theme, p.content, p.images, p.category, p.userID, u.userName, p.dateCreate, p.commentsQuantity, p.groupID,
	// count(CASE WHEN pl.like THEN TRUE END), count(CASE WHEN NOT pl.like THEN TRUE END),
	// (CASE WHEN ul.like is NULL THEN -1 WHEN ul.like THEN 1 WHEN NOT  ul.like THEN 0 END)

	// FROM posts p
	// LEFT JOIN users u ON u.id=p.userID
	// LEFT JOIN posts_likes pl ON pl.messageID=p.id
	// LEFT JOIN (SELECT messageID, like FROM posts_likes pl  WHERE pl.userID = ?) ul ON ul.messageID=p.id
	// `
	SelectPostsByConditionQuery = `SELECT p.id, p.theme, p.content, p.images, p.category, p.userID, u.userName, 
	p.dateCreate, p.commentsQuantity,
	p.groupID,  
	CASE WHEN p.userID=? THEN p.postType
		 WHEN p.postType=0 THEN 0 
		 WHEN p.postType=1 THEN 
		 	(SELECT 1 FROM followers 
			WHERE ((followerID=? AND followingID=p.userID) OR (followerID=p.userID AND followingID=?)) AND followStatus='following')
		 WHEN p.postType=2 THEN 
			(SELECT 2 FROM friends 
			WHERE mainUser=p.userID AND friendUser=?)
	END AS pType,
	count(CASE WHEN pl.like THEN TRUE END), count(CASE WHEN NOT pl.like THEN TRUE END),
	(CASE WHEN ul.like is NULL THEN -1 WHEN ul.like THEN 1 WHEN NOT  ul.like THEN 0 END)
	FROM posts p
	LEFT JOIN users u ON u.id=p.userID
	LEFT JOIN posts_likes pl ON pl.messageID=p.id 
	LEFT JOIN (SELECT messageID, like FROM posts_likes pl  WHERE pl.userID = ?) ul ON ul.messageID=p.id 
`
)

/*
inserts a new post into DB, returns an ID for the post
*/
func (dbm *DBModel) AddPost(post *models.Post) (*models.Post, error) {
	strOfCategories := helpers.JoinToNullString(post.Categories)

	strOfImages := helpers.JoinToNullString(post.Content.Images)

	groupID := helpers.StrToNullString(post.GroupID)

	postID, err := helpers.GenerateNewUUID()
	if err != nil {
		return nil, err
	}
	// TODO may be get date from frontend - it could be more precise

	q := `INSERT INTO posts (id, theme, content, images, category, postType, userID, dateCreate, commentsQuantity, groupID) VALUES (?,?,?,?,?,?,?,?,?,?)`
	_, err = dbm.DB.Exec(q, postID, post.Theme, post.Content.Text, strOfImages, strOfCategories, post.Privacy, post.Content.UserID, post.Content.DateCreate, 0, groupID)
	if err != nil {
		return nil, err
	}

	post.ID = postID
	post.Comments = nil
	post.CommentsQuantity = 0

	return post, nil
}

/*
search in the DB a post by the given ID
*/
func (dbm *DBModel) GetPostByID(id string, userIDForReaction string) (*models.Post, error) {
	// use Query instead of QueryRow to get using scanRowToPost function
	rows, err := dbm.DB.Query(SelectPostByIDQuery, userIDForReaction, userIDForReaction, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var post *models.Post
	if rows.Next() {
		post, err = scanRowToPost(rows)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, models.ErrNoRecords
	}

	if post.CommentsQuantity > 0 {
		post.Comments, err = dbm.GetCommentsByPostID(id, userIDForReaction, post.CommentsQuantity, 0)
		if err != nil {
			return nil, err
		}
	}

	return post, nil
}

/*
returns 'limit' messages of the chat with the given 'id', skips 'offset' last messages
*/
/*
returns 'limit' posts from DB, skipping 'offset' last posts, and matching the filter.
The parametr userID is used to mark the user's reaction to a post.
*/
func (dbm *DBModel) GetFiltredPosts(filter models.Filter, userIDForReaction string, limit, offset int) ([]*models.Post, error) {
	condition := ""
	var arguments []any
	v := reflect.ValueOf(filter)
	for _, field := range reflect.VisibleFields(reflect.TypeOf(filter)) {
		// if either of the fields !=0 add conditions to the query
		if !v.FieldByIndex(field.Index).IsZero() {
			condition = ` WHERE `
			arguments = []any{}
			if filter.GroupID != "" {
				condition += ` p.groupID= ? AND `
				arguments = append(arguments, filter.GroupID)
			}

			if filter.AuthorID != "" {
				condition += ` p.userID= ? AND `
				arguments = append(arguments, filter.AuthorID)
			}

			if filter.LikedByUserID != "" {
				condition += ` p.id IN (SELECT messageID FROM posts_likes pl  WHERE pl.userID = ? AND pl.like=true) AND `
				arguments = append(arguments, filter.LikedByUserID)
			}

			if filter.DisLikedByUserID != "" {
				condition += ` p.id IN (SELECT messageID FROM posts_likes pl  WHERE pl.userID = ? AND pl.like=false) AND `
				arguments = append(arguments, filter.DisLikedByUserID)
			}

			condition = strings.TrimSuffix(condition, `AND `)
			break
		}
	}

	return dbm.GetPostsByCondition(condition, arguments, userIDForReaction, limit, offset)
}

/*
returns posts that have got the given category
*/
func (dbm *DBModel) GetPostsLikedByUser(userID string, limit, offset int) ([]*models.Post, error) {
	condition := ` WHERE p.id IN (SELECT messageID FROM posts_likes pl  WHERE pl.userID = ? AND pl.like=true) `
	arguments := []any{userID}

	return dbm.GetPostsByCondition(condition, arguments, userID, limit, offset)
}

/*
returns posts belonging to no group
*/
func (dbm *DBModel) GetPostsNoGroup(userIDForReaction string, limit, offset int) ([]*models.Post, error) {
	// TODO: get list of close friends, then also get posts by them, ignore other 'friends' posts who are not in the list of 'friends' (postType==2)
	condition := ` WHERE p.groupID is NULL AND pType is NOT NULL` //see pType in the const SelectPostsByConditionQuery, it depends on postType, userIDForReaction, follows and friends tables
	arguments := []any{}
	return dbm.GetPostsByCondition(condition, arguments, userIDForReaction, limit, offset)
}

/*
returns posts belonging to no group
*/
func (dbm *DBModel) GetPostsInGroup(groupID string, userIDForReaction string, limit, offset int) ([]*models.Post, error) {
	condition := ` WHERE p.groupID=? `
	arguments := []any{groupID}

	return dbm.GetPostsByCondition(condition, arguments, userIDForReaction, limit, offset)
}

func (dbm *DBModel) GetUserPosts(userID string, userIDForReaction string, limit, offset int) ([]*models.Post, error) {
	condition := ` WHERE p.userID=? AND pType is NOT NULL `
	arguments := []any{userID}

	return dbm.GetPostsByCondition(condition, arguments, userIDForReaction, limit, offset)
}

/*
addes the condition to a query and run it. Returnes found posts
*/
func (dbm *DBModel) GetPostsByCondition(condition string, argumentsForCondition []any, userIDForReaction string, limit, offset int) ([]*models.Post, error) {
	query := SelectPostsByConditionQuery + condition +
		` GROUP BY p.id 
		ORDER BY p.dateCreate DESC, p.id
		LIMIT ? OFFSET ?
		`
	// exequting the query
	var rows *sql.Rows
	var err error
	arguments := append([]any{userIDForReaction,userIDForReaction, userIDForReaction, userIDForReaction, userIDForReaction}, argumentsForCondition...)
	arguments = append(arguments, limit, offset)
	rows, err = dbm.DB.Query(query, arguments...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// parsing the query's result
	var posts []*models.Post

	// add the first post without condition
	for rows.Next() {
		post, err := scanRowToPost(rows)
		if err != nil {
			return nil, err
		}
		post.Image = helpers.GetPostImgUrl(post.ID)
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

/*
scans and prefilles an item of modelPost for getPosts and getPostById
*/
func scanRowToPost(rows *sql.Rows) (*models.Post, error) {
	post := &models.Post{}
	post.Content.Likes = make([]int, models.N_REACTIONS)

	var postID, theme, content, images, category, userID, userName, groupID sql.NullString
	var commentsQuantity, privacy sql.NullInt64
	var postCreate sql.NullTime

	// parse the row with fields:
	// p.id, p.theme, p.content, p.images, p.category, p.userID, u.userName,
	// p.dateCreate, p.commentsQuantity,
	// p.groupID,  p.postType,
	// count(CASE WHEN pl.like THEN TRUE END), count(CASE WHEN NOT pl.like THEN TRUE END),
	// (CASE WHEN ul.like is NULL THEN -1 WHEN ul.like THEN 1 WHEN NOT  ul.like THEN 0 END)
	err := rows.Scan(&postID, &theme, &content, &images, &category,
		&userID, &userName,
		&postCreate, &commentsQuantity,
		&groupID, &privacy, 
		&post.Content.Likes[models.LIKE], &post.Content.Likes[models.DISLIKE],
		&post.Content.UserReaction,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecords
		}
		return nil, err
	}

	if !postID.Valid {
		return nil, models.ErrNoRecords
	}

	post = assemblePost(post, postID, theme, content, images, category, userID, userName, groupID, privacy, commentsQuantity, postCreate)

	return post, err
}

func assemblePost(post *models.Post, postID, theme, content, images, category, userID, userName, groupID sql.NullString,
	privacy, commentsQuantity sql.NullInt64, postCreate sql.NullTime,
) *models.Post {
	post.ID = postID.String
	post.Theme = theme.String
	post.Content.Text = content.String
	post.Content.Images = helpers.SplitNullString(images)
	post.Categories = helpers.SplitNullString(category)
	post.Content.UserID = userID.String
	post.Content.UserName = userName.String
	post.CommentsQuantity = int(commentsQuantity.Int64)
	post.Content.DateCreate = postCreate.Time
	post.GroupID = groupID.String
	post.Privacy = int(privacy.Int64)


	return post
}

/*
modify a post with the given id
*/
func (dbm *DBModel) ModifyPost(id string, theme, content string, images []string) error {
	fields := ""
	fieldsValues := []any{}
	if theme != "" {
		fields += "theme=?, "
		fieldsValues = append(fieldsValues, theme)
	}
	if content != "" {
		fields += "content=?, "
		fieldsValues = append(fieldsValues, content)
	}
	if len(images) != 0 {
		fields += "images=?, "
		fieldsValues = append(fieldsValues, helpers.JoinToNullString(images)) // TODO check when img added to a post with other images (seems it'll del old records)
	}
	fields, ok := strings.CutSuffix(fields, ", ")
	if !ok {
		panic("cant cut the , after fields list in func modufyPost")
	}
	fieldsValues = append(fieldsValues, id)

	q := fmt.Sprintf("UPDATE posts SET %s WHERE id=?", fields)
	_, err := dbm.DB.Exec(q, fieldsValues...)
	if err != nil {
		return err
	}

	return nil
}

/*
increace the post's commentsQuantity
*/
func (dbm *DBModel) increaseCommentsQuantityForPost(id string) error {
	_, err := dbm.DB.Exec(`UPDATE posts SET commentsQuantity=commentsQuantity+1 WHERE id=?`, id)
	if err != nil {
		return err
	}

	return nil
}

/*
decreace the post's commentsQuantity
*/
func (dbm *DBModel) decreaseCommentsQuantityForPost(id string) error {
	_, err := dbm.DB.Exec(`UPDATE posts SET commentsQuantity=commentsQuantity-1 WHERE id=?`, id)
	if err != nil {
		return err
	}

	return nil
}

/*
Delete a comment from the database (this is "hard" delete, we should use "soft" delete instead)
*/
func (dbm *DBModel) DeletePost(id string) error {
	res, err := dbm.DB.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("post delete failed: %w", err)
	}

	rowsNumber, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("post delete failed: %w", err)
	}
	if rowsNumber != 1 {
		return fmt.Errorf("post delete failed: deleted %d rows, expected 1", rowsNumber)
	}

	return nil
}
