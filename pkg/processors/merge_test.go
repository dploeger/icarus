package processors

import (
	"regexp"
	"testing"
	"time"

	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
)

func TestMergeProcessor_Process(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC).Truncate(time.Second)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "test2")
	now = time.Now().In(time.UTC).Truncate(time.Second)
	event2.Props.SetDate(ical.PropDateTimeStart, now)
	event2.Props.SetDate(ical.PropDateTimeEnd, now)
	merge := ical.NewCalendar()
	merge.Children = append(merge.Children, event2.Component)
	subject := NewMergeProcessor(*merge)
	subject.SetToolbox(NewToolbox())
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process got an error") {
		assert.Len(t, output.Children, 2, "Invalid number of events")
		assert.Equal(t, output.Children[0].Props.Get(ical.PropSummary).Value, "test", "Invalid first event")
		assert.Equal(t, output.Children[1].Props.Get(ical.PropSummary).Value, "test2", "Invalid second event")
	}
}

func TestMergeProcessor_Overwrite(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC).Truncate(time.Second)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now)
	event1.Props.SetText(ical.PropComment, "input")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "test")
	event2.Props.SetDate(ical.PropDateTimeStart, now)
	event2.Props.SetDate(ical.PropDateTimeEnd, now)
	event2.Props.SetText(ical.PropComment, "merge")
	merge := ical.NewCalendar()
	merge.Children = append(merge.Children, event2.Component)
	subject := NewMergeProcessor(*merge)
	subject.SetToolbox(NewToolbox())
	subject.MergeOption = MergeOptionOverWrite
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process got an error") {
		assert.Len(t, output.Children, 1, "Invalid number of events")
		assert.Equal(t, output.Children[0].Props.Get(ical.PropComment).Value, "merge", "Event from wrong calendar")
	}
}

func TestMergeProcessor_Skip(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC).Truncate(time.Second)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now)
	event1.Props.SetText(ical.PropComment, "input")
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "test")
	event2.Props.SetDate(ical.PropDateTimeStart, now)
	event2.Props.SetDate(ical.PropDateTimeEnd, now)
	event2.Props.SetText(ical.PropComment, "merge")
	merge := ical.NewCalendar()
	merge.Children = append(merge.Children, event2.Component)
	subject := NewMergeProcessor(*merge)
	subject.SetToolbox(NewToolbox())
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process got an error") {
		assert.Len(t, output.Children, 1, "Invalid number of events")
		assert.Equal(t, output.Children[0].Props.Get(ical.PropComment).Value, "input", "Event from wrong calendar")
	}
}

func TestMergeProcessor_Filter(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC).Truncate(time.Second)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "test2")
	now = time.Now().In(time.UTC).Truncate(time.Second)
	event2.Props.SetDate(ical.PropDateTimeStart, now)
	event2.Props.SetDate(ical.PropDateTimeEnd, now)
	event3 := ical.NewEvent()
	event3.Props.SetText(ical.PropSummary, "test3")
	now = time.Now().In(time.UTC).Truncate(time.Second)
	event3.Props.SetDate(ical.PropDateTimeStart, now)
	event3.Props.SetDate(ical.PropDateTimeEnd, now)
	merge := ical.NewCalendar()
	merge.Children = append(merge.Children, event2.Component, event3.Component)
	subject := NewMergeProcessor(*merge)
	toolbox := NewToolbox()
	toolbox.TextSelectorPattern = regexp.MustCompile("test[^3]")
	subject.SetToolbox(toolbox)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process got an error") {
		assert.Len(t, output.Children, 2, "Invalid number of events")
		assert.Equal(t, output.Children[0].Props.Get(ical.PropSummary).Value, "test", "Invalid first event")
		assert.Equal(t, output.Children[1].Props.Get(ical.PropSummary).Value, "test2", "Invalid second event")
	}
}
