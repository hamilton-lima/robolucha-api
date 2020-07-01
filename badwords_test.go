package main

import (
	"gitlab.com/robolucha/robolucha-api/utility"
	"testing"
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

func TestWhenCunhado_ShouldBeFalse(t *testing.T) {
	if utility.ContainsBadWord("CUnhado") {
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
