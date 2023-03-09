package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type jsondata struct {
	Data []string `json:"data"`
}

// watchFile() каждую секунду проверяет состояние файла, обновился ли
func (app *config) watchFile(filePath string) error {
	initialStat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	for {
		stat, err := os.Stat(filePath)
		if err != nil {
			return err
		}

		if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {
			break
		}

		time.Sleep(1 * time.Second)
	}

	return nil
}

// getLastLineNumber() получает номер последней строки текстового файла сервера
func (app *config) getLastLineNumber() (int, error) {
	res, err := http.Get(fmt.Sprintf("%s/getUnsentLineCount", app.server))
	if err != nil {
		log.Printf("error making http request: %s\n", err)
		return 0, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	result, err := strconv.Atoi(strings.Trim(string(body), "\n"))
	if err != nil {
		return 0, err
	}

	return result, nil
}

// sendLinesToServer() отпарвляет новые сохраненные строки на сервер
func (app *config) sendLinesToServer(fileLines []string) error {
	var data jsondata
	data.Data = fileLines

	json, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return err
	}

	request, err := http.NewRequest("POST", fmt.Sprintf("%s/uploadChages", app.server), bytes.NewBuffer(json))
	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
		return err
	}

	defer response.Body.Close()

	return err
}

// checkForUpdate() проверяет, есть ли на клиентском файле изменения,
// которых еще нет на серверном, если есть - догружает
func (app *config) checkForUpdate() (int, error) {
	serverfile, err := app.getLastLineNumber()
	if err != nil {
		return 0, err
	}

	f, err := os.OpenFile(app.fileLocate, os.O_RDWR, 0644)
	if err != nil {
		return 0, err
	}

	var fileLines []string

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	count := 0
	for fileScanner.Scan() {
		if serverfile <= count {
			fileLines = append(fileLines, fileScanner.Text()+"\n")
		}
		count++
	}

	if err := f.Close(); err != nil {
		return 0, err
	}

	if serverfile > count {
		return 0, errors.New("fileServerError: Лишние записи в файле на сервере")
	}

	err = app.sendLinesToServer(fileLines)
	if err != nil {
		return 0, err
	}

	return count, err
}

// sendUpdates() собирает новые строки и вызывает sendLinesToServer() для отправки данных
func (app *config) sendUpdates(lastline int) (int, error) {
	f, err := os.OpenFile(app.fileLocate, os.O_RDWR, 0644)
	if err != nil {
		log.Println(err)
		return 0, err
	}

	var fileLines []string

	fileScanner := bufio.NewScanner(f)
	fileScanner.Split(bufio.ScanLines)

	count := 0
	for fileScanner.Scan() {
		if lastline <= count {
			fileLines = append(fileLines, fileScanner.Text()+"\n")
		}
		count++
	}

	if err := f.Close(); err != nil {
		return 0, err
	}

	err = app.sendLinesToServer(fileLines)
	if err != nil {
		return 0, err
	}

	return count, err
}
