package main

import (
	"flag"
	"fmt"
	"log"
)

type config struct {
	server     string
	fileLocate string
}

func main() {
	server := flag.String("server", "http://127.0.0.1:8000", "Server address")
	fileLocate := flag.String("fileLocate", "./folder/file.txt", "File location")
	flag.Parse()

	app := &config{
		server:     *server,
		fileLocate: *fileLocate,
	}

	fmt.Printf("Сервер - %s, локальный файл находится в %s\n", app.server, app.fileLocate)
	
	fmt.Println("Проверка файлов на клиенте и сервере, если были внесены изменения в файл клиента, они подтянутся и на сервер.")
	count, err := app.checkForUpdate()
	if err != nil{
		log.Fatalln(err)
	}

	fmt.Println("\nДанные обновлены.\nПроверка изменений файла...")

	go func() {
		for i := 1;; i++{
			err := app.watchFile(app.fileLocate)
			if err != nil {
				fmt.Println(err)
			}

			count, err = app.sendUpdates(count)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Printf("Файл был сохранен в %d раз.\n", i)
		}

	}()

	<-make(chan struct{})
}
