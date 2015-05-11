//Zeke Snider

package main

import (
	//"io"
	"log"
	"net"
	"fmt"
	//"time"
	"io"
	"strings"
)

type HostGame struct {
	Name string
	Color byte
}

func main() {
	TCPAddress, err := net.ResolveTCPAddr("tcp4", ":8080")

	checkForError(err)

	listener, err := net.ListenTCP("tcp", TCPAddress)

	checkForError(err)

	for {
		conn, err := listener.Accept()
		
		checkForError(err)
		
		go handleConnection(conn)

	}

	
}

func checkForError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleConnection(conn net.Conn) {
	request := make([]byte, 120)
	defer conn.Close()

	_, err := conn.Read(request)

	requestString := string(request[:120])

	fmt.Printf("%v\n", requestString) 

	
	if err != nil {
		if err != io.EOF {
			fmt.Printf("error!!")
			log.Fatal(err)
		}
		//break
	}


	if strings.HasPrefix(requestString, "LISTGAME") {
		response = listGames()
	}

    conn.Write([]byte("GAMELIST  \n"))

    _, err = conn.Read(request)
    requestString = string(request[:120])
    fmt.Printf("%v\n", requestString) 


}

func handleGame(blackPlayer net.Conn, whitePlayer net.Conn) {
	defer blackPlayer.Close()
	defer whitePlayer.Close()
	turnCount := 1
	endGame := false
	while !endGame {
		if turnCount % 2 == 1 {

			moveRequest := make([]byte, 120)
			
			_, err := blackPlayer.Read(moveRequest)

			requestString := string(moveRequest[:120])

			Move := moveRequest[9:10]

			MoveDoneMessageString = strings.Join([]string{"MOVEDONE ", moveRequest}, "")

			whitePlayer.Write([]byte(MoveDoneMessageString))

			turnCount++

			if Move=="YY" {
				endGame = true
				break
			}

		} else {

			moveRequest := make([]byte, 120)
			
			_, err := whitePlayer.Read(moveRequest)

			requestString := string(moveRequest[:120])

			Move := moveRequest[9:10]

			MoveDoneMessageString = strings.Join([]string{"MOVEDONE ", moveRequest}, "")

			blackPlayer.Write([]byte(MoveDoneMessageString))

			turnCount++

			if Move=="YY" {
				endGame = true
				break
			}
		}
	}


}

func listGames() string {
	return "list the games!!"
}
