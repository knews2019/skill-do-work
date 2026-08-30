# REQ-407 Hand-Back

## Branch

`worktree-agent-REQ-407-migrate-install-update-bootstrap-to-go`, branched from `f38537e`. Twelve commits, oldest first:

| Hash | Subject |
|---|---|
| `1611116a8cdb2feb3cbcf449513b4420aff0e3d1` | port the managed-section replacer from Python to Go |
| `e494febed514eba020203c89b88b0472357a8e54` | replace the jq and Python settings reconcilers with order-preserving Go |
| `a1a08fd9bb66f4d4e35129d70004fc6104913803` | port the suite manifest validator to Go |
| `f27f564dbc928f70f1b75f6e1890e294fb9522da` | port the two-route upstream archive fetch to Go |
| `090a993118c2faec9647a1f845f3f74ae52b18e1` | close the five gaps folded in from REQ-406 |
| `73283d03a4d8cebf0165927967f4fb50a52905dd` | move the install and update transactions into Go and register five commands |
| `916d27061ab84c7b047ebe20b56d0b5c73dd4c71` | turn the five public shell entry points into compatibility launchers |
| `06d5d22a2558a178db4a2fb60fbbd3b58bdb8158` | retarget the installer behaviour suite at the Go implementation |
| `b61aee4cf40b9643d12f62d55feb998b2c93d87b` | retarget the update behaviour suite at the Go implementation |
| `7685de47a7c6dff628089be4a97bd8235e016dc7` | stop the install signal handler racing the writes it must undo |
| `2a4b28078d94b01b173bee26f393f9c9eae90a78` | retarget contract-regressions at the Go sources that now own the behaviour |
| `bda2f2b121d25cfb8402916e47c0f7d58918d47a` | document the Go 1.26.1+ prerequisite and correct the shipped claims |

Nothing pushed, nothing merged. Working tree clean. `git diff --stat f38537e..HEAD` touches exactly the 39 files in the Scope declaration and no others: 6662 insertions, 2636 deletions.

## Implementation Summary

### New Go packages

- **CREATED** `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go` — byte-exact port of the 344 lines of Python `replace-text-section.sh` carried. Line splitting on LF, CRLF and a bare CR; marker bodies stripped of a trailing SET of CR/LF bytes; span located as exact byte offsets with the end marker's own terminator inside it; four-case append separator; conditional write; atomic temp-in-target-directory plus rename with best-effort directory fsync; target mode preserved on replace, template mode on create.
- **CREATED** `skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions.go` — port of the Just recipe/alias scanner: five literal forms in priority order, backtick-adjacency rule for triple backticks, backslash escaping honoured only by the cooked forms, multiline state carried across lines, BOM stripped from line zero's classification view only.
- **CREATED** `.../internal/managedsection/managed_section_test.go`, `just_definitions_test.go` — 40 cases covering lone-CR splitting, CRLF, NUL, BOM asymmetry, malformed and reversed markers, symlink and dangling-symlink refusal, mode preservation, byte-idempotence, custom markers, collisions, and every multiline-literal boundary.
- **CREATED** `.../internal/settingshooks/settings_hooks.go` — order-preserving JSON settings composition replacing both the jq program and the embedded Python. Values travel as an `orderedObject` that remembers decoded key order; retired pipeline-guard command objects removed, wrappers emptied by that removal dropped, wrappers that lost nothing preserved, `Stop` deleted when it empties; fragment events appended once each in sorted order; two-space indent, HTML escaping off, trailing newline.
- **CREATED** `.../internal/settingshooks/settings_hooks_test.go` — key-order preservation, append-unique, guard removal, custom-hook preservation, empty-wrapper preservation, raw UTF-8, HTML-significant characters, and five malformed-input refusals.
- **CREATED** `.../internal/suitemanifest/suite_manifest.go` — port of `validate-suite-manifest.sh`, fail-fast. Manifest rows split on LF only so a CR stays content.
- **CREATED** `.../internal/suitemanifest/suite_manifest_test.go` — 16 rejection cases plus a non-mutation check and the Current-version parser's exact semantics.
- **CREATED** `.../internal/archivefetch/archive_fetch.go` — two-route fetch. HTTP through `bash scripts/atomic-download.sh` located by probing the two mirror-relative depths; git route via shallow clone plus `git archive --prefix`, staged beside the target and renamed only after it reads back.
- **CREATED** `.../internal/archivefetch/archive_fetch_test.go` — branch derivation for five URL shapes, HTTP win, unreadable-body fallthrough, export-ignore, preserved-target-on-total-failure with no scratch, underivable repository URL, and the `DO_WORK_UPSTREAM_URL` escape hatch.
- **CREATED** `.../internal/suiteinstall/install_transaction.go` — the whole install transaction: Git-worktree-root gate, archive resolution, manifest validation, four-module plan, symlink/non-directory/escape guards, source symlink guard, Justfile discovery preserving the real directory-entry spelling, candidate construction and validation for Justfile, CLAUDE.md and settings, full review diff, dirty-managed warning, one confirmation, backups including the Git index, module unstaging, writes, post-write byte and version verification, and exact recovery of every managed path plus the index.
- **CREATED** `.../internal/suiteinstall/update_transaction.go` — fetch, extract, validate, compare versions, delegate to the install transaction in-process, verify the installed version. Shared-install refusal and Git-worktree-root gate preserved.
- **CREATED** `.../internal/suiteinstall/suite_commands.go` — the five handlers, their long-form-only option parsing, and the single place a transaction outcome becomes a `CommandResult` for this family.
- **CREATED** `.../internal/suiteinstall/install_transaction_test.go`, `update_transaction_test.go`, `suite_commands_test.go` — fresh install, reinstall with mode and byte preservation plus idempotence, cancellation, post-write-failure recovery through a flaky `just`, collision refusal before confirmation, non-Git and subdirectory roots, symlinked destination, up-to-date/older/cancelled/successful update, shared-install refusal, version ordering, each command's argv contract and exit code, and the stdout-result / stderr-narration split.

