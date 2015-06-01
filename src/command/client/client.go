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


//stores a move
type Coordinate struct {
  PosX, PosY int
}

//stores a game
type HostGame struct {
	Name string
	Color byte
}

//stores the board layout
type BoardLayout [10][10]string
var board BoardLayout = BoardLayout{}

var IPAddress string = "localhost:8080"


func main() {
	mainMenu()
}

//lists the open games being hosted by the server
func listGames() {
	//start connection with the server
	conn, err := net.Dial("tcp", IPAddress)

	//check for an error
	checkForError(err)
	defer conn.Close()

	//send LISTGAME command to the server
	conn.Write([]byte("LISTGAME"))

	//get a response
	response := make([]byte, 120)
	_, err = conn.Read(response)
	responseString := string(response[:120])

	fmt.Printf("Game List: \n")
	//Print the game list
	if strings.HasPrefix(responseString, "GAMELIST") {
		//get number of games
		numElements, err := strconv.Atoi(responseString[9:10])
		checkForError(err)

		//loop through each element, get the name and color
		for i := 0; i<numElements; i++ {
			fmt.Printf("%v. %v %v\n", i+1, strings.TrimSpace(responseString[11+i*28:36+i*28]), responseString[37+i*28:38+i*28])
		}

	}

	//close the connection
	conn.Close()

}

//gets the user's name
func getName() string {
	inputReader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("Enter your name: ")
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)
		if len(text) < 25 || len(text) != 0 {
			name := text
			return name
		}
	}
	return "bob"
}

