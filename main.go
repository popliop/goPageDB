package main

import (
	"fmt"

	server "git.dhl.com/marsorm/goPageDB/pkg"
)

func main() {
	fmt.Println("Hello World")

	apiServer := server.NewAPIServer(":8080")
	apiServer.Run()
}
