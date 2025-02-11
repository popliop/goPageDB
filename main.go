package main

import (
	"fmt"

	"github.com/marsorm/goPageDB/pkg/server"
)

func main() {
	fmt.Println("Hello World")

	apiServer := server.NewAPIServer(":8080")
	apiServer.Run()
}
