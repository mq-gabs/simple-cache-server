package main

import (
	"log"
	"net"
	"net/http"
	"text/template"
)

func main() {
	l, err := net.Listen("tcp", ":9012")

	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}

	store := NewStore()
	go httpServer(store)

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Printf("Error when accepting connection: %v", err)
			continue
		}

		h := NewHandler(conn, store)

		go h.Handle()
	}
}

type IndexTemplate struct {
	Table []TableItem
}

func httpServer(s *Store) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t, _ := template.ParseFiles("templates/index.html")

		table := s.GetTable()

		indexTemplate := IndexTemplate{
			Table: table,
		}

		t.Execute(w, indexTemplate)
	})

	http.ListenAndServe(":9013", nil)
}
