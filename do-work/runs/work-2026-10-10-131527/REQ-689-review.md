# REQ-689 review: do-work run --coordinate (orchestrated, Route B, wave end)

Reviewer for run work-2026-10-10-131527. Diff: `cc049a01..f8f32682` (7 files, +41 -9; the builder commit plus two integrator seams in the merge: the stall loop reads `do-work status --watch` and the progress-log tails, and `status-guide.md:34` drops progress logs from its example). Line numbers below are at HEAD `f8f32682`.

## Verdict

86%. Approve after four one-sentence prose fixes (F1, F2, F4, F5). The flag, route, rules, preflight, stall loop, resume line and teardown are all there and match the REQ almost word for word. The four defects are seams between the new rules and rules that already existed: the builder write boundary, the one-writer rule, the handoff's claim lines, and a long full gate inside a 20-minute silence window.

## Evidence

- GREEN probe at HEAD: `bash do-work/runs/work-2026-10-10-131527/REQ-689-probe.sh` exit 0, 2.6 s (`shipped package reference contract: PASS`, `Shell-block lint passed: 81 fenced blocks`, `REQ-689 GREEN probe: all checks passed`).
- `grep -niw pid skills/do-work/actions/fan-out-reference.md`: no output (exit 1). No liveness check.
- `run_status.go:363-373`: `listRunLocalFiles` skips `manifest.md` and every `REQ-NNN-*` name, so the integrator seam "status does not read progress logs" is true, and `full-gate.lock` is listed as run-local.
- `next_commands.go:39-83` (read by the Exploration): queue-mode `advance` refuses unknown tokens, so consuming `--coordinate` in the action is required. Confirmed by the `## Input` text at `work.md:105`.
- Not run: the full test gate (not the reviewer's).

## Requirements check (REQ items 1 to 18)

| # | Item | Result |
|---|---|---|
| 1 | flag in Input, stripped list, usage, Step 0 line; bare form implies `--fan-out` | Delivered (`work.md:105`, `:108`, `:111`, `:482`) |
| 2 | route row above plain run row, phrases only | Delivered (`SKILL.md:33` above `:34`; no collision with `:31`, `:32`) |
| 3 | mandatory-shape paragraph | Delivered word for word (`fan-out-reference.md:134`); wording seam F7 |
| 4 | plain `--fan-out 3` unchanged | Delivered (Input entry and silent degrade byte-unchanged) |
| 5-10 | progress log, full-gate lock, red gate, manifest columns, hand-back coverage, run-dir rows, REQ-069/073 boundary | Delivered (`:146-152`, `:175-177`); seam F1 |
| 11-12 | four-line preflight, optional `run-policy.md` | Delivered (`:154-161`) |
| 13-14 | stall loop | Delivered (`:163`); defect F2, wording F3 |
| 15 | coordinated resume line, automatic context trigger | Delivered (`restart-with-parallel-handoff.md:14`, `:65`); seams F4, F5 |
| 16 | teardown stops idle agents, lists merged worktrees | Delivered (`work.md:471`, `:502`) |
| 17 | restatement sweep | Partial: F1, F4, F5, F6 are restatements the sweep missed |
| 18 | release | At finalization (N/A here) |

Requirements 90%. Every touched file is in Scope (status-guide.md added as the integrator seam); `work-reference.md` and `background-agents.md` were declared "only if stale" and are correctly unchanged (D-07, D-08 confirmed). Decisions D-01 to D-11 cover every choice visible in the diff.

## Findings

### Important

- **F1. The progress log breaks the builder write boundary.** `fan-out-reference.md:146` makes every builder append to `REQ-NNN-progress.log` in the run directory. But `:61` (Sole integrator) says the builder writes the main tree "with exactly one exception", its hand-back, and that "anything else under `do-work/runs/`" is a violation. The run-directory table row at `:173` restates "The one main-tree path a builder may write". The rule also does not say the path must be the absolute main-tree path from the brief, so a builder that writes it repo-relative writes into its worktree and the coordinator reads nothing (`:182` warns about this for the hand-back). The hand-back merge step 0 (`:73-74`) does not classify progress logs or the lock, so every merge either names them as foreign files or, if an integrator reads its own log as "a run artifact this REQ owns", stops on it. impact-user-visible. **Fix:**
  - `:61` replace `**with exactly one exception: its own \`do-work/runs/work-<YYYY-MM-DD-HHMMSS>/REQ-NNN-handback.md\`**, reached by the absolute main-tree path the orchestrator hands it and **never staged, committed, or merged** — it is an orchestrator-owned working file, not branch content.` with `**with exactly two exceptions: its own \`do-work/runs/work-<YYYY-MM-DD-HHMMSS>/REQ-NNN-handback.md\` and, when the coordinator shape is used, its own \`REQ-NNN-progress.log\` beside it (*Coordinated run rules*)**, each reached by the absolute main-tree path the orchestrator hands it and **never staged, committed, or merged** — each is an orchestrator-owned working file, not branch content.` and replace `The exception exists because the hand-back has to survive a dead transcript` with `The exceptions exist because both files have to survive a dead transcript`.
  - `:146` after `to \`REQ-NNN-progress.log\` in the run directory` insert `, by the absolute main-tree path its brief names`.
  - `:73` replace `**allow but never stage** each expected \`REQ-NNN-handback.md\` the run names` with `**allow but never stage** each expected \`REQ-NNN-handback.md\` the run names, every \`REQ-NNN-progress.log\`, and \`full-gate.lock\``.
  - `:173` replace `The one main-tree path a builder may write (*Sole integrator*)` with `One of the two main-tree paths a builder may write (*Sole integrator*)`.
- **F2. The stall loop stops a healthy integrator during a long full gate, and the lock it held blocks the next one.** `:146` logs a line only before and after a long command, so a full gate longer than 20 minutes (the REQ's Why measured 11 to 57 minutes per integration) leaves no new line. `:163` then stops that integrator. Its `full-gate.lock` stays, because only the integrator removes it (`:147`, D-11), the preflight says never delete one (`:158`), and the fresh integrator cannot take the lock. impact-user-visible. **Fix** at `:163`, replace `An integrator silent about 20 minutes or more (no new progress line, no new commit) is stopped, and a fresh integrator resumes` with `An integrator silent about 20 minutes or more (no new progress line, no new commit, and its last progress line does not start a command still inside its usual run time, such as a full gate) is stopped. If it held \`full-gate.lock\`, the coordinator removes that lock in the next writing gap and says so, because it stopped the owner itself. A fresh integrator then resumes`.
- **F4. A coordinated handoff makes the next session integrate leftover claims itself.** `restart-with-parallel-handoff.md:68` puts `advance REQ-NNN`, "then do the phase it names", above the resume command. In a coordinated handoff the fresh session runs those phases in the main session before `--coordinate` is in force. That is the failure the REQ's Why describes (main-session integration growing to about 670k tokens). impact-user-visible. **Fix:** append to `:68`: `In a handoff from a \`--coordinate\` run, write each claim line as \`hand REQ-NNN to an integrator (it enters with advance REQ-NNN)\`: the next session coordinates from its first line and never runs a claim's phase itself (\`actions/fan-out-reference.md\` → **Delegated integration — the coordinator shape**).`
- **F5. The automatic handoff breaks the one-writer rule and refires every turn.** `restart-with-parallel-handoff.md:14` writes queue fields (Step 1) and commits `RESTART-PROMPT.md` (Step 4) whenever context crosses the threshold. Under the coordinator shape the coordinator writes nothing under the project root while an integrator runs (`fan-out-reference.md:139`). Once above the threshold, every later turn is also above it. Its two announcement lines (Step 6) also collide with the stall loop's one closing status line (`:163`). impact-user-visible. **Fix** at `:14`, replace `This automatic handoff writes and commits the handoff only and does not end the session.` with `This automatic handoff waits for the next writing gap (no integrator running), runs once per run (later work falls under the rewrite rule in **Rules**), writes and commits the handoff only, and does not end the session. Its two announcement lines come before the coordinator's turn-status line.`
- **F6. `do-work status` can print an `rm` line for a live `full-gate.lock`.** `status.md:63` (Step 3) looks up "the owner process" named on a run-local file's first line and prints `rm -- '<path>'` when none is found. `full-gate.lock` is now suite-defined and its first line is a writer label, not a process (`fan-out-reference.md:147`), so a `pgrep` on the label finds nothing and the action offers to delete a lock that is in use. `status-guide.md:34` restates the same rule. The Go comment `run_status.go:363-364` ("every run-directory entry the suite does not define") is now inexact for `full-gate.lock`. impact-user-visible. **Fix:** append to `status.md:63`: ``\`full-gate.lock\` (`actions/fan-out-reference.md` → **Coordinated run rules**) names a writer label, not a process: report its holder and age, never an `rm` line.`` Append to `status-guide.md:34`: `A coordinated run's \`full-gate.lock\` names a session, not a process, so it never gets an \`rm\` line.` Go comment → report only.

### Minor

- **F3.** `fan-out-reference.md:163` reads status rows without the label the Constraints require ("coordinator-side judgment ... on this machine, and the prose says so"), and without saying how to read a C3 row, whose remedy is `do-work run` (REQ-690 F4). impact-negligible. **Fix:** after `which status does not read.` insert `Act on C2 rows only for agents this session started, on this machine, as the coordinator's judgment; read a C3 row as "brief the next integrator", never as a \`do-work run\` to execute.`
- **F7.** `fan-out-reference.md:134` "never runs Step 6 to Step 9 itself" conflicts with the coordinator dispatching builders and stamping `dispatch_at`, which is Step 6 (`work.md:289`), and with `:132` ("from the hand-back merge through Step 9"). The text is the REQ's own wording. impact-negligible. **Fix:** replace `never runs Step 6 to Step 9 itself` with `never runs the span from the hand-back merge through Step 9 itself`.
- **F8.** `SKILL.md:4` argument-hint `run [REQ|UR]` shows no run flag, while `status [REQ] [--watch]` does. The builder deferred it because REQ-690, REQ-691 and REQ-692 were editing that line. All three have landed, so the reason is gone. impact-negligible. **Fix:** replace `run [REQ|UR]` with `run [REQ|UR] [--fan-out [N]] [--coordinate]` (re-run `_dev/tests/staged-skills-contract.sh` with the gate).
- **F9.** `work-guide.md:146` Trigger aliases, which `README.md:87` calls "the canonical public list", omits the two new phrases. impact-negligible. **Fix:** after the closing fence of the alias block, before `Use whichever feels natural.`, add the line `\`do-work drive the queue\` and \`do-work use the main session as a coordinator\` run the same loop as \`do-work run --coordinate\`: the main session only coordinates (see **Building several REQs at once** above).`
- **F15 (nit, anti-bloat count).** New surface in the diff: 1 flag, 1 route row, 1 optional file (`do-work/run-policy.md`), 2 run-directory files, 3 manifest columns, 1 `####` subsection, 1 printed degradation line, 1 unrequested standing-preferences row (PD-5, justified by the REQ's Why). 0 helpers, 0 scripts, 0 Go, 0 shipped tests, 0 decorative tests. All within the Constraints' budget. impact-negligible → report only.

### Wave-end sweep findings (not caused by this diff)

- **F10.** `work.md:575` and `:576` cite "Step 8 substep 8" for worktree cleanup. Step 8 has five items, and cleanup is Step 9's last clause, wrapped by `do-work-cli worktree cleanup REQ-NNN` (REQ-660). Present before this wave. impact-negligible. **Fix:** `:575` replace `(Step 8 substep 8)` with `(Step 9; \`do-work-cli worktree cleanup REQ-NNN\` runs it)`; `:576` replace `(Step 6 hand-back, Step 8 substep 8)` with `(Step 6 hand-back, Step 9 cleanup)`.
- **F11.** `fan-out-reference.md:47` "Step 8's `git worktree remove` and `git branch -d` (*Cleanup — happy path*)": that heading (`:99`) says Step 9. impact-negligible. **Fix:** replace `Step 8's` with `Step 9's`.
- **F12.** `do-work-toolbox/actions/ai-report.md:32` (REQ-687 recorded it as `:31`, F6). One-line truth now false: `$ARGUMENTS` can be `--kind proposal|root-cause <topic>`, which is none of the listed forms and not one UR or REQ. impact-negligible. **Fix:** replace the sentence with `` `$ARGUMENTS` is `UR-NNN`, `REQ-NNN`, `most recent`, or blank; the `judge`, `index` and `find` forms below create no report, the `revise` form below writes a new revision of an existing report, and the `--kind` form below writes a decision brief about open work. One report invocation covers one UR or one REQ (a `--kind` brief covers one topic, REQ or UR); a revise covers one existing bundle; blank is the explicit `most recent` form. ``
- **F13.** `do-work-toolbox/actions/architecture-report.md:135` (REQ-687 F7). One-line truth now false: "`ai-report` takes a UR or REQ and presents completed work". It now also takes `--kind` with a topic for open work. impact-negligible. **Fix:** replace the STOP cell with `Keep it here; \`ai-report\` takes a UR or REQ and presents completed work, or with \`--kind\` writes a decision brief about open work`, and in the Because cell replace `a second, incompatible input contract` with `another, incompatible input contract`.
- **F14 (nit).** `do-work-toolbox/docs/ai-report-guide.md:3` says only "one completed UR or REQ" (REQ-687 noted it). The action's own `:3` already has the `--kind` clause. impact-negligible. **Fix:** after `for one completed UR or REQ` insert `; \`--kind proposal|root-cause\` writes a decision brief about open work instead`.
- Still stale, already recorded by their REQ with fix text, unchanged by this diff → report only: REQ-688 F2 `clarify.md:106` and F3 `capture.md:136,138` (fix text in `REQ-688-review.md:33-34`); REQ-690 F8 `prime-do-kanban.md:7` (no `open-work --format json`), F9 `forensics.md:10` and `README.md:176` (still send "stuck" to forensics; fixes: `User suspects something is broken or producing confusing results (a stuck in-flight REQ → \`actions/status.md\`)` and `Run \`do-work status\` to see why in-flight work looks stuck, and \`do-work forensics\` to diagnose failed or broken work.`); REQ-692 M2 `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv:1` and `staged-skills-contract.sh:792` (fix text in `REQ-692-review.md:23`); REQ-660 F3 linked worktrees read as dirty in `verify.go`, `handoff.go`, `cleanup_git.go` (Go). REQ-660 F3 also means this REQ's teardown list and the handoff's REMOVABLE verdict never mark a linked worktree removable.

## Restatement sweep (wave end)

Trigger set: this diff's redefinitions plus the elements on the `**Restatement sweep:**` lines of all 11 siblings (all 11 have one; none unread).

1. **Worktree add, merge and cleanup now wrapped by `do-work-cli worktree` (REQ-660).** Raw git steps are still a valid by-hand route, because `fan-out-reference.md:43`, `:83` and `:101` keep both. Consistent: `fan-out-reference.md:15`, `:25`, `:65`, `:79`, `:101`; `work-reference.md:354` (defer-gate cleanup by hand, same sequence); `restart-with-parallel-handoff.md:56`, `:114`; `cleanup.md:138-140`, `:333`; `background-agents.md:186`. REQ-660 F1 (ignore-line form, `:53`) and F2 (keep-first-`<pre>`, `:83`) stay fixed. Stale: F10, F11 (step pointers, older than this wave).
2. **Stamp and section writers (REQ-659).** This diff writes no stamp and no section. A restarted integrator re-enters through `advance REQ-NNN`, whose phases stamp only absent fields (`work-reference.md:73`). Consistent: `work.md:21`, `:202`, `:289`, `:310`, `:343`, `:376`; `fan-out-reference.md` Landed hand-back and Dispatch instant (hand writes, still valid). REQ-659 F4 (`work-reference.md:61`) and F5 (`:73`) stay fixed.
3. **Finalization manifest (REQ-658).** `work.md:463` (Step 8 item 5), `:467` (Step 9), `:500-501` (checklist) and `work-reference.md:719` agree on `finalize --auto-manifest ... --emit`. REQ-658 F3 is fixed. The coordinator never runs Steps 8 and 9, so nothing here restates them.
4. **ai-report completed-work-only wording (REQ-687).** `ai-report.md:32` and `architecture-report.md:135` are both one-line truths that are now false: F12, F13 (and F14 for the guide). `ai-report-reference.md:3` is still true: it describes the visual machinery and points completed-work resolution elsewhere. Also checked: `ai-report.md:28`, `:56`, `:157`, `:173`, guide `:21`, `:55`. Earlier ai-report siblings' items are now fixed: REQ-656 F1 (revise `fail` allowed at `:157`), F2 (`:56` now has kinds to refer to), F3 (`:28`).
5. **Core routing table and argument-hint (REQ-690, 691, 692, 689).** First match wins. The new row (`:33`) sits under run-with-recovery (`:32`) and above run (`:34`). None of `:30-32` contains a trigger the new phrases match, and the validate-feedback (`:36`), status (`:40`) and trace (`:49`) rows do not overlap with them. Argument-hint carries `status [REQ] [--watch]` and `trace ...`. `validate-feedback` is a toolbox action, so it is left out of the core hint by design. Stale: F8.
6. **Coordinator-shape rules.** Consistent: `work-reference.md:396` (D-07: its integrator rules did not change and it points to fan-out-reference for the shape); `background-agents.md` (D-08: no argument list, no run-directory table, ceiling note kept); `work.md:37`, `:105`, `:108`, `:111`, `:471`, `:482`, `:502`; `work-guide.md:134`; `standing-preferences.md:14`; `status.md:11`; `fan-out-reference.md` Mid-Run Messages (the coordinator keeps words in scratch until a gap). Stale: `fan-out-reference.md:61`, `:73`, `:173` (F1); `restart-with-parallel-handoff.md:14` (F5), `:68` (F4); `status.md:63`, `status-guide.md:34` (F6); `work-guide.md:146` (F9).
7. **Other sibling elements.** REQ-688's fence rule: `capture.md:129,236`, `abandon.md:66,99`, `fan-out-reference.md` Mid-Run Messages (`:201`, previously `:174`) are consistent; `clarify.md:106` and `capture.md:136,138` are still stale (recorded). REQ-690's routing words: `forensics.md:10` and `README.md:176` still stale (recorded). REQ-691's subcommand set: F1 and F2 were fixed at integration (`main.go` header, `prime-do-kanban.md:3`). REQ-692's promotion authority: I2 fixed (`capture-reference.md:42`, `:55`); M2 test-side text still stale (recorded).

## Review

**Overall: 86%** | 2026-10-10T18:38:00Z

**Verdict:** Approve after fixes F1, F2, F4, F5 (one sentence each; exact text above).

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 75% |
| Test Adequacy | 85% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 progress log contradicts the builder's one-exception write rule (`fan-out-reference.md:61`, `:73`, `:146`, `:173`); fix text in report — impact-user-visible → fix
- F2 stall loop stops an integrator mid full gate and strands `full-gate.lock` (`fan-out-reference.md:163`); fix text in report — impact-user-visible → fix
- F4 coordinated handoff's claim lines make the next session integrate in the main session (`restart-with-parallel-handoff.md:68`) — impact-user-visible → fix
- F5 automatic handoff writes during an integrator's span, refires each turn (`restart-with-parallel-handoff.md:14`) — impact-user-visible → fix
- F6 `status` offers `rm` for a live `full-gate.lock` (`status.md:63`, `status-guide.md:34`); Go comment `run_status.go:363` report only — impact-user-visible → fix

**Minor findings:** F3 stall loop lacks the this-machine label and C3 handling (`fan-out-reference.md:163`) — impact-negligible → fix; F7 "never runs Step 6 to Step 9" vs dispatch in Step 6 (`:134`) — impact-negligible → fix; F8 argument-hint lacks run flags (`SKILL.md:4`) — impact-negligible → fix; F9 trigger-alias list omits the new phrases (`work-guide.md:146`) — impact-negligible → fix; F10 "Step 8 substep 8" (`work.md:575-576`) — impact-negligible → fix; F11 "Step 8's" cleanup (`fan-out-reference.md:47`) — impact-negligible → fix; F12 `ai-report.md:32` Input sentence — impact-negligible → fix; F13 `architecture-report.md:135` — impact-negligible → fix; F14 (nit) `ai-report-guide.md:3` — impact-negligible → fix; F15 (nit) anti-bloat count, all within budget — impact-negligible → report only
**Acceptance:** Pass — GREEN probe at HEAD exit 0 (flag, route, rules, rows, columns, preflight, stall loop, resume line, teardown present; no `pid` in fan-out-reference.md)
**Restatement sweep:** redefined the run argument set (`--coordinate`), the coordinator shape (mandatory under the flag), the run-directory file set (`REQ-NNN-progress.log`, `full-gate.lock`) and three manifest columns, the builder's main-tree writes, the handoff triggers and resume command, the Step 10 teardown, and the core run route phrases; stale: F1, F4, F5, F6, F8, F9; inherited elements of all 11 siblings swept (none unread): F10, F11, F12, F13, F14 new, earlier recorded REQ-688 F2/F3, REQ-690 F8/F9, REQ-692 M2, REQ-660 F3 still stale
**Suggested testing:** 2 items — (1) read-through of a coordinated handoff paste block with one leftover claim after F4; (2) `do-work status` against a run directory holding a `full-gate.lock` with a label first line after F6
**Follow-ups created:** None (15 findings report only or fix-in-integration)

*Reviewed by review-work action*

## Delta re-check

**Overall: 95%** | 2026-10-10T18:38:00Z | delta `f8f32682..9cd669eb` (10 files, +19 -17)

**Verdict:** Approve. F1 to F14 say what the review meant. The GREEN probe at `9cd669eb` exits 0. No core test reads the core `argument-hint` (`staged-skills-contract.sh:829` checks only the board's), so F8 breaks no test. The changed F5 wording is right: Step 6 and the Red Flags require the message to end with the two announcement lines, so the status line goes before them.

Cross-checks that hold: step 0 (`fan-out-reference.md:73`) now allows the progress logs and the lock, and the "this REQ owns" stop at `:76` no longer catches them. The preflight's "never delete" (`:158`) covers locks left by earlier runs, before the first spawn. The stall loop's lock removal (`:163`) covers only an owner this coordinator stopped itself, so the two do not conflict. `work.md` Step 10 (`:471`) is unaffected. `status-guide.md:34` and `status.md:63` agree. The Verification Checklist in restart-with-parallel-handoff ("one `advance REQ-NNN` line per working claim") still holds for the coordinated claim line, which contains `advance REQ-NNN`.

**New findings:**
- N1 `fan-out-reference.md:163` still says "Each coordinator turn ends with one line", which contradicts the F5 wording at `restart-with-parallel-handoff.md:14` for the automatic-handoff turn — impact-negligible → fix: replace `Each coordinator turn ends with one line: done/total, active lanes, ETA or unknown.` with `Each coordinator turn ends with one line: done/total, active lanes, ETA or unknown (in a turn that writes the automatic handoff, the handoff's two announcement lines follow it).`
- N2 (nit) `fan-out-reference.md:73` calls `full-gate.lock` "builder scratch", but an integrator writes it — impact-negligible → fix: replace `— builder scratch in the main tree, not an error;` with `— run scratch in the main tree, not an error;`
- N3 (nit) The progress-log rule (`:146`) cannot apply to a remote builder, which is never given a main-tree path (`:26`). Such a builder has no log, and the stall loop reads only its commits. This gap existed before the delta — impact-negligible → report only
