package main

import(
	"testing"
	"net"
)

func connect(t *testing.T) {
	_, err := net.Dial("tcp", ":8001")

	if err != nil {
		t.Error(err)
	}
}

func TestConnect(t *testing.T) {
	//client connects to server
	go connect(t)

	/*	//write to connection
	var b []byte
	b = []byte("1")		//indicates ip of user to connect too
	_, err := conn.Write(b)
	if err != nil {
		t.Error("Error on write to connection")
	}
	*/
}
