package de.dieploegers.icarus.modifier;

import de.dieploegers.icarus.ModifierOption;
import de.dieploegers.icarus.OptionStore;
import de.dieploegers.icarus.exceptions.ProcessException;
import net.fortuna.ical4j.model.Calendar;
import net.fortuna.ical4j.model.component.VEvent;

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

    List<ModifierOption> getOptions();

    /**
     * Process an event
     *
     * @param options  Options set for this modifier
     * @param calendar The calendar in process
     * @param event    The event to process
     * @throws ProcessException Exception processing the events
     */

    void process(OptionStore options, Calendar calendar, VEvent event)
        throws ProcessException;

    /**
     * Finalize (e.g. change the complete calendar) after all events have been
     * processed.
     *
     * The method is called regardless of query, from and until arguments.
     *
     * @param options   Options set for this modifier
     * @param calendar      The calendar in process
     * @param matchedEvents A list of matched events
     * @throws ProcessException Exception processing the events
     */
    void finalize(
        OptionStore options, Calendar calendar, List<VEvent>
        matchedEvents
    ) throws
        ProcessException;

}
