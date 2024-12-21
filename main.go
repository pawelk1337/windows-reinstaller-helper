package main

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Name:  "win-reinstaller-helper",
		Usage: "A tool that helps you during Windows reinstallation. It helps you save important files and folders.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "basepath",
				Usage: "The base path of the save location. (can also be set in env as WRH_PATH)",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "restore",
				Usage:   "Restore the file/folder to its original location.",
				Aliases: []string{"r"},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					base := os.Getenv("WRH_PATH")
					if base == "" {
						if cmd.String("basepath") == "" {
							panic("basepath is required")
						} else {
							base = cmd.String("basepath")
						}
					}

					Resotre(base)
					return nil
				},
			},
			{
				Name:    "save",
				Usage:   "Save the file/folder to the database.",
				Aliases: []string{"s"},
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "same-location",
						Usage:   "When restoring the files and folders, they will be saved in the same location.",
						Value:   false,
						Aliases: []string{"s"},
					},
					&cli.StringFlag{
						Required: true,
						Name:     "type",
						Usage:    "file/folder/name",
						Aliases:  []string{"t"},
					},
					&cli.StringFlag{
						Required: true,
						Name:     "path",
						Usage:    "path of the file/folder",
						Aliases:  []string{"p"},
					},
					&cli.BoolFlag{
						Name:    "no-saved",
						Usage:   "disable the .saved file creation",
						Aliases: []string{"n"},
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					base := os.Getenv("WRH_PATH")
					if base == "" {
						if cmd.String("basepath") == "" {
							panic("basepath is required")
						} else {
							base = cmd.String("basepath")
						}
					}

					switch cmd.String("type") {
					case "file":
						SaveFile(base, cmd.String("path"), cmd.Bool("same-location"), cmd.Bool("no-saved"))
					case "folder":
						SaveDir(base, cmd.String("path"), cmd.Bool("same-location"), cmd.Bool("no-saved"))
					case "name":
						SaveName(base, cmd.String("path"))
					default:
						panic("type not found")
					}
					return nil
				},
			},
			{
				Name:    "delete",
				Usage:   "Delete a file or folder",
				Aliases: []string{"d"},
				Flags: []cli.Flag{
					&cli.StringFlag{
						Required: true,
						Name:     "path",
						Usage:    "path of the file/folder",
						Aliases:  []string{"p"},
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					base := os.Getenv("WRH_PATH")
					if base == "" {
						if cmd.String("basepath") == "" {
							panic("basepath is required")
						} else {
							base = cmd.String("basepath")
						}
					}

					Delete(base, cmd.String("path"))
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
