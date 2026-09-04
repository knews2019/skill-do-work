---
id: REQ-258
title: Split the prescribed shell behavior suite per script
status: completed
completed_at: 2026-08-20T08:39:00Z
commit: 1cc1836
claimed_at: 2026-08-20T08:17:28Z
route: B
created_at: 2026-08-18T17:49:24Z
status_changed_at: 2026-08-18T20:55:14Z
user_request: UR-056
addendum_to: REQ-246
domain: general
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-20T08:17:28Z
  basis:
    - trivial short-circuit
write_set:
- _dev/tests/prescribed-shell-scripts-behavior.sh
- _dev/tests/prescribed-shell-harness.sh
- _dev/tests/prescribed-shell-cases/add-local-git-exclude.sh
- _dev/tests/prescribed-shell-cases/atomic-download.sh
- _dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh
- _dev/tests/prescribed-shell-cases/capture-screenshot.sh
- _dev/tests/prescribed-shell-cases/cleanup-req-reservations.sh
- _dev/tests/prescribed-shell-cases/generate-report-image-batch.sh
- _dev/tests/prescribed-shell-cases/generate-report-image.sh
- _dev/tests/prescribed-shell-cases/install-last30days.sh
- _dev/tests/prescribed-shell-cases/install-memory-hooks.sh
- _dev/tests/prescribed-shell-cases/lexical-memory-recall.sh
- _dev/tests/prescribed-shell-cases/protected-inventory.sh
- _dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh
- _dev/tests/prescribed-shell-cases/qualify.sh
- _dev/tests/prescribed-shell-cases/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-cases/run-blocked-check.sh
- _dev/tests/prescribed-shell-cases/show-commit-diff.sh
- _dev/tests/prescribed-shell-cases/stage-exact-deletion.sh
kb_status: promoted
kb_entry: REQ-258-split-the-prescribed-shell-behavior-suit.md
---

# Split the Prescribed Shell Behavior Suite Per Script

## What

