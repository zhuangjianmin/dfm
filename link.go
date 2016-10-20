package dfm

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	cli "gopkg.in/urfave/cli.v1"
)

// Link will generate and create the symlinks to the dotfiles in the repo.
func Link(c *cli.Context) error {
	userDir := filepath.Join(getProfileDir(), c.Args().First())
	links := GenerateSymlinks(userDir)
	if err := CreateSymlinks(links, c.Bool("overwrite")); err != nil {
		return cli.NewExitError(err.Error(), 2)
	}

	CONFIG.CurrentProfile = c.Args().First()
	return nil
}

// LinkInfo simulates a tuple for our symbolic link
type LinkInfo struct {
	Src  string
	Dest string
}

func (l *LinkInfo) String() string {
	return fmt.Sprintf("Link( %s, %s )", l.Src, l.Dest)
}

func GenerateSymlinks(profileDir string) []LinkInfo {
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

			if CONFIG.Verbose {
				fmt.Printf("Generated symlink %s\n", ln.String())
			}

			links = append(links, ln)
		}
	}

	return links
}

func CreateSymlinks(l []LinkInfo, overwrite bool) error {
	ok := true

	for _, link := range l {
		if _, err := os.Stat(link.Dest); err == nil {
			if overwrite {
				if CONFIG.Verbose || DRYRUN {
					fmt.Printf("%s already exists, removing.\n", link.Dest)
				}

				if !DRYRUN {
					if rmerr := os.Remove(link.Dest); rmerr != nil {
						fmt.Printf("Unable to remove %s: %s\n",
							link.Dest,
							rmerr.Error())
					}
				}
			} else {
				fmt.Printf("%s already exists.\n", link.Dest)
				ok = false
			}
		}
	}

	if ok {
		for _, link := range l {
			if CONFIG.Verbose || DRYRUN {
				fmt.Printf("Creating symlink %s -> %s\n", link.Src, link.Dest)
			}

			if !DRYRUN {
				if err := os.Symlink(link.Src, link.Dest); err != nil {
					return err
				}
			}
		}

		return nil
	}

	return cli.NewExitError("Symlink targets exist. Refusing to create a broken state please remove the targets then rerun command.", 68)
}
