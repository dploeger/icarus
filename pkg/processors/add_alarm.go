package processors

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"time"
)

type AddAlarmProcessor struct {
	alarmBefore *int
	toolbox     Toolbox
}

func (a *AddAlarmProcessor) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("addAlarm", "Add an alarm to all selected events")
	a.alarmBefore = command.Int("A", "alarm-before", &argparse.Options{
		Help:     "Alarm should be set number of minutes before the event",
		Required: true,
	})
	return command, nil
}

func (a *AddAlarmProcessor) SetToolbox(toolbox Toolbox) {
	a.toolbox = toolbox
}

func (a *AddAlarmProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	for _, event := range input.Events() {
		if a.toolbox.EventMatchesSelector(event) {
			valarm := ical.NewComponent(ical.CompAlarm)
			valarmTrigger := ical.NewProp(ical.PropTrigger)
			valarmTrigger.SetDuration(time.Duration(*a.alarmBefore) * time.Minute * -1)
			valarm.Props.Add(valarmTrigger)
			valarm.Props.SetText(ical.PropAction, ical.ParamDisplay)
			event.Children = append(event.Children, valarm)
		}
		output.Children = append(output.Children, event.Component)
	}
	return nil
}
