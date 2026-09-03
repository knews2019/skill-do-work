---
id: REQ-560
title: '[impact-rule-change] Hand-back and finalize check cleanliness only on the REQ''s own paths'
status: pending
created_at: 2026-09-03T20:05:46Z
user_request: UR-106
domain: backend
prime_files: [_dev/primes/prime-action-files.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-531, REQ-559, REQ-503]
batch: lifecycle-overhead
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/do-work-cli/internal/finalization/
---

# Hand-Back and Finalize Check Cleanliness Only on the REQ's Own Paths

## What

A path the active REQ does not own, whether untracked, modified, or staged by another session, is never a reason to stop, never surfaced as a blocker, and never committed by the pipeline. Step 7's hand-back settlement and Step 9's finalization check the index and tree only for the REQ's own paths: its run artifacts, its lifecycle files, its write set, and its release paths. Everything else is left exactly as found and named in one progress line. The "preserve in a separate unrelated-work commit" behaviour goes away.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

"Finalization was briefly blocked by an unrelated existing calibration change, which had to be preserved separately." The pipeline authored two commits of other sessions' work to get past its own cleanliness check: 061b7dbf "Preserve corrected calibration history" and 83594c5e "Preserve concurrent maintainability audit draft".

## Context

- `work-reference.md` already states the rule as judgment ("Current-REQ relevance": preserve it, exclude it from this REQ's staging, continue, spend no time on it). Step 7 step 0 contradicts it mechanically: "stop and surface every other `do-work/` path" and "step 0 ends with a clean index". Under four concurrent sessions that condition is rarely true, so the orchestrator invents a way to make it true, which on 2026-09-03 meant committing another session's untracked audit draft under the pipeline's own name.
- The claim transaction already treats "a dirty claim target or index" as shared-target dirt with a typed refusal. This REQ keeps that for the REQ's own target and the index, and drops it for everything else.
- Committing a foreign file is worse than leaving it: it strips the owner's chance to finish it, attributes it to the wrong REQ, and can land a half-written draft on main. Leaving it costs nothing.
- REQ-503 to REQ-510 (the advance command) will absorb the hand-back mechanics later; this REQ changes the rule now so the chain inherits it.

## Detailed Requirements

- Step 7 step 0: the stage category and the allow category stay; the "stop and surface" category becomes "leave alone and name": the path is listed once in the progress output and excluded from staging. Only a dirty path the REQ itself owns still stops.
- The index must be clean of the REQ's own paths before the merge, not of every path; a foreign staged path is unstaged from the index only if it was staged by this run, otherwise it is left and named.
- Finalization: the manifest's exact-allowlist validation is unchanged; a foreign modified or untracked path outside the manifest never refuses the transaction. If a do-work-cli command today refuses on tree dirt outside its declared paths, narrow that check to its declared paths and pin the narrowing with a test in that package.
- Delete the sentences that authorize or describe a pipeline-authored "unrelated work" or "preserve" commit; the pipeline commits only what the REQ declares.
- One line in the run's progress output per foreign path left alone, so the maintainer sees what was skipped without the run stopping.

## Constraints

- Mechanics in programs, judgment in prose (CLAUDE.md). A do-work-cli change carries its own test; no sentence pins.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.
- The claim transaction's refusal on a dirty claim target is unchanged.

## Red-Green Proof
**RED prompt/case:** With an untracked file from another session under `do-work/audits/` and a modified `do-work/calibration-log.tsv` that this REQ did not touch, run a REQ through Step 7 hand-back and Step 9 finalize.
**Why RED now:** step 0 stops on the foreign paths or the orchestrator commits them under the pipeline's name to reach a clean index (061b7dbf, 83594c5e on 2026-09-03).
**GREEN when:** hand-back and finalize complete; the foreign paths are still untracked or modified in the working tree exactly as before; every commit the run made contains only the REQ's declared paths; the progress output names each foreign path once; and a dirty path the REQ itself owns still stops with the typed refusal.
**Validation:** Inferred during capture; the maintainer approved the capture ("do 1, 2 and 3").

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes a pipeline step contract.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 5660 tokens, over budget and `slugged: partial`. Matched because this REQ may change finalization validation in do-work-cli internals.

## Full Context
See `do-work/user-requests/UR-106/input.md` for complete verbatim input.
