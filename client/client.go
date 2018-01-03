package main

import(
	"fmt"
	"log"
	"net"
)

func main() {
	_, err := net.Dial("tcp", ":8001")

	if err != nil {
		log.Fatal("Error on Connection")
	}


	fmt.Printf("MENU:\n\n")
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
