package de.dieploegers.icarus;

import net.fortuna.ical4j.data.ParserException;
import org.apache.commons.io.IOUtils;
import org.junit.After;
import org.junit.Assert;
import org.junit.Before;
import org.junit.Test;

import java.io.*;
import java.text.ParseException;
import java.util.regex.Pattern;

/**
 * Processor Tester.
 */
public class ProcessorTest {

    private String iCalData = "";

    @Before
    public void before() throws Exception {
        InputStream stream = new FileInputStream(
            "src/test/resources/test.ics"
        );

        StringWriter writer = new StringWriter();
        IOUtils.copy(stream, writer, "UTF-8");
        this.iCalData = writer.toString();
    }

    @After
    public void after() throws Exception {
    }

    @Test()
    public void testAddAlarm() throws IllegalAccessException, InstantiationException, ClassNotFoundException, ParseException, ParserException, IOException {

        OptionStore options = new OptionStore();

        options.addOption(
            ModifierOption.withValue("alarmBefore", "1")
        );

        Processor processor = new Processor();

        String processed = processor.process(options, this.iCalData);

        Assert.assertTrue(
            "AddAlarm wasn't processed",
            Pattern.compile(
                "TRIGGER:-PT1H"
            ).matcher(processed).find()
        );

        Assert.assertTrue(
            "Alarm message not found",
            Pattern.compile(
                "Alarm provided by icarus"
            ).matcher(processed).find()
        );

    }

    @Test
    public void testAddAlarmMessage() throws IllegalAccessException, InstantiationException, ClassNotFoundException, ParseException, ParserException, IOException {

        OptionStore options = new OptionStore();

        options.addOption(
            ModifierOption.withValue("alarmBefore", "1")
        );

        options.addOption(
            ModifierOption.withValue("alarmMessage", "TESTALARM")
        );

        Processor processor = new Processor();

        String processed = processor.process(options, this.iCalData);

        Assert.assertTrue(
            "Alarm message not set",
            Pattern.compile(
                "TESTALARM"
            ).matcher(processed).find()
        );

    }

    @Test
    public void testAllDayTo() throws IllegalAccessException, InstantiationException, ClassNotFoundException, ParseException, ParserException, IOException {

        OptionStore options = new OptionStore();

        options.addOption(
            ModifierOption.withValue("allDayTo", "19:00-19:05")
        );

        options.addOption(
            ModifierOption.withValue("timezone", "Europe/Berlin")
        );

        Processor processor = new Processor();

        String processed = processor.process(options, this.iCalData);

        Assert.assertTrue(
            "AllDayTo wasn't processed",
            Pattern.compile(
                "DTSTART;TZID=Europe/Berlin:20171227T190000"
            ).matcher(processed).find()
        );

    }

    @Test
    public void testRemoveEvent() throws IllegalAccessException, InstantiationException, ClassNotFoundException, ParseException, ParserException, IOException {

        OptionStore options = new OptionStore();

        options.addOption(
            ModifierOption.withValue("from", "20170410000000")
        );

        options.addOption(
            ModifierOption.withValue("until", "20170423235959")
        );

        options.addOption(
            new ModifierOption("removeEvent")
        );

        Processor processor = new Processor();

        String processed = processor.process(options, this.iCalData);

        Assert.assertFalse(
            "RemoveEvent wasn't processed",
            Pattern.compile(
                "Osterferien 2017 Nordrhein-Westfalen"
            ).matcher(processed).find()
        );

    }

    @Test
    public void testQuery() throws IllegalAccessException, InstantiationException, ClassNotFoundException, ParseException, ParserException, IOException {
        OptionStore options = new OptionStore();

        options.addOption(
            ModifierOption.withValue(
                "query",
                "Osterferien 2017 Nordrhein-Westfalen"
            )
        );

        options.addOption(
            new ModifierOption("removeEvent")
        );

        Processor processor = new Processor();

        String processed = processor.process(options, this.iCalData);

        Assert.assertFalse(
            "RemoveEvent wasn't processed",
            Pattern.compile(
                "Osterferien 2017 Nordrhein-Westfalen"
            ).matcher(processed).find()
        );
    }

}
