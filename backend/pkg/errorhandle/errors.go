package errorhandle

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
)

func NotFound(app *application.Application, w http.ResponseWriter, r *http.Request) {
	app.ErrLog.Output(2, fmt.Sprintf("wrong path: %s", r.URL.Path))
	http.NotFound(w, r)
}

func ServerError(app *application.Application, w http.ResponseWriter, r *http.Request, message string, err error) {
	app.ErrLog.Output(2, fmt.Sprintf("fail handling the page %s: %s. ERR: %v\nDebug Stack:  %s", r.URL.Path, message, err, debug.Stack()))
	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
}

func ClientError(app *application.Application, w http.ResponseWriter, r *http.Request, errStatus int, logTexterr string) {
	app.ErrLog.Output(2, logTexterr)
	http.Error(w, "ERROR: "+http.StatusText(errStatus)+logTexterr, errStatus)
}

func BadRequestError(app *application.Application, w http.ResponseWriter, r *http.Request, logTexterr string) {
	ClientError(app, w, r, http.StatusBadRequest, logTexterr)
}

func MethodNotAllowed(app *application.Application, w http.ResponseWriter, r *http.Request, allowedMethods ...string) {
	if allowedMethods == nil {
		panic("no methods is given to func MethodNotAllowed")
	}
	allowdeString := allowedMethods[0]
	for i := 1; i < len(allowedMethods); i++ {
		allowdeString += ", " + allowedMethods[i]
	}

	w.Header().Set("Allow", allowdeString)
	ClientError(app, w, r, http.StatusMethodNotAllowed, fmt.Sprintf("using the method %s to go to a page %s", r.Method, r.URL))
}

func Forbidden(app *application.Application, w http.ResponseWriter, r *http.Request, reason string) {
	app.ErrLog.Output(2, fmt.Sprintf("access to '%s' was forbidden: '%s'", r.URL.Path, reason))

	w.WriteHeader(http.StatusForbidden) // Sets status code at 403
	w.Write([]byte("Forbiddent: Access denied"))
}
