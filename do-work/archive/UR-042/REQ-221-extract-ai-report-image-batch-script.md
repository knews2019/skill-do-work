---
id: REQ-221
title: Extract the ai-report image batch into a shipped script
status: completed
claimed_at: 2026-08-17T20:56:14Z
status_changed_at: 2026-08-17T19:58:23Z
route: B
completed_at: 2026-08-17T21:21:18Z
commit: e2b45bf
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-17T20:56:14Z
  basis:
  - Route B
  - 5-file write set
  - 1 new files
  - 2 subsystems involved
  - 5 acceptance criteria
  - async lifecycle behavior
  - full-suite verification
domain: general
created_at: 2026-08-17T18:37:31Z
user_request: UR-042
addendum_to: REQ-204
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-action-files.md]
tdd: true
maintenance: true
---

# Discovered Task: Extract the AI-Report Image Batch Into a Shipped Script

## What

Move the prescribed image-batch block out of `ai-report-reference.md` and into `<skill-root>/scripts/generate-report-image-batch.sh`, leaving the action file with a pointer plus the per-report prompts.

## Context

Found while implementing REQ-204. After adding process-tree ownership and publication verification, the fenced block is roughly 110 lines of purely mechanical shell — staging, job control, group verification, signal handling, wait-all, freshness evaluation, publication, rollback — living inside a markdown action file. The only per-report content in it is the two `launch_report_image … "<prompt N>"` lines.

The test suite already treats it as a program: `_dev/tests/prescribed-shell-scripts-behavior.sh` extracts the block with `awk`, rewrites `<skill-root>`, and executes it. That extraction step exists only because the code is embedded in prose.

## Requirements

- The script owns the mechanics: staging, launch, isolation verification, signal ownership, wait-all, per-status freshness, publication, and rollback.
- The action file keeps only what is per-report — the style brief, the prompts, and how the results are referenced from the HTML — plus a pointer to the script.
- Preserve every behavior REQ-198 and REQ-204 locked in; the four existing replays must pass against the script with their assertions unchanged in meaning.
- The replays call the script directly; the `awk` block-extraction harness goes away with the block it was extracting.
- Keep the shipped-package manifest, `suite/modules.tsv`, and the shipped reference contract consistent with the new script.

## Red-Green Proof

**RED prompt/case:** Point the four existing batch replays at a shipped `generate-report-image-batch.sh`; there is no such script, so they cannot run.
**Why RED now:** The mechanics have no executable home, so the only way to test them is to reconstruct them from prose at test time.
**GREEN when:** All four replays invoke the script directly and pass, the action file no longer contains a fenced batch implementation, and the full maintainer baseline is green.
**Validation:** Discovered task; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] The image-batch mechanics are now ~110 lines of shell embedded in a markdown action file, and the tests already have to extract them to run them. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a structural move, not a defect fix — nothing is broken today, so it is a judgment call about where the mechanics should live rather than something the builder should decide alone. It also touches the shipped package inventory, which is a user-visible surface.

<!-- D-XX counter: last used D-04. Next decision: D-05. -->

---

## Triage

**Route: B** - Medium

**Reasoning:** The move itself is well specified — the block, its destination, and the four replays that must keep passing are all named. What needed discovery was the consistency surface: which manifests, contract suites, and shipped inventories enumerate a shipped script, and which existing contract assertions were written against the block's *address* rather than its behavior.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**What the block actually was.** `skills/do-work-toolbox/actions/ai-report-reference.md` carried ~110 fenced lines under "**Fire in parallel, retain every status, then verify.**" — staging, job control, process-group verification, signal ownership, wait-all, per-status freshness, publication with post-rename nesting verification, and rollback. The only per-report content was two `launch_report_image <target> "<prompt N>"` lines.

**How the tests reached it.** `_dev/tests/prescribed-shell-scripts-behavior.sh` extracted the block with `awk`, rewrote `<skill-root>` with `sed`, wrote it to a scratch file, and ran it with `bash`. That harness existed only because the code lived in prose.