`_dev/tests/prescribed-shell-scripts-behavior.sh` now carries 47 named cases, and the ten reservation-cleanup + timestamp-repair cases dominate its tail. If it keeps growing, per-script files may read and fail more legibly. Organizational only — no case changes.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md` plus general/coding-guardrails/communication-style/maintenance crew. Approach: (1) lift the existing 30-line preamble into a new `_dev/tests/prescribed-shell-harness.sh` and give it a `prescribed_shell_finish` closer that derives its own case count at run time; (2) split the 76 case blocks into `_dev/tests/prescribed-shell-cases/<script>.sh` mechanically, gathering the two interleaved groups and moving the misfiled portfolio setup block by content match; (3) keep `prescribed-shell-scripts-behavior.sh` at its path as a runner that executes each case file as its own process, aggregates failures, and greps the case files for the total. Verify: line-multiset parity against the original case region, full suite green, each file green standalone, a deliberate breakage turns the runner red, `shellcheck --severity=warning` clean, `maintainer-verify.sh` exit 0.
- [x] **[APPLY]:** 19 files, all inside the declared Scope. Split performed by a one-shot extraction script rather than by hand, so no case block was retyped.
- [x] **[UNIFY]:** `git diff --stat` → 19 files, +2025/−1870, matching the Scope list exactly. Verified: line-multiset comparison of every non-blank case line, original vs split — 1756 in, 1756 out, zero missing, zero added. `shellcheck --severity=warning` clean across the harness, the runner, and all 17 case files. No debug artifacts introduced (the only `TODO`/`console` strings in the new files are pre-existing fixture *data* inside the `qualify` cases, moved verbatim). Deliberate-breakage probe: mutating one assertion in `qualify.sh` made the runner print the failing case, name the failing file, and exit 1; restored and re-verified green.

## Context

Builder-discovered on REQ-246 (Discovered Tasks, third item), classified [low]. Touches only a test file but is a reorganization, not mechanical hygiene, so it takes the consent flow rather than the test-hygiene carve-out.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-246: the shell behavior suite is one growing file and could split per-script. Should I process this as a new task? → Yes, add to queue — user chose the non-recommended branch: the suite is worth splitting per script now rather than later
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — one file is fine until it actually hurts.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.

<!-- D-XX counter: none used. Next decision: D-01. -->

---

## Triage

**Route: B** - Medium

**Reasoning:** The "what" is fully specified and the target file is named, but the split boundaries, the shared fixture preamble every case block depends on, the runner wiring, and the self-counting tail all need discovery before any file is created. Not Route A — this is a wholesale restructure of 1882 lines, not a value change in a named file.

**Planning:** Not required

---

## Exploration

Run inline in the orchestrator session (this harness has no separate explore subagent).

- **F1 — the file holds 76 named cases across 17 script groups, not 47.** The REQ's figure is stale; the file self-counts at run time, so no restatement had to be updated to make it stale. Group sizes: `repair-req-timestamps` 20, `audit-archive-timestamps` 11, `install-last30days` 7, `publish-portfolio-summary` 7, `cleanup-req-reservations` 6, `generate-report-image` 5, `atomic-download` 4, `qualify` 3, `capture-screenshot`/`run-blocked-check`/`stage-exact-deletion`/`generate-report-image-batch` 2 each, and `show-commit-diff`/`add-local-git-exclude`/`protected-inventory`/`lexical-memory-recall`/`install-memory-hooks` 1 each.
- **F2 — zero cross-group coupling.** No group reads a variable another group assigns, and both helper functions (`process_runs_unreaped_excluded`, `run_ai_report_batch_replay`) are used only inside the group that defines them. The one boundary defect: the six-line `publish-portfolio-summary` fixture setup (`portfolio_root`/`portfolio_source`/`portfolio_canonical`/`portfolio_candidate`) sits *below* the `install-memory-hooks` case header, so a naive header-boundary cut files it under the wrong script. It moves with the portfolio group.
- **F3 — two groups are interleaved and must be gathered.** `generate-report-image` (lines 465-529, 721-869) is split around `generate-report-image-batch` (530-720), and `repair-req-timestamps` (1098-1410, 1659-1805) is split around `audit-archive-timestamps` (1411-1658). Order *within* each group is preserved; the groups themselves get gathered. This is the reordering the split makes visible.
- **F4 — one consumer, exit status only.** `_dev/tests/staged-skills-contract.sh:183` is the sole live invocation (reached from `maintainer-verify.sh` → `contract-regressions.sh` → `staged-skills-contract.sh`), and REQ-234 already established that nothing parses the closing count line. Keeping the runner at its current path means no caller edit.
- **F5 — every new file is linted.** `maintainer-verify.sh` runs `shellcheck --severity=warning` over every tracked `*.sh`, so each new case file must pass. The existing `# shellcheck source=` directive is the pattern for the sourced-harness line.
- **F6 — the self-count is a stated contract.** The tail comment says the number "is that shape grepped out of this file at run time — so the reported number and the file cannot disagree, and nothing here is a remembered figure." The split has to keep that property by grepping the case files instead of the runner.

*Generated inline (Route B exploration)*

---

## Scope

**Files I will touch:**

- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify — keeps its path and exit-status contract, becomes the runner)
- `_dev/tests/prescribed-shell-harness.sh` (new — the preamble every case file shares)
- `_dev/tests/prescribed-shell-cases/add-local-git-exclude.sh` (new)
- `_dev/tests/prescribed-shell-cases/atomic-download.sh` (new)
- `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh` (new)
- `_dev/tests/prescribed-shell-cases/capture-screenshot.sh` (new)
- `_dev/tests/prescribed-shell-cases/cleanup-req-reservations.sh` (new)
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh` (new)
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` (new)
- `_dev/tests/prescribed-shell-cases/install-last30days.sh` (new)
- `_dev/tests/prescribed-shell-cases/install-memory-hooks.sh` (new)
- `_dev/tests/prescribed-shell-cases/lexical-memory-recall.sh` (new)
- `_dev/tests/prescribed-shell-cases/protected-inventory.sh` (new)
- `_dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh` (new)
- `_dev/tests/prescribed-shell-cases/qualify.sh` (new)
- `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` (new)
- `_dev/tests/prescribed-shell-cases/run-blocked-check.sh` (new)
- `_dev/tests/prescribed-shell-cases/show-commit-diff.sh` (new)
- `_dev/tests/prescribed-shell-cases/stage-exact-deletion.sh` (new)

**Acceptance criteria:**