### Modified Go

- **MODIFIED** `.../cmd/do-work-cli/main.go` — registers the five handlers with `os.Stderr` narration and `os.Stdin` confirmation. One import and one map entry per command family.
- **MODIFIED** `.../internal/commandruntime/command_runtime.go` — `parseGlobalOptions` became a method and scans ahead for the command token; `usageFinding` takes the command name; added `registeredCommands`, `nameableCommand`, `scanForCommandToken`. No finding emits `<command>` any more, and unknown/missing-command findings list what is registered.
- **MODIFIED** `.../internal/gittransaction/transaction_findings.go` — `BuildCommandResult` takes the command name and sets `Command`; every remediation template threads it through.
- **MODIFIED** `.../internal/gittransaction/git_transaction.go` — the success path now consults `state.existed`: a changed path that did not exist before and was not recorded as created is an unrecorded mutation and rolls back.
- **MODIFIED** `.../internal/gittransaction/git_transaction_test.go`, `transaction_findings_test.go`, `.../internal/commandruntime/command_runtime_test.go`, `.../internal/resultmodel/result_model_test.go` — new coverage described under Testing.

### Shell entry points, now launchers

- **MODIFIED** `tools/install-do-work-suite.sh` and `skills/do-work/tools/install-do-work-suite.sh` (byte-identical mirrors) — 621 lines to 96. Keeps the `--print-bootstrap-command` heredoc verbatim and its `$#`-gated dispatch, parses the same argv, maps `--project-root` onto the global `--repo-root`, locates `do-work-cli.sh` in either layout, invokes rather than execs.
- **MODIFIED** `tools/replace-text-section.sh` and its mirror — 355 lines to 40. The embedded Python is deleted; `--help` still prints the literal USAGE with no toolchain.
- **MODIFIED** `tools/validate-suite-manifest.sh` and its mirror — 151 lines to 30.
- **MODIFIED** `tools/fetch-upstream-archive.sh` and its mirror — 115 lines to 51. Keeps the HUP/INT/TERM 129/130/143 traps, the usage exit 2, and the exit 1 total-failure status; maps the three positionals onto flags.
- **MODIFIED** `skills/do-work/tools/do-work-update.sh` — 115 lines to 33. Derives the installed skill root from its own location and passes it as `--skill-root` so the shared-install guard measures the tree the user is actually running.

### Tests and documentation

- **MODIFIED** `_dev/tests/install-suite-behavior.sh` — deleted the Python-fallback and no-JSON-tool lanes (99 lines) and the flaky-jq post-write lane (37 lines); shrank the surviving restricted-PATH list; added a one-time CLI pre-build with explicit mtimes in the repo module and the fixture archive; replaced the Python settings probe with a whole-document byte comparison.
- **MODIFIED** `_dev/tests/update-script-behavior.sh` — both fixture builders carry the Go module and launcher with pinned mtimes; the hostile-archive case now also plants a hostile `do-work-cli.sh` and a hostile `main.go` and proves neither is built nor run; the two shell-grep assertions follow the behaviour into the Go sources; entry-point parity filters build output out of the diff's output by path; `fixture_root` honours `DO_WORK_TEST_FIXTURE_ROOT`.
- **MODIFIED** `_dev/tests/contract-regressions.sh` — eleven grep assertions retargeted at the Go files that now own each behaviour (two kept on the launcher as well); the four whole-output equality comparisons converted to exact whole-line `grep -Fxq` matches on the rendered finding line.
- **MODIFIED** `README.md` — Go 1.26.1+ prerequisite for install, update and runtime. The bootstrap fence is byte-unchanged.
- **MODIFIED** `skills/do-work/docs/prescribed-shell-primitives.md` — the "remains self-contained" claim corrected; only `--print-bootstrap-command` still is.
- **MODIFIED** `skills/do-work/tools/prime-do-work-update.md` — Read-first list points at the Go packages; records the prerequisite; gains a Traps section for the three hazards this migration introduced.

