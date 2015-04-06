package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"regexp"
	"runtime"
	"strings"
	"syscall"

	log "github.com/Sirupsen/logrus"
)

var (
	skip_preinit  = flag.Bool("skip-preinit", false, "Do not execute preinit scripts")
	preinit_dir   = flag.String("preinit", "/etc/preinit.d", "Path to preinit scripts")
	skip_postinit = flag.Bool("skip-postinit", false, "Do not execute postinit scripts")
	postinit_dir  = flag.String("postinit", "/etc/postinit.d", "Path to postinit scripts")
	skip_env      = flag.Bool("skip-env", false, "Do not load additional env from files")
	env_file      = flag.String("env", "/etc/env", "Path to additional env files")
	verbose       = flag.String("verbose", "warning", "Verbosity level")
)

func LoadEnv() []string {
	env := os.Environ()
	if !*skip_env {
		if buff, err := ioutil.ReadFile(*env_file); err == nil {
			for _, line := range strings.Split(strings.TrimSpace(string(buff)), "\n") {
				if ok, err := regexp.MatchString(`^\w+=.*$`, line); ok {
					env = append(env, line)
					continue
				} else if !ok {
					log.WithField("line", line).Warning("Line is not match")
				} else if err != nil {
					log.Warning(err.Error())
				}
			}
		} else {
			log.WithField("EnvFile", *env_file).Warning(err.Error())
		}
	}
	return env
}

func StartAndWait(command string, args []string) {
	cmd := exec.Command(command, args...)
	cmd.Env = LoadEnv()
	out, err := cmd.CombinedOutput()
	fields := log.Fields{
		"Cmd":  args[1],
		"Args": args[2:],
		"Env":  cmd.Env,
	}
	if err != nil {
		log.WithFields(fields).Warning(err.Error())
	} else {
		fields["Output"] = string(out)
		log.WithFields(fields).Info("Run Command Result")
	}
}

func RunDir(dir string) error {
	listdir, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, ent := range listdir {
		if !ent.IsDir() {
			StartAndWait("/bin/sh", []string{"-c", path.Join(dir, ent.Name())})
		}
	}

	return nil
}

func RunPreInit() {
	if !*skip_preinit {
		log.Info("Running PreInit")
		if err := RunDir(*preinit_dir); err != nil {
			log.WithFields(log.Fields{
				"Task": "RunPreInit",
				"Dir":  *preinit_dir,
			}).Warning(err.Error())
		}
	}
}
func RunPostInit() {
	if !*skip_postinit {
		log.Info("Running PostInit")
		if err := RunDir(*postinit_dir); err != nil {
			log.WithFields(log.Fields{
				"Task": "RunPostInit",
				"Dir":  *postinit_dir,
			}).Warning(err.Error())
		}
	}
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	log.SetOutput(os.Stderr)
	l, err := log.ParseLevel(*verbose)
	if err != nil {
		log.SetLevel(log.WarnLevel)
	} else {
		log.SetLevel(l)
	}
}

func ReapChildren() {
	var (
		status syscall.WaitStatus
		usage  syscall.Rusage
	)
	for {
		if pid, err := syscall.Wait4(-1, &status, syscall.WNOHANG, &usage); err != nil {
			log.WithField("status", status).Warning(err.Error())
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

func main() {

	if !*skip_preinit {
		RunPreInit()
	}

	notify := make(chan os.Signal, 1024)
	signal.Notify(notify, syscall.SIGCHLD, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	args := flag.Args()
	cmd := exec.Command(args[0], args[1:]...)

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
				ReapChildren()
			}
		}
	}()

	if !*skip_env {
		cmd.Env = LoadEnv()
	}
	log.WithFields(log.Fields{
		"Cmd":  args[0],
		"Args": args[1:],
		"Env":  cmd.Env,
	}).Info("Start Main Command")
	cmd.Start()
	if !*skip_postinit {
		RunPostInit()
	}
	cmd.Wait()
	log.Info("Stopping")
}
