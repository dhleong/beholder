package beholder

type ChallengeRatingQuery struct {
	cr string
}

func NewChallengeRatingQuery(cr string) *ChallengeRatingQuery {
	return &ChallengeRatingQuery{cr}
}

func ExtractChallengeRatingQueries(query string) (remaining string, extracted []Query) {
	return query, nil
}

func (q *ChallengeRatingQuery) Match(entity Entity) (result MatchResult) {
	result = MatchResult{Matched: false, Score: 0, Sequences: nil}
	if entity.GetKind() != MonsterEntity {
		return
	}

	m := entity.(Monster)
	if m.Challenge != q.cr {
		return
	}

	return MatchResult{Matched: true, Score: 1, Sequences: nil}
}
