package main

import (
	"os"
)

func main() {

	go queue_receive()

	go cleanCache()

	LR_port := os.Getenv("PORT")
	if len(LR_port) == 0 {
		LR_port = "7070"
	}

	LR_server := LinkRedirectServer()
	LR_server.Run(":" + LR_port)

}
