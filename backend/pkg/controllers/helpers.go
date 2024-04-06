package controllers

import (
	"errors"
	"fmt"
	"strings"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

// All functions of this package send an error message to the websocket connection given in 'currConnection' parameter,
// so it is not necessary to send an error message using these functions in a function of this package.

/*
gets a post from DB by its ID
*/
func getPost(app *application.Application, currConnection *wsconnection.UsersConnection, postId string, message webmodel.WSMessage) (*models.Post, error) {
	post, err := app.DBModel.GetPostByID(postId, currConnection.Client.UserID)
	if err != nil {
		return nil, handleErrGetFromDBByID(app, currConnection, message, "post", postId, err)
	}

	return post, nil
}

func createPostsPreview(posts []*models.Post) {
	for _, post := range posts {
		if len(post.Content.Text) < POST_PREVIEW_LENGTH {
			continue
		}
		for i := POST_PREVIEW_LENGTH - 10; i < len(post.Content.Text); i++ { // find first space after 440 char
			if string(post.Content.Text[i]) == " " {
				post.Content.Text = post.Content.Text[0:i] + "..."
				break
			}
		}
	}
}

func handleErrAddToDB(app *application.Application, uc *wsconnection.UsersConnection, requestMessage webmodel.WSMessage, whatNotSaved string, err error) error {
	errmessage := fmt.Sprintf("save %s in DB failed", whatNotSaved)
	if strings.Contains(err.Error(), "constraint") || errors.Is(err, models.ErrNoRecords) {
		return uc.WSBadRequest(requestMessage, fmt.Sprintf("%s: %v", errmessage, err))
	}
	return uc.WSError(errmessage, err)
}

func handleErrGetFromDBByID(app *application.Application, uc *wsconnection.UsersConnection, requestMessage webmodel.WSMessage, essence string, id string, err error) error {
	if err == models.ErrNoRecords {
		return uc.WSBadRequest(requestMessage, fmt.Sprintf("no %s with id '%s'", essence, id))
	}

	return uc.WSError(fmt.Sprintf("getting %s with id '%s' failed", essence, id), err)
}
