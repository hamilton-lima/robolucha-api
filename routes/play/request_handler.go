package play

import (
	"sync"
)

// Request definition
type Request struct {
	Data string
}

// Response definition
type Response struct {
	Data string
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
			go handler.process()
		}
	}()

	return &handler
}

// process handles one request from the handler.input channel
func (handler *RequestHandler) process() {
	select {
	case next := <-handler.messages:
		next.output <- Response{Data: next.input.Data}
	}
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
