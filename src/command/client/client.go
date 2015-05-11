//Zeke Snider

package main

import (
	"bufio"
    "fmt"
    "net"
    "os"
    "log"
)

type HostGame struct {
	Name string
	Color byte
}

type BoardLayout [10][10]byte
var board BoardLayout = BoardLayout{}

func main() {

	listGames()


}

func listGames() {


	conn, err := net.Dial("tcp", "localhost:8080")

	checkForError(err)

	defer conn.Close()

	checkForError(err)

	conn.Write([]byte("LISTGAME"))

	response := make([]byte, 120)
	_, err = conn.Read(response)

	responseString := string(response[:120])

	fmt.Printf("%v", responseString)

	conn.Write([]byte("JOINGAME"))

}

func playGame(playerColor byte, serverConnection net.Conn) {
	defer serverConnection.Close()
	turnCount := 1
	endGame := false
	while endGame != true {
		if (playerColor == 'B' && turnCount % 2 == 1) || (playerColor == 'W' && turnCount % 2 == 0) {
			var move string
			moveOK := false
			while moveOK != true {
				printBoard()
				inputReader := bufio.NewReader(os.Stdin)
				fmt.Print("Enter your move: ")
				move, _ = inputReader.ReadString('\n')

				moveOK = validateInput(move,playerColor)
			}

			MoveDoneMessageString = strings.Join([]string{"DOMOVE ", move}, "")
			serverConnection.Write([]byte(MoveDoneMessageString))

			turnCount++

		} else {
			playerMoveRequest := make([]byte, 120)
			_, err := serverConnection.Read(playerMoveRequest)
			fullMoveString := string(playerMove[:120])
			moveString := fullMoveString[9:10]

			if playerColor == 'W' {
				_ = validateInput(moveString, 'B')
			} else if playerColor == 'B' {
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

func mainMenu() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to continue... ")
	_, _ := inputReader.ReadString('\n')
	
	while true {
		fmt.Print("\n\nWelcome! Press 1 to list games or 2 to join a game! 3 to quit")
		text, _ := inputReader.ReadString('\n')
		if text == 1 {
			listGames()
		} else if text ==2 {
			hostGame()
		} else if text == 3 {
			exit(0)
		}
	}
}


func checkForError(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
	
	if len(moveInput)!=2 {
		return false
	} else if int(moveInput[0])>=75 || int(moveInput[0])<65 {
		return false
	} else if (int(moveInput[1])-48) < 0 || (int(moveInput[1])-48) > 9 {
		return false
	} 

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

