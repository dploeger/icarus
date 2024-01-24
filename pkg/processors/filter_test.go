package processors

import (
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestFilterProcessor_Process(t *testing.T) {
	toolbox := NewToolbox()
	toolbox.TextSelectorPattern = regexp.MustCompile("test")
	subject := FilterProcessor{Inverse: false}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "not")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component, event2.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process yielded error") {
		assert.Len(t, output.Children, 1, "Output calendar was filtered")
		assert.Equal(t, "test", output.Children[0].Props.Get(ical.PropSummary).Value, "Output calendar had the wrong event")
	}
}

func TestFilterProcessor_ProcessInverse(t *testing.T) {
	toolbox := NewToolbox()
	toolbox.TextSelectorPattern = regexp.MustCompile("test")
	subject := FilterProcessor{Inverse: true}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "not")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component, event2.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process yielded error") {
		assert.Len(t, output.Children, 1, "Output calendar was filtered")
		assert.Equal(t, "not", output.Children[0].Props.Get(ical.PropSummary).Value, "Output calendar had the wrong event")
	}
}
