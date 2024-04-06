package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/application"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/controllers/wsconnection"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/db/sqlite"
	"01.kood.tech/git/Hems_Chrisworth/social-network/backend/pkg/routes"
)

var port = "8000" // DB   *sqlite.DBModel

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile) // set flags in log package to print out file/line number on errors
	var err error

	// app keeps all dependences used by handlers

	testDB, versionDB, err := parseArgs() // parse arguments,  default testDB is false
	if err != nil {
		log.Fatalln(err)
	}

	dbModel, err := sqlite.InitDB(testDB, versionDB)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbModel.DB.Close()

	addr := fmt.Sprintf(":%s", port) // localhost
	app := application.New(dbModel, addr)
	if err != nil {
		log.Fatalln(err)
	}

	var wsHandlers [2]wsconnection.WSmux
	wsHandlers[routes.WS_REPLYERS_UNI] = routes.CreateUniWsRoutes(app)
	wsHandlers[routes.WS_REPLYERS_CHAT] = routes.CreateChatWsRoutes(app)

	mux := http.NewServeMux()
	app.Server.Handler = routes.CreateAPIroutes(mux, app, wsHandlers)

	go app.Hub.Run()
	app.InfoLog.Println("The chat Hub is running...")

	log.Println("main: running server on port", port)
	if err := app.Server.ListenAndServe(); err != nil {
		app.ErrLog.Fatalf("main: couldn't start server: %v\n", err)
	}
}

// Parses the program's arguments to obtain the server port. If no arguments found, it uses the 8000 port by default
// Usage: go run . --testdb
func parseArgs() (testDB bool, versionDB int, err error) {
	usage := `wrong arguments
	Usage: go run ./app [OPTIONS]
	OPTIONS: 
			--testdb start with the test DB
			--versiondb=desired version of DB
			--migrpath=path to the migration folder`
	flag.BoolVar(&testDB, "testdb", false, "--testdb if you want to start with the test DB")
	flag.IntVar(&versionDB, "versiondb", 0, "--versiondb=desired version of DB")
	flag.Parse()
	if flag.NArg() > 0 {
		return false, 0, fmt.Errorf(usage)
	}

	return
}
