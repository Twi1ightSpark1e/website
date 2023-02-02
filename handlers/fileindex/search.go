package fileindex

import (
	"regexp"
	"strings"
)

type findParams struct {
	FindQuery string
	FindMatchCase bool
	FindRegex bool
}

func (h *handler) nameMatchesSearchParams(name string, params findParams) (bool, error) {
	if len(params.FindQuery) == 0 {
		return true, nil
	}

	if params.FindMatchCase {
		return strings.Contains(name, params.FindQuery), nil
	}

	if params.FindRegex {
		return regexp.Match(params.FindQuery, []byte(name))
	}

	return strings.Contains(strings.ToLower(name), strings.ToLower(params.FindQuery)), nil
}
