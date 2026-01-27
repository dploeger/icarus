package processors

import (
	"testing"

	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
)

func TestCopyPropertyProcessor_Process(t *testing.T) {
	toolbox := NewToolbox()
	subject := CopyPropertyProcessor{
		SourceProperty: ical.PropSummary,
		TargetProperty: ical.PropDescription,
	}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process yielded error") {
		source := output.Children[0].Props.Get(ical.PropSummary)
		target := output.Children[0].Props.Get(ical.PropDescription)
		assert.NotNil(t, target, "Target property was not created")
		if source != nil && target != nil {
			assert.Equal(t, source.Value, target.Value, "Property was not copied")
		}
	}
}

func TestCopyPropertyProcessor_ProcessNonOverwrite(t *testing.T) {
	toolbox := NewToolbox()
	subject := CopyPropertyProcessor{
		SourceProperty: ical.PropSummary,
		TargetProperty: ical.PropDescription,
		Overwrite:      false,
	}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	event1.Props.SetText(ical.PropDescription, "test2")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process yielded error") {
		source := output.Children[0].Props.Get(ical.PropSummary)
		target := output.Children[0].Props.Get(ical.PropDescription)
		assert.NotNil(t, target, "Target property was not created")
		if source != nil && target != nil {
			assert.Equal(t, "test2", target.Value, "Property was copied although overwrite was set to false")
		}
	}
}

func TestCopyPropertyProcessor_ProcessOverwrite(t *testing.T) {
	toolbox := NewToolbox()
	subject := CopyPropertyProcessor{
		SourceProperty: ical.PropSummary,
		TargetProperty: ical.PropDescription,
		Overwrite:      true,
	}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	event1.Props.SetText(ical.PropDescription, "test2")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process yielded error") {
		source := output.Children[0].Props.Get(ical.PropSummary)
		target := output.Children[0].Props.Get(ical.PropDescription)
		assert.NotNil(t, target, "Target property was not created")
		if source != nil && target != nil {
			assert.Equal(t, source.Value, target.Value, "Property was not copied")
		}
	}
}
