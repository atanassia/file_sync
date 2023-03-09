package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	errorLog   *log.Logger
	infoLog    *log.Logger
	fileLocate string
}

func main() {
	addr := flag.String("addr", ":8000", "HTTP network address")
	fileLocate := flag.String("fileLocate", "./folder/file.txt", "File location")
	flag.Parse()

	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime|log.Lshortfile)
	errorLog := log.New(os.Stderr, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)

	app := &config{
		errorLog:   errorLog,
		infoLog:    infoLog,
		fileLocate: *fileLocate,
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	infoLog.Printf("Starting server on %s", *addr)
	err := srv.ListenAndServe()
	errorLog.Fatal(err)
}
