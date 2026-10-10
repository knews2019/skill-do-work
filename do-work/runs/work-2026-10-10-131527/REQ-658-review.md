## Review: REQ-658

**Approve** — `finalize --auto-manifest` works, and its `--emit` file passes both `finalize --manifest` and `advance --finalization-manifest` unchanged. One risk is worth fixing before release: without `--emit`, the command finalizes outside the `advance` phase gate (F1).
Route B | merge `457608cc` (range `84d56197..457608cc`)

### What's built
- `do-work-cli finalize --auto-manifest REQ-NNN` takes the judged inputs and fills the mechanical fields: request path, both digests, `completed_at`, `release_at`, writer, and the planner's required `commit_paths` plus each `--extra-path`. It refuses with exit 1 on a missing judged input, an unfinished journal, a staged index, or a planner refusal such as `RELEASE-PREIMAGE-STALE`.
- `--emit <path>` runs the shared planner core in dry mode and writes the manifest. Without `--emit`, the command finalizes directly through the same journal path as `finalize --manifest`.
- `work.md` Step 8 item 5, Step 9 and the `work-reference.md` Commit paragraph name the command and still describe a hand-built manifest as valid.

### Decisions / risks for you
- F1: keep or drop the non-emit mode. The REQ's requirements 4 and 5 imply it. The work loop never uses it (D-01), and it skips the `advance` phase gate. Dropping it, or refusing without `--emit`, removes about 6 lines and the risk. Keeping it matches the REQ text but leaves a one-flag path to archive an unreviewed REQ.

### Checks asked for by the orchestrator

1. **Shared core split (`finalization_prepare.go:50-222`).** The non-dry path is the same as before. Every change in non-dry mode is a fourth return value. The order of checks and the journal write are unchanged. `prepareBoundJournal` keeps decode and binding checks, so `FinalizeBound` (advance) behaves the same. Dry mode cannot reach the Git-private payload directory. `payloadDirectory` is set to `""` (`:73-75`), and `os.RemoveAll("")` does nothing. Release payloads go to an `os.MkdirTemp` directory that is removed on return (`:139-145`). Dry mode returns before the journal is built (`:201-203`). It still calls `journalLocations`, which creates the empty `.git/do-work-finalization/` directory (D-12, see F4). D-13 hides nothing else. The only error returned with a non-nil required set is the "commit_paths omits" refusal (`:197-199`). Every earlier refusal returns `nil`, and dry success returns before any later error. One unreachable edge: the resume branch (`:85`) returns a nil set with a nil error, and the caller would then panic on `err.Error()` (`finalization_auto_manifest.go:187`). Reaching it needs a journal whose digest equals SHA-256 of empty bytes, which cannot happen in practice. No finding.
2. **`finalization_auto_manifest.go`.** No judged input is invented. `judgedInputProblem` (`:75-97`) refuses a missing transition, terminal status, failure pair, message file, provenance, or `supplied_commit` hash. Every refusal comes before any write except the empty journal directory (F4). Advance path proven in a scratch copy of the module (`.../scratchpad/int-658/review/cli`, reviewer-only test files). Three cases each emitted a manifest and then ran `advance REQ --request-path ... --finalization-manifest <emitted>` to exit 0, `cleanup_complete`, and a clean tree: primary_commit with release and `--extra-path implementation.txt`, supplied_commit, and already-green. The emitted `commit_paths` equalled the existing hand-built fixture's set in all three. Non-emit mode also finalized to `cleanup_complete`. Seven missing or invalid input cases all exited 1, and the journal directory stayed empty.
3. **Resolver seam.** This is a second resolver (F2). The functional impact is nil because the core's `BuildPlan` calls `ResolveTarget` again. Only the refusal code differs.
4. **Anti-bloat count.** Production: 1 new file (248 lines), 1 options type, 5 functions, 1 newly exported function (`requeststate.DefaultWriterLabel`), and 12 flags. Nine flags are named by the REQ. `--failure-type` and `--failure-error-file` come from a capture assumption. `--writer` comes from pre-dispatch D-02 and the "writer label comes from an input flag" assumption. So no flag is unasked. 7 new result codes: 5 refusals, `FINALIZATION-MANIFEST-EMITTED` and `FINALIZATION-EMIT`. The REQ named none of them and nothing reads them (F5). Preflights (a) journal and (b) index repeat checks the core already makes. They exist to give exit 1 plus named paths instead of exit 2, which is fair. The non-emit mode is REQ-implied but unused by the loop (F1). Tests: 4 tests and 4 helpers, none decorative. Each test names the failure it pins. Total: small. No Constraint is broken.
5. **Prose.** Accurate to the code: the inputs, the filled fields, "runs the finalizer's checks without a journal", the staged-index and stale-version refusals, and "a hand-built manifest stays valid" (`work.md:462`, `:467`, `work-reference.md:719`). Step 8 item 5 grew from 486 to 1080 characters and lists the mechanical fields twice (F3). Telling the action to pass `--extra-path` per implementation path is harmless under `supplied_commit`. The reviewer checked that case: advance succeeds and the clean path is listed.
6. **Restatement sweep.** See the line in the Review block below.

