package adapters

import (
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
	"regexp"
	"time"
)

// The ConvertAllDayAdapter converts all-day events to timed events or vice versa
type ConvertAllDayAdapter struct {
	toolbox      processors.Toolbox
	targetAllDay *bool
	startTime    *string
	endTime      *string
	compress     *bool
	timezone     *string
}

func (c *ConvertAllDayAdapter) Initialize(parser *argparse.Parser) (*argparse.Command, error) {
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

func (c *ConvertAllDayAdapter) SetToolbox(toolbox processors.Toolbox) {
	c.toolbox = toolbox
}

func (c *ConvertAllDayAdapter) Process(input ical.Calendar, output *ical.Calendar) error {
	loc := time.UTC
	if c.timezone != nil {
		if l, err := time.LoadLocation(*c.timezone); err != nil {
			return err
		} else {
			loc = l
		}
	}
	p := processors.ConvertAllDayProcessor{
		TargetAllDay: *c.targetAllDay,
		StartTime:    *c.startTime,
		EndTime:      *c.endTime,
		Compress:     *c.compress,
		Timezone:     *loc,
	}
	p.SetToolbox(c.toolbox)
	return p.Process(input, output)
}

var _ Adapter = &ConvertAllDayAdapter{}
