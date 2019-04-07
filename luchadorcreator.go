package main

import (
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

func defaultCode() []Code {
	repeat := Code{Event: "onRepeat", Script: "move(20)\nfire(1)"}
	onHitWall := Code{Event: "onHitWall", Script: "turn(45)"}
	return []Code{repeat, onHitWall}
}

var background2 = []string{"background2_0001.png", "background2_0002.png", "background2_0003.png",
	"background2_0004.png", "background2_0005.png", "background2_0006.png", "background2_0007.png",
	"background2_0008.png", "background2_0009.png", "background2_0010.png", "background2_0011.png",
	"background2_0012.png", "background2_0013.png", "background2_0014.png", "background2_0015.png"}

var ornamentTop = []string{"ornamento_cima0001.png", "ornamento_cima0002.png",
	"ornamento_cima0003.png", "ornamento_cima0004.png", "ornamento_cima0005.png", "ornamento_cima0006.png",
	"ornamento_cima0007.png", "ornamento_cima0008.png", "ornamento_cima0009.png", "ornamento_cima0010.png",
	"ornamento_cima0011.png"}

var ornamentBottom = []string{"ornamento_baixo0001.png", "ornamento_baixo0002.png",
	"ornamento_baixo0003.png", "ornamento_baixo0004.png", "ornamento_baixo0005.png", "ornamento_baixo0006.png",
	"ornamento_baixo0007.png", "ornamento_baixo0008.png", "ornamento_baixo0009.png",
	"ornamento_baixo0010.png"}

var face = []string{"rosto0001.png", "rosto0002.png", "rosto0003.png", "rosto0004.png",
	"rosto0005.png", "rosto0006.png", "rosto0007.png", "rosto0008.png", "rosto0009.png", "rosto0010.png"}

var mouth = []string{"boca0001.png", "boca0002.png", "boca0003.png", "boca0004.png",
	"boca0005.png", "boca0006.png", "boca0007.png", "boca0008.png", "boca0009.png", "boca0010.png",
	"boca0011.png", "boca0012.png", "boca0013.png", "boca0014.png", "boca0015.png", "boca0016.png",
	"boca0017.png", "boca0018.png", "boca0019.png", "boca0020.png"}

var MASK_CONFIG_KEYS = []string{"background", "background.color",
	"background2", "background2.color",
	"ornamentTop", "ornamentTop.color",
	"ornamentBottom", "ornamentBottom.color",
	"face", "face.color",
	"mouth", "mouth.color",
	"eye", "eye.color"}

func randomConfig() []Config {

	list := []Config{}
	list = add2ConfigList(list, "background", "backgroud.png")
	list = add2ConfigList(list, "background.color", randomColor())

	list = add2ConfigList(list, "background2", randomString(background2))
	list = add2ConfigList(list, "background2.color", randomColor())

	list = add2ConfigList(list, "ornamentTop", randomString(ornamentTop))
	list = add2ConfigList(list, "ornamentTop.color", randomColor())

	list = add2ConfigList(list, "ornamentBottom", randomString(ornamentBottom))
	list = add2ConfigList(list, "ornamentBottom.color", randomColor())

	list = add2ConfigList(list, "face", randomString(face))
	list = add2ConfigList(list, "face.color", randomColor())

	list = add2ConfigList(list, "mouth", randomString(mouth))
	list = add2ConfigList(list, "mouth.color", randomColor())

	list = add2ConfigList(list, "eye", "eye.png")
	list = add2ConfigList(list, "eye.color", randomColor())

	log.WithFields(log.Fields{
		"config": list,
	}).Info("generated config")

	return list
}

func add2ConfigList(list []Config, key string, value string) []Config {
	return append(list, Config{Key: key, Value: value})
}

func randomString(list []string) string {
	return list[random(len(list))]
}

func randomColor() string {
	return NMSCOLORS[random(len(NMSCOLORS))]
}

func random(max int) int {
	source := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(source)
	return randomizer.Intn(max)
}
