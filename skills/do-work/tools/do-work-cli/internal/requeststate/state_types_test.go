package requeststate

import "testing"

func TestCommitHashValidation(t *testing.T) {
	for _, value := range []string{"1234567", "abcdef0123456789", "0123456789012345678901234567890123456789"} {
		if !validCommitHash(value) {
			t.Errorf("valid hash rejected: %q", value)
		}
	}
	for _, value := range []string{"123456", "ABCDEF0", "not-a-hash", "01234567890123456789012345678901234567890"} {
		if validCommitHash(value) {
			t.Errorf("invalid hash accepted: %q", value)
		}
	}
}

func TestStateOptionsRejectUnknownProvenance(t *testing.T) {
	if _, err := parseStateOptions(TransitionClaim, []string{"REQ-001", "--provenance", "invented"}); err == nil {
		t.Fatal("unknown selection provenance was accepted")
	}
}
