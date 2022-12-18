package recommendation

import (
	"log"
	"sort"
	"strings"

	"github.com/therealnemethmate/wordlesolver-cli/internal/utils"
	"golang.org/x/exp/slices"
)

type Recommender struct {
	guess               string
	words               []string
	excludeLetters      []string
	includeLetters      []string
	excludeCombinations map[string][]int
	letterValues        map[string][]int
}

func NewRecommender() *Recommender {
	return &Recommender{guess: ""}
}

func (rec *Recommender) getWordValue(word string) int {
	value := 0
	for i, v := range string(word) {
		letter := string(v)
		value += rec.letterValues[letter][i]
	}
	return value
}

func (rec *Recommender) sortWords() {
	sort.Slice(rec.words, func(a, b int) bool {
		var (
			valueA = rec.getWordValue(rec.words[a])
			valueB = rec.getWordValue(rec.words[b])
		)
		return valueA >= valueB
	})
}

func (rec *Recommender) calculateLetterValues() {
	rec.letterValues = map[string][]int{}
	for _, word := range rec.words {
		for i, char := range word {
			letter := string(char)
			if len(rec.letterValues[letter]) == 0 {
				rec.letterValues[letter] = []int{0, 0, 0, 0, 0}
			}
			rec.letterValues[letter][i] += 1
		}
	}
	log.Println(rec.letterValues)
}

func (rec *Recommender) loadWords() error {
	text, err := utils.Readfile("meta/words.txt")
	if err != nil {
		return err
	}

	rec.words = strings.Split(text, "\n")
	return nil
}

func (rec *Recommender) mergeExcludeCombinations(excludeCombinations map[string]int) {
	if rec.excludeCombinations == nil {
		rec.excludeCombinations = map[string][]int{}
	}
	for k, v := range excludeCombinations {
		if rec.excludeCombinations[k] == nil {
			rec.excludeCombinations[k] = []int{v}
		} else {
			rec.excludeCombinations[k] = append(rec.excludeCombinations[k], v)
		}
	}
}

func (rec *Recommender) filterWords(wordState string) []string {
	words := []string{}

	for _, word := range rec.words {
		include := true
		for i, char := range word {
			letter := string(char)
			if wordState[i] != '_' && string(wordState[i]) != letter {
				include = false
				break
			} else if slices.Contains(rec.excludeCombinations[letter], i) {
				include = false
				break
			} else if slices.Contains(rec.excludeLetters, letter) {
				include = false
				break
			}
		}
		if include {
			words = append(words, word)
		}
	}

	for _, letter := range rec.includeLetters {
		for _, word := range words {
			if strings.Contains(word, letter) {
				words = append(words, word)
			}
		}
	}

	return words
}

func (rec *Recommender) addExcludedLetters(wordState string) {
	for i, wordStateLetter := range wordState {
		letter := string(rec.guess[i])

		if wordStateLetter != '_' {
			if !slices.Contains(rec.includeLetters, letter) {
				rec.includeLetters = append(rec.includeLetters, letter)
			}
			break
		}
		if slices.Contains(rec.includeLetters, letter) {
			if !strings.Contains(wordState, letter) {
				break
			}
			if slices.Contains(rec.excludeCombinations[letter], i) {
				break
			}
			if rec.excludeCombinations[letter] == nil {
				rec.excludeCombinations[letter] = []int{i}
			}
			rec.excludeCombinations[letter] = append(rec.excludeCombinations[letter], i)
			break
		}
		if slices.Contains(rec.excludeLetters, letter) {
			break
		}
		if letter == string(wordStateLetter) {
			break
		}
		rec.excludeLetters = append(rec.excludeLetters, letter)
	}
}

func (rec *Recommender) addIncludedLetters(excludeCombinations map[string]int) {
	for k := range excludeCombinations {
		if !slices.Contains(rec.includeLetters, k) {
			rec.includeLetters = append(rec.includeLetters, k)
		}
	}
}

func (rec *Recommender) GetNext(wordState string, excludeCombinations map[string]int) (string, error) {
	if rec.words == nil {
		err := rec.loadWords()
		if err != nil {
			return "", err
		}
		if rec.letterValues == nil {
			rec.calculateLetterValues()
		}
		rec.sortWords()
	}
	if rec.guess == "" {
		rec.guess = rec.words[0]
		return rec.guess, nil
	}

	rec.addIncludedLetters(excludeCombinations)
	rec.addExcludedLetters(wordState)
	rec.mergeExcludeCombinations(excludeCombinations)
	rec.words = rec.filterWords(wordState)
	log.Printf("Include: %v, Exclude: %v, Comb: %v", rec.includeLetters, rec.excludeLetters, rec.excludeCombinations)
	rec.guess = rec.words[0]
	log.Println(len(rec.words))
	return rec.guess, nil
}
