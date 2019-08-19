package play_test

import (
	"gitlab.com/robolucha/robolucha-api/model"
	"gitlab.com/robolucha/robolucha-api/routes/play"
	"gotest.tools/assert"
	"testing"
)

func TestPlayRequestHandler(t *testing.T) {
	handler := play.Listen()

	s1 := handler.Send(play.Request{AvailableMatch: &model.AvailableMatch{ID: 42}})
	s2 := handler.Send(play.Request{AvailableMatch: &model.AvailableMatch{ID: 43}})

	r1 := <-s1
	r2 := <-s2

	assert.Equal(t, uint(42), r1.Match.ID)
	assert.Equal(t, uint(43), r2.Match.ID)

	s3 := handler.Send(play.Request{AvailableMatch: &model.AvailableMatch{ID: 3}})
	r3 := <-s3
	assert.Equal(t, uint(3), r3.Match.ID)

}
