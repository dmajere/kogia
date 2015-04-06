package main

import (
	"os"
	"runtime"

	"github.com/codegangsta/cli"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	app := cli.NewApp()
	app.Name = "kogia"
	app.Usage = "Small and simple init system for docker containers"
	app.Version = "0.11"
	app.Author = "Dmitry Fedorov"
	app.Email = "gmajere@gmail.com"

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "skip-preinit, S",
			Usage:  "Do not execute preinit scripts",
			EnvVar: "KOGIA_SKIP_PREINIT",
		},
		cli.StringFlag{
			Name:   "preinit, s",
			Value:  "/etc/preinit.d",
			Usage:  "Path to preinit scripts",
			EnvVar: "KOGIA_PREINIT_DIR",
		},
		cli.BoolFlag{
			Name:   "skip-postinit, P",
			Usage:  "Do not execute postinit scripts",
			EnvVar: "KOGIA_SKIP_POSTINIT",
		},
		cli.StringFlag{
			Name:   "postinit, p",
			Value:  "/etc/postinit.d",
			Usage:  "Path to postinit scripts",
			EnvVar: "KOGIA_POSTINIT_DIR",
		},
		cli.BoolFlag{
			Name:   "skip-env, E",
			Usage:  "Do not load additional env from files",
			EnvVar: "KOGIA_SKIP_ENV",
		},
		cli.StringFlag{
			Name:   "env, e",
			Value:  "/etc/env",
			Usage:  "Path to additional env file",
			EnvVar: "KOGIA_ENV",
		},
		cli.StringFlag{
			Name:   "level, l",
			Value:  "warning",
			Usage:  "Verbosity level",
			EnvVar: "KOGIA_LEVEL",
		},
	}
	app.Action = kogia_init
	app.Run(os.Args)
}
