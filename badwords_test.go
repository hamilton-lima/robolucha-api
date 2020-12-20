package main

import (
	"testing"

	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/utility"
)

func init() {
	utility.SetupBadWordListFromFolder("metadata")
}
func TestWhenAss_ShouldBeTrue(t *testing.T) {
	if !utility.ContainsBadWord("ass") {
		t.Fail()
	}
}
func TestWhenBunda_ShouldBeTrue(t *testing.T) {
	if !utility.ContainsBadWord("bunda") {
		t.Fail()
	}
}
func TestWhenPartial_ShouldBeTrue(t *testing.T) {
	if !utility.ContainsBadWord("na  bunda") {
		t.Fail()
	}
}

func TestWhenCu_ShouldBeTrue(t *testing.T) {
	if !utility.ContainsBadWord("no cu") {
		t.Fail()
	}
}

func TestWhenCunhado_ShouldBeTrue(t *testing.T) {
	if !utility.ContainsBadWord("CUnhado") {
		t.Fail()
	}
}

func TestWhenContains_ShouldBeTrue(t *testing.T) {
	if !utility.ContainsBadWord("nabunda") {
		t.Fail()
	}
}

func TestWhenContainsDuplication_ShouldBeTrue(t *testing.T) {
	if !utility.ContainsBadWord("Asss$$") {
		t.Fail()
	}
}

func TestCompositeWord(t *testing.T) {
	if !utility.ContainsBadWord("maior cuzao") {
		t.Fail()
	}
}

func TestValidNamesFromPublicTest(t *testing.T) {
	log.SetLevel(log.InfoLevel)

	names := []string{"kaduzin", "caduzin", "caduzinho", "melhordetodos", "melhor de todos", "thebest"}
	for _, name := range names {
		if utility.ContainsBadWord(name) {
			log.WithFields(log.Fields{
				"failed name": name,
			}).Info("TestValidNamesFromPublicTest")
			t.Fail()
		}
	}
}
