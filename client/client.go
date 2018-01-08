package main

import(
	"fmt"
	"log"
	"net"
	"bufio"
	"os"
	"strings"
)

//goroutine that accepts reads
func communicationWithServer(ch chan string, conn net.Conn) {
	for {
		data := make([]byte, 100)

		n, err := conn.Read(data)
		if err != nil {
			log.Fatal(err)
		}

		ch <- string(data[:n])
	}
}

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

	//channel to communicate with goroutine that accepts reads
	ch := make(chan string)
	go communicationWithServer(ch, conn)

	for {
		fmt.Printf("\nMENU:\n\n")
		fmt.Printf("1) See contacts currently online on network\n")
		fmt.Printf("2) Send chat request\n")
		fmt.Printf("3) Logoff\n")
		fmt.Printf("Enter Option: ")

		//read from standard input
		var option int
		_, err = fmt.Scanf("%d", &option)

		if err != nil {
			log.Fatal(err)
		}

		if option == 1 {
			//write list action to indicate want list of all clients
			action := []byte("list")

			_, err = conn.Write(action)
			if err != nil {
				log.Fatal(err)
			}

			//gets result. Splits result into list of all clients and then prints it
			res := <-ch
			clients := strings.Split(res, ",")

			for _, elem := range clients {
				fmt.Printf("%s", elem)
			}
		} else if option == 2 {
			fmt.Println("OPTION 2")
		} else if option == 3 {
			//Send close signal and exit
			fmt.Println("OPTION 3")
			//write exit action to indicate we're done
			action := []byte("exit")

			_, err = conn.Write(action)
			if err != nil {
				log.Fatal(err)
			}

			break

		}
	}
}
