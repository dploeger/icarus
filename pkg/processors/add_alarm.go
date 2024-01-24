package processors

import (
	"github.com/emersion/go-ical"
	"time"
)

// The AddAlarmProcessor adds an alarm definition to all selected events
type AddAlarmProcessor struct {
	AlarmBefore int
	toolbox     Toolbox
}

func (a *AddAlarmProcessor) SetToolbox(toolbox Toolbox) {
	a.toolbox = toolbox
}

func (a *AddAlarmProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	for _, event := range input.Events() {
		if a.toolbox.EventMatchesSelector(event) {
			valarm := ical.NewComponent(ical.CompAlarm)
			valarmTrigger := ical.NewProp(ical.PropTrigger)
			valarmTrigger.SetDuration(time.Duration(a.AlarmBefore) * time.Minute * -1)
			valarm.Props.Add(valarmTrigger)
			valarm.Props.SetText(ical.PropAction, ical.ParamDisplay)
			event.Children = append(event.Children, valarm)
		}
		output.Children = append(output.Children, event.Component)
	}
	return nil
}

var _ BaseProcessor = &AddAlarmProcessor{}
