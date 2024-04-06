package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/logger"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/errorhandle"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/handlers"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/session"
)

func routerMiddleware(app *application.Application, next http.Handler) http.Handler { // middleware, to check authentication
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		w.Header().Add("Access-Control-Allow-Headers", "*")
		w.Header().Add("Access-Control-Allow-Origin", "http://localhost:8080")

		start := time.Now()
		req := fmt.Sprintf("%s %s", r.Method, r.URL)
		app.InfoLog.Println(req)

		sess, respCookie, err := session.Get(app, r)
		if err != nil {
			errorhandle.ServerError(app, w, r, fmt.Sprintf("func %s failed ", logger.GetCurrentFuncName()), err)
			return
		}

		http.SetCookie(w, respCookie)

		if !sess.IsLoggedin() { // if user is not logged in(no cookie or session is expired), direct to main page
			app.InfoLog.Println("User is not logged in")
			w.WriteHeader(http.StatusUnauthorized)

		}

		// Store user information in the request context
		app.InfoLog.Printf("session for user '%v' is added to the context", sess.User)
		ctx := context.WithValue(r.Context(), handlers.CTX_USER, sess)
		// Call the next handler in the chain with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Println(req, "completed in", time.Since(start))
	})
}

func CreateAPIroutes(mux *http.ServeMux, app *application.Application, wsHandlers [2]wsconnection.WSmux) *http.ServeMux {
	// no auth stuff
	mux.Handle("/logout", handlers.LogoutUser(app))
	mux.Handle("/loginUser", handlers.LoginUser(app))
	mux.Handle("/registerUser", handlers.RegisterUser(app))
	mux.Handle("/submitPostImg", handlers.SubmitPostImg(app))
	mux.Handle("/submitCommentImg", handlers.SubmitCommentImg(app))

	mux.Handle("/img/", http.StripPrefix("/img/", http.FileServer(http.Dir("data/img"))))
	// auth stuff
	mux.Handle("/getUserPosts", routerMiddleware(app, handlers.GetUserPosts(app))) // for when user profile is clicked
	mux.Handle("/submitPost", routerMiddleware(app, handlers.SubmitPost(app)))
	mux.Handle("/getAllPosts", routerMiddleware(app, handlers.GetAllPosts(app)))
	mux.Handle("/submitComment", routerMiddleware(app, handlers.SubmitComment(app)))
	mux.Handle("/websocket", routerMiddleware(app, handlers.IndexWs(app, wsHandlers[WS_REPLYERS_UNI])))               // websocket stuff here
	mux.Handle(handlers.CREATE_CHAT_URL, routerMiddleware(app, handlers.OpenChat(app, wsHandlers[WS_REPLYERS_CHAT]))) // websocket stuff here
	//mux.Handle(handlers.JOIN_CHAT_URL, routerMiddleware(app, handlers.JoinGroupChat(app, wsHandlers[WS_REPLYERS_CHAT]))) // websocket stuff here
	mux.Handle("/checkAuth", routerMiddleware(app, passAuth())) // for initial page loading to check if a user is logged in
	return mux
}

func passAuth() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(nil)
	})
}
