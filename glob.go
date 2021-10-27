package fscopy

import (
	"path/filepath"

	"github.com/nofeaturesonlybugs/errors"
)

// Glob calls filepath.Glob() on each glob.  If any glob returns an error it is skipped and the next
// glob is processed; if an error is returned it contains the globs that failed.
//
// Matches from globs are deduped while adding them to the return value.
func Glob(globs ...string) ([]string, error) {
	var rv []string
	seen := map[string]struct{}{}
	bad := []string{}
	for _, glob := range globs {
		matches, err := filepath.Glob(glob)
		if err != nil {
			bad = append(bad, glob)
			continue
		}
		if len(matches) == 0 {
			continue
		}
		unseen := []string{}
		for _, match := range matches {
			if _, ok := seen[match]; ok {
				continue
			}
			unseen = append(unseen, match)
			seen[match] = struct{}{}
		}
		rv = append(rv, unseen...)
	}
	if len(bad) > 0 {
		return rv, errors.Errorf("invalid patterns: %v", bad)
	}
	return rv, nil
}
