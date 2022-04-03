package main

import (
	"ru.arrowinaknee.vk-chat/server"
)

func main() {
	var s = server.Server{
		Addr:   ":8089",
		DBLink: "mongodb://91.122.53.183:25560",
		// DBLink: os.Getenv("ARROWCHAT_DB_URI"),
		DBName: "arrowchat",
	}

	s.Start()
}
