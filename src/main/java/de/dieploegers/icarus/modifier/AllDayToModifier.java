package de.dieploegers.icarus.modifier;

import de.dieploegers.icarus.ModifierOption;
import de.dieploegers.icarus.OptionStore;
import de.dieploegers.icarus.exceptions.ProcessException;
import net.fortuna.ical4j.model.*;
import net.fortuna.ical4j.model.Calendar;
import net.fortuna.ical4j.model.TimeZone;
import net.fortuna.ical4j.model.component.VEvent;
import net.fortuna.ical4j.model.property.DtEnd;
import net.fortuna.ical4j.model.property.DtStart;
import net.fortuna.ical4j.util.Dates;

import java.util.*;

/**
 * Rework all day events to non-all day events
 */
public class AllDayToModifier implements Modifier {
    @Override
    public List<ModifierOption> getOptions() {
        List<ModifierOption> options = new ArrayList<>();
        options.add(
            new ModifierOption(
                "allDayTo",
                "Change an all day event to this time range (hh:mm-hh:mm)",
                true
            )
        );

        options.add(
            new ModifierOption(
                "timezone",
                "Timezone to use for the time. (e.g. Europe/Berlin)",
                true
            )
        );
        return options;
    }

    @Override
    public void process(
        OptionStore options, Calendar calendar, VEvent event
    ) throws ProcessException {
        if (
            options.isSet("allDayTo") &&
                event.getStartDate().toString().contains("VALUE=DATE:")
            ) {

            TimeZone timezone = null;

            if (options.isSet("timezone")) {
                TimeZoneRegistry registry = TimeZoneRegistryFactory.getInstance().createRegistry();
                timezone = registry.getTimeZone(
                    options.get("timezone")
                );
            }

            String[] times = options.get(
                "allDayTo"
            ).split("-");

            String[] startTimeParts = times[0].split(":");
            String[] endTimeParts = times[1].split(":");

            java.util.Calendar startCalendar = Dates.getCalendarInstance(
                event.getStartDate().getDate()
            );

            if (timezone != null) {
                startCalendar.setTimeZone(timezone);
            }

            if (event.getStartDate().getTimeZone() != null) {
                startCalendar.setTimeZone(event.getStartDate().getTimeZone());
            }

            startCalendar.setTime(event.getStartDate().getDate());

            java.util.Calendar endCalendar = Dates.getCalendarInstance(
                event.getStartDate().getDate()
            );

            if (timezone != null) {
                endCalendar.setTimeZone(timezone);
            }

            endCalendar.setTime(event.getStartDate().getDate());

            if (event.getEndDate().getTimeZone() != null) {
                endCalendar.setTimeZone(event.getEndDate().getTimeZone());
            }

            startCalendar.set(
                java.util.Calendar.HOUR_OF_DAY,
                Integer.valueOf(startTimeParts[0])
            );
            startCalendar.set(
                java.util.Calendar.MINUTE,
                Integer.valueOf(startTimeParts[1])
            );

            endCalendar.set(
                java.util.Calendar.HOUR_OF_DAY,
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
        OptionStore options, Calendar calendar, List<VEvent> matchedEvents
    ) throws ProcessException {

    }
}
