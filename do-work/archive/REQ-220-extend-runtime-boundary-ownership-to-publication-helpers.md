---
id: REQ-220
title: Extend runtime-boundary ownership to the remaining publication helpers
status: completed
claimed_at: 2026-08-17T20:35:52Z
status_changed_at: 2026-08-17T19:58:23Z
route: B
completed_at: 2026-08-17T20:54:49Z
commit: 
domain: general
created_at: 2026-08-17T18:37:31Z
user_request: UR-042
addendum_to: REQ-204
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: runtime-boundary-ownership
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
maintenance: false
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-17T20:36:40Z
  basis:
  - Route B
  - 3-file write set
  - 2 subsystems involved
  - 5 acceptance criteria
  - async lifecycle behavior
  - full-suite verification
write_set:
- skills/do-work-toolbox/scripts/generate-report-image.sh
- skills/do-work-toolbox/scripts/install-last30days.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
- skills/do-work/docs/prescribed-shell-primitives.md
- skills/do-work-toolbox/actions/ai-report-reference.md
---

# Review Fix: Extend Runtime-Boundary Ownership to the Remaining Publication Helpers

## What

Close the same runtime-boundary defect class in the two shipped scripts REQ-204's copy-paste sweep found still carrying it: a helper that cleans files without owning the processes it started, and a directory publication that reads `mv`'s exit status as proof the rename landed where it was aimed.

## Context

REQ-204 fixed both boundaries in the `ai-report` prescribed batch block. `_dev/primes/prime-shell-commands.md` requires grepping a fixed primitive across every action before calling it fixed, because these blocks are copy-pasted. That sweep found the same root cause — verifying a convenient proxy instead of owning the complete boundary, exactly as `_dev/lessons/validated-runtime-boundaries.md` frames it — alive in two more places. Neither was in REQ-204's declared scope, so both were routed here rather than fixed inline.

## Instances

- [x] `skills/do-work-toolbox/scripts/generate-report-image.sh`: backgrounds `imagegen` (line 54) and a watchdog subshell (line 73), but its HUP/INT/TERM traps are bare `trap 'exit N'`. Interrupted through REQ-204's batch it is killed by process group; interrupted as a direct invocation its backend outlives it.
- [x] `skills/do-work-toolbox/scripts/install-last30days.sh`: `mv "$staging_directory" "$target_directory"` (line 94) follows a check-and-backup window. A target that reappears inside that window swallows the staging tree, `mv` exits 0, and the script sets `publication_complete=1` over a nested tree.
- [x] `skills/do-work-toolbox/scripts/generate-report-image.sh` (line 48): `mv "$staged_output_path" "$output_path"` nests if the output path is ever a directory. Minor — both paths are caller-controlled inside private staging — but it belongs to the same audit.

## Requirements

- Reuse the already-proven idioms rather than inventing new ones: the verified-process-group pattern from `skills/do-work/scripts/run-blocked-check.sh`, and REQ-204's post-rename nesting verification.
- A helper that cannot prove process-group isolation must signal only the bare PID, never a group.
- A publication that nests must exit nonzero, leave the colliding destination byte-for-byte, and leave no private stage behind.
- Preserve every existing success path and exit status of both scripts.
- Add fixture replays to `_dev/tests/prescribed-shell-scripts-behavior.sh` for each instance, following the two REQ-204 added.

## Red-Green Proof

**RED prompt/case:** Interrupt `generate-report-image.sh` directly while its backend runs and observe the backend survive; run `install-last30days.sh` with a `mv` shim that recreates the target inside the check-then-rename window and observe it report a complete publication over a nested tree.
**Why RED now:** File cleanup does not own a process tree, and `mv`'s exit status does not prove exact-path publication when the destination is a directory. Both scripts ship to users today.
**GREEN when:** The interruption replay proves no launched process survives a direct invocation; the collision replay returns nonzero, preserves the destination byte-for-byte, and leaves no nested or private stage; and every pre-existing named case in the suite still passes.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] REQ-204 fixed both runtime boundaries in the ai-report batch, and the required copy-paste sweep found the same root cause in two more shipped scripts. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: REQ-204 was itself a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle. Worth knowing before you decide: these two scripts ship to users, and neither defect is reachable from the path REQ-204 already fixed — the first only bites a direct invocation, the second only a reappearing target.

<!-- D-XX counter: last used D-03. Next decision: D-04. -->

---

## Triage

**Route: B** - Medium

**Reasoning:** The REQ names all three defect sites with line numbers and names the two idioms to reuse, but the exact shape of those idioms — how much of `run-blocked-check.sh`'s verified-process-group pattern applies to a helper that already owns its backend, and how the nesting-failure branch must interact with `install-last30days.sh`'s existing rollback machinery — needed discovery before coding.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Defect sites (verified on disk):**

