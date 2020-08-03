package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// docker  run <container> cmd args
//go run main.go run cmd args
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("what?")
	}
}

func run() {
	fmt.Printf("running %v as PID %d \n", os.Args[2:], os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("running %v as PID %d \n", os.Args[2:], os.Getpid())

	syscall.Sethostname([]byte("inside-container"))
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// rootfs := "/rootfs-ubuntu"
	// must(syscall.Chroot(rootfs))
	// must(os.Chdir("/"))
	syscall.Chroot("/home/rootfs-ubuntu")
	os.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")
	must(cmd.Run())
}

func must(err error) {
	if err != nil {

		panic(err)
	}
}
