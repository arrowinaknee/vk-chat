package main

import "os"

var db db_t = db_t(os.Getenv("ARROW_DB_URI"))

func main() {
	serveRoutes(":8089")
}
