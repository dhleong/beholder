package beholder

import (
	"strings"
	"unicode"
)

// QueryMatcher .
type QueryMatcher struct {
	runes      []rune
	upperRunes []rune
}

// MatchResult is the result of
type MatchResult struct {
	Matched bool
	Score   float32
}

type matchedSequence struct {
	start        int
	end          int
	startsOnWord bool
}

func (s *matchedSequence) length() int {
	return s.end - s.start
}

// NewQueryMatcher .
func NewQueryMatcher(query string) *QueryMatcher {
	return &QueryMatcher{
		[]rune(query),
		[]rune(strings.ToUpper(query)),
	}
}

// Match .
func (q *QueryMatcher) Match(value string) MatchResult {
	runes := []rune(value)

	sequences := make([]*matchedSequence, 0, 8)

	longestSubsequence := 0
	currentSequence := &matchedSequence{}

	words := 1
	wordStartsMatched := 0
	inWord := true

	j := 0
	for i := 0; i < len(runes); i++ {

		enteredWord := i == 0
		if !unicode.IsLetter(runes[i]) {
			if inWord {
				inWord = false
			}
		} else if !inWord {
			inWord = true
			enteredWord = true
			words++
		}

		if unicode.ToUpper(runes[i]) == q.upperRunes[j] {
			if currentSequence != nil {
				currentSequence.end++
			} else {
				currentSequence := &matchedSequence{
					start:        j,
					end:          j + 1,
					startsOnWord: enteredWord,
				}
				sequences = append(sequences, currentSequence)

				if enteredWord {
					wordStartsMatched++
				}
			}

			// advance
			j++

			if j >= len(q.runes) {
				// matched everything from query; shortcut!
				break
			}
		} else {
			if currentSequence != nil && currentSequence.length() > longestSubsequence {
				longestSubsequence = currentSequence.length()
			}
			currentSequence = nil
		}
	}

	// check again, in case the longest subsequence ends
	// at the end of the match
	if currentSequence != nil && currentSequence.length() > longestSubsequence {
		longestSubsequence = currentSequence.length()
	}
	currentSequence = nil

	// base score on longest subsequence
	score := float32(longestSubsequence) / (float32(len(runes)) / 2.0)

	if len(sequences) == 1 && sequences[0].start == 0 {
		// bonus for matching at the beginning
		score *= float32(sequences[0].length())
	}

	// augment by number of word starts matched
	if words > 0 && wordStartsMatched > 0 {
		score *= float32(wordStartsMatched) / float32(words)
	}

	return MatchResult{
		Matched: j == len(q.runes),
		Score:   score,
	}
}
