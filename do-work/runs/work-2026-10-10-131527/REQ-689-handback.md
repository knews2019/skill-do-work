# REQ-689 hand-back: do-work run --coordinate

- Branch: `worktree-agent-REQ-689-run-coordinate-mode`
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-689-run-coordinate-mode
- Base commit: bd56c4b0
- Commits: 71c1ee5b `[REQ-689] add do-work run --coordinate: mandatory coordinator shape, run rules, preflight, stall loop, coordinated handoff` (one commit)

## File manifest

- `skills/do-work/SKILL.md` (modified): one route row between the run-with-recovery row and the plain run row: `drive the queue`, `use the main session as a coordinator`, `run --coordinate` → `./actions/work.md` with `--coordinate`. `argument-hint` untouched.
- `skills/do-work/actions/work.md` (modified): Architecture `:37` one sentence pointing at `--coordinate` (who integrates); `## Input` gains the `--coordinate` entry; stripped-token list, usage line and Step 0 checklist gain `--coordinate`; Step 10 teardown and its checklist line gain idle background agents, the stall check and the merged-worktree list. "degrades silently to the serial loop" and `[--fan-out [N]]` kept.
- `skills/do-work/actions/fan-out-reference.md` (modified): inside Delegated integration, the mandatory-shape paragraph and a `#### Coordinated run rules` subsection (five rules, REQ-069/REQ-073 boundary, preflight, run-policy, stall loop); run-directory table gains the `REQ-NNN-progress.log` and `full-gate.lock` rows; the `manifest.md` row gains the three columns. File header and other sections unchanged.
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modified): When to Use gains the automatic context-threshold trigger under `--coordinate`; resume list gains `do-work run --coordinate --fan-out N` beside the kept `do-work run --fan-out N`.
- `skills/do-work/docs/work-guide.md` (modified): one sentence in "Building several REQs at once" describing `--coordinate` in plain words.
- `skills/do-work/docs/standing-preferences.md` (modified): one table row for the pasted coordinator directive (PD-5).
- `skills/do-work/actions/work-reference.md`: unchanged (D-07).
- `skills/do-work/crew-members/background-agents.md`: unchanged (D-08). Ceiling note kept.

## AI Execution State (P-A-U Loop)

- **[PLAN]:** Read the brief, REQ (all sections incl. Exploration PD-1 to PD-6 and Scope), UR-153, archived REQ-662/666/667 (as data), crew members general, coding-guardrails, shared-principles, communication-style, primes prime-action-files and prime-releases, lessons-releases. Verified before restating: queue-mode `advance` refuses unknown tokens (per the Exploration's cited `next_targets.go`), `verify.go` `diskSpaceLevelFor` returns critical/warning/neutral and emits a `low-disk-space` finding only for the first two, `actions/forensics.md` heading is `### 14. Release and Queue Invariants (board-owned)`, REQ-690 defines `--watch` as an output shape for a scheduler body. Plan: insert-only edits at the exact anchors in the brief; one subsection under Delegated integration holding the rules, preflight, run policy and stall loop; teardown in work.md Step 10 only (fan-out-reference's stall loop points there).
- **[APPLY]:** Edits as planned, all insertions beside existing lines, no reflow. Six files changed, two named files left unchanged after the sweep.
- **[UNIFY]:** `git diff bd56c4b0 --stat`:
  ```
   skills/do-work/SKILL.md                            |  1 +
   skills/do-work/actions/fan-out-reference.md        | 29 +++++++++++++++++++++-
   .../actions/restart-with-parallel-handoff.md       |  2 ++
   skills/do-work/actions/work.md                     | 13 +++++-----
   skills/do-work/docs/standing-preferences.md        |  1 +
   skills/do-work/docs/work-guide.md                  |  2 +-
   6 files changed, 40 insertions(+), 8 deletions(-)
  ```
  Checks (from the worktree root):
  - RED probe at base: exit 1, 0.14 s wall.
  - GREEN probe (runs shipped-package-reference-contract.sh and action-shell-blocks.sh): exit 0, 3.96 s wall before commit; exit 0, 8.90 s wall after commit.
  - `grep -rniw pid skills/do-work/actions/fan-out-reference.md`: no output (exit 1 = no match).
  - `git diff --check`: exit 0, clean.
  - `_dev/tests/contract-regressions.sh`: exit 0, 34 s wall.
  - `_dev/tests/staged-skills-contract.sh`: refuses to run outside `maintainer-verify.sh --heavy` ("heavy-only"); not run, the gate is the integrator's.
  - Citation check is live: temporarily renaming my forensics citation to a nonexistent section made `shipped-package-reference-contract.sh` FAIL on that line; restored before commit.
  - Files checked by reading the diff: all six changed files; `work-reference.md:396`, `background-agents.md` (whole-file grep for argument list, run-directory table, builder test rule), `run-simple-reqs.md` (PD-6, not edited).
  - No debug artifacts; nothing under `do-work/` staged or committed.

