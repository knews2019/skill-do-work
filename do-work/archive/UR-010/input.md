---
id: UR-010
title: Harden the Step 9 commit-hash write-back so it can never blank a REQ
created_at: 2026-07-30T21:57:34Z
requests: [REQ-062, REQ-063, REQ-064]
word_count: 1010
---

# Harden the Step 9 commit-hash write-back so it can never blank a REQ

## Summary

A downstream consumer repo (`game-find-the-difference`) lost six archived REQ files —
`REQ-1282, 1285, 1287, 1288, 1289, 1290`, 8 KB to 26 KB each — to do-work's own Step 9
commit-hash write-back. Each file was complete at its implementation commit; the immediately
following `[REQ-NNN] record commit hash <hash>` metadata commit replaced the whole file with
nothing. Every one of those commits read `1 file changed, N deletions(-)` and nothing looked at
that number before it landed. The commit message claimed success, so the loss stayed invisible
until a board regeneration surfaced `unrecognized status ""` warnings much later.

Root cause: `actions/work-reference.md:860-873` prescribes the write-back as prose with no
deterministic command and no post-edit verification. The mechanics of the single-line edit are
left entirely to the agent.

## Extracted Requests

| REQ | Title | Covers |
|-----|-------|--------|
| REQ-062 | Deterministic, self-verifying commit-hash write-back | Upstream items 1, 2, 3, 4, 5 + acceptance 1 and 3 |
| REQ-063 | Detect blanked and unparseable REQ files in forensics | Upstream "detection check" + acceptance 2 (detection half) |
| REQ-064 | Restore blanked archived REQs from git history | Upstream "guided repair" + acceptance 2 (repair half) |

## Batch Constraints

- **Four corrections to the upstream report, established during capture. Implement the
  corrected form, not the literal ask:**
  1. **There is no `commit: PENDING` convention in this skill.** `grep` finds zero occurrences
     repo-wide; `commit:` is written once, post-commit, and may be entirely absent. The guard
     "no longer says PENDING" is meaningless here — the real guards are *the field now equals
     the hash* and *nothing else in the file changed*. Tolerate a literal `PENDING` as an
     existing value for consumer-repo compatibility; do not introduce the placeholder.
  2. **Guarding on "terminal-*success*" is a bug.** `actions/work-reference.md:128-141` stamps
     `commit:` on *every* terminal flip, and `:854` says "Failed requests get committed too." A
     terminal-success guard would hard-FAIL every legitimate `failed`/`cancelled` REQ. Correct
     guard: `status:` is byte-identical pre→post and non-empty; WARN — never FAIL — when the
     value sits outside `{completed, completed-with-issues, failed, cancelled}`.
  3. **`git diff --numstat` cannot be the primary guard.** It compares index↔worktree and
     yields nothing when the REQ is untracked or the repo is not git — which is this repo's own
     case (`do-work/` is in `.git/info/exclude`). Primary guard is exact pre/post byte-and-line
     arithmetic, computable in closed form because the script holds the pre-image. numstat's
     unique contribution is that it is the only guard that sees damage done *before* the script
     ran; that plus a size floor against `HEAD` is what encodes the 26 KB → 0 B signature.
  4. **`do-work forensics --repair-blanked` cannot exist.** `actions/forensics.md:19-21` is a
     hard contract: "Read-only. This action never modifies files… Report, don't fix." Detection
     goes in forensics; the repair goes to `actions/cleanup.md`, which already writes to the
     archive and already has the consent-gated-pass precedent (Pass 5, orphaned worktrees).

- **Follow the shipped-guard-script pattern, not more prose.** `tools/checks/qualify.sh`,
  `preflight.sh`, `scope-drift.sh`, and `archive-collision.sh` are the house pattern for turning
  fragile prescribed shell into hardened behavior; `_dev/tests/contract-regressions.sh:191-208`
  ratchets each one's existence *and* its reference from `actions/work.md`. Match that style:
  `set -uo pipefail`, `OK:`/`NOOP:`/`WARN:`/`INFO:`/`FAIL:` stdout vocabulary, exit 2 for usage
  errors, and header comments that state *why* each guard exists so nobody simplifies it away.

- **One canonical implementation per concern, callers point at it.** The blanked-REQ detector and
  its restore path are one script with a flag, not two implementations; the restore reuses the
  write-back script rather than re-implementing the frontmatter edit.

- **`SKILL.md` has an enforced word budget (2650) and must not grow.** No new action, no new
  routing row, no new dispatch row — this work extends existing actions only.

- **Shipped files must not cite this repo's `CLAUDE.md`/`AGENTS.md`** — both are `export-ignore`d,
  so the citation dangles downstream.

- Every REQ gets its own commit, `CHANGELOG.md` entry, and version bump (currently 0.150.15).

- **Local override:** `do-work/` is git-excluded in this repo via `.git/info/exclude`, so commit
  steps stage nothing under it and metadata commits are skipped here.

## Full Verbatim Input

# Upstream prompt — harden do-work's "record commit hash" write-back so it can never blank a REQ

Paste the block below to Claude working in the `knews2019/skill-do-work` repo. It reports a real
data-loss bug we hit, its root cause, and asks for a robust (deterministic + guarded) fix plus a
forensics repair path for already-blanked files.

---

## Prompt to paste upstream

There is a data-loss bug in Step 9's commit-hash write-back. In one work session it silently
truncated **six archived REQ files to 0 bytes** (`REQ-1282, 1285, 1287, 1288, 1289, 1290`). Each
file was complete at its implementation commit, then the immediately-following
`[REQ-NNN] record commit hash <hash>` metadata commit replaced the entire ~120-line file with
nothing. The only reason we recovered them was that the content still existed in git history at the
implementation commit; a `do-work cleanup` or a `git gc` away from unrecoverable. On the queue-kanban
board the blanked files surfaced as `unrecognized status ""` warnings, parked under Needs input /
Blocked as *untitled*, because an empty file has no `status:` frontmatter.

