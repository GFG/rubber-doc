package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/command"
	"github.com/urfave/cli"
)

var Version = "0.1.0-Dev"
var ApplicationName = "Rubber Doc - A simple documentation generator"

var (
	GenerateCommand = &command.GenerateCommand{}
)

func main() {
	app := cli.NewApp()
	app.Name = ApplicationName
	app.Version = Version
	app.Description = "Using a specification written in RAML, Blueprint, etc, generates an output based on the given format (html, phtml, markdown, etc)."

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})

	var debugLogging bool
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
	}
	app.Before = func(c *cli.Context) error {
		if c.GlobalBool("debug") {
			log.SetLevel(log.DebugLevel)
			log.Debug("Debug logging enabled")
			log.Debug(ApplicationName, "-", Version)
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:  "generate",
			Usage: "Generates an output (e.g html file) based on input format (e.g. blueprint)",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "src",
					Value:       "",
					Usage:       "Specification file location.",
					Destination: &GenerateCommand.Src,
				},
				cli.StringFlag{
					Name:        "output-dir",
					Value:       ".",
					Usage:       "Directory to output the generated file(s).",
					Destination: &GenerateCommand.OutputDir,
				},
				cli.StringFlag{
					Name:        "output-format",
					Value:       "html",
					Usage:       "Format of the generated file(s). e.g. html, phtml or markdown.",
					Destination: &GenerateCommand.OutputFormat,
				},
			},
			Action: func(c *cli.Context) {
				if err := GenerateCommand.Execute(); err != nil {
					log.Error(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
