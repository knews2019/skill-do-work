---
id: REQ-279
title: Scope the blanked-REQ scan to where REQ records actually live
status: completed
claimed_at: 2026-08-20T22:38:28Z
completed_at: 2026-08-20T22:42:34Z
commit: f487f04
created_at: 2026-08-19T13:42:45Z
user_request: UR-057
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-280, REQ-281, REQ-282, REQ-283]
batch: upstream-consumer-report-2026-08-19
effort_estimate: trivial
write_set:
- skills/do-work/tools/checks/blanked-req-scan.sh
- _dev/tests/record-commit-hash-guards.sh
---

# Scope the Blanked-REQ Scan to Where REQ Records Actually Live

## What

`skills/do-work/tools/checks/blanked-req-scan.sh:285` enumerates candidates with an unscoped `find do-work/archive -type f \( -name 'REQ-*.md' -o -name 'UR-*.md' \)`. That reaches `do-work/archive/UR-NNN/assets/`, where the capture flow parks screenshot-description sidecars named after the REQ they illustrate. Those are prose documents and correctly carry no frontmatter, so `has_parseable_frontmatter` (line 88) classifies each one as destroyed content, `resolve_recovery_source` finds no version with frontmatter because there never was one, and the scanner reports permanent data loss on an intact file.

Exclude `assets/` at any depth from the enumeration.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Narrow only the archived arm of the scanner's candidate `find` with an `assets/` path exclusion, update Check 13's scope sentence, and add a fixture proving asset prose is ignored while a real blanked REQ remains visible.
- [x] **[APPLY]:** Code written exactly as planned. Scope stayed within the scanner, Check 13, and the behavioral fixture.
- [x] **[UNIFY]:** Reviewed the 3-file diff and git diff --stat; ran bash -n on the changed shell files and the behavioral probe. No debug artifacts are present.

## Why

`actions/forensics.md` Check 13 reports each finding as **Critical**, and `actions/cleanup.md` Pass 6 then asks the user to approve a restore that "overwrites each file's content with the recovered blob," under a `git gc` urgency warning that pressures toward yes. For an asset sidecar the recovered blob is nothing, and the scanner's own refusal to write empty content is the only thing between the user and a deleted file. In the reporting consumer repo every finding was an `assets/` sidecar — four of them, all intact, with no true positive to hide behind.

## Context

The suite already fixed this exact class once, in the other tool. `skills/do-work-board/tools/queue-kanban/walk.go:195` prunes any directory named `assets` at any depth, and `isSkippedSection`'s comment records why: "a deliverable copy named `REQ-NNN-*.md` under `assets/` has no frontmatter, so its id falls back to the filename and collides with the real `REQ-NNN` ticket." `CHANGELOG.md` 0.150.12 is the entry. Same file shape, same wrong conclusion, different tool.

This repo has eight `assets/` directories but none holds a `REQ-*.md` (only `.png`), which is why the scan exits 0 here and the defect never surfaced in maintainer verification.

## Detailed Requirements

- Restrict the `find` at `blanked-req-scan.sh:285` so `do-work/archive/**/assets/**` cannot produce a candidate. Prefer `-not -path '*/assets/*'` over a `-maxdepth` cap.
- Do not change the `do-work/queue` / `do-work/working` arm of the enumeration — neither has an `assets/` subtree, and narrowing it would be scope the finding did not earn.
- Leave `has_parseable_frontmatter`, `resolve_recovery_source`, `resolve_recorded_hash`, the restore path, and every exit code untouched.
- Update the script's header block if it states where it walks, and `actions/forensics.md` Check 13's sentence "walks `do-work/archive/`, `do-work/queue/`, and `do-work/working/` for `REQ-*.md` and `UR-*.md` files" so the documented scope matches the code.

## Constraints

