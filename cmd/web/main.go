package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"
	"os"
	"vis/note/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	errorlog *log.Logger
	infolog  *log.Logger
	snippets *mysql.Snippetmodel
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:password@/sptbox?parseTime=true", "Mysql data source name")

	flag.Parse()

	infolog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorlog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

	db, err := openDB(*dsn)
	if err != nil {
		errorlog.Fatal(err)
	}
	defer db.Close()

	app := &application{
		errorlog: errorlog,
		infolog:  infolog,
		snippets: &mysql.Snippetmodel{DB: db},
	}

	infolog.Printf("Starting server %s", *addr)

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorlog,
		Handler:  app.routes(),
	}
	err = srv.ListenAndServe()
	errorlog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