## P-A-U

### [PLAN]

Turn the five public shell entry points into thin launchers over five newly registered `do-work-cli` subcommands. Move only what is genuinely Python-or-shell domain logic; keep every external tool the scripts already shell out to (`cp -Rp`, `tar`, `diff`, `git`, `just`, `bash scripts/atomic-download.sh`), because those preserve byte, mode and symlink semantics for free and keep the existing PATH-stub fixtures working verbatim. Build bottom-up — managedsection, settingshooks, suitemanifest, archivefetch, then the install and update transactions, then the command handlers — running each package's tests before moving on. Then the launchers, then the three shell suites, then documentation, then the gate.

Characterize before replacing: capture the current Python replacer's behaviour on 28 edge cases as golden fixtures, and capture the incumbent jq program's output on the reinstall fixture, so the port is measured rather than assumed.

### [APPLY]

Every file written is inside the Scope declaration's "Files I will touch" list. The diff touches exactly those 39 files. No scope expansion was needed, and nothing outside the worktree was written except this hand-back.

### [UNIFY]

`git diff --stat f38537e..HEAD` reviewed in full; 39 files, all in scope.

Checks run, and what each covered:

- **`gofmt -l` over the module** — empty output. All 20 Go files formatted.
- **`go vet ./...`** — clean across all nine packages.
- **`go test -count=1 ./...`** — all nine packages pass.
- **ShellCheck 0.11.0 at `--severity=warning`** (the gate's severity) on every changed shell file: `tools/install-do-work-suite.sh`, `tools/replace-text-section.sh`, `tools/validate-suite-manifest.sh`, `tools/fetch-upstream-archive.sh`, the four `skills/do-work/tools/` mirrors, `skills/do-work/tools/do-work-update.sh`, `_dev/tests/install-suite-behavior.sh`, `_dev/tests/update-script-behavior.sh`, `_dev/tests/contract-regressions.sh` — 12 files, clean. (`install-suite-behavior.sh` carries one pre-existing SC2016 *info* at line 200, present before this REQ and below the gate's severity.)
- **`bash -n`** on all five launchers and the three suites — all parse.
- **Mirror byte-identity** — `cmp -s` on all four mirrored tool pairs: identical. This is what `staged-skills-contract.sh:812-820` derives its check from.
- **Debug-artifact scan** over the whole diff for `fmt.Print`, `TODO`, `FIXME`, `XXX`, `DEBUG`, `console.log`, `panic(`, `t.Skip`, `debugger` — the only hits are substring false positives (`resul`+`t.Skip`+`pedWork`, and `\uXXXX` inside a comment). No debug artifacts.
- **Working tree** — `git status --short` empty; the built `do-work-cli` binary is covered by the module's shipped `.gitignore`.

File-by-file review notes worth recording:

- `install_transaction.go` — read against the 621-line shell original phase by phase; field names deliberately mirror the shell's variables so the two can be diffed by eye. Confirmed `writeStarted` is set BEFORE the unstage loop and `installVerified` only after every verification, matching the shell's rollback window.
- `managed_section.go` — every one of the three easy-to-"fix" rules carries a comment saying why it is what it is.
- `settings_hooks.go` — confirmed no `map[string]any` is ever marshalled.
- The four launcher mirrors — diffed against each other, not just `cmp`'d, to confirm the CLI-locator block is identical in all four.

## Testing

### RED (carried forward from before the interruption, not re-run)

Captured with today's shell installer under a PATH holding only `awk bash cat chmod cmp cp diff dirname find git grep gzip head just mkdir mktemp mv rm sed stat tar tr wc` — `python3` and `jq` both absent and confirmed absent by probe. Saved at `/tmp/claude-0/-home-user-skill-do-work/cee71dc2-6250-5a0b-9c51-9822eef12052/scratchpad/red/RED-EVIDENCE.txt`.

**RED case A** — existing project with a managed Justfile and custom Stop hooks:

```
exit status: 1
suite manifest valid: v0.245.0 (4 modules)
do-work suite install: python3 is required to reconcile an existing Justfile safely; no client files were changed
```

**RED case B** — fresh project:

```
exit status: 0
  settings reconciler: manual
--- manual instruction occurrences --- 6
--- settings.json written? --- settings.json ABSENT (no hooks composed)
```

### GREEN (re-verified after resuming)

Same two cases, same fixture shape, with the restricted PATH now also omitting `go` (the binary is pre-built and its mtime pinned). Saved at `.../scratchpad/red/GREEN-EVIDENCE.txt`.

