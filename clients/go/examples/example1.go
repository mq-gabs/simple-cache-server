package main

import (
	"scas"
	"log"
)

func main() {
	c, err := scas.CreateConnection(&scas.Config{})

	if err != nil {
		log.Fatalf("cannot connect: %v", err)	
	}

	c.Set("name", "John")

	res, err := c.Get("name")

	log.Println(res)
}