**The consistency surface (four registries, all of which needed the new script):**
- `_dev/tests/staged-skills-contract.sh` — a `toolbox_files` array enumerating every shipped toolbox file, plus a `sibling_script_contract` table pairing each `<skill-root>/scripts/…` reference in an action file with the real path it must resolve to.
- `_dev/tests/prescribed-shell-canonicalization.sh` — a list of prescribed scripts to canonicalize, and a list of required headings in `skills/do-work/docs/prescribed-shell-primitives.md`.
- `skills/do-work/docs/prescribed-shell-primitives.md` — the shipped "Shipped executable homes" table and per-primitive sections.
- `_dev/tests/contract-regressions.sh` — twelve assertions that pinned REQ-198 and REQ-204 guarantees *by matching text inside the markdown block*. These are the subtle ones: they had to be repointed at the script rather than deleted, or the move would silently un-ratchet both prior REQs.

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/scripts/generate-report-image-batch.sh` (new) — the extracted mechanics behind an argument interface
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modify) — fenced implementation replaced by a pointer plus the per-report invocation
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — replays call the script directly; `awk` harness deleted
- `_dev/tests/contract-regressions.sh` (modify) — REQ-198/204 assertions repointed from the block to the script
- `_dev/tests/staged-skills-contract.sh` (modify) — register the new shipped file and its sibling-reference contract
- `_dev/tests/prescribed-shell-canonicalization.sh` (modify) — add the script and the new guide heading
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — executable-home row and a section for the batch primitive

**Files I will NOT touch:** `skills/do-work-toolbox/scripts/generate-report-image.sh` (the per-image helper is unchanged; the batch calls it), `suite/modules.tsv` (declares packages, not individual scripts — a new file under an already-declared directory needs no row).

**Acceptance criteria (restated from REQ):**
- [x] The script owns staging, launch, isolation verification, signal ownership, wait-all, per-status freshness, publication, and rollback
- [x] The action file keeps only per-report material plus a pointer, and contains no fenced batch implementation
- [x] Every behavior REQ-198 and REQ-204 locked in is preserved; the four replays pass with their assertions unchanged in meaning
- [x] The replays call the script directly and the `awk` block-extraction harness is gone
- [x] The shipped-package manifest, `suite/modules.tsv`, and the shipped reference contract stay consistent

## Pre-Flight

**Repository state:** OK — clean outside `do-work/` at claim (REQ-220 committed at `1fb3635`/`a980e4f`)
**Test baseline:** OK — `bash _dev/tests/maintainer-verify.sh` exit 0 before implementation
**Dependencies:** OK

## Decisions

- **D-01 — DECIDE & STATE.** The batch's success signal is the published directory printed on **stdout**; per-image `MISSING:` and `REFUSING:` diagnostics moved to **stderr**. Reasoning: the block set a caller variable (`GEN=…`), which a script cannot do, so the caller now does `GEN="$(…)"`. That makes stdout a value channel, and a diagnostic written there would be captured into the path. The `MISSING:` line was also shortened to the target's basename — the full staged path names an invocation-private directory the reader cannot act on.
- **D-02 — DECIDE & STATE.** The per-image helper is resolved as a sibling of the script (`$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)`), not through `<skill-root>` or `PATH`, matching how `install-last30days.sh` reaches `add-local-git-exclude.sh`. No override environment variable was added — the replays inject their fake backend through `PATH` on `imagegen`, one level below, so nothing needed it (YAGNI).
- **D-03 — DECIDE & STATE.** `ps` is resolved by absolute path with the no-`ps` case degrading to bare-PID signalling, carrying REQ-220's lesson across rather than reproducing the block's bare `ps`. Reasoning: a caller may hand the batch a minimal `PATH`, where a bare `ps` fails silently and quietly downgrades every launch's isolation.
- **D-04 — DECIDE & STATE (orchestrator, at review).** Restored one REQ-198 assertion the repointing dropped rather than moved: `ai-report-reference.md must not publish generated/ before a current image succeeds`. Reasoning: the other eleven assertions describe mechanics that genuinely relocated to the script, but this one guards the *action file* against re-acquiring a naive `mkdir -p "$GEN"` — delegation does not prevent that, so the guard still has a live failure mode at its original address. Mutation-tested: reinserting the offending line into the action file makes it fail, and removing it makes the suite pass again.

## Implementation Summary

**What was done:** Moved the ai-report image batch out of prose and into a shipped executable home. `generate-report-image-batch.sh` takes the report directory, the shared style brief, and one `<target-name>:<prompt>` pair per section (split on the first colon so prompts may contain colons; a target containing `/` is a usage error), and owns every mechanic the fenced block owned. The action file now carries the invocation and the per-report material only. The four batch replays call the script directly, and the `awk`/`sed` harness that reconstructed the block at test time is deleted along with the block it extracted.

**Files changed:**
- `skills/do-work-toolbox/scripts/generate-report-image-batch.sh` (new) — 189 lines; staging, per-image launch under verified process-group isolation, signal ownership, wait-all with retained statuses, per-status freshness, verified publication, rollback
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified) — ~137 lines of fenced implementation replaced by the invocation, the stdout/stderr contract, and a pointer to the shipped guide
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified) — `awk` extraction harness removed; four replays repointed at the script; usage-validation case added (44 → 45 named cases)
- `_dev/tests/contract-regressions.sh` (modified) — eleven REQ-198/204 assertions repointed from the markdown block to the script source, one added for the nesting verification, one added proving the action file no longer restates the mechanics, one restored (D-04)
- `_dev/tests/staged-skills-contract.sh` (modified) — new script added to the shipped toolbox inventory and to the sibling-reference contract table
- `_dev/tests/prescribed-shell-canonicalization.sh` (modified) — new script added to the canonicalization list; `## Report image batch publication` added to the required-headings list
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified) — executable-home row plus a section documenting the batch primitive and its stdout/stderr contract

