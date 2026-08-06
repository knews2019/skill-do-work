---
id: REQ-121
title: One uncommitted-changes inventory, one REQ-association pass
status: completed
created_at: 2026-08-06T13:19:18Z
claimed_at: 2026-08-06T13:21:00Z
route: B
completed_at: 2026-08-06T13:35:00Z
commit: 167b0ae
user_request: UR-022
domain: general
prime_files: []
tdd: false
depends_on: []
maintenance: false
related: [REQ-111, REQ-112, REQ-114]
batch: census-durable-findings
---

# One Uncommitted-Changes Inventory, One REQ-Association Pass

## What

Candidate B of REQ-114, split out and approved. Two primitives are currently copy-pasted as prose across several action files; each becomes one shipped script under `tools/checks/`, and the action prose calls it with its manual procedure documented as the fallback.

**Primitive 1 — the uncommitted-changes inventory.** The `git rev-parse --git-dir` gate, then `git status --porcelain --untracked-files=all`, then M/A/D categorization, then four secret-shaped exclusion globs (`.env*`, `credentials*`, `*.pem`/`*.key`/`*.p12`/`*.pfx`, `*secret*`).

**Primitive 2 — the REQ-association pass.** Glob archived REQs, read `commit:` and a terminal-success `status`, parse the `## Implementation Summary` file list, path-match against a candidate set, tie-break on the latest `completed_at`.

## Why

Re-ran REQ-114's greps on 2026-08-06 rather than trusting the census figures, exactly as REQ-114 requires:

- `grep -rln "untracked-files=all" actions/` → **5 files** (`commit.md`, `inspect.md`, `stray-check.md`, `tidy-repo.md`, `work.md`). Census said 4, and the membership moved: `validate-feedback.md` no longer matches, `stray-check.md` and `work.md` now do.
- `commit.md` Step 1 and `inspect.md` Step 1 carry the `-uall` rationale **word-for-word identical**, and the same four exclusion globs. Their Step 3 association passes are near-verbatim too.

The `-uall` flag is load-bearing, not cosmetic: without it `git status --porcelain` collapses a wholly-untracked directory to a single `?? dir/` row and every file inside escapes the secret-shaped exclusion scan. That is a secret-leak path and `actions/stray-check.md`'s Red Flags record that it has been hit.

**The copies have already drifted, which is the argument.** `commit.md`'s status check accepts both `completed` and `completed-with-issues`, names the accelerator and its floor, and warns in its own Red Flags that testing only for the literal `completed` drops every remediated-with-issues REQ. `inspect.md`'s copy names both values but carries none of the rest. One of the two copies learned something the other did not — with two more sites now touching the same primitive, that gap only widens.

## Detailed Requirements

- `tools/checks/uncommitted-inventory.sh` — emits one machine-readable row per changed path, tagged `M`/`A`/`D`, with secret-shaped paths tagged `X` rather than dropped (both callers must still report them; silently skipping is the failure mode both prose copies already warn about).
- `tools/checks/associate-files.sh` — takes candidate paths, emits `REQ-NNN<TAB>path` for associated and `-<TAB>path` for unassociated.
- Terminal-success matching must honor the Schema Read Contract's `status` aliases (`done`/`finished`/`closed` → `completed`), not just the two literal values. This is the drift `commit.md` warns about, fixed once in the script instead of twice in prose.
- `do-work/` metadata paths stay excluded from association candidates, matching `tools/checks/scope-drift.sh`.
- Update `actions/commit.md` and `actions/inspect.md` to call the scripts, each keeping its manual procedure as the documented fallback. Leave `stray-check.md`, `tidy-repo.md` and `work.md` alone — they use `-uall` for their own narrower purposes and are not copies of this inventory.
- Register both scripts in `_dev/tests/contract-regressions.sh`'s `hardened_check_scripts` list so a prose pointer at a missing script fails the suite.

## Constraints

- Shell only — `tools/checks/*.sh` needs no compiled-tooling exception. `actions/board.md` remains the only capability allowed to need a compiler.
- No action prose may lose its documented fallback procedure.
- Neither script may write anything. Both are read-only reporters.

## Red-Green Proof

**RED:** create a wholly-untracked directory containing a secret-shaped file (`newdir/.env`). Under plain `git status --porcelain` the inventory sees one row, `?? newdir/`, and `newdir/.env` never reaches the exclusion scan — the secret-leak path.

**GREEN:** `tools/checks/uncommitted-inventory.sh` lists `newdir/.env` individually and tags it `X`.

**Validation:** User approved Candidate B from REQ-114's three, on the safety argument.

## Full Context

Split from `do-work/queue/REQ-114-residual-extraction-candidates.md` (Candidate B). Candidates A and C remain there, unapproved.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Two read-only bash scripts in the established `tools/checks/` shape. `uncommitted-inventory.sh` emits `<tag>\t<path>` rows (M/A/D/X); `associate-files.sh` takes candidate paths on stdin and emits `<owner>\t<path>`. Rewire `commit.md` and `inspect.md` Steps 1 and 3 to call them, each keeping its manual procedure as the documented fallback. Register both in the regression suite's `hardened_check_scripts` list.
- [x] **[APPLY]:** Both scripts written and made executable; both action files rewired; suite updated. Stayed inside the declared scope.
- [x] **[UNIFY]:** `git diff --stat` reviewed. `bash -n` parses both scripts (shellcheck not installed in this container). Regression suite: 8 FAILs, identical to the clean-tree baseline — no new failures. No debug artifacts, no `set -x`, no TODOs. Verified per file: both scripts (behavior tested, see Testing), `commit.md` + `inspect.md` (fallback prose present in all four rewired seams), `contract-regressions.sh` (4 new entries, all passing).

