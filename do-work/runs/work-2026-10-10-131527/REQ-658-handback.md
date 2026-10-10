# REQ-658 hand-back (finalize --auto-manifest builds the mechanical manifest fields and preflights the tree)

- Branch: `worktree-agent-REQ-658-finalize-auto-manifest`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-658-finalize-auto-manifest`
- Base commit: `bd56c4b0`
- Commits: `09e6f186` `[REQ-658] add finalize --auto-manifest with --emit validation` (one commit)

## File manifest

- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest.go` (new): flag parsing, judged-input check, working-REQ lookup, mechanical fields, preflight (a) journal, (b) index, (c) dry core, `--emit` write, refusal helper.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_auto_manifest_test.go` (new): the four named Red-Green tests.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go` (modified): `prepareBoundJournal` keeps decode and binding checks, then calls the new shared core `prepareManifestJournal(repositoryRoot, manifest, manifestBytes, dryRun)`, which returns the required commit paths; dry mode adopts release payloads into `os.MkdirTemp` (removed on return), skips `writeJournal`.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_commands.go` (modified): `handleFinalize` routes to `handleAutoManifest` when `--auto-manifest` is present.
- `skills/do-work/tools/do-work-cli/internal/requeststate/state_plan.go` (modified): exported `DefaultWriterLabel(repositoryRoot)`; `planCheckpoint` calls it (same expression as before).
- `skills/do-work/actions/work.md` (modified): Step 8 item 5 names `finalize --auto-manifest REQ-NNN ... --emit <path>` and its inputs; Step 9 says the manifest is the Step 8 `--emit` file or a hand-built one. Same line counts.
- `skills/do-work/actions/work-reference.md` (modified): Commit & Metadata-Commit paragraph names the command, says a hand-built manifest stays valid.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modified): line 121 only, one clause on `--auto-manifest` and `--emit`.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Split `prepareBoundJournal` into decode plus a shared core with a dry mode (D-05). Auto mode: parse flags, refuse missing judged inputs (exit 1), preflight journal and index, find the one working REQ from the `DiscoverRepository` snapshot (`RequestsByID[id]` filtered by `TreeSection == "working"`), fill digests, `completed_at`, writer, then run the core dry on a draft (commit_paths = extra paths) to learn the required set, fill commit_paths, `validateManifest`, then the core again (dry for `--emit`, real otherwise followed by `advanceJournal`). Tests first.
- [x] **[APPLY]:** Implemented as planned within the eight write-boundary paths. RED recorded before any production edit.
- [x] **[UNIFY]:** `git diff bd56c4b0 --stat`:
  ```
   skills/do-work/actions/work-reference.md           |   2 +-
   skills/do-work/actions/work.md                     |   4 +-
   .../finalization/finalization_auto_manifest.go     | 248 +++++++++++++++++++++
   .../finalization_auto_manifest_test.go             | 155 +++++++++++++
   .../internal/finalization/finalization_commands.go |   4 +
   .../internal/finalization/finalization_prepare.go  |  80 ++++---
   .../internal/requeststate/state_plan.go            |  11 +-
   .../tools/do-work-cli/lessons-do-work-cli.md       |   2 +-
   8 files changed, 474 insertions(+), 32 deletions(-)
  ```
  Checks (from the worktree root):
  - `REQ-658-probe.sh`: exit 0, 2 s, "REQ-658 probe passed: 4 tests, both action files name the command." (re-run after commit: same)
  - `REQ-658-preflight-probe.sh`: exit 0, 13 s (`internal/finalization` ok 8.4 s, `internal/requeststate` ok, core-checks "contract probes passed").
  - `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/finalization/ ./internal/requeststate/`: ok (finalization 54.5 s, requeststate 7.7 s), wall 55 s. One run, no rerun needed.
  - `bash _dev/tests/contracts/core-checks.sh`: exit 0, 6.6 s.
  - `gofmt -l skills/do-work/tools/do-work-cli/internal`: empty. `go vet` on both packages: exit 0. `git diff --check`: exit 0.
  Files checked: all eight above; no debug output, no stray files; only write-boundary paths in the diff; nothing under `do-work/` staged or committed.

## Proof record

- RED (unchanged code, tests only), all four failed with `Outcome:"failure"`, `FINALIZATION-USAGE`, evidence `unknown finalize option "--auto-manifest"`:
  - TestFinalizeAutoManifestEmitsAManifestFinalizeAccepts
  - TestFinalizeAutoManifestRefusesAStagedPathBeforeAnyJournal
  - TestFinalizeAutoManifestRefusesWithoutAMessageFile
  - TestFinalizeAutoManifestRefusesAStaleReleaseVersionBeforeWriting
- GREEN:
  ```
  --- PASS: TestFinalizeAutoManifestEmitsAManifestFinalizeAccepts (1.20s)
  --- PASS: TestFinalizeAutoManifestRefusesAStagedPathBeforeAnyJournal (0.09s)
  --- PASS: TestFinalizeAutoManifestRefusesWithoutAMessageFile (0.04s)
  --- PASS: TestFinalizeAutoManifestRefusesAStaleReleaseVersionBeforeWriting (0.16s)
  ```
- Probe and preflight probe: both exit 0 (above).

## Decisions

- D-11 (DECIDE & STATE): Message-file and failure-error-file handling: trailing `\r` and `\n` are trimmed, everything else kept byte for byte. This matches what a hand-built manifest carries and git strips trailing newlines anyway.
- D-12 (DECIDE & STATE): The dry core still calls `journalLocations`, which creates the empty Git-private `.git/do-work-finalization/` directory, exactly as every `finalize --manifest` attempt (including refused ones) already does. No journal file and no payload directory are created; in dry mode the core blanks the Git-private payload path so its cleanup calls cannot reach it. A non-creating lookup would need a new helper for no observable gain.
- D-13 (DECIDE & STATE): The core returns the required commit paths alongside the "commit_paths omits planned lifecycle or release targets" refusal. The draft pass treats a non-nil set as "the planners passed" and ignores that one refusal; every other refusal returns a nil set and becomes the auto-mode refusal. The second pass with the filled manifest runs the full check, so nothing is skipped. No typed error was added.
- D-14 (DECIDE & STATE): The draft manifest carries `--extra-path` values as its commit_paths, so the release shipped-change guard (which reads commit_paths under `primary_commit`) sees the implementation paths in the draft pass as well.
- D-15 (DECIDE & STATE): The draft pass passes nil manifest bytes; their digest can never equal an existing journal's, so the dry core can never resume a journal (preflight (a) already refuses first).
- D-16 (DECIDE & STATE): Up-front judged-input checks cover presence and transition pairing only (`--transition` complete|fail, `--terminal-status` for complete only, `--failure-type` and `--failure-error-file` for fail only, `--message-file`, `--provenance`, `--implementation-hash` with `supplied_commit`). Value rules (terminal status words, failure type words, hash shape, primary_commit forbids a hash) stay in `validateManifest`, which runs on the filled manifest; its refusal is also OutcomeRefused with code `FINALIZATION-AUTO-INPUT`.
- D-17 (DECIDE & STATE): Refusal codes: `FINALIZATION-AUTO-INPUT` (missing or invalid judged input), `FINALIZATION-JOURNAL-UNFINISHED` (next argv `do-work-cli recover-finalization`), `FINALIZATION-INDEX-NOT-EMPTY` (staged paths in affected paths and evidence), `FINALIZATION-REQUEST-NOT-WORKING` (zero or several working files), `FINALIZATION-PREPARE-REFUSED` (dry core refusal, carries e.g. `RELEASE-PREIMAGE-STALE`). Success with `--emit` adds info finding `FINALIZATION-MANIFEST-EMITTED` with the path and commit_paths. Internal faults (git or discovery errors, emit write failure `FINALIZATION-EMIT`) stay `commandFailure`, exit 2, as today. The non-emit real core refusing after a clean dry pass keeps `FINALIZATION-PREPARE` failure, identical to `finalize --manifest`.
- D-18 (DECIDE & STATE): `--emit` overwrites an existing file at the path (mode 0o600). The action passes a fresh scratch path; refusing on existence would add a check nobody asked for.
- D-19 (DECIDE & STATE): Routing is `slices.Contains(arguments, "--auto-manifest")` in `handleFinalize`; the auto parser refuses `--manifest` as mutually exclusive (usage failure, exit 2). `parseFinalizeArguments` is unchanged.

## Discovered Tasks

- `prepareManifestJournal` in non-dry mode still runs `os.RemoveAll(payloadDirectory)` on early refusals before any payload was adopted, so a stray `REQ-N.payloads` folder with no journal is deleted by any refused `finalize --manifest` (pre-existing behaviour, unchanged) → report only
- `lifecycleadvance/queue_commands.go:424` `queueWriterLabel` is a near copy of the writer default that differs on purpose-or-drift: it substitutes `unknown-host` for an empty hostname, while `requeststate.DefaultWriterLabel` keeps `":"+root`. `lifecycleadvance` already imports `requeststate`, so the two could be unified; left untouched per the brief → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file).
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` bullets at lines 17, 24, 34, 77, 103 (families `commit-preflight-before-side-effect`, `guard-on-one-entry-path`) and the line 121 module bullet. Applied: every preflight refuses before the first write, and `--emit` validates through the same core the real path uses, so there is no second copy of any check.
- Primes: `prime-shell-commands.md`, `prime-action-files.md`, `prime-releases.md` (skimmed for the edited areas), `prime-do-work-cli.md`.

