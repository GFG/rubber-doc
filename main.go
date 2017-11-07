package main

import (
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/rocket-internet-berlin/RocketLabsRubberDoc/command"
	"github.com/urfave/cli"
)

func main() {
	cmd := &command.GenerateCommand{}

	logger := logrus.New()

	app := cli.NewApp()
	app.Name = "RubberDoc"
	app.Version = "v0.1-alpha-2"
	app.Description = "A documentation generator for RAML and Blueprint."
	app.Usage = ""

	var debugLogging bool
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Enable debug logging",
			Destination: &debugLogging,
		},
	}

	app.Before = func(c *cli.Context) error {
		logger.Formatter = &logrus.TextFormatter{FullTimestamp: true}
		if c.GlobalBool("debug") {
			logger.Level = logrus.DebugLevel
			logger.Debug("Debug logging enabled")
			logger.Debug(app.Name, "-", app.Version)
		}
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "generate",
			Usage: "This command receives a configuration file and a specification file written in RAML or Blueprint.",

			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "spec",
					Value:       "",
					Usage:       "Specify the Specification's file location.",
					Destination: &cmd.SpecFile,
				},
				cli.StringFlag{
					Name:        "config",
					Value:       "",
					Usage:       "Specify the configuration's file location.",
					Destination: &cmd.ConfigFile,
				},
			},
			Action: func(c *cli.Context) {
				if err := cmd.Execute(); err != nil {
					logger.Error(err)
				}
			},
		},
	}

	app.Run(os.Args)
}
