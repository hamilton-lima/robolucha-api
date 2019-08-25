package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"gitlab.com/robolucha/robolucha-api/model"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func TestGenerateSampleGameDefinition(t *testing.T) {

	gd := model.BuildDefaultGameDefinition()

	gd.Name = "ALL-AGAINST-ALL"
	gd.Label = "All against all"
	gd.Type = "multiplayer"
	gd.SortOrder = 0

	gd.GameComponents = make([]model.GameComponent, 2)
	gd.SceneComponents = make([]model.SceneComponent, 0)
	gd.Codes = make([]model.Code, 0)
	gd.LuchadorSuggestedCodes = make([]model.Code, 0)

	gd.GameComponents[0].Name = "otto"
	gd.GameComponents[0].Configs = make([]model.Config, 0)
	gd.GameComponents[0].Codes = make([]model.Code, 3)
	gd.GameComponents[0].Codes[0] = model.Code{Event: "onStart", Script: "turnGun(90)"}
	gd.GameComponents[0].Codes[1] = model.Code{Event: "onRepeat", Script: "move(20)\nfire(3)"}
	gd.GameComponents[0].Codes[2] = model.Code{Event: "onHitWall", Script: "turn(90)\nturnGun(90)"}

	gd.GameComponents[1].Name = "farol"
	gd.GameComponents[1].Configs = make([]model.Config, 0)
	gd.GameComponents[1].Codes = make([]model.Code, 1)
	gd.GameComponents[1].Codes[0] = model.Code{Event: "onRepeat", Script: "turn(10)\nturnGun(-10)\nfire(1)"}

	output, _ := json.Marshal(gd)
	result := string(output)

	err := ioutil.WriteFile("example-gamedefinition.json", []byte(result), 0644)
	check(err)
}
