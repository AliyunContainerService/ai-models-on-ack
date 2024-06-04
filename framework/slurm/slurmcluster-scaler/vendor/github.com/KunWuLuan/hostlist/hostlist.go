package hostlist

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrNestedBrackets       = errors.New("nested brackets")
	ErrUnbalancedBrackets   = errors.New("unbalanced brackets")
	ErrBadRange             = errors.New("bad range")
	ErrStartGreaterThanStop = errors.New("start > stop")
	ErrRangeTooLarge        = errors.New("range too large")
	ErrResultsTooLarge      = errors.New("results too large")
	MAX_SIZE                = 100000
)

// see https://www.nsc.liu.se/~kent/python-hostlist/
func ExpandHostlist(hostlist string, allow_duplicates, sort bool) ([]string, error) {
	// Expand a hostlist expression string to a Python list.

	// Example: expand_hostlist("n[9-11],d[01-02]") ==>
	//          ['n9', 'n10', 'n11', 'd01', 'd02']

	// Unless allow_duplicates is true, duplicates will be purged
	// from the results. If sort is true, the output will be sorted.

	results := []string{}
	bracket_level := 0
	part := ""
	for _, c := range hostlist + "," {
		if c == ',' && bracket_level == 0 {
			// Comma at top level, split!
			if part != "" {
				partResults, err := expand_part(part)
				if err != nil {
					return nil, err
				}
				results = append(results, partResults...)
			}
			part = ""
		} else {
			part += string(c)
		}
		if string(c) == "[" {
			bracket_level += 1
		} else if string(c) == "]" {
			bracket_level -= 1
		}
		if bracket_level > 1 {
			return nil, ErrNestedBrackets
		} else if bracket_level < 0 {
			return nil, ErrUnbalancedBrackets
		}
	}
	if bracket_level > 0 {
		return nil, ErrUnbalancedBrackets
	}
	if !allow_duplicates {
		results = remove_duplicates(results)
	}
	if sort {
		slices.Sort(results)
	}
	return results, nil
}

func expand_part(s string) ([]string, error) {
	// Expand a part (e.g. "x[1-2]y[1-3][1-3]") (no outer level commas).
	if s == "" {
		return []string{""}, nil
	}

	// Split into:
	// 1) prefix string (may be empty)
	// 2) rangelist in brackets (may be missing)
	// 3) the rest
	re := regexp.MustCompile(`([^,\[]*)(\[[^\]]*\])?(.*)`)
	match := re.FindStringSubmatch(s)
	prefix := match[1]
	rangelist := match[2]
	rest := match[3]

	rest_expanded, err := expand_part(rest)
	if err != nil {
		return nil, err
	}

	if rangelist == "" {
		us_expanded := []string{prefix}
		return combine_lists(us_expanded, rest_expanded), nil
	}

	us_expanded, err := expand_rangelist(prefix, rangelist[1:len(rangelist)-1])
	if err != nil {
		return nil, err
	}

	return combine_lists(us_expanded, rest_expanded), nil
}

func expand_rangelist(prefix, rangelist string) ([]string, error) {
	// Expand a rangelist (e.g. "1-10,14"), putting a prefix before.
	results := []string{}
	rangeParts := strings.Split(rangelist, ",")
	for _, rangePart := range rangeParts {
		rangeExp, err := expand_range(prefix, rangePart)
		if err != nil {
			return nil, err
		}
		results = append(results, rangeExp...)
	}
	return results, nil
}

func expand_range(prefix, rangeStr string) ([]string, error) {
	if regexp.MustCompile(`^[0-9]+$`).MatchString(rangeStr) {
		return []string{prefix + rangeStr}, nil
	}
	re := regexp.MustCompile(`^([0-9]+)-([0-9]+)$`)
	match := re.FindStringSubmatch(rangeStr)
	if match == nil {
		return nil, ErrBadRange
	}
	startStr := match[1]
	endStr := match[2]
	start, err := strconv.Atoi(startStr)
	if err != nil {
		return nil, err
	}
	end, err := strconv.Atoi(endStr)
	if err != nil {
		return nil, err
	}
	width := len(startStr)
	if end < start {
		return nil, ErrStartGreaterThanStop
	}
	if end-start > MAX_SIZE {
		return nil, ErrRangeTooLarge
	}
	results := []string{}
	for i := start; i <= end; i++ {
		results = append(results, fmt.Sprintf("%s%0*d", prefix, width, i))
	}
	return results, nil
}

func remove_duplicates(slice []string) []string {
	encountered := map[string]bool{}
	result := []string{}
	for _, v := range slice {
		if encountered[v] {
			continue
		} else {
			encountered[v] = true
			result = append(result, v)
		}
	}
	return result
}

func combine_lists(list1, list2 []string) []string {
	if len(list1)*len(list2) > MAX_SIZE {
		panic(ErrResultsTooLarge)
	}
	result := []string{}
	for _, part1 := range list1 {
		for _, part2 := range list2 {
			result = append(result, part1+part2)
		}
	}
	return result
}
