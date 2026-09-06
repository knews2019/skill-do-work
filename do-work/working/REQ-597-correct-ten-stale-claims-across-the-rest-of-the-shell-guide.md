---
id: REQ-597
status: claimed
domain: general
created_at: 2026-09-06T05:49:08Z
user_request: UR-105
review_generated: true
impact: impact-negligible
effort_estimate: effort-mechanical
route: B
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-06T06:59:15Z
  basis:
    - Route B
    - 3-file write set
    - 4 acceptance criteria
maintenance: true
depends_on: [REQ-596]
related: [REQ-555, REQ-595, REQ-596]
write_set: [skills/do-work/docs/prescribed-shell-primitives.md, skills/do-work/actions/commit.md, skills/do-work-toolbox/actions/inspect.md]
title: 'Correct ten stale claims across the rest of the prescribed-shell guide'
claimed_at: 2026-09-06T06:59:15Z
---

# Correct Ten Stale Claims Across the Rest of the Prescribed-Shell Guide

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

A 69-claim sweep of `skills/do-work/docs/prescribed-shell-primitives.md` against the Go code found 21
claims that do not hold. REQ-596 corrected the eleven inside the two sections it owned. Ten remain, in
sections REQ-596 never opened, at guide lines 30, 32, 80, 92 (two), 106 (two) and 128 (two).

## Why

Same finding class as REQ-555, REQ-595 and REQ-596: the guide is the pointer target from sixteen shipped
files and describes implementations that do not exist. Two of the ten are marked `stale` rather than
`partly-stale`, meaning the sentence is false rather than incomplete.

## Context

Found by the section sweep run for REQ-596. The full report with each claim, the code that contradicts
it, a file:line citation and suggested replacement text is in
`do-work/runs/work-2026-09-05-231943/REQ-596-section-sweep.json`, third element of the array. That file
is the input to this request — do not re-derive the census, but do re-verify each claim against the code
before rewriting it, because a replacement sentence that is also wrong is worse than the one it replaces.

## Detailed Requirements

- Correct the ten claims the sweep names, in the lifecycle-timing, merge-aware-commit-diff,
  commit-file-listing, verified-exact-publication and portfolio-summary sections.
- Two are `stale` rather than `partly-stale` and should be checked first: the claim that the
  diff-tree form is the suite default for path-only consumers, and the claim about how `ln` and `mv`
  treat a directory standing in the destination's place.
- Check every other claim in each section you open, the way REQ-595 checked all fourteen Mechanics cells
  and REQ-596 checked its two sections. The sweep reports 37 claims checked across these sections; a
  claim it marked accurate is evidence, but a claim it did not reach is not.
- No sentence may generalize over two commands that behave differently. That shape has been caught three
  times in this file within the visible history — REQ-555, REQ-595 and REQ-596 — and the sentence
  REQ-596 replaced predates that history.
- **Two shipped callers now contradict the guide and are in this request's write set.**
  `skills/do-work/actions/commit.md:85` still names the `do-work/archive/**/REQ-*.md` glob REQ-596
  replaced, and still says the delegated check "prints one `<owner>\t<path>` row per candidate — a
  `REQ-NNN` id, or `-` for unassociated", with line 87 adding that files coming back `-` remain
  unassociated. `protected-inventory associate` never prints a `-` row; only the separate
  `associate-files.sh` entry point does. Check `../do-work-toolbox/actions/inspect.md` for the same two
  claims. The guide is the canonical home those actions point at, so a caller disagreeing with it is the
  same defect one file over.

## Constraints

- Guide prose only: no behaviour change, no code change, no route-column change, no Mechanics-column
  change, no heading change.
- `_dev/tests/audit-lockins.sh` Finding 7 pins the route column, the orchestration claim and the
  Mechanics column's shell vocabulary. Do not weaken any of the three.
- If a claim cannot be verified against code, say so in the request rather than guessing.

## Dependencies

Depends on REQ-596, which corrected the other eleven findings from the same sweep in the same file.

## Open Questions

None.

## Triage

**Route: B** — Explore then build.

**Reasoning:** The ten claims are named with file-and-line citations from a sweep that has already run,
so their existence is settled. What is not settled is their replacements. REQ-596 shipped a replacement
sentence that was false in its own right — it compressed a correct finding into an incorrect one — and
the review caught it. Every replacement here has to be re-derived from the code rather than transcribed
from the sweep's suggestion, and the sections these ten sit in have to be checked whole, because that is
the requirement that turned REQ-595 from one cell into three and REQ-596 from three claims into eleven.

**Planning:** Skipped. Three files, and the edit set is whatever the re-verification confirms.

**The caller drift is the newest part and the least examined.** `actions/commit.md` names a glob the
guide replaced one request ago and describes a `-` row the command never prints. Nobody has checked
`inspect.md` for the same two claims, or the other shipped callers for anything similar.

## Plan

**Planning not required** — Route B: three files, and the work is whatever the re-verification and the
section sweep confirm.

*Skipped by work action*

## Exploration

Three read-only agents: one re-derived the ten named claims from the code, one checked the five sections
they sit in whole, one checked every shipped caller that points at the guide. Each built fixtures and
ran the real commands. Full report in the run directory as `REQ-597-verification.json`.

**All ten named claims are confirmed, and the whole-section pass found six more in the same five
sections** — sixteen in the guide, not ten. Two of the ten are false outright: "the diff-tree form above
is the suite default for path-only consumers" (no shipped file contains `diff-tree` at all; the three Go
readers carry three different flag sets), and "the retained toolbox script is not a fallback" (no such
script exists in either skill). The rest generalize across a difference or credit one command with
another's behaviour: "the wrapper attaches the child's stdout to the console" (it is merged onto the
CLI's stderr, so redirecting stdout captures the rendered result and none of the child's bytes);
"timing evidence is metadata only" (the `--operation` label is free text, printed verbatim into the
folded `## Timing` section that is committed with the REQ — a token or a path in it lands in durable
evidence); "the publishing step verifies the path it actually wrote" (one of four shipped publications
does; the others take the proof from the call, and `rename` refuses only a directory while silently
replacing a file).

