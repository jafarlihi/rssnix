package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/gilliek/go-opml/opml"
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const Version = "0.3.0"

func addFeed(name string, url string) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get home path")
		os.Exit(1)
	}
	cfg, err := ini.Load(homePath + "/.config/rssnix/config.ini")
	for _, key := range cfg.Section("feeds").Keys() {
		if key.Name() == name {
			return errors.New("Feed named '" + name + "' already exists")
		}
	}
	file, err := os.OpenFile(homePath+"/.config/rssnix/config.ini", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString("\n" + name + " = " + url)
	return err
}

func bashCompleteFeeds(cCtx *cli.Context) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get home path")
		os.Exit(1)
	}
	cfg, err := ini.Load(homePath + "/.config/rssnix/config.ini")
	if err != nil {
		// config not foud, no feeds there
		return
	}

	args := cCtx.Args().Slice()
	presentFeeds := map[string]bool{}
	for _, feed := range args {
		presentFeeds[feed] = true
	}

	for _, key := range cfg.Section("feeds").Keys() {
		if _, ok := presentFeeds[key.Name()]; ok {
			continue
		}
		fmt.Println(key.Name())
	}
}

func main() {
	syscall.Umask(0)
	LoadConfig()

	app := &cli.App{
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "opens the config file with $EDITOR",
				Action: func(cCtx *cli.Context) error {
					editor, ok := os.LookupEnv("EDITOR")
					if len(editor) == 0 || !ok {
						return errors.New("$EDITOR environment variable is not set")
					}
					homePath, err := os.UserHomeDir()
					if err != nil {
						log.Error("Failed to get home path")
						os.Exit(1)
					}
					cmd := exec.Command(editor, homePath+"/.config/rssnix/config.ini")
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					return cmd.Run()
				},
			},
			{
				Name:    "refetch",
				Aliases: []string{"r"},
				Usage:   "delete and refetch given feed(s) or all feeds if no argument is given",
				Action: func(cCtx *cli.Context) error {
					InitialiseNewArticleDirectory()
					if cCtx.Args().Len() == 0 {
						UpdateAllFeeds(true)
					}
					for i := 0; i < cCtx.Args().Len(); i++ {
						UpdateFeed(cCtx.Args().Get(i), true)
					}
					return nil
				},
				BashComplete: bashCompleteFeeds,
			},
			{
				Name:    "update",
				Aliases: []string{"u"},
				Usage:   "update given feed(s) or all feeds if no argument is given",
				Action: func(cCtx *cli.Context) error {
					InitialiseNewArticleDirectory()
					if cCtx.Args().Len() == 0 {
						UpdateAllFeeds(false)
					}
					for i := 0; i < cCtx.Args().Len(); i++ {
						UpdateFeed(cCtx.Args().Get(i), false)
					}
					return nil
				},
				BashComplete: bashCompleteFeeds,
			},
			{
				Name:    "open",
				Aliases: []string{"o"},
				Usage:   "open given feed's directory or root feeds directory if no argument is given",
				Action: func(cCtx *cli.Context) error {
					var path string
					if cCtx.Args().Len() == 0 {
						path = Config.FeedDirectory
					} else {
						path = Config.FeedDirectory + "/" + cCtx.Args().Get(0)
					}
					cmd := exec.Command(Config.Viewer, path)
					cmd.Stdin = os.Stdin
					cmd.Stdout = os.Stdout
					return cmd.Run()
				},
				BashComplete: bashCompleteFeeds,
			},
			{
				Name:    "add",
				Aliases: []string{"a"},
				Usage:   "add a given feed to config",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() != 2 {
						return errors.New("exactly two arguments are required, first being feed name, second being URL")
					}
					return addFeed(cCtx.Args().Get(0), cCtx.Args().Get(1))
				},
			},
			{
				Name:    "import",
				Aliases: []string{"i"},
				Usage:   "import an OPML file",
				Action: func(cCtx *cli.Context) error {
					if cCtx.Args().Len() != 1 {
						return errors.New("argument specifying OPML file path or URL is required")
					}
					doc, err := opml.NewOPMLFromFile(cCtx.Args().Get(0))
					if err != nil {
						doc, err = opml.NewOPMLFromURL(cCtx.Args().Get(0))
						if err != nil {
							return err
						}
					}
					for _, outline := range doc.Body.Outlines {
						if len(outline.XMLURL) > 0 {
							var title string
							if len(outline.Title) > 0 {
								title = outline.Title
							} else if len(outline.Text) > 0 {
								title = outline.Text
							} else {
								continue
							}
							err = addFeed(strings.ReplaceAll(title, " ", "-"), outline.XMLURL)
							if err != nil {
								log.Error("Failed to add feed titled '" + title + "', error: " + err.Error())
								continue
							}
						}
						for _, innerOutline := range outline.Outlines {
							if len(innerOutline.XMLURL) > 0 {
								var title string
								if len(innerOutline.Title) > 0 {
									title = innerOutline.Title
								} else if len(innerOutline.Text) > 0 {
									title = innerOutline.Text
								} else {
									continue
								}
								err = addFeed(strings.ReplaceAll(title, " ", "-"), innerOutline.XMLURL)
								if err != nil {
									log.Error("Failed to add feed titled '" + title + "', error: " + err.Error())
									continue
								}
							}
						}
					}
					return nil
				},
			},
			{
				Name:    "version",
				Aliases: []string{"v"},
				Usage:   "display the version",
				Action: func(cCtx *cli.Context) error {
					log.Info(Version)
					return nil
				},
			},
			{
				Name:    "setup",
				Aliases: []string{"s"},
				Usage:   "sets up autocomplete",

				Action: func(cCtx *cli.Context) error {
					if cCtx.NArg() != 1 {
						return cli.Exit("needs 1 argument, either 'bash' or 'zsh'", 1)
					}
					t := cCtx.Args().Get(0)
					if t != "bash" && t != "zsh" {
						return cli.Exit("argument has to be either 'bash' or 'zsh'", 1)
					}
					err := SetupAutocomplete(t)
					if err != nil {
						return cli.Exit(err, 1)
					}
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Error(err)
	}
}
