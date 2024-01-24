package processors

import (
	"github.com/emersion/go-ical"
	"strings"
)

// The AddPropertyProcessor adds an ICS property to each selected event
type AddPropertyProcessor struct {
	PropertyName  string
	PropertyValue string
	Overwrite     bool
	toolbox       Toolbox
}

func (a *AddPropertyProcessor) SetToolbox(toolbox Toolbox) {
	a.toolbox = toolbox
}

func (a *AddPropertyProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	n := strings.ToUpper(a.PropertyName)
	for _, event := range input.Events() {
		if a.toolbox.EventMatchesSelector(event) {
			if event.Props.Get(n) != nil && a.Overwrite {
				event.Props.Del(n)
			}
			if event.Props.Get(n) == nil {
				event.Props.SetText(n, a.PropertyValue)
			}
		}
		output.Children = append(output.Children, event.Component)
	}
	return nil
}

var _ BaseProcessor = &AddPropertyProcessor{}
