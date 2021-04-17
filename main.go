package main

import (
	"os"
	"os/signal"
	"os/exec"
	"syscall"
	"time"
	"fmt"

	"rosetterc/rosettelog"

	"golang.org/x/sys/unix"
)

func main() {
	if os.Getpid() != 1 {
		rosettelog.Critical("Not run as PID 1... exiting")
		os.Exit(1)
	}

	fmt.Print("\u001b[2J\u001b[0;0H")
	rosettelog.Success("Welcome to.. RosetteRC")

	var procflags uintptr = unix.MS_SYNC | unix.MS_NOSUID | unix.MS_NOEXEC
	syscall.Mount("proc", "/proc", "proc", procflags, "")
	rosettelog.Success("Mounted proc filesystem to /proc")

	var sysflags uintptr = unix.MS_SYNC | unix.MS_NOSUID | unix.MS_NOEXEC
	syscall.Mount("sys", "/sys", "sysfs", sysflags, "")
	rosettelog.Success("Mounted sys filesystem to /sys")

	var runflags uintptr = unix.MS_SYNC | unix.MS_NOSUID
	syscall.Mount("run", "/run", "tmpfs", runflags, "mode=0755")
	rosettelog.Success("Mounted the tmp filesystem run to /run")

	var devflags uintptr = unix.MS_NOSUID
	syscall.Mount("dev", "/dev", "devtmpfs", devflags, "mode=0755")
	rosettelog.Success("Mounted dev filesystem to /dev")

	run("mount", "-o", "remount,rw", "/")
	rosettelog.Success("Remounted / as read-writable")

	run("hostname", "-F", "/etc/hostname")
	rosettelog.Success("Set system hostname")

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)

	go func(){
		for range c {
			rosettelog.Success("Goodbye o/")
			time.Sleep(time.Second * 5)
			syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
		}
	}()

	run("/sbin/getty", "38400", "tty1")
}

func run(cmdstr string, args ...string) {
	cmd := exec.Command(cmdstr, args...)

	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Run()
}
