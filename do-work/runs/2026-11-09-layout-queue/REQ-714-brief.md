# Builder Brief — REQ-714 (suite writes its knowledge base and ADR log under docs/)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The orchestrator merges, tests, reviews and releases.

- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-714-suite-writes-kb-and-decisions-under-docs`
- Branch: `worktree-agent-REQ-714-suite-writes-kb-and-decisions-under-docs` (already checked out there)
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its Plan, Exploration and Scope: read `do-work/working/REQ-714-suite-writes-kb-and-decisions-under-docs.md` from the MAIN tree (your worktree holds an older copy without those sections).
- Hand-back file (write it OUTSIDE both trees): `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/eb6add15-e307-4d23-8848-ad3f1445c6a7/scratchpad/run/714/handback.md`

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `_dev/primes/prime-action-files.md` and `_dev/primes/prime-releases.md`, and `_dev/primes/prime-shell-commands.md` before you edit a `.sh` file or a justfile recipe.
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).
4. `_dev/primes/lessons-action-files.md` before editing action files (it was dropped from the required reads only for budget; skim the family slugs and read what applies).

## What to build
Follow the `## Plan` in the request: four commits, in this order. Decisions already made by the orchestrator:

- **Compatibility rule** exactly as the Plan states it, including the warning finding when both `docs/kb` and a root `kb/` exist (`docs/kb` wins). Finding codes follow the existing `BKB-*` naming. An explicit `--kb` is always literal.
- **Recipes** become pass-through so they stop hard-coding the old path. Keep the recipe names. Check that every `NextJustRecipe` or next-step string the Go code emits still works with the new recipes.
- **Prose**: `skills/do-work-knowledge/actions/bkb.md` § Locating the Knowledge Base is the one canonical statement of the order. Other files name the new default or point there; do not restate the full rule in many places. In shipped Markdown write project paths as `<project-root>/docs/kb/` and `<project-root>/docs/decisions/`, because `_dev/tests/shipped-package-reference-contract.sh` and `_dev/tests/staged-skills-contract.sh` read a bare `docs/...` as a path inside the skill package. Add an exemption in the staged-skills contract only for spellings that cannot carry the prefix (a literal command argument, a folder tree).
- **ADR prompt**: resolve the ADR root in its first phase with the same order (docs first, a legacy root folder plus the one-line move hint, else create under docs).
- **Export-ignore assertion**: `_dev/tests/contract-regressions.sh` holds no such assertion any more and sits at its line ceiling. Put the assertion for the two new `.gitattributes` lines in `_dev/tests/staged-skills-contract.sh` beside its existing export-ignore checks, and correct the stale sentence in `.gitattributes` that names the old home. If you find the assertion already exists somewhere, update that one and add nothing.
- **ADR records**: rewrite the relative link depth in `docs/decisions/records/*.md` (one level deeper) so the 41 links that resolve today still resolve. Do it uniformly; do not repair the links that are already dead.
- **Commit 4 is atomic**: both `git mv` moves, `.gitattributes`, the staged-skills exclusions, `audit-lockins.sh`, the link fixes and the remaining repository references land in ONE commit.
- **Not yours**: `VERSION`, both `CHANGELOG.md` files and `skills/do-work/actions/version.md` (the finalizer writes the release); anything under `do-work/` and `ai-reports/`; the files in the Scope's "will NOT touch" list.

## Rules
- TDD: commit 1 holds only the failing tests. Run them and record the failure output. Commit 2 makes them pass. Keep the lock-in test for the legacy root folder.
- Write only inside the Scope's file list plus the two folder moves. If you need another file: when the request's own requirements need it, proceed and report it in the hand-back under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`; for the moves use `git mv`). Never `git add -A`. End each commit message with the line `Co-Authored-By: Claude Fable 5.2 <noreply@anthropic.com>`. Commit subjects start with `[REQ-714]`.
- Do not push. Do not merge. Do not start subagents.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l ./internal/knowledgecommands`, `go vet ./internal/knowledgecommands/`, `go test -count=1 ./internal/knowledgecommands/`.
2. `bash _dev/tests/shipped-package-reference-contract.sh`, `bash _dev/tests/contract-regressions.sh`, `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/staged-skills-contract.sh`, `bash _dev/tests/audit-lockins.sh`, and `shellcheck` on every `.sh` you changed.
3. With no `--kb`: `bkb-lint-structure` and `bkb-status` through `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format text <command>` report `location=docs/kb` and the same counts as the baseline in the request's Exploration. `folder-structure-check` no longer lists `kb/` or `decisions/`. `git archive HEAD | tar -t | grep -c '^docs/'` prints 0. A fresh `bkb-init` in a temp Git repository creates `docs/kb/`.
4. Last, the repository gate once: `bash _dev/tests/maintainer-verify.sh > /tmp/claude-501/req714-builder-gate.log 2>&1; echo "exit=$?"`. Never pipe it. It takes about two minutes. If it fails on something your change caused, fix it and rerun; if it fails only on a 30-second per-file time budget, rerun once.

## Hand-back file content
Use these headings: `## Commits` (hash and subject each), `## Files` (every file created, modified or deleted outside the two moved folders, with the verb; then one line for each folder move and the count of content-edited files inside), `## Scope changes`, `## Red-green evidence` (test names, the failing output before, the passing run after), `## Verification` (each command above with its result), `## P-A-U` (three short lines: PLAN, APPLY, UNIFY, the UNIFY line listing what you checked in `git diff --stat`), `## Decisions` (D-01.. each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE), `## Discovered Tasks`, `## Lessons read`.

Your final message: under 2500 characters. State the branch tip hash, the gate result, and anything the orchestrator must decide.
