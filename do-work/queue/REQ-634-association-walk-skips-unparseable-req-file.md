---
id: REQ-634
title: 'A REQ file the association walk cannot parse claims no paths, and do-work commit continues'
status: pending
created_at: 2026-10-06T15:27:35Z
user_request: UR-136
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md, _dev/primes/prime-shell-commands.md, _dev/primes/prime-action-files.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go, skills/do-work/tools/do-work-cli/internal/corehelpers/inventory_test.go, _dev/tests/contracts/core-checks.sh, skills/do-work/actions/commit.md, skills/do-work-toolbox/actions/inspect.md, skills/do-work/docs/prescribed-shell-primitives.md, skills/do-work/CHANGELOG.md, skills/do-work/tools/do-work-cli/lessons-do-work-cli.md, do-work/lessons-index.md]
---
# A REQ File the Association Walk Cannot Parse Claims No Paths, and do-work Commit Continues

## What
When `AssociateProjectPaths` (`skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go`) hits a REQ file whose Implementation Summary fails to parse, that one file claims no paths, the walk continues, and `protected-inventory associate` exits 0 with its owner rows on stdout and one warning on stderr naming the skipped file. The `PARSE-FAILED` exit-2 path goes away. Nothing about how a path is written changes.

## Why
A consumer reported (do-work 0.305.68) that one archived REQ with an odd number of backticks on one Summary line makes every `do-work commit` in a ~2,000-REQ repository exit 2 at Step 3, permanently, because archive records are immutable. The walk at `inventory.go:301-304` returns the parse error from the WalkDir callback, so a single-file problem becomes a walk-level failure. The maintainer's position, recorded at capture: the commit must not die on a formatting issue in one record, and the parser must not grow a grammar for every path representation a non-deterministic model can produce. So the fix is tolerance at the walk, not a smarter parser.

## Detailed Requirements
1. In `AssociateProjectPaths`, a parse error from `allBacktickedPaths` on one REQ file means that file claims no paths. Return nil from the callback for that file and keep walking. The function returns the skipped files alongside the association map so the caller can report them.
2. `handleAssociate` (`inventory.go:213-218`) drops the `PARSE-FAILED` branch and emits one warning finding per skipped file, code `ASSOCIATION-SUMMARY-UNPARSED`, with the REQ file's repository-relative path as the affected path and the parser's message as evidence. Outcome stays success, exit 0. In shim text mode the stdout rows (`<owner>\t<path>`) stay exactly as they are; the warning reaches stderr through the runtime's existing finding printer (`internal/commandruntime/command_runtime.go:95-104`), not through a direct `os.Stderr` write in the handler.
3. Flip the two lock-ins that pin the old behaviour: the `associate_unmatched` probe in `_dev/tests/contracts/core-checks.sh` (around lines 365-389) and the `preserves PARSE-FAILED` subtest in `inventory_test.go` (around lines 636-660). Both now assert: exit 0, the good file's owner row present on stdout, `PARSE-FAILED` absent, and the stderr finding naming the unparseable REQ file. Add a unit test on `AssociateProjectPaths` with one unparseable archived REQ and one claiming archived REQ: the claiming REQ's path is associated and the unparseable file is reported once.
4. Rewrite the sentence at `skills/do-work/actions/commit.md:85` and `skills/do-work-toolbox/actions/inspect.md:117` so `PARSE-FAILED` is no longer listed as an exit-2 case; say instead that a REQ file whose Summary cannot be parsed claims nothing and is named on stderr, and the agent reports it. Add one bullet under "What `associate` settles" in `skills/do-work/docs/prescribed-shell-primitives.md` stating the same rule.
5. Release: changelog entry and version bump per `_dev/primes/prime-releases.md` (shipped Go, action prose, and the docs guide all change).
6. Lessons: one entry in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (a per-record parse error inside a lookup walk must not become a walk-level failure) and the matching token refresh in `do-work/lessons-index.md`, in the same commit.

## Constraints
- No bullet-joining, no continuation-line grammar, no "first span is the path" rule. `allBacktickedPaths`, `firstBacktickedPaths` and `qualificationSummaryEntries` are untouched; qualify (`checks.go:321`) and scope-drift (`checks.go:152`) stay strict because there the REQ being parsed is the one being finalized and the author can fix it.
- The multi-path contract stays: `TestAssociationParserRetainsEveryClosedToken` and the core-checks.sh multi-path probe are unchanged.
- A skipped REQ's own files become unassociated and go to commit Step 4 grouping. Accepted.
- The reporter's optional in-flight marker (`ASSOCIATION-FOUND-INFLIGHT`) is dropped. In-flight `working/` claims stay included as documented at `inventory.go:257-263`; no marker, no confirm step.
- The reporter asked for no JSON code. The finding is the CLI's only stderr channel, so it also appears under `--format json`. The maintainer accepted this at plan approval.

## Builder Guidance
Certainty is high. The code change is small: the callback swallows the error into a skipped list, the handler turns the list into warning findings, one branch is deleted. Most of the work is the two test inversions, three prose sites, the changelog, and the lessons entry. Builder latitude: the exact evidence wording of the stderr line, and how `AssociateProjectPaths` returns the skipped list (second return value or a small result struct).

## Red-Green Proof
**RED prompt/case:** A fixture repository with an uncommitted `good-file.txt`, an archived completed REQ-502 whose Summary bullet is the existing core-checks.sh unmatched probe body (one closed `legacy-file.txt` span, then a second span opened with a backtick and never closed), and an archived completed REQ-503 claiming `good-file.txt`. Run `scripts/protected-inventory.sh start` then `scripts/protected-inventory.sh associate`.
**Why RED now:** `associate` exits 2, prints `PARSE-FAILED: unmatched backtick in Implementation Summary`, and prints no owner row for `good-file.txt`.
**GREEN when:** `associate` exits 0, stdout contains the row `REQ-503<TAB>good-file.txt`, stderr contains one finding line naming the REQ-502 file, and `PARSE-FAILED` appears nowhere.
**Validation:** User confirmed (2026-10-06, decisions D1 to D4 at capture)

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18027 tokens, over the 2000 budget; `slugged: partial`). Matched: this REQ changes a corehelpers classifier walk and the finding a command emits.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over the 2000 budget; `slugged: partial`). Matched: a prescribed command block's documented exit contract changes in commit.md and inspect.md.
- `_dev/primes/lessons-action-files.md` as a whole satellite (5879 tokens, over the 2000 budget; `slugged: partial`). Matched: action prose in commit.md and inspect.md changes.

## Full Context
See `do-work/user-requests/UR-136/input.md` for complete verbatim input.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: consumer bug report pasted into `do-work validate-feedback` on 2026-10-06, triaged with the maintainer; the full text is the UR's verbatim input.*
