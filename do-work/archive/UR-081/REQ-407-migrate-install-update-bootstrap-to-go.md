---
id: REQ-407
title: 'Migrate bootstrap, install, update, reconciliation, validation, and fetching into Go'
status: completed
created_at: 2026-08-29T20:28:26Z
route: C
write_set: [skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go, skills/do-work/tools/do-work-cli/internal/managedsection/managed_section_test.go, skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions.go, skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions_test.go, skills/do-work/tools/do-work-cli/internal/settingshooks/settings_hooks.go, skills/do-work/tools/do-work-cli/internal/settingshooks/settings_hooks_test.go, skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest.go, skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest_test.go, skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go, skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go, skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go, skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go, skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go, skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go, skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go, skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go, skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go, skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings_test.go, skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go, tools/install-do-work-suite.sh, skills/do-work/tools/install-do-work-suite.sh, tools/replace-text-section.sh, skills/do-work/tools/replace-text-section.sh, tools/validate-suite-manifest.sh, skills/do-work/tools/validate-suite-manifest.sh, tools/fetch-upstream-archive.sh, skills/do-work/tools/fetch-upstream-archive.sh, skills/do-work/tools/do-work-update.sh, _dev/tests/install-suite-behavior.sh, _dev/tests/update-script-behavior.sh, _dev/tests/contract-regressions.sh, README.md, skills/do-work/docs/prescribed-shell-primitives.md, skills/do-work/tools/prime-do-work-update.md]
estimate:
  p50_active_minutes: 75
  confidence: low
  calculated_at: 2026-08-30T07:23:00Z
  basis:
    - Route C
    - 16-file write set
    - 8 new files
    - 3 subsystems involved
    - 8 acceptance criteria
    - dependency depth 1
    - cross-route regression gates
    - full-suite verification
claimed_at: 2026-08-30T07:22:27Z
completed_at: 2026-08-30T17:36:06Z
commit: f45cdca
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, skills/do-work/tools/prime-do-work-update.md]
tdd: true
suggested_spec:
depends_on: [REQ-406]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
---

# Migrate Bootstrap, Install, Update, Reconciliation, Validation, and Fetching into Go

## What
Move installation and update domain logic into `do-work-cli` and remove Python/jq implementation branches.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements
- Migrate bootstrap/install/update, byte-safe managed-section replacement, settings reconciliation, suite validation, and archive fetching to Go.
- Preserve fresh/existing-project behavior, exact rollback, custom hooks, reserved recipe collision handling, and cancellation.
- Handle CRLF, BOM, NUL, symlinks, file modes, malformed markers/JSON, and existing Just/CLAUDE/settings content.
- Eliminate do-work Python and jq branches and document the Go 1.26.1+ prerequisite for installation, update, and runtime.
- Keep Python checks only when probing a Python target capability.

## Constraints
- Preserve the public scripts-and-Just installation shape through compatibility launchers.
- Installer/update tests must pass with Python and jq absent.

## Dependencies
Depends on REQ-406 (shared CLI and transaction foundation).

## Builder Guidance
Certainty level: Firm. Characterize current byte and filesystem behavior before replacing implementation branches.

## Red-Green Proof
**RED prompt/case:** Run installer/update fixtures with Python and jq removed from PATH, including an existing managed Justfile and custom settings hooks.
**Why RED now:** Existing reconciliation and managed-section branches conditionally depend on embedded Python or jq.
**GREEN when:** All fresh/existing fixtures succeed with Go alone, preserve exact unrelated bytes/state, and roll back failures according to the shared transaction contract.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Folded From REQ-406 (2026-08-30)

REQ-406 built the foundation and stopped at the command seam, so these became
testable only once a real command is registered — which is this REQ. Folded here
per the Fold-First Rule rather than minted as separate REQs.

- **The `<command>` placeholder is not a runnable argv.** `usageFinding` and the
  `invalid_options` template both emit `do-work-cli --format text <command>`, which
  shows the shape but cannot be pasted. Requirement 5 asks every finding for the
  *exact* next argv. Once this REQ registers commands the runtime knows the name and
  can thread it through.
- **No test observes a successful `--commit` transaction, and `commit_failed` has no
  behavioural test.** REQ-406's fixtures cover exit 0/1/3/4 but the committing success
  path and the commit-failure kind are asserted only through their finding templates.
- **`RollbackResult.Status` has a fourth wire value, the empty string,** for results
  that never ran a Git transaction. A consumer switching on it must handle `""`
  alongside the three constants. REQ-406's D-04 explains why normalising it to
  `not_needed` is not free: every read-only command would print a rollback line
  implying a mutation was possible.
- **The success path does not consult `state.existed`.** At
  `git_transaction.go:161-166` an unrecorded change to a declared target is detected,
  but a file created without `RecordCreated` still reports `succeeded`. Harmless while
  no command is registered; worth closing when one is.
- **Text rendering of changes, skipped work and rollback errors has no direct
  assertion.** The parity test covers findings; these three sections render unasserted.

---

## Triage

**Route: C** - Complex

**Reasoning:** Replaces the implementation language of the suite's installation path — bootstrap, install, update, byte-safe managed-section replacement, settings reconciliation, suite validation and archive fetching all move to Go, and the Python and jq branches are removed rather than left as fallbacks. The behaviour to preserve is byte-exact and filesystem-sensitive (CRLF, BOM, NUL, symlinks, modes, malformed markers), the rollback contract is inherited from REQ-406, and the public scripts-and-Just shape must keep working through compatibility launchers.

**Planning:** Required

## Plan

Turn the five public installation/update shell scripts into thin launchers over five newly registered `do-work-cli` subcommands, moving only the logic that is genuinely shell-or-Python domain code and keeping every external tool the current scripts already shell out to.

What actually moves to Go: the Python managed-section replacer and its Just-definition scanner (`replace-text-section.sh`, 350 lines of embedded Python), the jq/Python settings-hook reconciliation, the shell manifest validator, the two-route archive fetch orchestration, and the installer/updater control flow (module plan, symlink/mode guards, review diff, single confirmation, backup + exact recovery including the Git index, post-write verification, version comparison).

What deliberately does NOT move: the external tools the scripts invoke. The Go install command keeps shelling out to `cp -Rp`, `tar xzf`, `diff -ruN`/`diff -qr`, `git`, `just --list`, and `bash scripts/atomic-download.sh` exactly where the current script does. This preserves byte/mode/symlink semantics for free, keeps the existing PATH-stub fixtures (`cp`, `git`, `just`, `curl`) working verbatim, and avoids reimplementing `cp -R` in Go. The Python and jq branches are what get deleted, not every subprocess.

Biggest single deletion: the `settings_tool` three-way branch (`jq` / `python3` / `manual`) and the whole `manual_settings_instruction` path — roughly 130 lines of installer shell plus two test lanes. With Go always able to reconcile, "no JSON tool" ceases to exist as a state.

