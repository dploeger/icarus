package processors

import (
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestPrintProcessor_Process(t *testing.T) {
	toolbox := NewToolbox()
	toolbox.TextSelectorPattern = regexp.MustCompile("test")
	subject := PrintProcessor{}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "not")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component, event2.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process yielded error")
	assert.Len(t, output.Children, 2, "Output calendar was filtered")
	assert.Equal(t, "test", output.Children[0].Props.Get(ical.PropSummary).Value, "Output calendar had the wrong event")
	assert.Equal(t, "not", output.Children[1].Props.Get(ical.PropSummary).Value, "Output calendar had the wrong event")
}
