package main

import(
	"net"
	"log"
)

var connectedClients []string

func handleRequest(c net.Conn) {
	log.Println("REQUEST HANDLED BRA")

	//get IP address of client
	remoteaddr := c.RemoteAddr().String()

	//append IP address of client to connected clients slice
	connectedClients = append(connectedClients, remoteaddr)
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
