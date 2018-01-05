package main

import(
	"fmt"
	"log"
	"net"
	"bufio"
	"os"
)

func main() {
	//connect to chat server
	conn, err := net.Dial("tcp", ":8001")

	if err != nil {
		log.Fatal("Error on Connection")
	}

	fmt.Printf("ENTER USERNAME: ")

	//read username from user input
	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	//write username to connection
	n, err := conn.Write([]byte(name))
	log.Printf("%d\n", n)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nMENU:\n\n")
	fmt.Printf("1) See Contacts on Network:\n")
	fmt.Printf("2) Send Chat Request:\n")
	fmt.Printf("Enter Option: ")

	//read from standard input
	var option int
	_, err = fmt.Scanf("%d", &option)

	if err != nil {
		log.Fatal(err)
	}

	if option == 1 {
		fmt.Println("OPTION A: Will Do Later")

	} else if option == 2 {
		fmt.Println("OPTION B")
	}
	fmt.Println(option)
}
