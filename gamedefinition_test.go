package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func TestGenerateSampleGameDefinition(t *testing.T) {

	gd := BuildDefaultGameDefinition()

	gd.Name = "ALL-AGAINST-ALL"
	gd.Label = "All against all"
	gd.Type = "multiplayer"
	gd.SortOrder = 0

	gd.Participants = make([]Luchador, 2)
	gd.SceneComponents = make([]SceneComponent, 0)
	gd.Codes = make([]ServerCode, 0)
	gd.LuchadorSuggestedCodes = make([]ServerCode, 0)

	gd.Participants[0].Name = "otto"
	gd.Participants[0].Configs = make([]Config, 0)
	gd.Participants[0].Codes = make([]Code, 3)
	gd.Participants[0].Codes[0] = Code{Event: "onStart", Script: "turnGun(90)"}
	gd.Participants[0].Codes[1] = Code{Event: "onRepeat", Script: "move(20)\nfire(3)"}
	gd.Participants[0].Codes[2] = Code{Event: "onHitWall", Script: "turn(90)\nturnGun(90)"}

	gd.Participants[1].Name = "farol"
	gd.Participants[1].Configs = make([]Config, 0)
	gd.Participants[1].Codes = make([]Code, 1)
	gd.Participants[1].Codes[0] = Code{Event: "onRepeat", Script: "turn(10)\nturnGun(-10)\nfire(1)"}

	output, _ := json.Marshal(gd)
	result := string(output)

	err := ioutil.WriteFile("example-gamedefinition.json", []byte(result), 0644)
	check(err)
}
