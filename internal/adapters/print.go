package adapters

import (
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
)

// The PrintAdapter simply outputs the input calendar
type PrintAdapter struct {
	toolbox processors.Toolbox
}

func (f *PrintAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("print", "Outputs the calendar")
	return command, nil
}

func (f *PrintAdapter) Process(input ical.Calendar, output *ical.Calendar) error {
	p := processors.PrintProcessor{}
	p.SetToolbox(f.toolbox)
	return p.Process(input, output)
}

func (f *PrintAdapter) SetToolbox(toolbox processors.Toolbox) {
	f.toolbox = toolbox
}

var _ Adapter = &PrintAdapter{}
