package main

import(
	"net"
	"log"
	"bytes"
)

type client_struct struct {
	name string
	ch chan string
}

var connectedClients []client_struct

//function returns concatenated list of client names seperated by comma
func concatenateClientNames(name string) string {
	var buffer bytes.Buffer

	//append all names of clients to the client_name_list
	for _, elem := range connectedClients {
		if (elem.name == name) {
			//dont include name of the client
			continue
		}
		buffer.WriteString(elem.name)
		buffer.WriteString(",")
	}

	return buffer.String()
}

func parseChatRequest(str string) (string, string) {
	//look for first colon and seperate based on that
	var name string
	var msg string

	for i, elem := range str {
		if elem == ':' {
			name = str[:i]
			msg = str[i+1:]
			break
		}
	}

	return name, msg
}

func handleRequest(c net.Conn) {
	log.Println("REQUEST HANDLED BRA")
	validName := false
	var name string

	for (!validName) {
		//read from connection
		b := make([]byte, 10)
		n, err := c.Read(b)

		if err != nil {
			log.Fatal(err)
		}

		//client logging in
		name = string(b[:n])

		//make sure name isn't already taken in connectedClients
		validName = true
		for _, elem := range connectedClients {
			log.Printf("elem.name: %s, name: %s\n", elem.name, name)
			if (elem.name == name) {
				//name already taken
				validName = false

				_, err := c.Write([]byte("n"))		//indicates invalid name

				if err != nil {
					log.Fatal(err)
				}
				break
			}
		}
	}

	//write back to client to indicate name was accepted
	_, err := c.Write([]byte("v"))	//indicates valid name
	if err != nil {
		log.Fatal(err)
	}

	//create channel for chat with other users
	ch := make(chan string)

	//create struct for client
	new_client := client_struct{name: name, ch: ch}

	//append new client to connectedClients
	log.Println("APPENDING TO CONNECTED CLIENTS")
	connectedClients = append(connectedClients, new_client)
	log.Println(connectedClients[0].name)

	//channel for reporting action user wants
	ac := make(chan string)

	go func(ac chan string) {
		for {
			//action to be performed from client
			a := make([]byte, 250)

			//read message from client
			n1, err := c.Read(a)
			if err != nil {
				log.Fatal(err)
			}

			ac <- string(a[:n1])

			status := <-ac
			if status == "e" {
				break
			}
		}
	}(ac)


	//client is done once he/she sends exit message
	for {

		//check to see if received a chat request from other user or action from client
		var fullaction string

		select {
			case incomingChat := <-ch:
				//write to connection
				_, err = c.Write([]byte("c:"+incomingChat))
				if err != nil {
					log.Fatal(err)
				}

				continue
			case fullaction = <-ac:
				//action by the user
				log.Printf("Action by the user: %s\n", fullaction)
		}

		action := fullaction[:4]

		if action == "list" {
			//write back list of connected clients
			client_name_list := concatenateClientNames(name)

			_, err = c.Write([]byte("l:"+client_name_list))
			if err != nil {
				log.Fatal(err)
			}
			ac <- "c"	//continue
		} else if action == "chat" {
			//parse user to chat and msg to send
			nameToChat, msg := parseChatRequest(fullaction[5:])

			//look for name to chat in connected clients
			for _, elem := range connectedClients {
				if elem.name == nameToChat {
					//found entry corresponding to name request
					//send message through channel to indicate chat request
					elem.ch <- (name + ":" + msg)
					break
				}
			}

			ac <- "c"	//continue
		} else if action == "exit" {
			//remove name from connected clients then break
			for i, elem := range connectedClients {
				if elem.name == name {
					//found the one to remove
					connectedClients = append(connectedClients[0:i], connectedClients[i+1:]...)
					break
				}
			}
			ac <- "e"	//exit
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
