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
	app.Version = "lite"
	app.Author = "Dmitry Fedorov"
	app.Email = "gmajere@gmail.com"

	app.Flags = []cli.Flag{
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
