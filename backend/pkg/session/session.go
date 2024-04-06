package session

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/models"
)

type LoginStatus byte

const (
	Loggedin LoginStatus = iota
	Experied
	Notloggedin
)

const (
	EXP_SESSION               = 24 * time.Hour
	TIME_BEFORE_AFTER_REFRESH = 30 * time.Second
)

const SESSION_TOKEN = "userLoggedIn"

type Session struct {
	loginStatus   LoginStatus
	SessionUuid   string
	ExpirySession time.Time
	User          *models.UserBase
}

func New(app *application.Application, user *models.UserBase) (*Session, *http.Cookie, error) {
	expiresAt := time.Now().Add(EXP_SESSION)

	sessionID, err := app.DBModel.AddUserSession(user.ID, expiresAt)
	if err != nil {
		return nil, nil, fmt.Errorf("adding session failed: %w", err)
	}

	app.InfoLog.Printf("session tocken '%s' is added for the user ID '%s'", sessionID, user.ID)

	return &Session{loginStatus: Loggedin, SessionUuid: sessionID, ExpirySession: expiresAt, User: user}, createCookie(sessionID, expiresAt), nil
}

/*
returns session which contains status of login and uses's data if it's logged in.
If it is left lrss than 30 sec to expiried time, it will refresh the session
If an error occurs it will response to the client with error status and return the error
*/
func Get(app *application.Application, r *http.Request) (*Session, *http.Cookie, error) {
	session := GetNotloggedinSession()
	cookie, err := r.Cookie(SESSION_TOKEN)
	if err != nil && err != http.ErrNoCookie {
		return nil, nil, fmt.Errorf("getting cookie failed: '%s', url: '%s'", err, r.URL)
	}
	if err == http.ErrNoCookie || cookie.Value == "" {
		return session, createCookie("", time.Now()), nil // session status = notloggedin
	}

	// there is a sessionToken
	sessionToken := cookie.Value

	session.User, session.SessionUuid, session.ExpirySession, err = app.DBModel.GetUserBySession(sessionToken)
	if err != nil {
		if err == models.ErrNoRecords {
			return session, createCookie("", time.Now()), nil // session status = notloggedin
		}
		return nil, nil, fmt.Errorf("getting a user by uuid failed: %w", err)
	}

	if session.isExpired() {
		// delete the session & return expiried status
		session.User = nil
		err := app.DBModel.DeleteUsersSession(session.SessionUuid)
		if err != nil {
			return nil, nil, fmt.Errorf("deleting the expired session failed: %w", err)
		}

		session.loginStatus = Experied
		return session, createCookie("", time.Now()), nil
	}

	if session.timeToExpired() < TIME_BEFORE_AFTER_REFRESH {
		// refresh the session
		session, cookie, err = New(app, session.User)
		if err != nil {
			return nil, nil, fmt.Errorf("session creating failed: %w", err)
		}

		return session, cookie, nil
	}

	// user was found and their time was not expired:
	session.loginStatus = Loggedin
	return session, cookie, nil
}

/*
checks if a loggedin session is expired and change the session status
or if the session is expired then delete the session and change the session status to loggedout.
Returns the new status of the session.
*/
func (s *Session) Tidy(app *application.Application) (LoginStatus, error) {
	if s == nil {
		return Notloggedin, nil
	}
	switch s.loginStatus {
	case Loggedin:
		if s.User == nil {
			s.loginStatus = Notloggedin
			return Notloggedin, errors.New("exception for the session status: loggedIn status and nil User")
		}
		if s.isExpired() {
			s.loginStatus = Experied
			return Experied, nil
		}
	case Experied:
		if s.User == nil {
			s.loginStatus = Notloggedin
			return Notloggedin, errors.New("exception for the session status: expired status and nil User")
		}
		err := app.DBModel.DeleteUsersSession(s.SessionUuid)
		if err != nil {
			return Experied, fmt.Errorf("deleting the expired session failed: %w", err)
		}
		s.loginStatus = Notloggedin
		s.User = nil
		return Notloggedin, nil
	}
	return s.loginStatus, nil
}

func (s *Session) IsLoggedin() bool {
	return s != nil && s.loginStatus == Loggedin && !s.isExpired()
}

func (s *Session) GetStatus() string {
	var status string
	switch s.loginStatus {
	case Loggedin:
		status = "logged"
	case Experied:
		status = "experied"
	case Notloggedin:
		status = "not logged in"
	}
	return status
}

func (s *Session) isExpired() bool {
	exp := s.ExpirySession
	return exp.Before(time.Now())
}

func (s *Session) timeToExpired() time.Duration {
	exp := s.ExpirySession
	return time.Until(exp)
}

func (s *Session) String() string {
	return fmt.Sprintf("Status: %v\nSessionUuid: %s\nExp: %v\nUser:\n  ID: %s\n  UserName: %s\n", s.loginStatus, s.SessionUuid, s.ExpirySession, s.User.ID, s.User.UserName)
}

func GetNotloggedinSession() *Session {
	return &Session{loginStatus: Notloggedin, User: nil}
}

func createCookie(value string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     "userLoggedIn",
		Value:    value,
		MaxAge:   3600, // one hour
		Expires:  expiresAt,
		Path:     "http://localhost:8080", // TODO: path of the frontend?
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}
}
