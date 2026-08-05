---
id: REQ-114
title: The three remaining shell-logic extraction candidates, restated decay-free
status: pending
created_at: 2026-08-05T19:32:52Z
user_request: UR-022
domain: general
prime_files: []
tdd: false
depends_on: []
maintenance: false
related: [REQ-111, REQ-112]
batch: census-durable-findings
---

# The Three Remaining Shell-Logic Extraction Candidates, Restated Decay-Free

## What

Carries forward the three extraction candidates the shell-logic census ranked but never captured, restated so they survive without the census's line-number table. Each is a consolidation of a primitive that is currently copy-pasted across several action files. **None of these is approved work** — this REQ exists so the candidates are in the queue rather than in a decaying document, and each must clear the compiled-tooling and floor constraints before it becomes a change.

Every candidate below is described by **what to grep for**, not by line numbers. That is the point: the census's citations went stale within hours of a single merge to `actions/work-reference.md`, so a durable record has to name the search rather than the coordinate.

### Candidate A — one home for the merge-commit-aware diff

**The primitive:** `git rev-parse --verify -q '<sha>^2'` succeeding means the hash is a merge commit, in which case plain `git show` prints a combined diff that is usually empty — so the reader must use `git show --first-parent -m <sha>` instead.

**Find every site:** `grep -rn "verify -q" actions/` — at census time this returned 7 sites across `review-work.md`, `ai-report.md`, `present-work.md` (3 of them), `pipeline.md`, and `pipeline-reference.md`.

**Why it matters:** get it wrong and a reviewer reads an empty diff as an empty REQ and passes it. Seven hand copies of one primitive is also a standing violation of this repo's own rule that a prescribed-command fix be grepped across every action before being called fixed.

**Candidate shape:** `tools/checks/req-diff.sh <req-file>` resolving the REQ's diff, merge-aware, with the `<pre>..<merge_hash>` range handled too.

### Candidate B — one uncommitted-changes inventory, one REQ-association pass

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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

## Full Context

See `do-work/user-requests/UR-022/input.md`. The originating audit is `decisions/audits/2026-08-05-shell-logic-in-prose-census.md`, trimmed in the same change that created this REQ.
