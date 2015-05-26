//Zeke Snider

package main

import (
	//"io"
	"log"
	"net"
	"fmt"
	//"time"
	"os"
	"io"
	"strings"
	"container/list"
    "strconv"
)

type HostGame struct {
	Name string
	Color string
	Socket net.Conn
}

var hostList = list.New()

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
		listGames(conn)
	} else if strings.HasPrefix(requestString, "HOSTGAME") {
		name := requestString[9:34]
		color := requestString[35:36]
		hostGame(name, color, conn)
	} else if strings.HasPrefix(requestString, "JOINGAME") {
		name := requestString[9:34]
		gameNumber := requestString[35:36]

		fmt.Printf(".. %v ",gameNumber)

		gameNumberInt, err := strconv.Atoi(gameNumber)
		checkForError(err)

		joinGame(name, gameNumberInt, conn)
	}





}

func handleGame(blackPlayer net.Conn, whitePlayer net.Conn) {
	defer blackPlayer.Close()
	defer whitePlayer.Close()
	turnCount := 1
	endGame := false
	for !endGame {
		if turnCount % 2 == 1 {

			moveRequest := make([]byte, 120)

			_, err := blackPlayer.Read(moveRequest)

			//if the client disconnected, let the other player know, and stop this thread.
			if (err == io.EOF) {
				Move := "ZZ"
				MoveDoneMessageString := fmt.Sprintf("MOVEDONE %v", Move)
				whitePlayer.Write([]byte(MoveDoneMessageString))
				whitePlayer.Close()
				os.Exit(0)
			}

			Move := string(moveRequest[7:9])

			MoveDoneMessageString := fmt.Sprintf("MOVEDONE %v", Move)
			fmt.Printf("%v\n", MoveDoneMessageString)

			whitePlayer.Write([]byte(MoveDoneMessageString))

			turnCount++

			if Move=="YY" {
				endGame = true
				break
			}

		} else {

			moveRequest := make([]byte, 120)


			_, err := whitePlayer.Read(moveRequest)

			//if the client disconnected, let the other player know, and stop this thread.
			if (err == io.EOF) {
				Move := "ZZ"
				MoveDoneMessageString := fmt.Sprintf("MOVEDONE %v", Move)
				blackPlayer.Write([]byte(MoveDoneMessageString))
				blackPlayer.Close()
				os.Exit(0)
			}

			Move := string(moveRequest[7:9])

			MoveDoneMessageString := fmt.Sprintf("MOVEDONE %v", Move)
			fmt.Printf("%v\n",MoveDoneMessageString)

			blackPlayer.Write([]byte(MoveDoneMessageString))

			turnCount++

			if Move=="YY" {
				endGame = true
				break
			}
		}
	}
}

func joinGame(name string, gameNumber int, joiner net.Conn) {
	position := 1
	current := hostList.Front()

	fmt.Printf("game no = %v", gameNumber)

	//searching the list for the specified game number
	if position != gameNumber {
		for ; position != gameNumber; current= current.Next() {
			if current != nil {
				position++
			}
		}
	}

	var currentGame HostGame
	var joinResponse string
	var hostResponse string

	//getting the value from the host list and removing it from the list so other players
	//won't try to join the same game. 
	currentGame = current.Value.(HostGame)
	hostList.Remove(current)



	if currentGame.Color == "B" {
		joinResponse = fmt.Sprintf("GAMEJOIN W")
	} else if currentGame.Color == "W" {
		joinResponse = fmt.Sprintf("GAMEJOIN B")
	}

	hostResponse = fmt.Sprintf("GAMEPAIR %25v", name)

	joiner.Write([]byte(joinResponse))
	currentGame.Socket.Write([]byte(hostResponse))

	if currentGame.Color == "B" {
		handleGame(currentGame.Socket, joiner)
	} else {
		handleGame(joiner, currentGame.Socket)
	}




}

func hostGame(name string, color string, conn net.Conn){
	hostList.PushBack(HostGame{name, color, conn})

}
func listGames(conn net.Conn){
	gameList := fmt.Sprintf("GAMELIST %v ", hostList.Len())

	for e := hostList.Front(); e != nil; e = e.Next() {
		var currentGame HostGame
		currentGame = e.Value.(HostGame)
		gameList = fmt.Sprintf("%v%v %v ",gameList, currentGame.Name, currentGame.Color)
	}

	conn.Write([]byte(gameList))

	conn.Close()
	return
}
