package main

import (
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func reapChildren() {
	var (
		status syscall.WaitStatus
		usage  syscall.Rusage
	)
	for {
		if pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, &usage); err != nil {
			log.WithField("Status", status).Warning(err.Error())
			break
		} else {
			if pid == 0 {
				break
			}
			log.WithFields(log.Fields{
				"Pid":    pid,
				"Status": status,
			}).Info("Child Reaped")
		}

	}
}

func kogia_init(c *cli.Context) {

	log.SetOutput(os.Stderr)

	l, err := log.ParseLevel(c.String("level"))
	if err != nil {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(l)
	}

	args := c.Args()
	if len(args) <= 0 {
		log.Error("No command passed to run")
		os.Exit(1)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = os.Environ()

	notify := make(chan os.Signal, 1024)
	signal.Notify(notify, syscall.SIGCHLD, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	go func() {
		for sig := range notify {
			log.Info("Start Reaping")

			switch sig {
			case syscall.SIGHUP:
				cmd.Process.Signal(sig)
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				cmd.Process.Signal(sig)
				break
			case syscall.SIGCHLD:
				reapChildren()
			}
		}
	}()
	ticker := time.NewTicker(time.Second * 5)
	go func() {
		for _ := range ticker.C {
			reapChildren()
		}
	}()

	log.WithFields(log.Fields{
		"Cmd":  args[0],
		"Args": args[1:],
		"Env":  cmd.Env,
	}).Info("Start Main Command")
	err = cmd.Start()
	if err != nil {
		log.Error(err.Error())
		reapChildren()
		os.Exit(1)
	}
	err = cmd.Wait()
	if err != nil {
		log.Error(err.Error())
	}
	reapChildren()
	log.Info("Exit")
}
