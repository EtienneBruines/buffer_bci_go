package main

import (
	"log"

	"github.com/EtienneBruines/go_buffer_bci_client/bufferbci"
)

func main() {
	log.Println("Trying to connect to localhost:1972 ...")

	conn, err := bufferbci.Connect("localhost:1972")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("Connection established")

	header, err := conn.GetHeader()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Header:", header)
}
