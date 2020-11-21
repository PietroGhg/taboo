package main

import (
	"encoding/json"
	"reflect"
	"errors"
)

type Generic struct {
	T string
	Obj json.RawMessage
}

func genericJSON(i interface{}) ([]byte, error) {
	b, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	raw := json.RawMessage(b)
	t := reflect.TypeOf(i).String()
	g := Generic{t, raw}
	return json.Marshal(g)
}

func getMsg(b []byte) (interface{}, error) {
	var g Generic
	json.Unmarshal(b, &g)

	switch g.T {
	case "main.hello":
		var r hello
		json.Unmarshal(g.Obj, &r)
		return r, nil
	case "main.changeCard":
		var r changeCard
		json.Unmarshal(g.Obj, &r)
		return r, nil
	case "main.timerStart":
		var r timerStart
		json.Unmarshal(g.Obj, &r)
		return r, nil
	case "main.endTurn":
		var r endTurn
		json.Unmarshal(g.Obj, &r)
		return r, nil
	default:
		var r interface{}
		return r, errors.New("Cannot recognize type")
	}
}

type hello struct {
	Name string
}

type helloACK struct {
	Accepted bool
	Text string
	PlayerID int
	TeamID int
	CurrentPlayer int
}

type cardMsg struct {
	Img string
}

type changeCard struct {
	PlayerID int
}

type timerStart struct {
	PlayerID int
}

type endTurn struct {
	PlayerID int
}

type newTurn struct {
	CurrentPlayer int
}

type playerEntry struct {
	Name string
	Team int
	IsTurn bool
}

type playersList struct {
	Players []playerEntry
}

func (this *playersList) addPlayer(pe playerEntry) {
	this.Players = append(this.Players, pe)
}

		

	