1. One case file per script under test, each holding only that script's fixture blocks.
2. **No case changes.** Every one of the 76 named cases keeps its fixtures, assertions, and failure strings verbatim. The only permitted edits are the moves the split requires (gathering interleaved groups, relocating the portfolio setup block) and the per-file preamble/closer.
3. The suite still runs from `_dev/tests/prescribed-shell-scripts-behavior.sh` with the same exit-status contract, so `_dev/tests/staged-skills-contract.sh:183` needs no edit, and the full `maintainer-verify.sh` chain passes.
4. The reported case count stays computed at run time from the case files, never a remembered figure (F6).
5. Every new file passes `shellcheck --severity=warning` (F5).
6. Each case file is runnable on its own (`bash _dev/tests/prescribed-shell-cases/repair-req-timestamps.sh`) so a REQ touching one script can run that script's cases alone.

---

## Implementation Summary

**What was done:** Lifted the suite's shared fixture preamble into a new harness file, split all 76 named cases into one file per script under test, and reduced `prescribed-shell-scripts-behavior.sh` to a runner that executes each case file as its own process, aggregates failures, and derives the reported case count from the case files at run time. No case was changed: the assertions, fixtures, and failure strings all moved verbatim, verified by a line-multiset comparison (1756 non-blank case lines in, 1756 out, zero missing, zero added). Two groups that were interleaved in the monolith (`generate-report-image` around `generate-report-image-batch`, `repair-req-timestamps` around `audit-archive-timestamps`) are now each contiguous in their own file, and the `publish-portfolio-summary` fixture setup that had been sitting under the `install-memory-hooks` header moved to the file it belongs to. Each case file is runnable on its own.

**Files changed:**

- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified) — reduced from 1882 lines to a 35-line runner; keeps its path and exit-status contract so `staged-skills-contract.sh:183` needs no edit
- `_dev/tests/prescribed-shell-harness.sh` (new) — the shared preamble (fixture root, cleanup trap, `fail_case`, script roots) plus `prescribed_shell_finish`, which prints the sourcing file's run-time-derived case and failure counts and returns that file's status
- `_dev/tests/prescribed-shell-cases/add-local-git-exclude.sh` (new) — 1 case
- `_dev/tests/prescribed-shell-cases/atomic-download.sh` (new) — 4 cases
- `_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh` (new) — 11 cases
- `_dev/tests/prescribed-shell-cases/capture-screenshot.sh` (new) — 2 cases
- `_dev/tests/prescribed-shell-cases/cleanup-req-reservations.sh` (new) — 6 cases
- `_dev/tests/prescribed-shell-cases/generate-report-image-batch.sh` (new) — 2 cases
- `_dev/tests/prescribed-shell-cases/generate-report-image.sh` (new) — 5 cases
- `_dev/tests/prescribed-shell-cases/install-last30days.sh` (new) — 7 cases
- `_dev/tests/prescribed-shell-cases/install-memory-hooks.sh` (new) — 1 case
- `_dev/tests/prescribed-shell-cases/lexical-memory-recall.sh` (new) — 1 case
- `_dev/tests/prescribed-shell-cases/protected-inventory.sh` (new) — 1 case
- `_dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh` (new) — 7 cases, plus the relocated fixture setup
- `_dev/tests/prescribed-shell-cases/qualify.sh` (new) — 3 cases
- `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` (new) — 20 cases
- `_dev/tests/prescribed-shell-cases/run-blocked-check.sh` (new) — 2 cases
- `_dev/tests/prescribed-shell-cases/show-commit-diff.sh` (new) — 1 case
- `_dev/tests/prescribed-shell-cases/stage-exact-deletion.sh` (new) — 2 cases

## Decisions

- **D-01 — DECIDE & STATE. Case files are standalone processes over a sourced harness, not blocks sourced into one runner process.** Exploration proved zero cross-group coupling (F2), so per-file isolation costs nothing and buys the thing the REQ actually asked for: `bash _dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` runs 20 cases without the other 56, which is exactly the loop REQ-263, REQ-264 and REQ-271 will each want. Reversible — collapsing back to sourcing is a runner-only edit.
- **D-02 — DECIDE & STATE. The runner keeps its existing path and exit-status contract.** `staged-skills-contract.sh:183` and the `maintainer-verify.sh` chain above it need no edit, so the split cannot break the canonical baseline through a caller nobody remembered.
- **D-03 — DECIDE & STATE. Case files are named for the script they cover, `qualify.sh` included.** The single-word name mirrors the shipped `tools/checks/qualify.sh` rather than introducing a name, which is what keeps "grep the script name, find its cases" true (coding-guardrails § 5 precedence against § 3: § 5 governs names you introduce).
- **D-05 — DECIDE & STATE. Replaced the pluralization `&&`/`||` subshell chain in `prescribed_shell_finish` with plain `if` blocks.** Self-review flagged `$([ "$n" -eq 1 ] && printf case || printf cases)` as unnecessary cleverness in a test harness — the `&&`/`||` ternary is a known footgun when the first branch can fail. Same output, no chain.

