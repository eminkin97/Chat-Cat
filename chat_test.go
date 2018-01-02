package main

import(
	"testing"
	"net"
)

func TestConnect(t *testing.T) {
	//client connects to server
	_, err := net.Dial("tcp", ":8000")

	if err != nil {
		t.Error("Error on Connection")
	}

	/*	//write to connection
	var b []byte
	b = []byte("1")		//indicates ip of user to connect too
	_, err := conn.Write(b)
	if err != nil {
		t.Error("Error on write to connection")
	}
	*/
}
