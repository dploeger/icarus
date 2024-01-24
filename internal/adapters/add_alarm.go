package adapters

import (
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
)

// The AddAlarmAdapter adds an alarm definition to all selected events
type AddAlarmAdapter struct {
	alarmBefore *int
	toolbox     processors.Toolbox
}

func (a *AddAlarmAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("addAlarm", "Add an alarm to all selected events")
	a.alarmBefore = command.Int("A", "alarm-before", &argparse.Options{
		Help:     "Alarm should be set number of minutes before the event",
		Required: true,
	})
	return command, nil
}

func (a *AddAlarmAdapter) SetToolbox(toolbox processors.Toolbox) {
	a.toolbox = toolbox
}

func (a *AddAlarmAdapter) Process(input ical.Calendar, output *ical.Calendar) error {
	p := processors.AddAlarmProcessor{AlarmBefore: *a.alarmBefore}
	p.SetToolbox(a.toolbox)
	return p.Process(input, output)
}

var _ Adapter = &AddAlarmAdapter{}