- **D-04 — DECIDE & STATE. The `SC2034` disable in the harness is scoped to the three script-root assignments, not the file.** ShellCheck cannot see the sourcing files that consume them, so the warning is structural; a file-level disable would also hide a genuine unused variable added later.

## Discovered Tasks

- **qualify.sh cannot tell a moved line from an added one, so it misfires on every relocation REQ.** `tools/checks/qualify.sh` runs `git diff` with no `-M`/`-C` and greps `^+`, so relocated text reads as newly added. On this REQ it FAILed the P-A-U audit on four `TODO` strings that are deliberate fixture data in the REQ-254 `qualify` cases, byte-identical to their pre-change form in HEAD. Every future extract-to-module or file-split REQ hits the same false positive, and the failure mode is the dangerous direction: a builder under time pressure learns to un-check `[UNIFY]` or wave the FAIL away, which is exactly the reflex the audit exists to prevent. Candidate fix: enable rename/copy detection (`git diff -C --find-copies-harder`) for the debug-artifact scan, or subtract lines that appear verbatim in the pre-change tree.

- `effort_estimate: trivial` on this REQ produced a 5-minute P50 through Step 3.6's mechanical short-circuit, for a wholesale restructure of 1882 lines across 19 files. The field is size and was misjudged at capture. Worth a look at whether capture's effort judgment has a systematic bias on reorganization REQs, where the *per-case* change is trivial but the file count is not.
- The REQ text says the suite "now carries 47 named cases"; it carries 76. Nothing broke, because REQ-234 already made the count derived rather than remembered — but a REQ body that quotes a volatile figure will keep going stale between capture and build. No action needed on this REQ; noting the pattern.

---

## Qualification

**Passed with one FAIL overridden on evidence** — 19 files verified on disk, 6 acceptance criteria traced, P-A-U confirmed.

- **Mechanical (script):** `tools/checks/qualify.sh` FAILed the `[UNIFY]` audit on four `TODO` strings in `_dev/tests/prescribed-shell-cases/qualify.sh`. Overridden, with evidence, not with reasoning: those four lines exist byte-identically in `git show HEAD:_dev/tests/prescribed-shell-scripts-behavior.sh` (lines 1860, 1867, 1869, 1871), and `diff` over every `TODO`-bearing line in HEAD's monolith against every one in the split reports no difference. They are fixture *data* for the REQ-254 case that pins "a fresh TODO in a reporter file still FAILs". The script has no rename or copy detection, so a pure move reads as an addition — queued as a Discovered Task rather than worked around silently. The accompanying WARN (added output primitives in a file that owns its process exit) is the script's own reporter exemption and is correct.
- **Substantive (2):** every new file carries real fixture code — 14 to 466 lines — not boilerplate. No placeholders, no empty exports.
- **Requirements traced (3):** one file per script (17 files, criterion 1); no case changes, proven by line-multiset parity, 1756 in / 1756 out (criterion 2); entry point and exit-status contract unchanged, `staged-skills-contract.sh` untouched (criterion 3); the reported count is grepped from the case files at run time (criterion 4); `shellcheck --severity=warning` exits 0 (criterion 5); all 17 files pass standalone (criterion 6).
- **Flowing (6):** not a data-path change, but the runner was proven to actually run rather than pass vacuously — mutating one assertion in `qualify.sh` made it print the failing case, name the failing file, and exit 1.

---

## Testing

**Tests run:**

- `bash _dev/tests/prescribed-shell-scripts-behavior.sh` — exit 0, 76 named script cases across 17 per-script files
- `bash _dev/tests/prescribed-shell-cases/<each>.sh` — all 17 exit 0 standalone
- `shellcheck --severity=warning -- _dev/tests/prescribed-shell-harness.sh _dev/tests/prescribed-shell-scripts-behavior.sh _dev/tests/prescribed-shell-cases/*.sh` — exit 0
- `bash _dev/tests/maintainer-verify.sh` — **exit 0**, including the `contract-regressions.sh` → `staged-skills-contract.sh` → suite chain that consumes the runner's exit status

