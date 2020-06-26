package model

import (
	log "github.com/sirupsen/logrus"
)

// LogGameComponent build a simplified version of gameComponent for logging
func LogGameComponent(input *GameComponent) map[string]interface{} {
	if input == nil {
		return log.Fields{
			"ID": nil,
		}
	}

	return log.Fields{
		"ID":     input.ID,
		"name":   input.Name,
		"userID": input.UserID,
	}
}

// LogGameDefinition build a simplified version of gameDefinition for logging
func LogGameDefinition(input *GameDefinition) map[string]interface{} {
	if input == nil {
		return log.Fields{
			"ID": nil,
		}
	}

	return log.Fields{
		"ID":   input.ID,
		"name": input.Name,
		"type": input.Type,
	}
}

// LogAvailableMatches build a simplified version of []AvailableMatch for logging
func LogAvailableMatches(matches *[]AvailableMatch) map[string]interface{} {

	result := log.Fields{}

	for _, match := range *matches {
		result[string(match.ID)] = LogAvailableMatch(&match)
	}

	return result
}

// LogAvailableMatch build a simplified version of AvailableMatch for logging
func LogAvailableMatch(input *AvailableMatch) map[string]interface{} {
	if input == nil {
		return log.Fields{
			"ID": nil,
		}
	}

	return log.Fields{
		"ID":             input.ID,
		"name":           input.Name,
		"ClassroomID":    input.ClassroomID,
		"GameDefinition": LogGameDefinition(input.GameDefinition),
	}
}

// LogMatch build a simplified version of Match for logging
func LogMatch(input *Match) map[string]interface{} {
	if input == nil {
		return log.Fields{
			"ID": nil,
		}
	}

	return log.Fields{
		"ID":               input.ID,
		"AvailableMatchID": input.AvailableMatchID,
		"TimeStart":        input.TimeStart,
		"TimeEnd":          input.TimeEnd,
		"GameDefinition":   LogGameDefinition(&input.GameDefinition),
	}
}

// LogMatches build a simplified version of []Match for logging
func LogMatches(matches *[]Match) map[string]interface{} {

	result := log.Fields{}

	for _, match := range *matches {
		result[string(match.ID)] = LogMatch(&match)
	}

	return result
}
