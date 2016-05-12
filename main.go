package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
)

const (
	port = ":17000"
)

func main() {
	app := cli.NewApp()
	app.Name = "deadpool"
	app.Usage = "let the AIs fight"

	app.Commands = []cli.Command{
		{
			Name:    "new",
			Aliases: []string{"new"},
			Usage:   "create a new deadpool game folder",
			Action: func(c *cli.Context) error {
				if c.NArg() == 0 {
					err := fmt.Errorf("Name argument is mandatory")
					log.Println(err)
					return err
				}
				game := c.Args().First()
				gen := NewGenerator()
				err := gen.Generate(game)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("%s: game successfully initialized", game)
					log.Println("Have fun with your new game!!")
					log.Println("Don't forget to 'protoc' the game")
				}
				return err
			},
		},
		{
			Name:    "list",
			Aliases: []string{"list"},
			Usage:   "list the available games",
			Action: func(c *cli.Context) error {
				lister := NewLister()
				for _, game := range lister.List() {
					log.Println(game)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