```
--- absent under the restricted PATH ---
python3 absent
jq absent
go absent

GREEN CASE A: exit status: 0
  no python3 refusal
  custom-before preserved / custom-after preserved / stale managed content replaced
  retired guard removed / custom Stop hook preserved / empty wrapper preserved
  core SessionStart hook composed / unrelated custom key preserved

GREEN CASE B: exit status: 0
  settings.json WRITTEN (hooks composed)
  core SessionStart hook composed
  fresh justfile is the board template byte for byte
  CLAUDE.md written
  four modules: do-work do-work-board do-work-knowledge do-work-toolbox
```

The installer's own `MANUAL STEP` string and its `settings reconciler:` line are both gone (`grep -c` = 0 for each). Three `MANUAL STEP` hits remain in the captured review diff; all three are `do-work-knowledge` memory-hook prose and its install script, which are REQ-415/REQ-417 territory and out of scope.

### Byte-preservation characterization (re-verified after resuming)

A 28-case harness drove the OLD Python `replace-text-section.sh` and recorded, per case, the resulting target bytes, the target mode, the exit status and stderr. The same harness then drove the NEW Go-backed launcher.

**Result: all 28 cases produce byte-identical and mode-identical targets.** Cases: create-from-template (mode 750 preserved), replace (mode 640 preserved), lone-CR terminators, CRLF, embedded NUL, BOM asymmetry, all four append separators, duplicated begin marker, reversed markers, begin-only, symlink, dangling symlink, directory target, recipe+alias collision, BOM+CRLF collision, self-collision non-case, idempotence (mode 600, unchanged mtime), custom markers, `--help`, missing template, malformed section file, one-sided marker override, multiline-literal boundary, quadruple backtick, section with no recipes.

The lone-CR case is the one with no existing fixture anywhere in the suite. Python and Go both produce `custom:\r` + the LF section + `tail\r`.

Exit statuses: every success case and every collision refusal is unchanged (0 and 1). Eleven malformed-input cases moved from 1 to 2, which is the plan's intended contract — a refusal a consumer can resolve stays 1, a malformed input the caller must correct becomes 2.

### Settings composition parity (re-verified after resuming)

The Go composer's output was diffed against the incumbent jq program's output on the reinstall fixture (custom top-level key, two SessionStart wrappers, a guard-sharing Stop wrapper, a guard-only wrapper, a deliberately empty `preserve-empty` wrapper, and a memory hook): **byte-identical**. That exact document is now the expected file in `install-suite-behavior.sh`.

### Go package tests

`go test -count=1 ./...` in `skills/do-work/tools/do-work-cli` — all nine packages pass.

Red-before / green-after for the four gap-closure items that had a demonstrable RED:

| Test | Failure before | After |
|---|---|---|
| `TestLineSplittingMatchesPythonBytesSplitlines` and 25 siblings in `managedsection` | run against signature-only stubs: `expected a section failure, got nil`, `target bytes = ... want ...` — assertion failures, not compile errors | pass |
| `TestRetiredPipelineGuardRemovalKeepsEveryOtherEntry` and 7 siblings in `settingshooks` | run against a pass-through stub: `the retired pipeline guard survived composition`, `the core SessionStart hook was not composed exactly once` | pass |
| `TestAnUnrecordedCreationIsRolledBackRatherThanReportedAsSuccess` | with the `state.existed` gate removed: `an unrecorded creation reported success: ... Outcome:"success" ... CreatedPaths:[]string{}` | pass with the gate restored |
| `TestMalformedSuitesAreRejectedWithTheirCurrentDiagnostic/a_carriage_return_in_a_row` | `a malformed suite was accepted` — `bufio.Scanner` strips a trailing CR, so a CRLF manifest the shell rejects was passing | pass after splitting on LF only |

New behavioural tests closing the REQ-406 folds: `TestSuccessfulCommitLandsExactlyTheDeclaredTargets` (real repository, exact-path commit, empty index afterwards, revert argv, undeclared path untouched), `TestRefusedCommitRollsBackAndReportsCommitFailed` (a `pre-commit` hook that exits 1 drives a real `commit_failed` rollback), `TestEveryRemediationNamesTheRunningCommandRatherThanAPlaceholder`, `TestUsageFindingsNameARunnableCommand`, `TestUnknownAndMissingCommandFindingsListTheRegisteredCommands`, `TestTextRenderingNamesChangesSkippedWorkAndRollbackErrors`, `TestAResultThatRanNoTransactionRendersNoRollbackLine`.

### Shell suites

| Command | Exit |
|---|---|
| `bash _dev/tests/install-suite-behavior.sh` | 0 (three consecutive runs, 0 FAIL lines each) |
| `bash _dev/tests/update-script-behavior.sh` | 0 |
| `bash _dev/tests/contract-regressions.sh` | 0 |

