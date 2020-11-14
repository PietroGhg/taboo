package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

type Connection struct {
	playerName string
	isUp bool
	conn *websocket.Conn
}

func (this *Connection) setConn(conn *websocket.Conn) {
	this.conn = conn
}

var upgrader = websocket.Upgrader{
 ReadBufferSize:  1024,
 WriteBufferSize: 1024,
 CheckOrigin:     func(r *http.Request) bool { return true },
}

var m ModelImpl
var oneConnection bool
var connections map[int]*Connection

func homePage(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w,r, "index.html")
}

func sendCard(path string, playerID int, msgType int){
	if !m.isWatching(playerID){
		path = m.getBlackPath()
	}
	
	conn := connections[playerID].conn
	img, _ := ioutil.ReadFile(path)
	str := base64.StdEncoding.EncodeToString(img)
	cardMessage := cardMsg{str}
	i, _ := genericJSON(cardMessage)
	if err := conn.WriteMessage(msgType, i); err != nil {
		log.Println(err)
	}

}

func isUp(name string) bool {
	id, _ := m.getPlayerID(name)
	return connections[id].isUp
}

func checkAndSendCard(player *Player, msgType int) {
	//if it's the first connection, extract a card and send it
	//else, send the current card
	var path string
	if(!oneConnection){
		path = m.extractCard().getPath()
		oneConnection = true
	} else {
		path = m.currentCard.getPath()
	}	
	sendCard(path, player.id, msgType)
	
}

func handleHello(msg hello, conn *websocket.Conn, msgType int){
	m.Lock()
	//check if name already taken
	if !(m.checkName(msg.Name)) && isUp(msg.Name) {
		//TODO: handle reconnection
		log.Println("Attempt to log as ", msg.Name)
		ack := helloACK{
			Accepted :false,
			Text: "Username già in uso",
		}
		b, _ := genericJSON(ack)
		if err := conn.WriteMessage(msgType, b); err != nil {
			log.Println(err)
		}
		return
	} else if m.existsName(msg.Name) {
		if  !isUp(msg.Name){
			//reconnection
			log.Println("Reconnecting ", msg.Name)
			id, err := m.getPlayerID(msg.Name)
			if err != nil {
				log.Fatal("Error while retrieving player id")
			}
			
			//store new connection
			connections[id].conn = conn
			//prepare ack
			ack := helloACK{
			Accepted: true,
				Text : "Bentornato, " + msg.Name,
				PlayerID : id,
				TeamID : m.Players[id].team.id,
			}
			//send ack
			b, _ := genericJSON(ack)
			if err := conn.WriteMessage(msgType, b); err != nil {
				log.Println(err)
			}
			//send card
			checkAndSendCard(m.Players[id], msgType)
			return
		}
	}
	//name not taken: accept player
	//adds player to model
	name := string(msg.Name)
	player, err := m.addPlayer(name)
	if err != nil {
		fmt.Println("Cannot add player")
		//TODO: handle this case
	}
	//store connection
	connections[player.id] = &Connection{
		playerName: name,
		isUp: true,
		conn: conn,
	}

	//send ack
	ack := helloACK{
		Accepted : true,
		Text: "Benvenuto, " + msg.Name,
		PlayerID: player.id,
		TeamID: player.team.id,
	}
	b, _ := genericJSON(ack)
	if err := conn.WriteMessage(msgType, b); err != nil {
		log.Println(err)
	}

	//if it's the first connection, extract a card and send it
	//else, send the current card
	checkAndSendCard(player, msgType)
	m.Unlock()
}

func handleChangeCard(msg changeCard, conn *websocket.Conn, msgType int){
	m.Lock()
	if msg.PlayerID != m.currentPlayer {
		log.Println("Attempt to change turn from ", msg.PlayerID, " while it's ", m.currentPlayer, " turn")
		return 
	}
	c := m.extractCard()
	for playerID, conn := range connections{
		if(conn.isUp) {
			sendCard(c.getPath(), playerID, msgType)
		}
	}
	m.Unlock()
}

func handleStartTimer(msg timerStart, msgType int) {
	m.Lock()
	for _, c := range connections {
		if c.isUp {
			msg := timerStart{msg.PlayerID}
			js, _ := genericJSON(msg)
			if err := c.conn.WriteMessage(msgType, js); err != nil {
				log.Println(err)
			}
		}
	}
	m.Unlock()
}

func handleEndTurn(msg endTurn, msgType int) {
	m.Lock()
	//new turn
	m.nextTurn()
	p := m.currentPlayer
	
	//send a new turn message to all the connections
	nt := newTurn{p}
	toSend, _ := genericJSON(nt)
	for _, c := range connections {
		if err := c.conn.WriteMessage(msgType, toSend); err != nil {
			log.Println(err)
		}
	}

	//extract a new card and send it
	c := m.extractCard()
	for playerID := range connections {
		sendCard(c.getPath(), playerID, msgType)
	}
	m.Unlock()
}

func reader(conn *websocket.Conn){
	for {
		msgType, p, err := conn.ReadMessage()
		if err != nil {
			//handle disconnection
			for id, c := range connections {
				if c.conn == conn {
					c.isUp = false
					c.conn = nil
					m.Players[id].connected = false
					fmt.Println(m.Players[id].name, " disconnected");
					break
				}
			}
			
			return 
		}

		log.Println(string(p))
		msg, err := getMsg(p)
		if err != nil {
			fmt.Println(err)
		}
		switch msg.(type) {
		case hello:
			welcomeMsg := msg.(hello)
			handleHello(welcomeMsg, conn, msgType)
		case changeCard:
			changeCardMsg := msg.(changeCard)
			handleChangeCard(changeCardMsg, conn, msgType)
		case timerStart:
			startTimer := msg.(timerStart)
			handleStartTimer(startTimer, msgType)
		case endTurn:
			endT := msg.(endTurn)
			handleEndTurn(endT, msgType)
		default:
			log.Fatal("Unrecognized message", msg)
		}
		
	}
}

func wsEndPoint(w http.ResponseWriter, r *http.Request){
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Connection refused")
	} else {
		fmt.Println("Client,", r.RemoteAddr, " succesfully connected")
	}

	reader(ws)
}
	
func setupRoutes(){
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndPoint)
}

func main(){
	connections = make(map[int]*Connection)
	m.init(2)
	setupRoutes()
	fmt.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
