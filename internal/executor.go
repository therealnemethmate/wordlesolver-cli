package internal

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/therealnemethmate/wordlesolver-cli/internal/recommendation"
)

func NewExecutor(recommender *recommendation.Recommender) *Executor {
	return &Executor{"", recommender}
}

type Executor struct {
	nextGuess   string
	recommender *recommendation.Recommender
}

func (executor Executor) getTextInputFromUser() string {
	fmt.Println("Please input the current state of the game after you guessed again.")
	fmt.Println("Use '_' if the letter is not in the word and use '#' if the letter is in the word, but in the wrong spot.")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if text == "" {
		fmt.Println("Please enter a valid value!")
		return executor.getTextInputFromUser()
	}
	return text
}

func (executor Executor) Solve(wordState string) error {
	if len(wordState) != 5 {
		error := fmt.Errorf("you should provide exactly 5 letters and / or placeholders (_, *) for state %v", wordState)
		return error
	}

	if strings.Contains(wordState, "_") {
		nextGuess, err := executor.recommender.GetNext(wordState)
		if err != nil {
			return err
		}
		executor.nextGuess = nextGuess
		fmt.Printf("Your next guess should be: %v\n", nextGuess)
		text := executor.getTextInputFromUser()
		text = strings.ToLower(text)
		return executor.Solve(strings.Trim(text, "\n"))
	}
	fmt.Printf("Congratulations, your solution is: %v\n", wordState)
	return nil
}
