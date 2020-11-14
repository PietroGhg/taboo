package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"
	"sync"
)

const cardDir = "res/cards"
const blackPath = "res/monke.jpg"
	

type Team struct {
	id int
	points int
	players []*Player
}

func (t *Team) getID() int {
	return t.id
}

func (t *Team) addPlayer(p *Player){
	p.team = t
	t.players = append(t.players, p)
}
	

type Card struct {
	id int
}

func (c Card) getPath() string {
	return cardDir + "/" + fmt.Sprint(c.id) + ".png"
}

type Player struct {
	id int
	name string
	team *Team
	connected bool
}

type ModelImpl struct {
	sync.Mutex
	nextID int
	currentPlayer int
	currentCard Card
	teams []*Team
	Players map[int]*Player
	cards []Card
	extractableCards []int
}

func (this *ModelImpl) getNumCards() int {
	return len(this.cards)
}

func (this *ModelImpl) getBlackPath() string {
	return blackPath
}

func (this *ModelImpl) init(numTeams int){
	rand.Seed(time.Now().UTC().UnixNano())
	//initialize teams
	this.teams = make([]*Team, numTeams)
	for i := range this.teams {
		this.teams[i] = &Team{id: i, points: 0}
		this.teams[i].players = make([]*Player,0)
	}
	//load number of cards from cards dir
	cards, _ := ioutil.ReadDir(cardDir)
	//initialize extractableCards
	this.extractableCards = make([]int, len(cards))
	for i := range this.extractableCards {
		this.extractableCards[i] = i
	}	
	//initialize cards
	this.cards = make([]Card, len(cards))
	for i := range this.cards {
		this.cards[i] = Card{i}
	}
	//initilize players map
	this.Players = make(map[int]*Player)
}

func (this *ModelImpl) getTeamWithLessPlayers() *Team {
	min := len(this.teams[0].players)
	team := this.teams[0]
	for _, t := range this.teams {
		if len(t.players) < min {
			team = t
		}
	}
	return team
}

func (this *ModelImpl) setDisconnected(id int){
	this.Players[id].connected = false;
}

func (this *ModelImpl) setConnected(id int){
	this.Players[id].connected = true;
}

func (this *ModelImpl) existsName(name string) bool {
	fmt.Println("Checking ", name)
	for _, p := range this.Players {
		if p.name == name {
			return true
		}
	}
	return false
}

func (this *ModelImpl) getPlayerID(name string) (int, error) {
	for id, p := range this.Players {
		if p.name == name {
			return id, nil
		}
	}
	return -1, nil
}

//returns true if name is available
func (this *ModelImpl) checkName(name string) bool {
	for _, p := range this.Players {
		if p.name == name {
			return false
		}
	}
	return true
}

func (this *ModelImpl) addPlayer(name string) (*Player, error) {
	if !(this.checkName(name)){
		return nil, errors.New("Name already taken")
	}
	id := this.nextID
	this.nextID++
	team := this.getTeamWithLessPlayers()
	p := &Player{name: name,
		team: team,
		id: len(this.Players),
		connected: true,
	}
	team.addPlayer(p)
	if _, ok := this.Players[id]; ok {
		return nil, errors.New("Player already in game")
	}
	this.Players[id] = p
	return p, nil
}

func remove(s []int, i int) []int {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func (this *ModelImpl) extractCard() Card {
	//get random number
	n := rand.Intn(len(this.extractableCards))
	index := this.extractableCards[n]
	//remove element n from extractableCards
	this.extractableCards = remove(this.extractableCards, n)
	this.currentCard = this.cards[index]
	return this.currentCard
}

func (this *ModelImpl) nextTurn() {
	this.currentPlayer = (this.currentPlayer + 1) % len(this.Players)
		
}

func (this* ModelImpl) currentTeam() int {
	return this.Players[this.currentPlayer].team.id
}

func (this* ModelImpl) getTeam(playerID int) *Team {
	return this.Players[playerID].team
}

func (this *ModelImpl) isWatching(playerID int) bool{
	//a player is watching if he is playing or the current player is not from his team
	notFromTeam := this.getTeam(playerID).id != this.currentTeam()
	return this.isPlaying(playerID) || notFromTeam
}

func (this *ModelImpl) isPlaying(playerID int) bool{
	return playerID == this.currentPlayer
}






