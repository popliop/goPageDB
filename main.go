package main

import (
	"fmt"
	"log"

	"github.com/marsorm/goPageDB/pkg/db"
	"github.com/marsorm/goPageDB/pkg/server"
)

func main() {
	fmt.Println("Hello World")
	database, err := db.NewPostgressStorage()
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.CreateAccountTable(database)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("dbAnswer")

	apiServer := server.NewAPIServer(":8080")
	apiServer.Run()
}
