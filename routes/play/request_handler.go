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

// Request definition
type Request struct {
	AvailableMatch *model.AvailableMatch
	LuchadorID     uint
}

// Response definition
type Response struct {
	Match *model.Match
}

type message struct {
	input  Request
	output chan Response
}

// RequestHandler definition
type RequestHandler struct {
	messages  chan message
	wait      sync.WaitGroup
	ds        *datasource.DataSource
	publisher pubsub.Publisher
}

// Listen starts to process the input channel and returns the instance
func Listen(_ds *datasource.DataSource, _publisher pubsub.Publisher) *RequestHandler {
	handler := RequestHandler{
		messages:  make(chan message),
		ds:        _ds,
		publisher: _publisher,
	}

	// notify main goroutine to wait using the waitgroup from the handler
	handler.wait.Add(1)
	go func() {
		for {
			handler.wait.Add(1)
			go handler.process()
		}
	}()

	return &handler
}

// Send definition
func (handler *RequestHandler) Send(request Request) chan Response {
	response := make(chan Response)

	handler.messages <- message{
		input:  request,
		output: response,
	}

	return response
}

// process handles one request from the handler.input channel
func (handler *RequestHandler) process() {
	select {
	case next := <-handler.messages:
		next.output <- handler.buildResponse(next)
	}
	handler.wait.Done()
}

func (handler *RequestHandler) buildResponse(next message) Response {

	match := handler.findActiveMatch(next.input.AvailableMatch)
	if match == nil {
		match = handler.createMatch(next.input.AvailableMatch)
		handler.publishStartMatch(match)
		handler.publishJoinMatch(match, next.input.LuchadorID)
	} else {
		handler.publishJoinMatch(match, next.input.LuchadorID)
	}

	result := Response{Match: match}
	return result
}

func (handler *RequestHandler) findActiveMatch(availableMatch *model.AvailableMatch) *model.Match {

	var match model.Match

	if handler.ds.DB.
		Where("available_match_id = ?", availableMatch.ID).
		Where("time_end < time_start").
		Order("time_start desc").First(&match).
		RecordNotFound() {
		return nil
	}

	log.WithFields(log.Fields{
		"match": match,
	}).Info("findActiveMatch")

	return &match
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
