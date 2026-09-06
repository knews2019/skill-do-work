---
id: REQ-554
title: '[impact-negligible] Move the 46 lines commit.md and inspect.md share into the prescribed-shell guide'
status: claimed
priority: later
created_at: 2026-09-03T19:45:35Z
user_request: UR-105
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
related: [REQ-549, REQ-550, REQ-551, REQ-552, REQ-553, REQ-555, REQ-556, REQ-557, REQ-558]
batch: maintainability-audit-2026-09-03
maintenance: true
impact: impact-negligible
effort_estimate: effort-mechanical
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-09-06T00:43:32Z
  basis:
    - Route B
    - 5-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set: [skills/do-work/actions/commit.md, skills/do-work-toolbox/actions/inspect.md, skills/do-work/docs/prescribed-shell-primitives.md, _dev/tests/prescribed-shell-canonicalization.sh, _dev/tests/audit-lockins.sh]
route: B
claimed_at: 2026-09-06T00:38:56Z
---

# Move the 46 lines commit.md and inspect.md share into the prescribed-shell guide

## What
`skills/do-work/actions/commit.md` and `skills/do-work-toolbox/actions/inspect.md` share 46 byte-identical non-blank lines in runs of three or more: the M/A/D/X/XD classification legend, four file-reading bullets, and two complete "If the script is missing or will not run, do it by hand" fallbacks that restate the algorithms `internal/corehelpers/inventory.go` and `associate-files` implement. Move the legend, the two fallbacks, and the bullets into one section of `skills/do-work/docs/prescribed-shell-primitives.md` (which already owns the protected-inventory heading), cite it from both actions, and keep only the read-only-versus-staging deltas local.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Two prose files stating one rule drift; this is the only non-crew-member prose pair with duplicated three-line windows in the audited surface, and no INLINE RESIDUE row in `decisions/audits/2026-08-11-prescribed-shell-primitives.md` covers it.

## Context
Source: `do-work/audits/audit-2026-09-03.md` (Finding 6, sweep_key `commit-inspect-shared-body`, audited commit dc8a64e3, report committed at 83594c5e). Plan tag JUDGMENT; expected net line delta -60. Captured from the audit's §Plan paste-ready line after the maintainer said "capture the requests"; the validator step was skipped on the maintainer's instruction, so the builder treats the finding's Reproduce output as the claim to re-verify at claim time.

## Detailed Requirements
- `commit.md` and `inspect.md` — `If the script is missing or will not run, do it by hand from the complete NUL-delimited output` (207 words each, two relative-path fixups differ).
- `commit.md` and `inspect.md` — `If the script is missing or will not run, do it by hand: glob both directories` (byte-identical).
- `commit.md` and `inspect.md` — `- **Modified files**: Read the \`git diff\` for each file.` and the three bullets after it.
- `commit.md` and `inspect.md` — the `M`/`A`/`D`/`X`/`XD` legend and `Secret-shaped matching is case-insensitive` (X/XD suffixes reworded for read-only mode: keep that delta local).
- The canonicalization ratchet counts headings in the guide and pointers to it from named files: re-baseline those counts in `_dev/tests/prescribed-shell-canonicalization.sh` in the same commit.
- Reproduce at dc8a64e3 (prints 46, then the four fallback sites): `python3 -c "import difflib;a=[l.rstrip() for l in open('skills/do-work/actions/commit.md')];b=[l.rstrip() for l in open('skills/do-work-toolbox/actions/inspect.md')];print(sum(1 for i,j,s in difflib.SequenceMatcher(None,a,b).get_matching_blocks() if s>=3 for k in range(s) if a[i+k].strip()))" && rg -n 'If the script is missing or will not run' skills --glob '*.md'`

## Constraints
- Scope is exactly this finding class: do not fix nearby code, do not extend behaviour the finding does not name, no test files beyond the lock-in.
- The lock-in lands as one assertion in `_dev/tests/audit-lockins.sh` (the file already exists, is executable, and is already registered in the fast tier at `_dev/tests/contracts/probe-lanes.sh` -- add one assertion to it; do not create it and do not change its registration), pinned at today's value so it is green on day one and red the moment the number regrows; no other test file changes.
- Prime `_dev/primes/prime-action-files.md` before editing either action; the guide edit follows `_dev/primes/prime-shell-commands.md`.
- Lock-in limit: commit.md/inspect.md shared lines ≤ 10 after this REQ (today 46); fallback sentences in actions: 0 (today 4).

