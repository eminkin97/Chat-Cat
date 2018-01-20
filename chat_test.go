package main

import(
	"testing"
	"net"
	"log"
	"strings"
	"time"
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
		log.Println("1")
		t.Error(err)
	}

	if numWritten != 7 {
		log.Println("2")
		t.Logf("numwritten was not equal to y. Was %d", numWritten)
	}

	//make sure got valid
	name_status := make([]byte, 100)
	n, err := client.Read(name_status)
	if err != nil {
		log.Println("13.5")
		t.Error(err)
	}

	if (string(name_status[:n]) != "v") {
		t.Error("should have returned valid")
	}

	time.Sleep(time.Second * 2)

	//check to see if appended local ip address to connected clients slice
	log.Println(len(connectedClients))
	if len(connectedClients) <= 0 {
		log.Println("3")
		t.Error("connected Clients is empty")
	}

	//check to see if name was added correctly
	if (strings.Compare(connectedClients[0].name, "rudolph") != 0) {
		log.Println("4")
		t.Errorf("Name was not equal to rudolph. Name was: %s", connectedClients[0].name)
	}

	//write exit action to indicate we're done
	action := []byte("exit")

	_, err = client.Write(action)
	if err != nil {
		log.Println("5")
		t.Error(err)
	}

	//synchronize via channel
	<-done

	for _, elem := range connectedClients {
		if elem.name == "rudolph" {
			log.Println("6")
			t.Error("NAME WASN'T REMOVED AT END")
		}
	}

	//this routine acts as client in the test
	client.Close()
}

func listClientsTest(c net.Conn, client_name string, t *testing.T) {
	//write list action to indicate want list of all clients
	action := []byte("list")

	_, err := c.Write(action)
	if err != nil {
		log.Println("11")
		t.Error(err)
	}

	returned_bytes := make([]byte, 100)
	n, err := c.Read(returned_bytes)

	client_list := string(returned_bytes[2:n])

	if client_name == "jack" {
		if client_list != "rudolph," {
			log.Println("12")
			t.Errorf("Client list was wrong: %s\n", client_list)
		}
	} else if client_name == "rudolph" {
		if client_list != "jack," {
			log.Println("12")
			t.Errorf("Client list was wrong: %s\n", client_list)
		}
	}
}

func TestChatWithTwoClients(t *testing.T) {
	t.Log("STARTING CHAT WITH 2 CLIENTS TEST")
	server1, client1 := net.Pipe()
	server2, client2 := net.Pipe()

	done1 := make(chan bool, 1)
	done2 := make(chan bool, 1)

	go func(c net.Conn, t *testing.T, done chan bool) {
		//this routine acts as the server in the test
		handleRequest(c)

		done <- true
		c.Close()
		log.Println("MEOW")
	}(server1, t, done1)

	go func(c net.Conn, t *testing.T, done chan bool) {
		//this routine acts as the server in the test
		handleRequest(c)

		done <- true
		c.Close()
		log.Println("MEOW")
	}(server2, t, done2)

	//write from client connection to server connection
	name1 := []byte("rudolph")

	_, err := client1.Write(name1)
	if err != nil {
		log.Println("13")
		t.Error(err)
	}

	//make sure got valid
	name_status := make([]byte, 100)
	n, err := client1.Read(name_status)
	if err != nil {
		log.Println("13.5")
		t.Error(err)
	}

	if (string(name_status[:n]) != "v") {
		t.Error("should have returned valid")
	}

	name2 := []byte("jack")

	_, err = client2.Write(name2)
	if err != nil {
		log.Println("14")
		t.Error(err)
	}

	//make sure got valid
	name_status = make([]byte, 100)
	n, err = client2.Read(name_status)
	if err != nil {
		log.Println("13.5")
		t.Error(err)
	}

	if (string(name_status[:n]) != "v") {
		t.Error("should have returned valid")
	}

	listClientsTest(client1, "rudolph", t)
	listClientsTest(client2, "jack", t)

	//text chat from rudolph to jack
	action := []byte("chat:jack:meow")

	_, err = client1.Write(action)
	if err != nil {
		log.Println("15")
		t.Error(err)
	}

	//check to see if jack received it
	returned_bytes := make([]byte, 100)
	n, err = client2.Read(returned_bytes)
	if err != nil {
		log.Println("16")
		t.Error(err)
	}

	msg := string(returned_bytes[:n])
	if (msg != "c:rudolph:meow") {
		log.Println("17")
		t.Errorf("Message jack got was not correct. It was %s", msg)
	}

	//text chat from jack to rudolph
	action = []byte("chat:rudolph:bark")

	_, err = client2.Write(action)
	if err != nil {
		log.Println("15")
		t.Error(err)
	}

	//check to see if jack received it
	returned_bytes = make([]byte, 100)
	n, err = client1.Read(returned_bytes)
	if err != nil {
		log.Println("16")
		t.Error(err)
	}

	msg = string(returned_bytes[:n])
	if (msg != "c:jack:bark") {
		log.Println("17")
		t.Errorf("Message rudolph got was not correct. It was %s", msg)
	}

	//write exit action to indicate we're done
	action = []byte("exit")

	_, err = client1.Write(action)
	if err != nil {
		log.Println("18")
		t.Error(err)
	}

	_, err = client2.Write(action)
	if err != nil {
		log.Println("19")
		t.Error(err)
	}


	//synchronize via channel
	<-done1
	<-done2

	//this routine acts as client in the test
	client1.Close()
	client2.Close()
}

