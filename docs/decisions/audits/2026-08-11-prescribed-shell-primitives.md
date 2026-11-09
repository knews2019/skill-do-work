# Prescribed Shell Primitive Inventory

**Ran:** 2026-08-11 · **Scope:** shipped `skills/` tree · **Method:** grep-based inventory followed by reading every match; paths, not line numbers, are the durable coordinates.

## REQ-171 inventory gate

The operative starting census is **59 shell fences, 40 multi-line fences, 515 physical multi-line body lines (470 nonblank)**. The clean builder worktree reproduced 59 / 40 / 512 / 467; the required read-only main-tree screenshot diff adds exactly three nonblank body lines to `capture.md`, reconciling the independently verified 515 / 470 operative baseline before promotion. Dispositions below are exhaustive: **17 promote, 21 inline residue, 2 Go-owned**.

| Durable coordinate (file → heading → content label) | Disposition | Executable home | Behavior case / reason |
|---|---|---|---|
| `skills/do-work/docs/prescribed-shell-primitives.md` → Local Git ignore → resolve Git exclude | PROMOTE | `skills/do-work/scripts/add-local-git-exclude.sh` | subdirectory/worktree-safe exclude |
| same guide → Atomic download publication → curl temporary publication | PROMOTE | `skills/do-work/scripts/atomic-download.sh` | failed partial never becomes final |
| `skills/do-work/actions/work.md` → Step 1: Find Next Request → blocked-check timeout | PROMOTE | `skills/do-work/scripts/run-blocked-check.sh` | GNU-less timeout returns 124 |
| `skills/do-work/actions/capture.md` → File Naming → queue-kanban next-req accelerator | GO OWNER | `skills/do-work-board/tools/queue-kanban` | atomic reservation stays Go-only |
| `skills/do-work/actions/capture.md` → Step 4: Handle Screenshots → staged screenshot lifecycle | PROMOTE | `skills/do-work/scripts/capture-screenshot.sh` | coordinated two-writer no-clobber race |
| `skills/do-work/actions/capture.md` → Step 7: Commit → UR input/assets staging | INLINE RESIDUE | caller | staging/heredoc workflow example |
| `skills/do-work/actions/work-reference.md` → Commit & Metadata-Commit Procedure → implementation + archive staging | INLINE RESIDUE | caller | owner policy and commit heredoc |
| same file/heading → serial record-hash sequence | INLINE RESIDUE | `tools/checks/record-commit-hash.sh` invocation remains inline | existing cross-step orchestration |
| `skills/do-work/actions/commit.md` → Step 1: Preflight → protected inventory start | PROMOTE | `skills/do-work/scripts/protected-inventory.sh start` | once-X quarantine |
| same file → Step 3: Associate with REQs → protected association | PROMOTE | `skills/do-work/scripts/protected-inventory.sh associate` | safe owner, quarantined secret |
| same file → Step 5: Commit → exact XD deletion | PROMOTE | `skills/do-work/scripts/stage-exact-deletion.sh` | pathological exact path only |
| same file/heading → REQ-associated commit example | INLINE RESIDUE | caller | independent staging/commit recipe |
| same file/heading → unassociated commit example | INLINE RESIDUE | caller | independent staging/commit recipe |
| `skills/do-work/actions/review-work.md` → Commit (Standalone mode) → review metadata commit | INLINE RESIDUE | caller | independent staging/record-hash recipe |
| `skills/do-work/actions/cleanup.md` → Commit (Git repos only) → cleanup commit | INLINE RESIDUE | caller | action-specific staging/heredoc |
| `skills/do-work/actions/forensics.md` → Release and Queue Invariants → queue-kanban verify | GO OWNER | `skills/do-work-board/tools/queue-kanban` | verification remains Go-owned |
| `skills/do-work-toolbox/actions/ui-review.md` → Step 8.5 → viewport screenshot recipe | INLINE RESIDUE | caller | UI command recipe |
| same file/heading → accessibility snapshot recipe | INLINE RESIDUE | caller | UI command recipe |
| `skills/do-work-toolbox/actions/inspect.md` → Step 1: Preflight → protected inventory start | PROMOTE | `skills/do-work/scripts/protected-inventory.sh start` | inspect namespace quarantine |
| same file → Unscoped mode → protected association | PROMOTE | `skills/do-work/scripts/protected-inventory.sh associate` | safe association only |
| `skills/do-work-toolbox/actions/deep-explore.md` → Create the directory → session assignment | INLINE RESIDUE | caller | assignment-only invocation |
| `skills/do-work-toolbox/actions/install.md` → Bowser Phase 1 → detect components | INLINE RESIDUE | caller | independent detection commands |
| same file → Bowser Phase 4 → atomic skill download | PROMOTE | `skills/do-work/scripts/atomic-download.sh` | failed download preserves target |
| same file → Bowser Phase 5 → verify components | INLINE RESIDUE | caller | independent verification commands |
| same file → last30days Phase 2 → clone/copy transaction | PROMOTE | `skills/do-work-toolbox/scripts/install-last30days.sh` | fixture source installation |
| same file/phase → last30days local exclude | PROMOTE | same script via core `add-local-git-exclude.sh` | sibling helper resolves |
| same file → last30days Phase 3 → full guarantee verification | PROMOTE | same script, `check` mode | skill/ignore/Python gate |
| `skills/do-work-toolbox/actions/ai-report-reference.md` → Image Generation Backend → STYLE assignment | INLINE RESIDUE | caller | assignment-only style data |
| same file/heading → generation helper | PROMOTE | `skills/do-work-toolbox/scripts/generate-report-image.sh` | direct backend exact output/inert prompt |
| same file/heading → parallel image orchestration | INLINE RESIDUE | caller | wrapper owns mechanics, caller owns orchestration |
| `skills/do-work-toolbox/actions/ai-report.md` → Step 7: Render and Judge → local server assignment | INLINE RESIDUE | caller | independent UI recipe |
| `skills/do-work-knowledge/docs/interview-guide.md` → Integration with bkb → interview ingest example | INLINE RESIDUE | caller | independent example commands |
| same file/heading → bkb query example | INLINE RESIDUE | caller | independent example commands |
| `skills/do-work-knowledge/actions/setup-memory.md` → Preconditions → root assignments | INLINE RESIDUE | caller | assignment-only block |
| `skills/do-work-knowledge/actions/memory.md` → Input → root assignments | INLINE RESIDUE | caller | assignment-only block |
| `skills/do-work-knowledge/actions/memory-reference.md` → Lexical Recall → grep/scoring recall | PROMOTE | `skills/do-work-knowledge/scripts/lexical-memory-recall.sh` | raw apostrophe/query stays data |
| same file → Semantic Recall → backend probes | INLINE RESIDUE | caller | independent detect commands |
| same file → Usage-Ledger Contract → best-effort append | INLINE RESIDUE | caller | usage-ledger policy append |
| same file → Hook Install Internals → JSON hook merge | PROMOTE | `skills/do-work-knowledge/scripts/install-memory-hooks.sh` | partial prior hook merge |
| `skills/do-work-board/actions/board.md` → Step 5: Run selected mode → static artifact exclude | PROMOTE | core `add-local-git-exclude.sh` | sibling path resolution |

