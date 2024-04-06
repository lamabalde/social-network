package liker

import (
	"errors"
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

type LikePost struct {
	ID        int    `json:"id,omitempty"`
	UserID    string `json:"userID"`
	UserName  string `json:"userName,omitempty"`
	MessageID string `json:"messageID"`
	Reaction  bool   `json:"reaction"`
}

type LikeComment struct {
	ID        int    `json:"id,omitempty"`
	UserID    string `json:"userID"`
	UserName  string `json:"userName,omitempty"`
	MessageID string `json:"messageID"`
	Reaction  bool   `json:"reaction"`
}

type Liker interface {
	GetLike(*queries.DBModel) error
	InsertLike(*queries.DBModel, bool) error
	UpdateLike(*queries.DBModel, bool) error
	DeleteLike(*queries.DBModel) error
	CompareLike(bool) bool
	GetLikesNumbers(*queries.DBModel) (LikesNumbers, error)
}

type LikesNumbers struct {
	Likes            int    `json:"likes"`
	Dislikes         int    `json:"dislikes"`
	UserWithReaction string `json:"userWithReaction"`
	UserReaction     int8   `json:"userReactions"`
}

func NewLikeComment(user *models.User, reactData webmodel.Reaction) *LikeComment {
	var lc LikeComment
	lc.UserID = user.ID
	lc.UserName = user.UserName
	lc.MessageID = reactData.MessageID
	lc.Reaction = reactData.Reaction
	return &lc
}

func NewLikePost(userID, userName string, reactData webmodel.Reaction) *LikePost {
	var pc LikePost
	pc.UserID = userID
	pc.UserName = userName
	pc.MessageID = reactData.MessageID
	pc.Reaction = reactData.Reaction
	return &pc
}

func (pl *LikePost) GetLike(db *queries.DBModel) error {
	var err error
	pl.ID, pl.Reaction, err = db.GetUsersPostLike(pl.UserID, pl.MessageID)
	return err
}

func (pl *LikePost) GetLikesNumbers(db *queries.DBModel) (LikesNumbers, error) {
	var likesNum LikesNumbers
	likes, userReaction, err := db.GetPostLikes(pl.MessageID, pl.UserID)
	if err != nil {
		return likesNum, err
	}
	likesNum.Dislikes = likes[models.DISLIKE]
	likesNum.Likes = likes[models.LIKE]
	likesNum.UserReaction = userReaction
	likesNum.UserWithReaction = pl.UserID
	return likesNum, nil
}

func (pl *LikePost) InsertLike(db *queries.DBModel, like bool) error {
	var err error
	pl.Reaction = like
	pl.ID, err = db.InsertPostLike(pl.UserID, pl.MessageID, pl.Reaction)
	return err
}

func (pl *LikePost) UpdateLike(db *queries.DBModel, like bool) error {
	pl.Reaction = like
	return db.UpdatePostLike(pl.ID, pl.Reaction)
}

func (pl *LikePost) DeleteLike(db *queries.DBModel) error {
	return db.DeletePostLike(pl.ID)
}

func (pl *LikePost) CompareLike(like bool) bool {
	return pl.Reaction == like
}

func (cl *LikeComment) GetLike(db *queries.DBModel) error {
	var err error
	cl.ID, cl.Reaction, err = db.GetUsersCommentLike(cl.UserID, cl.MessageID)
	return err
}

func (cl *LikeComment) GetLikesNumbers(db *queries.DBModel) (LikesNumbers, error) {
	var likesNum LikesNumbers
	likes, userReaction, err := db.GetCommentLikes(cl.MessageID, cl.UserID)
	if err != nil {
		return likesNum, err
	}
	likesNum.Dislikes = likes[models.DISLIKE]
	likesNum.Likes = likes[models.LIKE]
	likesNum.UserReaction = userReaction
	likesNum.UserWithReaction = cl.UserID
	return likesNum, nil
}

func (cl *LikeComment) InsertLike(db *queries.DBModel, like bool) error {
	var err error
	cl.Reaction = like
	cl.ID, err = db.InsertCommentLike(cl.UserID, cl.MessageID, cl.Reaction)
	return err
}

func (cl *LikeComment) UpdateLike(db *queries.DBModel, like bool) error {
	cl.Reaction = like
	return db.UpdateCommentLike(cl.ID, cl.Reaction)
}

func (cl *LikeComment) DeleteLike(db *queries.DBModel) error {
	return db.DeleteCommentLike(cl.ID)
}

func (cl *LikeComment) CompareLike(like bool) bool {
	return cl.Reaction == like
}

func SetLike(db *queries.DBModel, liker Liker, newLike bool) error {
	err := liker.GetLike(db)
	if err != nil {
		// if there is no like/dislike made by the user, add a new one
		if errors.Is(err, models.ErrNoRecords) {
			err := liker.InsertLike(db, newLike)
			if err != nil {
				return fmt.Errorf("insert data to DB failed: %s", err)
			}
		} else {
			return fmt.Errorf("getting data from DB failed: %s", err)
		}
	} else {
		if liker.CompareLike(newLike) { // if it is the same like, delete it
			err := liker.DeleteLike(db)
			if err != nil {
				return fmt.Errorf("deleting data from DB failed: %s", err)
			}
		} else {
			err := liker.UpdateLike(db, newLike)
			if err != nil {
				return fmt.Errorf("updating data in DB failed: %s", err)
			}
		}
	}
	return nil
}
