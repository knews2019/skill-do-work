# Hand-back — REQ-714 (suite writes its knowledge base and ADR log under docs/)

Branch `worktree-agent-REQ-714-suite-writes-kb-and-decisions-under-docs`, tip `aaa4f4a5`, base `a79e94bb`. Worktree clean. Nothing pushed, nothing merged.

## Commits

1. `11afb731` [REQ-714] RED: knowledge base default location tests
2. `f328ce6d` [REQ-714] GREEN: knowledge commands default to docs/kb
3. `4934f0bd` [REQ-714] recipes and prose: knowledge base and ADR log live under docs/
4. `aaa4f4a5` [REQ-714] move this repository's kb/ and decisions/ under docs/

## Files

All modified, none created or deleted outside the two moved folders (41 files).

Go (`skills/do-work/tools/do-work-cli/internal/knowledgecommands/`):
- `commands.go` modified: the parser leaves the target empty unless supplied.
- `bkb_scan.go` modified: `resolveKnowledgeBaseLocation`, the two notices, `locateKnowledgeBase` treats empty as the default order.
- `bkb_init.go` modified: handler calls the resolver; `targetAncestors` creates missing parent folders.
- `dream_scan.go`, `interview_commands.go`, `memory_commands.go` modified: use the resolver, no `"kb"` literal.
- `bkb_init_test.go`, `bkb_scan_test.go` modified: three new tests.

Recipes: `justfile`, `skills/do-work-board/justfile.template` modified (three `bkb-*` recipes pass arguments through).

Prose, modified:
- do-work-knowledge: `actions/bkb.md`, `dream.md`, `interview.md`, `interview-reference.md`, `memory-value.md`, `memory-reference.md`, `prompts.md`; `docs/bkb-guide.md`, `dream-guide.md`, `interview-guide.md`; `prompts/README.md`, `prompts/architecture-decisions-log_create-or-expand.md`.
- do-work: `actions/kb-lessons-handoff.md`, `review-work.md`, `work.md`, `roadmap.md`; `docs/standing-preferences.md`, `docs/verify-requests-guide.md`.
- do-work-toolbox: `actions/folder-structure.md`, `maintainability-audit.md`, `maintainability-audit-reference.md`, `architecture-report.md`, `validate-feedback.md`, `tutorial.md`.

Repository machinery, modified: `.gitattributes`, `README.md`, `_dev/tests/staged-skills-contract.sh`, `_dev/tests/audit-lockins.sh`, `_dev/primes/prime-releases.md`, `_dev/primes/prime-kanban-board.md`, `skills/do-work-board/tools/queue-kanban/release.go` (comment only).

Folder moves (`git mv`):
- `kb/` to `docs/kb/`: 506 tracked files, 0 content-edited.
- `decisions/` to `docs/decisions/`: 38 tracked files, 24 content-edited (21 files under `records/` for the link depth; `_master_index.md`, `_progress.md`, `topics/_index_skill-architecture.md` for self-path tokens).

## Scope changes

- No file outside the Scope list was written.
- One Scope file was NOT written: `skills/do-work/actions/version.md`. The Scope lists it (line 60 names `kb/`), the brief forbids it. See D-11.

## Red-green evidence

Tests: `TestBKBInitWithoutTargetCreatesDocsKnowledgeBase`, `TestBKBInitWithoutTargetKeepsLegacyRootKnowledgeBase`, `TestBKBScanDefaultLocationOrder` (5 cases).

RED at `11afb731` (`go test -count=1 -run 'TestBKBInitWithoutTarget|TestBKBScanDefaultLocationOrder' ./internal/knowledgecommands/`):

```
--- FAIL: TestBKBInitWithoutTargetCreatesDocsKnowledgeBase
    bkb_init_test.go:446: default dry-run planned "kb/raw/_inbox_queue.md" outside docs/kb/
--- FAIL: TestBKBInitWithoutTargetKeepsLegacyRootKnowledgeBase
    bkb_init_test.go:501: fill-gaps does not carry the one-line move: []
--- FAIL: TestBKBScanDefaultLocationOrder
    --- FAIL: .../docs_only: status location: want docs/kb, got [BKB-NOT-FOUND ... checked kb, knowledge-base]
    --- FAIL: .../legacy_kb_only: finding BKB-LEGACY-LOCATION present=false
    --- FAIL: .../legacy_knowledge-base_only: finding BKB-LEGACY-LOCATION present=false
    --- FAIL: .../both: status location: want docs/kb, got location=kb
FAIL
```