- **Use the prune form, not `-maxdepth 2`.** The upstream report preferred the depth cap, arguing it states where records live rather than denylisting a known subdir. Take the prune anyway: the suite already carries one statement of "a `REQ-*.md` under `assets/` is not a record" in `walk.go`, and a second, differently-shaped rule in the scanner is exactly the drift CLAUDE.md's *state conditions, not lists* guidance targets. A depth cap also encodes a layout assumption that a future nesting change breaks silently, in the one script whose failure mode is proposing an overwrite.
- **Declined sub-remedy — the `sh` shebang note.** The report asked for a shebang-honoring note in Check 13 and Pass 6. Do not add it. The script is `#!/usr/bin/env bash`, mode 0755, and both documented call sites (`actions/forensics.md:170`, `actions/cleanup.md:148`) invoke it bare, which honors the shebang by construction. Running it through `sh` by hand is a caller error.
- **Deferred sub-remedy — splitting the two damage states.** The report also asked to separate "no frontmatter and no version ever had one" from "truncated, git holds the content." Out of scope here: the scanner already emits `-` as the recovery SHA for the first, `actions/cleanup.md` Pass 6 step 5 already handles it as reported-not-restored, and once `assets/` is excluded the residue is small. If the builder finds it cheap while in the file, raise it as a discovered task rather than widening this REQ.
- Write-set overlap: queued REQ-276 also edits `blanked-req-scan.sh` (line 91's closing-fence guard). Different concern, same file.

## Red-Green Proof

**RED prompt/case:** In the `blanked-req-scan.sh` fixture at `_dev/tests/record-commit-hash-guards.sh:381`, add a committed, intact prose file at `do-work/archive/UR-001/assets/REQ-001-screenshot-descriptions.md` with no frontmatter. Run the scanner.
**Why RED now:** It exits 1 and prints that file as damaged with `-` as its recovery source — a report of permanent data loss against a file that was never damaged, which `actions/cleanup.md` Pass 6 then offers to overwrite.
**GREEN when:** The same fixture exits 0 with "No blanked or unparseable REQ/UR files found," and a genuinely blanked `do-work/archive/UR-001/REQ-002-real.md` in the same tree is still reported. Both assertions live in the test, so the fix cannot be widened into blindness.
**Validation:** Inferred during capture from the upstream report's reproduce block; the sidecar shape is confirmed against this repo's `assets/` layout and the consumer's four reported paths.

## Full Context

See `do-work/user-requests/UR-057/input.md` for the complete verbatim upstream report.

---
*Source: upstream defect report D2, severity high, from `g1w-game-find-the-difference` running v0.212.25 — verbatim claim: "`blanked-req-scan.sh` reports intact asset sidecars as destroyed content … it invites a destructive repair against undamaged files." Accepted by `do-work-toolbox validate-feedback` triage (2026-08-19). Evidence: `skills/do-work/tools/checks/blanked-req-scan.sh:285` unscoped find; `:88` `has_parseable_frontmatter` is the damage test; `skills/do-work/actions/forensics.md:174` reports Critical; `skills/do-work/actions/cleanup.md:157` Pass 6 offers the overwrite; precedent fix at `skills/do-work-board/tools/queue-kanban/walk.go:195` and `CHANGELOG.md` 0.150.12. Surface-cost: N/A — narrowing an over-broad `find`, a direct fix.*

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names the scanner, its existing behavioral fixture, and one documentation site. The fix is a focused path predicate plus a regression assertion.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/checks/blanked-req-scan.sh` (modified) — excludes every `assets/` subtree only from the archive candidate scan.
- `skills/do-work/actions/forensics.md` (modified) — documents the archive `assets/` exclusion in Check 13.
- `_dev/tests/record-commit-hash-guards.sh` (modified) — verifies intact prose under archive assets is ignored while a genuinely blanked REQ remains reported.

## Qualification

Passed — 3 project files verified, all requirements traced, and the P-A-U change stayed within the declared write set.

## Testing

- `record-commit-hash-guards.sh` behavioral probes passed through the maintainer harness path, including the new archive-assets regression fixture.
- `bash -n skills/do-work/tools/checks/blanked-req-scan.sh _dev/tests/record-commit-hash-guards.sh` passed.
- The captured RED case is the intact no-frontmatter asset sidecar being reported as damaged; the GREEN assertion confirms it is ignored while the real blanked REQ still produces exit 1.

## Review

**Overall: 100%**

**Requirements:** Pass — the archive arm excludes `assets/` at any depth, the queue/working arm is unchanged, the scanner logic and exit codes are untouched, and Check 13 matches the runtime scope.

**Code quality:** Pass — the change is a single focused `find` predicate and a targeted fixture assertion.

**Test adequacy:** Pass — the regression covers both the false positive and the true positive.

**Scope:** Pass — only the two declared deliverables and their shared behavioral fixture changed.

**Acceptance:** Pass — the scanner probes pass and no additional testing is needed.

## Lessons Learned

**What worked:** The existing shared guard-probe fixture could exercise the scanner against a committed archive shape without adding another test runner.
**What didn't:** A direct invocation of the probe uses its maintainer-tree path assumption; the repository's harness substitution is required for the shipped script path.
**Worth knowing:** Archive assets are deliverables, not REQ records, even when their filenames resemble REQ files.

## Orientation

The blanked-record detector now reports only queue records and archive records outside `assets/`, while preserving its existing damage detection and recovery behavior.
