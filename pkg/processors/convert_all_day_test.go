package processors

import (
	"github.com/dploeger/icarus/v2/test"
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestConvertAllDayProcessor_Process(t *testing.T) {
	toolbox := NewToolbox()
	subject := ConvertAllDayProcessor{
		TargetAllDay: true,
		StartTime:    "",
		EndTime:      "",
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
	if assert.NoError(t, err, "Process yielded error") {
		assert.Equal(t, ical.ValueDate, output.Children[0].Props.Get(ical.PropDateTimeStart).ValueType(), "Event was not turned to all-day event")
		eventStart, _ := output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)
		eventEnd, _ := output.Children[0].Props.Get(ical.PropDateTimeEnd).DateTime(time.UTC)
		assert.Equal(t, start.Truncate(24*time.Hour), eventStart)
		assert.Equal(t, end.Truncate(24*time.Hour), eventEnd)
	}
}

func TestConvertAllDayProcessor_ProcessToTimed(t *testing.T) {
	toolbox := NewToolbox()
	subject := ConvertAllDayProcessor{
		TargetAllDay: false,
		StartTime:    "19:00",
		EndTime:      "20:00",
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
	if assert.NoError(t, err, "Process yielded error") {
		assert.Equal(t, ical.ValueDateTime, output.Children[0].Props.Get(ical.PropDateTimeStart).ValueType(), "Event was not turned to all-day event")
		start := now.Truncate(24 * time.Hour).Add(19 * time.Hour)
		end := now.Truncate(24 * time.Hour).Add(20 * time.Hour)
		startEvent, _ := output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)
		endEvent, _ := output.Children[0].Props.Get(ical.PropDateTimeEnd).DateTime(time.UTC)
		assert.Equal(t, start, startEvent, "Event did not start as expected")
		assert.Equal(t, end, endEvent, "Event did not end as expected")
	}
}

func TestConvertAllDayProcessor_ProcessSingleDay(t *testing.T) {
	toolbox := NewToolbox()
	subject := ConvertAllDayProcessor{
		TargetAllDay: true,
		StartTime:    "",
		EndTime:      "",
	}
	subject.SetToolbox(toolbox)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	start := time.Date(1970, 1, 1, 10, 00, 00, 00, time.UTC)
	event1.Props.SetDateTime(ical.PropDateTimeStart, start)
	end := start.Add(2 * time.Hour)
	event1.Props.SetDateTime(ical.PropDateTimeEnd, end)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process yielded error") {
		assert.Equal(t, ical.ValueDate, output.Children[0].Props.Get(ical.PropDateTimeStart).ValueType(), "Event was not turned to all-day event")
		eventStart, _ := output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)
		assert.Equal(t, start.Truncate(24*time.Hour), eventStart)
		assert.Nil(t, output.Children[0].Props.Get(ical.PropDateTimeEnd), "DTEND was present on a single day")
	}
}

func TestConvertAllDayProcessor_ProcessToTimedCompressed(t *testing.T) {
	toolbox := NewToolbox()
	subject := ConvertAllDayProcessor{
		TargetAllDay: false,
		StartTime:    "19:00",
		EndTime:      "20:00",
		Compress:     true,
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
	if assert.NoError(t, err, "Process yielded error") {
		dtStart := test.IgnoreError(output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)).(time.Time)
		dtEnd := test.IgnoreError(output.Children[0].Props.Get(ical.PropDateTimeEnd).DateTime(time.UTC)).(time.Time)
		assert.Equal(t, dtStart.Day(), dtEnd.Day(), "Event was not compressed")
	}
}

func TestConvertAllDayProcessor_ProcessToTimedTimezone(t *testing.T) {
	toolbox := NewToolbox()
	var loc time.Location
	if l, err := time.LoadLocation("Europe/Berlin"); err != nil {
		assert.Fail(t, "Can't load Europe/Berlin location")
	} else {
		loc = *l
	}
	subject := ConvertAllDayProcessor{
		TargetAllDay: false,
		StartTime:    "19:00",
		EndTime:      "20:00",
		Compress:     true,
		Timezone:     loc,
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
	if assert.NoError(t, err, "Process yielded error") {
		dtStart := test.IgnoreError(output.Children[0].Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)).(time.Time)
		assert.Equal(t, 19, dtStart.Hour(), "Timezone was not respected")
	}
}
