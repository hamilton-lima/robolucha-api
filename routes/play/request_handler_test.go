package play_test

import (
	"gitlab.com/robolucha/robolucha-api/routes/play"
	"gotest.tools/assert"
	"testing"
)

func TestPlayRequestHandler(t *testing.T) {
	handler := play.Listen()

	s1 := handler.Send(play.Request{Data: "foo"})
	s2 := handler.Send(play.Request{Data: "bar"})

	r1 := <-s1
	r2 := <-s2

	assert.Equal(t, "foo", r1.Data)
	assert.Equal(t, "bar", r2.Data)

	s3 := handler.Send(play.Request{Data: "other"})
	r3 := <-s3
	assert.Equal(t, "other", r3.Data)

}
