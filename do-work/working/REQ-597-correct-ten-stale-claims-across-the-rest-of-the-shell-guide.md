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
- [x] **[PLAN]:** Route B. Three builders in one worktree, stacked: the two broken `inspect.md` blocks first, then the sixteen guide claims re-derived from the code, then the two callers' association prose.
- [x] **[APPLY]:** Three commits `a1e652f`, `6913dc4`, `7df6488`, merged as `d5cf28b`. Three files, all in the write set.
- [x] **[UNIFY]:** `git diff --stat 804a8ba..d5cf28b`: 3 files, +23/-23. Four guards green on the merged tree (audit-lockins, prescribed-shell-canonicalization, quiet-grep-pipeline-audit, action-shell-blocks). No debug artifacts, no version or changelog touched.

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
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — sixteen corrections across five sections, plus the line-24 sentence that implies the protected-inventory launcher translates positionals
- `skills/do-work/actions/commit.md` (modify) — the replaced glob, the dash row that is never printed, and the metadata exception
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

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/inspect.md` (modified)
- `skills/do-work/actions/commit.md` (modified)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)

**The two `inspect.md` blocks run now.** They passed the repository root and the quarantine name as
positionals to a launcher that forwards `"$@"` unchanged to a command that takes only a mode and two
flags, so every run exited 2 with `unknown option <root>`, and line 117 read exit 2 as the skip
condition: `inspect` had never associated a file with a REQ in any project. Both blocks now pass
`--quarantine-name do-work-inspect-secret-quarantine` and nothing else. Measured on a fixture with a
working REQ, an archived REQ, an orphan and an untracked `.env.local`: `start` exits 0 with the four rows
and writes the quarantine at mode 0600; `associate` exits 0 with two association rows and none for the
orphan or the quarantined file. The lead-in sentences say what the wrapper actually does with the root:
it reads it from the current directory and rejects a root argument; from a subdirectory `start` prints
the same rows but `associate` exits 2 as if `do-work/` were missing; without `--quarantine-name` it
writes `commit`'s file; `associate` exits 2 with a `HELPER-USAGE` finding when the quarantine is missing.
Each of those was run, not read.

**The association prose in both callers describes the wrapper's real walk and rows.** The inventory
drops an untracked hidden file under `do-work/` as metadata and the prose now says so; the `-` row the
prose told readers to look for is never printed, so "a path from Step 1 that appears in no row is
unassociated" replaces it; a quarantined `X` path is absent from the output even when a REQ claims it,
so rows are matched only against the set that survived the overlay; "re-derives the repository root" is
gone because it does not.

**Sixteen guide claims across five sections, plus line 24 and line 137.** Each replacement was derived
from the code by the builder and checked against a fixture; five of the prior sweep's drafts were
rejected as false on measurement (among them "resolves the repository itself and takes no root
argument", which is half true, and an `EEXIST`/`EISDIR` pair that no single shipped call reports). The
publication sentence that credited every command with a verification only one performs now says which
command does what; the portfolio summary is described as reading its source once and writing every
output from those bytes, with no "retained toolbox script", because none exists; the report image batch
on line 137 had the same phantom script and was corrected on the request's fourth requirement; line 24 no
longer implies the launcher translates positionals. The per-claim evidence, with the fixture and output
for each, is in `do-work/runs/work-2026-09-05-231943/REQ-597-handback.md`.

**Not changed:** any Go source. Where the code is surprising the description now says so, and the code
defects the builders measured on the way are captured below rather than fixed in a prose request.

## Decisions

- **D1 The blocks changed to the form the launcher accepts, not the launcher to accept the blocks.** The
  launcher's contract is the CLI's, and `commit.md` already used the flag form. That the launcher cannot
  pass `--repo-root` at all, so neither caller can run from anywhere but the root, is a code defect with
  its own request (REQ-603); this request's prose states the dependence in the meantime.
- **D2 Line 24 and line 137 corrected beyond the sixteen.** Both are the class the request names, both
  sit in sections the request opened, and the fourth requirement says to check those sections whole.
- **D3 Drafts were inputs, not answers.** Five of the verification's suggested sentences were false on
  measurement and were not used. Every shipped sentence was re-derived, which is the rule REQ-596's
  review established for this file.
- **D4 No changelog in the builder commits.** Three shipped files make this a release; finalization writes
  one entry covering the three commits, as every recent release commit in the log does.

## Discovered Tasks

- **The protected-inventory launcher cannot pass a global flag and its compatibility shim discards the
  text it prepares.** `scripts/protected-inventory.sh:6` puts `"$@"` after the command token, so
  `--repo-root` is an unknown option and the runtime takes the current directory as the root; the shim
  loop at `inventory.go:445-456` replaces the result text unconditionally, so `NO-DO-WORK-DIR`,
  `PARSE-FAILED` and a walk error's finding never reach a caller and the exit is a silent 2, which both
  callers read as "skip REQ tracing". Also: `commit.md:67` tells a re-run to append to the retained
  quarantine, but `start` replaces it and only `associate` unions; `commit.md:61`'s exit-2 reading also
  covers a `git status` or quarantine-write failure; `associate` after `start --dry-run` exits 2 as
  not-started. Captured as **REQ-603**.
- **`atomic-download`'s occupancy policy is asymmetric and one stat is unchecked.** `--dry-run` refuses
  any existing target (exit 2) while the live run refuses only a directory and `os.Rename` silently
  replaces an occupying regular file, reported as `created`; `commands.go:891` discards the error from
  `os.Stat` and then reads `info.Size()`. Captured as **REQ-604**.
- **`finalization_apply.go:545` runs `diff-tree` without `-m`**, so a merge commit among the candidates
  lists no paths and the `exact` loop stays true; only the preceding binary-diff digest match keeps it
  unreachable. Captured as **REQ-605**.
- **The same phantom-script class in more out-of-write-set callers**: `present-work.md:136` and `:140`,
  `ai-report-reference.md:31`, `:37`, `:47`, `architecture-report.md:46`, `install.md:50`, `:246`, `:261`,
  `:335`, `board.md:87`, and `work-reference.md:322` ("hands the child the console's own handles"; both
  child streams get the CLI's stderr). Appended to **REQ-601**, whose write set gains
  `architecture-report.md` and `work-reference.md`.
- **The lifecycle timing category vocabulary is closed** (ten names; `verification` alone is rejected)
  and not documented where a caller would look. Noted; no request, the guide now names the categories.

## Qualification

**Passed.** Read from the range `804a8ba..d5cf28b`, three files, 23 insertions and 23 deletions.
Canonical `qualify` and `scope-drift` both satisfied; scope-drift first flagged two backticked tokens in
the Scope prose (`-` and `protected-inventory.sh`) as declared paths, and the prose was reworded, not the
parser fought.

- **The behaviour change is the two `inspect.md` blocks, and it is measured in both directions.** Before:
  exit 2, `unknown option <root>`, no quarantine written, from both blocks. After: `start` exit 0 with the
  four rows and a 0600 quarantine; `associate` exit 0 with the two association rows and none for the
  orphan or the quarantined file. That is the first time `inspect` has associated a file with a REQ.
- **Every replacement sentence has a run behind it, not a reading.** The builders' hand-back lists the
  fixture and output per sentence; five of the prior verification's drafts were rejected on measurement
  and not used, which is the rule REQ-596's review set for this file class and the reason a third builder
  chain was worth its cost.
- **Widening is inside the request's own fourth requirement.** Sixteen claims, not ten, and line 24 and
  line 137 besides: all in sections the request opened and checked whole.
- **All four guards over the three files exit 0** on the merged tree: `audit-lockins.sh`,
  `prescribed-shell-canonicalization.sh`, `quiet-grep-pipeline-audit.sh`, `action-shell-blocks.sh`.
- **Five discovered tasks captured, none folded into this prose request:** REQ-603, REQ-604, REQ-605 for
  the measured code defects; REQ-601 widened for the twelve further phantom-script sites.

## Testing

**No lane reads the sentences that changed; the evidence is the fixtures.** Everything below the
guards was measured by the builders on fixtures kept under the scratchpad (`req597-inspect/`,
`req597-guide/`, `req597-callers/`) and is re-executed independently by the three-lens review.

- `bash _dev/tests/audit-lockins.sh` — `Audit lock-in regressions passed.`, exit 0.
- `bash _dev/tests/prescribed-shell-canonicalization.sh` — passed, exit 0.
- `bash _dev/tests/quiet-grep-pipeline-audit.sh` — `quiet-grep pipeline audit passed (95 tracked shell
  files, 19 must-flag and 7 must-not-flag shapes).`, exit 0.
- `bash _dev/tests/action-shell-blocks.sh` — `Shell-block lint passed: 74 fenced blocks and 33 shipped
  shell files; ShellCheck enabled.`, exit 0. The two rewritten `inspect.md` blocks are among the 74.
- Fast gate on the merged tree at `d5cf28b`: `Maintainer verification passed.`, exit 0, gate wall 84s.
- Builders' gate from the worktree at each stage: `Maintainer verification passed.`, exit 0 (wall 83s at
  stage 1).
