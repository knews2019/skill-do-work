---
id: REQ-699
title: 'ai-report: commit the bundle as the project rules say, and expect one judge link finding per href'
status: completed
created_at: 2026-10-19T14:29:39Z
user_request: UR-157
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: false
maintenance: true
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set: ["skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/docs/ai-report-guide.md"]
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-10-19T14:33:53Z
  basis:
    - trivial short-circuit
claimed_at: 2026-10-19T14:32:39Z
route: A
dispatch_at: 2026-10-19T14:34:30Z
builder_handback_at: 2026-10-19T14:35:38Z
integration_at: 2026-10-19T14:36:16Z
review_at: 2026-10-19T14:47:38Z
kb_status: pending
commit: 88f8c6865b361916d766df4e3651957590d7dc9c
heavy_verified_at: 2026-10-19T14:47:39Z
heavy_verified_revision: 88f8c6865b361916d766df4e3651957590d7dc9c
completed_at: 2026-10-19T14:48:10Z
release_at: 2026-10-19T14:48:10Z
---
# ai-report: Commit the Bundle as the Project Rules Say, and Expect One Judge Link Finding per href
## What
In `skills/do-work-toolbox/actions/ai-report.md`, replace the revise form's absolute "Do not commit" step with a commit rule that defers to the consumer project's instructions (commit the new bundle and the regenerated catalog pair by exact path when the project asks for completed increments to be committed, otherwise print that exact commit command, and inside a `do-work run` leave the commit to the run), say the same on the creation path's Step 8 summary, and change the four places that say the judge reports "one" expected prior-bundle link finding to one finding per `href`. Update the guide sentence that still says committing is left to the user. Ship it as a release.
## Why
On 2026-10-19 a consumer whose instructions say to commit completed increments ran `ai-report revise` directly. The agent followed step 7 literally, handed back a commit command, the user deployed, and the live site showed the old report because the bundle and catalog pair were untracked. The rule traces to the consumer's own 2026-10-09 upstream report ("Committing stays the user's or the run's decision, as today"), implemented by REQ-656 (the revise form); the consumer now reverses that scope line. The suite is already inconsistent: `architecture-report.md` commits its own bundle, ai-report creation is silent, revise forbids. The judge wording contradicts the code: `report_judge.go` appends one `AI-REPORT-JUDGE-BROKEN-LINK` finding per href or src attribute with no de-duplication, so two links into the prior bundle give two expected findings.
## Detailed Requirements
Each requirement is one accepted validate-feedback finding, with its provenance (source: pasted text, identical to the inbox file named in the UR).
- Finding 1, source: pasted text. Verbatim claim: "`skills/do-work-toolbox/actions/ai-report.md:59` (Revise form, step 7): '7. **Do not commit.** Committing stays the user's or the run's decision.' An absolute instruction; an agent following the action literally cannot commit even when the consumer's instructions require it." Severity: none given. Evidence: `ai-report.md:59` reads exactly that; origin REQ-656 requirement 9 (`do-work/archive/UR-144/REQ-656-ai-report-revise-supersedes-sibling-bundle.md:56`) quoting the consumer's 2026-10-09 report; precedent `architecture-report.md:48` and `:113` commit the report bundle. Surface-cost: Earned (replaces one rule with another, no guard or tooling; incident: the 2026-10-19 deploy that shipped the stale report; test: the acceptance greps below plus `_dev/tests/shipped-package-reference-contract.sh` and `_dev/tests/contract-regressions.sh`). Remedy: replace step 7 with the proposed text: "**Commit as the project's rules say.** Inside a `do-work run`, the run commits the bundle with its other output; do nothing here. Invoked directly, follow the consumer project's instructions on committing (its `CLAUDE.md` or equivalent): when they ask for completed increments to be committed, commit the new bundle and the regenerated `ai-reports/catalog.json` and `ai-reports/index.html` by exact path in the same turn and report the hash; otherwise print that exact commit command. Never stage anything outside the bundle and the catalog pair, and never push." The builder may instead use the simpler form "commit the bundle and the catalog pair by exact path by default, as architecture-report does; the project's instructions or the user override" when it reads cleaner; either form must give the same result in the consumer case above. Also change `skills/do-work-toolbox/docs/ai-report-guide.md:99` ("committing is left to you") to match.
- Finding 2, source: pasted text. Verbatim claim: "`skills/do-work-toolbox/actions/ai-report.md:176-178` (Step 8, Verify and Report the Result): confirms the bundle, removes the judge directory and prints a summary. Nothing about committing, so the creation form is silent where the revise form forbids." Severity: none given. Evidence: `ai-report.md:176-178` has no commit sentence; the 2026-10-09 report quotes a consumer asking on 2026-10-02 "give me the link after you commit it". Surface-cost: N/A (prose completion of Finding 1). Remedy: append to the Step 8 summary paragraph: "The summary also names the commit hash when the bundle was committed under the project's rules, or the exact commit command (bundle by exact path, plus the catalog pair for a revise) when the project leaves the commit to the user. Inside a run, the run commits."
- Finding 3, source: pasted text. Verbatim claim: "`skills/do-work-toolbox/actions/ai-report.md:57` (Revise form, step 5): 'so it reports that one link as `AI-REPORT-JUDGE-BROKEN-LINK`: this finding is expected'. `:170` (Step 7): '(for a revise, until only the expected prior-bundle link finding remains)'. Observed: two findings for two hrefs into the same prior bundle, one with an `#anchor`." Severity: none given. Evidence: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go:212-228` loops over every `[href]` and `[src]` attribute, strips the fragment only for the fetch (line 217), and appends one finding per attribute with status 400 or above; no de-duplication. Surface-cost: N/A (direct prose fix). Remedy: at `:57` replace "so it reports that one link as `AI-REPORT-JUDGE-BROKEN-LINK`: this finding is expected" with "so it reports each `href` into the prior bundle as `AI-REPORT-JUDGE-BROKEN-LINK`, one finding per `href` (an anchor link counts separately): these findings are expected"; at `:170` replace "(for a revise, until only the expected prior-bundle link finding remains)" with "(for a revise, until only the expected prior-bundle link findings remain)"; and make the same plural fix in the two places the paste missed: Step 8 at `:176` ("whose only finding is the expected link to the prior bundle") and the checklist at `:206` ("only the expected prior-bundle link finding").
- Finding 4, source: pasted text. Verbatim claim: "`skills/do-work-toolbox/actions/ai-report.md:201` (Verification Checklist): 'A revise wrote a new `-rev<N>` sibling with the rev block and supersedes meta, regenerated the catalog, and left the prior bundle byte-identical.' No line on where the regenerated catalog pair ends up." Severity: none given. Evidence: that line is now `ai-report.md:206` (moved by REQ-698's second-opinion form) and ends at "byte-identical." Surface-cost: N/A. Remedy: append "and committed the bundle and the catalog pair or printed the commit command, as the project's rules say" to that checklist line, in the same edit as Finding 3's plural fix of it.
- Release: a change to shipped files under `skills/` is a release; add the changelog entry per `_dev/primes/prime-releases.md` (title says what shipped, mirror byte-identical).
## Constraints
- Unchanged: the prior bundle's bytes never change, the catalog is still regenerated in revise step 6, the run owns commits inside `do-work run`, and no form ever pushes.
- Never stage anything outside the bundle and the catalog pair in the new rule's commit.
- Shipped prose only; no Go change. The judge code already matches the new wording.
- Shipped files never cite this repo's own `CLAUDE.md`; naming the consumer project's `CLAUDE.md` (or similar instructions) generically is fine (precedent: `skills/do-work/actions/capture.md:102`).
## Builder Guidance
Certainty: high; the wording is proposed in the UR and verified against code. Delete-before-you-add applies: deleting step 7 alone is not enough, because the silent creation path already produced the "give me the link after you commit it" follow-up on 2026-10-02, so a stated rule is needed. Keep the condition as the rule (the project's instructions decide), not a list of instruction file names. Line numbers in the findings are from the current file; re-grep before editing.
## Red-Green Proof
**RED prompt/case:** `grep -n "Do not commit" skills/do-work-toolbox/actions/ai-report.md` prints line 59; `grep -n "one finding per" skills/do-work-toolbox/actions/ai-report.md` prints nothing; `grep -n "left to you" skills/do-work-toolbox/docs/ai-report-guide.md` prints line 99.
**Why RED now:** The revise form forbids committing outright, the creation path says nothing, and four lines say the judge reports one expected prior-bundle link finding when the code reports one per href.
**GREEN when:** `grep -n "Do not commit" skills/do-work-toolbox/actions/ai-report.md` prints nothing; `grep -n "Commit as the project's rules say" skills/do-work-toolbox/actions/ai-report.md` prints the step 7 line (or the chosen simpler default wording, in which case the REQ names it); `grep -n "one finding per" skills/do-work-toolbox/actions/ai-report.md` prints the revise step 5 line; `grep -n -i "link finding remains\|only finding is the expected link\|only the expected prior-bundle link finding" skills/do-work-toolbox/actions/ai-report.md` prints nothing; the guide line no longer says committing is left to the user; `_dev/tests/contract-regressions.sh` and `_dev/tests/shipped-package-reference-contract.sh` pass.
**Validation:** User confirmed (the acceptance checks come from the user's own upstream report; the user accepted the triage with "ok, capture the request and run them").
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` (8128 tokens, `slugged: partial` so no targeted form): matches because the request changes an action's stated rules and a downstream reader (the guide); over the 2000-token budget.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read prime-action-files, prime-releases, lessons-releases, the maintenance crew and coding-guardrails. Approach: state the commit rule once in Step 8 (reached by both forms), have revise step 7 apply it after the catalog regenerates, fix every singular "link finding" phrase (five sites), update the guide line.
- [x] **[APPLY]:** Seven replacements in `skills/do-work-toolbox/actions/ai-report.md` and one in `skills/do-work-toolbox/docs/ai-report-guide.md`, each asserted unique before replacing; no other file touched (builder commit 07e075d1).
- [x] **[UNIFY]:** `git diff --stat`: 2 files, +10/-8, both read in full. ai-report.md: acceptance greps (no "Do not commit", "Commit as the project's rules say" at step 7 and Step 8, "one finding per" at revise step 5, no singular "link finding" left). ai-report-guide.md: "left to you" gone. `_dev/tests/contract-regressions.sh` exit 0 (26s) and `_dev/tests/shipped-package-reference-contract.sh` exit 0 (1s) in the builder worktree. Prose only, no debug artifacts.
## Full Context
See `do-work/user-requests/UR-157/input.md` for complete verbatim input.
*Source: validate-feedback --capture --run on the pasted upstream suggestion "the ai-report revise form must not forbid committing" (do-work/inbox/2026-10-19_do-work-upstream-suggestion-ai-report-revise-commit-rule.md); Findings 1-4, source: pasted text.*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the exact files and lines (ai-report.md step 7, step 5, Step 7, Step 8, the checklist line, and guide line 99) and supplies the replacement wording; the judge code already behaves as the new wording says, so this is a prose edit plus a release.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- D-01 (DECIDE & STATE): the commit rule lives in Step 8, which both forms reach, and revise step 7 only applies it after the catalog regenerates. Reasoning: one statement of the rule instead of two copies; the revise needs its own step because its commit must follow step 6's catalog regeneration.
- D-02 (DECIDE & STATE): kept the consumer's "follow the project's instructions" form rather than "commit by default, project rules override". Reasoning: the REQ's GREEN pins the phrase "Commit as the project's rules say", and a default-commit form would change behaviour for projects whose instructions are silent, which no report asked for.
- D-03 (DECIDE & STATE): fixed the singular wording at five sites (the paste named three: lines 57, 170 and 201; the Step 8 verify sentence and the checklist render line carried it too).

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/ai-report.md` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)

**What was done:** Revise step 7 no longer says "Do not commit": Step 8 now holds the commit rule (follow the project's own instructions on committing, inside a run the run commits, never stage outside the bundle, never push), the revise form applies that rule at its step 7 after the catalog regenerates and commits the catalog pair with the bundle, and the Step 8 summary names the commit hash or the exact commit command. The five sentences that said the judge reports one expected prior-bundle link finding (revise step 5, Step 7, Step 8, and two checklist lines) now say one finding per `href`, matching `report_judge.go`. The guide's "committing is left to you" sentence says the same as the action. Builder commit 07e075d1, merged at 3fccc496 (range 6f6f8655..3fccc496).

## Qualification

**Result:** Pass for substantive implementation at `6f6f8655..3fccc496`. The canonical request/path-bound qualify record is satisfied with no findings; the three P-A-U boxes are ticked from the hand-back. Read the whole diff: two shipped prose files, +10/-8, every changed line traces to one of the four accepted findings (step 7 rule, Step 8 rule and summary, five plural link-finding sentences, checklist commit clause, guide sentence). No Go change, which matches the REQ's constraint that the judge already behaves as the new wording says (`report_judge.go:212-228` appends one finding per attribute).

**Findings judged:** None mechanical. Judgment: the rule is stated once (Step 8) and applied by the revise form at its step 7 after the catalog regenerates, so the two forms no longer disagree; "Inside a `do-work run`, the run commits" preserves the run's ownership; "never push" is kept in both places; the regenerated catalog pair rides the revise commit only, which is the only form that writes it.

**Acceptance trace:** `grep -n "Do not commit"` prints nothing; `grep -n "Commit as the project's rules say"` prints lines 59 (revise step 7) and 178 (Step 8); `grep -n "one finding per"` prints line 57; no singular "link finding" phrase remains (lines 170, 176, 198, 201 before the edit); the guide no longer says committing is left to the user. Shipped-file hygiene: `_dev/tests/contract-regressions.sh` and `_dev/tests/shipped-package-reference-contract.sh` both exit 0 on the builder branch. The consumer-side end-to-end check (a direct revise in a project whose instructions ask for commits ends with a hash and a clean `git status`) cannot run in this repository and is recorded as unrun.

*Qualified by work action*

## Testing

**Tests run:** `bash do-work/runs/work-2026-10-19-143229/REQ-699-probe.sh` (the REQ's acceptance greps, then `_dev/tests/shipped-package-reference-contract.sh` and `_dev/tests/contract-regressions.sh`) through the advance test gate; `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` directly in the detached checkout `.git/work-run-2026-10-19-143229/drain-head` at 88f8c686.
**Result:** ✓ Probe exit 0 (every grep as the Red-Green Proof states; both contract tests pass, 27s). Repository gate exit 0 on the rerun (157s). Every test-file invocation stayed under the 30s budget on the green run.

**Repository gate retry:** first run exited 1 after 141s only on the per-file Go test budget of an unrelated module under machine load (`go-test budget: FAIL module=skills/do-work/tools/do-work-cli wall=72s`, load average about 4.9); rerun exited 0 after 157s. Downstream reads the rerun.

**Red-green validation:** non-behavioral prose change; the captured RED greps (`Do not commit` at line 59, no `one finding per`, guide `left to you`) now return the GREEN results (none, line 57, none), traced to `## Red-Green Proof`.

