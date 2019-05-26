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

	gd.GameComponents = make([]GameComponent, 2)
	gd.SceneComponents = make([]SceneComponent, 0)
	gd.Codes = make([]ServerCode, 0)
	gd.LuchadorSuggestedCodes = make([]ServerCode, 0)

	gd.GameComponents[0].Name = "otto"
	gd.GameComponents[0].Configs = make([]ServerConfig, 0)
	gd.GameComponents[0].Codes = make([]ServerCode, 3)
	gd.GameComponents[0].Codes[0] = ServerCode{Event: "onStart", Script: "turnGun(90)"}
	gd.GameComponents[0].Codes[1] = ServerCode{Event: "onRepeat", Script: "move(20)\nfire(3)"}
	gd.GameComponents[0].Codes[2] = ServerCode{Event: "onHitWall", Script: "turn(90)\nturnGun(90)"}

	gd.GameComponents[1].Name = "farol"
	gd.GameComponents[1].Configs = make([]ServerConfig, 0)
	gd.GameComponents[1].Codes = make([]ServerCode, 1)
	gd.GameComponents[1].Codes[0] = ServerCode{Event: "onRepeat", Script: "turn(10)\nturnGun(-10)\nfire(1)"}

	output, _ := json.Marshal(gd)
	result := string(output)

	err := ioutil.WriteFile("example-gamedefinition.json", []byte(result), 0644)
	check(err)
}
