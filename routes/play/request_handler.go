package play

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"gitlab.com/robolucha/robolucha-api/datasource"
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/pubsub"
	"sync"
	"time"
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
	mutex     *sync.Mutex
}

// NewRequestHandler creates a new request handler
func NewRequestHandler(_ds *datasource.DataSource, _publisher pubsub.Publisher) *RequestHandler {
	handler := RequestHandler{
		ds:        _ds,
		publisher: _publisher,
		mutex:     &sync.Mutex{},
	}

	return &handler
}

// Play handler.mutex keeps this executation one by one
func (handler *RequestHandler) Play(availableMatch *model.AvailableMatch, luchadorID uint) *model.Match {
	defer handler.mutex.Unlock()

	handler.mutex.Lock()
	matches := *handler.ds.FindActiveMatches("available_match_id = ?", availableMatch.ID)
	var match *model.Match
	if len(matches) > 0 {
		match = &matches[0]
	}

	if match == nil {
		match = handler.createMatch(availableMatch)
		handler.publishStartMatch(match)
		handler.publishJoinMatch(match, luchadorID)
	} else {
		handler.publishJoinMatch(match, luchadorID)
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

// func (handler *RequestHandler) findActiveMatch(availableMatch *model.AvailableMatch) *model.Match {

// 	var match model.Match

// 	if handler.ds.DB.
// 		Joins("left join game_definitions on matches.game_definition_id = game_definitions.id").
// 		Where("available_match_id = ?", availableMatch.ID).
// 		Where(handler.ds.ActiveMatchesSQL()).
// 		Order("time_start desc").First(&match).
// 		RecordNotFound() {
// 		return nil
// 	}

// 	log.WithFields(log.Fields{
// 		"match": match,
// 	}).Info("findActiveMatch")

// 	return &match
// }

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

func (handler *RequestHandler) publishJoinMatch(match *model.Match, luchadorID uint) {

	join := model.JoinMatch{
		MatchID:    match.ID,
		LuchadorID: luchadorID,
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
	if handler.ds.DB.First(&result, id).RecordNotFound() {
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
