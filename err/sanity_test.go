package err

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

/*
 * Since CheckIfError uses os.Exit we need to go about testing this in a round-about way.
 * This test actually calls go itself to run itself with the TEST_MODE environment variable
 * set. If the called execution _does not_ fail then the test is a failure since it is expected
 * to fail.
 *
 * from: https://stackoverflow.com/questions/26225513/how-to-test-os-exit-scenarios-in-go
 */
func TestCheckIfError(t *testing.T)  {
	// do this on the "internal" test (TEST_MODE == 1)
	if os.Getenv("TEST_MODE") == "1" {
		CheckIfError(fmt.Errorf("This is an error condition"))
		return
	}

	// if not on the internal portion execute go to run the test
	cmd := exec.Command(os.Args[0], "-test.run=TestCheckIfError")
	cmd.Env = append(os.Environ(), "TEST_MODE=1")
	err := cmd.Run()
	// if the function does an os.Exit() > 0 then return (PASS)
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	// if the function does not have an error code (os.Exit() == 0) then the test is a failure
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// TestCheckIfErrorWithNoError is a simple test that ensures that no os.Exit is called when an error happens
func TestCheckIfErrorWithNoError(t *testing.T) {
	CheckIfError(nil)
}