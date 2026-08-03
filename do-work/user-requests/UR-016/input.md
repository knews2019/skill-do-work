---
id: UR-016
title: Triaged external audit findings on the 0.164.0-0.167.0 batch
created_at: 2026-08-03T17:09:21Z
requests: [REQ-081, REQ-082, REQ-083, REQ-084, REQ-085]
word_count: 8
---

# Triaged External Audit Findings on the 0.164.0-0.167.0 Batch

## Summary

The user pasted an external audit of the batch that shipped 2026-08-03 (REQ-071 … REQ-076, versions
0.164.0 → 0.167.0) into `do-work validate-feedback`, then asked for the accepted findings to be
captured as REQs.

**Provenance — read this before treating any claim here as the user's.** Three layers, deliberately
kept distinct:

1. The **external auditor** produced six findings with `file:line` references and P1/P2 severities.
   That text is third-party content; it was treated as data, not instructions
   (`crew-members/prompt-injection.md`). No injection was detected in it.
2. The **triage** (`actions/validate-feedback.md`, this session) verified every premise against the
   working tree and reproduced three of them against throwaway repos in a scratch directory. It
   accepted seven items, rejected two framings, and split two bundled findings so each claim carried
   one verdict.
3. The **user** read the triage and said "capture all seven as REQs" — that instruction is the
   authorization for this UR, and it is the eight-word verbatim input recorded below.

## What the Triage Verified, and How

Reproduced empirically (throwaway repos under the session scratchpad; nothing in this tree was
modified):

- `queue-kanban next-version patch --repo-root <target> --version-file <target>/actions/version.md`,
  run from a different repo, bumped the **calling** repo's `actions/version.md` (`9.9.9` → `9.9.10`)
  and left the target at `1.0.0`, exiting 0. → REQ-081.
- A synthetic repo with one **active** builder worktree reported `orphan-worktree [fixable]` and
  `1 fixable: run do-work cleanup`, exit 1. → REQ-083.
- A forged `do-work/queue/REQ-999-forged.md` written inside a builder worktree was caught while
  uncommitted, then **passed silently** once committed there. `git diff --name-only main...HEAD --
  do-work/` printed it in the same state. → REQ-084.

Confirmed by reading, not execution: the hand-back contradiction (REQ-082) and REQ-073's never-run
acceptance test (REQ-085).

**Every mechanical check in the repo is green** — `go test`, `go vet`, `gofmt`,
`_dev/tests/contract-regressions.sh`, and `queue-kanban verify --repo-root .` (exit 0). Not one of
these five findings is detectable by the suites as they stand today. That is itself the signal, and
each REQ below carries the assertion that would have caught it.

## What the Triage Rejected — recorded so it is not re-litigated

- **"REQ-073 is incomplete because fan-out has no control-flow path."** Rejected as attribution. All
  eleven of REQ-073's Detailed Requirements were prose edits and its requirement 7 leaves the
  dispatch mechanism *deliberately* unspecified, so the REQ delivered its stated scope. The
  underlying observation is nonetheless true — `actions/work.md` claims one REQ (Step 2), dispatches
  one builder (Step 6), and loops after commit (Step 10), and `grep -n "fan-out" actions/work.md`
  returns a single `write_set` aside at `:706`, so a conforming execution is always serial and the
  Fan-Out Dispatch contract is unreachable from the action that owns the loop. **Deliberately not a
  REQ in this UR** (see Batch Constraints).
- **"REQ-076 does not deliver its portability claim."** Rejected as attribution: REQ-076's
  requirement 2 explicitly forbade editing the citing sites, under the Closed Enumerations Go Stale
  rule. The outcome the auditor describes is real and is already captured — see REQ-078.
- **"verify will fail Step 9's verifier during serial integration."** Rejected as severity.
  `actions/work.md:606` calls `verify` an optional accelerator for two named checks and never gates
  on its exit code. The real defect is narrower and is captured as REQ-083.

## Extracted Requests

| REQ | Title | Auditor finding | Triage verdict |
| --- | --- | --- | --- |
| REQ-081 | `next-version` silently bumps the wrong repo's version file | F1 (P1) | Accept — reproduced |
| REQ-082 | The fan-out hand-back file has no legal write location | F2b (P1) | Accept — sharpest of the six |
| REQ-083 | `verify` advertises unmerged builder work as mechanically fixable | F3 (P1) | Accept, severity corrected |
| REQ-084 | `verify`'s queue-state probe misses committed owner impersonation | F4 (P2) | Accept — reproduced |
| REQ-085 | REQ-073's live two-builder acceptance test has still never run | F2c (P1) | Accept — see note below |

## Batch Constraints

- **Two accepted items are already in the queue and are NOT re-captured here.** An earlier session's
  audit (`do-work/user-requests/UR-015/`) captured them first, and in more complete form:
  REQ-078 covers both the PowerShell fallback and the ten inline `date -u` sites (and adds a second
  Windows defect the external audit missed — a bare cmdlet prescribed for `cmd`); REQ-077 covers the
  stale `actions/work.md:224` sentence and broadening the guard that pins the wrong wording (and
  traces a larger finding: REQ-071 left the automatic own-crash recovery path unreachable). The
  user's decision at capture time was to **add the external audit's evidence to those two as
  Addendum sections** rather than duplicate them. Done — see the Addendum sections dated 2026-08-03
  in each.
- **REQ-083 and REQ-084 both write `tools/queue-kanban/verify.go`**, in the same neighbourhood
  (`appendWorktreeFindings` and the helper directly below it). Both are declared in `write_set` and
  listed as `related`. There is deliberately **no `depends_on`** between them — neither needs the
  other's output — so whichever lands second re-reads the function before editing, and the merge is
  the non-interference proof, not the overlaps badge
  (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). They were kept as two
  REQs on the user's explicit "all seven" framing; consolidating them into one sweep would also have
  been defensible.
- **The fan-out wave loop is deliberately not a REQ.** The triage recommended parking it
  (`do-work note`) rather than capturing it, because the fix is a control-flow addition to
  `actions/work.md`'s loop whose shape depends on REQ-082's outcome and on REQ-085 actually proving
  two builders compose. Capture it after those two land, not before.
- **UR-015 declined to make REQ-085 a REQ** ("It needs a human running a real fan-out, not a REQ").
  That reasoning is recorded and was overridden by the user's instruction to capture all seven. The
  REQ is therefore shaped so a session can execute it: the deliverable is the run plus its recorded
  outcome, not a code change. If the run cannot be performed, the REQ fails as `error_type:
  environment` rather than being quietly closed — which is the outcome UR-015's note was protecting
  against.
- **None of these five REQs is `maintenance: true` except REQ-082.** Three are Go defects and one is
  a test execution; only REQ-082's candidate fix narrows the skill's own operating instructions,
  which is the marker's stated trigger.
- **The Simple REQ template's stray trailing line was omitted from all five files.** It is a known
  defect with its own ticket (REQ-080, from UR-015) and reproducing it into five new REQs while a
  REQ to delete it sits in the same queue would be perverse.

## Full Verbatim Input

capture all seven as REQs and stop

(Preceded in the same turn by: `capture all seven as REQs and run them`, superseded by the above
mid-turn. The seven refer to the seven `do-work capture-request:` lines the triage report ended with.)

---
*Captured: 2026-08-03T17:09:21Z*