**Regression evidence (this is a reorganization, so red-green does not apply to behavior):**

- **Content parity.** Every non-blank line of the original case region compared as a multiset against the union of the 17 case files: 1756 in, 1756 out, zero missing, zero added. No case's fixtures, assertions, or failure strings changed.
- **The suite still bites.** A disabled suite and a working one print the same thing on a clean tree, so the runner was proven to fail: one assertion in `qualify.sh` was mutated to compare against a string the output never contains, and the run printed the failing case's `FAIL:` line, reported `qualify: 3 cases, 1 failure.`, printed `FAIL: 1 of 17 per-script case files reported failures.`, and exited **1**. Restored and re-verified green. (Test shape borrowed from REQ-282's lesson.)
- **The count cannot go stale.** The reported figure is grepped from the case files at run time by the runner, and each file's own figure is grepped from that file by `prescribed_shell_finish` — no remembered number anywhere (REQ-234's lesson).

**Baseline:** pre-flight recorded the suite green before any change (`preflight.sh`, working tree clean outside `do-work/`), so every result above is a clean-to-clean comparison. No pre-existing failures to exclude.

---

## Review

**Overall: 94%** | Requirements 100% | Code Quality 92% | Test Adequacy 95% | Scope Discipline 100% | Risk: Low | **Acceptance: Pass**

**Approve with follow-ups** — the split is content-preserving and proven to still bite; the restatement sweep found live planning text and three queued `write_set` fields that still describe the file this REQ dissolved.

### What's built

`_dev/tests/prescribed-shell-scripts-behavior.sh` is now a 35-line runner over 17 per-script case files in `_dev/tests/prescribed-shell-cases/`, sharing one fixture preamble. Every one of the 76 cases moved verbatim, each file runs on its own, and the entry point and exit-status contract are unchanged so nothing upstream needed an edit.

### Findings

**Important**

1. **`do-work/RESTART-PROMPT.md:33` and the `write_set` of REQ-263, REQ-264 and REQ-271 all still name the dissolved monolith.** The restart prompt tells the next session that `prescribed-shell-scripts-behavior.sh` is "the whole bottleneck" written by five REQs at "at most one per wave" — the exact planning constraint this REQ removes. REQ-263 and REQ-264 now write `prescribed-shell-cases/qualify.sh`, REQ-271 writes `prescribed-shell-cases/repair-req-timestamps.sh`, and those three are mutually disjoint, so they can run in one wave rather than three. One root cause (the bottleneck-file assumption is obsolete), so one sweep follow-up rather than four. **Correctness is not at risk** — `write_set` is display-only and Step 5.5 overwrites it from the fresh Scope declaration at claim time — but a human planning waves off the board reads the stale value, which is precisely what the restart prompt is for.

**Minor**

2. **`decisions/audits/2026-08-11-defensive-surface.md` cites the runner as the coverage for eleven KEEP decisions.** Each row names "`prescribed-shell-scripts-behavior.sh` <name> case", and the named case now lives one hop away in the case directory. Left unedited on purpose: it is a dated decision record, and rewriting dated history is the call REQ-234 already declined for `CHANGELOG.md`. Folded into the same sweep follow-up.

**Nit**

3. `REQ-271`'s Red-Green Proof quotes "66 named cases" as today's observation; the figure was already stale before this REQ (76) and the command it belongs to still works. No action.

### Requirements Checklist

| Requirement | Status | Evidence |
|---|---|---|
| One case file per script under test | Delivered | 17 files, one per script group |
| No case changes | Delivered | line-multiset parity, 1756 in / 1756 out, zero missing, zero added |
| Same entry point and exit-status contract | Delivered | `staged-skills-contract.sh` untouched; `maintainer-verify.sh` exit 0 |
| Count derived at run time, never remembered | Delivered | runner greps the case files; `prescribed_shell_finish` greps its own |
| `shellcheck --severity=warning` clean | Delivered | exit 0 across all 19 files |
| Each case file runnable standalone | Delivered | all 17 exit 0 individually |

### Acceptance Testing

