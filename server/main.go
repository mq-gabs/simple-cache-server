package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"scas/handler"
	"scas/store"
	"scas/utils/worker"
	"text/template"
)

var logger = log.Default()

const (
	protoTCP   = "tcp"
	serverPort = "9012"

	NumWorkers     = 1000
	NumWorkerQueue = 100
)

func main() {
	ctx := context.Background()

	w := worker.New(ctx, NumWorkers, NumWorkerQueue)
	defer w.Stop()

	l, err := net.Listen(protoTCP, ":"+serverPort)
	if err != nil {
		log.Fatalf("[ERROR] cannot start server: %v", err)
	}

	store := store.New()
	go httpServer(store)

	for {
		conn, err := l.Accept()
		if err != nil {
			logger.Printf("[ERROR] cannot accept connection: %v\n", err)
			continue
		}

		h := handler.New(conn, store)

		w.Submit(func(poolCtx context.Context) {
			h.Handle(poolCtx)
		})
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
