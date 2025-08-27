package main

import (
	"log"
	"net"
	"net/http"
	"scas/handler"
	"scas/store"
	"text/template"
)

func main() {
	l, err := net.Listen("tcp", ":9012")

	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}

	store := store.New()
	go httpServer(store)

	for {
		conn, err := l.Accept()

		if err != nil {
			log.Printf("Error when accepting connection: %v", err)
			continue
		}

		h := handler.New(conn, store)

		go h.Handle()
	}
}

type IndexTemplate struct {
	Table []store.TableItem
}

func httpServer(s *store.Store) {
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