## Proof record

**RED at base bd56c4b0** (exit 1):
```
FAIL: work.md ## Input has no --coordinate entry
FAIL: work.md usage line lacks --coordinate
FAIL: work.md Step 0 checklist line lacks --coordinate
FAIL: SKILL.md: coordinate route row or plain run row not found
FAIL: SKILL.md coordinate row does not route to ./actions/work.md
FAIL: SKILL.md coordinate row lacks the coordinator phrase
FAIL: fan-out-reference.md lacks the mandatory-shape sentence
FAIL: fan-out-reference.md lacks the pending-answers rule
FAIL: run-directory table lacks the REQ-NNN-progress.log row
FAIL: run-directory table lacks the full-gate.lock row
FAIL: manifest.md row lacks the column: takeover-to-finalization minutes
FAIL: manifest.md row lacks the column: full gates run
FAIL: manifest.md row lacks the column: stall restarts
FAIL: run rules do not name the REQ-069/REQ-073 boundary
FAIL: preflight lacks the low-disk-space level
FAIL: fan-out-reference.md lacks do-work/run-policy.md
FAIL: stall loop does not read do-work status --watch
FAIL: handoff resume command lacks do-work run --coordinate --fan-out N
FAIL: teardown does not stop idle background agents
19 check(s) failed
```

**Read-through cases:**
1. `do-work run --coordinate --fan-out 3`. Before: `--coordinate` is residue after stripping, so the run stopped with "Unrecognized argument(s): --coordinate". After: `## Input` defines `--coordinate`; the stripped list now strips `--wave N`, `--fan-out [N]`, `--coordinate`, `--skip-impact-negligible`; residue is empty, so no error. The action consumes `--coordinate` (session keeps the mode) and passes `--fan-out 3` to queue-mode `advance`. Bare `do-work run --coordinate` passes a bare `--fan-out`.
2. Message "drive the queue". Before: no row matched it; as unmatched descriptive multi-word input it fell to capture. After: first-match routing reaches `SKILL.md:33` (the new row) before `:34` (plain run row) and routes to `./actions/work.md` with `--coordinate`. `run --coordinate` also hits `:33` first, since `:32` (run-with-recovery) contains none of its words as a trigger.
3. `grep -n "progress.log\|full-gate.lock\|run-policy" skills/do-work/actions/fan-out-reference.md`. Before: no output. After: lines 140 (Progress log rule), 141 (Full-gate lock rule), 146 (boundary sentence), 152 (preflight stale locks), 153 (preflight run policy), 155 (run-policy paragraph), 157 (stall loop), 169 (`REQ-NNN-progress.log` table row), 170 (`full-gate.lock` table row). The `manifest.md` row (line 171) carries the three columns.
4. `phandoff` inside a coordinated run. Before: the handoff's build resume line was `do-work run --fan-out N`. After: the resume list offers `- to build, when this run was started with --coordinate — do-work run --coordinate --fan-out N`, so the paste block carries `do-work run --coordinate --fan-out N`.
5. Plain `do-work run --fan-out 3`: the `--fan-out [N]` Input entry, the Architecture paragraph and "degrades silently to the serial loop" are byte-unchanged; the only change on its parse path is that `--coordinate` joins the stripped list, which does not match `--fan-out 3`. It parses and degrades exactly as before.