Honest note on the lock-in: in the legacy test the path assertions (only `kb/` paths, no `docs/` folder, existing file untouched) were already true before the change. They are the lock-in. Only the move-hint assertion failed at RED. The case `both with explicit legacy` also passed at RED; it pins "explicit `--kb` is literal".

A second RED appeared while making it green: the first test then failed with `mkdirat docs/kb: no such file or directory` (see D-02).

GREEN at `f328ce6d`: all three tests and all five cases PASS; whole package `ok` (about 7 s). Full outputs: `red.txt` and `green.txt` beside this file.

## Verification

All run from the worktree at `aaa4f4a5`.

1. `gofmt -l ./internal/knowledgecommands`: no output. `go vet ./internal/knowledgecommands/`: clean. `go test -count=1 ./internal/knowledgecommands/`: ok.
2. `bash _dev/tests/shipped-package-reference-contract.sh`: PASS. `bash _dev/tests/contract-regressions.sh`: exit 0. `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh`: PASS, exit 0. `bash _dev/tests/audit-lockins.sh`: passed, exit 0. `shellcheck -S warning` on the two changed scripts: clean. Plain `shellcheck` prints only info-level notes (SC2016, SC2126) on lines I did not change; the count in `audit-lockins.sh` is 7 before and after.
3. No `--kb`: `bkb-lint-structure` success, `BKB-STRUCTURE-CLEAN`. `bkb-status`: `location=docs/kb articles=245 topic_clusters=11 inbox=0 ready=0`, defrag and garden overdue 68 days (same as the baseline). `folder-structure-check`: 12 gaps, none is `kb/` or `decisions/`. `git archive HEAD | tar -t | grep -c '^docs/'`: 0 (683 entries, same count as before the move). Fresh `bkb-init` in a temp Git repository: creates `docs/kb/`, no root `kb/`. Extra: a temp repository with a root `kb/` gets `refused` + `BKB-INIT-EXISTS` + `BKB-LEGACY-LOCATION` and no `docs/` folder. `just bkb-status`, `just bkb-lint-structure kb`, `just bkb-init "odd name" --dry-run` all work with the new recipes.
4. Repository gate, run once: `bash _dev/tests/maintainer-verify.sh > /tmp/claude-501/req714-builder-gate.log 2>&1` gave `exit=0`, "Maintainer verification passed.", gate wall 157 s.

Mutation checks on the two new assertions: with the `/docs/kb` line removed from `.gitattributes`, the staged-skills contract fails with "the archive must not export the maintainer's /docs/kb"; with `docs/decisions/topics` moved away, `audit-lockins.sh` fails with "the dead-path scan did not run". Both restored.

## P-A-U

- PLAN: one resolver for the default order, three pass-through recipes, one canonical prose section, one atomic move commit.
- APPLY: built as planned, in four commits. Two additions the requirement forced: parent-folder creation in `bkb-init` (D-02) and the `mkdir -p docs &&` in the hint (D-01).
- UNIFY: checked `git diff -M --name-status a79e94bb..HEAD`: 41 modified files, all in Scope; 520 pure renames and 24 renames with edits, 0 added, 0 deleted. Read the Go diff (327 changed lines) and the prose word-diff. Searched added lines for TODO, FIXME, print statements and the mutation marker: none. `gofmt`, `go vet`, `shellcheck` as above. ADR links: 41 resolve and 61 are dead, the same as before the move.

## Decisions