One genuine bug was found by running these rather than by review: `install-suite-behavior.sh` passed standalone and failed under `contract-regressions.sh` with *"interruption did not complete the filesystem and index recovery path"*. The signal handler was recovering from its own goroutine while the main goroutine was still copying modules, so recovery intermittently reported itself incomplete after removing a directory the main goroutine was still writing into. Fixed in `7685de4` by cancelling a work context from the handler and letting the main goroutine run the single recovery. Verified with three consecutive clean runs.

### Canonical gate

`bash _dev/tests/maintainer-verify.sh` from the worktree root, judged by direct exit status, never piped, `QUEUE_KANBAN_BROWSER` unset so the browser lane stays in its default skipped state.

```
Maintainer verification passed.
=== MAINTAINER-VERIFY EXIT: 0 ===
```

Twelve stages recorded. The last lanes read:

```
SKIP: no browser is available; strict browser behavior lane was not run. Set QUEUE_KANBAN_BROWSER to name one.
maintainer-verify: do-work-cli go vet
maintainer-verify: do-work-cli uncached tests
ok  .../internal/archivefetch       0.113s
ok  .../internal/commandruntime     0.147s
ok  .../internal/gittransaction     0.331s
ok  .../internal/managedsection     0.025s
ok  .../internal/resultmodel        0.004s
ok  .../internal/settingshooks      0.003s
ok  .../internal/suiteinstall       0.710s
ok  .../internal/suitemanifest      0.017s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  github.com/knews2019/skill-do-work/audit-metrics 1.651s
Maintainer verification passed.
```

The ShellCheck lane (every tracked `*.sh`), the `gofmt -l` lane (every tracked `*.go`), the `contract-regressions.sh` aggregate, and the `queue-kanban` / `do-work-cli` / `audit-metrics` vet-and-test lanes all passed. The `--self-test` counting gate is untouched and passing, since no verification lane was added.

## Python and jq Elimination

Walking the exploration's 23-site list. "Removed" means the branch no longer exists anywhere; "replaced by Go" means the behaviour survives in a named Go file.

**Install/update path — python3 (sites 1-6):**

1. `tools/replace-text-section.sh:5-8` + mirror — the `command -v python3` guard and its `python3 is required` message. **REMOVED.** The launcher has no interpreter guard; `do-work-cli.sh` enforces the Go floor instead.
2. `tools/replace-text-section.sh:10` + mirror — `exec python3 - "$@" <<'PY'` and the entire 344-line implementation. **REPLACED BY GO**: `internal/managedsection/managed_section.go` and `just_definitions.go`.
3. `tools/install-do-work-suite.sh:267` + mirror — `if command -v python3` gating whether the Justfile candidate goes through the replacer. **REMOVED**; the install transaction always calls `managedsection.ReplaceSection`.
4. `tools/install-do-work-suite.sh:278-282` + mirror — the no-python3 Justfile branch, including `python3 is required to reconcile an existing Justfile safely`. **REMOVED.** This was half the REQ's RED.
5. `tools/install-do-work-suite.sh:308` + mirror — `if command -v python3` gating CLAUDE.md reconciliation. **REMOVED.**
6. `tools/install-do-work-suite.sh:316-320` + mirror — the no-python3 CLAUDE.md branch. **REMOVED.**

**Install/update path — jq and the settings composer (sites 7-12):**

7. `tools/install-do-work-suite.sh:329-333` + mirror — the `settings_tool` three-way `jq`/`python3`/`manual` selection. **REMOVED.**
8. `tools/install-do-work-suite.sh:356-388` + mirror — the jq reconciliation program. **REPLACED BY GO**: `internal/settingshooks/settings_hooks.go`, verified byte-identical on the reinstall fixture.
9. `tools/install-do-work-suite.sh:389-390` + mirror — `jq -e .` re-parse of the candidate. **REPLACED BY GO**: composition decodes and re-encodes, so an unparseable candidate cannot be produced; malformed input is refused before any write.
10. `tools/install-do-work-suite.sh:392-448` + mirror — the embedded python3 settings composer. **REPLACED BY GO**, same file.
11. `tools/install-do-work-suite.sh:449-450` + mirror — `python3 -m json.tool` re-parse. **REPLACED BY GO**, as site 9.
12. `tools/install-do-work-suite.sh:600-604` + mirror — post-write validation of the INSTALLED settings (`jq -e .` / `python3 -m json.tool`). **REPLACED BY GO**: `verifyInstalledBytes` re-composes the installed file, which is a strictly stronger check than a re-parse. This was the flaky-jq lane's only injection point; see Decisions D-06.

**Consequent dead surface (sites 13-14):**

13. `tools/install-do-work-suite.sh:6` + mirror — `manual_settings_instruction`, the 380-character MANUAL STEP string, and its references at :494 and :620. **REMOVED.** Verified absent by grep across the whole repo.
14. `tools/install-do-work-suite.sh:343`/`:458` and the branches at :493-503, :575-580, :597-605, :619-621 — the whole `settings_tool != manual` block. **REMOVED.** `settings_tool` appears nowhere in the repository.