### Findings

**Important:**
- F1. Without `--emit`, `finalize --auto-manifest` finalizes outside the `advance` phase gate (`finalization_auto_manifest.go:200-206`). Reproduced on a Triage-only REQ. `advance --finalization-manifest <emitted>` refused (exit 2, `FINALIZATION-PREPARE`, the "only permits transition fail" gate). The same command without `--emit` exited 0 and archived the REQ as completed. `finalize --manifest` already had this power, but `core-checks.sh` keeps it out of action prose. `work.md:462` now names `finalize --auto-manifest`, so dropping one flag lands on the ungated path. This is the bare-subcommand-executes trap class. Fix: refuse without `--emit` (the loop uses emit then advance, D-01), or bind non-emit to the same phase check `FinalizeBound` gets from `advance`. — impact-user-visible → report only

**Minor:**
- F2. Second REQ-id resolver: `finalization_auto_manifest.go:148-156` filters `RequestsByID[id]` by `TreeSection == "working"` instead of `requeststate.ResolveTarget` (`state_targets.go:11`). It skips the identity check and picks the working copy when the same id also exists elsewhere, where `ResolveTarget` refuses `REQUEST-AMBIGUOUS`. The outcome converges because the core's `BuildPlan` re-resolves and refuses. Only the code differs (`FINALIZATION-PREPARE-REFUSED` wrapping `REQUEST-AMBIGUOUS` instead of `FINALIZATION-REQUEST-NOT-WORKING`). Fix: call `ResolveTarget(snapshot, id, "")`, then check `TreeSection == "working"`, as `corehelpers/request_writers.go:196-204` does. — impact-negligible → report only
- F3. The prose lists the mechanical fields twice and still frames them as authored. `work.md:462` says the action authors "writer and timestamp, exact request/checkpoint preimages", then says the command fills them. `work-reference.md:719` says "The action judges and authors" the digests, timestamp and writer, then says the command fills them. The judged-versus-mechanical split, which is the point of UR-145, is never stated. Fix: one sentence naming only the judged inputs plus "the command fills the rest", keeping the `core-checks.sh` tokens. — impact-negligible → report only
- F4. Preflight (a) calls `journalLocations` (MkdirAll, `finalization_journal.go:110`) before the index check. So a staged-index refusal that prints "nothing was written" creates `.git/do-work-finalization/` (reproduced in a scratch repo). `finalize --manifest` checks the index first (`finalization_prepare.go:66-68`) and creates nothing on that refusal, so D-12's "exactly as every finalize --manifest attempt" is wrong for this case. Fix: run preflight (b) before (a). — impact-negligible → report only
- F5. Refusal-code boundaries do not match D-16. A bad `--terminal-status` (for example `done`) or a bad `--extra-path` refuses `FINALIZATION-PREPARE-REFUSED`, not `FINALIZATION-AUTO-INPUT`, because the draft core runs before `validateManifest` (`finalization_auto_manifest.go:185-194`). Exit 1 either way. The 7 new codes have no reader. Fix: correct D-16's wording, or fold the five refusal codes into fewer. — impact-negligible → report only
- F6. Test gap: no repository test pins the main consumer path, `advance --finalization-manifest <emitted file>` (REQ Assumptions, acceptance criterion 4). Non-emit mode, `supplied_commit`, the `fail` transition, and the unfinished-journal and not-working refusals are also untested. The reviewer proved the advance path in a scratch copy only. Fix: add one `lifecycleadvance` test that emits under `supplied_commit` and then advances (about 1 s). — impact-negligible → report only
- F7. Inherited from REQ-660 (builder worktree lifecycle command), its F3: the dirt readers still count owned worktree links as uncommitted changes. `corehelpers/handoff.go:87` and `cleanup/cleanup_git.go:345` are confirmed unchanged at `457608cc`. queue-kanban `verify.go` and the Crash Recovery prose are as REQ-660 recorded them. This merge did not touch them. Fix: as REQ-660 F3 recorded (teach those readers `do-work/worktree-links`). — impact-user-visible → report only

