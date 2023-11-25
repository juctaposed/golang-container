package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall" // https://pkg.go.dev/syscall#pkg-constants
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
	// added  UTS, PID, MNT namespaces ... add IPC/user/net later
	cmd.SysProcAttr = &syscall.SysProcAttr{ // system attributes for new process
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
	// swap to root fs by performing a bind mount - attach rootfs directory to itself, making it available in more than one location
	must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
		must(os.MkdirAll("rootfs/oldrootfs", 0700)) // new directory "oldrootfs" w/ permission to read/write/execute only for owner 
		// move current directory at "/" to "rootfs/oldrootfs". Pivotroot will swap "rootfs" to "/" -> OS requires for bind mount call that swapping FS are not part of same tree , we can mount to self
		must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
		must(od.Chdir("/")) // change pwd to new root fs
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