## Scope

**Files I will touch:**
- `tools/checks/uncommitted-inventory.sh`
- `tools/checks/associate-files.sh`
- `actions/commit.md`
- `actions/inspect.md`
- `_dev/tests/contract-regressions.sh`

**Files I will NOT touch:** `actions/stray-check.md`, `actions/tidy-repo.md`, `actions/work.md` — they use `-uall` for their own narrower purposes and are not copies of this inventory.

**Acceptance criteria:** the `-uall` secret-leak path is closed by a tested script; both callers keep a documented fallback; no new regression-suite failures.

## Implementation Summary

**Files changed:**
- `tools/checks/uncommitted-inventory.sh` (new)
- `tools/checks/associate-files.sh` (new)
- `actions/commit.md` (modified)
- `actions/inspect.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Extracted the two copy-pasted primitives into shipped scripts. The inventory script gates on `git rev-parse --git-dir`, reads `git status --porcelain --untracked-files=all -z`, and tags each path M/A/D/X — secret-shaped names are tagged `X` and still reported, never dropped. The association script scans archived and in-flight REQs, honors the Schema Read Contract's `status` aliases, and tie-breaks contended paths on `completed_at`. Steps 1 and 3 of both `commit.md` and `inspect.md` now call the scripts with their manual procedures retained as fallbacks. Four entries added to the suite's `hardened_check_scripts` list.

## Decisions

- **D-01** — A renamed path is tagged `M`, not `A`. DECIDE & STATE. The callers treat `M` as "read the diff" and `A` as "read the whole file as new content"; a move is the former. Reversible, single-line.
- **D-02** — Output is line-oriented, so a filename containing a literal newline produces a row spanning lines. DECIDE & STATE. The `-z` read means such a path is never corrupted, only unrenderable on one line; the consumers are prose steps read by an agent, for which line output is right. Documented as a known limit in the script header rather than papered over.
- **D-03** — Read `git status` via process substitution rather than `$(...)`. Forced, not chosen: bash strips NUL bytes from command substitution, which collapses every `-z` record into one. A pipe would also have worked but runs the loop in a subshell, so the emitted-row counter would always read 0.
- **D-04** — An archived REQ outranks an in-flight one for a contended path. DECIDE & STATE. A finished REQ is stronger evidence of ownership than a running one; implemented by stamping in-flight REQs with an empty `completed_at`, which sorts below every archived stamp.

## Testing

**Tests run:** behavior probes against purpose-built scratch repos, plus an end-to-end run against this repo's live working tree.

**Red-Green Proof — the secret-leak path:**
- **RED:** in a repo with an untracked `newdir/` containing `.env`, plain `git status --porcelain` returns two rows — `?? newdir/` and `?? top.txt`. `newdir/.env` is invisible and never reaches the exclusion scan.
- **GREEN:** `uncommitted-inventory.sh` returns `X newdir/.env`, `A newdir/app.js`, `A top.txt`. The secret is listed individually and tagged.

**Inventory — other paths verified:** unstaged modify → `M`; staged rename → `M` on the new path with the origin field correctly consumed (no phantom row for the old path); unstaged delete → `D`; secret-shaped tracked file → `X`; paths containing spaces → intact; clean tree → exit 1; non-git directory → exit 2.

**Association — verified:** contended path goes to the later `completed_at` (REQ-002 over REQ-001); `status: done` alias associates (REQ-003) — the case both prose copies would have dropped; `status: failed` does not associate; in-flight `working/` REQ associates despite `status: claimed`; archived beats in-flight on a contended path; `do-work/` metadata never associates; empty stdin → exit 1; missing `do-work/` → exit 2; bad flag → exit 2.

**End-to-end:** run against this repo's live dirty tree, the inventory returned 7 rows and the association correctly attributed `_dev/tests/contract-regressions.sh` → REQ-104 and `actions/commit.md` → REQ-113 from real archived REQs.

**Result:** all pass. Regression suite at 8 FAILs, unchanged from the clean-tree baseline.

## Lessons Learned

- **`$(...)` silently destroys NUL-delimited output in bash.** Any prescribed command that reaches for `git ... -z` for path safety must read it through process substitution or a pipe; capturing it in a variable removes the very delimiter that made it safe. Worth checking wherever the skill prescribes `-z`.
- **Porcelain rename records carry a second field.** `R  new\0old\0` — a consumer that reads one field per record parses the origin path as the next record's status bytes and shifts every subsequent row. Any hand-rolled `-z` parse in prose has this bug unless it explicitly discards the origin.
- **The drift was already visible before the extraction.** `commit.md` had learned that a `status: done` REQ must still associate; `inspect.md` had not. Two copies of one primitive do not stay two copies of the same primitive — that is the argument for extraction, and it was observable in the diff rather than hypothetical.
