package converters

import (
	"encoding/csv"
	"fmt"
	"github.com/emersion/go-ical"
	"github.com/google/uuid"
	"golang.org/x/exp/maps"
	"io"
	"slices"
	"time"
)

// CSVConverter converts an incoming CSV formatted file into a Calendar object
type CSVConverter struct {
	Separator       string
	TimestampFormat string
	FieldMap        map[string]string
	HasHeaders      bool
	Headers         []string
}

var _ BaseConverter = &CSVConverter{}

func (c *CSVConverter) Convert(input io.Reader, output *ical.Calendar) error {
	reader := csv.NewReader(input)
	reader.Comma = ([]rune(c.Separator))[0]

	var headers []string

	if c.HasHeaders {
		if h, err := reader.Read(); err != nil {
			return fmt.Errorf("can not read headers from CSV file: %w", err)
		} else {
			headers = h
		}
	} else if len(c.Headers) == 0 {
		return fmt.Errorf("missing headers configuration when no headers exist in the CSV file")
	} else {
		headers = c.Headers
	}

	if rows, err := reader.ReadAll(); err != nil {
		return fmt.Errorf("can not read rows from CSV file: %w", err)
	} else {
		for _, row := range rows {
			fieldValues := make(map[string]string)
			for i, col := range row {
				var field string
				if slices.Contains(maps.Keys(c.FieldMap), headers[i]) {
					field = c.FieldMap[headers[i]]
				} else {
					field = headers[i]
				}
				fieldValues[field] = col
			}
			e := ical.NewEvent()
			e.Props.SetText(ical.PropUID, uuid.NewString())
			if !slices.Contains(maps.Keys(fieldValues), ical.PropDateTimeStamp) {
				e.Props.SetDateTime(ical.PropDateTimeStamp, time.Now())
			}
			for key, value := range fieldValues {
				if slices.Contains([]string{ical.PropDateTimeStart, ical.PropDateTimeEnd, ical.PropDateTimeStamp}, key) {
					if t, err := time.Parse(c.TimestampFormat, value); err != nil {
						return fmt.Errorf("can not convert %s to timestamp using format %s: %w", value, c.TimestampFormat, err)
					} else {
						e.Props.SetDateTime(key, t)
					}
				} else {
					e.Props.SetText(key, value)
				}
			}
			output.Children = append(output.Children, e.Component)
		}
	}

	return nil
}
