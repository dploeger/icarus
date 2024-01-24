package processors

import (
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAddAlarmProcessor_Process(t *testing.T) {
	event1 := ical.NewEvent()
	event1.Props.SetText(ical.PropSummary, "test")
	now := time.Now().In(time.UTC).Truncate(time.Second)
	event1.Props.SetDate(ical.PropDateTimeStart, now)
	event1.Props.SetDate(ical.PropDateTimeEnd, now)
	input := ical.NewCalendar()
	input.Children = append(input.Children, event1.Component)
	subject := AddAlarmProcessor{}
	subject.SetToolbox(NewToolbox())
	subject.AlarmBefore = 15
	output := ical.NewCalendar()
	err := subject.Process(*input, output)
	if assert.NoError(t, err, "Process got an error") {
		assert.Len(t, output.Children, 1, "Invalid number of events")
		assert.Len(t, output.Children[0].Children, 1, "No alarm found")
		alarmDuration, _ := output.Children[0].Children[0].Props.Get(ical.PropTrigger).Duration()
		assert.Equal(t, 15*time.Minute*-1, alarmDuration)
	}
}
