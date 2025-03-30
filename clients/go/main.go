package main

import "log"

func main() {
	c, err := CreateConnection(&Config{})
	if err != nil {
		log.Fatal(err.Error())
	}

	
	err = c.Set("name", "John Doe")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Saved!")

	res, err := c.Get("name")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s", res)

	err = c.Set("age", "20")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Saved!")

	
	res, err = c.Get("age")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %s", res)

	err = c.Erase("age")
	if err != nil {
		log.Fatal(err)
	}

	res, err = c.Get("age")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Response: %v", res)
}