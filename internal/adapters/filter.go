package adapters

import (
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
)

// The FilterAdapter filters the calendar for selected events
type FilterAdapter struct {
	toolbox processors.Toolbox
	inverse *bool
}

func (f *FilterAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("filter", "Output only events that match the selector")
	f.inverse = command.Flag("I", "inverse", &argparse.Options{
		Help:    "Inverse the function. Output events that DO NOT match the selector",
		Default: false,
	})
	return command, nil
}

func (f *FilterAdapter) Process(input ical.Calendar, output *ical.Calendar) error {
	p := processors.FilterProcessor{Inverse: *f.inverse}
	p.SetToolbox(f.toolbox)
	return p.Process(input, output)
}

func (f *FilterAdapter) SetToolbox(toolbox processors.Toolbox) {
	f.toolbox = toolbox
}

var _ Adapter = &FilterAdapter{}
