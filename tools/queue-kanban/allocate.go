package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

// requestIdNumberPattern pulls the numeric part out of a "REQ-NNN" id or a
// "REQ-NNN-slug.md" filename. Zero-padding width is irrelevant — the value is
// what orders ids, so REQ-0142 and REQ-142 are the same number.
var requestIdNumberPattern = regexp.MustCompile(`^REQ-0*(\d+)`)

// nextRequestNumber returns the next free REQ number for a repo: one past the
// highest number already in use across do-work/queue/, do-work/working/, and the
// whole do-work/archive/ subtree (including nested archive/UR-NNN/ folders), and
// 1 when the tree holds no REQ at all.
//
// This is actions/capture.md's existing allocation rule, executed instead of
// eyeballed. It reuses enumerateDoWorkTree (walk.go) — the same walk the board
// builds on, with the same pruning of deliverables/, runs/, and assets/ — rather
// than introducing a second scan that could drift from it.
//
// Gaps are deliberately fine (REQ-072 requirement 7): max+1 tolerates them and
// nothing in the skill walks a contiguous sequence, so there is no gap-filling,
// gap-detection, or compaction here and none should be added.
//
// Read-only toward the queue: nothing under do-work/ is created, moved, or
// rewritten. Allocation is also NOT atomic — two allocators running at the same
// instant can both compute the same number and both succeed, because the number
// is not reserved anywhere. That is accepted (a number-keyed reservation would be
// new durable coordination state); allocation is human-initiated and runs in
// milliseconds, which is what makes the collision window negligible rather than
// impossible.
func nextRequestNumber(repoRootOverride string) (int, error) {
	repoRoot, resolveError := resolveRepoRootOrDefault(repoRootOverride)
	if resolveError != nil {
		return 0, resolveError
	}
	discovered, walkError := enumerateDoWorkTree(repoRoot)
	if walkError != nil {
		return 0, walkError
	}

	highestNumberInUse := 0
	for _, requestFile := range discovered.RequestFiles {
		for _, candidateNumber := range requestNumbersInFile(requestFile.AbsolutePath) {
			if candidateNumber > highestNumberInUse {
				highestNumberInUse = candidateNumber
			}
		}
	}
	return highestNumberInUse + 1, nil
}

// requestNumbersInFile returns every REQ number a single file lays claim to: the
// one in its filename and, when the frontmatter carries a different `id:`, that
// one too. Both count — a file renamed away from its id still owns its id's
// number, and handing that number out again would manufacture the duplicate this
// tool is meant to detect.
func requestNumbersInFile(absolutePath string) []int {
	var claimedNumbers []int
	if fileNameNumber, parsedOk := requestNumberFromText(baseNameOf(absolutePath)); parsedOk {
		claimedNumbers = append(claimedNumbers, fileNameNumber)
	}
	if frontmatterNumber, parsedOk := requestNumberFromFrontmatterId(absolutePath); parsedOk {
		claimedNumbers = append(claimedNumbers, frontmatterNumber)
	}
	return claimedNumbers
}

// requestNumberFromText parses the numeric part of a REQ id or REQ filename.
func requestNumberFromText(text string) (int, bool) {
	match := requestIdNumberPattern.FindStringSubmatch(strings.TrimSpace(text))
	if match == nil {
		return 0, false
	}
	parsedNumber, parseError := strconv.Atoi(match[1])
	if parseError != nil {
		return 0, false
	}
	return parsedNumber, true
}

// requestNumberFromFrontmatterId reads just the frontmatter `id:` field. It uses
// the existing lenient frontmatter reader so a REQ with slightly malformed YAML
// still contributes its number instead of silently dropping out of the max.
func requestNumberFromFrontmatterId(absolutePath string) (int, bool) {
	fileBytes, readError := os.ReadFile(absolutePath)
	if readError != nil {
		return 0, false
	}
	yamlText, _, hasFrontmatter := splitFrontmatter(string(fileBytes))
	if !hasFrontmatter {
		return 0, false
	}
	fields := lenientFrontmatterFields(yamlText)
	return requestNumberFromText(coerceScalarToString(fields["id"]))
}

// baseNameOf returns the last path element. Kept local and tiny so allocate.go
// does not pull filepath in for one call.
func baseNameOf(path string) string {
	if slashIndex := strings.LastIndexAny(path, `/\`); slashIndex >= 0 {
		return path[slashIndex+1:]
	}
	return path
}
