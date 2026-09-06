---
id: REQ-549
title: '[impact-negligible] Drop the eight dead path tokens from the decision indexes and the lessons prime'
status: completed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: []
related: [REQ-550, REQ-551, REQ-552, REQ-553, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: true
impact: impact-negligible
effort_estimate: effort-mechanical
write_set: [decisions/topics/_index_skill-architecture.md, decisions/topics/_index_knowledge-base.md, decisions/audits/2026-08-05-shell-logic-in-prose-census.md, decisions/imported-specs/2026-04-12_close-gaps-in-interview.md, _dev/primes/lessons-action-files.md, _dev/tests/audit-lockins.sh]
claimed_at: 2026-09-05T00:38:07Z
route: A
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-09-05T00:50:53Z
  basis:
    - Route A
    - 6-file write set
    - 2 subsystems involved
    - 3 acceptance criteria
completed_at: 2026-09-05T02:07:31Z
commit: cb4c67cd
---

# Drop the eight dead path tokens from the decision indexes and the lessons prime

## What
Eight path tokens cited in `decisions/` and `_dev/primes/` resolve to no tracked file (31 citation lines). Drop the dead entries from the two topic indexes, fix the bare `work-reference.md` in the lessons prime, and mark the census and imported-spec pointers retired the way the ADRs already do. Leave the ADRs untouched.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** — Read the REQ, `_dev/primes/prime-action-files.md`, `prime-shell-commands.md`, `anti-slop.md`, `coding-guardrails.md`, `maintenance.md` and `CLAUDE.md`, then verified all eight tokens with `git ls-files` before editing anything (V1 output above) and mapped all 32 citation lines to owners. The plan changed once, on evidence: the REQ's stated lock-in target is unreachable from the write set, so the assertion was re-aimed at the live-routing class (D-01).
- [x] **[APPLY]:** — Exactly the six write-set files changed; `git status --short` listed those six and nothing else, and the commit's diffstat is `6 files changed, 71 insertions(+), 5 deletions(-)`. No file under `do-work/` was written, staged or committed on the branch.
- [x] **[UNIFY]:** — Read the full `git diff` before committing. Per file: both topic indexes lost exactly one `sources:` line and gained a `updated:` date, nothing else; the census and the imported spec gained only a blockquote each with no surrounding prose touched; the prime changed one token on one line; `audit-lockins.sh` gained one block before the existing failure check and kept mode `-rwxr-xr-x`. Linters run: `shellcheck --severity=warning` exit 0, `bash -n` exit 0, `bash _dev/tests/audit-lockins.sh` exit 0 with the RED case proven first. No debug artifacts, no leftover scratch files …

## Why
Agents read the topic indexes and the prime as live routing; a pointer to a deleted file costs a search on every read and teaches the reader that the records are unreliable.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 8, sweep_key `dead-path-pointers-in-records`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag SIMPLE; expected net line delta -8. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `decisions/topics/_index_skill-architecture.md` lists `crew-members/karpathy.md` as a live source: remove the entry.
- `decisions/topics/_index_knowledge-base.md` lists `actions/build-knowledge-base.md` as a live source: remove the entry.
- `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` cites `hooks/pipeline-guard.sh`: mark retired (the hook is gone; `internal/settingshooks/settings_hooks.go` keeps the name only to remove it from consumer settings).
- `decisions/imported-specs/2026-04-12_close-gaps-in-interview.md` names `decisions/imported-specs/expand-skill-do-work-interview.md` as required reading on three lines: mark retired.
- `_dev/primes/lessons-action-files.md` cites a bare `work-reference.md`, unresolvable from `_dev/primes/`: write the full `skills/do-work/actions/work-reference.md` path like every other prime.
- Not touched: the ten ADRs (`adr-001,003,005,006,007,008,009,010,013,014`) are immutable history; `decisions/audits/2026-08-11-defensive-surface.md` is a frozen snapshot by its own header.
- Reproduce at dc8a64e3: `for t in crew-members/karpathy.md actions/build-knowledge-base.md hooks/pipeline-guard.sh actions/pipeline.md actions/pipeline-reference.md _dev/tests/record-commit-hash-guards.sh skills/do-work-knowledge/crew-members/security.md decisions/imported-specs/expand-skill-do-work-interview.md; do git ls-files | grep -q --fixed-strings "$t" || rg -n --fixed-strings "$t" _dev/primes decisions; done`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Prose only outside the lock-in; no shipped file under `skills/` changes.
- Lock-in limit: dead path tokens cited outside decisions/records: 0 after this REQ (today 8); red when the Reproduce command prints a non-ADR line.

## Dependencies
No dependency. First REQ of the audit batch (`batch: maintainability-audit-2026-09-03`); nothing waits on it.

## Builder Guidance
Firm. Delete or mark retired; do not rewrite surrounding prose.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints 31 lines, five of them outside the ADRs (two topic indexes, the census, the imported spec, the prime).
**GREEN when:** The command prints only ADR lines, and the lock-in asserts the count of dead tokens cited outside `decisions/records/` is 0 (pinned first at today's 8 and lowered to 0 in this REQ).
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3968 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for dead-path-pointers-in-records.*

---

## Triage

**Route: A** - Simple

**Reasoning:** The request names the eight dead tokens, the six files that cite them, the 31 citation lines, and the ADR retirement pattern to match. Verification is mechanical.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
## Implementation Summary

**Files changed:**
- `decisions/topics/_index_skill-architecture.md` (modified)
- `decisions/topics/_index_knowledge-base.md` (modified)
- `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` (modified)
- `decisions/imported-specs/2026-04-12_close-gaps-in-interview.md` (modified)
- `_dev/primes/lessons-action-files.md` (modified)
- `_dev/tests/audit-lockins.sh` (modified)

**What was done:** Two dead path entries were dropped from the decision topic indexes, retirement notes were added to the census and the imported spec, the bare filename in the lessons prime was written out as a full path, and one lock-in check was added that fails when a live routing surface cites a path matching no tracked file.

All eight tokens named by the request were re-verified dead before anything was edited. Two commands were run per token from the worktree root, the substring form the request's Reproduce line uses and the strict pathspec form, and both returned empty output for all eight. Every citation line was then mapped to its owner, which turned up a correction to the request: the request says 31 citation lines, the real number is 32, both in the worktree and at the audited commit dc8a64e3. The audit undercounted by one and nothing about the fix changes. Of those 32 lines, 19 sit in decisions/records/ and were left untouched as the request instructs.

Each topic index lost exactly one entry from its sources list, and both files had their updated date bumped to 2026-09-05 in the same commit. The bump was not asked for. It follows the directory's own precedent: REQ-145 (commit c42f228b, which removed the stateful pipeline) bumped the updated date in the same commit that dropped a dead source from a sibling index.

The two retirement notes copy the shape already used in adr-003: a dated note blockquote, the fact about the path, and a closing clause saying the record's own reasoning stands and the historical references are left as written. The sibling superseded-page shape used in adr-005, adr-006 and adr-008 was rejected because neither of these two pages is being retired as a whole. Both facts asserted inside the notes were checked in git rather than assumed. The retired hook was deleted by REQ-145 on 2026-08-08, and the Go settings-hooks file keeps the name only to strip it from consumer settings; the v1 interview spec was added in commit 049e6ff1 and deleted in commit d86a3418 on 2026-04-16. No surrounding prose was rewritten in either file, and no existing line changed.

The lessons prime changed one token on one line: the bare filename became the full repository path that the request specifies and that already appears twice elsewhere in the primes.

The lock-in added to the audit script keys on a condition rather than on a list of the eight tokens. It reads every sources entry in the topic indexes and every backticked directory-plus-extension token in the primes, then asks whether any tracked path contains that token. Dated records are not scanned, and that boundary is written into the check's header comment. The count it pins was 2 before this change and is 0 after. The audit script was not newly created and its registration in the fast tier was not touched, because both already existed.

**Implementation range:** `1d6aad51..cb4c67cd`. Builder commit `0e631d13`.

## Decisions

- **D-01 — the request's stated GREEN is not reachable from this write set, so the lock-in pins the class the finding is actually about:** The request asks that the Reproduce command print only ADR lines and that the count of dead tokens cited outside decisions/records/ be 0. Two independent reasons block that. First, three files outside decisions/records/ also cite dead tokens and are outside the write set: the decisions log, the phase-1 hashes audit, and the defensive-surface snapshot that the request itself declares frozen and not-touched. Second, "mark retired" deliberately keeps the token in the text, because that is what makes the note readable, so the census and the imported spec still match the grep afterwards and each in fact gains one line. The lock-in instead pins the live-routing class that the request's own Why names: a path token in a topic-index sources list or in a prime citation that matches no tracked file. That count was 2 before and is 0 after, and it fails loudly the moment a third appears. This is a deliberate deviation from the request's written GREEN.
- **D-02 — the lock-in keys on conditions, not a token list:** A pinned list of the eight tokens would go stale the moment a ninth file is deleted, which the maintainer doc calls out under Closed Enumerations Go Stale. The check reads every sources entry and every backticked directory-plus-extension token and asks whether any tracked path contains it. It carries three carve-outs, each stated as a condition and each already drawn by the shipped cross-referencing contract in the action-files prime: a glob pattern is not a path, a token rooted somewhere the reader resolves from instead is not a path from here, and a bare consuming-project queue token is state rather than a citation.
- **D-03 — the two index entries were removed rather than replaced, and the builder disagrees mildly:** The karpathy crew-member file was renamed, not deleted, so replacing the entry would have killed the dead pointer and kept the routing value, and the nearby precedent supports replacement. Removal was chosen anyway because the request says "remove the entry", marks Builder Guidance "Firm", and a mechanical batch request should not carry a judgment call. The thought is recorded as a discovered task rather than dropped. The knowledge-base entry has no live successor to point at, so removal is unambiguously right there.
- **D-04 — the updated date was bumped on both indexes:** Not asked for, but it is the same file's own convention and the directory's precedent. Leaving it stale would have made one index claim it was last touched in 2026-04-15 while its source list had changed.
- **D-05 — in-shell substring match instead of a piped quiet grep:** The first version piped the tracked-file list into `grep -qF` and reported all 48 tokens dead. `grep -q` exits on the first match and closes the pipe, and under the script's `set -uo pipefail` the writer's SIGPIPE status of 141 becomes the pipeline's status, so the following `&&` never fired. It was replaced with a bash `case` substring test, which is the same reading, has no subprocess, and cannot trip the trap. The reason is written in a comment at the site.
- **D-06 — no guard on the tracked-file listing returning empty:** If that listing ever failed, every token would report dead and the check would go loudly red, which is the safe direction. A guard would be defensive surface with no incident behind it, which the shipped coding guardrails forbid under Earned defense.
- **D-07 — the lock-in does not catch bare filenames:** The bare-filename defect this request fixes is a different class from a dead path, and its token still resolves by substring against a tracked file. Catching that class needs a different predicate with its own false positives, since several bare filenames are legitimate in these files, and the request asks for one assertion. Left out on purpose.

## Qualification

Passed the request-bound advance qualify gate for `1d6aad51..cb4c67cd` against the merged range. Exactly the six declared files changed, 71 insertions and 5 deletions, no undeclared touch, nothing under `do-work/` on the builder branch, and no file under `skills/`. The P-A-U boxes were reconciled from the builder hand-back, which is where worktree dispatch puts them.
## Testing

**Red-green validation:** RED was produced by temporarily restoring the two removed index entries and running the audit script. The hand-back records the output as:

```
$ bash _dev/tests/audit-lockins.sh
FAIL: decisions/topics/_index_knowledge-base.md cites actions/build-knowledge-base.md, which matches no tracked file
FAIL: decisions/topics/_index_skill-architecture.md cites crew-members/karpathy.md, which matches no tracked file
exit=1
```

GREEN is the same command on the committed tree:

```
$ bash _dev/tests/audit-lockins.sh
Audit lock-in regressions passed.
exit=0
```

There are no test function names to quote. The lock-in is a 59-line block added to the shell script `_dev/tests/audit-lockins.sh`, placed before the existing failure check, not a named function in a test framework. The hand-back names it only as "the Finding 8 lock-in".

**Controls preserved:** The fast-tier suite `_dev/tests/contract-regressions.sh` was re-run in full. The hand-back records every content probe passing:

```
shipped package reference contract: PASS
Audit lock-in regressions passed.
replace-text-section contract probes passed.
recovery set-aside contract probes passed.
Suite manifest contract probes passed.
Shell-block lint passed: 73 fenced blocks and 32 shipped shell files; ShellCheck enabled.
SessionStart hook behavior probes passed.
Prescribed shell primitive canonicalization checks passed.
Defensive-surface exact deletion regressions passed.
do-work-cli launcher behavior tests passed
p50 estimator suite: all probes passed.
select-simple-reqs suite: all probes passed.
Go test budget behavior probes passed.
```

Two of these protect the surfaces this change touches: the audit lock-in probe, which is the file the new assertion was added to, and the shell-block lint, which covers the 32 shipped shell files and runs ShellCheck over them. The hand-back does not describe what each of the remaining probes protects, so that is not stated here.

**Module verification:** ShellCheck version 0.11.0. At the floor the gate actually uses, recorded at `_dev/tests/maintainer-verify.sh:566`:

```
$ shellcheck --severity=warning -- _dev/tests/audit-lockins.sh
gate-severity exit=0
$ bash -n _dev/tests/audit-lockins.sh
syntax exit=0
```

Plain ShellCheck with no severity floor exits 1 with three sub-warning items: two pre-existing SC2126 style notes on lines 28 and 48, which are not builder code, and one SC2016 info on line 118, the single-quoted backtick-matching pattern passed to `grep -o`. The builder judged SC2016 a false positive there, since the single quotes are what the pattern needs, and added no disable directive because the item sits below the gate's floor.

The contract-regressions suite exited 1. The five failures are per-file duration limits only, with no content failure among them:

```
FAIL: _dev/tests/contracts/replace-text-section.sh took 36s; each fast test file must finish under 30s
FAIL: .../suite-manifest-contract.sh took 36s; each test file must finish under 30s
FAIL: .../action-shell-blocks.sh took 40s; each test file must finish under 30s
FAIL: .../session-start-hook-behavior.sh took 112s; each test file must finish under 30s
FAIL: .../prescribed-shell-canonicalization.sh took 64s; each test file must finish under 30s
```

These five failures were diagnosed as not belonging to this change. None of the five files is in the write set, and all five are machine-contention artifacts from several builders running in parallel on the same box. Two were re-measured standalone on the same tree immediately after: the suite-manifest file took 10s against the batch's 36s, and the action-shell-blocks file took 6s against the batch's 40s, both under the 30s limit. The changed file itself timed the same with and without the new block, 1s in both runs. The hand-back says the orchestrator's post-merge gate on a quiet box is where this should be confirmed.

The eight tokens were also re-verified dead in the worktree with a loop over `git ls-files`, printing "(no tracked file)" for all eight, and the strict pathspec form returned empty for all eight as well.

## Discovered Tasks

- **DT-1 — the removed skill-architecture entry has a live successor.** The dropped token was renamed, not deleted, and the successor is a tracked file. Re-adding it would keep the skill-architecture topic routing to the guardrails it was built from. One line, in the sources list of `decisions/topics/_index_skill-architecture.md`, which was line 24 before the edit. See D-03. → queue as follow-up
- **DT-2 — the declined-ADR-014 entry in the decisions log cites a file that no longer exists under that name.** At `decisions/log.md:94` it names the old crew-member path among "surviving canon". The ADR carries a rename note; the log entry does not. Outside this write set. → queue as follow-up
- **DT-3 — four rows in the frozen defensive-surface snapshot cite two untracked files.** At `decisions/audits/2026-08-11-defensive-surface.md` lines 27, 30, 46 and 98. The file's own header declares it a frozen snapshot, so this may be correct as it stands; it is the place to start if frozen snapshots are meant to carry retirement notes the way the ADRs do. → report only
- **DT-4 — the fast tier's 30s per-file limit produces false failures under parallel load.** Five test files exceeded it in this run and finish in 6 to 10 seconds standalone. If builders run concurrently as a matter of course, that limit will keep failing files that have nothing to do with the request under test. Evidence is in the V4 section of the hand-back. → report only

## Review

**Overall: 91%**
**Acceptance: Pass.** The reviewer re-verified all eight tokens dead with `git ls-files` rather than trusting the ledger, reproduced RED on an isolated copy, and tried to break the substituted lock-in three ways — a dead token in a third topic index and in two different prime files. All three were caught.

It judged the substituted acceptance target a legitimate scope judgment rather than a moved goalpost, and found a second and stronger reason than the builder gave: the request instructs the builder to mark citations retired, and a retirement note deliberately keeps the token in the text, so those lines still match the grep afterwards. The census went from one matching line to two. The request asks for two things that cannot both be true.

One false-positive class was found: `<skill-root>/…` and `.claude/skills/…` tokens produce failures the carve-out comment claims are excluded. The `<skill-root>/` spelling is already live in the repository, so the next prime that copies that idiom breaks the check.

The reviewer's closing point is not about this work at all: the request that produced it is internally contradictory in four ways, and one of those — a stale instruction to create a test file that already exists — sat in eight requests. That was corrected across all of them during this run.

## Lessons Learned

Two rules came out of this, both above the level of the files touched.

A check that counts textual matches of a token cannot be reconciled with an instruction to write retirement notes about that same token, because a readable retirement note keeps the token in the text and adds a matching line. When an audit request asks for both, the acceptance target has to be re-aimed at the class the finding is actually about, and the deviation stated plainly. Look for this whenever an acceptance criterion is a grep count and the work includes annotating history.

In a shell check that runs under `set -uo pipefail`, do not test membership with a pipeline into `grep -q`. The quiet grep exits at the first match and closes the pipe, the writer dies with SIGPIPE, and 141 becomes the pipeline's status, so the success branch never runs and every item reports as missing. Use a shell-native substring test on a variable instead, which has no subprocess and cannot trip the trap.

## Orientation

The two decision topic indexes and the lessons prime now route only to files that exist, and a fast-tier check fails the moment any topic-index source or prime citation stops resolving to a tracked path. Dead pointers that remain in dated records are now either annotated as retired or recorded as follow-up work, rather than being indistinguishable from live routing.
