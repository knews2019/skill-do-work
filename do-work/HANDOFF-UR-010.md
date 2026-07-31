# Handoff — UR-010: Harden the Step 9 commit-hash write-back

**Session ended:** 2026-07-30T22:40Z · **Stopped by:** user request, mid-REQ-064
**Repo:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` · branch `main`, nothing pushed
**Approved plan:** `/Users/t2/.claude/plans/upstream-prompt-prancy-falcon.md`

---

## Why this work exists

A downstream consumer repo (`game-find-the-difference`) lost six archived REQ files —
`REQ-1282, 1285, 1287, 1288, 1289, 1290`, 8 KB to 26 KB each — to do-work's own Step 9
commit-hash write-back. Each file was complete at its implementation commit; the immediately
following `[REQ-NNN] record commit hash <hash>` metadata commit replaced the whole file with
nothing. Every one of those commits read `1 file changed, N deletions(-)` and nothing looked at
that number. The commit message claimed success, so the loss stayed invisible until a board
regeneration surfaced `unrecognized status ""` warnings much later.

The upstream bug report is preserved verbatim in `do-work/user-requests/UR-010/input.md`.

**Four corrections to that report were established at capture time and are binding.** They are
written into the UR's Batch Constraints; the short version:

1. **No `commit: PENDING` convention exists here** — zero occurrences repo-wide. `commit:` is
   written once, post-commit, and may be absent. Don't introduce the placeholder.
2. **Never guard on terminal-*success*** — `failed` and `cancelled` REQs legitimately carry
   `commit:`. Guard on `status:` unchanged and non-empty; WARN, never FAIL, outside the
   terminal set.
3. **`git diff --numstat` is not the primary guard** — it is silent when the REQ is untracked,
   which is this repo's own case. Byte/line arithmetic is primary.
4. **`do-work forensics --repair-blanked` cannot exist** — forensics is contractually read-only.
   Detection lives in forensics; the repair lives in `actions/cleanup.md`.

---

## Status

| REQ | State | Commit |
|-----|-------|--------|
| REQ-062 — guarded write-back | ✅ done, archived, hash written back | `62a4188` |
| REQ-063 — forensics detection | ✅ done, archived, hash written back | `d91c567` |
| REQ-064 — cleanup restore | 🔶 **in progress**, claimed in `do-work/working/` | — |
| REQ-065 — HANDOFF.md guidance | ⏸ `pending-answers`, awaiting `do-work clarify` | — |

Version is at **0.152.0** (`actions/version.md`), changelog entries exist for 0.151.0 and
0.152.0. REQ-064 will need its own bump to 0.153.0 and its own entry.

**Baselines are green right now:**
```
bash _dev/tests/contract-regressions.sh     # passes, includes 25 behavioral probes
shellcheck tools/checks/*.sh                # clean
```

---

## What shipped (REQ-062 + REQ-063)

**`tools/checks/record-commit-hash.sh`** — the guarded write-back. Frontmatter-scoped awk edit
(a `commit:` in body prose is structurally unreachable), atomic rename behind closed-form
byte/line arithmetic, HEAD size floor that catches the 26 KB → 0 B signature, idempotent with a
stranded-edit exception, plus a read-only `--verify` mode that reads the committed blob back
(the only thing that catches a content-mutating pre-commit hook). Edits and guards — it does
**not** stage or commit; that stays with the caller.

**`tools/checks/blanked-req-scan.sh`** — read-only detector for REQ/UR files that are 0 bytes or
have lost their frontmatter. Resolves each one's recovery-source commit and the hash recorded in
the blanking commit's message. Emits a human report plus `BLANKED<TAB>path<TAB>sha<TAB>hash`
records; `--porcelain` emits only the records.

**Prose:** `actions/work-reference.md`'s Step 9 procedure now invokes the script (serial mode
collapsed to one command block, killing the hash-lost-between-blocks hazard);
`actions/work.md:627`/`:629`/`:686` mirror it; `actions/forensics.md` gained Critical check 13
and check 11 now skips files with no parseable frontmatter so damage is reported once, with the
remedy that fits.

**Tests:** `_dev/tests/record-commit-hash-guards.sh` — the repo's first behavioral git fixture,
invoked from `contract-regressions.sh` (there is **no auto-discovery**; an uninvoked probe file
is dead weight that reads as coverage). `hardened_check_scripts` entries are now
`script|referencing-action-file` pairs, because the scanner is referenced from forensics, not
`actions/work.md`.

---

## REQ-064 — exactly where it stands

**File:** `do-work/working/REQ-064-restore-blanked-reqs-in-cleanup.md` (`status: claimed`,
`route: B`, `claimed_at: 2026-07-30T22:33:10Z`).

### Done and verified

- **`--restore` / `--dry-run` implemented** in `tools/checks/blanked-req-scan.sh` (uncommitted).
  Restores content from the recovery commit via temp-file + atomic rename, refuses to write
  empty recovered content, then re-applies `commit:` **by calling
  `tools/checks/record-commit-hash.sh`** — never by hand-editing frontmatter, so the guards
  come along. Skips healthy files and files with no recovery source. `--restore` without
  `--dry-run` exits 0 on a complete repair (a fixed thing is not a finding); `--dry-run` keeps
  the finding exit code.
- **6 restore probes** added to `_dev/tests/record-commit-hash-guards.sh` (uncommitted): the
  full incident reproduction, byte-identity of the restored content against the pre-blanking
  blob, the untouched-healthy-neighbour case, `--dry-run` writes nothing, and a clean re-run.
- **Suite is green** with these changes in place.

### Not done

1. **`actions/cleanup.md` — `### Pass 6: Restore Blanked Archived REQs (consent-gated)`.**
   Model it on Pass 5 (orphaned worktrees): what it finds, what it asks, what it does. Keep it
   short — it delegates to the scanner, so it must not re-specify the recovery algorithm.
   Also update, per the REQ's Detailed Requirements:
   - `## When to Use` — the "cleanup only reorganizes… never deletes" carve-out list gains a
     **third** narrow exception, stated as *restores* content;
   - `## Steps` — the "Six passes, in order" preamble becomes seven;
   - `## Commit (Git repos only)` — stage restored paths;
   - `## Reporting`, `## Archive Structure After Cleanup`, `## What This Action Does NOT Do`.
2. **`docs/cleanup-guide.md`** — a short paragraph (it is in the REQ's declared Scope).
3. **Version bump to 0.153.0** + a `CHANGELOG.md` entry. Title must say what shipped and must
   not duplicate an existing title.
4. **REQ-064's pipeline sections** — Implementation Summary, Decisions, Qualification, Testing,
   Review (including the **mandatory Restatement Sweep**), Lessons Learned, Orientation. Then
   archive to `do-work/archive/` root, commit, and write the hash back with
   `tools/checks/record-commit-hash.sh`.

### A live decision the next session must make

`_dev/tests/contract-regressions.sh` pins `tools/checks/blanked-req-scan.sh` to
**`actions/forensics.md`** in `hardened_check_scripts`. Once Pass 6 lands, `actions/cleanup.md`
also references it. Either leave the pin as is (forensics is still a valid referencing file), or
extend the entry to assert both. Leaving it is defensible; asserting both is stronger. Don't
weaken the assertion to accommodate the second caller.

---

## ⚠ Read this before running `do-work run`

`REQ-064` sits in `do-work/working/` with `status: claimed`, and **the orchestrator lock has been
released** (there is no `do-work/orchestrator-lock.json`). A cold `do-work run` will therefore
treat it as an abandoned crash artifact at Step 1 and re-queue it — resetting it to `pending` and
**stripping its orchestrator-generated sections** (`## Triage`, `## Plan`, `## Exploration`,
`## Scope`).

The **code changes on disk survive that** — only the REQ's narrative sections are lost. Two
options:

- **Continue REQ-064 directly** (recommended): acquire a fresh lock, re-claim REQ-064 in
  `claimed_reqs`, and resume at Step 6 without a full Step 1 scan.
- **Let it re-queue** and re-run the pipeline from Step 1. Cheap, but you will re-derive triage
  and scope. The Scope section is reproduced below so nothing is actually lost either way.

REQ-064's declared Scope (mirror of its `write_set`):
`tools/checks/blanked-req-scan.sh`, `actions/cleanup.md`, `docs/cleanup-guide.md`,
`_dev/tests/record-commit-hash-guards.sh`, `CHANGELOG.md`, `actions/version.md`.

---

## Repo conventions that bit this session

- **`do-work/` is git-excluded here** via `.git/info/exclude`. Commit steps stage nothing under
  it, and metadata commits are skipped — the `commit:` hash is written into the archived REQ but
  never committed. `record-commit-hash.sh` detects the untracked path, says which git-dependent
  guards it skipped, and still runs every content guard. (`do-work/HANDOFF.md:35` still advises
  hand-editing the hash here — that is what REQ-065 asks about.)
- **Before every commit:** bump `actions/version.md`, add a `CHANGELOG.md` entry whose title says
  what shipped (no codenames), and verify the version and title are not already used.
- **Shipped files must never cite this repo's `CLAUDE.md`/`AGENTS.md`** — both are
  `export-ignore`d, so the citation dangles downstream. The suite greps for it.
- **`SKILL.md` has an enforced 2650-word budget.** None of this work may grow it — no new action,
  routing row, or dispatch row.
- **Never push.** Commit locally only.

## Traps worth not rediscovering

- The awk in `record-commit-hash.sh` **buffers the frontmatter**; delete its `END` flush and the
  script becomes the truncation it prevents (an unterminated block would emit nothing).
- Adding a missing trailing newline also rewrites the final line, so git reports `+2/-2` on an
  otherwise-correct one-line edit. Both the byte arithmetic and the numstat allowance carry an
  explicit `trailing_newline_added` term. Do not "fix" a failure here by loosening the bound.
- `git log` applies history simplification by default — the scanner uses `--full-history` so a
  file blanked on a side branch resolves to the right recovery source.
- `grep -q` on a `git show` pipe under `pipefail` reports failure via SIGPIPE. Both read-back
  greps use `grep -c … || true`.
- A recorded hash that does not resolve is **correctly** refused by `record-commit-hash.sh`. A
  restore probe using a fictional hash will fail for that reason — the fixture must use a real
  commit, as the field always does.
