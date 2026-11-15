---
id: REQ-732
title: '[impact-negligible] Allow docs/run-prompts/ in the staged-skills contract and move the run prompts there'
status: completed
created_at: 2026-11-14T22:47:58Z
user_request: UR-176
domain: general
prime_files: []
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
related: [REQ-726, REQ-727, REQ-728, REQ-729, REQ-730, REQ-731]
batch: tinka-follow-ups
write_set: ["_dev/tests/staged-skills-contract.sh", ".gitattributes", "docs/run-prompts/*.md", "_dev/run-prompts/*.md"]
claimed_at: 2026-11-14T23:28:20Z
route: A
required_lessons: ["_dev/primes/lessons-shell-commands.md#assertion-passes-when-feature-is-dead"]
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-11-14T23:31:31Z
  basis:
    - Route A
    - 8-file write set
    - 4 acceptance criteria
dispatch_at: 2026-11-14T23:33:30Z
builder_handback_at: 2026-11-15T02:33:50Z
integration_at: 2026-11-15T02:34:06Z
review_at: 2026-11-15T02:40:30Z
kb_status: pending
commit: ce3cb5dc3b4099574b9d5ca103fd6cffc426957a
heavy_verified_at: 2026-11-15T02:45:34Z
heavy_verified_revision: ce3cb5dc3b4099574b9d5ca103fd6cffc426957a
completed_at: 2026-11-15T02:45:52Z
---
# Allow docs/run-prompts/ in the Staged-Skills Contract and Move the Run Prompts There
## What
Let the staged-skills contract accept `docs/run-prompts/` under the root `docs` folder, keep that folder out of the shipped package, and move the prompt files from `_dev/run-prompts/` there.
## Why
The maintainer's global rule says a prompt for another session goes to `docs/run-prompts/<date>-<slug>.md`. In this repository that file broke a check, so the rule could not be followed here.
## Context
- `_dev/tests/staged-skills-contract.sh` (the loops near lines 75 and 99 to 114, message "legacy root runtime must be retired at modular cutover") allows only `docs/kb/`, `docs/decisions/` and `docs/changelog/` under the root `docs` folder.
- On 2026-11-14 the file `docs/run-prompts/2026-11-14-tinka-canvas-suggestion.md` turned the heavy lane `staged-skills` red five commits after it was added, and was removed in commit 00568685. The fast gate did not catch it.
- `.gitattributes` marks the three allowed folders `export-ignore`, and the contract asserts those lines.
- REQ-724 (move the maintainer's global standing rules into the skill) is queued and names this question as open. It does not depend on this request.
## Detailed Requirements
- Add `docs/run-prompts` to both lists in the contract: the export-ignore assertion and the allowed paths of the leftover scan.
- Add the matching `export-ignore` line to `.gitattributes`, with one comment line in the style of its neighbours.
- Move every file from `_dev/run-prompts/` to `docs/run-prompts/` with `git mv`, and remove the empty folder.
- Sweep the repository for links or paths that name `_dev/run-prompts/` in files that are not archived records, and repoint them. Archived request records keep their text.
## Constraints
- Nothing under `skills/` or `suite/` changes, so this is no release.
- The contract keeps refusing every other new file under the root `docs` folder.
## Dependencies
None. The other requests of this batch can run in any order.
## Red-Green Proof
**RED prompt/case:** With a file `docs/run-prompts/probe.md` present, run `bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`.
**Why RED now:** The lane fails with the legacy-root message that names the probe file.
**GREEN when:** The same lane passes with that file present and still fails with a file `docs/other/probe.md`; both probe files are removed afterwards; the moved prompts are in `docs/run-prompts/`.
**Validation:** Inferred during capture. The proof comes from the prompt file the maintainer asked to capture; it was not confirmed as a separate question.
## Open Questions
- [x] Where do run prompts live in this repository: `docs/run-prompts/` or `_dev/run-prompts/`? → Allow docs/run-prompts/ (Recommended).
  Note (2026-11-14): The maintainer chose to make the global rule work here as in every other project. The accepted cost: one more allowed folder under `docs` to keep out of the shipped package. Keeping `_dev/run-prompts/` with a line in `CLAUDE.md` is out of scope.
<!-- D-XX counter: last used D-03. Next decision: D-04. -->
## Builder Guidance
- Certainty: firm. The proof uses the heavy lane `staged-skills`, which the normal heavy drain of a run executes.
- Scope: the smallest change. Two list entries in the contract, one line in `.gitattributes`, and the moves.
## Full Context
See `do-work/user-requests/UR-176/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Two list entries plus a comment in `_dev/tests/staged-skills-contract.sh`, one `export-ignore` line plus a comment in `.gitattributes`, `git mv` of the six prompts, one commit (from the builder hand-back).
- [x] **[APPLY]:** Done as planned in builder commit `342a36a6`, plus `rmdir` of the empty `_dev/run-prompts` folder in the worktree. At the main-checkout merge git removed the empty folder itself, so no `rmdir` was needed there (from the builder hand-back, integrator note added).
- [x] **[UNIFY]:** `git show --stat HEAD`: 8 files, 15+/9-. The `.gitattributes` and contract hunks were reviewed line by line, the six renames show 0 changed lines, no probe files or debug artifacts left (from the builder hand-back).
*Source: do-work capture-request: capture the request in _dev/run-prompts/2026-11-15-run-prompts-folder-rule.md*

## Triage

**Route: A** - Simple

**Reasoning:** The request names every file and every change: two list entries in `_dev/tests/staged-skills-contract.sh`, one `export-ignore` line in `.gitattributes`, and a `git mv` of the six files in `_dev/run-prompts/`. The repository sweep is already done (below) and finds nothing to repoint outside records.

**Planning:** Not required

**Pre-dispatch decisions (unattended run, no question asked):**
- **D-01**: Records keep their text. Every remaining mention of `_dev/run-prompts/` sits in `do-work/` records (UR inputs, REQ `Source:` lines, sibling REQs REQ-724 to REQ-731) or inside the moved prompts themselves, which are dated historical prompts. The builder moves the prompts byte for byte (a pure rename) and edits no `do-work/` path. Value: git records a clean rename and no record is rewritten. Risk: low; a reader of an old record follows `_dev/run-prompts/` to a missing folder, and `git log --follow` finds the file.
- **D-02**: The GREEN probe runs the contract's own docs checks, not the whole contract. The full contract took 39 s at base under load from ten parallel sessions, over the 30-second probe budget. The probe cuts the export-ignore loop and the legacy-root scan out of the real script by their text markers and runs them against a temporary fixture tree, so it tests both `docs/run-prompts/probe.md` (must pass) and `docs/other/probe.md` (must fail) without writing into the checkout. The builder still runs the whole contract once in its worktree, and the integrator's heavy lane runs it. Value: a fast, exact probe of the changed lines. Risk: if the builder renames the marker lines (`for maintainer_docs_path in`, `board_files=(`), the probe fails closed and says so.
- **D-03**: Scope stays at the request. The heavy-lane `staged-skills` coverage in `_dev/tests/heavy-lanes.json` does not list root `docs/`, which is why the fast gate missed the 2026-11-14 break. Widening that coverage is a separate change and goes to `## Discovered Tasks` in the hand-back, not into this build.

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `_dev/tests/staged-skills-contract.sh` (modified) — docs/run-prompts added to the maintainer-docs export-ignore loop and to the legacy-root find exclusions; the loop comment names it (REQ-732).
- `.gitattributes` (modified) — an export-ignore line for /docs/run-prompts added after the /docs/decisions line, one comment paragraph explains the folder, and the header sentence that lists the asserted lines names it.
- `docs/run-prompts/2026-11-15-catalog-lists-second-opinion-pages.md`, `docs/run-prompts/2026-11-15-judge-renders-one-named-page.md`, `docs/run-prompts/2026-11-15-move-global-rules-into-skill.md`, `docs/run-prompts/2026-11-15-record-no-review-canvas-decision.md`, `docs/run-prompts/2026-11-15-relation-map-preset.md`, `docs/run-prompts/2026-11-15-run-prompts-folder-rule.md` (renamed from the _dev/run-prompts folder, 100% identical) — the six run prompts moved with git mv.

**What was done:** The staged-skills contract now accepts the root docs/run-prompts folder and still refuses any other new folder under the root docs folder; the folder is export-ignored, so it never ships. Hand-back merge range 8ad193d8..ce3cb5dc, no conflicts, no integration seams. The _dev/run-prompts folder is gone in the main checkout after the merge. The only remaining _dev/run-prompts mentions are inside the moved prompts themselves, which are historical records and keep their text.

## Decisions

(from the builder hand-back)
- D-04 DECIDE & STATE: reflow the two comments I touched to the neighbours' ~95-column width (including the pre-existing long line in `.gitattributes` header) instead of leaving a long line. Comment-only.
- D-05 DECIDE & STATE: the find exclusion went onto the existing `docs/changelog` line with the `-print -quit` moved to its own line, keeping each line short; behaviour is unchanged.

## Discovered Tasks

(from the builder hand-back)
- Heavy-lane `staged-skills` coverage in `_dev/tests/heavy-lanes.json` does not include root `docs/`, so the fast gate does not rerun the contract when a file lands under `docs/`. That is why the 2026-11-14 `docs/run-prompts/` file broke the lane unnoticed for five commits (D-03). → report only
- `git mv` leaves an empty source folder on disk; any probe or check that tests "folder is gone" with `-d` sees it until `rmdir`. Minor; no task needed unless repeated. → report only

## Qualification

**Gate records:** the qualify gate on range 8ad193d8..ce3cb5dc is satisfied. Its first run raised QUALIFY-CLAIMED-PATH-MISSING and QUALIFY-PATH-NOT-IN-DIFF only for backticked words in the Implementation Summary that are not file paths (a command name, a folder, an attribute line). I removed those backticks; no claim changed.

**Requirement trace** (`git diff --stat 8ad193d8..ce3cb5dc`: 8 files, 15+/9-):
- Both contract lists name the folder: the export-ignore loop in `_dev/tests/staged-skills-contract.sh` line 75 and the leftover-scan exclusion at line 107. Met.
- The export-ignore line in `.gitattributes` line 57 is present, and git reports `docs/run-prompts: export-ignore: set`. The request asked for one comment line; the builder wrote a three-line paragraph in the style of the neighbouring paragraphs and also extended the header sentence that lists the asserted lines. Comment-only, accepted.
- The six prompts are renamed with 0 changed lines, and the old folder no longer exists in the main checkout. `git archive ce3cb5dc` lists no run-prompts path. Met.
- Sweep: the only `_dev/run-prompts` mentions outside `do-work/` are inside two of the moved prompts, which are historical records and keep their text. Met.
- Constraint: nothing under `skills/` or `suite/` changed (0 paths), so this is no release.

**Scope:** Route A. The touched files equal the `write_set` exactly. No debug artifacts.

## Testing

**Tests run:**
- Repository gate `bash _dev/tests/maintainer-verify.sh` via `helpers/gate.sh` at merge `ce3cb5dc3b4099574b9d5ca103fd6cffc426957a` (detached drain-head checkout, stage reuse off), started at 1-minute load 2.90 → `GATE-EXIT=0 WALL=125s`.
- GREEN probe `do-work/runs/work-2026-11-15-012821/REQ-732-probe.sh` through `advance --probe-file` → probe status 0 (`REQ-732 probe PASS`: the real tree passes, a fixture with only a run-prompts probe file passes, and the same fixture plus a file in another new docs folder fails with the legacy-root message). Green gate recorded for the merge revision.
- The full staged-skills contract runs in the `staged-skills` heavy lane (see `## Heavy Verification Result`).

**Result:** ✓ All passing

**Red-green validation:** `tdd: false`. RED at the builder base (from the hand-back): with a run-prompts probe file present the contract exited 1 with `FAIL: legacy root runtime must be retired at modular cutover: docs (.../docs/run-prompts/probe.md)`, and the GREEN probe exited 1. GREEN at the builder tip `342a36a6`: the contract passes with the run-prompts probe file and still fails with a file in another docs folder. The GREEN probe passes again at the merge (re-run here).

**New tests added:** None. The REQ-732 probe is a run artifact, not a repository test.

**Heavy verification plan:** range 8ad193d8..ce3cb5dc. The planner marks coverage uncertain because no lane covers `.gitattributes`, so it selects all six lanes: `queue-kanban-javascript`, `queue-kanban-browser`, `do-work-cli-integrations`, `staged-skills`, `updater`, `installer`, each `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <lane>`, reason `coverage is uncertain for: .gitattributes`. The `staged-skills` lane is the one this REQ needs (it runs the changed contract).

*Verified by work action (integrator, run work-2026-11-15-012821)*

## Review

**Overall: 95%** | 2026-11-15T02:40:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 98% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F2: the auto-memory note `run-prompts-not-in-root-docs.md` and its `MEMORY.md` index line (outside the repository) still say to use `_dev/run-prompts/` and never create `docs/run-prompts/` — impact-rule-change → fixed before archive: the integrator deleted the retired memory note and its index line (outside the repository, no REQ, no commit)

**Minor findings:** F1 `.gitattributes:19` still 141 columns although D-04 claims the long line was reflowed — impact-negligible → report only; F3 non-archived queue `*Source:*` lines (REQ-727:79, REQ-725:62, REQ-726:67, REQ-732:68) still name `_dev/run-prompts/`, kept per D-01 — impact-negligible → report only; F4 (nit) `.gitattributes:57` attribute column one off — impact-negligible → report only
**Acceptance:** Pass — GREEN probe re-run at `ce3cb5dc` (allowed folder passes, `docs/other/` fails) plus the integrator's fast gate; `staged-skills` heavy lane not yet recorded
**Restatement sweep:** redefined the set of folders allowed under root `docs/` (kb, decisions, changelog, run-prompts) and the home of this repository's run prompts; grepped `docs/(kb|decisions|changelog)` enumerations and `_dev/run-prompts` outside and inside `do-work/`: `.gitattributes:17-18` and the contract comment at lines 71-74 agree, `skills/do-work-toolbox/actions/folder-structure.md:38` already lists `run-prompts/` under `docs/`, `CLAUDE.md` and `_dev/primes/` hold no restatement; stale: the out-of-repo auto-memory note (F2) and queue Source lines (F3)
**Anti-bloat:** 0 unrequested helpers, options, files or tests
**Suggested testing:** 1 items
**Follow-ups created:** None (3 findings report only, F2 is an orchestrator memory edit)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base: `8ad193d8c4c0c45d1a40c348f3f2341e1dc7d1f4`
- Target: `ce3cb5dc3b4099574b9d5ca103fd6cffc426957a`
- The planner marks coverage uncertain (no lane covers `.gitattributes`) and selects all six lanes. Each runs `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane <lane-id>`, reason `coverage is uncertain for: .gitattributes`:
  - `queue-kanban-javascript`
  - `queue-kanban-browser`
  - `do-work-cli-integrations`
  - `staged-skills` (the lane this REQ needs: it runs the changed contract)
  - `updater`
  - `installer`

## Heavy Verification Result

- Target: `ce3cb5dc3b4099574b9d5ca103fd6cffc426957a`; execution revision: `ce3cb5dc3b4099574b9d5ca103fd6cffc426957a` (detached drain-head checkout, `QUEUE_KANBAN_BROWSER` set, no lane skipped). `HEAVY-EXIT=0`, no findings.
  - `queue-kanban-javascript`: executed, exit 0, 7 s
  - `queue-kanban-browser`: executed, exit 0, 77 s
  - `do-work-cli-integrations`: executed, exit 0, 59 s
  - `staged-skills`: executed, exit 0, 31 s
  - `updater`: executed, exit 0, 65 s
  - `installer`: executed, exit 0, 26 s

## Lessons Learned

**What worked:** The pre-dispatch probe cut the contract's docs checks out of the real script by text markers, so it tested the shipped lines in under a second, in both directions (the allowed folder passes, another new docs folder fails).

**What didn't:** Decision D-04 says the long `.gitattributes` header line was reflowed, but the line is still 141 columns (review F1). Comment-only, report only.

**Worth knowing:** No heavy lane covers `.gitattributes`, so the planner marks coverage uncertain and selects all six lanes for any change to it. The `staged-skills` lane also does not cover root `docs/`, which is why the 2026-11-14 run-prompts file broke it unnoticed (Discovered Tasks). At the main-checkout merge git removed the emptied `_dev/run-prompts` folder by itself; only the builder's worktree needed `rmdir`. No lessons-satellite bullet: the builder proposed none.

## Orientation

The staged-skills contract (`_dev/tests/staged-skills-contract.sh`) and the `.gitattributes` export-ignore block together define which maintainer folders may live under root `docs/`; both now list kb, decisions, changelog and run-prompts.

## Timing

Observed 2026-11-14T23:33:30Z to 2026-11-15T02:45:34Z: 3h 12m 04s total, 3h 10m 24s attributed across 5 events, 1m 40s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 3h 00m 20s | 1 |
| verification-gate | 7m 26s | 2 |
| review | 2m 29s | 1 |
| handback-merge | 9s | 1 |

Slowest stage: builder-work / builder worktree build, 3h 00m 20s, outcome success.
Slowest command: verification-gate / heavy drain, 4m 39s, exit 0, .
