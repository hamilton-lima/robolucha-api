package play

import (
	"gotest.tools/assert"
	"testing"
)

var handler *RequestHandler

func TestPlayRequestHandler(t *testing.T) {
	handler = BuildRequestHandler()
	s1 := handler.Send(Request{data: "foo"})
	s2 := handler.Send(Request{data: "bar"})

	r1 := <-s1
	r2 := <-s2

	assert.Equal(t, "foo", r1.data)
	assert.Equal(t, "foo", r2.data)
}
