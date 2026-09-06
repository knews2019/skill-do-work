package finalization

import (
	"reflect"
	"testing"
)

// REQ-557 / D-02: the helper that computed the missing commit paths used to be
// `uniqueSorted`, whose whole body was `result, _ := normalizeRepositoryPaths(paths);
// return result`. Throwing the validator's error away made the guard on the next line
// unreachable: one required path that is empty, absolute, or escaping turned the WHOLE
// required set into nil, the subtraction then found nothing missing, and prepare
// continued as if commit_paths covered every planned target.
//
// This test fails the moment the error is discarded again, whatever the helper is called.
func TestMissingCommitPathsRefusesAnUnusableRequiredPathInsteadOfEmptyingTheSet(t *testing.T) {
	for _, unusableRequiredPath := range []string{"", "/etc/passwd", "../outside.md"} {
		requiredCommitPaths := []string{"do-work/archive/REQ-557.md", unusableRequiredPath}

		missing, missingError := missingCommitPaths(requiredCommitPaths, nil)

		if missingError == nil {
			t.Fatalf("missingCommitPaths accepted the unusable required path %q and returned missing=%#v; "+
				"the commit_paths guard is disabled again", unusableRequiredPath, missing)
		}
		if missing != nil {
			t.Fatalf("missingCommitPaths(%q) returned both an error and %#v", unusableRequiredPath, missing)
		}
	}
}

// The same helper on a usable set still answers the question the guard asks: which
// planned lifecycle or release targets commit_paths does not cover.
func TestMissingCommitPathsNamesTheRequiredTargetsCommitPathsOmits(t *testing.T) {
	missing, missingError := missingCommitPaths(
		[]string{"do-work/CHECKPOINT.md", "do-work/archive/REQ-557.md", "do-work/CHECKPOINT.md"},
		[]string{"do-work/CHECKPOINT.md"},
	)
	if missingError != nil {
		t.Fatalf("missingCommitPaths refused a usable set: %v", missingError)
	}
	if want := []string{"do-work/archive/REQ-557.md"}; !reflect.DeepEqual(missing, want) {
		t.Fatalf("missingCommitPaths = %#v, want %#v", missing, want)
	}
}