After promotion the shipped census is **59 shell fences, 23 multi-line fences, 200 multi-line body lines (167 nonblank)**. Against the operative 515 / 470 baseline, shipped multi-line Markdown loses **315 physical lines / 303 nonblank lines** while all 17 promoted callers remain visible as one-line invocations. The executable replacement is 563 script lines plus 210 lines of sourceable fixture helpers and behavior probes.

## Disposition

| Primitive | Canonical shipped home | Former rationale sites | Divergence / execution-only disposition |
|---|---|---|---|
| Per-file untracked inventory | `skills/do-work/docs/prescribed-shell-primitives.md` → Per-file untracked inventory; executable source remains `tools/checks/uncommitted-inventory.sh` | `skills/do-work/actions/commit.md`, `skills/do-work-toolbox/actions/inspect.md`, `skills/do-work-toolbox/actions/stray-check.md` | REQ-121's NUL/rename/secret-aware variant wins. `work.md`, `tidy-repo.md`, preflight, and queue-kanban uses are narrower executions, not fallback copies. |
| Merge-aware commit diff | core guide → Merge-aware commit diff | `skills/do-work/actions/review-work.md`, `skills/do-work-toolbox/actions/ai-report.md`, three template/instruction sites in `skills/do-work-toolbox/actions/present-work.md` | The fixed form is quoted `'<commit>^2'` plus `git show --first-parent -m`; worktree orchestration keeps its explicit `<pre>..<merge_hash>` range. |
| Commit file listing | core guide → Commit file listing | explanatory copy in `skills/do-work-toolbox/actions/ai-report.md` | `skills/do-work/tools/checks/blanked-req-scan.sh` is executable use/comment, not a second policy home. `git diff-tree ... -r -m` wins over message-bearing `git show --name-only`. |
| Local Git ignore | core guide → Local Git ignore | long copies in the core, knowledge, and toolbox `crew-members/background-agents.md`; rationale in board/install actions | Commands in board/Just/install/setup-memory remain at their execution sites. The fixed variant uses `git rev-parse --git-path`, `**/`, and a separate tracked-file check when “never committed” is the requirement. |
| Atomic download publication | core guide → Atomic download publication | explanatory and Red-Flag copies in `skills/do-work-toolbox/actions/install.md` | Download commands in install tables and shipped updater/installer scripts are required uses. The fixed temp-download, rename-on-success, cleanup-plus-`false` shape wins. |
| Raw text before shell quoting | core guide → Raw text before shell quoting | `skills/do-work-knowledge/actions/memory-reference.md` | `memory.md` already points at its reference file; query-specific tokenization remains there, while the injection rationale moves to the guide. |
| Diff output filtering | core guide → Diff output filtering | no duplicate shipped action statement found | The maintainer trap is now shipped once. Existing diffs keep caller-specific filters; no command required rewriting in this pass. |
| State across command blocks | core guide → State across command blocks | general rationale in `work.md`, `work-reference.md`, `commit.md`, and toolbox `inspect.md` | Deterministic re-derivation and literal merge endpoints remain at callers; only the repeated general explanation moves. |

## Verification searches

The regression probe `_dev/tests/prescribed-shell-canonicalization.sh` owns the executable search. It requires the guide headings and former-site pointers, then rejects the known high-risk rationale phrases outside the guide. `_dev/tests/action-shell-blocks.sh` separately parses every surviving fenced shell block.

Execution count is not duplication count. A command that must run in three callers remains three commands; this audit consolidates the shared failure-mode statement so a future semantic correction has one source.
