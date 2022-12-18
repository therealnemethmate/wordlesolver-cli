package recommendation

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestFilterWords(t *testing.T) {
	t.Run("should filter words which are present in excludeLetters", func(t *testing.T) {
		rec := Recommender{guess: "_____", words: []string{"asdfg", "qwert", "qkrkr"}, excludeLetters: []string{"r", "e"}}
		words := rec.filterWords("_____")
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
		words := rec.filterWords(guess)
		switch {
		case !slices.Contains(words, "qkrkr"):
			t.Errorf("filterWords shouldn't have filter: %v", "qkrkr")
		case !slices.Contains(words, "qwert"):
			t.Errorf("filterWords shouldn't have filter: %v", "qwert")
		case !slices.Contains(words, "asdfg"):
			t.Errorf("filterWords shouldn't have filter %v", "asdfg")
		}
	})

	t.Run("should exclude words based on one excludeCombination", func(t *testing.T) {
		guess := "_____"
		rec := Recommender{
			guess:               guess,
			words:               []string{"asdfg", "qwert", "qkrkr"},
			excludeCombinations: map[string][]int{"a": {0}},
		}
		words := rec.filterWords(guess)
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
			includeLetters:      []string{"r"},
			excludeCombinations: map[string][]int{"r": {3}},
		}
		words := rec.filterWords(guess)
		switch {
		case !slices.Contains(words, "qkrkr"):
			t.Errorf("filterWords shouldn't have filter: %v", "qkrkr")
		case slices.Contains(words, "qwert"):
			t.Errorf("filterWords should've filter: %v", "qwert")
		case !slices.Contains(words, "asdfg"):
			t.Errorf("filterWords shouldn't have filter %v", "asdfg")
		}
	})

	t.Run("should include only those words, where the letter is not in excludeCombinations", func(t *testing.T) {
		guess := "_____"
		rec := Recommender{
			guess:               guess,
			words:               []string{"asdfg", "qwert", "qkrkr"},
			includeLetters:      []string{"r"},
			excludeCombinations: map[string][]int{"r": {3}},
		}
		words := rec.filterWords(guess)
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
		guess := "a____"
		rec := Recommender{
			guess:               guess,
			words:               []string{"asdfg", "qwert", "qkrkr"},
			includeLetters:      []string{"r"},
			excludeCombinations: map[string][]int{"r": {3}},
		}
		words := rec.filterWords(guess)
		switch {
		case slices.Contains(words, "qkrkr"):
			t.Errorf("filterWords should've filter: %v", "qkrkr")
		case slices.Contains(words, "qwert"):
			t.Errorf("filterWords should've filter: %v", "qwert")
		case !slices.Contains(words, "asdfg"):
			t.Errorf("filterWords shouldn't have filter: %v", "asdfg")
		}
	})
}
