package converters

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/emersion/go-ical"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestConvert(t *testing.T) {
	startTime, _ := time.Parse(time.RFC3339, "2024-08-26T12:24:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2024-08-26T13:24:00Z")

	subject := heredoc.Docf(`
		DTSTART,DTEND,SUMMARY
		%s,%s,Test
	`,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
	)

	converter := CSVConverter{
		Separator:       ",",
		TimestampFormat: time.RFC3339,
		FieldMap:        nil,
		HasHeaders:      true,
		Headers:         nil,
		Location:        time.UTC,
	}

	calendar := ical.NewCalendar()

	err := converter.Convert(strings.NewReader(subject), calendar)
	assert.NoError(t, err, "Converter errored out")
	assert.Len(t, calendar.Children, 1, "Invalid number of children in calendar")
	event := calendar.Children[0]
	summary, _ := event.Props.Get(ical.PropSummary).Text()
	assert.Equal(t, "Test", summary)
	dtStart, _ := event.Props.Get(ical.PropDateTimeStart).DateTime(startTime.Location())
	assert.Equal(t, startTime, dtStart, "Wrong start time")
	dtEnd, _ := event.Props.Get(ical.PropDateTimeEnd).DateTime(endTime.Location())
	assert.Equal(t, endTime, dtEnd, "Wrong end time")
}

func TestConvertSeparater(t *testing.T) {
	startTime, _ := time.Parse(time.RFC3339, "2024-08-26T12:24:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2024-08-26T13:24:00Z")

	subject := heredoc.Docf(`
		DTSTART;DTEND;SUMMARY
		%s;%s;Test
	`,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
	)

	converter := CSVConverter{
		Separator:       ";",
		TimestampFormat: time.RFC3339,
		FieldMap:        nil,
		HasHeaders:      true,
		Headers:         nil,
		Location:        time.UTC,
	}

	calendar := ical.NewCalendar()

	err := converter.Convert(strings.NewReader(subject), calendar)
	assert.NoError(t, err, "Converter errored out")
	assert.Len(t, calendar.Children, 1, "Invalid number of children in calendar")
	event := calendar.Children[0]
	summary, _ := event.Props.Get(ical.PropSummary).Text()
	assert.Equal(t, "Test", summary)
	dtStart, _ := event.Props.Get(ical.PropDateTimeStart).DateTime(startTime.Location())
	assert.Equal(t, startTime, dtStart, "Wrong start time")
	dtEnd, _ := event.Props.Get(ical.PropDateTimeEnd).DateTime(endTime.Location())
	assert.Equal(t, endTime, dtEnd, "Wrong end time")
}

func TestConvertTimestampFormat(t *testing.T) {
	startTime, _ := time.Parse("2006-01-02 15:04", "2024-08-26 12:24")
	endTime, _ := time.Parse("2006-01-02 15:04", "2024-08-26 13:24")

	subject := heredoc.Docf(`
		DTSTART,DTEND,SUMMARY
		%s,%s,Test
	`,
		startTime.Format("2006-01-02 15:04"),
		endTime.Format("2006-01-02 15:04"),
	)

	converter := CSVConverter{
		Separator:       ",",
		TimestampFormat: "2006-01-02 15:04",
		FieldMap:        nil,
		HasHeaders:      true,
		Headers:         nil,
		Location:        time.UTC,
	}

	calendar := ical.NewCalendar()

	err := converter.Convert(strings.NewReader(subject), calendar)
	assert.NoError(t, err, "Converter errored out")
	assert.Len(t, calendar.Children, 1, "Invalid number of children in calendar")
	event := calendar.Children[0]
	summary, _ := event.Props.Get(ical.PropSummary).Text()
	assert.Equal(t, "Test", summary)
	dtStart, _ := event.Props.Get(ical.PropDateTimeStart).DateTime(startTime.Location())
	assert.Equal(t, startTime, dtStart, "Wrong start time")
	dtEnd, _ := event.Props.Get(ical.PropDateTimeEnd).DateTime(endTime.Location())
	assert.Equal(t, endTime, dtEnd, "Wrong end time")
}

