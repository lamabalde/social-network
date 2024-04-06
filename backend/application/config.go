package application

import (
	"log"
	"net/http"
	"time"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/logger"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wshub"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite/queries"

	"github.com/gorilla/websocket"
)

type Application struct {
	ErrLog  *log.Logger
	InfoLog *log.Logger
	Hub      *wshub.Hub
	DBModel  *queries.DBModel
	Upgrader websocket.Upgrader
	Server   *http.Server
}

func New(dbModel *queries.DBModel, serverAddress string) *Application {
	application := &Application{}

	application.ErrLog, application.InfoLog = logger.CreateLoggers()

	application.Hub = wshub.NewHub()

	application.DBModel = dbModel

	application.Upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "http://localhost:8080" || origin == "http://127.0.0.1:8080"
			//return true
		},
	}

	application.Server = &http.Server{
		Addr:         serverAddress,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
		ErrorLog:     application.ErrLog,
	}

	return application
}
