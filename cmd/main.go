package main

import (
	"context"
	"fmt"
	"getir-case-study/internal/cache"
	"getir-case-study/internal/driver"
	"getir-case-study/internal/models"
	"log"
	"net/http"
	"os"
)

type application struct {
	port     int
	infoLog  *log.Logger
	errorLog *log.Logger
	DB       models.DBModel
	cache    *cache.LocalCache
}

func (app *application) serve() error {

	// routing implementation with standard lib
	http.HandleFunc("/in-memory", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			app.getEntry(w, r)
		case "POST":
			app.setEntry(w, r)
		}
	})
	http.HandleFunc("/records", app.getRecords)

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.port),
	}

	app.infoLog.Printf("Starting Back end server on port %d\n", app.port)

	return srv.ListenAndServe()
}

func main() {
	var (
		port = 4001
		dsn  = "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getircase-study?retryWrites=true"
	)

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// initialize local cache for in-memory database
	lc := cache.NewLocalCache()

	conn, err := driver.ConnectMongoClient(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer conn.Disconnect(context.TODO())

	app := &application{
		port:     port,
		infoLog:  infoLog,
		errorLog: errorLog,
		DB:       models.DBModel{DB: conn.Database("getircase-study")},
		cache:    lc,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
