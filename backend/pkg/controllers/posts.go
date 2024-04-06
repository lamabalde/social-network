package controllers

import (
	"fmt"
	"os"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

func ReplyDeletePost(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		postID, err := parse.PayloadToString(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for the portion of messages '%s'", message.Payload), err)
		}

		if webmodel.IsEmpty(postID) {
			return nil, currConnection.WSBadRequest(message, "postID is empty")
		}

		err = app.DBModel.DeletePost(postID)
		if err != nil {
			return nil, currConnection.WSError("delete posts from DB failed", err)
		}
		err = deletePostImg(postID)
		if err != nil {
			if !os.IsNotExist(err) { // if post doesn't have an image, ignore error
				return nil, currConnection.WSError("delete post image failed", err)
			}
		}
		return postID, nil
	}
}

func ReplyPosts(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		offset, err := parse.PayloadToInt(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for the portion of messages '%s'", message.Payload), err)
		}

		posts, err := app.DBModel.GetPostsNoGroup(currConnection.Client.UserID, webmodel.POSTS_ON_POSTSVIEW, offset)
		if err != nil {
			return nil, currConnection.WSError("getting posts from DB failed", err)
		}

		createPostsPreview(posts)
		return posts, nil
	}
}

func ReplyPostsInGroup(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		postsInGroupPortion, err := parse.PayloadToPostsInGroupPortion(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for the portion of messages '%s'", message.Payload), err)
		}

		errmessage := postsInGroupPortion.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		posts, err := app.DBModel.GetPostsInGroup(postsInGroupPortion.GroupID, currConnection.Client.UserID, webmodel.POSTS_ON_POSTSVIEW, postsInGroupPortion.Offset)
		if err != nil {
			return nil, currConnection.WSError("getting posts from DB failed", err)
		}

		createPostsPreview(posts)
		return posts, nil
	}
}

func ReplyUserPostsPortion(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		userPostsPortion, err := parse.PayloadToUserPostsPortion(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for the portion of messages '%s'", message.Payload), err)
		}

		errmessage := userPostsPortion.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		posts, err := app.DBModel.GetUserPosts(userPostsPortion.UserID, currConnection.Client.UserID, webmodel.POSTS_ON_POSTSVIEW, userPostsPortion.Offset)
		if err != nil {
			return nil, currConnection.WSError("getting posts from DB failed", err)
		}

		createPostsPreview(posts)
		return posts, nil
	}
}

func deletePostImg(postID string) error {
	imgUrl := helpers.FindFile("data/img/post/", postID+".*")
	err := os.Remove(imgUrl)
	if err != nil {
		return err
	}
	return nil
}
