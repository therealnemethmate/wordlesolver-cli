package internal

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func NewSolver(recommender *Recommender) *Solver {
	return &Solver{"", recommender}
}

type Solver struct {
	nextGuess   string
	recommender *Recommender
}

func (solver Solver) getTextInputFromUser() string {
	fmt.Println("Please input the current state of the game after you guessed again:")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	if text == "" {
		fmt.Println("Please enter a value")
		return solver.getTextInputFromUser()
	}
	return text
}

func (solver Solver) getExcludedCombinations() map[string]int {
	fmt.Println("Please input which letters were correct with their position in the word (positions start at 0). Please separate all combinations with commas.")
	fmt.Println("For example: a0,b2")
	reader := bufio.NewReader(os.Stdin)
	//TODO error handling
	excludeCombinations := map[string]int{}
	text, _ := reader.ReadString('\n')
	if text == "" {
		return excludeCombinations
	}
	inputs := strings.Split(text, ",")

	for _, v := range inputs {
		input := strings.Split(v, "")
		position, _ := strconv.Atoi(input[1])
		excludeCombinations[input[0]] = position
	}
	return excludeCombinations
}

func (solver Solver) Execute(wordState string, excludeCombinations map[string]int) error {
	if len(wordState) != 5 {
		error := fmt.Errorf("you should provide exactly 5 letters and / or placeholders for state %v", wordState)
		return error
	}

	if strings.Contains(wordState, "_") {
		nextGuess, err := solver.recommender.GetNext(wordState, excludeCombinations)
		if err != nil {
			return err
		}
		solver.nextGuess = nextGuess
		fmt.Printf("Your next guess should be: %v\n", nextGuess)
		text := solver.getTextInputFromUser()
		return solver.Execute(strings.Trim(text, "\n"), solver.getExcludedCombinations())
	}
	fmt.Printf("Congratulations, your solution is: %v\n", wordState)
	return nil
}
