package toolboxcommands

import "testing"

func TestParseMutationFlags(t *testing.T) {
	rest, dry, commit, err := parseMutationFlags([]string{"--dry-run", "x"})
	if err != nil || !dry || commit || len(rest) != 1 || rest[0] != "x" {
		t.Fatalf("unexpected parse: %v %v %v %v", rest, dry, commit, err)
	}
	if _, _, _, err := parseMutationFlags([]string{"--dry-run", "--commit"}); err == nil {
		t.Fatal("combined flags accepted")
	}
}
