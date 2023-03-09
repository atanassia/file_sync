package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
)

// serverError() возвращает клиенту ошибку сервера (500) и выводит стек ошибок в логи
func (app *config) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.errorLog.Println(err)
	app.errorLog.Output(2, trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// clientError() возвращает ошибку клиента (400)
func (app *config) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// notFound() возвращает 404 ошибку, указанный путь не найден
func (app *config) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// lineCounter возвращает кол-во строк в файле сервера
func (app *config) lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 1
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
