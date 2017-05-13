# CSS 432 Final Project: Othello
Zeke Snider  
5/31/15

## Build Instructions
The game was programmed using the [Go programming language](https://golang.org). To build it you must first [install](https://golang.org/doc/install) the Go runtime. It is installed by default on the linux lab machines, however it is installed in a non default location. You will have to change the makefiles to run the game on those machines.

To run the server, cd into the main directory and type "make server", then "make rserver". The server will start on port 8080. To connect, type "make client", then "make rclient". This will launch a client which will attempt to connect to a server on localhost:8080. To change where it connects, open client.go and edit line 31 to change the "IPAddress" variable.

All source code is located within the /src/ directory.
    
## Socket Design
The game uses a client server model. Each client communicates to eachother by use of the server as a middle man.

Requests:

* JOINGAME [Name 25 Chars] | [Game No 2 Chars] 
* HOSTGAME [Name 25 Chars] | [Color: 1 char]
* LISTGAME [Page 1 char]
* NUMGAMES 
* DOMOVE [2 chars]
	(XX = cannot move)
	(YY = END GAME)


Responses:
* MOVEDONE [move 2 chars]
	(ZZ = oponent disconnected.)
* GAMEPAIR [player name 25 chars]
* GAMEJOIN [color 2 chars]
* GAMELIST [Name 25 chars] | [Color 1 char] | [Name 25 chars] | [Color 1 char] | [Name 25 chars] | [Color 1 char] | [Name 25 chars] | [Color 1 char]


The initial requests to the serve are:

### LISTGAME
Lists the games, and returns the list to the client. Then the connection is terminated. The server's response is formatted in the GAMELIST format detailed above. The name is padded to be 25 chars, and unpadded when received on the client side.

### HOSTGAME
Hosts a game, then waits for a response when another player joins the room. Once another player joins, GAMEJOIN response is sent. The client and server then enter the gameloop.

### JOINGAME
Joins another player's room. The server responds with GAMEJOIN reponse. The client and server then enter the gameloop.

### Game Loop
The board is displayed to both players, and the player color Black is prompted to move first. Once he/she moves, DOMOVE is sent to the server, and the server forwards on the request as MOVEDONE to the other player. Then the other player is prompted for their move and so on. The validity of the move is checked on the client's side before making the move and sending it. If a player has no moves possible, it sends "XX" as the move which is a code for skip. The next player then moves with no changes to the board. If "YY" is sent, the server disconnects from both clients as the game is over. If a client disconnects midmatch, an error messages is displayed to the other player and the server terminates the connection. After a game is over, the score is displayed to both players with a winner by the client side, then the user is redirected to the main menu. 

