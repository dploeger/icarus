package processors

import (
	"fmt"
	"github.com/emersion/go-ical"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

// The ConvertAllDayProcessor converts all-day events to timed events or vice versa
type ConvertAllDayProcessor struct {
	TargetAllDay bool
	StartTime    string
	EndTime      string
	Compress     bool
	Timezone     time.Location
	toolbox      Toolbox
}

func (c *ConvertAllDayProcessor) SetToolbox(toolbox Toolbox) {
	c.toolbox = toolbox
}

func (c *ConvertAllDayProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	if !c.TargetAllDay {
		t := regexp.MustCompile(`\d{2}:\d{2}`)
		if !t.MatchString(c.StartTime) {
			return fmt.Errorf("can not interpret start time %s", c.StartTime)
		}
		if !t.MatchString(c.EndTime) {
			return fmt.Errorf("can not interpret start time %s", c.EndTime)
		}
	}
	for _, event := range input.Events() {
		if c.toolbox.EventMatchesSelector(event) {
			isAllDay := event.Props.Get(ical.PropDateTimeStart).ValueType() == ical.ValueDate
			if c.TargetAllDay && !isAllDay {
				logrus.Debugf("Converting event %s to an all day event", event.Name)
				if startTimestamp, endTimestamp, err := c.getTimestamps(event, c.Timezone); err != nil {
					return err
				} else {
					event.Props.SetDate(ical.PropDateTimeStart, startTimestamp)
					if startTimestamp.Truncate(24*time.Hour).Compare(endTimestamp.Truncate(24*time.Hour)) == 0 {
						logrus.Debugf("Start and time end match, removing DTEND property")
						event.Props.Del(ical.PropDateTimeEnd)
					} else {
						logrus.Debugf("Start and time end differ (%s != %s), setting DTEND to last day", startTimestamp.Truncate(24*time.Hour), endTimestamp.Truncate(24*time.Hour))
						event.Props.SetDate(ical.PropDateTimeEnd, endTimestamp)
					}
				}
			} else if !c.TargetAllDay && isAllDay {
				logrus.Debugf("Converting event %s to a timed event", event.Name)
				if startTimestamp, endTimestamp, err := c.getTimestamps(event, c.Timezone); err != nil {
					return err
				} else if startDuration, endDuration, err := c.getStartAndEndDuration(); err != nil {
					return err
				} else {
					if c.Compress {
						endTimestamp = startTimestamp
					}
					event.Props.SetDateTime(ical.PropDateTimeStart, startTimestamp.Add(startDuration))
					event.Props.SetDateTime(ical.PropDateTimeEnd, endTimestamp.Add(endDuration))
				}
			}
			output.Children = append(output.Children, event.Component)
		}
	}
	return nil
}

// converts the start and end timestamps to time.Time objects
func (c *ConvertAllDayProcessor) getTimestamps(event ical.Event, loc time.Location) (time.Time, time.Time, error) {
	var startTimestamp time.Time
	var endTimestamp time.Time

	if timestamp, err := event.Props.Get(ical.PropDateTimeStart).DateTime(&loc); err != nil {
		return startTimestamp, endTimestamp, err
	} else {
		startTimestamp = timestamp
	}
	if timestamp, err := event.Props.Get(ical.PropDateTimeEnd).DateTime(&loc); err != nil {
		return startTimestamp, endTimestamp, err
	} else {
		endTimestamp = timestamp
	}
	return startTimestamp, endTimestamp, nil
}

// converts the given HH:MM Timestamp to a time.Duration
func (c *ConvertAllDayProcessor) getDurationFromTimeString(timeString string) (time.Duration, error) {
	timeParts := strings.Split(timeString, ":")
	return time.ParseDuration(fmt.Sprintf("%sh%sm", timeParts[0], timeParts[1]))
}

// gets the start and end time HH:MM timestamps to time.Duration
func (c *ConvertAllDayProcessor) getStartAndEndDuration() (time.Duration, time.Duration, error) {
	var startDuration time.Duration
	var endDuration time.Duration
	if duration, err := c.getDurationFromTimeString(c.StartTime); err != nil {
		return startDuration, endDuration, err
	} else {
		startDuration = duration
	}
	if duration, err := c.getDurationFromTimeString(c.EndTime); err != nil {
		return startDuration, endDuration, err
	} else {
		endDuration = duration
	}
	return startDuration, endDuration, nil
}

var _ BaseProcessor = &ConvertAllDayProcessor{}
