package beholder

import (
	"strings"
	"unicode"
)

type Query interface {
	Match(Entity) MatchResult
}

// QueryMatcher .
type QueryMatcher struct {
	queries    []Query
	runes      []rune
	upperRunes []rune
}

type MatchResult struct {
	Matched   bool
	Score     float32
	Sequences []*MatchedSequence
}

// MatchedSequence represents a range of an Entity's
// name that matched the query string
type MatchedSequence struct {
	Start        int
	End          int
	startsOnWord bool
}

func (s *MatchedSequence) length() int {
	return s.End - s.Start
}

// NewQueryMatcher .
func NewQueryMatcher(query string) *QueryMatcher {
	remainingQuery, extracted := ExtractQueries(query)
	return &QueryMatcher{
		extracted,
		[]rune(remainingQuery),
		[]rune(strings.ToUpper(remainingQuery)),
	}
}

func (q *QueryMatcher) Match(entity Entity) MatchResult {
	var result *MatchResult = nil

	if len(q.runes) > 0 {
		r := q.MatchName(entity.GetName())
		result = &r
	}

	if len(q.queries) > 0 {
		// If we have any Queries, the Entity must match ALL
		// *in addition to* any provided text query
		for _, query := range q.queries {
			m := query.Match(entity)
			if !m.Matched {
				if result != nil {
					result.Matched = false
				}
				break
			} else if result == nil {
				// If we don't have any result yet, then we did
				// not have a query string. We found a match here
				result = &m
			}
		}
	}

	if result == nil {
		return MatchResult{Matched: false}
	}
	return *result
}

func (q *QueryMatcher) MatchName(value string) MatchResult {
	runes := []rune(value)

	sequences := make([]*MatchedSequence, 0, 8)

	longestSubsequence := 0
	var currentSequence *MatchedSequence

	words := 0
	wordStartsMatched := 0
	inWord := true

	j := 0
	for i := range runes {

		enteredWord := i == 0
		if !unicode.IsLetter(runes[i]) {
			if inWord {
				inWord = false
			}
		} else if !inWord {
			inWord = true
			enteredWord = true
		}

		if enteredWord {
			words++
		}

		if unicode.ToUpper(runes[i]) == q.upperRunes[j] {
			if currentSequence != nil {
				currentSequence.End++
			} else {
				currentSequence = &MatchedSequence{
					Start:        i,
					End:          i + 1,
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

	if len(sequences) == 1 && sequences[0].Start == 0 {
		// bonus for matching at the beginning
		score *= float32(sequences[0].length())
	}

	// augment by number of word starts matched
	if words > 0 && wordStartsMatched > 0 {
		score *= float32(wordStartsMatched) / float32(words)
	}

	return MatchResult{
		Matched:   j == len(q.runes),
		Score:     score,
		Sequences: sequences,
	}
}
