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

	//establish a socket
	TCPAddress, err := net.ResolveTCPAddr("tcp4", ":8080")
	checkForError(err)
	listener, err := net.ListenTCP("tcp", TCPAddress)
	checkForError(err)

	//repeat checking for a new connection
	for true {
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

//handle a connection
func handleConnection(conn net.Conn) {

	//get a request
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

	//check which type of request it is, and handle accordingly
	if strings.HasPrefix(requestString, "LISTGAME") {
		listGames(conn)
	} else if strings.HasPrefix(requestString, "HOSTGAME") {
		name := requestString[9:34]
		color := requestString[35:36]
		hostGame(name, color, conn)
	} else if strings.HasPrefix(requestString, "JOINGAME") {
		name := requestString[9:34]
		gameNumber := requestString[35:36]

		fmt.Printf("JOINGAME %v \n",gameNumber)

		gameNumberInt, err := strconv.Atoi(gameNumber)
		checkForError(err)
		joinGame(name, gameNumberInt, conn)
	}
}

//handles a game between two players
func handleGame(blackPlayer net.Conn, whitePlayer net.Conn) {
	defer blackPlayer.Close()
	defer whitePlayer.Close()
	turnCount := 1
	endGame := false

	//repeat until the game is over. 
	for !endGame {
		//if it's the black player's turn
		if turnCount % 2 == 1 {

			//read the request
			moveRequest := make([]byte, 120)
			_, err := blackPlayer.Read(moveRequest)

			//if the client disconnected, let the other player know, and stop this thread.
			if (err == io.EOF) {
				Move := "ZZ"
				MoveDoneMessageString := fmt.Sprintf("MOVEDONE %v", Move)
				whitePlayer.Write([]byte(MoveDoneMessageString))
				whitePlayer.Close()
				
				endGame = true
				break
			}

			//parse the move, create request to send to the other player
			Move := string(moveRequest[7:9])
			MoveDoneMessageString := fmt.Sprintf("MOVEDONE %v", Move)
			fmt.Printf("%v\n", MoveDoneMessageString)

			//write the move done by the blackplayer to the whiteplayer
			whitePlayer.Write([]byte(MoveDoneMessageString))

			//advance the turn count
			turnCount++

			//if the move is YY, the game is over. 
			if Move=="YY" {
				endGame = true
				break
			}

		//same functionality for the white player's turn. roles are reversed
		} else {

			moveRequest := make([]byte, 120)
			_, err := whitePlayer.Read(moveRequest)

			if (err == io.EOF) {
				Move := "ZZ"
				MoveDoneMessageString := fmt.Sprintf("MOVEDONE %v", Move)
				blackPlayer.Write([]byte(MoveDoneMessageString))
				blackPlayer.Close()

				endGame = true
				break
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

//joins a player to another player's game that is being joined
func joinGame(name string, gameNumber int, joiner net.Conn) {

	//for going through host list
	position := 1
	current := hostList.Front()

	//fmt.Printf("game no = %v", gameNumber)

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


	//send a response back to the joiner with their color
	if currentGame.Color == "B" {
		joinResponse = fmt.Sprintf("GAMEJOIN W")
	} else if currentGame.Color == "W" {
		joinResponse = fmt.Sprintf("GAMEJOIN B")
	}
	joiner.Write([]byte(joinResponse))

	//send a response to the host with the joiner's name 
	hostResponse = fmt.Sprintf("GAMEPAIR %25v", name)
	currentGame.Socket.Write([]byte(hostResponse))

	//call handlegame with correct order on sockets for B/W colors
	if currentGame.Color == "B" {
		handleGame(currentGame.Socket, joiner)
	} else {
		handleGame(joiner, currentGame.Socket)
	}
}

//add the host to the list of hosts
func hostGame(name string, color string, conn net.Conn){
	hostList.PushBack(HostGame{name, color, conn})

}

//lists the games and send it back to the client
func listGames(conn net.Conn){
	gameList := fmt.Sprintf("GAMELIST %v ", hostList.Len())

	//add all hosts's names and colors to the response with padding
	for e := hostList.Front(); e != nil; e = e.Next() {
		var currentGame HostGame
		currentGame = e.Value.(HostGame)
		gameList = fmt.Sprintf("%v%v %v ",gameList, currentGame.Name, currentGame.Color)
	}

	//writ the response, then close the connection
	conn.Write([]byte(gameList))
	conn.Close()

	return
}
