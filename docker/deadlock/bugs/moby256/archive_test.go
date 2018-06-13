package docker

import (
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"fmt"
)

func TestCmdStreamLargeStderr(t *testing.T) {
	// This test checks for deadlock; thus, the main failure mode of this test is deadlocking.
        // If we change count=63, then this test case could be passed. If count >= 64, it would be failed. 
	cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero bs=1k count=1000 of=/dev/stderr; echo hello")
	out, err := CmdStream(cmd)
	if err != nil {
		t.Fatalf("Failed to start command: " + err.Error())
	}
	fmt.Println("start call io copy in test")
	_, err = io.Copy(ioutil.Discard, out)
	fmt.Println("end call io copy in test")
	if err != nil {
		t.Fatalf("Command should not have failed (err=%s...)", err.Error()[:100])
	}
}
