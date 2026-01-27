package adapters

import (
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
)

// The CopyPropertyAdapter adds an ICS property to each selected event
type CopyPropertyAdapter struct {
	sourceProperty *string
	targetProperty *string
	overwrite      *bool
	toolbox        processors.Toolbox
}

func (a *CopyPropertyAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	c := parser.NewCommand("copyProperty", "Copies a property value to another property")
	a.sourceProperty = c.String("S", "source", &argparse.Options{
		Help:     "Name of the source property to copy",
		Required: true,
	})
	a.targetProperty = c.String("T", "target", &argparse.Options{
		Help:     "Name of the target property",
		Required: true,
	})
	a.overwrite = c.Flag("O", "overwrite", &argparse.Options{
		Help:     "Overwrite property if it exists",
		Required: false,
		Default:  true,
	})
	return c, nil
}

func (a *CopyPropertyAdapter) SetToolbox(toolbox processors.Toolbox) {
	a.toolbox = toolbox
}

func (a *CopyPropertyAdapter) Process(input ical.Calendar, output *ical.Calendar) error {
	p := processors.CopyPropertyProcessor{
		SourceProperty: *a.sourceProperty,
		TargetProperty: *a.targetProperty,
		Overwrite:      *a.overwrite,
	}
	p.SetToolbox(a.toolbox)
	return p.Process(input, output)
}

var _ Adapter = &CopyPropertyAdapter{}
