package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	// Check first command line argument, then conditionals for run/child
	switch os.Args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		panic("oh no")
	}
}
// create "Command" too run the same executable -> proc self exe is the currently runing program 
func parent() {
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...) // append rest of cli arguments after "child"
	// add  UTS, PID, MNT namespaces
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	// set input, output, and error streams to the streams of parent process
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// try to run via "Run" command
	if err := cmd.Run(); err!= nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}
// parent calls child, creating a command based on subsequent arguments
func child() {
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err!= nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