**Deliberately retained (sites 15-20):**

15. `skills/do-work/tools/checks/preflight.sh:120-130` — the `python3` baseline.json writer with its printf fallback. **RETAINED**, untouched. REQ-414 territory; not on the install/update path.
16. `skills/do-work/tools/checks/preflight.sh:143` — the virtualenv probe. **RETAINED**, untouched. This is the REQ's own "keep Python checks only when probing a Python target capability".
17. `skills/do-work-knowledge/scripts/install-memory-hooks.sh:20,28,46,51` — the jq guard, validation, merge program and merge validation. **RETAINED** (4 occurrences confirmed present). REQ-417 territory.
18. `skills/do-work-knowledge/hooks/memory-stop-capture.sh:49-50,98,112,163-175` — jq transcript extraction, capture program, documented no-jq fallback. **RETAINED** (9 occurrences confirmed present). REQ-415 territory.
19. `skills/do-work-toolbox/scripts/install-last30days.sh:159-169` — the python3.13/3.12/3/python probe for a Python target. **RETAINED**, untouched.
20. `skills/do-work-knowledge/actions/setup-memory.md:57,59,64` and `memory-reference.md:157,158` — prose describing the jq-or-manual contract for memory hooks. **RETAINED**, untouched.

**Harness-side (sites 21-23):**

21. Maintainer-test python3 in `contract-regressions.sh` (17 sites), `staged-skills-contract.sh` (6 sites) and `shipped-package-reference-contract.sh:6`. **RETAINED.** These run in the outer harness process, are not shipped, and are not what the REQ's "Python absent" constraint targets. Untouched.
22. `install-suite-behavior.sh:338-354` and `:538-548` — the two python3 settings-structure assertions. The first is **REPLACED** by a whole-document byte comparison against an expected JSON file; the second went with the Python-fallback lane it belonged to. `install-suite-behavior.sh` now contains zero `python3` references.
23. The three restricted-PATH command lists at `:380`, `:506` and `:554`. `:506` and `:554` went with their lanes. `:380` **SHRUNK** to `awk bash cat cp diff dirname find git grep gzip mkdir mktemp mv rm sed tar` — python3 removed, and `chmod cmp head stat tr wc` dropped as no longer needed. One list remains, and python3's absence from it is the standing proof for the constraint.

All 23 sites reached. No gaps.

## Integration Seams

Lines that must be added to files outside this Scope. None of these were written.

**1. `skills/do-work/actions/version.md` — the Go prerequisite for the update path.** This is the plan's declared coverage gap 1 and belongs to the integrator at the version bump. Suggested line, to go in the prerequisites/safeguards area near the existing engine-delegation text (lines 37-58):

```markdown
**Prerequisite:** Go 1.26.1 or newer. The update engine is the `do-work-cli` command, built from source on first use; `tools/do-work-cli.sh` refuses with an actionable message when the toolchain is missing or too old.
```

**2. Release ritual — not performed, by instruction.** `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md` and `skills/do-work/CHANGELOG.md` are untouched. Suggested changelog entry, for the integrator to place and version:

```markdown
## X.Y.Z — Go-Powered Install and Update (2026-08-30)

Installing and updating do-work now runs on a single Go command instead of a mix of shell, embedded Python and jq. The public commands and the bootstrap snippet are unchanged, but installs no longer depend on which JSON tool happens to be on your machine — settings are always reconciled, and your existing settings.json keeps its exact key order.

- Bootstrap, install, update, managed-section replacement, settings reconciliation, suite validation and archive fetching all run through `do-work-cli`.
- The "no JSON tool available" state is gone, along with the manual-step instruction it used to print.
- **Installing and updating now require Go 1.26.1 or newer.**
```

## Decisions

**D-01 — Keep the subprocesses, delete only the language branches. DECIDE & STATE.**
The Go install still invokes `cp -Rp`, `tar`, `diff`, `git`, `just` and `bash scripts/atomic-download.sh` exactly where the shell did. This preserves byte, mode and symlink semantics without reimplementing `cp -R`, keeps every existing `cp`/`git`/`just`/`curl` PATH-stub fixture working verbatim, and confines the migration to the branches the REQ actually names. Honours plan D1.

**D-02 — Malformed input to `replace-section` exits 2, not 1. DECIDE & STATE.**
The shell tool exited 1 for everything. The plan's command surface assigns exit 1 to refusals a consumer can resolve (the reserved-recipe collision) and exit 2 to malformed inputs the caller must correct. Every collision case keeps exit 1; eleven malformed-marker/symlink/missing-template cases move to 2. No test asserted the specific status for those cases — all check non-zero.

