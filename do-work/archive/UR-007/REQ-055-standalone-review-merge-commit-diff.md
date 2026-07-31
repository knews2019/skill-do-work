---
id: REQ-055
title: Standalone review of a worktree-merged REQ must diff the merge commit against its first parent
status: completed
claimed_at: 2026-07-29T14:46:44Z
commit: 93424eb
status_changed_at: 2026-07-29T15:05:02Z
route: B
domain: general
tdd: false
maintenance: false
prime_files: []
created_at: 2026-07-29T09:32:07Z
user_request: UR-007
addendum_to: REQ-037
depends_on: []
write_set: ["actions/review-work.md", "actions/present-work.md", "actions/pipeline.md", "actions/ai-report.md", "actions/pipeline-reference.md"]
---

# Standalone review of a worktree-merged REQ must diff the merge commit against its first parent

## What

In `actions/review-work.md`'s "Get the Diff" step, detect when the REQ's `commit:` hash is a merge commit and diff it as `git show --first-parent -m <commit>` (or `git diff <commit>^1..<commit>`) instead of plain `git show <commit>`.

## Why

In worktree dispatch mode a finished REQ's `commit:` frontmatter holds the `--no-ff` merge commit (established by REQ-037). Plain `git show` on a merge commit prints a combined diff that is usually empty, so a standalone re-review (`do-work review REQ-NNN`) of worktree-merged work finds nothing to review. Pipeline-mode review reads the merge range correctly and is unaffected. Approved by the user via `do-work clarify` on 2026-07-29 (follow-up from REQ-037's review, surfaced in REQ-042).

## Constraints

- Scope the change to the merge-commit case — serial-mode `commit:` hashes are ordinary commits and must keep today's behavior.
- Grep for other consumers of the `commit:` hash as a diff source before calling it fixed (prescribed-command fixes are rarely local — restated from the repo's maintainer doc).

## Acceptance

- Standalone review of a worktree-merged REQ shows the merged change; standalone review of a serial REQ is unchanged.

## Addendum (2026-07-29, UR-008)

A follow-up deep review of the REQ-035–040 batch ran the "grep for other consumers" sweep this REQ's constraint asks for, so the builder doesn't have to rediscover the sites. Known `git show <commit>` consumers that go empty on a merge hash, all in scope for the same fix pattern:

- `actions/review-work.md` — the "Get the Diff" step already named above.
- `actions/present-work.md` — two sites (the explainer's diff read and the receipt block).
- `actions/pipeline.md` — one site.
- `actions/ai-report.md` — the "Verify It Yourself" receipt block (~:212). Note `ai-report.md` ~:70 already uses `git diff-tree … -m` and is fine.

Verified fine as-is: `actions/inspect.md` uses `git show <commit>:<path>` (file-at-commit, works on merges). The `write_set` should grow to cover the files above when Step 5.5 firms it. This extends the original scope; it contradicts nothing.

## Implementation Summary

Every site that reads a REQ's `commit:` hash as a diff source now branches on the merge case using one shared idiom: detect with `git rev-parse --verify -q <sha>^2` (succeeds ⇒ the commit has a second parent ⇒ it is worktree dispatch mode's `--no-ff` merge), and in that case diff against the first parent with `git show --first-parent -m <sha>`. The addition is purely conditional prose appended to each existing `git show <sha>` prescription, so serial-mode ordinary commits keep today's behavior — verified on this repo: `git rev-parse ^2` fails on an ordinary commit, and on merge `e1f6cf6` plain `git show` emits 7 header-only lines against 4028 lines for the first-parent form. The Addendum's sweep listed two sites in `present-work.md`; a re-verify found a third (the interactive explainer's "For the developer" receipt block at the Content bullet), which carries the identical defect and is fixed the same way. `_dev/tests/contract-regressions.sh` passes.

Adversarial review returned FIX-THEN-PASS and four fixes landed on top: every `^2` argument is now single-quoted (`git rev-parse --verify -q '<sha>^2'`) because bare `^` is a glob operator under zsh `EXTENDED_GLOB` and an escape character in cmd.exe — either one makes the probe fail and silently read as "not a merge", fail-closed into exactly the bug this REQ fixes; review-work's Two-Modes Standalone cell now points at Step 4 for the merge case; `actions/pipeline-reference.md`'s "How to verify" rendering template was brought into agreement with pipeline.md's data bullet; and the trap was added to the maintainer doc's prescribed-shell-command catalog.

**Files changed:**

- `actions/review-work.md` — Step 4's Standalone-mode paragraph gains the merge branch, citing `actions/work-reference.md` → **Worktree Dispatch Mode (Step 1)** for why the hash is a merge; the Two-Modes table's Standalone cell now points at Step 4 for the merge case, mirroring how the Pipeline row carries its own worktree variant inline.
- `actions/present-work.md` — three sites: Step 3's diff read, the client-brief template's "How to Verify" bracketed instruction (kept inside the brackets so the note guides the generating agent rather than landing literally in the brief), and the interactive explainer's `git show <sha>` receipt block.
- `actions/pipeline.md` — the Completion Report's **How to verify** data bullet now tells the renderer to emit the first-parent form for merge SHAs.
- `actions/ai-report.md` — the **Verify It Yourself** required-section spec (~:212) gains the same conditional; the already-correct `git diff-tree … -m` idiom at ~:70 was left alone.
- `actions/pipeline-reference.md` — the **How to verify** rendering template's "Inspect each commit" step gains a bracketed conditional instruction, so the template no longer shows a bare `git show <sha>` that contradicts pipeline.md's data bullet.
- `CLAUDE.md` — one new bullet in **Prescribed Shell Commands Must Surface What the Steps Consume** records the merge-commit/empty-combined-diff trap and the quoted detection idiom, next to the related `git show --name-only` bullet.

## Review

Adversarial review workflow (4 Opus lenses -> 2 diverse refuters per Important+ finding; 10 agents): verdict FIX-THEN-PASS -> fixes applied -> PASS.

- Upheld + fixed: unquoted `^2` merge-detection was fail-closed under zsh EXTENDED_GLOB / cmd.exe (quoted at every site, repo-wide grep clean); Two-Modes table's Standalone "How to get the diff" cell still prescribed the bare `git show` (now points at Step 4's first-parent form).
- Adopted minors: pipeline-reference.md's rendering template agreed with pipeline.md's data bullet; CLAUDE.md shell-trap catalog gained the merge-commit trap.
- Skipped (recorded, deliberate): `--first-parent -m` git-version-floor note (acceptable floor); contract-regressions pin for the four-file idiom (YAGNI).
- 1 further Important finding killed by refutation (duplicate of the table-cell finding from another lens).
