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
	return &Recommender{
		guess:               "",
		excludeCombinations: make(map[string][]int),
		excludeLetters:      []string{},
		includeLetters:      []string{},
	}
}

func (rec *Recommender) getWordValue(word string) int {
	value := 0
	for i, v := range string(word) {
		letter := string(v)
		value += rec.letterValues[letter][i]
	}
	return value
}

func (rec *Recommender) sortWords(words []string) []string {
	sort.Slice(rec.words, func(a, b int) bool {
		var (
			valueA = rec.getWordValue(rec.words[a])
			valueB = rec.getWordValue(rec.words[b])
		)
		return valueA >= valueB
	})
	return words
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
}

func (rec *Recommender) loadWords() error {
	text, err := utils.Readfile("meta/words.txt")
	if err != nil {
		return err
	}

	rec.words = strings.Split(text, "\n")
	return nil
}

func (rec *Recommender) filterExcludedWords(wordState string, words []string) []string {
	result := []string{}

	for _, word := range words {
		include := true
		for i, char := range word {
			letter := string(char)
			if slices.Contains(rec.excludeCombinations[letter], i) {
				include = false
			} else if slices.Contains(rec.excludeLetters, letter) {
				if strings.Contains(wordState, letter) {
					// TODO handle what happens when it contains it twice
				}
				include = false
			}
		}
		if include {
			result = append(result, word)
		}
	}

	return result
}

func (rec *Recommender) filterWords(wordState string, originalWords []string) []string {
	words := []string{}

	for _, word := range originalWords {
		if len(word) != 5 {
			break
		}
		include := true
		for i, c := range wordState {
			if c != '_' && c != '*' && string(word[i]) != string(c) {
				include = false
			}
		}
		if include {
			words = append(words, word)
		}
	}

	for i, char := range wordState {
		if char != '_' && char != '*' && slices.Contains(rec.includeLetters, string(char)) {
			slices.Delete(rec.includeLetters, i, i+1)
		}
	}

	if len(rec.includeLetters) != 0 {
		for _, word := range words {
			include := true
			for _, letter := range rec.includeLetters {
				if !strings.Contains(word, letter) {
					include = false
					break
				}
			}
			if include {
				words = append(words, word)
			}
		}
	}

	return rec.filterExcludedWords(wordState, words)
}

func (rec *Recommender) addExcludedLetters(wordState string) {
	for i, char := range wordState {
		guessedLetter := string(rec.guess[i])
		if guessedLetter == string(char) {
			continue
		}
		if char == '_' && !slices.Contains(rec.includeLetters, guessedLetter) {
			rec.excludeLetters = append(rec.excludeLetters, guessedLetter)
		}
	}
}

func (rec *Recommender) addIncludedLetters(wordState string) {
	for _, char := range wordState {
		if char == '*' && !slices.Contains(rec.includeLetters, string(char)) {
			rec.includeLetters = append(rec.includeLetters, string(char))
		}
	}
}

func (rec *Recommender) addCombinations(wordState string) {
	for i, char := range wordState {
		key := string(char)
		if char == '*' && !slices.Contains(rec.excludeCombinations[key], i) {
			rec.excludeCombinations[key] = append(rec.excludeCombinations[key], i)
		}
	}
}

func (rec *Recommender) GetNext(wordState string) (string, error) {
	if rec.words == nil {
		err := rec.loadWords()
		if err != nil {
			return "", err
		}
		if rec.letterValues == nil {
			rec.calculateLetterValues()
		}
		rec.words = rec.sortWords(rec.words)
	}
	if rec.guess == "" {
		rec.guess = rec.words[0]
		return rec.guess, nil
	}

	rec.addIncludedLetters(wordState)
	rec.addCombinations(wordState)
	rec.addExcludedLetters(wordState)

	rec.words = rec.filterWords(wordState, rec.words)

	log.Printf("Include: %v, Exclude: %v, Comb: %v", rec.includeLetters, rec.excludeLetters, rec.excludeCombinations)
	log.Println(rec.words)
	rec.guess = rec.words[0]
	return rec.guess, nil
}
