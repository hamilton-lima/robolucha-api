package main

import (
	"math/rand"

	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/model"
)

var maskShape = []string{
	"segunda_cor0001.png",
	"segunda_cor0002.png",
	"segunda_cor0003.png",
	"segunda_cor0004.png",
	"segunda_cor0005.png",
	"segunda_cor0006.png",
	"segunda_cor0007.png",
	"segunda_cor0008.png",
	"segunda_cor0009.png",
	"segunda_cor0010.png",
	"segunda_cor0011.png",
	"segunda_cor0012.png",
	"segunda_cor0013.png",
	"segunda_cor0014.png",
	"segunda_cor0015.png"}

var ornamentTop = []string{
	"ornamento_cima0001.png",
	"ornamento_cima0002.png",
	"ornamento_cima0003.png",
	"ornamento_cima0004.png",
	"ornamento_cima0005.png",
	"ornamento_cima0006.png",
	"ornamento_cima0007.png",
	"ornamento_cima0008.png",
	"ornamento_cima0009.png",
	"ornamento_cima0010.png",
	"ornamento_cima0011.png"}

var ornamentBottom = []string{
	"ornamento_baixo0001.png",
	"ornamento_baixo0002.png",
	"ornamento_baixo0003.png",
	"ornamento_baixo0004.png",
	"ornamento_baixo0005.png",
	"ornamento_baixo0006.png",
	"ornamento_baixo0007.png",
	"ornamento_baixo0008.png",
	"ornamento_baixo0009.png",
	"ornamento_baixo0010.png"}

var face = []string{
	"rosto0001.png",
	"rosto0002.png",
	"rosto0003.png",
	"rosto0004.png",
	"rosto0005.png",
	"rosto0006.png",
	"rosto0007.png",
	"rosto0008.png",
	"rosto0009.png",
	"rosto0010.png"}

var mouth = []string{
	"boca0001.png",
	"boca0002.png",
	"boca0003.png",
	"boca0004.png",
	"boca0005.png",
	"boca0006.png",
	"boca0007.png",
	"boca0008.png",
	"boca0009.png",
	"boca0010.png",
	"boca0011.png",
	"boca0012.png",
	"boca0013.png",
	"boca0014.png",
	"boca0015.png",
	"boca0016.png",
	"boca0017.png",
	"boca0018.png",
	"boca0019.png",
	"boca0020.png"}

var eyes = []string{
	"olho0001.png",
	"olho0002.png"}

var maskColors = []string{
	"mask.primary.color",
	"mask.secondary.color",
	"mask.decoration.top.color",
	"mask.decoration.bottom.color",
	"eyes.color",
	"feet.color",
	"wrist.color",
	"ankle.color",
	"skin.color"}

var maskShapes = map[string][]string{
	"mask.shape":                   maskShape,
	"mask.decoration.top.shape":    ornamentTop,
	"mask.decoration.bottom.shape": ornamentBottom,
	"face.shape":                   face,
	"mouth.shape":                  mouth,
	"eyes.shape":                   eyes,
}

func randomConfig() []model.Config {

	list := []model.Config{}

	for _, color := range maskColors {
		randomizedColor := NMSColor{"Black", "#000000"} //placeholder
		randomizedColor = randomColor()
		list = add2ConfigList(list, color, randomizedColor.hex)
		list = add2ConfigList(list, color+".name", randomizedColor.name)
	}

	for shape, options := range maskShapes {
		list = add2ConfigList(list, shape, randomString(options))
	}

	log.WithFields(log.Fields{
		"config": list,
	}).Info("generated config")

	return list
}

func randomName(list []model.Config) string {
	var primaryColor = getFromConfigList(list, "mask.primary.color.name").Value
	var nounList = []string{
		"Abismo", "Comando", "Perro", "Cabeza", "Gato", "Toro", "Chupacabra", "Taco", "Soldado", "Huracán", "Rey", "Pirata",
	} //change to a JSON file?
	var adjectiveList = []string{
		"Grande", "Insano", "Fuerte", "Afortunado", "Ligero", "Muerto", "Peligroso", "Furioso", "Terrible",
	} //change to a JSON file?

	return "El " + nounList[random(len(nounList))] + " " + primaryColor + " " + adjectiveList[random(len(adjectiveList))]
}

func getFromConfigList(list []model.Config, key string) model.Config {
	for i := range list {
		if list[i].Key == key {
			// Found!
			return list[i]
		}
	}
	return model.Config{}
}

func add2ConfigList(list []model.Config, key string, value string) []model.Config {
	return append(list, model.Config{Key: key, Value: value})
}

func randomString(list []string) string {
	return list[random(len(list))]
}

func randomColor() NMSColor {
	return NMSCOLORS[random(len(NMSCOLORS))]
}

func random(max int) int {
	return rand.Intn(max)
}
