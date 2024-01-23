package processors

import (
	"github.com/dploeger/icarus/internal"
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddDTStampProcessor_Process(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC).Truncate(time.Second)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	subject := AddDTStampProcessor{}
	subject.SetToolbox(NewToolbox())
	subject.timestamp = internal.StringAddr(now.Format("20060102T150405Z"))
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process got an error")
	assert.Len(t, output.Children, 1, "Invalid number of events")
	dtstamptime, _ := output.Children[0].Props.Get(ical.PropDateTimeStamp).DateTime(time.UTC)
	assert.Equal(t, now, dtstamptime)
}

func TestAddDTStampProcessor_ProcessNoOverwrite(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC).Truncate(time.Second)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now)
	event1.Props.SetDateTime(ical.PropDateTimeStamp, now)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	subject := AddDTStampProcessor{}
	subject.overwrite = internal.BoolAddr(false)
	subject.SetToolbox(NewToolbox())
	modifiedTime := now.Add(-1 * time.Hour)
	subject.timestamp = internal.StringAddr(modifiedTime.Format("20060102T150405Z"))
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process got an error")
	assert.Len(t, output.Children, 1, "Invalid number of events")
	dtstamptime, _ := output.Children[0].Props.Get(ical.PropDateTimeStamp).DateTime(time.UTC)
	assert.Equal(t, now, dtstamptime)
}

func TestAddDTStampProcessor_ProcessOverwrite(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC).Truncate(time.Second)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now)
	event1.Props.SetDateTime(ical.PropDateTimeStamp, now)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	subject := AddDTStampProcessor{}
	subject.overwrite = internal.BoolAddr(true)
	subject.SetToolbox(NewToolbox())
	modifiedTime := now.Add(-1 * time.Hour)
	subject.timestamp = internal.StringAddr(modifiedTime.Format("20060102T150405Z"))
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	assert.NoError(t, err, "Process got an error")
	assert.Len(t, output.Children, 1, "Invalid number of events")
	dtstamptime, _ := output.Children[0].Props.Get(ical.PropDateTimeStamp).DateTime(time.UTC)
	assert.Equal(t, modifiedTime, dtstamptime)
}
