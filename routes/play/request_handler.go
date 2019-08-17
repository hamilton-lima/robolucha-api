package play

import (
	"sync"
)

// Request definition
type Request struct {
	data string
}

// Response definition
type Response struct {
	data string
}

type internalRequest struct {
	data   Request
	output chan Response
}

// RequestHandler definition
type RequestHandler struct {
	input chan internalRequest
	wait  sync.WaitGroup
}

// BuildRequestHandler creates an instance of the request handler
func BuildRequestHandler() *RequestHandler {
	handler := RequestHandler{
		input: make(chan internalRequest),
	}

	// starts the listener and notify main goroutine to wait
	// using the waitgroup from the handler
	go func() {
		for {
			handler.wait.Add(1)
			go handler.process()
			handler.wait.Wait()
		}
	}()

	return &handler
}

// process handles one request from the handler.input channel
func (handler *RequestHandler) process() {
	next := <-handler.input
	next.output <- Response{data: next.data.data}
	handler.wait.Done()
}

// Send definition
func (handler *RequestHandler) Send(input Request) chan Response {
	response := make(chan Response)

	request := internalRequest{
		data:   input,
		output: response,
	}

	handler.input <- request
	return request.output
}
