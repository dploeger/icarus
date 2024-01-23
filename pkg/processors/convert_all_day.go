package processors

import (
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/emersion/go-ical"
	"github.com/sirupsen/logrus"
	"regexp"
	"strings"
	"time"
)

// The ConvertAllDayProcessor converts all-day events to timed events or vice versa
type ConvertAllDayProcessor struct {
	toolbox      Toolbox
	targetAllDay *bool
	startTime    *string
	endTime      *string
	compress     *bool
	timezone     *string
}

func (c *ConvertAllDayProcessor) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
	command := parser.NewCommand("convertAllDay", "Convert an all day ")
	c.targetAllDay = command.Flag("A", "all-day", &argparse.Options{
		Help:    "If set, will convert from timed to all-day events. If not, will convert from all-day to timed events",
		Default: true,
	})
	hourRegExp := regexp.MustCompile("[0-9]{2}:[0-9]{2}")
	c.startTime = command.String("S", "start", &argparse.Options{
		Help: "Start time in a 24 hour format (hh:mm)",
		Validate: func(args []string) error {
			if len(args) != 1 {
				return errors.New("start only accepts one argument")
			}
			if !hourRegExp.MatchString(args[0]) {
				return errors.New("start requires a 24 hour format (hh:mm)")
			}
			return nil
		},
	})
	c.endTime = command.String("E", "end", &argparse.Options{
		Help: "End time in a 24 hour format (hh:mm)",
		Validate: func(args []string) error {
			if len(args) != 1 {
				return errors.New("end only accepts one argument")
			}
			if !hourRegExp.MatchString(args[0]) {
				return errors.New("end requires a 24 hour format (hh:mm)")
			}
			return nil
		},
	})
	c.compress = command.Flag("C", "compress", &argparse.Options{
		Help:    "If the event spans multiple dates, only use the start date when converting to a timed event",
		Default: false,
	})
	c.timezone = command.String("T", "timezone", &argparse.Options{
		Help: "IANA timezone name for the new times",
		Validate: func(args []string) error {
			_, err := time.LoadLocation(args[0])
			if err != nil {
				return errors.New(fmt.Sprintf("%s is not a valid timezone name", args[0]))
			}
			return nil
		},
	})
	return command, nil
}

func (c *ConvertAllDayProcessor) SetToolbox(toolbox Toolbox) {
	c.toolbox = toolbox
}

func (c *ConvertAllDayProcessor) Process(input ical.Calendar, output *ical.Calendar) error {
	loc := time.UTC
	if c.timezone != nil {
		if l, err := time.LoadLocation(*c.timezone); err != nil {
			return err
		} else {
			loc = l
		}
	}
	for _, event := range input.Events() {
		if c.toolbox.EventMatchesSelector(event) {
			isAllDay := event.Props.Get(ical.PropDateTimeStart).ValueType() == ical.ValueDate
			if *c.targetAllDay && !isAllDay {
				logrus.Debugf("Converting event %s to an all day event", event.Name)
				if startTimestamp, endTimestamp, err := c.getTimestamps(event, loc); err != nil {
					return err
				} else {
					event.Props.SetDate(ical.PropDateTimeStart, startTimestamp)
					if startTimestamp.Truncate(24*time.Hour) == endTimestamp.Truncate(24*time.Hour) {
						event.Props.Del(ical.PropDateTimeEnd)
					} else {
						event.Props.SetDate(ical.PropDateTimeEnd, endTimestamp)
					}
				}
			} else if !*c.targetAllDay && isAllDay {
				logrus.Debugf("Converting event %s to a timed event", event.Name)
				if startTimestamp, endTimestamp, err := c.getTimestamps(event, loc); err != nil {
					return err
				} else if startDuration, endDuration, err := c.getStartAndEndDuration(); err != nil {
					return err
				} else {
					if c.compress != nil && *c.compress {
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
func (c *ConvertAllDayProcessor) getTimestamps(event ical.Event, loc *time.Location) (time.Time, time.Time, error) {
	var startTimestamp time.Time
	var endTimestamp time.Time

	if timestamp, err := event.Props.Get(ical.PropDateTimeStart).DateTime(loc); err != nil {
		return startTimestamp, endTimestamp, err
	} else {
		startTimestamp = timestamp
	}
	if timestamp, err := event.Props.Get(ical.PropDateTimeEnd).DateTime(loc); err != nil {
		return startTimestamp, endTimestamp, err
	} else {
		endTimestamp = timestamp
	}
	return startTimestamp, endTimestamp, nil
}

// converts the given HH:MM timestamp to a time.Duration
func (c *ConvertAllDayProcessor) getDurationFromTimeString(timeString string) (time.Duration, error) {
	timeParts := strings.Split(timeString, ":")
	return time.ParseDuration(fmt.Sprintf("%sh%sm", timeParts[0], timeParts[1]))
}

// gets the start and end time HH:MM timestamps to time.Duration
func (c *ConvertAllDayProcessor) getStartAndEndDuration() (time.Duration, time.Duration, error) {
	var startDuration time.Duration
	var endDuration time.Duration
	if duration, err := c.getDurationFromTimeString(*c.startTime); err != nil {
		return startDuration, endDuration, err
	} else {
		startDuration = duration
	}
	if duration, err := c.getDurationFromTimeString(*c.endTime); err != nil {
		return startDuration, endDuration, err
	} else {
		endDuration = duration
	}
	return startDuration, endDuration, nil
}

var _ BaseProcessor = &ConvertAllDayProcessor{}
