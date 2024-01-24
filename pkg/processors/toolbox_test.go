package processors

import (
	"github.com/dploeger/icarus/v2/test"
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
	"time"
)

func TestToolbox_EventMatchesSelector(t *testing.T) {
	subject := NewToolbox()
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	assert.True(t, subject.EventMatchesSelector(*event1), "Selector filtered out test event")
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "something")
	assert.True(t, subject.EventMatchesSelector(*event2), "Selector filtered out test event")
}

func TestToolbox_EventMatchesSelectorSelectsEvents(t *testing.T) {
	subject := NewToolbox()
	subject.TextSelectorPattern = regexp.MustCompile("test")
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	assert.True(t, subject.EventMatchesSelector(*event1), "Selector filtered out test event")
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "not")
	assert.False(t, subject.EventMatchesSelector(*event2), "Selector did not filter out test event")
}

func TestToolbox_EventMatchesSelectorSelectsDates(t *testing.T) {
	subject := NewToolbox()
	subject.DateRangeSelectorStart = test.IgnoreError(time.Parse("20060102T150405Z", "20220126T080000Z")).(time.Time)
	subject.DateRangeSelectorEnd = test.IgnoreError(time.Parse("20060102T150405Z", "20220127T080000Z")).(time.Time)
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	event1.Props.SetDateTime(ical.PropDateTimeStart, test.IgnoreError(time.Parse("20060102T150405Z", "20220126T090000Z")).(time.Time))
	event1.Props.SetDateTime(ical.PropDateTimeEnd, test.IgnoreError(time.Parse("20060102T150405Z", "20220126T100000Z")).(time.Time))
	assert.True(t, subject.EventMatchesSelector(*event1), "Selector filtered out test event")
	event2 := ical.NewEvent()
	event2.Props.SetText(ical.PropSummary, "not")
	event2.Props.SetDateTime(ical.PropDateTimeStart, test.IgnoreError(time.Parse("20060102T150405Z", "20220127T090000Z")).(time.Time))
	event2.Props.SetDateTime(ical.PropDateTimeEnd, test.IgnoreError(time.Parse("20060102T150405Z", "20220127T100000Z")).(time.Time))
	assert.False(t, subject.EventMatchesSelector(*event2), "Selector did not filter out test event")
}
