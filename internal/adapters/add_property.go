package adapters

import (
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
)

// The AddPropertyAdapter adds an ICS property to each selected event
type AddPropertyAdapter struct {
	propertyName  *string
	propertyValue *string
	overwrite     *bool
	toolbox       processors.Toolbox
}

func (a *AddPropertyAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	c := parser.NewCommand("addProperty", "Adds a new property to each selected event")
	a.propertyName = c.String("N", "name", &argparse.Options{
		Help:     "Name of the new property",
		Required: true,
	})
	a.propertyValue = c.String("V", "value", &argparse.Options{
		Help:     "Value of the new property (only text values allowed)",
		Required: true,
	})
	a.overwrite = c.Flag("O", "overwrite", &argparse.Options{
		Help:     "Overwrite property if it exists",
		Required: false,
		Default:  true,
	})
	return c, nil
}

func (a *AddPropertyAdapter) SetToolbox(toolbox processors.Toolbox) {
	a.toolbox = toolbox
}

func (a *AddPropertyAdapter) Process(input ical.Calendar, output *ical.Calendar) error {
	p := processors.AddPropertyProcessor{
		PropertyName:  *a.propertyName,
		PropertyValue: *a.propertyValue,
		Overwrite:     *a.overwrite,
	}
	p.SetToolbox(a.toolbox)
	return p.Process(input, output)
}

var _ Adapter = &AddPropertyAdapter{}
