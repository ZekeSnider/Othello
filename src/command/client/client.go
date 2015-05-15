//Zeke Snider

package main

import (
	"bufio"
    "fmt"
    "net"
    "os"
    "log"
    "strings"
    "strconv"
)

type HostGame struct {
	Name string
	Color byte
}

type BoardLayout [10][10]byte
var board BoardLayout = BoardLayout{}

func main() {

	mainMenu()


}

func listGames() {
	conn, err := net.Dial("tcp", "localhost:8080")

	checkForError(err)

	defer conn.Close()

	conn.Write([]byte("LISTGAME"))

	response := make([]byte, 120)
	_, err = conn.Read(response)

	responseString := string(response[:120])

	if strings.HasPrefix(responseString, "GAMELIST") {
		numElements, err := strconv.Atoi(responseString[9:10])
		checkForError(err)

		for i := 0; i<numElements; i++ {
			fmt.Printf("%v. %v %v\n", i+1, responseString[11+i*28:36+i*28], responseString[37+i*28:38+i*28])
		}
		
	}
	

	conn.Close()

}

func getName() string {
	inputReader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("\n\nEnter your name: ")
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) < 25 || len(text) != 0 {
			name := text
			return name
		}
	}
	return "bob"
}

func getGameNumber() int {
	inputReader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("\n\nWhat game would you like to join?: ")
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)
		textInt, err := strconv.Atoi(text)
		checkForError(err)
		if textInt<=8 && textInt>0 {
			return textInt
		}
	}
	return 1
}

func getColor() byte {
	inputReader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("\n\nChoose your color (W for white, B for Black): ")
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)
		if text == "W" {
			playerColor := byte('W')
			return playerColor
		} else if text == "B" {
			playerColor := byte('B')
			return playerColor
		}
	}
	return byte('B')
}
func joinGame(gameNumber int, name string) {

	//check if gamenumber size is greater than 8
	serverConnection, err := net.Dial("tcp", "localhost:8080")

	checkForError(err)

	defer serverConnection.Close()
	
	requestString := fmt.Sprintf("JOINGAME %25v %v", name, gameNumber)
	fmt.Printf("Game number is.. %v ",requestString)

	serverConnection.Write([]byte(requestString))
	response := make([]byte, 120)

	fmt.Printf("Joining... \n")

	_, _ = serverConnection.Read(response)
	responseString := string(response[:120])

	if strings.HasPrefix(responseString, "GAMEJOIN") {
		playerColor := responseString[9:10]
		playGame(playerColor, serverConnection)
	} else {
		fmt.Printf("ERROR: That game cannot be joined!")
		os.Exit(0)
	}

}
func hostGame(playerColor byte, name string) {
	serverConnection, err := net.Dial("tcp", "localhost:8080")

	checkForError(err)

	defer serverConnection.Close()
	
	requestString := fmt.Sprintf("HOSTGAME %25v %v", name, string(playerColor))
	fmt.Printf("your color is.. %v", playerColor)

	serverConnection.Write([]byte(requestString))
	response := make([]byte, 120)

	fmt.Printf("Waiting for pair...\n")

	_, _ = serverConnection.Read(response)
	responseString := string(response[:120])
	fmt.Printf("Got a Response!!!! %v", responseString)

	if strings.HasPrefix(responseString, "GAMEPAIR") {
		playGame(string(playerColor), serverConnection)
	} else if strings.HasPrefix(responseString, "HOSTFULL") {
		fmt.Printf("ERROR: No more host slots available!!")
		os.Exit(0)
	}


}

