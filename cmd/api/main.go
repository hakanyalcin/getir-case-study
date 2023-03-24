package main

import (
	"flag"
	"fmt"
	"getir-case-study/internal/cache"
	"getir-case-study/internal/driver"
	"getir-case-study/internal/models"
	"log"
	"net/http"
	"os"
)

type config struct {
	port int
	db   struct {
		dsn string
	}
}

type application struct {
	config   config
	infoLog  *log.Logger
	errorLog *log.Logger
	DB       models.DBModel
	cache    *cache.LocalCache
}

func (app *application) serve() error {

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
		Addr: fmt.Sprintf(":%d", app.config.port),
	}

	app.infoLog.Printf("Starting Back end server on port %d\n", app.config.port)

	return srv.ListenAndServe()
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4001, "Server port to listen on")
	flag.StringVar(&cfg.db.dsn, "dsn", "mongodb+srv://challengeUser:WUMglwNBaydH8Yvu@challenge-xzwqd.mongodb.net/getircase-study?retryWrites=true", "DSN")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	lc := cache.NewLocalCache()

	conn, err := driver.ConnectDB(cfg.db.dsn)
	if err != nil {
		errorLog.Fatal(err)
	}

	app := &application{
		config:   cfg,
		infoLog:  infoLog,
		errorLog: errorLog,
		DB:       models.DBModel{DB: conn},
		cache:    lc,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
