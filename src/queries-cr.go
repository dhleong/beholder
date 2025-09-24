package beholder

import (
	"regexp"
)

type ChallengeRatingQuery struct {
	cr string
}

func NewChallengeRatingQuery(cr string) *ChallengeRatingQuery {
	return &ChallengeRatingQuery{cr}
}

var crQueryRegex *regexp.Regexp = regexp.MustCompile(`(?i)\s*CR([0-9/]+)\s*`)

func ExtractChallengeRatingQueries(query string) (remaining string, extracted []Query) {
	matches := crQueryRegex.FindAllStringSubmatch(query, -1)
	if len(matches) > 0 {
		extracted = make([]Query, 0, len(matches))
		for _, match := range matches {
			cr := match[1]
			q := NewChallengeRatingQuery(cr)
			extracted = append(extracted, q)
		}
	}
	remaining = crQueryRegex.ReplaceAllLiteralString(query, "")
	return
}

func (q *ChallengeRatingQuery) Match(entity Entity) MatchResult {
	if entity.GetKind() != MonsterEntity {
		return MatchResult{Matched: false}
	}

	m := entity.(Monster)
	if m.Challenge != q.cr {
		return MatchResult{Matched: false}
	}

	return MatchResult{Matched: true, Score: 1, Sequences: nil}
}
