package main

import (
	"os"
)

func main() {

	CP_port := os.Getenv("PORT")
	if len(CP_port) == 0 {
		CP_port = "6060"
	}

	CPserver := ControlPanelServer()
	CPserver.Run(":" + CP_port)
}
