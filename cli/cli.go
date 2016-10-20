package cli

import (
	"os"

	cli "gopkg.in/urfave/cli.v1"

	"github.com/chasinglogic/dfm"
)

// Added this to make testing easier.
func buildApp() *cli.App {
	app := cli.NewApp()
	app.Name = "dfm"
	app.Usage = "Manage dotfiles."
	app.Version = "1.0-dev"
	app.Authors = []cli.Author{
		{
			Name:  "Mathew Robinson",
			Email: "mathew.robinson3114@gmail.com",
		},
	}

	app.Before = dfm.LoadConfig
	app.After = dfm.SaveConfig

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Use `DIR` for storing dfm configuration and profiles",
			Value: dfm.DefaultConfigDir(),
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Print verbose messaging.",
		},
		cli.BoolFlag{
			Name:  "dry-run",
			Usage: "Don't create symlinks just print what would be done.",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "Add a file to the current profile.",
			Action:  dfm.Add,
		},
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "Create a dotfiles profile from a git repo.",
			Action:  dfm.Create,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "alias, a",
					Usage: "Creates `ALIAS` for the profile instead of username",
				},
				cli.BoolFlag{
					Name:  "overwrite, o",
					Usage: "Overwrites existing files when creating links.",
				},
				cli.BoolFlag{
					Name:  "link, l",
					Usage: "Links the profile after creation.",
				},
			},
		},
		{
			Name:    "link",
			Aliases: []string{"l"},
			Usage:   "Recreate the links from the dotfiles profile.",
			Action:  dfm.Link,
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "overwrite, o",
					Usage: "Overwrites existing files when creating links.",
				},
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "List available profiles",
			Action:  dfm.List,
		},
		{
			Name:    "update",
			Aliases: []string{"up"},
			Usage:   "Pull the latest version of the profile from origin master.",
			Action:  dfm.Update,
		},
		{
			Name:        "remove",
			Aliases:     []string{"rm"},
			Usage:       "Remove the profile and all it's symlinks.",
			Description: "Removes the profile and all it's symlinks, if there is another profile on this system we will switch to it. Otherwise will do nothing.",
			Action:      dfm.Remove,
		},
	}

	return app
}

// Run is the entry point for the app
func Run() int {
	a := buildApp()
	a.Run(os.Args)
	return 0
}