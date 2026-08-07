---
id: REQ-114
title: The three remaining shell-logic extraction candidates, restated decay-free
status: completed
created_at: 2026-08-05T19:32:52Z
claimed_at: 2026-08-07T07:25:25Z
completed_at: 2026-08-07T07:29:24Z
route: A
user_request: UR-022
domain: general
prime_files: []
tdd: false
depends_on: []
maintenance: false
kb_status: pending
kb_entry: ""
related: [REQ-111, REQ-112]
batch: census-durable-findings
---

# The Three Remaining Shell-Logic Extraction Candidates, Restated Decay-Free

## What

Carries forward the three extraction candidates the shell-logic census ranked but never captured, restated so they survive without the census's line-number table. Each is a consolidation of a primitive that is currently copy-pasted across several action files. Candidate B was split out, approved, and delivered as REQ-121; **Candidates A and C are not approved work** and must each clear the floor constraints before becoming a change.

Every candidate below is described by **what to grep for**, not by line numbers. That is the point: the census's citations went stale within hours of a single merge to `actions/work-reference.md`, so a durable record has to name the search rather than the coordinate.

### Candidate A — one home for the merge-commit-aware diff

**The primitive:** `git rev-parse --verify -q '<sha>^2'` succeeding means the hash is a merge commit, in which case plain `git show` prints a combined diff that is usually empty — so the reader must use `git show --first-parent -m <sha>` instead.

**Find every site:** `grep -rn "verify -q" actions/` — at census time this returned 7 sites across `review-work.md`, `ai-report.md`, `present-work.md` (3 of them), `pipeline.md`, and `pipeline-reference.md`.

**Why it matters:** get it wrong and a reviewer reads an empty diff as an empty REQ and passes it. Seven hand copies of one primitive is also a standing violation of this repo's own rule that a prescribed-command fix be grepped across every action before being called fixed.

**Candidate shape:** `tools/checks/req-diff.sh <req-file>` resolving the REQ's diff, merge-aware, with the `<pre>..<merge_hash>` range handled too.

### Candidate B — one uncommitted-changes inventory, one REQ-association pass

> **Split out and approved as REQ-121 on 2026-08-06.** Kept here for the record; do not implement it under this REQ.

**The primitive:** `git rev-parse --git-dir` gate, then `git status --porcelain --untracked-files=all`, then M/A/D categorization, then four secret-shaped exclusion globs (`.env*`, `credentials*`, `*.pem`/`*.key`/`*.p12`/`*.pfx`, `*secret*`). Separately: glob archived REQs, read `commit:` and a terminal-success `status`, parse `## Implementation Summary` file lists, path-match, tie-break on latest `completed_at`.

**Find every site:** `grep -rln "untracked-files=all" actions/` and `grep -rln "Implementation Summary" actions/ | xargs grep -ln "completed_at"`. At census time `commit.md` and `inspect.md` carried a near-verbatim copy of both halves, and `tidy-repo.md` and `validate-feedback.md` carried parts of the first.

**Why it matters:** the `-uall` flag is load-bearing, not cosmetic — without it `git status --porcelain` collapses a new directory to `?? dir/` and every file inside escapes the secret-exclusion scan. That is a secret-leak path, and `stray-check.md`'s Red Flags record that it has been hit. Both `commit.md` and `inspect.md` currently spend a paragraph *explaining* the flag, which is prose doing a script's job.

**Candidate shape:** `tools/checks/uncommitted-inventory.sh` plus `tools/checks/associate-files.sh`.

### Candidate C — writer-label claim classification

**The primitive:** derive this checkout's identity (`hostname -s`, falling back to plain `hostname`, plus `git rev-parse --show-toplevel`), compare it against each `## In Progress (interrupted)` entry's `writer:` label in `do-work/CHECKPOINT.md`, and classify each `working/` REQ as own-crash / foreign-label / label-less.

**Find the site:** `grep -n "writer:" actions/work-reference.md` — the Crash Recovery section. `queue-kanban verify` already holds adjacent pieces (a 3-hour stale-claim threshold, a `REQ-\d+` checkpoint-mention scan) but never derives the label.

**Why it matters:** lowest frequency of the three, highest blast radius. Guessing wrong runs recovery's strip of thirteen generated sections against a REQ another checkout may be building right now. `work-reference.md` says so itself — reading local checkpoint state as authorship "strips a live foreign claim." This is the one place in the census where a prose slip destroys work rather than misreporting it.

**Candidate shape:** `tools/checks/classify-claims.sh`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Re-run the durable discovery searches, verify Candidate B's separate delivery, and update only the census's concise disposition note. Verify: Candidate A's merge-aware-diff copies remain discoverable by grep, Candidate B's two scripts and callers exist, and Candidate C remains a single crash-recovery classification surface.
- [x] **[APPLY]:** Updated the audit and this request's status wording to distinguish Candidate B's completed REQ-121 delivery from still-unapproved Candidates A and C. No extraction, script, or action prose was changed.
- [x] **[UNIFY]:** Reviewed the audit and request status wording: both preserve the no-table approach and state no new implementation as fact. Ran `git diff --check`, `bash -n _dev/tests/contract-regressions.sh`, and `bash _dev/tests/contract-regressions.sh` (all pass). Verified the Candidate B inventory script still reports the current changed paths, and re-ran the Candidate A/C discovery searches. No debug artifacts or unrelated source changes.

## Why (if provided)

The census's two durable findings are implemented (REQ-111, REQ-112). These three are what was left, and they were deliberately not captured at the time because each depends on call sites still being where the table said — the perishable half. Restating them as greps rather than coordinates is what makes them keepable.

## Detailed Requirements

