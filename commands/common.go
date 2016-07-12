package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

// Config is used to store global options
type Config struct {
	Verbose        bool
	DryRun         bool
	ConfigDir      string
	CurrentProfile string
}

// VERBOSE is used to globally set verbosity
var VERBOSE = false

// DRYRUN is used to globally set whether this is a dry run
var DRYRUN = false

func loadConfig(c *cli.Context) (*Config, error) {
	var config Config

	configJSON, rerr := ioutil.ReadFile(filepath.Join(c.Parent().String("config"), "config.json"))
	if rerr != nil {
		return &config, rerr
	}

	err := json.Unmarshal(configJSON, &config)
	if err != nil {
		return &config, err
	}

	config.ConfigDir = c.String("config")
	VERBOSE = c.Bool("verbose")
	DRYRUN = c.Bool("dry-run")

	if config.Verbose {
		VERBOSE = config.Verbose
	}

	if config.DryRun {
		DRYRUN = config.DryRun
	}

	return &config, nil
}

// Save will save the config to the configDir/config.json
func (c *Config) Save() error {
	JSON, merr := json.MarshalIndent(c, "", "\t")
	if merr != nil {
		return merr
	}

	return ioutil.WriteFile(c.ConfigDir+"config.json", JSON, 0644)
}

// LinkInfo simulates a tuple for our symbolic link
type LinkInfo struct {
	Src  string
	Dest string
}

func (l *LinkInfo) String() string {
	return fmt.Sprintf("Link( %s, %s )", l.Src, l.Dest)
}

func getProfileDir(c *cli.Context) string {
	return filepath.Join(c.Parent().String("config"), "profiles")
}

func getUser(c *cli.Context) string {
	// This handles the case when create passes us it's context
	if len(strings.Split(c.Args().First(), "/")) > 1 {
		_, user := createURL(strings.Split(c.Args().First(), "/"))
		return user
	}

	return c.Args().First()
}

func generateSymlinks(profileDir string) []LinkInfo {
	links := []LinkInfo{}
	// TODO: Handle the config dir special case
	files, err := ioutil.ReadDir(profileDir)
	if err != nil {
		return links
	}

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), ".") {
			ln := LinkInfo{
				filepath.Join(profileDir, file.Name()),
				filepath.Join(os.Getenv("HOME"), "."+file.Name()),
			}

			if VERBOSE {
				fmt.Printf("Generated symlink %s\n", ln.String())
			}

			links = append(links, ln)
		}
	}

	return links
}
