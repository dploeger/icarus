package de.dieploegers.icarus.modifier;

import de.dieploegers.icarus.exceptions.ProcessException;
import net.fortuna.ical4j.model.*;
import net.fortuna.ical4j.model.component.VEvent;
import net.fortuna.ical4j.model.property.DtEnd;
import net.fortuna.ical4j.model.property.DtStart;
import net.fortuna.ical4j.util.Dates;
import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.Option;

import java.util.ArrayList;
import java.util.List;

/**
 * Rework all day events to non-all day events
 */
public class AllDayToModifier implements Modifier {
    @Override
    public List<Option> getOptions() {
        List<Option> options = new ArrayList<>();
        options.add(
            Option.builder()
                .longOpt("allDayTo")
                .hasArg()
                .desc("Change an all day event to this time range (hh:mm-hh:mm)")
                .build()
        );

        options.add(
            Option.builder()
                .longOpt("timezone")
                .hasArg()
                .desc("Timezone to use for the time. (e.g. Europe/Berlin)")
                .build()
        );
        return options;
    }

    @Override
    public void process(
        CommandLine commandLine, Calendar calendar, VEvent event
    ) throws ProcessException {
        if (
            commandLine.hasOption("allDayTo") &&
                event.getStartDate().toString().contains("VALUE=DATE:")
            ) {

            TimeZone timezone = null;

            if (commandLine.hasOption("timezone")) {
                TimeZoneRegistry registry = TimeZoneRegistryFactory.getInstance().createRegistry();
                timezone = registry.getTimeZone(
                    commandLine.getOptionValue("timezone")
                );
            }

            String[] times = commandLine.getOptionValue(
                "allDayTo"
            ).split("-");

            String[] startTimeParts = times[0].split(":");
            String[] endTimeParts = times[1].split(":");

            java.util.Calendar startCalendar = Dates.getCalendarInstance(
                event.getStartDate().getDate()
            );

            if (event.getStartDate().getTimeZone() != null) {
                startCalendar.setTimeZone(event.getStartDate().getTimeZone());
            }

            startCalendar.setTime(event.getStartDate().getDate());

            java.util.Calendar endCalendar = Dates.getCalendarInstance(
                event.getStartDate().getDate()
            );

            endCalendar.setTime(event.getStartDate().getDate());

            if (event.getEndDate().getTimeZone() != null) {
                endCalendar.setTimeZone(event.getEndDate().getTimeZone());
            }

            startCalendar.set(
                java.util.Calendar.HOUR,
                Integer.valueOf(startTimeParts[0])
            );
            startCalendar.set(
                java.util.Calendar.MINUTE,
                Integer.valueOf(startTimeParts[1])
            );

            endCalendar.set(
                java.util.Calendar.HOUR,
                Integer.valueOf(endTimeParts[0])
            );

            endCalendar.set(
                java.util.Calendar.MINUTE,
                Integer.valueOf(endTimeParts[1])
            );

            DateTime startDateTime = new DateTime(startCalendar.getTime());
            if (timezone != null) {
                startDateTime.setTimeZone(timezone);
            }

            Property dtStartProperty = event.getProperty(Property.DTSTART);
            if (dtStartProperty != null) {
                event.getProperties().remove(dtStartProperty);
                event.getProperties().add(new DtStart(startDateTime));
            }

            DateTime endDateTime = new DateTime(endCalendar.getTime());
            if (timezone != null) {
                endDateTime.setTimeZone(timezone);
            }

            Property dtEndProperty = event.getProperty(Property.DTEND);
            if (dtEndProperty != null) {
                event.getProperties().remove(dtEndProperty);
                event.getProperties().add(new DtEnd(endDateTime));
            }

        }
    }

    @Override
    public void finalize(
        CommandLine commandLine, Calendar calendar, List<VEvent> matchedEvents
    ) throws ProcessException {
        // Nothing to do
    }
}
