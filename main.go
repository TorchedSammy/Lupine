package main

import (
	"os"
	"os/signal"
	"os/exec"
	"syscall"
	"time"
	"fmt"

	"github.com/TorchedSammy/Lupine/lupinelog"

	"golang.org/x/sys/unix"
	"github.com/yuin/gopher-lua"
)

func main() {
	if os.Getpid() != 1 {
		lupinelog.Critical("Not run as PID 1... exiting")
		os.Exit(1)
	}

	fmt.Print("\u001b[2J\u001b[0;0H")
	lupinelog.Success("Welcome to Lupine!")

	l := lua.NewState()
	l.OpenLibs()
	l.PreloadModule("lupinelog", lupinelog.Loader)

	var procflags uintptr = unix.MS_SYNC | unix.MS_NOSUID | unix.MS_NOEXEC
	syscall.Mount("proc", "/proc", "proc", procflags, "")
	lupinelog.Success("Mounted proc filesystem to /proc")

	var sysflags uintptr = unix.MS_SYNC | unix.MS_NOSUID | unix.MS_NOEXEC
	syscall.Mount("sys", "/sys", "sysfs", sysflags, "")
	lupinelog.Success("Mounted sys filesystem to /sys")

	var runflags uintptr = unix.MS_SYNC | unix.MS_NOSUID
	syscall.Mount("run", "/run", "tmpfs", runflags, "mode=0755")
	lupinelog.Success("Mounted the tmp filesystem run to /run")

	var devflags uintptr = unix.MS_NOSUID
	syscall.Mount("dev", "/dev", "devtmpfs", devflags, "mode=0755")
	lupinelog.Success("Mounted dev filesystem to /dev")

	run("mount", "-o", "remount,rw", "/")
	lupinelog.Success("Remounted / as read-writable")

	files, err := os.ReadDir("/etc/rosetterc/services")
    if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

    for _, f := range files {
		err := l.DoFile("/etc/rosetterc/services/" + f.Name())
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)

	go func(){
		for range c {
			lupinelog.Success("Goodbye o/")
			time.Sleep(time.Second * 5)
			syscall.Sync()
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