### Requirements Checklist
- [x] 1. `--auto-manifest REQ-N` with judged inputs (`--transition`, `--terminal-status`, `--message-file`, `--provenance` + `--implementation-hash`, optional `--release-manifest`, failure fields) — delivered
- [x] 2. Mechanical fields filled (path, digests, `completed_at` from the `now` source, writer, `commit_paths` = planner set + `--extra-path`) — delivered; emitted set equals hand-built fixture sets in 4 fixtures
- [x] 3. Preflight before writing (version via `RELEASE-PREIMAGE-STALE`, empty index) — delivered; "nothing staged outside commit_paths" subsumed by the empty-index rule (D-03); empty journal directory caveat (F4)
- [x] 4. `--emit` writes and stops — delivered
- [x] 5. Outcome, phase, archive path, commit hash, findings — delivered through the existing text renderer (D-07)
- [x] 6. Exit codes unchanged — delivered (refusals 1, usage/internal 2)
- [x] 7. Judged fields never invented — delivered
- [x] 8. Action prose names the command, hand path valid — delivered (F3 on wording)
- [x] 9. Release — left to the integrator at finalization

### Acceptance Testing

**Result: Pass** (implementation and integration stages)
- `go test -count=1 -run 'AutoManifest|Finalize' ./internal/finalization/`: ok (7.0 s).
- Scratch copy of the module: emit then `advance --finalization-manifest` succeeded for 3 fixtures. Non-emit finalized. 7 refusal cases exited 1. `supplied_commit` with an implementation `--extra-path` succeeded. A Triage-only REQ was archived by non-emit while advance refused it (F1).
- Real binary in a scratch repo: staged-index refusal text output, exit 1, empty `.git/do-work-finalization/` created (F4).
- Integrator's full gate at `457608cc`: exit 0 (not re-run, per brief).

### Suggested Additional Testing
- Live acceptance: unassessed. In the next real work-loop finalization, run Step 8 `--emit` and then Step 9 `advance`, and confirm one `cleanup_complete` record.
- Deployment: unassessed. After `update-suite` in a consumer project, run `finalize --auto-manifest` once with `--emit` on a claimed REQ.
- Edge case: the same REQ id in both `working/` and `archive/` (F2 code path).
- Edge case: the `fail` transition end to end (`--failure-type`, `--failure-error-file`) through advance.
- Regression: a session that still hand-builds a manifest through advance (unchanged path).

### Scores (on the record — not the headline)