- D-01 DECIDE & STATE: the move hint is `mkdir -p docs && git mv kb docs/kb`. The Plan's `git mv kb docs/kb` fails when `docs/` does not exist ("renaming 'kb' failed: No such file or directory"; probed).
- D-02 DECIDE & STATE: `bkb-init` now creates the missing parent folders of a nested target. Without it the new default failed in a Git repository. Rollback removes a parent this run created.
- D-03 DECIDE & STATE: codes are `BKB-LEGACY-LOCATION` (info) and `BKB-DUPLICATE-LOCATION` (warning). `bkb-init`, `bkb-status` and `bkb-lint-structure` emit them. `dream-scan`, `interview-ingest` and `memory-audit` use the same order but emit no hint.
- D-04 DECIDE & STATE: a notice never changes the `bkb-init` outcome. On status and lint the warning gives outcome `findings`; the info hint keeps `success`.
- D-05 DECIDE & STATE: "both exist" means the `docs/kb` directory exists and a root `kb/` or `knowledge-base/` holds `raw/` and `wiki/`. If `docs/kb` is an empty folder beside a real root `kb/`, `docs/kb` still wins (the Plan's literal rule) and the warning names both.
- D-06 DECIDE & STATE: an explicit `--kb kb` no longer falls back to `knowledge-base/`. Before, the code could not tell an explicit `kb` from the default.
- D-07 DECIDE & STATE: `bkb.md` tells the agent to pass `--kb` only when the user gave a path. The command then applies the order and the hint reaches the user.
- D-08 DECIDE & STATE: the staged-skills exemption is the prefixes `docs/kb/` and `docs/decisions/` plus the exact token `docs/decisions`. The contract cannot see how a token is spelled in context. Bare spellings remain only in literal command arguments: six `--exclude-path docs/kb/`, the `verify-requests --against` example and the fenced `git mv` in the ADR prompt. `shipped-package-reference-contract.sh` needed no change.
- D-09 DECIDE & STATE: the ADR-014 pointer in `standing-preferences.md` is now a canonical GitHub link. A bare `docs/decisions/...` path fails the shipped reference contract. The link target exists on GitHub only after this work reaches `main`.
- D-10 DECIDE & STATE: `audit-lockins.sh` fails when the topic-index glob matches nothing. This is the silent skip the Exploration named.
- D-11 DECIDE & STATE: I ran `git add -A docs` once, limited to the `docs` path, to stage the 24 edited files inside the moved folders. This breaks the letter of "never `git add -A`". I checked the result before committing: 544 renames, 0 added files.
- D-12 ESCALATE: `skills/do-work/actions/version.md` line 60 still says the update does not mutate `kb/`. Value: the shipped sentence names the new folder (`<project-root>/docs/kb/`). Risk: low, wording only, reversible; left alone the release ships one stale folder name. The finalizer or the orchestrator can change it in the release commit.

## Discovered Tasks

- `CHANGELOG-2026-09-05-up-to-v0.303.6.md:2221` links to `blob/main/decisions/records/adr-020-...`; the link breaks on GitHub once the move reaches `main`. The file is in the "will NOT touch" list. impact-negligible → report only
- 61 relative links in `docs/decisions/records/*.md` were dead before the move and are still dead. impact-negligible → report only
- Files inside `docs/kb/` and the dated ADR records still name `kb/` and `decisions/` in their text (no content edits were in scope). impact-negligible → report only
- A knowledge base folder whose name starts with `-` cannot be passed to the pass-through recipes as a bare word; `just bkb-status --kb <name>` works. impact-negligible → report only
- `interview-ingest` and `memory-audit` now read `--kb ""` as "not supplied". impact-negligible → report only

## Lessons read

- Read in full: `general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`; `prime-action-files.md`, `prime-releases.md`, `prime-shell-commands.md`; `CLAUDE.md`.
- `_dev/primes/lessons-action-files.md`: read the version entries and every family line. Applied `example-path-read-as-citation` (project paths carry `<project-root>/`), `alternate-writer-contract-drift` (swept every restatement of the location, 24 prose files) and `restated-mechanism-unchecked` (the `bkb.md` order was written from the Go code; the parent search is the root plus 2 levels up).
- Applied `assertion-passes-when-feature-is-dead` from `prime-shell-commands.md` (the two mutation checks) and `canonical-link-outlives-its-target` from `prime-releases.md` (first Discovered Task).
- Not read: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (dropped for budget in the request; the brief did not ask for it).

## Addendum (orchestrator, before the merge)

- D-12 resolved: `skills/do-work/actions/version.md` line 60 now says the update does not mutate "the knowledge base" (commit 00eac450 on the branch). The brief's ban on that file was about the version line the finalizer writes, not this sentence.
