package de.dieploegers.icarus;

import de.dieploegers.icarus.exceptions.ProcessException;
import de.dieploegers.icarus.modifier.Modifier;
import net.fortuna.ical4j.data.CalendarBuilder;
import net.fortuna.ical4j.data.ParserException;
import net.fortuna.ical4j.model.Calendar;
import net.fortuna.ical4j.model.Property;
import net.fortuna.ical4j.model.component.CalendarComponent;
import net.fortuna.ical4j.model.component.VEvent;
import org.reflections.Reflections;

import java.io.IOException;
import java.io.StringReader;
import java.text.ParseException;
import java.text.SimpleDateFormat;
import java.util.ArrayList;
import java.util.Date;
import java.util.List;
import java.util.regex.Pattern;

/**
 * The icarus processor
 */
public class Processor {

    private final List<Modifier> modifiers;

    /**
     * Creates the processor and scans available modifiers
     *
     * @throws ClassNotFoundException Error instantiating modifier
     * @throws IllegalAccessException Error instantiating modifier
     * @throws InstantiationException Error instantiating modifier
     */

    public Processor(
    ) throws ClassNotFoundException, IllegalAccessException, InstantiationException {

        this.modifiers = new ArrayList<>();

        final Reflections reflections = new Reflections(
            "de.dieploegers.icarus.modifier"
        );

        for (
            final Class<? extends Modifier> modifier :
            reflections.getSubTypesOf(Modifier.class)
            ) {
            this.modifiers.add(
                (Modifier) Class.forName(modifier.getName()).newInstance()
            );
        }
    }

    /**
     * Process the ICAL data
     *
     * @param options  All available options
     * @param iCalData The raw ical data
     * @return The processed ical data
     * @throws IOException     Error reading calendar
     * @throws ParserException Error parsing calendar
     * @throws ParseException  Error parsing from or until date
     */

    public String process(final OptionStore options, final String iCalData) throws
        IOException,
        ParserException, ParseException {
        final CalendarBuilder builder = new CalendarBuilder();
        final StringReader stringReader = new StringReader(iCalData);
        final Calendar calendar = builder.build(stringReader);

        final Pattern queryPattern;
        if (options.isSet("query")) {
            queryPattern = Pattern.compile(
                options.get("query")
            );
        } else {
            queryPattern = Pattern.compile(".*");
        }

        Date dateFrom = null;

        if (options.isSet("from")) {
            dateFrom = new SimpleDateFormat("yyyyMMddHHmmss").parse(
                options.get("from")
            );
        }

        Date dateUntil = null;

        if (options.isSet("until")) {
            dateUntil = new SimpleDateFormat("yyyyMMddHHmmss").parse(
                options.get("until")
            );
        }

        final List<VEvent> matchedEvents = new ArrayList<>();

        for (final CalendarComponent component : calendar.getComponents()) {
            if (component.getClass().getSimpleName().equals("VEvent")) {
                final VEvent event = (VEvent) component;

                if (queryPattern.matcher(
                    event.getProperty(Property.SUMMARY).getValue()
                ).find()) {

                    // Skip events not in the requested timeline

                    boolean dateFiltered = false;

                    if (dateFrom == null && dateUntil == null) {
                        dateFiltered = true;
                    } else if (dateFrom != null && dateUntil == null) {
                        if (event.getStartDate().getDate().after(dateFrom)) {
                            dateFiltered = true;
                        }
                    } else if (dateFrom == null) {
                        if (event.getEndDate().getDate().before(dateUntil)) {
                            dateFiltered = true;
                        }
                    } else {
                        if (
                            event.getStartDate().getDate().after(dateFrom) &&
                                event.getEndDate().getDate().before(dateUntil)
                            ) {
                            dateFiltered = true;
                        }
                    }

                    if (dateFiltered) {

                        matchedEvents.add(event);

                        for (final Modifier modifier : this.modifiers) {
                            try {
                                modifier.process(
                                    options,
                                    calendar,
                                    event
                                );
                            } catch (final ProcessException e) {
                                System.out.println(
                                    "Error modifying event: " + e.toString()
                                );
                                System.exit(1);
                            }
                        }

                    }

                }
            }
        }

        for (final Modifier modifier : this.modifiers) {
            try {
                modifier.finalize(
                    options,
                    calendar,
                    matchedEvents
                );
            } catch (final ProcessException e) {
                System.out.println(
                    "Error finalizing events: " + e.toString()
                );
                System.exit(1);
            }
        }

        return calendar.toString();

    }

    /**
     * Retrieve scanned modifiers
     *
     * @return A list of found modifiers
     */

    public List<Modifier> getModifiers() {
        return this.modifiers;
    }
}