## Qualification

Passed — 7 files verified in the diff, 5 requirements traced, no debug artifacts. Qualified against git state, not the builder's report.

The scope question worth naming: the builder changed four `_dev/tests/` files I did not enumerate up front, because I asked it to investigate the consistency surface and change what was actually required. Each was checked individually and each is required — three are registries that must list a new shipped script, and the fourth (`contract-regressions.sh`) held assertions written against the block's address. That last one is where a careless extraction would have quietly un-ratcheted REQ-198 and REQ-204: deleting those assertions would leave a green suite and no lock-in. They were repointed at the script, so the guarantees survive the move; one that belonged at the old address was restored (D-04).

## Testing

**Commands run:**
- `bash _dev/tests/prescribed-shell-scripts-behavior.sh` — exit 0 (45 named script cases)
- `bash _dev/tests/contract-regressions.sh` — exit 0
- `bash _dev/tests/maintainer-verify.sh` — exit 0, zero FAIL lines

**Red-green validation** (traced to `## Red-Green Proof`; RED re-derived by the orchestrator by moving the script aside and re-running the suite, not taken from the builder's report):

RED (exit 1) — every replay fails at `127` (command not found), which is exactly the captured proof "there is no such script, so they cannot run":
- `FAIL: ai-report all-failed batch replay returned nonzero instead of falling back`
- `FAIL: ai-report mixed batch replay did not publish only the status-backed successful image`
- `FAIL: ai-report mixed batch replay did not report the published generated/ directory on stdout`
- `FAIL: ai-report interrupted batch replay exited 127 instead of the TERM status 143`
- `FAIL: ai-report publish-collision replay did not preserve the colliding destination byte-for-byte`
- `FAIL: ai-report publish-collision replay left its staged batch nested inside the colliding destination`
- plus the two new usage cases (`returned 127 instead of the usage status 2`)

GREEN: with the script restored the suite exits 0 at 45 named cases.

**Separate mutation proof for D-04:** reinserting `GEN="ai-reports/<report-slug>/generated"; mkdir -p "$GEN"` into the action file makes `contract-regressions.sh` exit 1 with `FAIL: ai-report-reference.md must not publish generated/ before a current image succeeds.`; removing it returns the suite to exit 0. The restored assertion pins a real failure rather than decorating the file.

**Existing tests updated:** twelve `contract-regressions.sh` assertions repointed (behavior unchanged, address changed); the `awk` extraction harness deleted.

## Review

**Requirements check — Pass.** All five requirements trace to changes and are ticked in Scope. The one that carried real risk — "preserve every behavior REQ-198 and REQ-204 locked in" — was verified assertion by assertion rather than by the suite being green, because a deleted assertion also leaves the suite green.

**Code review — Pass.** The script is a faithful move: the retained-PID/retained-status structure, the fail-closed `generated/` checks, the nesting verification, and the reap-before-cleanup trap ordering are all carried across intact. Two improvements arrived with it rather than being invented: absolute `ps` resolution (REQ-220's lesson, which the original block predates) and the stdout/stderr split forced by turning a caller variable into a return value. Argument validation rejects an unpaired argument and a target containing `/` with exit 2, matching the sibling scripts' usage convention. Empty-array expansion under `set -u` was checked against REQ-216's macOS bash 3.2 lesson: every array access is index-bounded by its own count, and the argument floor guarantees at least one pair, so no bare `"${array[@]}"` expansion of an empty array occurs.

**Restatement sweep — clean after the fact.** The action file's own prose was rewritten as part of the move; `prescribed-shell-primitives.md` gained the matching primitive section in the same change; no other file restates the batch mechanics (verified by grep for `launch_report_image`, `image_generation_stage`, `Fire in parallel`).

**Acceptance = Pass. Overall: 96%.** Scope drift: four `_dev/tests/` files beyond the initial enumeration, all required by the consistency investigation the REQ mandates — Minor.

## Lessons Learned

**What worked:** Treating "which assertions were written against the code's *address* rather than its *behavior*" as the first question of any extraction. Twelve `contract-regressions.sh` assertions matched text inside the markdown block; deleting them would have left a green suite with two REQs' guarantees silently un-ratcheted. Repointing them at the script is what makes the move behavior-preserving rather than merely test-passing.

**What didn't:** Assuming a green suite proves an extraction preserved its lock-ins. It cannot — removing an assertion and removing the behavior it guards look identical from the exit code. The check that works is per-assertion: for each one, decide whether it moved with the code or belonged at the old address. One belonged at the old address (D-04) and had been dropped.

**Worth knowing:** A prescribed block that the test suite has to `awk` out of prose to execute is already a program with no home — the harness is the tell. Also, extracting a block into a script converts every caller variable it set into a return-channel problem; here `GEN=…` became stdout, which forced every diagnostic onto stderr. That is a real interface decision, not a mechanical translation.

## Orientation

The ai-report image batch now has an executable home: `generate-report-image-batch.sh` takes a report directory, a style brief, and `<target>:<prompt>` pairs, and prints the published `generated/` directory. The action file states what to draw; the script owns how the batch runs. Lives in the do-work-toolbox shipped-scripts subsystem, documented in `skills/do-work/docs/prescribed-shell-primitives.md` → **Report image batch publication**. **[MAP CHANGED]** — this adds a shipped executable and moves a contract boundary: the caller's relationship to the batch is now an argument interface and a stdout value rather than inlined shell sharing the caller's variables, and four registries (shipped inventory, sibling-reference contract, canonicalization list, executable-homes table) now name it. Prime spot-check: `_dev/primes/prime-shell-commands.md` and `_dev/primes/prime-action-files.md` referenced paths all still resolve; neither describes the batch by address, so neither went stale.

## Discovered Tasks

None.