**GREEN** (exit 0):
```
shipped package reference contract: PASS
Shell-block lint passed: 75 fenced blocks and 33 shipped shell files; ShellCheck enabled.
REQ-689 GREEN probe: all checks passed
```

**Stall-loop path documented:** both. The tick reads `do-work status --watch`; where that command is absent it reads each `REQ-NNN-progress.log` tail and each worktree's last commit time (PD-3).

## Decisions

- D-01 DECIDE & STATE: route phrasing is exactly `drive the queue`, `use the main session as a coordinator`, `run --coordinate`; no bare `coordinator` (PD-1). Route cell `./actions/work.md` with `--coordinate`.
- D-02 DECIDE & STATE: subsection title `#### Coordinated run rules`. Preflight, run policy and stall loop sit inside it as bold-led paragraphs, each stating its own condition (`--coordinate` only), while the five rules apply whenever the coordinator shape is used. Keeps the new heading count at one.
- D-03 DECIDE & STATE: the heartbeat rule is titled **Progress log**, not "Heartbeat", so it does not read as the `heartbeat_at` machinery REQ-069/REQ-073 deleted, which the next sentence of the same subsection rules out.
- D-04 DECIDE & STATE: degradation wording under `--coordinate`: it prints one line, `--coordinate: no worktree or agent dispatch here, running the serial loop`, and runs the serial loop (PD-4). Plain `--fan-out` still degrades silently.
- D-05 DECIDE & STATE: the Input entry adds one clause that `advance`'s continuation does not carry the flag, so the session keeps the mode for the whole run (fact from the Exploration, `lifecycleadvance/queue_commands.go:306-326`). Without it a coordinator following the continuation argv could drop the mode.
- D-06 DECIDE & STATE: "Each tick re-runs the selector" gained "in a writing gap", because a queue-mode selection commits a claim and the existing one-writer rule forbids coordinator writes while an integrator runs.
- D-07 DECIDE & STATE: `work-reference.md:396` left unchanged. It restates the integrator rules (entry, never recover, never Step 10), none of which changed; it says nothing about the shape being optional, so it is not stale.
- D-08 DECIDE & STATE: `background-agents.md` left unchanged. It restates neither the run argument list nor the run-directory table nor a builder test rule; its manifest example is a generic code-review manifest. Ceiling note kept.
- D-09 DECIDE & STATE: the preflight disk line cites `actions/forensics.md` → **14. Release and Queue Invariants (board-owned)** (the heading's full text; the brief's shorter form would fail the section-citation check). It quotes the level as warning, critical or none: `verify.go` emits a `low-disk-space` finding only at warning or critical, and "neutral" otherwise.
- D-10 DECIDE & STATE: the manifest-row rule names the coordinator as the writer of the three columns, in the gap after each integration, consistent with the existing "the orchestrator's, never written by a builder" and the one-writer rule.
- D-11 DECIDE & STATE: the full-gate lock is removed by the integrator after its suite; a lock left behind is what the preflight reports as stale and never deletes.

## Discovered Tasks

- `skills/do-work/actions/run-simple-reqs.md:24,64,67` forwards only `--fan-out` to `do-work run`, so a simple-REQs run cannot be coordinated (PD-6). impact-user-visible → report only
- `skills/do-work/SKILL.md:4` `argument-hint` lists no run flags at all (`run [REQ|UR]`), so neither `--fan-out` nor `--coordinate` is discoverable there; left alone because REQ-690/691/692 edit that line. impact-negligible → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file): family `canonical-link-outlives-its-target` (two bullets), `manifest-ownership-vs-edit-content`. No archive links added in shipped prose, so nothing to sweep.
- Quoted in the brief: `lessons-action-files.md` [family: restated-mechanism-unchecked] (REQ-639: Verified Facts are claims; read the command before shipped prose restates it). Applied: read `verify.go` for the level names and the forensics heading before citing them.
- Quoted in the brief: `lessons-do-kanban.md` [family: disk-space-blind-spot] (the probe reports and never deletes; repo root only). Applied: the preflight line says the preflight never deletes either.
- `prime-action-files.md` Traps: `alternate-writer-contract-drift` (swept restatements in the four named files plus `run-simple-reqs.md` and help surfaces by grepping `--fan-out`), `budgeted-context-routing` (not applicable).

