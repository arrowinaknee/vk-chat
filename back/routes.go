package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func handlePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 page not found")
		return
	}

	bytes, err := ioutil.ReadFile("../front/index.html")
	if err != nil {
		log.Printf("can't load page: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not load the page")
		return
	}
	fmt.Fprint(w, string(bytes))
}

func serveRoutes(addr string) {
	s := websocket.Server{Handler: websocket.Handler(handleConn)}
	http.HandleFunc("/", handlePage)
	http.HandleFunc("/users", HandleUsers)
	http.HandleFunc("/messages/ws", s.ServeHTTP)
	http.HandleFunc("/messages/history", handleMessageHistory)
	log.Fatal(http.ListenAndServe(addr, nil))
}
