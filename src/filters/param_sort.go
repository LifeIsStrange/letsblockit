package filters

import (
	"fmt"
	"sort"
	"strings"
)

func (f *Filter) SortValues(paramName string, params map[string]interface{}) error {
	var sortType ParamSortType
	for _, p := range f.Params {
		if p.Name == paramName {
			sortType = p.Sort
		}
	}
	switch sortType {
	case NaturalSliceSort:
		v, ok := params[paramName].([]string)
		if !ok {
			return fmt.Errorf("invalid value type for param %s", paramName)
		}
		sort.Slice(v, func(i, j int) bool {
			return strings.Compare(v[i], v[j]) < 0
		})
		params[paramName] = v
	case DomainSliceSort:
		v, ok := params[paramName].([]string)
		if !ok {
			return fmt.Errorf("invalid value type for param %s", paramName)
		}
		sort.Slice(v, func(i, j int) bool {
			return strings.Compare(stripToDomain(v[i]), stripToDomain(v[j])) < 0
		})
		params[paramName] = v
	case DomainRuleSort:
		v, ok := params[paramName].(string)
		if !ok {
			return fmt.Errorf("invalid value type for param %s", paramName)
		}
		var lines []string
		for _, line := range strings.Split(v, "\n") {
			if strings.HasPrefix(line, "!") || len(line) == 0 {
				continue
			}
			lines = append(lines, line)
		}
		sort.Slice(lines, func(i, j int) bool {
			return strings.Compare(stripToDomain(lines[i]), stripToDomain(lines[j])) < 0
		})
		params[paramName] = strings.Join(lines, "\n")
	default:
		return fmt.Errorf("unknown sort type %s for param %s", sortType, paramName)
	}
	return nil
}

func stripToDomain(input string) string {
	output := input
	for _, p := range []string{"https://", "http://", "www.", "*.", "."} {
		output = strings.TrimPrefix(output, p)
	}
	return output
}