Byte-preservation obligations the port must honour exactly: `bytes.splitlines` semantics (a lone CR ends a line), `rstrip(b"\r\n")` for marker-body comparison, BOM handling on line 0 of the Just scanner, target mode preserved on replace and template mode on create, atomic temp-in-target-directory + rename, and order-preserving JSON so an existing `settings.json` keeps its key order (Go's default map marshalling would silently reorder unrelated user state).

Output contract that fourteen later REQs inherit: stdout carries only the rendered `CommandResult`; every progress line, review diff, and confirmation prompt goes to stderr. Existing diagnostic phrases are carried through as finding evidence and narration text verbatim, so the substring assertions in the three behavioural suites survive.

Test posture is RED-first per the REQ: build the restricted-PATH fixture with `python3` and `jq` absent and show today's installer failing (existing Justfile refused, settings left to the manual instruction), then make it green with Go alone.

**Command surface (the pattern fourteen later REQs inherit):**

Binary `do-work-cli`, reached through `skills/do-work/tools/do-work-cli.sh`. Global options are unchanged and must precede the command: `--repo-root <path>` (defaults to cwd) and `--format text|json` (defaults to text). Five commands are registered; each is a `commandruntime.CommandHandler` returning a `resultmodel.CommandResult`, rendered by the runtime. stdout carries only that rendered result; narration goes to stderr. Exit codes come solely from `resultmodel.ExitCode`: 0 success, 1 findings or safely refused, 2 usage/precondition/runtime failure, 3 rolled back, 4 incomplete rollback or committed-state risk.

1. `install-suite [--archive <path>]` — target project is the global `--repo-root`, which must be a Git worktree root. Without `--archive`, fetches `DO_WORK_UPSTREAM_URL` (default the main-branch tarball) through the archivefetch package. Reads one line from stdin for its single confirmation.
   Returns: success with `changes` naming the four module destinations, the Justfile, CLAUDE.md and .claude/settings.json; success with `skipped_work` code `INSTALL-CANCELLED` when the confirmation is declined (exit 0, nothing written); refused (exit 1) for a symlinked or escaping managed destination, a reserved Just recipe collision, or a source tree containing a symlink; failure (exit 2) for a non-Git root, a bad archive, or manifest validation; rolled_back (exit 3) when a write or post-write verification fails and every managed path plus the Git index is restored; committed_state_risk (exit 4) when recovery is incomplete.
   Launcher: `install-do-work-suite.sh --project-root <root> [--archive <file>]` maps `--project-root` onto the global `--repo-root`, and still answers `--print-bootstrap-command` from its own heredoc without invoking Go.

2. `update-suite` — target project is the global `--repo-root`. Fetches, extracts, validates the manifest, compares `actions/version.md` against the upstream VERSION, calls the install transaction in-process, then re-reads the installed version.
   Returns: success with `skipped_work` code `UPDATE-ALREADY-CURRENT` when versions match; failure (exit 2) when upstream is older, when the fetch or validation fails, or when post-update verification disagrees; otherwise whatever install returns, including the cancelled-as-success shape. Launcher: `do-work-update.sh --project-root <root>`.

3. `replace-section --target <path> --section-file <path> [--template-file <path>] [--reject-recipe-collisions] [--begin-marker <line> --end-marker <line>]` — repository-independent; operates on the named files.
   Returns: success with one `changes` entry (kind created or modified) or with no changes when the target already matches byte for byte; refused (exit 1) with code `SECTION-RESERVED-RECIPE-COLLISION` whose evidence is `target defines reserved Just recipe or alias outside managed section: <sorted, comma-separated names>`; failure (exit 2) for malformed or duplicated markers, a symlinked target, a missing `--template-file` for an absent target, or a section file that is not one newline-terminated managed section. Launcher: `replace-text-section.sh` with the same argv, unchanged.

4. `validate-manifest --root <archive-root>` — read-only.
   Returns: success, evidence `suite manifest valid: v<version> (4 modules)`; failure (exit 2) with one finding per violation, evidence carrying the current phrasing verbatim (`manifest header must be exactly: source<TAB>destination`, `line N source traverses directories`, `VERSION must be a plain semantic version (X.Y.Z)`, and the rest). Launcher: `validate-suite-manifest.sh --root <archive-root>`, unchanged.

5. `fetch-archive --target <archive-path> --url <tarball-url> [--repo-url <git-url>]` — writes only the named archive path.
   Returns: success with one `changes` entry on the archive path whose detail names the winning route (`fetched over HTTP`, or `fetched with git (HTTP route failed …)`), so the route stays machine-readable and greppable on stdout; failure (exit 2) with a finding whose evidence names both route outcomes and the `DO_WORK_UPSTREAM_URL` escape hatch, leaving any pre-existing target untouched and no scratch behind. Launcher: `fetch-upstream-archive.sh <target> <url> [repo-url]` maps the three positionals to flags, keeps its 129/130/143 signal traps, and keeps exit 2 for its own usage error.

The pattern the fourteen later REQs copy: one internal package per domain exposing pure functions plus a handler file; `cmd/do-work-cli/main.go` gains one import and one map entry per command; the global `--repo-root` is the only way to name the repository; command flags are long-form only, with no single-letter aliases and no positional arguments (positionals stay in the compatibility launcher where a legacy argv demands them); every handler returns a `CommandResult` and never calls `os.Exit` or writes to stdout itself; each command passes its own name into `gittransaction.BuildCommandResult` so every finding's `next_argv` is a runnable command line.

**Tasks:**

1. Add internal/managedsection (byte-exact port of the replace-text-section Python: marker span location with bytes.splitlines semantics, create-from-template, append, in-place replace, atomic temp-in-parent + rename, mode preservation) and internal/settingshooks (order-preserving JSON decode/encode, hooks.json event composition with append-unique, removal of only retired pipeline-guard command objects, empty-wrapper pruning, Stop-key deletion when empty). Register the `replace-section` command and rewrite both mirrors of replace-text-section.sh as launchers. Carry the existing diagnostic phrases (`target defines reserved Just recipe or alias outside managed section: ...`, `must contain exactly one begin marker and one end marker`, `target must not be a symlink`) through as finding evidence verbatim.
   - *Serves:* Detailed Requirements: byte-safe managed-section replacement and settings reconciliation move to Go; CRLF/BOM/NUL/symlinks/file modes/malformed markers/JSON handled; eliminate Python and jq branches. Constraint: public scripts shape preserved through compatibility launchers.
2. Add internal/suitemanifest (VERSION and skills/do-work/VERSION semver + single-newline checks, modules.tsv header/tab/CR/duplicate/traversal/absolute checks, the four required source→destination pairs, real-directory and non-empty SKILL.md checks, actions/version.md Current-version agreement) and internal/archivefetch (HTTP route via `bash scripts/atomic-download.sh` located relative to os.Executable at ../../scripts/atomic-download.sh, git route via shallow clone + `git archive --prefix`, branch derivation from the tarball URL, stage-then-rename publication, scratch cleanup, DO_WORK_UPSTREAM_URL default). Register `validate-manifest` and `fetch-archive`; rewrite both mirrors of validate-suite-manifest.sh and fetch-upstream-archive.sh as launchers that keep their current HUP/INT/TERM traps and exit statuses. Report the winning fetch route as a RecordedChange on the archive path so the route name stays on stdout.
   - *Serves:* Detailed Requirements: suite validation and archive fetching move to Go. Constraint: public scripts-and-Just installation shape preserved through compatibility launchers.
3. Add internal/suiteinstall with the install and update transactions and the five command handlers. Install: resolve and require the Git worktree root, extract the supplied or fetched archive, validate the manifest, build the four-module plan, reject symlinked/non-directory/escaping destinations and symlink-bearing sources, build and validate the Justfile and CLAUDE.md candidates through managedsection, compose settings through settingshooks, print the full review diff plus the dirty-managed warning, take the single confirmation, snapshot originals and the Git index, unstage module paths, write modules, replace Just/CLAUDE/settings atomically, verify installed bytes and versions, and recover every managed path plus the index on any failure or on SIGHUP/SIGINT/SIGTERM. Update: fetch, extract, validate, compare versions, delegate to install in-process, verify the installed version. Delete the settings_tool branch and manual_settings_instruction entirely. Rewrite both installer mirrors and do-work-update.sh as launchers, keeping the --print-bootstrap-command heredoc byte-identical in the installer launchers.
   - *Serves:* Detailed Requirements: bootstrap/install/update move to Go; preserve fresh/existing-project behaviour, exact rollback, custom hooks, reserved recipe collision handling, and cancellation; eliminate Python and jq branches. Constraints: both.
4. Close the five gaps folded in from REQ-406. Thread the real command name into every finding: parseGlobalOptions scans ahead for the command token so option errors name it, usageFinding takes the name, and BuildCommandResult takes a command name so no template emits `<command>`; when no command is known, the finding lists the registered commands and names a real one. Gate the transaction success path on state.existed so a target created without RecordCreated is treated as an unrecorded mutation and rolled back. Add behavioural fixture-repo tests for a successful --commit transaction and for the commit_failed kind. Assert that a result which ran no Git transaction renders no rollback line in text and carries the empty-string status in JSON. Assert the text rendering of changes, skipped work, and rollback errors directly.
   - *Serves:* Folded From REQ-406: the `<command>` placeholder is not a runnable argv; no test observes a successful --commit or commit_failed; RollbackStatus has a fourth empty-string wire value; the success path does not consult state.existed; changes/skipped/rollback-error rendering is unasserted.
5. Retarget the test and documentation surface. install-suite-behavior.sh: delete the Python-fallback and no-JSON-tool lanes with their manual-instruction assertions, shrink the restricted-PATH command lists to what Go actually invokes so the absence of python3 and jq is the proof, pre-build the CLI once and copy it into fixtures so restricted-PATH runs never need `go`, replace the flaky-jq post-write lane with the flaky-just lane's assertions, and rewrite the Python-based settings structure assertion as a byte comparison against an expected JSON file. update-script-behavior.sh: carry the do-work-cli module and launcher into build_suite_install and build_suite_tree, extend the hostile-archive case to plant a hostile do-work-cli.sh and Go source and prove neither is built or run, capture the total-failure fetcher report with 2>&1, and retarget the two shell-grep assertions (fetch-upstream-archive.sh delegation, DO_WORK_UPSTREAM_URL) at the Go sources. contract-regressions.sh: retarget the eleven assertion blocks that pin implementation strings inside the five scripts, and convert the four exact-whole-output replace-text-section comparisons to exact `grep -Fx` matches on the rendered finding line. Document the Go 1.26.1+ prerequisite in README.md, correct the self-contained-installer claim in prescribed-shell-primitives.md, and update prime-do-work-update.md's Read-first list to point at the Go packages.
   - *Serves:* Detailed Requirements: document the Go 1.26.1+ prerequisite for installation, update, and runtime. Constraint: installer/update tests must pass with Python and jq absent. Red-Green Proof: all fresh/existing fixtures succeed with Go alone.

**Files to touch:**

- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go` (new) — Byte-exact port of the replace-text-section Python: marker span location, create/append/replace, atomic temp-in-parent rename, mode preservation.
- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section_test.go` (new) — Unit coverage for CRLF, BOM, NUL, lone-CR line splitting, malformed markers, symlink refusal, mode preservation, and idempotence.
- `skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions.go` (new) — Port of the Just recipe/alias scanner (quote, triple-quote, backtick and triple-backtick state machine, alias regex, BOM) used for reserved-recipe collision rejection.
- `skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions_test.go` (new) — Pins the multiline-literal boundaries that contract-regressions currently proves only end to end.
- `skills/do-work/tools/do-work-cli/internal/settingshooks/settings_hooks.go` (new) — Order-preserving JSON settings reconciliation replacing both the jq program and the embedded Python.
- `skills/do-work/tools/do-work-cli/internal/settingshooks/settings_hooks_test.go` (new) — Pins key-order preservation, append-unique composition, retired-guard removal, custom-hook preservation, and empty-wrapper pruning.
- `skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest.go` (new) — Suite validation moved from validate-suite-manifest.sh.
- `skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest_test.go` (new) — Covers the manifest rejection cases suite-manifest-contract.sh drives end to end.
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go` (new) — Two-route archive fetch moved from fetch-upstream-archive.sh, still delegating the HTTP route to atomic-download.sh.
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go` (new) — Covers branch derivation, git-route export-ignore honouring, stage-then-rename, and preserved-target-on-total-failure.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (new) — The install transaction: plan, guards, candidates, review, confirmation, backups, writes, verification, recovery, signal handling.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (new) — Fresh and existing-project fixtures, recovery on post-write failure, and cancellation as a success-with-skipped-work outcome.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go` (new) — The update transaction: fetch, validate, version comparison, in-process install delegation, post-update verification.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go` (new) — Covers up-to-date, older-upstream refusal, and cancelled-update outcomes without the removed cancel-status env var.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` (new) — The five command handlers and their argv parsing; the single place a FailureKind or transaction outcome becomes a CommandResult for this family.
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (new) — Pins each command's argv contract, exit codes, and the stdout-result / stderr-narration split.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modify) — Register the five handlers; this one-line-per-command wiring is the pattern the fourteen later REQs copy.
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` (modify) — Thread the real command name into usage findings and list registered commands, replacing the `<command>` placeholder.
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` (modify) — Assert usage findings now name a runnable argv.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modify) — Gate the success path on state.existed so an unrecorded creation is not reported as succeeded.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modify) — Add the successful --commit and commit_failed behavioural cases folded in from REQ-406.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go` (modify) — BuildCommandResult takes the command name so no remediation template emits `<command>`.
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings_test.go` (modify) — Update for the new signature and assert the exact next argv.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify) — Assert text rendering of changes, skipped work, and rollback errors, and the empty RollbackStatus wire value.
- `tools/install-do-work-suite.sh` (modify) — Becomes a launcher: locate do-work-cli.sh in either layout, keep the --print-bootstrap-command heredoc and the signal traps, translate --project-root to the global --repo-root.
- `skills/do-work/tools/install-do-work-suite.sh` (modify) — Byte-identical mirror of the root installer launcher.
- `tools/replace-text-section.sh` (modify) — Becomes a launcher over `replace-section`; the embedded Python is deleted.
- `skills/do-work/tools/replace-text-section.sh` (modify) — Byte-identical mirror.
- `tools/validate-suite-manifest.sh` (modify) — Becomes a launcher over `validate-manifest`, preserving --root.
- `skills/do-work/tools/validate-suite-manifest.sh` (modify) — Byte-identical mirror.
- `tools/fetch-upstream-archive.sh` (modify) — Becomes a launcher over `fetch-archive`, mapping the three positionals to flags and keeping the 129/130/143 traps and the usage exit 2.
- `skills/do-work/tools/fetch-upstream-archive.sh` (modify) — Byte-identical mirror.
- `skills/do-work/tools/do-work-update.sh` (modify) — Becomes a launcher over `update-suite`; the version-order awk, fetch, extract, validate, and delegation logic move to Go.
- `_dev/tests/install-suite-behavior.sh` (modify) — Drop the jq/Python/manual lanes, shrink the restricted-PATH lists to prove python3 and jq are unused, pre-build and copy the CLI into fixtures, and replace the Python-based settings assertion.
- `_dev/tests/update-script-behavior.sh` (modify) — Fixtures must carry the Go module and launcher; extend the hostile-archive case to the CLI; retarget the two shell-grep assertions; capture the fetcher total-failure report with 2>&1.
- `_dev/tests/contract-regressions.sh` (modify) — Retarget the eleven assertion blocks that pin implementation strings inside the five scripts, and convert the exact-whole-output replace-text-section comparisons to exact matches on the rendered finding line.
- `README.md` (modify) — Document the Go 1.26.1+ prerequisite for installation and update; the bootstrap block itself stays byte-identical to --print-bootstrap-command.
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — Line 23 claims the installer 'remains self-contained'; that becomes false once it delegates to do-work-cli.
- `skills/do-work/tools/prime-do-work-update.md` (modify) — Its Read-first list names the shell implementations; it must point at the Go packages and record the Go prerequisite.

**Decisions the builder must honour:**

- D1 Keep the subprocesses, delete the language branches. The Go install still invokes cp -Rp, tar, diff, git, just and bash scripts/atomic-download.sh exactly where the shell does. This preserves byte/mode/symlink semantics without reimplementing cp -R, keeps the existing cp/git/just/curl PATH-stub fixtures working verbatim, and confines the migration to the Python and jq branches the REQ actually names.
- D2 Global --repo-root names the target project; command flags name only command-specific inputs. The launchers translate --project-root <path> into --repo-root <path>, so the CLI keeps one repository-root concept. This is the pattern the fourteen later REQs copy.
- D3 stdout carries only the rendered CommandResult; all narration (progress, review diff, dirty-path warning, confirmation prompt, success line) goes to stderr. Every behavioural suite already captures 2>&1, so the existing substring assertions survive, and --format json stays parseable. The one fetcher assertion that reads stderr only changes its capture.
- D4 Existing diagnostic phrases become finding evidence verbatim ('manifest header must be exactly', 'source traverses directories', 'managed destination must not be a symlink', 'target defines reserved Just recipe or alias outside managed section: ', 'discards those changes', 'restored every managed path and the Git index to their exact pre-install state'). Substring assertions in three suites depend on them; only the four whole-output equality checks are restated as exact matches on the rendered finding line.
- D5 install-suite keeps its own snapshot-and-restore recovery rather than gittransaction. gittransaction refuses directory targets ('declare exact files instead') and cannot restore the Git index; the installer replaces four whole module trees and must restore the index. It maps its outcome into the same resultmodel vocabulary (success / refused / rolled_back / committed_state_risk) so callers see one contract.
- D6 Cancellation is outcome success with a SkippedWork entry, not a refusal. That keeps today's exit 0 through the public shell path without a second outcome-to-number table, and lets update-suite distinguish cancelled from installed in-process. DO_WORK_INSTALL_CANCEL_EXIT_STATUS is deleted — its only caller was do-work-update.sh, which no longer runs the installer as a subprocess.
- D7 The settings_tool three-way branch and manual_settings_instruction are deleted outright. With Go always able to reconcile, 'no JSON tool available' stops being a reachable state; roughly 130 lines of installer shell and two test lanes go with it.
- D8 Settings JSON is decoded and re-encoded order-preserving (2-space indent, HTML escaping off, trailing newline). Go's default map marshalling sorts keys, which would silently reorder a user's settings.json and violate 'preserve exact unrelated bytes/state'.
- D9 The --print-bootstrap-command snippet stays a literal heredoc in the installer launcher, byte-identical in both mirrors and in README. It is static text, it must run before anything is installed, and moving it would force a non-CommandResult stdout exception into the contract that fourteen later REQs inherit.
- D10 The winning fetch route is reported as a RecordedChange on the archive path, not a finding. That keeps the route name on stdout for the existing assertion without inventing an info-severity finding on a success outcome.
- D11 Install handles SIGHUP/SIGINT/SIGTERM itself and runs the same recovery path, so the existing cp-shim TERM-during-write fixture keeps working unchanged. The launchers keep their current traps and pass the Go exit status through, which is why they invoke rather than exec the CLI.

**Testing approach:**

RED first, exactly as the REQ's Red-Green Proof states. Build the restricted-PATH fixture (python3 and jq absent, `just` present) against an existing managed Justfile and a settings.json carrying custom Stop hooks plus the retired pipeline guard, and record today's failures: the installer refuses with 'python3 is required to reconcile an existing Justfile safely', and on a fresh project it falls through to the manual settings instruction instead of composing hooks. Those two are the RED.

Go-side unit tests carry the byte-level burden, because that is where the edge cases are cheapest to enumerate. managedsection: lone-CR line splitting (Python's bytes.splitlines ends a line on a bare \r — a Go strings.Split on \n silently diverges here), CRLF marker bodies, a BOM-prefixed first line, an embedded NUL outside the section, reversed and duplicated markers, symlinked target refusal, target mode preserved on replace, template mode preserved on create, the four separator cases when appending to a marker-free file, and byte-idempotence on a second run. just_definitions: every multiline-literal opener form, an escaped quote inside a cooked string, a triple-backtick that is not a delimiter, aliases, @-prefixed recipes, and := assignments that are not recipes. settingshooks: key order preserved across a round trip, append-unique per event, removal of only pipeline-guard command objects, a same-wrapper custom hook preserved, an emptied wrapper dropped, a deliberately empty unrelated wrapper preserved, and Stop deleted when nothing remains.

The three shell suites stay the end-to-end characterization layer and keep every fixture whose injection point survives — the flaky `just` stub, the failing `git restore --staged` stub, the TERM-sending `cp` stub, and the 429 `curl` stub all still work because Go invokes those same binaries. The restricted-PATH command lists shrink to what Go actually calls (bash, cp, diff, git, tar, gzip, curl, and `just` where the case wants it); python3 and jq never reappear, and their absence is the standing proof for the constraint. The CLI binary is built once by the harness and copied into fixtures so restricted-PATH runs never need `go` on PATH.

Two fixtures need genuine rework rather than a copy. The hostile-archive case gains a planted hostile `skills/do-work/tools/do-work-cli.sh` and Go source in the upstream tarball, and asserts neither is built nor run — the trusted-engine guarantee now covers the Go module, and this is the case that would catch an update that rebuilds from the archive before the write boundary. The flaky-jq post-write lane loses its injection point; its assertions (modules, Justfile and settings all restored to exact originals) fold into the flaky-`just` lane, which already exercises the same recovery over the same three surfaces.

Gate: `bash _dev/tests/maintainer-verify.sh` from the repo root, exit 0, browser lane left in its default skipped state. It already runs the module's go vet and uncached tests plus the launcher probe, so no lane needs adding.

**Risks:**

- R1 Scope. Thirty-nine files, sixteen of them new Go, across five public shell entry points and three behavioural suites. The REQ's own p50 of 75 minutes is not reachable; the managed-section port and the installer transaction are each a multi-hour job on their own. Five tasks is the honest floor only because tasks 1-3 are each a whole subsystem. Expect the write set to be the collision surface for any parallel builder.
- R2 The lone-CR line split. Python's bytes.splitlines ends a line on a bare \r; a Go port that splits on \n only will mislocate marker spans in any file containing an old-Mac line ending, and no existing fixture covers it. This is the single most likely silent byte-level regression.
- R3 JSON key order. Go's encoding/json sorts map keys. If the port marshals a map[string]any anywhere, every consumer's settings.json is reordered on the next install — unrelated user state changed, which the GREEN criterion forbids, and no current test would catch it because the existing assertions check content, not order.
- R4 Test injection points that no longer exist. The flaky-jq stub and any assertion that greps a shell file for its own implementation strings (eleven blocks in contract-regressions.sh, two in update-script-behavior.sh) have to be retargeted or retired. Retiring one that still had value is easy to do by accident; each retarget should name the Go file that now owns the behaviour.
- R5 The running binary deletes itself. In the update flow the CLI lives at .claude/skills/do-work/tools/do-work-cli/do-work-cli and install rm -rf's that directory. Unlinking a running executable is fine on Linux and macOS, but the recovery path restores the backed-up module tree including the old binary, and the launcher's mtime staleness check then decides whether to rebuild. Worth an explicit fixture on the recovery-after-failure path.
- R6 Go becomes a hard installation prerequisite. --print-bootstrap-command, manifest validation and section replacement all work today with no toolchain; after this REQ every one of them needs Go 1.26.1+. The launcher already fails with an actionable message, but anyone bootstrapping on a machine without Go now stops at the first step rather than at the install.
- R7 First build during install. The launcher builds the CLI before the installer runs, creating an untracked binary inside the module directory. It is covered by the module's own .gitignore, which ships — but if that .gitignore is ever export-ignored or dropped, every install would report its own build output as a dirty managed path.
- R8 The dry-run/commit surface of the shared transaction is exercised by none of this REQ's commands, since install keeps its own recovery. Task 4's coverage for --commit success and commit_failed therefore lives entirely in gittransaction's package tests, with no command-level path behind it until REQ-409.

**Declared coverage gaps (accepted, not hidden):**

- The Go 1.26.1+ prerequisite line for the update path belongs in skills/do-work/actions/version.md, the canonical agent-driven update contract, but that file is the integrator's and is excluded from this write set. Covered here in README.md and prime-do-work-update.md only; the integrator should add the matching line when bumping the version.
- The bootstrap snippet's own text stays a shell heredoc in the installer launcher (D9). If 'migrate bootstrap into Go' is read to include emitting that string from Go, this is uncovered and needs a stdout exception in the result contract to close.
- End-to-end signal coverage narrows to what the cp-shim reaches. The TERM-during-module-write case survives, but there is no end-to-end probe of a signal arriving during the Justfile or settings write phase; those are covered only by the Go install transaction's own tests.
- The second post-write-failure lane is lost with the flaky-jq stub. Post-write settings verification failure is no longer independently injectable from outside the process; its recovery assertions fold into the flaky-just lane, which restores the same three surfaces but reaches them through a different failure.
- R4 is scoped to the installation and update path. The remaining do-work Python and jq branches — preflight.sh's baseline.json writer, the knowledge memory-hook jq paths, and install-memory-hooks.sh's jq requirement — stay put; they are REQ-414, REQ-415 and REQ-417 territory, and the REQ's own 16-file estimate basis confirms that reading. R5 (Python checks only when probing a Python target) is satisfied by leaving preflight.sh's requirements.txt/virtualenv probe untouched.
- Cross-platform behaviour of the retained subprocesses is assumed, not tested. cp -Rp, diff -ruN and tar are invoked with the same flags the shell used, so BSD/GNU differences are inherited unchanged rather than newly introduced, but nothing in the suite runs on a BSD userland.

**Plan validation:** 5 tasks, at the quality cap. Every task traces to a stated requirement. The six coverage gaps above are declared rather than papered over; gap 1 (the Go prerequisite line in `actions/version.md`) is the integrator's to close at the Commit Phase, and gap 5 correctly scopes the remaining Python and jq branches to REQ-414, REQ-415 and REQ-417.

*Generated by Plan agent*

## Exploration

**Key files:**

- `/home/user/skill-do-work/tools/replace-text-section.sh` — The managed-section replacer. A 9-line bash preamble (python3 guard) then `exec python3 - "$@" <<'PY'` holding 344 lines of embedded Python: marker-span location, Just recipe/alias scanner with a multiline-literal state machine, atomic replace. Byte-identical mirror at skills/do-work/tools/replace-text-section.sh (md5 a14a442869be106cb9d9fc94e2ab378e). Mirror identity is enforced by _dev/tests/staged-skills-contract.sh:812-820, which derives the mirrored set from tools/*.sh on disk rather than a hand list. (1-8 python3 guard; 10 exec; 21-24 BEGIN/END/JUST_IDENTIFIER/JUST_ALIAS constants; 27-29 die(); 32-45 read_regular; 48-55 lines_with_offsets; 58-59 line_body; 62-81 marker_span; 84-103 just_delimiter_matches; 106-118 just_opening_delimiter; 121-157 just_definition_name; 160-186 just_multiline_string_state; 189-218 just_definition_names; 221-260 atomic_replace; 263-267 USAGE; 269-299 argv parsing + marker override validation; 301-304 section-file validation; 306-316 create-from-template; 318-322 existing-target read; 324-337 collision rejection; 339-350 replace-or-append; 352-353 conditional write)
- `/home/user/skill-do-work/tools/install-do-work-suite.sh` — The whole install transaction: bootstrap heredoc, argv, Git-root gate, archive fetch/extract/validate, four-module plan, symlink and escape guards, Justfile and CLAUDE.md candidate construction, jq/python3/manual settings reconciliation, review diff, single confirmation, backups (incl. Git index), unstage, writes, post-write verification, and exact recovery. 621 lines. Byte-identical mirror at skills/do-work/tools/install-do-work-suite.sh (md5 084aba1413670cfefeec19ff4f57b3e5). (5 DO_WORK_UPSTREAM_URL default; 6 manual_settings_instruction; 7-10 DO_WORK_INSTALL_CANCEL_EXIT_STATUS; 12-15 fail(); 17-32 print_bootstrap_command heredoc; 34-37 --print-bootstrap-command gate (only when $#==1); 39-57 argv; 59-72 project-root/Git-root/index-path resolution; 74-84 sibling tool location + install_tmp; 86-100 state vars incl. instructions markers at 98-99; 102-159 recover_install; 161-171 cleanup; 172-173 traps (EXIT cleanup; HUP INT TERM all -> exit 130); 175-188 archive selection/fetch; 190-194 extract + validate + suite_version; 196-202 module plan from modules.tsv; 204-223 destination symlink/dir/escape guards; 225-231 source symlink guard; 233-238 template/hook/instructions presence; 240-255 Justfile discovery (justfile, Justfile, .justfile); 257-264 awk section extraction from board template; 266-294 Justfile candidate + marker/recipe/just --list validation; 296-325 CLAUDE.md candidate + validation; 327-341 settings_tool selection + settings target probe; 343-458 settings composition (jq 356-390, python3 391-450, chmod/grep checks 452-457); 460-464 plan preamble on stdout; 466-505 review diff build + cat; 507-513 dirty-managed warning (stderr); 514-519 single confirmation; 521-548 backups incl. Git index; 550-558 unstage loop, write_started=1 at 553; 560-584 writes; 586-615 post-write verification; 617-621 success)
- `/home/user/skill-do-work/tools/validate-suite-manifest.sh` — Pure-shell suite validator. Reads --root, checks VERSION and skills/do-work/VERSION semver + single-newline, modules.tsv header/tab/CR/blank/empty-column/traversal/absolute/duplicate checks, the four required source->destination pairs, real-directory + non-empty SKILL.md per module, physical containment under <root>/skills/, exactly four rows, and actions/version.md Current-version agreement. No python/jq. Byte-identical mirror at skills/do-work/tools/. (5-8 fail(); 10-12 usage; 14-18 root resolution; 20-26 VERSION checks; 28-33 manifest header; 35-42 accumulators; 44-118 per-row loop (46 blank, 47-49 CR, 51-59 tab/columns, 61-72 traversal/absolute, 74-81 duplicates, 83-103 required pairs, 105-115 module dir/SKILL.md/containment); 120-123 four-module completeness; 125-134 core VERSION; 136-149 actions/version.md; 151 success line on stdout)
- `/home/user/skill-do-work/tools/fetch-upstream-archive.sh` — Two-route archive fetch. Route 1 delegates to skills/do-work/scripts/atomic-download.sh (located by probing two mirror-relative depths); route 2 shallow-clones and repacks with `git archive --prefix`. Verifies each candidate with `tar tzf` before publishing. Byte-identical mirror at skills/do-work/tools/. Exempt from nothing: _dev/tests/prescribed-shell-canonicalization.sh:44-63 fails any skills/*/tools/*.sh that calls curl outside a quoted heredoc. (13 `set -u` only (no -e, no pipefail); 15-18 usage -> stderr, exit 2; 20-24 positionals; 26-37 fetch_cleanup on EXIT + `exit 129` HUP / `exit 130` INT / `exit 143` TERM; 39-50 atomic-download.sh probe of ../scripts and ../skills/do-work/scripts; 52-67 branch/prefix/repo-URL derivation from */archive/refs/heads/*.tar.gz; 69-81 HTTP route (75 stdout success line); 83-110 git route (99-107, 105 stdout success line); 112-115 total failure (two stderr lines, exit 1))
- `/home/user/skill-do-work/skills/do-work/tools/do-work-update.sh` — The updater. Project-root/Git-root gate, skill-inside-project gate, awk semver comparison, fetch + extract + validate through the INSTALLED siblings, then delegates the write to the installed installer with DO_WORK_INSTALL_CANCEL_EXIT_STATUS=3 and verifies the installed version afterward. Not mirrored into tools/ (contract-regressions.sh:37 forbids a root tools/do-work-update.sh). (5 DO_WORK_UPSTREAM_URL; 12-26 version_order awk (prints 1/-1/0); 28-30 read_action_version (regex differs from the validator's); 32-37 argv; 39-45 paths; 47-50 skill-inside-project gate; 52-56 local version; 58-62 Git worktree root; 64-70 temp dirs then EXIT trap (trap set AFTER mkdir); 72-78 fetch + extract; 80-89 validate + remote version; 91-97 version comparison (0 up-to-date exit 0; -1 fail); 99-108 delegate with cancel status 3; 109-113 post-update verification; 115 success line)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli.sh` — The REQ-406 launcher. Rebuilds when the binary is missing or when any *.go/go.mod/go.sum under the module is `-newer` the binary; enforces a Go 1.26.1 floor via awk; builds to a mktemp inside the module dir then `mv -f`; `exec`s the binary. Its own traps are 129/130/143. (7 minimum_go_version=1.26.1; 9-21 version_at_least awk; 23-28 staleness check via `find ... -newer "$binary_path" -print -quit`; 30-43 Go presence/version refusal (exit 2); 45-62 temp build + chmod + mv; 65 exec)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` — CommandResult, findings, RecordedChange, SkippedWork, RollbackResult/RollbackStatus, ExitCode (the sole outcome->number authority), NormalizeResult, RenderResult, renderText, joinArgv. (19-28 CommandOutcome constants; 70-76 RollbackStatus (3 constants; the empty string is a de-facto 4th wire value); 84-93 CommandResult; 97-112 ExitCode (0/1/1/3/4/2, default 2); 114-150 NormalizeResult; 152-166 RenderResult; 168-210 renderText (194-196 changes, 197-199 skipped, 200-208 rollback); 212-225 joinArgv)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` — ExecutionContext, CommandHandler, CommandRuntime.Run, writeResult, parseGlobalOptions (--repo-root, --format, both space- and =-separated), usageFinding, absolutePath. Handlers write nothing themselves; the runtime renders to the writer given to NewRuntime (os.Stdout in main). (13-18 types; 35-60 Run; 62-78 writeResult; 80-133 parseGlobalOptions (88-91 first non-dash argument becomes the command; 131 MISSING-COMMAND); 135-145 usageFinding — 142-143 emit the non-runnable `<command>` placeholder; 147-153 absolutePath)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` — FailureKind constants, TransactionResult/Options, MutationRecorder, ExecuteTransaction, preflight/dirty checks, rollbackFailure, committedRisk, runGit. (19-28 FailureKind constants; 95-202 ExecuteTransaction; 121-132 per-target dirty/tracked preflight; 145-147 dry-run early return; 155-166 changed-vs-recorded check (the fold note's `state.existed` gap: a created-but-only-RecordTouched path still reports success); 167-169 non-commit success return; 170-201 commit + post-commit verification; 269-285 inspectTargets (279 refuses a directory target: "declare exact files instead"); 330-380 rollbackFailure; 408-419 runGit)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go` — failureTemplates table, FindingCode (derived from the kind), BuildCommandResult, buildFinding, survivingChanges, failurePaths, gitPathArgv. (30-102 failureTemplates (36-37 and 72, 81 emit the `<command>` placeholder); 106-108 FindingCode; 113-124 BuildCommandResult (takes only a TransactionResult today — no command name); 126-153 buildFinding; 157-178 survivingChanges)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` — 12 lines. `commandruntime.NewRuntime(os.Stdout, nil)` then `os.Exit(runtime.Run(os.Args[1:]))`. No command is registered. (9-12)
- `/home/user/skill-do-work/skills/do-work-board/justfile.template` — The Justfile template. The ENTIRE file is one managed section: line 1 is the begin marker and line 26 is the end marker, so the awk extraction at install-do-work-suite.sh:258-263 yields a byte-identical copy and a fresh install's justfile equals this file byte for byte (asserted at install-suite-behavior.sh:238-240). Defines the five reserved recipes: run-kanban, run-kanban-cli, kanban-static, kanban-summary, run-do-work-update. Mode 0644. (1 `# >>> do-work:recipes >>>`; 5 run-kanban; 11 run-kanban-cli; 15 kanban-static; 20 kanban-summary; 24-25 run-do-work-update; 26 `# <<< do-work:recipes <<<`)
- `/home/user/skill-do-work/skills/do-work/agent-instructions.template.md` — 5 lines, 274 bytes, mode 0644. The entire file is one managed section bounded by the HTML-comment markers, so it doubles as the section file AND the create-from-template file in install-do-work-suite.sh:312-314. Line 4 carries the `crew-members/communication-style.md` link the installer greps for. (1 `<!-- >>> do-work:communication-style >>> -->`; 5 `<!-- <<< do-work:communication-style <<< -->`)
- `/home/user/skill-do-work/skills/do-work/hooks/hooks.json` — The hook fragment merged into .claude/settings.json. One event (SessionStart) with one wrapper containing one command object: `bash "${CLAUDE_PROJECT_DIR:-.}/.claude/skills/do-work/hooks/session-start.sh"`. Contains NO pipeline-guard entry (contract-regressions.sh:454-457). (1-14)
- `/home/user/skill-do-work/suite/modules.tsv` — Header `source<TAB>destination` plus exactly four rows mapping skills/do-work{,-board,-knowledge,-toolbox} to .claude/skills/<same>. LF terminated, no CR. (1-5)
- `/home/user/skill-do-work/skills/do-work/scripts/atomic-download.sh` — The HTTP primitive route 1 delegates to. mktemp beside the target, curl with retry flags and optional GH_TOKEN/GITHUB_TOKEN bearer, mv on success, and a nested-download guard for a directory-occupied target. Exits 2 on usage/alloc/publish failure, the curl status on download failure, 1 on the directory case. Must keep working from the Go install command unchanged. (5-8 usage; 14-19 cleanup trap on EXIT HUP INT TERM; 21-24 mktemp beside target; 26-34 opt-in credentials via `set --`; 36-43 curl; 44-59 publish + nested-download guard)
- `/home/user/skill-do-work/_dev/tests/install-suite-behavior.sh` — 890 lines. The installer's whole behavioural suite: archive construction, version-mismatch rejection, bootstrap parity with README, fresh install, reinstall with custom Just/settings/CLAUDE.md, marker-free append, BOM+CRLF reserved-recipe collision under a restricted PATH, delimiter-bearing collision, invalid Just/settings, corrupt archives, the Python-fallback lane, the no-JSON-tool manual lane, post-write Just and settings rollback lanes, the case-preserving Justfile recovery lane, TERM-during-write, non-Git roots, cancellation, three dirty-state recovery shapes, the unstage-loop failure lane, and the confirmed-discard success lane. (23-31 managed_state_paths; 33-80 snapshot/assert_install_state_unchanged; 82-88 new_git_project; 90-96 run_installer (pipes `y`, captures 2>&1); 144-158 archive build (`cp -R $repo_root/skills`); 163-183 version-mismatch; 185-197 bootstrap == README block; 199-260 fresh install through the bootstrap with a curl stub; 262-362 reinstall + idempotence (338-354 the python3 settings-structure assertion); 364-375 marker-free append; 377-417 restricted-PATH collision (380 the command list); 419-444 delimiter-bearing collision; 446-469 invalid Just / invalid settings; 471-501 corrupt archives; 503-549 python3 fallback lane (506 command list, 536 `settings reconciler: python3`); 551-600 no-JSON-tool manual lane (554 command list, 592 the exact manual-instruction regex, 594 requires it twice); 602-657 flaky-`just` post-write rollback + case-recovery; 659-694 flaky-`jq` post-write settings rollback; 696-734 TERM-during-cp; 736-745 non-Git roots; 747-761 cancellation; 763-805 three dirty-state recoveries; 807-844 unstage-loop failure; 846-864 dirty cancellation; 866-884 confirmed discard)
- `/home/user/skill-do-work/_dev/tests/update-script-behavior.sh` — 631 lines. Builds a synthetic installed suite and a synthetic upstream tree, stubs curl, and drives the real updater. Covers the four-module update, Just-recipe entry-point parity, the trusted-engine guarantee against a hostile archive, malformed/traversing manifests, a client-side destination symlink, dirty-consent cancel/confirm, mid-copy failure recovery, and the whole fetcher contract (rate-limited fallback, export-ignore, branch selection, signals, total failure). (8-13 script paths; 114-145 build_suite_install (the 'existing project' simulator); 147-187 build_suite_tree (the upstream simulator); 201-215 curl stub; 217-225 run_updater (captures 2>&1); 232-294 four-module update; 296-313 Just entry-point parity; 315-338 hostile-archive trusted-engine case; 340-374 malformed/traversal manifests; 376-389 destination symlink; 391-412 dirty cancel/confirm; 414-446 mid-copy failure recovery; 448-606 fetcher cases (487-489 stdout-only capture of the winning route, 554-576 HUP/INT/TERM 129/130/143, 595-597 stderr-only capture of the total-failure report); 608-619 the two shell-grep assertions; 621-624 version.md delegation grep)
- `/home/user/skill-do-work/_dev/tests/contract-regressions.sh` — The aggregate suite maintainer-verify runs. Holds every replace-text-section behavioural probe plus the grep-level assertions that pin implementation strings inside the five scripts, and invokes every sibling probe file. (429-437 two file_not_contains on updater/installer; 4569-4592 five updater/installer assertions; 4796-4799 module_relatives assertion; 4823 installer on the maintainer-doc-mention allowlist; 6438-6445 suite-manifest probe; 6568-6575 update-script probe; 6584-6591 do-work-cli launcher probe; 6608-6615 install-suite probe; 6644-7157 the replace-text-section probe block (6655-6695 replace/create/mode/idempotence, 6697-6709 filename variants, 6711-6753 five reserved-name collisions incl. a CRLF one, 6755-6807 six multiline-default header forms, 6809-6935 four accept-cases, 6937-7089 four EXACT whole-output comparisons at 6966/6971, 7012/7017, 7043/7048, 7078/7083, 7091-7119 non-collisions and external-after-managed, 7121-7130 retired flag, 7132-7149 four malformed-marker cases, 7151-7155 just parse); 7163-7205 python3 justfile-casing replay; 7207-7234 late assertions incl. 7227-7234 the two `dir=parent` / `os.replace(...)` implementation greps)
- `/home/user/skill-do-work/_dev/tests/maintainer-verify.sh` — The canonical gate. Floors: Go 1.26.1, ShellCheck 0.11.0. Runs ShellCheck over every tracked *.sh, gofmt -l over every tracked *.go, contract-regressions.sh, then queue-kanban / do-work-cli / audit-metrics vet+test. Its --self-test asserts EXACTLY 11 stages (12 with Node), so adding a lane requires editing the self-test shim and the expected count together. (10-11 floors; 53-193 write_command_shim (110-119 the do-work-cli stage arms); 209-243 assert_success_stages (expected_count 11/12); 221-223 the stage list; 344-347 the failure-injection stage list; 406-411 required commands; 446-458 ShellCheck lane; 460-477 gofmt lane; 479-480 aggregate; 527-536 do-work-cli vet + `go test -count=1 ./...`)
- `/home/user/skill-do-work/_dev/tests/staged-skills-contract.sh` — Package boundary contract. require_file lists the three mirrored tools under skills/do-work/tools/ AND the three under tools/, and derives the mirror byte-identity check from tools/*.sh on disk so a fourth mirrored tool is covered automatically. (12-18 require_file; 44-50 core file list incl. the three mirrored tools + do-work-update.sh; 85-91 the retained root bootstrap tools; 809-820 the derived mirror byte-identity check)
- `/home/user/skill-do-work/_dev/tests/prescribed-shell-canonicalization.sh` — Requires the listed prescribed scripts to be executable (including skills/do-work/tools/fetch-upstream-archive.sh at line 23) and fails ANY skills/*/tools/*.sh that invokes curl outside a quoted heredoc. The installer's BOOTSTRAP heredoc is legal only because it is a quoted heredoc. (9-29 executable list; 31-63 the direct-curl scanner (its awk only recognizes `<<'DELIM'` / `<<-'DELIM'` quoted heredocs); 65-87 required headings in skills/do-work/docs/prescribed-shell-primitives.md)
- `/home/user/skill-do-work/_dev/tests/suite-manifest-contract.sh` — 179 lines. Drives tools/validate-suite-manifest.sh over a canonical fixture plus 16 rejection cases, and asserts the validator does not modify the fixture (cksum before/after). (26-42 make_valid_fixture; 52-76 run/expect helpers; 83-91 canonical fixture + non-mutation check; 93-173 the rejection cases)
- `/home/user/skill-do-work/skills/do-work/docs/prescribed-shell-primitives.md` — Shipped canonical rationale. Line 23 asserts `tools/install-do-work-suite.sh` remains self-contained because it is the bootstrap that installs these packages — false once it delegates to do-work-cli. (5-23 shipped executable homes table + the self-contained claim at 23)
- `/home/user/skill-do-work/skills/do-work/tools/prime-do-work-update.md` — The REQ's second named prime. Its Read-first list names the four shell implementations; its Stakes paragraph states the updater's guard contract and names _dev/tests/update-script-behavior.sh as the holder of current behaviour. (5-10 Read first; 12-16 Do not edit; 18-24 Stakes; 26 Lessons pointer to lessons-do-work-update.md)
- `/home/user/skill-do-work/README.md` — The Installation section's bash fence must stay byte-identical to `--print-bootstrap-command` (asserted at install-suite-behavior.sh:189-197). No Go prerequisite is documented anywhere in the file today. (Installation section fence lines 6-21 (the same 13 lines as the installer heredoc); Updating paragraph; 'Upgrade an existing installation with an AI agent' block)
- `/home/user/skill-do-work/skills/do-work/actions/version.md` — The canonical agent-driven update contract. Line 5 carries the `**Current version**:` line the validator and updater both parse. Lines 37-58 forbid duplicating the engine and name DO_WORK_UPSTREAM_URL. Serial-owner file per CLAUDE.md; the plan excludes it from the write set. (5 Current version; 37-48 engine delegation; 50-56 DO_WORK_UPSTREAM_URL escape hatch; 85, 94 the update-mode checklist items)
- `/home/user/skill-do-work/.gitattributes` — export-ignore list. `/.gitignore` is ROOT-ANCHORED, so skills/do-work/tools/do-work-cli/.gitignore ships. /_dev, /do-work, /kb, /decisions, /CLAUDE.md, /AGENTS.md, /ai-reports, /.vscode, /CHANGELOG-20*.md do not ship. (43-53)
- `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/.gitignore` — `/do-work-cli` and `/do-work-cli.build.*`. Tracked, not export-ignored, so it ships and keeps a consumer's built binary out of `git status` — which is what stops the installer's own dirty-managed check from flagging the build output. (1-2)

**Current behaviour that must be preserved byte-for-byte:**

MARKER FORMAT. Default pair (replace-text-section.sh:21-22): BEGIN `# >>> do-work:recipes >>>`, END `# <<< do-work:recipes <<<`. The installer overrides them for CLAUDE.md (install-do-work-suite.sh:98-99): `<!-- >>> do-work:communication-style >>> -->` / `<!-- <<< do-work:communication-style <<< -->`. Override validation (:287-299): both flags together or neither; each must be non-empty and contain no \n or \r; they must differ; the values are UTF-8 encoded verbatim.

MARKER SPAN ALGORITHM (marker_span, :62-81), verified empirically:
1. `data.splitlines(keepends=True)`. For bytes this splits ONLY on \n, \r and \r\n — confirmed: b'a\\vb\\fc\\x1cd\\x1ee\\x85f' is ONE line, b'a\\rb\\nc\\r\\nd' is four. A LONE CR therefore ends a line. Byte offsets are accumulated cumulatively, so the span is exact byte positions, not line numbers.
2. Comparison body = `line.rstrip(b"\r\n")`, which strips a SET of trailing CR/LF bytes, not a suffix (b'line\\n\\r\\n\\r' -> b'line'). A marker terminated by \r\n or by a bare \r matches the same constant as one terminated by \n.
3. No begin AND no end -> None (means "no managed section", not an error).
4. Anything other than exactly one of each -> die `<label> must contain exactly one begin marker and one end marker`. This covers duplicated begin, begin-only, end-only.
5. begin_index >= end_index -> die `<label> has reversed or nested managed markers`.
6. require_section_only (section file only) and (begin != 0 or end != last line) -> die `<label> must contain only the complete managed section`.
7. span = (offset[begin], offset[end] + len(lines[end])) — the end marker's own terminator is INSIDE the span.

REPLACE vs APPEND (:339-350). Span found -> `data[:start] + section + data[end:]`. Note the consequence: a CRLF- or lone-CR-terminated marker region is replaced with the section file's LF bytes, while surrounding lines keep their original terminators. Verified: `custom:\r<CR-markers>\rtail\r` came back as `custom:\r` + LF section + `tail\r`. No span -> append with a separator chosen by four cases, verified: target empty -> ``; target ends `\n\n` -> ``; target ends `\n` -> `\n`; otherwise -> `\n\n`.

CONDITIONAL WRITE (:352-353). If replacement == target, NOTHING is written — verified same inode and same mtime across a repeat run. This is what makes reinstall byte-idempotent.

ATOMIC REPLACE (:221-260). Parent = dirname(abspath(target)); must be an existing directory else die `target parent directory does not exist: <parent>`. `tempfile.mkstemp(prefix=".<basename>.", suffix=".tmp", dir=<parent>)` — the temp is in the TARGET's directory (pinned by contract-regressions.sh:7227-7230). Write, flush, os.fsync, then `os.fchmod(fd, mode)`, then `os.replace(temp, path)` (pinned by contract-regressions.sh:7231-7234), then a best-effort directory fsync whose OSError is deliberately swallowed. On OSError -> die `atomic replacement failed for <path>: <err>`; the finally block unlinks any surviving temp.

FILE MODES. Replace path (:321): mode = `stat.S_IMODE(os.stat(target).st_mode)` of the EXISTING target — verified 640 preserved. Create path (:314-315): mode = S_IMODE of the TEMPLATE file — verified 750 preserved. The installer then propagates modes through `cp -p` onto an already-existing mktemp destination, which GNU cp does honour (verified: mktemp 600 became 640 after `cp -p` from a 640 source).

SYMLINKS. read_regular (:32-45) refuses `os.path.islink` first with `<label> must not be a symlink: <path>`, then non-regular with `<label> must be a regular file: <path>` (a directory target hits this). Target existence is probed as `os.path.exists(target) or os.path.islink(target)` (:306), so a DANGLING symlink counts as existing and then dies at :318-319 `target must not be a symlink: <path>`. Installer-level symlink guards: managed destination must not be a symlink (:207), must be a directory when it exists (:208-210), its nearest existing parent resolved with `pwd -P` must stay under the project root (:211-222); every module source must be a real directory (:226-227) and must contain no symlink at all — `find "$source_path" -type l -print -quit | grep -q .` (:228-230). Justfile target must be a regular non-symlink file (:245-246); CLAUDE.md (:299-305) and settings.json (:335-341) the same. Git index must be a regular non-symlink file (:543-548). validate-suite-manifest.sh refuses symlinked VERSION (:20), modules.tsv (:28), module roots (:106), SKILL.md (:113), core VERSION (:126), actions dir (:138), version.md (:140), and resolves module roots physically to reject a symlink escape (:108-112).

NUL BYTES AND BOM. NUL is fully transparent to the Python path: the byte-preservation fixture at contract-regressions.sh:6661 puts `prefix\000byte` before the section and requires an exact cmp afterwards. BOM is handled in exactly ONE place: just_definition_names (:193-198) strips a leading \xef\xbb\xbf from the classification view of line index 0 only; the raw bytes are untouched and marker_span does NOT strip it. Verified consequence: a target whose first line is `<BOM># >>> do-work:recipes >>>` fails with `must contain exactly one begin marker and one end marker`, because the begin marker no longer compares equal while the end marker still does. That asymmetry is deliberate (REQ-173) and is what install-suite-behavior.sh:377-417 exercises via a BOM+CRLF `run-kanban:` collision.

JUST DEFINITION SCANNER (:84-218). Per line: skip empty or leading space/tab/`#`. Alias regex `alias[ \t]+([A-Za-z_][A-Za-z0-9_-]*)[ \t]*:=` wins first. Otherwise an optional leading `@`, then `[A-Za-z_][A-Za-z0-9_-]*`; scan the remainder skipping quoted regions; the first bare `:` that is not `:=` makes it a recipe, a `:=` makes it an assignment (None). Opening delimiters recognised in priority order: `'''`, `"""`, ``` ``` ``` (rejected when adjacent to another backtick on either side), then a single `"`, `'` or `` ` ``. Closing: same token, with the backtick-adjacency rule for ``` ``` ```, plus an odd-backslash escape rule that applies ONLY to `"` and `"""` (raw `'`/`'''` and backticks are not escapable). Multiline state carries across lines; a line that starts OUTSIDE a literal and begins with space/tab/`#` or is empty resets the state to None, which is how recipe bodies and comments cannot open a literal; inside a scan a bare `#` outside a literal ends the line. When a literal spans lines, the accumulated lines are joined and classified once as a single definition source (:199-217).

COLLISION REJECTION (:324-337). reserved = definitions in the section file; empty -> die `section file defines no Just recipes or aliases for collision validation`. unmanaged target = target bytes minus the managed span (or all of it when there is no span). collisions = sorted(unmanaged ∩ reserved); non-empty -> die `target defines reserved Just recipe or alias outside managed section: ` + `", ".join(...)`. The target is NOT written when this fires.

STREAMS AND EXIT CODES TODAY.
- replace-text-section: success is silent, exit 0. `--help` writes USAGE to STDOUT, exit 0. Every other error writes `replace-text-section: <message>` to STDERR, exit 1.
- validate-suite-manifest: success writes `suite manifest valid: v<version> (<n> modules)` to STDOUT, exit 0. Every failure writes `suite manifest: <message>` to STDERR, exit 1. Usage error uses the same shape.
- fetch-upstream-archive: success writes `upstream archive fetched over HTTP` or `upstream archive fetched with git (HTTP route <outcome>)` to STDOUT, exit 0. Total failure writes two STDERR lines naming both route outcomes and DO_WORK_UPSTREAM_URL, exit 1. Usage -> STDERR, exit 2. HUP 129, INT 130, TERM 143 (three distinct statuses).
- install-do-work-suite: STDOUT carries the plan preamble (`Ready to install do-work suite vX into <root>:`, the four module relatives, `  Justfile: <rel>`, `  agent instructions: <rel>`, `  settings reconciler: <jq|python3|manual>`), `Reviewing the complete managed install before overwrite:`, the whole review diff, the prompt `Install this complete four-skill suite? [y/N] `, the cancel line `Installation cancelled; no files were changed.`, the success line `Installed do-work suite vX with four verified modules.`, and the manual instruction in the manual case. STDERR carries every `fail` (`do-work suite install: <msg>`), the dirty-managed warning, and the two recovery lines. Exit 0 success; 1 any failure and any incomplete recovery; `$DO_WORK_INSTALL_CANCEL_EXIT_STATUS` (default 0) on decline; 130 for HUP, INT and TERM alike (`trap 'exit 130' HUP INT TERM` at :173 — note this does NOT match the fetcher's 129/130/143).
- do-work-update: STDOUT `Checking do-work updates…` (UTF-8 ellipsis), `You're up to date (vX)`, `Update available: vR (you have vL), archive layout: four-module suite.`, `Update cancelled; no files were changed.`, `Updated to vR at <root> using the four-module suite.`. STDERR `do-work update: <msg>`, exit 1.

INSTALL ORDER OF OPERATIONS (all guards before the single confirmation, all writes after it).
Pre-confirmation: argv -> project root exists and is the Git worktree root -> git index path resolved -> sibling tools located -> install_tmp created (traps armed at :172-173) -> archive supplied or fetched -> `tar xzf <archive> -C <source_root> --strip-components=1` -> `bash validate-suite-manifest.sh --root <source_root>` -> suite_version from `sed -n '1p' VERSION` -> module plan read from modules.tsv (header row skipped by `[ "$module_source" = source ] && continue`) -> destination guards -> source symlink guards -> board template / hooks.json / instructions template presence -> Justfile discovery in order justfile, Justfile, .justfile via `find -mindepth 1 -maxdepth 1 -name <n> -print -quit` (first hit wins; absent -> `$project_root/justfile`, just_existed=0) -> awk extracts the managed section from the board template -> Justfile candidate built in install_tmp -> candidate marker count == 1 each, five reserved recipe names present, `just --justfile <candidate> --list` when just is on PATH -> CLAUDE.md candidate built and validated -> settings candidate composed and validated -> plan preamble -> review diff assembled with `diff -ruN <dest> <source>` per module and `diff -u` for each configuration file (a status above 1 is a hard failure) -> dirty-managed warning -> ONE confirmation accepting y|Y|yes|YES.
Post-confirmation: backups into install_tmp/originals (`cp -Rp <dest>/. ...` per existing module, `cp -p` for justfile/settings/CLAUDE.md, `cp -p` for the Git index) -> `write_started=1` -> unstage loop (`git ls-files -- <rel>` non-empty then `git restore --staged -- <rel>`, module paths only) -> module writes `rm -rf dest; mkdir -p dest; cp -Rp <src>/. <dest>/` -> Justfile via mktemp `.do-work-just.install.XXXXXX` in the target's directory + `cp -p` + `mv` -> settings via mktemp `.settings.json.install.XXXXXX` + `cp -p` + `mv` -> CLAUDE.md via mktemp `.do-work-instructions.install.XXXXXX` at the project root + `cp -p` + `mv` -> verification: `diff -qr src dest` per module, `cmp -s` for justfile and CLAUDE.md, `just --list` on the installed Justfile when just exists, `cmp -s` plus a JSON re-parse for settings, installed VERSION == suite_version, installed actions/version.md has exactly one Current-version line and it equals suite_version -> `install_verified=1` -> success line.

ROLLBACK POINTS. `write_started` is set at :553, BEFORE the unstage loop, so an unstage failure recovers. `install_verified` is set at :617, AFTER every verification, so ANY verification failure recovers. cleanup() (:161-171) unhooks EXIT, IGNORES HUP/INT/TERM for the duration of recovery, runs recover_install when write_started and not install_verified, removes install_tmp, and exits with the original status (forced to 1 if recovery reported a problem). recover_install (:102-159) restores in a fixed order — modules (rm -rf then mkdir -p + `cp -Rp backup/.`), justfile, settings, CLAUDE.md, then the Git index (mktemp in the index's own directory, `cp -p`, `mv`; or `rm -f` the index when none existed) — under `set +e`, accumulating a single recovery_failed flag. Success message: `do-work suite install: restored every managed path and the Git index to their exact pre-install state.` Failure message names the four skill directories plus the three configuration paths and the index path, and returns 1.

SETTINGS RECONCILIATION, current semantics (identical intent in both branches): the settings root must be an object; hooks must be an object; hooks.Stop must be an array if present. Inside Stop, only entries that are objects with an array `hooks` are inspected; from those, remove every hook object whose `command` is a STRING CONTAINING `.claude/skills/do-work/hooks/pipeline-guard.sh`; if that removal emptied the wrapper, drop the wrapper entirely; if it did not empty it, keep the wrapper with the reduced list; a wrapper that was already empty and lost nothing is PRESERVED (the `{"matcher":"preserve-empty","hooks":[]}` fixture). If Stop ends up empty, delete the Stop key. Then for each event in the fragment, create a missing array and append each entry that is not already deep-equal to an existing one. Output is 2-space indented with a trailing newline. VERIFIED DIVERGENCE the port must choose between: jq emits raw UTF-8 (`"héllo→"`), Python's default ensure_ascii emits `"h\\u00e9llo\\u2192"`. Both preserve document key order. jq iterates fragment events through `keys` (SORTED); Python iterates insertion order — identical today with one event. jq is preferred when present, so raw UTF-8 is the incumbent behaviour.

FETCHER SEMANTICS. Branch derivation applies ONLY to the `*/archive/refs/heads/*.tar.gz` shape: branch = the text after the marker minus `.tar.gz`; base = the text before; `archive_prefix` = `<last segment of base>-<branch>/`; a missing repo URL is derived as `<base>.git`. Anything else keeps `archive_prefix='upstream-main/'`, derives no repo URL, and the git route reports `unavailable (no repository URL supplied and none derivable from the tarball URL)`. The clone uses `--depth 1 --quiet` plus `--single-branch --branch <b>` when a branch was derived, so a named branch is selected exactly and a missing one FAILS rather than silently substituting HEAD (REQ-424). The git route stages into `mktemp "${archive_target_path}.fetching.XXXXXX"` and only `mv`s after `[ -s stage ]` and `tar tzf stage` pass, so a pre-existing target survives total failure untouched and no scratch is left behind.

VERSION PARSING DIVERGENCE that must be preserved or deliberately unified: validate-suite-manifest.sh:142-145 requires EXACTLY ONE line matching `^\*\*Current version\*\*:` and extracts with `sed -n 's/^\*\*Current version\*\*:[[:space:]]*//p'` (the whole remainder of the line). do-work-update.sh:29 extracts with `sed -n 's/^\*\*Current version\*\*: *\([0-9][0-9.]*\).*/\1/p' | head -n 1` (leading digits only, first match). VERSION files must be exactly one newline-terminated semver line, enforced by `printf '%s\n' "$v" | cmp -s - <file>`.

**Python and jq elimination list (23 sites):**

- ELIMINATION LIST — python3, install/update path. tools/replace-text-section.sh:5-8 — `command -v python3` guard; on absence prints `replace-text-section: python3 is required` to stderr and exits 1. Mirror: skills/do-work/tools/replace-text-section.sh:5-8.
- tools/replace-text-section.sh:10 — `exec python3 - "$@" <<'PY'`; the entire 344-line implementation (marker span, Just scanner, atomic replace, argv parsing) is this one heredoc. Mirror: skills/do-work/tools/replace-text-section.sh:10.
- ELIMINATION LIST — python3, installer. tools/install-do-work-suite.sh:267 — `if command -v python3` gates whether the Justfile candidate goes through replace-text-section.sh at all. Present: existing target reconciled with --reject-recipe-collisions (:268-272), absent target created from the board template (:273-277). Mirror: skills/do-work/tools/install-do-work-suite.sh:267.
- tools/install-do-work-suite.sh:278-282 — the no-python3 Justfile branch: refuses outright when a Justfile already exists (`python3 is required to reconcile an existing Justfile safely; no client files were changed`, line 280) and otherwise just `cp -p`s the board template. This is one half of the REQ's RED. Mirror line 278-282.
- tools/install-do-work-suite.sh:308 — `if command -v python3` gates the CLAUDE.md managed-section reconciliation through replace-text-section.sh with the communication-style markers (:309-315). Mirror line 308.
- tools/install-do-work-suite.sh:316-320 — the no-python3 CLAUDE.md branch: refuses when CLAUDE.md exists (`python3 is required to reconcile an existing CLAUDE.md safely; no client files were changed`, line 318), otherwise copies the template. Mirror line 316-320.
- tools/install-do-work-suite.sh:329-333 — `settings_tool` three-way selection: `jq` if `command -v jq`, else `python3` if `command -v python3`, else the literal `manual` initialised at line 327. Mirror line 329-333.
- tools/install-do-work-suite.sh:356-388 — the jq reconciliation program invoked as `jq --slurpfile fragment "$core_hooks" '...' "$settings_input"`. Computes: coerce/validate hooks to an object; validate Stop to an array; per Stop wrapper that is an object with an array `hooks`, drop every hook object whose `.command` is a string containing `.claude/skills/do-work/hooks/pipeline-guard.sh`, drop the wrapper entirely only when that removal emptied it, keep an already-empty wrapper; delete Stop when it ends empty; then for each fragment event (iterated via `keys`, i.e. SORTED) create a missing array and append each entry not already present by deep value equality. Mirror line 356-388.
- tools/install-do-work-suite.sh:389-390 — `jq -e . "$settings_candidate"` re-parses the composed candidate; failure aborts before any write. Mirror line 389-390.
- tools/install-do-work-suite.sh:392-448 — the embedded python3 settings composer, invoked as `python3 - "$settings_input" "$core_hooks" "$settings_candidate"`. Same semantics as the jq program (retired-guard removal at 416-424, emptied-wrapper drop at 425-431, Stop deletion at 434-437, fragment append-unique at 438-444) and writes with `json.dump(settings, handle, indent=2)` plus an explicit trailing newline (445-447). Iterates fragment events in insertion order, and escapes non-ASCII as \uXXXX. Mirror line 392-448.
- tools/install-do-work-suite.sh:449-450 — `python3 -m json.tool "$settings_candidate"` re-parses the composed candidate in the Python branch. Mirror line 449-450.
- tools/install-do-work-suite.sh:600-604 — post-write validation of the INSTALLED settings file: `jq -e . "$settings_target"` in the jq branch (line 601), `python3 -m json.tool "$settings_target"` in the Python branch (line 603). This is the only injection point the flaky-jq rollback lane at _dev/tests/install-suite-behavior.sh:659-694 has. Mirror line 600-604.
- CONSEQUENT DEAD SURFACE once the branches go. tools/install-do-work-suite.sh:6 — `manual_settings_instruction`, the 380-character MANUAL STEP string. Referenced at :494 (into the review diff) and :620 (after success). Its exact wording is pinned by _dev/tests/install-suite-behavior.sh:592 and required to appear TWICE (:594).
- tools/install-do-work-suite.sh:343 / :458 — `if [ "$settings_tool" != manual ]` opens and closes the whole compose block; :493-503 branches the review diff; :575-580 skips the settings write; :597-605 skips the settings verification; :619-621 prints the manual instruction on success. All become unreachable once Go always reconciles.
- KEEP — python3 outside the install/update path. skills/do-work/tools/checks/preflight.sh:120-130 — writes do-work/.../baseline.json with `python3 - ... <<'PYEOF'`, with a printf fallback on the `||` branch. REQ-414 territory per the plan; not part of install/update.
- KEEP — explicitly named by the REQ. skills/do-work/tools/checks/preflight.sh:143 — `python3 -c "import sys; sys.exit(0 if sys.prefix != sys.base_prefix else 1)"`, the virtualenv probe for a Python TARGET project. The REQ says: keep Python checks only when probing a Python target capability.
- KEEP — outside scope. skills/do-work-knowledge/scripts/install-memory-hooks.sh:20 (jq required guard), :28 (`jq .` validation of the existing settings), :46 (the jq merge program), :51 (`jq .` validation of the merge output). Knowledge memory hooks, REQ-417 territory.
- KEEP — outside scope. skills/do-work-knowledge/hooks/memory-stop-capture.sh:49-50 (jq transcript_path extraction with a sed fallback), :98, :112 (`jq -rs` capture program), :163-175 (the documented no-jq fallback). REQ-415 territory.
- KEEP — outside scope. skills/do-work-toolbox/scripts/install-last30days.sh:159-169 — probes python3.13/python3.12/python3/python for a Python 3.12+ target capability.
- KEEP — prose, outside scope. skills/do-work-knowledge/actions/setup-memory.md:57,59,64 and skills/do-work-knowledge/actions/memory-reference.md:157,158 describe the jq-or-manual contract for memory hooks.
- HARNESS-SIDE python3 (maintainer tests, not shipped, and not what the REQ's 'Python absent' constraint targets): _dev/tests/contract-regressions.sh at 112, 188, 249, 459, 487, 1648, 1753, 2149, 2253, 2296, 2405, 2642, 3027, 3169, 5293, 5728, 7163; _dev/tests/staged-skills-contract.sh at 208, 277, 322, 882, 987, 1059; _dev/tests/shipped-package-reference-contract.sh:6.
- HARNESS-SIDE python3 that IS in the installer suite and must be retargeted: _dev/tests/install-suite-behavior.sh:338-354 (asserts the reinstalled settings structure, including the exact expected hooks.Stop list) and :538-548 (the same assertion for the Python-fallback lane). These run in the outer harness process, not under a restricted PATH.
- RESTRICTED-PATH COMMAND LISTS that encode today's dependency set: _dev/tests/install-suite-behavior.sh:380 and :506 both read `awk bash cat chmod cmp cp diff dirname find git grep gzip head mkdir mktemp mv python3 rm sed stat tar tr wc` (python3 present, jq and just absent); :554 reads the same list MINUS python3 (the no-JSON-tool lane). None of the three includes `go`, `curl` or `just`.

**Test conventions:**

HOW THE SUITES RUN. Nothing auto-discovers _dev/tests/*.sh. `_dev/tests/contract-regressions.sh` is the aggregate and invokes each probe explicitly with `bash "$probe"` inside an `if [ ! -f ... ]; then FAIL; elif ! bash ...; then FAIL; fi` shape: suite-manifest-contract.sh at 6438-6445, shipped-package-reference-contract.sh at 6451-6458, action-shell-blocks.sh at 6463-6473, session-start-hook-behavior.sh at 6477-6484, prescribed-shell-canonicalization.sh at 6488-6495, defensive-surface-audit.sh at 6539-6546, record-commit-hash-guards.sh at 6552-6562, update-script-behavior.sh at 6568-6575, do-work-cli-launcher-behavior.sh at 6584-6591, staged-skills-contract.sh at 6596-6603, install-suite-behavior.sh at 6608-6615, p50-estimator-determinism.sh at 6620-6627, select-simple-reqs-behavior.sh at 6635-6642. `_dev/tests/maintainer-verify.sh` runs contract-regressions.sh as its single "aggregate" stage (line 480), then queue-kanban vet+test (+ optional Node and browser lanes), do-work-cli vet+test (527-536), and audit-metrics vet+test. Baseline confirmed green on this tree: install-suite-behavior.sh exits 0 in ~16s, update-script-behavior.sh in ~4s, and `go vet ./... && go test -count=1 ./...` in the CLI module passes.

MAINTAINER-VERIFY SELF-TEST IS A COUNTING GATE. `maintainer-verify.sh --self-test` builds a fake toolchain and asserts EXACTLY 11 recorded stages (12 with Node) at :209-243, with the stage names listed at :221-223 and the failure-injection sweep at :344-347. Its `go` shim (:110-119) only accepts `vet ./...` and `test -count=1 ./...` from `*/skills/do-work/tools/do-work-cli`. Adding a verification lane means editing the shim, the expected list and the count in the same change; the plan's "no lane needs adding" is the cheap path.

INSTALLER SUITE: FRESH vs EXISTING. Fresh = `new_git_project` (install-suite-behavior.sh:82-88): `git init -q` plus user.email/user.name, and nothing else — no justfile, no .claude, no CLAUDE.md. Existing = the same project after one install, then hand-written custom content dropped on top: a justfile with `custom-before:` / stale managed section / `custom-after:` at 640 (lines 263-276), a settings.json with a `custom` key, two SessionStart wrappers, and a Stop array carrying the retired pipeline guard beside a same-wrapper custom hook, a guard-only wrapper, a deliberately empty `{"matcher":"preserve-empty","hooks":[]}` wrapper and a memory hook, at 600 (277-299), and a CLAUDE.md with prose before and after a stale managed section (300-310). The archive is a real tarball built at 144-158 from `$repo_root/VERSION`, `cp -R $repo_root/suite`, `cp -R $repo_root/skills`, and the three root tools, tarred from a `skill-do-work-main/` parent. `run_installer` (90-96) pipes `y\n` and captures stdout+stderr together into one file.

UPDATER SUITE: FRESH vs EXISTING. `build_suite_install` (update-script-behavior.sh:114-145) is the "existing project" simulator: it hand-builds `.claude/skills/<four modules>` with a SKILL.md and payload.txt each, an `actions/version.md` pinned at `**Current version**: 0.0.1`, copies the five real scripts (updater, validator, installer, replacer, fetcher) plus atomic-download.sh into the installed tree, and plants `do-work/queue/sentinel.txt`, `kb/sentinel.txt`, a `Justfile`, and `{"hooks":{}}`. `build_suite_tree` (147-187) is the upstream simulator at VERSION 0.0.2 with the same five scripts, the real hooks.json / agent-instructions template / justfile.template, and a `new-core.txt` that must appear on success and vanish on rollback. Both helpers must gain the do-work-cli module and launcher for the Go port, and `run_updater` (217-225) captures 2>&1.

STUB/INJECTION POINTS AND WHAT EACH PROVES. curl stub copying a fixture tarball and counting invocations (install 201-217, update 204-215) — proves exactly one download. Flaky `just` that succeeds once then fails (install 612-623) — drives post-write Justfile validation failure and the case-preserving `Justfile` recovery (627-643) and the three dirty-state recoveries (772-805). Flaky `jq` that succeeds twice then fails (install 669-681) — the ONLY injection point for post-write settings verification failure; it disappears with the jq branch. `cp` shim that TERMs its parent after copying into `.claude/skills/do-work-board/` (install 708-717) — proves the signal path runs the same recovery; the update suite has an analogous cp shim that fails once (419-423). `git` shim that fails the SECOND `restore --staged` (install 818-830) — proves an index mutation already made is recovered. `curl` that always exits 22 (update 478-484) — proves the git fallback wins. Signal fixtures overriding `bash` as a function to kill the shell (update 554-576) — pin 129/130/143 exactly.

WHAT SURVIVES A GO PORT AND WHAT DOES NOT. Everything that stubs `cp`, `git`, `just`, `curl` or `tar` keeps working provided the Go command still invokes those binaries. The flaky-`jq` lane loses its injection point. The three restricted-PATH command lists (380, 506, 554) omit `go`, so any Go-backed path under those PATHs needs a pre-built binary copied into the fixture; note the launcher rebuilds whenever any *.go or go.mod is `-newer` the binary, and `cp -R` copies in readdir order, so whether the binary ends up newer than the sources is NOT deterministic (observed both ways on this tree) — the fixture must set mtimes explicitly rather than rely on copy order.

GREP-LEVEL ASSERTIONS THAT PIN IMPLEMENTATION STRINGS (all must be retargeted or retired). contract-regressions.sh: 429-432 (updater must not contain suite-layout-v2/--capabilities/legacy_shipped_paths), 434-437 (installer must not contain --migrate-legacy-do-work or the old memory hook paths), 4569-4572 (updater executable), 4573-4576 (`--project-root`), 4577-4580 (`tar xzf "$upstream_tarball" -C "$fresh_upstream"`), 4581-4584 (`Install this complete four-skill suite`), 4585-4588 (no `cp -R "$skill_root"`), 4589-4592 (`recover_install`), 4796-4799 (`module_relatives`), 7215-7218 (installer names replace-text-section.sh), 7227-7230 (`suffix=.*dir=parent`), 7231-7234 (`os\\.replace\\(temporary_path, path\\)`). update-script-behavior.sh: 28-38 (retired updater tokens), 610-619 (both callers grep for `fetch-upstream-archive.sh` and `DO_WORK_UPSTREAM_URL`), 621-624 (version.md delegates to tools/do-work-update.sh). install-suite-behavior.sh: 117-125 (retired installer tokens).

FOUR EXACT WHOLE-OUTPUT COMPARISONS in contract-regressions.sh that compare `$(cat <output-file>)` against a single literal line beginning `replace-text-section: target defines reserved Just recipe or alias outside managed section: ...` — 6966/6971, 7012/7017, 7043/7048, 7078/7083. These break the moment the tool prints anything else (a CommandResult header, a repository line), which is why the plan restates them as exact matches on the rendered finding line.

GO TEST CONVENTIONS (skills/do-work/tools/do-work-cli). Tests live in the package under test (`package commandruntime`, `package gittransaction`, `package resultmodel`), not `_test` packages. Fixture helpers are lowercase functions taking `*testing.T` first with `t.Helper()`: `newFixtureRepository` (command_runtime_test.go:323-333) uses `t.TempDir()` plus `t.Setenv` for GIT_CONFIG_GLOBAL/GIT_CONFIG_SYSTEM/GIT_TERMINAL_PROMPT, `commitFixture`, `runFixtureGit`, `writeFixtureFile`. Table tests use a `tests := []struct{ name string; ... }` slice with `t.Run(test.name, ...)`. Test names are long declarative sentences (TestExitCodeContractThroughRealGitTransactions, TestEveryDeclaredFailureKindProducesACompleteFinding, TestCompletedRollbackReportsNoSurvivingChanges). Enumerations are read back out of the source rather than restated: `declaredFailureKinds` (transaction_findings_test.go:14-29) regex-scans git_transaction.go. Expected rendered text is produced by the production renderer (`renderedNextLine`, command_runtime_test.go:296-311) rather than hand-joined. Identifier style throughout the module is multi-word and spelled out (`repositoryRoot`, `commandArgs`, `standardOutput`, `unformatted_go_files`), matching CLAUDE.md § Naming Conventions.

**Concerns and traps:**

- C1 The lone-CR line split is real and currently load-bearing. Verified: Python's bytes.splitlines ends a line on a bare \r (and only on \r, \n, \r\n — not \v, \f, \x1c, \x1d, \x1e, \x85). A target written entirely with CR terminators has its markers located correctly today and its CR-terminated span replaced with the section file's LF bytes while the surrounding CRs survive. A Go port that splits on \n alone silently mislocates the span. No existing fixture covers a lone-CR file; the CRLF cases at contract-regressions.sh:6737 and 6835-6850 and install-suite-behavior.sh:387 do not exercise it.
- C2 `rstrip(b"\r\n")` strips a SET, not a suffix. Verified b'line\n\r\n\r' -> b'line'. Combined with splitlines this is equivalent to dropping one terminator in practice, but a naive Go `strings.TrimSuffix(line, "\r\n")` is not the same function.
- C3 BOM handling is deliberately ASYMMETRIC and a port that "fixes" it breaks a shipped test. just_definition_names strips the BOM from the classification view of line 0 only (replace-text-section.sh:193-198); marker_span never strips it. Verified consequence: a BOM before the begin marker yields `must contain exactly one begin marker and one end marker`, because the end marker still matches while the begin marker does not. install-suite-behavior.sh:377-417 depends on the BOM+CRLF file being scanned for `run-kanban` while its markers stay absent.
- C4 JSON key order and non-ASCII escaping. Go's encoding/json sorts map keys, which would silently reorder a consumer's settings.json — the GREEN criterion forbids changing unrelated bytes and no current assertion checks order. Separately, the two incumbent reconcilers DISAGREE on non-ASCII: verified jq emits raw UTF-8 while Python emits \uXXXX. jq is the preferred branch today, so raw UTF-8 is the behaviour to preserve; that means SetEscapeHTML(false) and no \u escaping of ordinary runes.
- C5 The four whole-output equality comparisons (contract-regressions.sh:6966/6971, 7012/7017, 7043/7048, 7078/7083) compare the ENTIRE captured stdout+stderr against one literal line. Any CommandResult header, repository line or trailing rollback line makes all four fail. The eleven grep assertions listed under testConventions plus the two in update-script-behavior.sh:613,616 pin implementation strings inside the five shell files and will pass vacuously or fail loudly depending on how the launchers are written — retarget each at the Go file that now owns the behaviour rather than deleting it.
- C6 The flaky-`jq` post-write lane (install-suite-behavior.sh:659-694) loses its only injection point. It is the second of two independently-injectable post-write recovery proofs; folding its assertions into the flaky-`just` lane keeps the surface coverage but removes the second failure mode.
- C7 `prescribed-shell-canonicalization.sh:44-63` fails ANY skills/*/tools/*.sh that invokes curl outside a QUOTED heredoc, and its awk only recognises `<<'DELIM'` / `<<-'DELIM'`. The installer launcher must keep the BOOTSTRAP heredoc in exactly that form. The same file (line 23) requires skills/do-work/tools/fetch-upstream-archive.sh to remain executable.
- C8 Mirror byte-identity is enforced by derivation, not by a list: staged-skills-contract.sh:812-820 walks tools/*.sh and requires each to equal skills/do-work/tools/<same name>. A new root tools/*.sh automatically joins that check. Conversely contract-regressions.sh:37 and staged-skills-contract.sh:60-83 FORBID a root tools/do-work-update.sh, tools/queue-kanban, tools/checks, tools/prime-do-work-update.md — the updater is core-only.
- C9 Signal-status inconsistency already exists and is asserted. install-do-work-suite.sh:173 maps HUP, INT and TERM all to 130. fetch-upstream-archive.sh:35-37 and do-work-cli.sh:52-54 map them to 129/130/143, and update-script-behavior.sh:554-576 asserts the fetcher's three statuses exactly. Launchers must invoke rather than exec if they need to preserve their own trap behaviour, and must not accidentally unify the installer's 130 with the fetcher's triple.
- C10 The launcher's staleness check is mtime-based and copy-order dependent. `find <module> -type f \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \) -newer <binary> -print -quit` triggers a rebuild whenever any source is strictly newer. `cp -R` copies in readdir order, not alphabetical order — observed on this tree copying the binary LAST (no rebuild), but the ordering is not guaranteed. In a restricted-PATH fixture with no `go`, a spurious rebuild is exit 2. Set mtimes explicitly in fixtures rather than relying on copy order. Same hazard on the recovery path: recover_install restores the backed-up module tree including whatever binary it held.
- C11 The running binary lives inside a directory the install `rm -rf`s. In the update flow the CLI is at .claude/skills/do-work/tools/do-work-cli/do-work-cli, and install-do-work-suite.sh:563 does `rm -rf -- "$destination_path"` on that module. Unlinking a running executable is fine on Linux/macOS, but any code path that re-reads its own file after the module write (os.Executable-relative lookups, for one) sees a deleted path. fetch-upstream-archive.sh already locates atomic-download.sh relative to its own script directory, so an os.Executable-relative equivalent inherits this.
- C12 `skills/do-work/tools/do-work-cli/.gitignore` is what keeps the built binary out of a consumer's `git status`, and therefore out of the installer's own dirty-managed check at install-do-work-suite.sh:507-513. It is tracked and NOT export-ignored (the .gitattributes entry is root-anchored `/.gitignore`), so it ships. Losing it would make every install report its own build output as a dirty managed path.
- C13 Go becomes a hard prerequisite for surfaces that need no toolchain today: `--print-bootstrap-command`, manifest validation and section replacement. README.md documents no Go prerequisite anywhere, and skills/do-work/docs/prescribed-shell-primitives.md:23 explicitly claims the installer `remains self-contained because it is the bootstrap that installs these packages` — that sentence becomes false and must change in the same commit.
- C14 skills/do-work/actions/version.md is the canonical agent-driven update contract (lines 37-58 forbid duplicating the engine; line 5 is the Current-version line the validator parses) but it is a serial-only integrator file per CLAUDE.md § Before Every Commit. The Go-prerequisite line for the update path belongs there and is outside the plan's write set; the integrator has to add it during the version bump.
- C15 Two different Current-version parsers exist: validate-suite-manifest.sh:142-145 requires exactly one matching line and takes the whole remainder; do-work-update.sh:29 takes only leading digits from the first match. Unifying them in Go is easy to do accidentally and changes what a malformed version.md does on each path.
- C16 Cancellation currently exits with `$DO_WORK_INSTALL_CANCEL_EXIT_STATUS` (install-do-work-suite.sh:7, default 0), and its ONLY non-test caller is do-work-update.sh:102, which sets 3 to distinguish cancelled from installed. Deleting the variable is correct only if the update path stops running the installer as a subprocess. Note also the variable's own validation at :8-10 (`invalid cancellation status override`) is a public error path.
- C17 `install-suite-behavior.sh:144-158` builds its archive with `cp -R "$repo_root/skills"`, which copies UNTRACKED files — including the 3 MB built do-work-cli binary if one is present in the maintainer's tree. That makes the fixture archive differ from a real tarball install (which ships source only) in a way that could mask a missing-binary bug. Decide deliberately whether the fixture archive should carry a binary.
- C18 The Justfile discovery order is `justfile, Justfile, .justfile` with `find -mindepth 1 -maxdepth 1 -name <n> -print -quit` (install-do-work-suite.sh:240-251), and the case-recovery lane at install-suite-behavior.sh:627-643 asserts a failed install restores `Justfile` as `Justfile` and does NOT create `justfile`. On a case-insensitive filesystem this only works because the probe records the real directory-entry spelling. REQ-180's lesson (use the tracked filename's exact case) applies to any Go reimplementation of this probe.
- C19 The board justfile.template and the agent-instructions template are each ENTIRELY one managed section, so section-file == template-file on both paths and a fresh install's justfile is byte-identical to the template (asserted at install-suite-behavior.sh:238-240 and 256-258). Any Go path that synthesises rather than copies those bytes breaks the cmp.
- C20 CONCURRENT WRITER — not mine. `git status --short` was clean when this exploration began and now shows ` M do-work/runs/work-2026-08-29-213539/REQ-390-handback.md`, mtime 2026-08-30 07:43:26Z, written ~50 seconds before I observed it. do-work/working/REQ-407-*.md and do-work/CHECKPOINT.md were also touched in the last 30 minutes. Another session in this orchestration is writing to the tree. I read files, ran three read-only test suites (which use mktemp workdirs) and did all experiments under the scratchpad; I did not modify the repository and did not revert that file, since reverting would destroy another agent's work.

*Generated by Explore agent*

## Scope

**Files I will touch:**

- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go` (new)
- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions.go` (new)
- `skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/settingshooks/settings_hooks.go` (new)
- `skills/do-work/tools/do-work-cli/internal/settingshooks/settings_hooks_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go` (new)
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings_test.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify)
- `tools/install-do-work-suite.sh` (modify)
- `skills/do-work/tools/install-do-work-suite.sh` (modify)
- `tools/replace-text-section.sh` (modify)
- `skills/do-work/tools/replace-text-section.sh` (modify)
- `tools/validate-suite-manifest.sh` (modify)
- `skills/do-work/tools/validate-suite-manifest.sh` (modify)
- `tools/fetch-upstream-archive.sh` (modify)
- `skills/do-work/tools/fetch-upstream-archive.sh` (modify)
- `skills/do-work/tools/do-work-update.sh` (modify)
- `_dev/tests/install-suite-behavior.sh` (modify)
- `_dev/tests/update-script-behavior.sh` (modify)
- `_dev/tests/contract-regressions.sh` (modify)
- `README.md` (modify)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify)
- `skills/do-work/tools/prime-do-work-update.md` (modify)

**Files I will NOT touch:**

- VERSION, skills/do-work/VERSION, skills/do-work/actions/version.md, CHANGELOG.md, skills/do-work/CHANGELOG.md — serial-only files owned by the integrator (CLAUDE.md, Before Every Commit). The plan's own coverage gap 1 correctly routes the Go prerequisite line in the version action to the integrator.
- Anything under do-work/ — queue state belongs to the orchestrator alone.
- The Python and jq branches outside the installation and update path (preflight baseline writer, knowledge memory hooks, install-memory-hooks) — REQ-414, REQ-415 and REQ-417 territory, per coverage gap 5.

**Acceptance criteria (restated from REQ):**

- [ ] Bootstrap, install, update, byte-safe managed-section replacement, settings reconciliation, suite validation and archive fetching all run through Go.
- [ ] Fresh-project and existing-project behaviour, exact rollback, custom hooks, reserved recipe collision handling and cancellation are all preserved.
- [ ] CRLF, BOM, NUL, symlinks, file modes, malformed markers and malformed JSON are handled as they are today.
- [ ] No do-work Python or jq branch remains in the installation or update path, and the Go 1.26.1+ prerequisite is documented for installation, update and runtime.
- [ ] Python checks remain only where a Python target capability is being probed.
- [ ] The public scripts-and-Just installation shape still works, through compatibility launchers.
- [ ] Installer and update tests pass with python3 and jq absent from PATH.

## Implementation Summary

**Files changed (39), taken from the merge range rather than from prose:**

- `README.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch.go` (new)
- `skills/do-work/tools/do-work-cli/internal/archivefetch/archive_fetch_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/gittransaction/transaction_findings_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions.go` (new)
- `skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section.go` (new)
- `skills/do-work/tools/do-work-cli/internal/managedsection/managed_section_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/settingshooks/settings_hooks.go` (new)
- `skills/do-work/tools/do-work-cli/internal/settingshooks/settings_hooks_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest.go` (new)
- `skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest_test.go` (new)
- `skills/do-work/tools/do-work-update.sh` (modified)
- `skills/do-work/tools/fetch-upstream-archive.sh` (modified)
- `skills/do-work/tools/install-do-work-suite.sh` (modified)
- `skills/do-work/tools/prime-do-work-update.md` (modified)
- `skills/do-work/tools/replace-text-section.sh` (modified)
- `skills/do-work/tools/validate-suite-manifest.sh` (modified)
- `tools/fetch-upstream-archive.sh` (modified)
- `tools/install-do-work-suite.sh` (modified)
- `tools/replace-text-section.sh` (modified)
- `tools/validate-suite-manifest.sh` (modified)

**What was done:** Moved bootstrap, install, update, byte-safe managed-section replacement, settings reconciliation, suite validation and archive fetching into the `do-work-cli` Go module, and deleted the Python and jq branches they used to run through. Five new internal packages (`managedsection`, `settingshooks`, `suitemanifest`, `archivefetch`, `suiteinstall`) back five registered commands, and the five public shell entry points shrink to compatibility launchers — the installer from 621 lines to 96, the section replacer from 355 to 40, the manifest validator from 151 to 30. The external tools the scripts already invoked (`cp -Rp`, `tar`, `diff`, `git`, `just`, `atomic-download.sh`) stay as subprocesses, which preserves byte, mode and symlink semantics rather than reimplementing them. The `settings_tool` three-way branch and its manual-instruction path are gone: with Go always able to reconcile, "no JSON tool available" ceases to be a state. Also closes the five gaps folded in from REQ-406 — no finding emits a `<command>` placeholder any more, the transaction success path consults `state.existed`, and the `--commit`, `commit_failed`, empty-rollback-status and text-rendering cases gained tests.

**Builder branch:** `worktree-agent-REQ-407-migrate-install-update-bootstrap-to-go` — twelve commits, `1611116` → `bda2f2b`.
**Merge range:** `0bc1480..f45cdca` (cumulative — includes the BOM remediation).
**Hand-back:** `do-work/runs/work-2026-08-29-213539/REQ-407-handback.md` — carries the full per-file manifest, the 23-site Python and jq elimination walk, `## Decisions`, `## Discovered Tasks` and `## Integration Seams`.

## Testing

**Tests run (merged tree, merge range `0bc1480..acf6b73`):**
- `cd skills/do-work/tools/do-work-cli && go test -count=1 ./...` — ✓ all nine packages
- `bash _dev/tests/install-suite-behavior.sh`, `update-script-behavior.sh`, `contract-regressions.sh` — ✓ exit 0 each; the installer suite run three consecutive times clean
- `bash _dev/tests/maintainer-verify.sh` (canonical repository gate, browser lane in its default skipped state) — ✓ **exit 0** on the merged tree, 12 stages. Judged by direct exit status, never piped.

**Red-green validation:** traced to the REQ's `## Red-Green Proof`, whose RED is the installer/update fixtures with `python3` and `jq` removed from PATH, including an existing managed Justfile and custom settings hooks.

- Restricted-PATH installer, existing Justfile: ✗ before — `python3 is required to reconcile an existing Justfile safely` → ✓ after
- Restricted-PATH installer, fresh project: ✗ before — falls through to `settings reconciler: manual` and never writes `settings.json` → ✓ after, with the core hook composed and the retired guard removed
- GREEN was re-verified after the session resumed, with **`go` also absent from PATH**, proving the pre-built binary path rather than an implicit toolchain dependency.

**Byte-preservation proofs (re-verified after resuming):**
- A 28-case characterization of the old embedded Python replacer produces **byte-identical and mode-identical** targets through the Go port — including a lone-CR case no existing fixture covered.
- The Go settings composer is **byte-identical** to the incumbent jq program on the reinstall fixture, key order included.

**A real race, found by running rather than reviewing.** The install signal handler recovered from its own goroutine while the main goroutine was still copying, so an interrupted install intermittently reported its recovery incomplete. It passed standalone and failed under the aggregate — the shape of a race, not a flake. Fixed in `7685de4`: the handler now cancels a work context and waits for the single main-goroutine recovery.

**New tests added:** 40 managed-section and Just-scanner cases (lone-CR splitting, CRLF, NUL, BOM asymmetry, malformed and reversed markers, symlink and dangling-symlink refusal, mode preservation, byte-idempotence, multiline-literal boundaries); settings-hook key-order, append-unique, guard-removal, custom-hook preservation and five malformed-input refusals; 16 manifest rejection cases; archive-fetch branch derivation across five URL shapes plus route fallthrough and preserved-target-on-total-failure; and the full install/update transaction matrix including cancellation, post-write-failure recovery, collision refusal before confirmation, and the stdout-result / stderr-narration split.

**Existing tests updated (cross-REQ impact):** the three behavioural shell suites were retargeted at the Go implementation. `install-suite-behavior.sh` lost the Python-fallback and no-JSON-tool lanes (99 lines) and the flaky-jq post-write lane (37 lines) because the states they asserted no longer exist; `contract-regressions.sh` moved eleven grep assertions onto the Go sources that now own each behaviour. The hand-back records each with its reason.

*Verified by work action*

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

## Remediation — BOM-prefixed settings.json

**A regression this REQ introduced, found in review and fixed inside it.** A `.claude/settings.json`
carrying a UTF-8 BOM installed fine before REQ-407 and afterwards hard-failed the whole install with
`settings are not valid JSON: invalid character 'ï' looking for beginning of value` — the first BOM byte
read as Latin-1, exit 2.

The regression is against real incumbent behaviour, measured both ways:

| Incumbent | BOM-prefixed input |
|---|---|
| `jq -e .` | **accepts and strips it**, exit 0 |
| `python3 -m json.tool` | rejects it |

The old installer's three-way branch preferred `jq` whenever it was present, so the common case worked
and stopped working. Any project whose settings file came from a Windows editor or a PowerShell redirect
was hard-blocked from installing.

**Fixed in `8b4e1b1`:** `decodeOrderedJSON` strips one leading BOM before decoding, reproducing jq's
behaviour including that the mark never re-enters the encoded output. Deliberately narrow — a doubled
mark and a mark at a non-zero offset both stay refused, because those are malformed input rather than a
byte-order mark, and all five pre-existing malformed-input refusals keep their original messages.

**Evidence.** RED committed before the fix: `TestLeadingByteOrderMarkIsStrippedLikeJq` failing with the
exact `invalid character 'ï'` text. Independently re-verified by the orchestrator on the merged tree
against a throwaway fixture repo carrying a BOM-prefixed settings file:

```
install-suite: success
change .claude/settings.json [modified]: core hooks composed into existing settings
rollback: not_needed
EXIT=0
written file first bytes: 7b 0a 20 20   (no ef bb bf)
```

*(One earlier orchestrator run appeared to contradict this. It was an instrumentation error — a
`${PIPESTATUS[0]}` reading past an intervening `&&`, against a fixture the installer never actually ran
on. The re-run above separates the streams and chains nothing.)*

## Review

**Acceptance: Partial.** The REQ's GREEN condition holds: installer and update fixtures succeed with `python3` and `jq` absent from PATH, existing managed Justfiles and custom settings hooks are preserved, and the canonical gate exits 0 on the merged tree. One regression this REQ introduced was found and fixed inside it; one narrow regression remains and is carried by a follow-up.

**Method.** Five reviewers, one per dimension: byte-and-filesystem fidelity (asked to recover the pre-REQ Python and jq from git history and re-derive the byte-identical claim against adversarial inputs of its own devising, not to trust the builder's 28 cases), install safety, the shell launchers and the public shape, elimination-and-scope, and consumer impact. Every finding then went to three skeptics prompted to refute it, each asked to state the severity its evidence actually supports.

**The first verification pass was cut short by a session usage limit that killed 89 of its 98 agents.** The eleven Important findings were re-adjudicated to completion afterwards, and two were settled directly by the orchestrator. Minor and Nit findings go in this report only.

**Dimension verdicts:** all five partial.

**Findings that stand.**

- **Important, and fixed inside this REQ — a BOM-prefixed `.claude/settings.json` no longer installed.** Full account in `## Remediation — BOM` above. Fixed in `8b4e1b1`.
- **Minor — the three special mode bits are stripped from managed files.** `permissionsOf` and the settings-mode read use `info.Mode().Perm()` (mask `0o777`) where the Python they replaced used `stat.S_IMODE` (mask `0o7777`), so setuid, setgid and sticky are dropped from `Justfile`, `CLAUDE.md` and `.claude/settings.json` on every install. Reproduced three ways: unit level (`7644` → `644`, `6755` → `755` — execute bits survive, only the special bits are lost), A/B against the pre-REQ Python replacer, and end-to-end through both installers. Adjudicated 2–1; **both surviving verifiers judged it Minor rather than Important**, because git records only `100644`/`100755`, `umask` cannot set special bits, and a setgid *directory* propagates the group rather than the bit — so reaching it requires a consumer to have hand-run `chmod` on one of those three files. Carried by [REQ-426](../queue/REQ-426-preserve-special-mode-bits.md).

**Findings refuted, nine of ten Important ones knocked down by measurement rather than argument.** The verifiers built real install fixtures, drove signals at them on a real pty, and compared against the pre-REQ implementations recovered from git history:

- *HUP/INT/TERM at the confirmation prompt deadlocks the installer* — refuted 3–0. The mechanism is real; the deadlock is not, measured on a real pty.
- *The launchers dropped their signal traps, so a signal returns while the Go install keeps writing* — refuted 3–0. The claimed consequence needs a signal reaching the launcher but not the CLI; bash runs the CLI in the same process group, measured.
- *The fetch error is swallowed, losing the route report and the `DO_WORK_UPSTREAM_URL` escape hatch* — refuted 2–1. The one-line output difference is real; every consequence claimed from it is not.
- *`--upstream-url` is dead surface* — refuted 3–0, and the proposed fix does not compile.
- *The missing-Go message names no download source and drops its no-files-changed assurance* — refuted 3–0; four of its five factual claims are false, and the measured consequence is one stderr line on a path that writes nothing.
- *Shipped board docs still say Go is needed by exactly one action* — refuted 3–0 on its factual premise.
- *An installed suite carries no record of the Go prerequisite* — refuted 3–0.
- *Every finding's next/verify argv names a binary not on a consumer's PATH* — refuted 3–0.
- *The Go 1.26.1 floor is higher than the module needs* — refuted 3–0 **as a defect in this REQ**, and correctly so: `REQ-406` and `REQ-407` both take that floor verbatim from UR-081, so implementing it is not a bug. The measurement behind the finding is nonetheless real — the module builds and passes all six packages' tests at `go 1.23.0` — and because the floor now decides who can install at all, the question is put to the maintainer as [REQ-427](../queue/REQ-427-confirm-go-version-floor.md) (`pending-answers`) rather than dropped.

**Scope drift:** none. 39 files declared, 39 touched. The first Implementation Summary listed only 32 because it was assembled by parsing the hand-back's prose; it was rebuilt from the merge range itself, which is the right source.

*Reviewed by work action (five dimensions, three refutation lenses per finding, plus direct orchestrator measurement)*

## Lessons Learned

**What worked:** Keeping the external tools the scripts already invoked — `cp -Rp`, `tar`, `diff`, `git`, `just`, `atomic-download.sh` — as subprocesses instead of reimplementing them in Go. That is what made a 6600-line port preserve byte, mode and symlink semantics for free, and it kept the existing PATH-stub fixtures working verbatim. Deleting the `settings_tool` three-way branch outright was the right shape: with Go always able to reconcile, "no JSON tool available" stopped being a state rather than becoming a better-handled one. Committing each package as its tests passed saved the work when a usage limit killed the builder mid-run — the earlier attempt had eight modified files and five new packages uncommitted at the moment it died.

**What didn't:** Two claims of parity that were narrower than they sounded.

"Byte-identical" was asserted from 28 characterization cases whose fixtures used modes `750`, `640` and `600` — none with a special bit — so the mode regression sat inside a claim that specifically said "mode-identical". **A characterization suite proves parity over the inputs it contains, and its fixtures are where you look for what it silently excludes.** The same applies to the settings composer, which is byte-identical to jq except for `U+2028`, `U+2029` and `U+007F`.

The BOM case is the sharper lesson: the port replaced *two* incumbents with different behaviour — jq accepted and stripped a BOM, `python3 -m json.tool` rejected it — and matching one of them silently regressed the path that preferred the other. **When a branch is being deleted, parity is against the branch that actually ran, not the one that is easiest to compare with.**

**Worth knowing:** `permissionsOf` in `managed_section.go` and the settings-mode read in `install_transaction.go` are the two places file modes enter the write path; both mask to `0o777` today. `decodeOrderedJSON` is the single door every settings input goes through, which is why one `TrimPrefix` there covers both `ComposeSettings` inputs. The launcher only builds when the binary is missing or older than its sources, so a shipped prebuilt binary makes the Go prerequisite moot at runtime — that is why the GREEN evidence was re-run with `go` itself absent from PATH.

## Orientation

Installing and updating do-work now runs on one Go command instead of shell plus whichever of `jq` or `python3` happened to be present. `do-work-cli` gained five commands — install, update, managed-section replacement, manifest validation and archive fetching — and the five public shell entry points became thin launchers over them, the installer dropping from 621 lines to 96. What a consumer notices is that settings are always reconciled and their `settings.json` keeps its exact key order, where before the outcome depended on their toolchain and a "no JSON tool" path printed manual instructions instead of installing.

`[MAP CHANGED]` — the installation path changed language and gained a prerequisite. Go 1.26.1+ is now required to install or update (see REQ-427, which asks whether that floor is the one you want). The `settings_tool` three-way branch and its manual-instruction path are gone from the map entirely. `skills/do-work/tools/prime-do-work-update.md` was updated by this REQ: its Read-first list now points at the Go packages and it gained a Traps section for the three hazards the migration introduced.

This REQ also registered the first real `do-work-cli` commands, so the pattern the remaining batch copies is now fixed: stdout carries only the rendered `CommandResult`, narration goes to stderr, and exit codes come solely from `resultmodel.ExitCode`.

