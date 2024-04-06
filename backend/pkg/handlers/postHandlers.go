package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/errorhandle"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/session"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

type contestKey string

const CTX_USER = contestKey("userSession")

// page functions
func GetAllPosts(app *application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// get user ID from r.context for getting all posts from people who he/she follows

		userID, err := getUserIDFromRequest(r)
		if err != nil {
			errorhandle.ServerError(app, w, r, "GetAllPosts: getting userID from context failed", err)
			return
		}

		posts, err := app.DBModel.GetFiltredPosts(models.Filter{}, userID, 20, 0)
		if err != nil {
			errorhandle.ServerError(app, w, r, "error getting posts from db", err)
			return
		}

		payload, _ := json.Marshal(posts)
		w.Write(payload)
		// not needed: w.WriteHeader(http.StatusOK)
	})
}

func SubmitPost(app *application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		postData := &webmodel.Post{}
		json.NewDecoder(r.Body).Decode(postData)

		userID, err := getUserIDFromRequest(r)
		if err != nil {
			errorhandle.ServerError(app, w, r, "SubmitPost: getting userID from context failed", err)
			return
		}

		post := &models.Post{
			Theme:      postData.Theme,
			Categories: postData.Categories,
			Content: models.Content{
				Text:       postData.Content,
				DateCreate: postData.DateCreate,
				UserID:     userID,
			},
		}

		post, err = app.DBModel.AddPost(post) // the response post with timestamp, postID etc
		if err != nil {
			errorhandle.ServerError(app, w, r, "SubmitPost: error adding post to db:", err)
			return
		}

		payload, err := json.Marshal(post)
		if err != nil {
			errorhandle.ServerError(app, w, r, "GetUserPosts: marshaling posts failed", err)
			return
		}
		// not needed w.WriteHeader(http.StatusOK)
		w.Write(payload)
	})
}

func SubmitPostImg(app *application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		imageData := &webmodel.PostImg{}
		json.NewDecoder(r.Body).Decode(&imageData)

		if imageData.Image != "" {
			img := strings.Split(imageData.Image, ",")
			bytes, _ := helpers.ConvertBase64ToImg(img[1])
			fileName, _ := helpers.CreateImgFromBytes(bytes, imageData.PostID, img[0], "post")
			app.DBModel.ModifyPost(imageData.PostID, "", "", []string{fileName})
			payload, _ := json.Marshal(fileName)
			w.WriteHeader(http.StatusOK)
			w.Write(payload)
			return
		}

		w.WriteHeader(http.StatusOK)
	})
}

func SubmitComment(app *application.Application) http.Handler { // TODO
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID, ok := r.Context().Value(CTX_USER).(string) // ok -> boolean true if userID string was found in the context
		if !ok {
			log.Println("error getting user id from request context")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		fmt.Println(userID)
		// add post to database

		// if ok, return response post with timestamp etc
	})
}

func GetUserPosts(app *application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("get user posts")

		userID, err := getUserIDFromRequest(r)
		if err != nil {
			errorhandle.ServerError(app, w, r, "GetUserPosts: getting userID from context failed", err)
			return
		}

		filter := models.Filter{
			AuthorID: userID,
		}

		posts, err := app.DBModel.GetFiltredPosts(filter, userID, 10, 0)
		if err != nil {
			errorhandle.ServerError(app, w, r, "GetUserPosts: getting posts from db failed", err)
			return
		}

		payload, err := json.Marshal(posts)
		if err != nil {
			errorhandle.ServerError(app, w, r, "GetUserPosts: marshaling posts failed", err)
			return
		}

		// not needed w.WriteHeader(http.StatusOK)
		w.Write(payload)
	})
}

func getUserIDFromRequest(r *http.Request) (string, error) {
	sess, ok := r.Context().Value(CTX_USER).(*session.Session)
	if !ok {
		return "", errors.New("error getting session from request context")
	}

	if !sess.IsLoggedin() {
		return "", errors.New("unauthorized session in request context")
	}
	return sess.User.ID, nil
}
