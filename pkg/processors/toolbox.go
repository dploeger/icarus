package processors

// a toolbox of common functions

import (
	"github.com/emersion/go-ical"
	"github.com/sirupsen/logrus"
	"regexp"
	"time"
)

type Toolbox struct {
	// TextSelectorPattern is the user specified regexp pattern to select events by
	TextSelectorPattern *regexp.Regexp
	// TextSelectorProps is a list of property names that are checked for the selector
	TextSelectorProps []string
	// DateRangeSelectorStart is the start of a range of dates events have to match
	DateRangeSelectorStart time.Time
	// DateRangeSelectorEnd is the end of a range of dates events have to match
	DateRangeSelectorEnd time.Time
}

// EventMatchesSelector checks if the given event matches the configured filter
func (t Toolbox) EventMatchesSelector(event ical.Event) bool {
	if event.Props.Get(ical.PropDateTimeStart) != nil {
		dtStart, _ := event.Props.Get(ical.PropDateTimeStart).DateTime(time.UTC)
		if !t.DateRangeSelectorStart.IsZero() && dtStart.Before(t.DateRangeSelectorStart) {
			return false
		}
	}

	if event.Props.Get(ical.PropDateTimeEnd) != nil {
		dtEnd, _ := event.Props.Get(ical.PropDateTimeEnd).DateTime(time.UTC)
		if !t.DateRangeSelectorEnd.IsZero() && dtEnd.After(t.DateRangeSelectorEnd) {
			return false
		}
	}

	if t.TextSelectorPattern != nil {
		for _, prop := range t.TextSelectorProps {
			if event.Props.Get(prop) != nil && t.TextSelectorPattern.MatchString(event.Props.Get(prop).Value) {
				logrus.Debugf("Event %s matched selector", event.Name)
				return true
			}
		}
	}
	logrus.Debugf("Event %s DID NOT match selector", event.Name)
	return false
}

// NewToolbox returns a toolbox with default values
func NewToolbox() Toolbox {
	return Toolbox{
		TextSelectorPattern: regexp.MustCompile(".*"),
		TextSelectorProps:   []string{ical.PropSummary, ical.PropDescription},
	}
}
