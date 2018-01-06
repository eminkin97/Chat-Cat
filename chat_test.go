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

	//check to see if appended local ip address to connected clients slice
	if len(connectedClients) <= 0 {
		t.Error("connected Clients is empty")
	}

	//check to see if name was added correctly
	if (strings.Compare(connectedClients[0].name, "rudolph") != 0) {
		t.Errorf("Name was not equal to rudolph. Name was: %s", connectedClients[0].name)
	}

	//write exit action to indicate we're done
	action := []byte("exit")

	_, err = client.Write(action)
	if err != nil {
		t.Error(err)
	}

	//synchronize via channel
	<-done

	for _, elem := range connectedClients {
		if elem.name == "rudolph" {
			t.Error("NAME WASN'T REMOVED AT END")
		}
	}

	//this routine acts as client in the test
	client.Close()
}


func TestListConnectedClients(t *testing.T) {
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

	_, err := client.Write(name)
	if err != nil {
		t.Error(err)
	}

	//write list action to indicate want list of all clients
	action := []byte("list")

	_, err = client.Write(action)
	if err != nil {
		t.Error(err)
	}

	returned_bytes := make([]byte, 100)
	n, err := client.Read(returned_bytes)

	client_list := string(returned_bytes[:n])

	if client_list != "rudolph," {
		t.Error("Client list was wrong")
	}

	//write exit action to indicate we're done
	action = []byte("exit")

	_, err = client.Write(action)
	if err != nil {
		t.Error(err)
	}


	//synchronize via channel
	<-done

	//this routine acts as client in the test
	client.Close()

}
