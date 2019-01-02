package de.dieploegers.icarus;

import org.junit.Assert;
import org.junit.Rule;
import org.junit.Test;
import org.junit.contrib.java.lang.system.ExpectedSystemExit;

import java.io.FileInputStream;
import java.io.FileNotFoundException;

/**
 * Main Tester.
 */
public class MainTest {

    @Rule
    public final ExpectedSystemExit exit = ExpectedSystemExit.none();

    @Test()
    public void testCall() {

        this.exit.expectSystemExitWithStatus(0);

        Main.main(new String[]{
            "--query=.*",
            "--alarmBefore=1",
            "src/test/resources/test.ics"
        });

    }

    @Test()
    public void testInvalidFile() {

        this.exit.expectSystemExitWithStatus(1);

        Main.main(new String[]{
            "--query=.*",
            "--alarmBefore=1",
            "src/test/resources/none.ics"
        });

    }

    @Test
    public void testStdin() {

        FileInputStream fileInputStream = null;
        try {
            fileInputStream = new FileInputStream(
                "src/test/resources/test.ics"
            );
        } catch (final FileNotFoundException e) {
            Assert.fail("Test file not found");
        }

        System.setIn(fileInputStream);

        this.exit.expectSystemExitWithStatus(0);

        Main.main(new String[]{
            "--query=.*",
            "--from=20170410000000",
            "--until=20170423235959",
            "--removeEvent"
        });

    }

} 
