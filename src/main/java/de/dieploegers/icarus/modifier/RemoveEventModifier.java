package de.dieploegers.icarus.modifier;

import de.dieploegers.icarus.exceptions.ProcessException;
import net.fortuna.ical4j.model.Calendar;
import net.fortuna.ical4j.model.component.VEvent;
import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.Option;

import java.util.ArrayList;
import java.util.List;

/**
 * Removes an Event
 */
public class RemoveEventModifier implements Modifier {

    @Override
    public List<Option> getOptions() {
        List<Option> options = new ArrayList<>();
        options.add(
            Option.builder()
                .longOpt("removeEvent")
                .build()
        );
        return options;
    }

    @Override
    public void process(
        CommandLine commandLine, Calendar calendar, VEvent event
    ) throws ProcessException {
        // nothing to do
    }

    @Override
    public void finalize(
        CommandLine commandLine, Calendar calendar, List<VEvent> matchedEvents
    ) throws ProcessException {
        if (commandLine.hasOption("removeEvent")) {
            calendar.getComponents().removeAll(matchedEvents);
        }
    }
}
