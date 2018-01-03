package main

import(
	"testing"
	"net"
)

func TestClientConnect(t *testing.T) {
	server, client := net.Pipe()
	done := make(chan bool, 1)

	go func(c net.Conn, t *testing.T, done chan bool) {
		//this routine acts as the server in the test
		handleRequest(c)

		if len(connectedClients) <= 0 {
			t.Error("YOO connected Clients is empty")
		}

		done <- true
	}(server, t, done)

	//synchronize via channel
	<-done

	//this routine acts as client in the test
	client.Close()

	//check to see if appended local ip address to connected clients slice
	if len(connectedClients) <= 0 {
		t.Error("connected Clients is empty")
	}
}
