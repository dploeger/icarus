package processors

import (
	"github.com/emersion/go-ical"
	"strings"
)

// The DeletePropertyProcessor deletes an ICS property from all selected events
type DeletePropertyProcessor struct {
	PropertyName string
	toolbox      Toolbox
}

func (d *DeletePropertyProcessor) SetToolbox(toolbox Toolbox) {
	d.toolbox = toolbox
}

func (d *DeletePropertyProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	for _, event := range input.Events() {
		if d.toolbox.EventMatchesSelector(event) {
			if event.Props.Get(strings.ToUpper(d.PropertyName)) != nil {
				event.Props.Del(strings.ToUpper(d.PropertyName))
			}
		}
		output.Children = append(output.Children, event.Component)
	}
	return nil
}

var _ BaseProcessor = &DeletePropertyProcessor{}
