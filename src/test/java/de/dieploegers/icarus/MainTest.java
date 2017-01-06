package de.dieploegers.icarus;

import org.junit.*;
import org.junit.contrib.java.lang.system.ExpectedSystemExit;

import java.io.ByteArrayOutputStream;
import java.io.PrintStream;
import java.util.regex.Matcher;
import java.util.regex.Pattern;

/**
 * Main Tester.
 *
 * @author <Authors name>
 * @version 1.0
 * @since <pre>Jan 5, 2017</pre>
 */
public class MainTest {

    @Rule
    public final ExpectedSystemExit exit = ExpectedSystemExit.none();

    @Before
    public void before() throws Exception {
    }

    @After
    public void after() throws Exception {
    }

    /**
     * Method: main(String[] args)
     */
    @Test
    public void testAddAlarm() throws Exception {
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        PrintStream ps = new PrintStream(baos);

        System.setOut(ps);

        exit.expectSystemExitWithStatus(0);

        Main.main(new String[]{
            "--query=.*",
            "--alarmBefore=1",
            "src/test/resources/test.ics"
        });

        String result = baos.toString();

        Matcher addAlarmMatcher = Pattern.compile(
            "TRIGGER:PT1H"
        ).matcher(result);

        Assert.assertTrue(
            "AddAlarm wasn't processed",
            addAlarmMatcher.find()
        );

    }

    /**
     * Method: main(String[] args)
     */
    @Test
    public void testAllDayTo() throws Exception {
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        PrintStream ps = new PrintStream(baos);

        System.setOut(ps);

        exit.expectSystemExitWithStatus(0);

        Main.main(new String[]{
            "--query=.*",
            "--allDayTo=19:00-19:05",
            "--timezone=Europe/Berlin",
            "src/test/resources/test.ics"
        });

        String result = baos.toString();

        Matcher allDayMatcher = Pattern.compile(
            "DTSTART;TZID=Europe/Berlin:20171227T200000"
        ).matcher(result);

        Assert.assertTrue(
            "AllDayTo wasn't processed",
            allDayMatcher.find()
        );

    }

    /**
     * Method: main(String[] args)
     */
    @Test
    public void testRemoveEvent() throws Exception {
        ByteArrayOutputStream baos = new ByteArrayOutputStream();
        PrintStream ps = new PrintStream(baos);

        System.setOut(ps);

        exit.expectSystemExitWithStatus(0);

        Main.main(new String[]{
            "--query=.*",
            "--from=20170410000000",
            "--until=20170423235959",
            "--removeEvent",
            "src/test/resources/test.ics"
        });

        String result = baos.toString();

        Matcher allDayMatcher = Pattern.compile(
            "Osterferien 2017 Nordrhein-Westfalen"
        ).matcher(result);

        Assert.assertFalse(
            "RemoveEvent wasn't processed",
            allDayMatcher.find()
        );

    }

} 
