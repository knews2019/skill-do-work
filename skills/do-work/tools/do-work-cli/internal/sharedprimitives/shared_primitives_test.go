package sharedprimitives

import (
	"errors"
	"reflect"
	"testing"
)

// The four copies this package replaces disagreed about the empty string: two
// filtered it, two did not. The canonical helper keeps it, so a blank that ever
// reaches an evidence list is visible instead of vanishing.
func TestUniqueSortedStringsKeepsTheEmptyStringAndDropsDuplicates(t *testing.T) {
	sortedValues := UniqueSortedStrings([]string{"beta", "", "alpha", "beta", ""})
	if want := []string{"", "alpha", "beta"}; !reflect.DeepEqual(sortedValues, want) {
		t.Fatalf("UniqueSortedStrings = %#v, want %#v", sortedValues, want)
	}
}

// Callers marshal the result straight into JSON evidence, so an empty input must
// produce [] and never null.
func TestUniqueSortedStringsReturnsAnEmptySliceRatherThanNil(t *testing.T) {
	sortedValues := UniqueSortedStrings(nil)
	if sortedValues == nil || len(sortedValues) != 0 {
		t.Fatalf("UniqueSortedStrings(nil) = %#v, want a non-nil empty slice", sortedValues)
	}
}

func TestSubtractStringValuesKeepsLeftOrderAndRemovesRightMembers(t *testing.T) {
	remainingValues := SubtractStringValues([]string{"c", "a", "b", "a"}, []string{"b"})
	if want := []string{"c", "a", "a"}; !reflect.DeepEqual(remainingValues, want) {
		t.Fatalf("SubtractStringValues = %#v, want %#v", remainingValues, want)
	}
}

func TestSubtractStringValuesReturnsAnEmptySliceRatherThanNil(t *testing.T) {
	remainingValues := SubtractStringValues([]string{"a"}, []string{"a"})
	if remainingValues == nil || len(remainingValues) != 0 {
		t.Fatalf("SubtractStringValues = %#v, want a non-nil empty slice", remainingValues)
	}
}

func TestFirstNonNilErrorPrefersTheFirstCandidate(t *testing.T) {
	firstCandidate, secondCandidate := errors.New("first"), errors.New("second")
	if got := FirstNonNilError(firstCandidate, secondCandidate); got != firstCandidate {
		t.Fatalf("FirstNonNilError(first, second) = %v, want first", got)
	}
	if got := FirstNonNilError(nil, secondCandidate); got != secondCandidate {
		t.Fatalf("FirstNonNilError(nil, second) = %v, want second", got)
	}
	if got := FirstNonNilError(nil, nil); got != nil {
		t.Fatalf("FirstNonNilError(nil, nil) = %v, want nil", got)
	}
}

func TestCompareSemanticVersionsUsesTheStandardOrientation(t *testing.T) {
	for _, testCase := range []struct {
		leftVersion  string
		rightVersion string
		wantOrdering int
	}{
		{"1.0.0", "1.0.1", -1},
		{"1.0.1", "1.0.0", 1},
		{"1.2.3", "1.2.3", 0},
		{"0.9.9", "1.0.0", -1},
		{"1.10.0", "1.9.0", 1},
	} {
		ordering, parsed := CompareSemanticVersions(testCase.leftVersion, testCase.rightVersion)
		if !parsed || ordering != testCase.wantOrdering {
			t.Fatalf("CompareSemanticVersions(%q, %q) = (%d, %v), want (%d, true)",
				testCase.leftVersion, testCase.rightVersion, ordering, parsed, testCase.wantOrdering)
		}
	}
}

// The reason the comparator reports parsing separately: an unparseable version
// must not be readable as "equal", which is how one of the two deleted copies
// behaved and how a release guard could have been talked into passing.
func TestCompareSemanticVersionsReportsUnparseableInputInsteadOfScoringItEqual(t *testing.T) {
	for _, unparseableVersion := range []string{"", "1.0", "1.0.0.0", "1.0.x", "1.01.0", "v1.0.0", "1.0.0-beta"} {
		ordering, parsed := CompareSemanticVersions(unparseableVersion, "1.0.0")
		if parsed {
			t.Fatalf("CompareSemanticVersions(%q, \"1.0.0\") reported parsed", unparseableVersion)
		}
		if ordering != 0 {
			t.Fatalf("CompareSemanticVersions(%q, \"1.0.0\") ordering = %d, want 0", unparseableVersion, ordering)
		}
	}
}

func TestParseSemanticVersionAcceptsBareThreePartVersionsOnly(t *testing.T) {
	components, parsed := ParseSemanticVersion("10.0.3")
	if !parsed || components != [3]int{10, 0, 3} {
		t.Fatalf("ParseSemanticVersion(\"10.0.3\") = (%v, %v), want ([10 0 3], true)", components, parsed)
	}
	if _, parsed := ParseSemanticVersion("0.0.0"); !parsed {
		t.Fatalf("ParseSemanticVersion(\"0.0.0\") rejected a single zero component")
	}
	for _, rejected := range []string{"1.0", "1.0.0.0", "1..0", "1.0.0-rc1", "01.0.0", "1.0.0 ", "-1.0.0"} {
		if _, parsed := ParseSemanticVersion(rejected); parsed {
			t.Fatalf("ParseSemanticVersion(%q) accepted a version the strict rules reject", rejected)
		}
	}
}
