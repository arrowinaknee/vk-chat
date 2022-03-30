package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

var connections []*websocket.Conn
var messages []string

func handleConn(ws *websocket.Conn) {
	connections = append(connections, ws)
	log.Printf("New connection, total: %d", len(connections))
	for {
		msg := make([]byte, 512)
		n, err := ws.Read(msg)
		if err != nil {
			if err == io.EOF {
				for i, o := range connections {
					if o == ws {
						connections = append(connections[:i], connections[i+1:]...)
					}
				}
				log.Printf("Connection closed, total: %d", len(connections))
				return
			}
			log.Fatal(err)
		}
		messages = append(messages, string(msg[:n]))
		for _, ows := range connections {
			fmt.Fprintf(ows, "%s\n", msg[:n])
		}
	}
}

func handleMessageHistory(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(messages)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		log.Printf("Error serving /messages/history: %s\n", err)
		return
	}
	w.Write(bytes)
}

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

func main() {
	s := websocket.Server{Handler: websocket.Handler(handleConn)}
	http.HandleFunc("/", handlePage)
	http.HandleFunc("/messages/ws", s.ServeHTTP)
	http.HandleFunc("/messages/history", handleMessageHistory)
	log.Fatal(http.ListenAndServe(":8089", nil))
}
