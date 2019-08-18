package play

import (
	"gitlab.com/robolucha/robolucha-api/model"
	"sync"
)

// Request definition
type Request struct {
	Data *model.AvailableMatch
}

// Response definition
type Response struct {
	Data *Match
}

type message struct {
	input  Request
	output chan Response
}

// RequestHandler definition
type RequestHandler struct {
	messages chan message
	wait     sync.WaitGroup
}

// Listen starts to process the input channel and returns the instance
func Listen() *RequestHandler {
	handler := RequestHandler{
		messages: make(chan message),
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
		next.output <- buildResponse(next)
	}
	handler.wait.Done()
}

func buildResponse(next message) Response {
	availableMatch := next.input.Data

	match := findActiveMatch(availableMatch)
	if match == nil {
		createMatch(availableMatch)
		// send message to start match
		// send message to join
	} else {
		// send message to join
	}

	result := Response{Data: &match}
	return result
}

func findActiveMatch(availableMatch *model.AvailableMatch) *model.Match {

	var match model.Match
	ds.DB.
		Where("game_definition_id = ?", availableMatch.ID).
		Where("time_end < time_start").
		Order("time_start desc").First(&match)

	log.WithFields(log.Fields{
		"match": match,
	}).Info("findActiveMatch")

	return match
}

func createMatch(availableMatch *model.AvailableMatch) *model.Match {

	gameDefinition := ds.FindGameDefinition(availableMatch.gameDefinitionID)
	output, _ := json.Marshal(gameDefinition)
	gameDefinitionData := string(output)

	match := model.Match{
		TimeStart:          time.Now(),
		GameDefinitionID:   gameDefinitionID,
		GameDefinitionData: gameDefinitionData,
		AvailableMatchID: availableMatch.ID
	}

	ds.DB.Create(&match)

	log.WithFields(log.Fields{
		"match.id":         match.ID,
		"match": match,
	}).Info("Match created")

	return &match
}
