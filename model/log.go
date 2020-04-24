package model

import (
	log "github.com/sirupsen/logrus"
)

// LogGameComponent build a simplified version of gameComponent for logging
func LogGameComponent(input GameComponent) map[string]interface{} {
	return log.Fields{
		"ID":     input.ID,
		"name":   input.Name,
		"userID": input.UserID,
	}
}

// LogGameDefinition build a simplified version of gameDefinition for logging
func LogGameDefinition(input GameDefinition) map[string]interface{} {
	return log.Fields{
		"ID":   input.ID,
		"name": input.Name,
		"type": input.Type,
	}
}

// LogMatch build a simplified version of Match for logging
func LogMatch(input Match) map[string]interface{} {
	return log.Fields{
		"ID":               input.ID,
		"AvailableMatchID": input.AvailableMatchID,
		"GameDefinition":   LogGameDefinition(input.GameDefinition),
	}
}

// LogMatches build a simplified version of []Match for logging
func LogMatches(matches *[]Match) map[string]interface{} {

	result := log.Fields{}

	for _, match := range *matches {
		result[string(match.ID)] = LogMatch(match)
	}

	return result
}