//gets the game number from the user
func getGameNumber() int {
	inputReader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("What game would you like to join?: ")
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

//gets the user's color
func getColor() byte {
	inputReader := bufio.NewReader(os.Stdin)
	for true {
		fmt.Print("Choose your color (W for white, B for Black): ")
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

//joins another player's game
func joinGame(gameNumber int, name string) {

	//start connection
	serverConnection, err := net.Dial("tcp", IPAddress)
	checkForError(err)
	defer serverConnection.Close()

	//sending request
	requestString := fmt.Sprintf("JOINGAME %25v %v", name, gameNumber)
	//fmt.Printf("Game number is.. %v ",requestString)
	serverConnection.Write([]byte(requestString))

	//get response
	response := make([]byte, 120)
	fmt.Printf("Joining... \n")
	_, _ = serverConnection.Read(response)
	responseString := string(response[:120])


	//Join the game if the response was positive, otherwise quit
	if strings.HasPrefix(responseString, "GAMEJOIN") {
		playerColor := responseString[9:10]
		playGame(playerColor, serverConnection)
	} else {
		fmt.Printf("ERROR: That game cannot be joined!")
		os.Exit(0)
	}

}

//hosts a game 
func hostGame(playerColor byte, name string) {
	//starting connection with the server
	serverConnection, err := net.Dial("tcp", IPAddress)
	checkForError(err)
	defer serverConnection.Close()

	//sending a hostagme request
	requestString := fmt.Sprintf("HOSTGAME %25v %v", name, string(playerColor))
	//fmt.Printf("your color is.. %v", playerColor)
	serverConnection.Write([]byte(requestString))

	//waiting for a pair to be found
	response := make([]byte, 120)
	fmt.Printf("Waiting for pair...\n")
	_, _ = serverConnection.Read(response)
	responseString := string(response[:120])
	//fmt.Printf("Got a Response!!!! %v", responseString)

	//making a sound when a response is received
	fmt.Print("\x07")
	fmt.Print("\x07")
	fmt.Print("\x07")

	//if the response is GAMEPAIR, start the game
	if strings.HasPrefix(responseString, "GAMEPAIR") {
		Name := responseString[9:34]
		fmt.Printf("Game pair found! You are playing with %v.\n", strings.TrimSpace(Name))
		playGame(string(playerColor), serverConnection)

	//if there are no host slots available, print an error
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

	var otherPlayer string

	if playerColor == "W" {
		otherPlayer = "B"
	} else {
		otherPlayer = "W"
	}

	//repeat until the game is over
	for endGame != true {
		//if it is your turn
		if (playerColor == "B" && turnCount % 2 == 1) || (playerColor == "W" && turnCount % 2 == 0) {
			var move string

			//if a move is possible for you
			if isMovePossible(playerColor) {
				var moveList []Coordinate
				moveOK := false

				//repeats until the move can flip more than 0 tiles
				for moveOK != true {
					//prints the board, prompts user for their choice, checks the move
					printBoard()
					inputReader := bufio.NewReader(os.Stdin)
					fmt.Print("Enter your move: ")
					move, _ = inputReader.ReadString('\n')
					fmt.Printf("move = %v\n", move)

					//gets list of tiles to flip for the move
					moveList = validateInput(move[0:2],playerColor)

					//determine whether its empty or not
					moveOK = !isMoveEmpty(moveList)

				}

				//sending the move to the server

				//fmt.Printf("move = %v\n", move[0:2])
				MoveDoneMessageString := fmt.Sprintf("DOMOVE %v", move[0:2])
				//fmt.Printf("sending... %v", MoveDoneMessageString)
				serverConnection.Write([]byte(MoveDoneMessageString))

				//executing the move client side
				doMove(moveList, playerColor)

				//advancing the game
				turnCount++

				//checking if the game is over
				if isGameOver() {
					endGame = true
					break
				}
			} else { //if no moves are possible, tell the server
				fmt.Printf("No moves are possible for you! \n")
				move = "XX"
				MoveDoneMessageString := fmt.Sprintf("DOMOVE %v", move[0:2])
				//fmt.Printf("sending... %v", MoveDoneMessageString)
				serverConnection.Write([]byte(MoveDoneMessageString))
				turnCount++
			}

		} else { //if it is not your turn
			var moveList []Coordinate

			//print the board, wait for oponent to move
			printBoard()
			fmt.Print("Waiting for opponent to move...\n")

			//get response from server
			playerMoveRequest := make([]byte, 120)
			_, err := serverConnection.Read(playerMoveRequest)
			checkForError(err)
			fullMoveString := string(playerMoveRequest[:120])
			moveString := fullMoveString[9:11]

			//fmt.Printf("full string: '%s'\n", playerMoveRequest)
			//fmt.Printf("movestring = '%v' ", moveString)

			//if the opnent can't move, print that and move on
			if moveString == "XX" {
				fmt.Printf("The opponent could not move.\n")
			} else if moveString == "ZZ"{ 
				fmt.Printf("Error! The oponent disconnected.\n")
				serverConnection.Close()
				mainMenu()
				
			} else { //otherwise, execute the oponent's move. 
				moveList = validateInput(moveString, otherPlayer)
				doMove(moveList, otherPlayer)
			}

			turnCount++

			//if the game is over, let the server know and set a flag.
			if isGameOver() {
				move := "YY"
				MoveDoneMessageString := fmt.Sprintf("DOMOVE %v", move[0:2])
				fmt.Printf("sending... %v", MoveDoneMessageString)
				serverConnection.Write([]byte(MoveDoneMessageString))

				endGame = true
				break
			}
		}	
	}
	fmt.Printf(getGameResults())
}


//displays the main menu
func mainMenu() {
	//Press enter to continue
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to continue... ")
	_, _ = inputReader.ReadString('\n')

	//repeats until a valid input is given
	for true {
		//prompts the user to choose an option
		fmt.Print("\nWelcome! Press 1 to list games, 2 to join a game, 3 to host a game! 4 to quit: ")
		text, _ := inputReader.ReadString('\n')
		text = strings.TrimSpace(text)

		textInt, err := strconv.Atoi(text)
		checkForError(err)

		//parses the response and calls the appropriate function
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
			board[i][j] = "0"
		}
	}
	board[4][4] = "B"
	board[4][5] = "W"
	board[5][4] = "B"
	board[5][5] = "W"
}

func printBoard() {
	fmt.Printf("\n   A B C D E F G H I J\n")
	fmt.Printf("  ┌-------------------\n")

	for i := range board {
 		fmt.Printf("%d |",i)
		for j := range board[i] {
			if board[j][i] != "0" {
				fmt.Printf("%v ", board[j][i])
			} else {
				fmt.Printf("  ")
			}
		}
		fmt.Print("\n")
	}
}

func validateInput(moveInput string, playerColor string) []Coordinate{
	var flipList []Coordinate
	var currentX int
	var currentY int

	//creating list of different move shifts for each 8 directions
	shiftList := []Coordinate {Coordinate{PosX: 0, PosY: 1},
	 Coordinate{PosX: 1, PosY:1}, Coordinate{PosX: 1, PosY:0}, Coordinate{PosX: 1, PosY:-1},
	 Coordinate{PosX: 0, PosY:-1},Coordinate{PosX: -1, PosY:-1},Coordinate{PosX: -1, PosY:0},
	 Coordinate{PosX: -1, PosY: 1}}


	//declaring other player's color
	var otherColor string
	if playerColor == "W" {
		otherColor = "B"
	} else if playerColor == "B"{
		otherColor = "W"
	} else {
		fmt.Printf("error.. playerColor = %v", playerColor)
	}

	//converting string input to ints
	startCoords := convertInput(moveInput)
	startX := startCoords.PosX
	startY := startCoords.PosY

	//fmt.Printf("startX = %v, startY = %v", startX, startY)

	//if the move is not on the board or the space is not empty, return an empty fliplist
	if !isOnBoard(startX, startY) || board[startX][startY] != "0" {
		//fmt.Printf("error 1 ")
		return flipList
	}
	
	//set the current position to the player's color temporarily
	board[startX][startY] = playerColor


	//looping once for each direction on shiftlist
	for _, currentShift := range shiftList {
		//setting current position, start position, shift amounts
		currentX = startX
		currentY = startY
		shiftX := currentShift.PosX
		shiftY := currentShift.PosY

		//do first shift
		currentX +=  shiftX
		currentY +=  shiftY

		//if it is on the board
		if (isOnBoard(currentX, currentY)) {
			//if it is the other player's color
			if (board[currentX][currentY] == otherColor) {
				//fmt.Printf("made it in sx = %v sy = %v\n", shiftX, shiftY)

				//The tile next to the start position is another player's tile so we should
				//shift again.
				currentX +=  shiftX
				currentY +=  shiftY

				//if this piece is not ont the board, continue to the next shift
				if !isOnBoard(currentX, currentY) {
					//fmt.Printf("not on board %v %v\n", currentX, currentY)
					continue
				}

				//repeat while the tile is still the other color
				for board[currentX][currentY] == otherColor {
					//shift
					currentX +=  shiftX
					currentY +=  shiftY

					//stop if its off the board
					if !isOnBoard(currentX, currentY) {
						//fmt.Printf("not on board %v %v\n", currentX, currentY)
						break
					}

				}

				//stop if its off the board
				if !isOnBoard(currentX, currentY) {
					//fmt.Printf("not on board %v %v\n", currentX, currentY)
					continue
				}

				//found a connecting tile to make a move, start going backwards
				//and adding tiles to the move list
				if (board[currentX][currentY] == playerColor) {
					//fmt.Printf("backtracking.. %v %v\n", shiftX, shiftY)
					//repeat until we reach the start
					for true {
						//go backwards one tile
						currentX -= shiftX
						currentY -= shiftY

						//if this is the last tile, end.
		 				if (currentX == startX && currentY == startY) {
							break
						}

						//add the tile to the flip list
						newCoordiante := Coordinate{PosX:currentX, PosY:currentY}
		        		flipList = append(flipList, newCoordiante)
					}
				}
			} else {
				//fmt.Printf("other color [%v %v] %v!=%v\n", currentX, currentY,  board[currentX][currentY], otherColor)
			}
		} else {
			//fmt.Printf("not on board %v %v\n", currentX, currentY)
		}
	}

	//set the start position back to null
	board[startX][startY] = "0"

	//if the move is valid, add the start position to the list of tiles to flips
	if len(flipList) > 0 {
		newCoordiante := Coordinate{PosX:startX, PosY:startY}
		flipList = append(flipList, newCoordiante)
	}


	//fmt.Printf("FLIP LIST %v", flipList)

	//return the list of tiles to flip
	return flipList
}

//converts an input string to int tuple 
func convertInput(moveInput string) Coordinate {
	//if the length is not correct, return -1s
	if len(moveInput)!=2 {

		fmt.Printf("fail1")
		var convertedMove Coordinate
		convertedMove.PosX = -1
		convertedMove.PosY = -1
		return convertedMove

	}

	//otherwise, convert the move and return it
	var convertedMove Coordinate
	convertedMove.PosX = int(moveInput[0]-65)
	convertedMove.PosY = int(moveInput[1]-48)

	return convertedMove
}

//determines whether or not it is possible for the player to make a move
func isMovePossible(playerColor string) bool {
	//loop over board
	for i := range board {
		for j := range board[i] {
			var inputString string
			inputString = fmt.Sprintf("%s%s",string(i+65),string(j+48))
			//fmt.Printf("%v\n",inputString)

			//check yoyal flips possible for the space
			moveList := validateInput(inputString, playerColor)

			//if the flip list is not empty, return true 
			//because there are flips possible
			if !isMoveEmpty(moveList) {
				return true
			}
		}
	}

	//return false if no flips were found
	return false;

}

//determines whether or not the game is over
func isGameOver() bool {
	//loop over board
	for i := range board {
		for j := range board[i] {
			//if a space is empty return false
			if board[i][j] == "0" {
				return false
			}
		}
	}

	//if no spaces are empty return true
	return true
}

//gets a summary of the results of the game in a string
func getGameResults() string {
	whiteScore := 0
	blackScore := 0

	//tallying the score by looping through all tiles
	for i := range board {
		for j := range board[i] {
			//incrementing the correct count
			if board[i][j] == "W" {
				whiteScore++
			} else if board[i][j] == "B" {
				blackScore++
			}
		}
	}

	//printing the board and outputting results
	fmt.Printf("\n")
	printBoard() 
	fmt.Printf("\n")

	if whiteScore > blackScore {
		returnString := fmt.Sprintf("White won! the final score was %v - %v\n", whiteScore, blackScore)
		return returnString
	} else {
		returnString := fmt.Sprintf("Black won! the final score was %v - %v\n", blackScore, whiteScore)
		return returnString
	}
}

//tells whether or now the move is on the board
func isOnBoard(inputX int, inputY int) bool {
	if (inputX >= 0 && inputX <=9 && inputY >= 0 && inputY <= 9) {
		return true
	} else {
		return false
	}
}

//determines whether the move list is empty or not, returns a bool
func isMoveEmpty(moveList []Coordinate) bool {
	if (len(moveList) == 0) {
		return true
	} else {
		return false
	}
}


//Flips list of moves in array of Coordinate structs passed in
//Changes board values for each position to the player color specified
func doMove(moveList []Coordinate, playerColor string) {
	for _, i:= range moveList {
		board[i.PosX][i.PosY] = playerColor
	}
}