**D-03 — Fragment hook events are iterated in sorted order. DECIDE & STATE.**
jq iterated `keys` (sorted) and was the preferred incumbent branch; Python iterated insertion order. They agree today because the fragment has one event. Sorting matches the incumbent and makes the composed output independent of how the shipped fragment happens to be written.

**D-04 — Settings strings are escaped by `encoding/json` with HTML escaping off, which escapes U+2028 and U+2029 where jq would emit them raw. DECIDE & STATE.**
Everything else matches jq byte for byte, including raw UTF-8 for ordinary non-ASCII and unescaped `<`, `>`, `&`. Hand-rolling a string escaper to close a two-codepoint gap would trade a well-tested encoder for a bespoke one; both forms are valid JSON with identical semantics, and no settings file in the wild is likely to carry a bare line separator inside a string. Recorded rather than left to the diff because it is the one place the port is not byte-exact.

**D-05 — The manifest validator stays fail-fast. DECIDE & STATE.**
The plan says "one finding per violation"; the shell exits at the first `fail`. Fail-fast preserves current behaviour exactly and still satisfies the wording — there is one violation reported because reporting stops there. `suite-manifest-contract.sh` asserts only pass/fail, so nothing depended on seeing more than one.

**D-06 — The flaky-jq post-write lane is deleted rather than replaced with an equivalent injection. DECIDE & STATE.**
Its only injection point was the jq branch. Its three assertions (modules, Justfile and settings all restored to exact originals) are already made by the flaky-`just` lane over the same three surfaces. This is the plan's declared coverage gap 4: post-write settings verification failure is no longer independently injectable from outside the process, and the Go transaction's own tests carry that case.

**D-07 — Interrupted installs exit 130 for HUP, INT and TERM alike. DECIDE & STATE.**
That is what `install-do-work-suite.sh:173` did. The fetcher's 129/130/143 triple is asserted exactly by `update-script-behavior.sh` and stays in the fetcher launcher. Unifying the two would change a status a shipped test pins. This is a signal status, not an outcome, so `resultmodel.ExitCode` remains the only outcome-to-number authority.

**D-08 — The signal handler cancels work and waits; it does not recover. DECIDE & STATE.**
My first implementation recovered from the handler's goroutine and raced the writes it was undoing, producing an intermittent failure that only appeared under the aggregate suite. Recovery now has exactly one caller on one goroutine.

**D-09 — Entry-point parity filters the built binary out of the diff's OUTPUT rather than using `diff -x`. DECIDE & STATE.**
Go embeds the build directory, so a rebuilt binary is path-dependent by construction and cannot be compared as suite content. `diff -x do-work-cli` would also match the module *directory* of the same name and silently blind the check to the entire Go source tree — the exact trap `_dev/primes/prime-shell-commands.md` records. A separate assertion keeps the coverage that both entry points leave a runnable command behind.

**D-10 — The fixture archives carry a pre-built binary, and the mtimes are set explicitly. DECIDE & STATE.**
C17 asked for a deliberate choice. The restricted-PATH lanes have no `go`, so a fixture that ships source only would fail on the launcher's first rebuild. Carrying a binary makes those lanes prove what they exist to prove. The cost — that the fixture archive differs from a real tarball, which ships source only — is bounded by D-09's exclusion and by the update suite's hostile-archive case, which proves an archive's own Go source is never built.

**D-11 — `--print-bootstrap-command` stays a literal heredoc in the launcher. DECIDE & STATE.** Honours plan D9: it is static text, it must run before anything is installed, and moving it would force a non-`CommandResult` stdout exception into the contract fourteen later REQs inherit.

**D-12 — `do-work-update.sh` passes its own `--skill-root`. DECIDE & STATE.**
The shell updater derived the skill root from its own script location, which is what makes the shared-install refusal measure the tree the user is actually running rather than a conventional path. The launcher preserves that by passing it explicitly; the command defaults to `<repo-root>/.claude/skills/do-work` when it is omitted.

**D-13 — `_dev/tests/update-script-behavior.sh` gained a `DO_WORK_TEST_FIXTURE_ROOT` escape hatch. DECIDE & STATE.**
The suite deleted its fixtures on exit, so a failure could only be diagnosed by re-deriving it. This is maintainer-side test scaffolding, not shipped surface, so the "every flag needs a non-test caller" rule in `prime-shell-commands.md` does not apply — that rule scopes to `skills/*/tools/` and `skills/*/scripts/`. It is what let me find the D-08 race in one run instead of several.

Nothing reached the ESCALATE tier. The one judgment call that came closest is D-04, and it is stated rather than escalated because both outputs are valid JSON with identical semantics and the divergence is unreachable for realistic settings files.

## Discovered Tasks