## Anti-bloat check

`git diff bd56c4b0 --stat` is in [UNIFY] above. Additions the REQ did not name:
- `docs/standing-preferences.md` row: PD-5, the pasted directive is the REQ's Why.
- Input clause "`advance`'s continuation does not carry the flag": D-05, needed so the mode survives continuations.
- "in a writing gap" on the selector re-run: D-06, keeps the one-writer rule true.
- "removing it after" on the full-gate lock: D-11, one clause so the stale-lock preflight line has a meaning.
No new action file, Go, script, helper, frontmatter field or status. New surface: one flag, one optional policy file, two run-directory files, three manifest columns.

## Proposed CHANGELOG entry

**Coordinated Runs With `do-work run --coordinate`**

The main session can now be told, with one flag or the phrase "drive the queue", to act only as a coordinator, instead of the directive users pasted at the start of about 13 sessions. Coordinated runs carry their run rules, a preflight and a stall check in the skill itself, and a handoff from one resumes coordinated.

- `do-work run --coordinate [--fan-out N] [REQ-NNN ...]`: the main session dispatches, writes briefs and reports, and never builds, merges or runs Steps 6 to 9; no integrator pushes; a blocking question becomes a `pending-answers` follow-up. Bare `--coordinate` implies `--fan-out`. Plain `--fan-out` is unchanged.
- Run rules whenever the coordinator shape is used: a per-REQ `REQ-NNN-progress.log`, one full suite at a time under `full-gate.lock`, rerun-alone triage for a red full gate, three new `manifest.md` columns, and a hand-back coverage check.
- A four-line preflight (board disk level, leftover processes, stale locks, optional `do-work/run-policy.md`) and a stall loop every 15 to 20 minutes that reads `do-work status --watch`, or progress logs and worktree commit times.
- Handoffs from a coordinated run resume with `do-work run --coordinate --fan-out N`, and the coordinator writes one on its own when the harness reports high context use. Teardown stops idle background agents and lists merged worktrees for removal.

## Proposed lesson bullet

Satellite: `_dev/primes/lessons-action-files.md`.

- [family: restated-mechanism-unchecked] [REQ-689: a brief's section citation is a claim too; `→ **Release and Queue Invariants (board-owned)**` fails the citation check because the heading carries its number, `14. `, so cite the heading's full text and let `shipped-package-reference-contract.sh` prove it](../../do-work/archive/UR-153/REQ-689-run-coordinate-mode.md)

## Integration seams

- `skills/do-work/SKILL.md`: one row at `:33`, between run-with-recovery and run. REQ-690/691/692 add rows elsewhere in the table; a conflict, if any, is adjacent-row only.
- `skills/do-work/actions/work.md`: edits at `:37`, `## Input` (`:104-111`), Step 10 (`:471`) and checklist lines Step 0 (`:482`) and Step 10 (`:502`). REQ-658/659 edit Steps 6 to 9 and their checklist lines; the Step 10 checklist line is the line most likely to sit next to theirs.
- `skills/do-work/actions/fan-out-reference.md`: only Delegated integration and the run-directory table. REQ-660 edits Naming, merge and Cleanup.
- `skills/do-work/docs/work-guide.md`: one sentence inside `:134`. If REQ-690 rewrites that paragraph, re-apply the sentence.
- REQ-690: the stall loop names `do-work status --watch` and keeps the fallback, so it reads correctly whether REQ-690 lands or is set aside.

Test wall times: RED probe 0.14 s; GREEN probe 3.96 s and 8.90 s (two runs, both exit 0); contract-regressions.sh 34 s.
