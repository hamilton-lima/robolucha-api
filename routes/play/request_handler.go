package play

import (
	"encoding/json"
	"time"

	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"
)

// PlayRequest definition
type PlayRequest struct {
	AvailableMatch *model.AvailableMatch
	LuchadorID     uint
}

// RequestHandler definition
type RequestHandler struct {
	ds        *datasource.DataSource
	publisher pubsub.Publisher
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(_ds *datasource.DataSource, _publisher pubsub.Publisher) *RequestHandler {
	handler := RequestHandler{
		ds:        _ds,
		publisher: _publisher,
	}

	return &handler
}

func (handler *RequestHandler) findMatch(availableMatch *model.AvailableMatch) *model.Match {
	matches := *handler.ds.FindActiveMatches("available_match_id = ?", availableMatch.ID)

	log.WithFields(log.Fields{
		"active matches": model.LogMatches(&matches),
	}).Info("Play")

	if len(matches) > 0 {
		return &matches[0]
	}
	return nil
}

func isParticipating(match *model.Match, luchadorID uint) bool {
	for _, participant := range match.Participants {
		if participant.ID == luchadorID {
			return true
		}
	}
	return false
}

// Play definition
func (handler *RequestHandler) Play(
	availableMatch *model.AvailableMatch,
	luchadorID uint,
	teamID uint) *model.Match {

	log.WithFields(log.Fields{
		"availableMatch": availableMatch,
		"luchadorID":     luchadorID,
		"teamID":         teamID,
	}).Info("Play")

	match := handler.findMatch(availableMatch)

	// Match dont exist TRY to create
	if match == nil {
		log.WithFields(log.Fields{
			"status": "match not found",
		}).Info("Play")

		match = handler.createMatch(availableMatch)
		log.WithFields(log.Fields{
			"status":  "match created",
			"matchID": match.ID,
		}).Info("Play")

		handler.publishStartMatch(match)
		handler.publishJoinMatch(match, luchadorID, teamID)
	} else {
		log.WithFields(log.Fields{
			"status":  "match found",
			"matchID": match.ID,
		}).Info("Play")

		// Match is a tutorial, reset if active
		if match.GameDefinition.Type == model.GAMEDEFINITION_TYPE_TUTORIAL {
			log.WithFields(log.Fields{
				"status":  "Match is an tutorial, end and create again",
				"matchID": match.ID,
			}).Info("Play")

			// if is participating on a tutorial match restart it
			if isParticipating(match, luchadorID) {
				handler.ds.EndMatch(match)
			}

			match = handler.createMatch(availableMatch)
			log.WithFields(log.Fields{
				"status":  "Tutorial recreated",
				"matchID": match.ID,
			}).Info("Play")

			handler.publishStartMatch(match)
		}

		handler.publishJoinMatch(match, luchadorID, teamID)
	}

	return match
}

// FindTutorialMatchesByParticipant definition
func (handler *RequestHandler) FindTutorialMatchesByParticipant(gameComponent *model.GameComponent) []model.Match {

	matches := handler.ds.FindActiveMatches("game_definitions.type = ?", model.GAMEDEFINITION_TYPE_TUTORIAL)
	log.WithFields(log.Fields{
		"matches": model.LogMatches(matches),
	}).Info("FindTutorialMatchesByParticipant")

	result := make([]model.Match, 0)

	for _, match := range *matches {

		log.WithFields(log.Fields{
			"match":        model.LogMatch(&match),
			"participants": match.Participants,
		}).Info("FindTutorialMatchesByParticipant/filter participants")

		for _, participant := range match.Participants {
			if participant.ID == gameComponent.ID {
				result = append(result, match)
			}
		}
	}

	return result
}

func (handler *RequestHandler) createMatch(availableMatch *model.AvailableMatch) *model.Match {

	gameDefinition := handler.ds.FindGameDefinition(availableMatch.GameDefinitionID)
	output, _ := json.Marshal(gameDefinition)
	gameDefinitionData := string(output)

	match := model.Match{
		TimeStart:          time.Now(),
		GameDefinitionID:   gameDefinition.ID,
		GameDefinitionData: gameDefinitionData,
		AvailableMatchID:   availableMatch.ID,
	}

	handler.ds.DB.Create(&match)

	log.WithFields(log.Fields{
		"match.id": match.ID,
		"match":    match,
	}).Info("Match created")

	return &match
}

func (handler *RequestHandler) publishStartMatch(match *model.Match) {
	// publish event to run the match
	resultJSON, _ := json.Marshal(match)
	result := string(resultJSON)
	handler.publisher.Publish("start.match", result)

	log.WithFields(log.Fields{
		"start.match": result,
	}).Info("publishStartMatch")

}

func (handler *RequestHandler) publishJoinMatch(match *model.Match, luchadorID uint, teamID uint) {

	join := model.JoinMatch{
		MatchID:    match.ID,
		LuchadorID: luchadorID,
		TeamID:     teamID,
	}

	// publish event to run the match
	resultJSON, _ := json.Marshal(join)
	result := string(resultJSON)
	handler.publisher.Publish("join.match", result)

	log.WithFields(log.Fields{
		"join.match": result,
	}).Info("publishJoinMatch")

}

// FindAvailableMatchByID definition
func (handler *RequestHandler) FindAvailableMatchByID(id uint) *model.AvailableMatch {
	var result model.AvailableMatch
	if handler.ds.DB.Preload("GameDefinition").First(&result, id).RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"result": result,
	}).Info("FindAvailableMatchByID")

	return &result
}

// LeaveTutorialMatches definition
func (handler *RequestHandler) LeaveTutorialMatches(gameComponent *model.GameComponent) {

	matches := handler.FindTutorialMatchesByParticipant(gameComponent)
	log.WithFields(log.Fields{
		"matches": model.LogMatches(&matches),
	}).Info("tutorial matches")

	channel := "end.match"

	for _, match := range matches {
		handler.ds.EndMatch(&match)

		matchJSON, _ := json.Marshal(match)
		message := string(matchJSON)
		handler.publisher.Publish(channel, message)
	}

}

func (handler *RequestHandler) UserHasLevelToPlay(user *model.UserLevel, gameDefinition *model.GameDefinition) bool {
	level := user.Level
	min := gameDefinition.MinLevel
	max := gameDefinition.MaxLevel
	canPlay := (level >= min) && (max == 0 || level <= max)
	return canPlay
}