- `skills/do-work-toolbox/scripts/generate-report-image.sh:54` backgrounds `imagegen` and `:69-73` backgrounds a subshell that runs `codex exec`. Neither launch enables job control, so both share the script's own process group, and `cleanup_report_image_paths` signals only `$backend_process_id` — the direct child. For the agentic path that PID is the *subshell*, so the sandbox-bypassed `codex` process it wrapped is never signalled and outlives the script. `cleanup_report_image_paths` then `rm -rf`s `$agent_directory` while that process may still be writing into it.
- `skills/do-work-toolbox/scripts/generate-report-image.sh:48` — `mv "$staged_output_path" "$output_path"` exits zero when `$output_path` is a directory, nesting the staged PNG inside it, so `publish_report_image` reports a publication that never happened.
- `skills/do-work-toolbox/scripts/install-last30days.sh:94` — same primitive. `:85-91` moves a pre-existing target into a private backup, then `:94` renames staging onto the now-absent path. A target recreated inside that window turns the rename into a nesting `mv` that exits zero, so `publication_complete=1` is set over a tree that was never published.

**Idioms to reuse (both already locked in by tests):**

- `skills/do-work/scripts/run-blocked-check.sh:43-73` — the verified-process-group pattern: establish isolation with job control, *verify* `pgid == pid` and `pgid != caller pgid` before signalling anything, TERM the group, grace-loop, escalate to KILL, then reap the leader. The verification is what makes it safe; the STOP/CONT wrapper there exists only because that script must fail closed (exit 125) rather than degrade, which is not this REQ's requirement.
- `skills/do-work-toolbox/actions/ai-report-reference.md:148-159` — REQ-204's post-rename nesting verification: after the `mv`, probe `<destination>/<stage-basename>`; if it exists the rename nested, so remove only the nested stage, leave the destination untouched, and exit nonzero.

