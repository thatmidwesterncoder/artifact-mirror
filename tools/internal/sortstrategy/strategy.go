package sortstrategy

import (
	"fmt"
	"slices"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/maruel/natural"
)

func Natural(t []string) ([]string, error) {
	st := make([]string, len(t))
	copy(st, t)
	slices.SortFunc(st, natural.Compare)
	return st, nil
}

func Semver(t []string) ([]string, error) {
	vs := make([]*semver.Version, len(t))
	for i, r := range t {
		v, err := semver.NewVersion(r)
		if err != nil {
			return nil, fmt.Errorf("error parsing version: %s", err)
		}
		vs[i] = v
	}
	sort.Sort(semver.Collection(vs))

	st := make([]string, len(t))
	for _, t := range vs {
		st = append(st, t.String())
	}
	return st, nil
}
