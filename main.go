package main

import (
	"github.com/GiladLeef/SharkCrack/userinput"
	"os"
)

func main() {
	args := os.Args[1:]
	if userinput.IsClient(args) {
		client := InitializeClient()
		client.Start()
	} else {
		server := InitializeServer()
		server.Start()
	}
}
