package processors

import (
	"github.com/emersion/go-ical"
)

// The FilterProcessor filters the calendar for selected events
type FilterProcessor struct {
	Inverse bool
	toolbox Toolbox
}

func (f *FilterProcessor) Process(inputCalendar ical.Calendar, outputCalendar *ical.Calendar) error {
	for _, event := range inputCalendar.Events() {
		if (f.Inverse && !f.toolbox.EventMatchesSelector(event)) ||
			(!f.Inverse && f.toolbox.EventMatchesSelector(event)) {
			outputCalendar.Children = append(outputCalendar.Children, event.Component)
			continue
		}
	}
	return nil
}

func (f *FilterProcessor) SetToolbox(toolbox Toolbox) {
	f.toolbox = toolbox
}

var _ BaseProcessor = &FilterProcessor{}
