// Package sharedprimitives holds the small, dependency-free helpers that several
// internal packages had each grown their own copy of. It imports nothing else
// inside this module on purpose: a leaf package can be called from anywhere
// without an import cycle, which is why the duplicated copies existed at all.
package sharedprimitives

import (
	"sort"
	"strconv"
	"strings"
)

// UniqueSortedStrings returns the distinct values of the input in ascending
// order. It keeps every value it is given, including the empty string: two of
// the copies this replaces filtered blanks, and a blank that reaches an evidence
// list should be visible there rather than silently dropped. The result is
// always non-nil so a caller marshalling it emits [] rather than null.
func UniqueSortedStrings(values []string) []string {
	distinctValues := make(map[string]bool, len(values))
	for _, value := range values {
		distinctValues[value] = true
	}
	sortedValues := make([]string, 0, len(distinctValues))
	for value := range distinctValues {
		sortedValues = append(sortedValues, value)
	}
	sort.Strings(sortedValues)
	return sortedValues
}

// SubtractStringValues returns the left values that do not appear among the
// right values, in their original order and with duplicates preserved. The
// result is always non-nil.
func SubtractStringValues(leftValues, rightValues []string) []string {
	rightSet := make(map[string]bool, len(rightValues))
	for _, value := range rightValues {
		rightSet[value] = true
	}
	remainingValues := []string{}
	for _, value := range leftValues {
		if !rightSet[value] {
			remainingValues = append(remainingValues, value)
		}
	}
	return remainingValues
}

// FirstNonNilError returns the first candidate when it is set, otherwise the
// second. It exists so a caller holding two independent errors can report one
// without an if-chain at every site.
func FirstNonNilError(firstCandidate, secondCandidate error) error {
	if firstCandidate != nil {
		return firstCandidate
	}
	return secondCandidate
}

// CompareSemanticVersions orders two bare semantic versions in the standard Go
// orientation: negative when the left version is older, positive when it is
// newer, zero when the three components match. The second result reports
// whether BOTH inputs parsed. It is returned separately, and not folded into the
// ordering, because a comparator that scored an unparseable version as "equal"
// is the exact drift this helper was written to remove: a caller must not be
// able to read "I could not parse this" as "these are the same version".
func CompareSemanticVersions(leftVersion, rightVersion string) (int, bool) {
	leftComponents, leftParsed := ParseSemanticVersion(leftVersion)
	rightComponents, rightParsed := ParseSemanticVersion(rightVersion)
	if !leftParsed || !rightParsed {
		return 0, false
	}
	for index := range leftComponents {
		if leftComponents[index] < rightComponents[index] {
			return -1, true
		}
		if leftComponents[index] > rightComponents[index] {
			return 1, true
		}
	}
	return 0, true
}

// ParseSemanticVersion reads a bare major.minor.patch version. It is strict on
// purpose: exactly three dot-separated parts, no empty part, no leading zero on
// a multi-character part, digits only, and no pre-release or build suffix. A
// version that fails any of those is reported as unparsed rather than coerced.
func ParseSemanticVersion(value string) ([3]int, bool) {
	var components [3]int
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return components, false
	}
	for index, part := range parts {
		if part == "" || len(part) > 1 && part[0] == '0' {
			return components, false
		}
		parsed, parseError := strconv.Atoi(part)
		if parseError != nil || parsed < 0 {
			return components, false
		}
		components[index] = parsed
	}
	return components, true
}