**Heavy verification plan:**
- Range: 6f6f8655..88f8c686
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work-toolbox/actions/ai-report.md and skills/do-work-toolbox/docs/ai-report-guide.md matched subtree skills

*Verified by work action*

## Review

**Overall: 100%** | 2026-10-19T14:47:38Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None open. First pass (98%, range 6f6f8655..3fccc496) raised M1 (Step 8's "otherwise" branch said "leave the commit to the user" while step 7 and the summary said "print that exact commit command") and N1 (the summary's "when the project leaves the commit to the user" could qualify both alternatives), both impact-negligible; fixed on the builder branch (e97d4aa6), re-merged at 88f8c686, re-check confirmed both resolved and nothing else changed.
**Acceptance:** Pass — all four accepted findings implemented, acceptance greps hold on the merged tree, constraints hold (prior-bundle bytes sentence untouched, revise step 6 untouched, run owns commits, never push, no Go change), prose matches `report_judge.go:212-228`.
**Restatement sweep:** redefined the ai-report commit rule and the expected prior-bundle link-finding count; swept `skills/` and `_dev/primes`, no stale restatement outside the diff (CHANGELOG is history, not flagged).
**Suggested testing:** 2 items (a direct `ai-report revise` in a consumer whose instructions ask for commits ends with a hash and a clean `git status`; the release entry title and mirrors at finalization), both report only.
**Follow-ups created:** None (2 findings report only, both resolved before merge)

