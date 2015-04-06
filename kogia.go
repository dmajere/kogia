package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
)

func loadEnv(env_file string) ([]string, error) {
	var env []string
	if env_file == "" {
		return []string{}, nil
	}

	if buff, err := ioutil.ReadFile(env_file); err == nil {
		for _, line := range strings.Split(strings.TrimSpace(string(buff)), "\n") {
			if ok, err := regexp.MatchString(`^\w+=.*$`, line); ok {
				env = append(env, line)
				continue
			} else if !ok {
				log.WithField("line", line).Warning("Line is not match")
			} else if err != nil {
				return []string{}, err
			}
		}
	} else {
		return env, err
	}
	return env, nil
}

func startAndWait(command string, env_file string) ([]byte, error) {
	cmd := exec.Command(command)
	cmd.Env = os.Environ()

	if data, err := loadEnv(env_file); err != nil {
		log.WithField("EnvFile", env_file).Warning(err.Error())
	} else {
		cmd.Env = append(cmd.Env, data...)
	}

	out, err := cmd.CombinedOutput()
	return out, err
}

func runDir(dir, env_file string) error {
	listdir, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, ent := range listdir {
		if !ent.IsDir() {
			cmd := path.Join(dir, ent.Name())
			out, err := startAndWait(cmd, env_file)
			fields := log.Fields{
				"Cmd": cmd,
				"Env": env_file,
			}
			if err != nil {
				log.WithFields(fields).Warning(err.Error())
			} else {
				log.WithFields(fields).Infof("Run Command Result %s", string(out))
			}

		}
	}

	return nil
}

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

	var env_file string

	if c.Bool("skip-env") {
		env_file = ""
	} else {
		env_file = c.String("env")
		if env, err := loadEnv(env_file); err == nil {
			cmd.Env = append(cmd.Env, env...)
		}
	}

	if !c.Bool("skip-preinit") {
		log.Info("Running PreInit")
		if err := runDir(c.String("preinit"), env_file); err != nil {
			log.WithFields(log.Fields{
				"Task": "RunPreInit",
				"Dir":  c.String("preinit"),
			}).Warning(err.Error())
		}
	}

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

	log.WithFields(log.Fields{
		"Cmd":  args[0],
		"Args": args[1:],
		"Env":  cmd.Env,
	}).Info("Start Main Command")
	err = cmd.Start()
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}

	if !c.Bool("skip-postinit") {
		log.Info("Running PostInit")
		if err := runDir(c.String("postinit"), env_file); err != nil {
			log.WithFields(log.Fields{
				"Task": "RunPostInit",
				"Dir":  c.String("postinit"),
			}).Warning(err.Error())
		}
	}

	err = cmd.Wait()
	if err != nil {
		log.Error(err.Error())
	}
	log.Info("Exit")
}
