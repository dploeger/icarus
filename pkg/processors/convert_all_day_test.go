package processors

import (
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	"icarus/internal"
	"testing"
	"time"
)

func TestConvertAllDayProcessor_Process(t *testing.T) {
	toolbox := NewToolbox()
	subject := ConvertAllDayProcessor{
		targetAllDay: internal.BoolAddr(true),
		startTime:    nil,
		endTime:      nil,
	}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	start := time.Now().In(time.UTC)
	event1.Props.SetDateTime(ical.PropDateTimeStart, start)
	end := start.Add(2 * 24 * time.Hour)
	event1.Props.SetDateTime(ical.PropDateTimeEnd, end)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process yielded error")
	assert.Equal(t, ical.ValueDate, output.Children[0].Props.Get(ical.PropDateTimeStart).ValueType(), "Event was not turned to all-day event")
	eventStart, _ := output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)
	eventEnd, _ := output.Children[0].Props.Get(ical.PropDateTimeEnd).DateTime(time.UTC)
	assert.Equal(t, start.Truncate(24*time.Hour), eventStart)
	assert.Equal(t, end.Truncate(24*time.Hour), eventEnd)
}

func TestConvertAllDayProcessor_ProcessToTimed(t *testing.T) {
	toolbox := NewToolbox()
	subject := ConvertAllDayProcessor{
		targetAllDay: internal.BoolAddr(false),
		startTime:    internal.StringAddr("19:00"),
		endTime:      internal.StringAddr("20:00"),
	}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process yielded error")
	assert.Equal(t, ical.ValueDateTime, output.Children[0].Props.Get(ical.PropDateTimeStart).ValueType(), "Event was not turned to all-day event")
	start := now.Truncate(24 * time.Hour).Add(19 * time.Hour)
	end := now.Truncate(24 * time.Hour).Add(20 * time.Hour)
	startEvent, _ := output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)
	endEvent, _ := output.Children[0].Props.Get(ical.PropDateTimeEnd).DateTime(time.UTC)
	assert.Equal(t, start, startEvent, "Event did not start as expected")
	assert.Equal(t, end, endEvent, "Event did not end as expected")
}

func TestConvertAllDayProcessor_ProcessSingleDay(t *testing.T) {
	toolbox := NewToolbox()
	subject := ConvertAllDayProcessor{
		targetAllDay: internal.BoolAddr(true),
		startTime:    nil,
		endTime:      nil,
	}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	start := time.Now().In(time.UTC)
	event1.Props.SetDateTime(ical.PropDateTimeStart, start)
	end := start.Add(2 * time.Hour)
	event1.Props.SetDateTime(ical.PropDateTimeEnd, end)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process yielded error")
	assert.Equal(t, ical.ValueDate, output.Children[0].Props.Get(ical.PropDateTimeStart).ValueType(), "Event was not turned to all-day event")
	eventStart, _ := output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)
	assert.Equal(t, start.Truncate(24*time.Hour), eventStart)
	assert.Nil(t, output.Children[0].Props.Get(ical.PropDateTimeEnd), "DTEND was present on a single day")
}

func TestConvertAllDayProcessor_ProcessToTimedCompressed(t *testing.T) {
	toolbox := NewToolbox()
	subject := ConvertAllDayProcessor{
		targetAllDay: internal.BoolAddr(false),
		startTime:    internal.StringAddr("19:00"),
		endTime:      internal.StringAddr("20:00"),
		compress:     internal.BoolAddr(true),
	}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now.Add(1*24*time.Hour))
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process yielded error")
	dtStart := internal.IgnoreError(output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)).(time.Time)
	dtEnd := internal.IgnoreError(output.Children[0].Props.Get(ical.PropDateTimeEnd).DateTime(time.UTC)).(time.Time)
	assert.Equal(t, dtStart.Day(), dtEnd.Day(), "Event was not compressed")
}

func TestConvertAllDayProcessor_ProcessToTimedTimezone(t *testing.T) {
	toolbox := NewToolbox()
	subject := ConvertAllDayProcessor{
		targetAllDay: internal.BoolAddr(false),
		startTime:    internal.StringAddr("19:00"),
		endTime:      internal.StringAddr("20:00"),
		compress:     internal.BoolAddr(true),
		timezone:     internal.StringAddr("Europe/Berlin"),
	}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now.Add(1*24*time.Hour))
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process yielded error")
	dtStart := internal.IgnoreError(output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)).(time.Time)
	assert.Equal(t, 19, dtStart.Hour(), "Timezone was not respected")
}