- Treat each candidate as **its own** unit of work. Do not implement all three under this REQ; split per candidate when one is approved, so each gets its own review and its own floor decision.
- **Re-run each candidate's grep before implementing it.** The site counts above are as-of-census figures and are explicitly not to be trusted — a merge to `actions/` between then and now may have added or removed sites.
- Every candidate must keep its action's shell floor. A `tools/checks/*.sh` script is shell and needs no exception; a `queue-kanban` subcommand does, and only in the accelerator shape (gated on an already-built binary, prose fallback documented).
- Candidate B carries a security dimension the other two do not: the `-uall` behaviour is what keeps secret-shaped files inside the exclusion scan. Any consolidation must preserve it and should test it.

## Constraints

- No action prose may lose its documented fallback procedure.
- `actions/board.md` remains the only capability allowed to *need* a compiler.

## Builder Guidance

**Exploratory.** Whether any of these three is worth doing is an open question, not a settled plan — the census ranked them by frequency × bug risk and stopped there. Candidate B has the clearest safety argument; Candidate A has the highest copy count; Candidate C has the worst failure mode. A reasonable outcome is to do one and re-evaluate.

## Red-Green Proof

Not applicable as written — this REQ carries three candidates rather than one behaviour change, and each will state its own proof when it is split out and approved. Recording that explicitly rather than inventing a proof target for work nobody has approved.

**Validation:** User confirmed the *intent* (move the census's value into the queue); the candidate contents are restated from the audit, not separately confirmed.

## Grep Refresh — 2026-08-06

Re-ran all three greps per the Detailed Requirements. The census figures were stale, as predicted:

| Candidate | Census | Now | Movement |
|---|---|---|---|
| A — merge-aware diff (`grep -rn "verify -q" actions/`) | 7 sites | **8 sites** | `actions/work.md` gained one; the other 7 hold (`present-work.md` ×3, `review-work.md`, `ai-report.md`, `pipeline.md`, `pipeline-reference.md`) |
| B — inventory (`grep -rln "untracked-files=all" actions/`) | 4 files | **5 files** | `validate-feedback.md` dropped out; `stray-check.md` and `work.md` joined. `commit.md` + `inspect.md` still carry the word-for-word copy |
| C — writer label (`grep -n "writer:" actions/work-reference.md`) | 1 site | **1 site** | Unchanged. The `hostname -s` derivation is at line 463; no second site anywhere in `actions/` |

**One constraint is now cheaper than the REQ assumed.** `tools/checks/` already ships six scripts (`archive-collision.sh`, `blanked-req-scan.sh`, `preflight.sh`, `qualify.sh`, `record-commit-hash.sh`, `scope-drift.sh`) and `_dev/tests/contract-regressions.sh` already pins each to its calling action. So the `tools/checks/*.sh` shape is an established, floor-compatible pattern for all three candidates — no compiled-tooling exception is in play for any of them, and the floor decision the REQ reserves is settled in advance.

## Full Context

See `do-work/user-requests/UR-022/input.md`. The originating audit is `decisions/audits/2026-08-05-shell-logic-in-prose-census.md`, trimmed in the same change that created this REQ.

---

## Triage

**Route: A** - Simple

**Reasoning:** This is a bounded documentation close-out: its durable candidate record already exists, Candidate B has a separately approved delivery, and the two remaining candidates explicitly lack implementation approval. The only source change is a concise status correction in the originating audit.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- **D-01: Close the inventory without implementing Candidate A or Candidate C.** DECIDE & STATE. The REQ and its source audit explicitly reserve both as unapproved work; treating a queue run as permission to choose an extraction would violate that boundary. Candidate B is already independently delivered by REQ-121.

## Implementation Summary

**Files changed:**
- `decisions/audits/2026-08-05-shell-logic-in-prose-census.md` (modified)

**What was done:** Recorded the durable current disposition: Candidate B is delivered by REQ-121, while Candidates A and C remain separate, unapproved candidates. Corrected this request's now-false original "None of these" status wording; no extraction implementation or action-prose change was made.

## Qualification

Passed — 1 project file verified in the diff, the P-A-U record is complete, and the requirements trace directly to the audit's corrected candidate disposition. The change is substantive documentation: it removes the now-false "None is approved" claim without reviving the perishable census table.

## Testing

**Tests run:** `git diff --check`; `bash -n _dev/tests/contract-regressions.sh`; `bash _dev/tests/contract-regressions.sh`; `tools/checks/uncommitted-inventory.sh .`; Candidate A/C discovery searches.

**Result:** All checks pass. The regression suite reports no failures; the inventory script still emits the individual changed paths; Candidate A's search still finds the expected merge-aware-diff pattern, and Candidate C's classification surface remains in Crash Recovery / In-Progress Record.

*Verified by work action*

## Review

**Overall: 100%** | 2026-08-07T07:25:25Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
None

**Minor findings:** 0 (report only)
**Acceptance:** Pass — the audit and REQ now distinguish Candidate B's completed REQ-121 delivery from unapproved Candidates A and C, without reviving the decaying table or implementing an unapproved extraction.
**Suggested testing:** 0 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Keeping the candidate record grep-based made it possible to verify the current state without reviving the census's stale line-number table.

**What didn't:** The inventory's original blanket statement that none of its candidates was approved became stale after Candidate B split into REQ-121; a disposition close-out needs to update that statement in both the audit and the REQ.

**Worth knowing:** A queue run authorizes processing the inventory REQ, not selecting an unapproved candidate for implementation. Candidate A and Candidate C remain separate decisions.

## Orientation

The shell-logic census now records its residual candidate status accurately: Candidate B is shipped under REQ-121, while Candidate A and Candidate C remain explicit future decisions. The durable candidate record lives in `decisions/audits/2026-08-05-shell-logic-in-prose-census.md`; no system contract changed.
