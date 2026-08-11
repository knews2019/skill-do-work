# REQ-171 Implementation Plan

## Discovery Boundary

`prime_files` is empty, so this plan does not load or infer a global architecture. Discovery is limited to the REQ and REQ-165 record, the prescribed-shell guide and its audit/ratchet, the shipped Markdown shell fences, the cited REQ-166 fixture test, the aggregate test seam, and tests that currently pin a block being promoted.

## Promotion Inventory Approach

The first implementation artifact is an update to `decisions/audits/2026-08-11-prescribed-shell-primitives.md`. Extend its current primitive table with a path-and-heading inventory of every multi-line `bash`/`sh` fence in shipped `skills/` Markdown. Use the same Markdown fence boundary as `_dev/tests/action-shell-blocks.sh` (zero to three leading spaces), and use durable coordinates (file plus enclosing heading, not line numbers). Give every block exactly one disposition:

1. **Promote** — the block contains control flow, failure recovery, temporary-state management, parsing, concurrency, or an ordered filesystem mutation whose semantics can regress while syntax still passes.
2. **Already executable** — the normative behavior already lives in a shipped, behavior-tested script; retain only its invocation and caller-specific policy.
3. **Inline residue** — the block is one command, a sequence of independent illustrative commands, a template/transcript, or pseudocode with placeholders rather than a reusable executable contract.
4. **Non-shell owner** — the capability belongs to the shipped Go board tool (notably atomic REQ reservation), so REQ-171 must not create a shell twin.

The inventory is a gate, not a report written after coding: no script implementation starts until every current multi-line fence has a disposition and every `Promote` row names its target script and behavior probe. Re-run the same census after migration; any surviving promoted block is a failure. Record the before/after count of executable shell lines inside shipped Markdown fences and the added script/test lines in the REQ's final implementation report. The Markdown shell-line count must decrease.

### Expected promotion set from the bounded census

The builder must verify this table against the inventory before editing; changes require a recorded disposition, not silent omission.

| Existing surface | Planned executable home | Call-site result |
|---|---|---|
| Merge-aware ordinary/merge commit display in the canonical guide and review/report/presentation actions | `skills/do-work/scripts/show-commit-diff.sh` | Pass the literal commit as one argument; script detects parent 2 and chooses ordinary or first-parent merge output. |
| Local Git exclude blocks in the guide, board static mode, memory setup/background guidance, and toolbox install | `skills/do-work/scripts/add-local-git-exclude.sh` | Pass repository root, probe path, and root-aligned exclude pattern as separate arguments; callers retain whether an already-tracked path is forbidden. |
| Atomic curl publication in the guide and toolbox skill downloads | `skills/do-work/scripts/atomic-download.sh` | Pass URL and final path as separate arguments; no caller writes the final file incrementally. The self-bootstrapping suite installer remains a documented inline exception because the helper does not exist until the archive is fetched. |
| Capture Step 4 screenshot copy/verify/no-clobber/source-cleanup block | `skills/do-work/scripts/capture-screenshot.sh` | Pass the dispatcher-supplied staged source and permanent destination; output/exit status tells the action whether it may link the asset. |
| Work's bounded `blocked_check` runner | `skills/do-work/scripts/run-blocked-check.sh` | Pass the deterministic probe file; script owns `timeout`/`gtimeout`/portable polling and returns the probe status (124 on timeout). |
| Commit/inspect inventory setup and later quarantine overlay/association blocks | `skills/do-work/scripts/protected-inventory.sh` | Use explicit `start` and `associate` modes around the existing `tools/checks/uncommitted-inventory.sh` and `associate-files.sh`; both core commit and toolbox inspect call the core script. |
| Commit's exact-deletion staging guard | `skills/do-work/scripts/stage-exact-deletion.sh` | Pass one literal path; script verifies the cached status before and after `git add -u` and refuses extra records. |
| Knowledge lexical recall pseudocode/scoring block | `skills/do-work-knowledge/scripts/lexical-memory-recall.sh` | Supply the raw query as data (stdin), never interpolated shell; script emits the bounded, attributed ranked results. |
| Knowledge settings hook merge block | `skills/do-work-knowledge/scripts/install-memory-hooks.sh` | Pass project root and knowledge package root; script performs independent hook gating, validated temporary publication, verification, and rollback. |
| Toolbox report-image generator and its output verification | `skills/do-work-toolbox/scripts/generate-report-image.sh` | Pass absolute output and sanitized description separately; script owns backend selection, isolated temporary directory, and usable-file exit status. Dynamic multi-image orchestration may remain a short inline invocation sequence. |
| Toolbox `last30days` clone/copy/ignore/verification workflow | `skills/do-work-toolbox/scripts/install-last30days.sh` | One script owns detect/repair/install/verify and calls core local-ignore; the action retains consent, trust, and reporting policy. |

