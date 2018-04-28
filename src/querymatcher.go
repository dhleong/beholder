package beholder

import "strings"

// QueryMatcher .
type QueryMatcher struct {
	query string
}

// NewQueryMatcher .
func NewQueryMatcher(query string) *QueryMatcher {
	return &QueryMatcher{strings.ToUpper(query)}
}

// Matches .
func (q *QueryMatcher) Matches(value string) bool {
	return strings.Contains(strings.ToUpper(value), q.query)
}