func TestSameName(t *testing.T) {
	server1, client1 := net.Pipe()
	server2, client2 := net.Pipe()

	done1 := make(chan bool, 1)
	done2 := make(chan bool, 1)

	go func(c net.Conn, t *testing.T, done chan bool) {
		//this routine acts as the server in the test
		handleRequest(c)

		done <- true
		c.Close()
		log.Println("MEOW")
	}(server1, t, done1)

	go func(c net.Conn, t *testing.T, done chan bool) {
		//this routine acts as the server in the test
		handleRequest(c)

		done <- true
		c.Close()
		log.Println("MEOW")
	}(server2, t, done2)

	//write from client connection to server connection
	name1 := []byte("rudolph")

	_, err := client1.Write(name1)
	if err != nil {
		log.Println("13")
		t.Error(err)
	}

	//make sure got valid
	status := make([]byte, 100)
	n, err := client1.Read(status)
	if err != nil {
		log.Println("13.5")
		t.Error(err)
	}

	if (string(status[:n]) != "v") {
		t.Error("should have returned valid")
	}

	time.Sleep(time.Second * 2)

	//write same name from different client
	name2 := []byte("rudolph")

	_, err = client2.Write(name2)
	if err != nil {
		log.Println("14")
		t.Error(err)
	}

	log.Println("MEOW2")
	//make sure got invalid
	status = make([]byte, 100)
	n, err = client2.Read(status)
	if err != nil {
		log.Println("14.5")
		t.Error(err)
	}

	log.Println(string(status[:n]))
	if (string(status[:n]) != "n") {
		t.Error("should have returned invalid")
	}

	log.Println("MEOW3")
	//write correct name for client 2
	name2 = []byte("jack")

	_, err = client2.Write(name2)
	if err != nil {
		log.Println("14")
		t.Error(err)
	}

	log.Println("MEOW4")
	//make sure got valid
	status = make([]byte, 100)
	n, err = client2.Read(status)
	if err != nil {
		log.Println("14.5")
		t.Error(err)
	}

	log.Println("MEOW5")

	if (string(status[:n]) != "v") {
		t.Error("should have returned valid pt 2")
	}

	//write exit action to indicate we're done
	action := []byte("exit")

	_, err = client1.Write(action)
	if err != nil {
		log.Println("18")
		t.Error(err)
	}

	_, err = client2.Write(action)
	if err != nil {
		log.Println("19")
		t.Error(err)
	}


	//synchronize via channel
	<-done1
	<-done2

	//this routine acts as client in the test
	client1.Close()
	client2.Close()
}
