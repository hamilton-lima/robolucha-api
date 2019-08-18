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

	match := Match{MatchID: availableMatch.ID}
	result := Response{Data: &match}
	return result
}
