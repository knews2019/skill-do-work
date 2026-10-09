---
id: REQ-672
title: 'validate-feedback --capture captures the accepted findings as one UR with one REQ per finding and its provenance'
status: pending
created_at: 2026-10-09T21:26:55Z
user_request: UR-149
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: [REQ-671]
related: [REQ-670, REQ-671, REQ-673, REQ-674]
batch: validate-feedback-capture
---
# validate-feedback --capture Captures the Accepted Findings as One UR
## What
With `--capture`, `validate-feedback` builds one capture payload from the Accept items plus the Discuss items the user accepted (REQ-671), and runs `skills/do-work/actions/capture.md` on it once. Result: one UR, one REQ per finding, each REQ citing its finding id and its source (pasted text, or the file path from REQ-670) in the provenance block the handoff already asks for. Findings that duplicate a queued REQ fold through capture's existing fold-first scan.
## Why
Report item C3 (UR-149 input): users retyped accepted finding ids into a capture by hand after each triage ("capture F1, F2, F6, F9 and C1" on 2026-09-17; "capture the ones that are accepted" on 2026-10-02; "capture to fix accepted" on 2026-10-07).
## Verified Facts (checked at 0.305.87)
- `skills/do-work-toolbox/actions/validate-feedback.md:113-117`: the handoff already says to keep "the accepted finding's **verbatim claim**, **original severity/source**, **Evidence**, and **Surface-cost** result together in the capture payload", then prints `do-work capture-request:` and `do-work run` as lines for the user to type.
- `skills/do-work/actions/capture.md:119`: "A review, build, triage, or consumer-report finding reaches capture only when the user invokes `do-work capture` and quotes the complete report-only finding line as the source". `--capture` is a user invocation; the fold-first scan named on the same line still runs.
- `skills/do-work/actions/capture-reference.md` → Fold-First Rule: an explicit capture enters the scan and folds into a matching queued REQ instead of minting a new one.
## Detailed Requirements
1. Add a new Step 7, capture, that runs only with `--capture`, after REQ-671's Discuss questions.
2. Build one payload from the Accept items plus the Discuss items the user accepted, each as the existing `:113-114` provenance block (verbatim claim, original severity/source, Evidence, Surface-cost).
3. Run `skills/do-work/actions/capture.md` on that payload once. Result: one UR, one REQ per finding.
4. Each REQ's source line names the finding id (for example `F3`, "Finding 3") and the input kind: pasted text, or the file path from REQ-670.
5. A finding that duplicates a queued REQ folds into it through the existing fold-first scan; it does not mint a new REQ.
6. If the set is empty (no Accept items and no accepted Discuss items), say so and stop without capturing.
7. Edit the handoff block at `:113-117` so it describes what happens with `--capture` and keeps the typed lines for the no-flag case.
8. Release per `_dev/primes/prime-releases.md`.
## Constraints
- Capture ≠ Execute: without `--capture` or an explicit phrase, nothing is captured, for any verdict, including impact-critical ones.
- No change to capture's REQ templates, numbering, reservations or transaction. No new REQ fields, no new statuses.
- Out of scope: a `--review N` front end; changes to verdict rules or to how `do-work run` selects work.
## Assumptions (recorded at capture, no questions asked)
- `--capture` counts as the user invocation `capture.md:119` requires, because the user typed the flag. If the builder finds that line's wording would refuse a capture started from validate-feedback, add one clause naming `validate-feedback --capture` as such an invocation; otherwise leave `capture.md` unchanged.
- Capture's own Step 3 clarification still applies, but the Discuss decisions from REQ-671 count as resolved and are not asked again.
- An empty set with `--run` also starts no run (REQ-673 never runs).
- The UR's verbatim input is the triage input (the paste, or the file's bytes), with the accepted subset and the per-finding verdicts recorded in the UR Summary; capture's existing raw-containment rule decides the exact bytes.
## Dependencies
Depends on REQ-671 (the Discuss answers that decide which Discuss items are captured). REQ-673 (verify, then optionally run) depends on this REQ.
## Builder Guidance
Certainty is high on outcome (one UR, one REQ per finding, provenance per REQ, folds). Latitude: how the payload is laid out and how the finding id is phrased in each REQ's source line.
## Red-Green Proof
**RED prompt/case:** `do-work-toolbox validate-feedback --capture` on a paste with 2 Accept, 1 Discuss and 1 Push back items; answer Accept to the Discuss question.
**Why RED now:** `validate-feedback.md:122` forbids creating REQs, and the handoff at `:115` only prints a capture line for the user to type.
**GREEN when:** exactly one new UR with 3 REQs exists, each REQ naming its finding id and source. A fourth finding that duplicates a queued REQ folds into it and appears under `## Folded Requests`. A paste with no Accept items and every Discuss item dropped captures nothing and says so.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers alternate artifact writers and downstream readers; this REQ makes a toolbox action start capture.
## Full Context
See `do-work/user-requests/UR-149/input.md` for complete verbatim input (sections Request C3, What happened, Where the behaviour lives today, Proposed direction C3, Acceptance check). No queued candidate shares this root cause.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-validate-feedback-capture.md`, Request item C3: "Capture the Accept items plus the Discuss items the user accepted as one UR (one user request record) with one REQ per finding. Each REQ cites its finding id and the source (pasted text, or the path of a review or audit file) in the provenance block the handoff already asks for (`validate-feedback.md:114`)."*
