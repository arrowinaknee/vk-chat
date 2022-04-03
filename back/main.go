package main

import (
	"os"

	"ru.arrowinaknee.vk-chat/server"
)

func main() {
	var s = server.Server{
		Addr:   ":8089",
		DBLink: os.Getenv("ARROWCHAT_DB_URI"),
		DBName: "arrowchat",
	}

	s.Start()
}
