package processors

import (
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"strings"
)

type AddPropertyProcessor struct {
	propertyName  *string
	propertyValue *string
	overwrite     *bool
	toolbox       Toolbox
}

func (a *AddPropertyProcessor) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
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

func (a *AddPropertyProcessor) SetToolbox(toolbox Toolbox) {
	a.toolbox = toolbox
}

func (a *AddPropertyProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	n := strings.ToUpper(*a.propertyName)
	for _, event := range input.Events() {
		if a.toolbox.EventMatchesSelector(event) {
			if event.Props.Get(n) != nil && *a.overwrite {
				event.Props.Del(n)
			}
			if event.Props.Get(n) == nil {
				event.Props.SetText(n, *a.propertyValue)
			}
		}
		output.Children = append(output.Children, event.Component)
	}
	return nil
}
