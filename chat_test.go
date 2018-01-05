package main

import(
	"testing"
	"net"
	"log"
	"strings"
)

//check connection and that name is being added correctly
func TestClientConnect(t *testing.T) {
	server, client := net.Pipe()
	done := make(chan bool, 1)

	go func(c net.Conn, t *testing.T, done chan bool) {
		//this routine acts as the server in the test
		handleRequest(c)

		done <- true
		c.Close()
		log.Println("MEOW")
	}(server, t, done)

	//write from client connection to server connection
	name := []byte("rudolph")

	numWritten, err := client.Write(name)
	if err != nil {
		t.Error(err)
	}

	if numWritten != 7 {
		t.Logf("numwritten was not equal to y. Was %d", numWritten)
	}

	//synchronize via channel
	<-done

	//this routine acts as client in the test
	client.Close()

	//check to see if appended local ip address to connected clients slice
	if len(connectedClients) <= 0 {
		t.Error("connected Clients is empty")
	}

	//check to see if name was added correctly
	if (strings.Compare(connectedClients[0].name, "rudolph") != 0) {
		t.Errorf("Name was not equal to rudolph. Name was: %s", connectedClients[0].name)
	}
}