## Dependencies
No dependency. REQ-555 (rewrite the guide's executable-homes table) depends on this REQ because it re-baselines the same guide and ratchet counts.

## Builder Guidance
Firm on one home for the shared text; latitude on the section title and on how the read-only delta is phrased.

## Red-Green Proof
**RED prompt/case:** Run the Reproduce command from Detailed Requirements.
**Why RED now:** It prints 46 shared lines and four fallback sites.
**GREEN when:** Shared lines are at most 10 and the fallback sentence appears 0 times in actions; the lock-in pins the difflib count at the post-fix value and the fallback sentence count at 0.
**Validation:** Inferred during capture from the audit report's Reproduce output; the maintainer approved the plan line without adjusting it.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3968 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing, pipeline fields, status contracts, downstream readers".
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing shipped shell, argv/quoting, prescribed command blocks, publication scripts".

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-105/input.md` for complete verbatim input.

---
*Source: `do-work/audits/audit-2026-09-03.md` §Plan, capture-request line for commit-inspect-shared-body.*

## Triage

**Route: B** — Explore then build.

**Reasoning:** The change itself is mechanical prose movement across five named files. What needed
discovery is that three of the request's own statements are wrong or unreachable, and a builder
following them literally writes an assertion that is red forever. Discovery, not design; the
exploration names every file, every line and the exact assertion text, so there is nothing left to
plan.

**Planning:** Skipped.

**`maintenance: true`.** This request edits shipped instruction files, so
`skills/do-work/crew-members/maintenance.md` loads for the builder.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Explore agent, read-only, re-verified against HEAD. Full report in the run directory as
`do-work/runs/work-2026-09-05-231943/REQ-554-exploration.md`.

**Three of the request's statements do not survive contact with HEAD, and one of them would produce a
permanently red gate.**

1. **The lock-in limit of "shared lines ≤ 10" is unreachable.** Simulated at HEAD: deleting every
   shared *content* line from both files still leaves 17, and all 17 are template scaffolding that
   `_dev/primes/prime-action-files.md:70-80` requires both actions to carry — `## When to Use`,
   `## Steps`, code fences, `## Error Handling`, the `| Situation | Action |` header, `## Rules`,
   `## Red Flags`. A ceiling of 10 cannot be met without violating the prime that mandates the
   scaffold. The prior exploration's suggested ceiling of 20 is also unreachable: keeping the ASCII
   flow diagram and the Error Handling table row — both structural, neither movable prose — the lowest
   reachable count is 21.
2. **`_dev/tests/prescribed-shell-canonicalization.sh` counts nothing.** The request says it "counts
   headings and pointers" and that those counts need re-baselining. The file has zero numeric
   constants: lines 66-83 are `grep -Fqx` membership checks over eleven hardcoded headings, and
   85-117 are `grep -Fq` membership checks over sixteen pointer sites. Adding a heading or a pointer
   cannot fail it. The only worthwhile edit there is adding the new section's heading to the
   membership list.
3. **The prior exploration's paste-ready assertion is green on day one and never fires.** Its
   `--glob '*/actions/*.md'` matches zero files, so its pipeline prints 0 today when the true count
   is 4.

**And one thing the request feared does not happen.** The canonicalization scan's stale-pattern check
skips the canonical guide itself (`prescribed-shell-canonicalization.sh:146`), so prose moved *into*
the guide is exempt; none of its nine stale patterns or seven old-implementation fragments matches any
moved line.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — one new `##` section holding the
  inventory tag legend, the secret-shaped matching sentence, the four file-reading bullets, the
  association semantics and both manual fallbacks, inserted before `## Merge-aware commit diff`
- `skills/do-work/actions/commit.md` (modify) — delete the moved prose and point at the new section;
  keep the two wordings that differ between the two actions rather than forcing them together
- `skills/do-work-toolbox/actions/inspect.md` (modify) — the mirror of the same deletion at its own
  line numbers, keeping its own read-only wordings
- `_dev/tests/prescribed-shell-canonicalization.sh` (modify) — one line added to the required-heading
  membership list; nothing to re-baseline, because nothing there counts
- `_dev/tests/audit-lockins.sh` (modify) — one new assertion block in the file's existing shape

**Files I will NOT touch:** every other action file. The two wordings that legitimately differ between
a writing action and a read-only one stay where they are — collapsing them would be a behaviour claim
dressed as deduplication.

**Acceptance criteria (restated from REQ):**
- [ ] The shared body lives once, in the prescribed-shell guide, and both actions point at it
- [ ] Each action keeps the wording that is genuinely its own
- [ ] A lock-in assertion fails if the shared prose comes back, and passes on the template scaffold
      both actions are required to carry
- [ ] The canonicalization probe accepts the new section
- [ ] The gate is green

## Decisions

- **D-01 — the lock-in ceiling is the measured post-move count, not the request's 10. DECIDE &
  STATE.** Ten is unreachable: seventeen identical lines survive after every shared sentence is gone,
  and all seventeen are scaffold `_dev/primes/prime-action-files.md` requires. Twenty is unreachable
  too, for two more structural lines. The ceiling is set to whatever the move actually lands on, with
  the reason written into the assertion's comment so the next reader sees why it is not zero.
  **Value:** the assertion fires when prose comes back and never fires on the scaffold. **Risk:** a
  ceiling set at the current value ratchets nothing further; that is the honest trade, and it is what
  the request wanted the number to mean.
- **D-02 — nothing is re-baselined in `prescribed-shell-canonicalization.sh`, because it counts
  nothing. DECIDE & STATE.** The request asks for a re-baseline there. The file has no numeric
  constants at all; the only change it needs is the new heading in its membership list.

## Pre-Flight

**Git:** ✓ Clean. Canonical `recover` reports `FINALIZATION-NONE`. Six claims sit in `do-work/working/`:
REQ-583, REQ-587, REQ-591 and REQ-592 held at Step 7.7 for the heavy drain; REQ-486 in flight on the
board module; REQ-552 merged and awaiting review. None of them touches this request's five files, and
REQ-552's own edit to `_dev/tests/audit-lockins.sh` is already merged, so this request appends beside
it rather than colliding with it.

**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exited 0 at this REQ's claim revision
`2876d96` — **65s wall**, exit status read directly from `$?`.

**Tests baseline:** ✓ `bash _dev/tests/prescribed-shell-canonicalization.sh` exited 0, launched true.
That is the probe this request's new heading has to satisfy.

**Dependencies:** ✓ `depends_on` is empty. Go 1.26.1, ShellCheck 0.11.0, `just` 1.43.0 all at or above
the gate's floors.

**Machine condition:** 4 CPUs, load average around 5 with a sibling builder running. The brief caps
this request at one full canonical gate run for that reason and names the targeted probes to iterate
on instead.

*Checked by work action*
