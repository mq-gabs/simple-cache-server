package main

import (
	scas "github.com/mq-gabs/simple-cache-server/clients/go"
	"log"
)

func main() {
	c, err := scas.CreateConnection(&scas.Config{})
	defer c.Close()

	if err != nil {
		log.Fatalf("cannot connect: %v", err)
	}

	err = c.Set("name", "John")

	if err != nil {
		log.Fatal(err)
	}

	res, err := c.Get("name")

	log.Println(res)
}

