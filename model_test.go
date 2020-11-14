package main

import "testing"

func TestModelImplInit(t *testing.T){
	m := ModelImpl{}
	m.init(2)
	if len( m.extractableCards) != 504 {
		t.Errorf("wrong number of cards")
	}
}

func contains(s []Card, el Card) bool {
	for _, e := range s {
		if e == el {
			return true
		}
	}
	return false
}

func TestModelImplExtract(t * testing.T){
	m := ModelImpl{}
	m.init(2)
	extracted := make([]Card, 0)
	for i := 0; i < m.getNumCards(); i++ {
		c := m.extractCard()
		if contains(extracted, c){
			t.Errorf("card already extracted")
		}
		extracted = append(extracted, c)
	}

	if len(m.extractableCards) != 0 {
		t.Log("There are ", len(m.extractableCards))
		t.Errorf("number of cards is not 0")
	}
}

func containsP(s []*Player, p *Player) bool {
	for _, el := range s {
		if el == p {
			return true
		}
	}
	return false
}

func TestAddPlayer(t *testing.T){
	m := ModelImpl{}
	m.init(2)


	p, _ := m.addPlayer("pietro")

	if !containsP(m.teams[0].players, p) {
		t.Errorf("team does not contain player")
	}

	if p.team != m.teams[0] {
		t.Errorf("player does not have team")
	}

	p2, _ := m.addPlayer("marco")

	if !containsP(m.teams[1].players, p2) {
		t.Errorf("p2 not in team 1")
	}

	if p2.team != m.teams[1] {
		t.Errorf("team 1 is not p2's team")
	}

	
}
	