### Root cause

`actions/work-reference.md` → **Commit & Metadata-Commit Procedure (Step 9)**, lines ~860–875. The
step is a **prose instruction with no deterministic command and no post-edit verification**:

> "Write that hash into the archived REQ's `commit:` field (replace the "commit:" line, or add it if
> missing), then commit the edit."

The mechanics of the single-line edit are left entirely to the agent. Any of these plausible agent
moves blanks or corrupts the file, and **nothing catches it before the metadata commit lands**:

- a `Write`/heredoc/`echo >` that intends to rewrite one line but emits an empty or partial body;
- a `sed -i`/`perl -i` whose pattern doesn't match, combined with a shell quoting error;
- an editor that reads a stale/empty buffer.

The metadata commit then stages and commits the damage with a message that *claims success*
(`record commit hash <hash>`), so the failure is invisible until a board regeneration much later.

### What I want: make this step impossible to get wrong

Please redesign the write-back so it is **deterministic and self-verifying**, and so the metadata
commit **refuses to land** if the edit did anything other than change the one `commit:` line. Concretely:

1. **Deterministic single-line edit, not free-form.** Specify one canonical way to set the field that
   edits *only* the `commit:` line and cannot rewrite the body. Prefer a tool-based exact-string edit
   (match `commit: PENDING` — or the existing value — on the frontmatter line and replace only that),
   or if the skill must stay shell-only, a guarded in-place edit that operates on a temp file and is
   swapped in only after the guards below pass. Whichever is chosen, the hash is still re-typed as a
   literal (shell vars don't survive across the model round-trip between the edit and the commit —
   keep that existing warning).

2. **Pre-commit guards that abort on any anomaly.** Before `git add`/`git commit` of the metadata
   commit, assert ALL of:
   - the REQ file is **non-empty** (`test -s`);
   - `git diff --numstat -- <req>` shows the change is **tiny** — on the order of `1 insertion, 1
     deletion` (or `1 insertion, 0 deletions` when adding a missing field). A mass deletion
     (e.g. `0  126`) means the body was destroyed → **abort, do not commit**;
   - the frontmatter still contains `status:` with a terminal-success value and the same `id:`;
   - the `commit:` line now equals the intended `<hash>` and no longer says `PENDING`.
   If any check fails, stop and report — never commit a REQ whose diff deletes its body.

3. **Idempotency.** Running the write-back twice (or on a REQ whose `commit:` is already the hash)
   must be a no-op, not a second corruption. If `commit:` already equals `<hash>`, skip the edit and
   the metadata commit entirely.

4. **Apply the same guards in worktree-dispatch mode**, where the hash written is `<merge_hash>`
   rather than `git rev-parse --short HEAD` — that path uses the identical write-back and is equally
   exposed.

5. **Mirror the one-line summaries** in `actions/work.md:627` and `:629` so the skeleton and the
   reference don't drift.

### Also: a forensics check + repair for already-blanked REQs

`actions/forensics.md` currently warns on unrecognized `status` values but has no notion of a
**0-byte / body-destroyed archived REQ**, and no repair. Please add:

- A **detection check**: any `do-work/archive/**/REQ-*.md` (or UR-*.md) that is 0 bytes, or that has
  no parseable frontmatter, is reported as a data-loss anomaly — distinct from the generic
  "unrecognized status" warning, because the remedy is different (recover content, not edit a status).
- A **guided repair**: for each blanked file, locate the last commit where it was non-empty
  (`git log --diff-filter=A` for creation, or walk `git log --format=%H -- <file>` and pick the last
  commit whose `git cat-file -s <commit>:<file>` is > 0), restore that content, and re-apply the
  `commit:` field using the hash recorded in the offending `record commit hash` commit message. This
  is exactly the manual recovery we had to perform by hand for the six files above; it should be a
  one-command `do-work forensics --repair-blanked` (or similar), consent-gated.

### Acceptance

- A simulated bad edit (empty body) during Step 9 is caught by the guards and the metadata commit is
  refused, with a clear operator message — add a test that asserts this.
- `do-work forensics` flags a 0-byte archived REQ and its repair restores the exact prior content
  plus the correct `commit:` hash — add a test with a fixture.
- Normal (correct) write-back still produces a `1 insertion, 1 deletion` metadata commit and passes.

---

## Appendix — evidence from our repo (game-find-the-difference)

| REQ | impl commit (content) | metadata commit that blanked it | prior size |
|-----|-----------------------|----------------------------------|-----------|
| 1282 | `9617040b` | `091a67cc` "record commit hash 9617040b" | 9078 B → 0 B |
| 1285 | `dcca5026` | `004dab14` | ~8.8 KB → 0 B |
| 1287 | `43a1694c` | `fe18ebe5` | ~8.0 KB → 0 B |
| 1288 | `e2554361` | `7ae09da9` | ~11.3 KB → 0 B |
| 1289 | `6b78f855` | `a5a76a48` | ~26 KB → 0 B |
| 1290 | `c5091502` | `7e2d0e78` | ~14.2 KB → 0 B |

Each metadata commit's diff was `1 file changed, N deletions(-)` (whole file removed) — the exact
signature guard #2 above would have caught.

Recovery we ran by hand (for reference — this is what the `--repair-blanked` path should automate):

```bash
# restore body from the implementation commit, then fill commit: with the recorded hash
git show <impl_commit>:<archive/REQ-file> > <archive/REQ-file>
# replace only the frontmatter line:  commit: PENDING  ->  commit: <impl_commit>
```

---
*Captured: 2026-07-30T21:57:34Z*