*Reviewed by review-work action (background reviewer, report at do-work/runs/work-2026-10-19-143229/REQ-699-review.md); re-reviewed after the review-fix re-merge against cumulative range 6f6f8655..88f8c686*

## Heavy Verification Plan

Base: `6f6f865557a91e9ce9befee71fd5c84b63ddabc4`
Target: `88f8c6865b361916d766df4e3651957590d7dc9c`

- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work-toolbox/actions/ai-report.md matched subtree skills; skills/do-work-toolbox/docs/ai-report-guide.md matched subtree skills

## Heavy Verification Result

Target revision: `88f8c6865b361916d766df4e3651957590d7dc9c`
Execution revision: `88f8c6865b361916d766df4e3651957590d7dc9c` (detached checkout `.git/work-run-2026-10-19-143229/drain-head`, run at 2026-10-19T14:47:39Z)

- staged-skills: exit 0, 38s wall, disposition executed (no prior evidence reused)

## Lessons Learned

**What worked:** Verifying the judge's loop (`report_judge.go:212-228`) before accepting the "one finding per href" claim turned a wording request into a code-backed fix and found two more singular sentences than the paste named. Stating the commit rule once in Step 8 and having the revise step apply it after the catalog regenerates kept the two forms from drifting again.
**What didn't:** The first wording left Step 8's "otherwise" branch as "leave the commit to the user" while the other two sites said "print that exact commit command"; the reviewer caught it before the gate ran, so the fix cost one re-merge and one gate.
**Worth knowing:** The "Do not commit" step came from the consumer's own 2026-10-09 report ("as today") via REQ-656 requirement 9, so a scope line in an upstream report can become an absolute rule downstream; `architecture-report.md:48` and `:113` already committed their own bundle, so the toolbox had two commit postures for one kind of artifact.

## Orientation

Now an `ai-report` bundle written outside a run reaches the project's deployable git state when the project's instructions ask for committed increments: both forms of the toolbox `ai-report` action (creation and revise) follow one commit rule, the revise commit carries the regenerated catalog pair, and the summary names the hash or the exact command. The render-check prose counts expected prior-bundle link findings the way the judge code does. Lives in the do-work-toolbox report subsystem (`skills/do-work-toolbox/actions/ai-report.md` and its guide); no map change. Touched primes `_dev/primes/prime-action-files.md` and `_dev/primes/prime-releases.md` spot-checked: every path they name still exists.

## Timing

Observed 2026-10-19T14:34:30Z to 2026-10-19T14:46:33Z: 12m 03s total, 7m 21s attributed across 2 events, 4m 42s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 6m 13s | 1 |
| builder-work | 1m 08s | 1 |

Slowest stage: verification-gate / repository gate and staged-skills lane at 88f8c686, 6m 13s, outcome success.
