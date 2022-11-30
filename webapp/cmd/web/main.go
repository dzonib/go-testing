package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"webapp/pkg/data"
	"webapp/pkg/repository"
	"webapp/pkg/repository/dbrepo"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	DSN     string
	DB      repository.DatabaseRepo
	Session *scs.SessionManager
}

func main() {
	// register type with application
	gob.Register(data.User{})

	// set up an app config
	app := application{}

	flag.StringVar(&app.DSN, "dsn", "host=localhost port=5434 user=postgres password=postgres dbname=users sslmode=disable timezone=UTC connect_timeout=5", "Postgres connection")

	flag.Parse()

	conn, err := app.connectToDb()

	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	// now we can use all db methods
	app.DB = &dbrepo.PostgresDBRepo{
		DB: conn,
	}

	// get a session manager
	app.Session = getSession()

	// get application routes
	mux := app.routes()

	// print out a message
	log.Println("Starting Server on port 8081..")

	// start a server
	err = http.ListenAndServe(":8081", mux)

	if err != nil {
		log.Fatal(err)
	}
}
