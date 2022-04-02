package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Temorary file, handles app page serving before a proper front-end app is built

func (s *Server) handlePage(w http.ResponseWriter, r *http.Request) {
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
