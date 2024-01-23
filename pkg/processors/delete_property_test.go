package processors

import (
	"github.com/dploeger/icarus/v2/internal"
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteFieldProcessor_Process(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	event1.Props.SetText("X-TEST", "test")
	assert.NotNil(t, event1.Props.Get("X-TEST"), "Test property was not set")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	subject := DeletePropertyProcessor{}
	subject.SetToolbox(NewToolbox())
	subject.propertyName = internal.StringAddr("X-TEST")
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process got an error")
	assert.Len(t, output.Children, 1, "Invalid number of events")
	assert.Nil(t, output.Children[0].Props.Get("X-TEST"), "Test property still exists")
}
