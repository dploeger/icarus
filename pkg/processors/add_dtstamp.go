package processors

import (
	"github.com/emersion/go-ical"
	"time"
)

// The AddDTStampProcessor adds a DTSTAMP field to all selected events
type AddDTStampProcessor struct {
	Timestamp time.Time
	Overwrite bool
	toolbox   Toolbox
}

func (t *AddDTStampProcessor) SetToolbox(toolbox Toolbox) {
	t.toolbox = toolbox
}

func (t *AddDTStampProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	for _, event := range input.Events() {
		if t.toolbox.EventMatchesSelector(event) {
			if event.Props.Get(ical.PropDateTimeStamp) == nil || t.Overwrite {
				event.Props.SetDateTime(ical.PropDateTimeStamp, t.Timestamp)
			}
		}
		output.Children = append(output.Children, event.Component)
	}
	return nil
}

var _ BaseProcessor = &AddDTStampProcessor{}
