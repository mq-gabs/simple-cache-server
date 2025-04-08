package main

import (
	"log"
)

func main() {
	conn, err := CreateConnection(&Config{})

	if err != nil {
		log.Fatal(err)
	}

	if err := conn.SetWithTTL("name", "John Doe", 10000); err != nil {
		log.Fatal(err)
	}

    v, err := conn.Get("name")
	if err != nil {
		log.Fatal(err)
	}

	log.Println(v)

	if err := conn.Erase("name"); err != nil {
		log.Fatal(err)
	}

	log.Println("End")
}
