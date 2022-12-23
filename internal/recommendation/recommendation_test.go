package recommendation

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestFilterWords(t *testing.T) {
	t.Run("should filter words which are present in excludeLetters", func(t *testing.T) {
		rec := Recommender{guess: "_____", words: []string{"asdfg", "qwert", "qkrkr"}, excludeLetters: []string{"r", "e"}}
		words := rec.filterWords("_____", rec.words)
		switch {
		case slices.Contains(words, "qkrkr"):
			t.Errorf("filterWords should've filter: %v", "qkrkr")
		case slices.Contains(words, "qwert"):
			t.Errorf("filterWords should've filter: %v", "qwert")
		case !slices.Contains(words, "asdfg"):
			t.Errorf("filterWords shouldn't have filter %v", "asdfg")
		}
	})

	t.Run("should include words with letters in includeLetters", func(t *testing.T) {
		guess := "_____"
		rec := Recommender{guess: guess, words: []string{"asdfg", "qwert", "qkrkr"}, includeLetters: []string{"r", "e"}}
		words := rec.filterWords(guess, rec.words)
		switch {
		case slices.Contains(words, "qkrkr"):
			t.Errorf("filterWords should've filter: %v", "qkrkr")
		case !slices.Contains(words, "qwert"):
			t.Errorf("filterWords shouldn't have filter: %v", "qwert")
		case slices.Contains(words, "asdfg"):
			t.Errorf("filterWords should've filter %v", "asdfg")
		}
	})

	t.Run("should exclude words based on one excludeCombination", func(t *testing.T) {
		guess := "_____"
		rec := Recommender{
			guess:               guess,
			words:               []string{"asdfg", "qwert", "qkrkr"},
			excludeCombinations: map[string][]int{"a": {0}},
		}
		words := rec.filterWords(guess, rec.words)
		switch {
		case !slices.Contains(words, "qkrkr"):
			t.Errorf("filterWords shouldn't have filter: %v", "qkrkr")
		case !slices.Contains(words, "qwert"):
			t.Errorf("filterWords shouldn't have filter: %v", "qwert")
		case slices.Contains(words, "asdfg"):
			t.Errorf("filterWords should've filter %v", "asdfg")
		}
	})

	t.Run("should include only those words, where the letter is not in excludeCombinations", func(t *testing.T) {
		guess := "_____"
		rec := Recommender{
			guess:               guess,
			words:               []string{"asdfg", "qwert", "qkrkr"},
			includeLetters:      []string{},
			excludeCombinations: map[string][]int{"r": {3}},
		}
		words := rec.filterWords(guess, rec.words)
		switch {
		case !slices.Contains(words, "qkrkr"):
			t.Errorf("filterWords shouldn't have filter: %v", "qkrkr")
		case slices.Contains(words, "qwert"):
			t.Errorf("filterWords should've filter: %v", "qwert")
		case !slices.Contains(words, "asdfg"):
			t.Errorf("filterWords shouldn't have filter %v", "asdfg")
		}
	})

	t.Run("should include only those words where the letter is not in excludeCombinations - with letter from state", func(t *testing.T) {
		guess := "_a_e_"
		rec := Recommender{
			guess:               guess,
			words:               []string{"caaed", "waver", "tamer", "ramen"},
			includeLetters:      []string{"r"},
			excludeLetters:      []string{"s", "o"},
			excludeCombinations: map[string][]int{"r": {0, 2}},
		}
		words := rec.filterWords(guess, rec.words)
		switch {
		case slices.Contains(words, "caaed"):
			t.Errorf("filterWords should've filter: %v", "caaed")
		case slices.Contains(words, "ramen"):
			t.Errorf("filterWords should've filter: %v", "ramen")
		case !slices.Contains(words, "tamer"):
			t.Errorf("filterWords shouldn't have filter: %v", "tamer")
		}
	})

	t.Run("should be ok", func(t *testing.T) {
		guess := "torta"
		wordState := "_orta"
		rec := Recommender{
			guess:               guess,
			words:               []string{"torta", "aorta"},
			includeLetters:      []string{},
			excludeLetters:      []string{"s", "e", "b", "p", "y", "t"},
			excludeCombinations: map[string][]int{},
		}
		words := rec.filterWords(wordState, rec.words)
		switch {
		case slices.Contains(words, "torta"):
			t.Errorf("filterWords should've filter: %v", "torta")
		case !slices.Contains(words, "aorta"):
			t.Errorf("filterWords shouldn't have filter: %v", "aorta")
		}
	})
}
