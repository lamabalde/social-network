package controllers

import (
	"fmt"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel/parse"
)

// proccess the comment creation: adds the commnet to DB and replyes with the full post
func ReplyNewComment(app *application.Application) wsconnection.FuncReplyCreator {
	return func(currConnection *wsconnection.UsersConnection, message webmodel.WSMessage) (any, error) {
		comment, err := parse.PayloadToComment(message.Payload)
		if err != nil {
			return nil, currConnection.WSError(fmt.Sprintf("Invalid payload for a new comment: %s", message.Payload), err)
		}

		errmessage := comment.Validate()
		if errmessage != "" {
			return nil, currConnection.WSBadRequest(message, errmessage)
		}

		err = saveCommentToDB(app, currConnection, message, comment, currConnection.Client.UserID)
		if err != nil {
			return nil, err
		}

		post, err := getPost(app, currConnection, comment.PostID, message)
		if err != nil {
			return nil, err
		}

		return post, nil
	}
}

func saveCommentToDB(app *application.Application, currConnection *wsconnection.UsersConnection, message webmodel.WSMessage, comment webmodel.Comment, authorID string) error {
	// dateCreate := comment.DateCreate //TODO get date from frontEnd

	id, err := app.DBModel.AddComment(comment.PostID, comment.Content, nil, authorID)
	if err != nil {
		return handleErrAddToDB(app, currConnection, message,
			fmt.Sprintf("a new comment to post '%s' from userID '%s'", comment.PostID, authorID), err)
	}

	app.InfoLog.Printf("A comment is added to DB. id: '%s'", id)
	return nil
}

// Delete the comment -not WS version
/*
func DeleteComment(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var id int
		var err error
		var comment *models.Comment
		var user *models.User

		// Only authenticate users can delete comments
		if sess := r.Context().Value(acl.SessionKey); !sess.(*session.Session).IsLoggedin() { // we might want to store user id instead
			errMsg := "Authentication needed"
			json.NewEncoder(w).Encode(map[string]string{"status": "failure", "error": errMsg})
			app.InfoLog.Printf("Failed to delete comment: %s", errMsg)
			return
		}

		// Check if comment id is provided
		if r.FormValue("id") == "" {
			errMsg := "Missing Comment Id"
			json.NewEncoder(w).Encode(map[string]string{"status": "failure", "error": errMsg})
			app.InfoLog.Printf("Failed to delete Comment: %s", errMsg)
			return
		}

		if id, err = strconv.Atoi(r.FormValue("id")); err != nil {
			errMsg := "Invalid comment Id"
			json.NewEncoder(w).Encode(map[string]string{"status": "failure", "error": errMsg})
			app.InfoLog.Printf("Failed to delete comment: %s", errMsg)
			return
		}

		// Check if we have to post
		if comment, err = app.DBModel.GetCommentByID(id); err != nil {
			errMsg := "No such comment in DB"
			json.NewEncoder(w).Encode(map[string]string{"status": "failure", "error": errMsg})
			app.InfoLog.Printf("Failed to delete comment: %s", errMsg)
			return
		}

		// Check if we are the author of this post
		if comment.Message.Author.ID != user.ID {
			errMsg := "User is not the author of this comment"
			json.NewEncoder(w).Encode(map[string]string{"status": "failure", "error": errMsg})
			app.InfoLog.Printf("Failed to delete comment: %s", errMsg)
			return
		}

		if err = app.DBModel.Delete(id); err != nil {
			errorhandle.ServerError(app, w, r, fmt.Sprintf("comment delete failed: func %s:", logger.GetCurrentFuncName()), err)
			return
		}

		if err = app.DBModel.DeleteCommentLikeByMessageID(id); err != nil {
			errorhandle.ServerError(app, w, r, fmt.Sprintf("comment delete failed: func %s:", logger.GetCurrentFuncName()), err)
			return
		}

		app.InfoLog.Printf("Post deleted successfully '%v'", comment.ID)
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "redirect": fmt.Sprintf("/post/%v", comment.PostID)})
	}
}
*/
