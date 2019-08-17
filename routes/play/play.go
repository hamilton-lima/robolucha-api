package play

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
}

// BuildRequestHandler creates an instance of the request handler
func BuildRequestHandler() *RequestHandler {
	handler := RequestHandler{
		input: make(chan internalRequest),
	}

	go handler.listen()
	return &handler
}

func (handler *RequestHandler) listen() {
	for {
		// process one by one
		next := <-handler.input
		// adds the response to the output channel
		next.output <- Response{data: next.data.data}
	}
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
