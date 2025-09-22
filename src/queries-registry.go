package beholder

var queryExtractors = []func(string) (string, []Query){
	ExtractChallengeRatingQueries,
}

func ExtractQueries(query string) (remaining string, extracted []Query) {
	remaining = query
	extracted = make([]Query, 0)
	for _, extractor := range queryExtractors {
		newRemaining, newExtracted := extractor(remaining)
		remaining = newRemaining
		extracted = append(extracted, newExtracted...)
	}
	return
}
