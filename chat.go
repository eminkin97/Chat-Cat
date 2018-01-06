package main

import(
	"net"
	"log"
	"bytes"
)

type client_struct struct {
	name string
	conn net.Conn
	ch chan string
}

var connectedClients []client_struct

//function returns concatenated list of client names seperated by comma
func concatenateClientNames() string {
	var buffer bytes.Buffer

	//append all names of clients to the client_name_list
	for _, elem := range connectedClients {
		buffer.WriteString(elem.name)
		buffer.WriteString(",")
	}

	return buffer.String()
}

func handleRequest(c net.Conn) {
	log.Println("REQUEST HANDLED BRA")

	//read from connection
	b := make([]byte, 10)
	n, err := c.Read(b)

	if err != nil {
		log.Fatal(err)
	}

	//client logging in
	name := string(b[:n])

	//create channel for hat with other users
	ch := make(chan string)

	//create struct for client
	new_client := client_struct{name: name, conn: c, ch: ch}

	//append new client to connectedClients
	connectedClients = append(connectedClients, new_client)

	//client is done once he/she sends exit message
	for {

		//check to see if received a chat request from other user
		/*select {
			case externalRequest := <-c2:
				log.Println(res)
			case <-time.After(time.Second * .1):
				log.Println("no chat request")
		}*/

		//action to be performed from client
		a := make([]byte, 20)

		//read message from client
		n1, err := c.Read(a)
		if err != nil {
			log.Fatal(err)
		}

		action := string(a[:4])

		if action == "list" {
			//write back list of connected clients
			client_name_list := concatenateClientNames()

			_, err = c.Write([]byte(client_name_list))
			if err != nil {
				log.Fatal(err)
			}
		} else if action == "chat" {
			//start chat with name nameToChat
			nameToChat := string(a[5:n1])
			log.Println(nameToChat)
		} else if action == "exit" {
			//remove name from connected clients then break
			for i, elem := range connectedClients {
				if elem.name == name {
					//found the one to remove
					connectedClients = append(connectedClients[0:i], connectedClients[i+1:]...)
					break
				}
			}
			break
		}
	}

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
