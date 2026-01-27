package processors

import (
	"github.com/emersion/go-ical"
)

// The CopyPropertyProcessor converts all-day events to timed events or vice versa
type CopyPropertyProcessor struct {
	SourceProperty string
	TargetProperty string
	Overwrite      bool
	toolbox        Toolbox
}

func (c *CopyPropertyProcessor) SetToolbox(toolbox Toolbox) {
	c.toolbox = toolbox
}

func (c *CopyPropertyProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	for _, event := range input.Events() {
		if c.toolbox.EventMatchesSelector(event) {
			source := event.Props.Get(c.SourceProperty)
			target := event.Props.Get(c.TargetProperty)
			if source != nil && (target == nil || c.Overwrite) {
				source.Name = c.TargetProperty
				event.Props.Set(source)
			}
		}
		output.Children = append(output.Children, event.Component)
	}
	return nil
}

var _ BaseProcessor = &CopyPropertyProcessor{}
