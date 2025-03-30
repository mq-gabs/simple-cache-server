package main

import (
	"log"
	"net"
)

func main() {
	l, err := net.Listen("tcp", ":9012")

	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	
	store := NewStore()

	for {
		log.Println("Accepting...")
		conn, err := l.Accept()
		
		if err != nil {
			log.Printf("Error when accepting connection: %v", err)
			continue
		}
		
		h := NewHandler(conn, store)

		go h.Handle()
	}
}