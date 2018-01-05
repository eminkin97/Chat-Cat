package main

import(
	"net"
	"log"
)

type client_struct struct {
	name string
	conn net.Conn
}

var connectedClients []client_struct

func handleRequest(c net.Conn) {
	log.Println("REQUEST HANDLED BRA")

	//read name from connection
	b := make([]byte, 10)
	n, err := c.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	name := string(b[:n])

	//create struct for client
	new_client := client_struct{name: name, conn: c}

	//append new client to connectedClients
	connectedClients = append(connectedClients, new_client)
	log.Println("REQUEST DONE BRA")
}

func main() {
	//listen for connections
	ln, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatal(err)
	}

	//accept connections and handle them with a seperate goroutine
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go handleRequest(conn)
	}
}