- **`replace-section` is now noisy on success.** The shell tool was silent and exited 0; the command prints a four-line rendered `CommandResult` to stdout for every successful replacement. That is the intended contract, but it makes the `contract-regressions.sh` log substantially longer, and any future caller that captured the old silence as "nothing happened" would now see output. Worth a `--quiet` or a `success` render that prints nothing, decided once for the whole command platform rather than per command.
- **`.claude/settings.json` is always written now, even when composition is a no-op.** The old installer skipped the settings write entirely in the `manual` case. It now always writes, and the write is byte-idempotent, so a reinstall leaves identical bytes — but it does touch the file's mtime on every install. Cheap to make conditional on the composed bytes differing, matching what `managedsection` already does for the Justfile.
- **`skills/do-work/tools/checks/preflight.sh` still writes `baseline.json` through embedded Python** (REQ-414's territory, correctly out of scope here). Now that `settingshooks` exists, that writer has a natural Go home and the printf fallback beside it could go.
- **The two Current-version parsers were unified, in the validator's favour.** `validate-suite-manifest.sh:142-145` required exactly one matching line and took the whole remainder; `do-work-update.sh:29` took only leading digits from the first match. `suitemanifest.ReadActionVersion` implements the validator's stricter form for both paths, and `readInstalledVersion` additionally requires a plain semver. C15 flagged this as easy to do accidentally; it was done deliberately, and the looser parse has no remaining caller. Flagging it here because it is a behaviour change on the update path that no test distinguishes: a malformed `version.md` that the old updater would have partially parsed now fails cleanly.
- **`_dev/tests/maintainer-verify.sh --self-test` still asserts exactly 11 stages.** No lane was added, so it is untouched and passing — but the plan's "no lane needs adding" is now load-bearing for anyone adding a Go package with its own gate.

## Remediation — BOM

**Commit:** `8b4e1b1fd9ba2a47ad5ac54746f6a890843b2102` — "REQ-407: accept a BOM-prefixed settings.json again, as jq did", on branch `worktree-agent-REQ-407-migrate-install-update-bootstrap-to-go`. Two files, staged by explicit path: `skills/do-work/tools/do-work-cli/internal/settingshooks/settings_hooks.go` and its test. Not pushed, not merged.

**The regression.** `decodeOrderedJSON` handed the consumer's bytes straight to `encoding/json`, which refuses a UTF-8 byte-order mark. The incumbent three-way branch preferred `jq` whenever it was installed, and `jq` accepts a BOM and drops it, so a `settings.json` written by a Windows editor or a PowerShell redirect installed before the port and hard-failed the entire install after it.

**RED first.** `TestLeadingByteOrderMarkIsStrippedLikeJq` was added before the fix and failed with the exact production text:

```
--- FAIL: TestLeadingByteOrderMarkIsStrippedLikeJq (0.00s)
    settings_hooks_test.go:182: ComposeSettings: settings are not valid JSON: invalid character 'ï' looking for beginning of value
```

**The fix.** One line in `decodeOrderedJSON`: `data = bytes.TrimPrefix(data, []byte(utf8ByteOrderMark))`, with `utf8ByteOrderMark` a `const` `"﻿"`. Trim, not general leniency — the decoder is otherwise untouched, so every other refusal in the package keeps its current message. Because the strip happens before decoding, the mark is never in the value tree and cannot be re-encoded; the test asserts the output contains no `﻿` anywhere and starts with `{\n`.

**Bounds locked.** Two cases added to `TestMalformedSettingsAreRefusedWithoutProducingOutput`: a doubled BOM (`"﻿﻿{}"`) and a BOM at a non-zero offset (`"{﻿\"hooks\": {}}"`). Both must stay `settings are not valid JSON`. A byte-order mark is one mark at position zero; anything else is malformed input.

**Package gates.** `go vet ./...` and `go test -count=1 ./...` in `skills/do-work/tools/do-work-cli` — exit 0. `gofmt -l` over both touched files — silent.

**End-to-end.** Throwaway git fixture with a BOM-prefixed `.claude/settings.json` carrying one consumer `SessionStart` hook (`echo consumer-own-start`) and a `permissions.allow` block. `printf 'y\n' | bash tools/install-do-work-suite.sh --project-root <fixture>` — exit 0, `install-suite: success`, `.claude/settings.json [modified]: core hooks composed into existing settings`, `rollback: not_needed`. The written file starts `7b 0a` (`{\n`), contains no `EF BB BF` anywhere, parses under both `jq` and `python3 -m json.tool`, and holds two `SessionStart` wrappers: the consumer's own hook once, the core `session-start.sh` hook once, with `permissions.allow` intact.

**Canonical gate.** `bash _dev/tests/maintainer-verify.sh` from the worktree root — **exit 0**, "Maintainer verification passed." Run directly, unpiped, with `QUEUE_KANBAN_BROWSER` unset (the strict browser lane self-skipped, as it does on this host). Run twice: once after the fix, and again after the `var` → `const` tidy-up, so the passing run covers the committed source exactly.
