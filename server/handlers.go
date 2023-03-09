package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type jsondata struct {
	Data []string `json:"data"`
}

// uploadChages() загружает изменения в файл
func (app *config) uploadChages(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	var data jsondata
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		app.errorLog.Println(err)
	}

	f, err := os.OpenFile(app.fileLocate, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		app.errorLog.Println(err)
	}

	for _, value := range data.Data {
		if _, err := f.WriteString(value); err != nil {
			app.serverError(w, err)
			return
		}
	}

	if err := f.Close(); err != nil {
		app.serverError(w, err)
		return
	}

	app.infoLog.Println("Данные обновлены.")

	fmt.Fprintln(w, "Данные загружены.")
}

// getUnsentLineCount получает новые строки от клиента (когда клиент только-только запустился) 
func (app *config) getUnsentLineCount(w http.ResponseWriter, r *http.Request) {
	f, err := os.OpenFile(app.fileLocate, os.O_RDONLY, 0644)

	if err != nil {
		app.errorLog.Println(err)
	}

	f.Stat()

	fileStat, err := f.Stat()
	if err != nil {
		app.errorLog.Println(err)
	}

	if fileStat.Size() == 0 {
		fmt.Fprintln(w, 0)
		return
	}

	n, err := app.lineCounter(f)
	if err != nil {
		app.serverError(w, err)
	}

	if err := f.Close(); err != nil {
		app.serverError(w, err)
		return
	}

	fmt.Fprintln(w, n)
}