**Pass.** Beyond the green run: a mutated assertion in `qualify.sh` produced the failing case line, `qualify: 3 cases, 1 failure.`, `FAIL: 1 of 17 per-script case files reported failures.`, and exit 1. Both runner guards were exercised rather than assumed — an empty case directory prints `FAIL: no per-script case files found` and exits 1, and a case file with a syntax error is counted as a failing file and exits 1.

### Suggested Additional Testing

- **Regression scenario:** the next REQ that adds a case (REQ-263 or REQ-264) is the real proof that the new layout is usable, not just correct. Watch whether its Step 5.5 Scope lands on the case file without prompting.
- Nothing else — this change touches test infrastructure only, has no runtime surface, and handles no user input.

### Follow-up REQs Created

- One sweep REQ for findings 1 and 2 (same root cause: text that still describes the pre-split layout).

*Restatement sweep: run. The diff redefines where the prescribed-shell cases live. Swept every live reference to the path — `staged-skills-contract.sh:183` (invokes the runner, still correct), `skills/do-work/scripts/repair-req-timestamps.sh:70` (still true), `do-work/audits/audit-2026-08-14.md` (claim is about invocation, still accurate), `do-work/runs/*` and `do-work/archive/*` and `ai-reports/*` (history, correctly left alone). The stale live ones are findings 1 and 2.*

---

## Lessons Learned

**What worked:**
- Splitting by a one-shot extraction script instead of by hand, then proving the result with a **line-multiset comparison** against the original. 1756 in, 1756 out, zero drift — a claim of "no case changes" that is checkable rather than asserted. Hand-splitting 1882 lines across 17 files would have made the same claim unverifiable.
- Auditing cross-group variable and function coupling *before* choosing the architecture. The answer (zero coupling) is what made per-file processes safe, and it took one script to learn rather than a failed run to discover.

**What didn't:**
- The first pass at `prescribed_shell_finish` used `$([ "$n" -eq 1 ] && printf case || printf cases)` for pluralization. It works, but the `&&`/`||` ternary breaks silently the moment the first branch can fail, and it was cleverness bought for a cosmetic. Replaced with plain `if` blocks (D-05).
- Hoisting `core_checks` into the shared harness looked tidier than leaving it inside the `qualify` block. It would have been a case change and an unused-variable warning, for nothing. Reverted before it shipped: in a "no case changes" REQ, the tidying reflex is the failure mode.

**Worth knowing:**
- **`tools/checks/qualify.sh` cannot tell a moved line from an added one.** It runs `git diff` with no `-M`/`-C` and greps `^+`, so a relocation REQ gets FAILed on every pre-existing `TODO`/debug marker inside the moved text. This REQ hit it on four fixture strings that are byte-identical in HEAD. The correct response is to prove the lines pre-exist (`git show HEAD:<file> | grep`) and record the override with that evidence — not to un-check `[UNIFY]`, which is the reflex the audit exists to prevent. Queued as a discovered task.
- **A restructure's real blast radius is the text that plans around the old structure, not the code.** Every consumer of this file was fine (one caller, exit status only). What went stale was `RESTART-PROMPT.md` and three queued `write_set` fields describing a scheduling bottleneck that no longer exists. `write_set` self-heals at Step 5.5, so the damage is to human wave-planning, not to the build.
- Groups were interleaved in the monolith (`generate-report-image` around its `-batch` sibling, `repair-req-timestamps` around `audit-archive-timestamps`) and one fixture-setup block sat under the wrong case header entirely. Nobody could see either while it was one file. That is the concrete argument for the split, and it is worth stating more strongly than the REQ's "may read more legibly".

## Orientation

You can now run one script's prescribed-shell proofs on their own — `bash _dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` runs 20 cases instead of all 76 — and a failing run names the file as well as the case. Lives in the maintainer test suite (`_dev/tests/`), under `_dev/primes/prime-shell-commands.md`'s domain.

**[MAP CHANGED]** The prescribed-shell behavior suite is no longer one file. `_dev/tests/prescribed-shell-scripts-behavior.sh` keeps its name and its exit-status contract but is now a dispatcher; the cases live one file per script in `_dev/tests/prescribed-shell-cases/`, over a shared preamble in `_dev/tests/prescribed-shell-harness.sh`. **A REQ that adds a case now writes that script's case file, not the runner** — which also means several such REQs no longer collide, so they can share a wave. Prime staleness spot-check: `_dev/primes/prime-shell-commands.md` names no path this REQ moved, so it is not stale.
