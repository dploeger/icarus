package processors

import (
	"strings"

	"github.com/emersion/go-ical"
	"github.com/thoas/go-funk"
)

type MergeOption int

const (
	MergeOptionOverWrite MergeOption = iota
	MergeOptionSkip
)

// The AddAlarmProcessor adds an alarm definition to all selected events
type MergeProcessor struct {
	MergeCalendar ical.Calendar
	MergeOption   MergeOption
	MergeProps    []string
	toolbox       Toolbox
}

func NewMergeProcessor(calendar ical.Calendar) MergeProcessor {
	return MergeProcessor{
		MergeCalendar: calendar,
		MergeOption:   MergeOptionSkip,
		MergeProps:    []string{ical.PropSummary, ical.PropDateTimeStart, ical.PropDateTimeEnd},
	}
}

func (a *MergeProcessor) SetToolbox(toolbox Toolbox) {
	a.toolbox = toolbox
}

func (a *MergeProcessor) eventMatches(eva ical.Event, evb ical.Event) bool {
	return a.getEventID(eva) == a.getEventID(evb)
}

func (a *MergeProcessor) getEventID(ev ical.Event) string {
	var propValues []string
	for _, mergeProp := range a.MergeProps {
		propValue := ev.Props.Get(mergeProp)
		if propValue != nil {
			propValues = append(propValues, propValue.Value)
		}
	}
	return strings.Join(propValues, ":")
}

func (a *MergeProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	var skipEventsInput []string
	var skipEventsMerge []string
	for _, mergeEvent := range a.MergeCalendar.Events() {
		if a.toolbox.EventMatchesSelector(mergeEvent) {
			for _, inputEvent := range input.Events() {
				if a.eventMatches(mergeEvent, inputEvent) {
					if a.MergeOption == MergeOptionOverWrite {
						skipEventsInput = append(skipEventsInput, a.getEventID(inputEvent))
					} else {
						skipEventsMerge = append(skipEventsMerge, a.getEventID(mergeEvent))
					}

				}
			}
		}
	}
	for _, event := range input.Events() {
		if !funk.ContainsString(skipEventsInput, a.getEventID(event)) {
			output.Children = append(output.Children, event.Component)
		}
	}
	for _, event := range a.MergeCalendar.Events() {
		if a.toolbox.EventMatchesSelector(event) && !funk.ContainsString(skipEventsMerge, a.getEventID(event)) {
			output.Children = append(output.Children, event.Component)
		}
	}
	return nil
}

var _ BaseProcessor = &AddAlarmProcessor{}
