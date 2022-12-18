package cmd

import (
	"os"

	"github.com/therealnemethmate/wordlesolver-cli/internal"
	"github.com/therealnemethmate/wordlesolver-cli/internal/recommendation"
	"github.com/urfave/cli"
)

const (
	Version     = "0.0.1"
	Description = "wordlesolver is a cli application to help you solve your daily worlde at https://www.nytimes.com/games/wordle/index.html"
)

func Start() error {
	(&cli.App{
		Name:        "Wordle Solver",
		HelpName:    "wordlesolver",
		Version:     Version,
		Description: Description,
		Commands:    []cli.Command{},
		Action: func(context *cli.Context) error {
			recommender := recommendation.NewRecommender()
			executor := internal.NewExecutor(recommender)
			return executor.Solve("_____", map[string]int{})
		},
	}).Run(os.Args)

	return nil
}
