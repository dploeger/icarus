package adapters

import (
	"os"

	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
)

// The FilterAdapter filters the calendar for selected events
type MergeAdapter struct {
	toolbox       processors.Toolbox
	mergeCalendar *os.File
	overwrite     *bool
	mergeProps    *[]string
}

func (a *MergeAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("merge", "Merge another calendar into the input calendar. The filter options work on the additional calendar for this comand.")
	a.mergeCalendar = command.File("I", "merge-calendar", os.O_RDONLY, 0600, &argparse.Options{
		Help:     "Calendar to merge into the input calendar",
		Required: true,
	})
	a.overwrite = command.Flag("O", "overwrite", &argparse.Options{
		Help: "Overwrite events in the input calendar with those from the other calendar",
	})
	a.mergeProps = command.StringList("P", "merge-props", &argparse.Options{
		Help:    "The ical event properties which decide if two events are the same",
		Default: []string{ical.PropSummary, ical.PropDateTimeStart, ical.PropDateTimeEnd},
	})
	return command, nil
}

func (a *MergeAdapter) Process(input ical.Calendar, output *ical.Calendar) error {
	var mergeCalendar ical.Calendar
	dec := ical.NewDecoder(a.mergeCalendar)
	if cal, err := dec.Decode(); err != nil {
		return err
	} else {
		mergeCalendar = *cal
	}
	p := processors.NewMergeProcessor(mergeCalendar)
	p.MergeProps = *a.mergeProps
	if *a.overwrite {
		p.MergeOption = processors.MergeOptionOverWrite
	}
	p.SetToolbox(a.toolbox)
	return p.Process(input, output)
}

func (f *MergeAdapter) SetToolbox(toolbox processors.Toolbox) {
	f.toolbox = toolbox
}

var _ Adapter = &MergeAdapter{}
