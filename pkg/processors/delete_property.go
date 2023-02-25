package processors

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"strings"
)

type DeletePropertyProcessor struct {
	propertyName *string
	toolbox      Toolbox
}

func (d *DeletePropertyProcessor) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	c := parser.NewCommand("deleteProperty", "Deletes a property from all selected events")
	d.propertyName = c.String("P", "property", &argparse.Options{
		Help:     "Property to delete (e.g. X-SPECIAL-PROP, CATEGORIES)",
		Required: true,
	})
	return c, nil
}

func (d *DeletePropertyProcessor) SetToolbox(toolbox Toolbox) {
	d.toolbox = toolbox
}

func (d *DeletePropertyProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	for _, event := range input.Events() {
		if d.toolbox.EventMatchesSelector(event) {
			if event.Props.Get(strings.ToUpper(*d.propertyName)) != nil {
				event.Props.Del(strings.ToUpper(*d.propertyName))
			}
		}
		output.Children = append(output.Children, event.Component)
	}
	return nil
}
