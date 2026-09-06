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
dispatch_at: 2026-09-06T01:19:09Z
builder_handback_at: 2026-09-06T01:19:09Z
claimed_at: 2026-09-06T00:38:56Z
---

# Move the 46 lines commit.md and inspect.md share into the prescribed-shell guide

## What
`skills/do-work/actions/commit.md` and `skills/do-work-toolbox/actions/inspect.md` share 46 byte-identical non-blank lines in runs of three or more: the M/A/D/X/XD classification legend, four file-reading bullets, and two complete "If the script is missing or will not run, do it by hand" fallbacks that restate the algorithms `internal/corehelpers/inventory.go` and `associate-files` implement. Move the legend, the two fallbacks, and the bullets into one section of `skills/do-work/docs/prescribed-shell-primitives.md` (which already owns the protected-inventory heading), cite it from both actions, and keep only the read-only-versus-staging deltas local.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Both `prime_files` read — `_dev/primes/prime-action-files.md` is what establishes the
  seventeen-line scaffold floor that makes the request's own ceiling unreachable — plus
  `crew-members/maintenance.md`, which this REQ's `maintenance: true` marker requires.
- [x] **[APPLY]:** Five files, exactly the declared `write_set`.
- [x] **[UNIFY]:** `git diff --stat` on the merge range reports five files, 92 insertions, 56
  deletions. Linters: `bash _dev/tests/action-shell-blocks.sh` — `Shell-block lint passed: 74 fenced
  blocks and 33 shipped shell files; ShellCheck enabled.`, exit 0. No debug artifacts. Per file:
  the guide gained one section and the back-reference it needed to keep resolving; `commit.md` and
  `inspect.md` each lost the moved prose and gained a pointer, keeping their own mode-specific
  wordings; the canonicalization probe gained one heading in its membership list; the lock-in gained
  one block beside REQ-552's.

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
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — one new heading-level section holding the
  inventory tag legend, the secret-shaped matching sentence, the four file-reading bullets, the
  association semantics and both manual fallbacks, inserted before the merge-aware commit diff section
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

## Implementation Summary

**Files changed:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `skills/do-work-toolbox/actions/inspect.md` (modified)
- `_dev/tests/prescribed-shell-canonicalization.sh` (modified)
- `_dev/tests/audit-lockins.sh` (modified)

**What was done:** The inventory tag legend, the secret-shaped matching sentence, the four
file-reading bullets, the association semantics and both manual fallbacks now live once, in a new
section of the prescribed-shell guide, with both actions pointing at it. What each action *does* with
an excluded row stayed where it was — a writing action lets only the deletion proceed, a read-only one
inspects only the path and deletion state — because collapsing those would have been a behaviour claim
dressed as deduplication. The guide says in one sentence that this part is caller policy it does not
own.

**The ceiling is 30, measured, and its composition is written into the assertion.** Seventeen lines are
scaffold `_dev/primes/prime-action-files.md:70-80` requires both actions to carry; four are structural
rows inside each action's own flow diagram and error table; eight are the Step 4 semantic clustering
algorithm, deliberately left duplicated because it is caller policy the guide's own charter excludes;
one routes unassociated files to each action's own step number. The request asked for ten, which is
below the scaffold floor and would have been red on day one.

**The assertion was proven red before it was accepted**, which matters because the version this
request inherited was green on day one: its `--glob '*/actions/*.md'` matched zero files, since a
single `*` does not cross a directory separator. The shipped assertion uses `**/actions/*.md` and
reads rg's own exit status rather than a piped total; run against copies of both action files at the
base revision it printed five FAIL lines and exited 1, naming all four manual-fallback sites.

Merge range `b2ba3ea2..a0be855c`, five files, 92 insertions, 56 deletions — identical to the builder
branch's own diff against its base. Builder branch head `d2609378`.

## Decisions — implementation

- **D-03 — the Step 4 semantic clustering algorithm stays duplicated. DECIDE & STATE.** It is caller
  policy, not a shell primitive, and `prescribed-shell-primitives.md:3` states a charter that excludes
  it. This is the deliberate reason the ceiling is 30 rather than 22: moving those eight lines would
  have bought a smaller metric by weakening the guide's own definition of what it holds. Recorded as a
  discovered task instead — it needs a home, and that home is not this guide.
- **D-04 — the wordings that differ between the writing action and the read-only one are preserved.
  DECIDE & STATE.** The guide's legend states classification only.
- **D-05 — a dangling back-reference is fixed rather than worked around. DECIDE & STATE.** Moving the
  prose left an "accepting every alias above" whose "above" no longer resolved; the terminal-success
  alias bullet moved into the same section so it resolves again, tightened to name what it means.