**Overall: 88%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 93% | All 9 delivered; requirement 3 subsumed by an existing rule |
| Code Quality | 85% | Clean core split; F1, F2, F4, F5 |
| Test Adequacy | 78% | 4 focused tests; main consumer path untested in repo (F6) |
| Scope | 96% | Exactly the 8 declared paths |
| Risk | Low | F1 needs an agent to drop `--emit` |
| Acceptance | Pass | Implementation + integration, incl. scratch advance proof |

### Follow-ups created
- None (7 findings report only)

## Review

**Overall: 88%** | 2026-10-10T15:18:24Z

| Dimension | Score |
|-----------|-------|
| Requirements | 93% |
| Code Quality | 85% |
| Test Adequacy | 78% |
| Scope | 96% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `finalize --auto-manifest` without `--emit` finalizes outside the `advance` phase gate. Reproduced: a Triage-only REQ was archived as completed while advance refused the same manifest (`finalization_auto_manifest.go:200-206`, named at `work.md:462`). Fix: refuse without `--emit`, or apply the advance phase check. — impact-user-visible → report only

**Minor findings:**
- F2 second REQ-id resolver (`finalization_auto_manifest.go:148-156`) instead of `requeststate.ResolveTarget` + a `TreeSection` check; the outcome converges, only the refusal code differs — impact-negligible → report only
- F3 `work.md:462` and `work-reference.md:719` list the mechanical fields as authored and then as filled; Step 8 item 5 grew from 486 to 1080 characters; the judged/mechanical split is never stated — impact-negligible → report only
- F4 preflight (a) creates `.git/do-work-finalization/` before the index refusal that says "nothing was written"; `finalize --manifest` does not; run the index check first — impact-negligible → report only
- F5 a bad `--terminal-status` or `--extra-path` refuses `FINALIZATION-PREPARE-REFUSED`, not `FINALIZATION-AUTO-INPUT` as D-16 says; the 7 new codes have no reader — impact-negligible → report only
- F6 no repository test for `advance --finalization-manifest <emitted>`, non-emit mode, `supplied_commit`, `fail`, or the journal and not-working refusals; add one lifecycleadvance emit-then-advance test — impact-negligible → report only
- F7 inherited from REQ-660 (worktree lifecycle command) F3: `handoff.go:87`, `cleanup_git.go:345`, queue-kanban `verify.go` and Crash Recovery still read owned worktree links as dirt; untouched by this merge — impact-user-visible → report only

**Acceptance:** Pass — implementation and integration stages: focused finalization tests ok; the emitted manifest passes `finalize --manifest` and (scratch copy) `advance --finalization-manifest` for primary+release, supplied and already-green fixtures; 7 refusals exit 1; deployment and live acceptance unassessed
**Restatement sweep:** redefined (E1) how the finalization manifest is built (`finalize --auto-manifest ... --emit`), (E2) which manifest fields are mechanical vs judged, (E3) `work.md` Step 8 item 5 / Step 9 and `work-reference.md` Commit wording. Stale: `work.md:462` and `work-reference.md:719` first clauses still frame mechanical fields as action-authored (F3). Consistent: `work.md:256`, `:467`, `:500`; `work-reference.md:352`, `:388`, `:715`; `lessons-do-work-cli.md:127`; `advance_commands.go:264-265` placeholder `<action-authored-finalization-manifest>` (kept by D-10, still valid); no shipped text describes finding `commit_paths` by refusal; `docs/` and `skills/do-work/docs/` carry no `finalize` usage. Inherited: REQ-659 (frontmatter set and req append-section) F4 `work-reference.md:61` and F5 `:73` fixed and still fixed; REQ-660 (worktree lifecycle command) F1 and F2 (`fan-out-reference.md:78`) fixed and still fixed, F3 dirt readers still stale (F7); this merge made none of them stale again (it touched only `work.md:462`, `:467`, `work-reference.md:719`, `lessons-do-work-cli.md:127`)
**Suggested testing:** 5 items
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*
