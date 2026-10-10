# Builder brief: REQ-658 (finalize --auto-manifest builds the mechanical fields of the finalization manifest and preflights the tree)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-658-finalize-auto-manifest
- Branch: worktree-agent-REQ-658-finalize-auto-manifest, created by the coordinator from main HEAD with `git worktree add -b` right before dispatch (do not create it; commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-658-finalize-auto-manifest.md. Read it fully: What, Why, Verified Facts, Detailed Requirements, Constraints, Assumptions, Builder Guidance, Red-Green Proof, Required Lessons Dropped for Budget, and the orchestrator's `## Exploration` (pre-dispatch decisions D-01 to D-10) and `## Scope`. Requirement 9 (the release) belongs to the integrator; you PROPOSE the changelog entry and a lesson bullet in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-145/input.md (sections Request item 1, What happened, Where the behaviour lives today Item 1, Proposed direction 1, Acceptance check).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-131527/REQ-658-handback.md
- Route B, tdd: true, impact-user-visible, effort-substantive, domain backend. Go change in `internal/finalization/` plus two action paragraphs.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, backend.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-shell-commands.md and prime-action-files.md (you edit prescribed prose in two action files), prime-releases.md (read; the release itself is the integrator's), and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/prime-do-work-cli.md. Required lesson: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-releases.md (whole file, small). The do-work-cli satellite was dropped for budget; read at least its bullets at lines 17, 24, 34, 77 and 103 of /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/lessons-do-work-cli.md (families `commit-preflight-before-side-effect` and `guard-on-one-entry-path`): both describe exactly the traps this command can fall into.

## The change (decided at pre-dispatch, do not reopen)
Paths are relative to your worktree root. Go code lives under `skills/do-work/tools/do-work-cli/` (below: `cli/`).

1. **Flag set** on `finalize` (dispatch in `cli/internal/finalization/finalization_commands.go:33-42` `handleFinalize`; today `:315-331` `parseFinalizeArguments` accepts only `--manifest`). `--auto-manifest REQ-N` is mutually exclusive with `--manifest`. Judged inputs: `--transition complete|fail`, `--terminal-status` (required for complete, refused for fail), `--message-file <path>`, `--provenance primary_commit` or `--provenance supplied_commit --implementation-hash <hash>`, optional `--release-manifest <path>`, and for fail `--failure-type` plus `--failure-error-file <path>` (D-09). Mechanical overrides: optional `--writer <label>`, repeatable `--extra-path <repo-relative path>`, optional `--emit <path>`. A missing judged input refuses with OutcomeRefused (exit 1) and a typed finding naming the flag; an unknown option stays `FINALIZATION-USAGE` failure (exit 2) as today (D-06). Never invent a judged value.
2. **Mechanical fields**, in a new file `cli/internal/finalization/finalization_auto_manifest.go`:
   - Request path: exactly one file for REQ-N under `do-work/working/` (use the `repositorymodel.DiscoverRepository` snapshot's request files filtered by tree section `working` and exact id, not a filename glob; `REQ-65` must not match `REQ-658`). Zero or several refuse (exit 1).
   - `expected_request_sha256` and `expected_checkpoint_sha256`: `digestBytes` of the live request file and `do-work/CHECKPOINT.md`, the same reads `finalization_prepare.go:71-78` checks.
   - `completed_at`: `time.Now().UTC().Truncate(time.Second)` formatted `2006-01-02T15:04:05Z`, the `now` source (`cli/internal/corehelpers/commands.go:73`). `release_at` equals `completed_at` when `--release-manifest` is given, else empty.
   - `writer_label`: `--writer`, else the existing default. Export the expression at `cli/internal/requeststate/state_plan.go:342-346` once as `requeststate.DefaultWriterLabel(repositoryRoot string) string` and call it there and in the new code (D-02). Do not touch `lifecycleadvance/queue_commands.go:424`.
   - `commit_message`: the message file's content (your call on trailing-newline handling; record it as a decision).
   - `commit_paths`: the planner's required set (lifecycle `statePlan.TargetPaths` plus every release postimage path, `finalization_prepare.go:155-171`) plus every `--extra-path`, normalized and sorted with the existing `normalizeRepositoryPaths`.
3. **One planner path, with a dry mode** (D-05). Split `prepareBoundJournal` (`finalization_prepare.go:39-188`) into the file decode plus a shared core that takes a decoded `Manifest` and its bytes. The core gets a dry mode that: adopts release payloads into an `os.MkdirTemp` directory removed on return (never the Git-private payload directory), skips `writeJournal`, and returns the required commit paths. Auto mode calls the core twice: once in dry mode with a draft manifest to learn the required set (the draft cannot pass `validateManifest` yet because `commit_paths` is empty, so validate after filling), then, with the filled manifest, through `validateManifest` and the same core, dry for `--emit`, real otherwise. No second copy of any check. `finalize --manifest` and `FinalizeBound` keep their exact current behaviour and outcomes.
4. **Preflight before any write**, in this order, each a refusal with exit 1 that leaves no journal, no payload directory and no emitted file: (a) an unfinished journal for REQ-N exists: refuse and name `recover-finalization` in the finding's next argv (D-08); (b) the index is not empty (`git diff --cached --name-only`; the finding lists the staged paths) (D-03); (c) the dry core refuses: return its reason, which includes the release planner's `RELEASE-PREIMAGE-STALE` when the project version no longer matches the release manifest's old version (D-04).
5. **`--emit <path>`** writes `json.Marshal` of the filled `Manifest` (resolve the path with `containedOrAbsolute`, write mode 0o600) and stops with OutcomeSuccess plus one info finding carrying the emitted path and the commit_paths (D-07). Without `--emit`, marshal the same bytes and run the real core and `advanceJournal`, exactly as `handleFinalize` does after `prepareJournal`.
6. **Output**: no new renderer. The existing text rendering (`cli/internal/resultmodel/result_model.go:1023-1110`) already prints outcome, each finding's code and evidence, and each finalization record's phase, archive path and commit hashes (D-07). Exit codes stay `result_model.go:648-662`.
7. **Action prose** (requirement 8). `skills/do-work/actions/work.md` Step 8 item 5 (`:462`) and Step 9 (`:467`), and `skills/do-work/actions/work-reference.md` Commit & Metadata-Commit paragraph (`:719`): name `finalize --auto-manifest REQ-NNN ... --emit <path>` as the way to build the manifest from the judged inputs (the judged fields stay the action's), then pass that file to the same `advance` continuation; say the hand-built manifest stays valid. Hard constraints from `_dev/tests/contracts/core-checks.sh:752-788`, which you must keep green:
   - keep verbatim: "author exactly one strict manifest" (Step 8); "exact `advance` continuation", "ordered `finalizations`", "`phase: cleanup_complete`", "empty `blocked_paths`/`reason_codes`" (Step 9); "Pass that one manifest to the current `advance` continuation." (work-reference Commit section);
   - Step 9 (heading through the line before `### Step 10`) stays at 8 lines or fewer, so add words inside the existing paragraph, not new lines;
   - in those four sections no line may match `do-work-cli\.sh.*finalize` (this also matches "finalization" later on any line that names the `.sh` wrapper), `finalize --manifest`, `git add`, `git commit` or `record-commit-hash`. So name the command as `finalize --auto-manifest`, never with the wrapper path in front, and describe the hand path without spelling "finalize --manifest".
8. **Module map**: in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:121` (the `internal/finalization/` bullet), add one clause saying `finalize --auto-manifest REQ-NNN` builds the mechanical manifest fields and `--emit` validates without a journal. Edit only that line.

## Anti-bloat (YAGNI); the maintainer asked for this to be watched
- Smallest change that delivers the named behaviour. Prefer reusing the existing planner, validator and renderer over adding.
- No new helpers, flags, options, config, abstractions or files beyond the ones named above. No `advance` auto mode (D-01), no new output renderer (D-07), no change to the `advance` missing-evidence hint (D-10).
- Tests pin only the named failures: exactly the four tests below, one RED case each. No decorative or "while we're here" tests.
- Do not fix adjacent things you notice; record them as discovered tasks instead.
In the hand-back, paste `git diff --stat` and list every function, constant, option, file or test you added that this brief did not name (expected: none). Each one gets a one-line reason, or you remove it.

## Write boundary
Exactly these paths in your worktree:
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest.go (new)
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest_test.go (new)
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go
- skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go
- skills/do-work/actions/work.md
- skills/do-work/actions/work-reference.md
- skills/do-work/tools/do-work-cli/lessons-do-work-cli.md
Anything else: stop and say so in the hand-back.

## Hard rules
- Every commit subject on your branch starts with `[REQ-658]` (for example `[REQ-658] add finalize --auto-manifest`). The board credits commits by that prefix only.
- Never create, edit, stage or commit any path under `do-work/`, in the worktree or anywhere else, except the one hand-back file above (which you write but never stage or commit). The integrator's queue guard refuses a branch that touches `do-work/`.
- Never merge, rebase, push, or check out another branch. Never run `recover`, `advance`, any `do-work-cli` lifecycle command against a real repository (your Go tests use temp fixture repositories), `just do-work-update`, or the repository gate (`_dev/tests/maintainer-verify.sh`); the integrator owns the gate.
- Do not touch `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` or any version mirror: the release belongs to the integrator.
- Scratch files and fixtures go in `mktemp -d` or `t.TempDir()` directories outside both trees and are never committed. Clean up anything you start.
- Go: run focused `go test -run` and `go vet` on the touched packages only, never the repository gate and never a whole-repository test run.
- Load is high while eleven builders run at once: if a wall-time budget or a test fails once under load, rerun it once and record both runs.

## Integration seam
Integration order puts REQ-658 third, after REQ-660 (worktree lifecycle command) and REQ-659 (frontmatter set and req append-section). Shared files:
- `skills/do-work/actions/work.md`: REQ-659 edits the stamp lines near Step 6 and Step 7, REQ-689 (run --coordinate) edits coordinator prose. You edit only Step 8 item 5 and Step 9.
- `skills/do-work/actions/work-reference.md`: REQ-659 and REQ-690 (status action) edit other sections. You edit only the Commit & Metadata-Commit paragraph.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: REQ-659 and REQ-660 add bullets. You edit only line 121.
- `internal/finalization/` and `internal/requeststate/state_plan.go`: no sibling touches them.
Keep your edits local to those spots (do not reflow neighbouring prose or re-wrap paragraphs) so the serial merges stay clean.

## Proof to run and record (from the REQ's Red-Green Proof)
All four tests go in `finalization_auto_manifest_test.go` and reuse the package fixtures: `newFinalizationRepository`, `writeFinalizationFile`, `runFinalizationGit`, `digestFile` (`finalization_commands_test.go:231-277`) and `seedPlannedReleaseFinalization` (`finalization_req499_test.go:638-675`, a claimed REQ-760 with a VERSION 1.0.0 to 1.0.1 release manifest; read the release manifest path out of the manifest file it returns). Call the handler in-process (`Handlers()[CommandFinalize]`), as the existing tests do.
1. `TestFinalizeAutoManifestEmitsAManifestFinalizeAccepts`: on the REQ-760 fixture, run `--auto-manifest REQ-760 --transition complete --terminal-status completed --message-file <f> --provenance primary_commit --release-manifest <r> --emit <m.json>`. Assert success, no journal exists, the emitted `commit_paths` equal the planner's set (request path, `do-work/archive/REQ-760.md`, `do-work/CHECKPOINT.md`, `VERSION`, `CHANGELOG.md`) and `release_at` equals `completed_at`. Then run `finalize --manifest <m.json>` and assert success with one finalization record at `cleanup_complete`.
2. `TestFinalizeAutoManifestRefusesAStagedPathBeforeAnyJournal`: stage one unrelated file; assert OutcomeRefused, `resultmodel.ExitCode` 1, the finding names the staged path, and neither a journal (`journalLocations`) nor the emit file exists.
3. `TestFinalizeAutoManifestRefusesWithoutAMessageFile`: omit `--message-file`; assert OutcomeRefused, exit 1, the finding names the flag, no emit file.
4. `TestFinalizeAutoManifestRefusesAStaleReleaseVersionBeforeWriting`: on the REQ-760 fixture, rewrite and commit `VERSION` as `1.0.1`; assert OutcomeRefused, exit 1, evidence carries `RELEASE-PREIMAGE-STALE`, and no journal, no payload directory and no emit file exist.
RED: write the four tests first and run them against the unchanged code; record the failure text (expected: "unknown finalize option \"--auto-manifest\"" for each). GREEN: the same run after the change.

## Verify before hand-back (from the worktree root, wall times recorded)
- `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-131527/REQ-658-probe.sh` (it reads relative paths, so it checks your tree): exit 0, four PASS lines, both action files name the command. It fails at the base commit by design.
- `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-10-131527/REQ-658-preflight-probe.sh`: exit 0 (the existing finalize-path tests, requeststate writer tests and `_dev/tests/contracts/core-checks.sh` stay green).
- `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/finalization/ ./internal/requeststate/` (about 70 s under load): ok.
- `gofmt -l skills/do-work/tools/do-work-cli/internal` prints nothing; `go vet -C skills/do-work/tools/do-work-cli ./internal/finalization/ ./internal/requeststate/` clean.
- `git diff <base> --stat` shows only the eight write-boundary paths; `git diff --check` clean.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path, base commit (`git -C <worktree> rev-parse --short HEAD` before your first edit) and every commit hash (one commit preferred).
- File manifest: each file with (new)/(modified)/(deleted) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists `git diff --stat`, every check's exit code and wall time, and each file checked.
- Proof record: the RED run (failure text per test), the GREEN run (four PASS lines), the probe and preflight-probe runs.
- `## Decisions` (D-11 onwards; D-01 to D-10 are in the REQ's Exploration; each DECIDE & STATE, or ESCALATE with Value and Risk). Include your message-file newline handling.
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellites and families).
- Anti-bloat check (see above): the `git diff --stat` and the list of anything added that this brief did not name.
- A proposed CHANGELOG entry: a descriptive title that says what shipped, in plain words, one or two sentences on why it matters, then specific bullets. The integrator writes it with the version.
- A proposed lesson bullet, naming the satellite you think it belongs in, in that file's bullet shape (`- [family: <slug>] [REQ-658: <one-line lesson>](<relative archive link; the integrator fixes the path>)`). Say "none" if nothing was learned that a later builder would need.
- Integration seams (see the note above) and test wall times.