**Interaction discovered during exploration:** in `install-last30days.sh` the nesting branch cannot simply fall into the existing rollback path. `restore_previous_tree` opens with an unconditional `rm -rf "$target_directory"`, which would destroy the very colliding tree the REQ requires be left byte-for-byte. The existing rollback-FAILED branch (`:41-43`) already models the right disposition — report that the prior tree remains at its backup path and clear `backup_parent` so cleanup neither restores over the collider nor deletes the backup.

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/scripts/generate-report-image.sh` (modify) — verified process-group ownership for both backends; post-rename nesting verification in `publish_report_image`
- `skills/do-work-toolbox/scripts/install-last30days.sh` (modify) — post-rename nesting verification at publication, with a non-destructive disposition for the collider
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — three fixture replays, one per instance

**Files I will NOT touch:** `skills/do-work-toolbox/actions/ai-report-reference.md` (REQ-204 already closed both boundaries there; REQ-221 owns moving that block), `skills/do-work/scripts/run-blocked-check.sh` (the source idiom, unchanged).

**Acceptance criteria (restated from REQ):**
- [x] A direct invocation of `generate-report-image.sh`, interrupted, leaves no launched process — backend or backend descendant — alive
- [x] A helper whose process-group isolation cannot be proved is signalled by bare PID only, never by group
- [x] `generate-report-image.sh` publishing onto a directory exits nonzero, leaves that directory byte-for-byte, and leaves no private stage
- [x] `install-last30days.sh` whose target reappears inside the check-then-rename window exits nonzero, leaves the reappeared tree byte-for-byte, and leaves no private stage
- [x] Every existing success path and exit status of both scripts is preserved; all pre-existing named cases in `_dev/tests/prescribed-shell-scripts-behavior.sh` still pass

## Pre-Flight

**Repository state:** OK — working tree clean outside `do-work/`
**Test baseline:** OK — `bash _dev/tests/prescribed-shell-scripts-behavior.sh` passing before implementation
**Dependencies:** OK

## Decisions

- **D-01 — DECIDE & STATE.** Extended the declared scope by two files to close restatement-sweep findings: `skills/do-work/docs/prescribed-shell-primitives.md` (the shipped executable-homes table understated both scripts, listing neither process-tree ownership nor verified exact publication, while the neighbouring `run-blocked-check.sh` row does name its process-group ownership) and `skills/do-work-toolbox/actions/ai-report-reference.md`'s **Generation helper** paragraph (it restated the helper's contract as file-cleanup-only and did not mention the process tree or the fail-closed directory case). Reasoning: both are prose restatements of behavior this REQ changed, and Step 7's Restatement Sweep is mandatory and explicitly covers files outside declared Scope. The prior session recorded the matching lesson — REQ-205's sweep missed exactly this shipped guide. Two sentences of prose sync is cheaper and safer than a follow-up REQ, and leaving them stale is the documented failure mode. `write_set` was extended with the Scope list, in that direction only.
- **D-02 — DECIDE & STATE.** In `generate-report-image.sh`, a nesting publication calls `exit 1` from inside `publish_report_image` rather than returning nonzero. Reasoning: the direct-backend branch reads `publish_report_image`'s failure as "try the agentic fallback", which would then run `: > "$staged_output_path"` against a variable the nesting branch has already cleared. No backend can repair a directory sitting at the output path, so the failure is terminal by nature. The `EXIT` trap still runs, so cleanup of the agent directory is unaffected. Reversible and local to one function.
- **D-03 — DECIDE & STATE.** `generate-report-image.sh` resolves `ps` by absolute path (`/bin/ps`, then `/usr/bin/ps`) instead of calling it bare as the prescribed batch block does, and treats its absence as "no group provable". Reasoning: this script is invoked by callers that hand it a minimal `PATH` (the suite's own agentic opt-in case sets `PATH="$agentic_bin"`), so a bare `ps` would silently fail there and degrade every invocation to bare-PID signalling. It matches `run-blocked-check.sh`, which already resolves `ps` this way. The degraded path is the documented fail-safe, so the absence case stays correct either way.

## Implementation Summary

**What was done:** Gave `generate-report-image.sh` ownership of the process tree it launches and gave both shipped publication helpers a post-rename nesting verification, closing the runtime-boundary defect class REQ-204's copy-paste sweep found alive outside its own scope. Both backends now launch under job control so each leads its own process group; the group is recorded only when it verifies as the backend's own and not the caller's, and an unverified group degrades to bare-PID signalling. Teardown TERMs, grace-waits, escalates to KILL, and reaps before the staged file and agent directory are removed. Both `mv` publications now probe for the nested-stage path that a directory destination produces, discard only their own nested stage, leave the destination byte-for-byte, and fail closed.

**Files changed:**
- `skills/do-work-toolbox/scripts/generate-report-image.sh` (modified) — verified-process-group ownership for the direct and agentic backends, teardown ordered before path removal, and post-rename nesting verification in `publish_report_image`
- `skills/do-work-toolbox/scripts/install-last30days.sh` (modified) — post-rename nesting verification at publication that clears `backup_parent` so cleanup neither restores over the collider nor deletes the prior tree
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified) — three fixture replays (process-tree ownership, output-is-a-directory, publication-collision); named-case count 41 → 44
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified) — executable-homes table rows for both scripts now name the boundaries they own (D-01)
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified) — Generation-helper paragraph restates the process-tree and fail-closed-publication contract (D-01)

## Qualification

Passed — 5 files verified in the diff, 5 requirements traced, no debug artifacts. Qualification was performed against git state rather than the builder's report: the diff was read line by line, and the red/green claim was re-derived independently (see Testing) rather than accepted.

One requirement is closed by construction rather than by a replay: *"A helper that cannot prove process-group isolation must signal only the bare PID, never a group."* `record_backend_process_group` clears the recorded group unless it is numeric, equals the backend PID, and differs from the caller's group, and every signal site branches on that emptiness. Forcing that branch needs a `ps`-less or job-control-less host, which the fixture cannot synthesize portably — the same reason the equivalent fail-safe in the prescribed batch block is untested. Noted rather than papered over.

## Testing

**Commands run:**
- `bash _dev/tests/prescribed-shell-scripts-behavior.sh` — exit 0 (44 named script cases)
- `bash _dev/tests/maintainer-verify.sh` — exit 0, zero FAIL lines

**Red-green validation** (traced to `## Red-Green Proof`; RED was re-derived by the orchestrator, not taken from the builder's report — both scripts were reverted with `git checkout --` and the suite re-run, leaving the new replays in place):

RED (exit 1) produced exactly the three new cases and no pre-existing regression:
- `FAIL: generate-report-image process-tree case left 1 backend process(es) or descendant(s) alive`
- `FAIL: generate-report-image output-is-a-directory case reported success`
- `FAIL: generate-report-image output-is-a-directory case left its staged image nested inside the occupying directory`
- `FAIL: generate-report-image output-is-a-directory case leaked private staging`
- `FAIL: install-last30days publication-collision case left its staging tree nested inside the reappeared target`
- `FAIL: install-last30days publication-collision case leaked private staging`
- `FAIL: install-last30days publication-collision case did not leave the prior tree recoverable at its backup path`

