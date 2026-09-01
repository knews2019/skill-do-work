---
id: REQ-440
title: '[impact-critical] Refuse non-file static board output targets'
status: completed
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
route: A
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-09-01T18:56:32Z
  basis:
    - trivial short-circuit
related: [REQ-437, REQ-438, REQ-439, REQ-441, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
claimed_at: 2026-09-01T18:55:29Z
completed_at: 2026-09-01T19:45:42Z
---

# Refuse Non-File Static Board Output Targets

## What

Refuse static-board publication when any generated output name already exists as a directory, symlink, or other non-regular object. Validate all three targets before the first backup rename so successful cleanup can never recursively delete user content hidden beneath an output filename.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Add a caller-seam table regression in `generate_test.go` that places a directory, symlink, or named pipe at one of the three generated target names, invokes static generation, and proves refusal preserves every pre-existing target and leaves no private publication directory. Then add one fixed-scope `Lstat` preflight in `generate.go` before private staging or backup renames; absent targets and regular non-symlink files continue into the existing publication and rollback path. Verify with the focused RED/GREEN test, the existing publication tests, the package suite, formatting/vetting, and a scoped diff review.
- [x] **[APPLY]:** Added the caller-seam directory/symlink/named-pipe regression and the planned fixed three-target `Lstat` preflight before private staging. Root cause: the backup loop treated every successful `Lstat` as replaceable without checking the inode type. The existing backup, publication, rollback, and regular-file replacement paths are unchanged.
- [x] **[UNIFY]:** Reviewed the scoped 167-line source/test diff and this P-A-U edit. `generate.go` contains only the fixed three-target preflight; `generate_test.go` covers all three non-regular kinds, no publisher calls, preserved targets, and no residue. `gofmt -d` and `git diff --check` were clean; the focused regression, existing publication/rollback tests, full Go package suite, and `go vet ./...` passed. No debug artifacts or unrelated source changes were present.

## Finding Provenance

- **Verbatim claim / severity:** `[P1] Refuse non-file static output targets before backing up.`
- **Evidence:** `generate.go:479-504` accepts any successful `Lstat`, renames the inode below the private directory, publishes a file, and later removes the private directory recursively.
- **Origin / earned by:** `803e4e77`/REQ-183 (Static Board Generation Can Publish a Mixed Three-File Bundle) introduced recovery backups but treated every existing object as replaceable. An isolated replay replaced an `index.html` directory and deleted its nested file while returning success.
- **Surface-cost:** Earned. The exact data-loss replay justifies one preflight over three fixed targets plus directory/symlink regressions; this is much cheaper than deleting arbitrary directory contents.

## Detailed Requirements

- Inspect all generated target paths before moving any existing target.
- Permit only absent targets and existing regular non-symlink files.
- Refuse directories, symlinks, and special files without mutating any target or leaving private publication residue.
- Preserve current regular-file publication and rollback behavior.

## Constraints

- Keep the existing three-output transactional publication shape.
- The validation is fixed-scope and must not become a configurable filesystem policy layer.

## Red-Green Proof

**RED prompt/case:** Replace `index.html` with a directory containing `kept.txt`, run static generation, and inspect the target and private scratch paths.
**Why RED now:** Publication moves the directory into its private backup and successful cleanup recursively deletes it.
**GREEN when:** Generation refuses before any rename, the directory and nested bytes remain unchanged, no private residue remains, and ordinary regular-file regeneration still succeeds.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. A preflight regular-file/no-symlink check before the first rename is the accepted remedy.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---

## Triage

**Route: A** - Simple

**Reasoning:** The destructive behavior, affected static-publication path, accepted fix, and regression replay are explicit. The smallest correct change is a focused preflight plus caller-seam tests.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-kanban-board.md` — 4707 tokens; matches static board publication, but the partial satellite cannot be narrowed below the 2000-token budget. Read anyway under the prime's touch-conditional Lessons rule.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 5083 tokens; matches queue-kanban static output, but the partial satellite cannot be narrowed below the 2000-token budget. Read anyway because `generate.go` is named by the utility prime's Read-first section.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified)

**What was done:** Static publication now preflights all three output targets with `Lstat` before creating private staging or renaming backups, refusing every existing non-regular object. A caller-seam regression covers directory, symlink, and named-pipe targets while preserving existing regular-file replacement and rollback behavior.

## Qualification

Passed — 2 files verified, 4 requirements traced, P-A-U confirmed. The fixed-scope preflight runs before private staging and any rename; the caller-seam test proves refusal across three non-regular object classes, no publisher invocation, preserved target bytes/types, and no scratch residue.

## Testing

**Tests run:** `go test -run '^TestGenerate(FirstPublicationAndSuccessfulReplacement|PublicationFailureRestoresThePreviousBundle|RefusesNonRegularOutputTargetsWithoutMutation)$' -count=1`, `go test ./...`, `go vet ./...`, and `bash _dev/tests/maintainer-verify.sh`
**Result:** Focused publication tests, the full queue-kanban module suite, and vet passed. The canonical maintainer gate first failed on a pre-existing ShellCheck `SC2034` in `_dev/tests/shipped-shell-thinness.sh:9`, outside this REQ's diff; that warning was fixed and committed separately as `2d140f63` (0.260.2), after which `bash _dev/tests/maintainer-verify.sh` exited 0 against the working tree containing this REQ's changes (2026-09-01).

**Red-green validation:**
- `TestGenerateRefusesNonRegularOutputTargetsWithoutMutation`: ✗ before implementation — directory, symlink, and named-pipe targets were replaced and the publisher ran three times → ✓ after implementation — publication refused before any publisher call and all targets remained unchanged.

**New tests added:**
- `TestGenerateRefusesNonRegularOutputTargetsWithoutMutation` in `generate_test.go`

*Verified by work action. The claim was held at this step on 2026-09-01 while the canonical gate was red for the unrelated warning above, then resumed in place once the gate was green; no generated sections were stripped.*

---
*Source: accepted Finding 19 from the validated external feedback.*

## Review

**Overall: 96%** | 2026-09-01T19:43:22Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 3 (report only) — stale doc comment above `publishStaticSiteOutputs` (routed to `do-work/prose-backlog.md`); the backup loop's own `Lstat` still lacks an `IsRegular` check inside the check-to-rename window that the REQ's Constraints leave by design; the regression asserts an error but not its text.
**Acceptance:** Pass — focused publication tests, full package suite, and `go vet` all green; the named regression was reproduced RED against `git show HEAD:generate.go` in an isolated copy, destroying the nested fixture exactly as the finding claimed.
**Suggested testing:** 2 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Orientation

Static board generation can no longer delete a directory, symlink, or special file that sits at one of its three output names; lives in the Kanban board tool's static publication path (`_dev/primes/prime-kanban-board.md`). Leaf change, map unchanged.

## Discovered Tasks

- **impact-critical** `do-work-cli` reads an empty inline frontmatter list (`depends_on: []`) as one dependency literally named `[]`, so the canonical `next` selector excluded 21 of 30 pending REQs with `DEPENDENCY-MISSING: missing dependencies: []` on 2026-09-01. Root cause: `RequestDocument.FieldValue` in `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` copies `ListValues` with `append([]string(nil), ...)`, which turns an empty slice into `nil`, and `listValue` then falls back to the scalar text `[]`.
