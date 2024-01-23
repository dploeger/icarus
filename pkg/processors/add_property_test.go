package processors

import (
	"github.com/dploeger/icarus/internal"
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPropertyProcessor_Process(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	subject := AddPropertyProcessor{}
	subject.SetToolbox(NewToolbox())
	subject.propertyName = internal.StringAddr("X-TEST")
	subject.propertyValue = internal.StringAddr("test")
	subject.overwrite = internal.BoolAddr(false)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process got an error")
	assert.Len(t, output.Children, 1, "Invalid number of events")
	assert.NotNil(t, output.Children[0].Props.Get("X-TEST"), "Test property does not exists")
	assert.Equal(t, "test", output.Children[0].Props.Get("X-TEST").Value, "Test property has wrong value")
}

func TestAddPropertyProcessor_NotOverwrite(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	event1.Props.SetText("X-TEST", "test")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	subject := AddPropertyProcessor{}
	subject.SetToolbox(NewToolbox())
	subject.propertyName = internal.StringAddr("X-TEST")
	subject.propertyValue = internal.StringAddr("test2")
	subject.overwrite = internal.BoolAddr(false)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process got an error")
	assert.Len(t, output.Children, 1, "Invalid number of events")
	assert.NotNil(t, output.Children[0].Props.Get("X-TEST"), "Test property does not exists")
	assert.Equal(t, "test", output.Children[0].Props.Get("X-TEST").Value, "Test property has wrong value")
}

func TestAddPropertyProcessor_Overwrite(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	event1.Props.SetText("X-TEST", "test")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	subject := AddPropertyProcessor{}
	subject.SetToolbox(NewToolbox())
	subject.propertyName = internal.StringAddr("X-TEST")
	subject.propertyValue = internal.StringAddr("test2")
	subject.overwrite = internal.BoolAddr(true)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process got an error")
	assert.Len(t, output.Children, 1, "Invalid number of events")
	assert.NotNil(t, output.Children[0].Props.Get("X-TEST"), "Test property does not exists")
	assert.Equal(t, "test2", output.Children[0].Props.Get("X-TEST").Value, "Test property has wrong value")
}
