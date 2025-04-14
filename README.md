# :rocket: Simple Cache Server (SCAS)

## :question: What is it?

Simple Cache Server (SCAS) is a lightweight in-memory cache server written in Golang. It operates over the TCP protocol and supports three primary actions:

- :white_check_mark: SET: Store a value associated with a key.

- :mag_right: GET: Retrieve a value by its key.

- :x: ERASE: Remove a key-value pair from memory.

## :star: Features

⚡ Lightweight and fast in-memory caching.

🎯 Simple protocol for easy communication.

🌐 TCP-based communication for flexibility.

🛠️ A dedicated Go client for interaction with the server.

## 📦 Installation

### :computer: Server

Clone the repository and build the server:

```bash
git clone https://github.com/mq-gabs/simple-cache-server.git
cd simple-cache-server/server
go build -o scas
```

Run the server:

```bash
./scas
```
You can view keys and its values by access port `9013` on browser.


## :link: Protocol

SCAS communicates over TCP with a custom binary protocol:

- :satellite: Messages are sent as bytes.

- :small_blue_diamond: Sections are separated by 0x0A.

- :ticket: The first byte represents the action (SET, GET, ERASE).

- :key: The second section contains the key.

- :outbox_tray: The third section (only for SET) contains the value.

## 🧑‍💻 Client example

### Go

Here is a simple of example of how to use the client in Go.

- Create a file named `main.go`
```go
package main

import (
	scas "github.com/mq-gabs/simple-cache-server/clients/go"
	"log"
)

func main() {
	// Creates conection
	c, err := scas.CreateConnection(&scas.Config{})
	// Closes connection at the end
	defer c.Close()
	if err != nil {
		log.Fatalf("cannot connect: %v", err)
		return
	}
	log.Println("connection created")

	// Creates a key 'name' and save the value 'John Doe'
	err = c.Set("name", "John Doe")
	if err != nil {
		log.Fatalf("cannot set: %v", err)
	}
	log.Println("value saved")

	// Retrieve the value saved in the key 'name'
	value, err := c.Get("name")
	if err != nil {
		log.Fatalf("cannot get: %v", err)
	}
	log.Printf("value: %v", value)
}
```

- Install the client
```bash
go mod tidy
```

- Run
```bash
go run main.go
```