func playGame(playerColor string, serverConnection net.Conn) {
	defer serverConnection.Close()
	startBoard()
	turnCount := 1
	endGame := false
	for endGame != true {
		if (playerColor == "B" && turnCount % 2 == 1) || (playerColor == "W" && turnCount % 2 == 0) {
			var move string
			moveOK := false
			for moveOK != true {
				printBoard()
				inputReader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter your move: ")
				move, _ = inputReader.ReadString('\n')
				fmt.Printf("move = %v\n", move)

				moveOK = validateInput(move[0:2],byte(playerColor[0]))
			}

			fmt.Printf("move = %v\n", move)
			MoveDoneMessageString := fmt.Sprintf("DOMOVE %v", move[0:2])
			fmt.Printf("sending... %v", MoveDoneMessageString)
			serverConnection.Write([]byte(MoveDoneMessageString))

			turnCount++

		} else {
			printBoard()
			fmt.Print("Waiting for opponent to move...\n")
			playerMoveRequest := make([]byte, 120)
			
			_, err := serverConnection.Read(playerMoveRequest)
			checkForError(err)
			
			fullMoveString := string(playerMoveRequest[:120])
			moveString := fullMoveString[7:9]
			fmt.Printf("movestring = '%v' ", moveString)

			if playerColor == "W" {
				_ = validateInput(moveString, 'B')
			} else if playerColor == "B" {
				_ = validateInput(moveString, 'W')
			}

			turnCount++
		}

		if checkGameOver() {
			endGame = true
			break
		}
	}


}
func checkGameOver() bool{
	blackCount := 0
	whiteCount := 0
	for i := range board {
		for j := range board[i] {
			if board[i][j] == 'W' {
				whiteCount++
			} else if board[i][j] == 'B' {
				blackCount++
			} else {
				return false
			}
		}
	}
	fmt.Printf("Game over! ")
	if whiteCount > blackCount {
		fmt.Printf("White wins! ")
	} else { 
		fmt.Printf("Black wins! ")
	}
	fmt.Printf("Final score: b:%d w:%d", whiteCount, blackCount)
	return true

}
func mainMenu() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to continue... ")
	_, _ = inputReader.ReadString('\n')
	
	for true {
		fmt.Print("\nWelcome! Press 1 to list games, 2 to join a game, 3 to host a game! 4 to quit")
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)

		textInt, err := strconv.Atoi(text)
		checkForError(err)
		if err == nil {

			if textInt == 1 {
				listGames()
			} else if textInt ==2 {
				name := getName()
				gameNumber := getGameNumber()
				joinGame(gameNumber, name)
			} else if textInt == 3 {
				name := getName()
				color := getColor()
				hostGame(color, name)
			}else if textInt == 4 {
				os.Exit(0)
			}
		}
	}
}



func checkForError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func startBoard() {
	for i := range board {
		for j := range board[i] {
			board[i][j] = 0
		}
	}
	board[4][4] = 'B'
	board[4][5] = 'W'
	board[5][4] = 'B'
	board[5][5] = 'W'
}

func printBoard() {
	fmt.Printf("\n   A B C D E F G H I J\n")
	fmt.Printf("  ┌-------------------\n")

	for i := range board {
 		fmt.Printf("%d |",i)
		for j := range board[i] {
			fmt.Printf("%c ", board[i][j])
		}
		fmt.Print("\n")
	}
}

func validateInput(moveInput string, playerColor byte) bool{
	
	fmt.Printf("moveinput[0] = %v, moveInput[1] = %v", moveInput[0], moveInput[1])
	fmt.Printf(" len = %v ",len(moveInput))
	if len(moveInput)!=2 {
		fmt.Printf("fail1")
		return false
	} else if int(moveInput[0])>=75 || int(moveInput[0])<65 {
		fmt.Printf("fail2")
		return false
	} else if (int(moveInput[1])-48) < 0 || (int(moveInput[1])-48) > 9 {
		fmt.Printf("fail3")
		return false
	} 
	fmt.Printf("made it here..")

	changeCount := 0
	for i:=1; i<=8; i++ {
		changeCount += doMove( int(moveInput[1])-48, int(moveInput[0])-65, i, playerColor)
	}

	if changeCount == 0 {
		return false
	} else {

		if board[int(moveInput[1])-48][int(moveInput[0])-65] != playerColor {
			board[int(moveInput[1])-48][int(moveInput[0])-65] = playerColor
			changeCount++
		}
		return true
	}
}