GREEN: with both fixes restored the suite exits 0. That RED set matches the captured proof on both halves — the interrupted direct invocation left a surviving descendant, and the reappearing target nested the staging tree. It also recorded damage the captured proof only implied: under the old code `install-last30days.sh` set `publication_complete=1`, so the EXIT trap deleted the backup and **destroyed the prior tree**; the case now pins that the prior tree stays recoverable.

**Existing tests updated:** none — all 41 pre-existing named cases pass unchanged.

## Review

**Requirements check — Pass.** Every stated requirement maps to a change: the two proven idioms were reused rather than reinvented (`run-blocked-check.sh`'s verified-process-group pattern, REQ-204's post-rename nesting probe); the unprovable-isolation case signals the bare PID only; both nesting publications exit nonzero, leave the destination byte-for-byte and leave no private stage; every existing success path and exit status is preserved (full suite green, including the pre-existing interruption case that still asserts TERM status 143); and a fixture replay was added per instance.

**Code review — Pass.** The teardown sits at the top of `cleanup_report_image_paths`, which already runs before the `rm`s, so ordering holds on both the EXIT path and the three `trap 'exit N'` signal paths without duplicating the reap into each trap. The `install-last30days.sh` nesting branch was checked against `cleanup_install_paths` line by line: with `publication_started=0`, `publication_complete=0`, `staging_directory=""`, `backup_parent=""` and `clone_directory=""`, every cleanup branch is a no-op — which is the point, because `restore_previous_tree` opens with an unconditional `rm -rf "$target_directory"` and would otherwise have destroyed the collider the requirement protects.

**Restatement sweep — 2 findings, both fixed in this REQ (D-01).** The shipped executable-homes table in `skills/do-work/docs/prescribed-shell-primitives.md` and the Generation-helper paragraph in `skills/do-work-toolbox/actions/ai-report-reference.md` both restated contracts this diff changed. Fixed rather than deferred; the prior session's lesson records a sweep that missed this same guide.

**Acceptance = Pass.** Both gates exit 0. Scope drift: two files added beyond the declaration, both restatement-sweep closures, declared as D-01 and mirrored into `write_set` — Minor.

**Overall: 95%.**

## Lessons Learned

**What worked:** Reverting the shipped scripts and re-running the suite is a cheap, decisive way to verify a builder's red/green claim instead of trusting it — it took one command and turned "the builder says RED" into a named FAIL list. It also surfaced damage the captured proof had not predicted (the destroyed backup tree), which became an assertion.

**What didn't:** Reading `publication_started=0` as sufficient to stop the rollback. `cleanup_install_paths` fires on `publication_started == 1` **or** a surviving backup, and `restore_previous_tree` opens with an unconditional `rm -rf "$target_directory"` — so the obvious fix would have satisfied "exit nonzero" while deleting the very tree the requirement says to leave byte-for-byte. A fail-closed branch has to be checked against the cleanup it hands control to, not just against its own postcondition.

**Worth knowing:** This defect class has now been fixed in four places (portfolio publication, the ai-report batch, and these two). Each fix was local because the guide states the boundary per script rather than as a condition; `skills/do-work/docs/prescribed-shell-primitives.md` names the nesting rule only inside its portfolio section. Also: `generate-report-image.sh` must resolve `ps` by absolute path — the suite's own agentic opt-in case runs it under `PATH="$agentic_bin"`, where a bare `ps` silently fails and degrades the fix to bare-PID signalling (D-03).

## Orientation

`generate-report-image.sh` and `install-last30days.sh` now own the boundaries they cross: an interrupted image generation leaves nothing it launched alive, and neither helper can report a publication that `mv` actually nested inside a directory. Lives in the do-work-toolbox shipped-scripts subsystem, alongside the prescribed-shell contract documented in `_dev/primes/prime-shell-commands.md`. No map change — this closes the last two known instances of a defect class already fixed twice elsewhere; no new module, data flow, or contract. Prime spot-check: `_dev/primes/prime-shell-commands.md`'s referenced paths all still resolve, and its pointer to `_dev/lessons/validated-runtime-boundaries.md` remains accurate.

## Discovered Tasks

- [normal] `skills/do-work/docs/prescribed-shell-primitives.md` states the `mv`/`ln`-onto-a-directory nesting rule only inside its **Portfolio summary publication** section, as a property of one script. The same boundary has now been closed in four scripts, each time locally, because the guide has no shared statement of the condition. Stating it once as a condition ("any publication whose destination could be occupied verifies the path it actually wrote") and having the per-script sections point at it would match CLAUDE.md's *State conditions, not lists* and would have made this REQ's two sites findable from the guide instead of from a review sweep. Not done here: it restructures a shipped guide beyond this REQ's fix-two-scripts remit.