**Four of the prior sweep's own suggested replacements were false**, which is the failure the request
was written to prevent and the reason every replacement here was re-derived rather than transcribed.
The worst: "Go's `os.Rename` and `os.Link` refuse a destination that already exists". `rename(2)`
replaces a regular file silently and refuses only a directory — measured — and the sweep had tested
only the directory case and generalized. Another said `--agent` and `--revision` reach the folded
section; they reach the stream only. Another stated a shipped convention as a rule the command enforces.
The fourth dropped the one consumer whose flags contradict the point being made.

**The callers are worse than the guide, and one of them is broken code, not prose.**
`skills/do-work-toolbox/actions/inspect.md` prescribes, in two ```bash blocks agents run,
`protected-inventory.sh start "$(git rev-parse --show-toplevel)" do-work-inspect-secret-quarantine`.
That launcher is a bare `exec` that translates nothing; the Go handler rejects any positional after the
mode and exits 2 with `unknown option`; and the action's Step 3 reads exit 2 as its skip condition. **As
shipped, `inspect` never associates a file with a REQ, in any project.** The working form is
`--quarantine-name do-work-inspect-secret-quarantine` with no root argument, confirmed in a fixture.

The sentence that plausibly let that survive review is in the guide, and it is one this run wrote:
line 24's "`scripts/protected-inventory.sh` goes further and sets `DO_WORK_COMPATIBILITY_SHIM=1`"
implies it also translates positionals, which it does not. Both `commit.md` and `inspect.md` also still
name the `**` glob the guide replaced one request ago and describe a `-` row the wrapper never prints —
the signal for "unassociated" is a missing row, which also covers a quarantined `X` path.

**Four more shipped files outside this request's write set carry the same class**, each verified:
`ai-report-reference.md` (staging "adjacent to `generated/`" — it is the system temp directory, twice),
`install.md` (a Red Flag names `SKILL.md.download`, a file `os.CreateTemp` never creates; the real
stray is `SKILL.md.download.<random>`), `present-work.md` (outputs "verified against the source" — no
output is ever compared to the source), `board.md` ("skips silently" — it prints a warning finding).
Captured as REQ-601.

Fourteen other caller sites were checked and hold.

*Generated by three Explore agents*

## Scope

**Files I will touch:**
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — sixteen corrections across five sections, plus the line-24 sentence that implies `protected-inventory.sh` translates positionals
- `skills/do-work/actions/commit.md` (modify) — the replaced glob, the `-` row that is never printed, and the metadata exception
- `skills/do-work-toolbox/actions/inspect.md` (modify) — the same three, plus the two prescribed blocks whose arguments the command rejects

**The two `inspect.md` blocks are the priority and the only behaviour change.** Every other edit here
corrects a description; those two correct an instruction that exits 2 on every run. That is in scope
because the request names the file and the class, and because leaving a known-broken prescribed command
in place to protect a prose-only framing would be the wrong call.

**Sixteen guide corrections, not ten**, on the request's own fourth requirement — the sections were
checked whole. Each replacement is derived from the code by the builder, never transcribed from the
prior sweep, four of whose drafts were false.

**Files I will NOT touch:** `ai-report-reference.md`, `install.md`, `present-work.md`, `board.md` — same
class, outside the write set, captured as REQ-601. Any Go source: the guide's job is to describe the
code, and where the code is surprising (rename replaces a file; `--operation` is unredacted) the
description says so rather than the code changing. `_dev/tests/`: no guard here is derivable from code
the way the pattern set was.

**Acceptance criteria:**
- [ ] The two `inspect.md` prescribed commands run to exit 0 against a fixture repository and produce
  association rows
- [ ] Every replacement sentence is derived from the code by the builder and checked against a fixture;
  none is transcribed from the prior sweep's drafts
- [ ] The guide's two sections that describe the same word two ways — shell `mv` nests, `rename(2)`
  refuses — are consistent with each other after the edit
- [ ] No sentence generalizes over commands that behave differently, and no sentence credits a command
  with another's work
- [ ] The three guards over the guide exit 0
- [ ] The four out-of-scope callers are captured, not fixed

## Pre-Flight

**Green gate at `34077d8`**, the revision the builder branches from.
`bash _dev/tests/maintainer-verify.sh` printed `Maintainer verification passed.` and exited 0, gate
wall 75s. One `SKIP` line, the heavy-only one every fast run prints.

**Three guards over the guide are green and must stay green**: `audit-lockins.sh` (the route column,
the orchestration claim, the Mechanics column's shell vocabulary, and the secret pattern set derived
from `secretPath`), `prescribed-shell-canonicalization.sh` (twelve headings, sixteen pointer sites),
and `quiet-grep-pipeline-audit.sh`. None of those columns, headings or pointers changes here.

**The one behaviour change has a before and an after that can be run.** The two `inspect.md` blocks
exit 2 with `unknown option` against any repository today; the corrected form is confirmed in a fixture
to exit 0 and print association rows. That is the builder's first proof, before any prose moves.

**The prior sweep's drafts are input, not answers.** Four of its ten guide suggestions were false when
re-derived. The builder is given those four by name, and every replacement it ships must cite the code
it was derived from and the fixture it was checked against.

**The builder works in an isolated worktree** at
`../skill-do-work-worktrees/worktree-agent-REQ-597-guide-and-callers`, branched from `34077d8`, and
hands back one file to the main checkout without staging or committing anything there.