func doMove(Starti int, Startj int, testType int, playerColor byte, ) int {
	var Currenti int
	var Currentj int
	changeCount := 0

	if testType == 1 {
		Currenti = Starti + 1
		Currentj = Startj
	} else if testType == 2 {
		Currenti = Starti - 1
		Currentj = Startj
	} else if testType == 3 {
		Currenti = Starti
		Currentj = Startj - 1
	} else if testType == 4 {
		Currenti = Starti
		Currentj = Startj + 1
	} else if testType == 5 {
		Currenti = Starti - 1
		Currentj = Startj - 1
	} else if testType == 6 {
		Currenti = Starti + 1
		Currentj = Startj - 1
	} else if testType == 7 {
		Currenti = Starti + 1
		Currentj = Startj + 1
	} else if testType == 8 {
		Currenti = Starti + 1
		Currentj = Startj - 1
	}

	for Currenti <= 9 && Currenti >= 0 && Currentj <= 9 && Currentj >= 0 && board[Currenti][Currentj] != '0'{
		fmt.Printf("%d %d\n", Currenti,Currentj)
		if board[Currenti][Currentj] == playerColor {
			if testType == 1 {
				for i:= Starti; i<Currenti; i++ {
					if board[i][Currentj] != playerColor && board[i][Currentj] != '0' {
						board[i][Currentj] = playerColor
						changeCount++
					}
				}
			} else if testType == 2 {
				for i:= Starti; i>Currenti; i-- {
					if board[i][Currentj] != playerColor && board[i][Currentj] != '0' {
						board[i][Currentj] = playerColor
						changeCount++
					}
				}
			} else if testType == 3 {
				for j:= Startj; j>Currentj; j-- {
					if board[Currentj][j] != playerColor && board[Currentj][j] != '0'{
						board[Currentj][j] = playerColor
						changeCount++
					}
				}
			} else if testType == 4 {
				for j:= Startj; j<Currentj; j++ {
					if board[Currentj][j] != playerColor && board[Currentj][j] != '0'{
						board[Currentj][j] = playerColor
						changeCount++
					}
				}
			} else if testType == 5 {
				i := Starti
				j := Startj
				for i!= Currenti && j != Currentj && i >= 0 && i<=9 && j >= 0 && j<=9 {
					if board[i][j] != playerColor && board[i][j] != '0'{
						board[i][j] = playerColor
						changeCount++
					}
					i--
					j--
				}
			} else if testType == 6 {
				i := Starti
				j := Startj
				for i!= Currenti && j != Currentj && i >= 0 && i<=9 && j >= 0 && j<=9 {
					fmt.Printf("i: %d j: %d\n",i,j)
					if board[i][j] != playerColor && board[i][j] != '0'{
						board[i][j] = playerColor
						changeCount++
					}
					i--
					j++
				}
			} else if testType == 7 {
				i := Starti
				j := Startj
				for i!= Currenti && j != Currentj && i >= 0 && i<=9 && j >= 0 && j<=9 {
					if board[i][j] != playerColor && board[i][j] != '0'{
						board[i][j] = playerColor
						changeCount++
					}
					i++
					j++
				}
			} else if testType == 8 {
				i := Starti
				j := Startj
				for i!= Currenti && j != Currentj && i >= 0 && i<=9 && j >= 0 && j<=9 {
					if board[i][j] != playerColor && board[i][j] != '0'{
						board[i][j] = playerColor
						changeCount++
					}
					i--
					j++
				}
			}
			return changeCount
		}

		if testType == 1 {
			Currenti++
		} else if testType == 2 {
			Currenti--
		} else if testType == 3 {
			Currentj--
		} else if testType == 4 {
			Currentj++
		} else if testType == 5 {
			Currenti++
			Currentj--
		} else if testType == 6 {
			Currenti--
			Currentj++
		} else if testType == 7 {
			Currenti++
			Currentj++
		} else if testType == 8 {
			Currenti--
			Currentj++
		}
	}
	return changeCount
}

