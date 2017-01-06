package de.dieploegers.icarus.modifier;

import de.dieploegers.icarus.exceptions.ProcessException;
import net.fortuna.ical4j.model.Calendar;
import net.fortuna.ical4j.model.component.VEvent;
import org.apache.commons.cli.CommandLine;
import org.apache.commons.cli.Option;

import java.util.List;

/**
 * Modifier interface
 */
public interface Modifier {

    /**
     * Return options to be added as command line arguments
     *
     * @return A list of options to add
     */

    List<Option> getOptions();

    /**
     * Process an event
     *
     * @param commandLine The commandline instance holding program arguments
     * @param calendar    The calendar in process
     * @param event       The event to process
     * @throws ProcessException Exception processing the events
     */

    void process(CommandLine commandLine, Calendar calendar, VEvent event)
        throws ProcessException;

    /**
     * Finalize (e.g. change the complete calendar) after all events have been
     * processed.
     *
     * The method is called regardless of query, from and until arguments.
     *
     * @param commandLine The commandline instance holding program arguments
     * @param calendar The calendar in process
     * @param matchedEvents A list of matched events
     * @throws ProcessException Exception processing the events
     */
    void finalize(
        CommandLine commandLine, Calendar calendar, List<VEvent>
        matchedEvents
    ) throws
        ProcessException;

}
