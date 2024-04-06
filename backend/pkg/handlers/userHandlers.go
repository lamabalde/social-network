package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/errorhandle"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/helpers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/session"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/webmodel"
)

func RegisterUser(app *application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		body := webmodel.UserCredentials{}
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			errorhandle.ServerError(app, w, r, "RegisterUser: unmarshalling failed", err)
			return
		}
		if body.UserName == "" { // if the user has not specified a user name
			body.UserName = body.FirstName
		}
		user, err := controllers.CreateUser(body)
		if err != nil {
			errorhandle.ServerError(app, w, r, "RegisterUser: creating user failed", err)
			return
		}
		userID, err := app.DBModel.AddUser(user)

		if err != nil {
			if errors.Is(err, models.ErrUniqueUserName) {
				app.InfoLog.Println("cannot add user to DB: UserName already exists.")

				payload, _ := json.Marshal(models.ErrorMessage{
					Errors: "User already exists.",
				})

				w.WriteHeader(http.StatusBadRequest)
				w.Write(payload)
				return
			}

			if errors.Is(err, models.ErrUniqueUserEmail) {
				app.InfoLog.Println("cannot add user to DB: Email already exists.")

				payload, _ := json.Marshal(models.ErrorMessage{
					Errors: "User with this email already exists.",
				})

				w.WriteHeader(http.StatusBadRequest)
				w.Write(payload)
				return
			}

			errorhandle.ServerError(app, w, r, "RegisterUser: adding user to DB failed", err)
			return
		}
		if body.Image != "" {
			img := strings.Split(body.Image, ",")
			bytes, _ := helpers.ConvertBase64ToImg(img[1])
			_, _ = helpers.CreateUserImgFromBytes(bytes, userID, img[0])
		}
		w.WriteHeader(http.StatusOK) // if register was a success, then we send 200OK resp
	})
}

func LoginUser(app *application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("logging in user")
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		fmt.Printf("req  method  %#v\nheader %#v\n", r.Method, r.Header)
		//TODO handle OPTIONS method if there is an unsafe request
		body := webmodel.UserCredentials{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			errorhandle.ServerError(app, w, r, "LoginUser: unmarshalling failed", err)
			return
		}
		if body.UserName == "" { // if when registering and username isn't specified
			body.UserName = body.FirstName
		}
		user, err := app.DBModel.GetUserByNameOrEmail(body.UserName)
		if err != nil {
			if err == models.ErrNoRecords {
				app.InfoLog.Printf("Login attemt failed: User '%s' doesn't exist in DB", body.UserName)

				w.WriteHeader(http.StatusUnauthorized)
				payload, _ := json.Marshal(models.ErrorMessage{
					Errors: "Invalid login",
				})
				w.Write(payload)
				return
			}
			errorhandle.ServerError(app, w, r, "LoginUser: getting user from DB failed", err)
			return
		}

		fmt.Println(user.Password, body.Password)
		if !helpers.CompareHashToPassword(user.Password, body.Password) { // check password
			app.InfoLog.Println("LoginUser: Invalid password")

			w.WriteHeader(http.StatusUnauthorized)
			payload, _ := json.Marshal(models.ErrorMessage{
				Errors: "Invalid password",
			})
			w.Write(payload)
			return
		}

		sess, cookie, err := session.New(app, &models.UserBase{ID: user.ID, UserName: user.UserName})
		if err != nil {
			errorhandle.ServerError(app, w, r, "LoginUser: session creation failed", err)
			return
		}

		payload, _ := json.Marshal(struct {
			UserID    string
			Username  string
			SessionID string
		}{sess.User.ID, sess.User.UserName, sess.SessionUuid})

		http.SetCookie(w, cookie)
		// not needed w.WriteHeader(http.StatusOK)
		w.Write(payload)
	})
}

func LogoutUser(app *application.Application) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		// TODO what about ws connection, will it be closed when the logout button is clicked?
		cookie, err := r.Cookie("userLoggedIn")

		if err == nil && cookie.Value != "" {
			if err := app.DBModel.DeleteUsersSession(cookie.Value); err != nil {
				errorhandle.ServerError(app, w, r, "LogoutUser: deleting session failed", err)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(""))
	})
}