## Anti-bloat check

`git diff --stat` above. Additions the brief did not name:
- `autoManifestOptions` type, `parseAutoManifestArguments`, `judgedInputProblem`, `handleAutoManifest`, `readJudgedText`, `autoManifestRefusal`: these are the parsing, input-refusal, handler, file-read and refusal-result code that requirement 1 and the preflight require; each is used from more than one call site or keeps the handler readable. Reason recorded here because the brief listed the file, not the function names.
- Test helpers `autoManifestFixture`, `runAutoManifest`, `assertNothingWritten`, `findingText`: shared by the four named tests, no extra tests.
- No new flags beyond the brief's list, no new renderer, no `advance` change, no config.

## Proposed CHANGELOG entry (the integrator writes the version)

**finalize --auto-manifest builds the finalization manifest for you**

Sessions hand-built the finalization manifest about 171 times in 15 days and found the required commit paths by reading refusals. The CLI now fills every mechanical field and checks the tree before writing anything.

- `do-work-cli finalize --auto-manifest REQ-NNN` takes the judged inputs (`--transition`, `--terminal-status` or `--failure-type` with `--failure-error-file`, `--message-file`, `--provenance` with `--implementation-hash` for `supplied_commit`, optional `--release-manifest`) and fills the request path, both digests, `completed_at`, the writer label and the planner's required `commit_paths` plus each `--extra-path`.
- It refuses with exit 1, writing nothing, when a judged input is missing, an unfinished journal exists (it names `recover-finalization`), the index has staged paths (it lists them), or the release version is stale (`RELEASE-PREIMAGE-STALE`).
- `--emit <path>` runs every finalizer check without a journal and writes the manifest, which `advance --finalization-manifest` and the existing manifest mode accept unchanged. Hand-built manifests stay valid.
- `actions/work.md` Step 8 and 9 and the work-reference Commit procedure name the new command.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`

- [family: guard-on-one-entry-path] [REQ-658: a validate-only mode earns trust only when it runs the real planner core in a dry mode (temporary payload directory, no journal), not a parallel copy of the checks; let the core report the computed set alongside its "missing" refusal so a draft can learn it in one pass](../../do-work/archive/REQ-658-finalize-auto-manifest.md#lessons-learned)

## Integration seams and wall times

- `work.md`: only Step 8 item 5 (one line) and Step 9 (one line) changed; no reflow. `work-reference.md`: only the Commit & Metadata-Commit first paragraph (one line). `lessons-do-work-cli.md`: only line 121. REQ-659, REQ-660, REQ-689 and REQ-690 edit other lines of these files, so merges should be clean.
- `internal/finalization/` and `internal/requeststate/state_plan.go`: no sibling touches them.
- Test wall times: four new tests 1.5 s total; full `internal/finalization` 54.5 s and `internal/requeststate` 7.7 s under load.