Expected inline/already-owned dispositions include: `git diff-tree --no-commit-id --name-only -r -m` (one-line path listing), simple root-variable assignments, independent detect/verify commands, UI-review command recipes, commit-message heredoc examples, report style constants, Go build/reservation invocations, and the existing secret-aware inventory implementation in `skills/do-work/tools/checks/`. If the completed inventory proves that one of these has reusable branching or recovery logic, move it to `Promote` before implementation and add its fixture case.

## Files to Modify

### Inventory and canonical contract

- `decisions/audits/2026-08-11-prescribed-shell-primitives.md` — promotion inventory, script homes, inline/non-shell rationales, and self-bootstrap exception.
- `skills/do-work/docs/prescribed-shell-primitives.md` — make scripts the normative implementation; keep concise intent, failure-mode rationale, arguments/exit meanings, and explicitly marked inline-only primitives.
- `_dev/tests/prescribed-shell-canonicalization.sh` — require every advertised script to exist and be executable, reject old promoted implementations/restatements at callers, preserve the single-home ratchet, and avoid depending on volatile line numbers.
- `decisions/audits/2026-08-11-defensive-surface.md` — add executable-layer dispositions for every new shipped `.sh`, required by the existing defensive-surface ratchet.

### New shipped scripts

- Core: `skills/do-work/scripts/show-commit-diff.sh`, `add-local-git-exclude.sh`, `atomic-download.sh`, `capture-screenshot.sh`, `run-blocked-check.sh`, `protected-inventory.sh`, and `stage-exact-deletion.sh`.
- Knowledge: `skills/do-work-knowledge/scripts/lexical-memory-recall.sh` and `install-memory-hooks.sh`.
- Toolbox: `skills/do-work-toolbox/scripts/generate-report-image.sh` and `install-last30days.sh`.

The inventory can narrow this list only by recording that the current block is inline residue or already executable and why; it can extend the list only together with a fixture case and call-site rewrite.

### Shipped call sites

- Core actions: `skills/do-work/actions/capture.md`, `commit.md`, `review-work.md`, and `work.md`.
- Board: `skills/do-work-board/actions/board.md`.
- Knowledge: `skills/do-work-knowledge/actions/memory-reference.md` and `setup-memory.md`.
- Toolbox: `skills/do-work-toolbox/actions/ai-report.md`, `ai-report-reference.md`, `inspect.md`, `install.md`, and `present-work.md`.
- Cross-package background guidance: `skills/do-work/crew-members/background-agents.md`, `skills/do-work-knowledge/crew-members/background-agents.md`, and `skills/do-work-toolbox/crew-members/background-agents.md` where the local-ignore primitive is invoked.
- Existing executable caller: `skills/do-work/tools/do-work-update.sh` should call the core atomic-download helper if its updater sequencing remains valid. `skills/do-work/tools/install-do-work-suite.sh` remains self-contained for bootstrap and is explicitly exempted in the inventory.

Each rewritten action keeps one sentence saying what the step intends, followed by one invocation. Core callers use `<skill-root>/scripts/...`; sibling packages use explicit `<skill-root>/../do-work/scripts/...` references for core primitives. Caller-specific gates, consent, and policy stay in the action rather than moving into generic scripts.

### Test surface

- Add `_dev/tests/fixture-repo.sh` as a sourceable helper for isolated repositories, Git identity/initial commit, seeded `do-work/queue`, version fixtures, nested untracked files, commits, branches, and merges. It owns cleanup registration but no assertions.
- Add `_dev/tests/prescribed-shell-scripts-behavior.sh` as the focused execution probe for all promoted scripts.
- Update `_dev/tests/contract-regressions.sh` to guard and invoke the new probe explicitly; `_dev/tests/*.sh` are not auto-discovered.
- Update `_dev/tests/staged-skills-contract.sh` so the screenshot contract pins the shipped script invocation and behavior test rather than the removed inline implementation.
- Optionally migrate `_dev/tests/session-start-hook-behavior.sh` to the shared fixture helper only where it removes duplicate setup without changing REQ-166 cases or assertions.
- Keep `_dev/tests/action-shell-blocks.sh`; its Markdown portion continues to lint the inline residue and its direct `.sh` scan syntax/lints every new shipped script. Only its wording needs adjustment if it still claims all prescribed behavior remains inline.

