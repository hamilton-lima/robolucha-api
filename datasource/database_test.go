package datasource

import (
	"gitlab.com/robolucha/robolucha-api/model"
	"gotest.tools/assert"
	"testing"
)

func TestRemoveDuplicated(t *testing.T) {

	codes := make([]model.Code, 4)
	codes[0] = model.Code{ID: 1, Event: "onStart", Script: "turnGun(90)", GameDefinitionID: 1}
	codes[1] = model.Code{ID: 2, Event: "onStart", Script: "turnGun(91)", GameDefinitionID: 1}
	codes[2] = model.Code{ID: 3, Event: "onStart", Script: "turnGun(92)", GameDefinitionID: 1}
	codes[3] = model.Code{ID: 4, Event: "onStart", Script: "turnGun(93)", GameDefinitionID: 1}
	result := removeDuplicates(codes)

	assert.Equal(t, len(result), 1)
	assert.Equal(t, result[0].Script, "turnGun(93)")
	assert.Equal(t, result[0].ID, uint(4))

}

func TestRemoveDuplicatedMixed(t *testing.T) {

	codes := make([]model.Code, 4)
	codes[0] = model.Code{ID: 1, Event: "onStart", Script: "turnGun(90)", GameDefinitionID: 1}
	codes[1] = model.Code{ID: 2, Event: "onStart", Script: "turnGun(91)", GameDefinitionID: 1}
	codes[2] = model.Code{ID: 3, Event: "onStart", Script: "turnGun(92)", GameDefinitionID: 1}
	codes[3] = model.Code{ID: 4, Event: "onRepeat", Script: "turnGun(93)", GameDefinitionID: 1}
	result := removeDuplicates(codes)

	assert.Equal(t, len(result), 2)
	assert.Equal(t, result[0].Script, "turnGun(92)")
	assert.Equal(t, result[0].ID, uint(3))
	assert.Equal(t, result[1].Script, "turnGun(93)")
	assert.Equal(t, result[1].ID, uint(4))

}
