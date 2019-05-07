package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func AssertConfigMatch(t *testing.T, a []Config, b []Config) {
	for _, configA := range a {
		found := false
		for _, configB := range b {

			if configA.Key == configB.Key && configA.Value == configB.Value {
				found = true
				break
			}
		}
		assert.True(t, found)
		log.WithFields(log.Fields{
			"config": configA,
			"found":  found,
		}).Info("match found for config")
	}
}

func CountChangesConfigMatch(t *testing.T, a []Config, b []Config) int {
	counter := 0
	for _, configA := range a {
		found := false
		for _, configB := range b {

			if configA.Key == configB.Key && configA.Value != configB.Value {
				found = true
				break
			}
		}

		if !found {
			counter++
		}
	}

	return counter
}
