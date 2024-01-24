package adapters

import (
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
)

// The DeletePropertyAdapter deletes an ICS property from all selected events
type DeletePropertyAdapter struct {
	propertyName *string
	toolbox      processors.Toolbox
}

func (d *DeletePropertyAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	c := parser.NewCommand("deleteProperty", "Deletes a property from all selected events")
	d.propertyName = c.String("P", "property", &argparse.Options{
		Help:     "Property to delete (e.g. X-SPECIAL-PROP, CATEGORIES)",
		Required: true,
	})
	return c, nil
}

func (d *DeletePropertyAdapter) SetToolbox(toolbox processors.Toolbox) {
	d.toolbox = toolbox
}

func (d *DeletePropertyAdapter) Process(input ical.Calendar, output *ical.Calendar) error {
	p := processors.DeletePropertyProcessor{PropertyName: *d.propertyName}
	p.SetToolbox(d.toolbox)
	return p.Process(input, output)
}

var _ Adapter = &DeletePropertyAdapter{}