### Integrator-only release files

Builders do not touch serial release files. After reconciliation, the integrator updates `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, root `CHANGELOG.md`, and the byte-identical `skills/do-work/CHANGELOG.md` mirror. Choose a version strictly above the then-current top changelog entry; do not preselect a number while parallel work is active.

## Ordered Changes

1. **Create and ratchet the inventory.** Extend the existing prescribed-shell audit with every current multi-line fence and its disposition. Add canonicalization assertions for the intended script homes so the first targeted run is RED because the scripts/caller forms do not yet exist.
2. **Build the shared fixture scaffold.** Add the sourceable fixture helper with deterministic setup functions and cleanup. Reuse it in the new behavior probe; migrate REQ-166's test only if the resulting diff is smaller and its four behavior cases remain unchanged.
3. **Promote the three shared canonical primitives first.** Implement merge-aware diff, local Git exclude, and atomic download scripts with fixture cases, then replace guide blocks and all discovered action callers. Run the focused probe and canonicalization after each script so each promotion is independently green and abortable.
4. **Promote core action-specific logic.** Move screenshot capture, blocked-check timeout, protected inventory orchestration, and exact-deletion staging into core scripts. Rewrite capture/work/commit/inspect call sites and update the staged screenshot contract. Preserve existing exit meanings and durable-path behavior.
5. **Promote package-owned logic.** Move lexical recall and hook installation into knowledge scripts; move report-image generation and `last30days` installation into toolbox scripts. Keep package-specific behavior in its package, but call shared local-ignore/atomic helpers from core by explicit sibling path.
6. **Finish canonical guide/caller contraction.** Replace promoted fences with one-line intent plus invocation, point every guide section at its normative implementation or mark it as deliberate inline residue, and ensure no action relies on variables or random paths from another command block.
7. **Update audit and suite seams.** Add all new scripts to the defensive-surface audit, wire the behavioral probe into `contract-regressions.sh`, and update comments/guards without weakening the lint harness.
8. **Re-run the inventory and account for net surface.** Prove every original multi-line fence has a disposition, no promoted implementation remains in prose, and executable shell lines in shipped Markdown decreased. Add the before/after prose count and script/test additions to the REQ report.
9. **Reconcile and release.** Review complete branch diffs, resolve only true overlapping call-site changes, run the full verification set on the integrated tree, then perform the serial version/changelog mirror update.

## Architectural Decisions

- Scripts accept paths, URLs, commits, and modes as separately quoted arguments; raw/imported text is read as data (stdin where appropriate). No script evaluates constructed shell strings.
- Stdout is reserved for result data/actions that callers consume; diagnostics go to stderr. Usage or environmental/setup errors return 2; semantic negative results preserve the existing caller contract; underlying Git/curl/probe failures are not converted to success.
- Generic scripts own mechanics only. Actions continue to own user confirmation, fallback/report wording, tracked-file policy, and workflow decisions.
- Core is the sole home for a primitive used by more than one package. Knowledge/toolbox scripts are only for package-specific capabilities. Cross-package references always name the core sibling explicitly.
- Scripts target the repository's portability floor: Bash 3.2/POSIX utilities, no associative arrays or `mapfile`, no GNU-only dependency without a tested fallback. Git-internal paths come from `git rev-parse --git-path`.
- Temporary publication uses deterministic sibling names where later blocks must re-derive them; single-invocation private scratch uses `mktemp` plus a trap. Cleanup never masks the failure it follows.
- The Go board tool remains the owner of atomic REQ reservation. No shell duplicate is introduced.
- `install-do-work-suite.sh` is a justified self-hosting exception to the shared atomic-download helper because it must fetch the package containing that helper.
- The suite ships whole package directories through `suite/modules.tsv`; adding `scripts/` requires no manifest row, but staged/install-layout tests must prove explicit sibling invocations resolve after installation.

## Requirement Mapping

| REQ-171 requirement | Planned coverage |
|---|---|
| Shared fixture scaffold first | Ordered changes 1–2; `_dev/tests/fixture-repo.sh`. |
| Promote each qualifying primitive | Gated inventory plus ordered changes 3–6 and the expected promotion table. |
| Guide points at script; ratchet survives | Guide rewrite and `prescribed-shell-canonicalization.sh` changes in steps 1 and 6. |
| Fixture test for each semantic trap | One named case per promoted script in `prescribed-shell-scripts-behavior.sh`; no promoted row can land without a case. |
| No shell twin for Go reservation | Non-shell-owner inventory disposition and architectural decision. |
| One-line intent plus invocation at call sites | Call-site contraction in steps 3–6. |
| Cross-package primitives live in core | Core ownership and explicit sibling paths in files/architecture sections. |
| Prose surface shrinks and is reported | Step 8 line accounting and a hard requirement that shipped Markdown shell lines decrease. |
| Existing lint harness remains | `action-shell-blocks.sh` retained for residue and direct script linting. |

## Testing Approach

Follow the captured RED/GREEN proof: add execution cases first and confirm they fail because the promoted script/canonical call form is absent. Then implement one script at a time until its cases pass.

Focused fixture cases:

- **Merge diff:** ordinary commit output contains its patch; a real two-parent merge returns the first-parent patch instead of an empty combined diff; invalid revision fails nonzero.
- **Local ignore:** run from a subdirectory and linked worktree; assert `git rev-parse --git-path` location, root-aligned pattern, idempotent single append, space-safe arguments, and non-Git behavior.
- **Atomic download:** fake curl writes a non-empty partial then fails; assert nonzero, old/final target is not published, and `.download` is removed. Success publishes exact bytes and leaves no temp file.
- **Screenshot capture:** successful bytes match, destination is no-clobber, staged source and now-empty dispatch directory are removed; an existing destination or copy/link failure preserves the source and fails; inline-source mode does not delete its source.
- **Blocked check:** success/failure propagation, GNU timeout path, forced portable fallback path, timeout 124, and deterministic probe cleanup.
- **Protected inventory:** wholly untracked nested secret reaches the inventory; quarantine persists across `start`/`associate`; an earlier `X` cannot reappear as readable A/M; missing deterministic quarantine fails closed.
- **Exact deletion:** an exact cached D is accepted; unstaged deletion is staged and rechecked; rename/multiple-record/unrelated-status shapes fail.
- **Lexical recall:** hostile quote/command text remains data, token filtering is correct, ranking honors distinct-token score and recency weights, output is capped and attributed.
- **Memory hooks:** empty/existing/partial settings, both hooks already present, invalid JSON, missing `jq`, preservation of unrelated hooks, validated temp publication, and rollback on verification failure.
- **Report image:** direct backend success, missing/empty output failure, disabled agentic fallback, enabled isolated fallback with fake backend, and temp cleanup.
- **Last30days:** fresh install, partial repair, clone/copy failure, local-ignore failure, Python-version failure, non-Git mode, idempotent rerun, and no project config/API-key write.

Verification commands on the reconciled tree:

- `bash -n` and ShellCheck at warning severity for every new/modified shell file.
- `bash _dev/tests/prescribed-shell-scripts-behavior.sh`.
- `bash _dev/tests/prescribed-shell-canonicalization.sh`.
- `bash _dev/tests/action-shell-blocks.sh --self-test` and `bash _dev/tests/action-shell-blocks.sh`, including a PATH without ShellCheck to preserve degraded behavior.
- `bash _dev/tests/session-start-hook-behavior.sh` if the fixture helper is reused there.
- `bash _dev/tests/staged-skills-contract.sh` and `bash _dev/tests/shipped-package-reference-contract.sh` to prove source and installed sibling paths.
- `bash _dev/tests/defensive-surface-audit.sh`.
- `bash _dev/tests/contract-regressions.sh`.
- `git diff --check`, a debug-artifact scan, and a final multi-line fence inventory/net-surface calculation.

Acceptance is GREEN only when temporarily regressing either `atomic-download.sh` to direct-final publication or `show-commit-diff.sh` to plain `git show` makes the focused probe fail with the script and case name, while the clean integrated tree passes the complete suite.
