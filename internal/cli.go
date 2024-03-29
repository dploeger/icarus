package internal

import (
	"fmt"
	"github.com/akamensky/argparse"
	"github.com/dploeger/icarus/v2/internal/adapters"
	"github.com/dploeger/icarus/v2/pkg/outputTypes"
	"github.com/dploeger/icarus/v2/pkg/processors"
	"github.com/emersion/go-ical"
	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"
	"os"
	"regexp"
	"time"
)

type processorCommand struct {
	adapter adapters.Adapter
	command *argparse.Command
}

func Main() int {
	availableOutputTypes := outputTypes.GetOutputTypes()
	parser := argparse.NewParser("icarus", "iCal file adapter")
	inputFile := parser.File("f", "file", os.O_RDONLY, 0444, &argparse.Options{
		Help: "File to read ics data from. Defaults to stdin",
	})
	outputFile := parser.File("o", "output", os.O_RDWR, 0644, &argparse.Options{
		Help: "File to write ics data to. Defaults to stdout",
	})
	selector := parser.String("s", "selector", &argparse.Options{
		Default: ".*",
		Help:    "Regular Expression pattern to select events by their summary, description, etc.",
	})
	selectorProps := parser.StringList("p", "selector-props", &argparse.Options{
		Help:    "Event properties that are searched using the text selector pattern",
		Default: []string{ical.PropSummary, ical.PropDescription},
	})
	dateSelectorStart := parser.String("b", "timestamp-start", &argparse.Options{
		Help: "An RFC3339-formatted (2006-01-02T15:04:05+07:00) timestamp that selects only events starting at or after that time",
		Validate: func(args []string) error {
			_, err := time.Parse(time.RFC3339, args[0])
			if err != nil {
				return fmt.Errorf("can not parse start timestamp: %w", err)
			}
			return nil
		},
	})
	dateSelectorEnd := parser.String("e", "timestamp-end", &argparse.Options{
		Help: "An RFC3339-formatted (2006-01-02T15:04:05+07:00) timestamp that selects only events ending at or before that time",
		Validate: func(args []string) error {
			_, err := time.Parse(time.RFC3339, args[0])
			if err != nil {
				return fmt.Errorf("can not parse end timestamp: %w", err)
			}
			return nil
		},
	})
	outputType := parser.Selector("t", "output-type", funk.Keys(availableOutputTypes).([]string), &argparse.Options{
		Help:    fmt.Sprintf("Type of output. Valid types:\n%s\n\t\t\t", outputTypes.GetOutputHelp()),
		Default: "ics",
	})
	columns := parser.StringList("c", "columns", &argparse.Options{
		Help:    "Columns to output (used by list and csv outputs)",
		Default: []string{"SUMMARY", "DTSTART", "DTEND", "DESCRIPTION"},
	})
	logLevel := parser.String("l", "loglevel", &argparse.Options{
		Help:    "Loglevel to use",
		Default: "error",
	})

	var processorCommands []processorCommand
	for _, adapter := range adapters.GetAdapters() {
		if command, err := adapter.Initialize(parser); err != nil {
			fmt.Print(parser.Usage(err))
			return 1
		} else {
			processorCommands = append(processorCommands, processorCommand{
				adapter: adapter,
				command: command,
			})
		}
	}

	for _, outputType := range availableOutputTypes {
		if err := outputType.Initialize(parser); err != nil {
			fmt.Print(parser.Usage(err))
			return 1
		}
	}

	if err := parser.Parse(os.Args); err != nil {
		fmt.Print(parser.Usage(err))
		return 1
	}

	if loggerLevel, err := logrus.ParseLevel(*logLevel); err != nil {
		logrus.Errorf("%s is not a valid log level", *logLevel)
		fmt.Print(parser.Usage(err))
		return 1
	} else {
		logrus.SetLevel(loggerLevel)
	}

	if (os.File{}) == *inputFile {
		inputFile = os.Stdin
	}
	if (os.File{}) == *outputFile {
		outputFile = os.Stdout
	}

	logrus.Debug("Parsing input calendar")

	var inputCalendar ical.Calendar
	dec := ical.NewDecoder(inputFile)
	if cal, err := dec.Decode(); err != nil {
		fmt.Print(parser.Usage(err))
		return 2
	} else {
		inputCalendar = *cal
	}

	outputCalendar := ical.NewCalendar()
	outputCalendar.Props = inputCalendar.Props

	var dStart time.Time
	if dateSelectorStart != nil {
		dStart, _ = time.Parse(time.RFC3339, *dateSelectorStart)
	}

	var dEnd time.Time
	if dateSelectorEnd != nil {
		dEnd, _ = time.Parse(time.RFC3339, *dateSelectorEnd)
	}

	toolbox := processors.Toolbox{
		TextSelectorPattern:    regexp.MustCompile(fmt.Sprintf("(?i)%s", *selector)),
		TextSelectorProps:      *selectorProps,
		DateRangeSelectorStart: dStart,
		DateRangeSelectorEnd:   dEnd,
	}

	for _, processorCommand := range processorCommands {
		if processorCommand.command.Happened() {
			logrus.Infof("Processor %s was selected. Starting process", processorCommand.command.GetName())
			processorCommand.adapter.SetToolbox(toolbox)
			if err := processorCommand.adapter.Process(inputCalendar, outputCalendar); err != nil {
				fmt.Print(parser.Usage(err))
				return 3
			}
		}
	}

	logrus.Infof("Generating output type %s", *outputType)

	if err := availableOutputTypes[*outputType].Generate(outputCalendar, outputFile, outputTypes.OutputOptions{
		Columns: *columns,
	}); err != nil {
		fmt.Print(parser.Usage(err))
		return 4
	}

	return 0
}
