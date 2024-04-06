package controllers

import (
	"fmt"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

// Proccess the post creation
func ReplyNewPost(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		postData, err := parse.PayloadToPost(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for a new post: %s", message.Payload), err)
		}

		errmessage := postData.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		post, err := savePostToDB(app, currConnection, message, postData)
		if err != nil {
			return nil, err
		}

		post.Content.Likes = []int{0, 0} // to prevent error in frontend
		return post, nil
	}
}

func savePostToDB(app *application.Application, currConnection *wsconnection.UsersConnection, message webmodel.WSMessage, postData webmodel.Post) (*models.Post, error) {
	// dateCreate := postData.DateCreate //TODO get date from FrontEnd
	post := createPostFromWebModelPost(postData)
	post.Content.UserID = currConnection.Client.UserID
	post.Content.UserName = currConnection.Client.UserName

	var err error
	post, err = app.DBModel.AddPost(post)
	if err != nil {
		return nil, handleErrAddToDB(app, currConnection, message,
			fmt.Sprintf("a new post from userID '%s'", post.Content.UserID), err)
	}

	app.InfoLog.Printf("Post is added to DB. id: '%s' categories: '%v' ", post.ID, postData.Categories)
	return post, nil
}

func createPostFromWebModelPost(postData webmodel.Post) *models.Post {
	return &models.Post{
		Theme:      postData.Theme,
		Categories: postData.Categories,
		Content: models.Content{
			Text:       postData.Content,
			Images:     nil,
			DateCreate: time.Now(),
		},
		GroupID: postData.GroupID,
		Privacy: postData.PostType,
	}
}