func TestConvertFieldMap(t *testing.T) {
	startTime, _ := time.Parse(time.RFC3339, "2024-08-26T12:24:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2024-08-26T13:24:00Z")

	subject := heredoc.Docf(`
		Start,End,Text
		%s,%s,Test
	`,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
	)

	converter := CSVConverter{
		Separator:       ",",
		TimestampFormat: time.RFC3339,
		FieldMap: map[string]string{
			"Start": ical.PropDateTimeStart,
			"End":   ical.PropDateTimeEnd,
			"Text":  ical.PropSummary,
		},
		HasHeaders: true,
		Headers:    nil,
		Location:   time.UTC,
	}

	calendar := ical.NewCalendar()

	err := converter.Convert(strings.NewReader(subject), calendar)
	assert.NoError(t, err, "Converter errored out")
	assert.Len(t, calendar.Children, 1, "Invalid number of children in calendar")
	event := calendar.Children[0]
	summary, _ := event.Props.Get(ical.PropSummary).Text()
	assert.Equal(t, "Test", summary)
	dtStart, _ := event.Props.Get(ical.PropDateTimeStart).DateTime(startTime.Location())
	assert.Equal(t, startTime, dtStart, "Wrong start time")
	dtEnd, _ := event.Props.Get(ical.PropDateTimeEnd).DateTime(endTime.Location())
	assert.Equal(t, endTime, dtEnd, "Wrong end time")
}

func TestConvertNoHeaders(t *testing.T) {
	startTime, _ := time.Parse(time.RFC3339, "2024-08-26T12:24:00Z")
	endTime, _ := time.Parse(time.RFC3339, "2024-08-26T13:24:00Z")

	subject := heredoc.Docf(`
		%s,%s,Test
	`,
		startTime.Format(time.RFC3339),
		endTime.Format(time.RFC3339),
	)

	converter := CSVConverter{
		Separator:       ",",
		TimestampFormat: time.RFC3339,
		FieldMap:        nil,
		HasHeaders:      false,
		Headers:         []string{"DTSTART", "DTEND", "SUMMARY"},
		Location:        time.UTC,
	}

	calendar := ical.NewCalendar()

	err := converter.Convert(strings.NewReader(subject), calendar)
	assert.NoError(t, err, "Converter errored out")
	assert.Len(t, calendar.Children, 1, "Invalid number of children in calendar")
	event := calendar.Children[0]
	summary, _ := event.Props.Get(ical.PropSummary).Text()
	assert.Equal(t, "Test", summary)
	dtStart, _ := event.Props.Get(ical.PropDateTimeStart).DateTime(startTime.Location())
	assert.Equal(t, startTime, dtStart, "Wrong start time")
	dtEnd, _ := event.Props.Get(ical.PropDateTimeEnd).DateTime(endTime.Location())
	assert.Equal(t, endTime, dtEnd, "Wrong end time")
}

func TestConvertLocation(t *testing.T) {

	subject := heredoc.Doc(`
		DTSTART,DTEND,SUMMARY
		2024-08-26 12:24,2024-08-26 13:24,Test
	`)

	cetLocation, _ := time.LoadLocation("CET")

	converter := CSVConverter{
		Separator:       ",",
		TimestampFormat: "2006-01-02 15:04",
		FieldMap:        nil,
		HasHeaders:      true,
		Headers:         nil,
		Location:        cetLocation,
	}

	calendar := ical.NewCalendar()

	err := converter.Convert(strings.NewReader(subject), calendar)
	assert.NoError(t, err, "Converter errored out")
	assert.Len(t, calendar.Children, 1, "Invalid number of children in calendar")
	event := calendar.Children[0]
	summary, _ := event.Props.Get(ical.PropSummary).Text()
	assert.Equal(t, "Test", summary)
	dtStart, _ := event.Props.Get(ical.PropDateTimeStart).DateTime(cetLocation)
	assert.Equal(t, "2024-08-26 12:24", dtStart.Format("2006-01-02 15:04"), "Wrong start time")
	dtEnd, _ := event.Props.Get(ical.PropDateTimeEnd).DateTime(cetLocation)
	assert.Equal(t, "2024-08-26 13:24", dtEnd.Format("2006-01-02 15:04"), "Wrong end time")
}