- **D-06 — one sentence in the guide is left untouched for REQ-555. DECIDE & STATE.** REQ-555 rewrites
  the guide's executable-homes table and declares `depends_on: [REQ-554]`; only the half of that
  paragraph this request owns was rewritten.

## Discovered Tasks

- **The Step 4 semantic clustering algorithm is still duplicated between the two actions** — eight
  identical lines, the largest remaining shared block. It needs a home that is not the prescribed-shell
  guide, whose charter excludes caller policy, so it is its own request.
- **The two actions' run-level quarantine paragraphs are near-identical**, differing only in a trailing
  example sentence — and the difflib metric scores them at 0 because that difference breaks the
  matching run, so this duplication is invisible to the new lock-in.
- **`_dev/tests/contract-regressions.sh` is exactly 77 lines against its own ceiling of 77.** Zero
  headroom: the next person who adds a line to it breaks the fast gate.

## Qualification

**Passed, with one reporter-output finding judged rather than obeyed.** Read from the merge range
`b2ba3ea2..a0be855c`; `scope-drift` satisfied, `qualify` satisfied after the judgment below.

- **`QUALIFY-REPORTER-OUTPUT` on `_dev/tests/audit-lockins.sh`: accepted, not a defect.** The gate
  flagged the new block's `print(...)` line. That line is the measurement itself — a `difflib`
  shared-line count the assertion compares against its ceiling — running inside a maintainer test
  script whose whole job is to print a number and a verdict. It is not left-over debug output, it is
  not in shipped code, and the file it lives in is `_dev/`, which is export-ignored. Recorded here
  rather than silenced, because the flag is doing its job and the answer is a judgment.
- **The declared set and the touched set are identical.** Five files, 92 insertions, 56 deletions —
  the same diff the builder branch carries against its base.
- **The assertion was proven red first, and that mattered more than usual here.** The version this
  request inherited was green on day one because its glob matched no files at all. The shipped one
  printed five FAIL lines against copies of both action files at the base revision — the shared-line
  count above the ceiling, plus all four manual-fallback sites by path and line.
- **The ceiling is not a number someone liked; it is composed and written down.** Seventeen scaffold
  lines the action-file prime requires, four structural rows in each action's own diagram and error
  table, eight lines of caller policy the guide's charter refuses to own, one routing line. A reader
  who wants to lower it now has the list of what they would have to change.
- **The one thing the request feared was checked rather than assumed.** The canonicalization scan skips
  the canonical guide itself, so prose moved into the guide is exempt from its stale-pattern list; none
  of its nine patterns or seven old-implementation fragments matches any moved line, and the probe
  passes with the new heading in its membership list.
- **`near_identical_cross_file_pairs` is 0** in `contract-regressions.sh` after the move, which is the
  independent check that the duplication the request targeted is actually gone rather than reworded.

Requirements traced: the shared body lives once and both actions point at it; each action keeps its own
wording; the lock-in fails when the prose returns and passes on the scaffold both actions must carry;
the canonicalization probe accepts the new section; the gate is green.

*Checked by work action*

## Testing

**Tests run:** the canonical gate, plus every probe that could plausibly notice prose moving between
shipped instruction files — `audit-lockins.sh`, `prescribed-shell-canonicalization.sh`,
`action-shell-blocks.sh`, `defensive-surface-audit.sh`, `contract-regressions.sh`, and the heavy-tier
`prescribed-shell-scripts-behavior.sh`.

**Result:** ✓ Green. The canonical gate exited 0 at the merge revision `f45a913` — **79s wall**, exit
status read directly from `$?`. Each probe's own line:
`Audit lock-in regressions passed.`;
`Prescribed shell primitive canonicalization checks passed.`;
`Shell-block lint passed: 74 fenced blocks and 33 shipped shell files; ShellCheck enabled.`;
`Defensive-surface exact deletion regressions passed.`;
`Contract regression checks passed.` — which also reports `near_identical_cross_file_pairs 0`, the
independent check that the duplication is gone rather than reworded;
and `Prescribed shell script behavior probes passed (110 named script cases across 18 per-script files).`

**The assertion's red was taken against the pre-move tree, not asserted.** Extracted standalone and run
against `git show`-restored copies of both action files at the base revision, it printed
`FAIL: commit.md and inspect.md share 46 identical lines; ceiling is 30.` plus four
`FAIL: manual "do it by hand" fallback remains in a shipped action:` lines naming each site by path
and line. Those four sites are what proves the glob is live — the version this request inherited used
`*/actions/*.md`, which matches nothing because a single `*` does not cross a directory separator.

*Verified by work action*
