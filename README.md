# icarus - iCal batch processor
[![Travis](https://img.shields.io/travis/dploeger/icarus)](https://travis-ci.org/dploeger/icarus)

![icarus](design/logo.png)

## Introduction

icarus does batch processing on iCal/ICS files. It's easy to use and
extend.

Currently, the following processors are available:

* Add an alarm
* Convert all day events to date time events
* Remove events

## Disclaimer

icarus is currently beta software. It will not trash your ics files, but
 you might experience problems when importing the resulting ics data.

## Requirements

icarus is based on Java 8, so a corresponding JRE is needed.

## Usage

icarus is a command line application, that is called with
arguments and the source ics file:

    java -jar icarus.jar [options] icsfile.ics

Use

    java -jar icarus.jar --help

to get detailed help.

Basically, icarus expects a query and/or a time range to search
events.

The query is a regular expression that is matched against the title of
the appointment. The arguments "from" and "until" describe a time range,
that the event has to be in.

The desired processor is configured using the corresponding
argument. For example, the argument "--removeEvent" will remove all events
in the filtered range or with the filtered query.

The resulting ICS data is written to stdout and can be piped to a new
ICS file to be used in your favorite calendar app.