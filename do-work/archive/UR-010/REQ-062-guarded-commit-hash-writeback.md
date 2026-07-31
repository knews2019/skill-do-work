---
id: REQ-062
title: Deterministic, self-verifying commit-hash write-back
status: completed
claimed_at: 2026-07-30T22:01:38Z
completed_at: 2026-07-30T22:23:05Z
commit: 62a4188
route: C
created_at: 2026-07-30T21:57:34Z
user_request: UR-010
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-063, REQ-064]
batch: commit-hash-writeback-hardening
write_set: [tools/checks/record-commit-hash.sh, actions/work-reference.md, actions/work.md, _dev/tests/record-commit-hash-guards.sh, _dev/tests/contract-regressions.sh, CHANGELOG.md, actions/version.md]
---

# Deterministic, self-verifying commit-hash write-back

## What

Replace the free-form prose write-back at `actions/work-reference.md:860-873` with a shipped guard
script, `tools/checks/record-commit-hash.sh <req-file> <hash>`, that makes the one-line `commit:`
edit deterministic and verifies it before anything is staged. A guard that trips leaves the REQ file
byte-identical to how it was found and tells the operator to stop rather than commit.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Approach written in `## Plan` — five-phase fail-closed script modeled on `tools/checks/qualify.sh`, with all rejection paths ahead of any write and an atomic rename behind the arithmetic. Crew rules loaded: `general.md`, `coding-guardrails.md`, `testing.md` (tdd: true).
- [x] **[APPLY]:** Written exactly as planned, inside the declared Scope — 7 files, no others touched. RED first: the probe file was written and confirmed failing (`FAIL: tools/checks/record-commit-hash.sh must exist and be executable`) before the script existed.
- [x] **[UNIFY]:** `git diff --stat` reviewed — 5 modified files (55 insertions, 14 deletions) plus 2 new untracked files, matching the Scope list exactly. Debug-artifact scan over the non-`do-work/` diff: none. `shellcheck` clean on both new scripts after silencing two info-level false positives (SC2329 — `remove_temp_files` is invoked via the EXIT trap; SC2016 — literal backticks in fixture prose, matching `contract-regressions.sh`'s existing file-level disable). Per file: `record-commit-hash.sh` (new, guards ordered so every rejection precedes the write); `record-commit-hash-guards.sh` (new, 13 probes, isolated git config, `rm -rf` cleanup); `work-reference.md` (procedure block replaced, three prose lines updated); `work.md` (:627/:629 summaries mirrored, script basename present for the contract assertion); `contract-regressions.sh` (script registered, probe file invoked); `CHANGELOG.md` + `version.md` (0.151.0).

## Why

Six archived REQ files in a consumer repo (8 KB–26 KB each) were truncated to 0 bytes by this exact
step. The metadata commit that destroyed each one carried a message claiming success, so the loss
was invisible for weeks. The step's mechanics are currently left entirely to the agent.

## Context

`tools/checks/qualify.sh` is the style reference: `set -uo pipefail`, `OK:`/`WARN:`/`FAIL:` stdout
vocabulary, exit 2 for usage errors, `DO_WORK_*` env options, and comments that explain why each
guard exists. `_dev/tests/contract-regressions.sh:191-208` ratchets every `tools/checks/*.sh` for
existence, executability, and a basename reference in `actions/work.md`.

## Detailed Requirements

**Deterministic edit.** awk scoped to the frontmatter block only — line 1 `---` up to the next bare
`---` — with keys anchored at column 0. A `commit:` in body prose or inside a fenced YAML sample must
be structurally unreachable; several archived REQs quote the schema, which is exactly how a
file-wide `sed` corrupts one. The awk buffers the frontmatter, so its `END` block **must** flush the
buffer and exit non-zero on an unterminated block — without that flush the script becomes the
truncation it exists to prevent. Insert point when `commit:` is absent: immediately after
`completed_at:` (the schema pairs them as the terminal-flip stamp), falling back to the last
frontmatter line.

**Atomic swap.** Write to `mktemp` in the REQ's own directory so `mv` is a same-filesystem rename;
seed it with `cp -p` from the original so `mktemp`'s 0600 cannot ride onto the REQ. `mv` only after
every guard passes. `trap` cleans up temps.

**Pre-edit guards.** Exactly 2 arguments (an unquoted path with a space arrives as 3 and would edit
the wrong file); hash matches `^[0-9a-f]{7,40}$` (this is what rejects a pasted literal `<hash>`);
not a symlink; **non-empty** — an already-0-byte REQ is the aftermath, not an input, so exit 2 and
name the recovery command; no CRLF (`$0 == "---"` cannot see a CRLF delimiter, so the edit would
silently no-op); exactly 0 or 1 `commit:` keys (2+ is ambiguous under last-wins YAML); `id:` matches
`^REQ-[0-9]+$` because it goes into the commit message; hash resolves via
`git rev-parse --verify -q "<hash>^{commit}"` with `^` quoted; **size floor**
`worktree_bytes * 2 >= git cat-file -s HEAD:<path>`; numstat baseline recorded.

**Post-edit guards, before `mv`.** `post_bytes` and `post_lines` equal the closed-form expected
values exactly; `diff | grep -c '^[<>]'` equals 2 (replace) or 1 (insert); re-parse confirms the
frontmatter still closes, exactly one `commit:` equal to the hash, and `id:`/`status:` unchanged.
After `mv`: the numstat delta grew by at most 1/1, else restore from backup and FAIL.

**Idempotency.** `commit:` already equals the hash → `NOOP:` on exit 0, telling the caller to make no
metadata commit. Condition the NOOP on `git diff --quiet HEAD -- <req>` where the path is tracked: if
the worktree records the hash but HEAD does not, an earlier commit was rejected (pre-commit hook) and
the edit is *stranded* — report `OK:` so the caller commits it. Compare the **parsed** value, not the
raw line, or a `commit: abc1234  # note` form rewrites identically forever.

**`--verify <req-file> <hash>` mode.** Read-only, run after the metadata commit: committed blob size
equals worktree size, and the committed blob holds exactly one `commit: <hash>` line. This is the
only thing that can catch a content-mutating pre-commit hook (formatter, lint `--fix`, whitespace
stripper) rewriting the file after every pre-commit guard passed. Use `grep -c … || true`, never
`grep -q` on a `git show` pipe — SIGPIPE under `pipefail` is the trap `qualify.sh` was already bitten
by.

**Degradation.** Non-git repo, or git-ignored REQ path: every content guard still runs; only the
git-dependent guards are skipped, each with a stated reason. Check `git ls-files --error-unmatch`
**before** `git check-ignore` so a tracked-but-also-ignored path is still treated as tracked.

**WARN, never FAIL:** hash is not an ancestor of `HEAD` — broken provenance, not data loss, and both
legitimate shapes pass (serially the hash *is* HEAD; in worktree mode `<merge_hash>` is HEAD's
ancestor because the changelog commit sits on top). **INFO** when the hash is a merge commit
(`git rev-parse --verify -q "<hash>^2"` succeeds), pointing consumers at `git show --first-parent -m`.
The script must never run `git show <hash>` itself — an empty combined diff would validate nothing;
say so in the header so nobody "improves" it back in.

**Wiring, same commit.**
- `actions/work-reference.md:862-873` — replace the fenced block with the script invocation. Serial
  mode collapses to **one command**, `… record-commit-hash.sh <req> "$(git rev-parse --short HEAD)"`,
  because resolve-and-consume now happen inside a single block. State that fix explicitly: the
  re-type-the-hash warning exists precisely because shell state does not survive between an agent's
  command blocks. Worktree mode passes the `<merge_hash>` literal. Document each stdout token and
  what the caller does with it. Keep a short fallback for consumers whose tarball predates the
  script: hand-edit, then `git diff --numstat HEAD -- <req>` must read `1 1` (or `1 0`) — anything
  else means stop and recover, not commit.
- `actions/work-reference.md:860` and `:875` — one-line prose updates naming the script.
- `actions/work.md:627` and `:629` — mirror the summaries, including the basename, so the skeleton
  and the reference do not drift and the contract test's reference assertion passes.
- `_dev/tests/contract-regressions.sh` — add the script to `hardened_check_scripts`.

## Constraints

- Do **not** guard on "terminal-*success*" status. `actions/work-reference.md:128-141` stamps
  `commit:` on every terminal flip and `:854` says failed requests get committed too — a
  terminal-success guard would hard-FAIL every legitimate `failed`/`cancelled` REQ. Guard on
  `status:` unchanged and non-empty; WARN when outside
  `{completed, completed-with-issues, failed, cancelled}`.
- Do **not** introduce a `commit: PENDING` placeholder. Zero occurrences exist repo-wide. Tolerate a
  literal `PENDING` as an existing value for consumer-repo compatibility only.
- `git diff --numstat` is a complementary guard, never the primary one — it yields nothing when the
  REQ is untracked or the repo is not git, which is this repo's own case.
- The script edits and guards; it does **not** commit. `git add`/`git commit` stay with the caller.
- Never `--no-verify`, `--no-gpg-sign`, or `--amend` anywhere in the prose this REQ touches.
- Shipped files must not cite this repo's `CLAUDE.md`/`AGENTS.md` — both are `export-ignore`d.
- `SKILL.md` must not grow: no new action, routing row, or dispatch row.

## Dependencies

None. REQ-063 and REQ-064 both build on this REQ's script and its test harness.

## Builder Guidance

Certainty level: **Firm.** The design was pressure-tested before capture and the guard set is
deliberate — a bug in this script is worse than the prose it replaces, so prefer an extra guard over
a clever one. Every guard listed above traces to a concrete failure mode; if you drop one, say which
and why in the Implementation Summary.

## Red-Green Proof

**RED prompt/case:** In a scratch git repo, create an archived REQ with a real ~9 KB body and a
`commit:` line, blank it to 0 bytes in the worktree (the incident's exact state), then run the
write-back. Today the prose offers nothing that inspects the file, so the damage proceeds to a
metadata commit that claims success.
**Why RED now:** `actions/work-reference.md:860-873` prescribes no command and no verification; six
consumer REQ files were destroyed this way.
**GREEN when:** `tools/checks/record-commit-hash.sh <req> <hash>` exits non-zero on that input,
leaves the file untouched, and prints an operator message naming the size mismatch and the recovery
command — while the same script on an intact REQ produces exactly a `1 1` numstat and exit 0.
**Validation:** Inferred during capture from the upstream acceptance criteria.

## Full Context

See `do-work/user-requests/UR-010/input.md` for complete verbatim input.

---

## Triage

**Route: C** - Complex

**Reasoning:** A new shipped guard script carrying ~15 interacting guards whose whole purpose is preventing data loss, plus prose rewrites across two action files that other actions read as contract, plus the repo's first git-fixture test harness. A bug in this script is worse than the prose it replaces, so it earns planning and exploration.

**Planning:** Required

---

## Plan

**Approach.** Build `tools/checks/record-commit-hash.sh` as a single-pass, fail-closed editor modeled on `tools/checks/qualify.sh`. Order of operations is the design: everything that can reject the input runs before any byte is written, the write goes to a temp file in the REQ's own directory, and the atomic `mv` happens only after the post-edit arithmetic proves exactly one line changed.

**Task 1 — the script.** Five phases in one file:
1. *Argument hygiene* — arg count, hash shape, symlink, non-empty, CRLF, frontmatter/`id:` shape, duplicate `commit:` keys. All exit 2 (usage) except the duplicate-key and unterminated-frontmatter cases, which are exit 1 (data-integrity).
2. *Git probes* — availability, tracked-before-ignored ordering, hash resolution, ancestry WARN, merge-commit INFO. Every probe optional-by-detection.
3. *Pre-edit content guards* — record `pre_edit_bytes`/`pre_edit_lines`, the trailing-newline flag, the HEAD size floor, and the numstat baseline.
4. *The edit* — frontmatter-buffering awk with an `END` flush, into a `cp -p`-seeded temp; then closed-form byte/line arithmetic, the changed-line count, and a re-parse of the result; then `mv`; then the post-`mv` numstat delta.
5. *Idempotency + reporting* — `NOOP:` only when the field already matches **and** HEAD agrees; `OK:` otherwise.

**Task 2 — `--verify` mode.** A separate, read-only entry point that re-reads the committed blob. Dispatch on `$1 == "--verify"` before the two-arg path so the usage error stays accurate.

**Task 3 — prose.** Rewrite `actions/work-reference.md:862-873` (the fenced block) plus the one-line updates at `:860` and `:875`; mirror into `actions/work.md:627`/`:629`. Keep a hand-edit fallback for consumers on an older tarball.

**Task 4 — tests.** New `_dev/tests/record-commit-hash-guards.sh` with a `mktemp -d` + `git init` fixture, invoked from `contract-regressions.sh`. Register the script in `hardened_check_scripts`.

**Plan validation:**
1. *Requirement coverage* — all five upstream items map: deterministic edit → Task 1 phase 4; guards → phases 1/3/4; idempotency → phase 5; worktree mode → the hash is an argument, so it is mode-agnostic by construction (validated, not special-cased); mirrored summaries → Task 3.
2. *No orphan tasks* — `--verify` (Task 2) is the one addition beyond the upstream ask. It is earned: a content-mutating pre-commit hook is a sufficient cause of the original incident and is invisible to every pre-commit guard.
3. *Scope sanity* — 4 tasks, at the flagging threshold but cohesive: one script, its mode, its callers, its tests. Splitting would ship a script no prose invokes.
4. *File conflicts* — none; `do-work/working/` holds only this REQ.

*Generated during planning*

---

## Exploration

**`tools/checks/qualify.sh`** is the style contract: `set -uo pipefail` (no `-e`, so checks accumulate), `failure_count` accumulation, `OK:`/`WARN:`/`FAIL:` line prefixes, exit 2 for usage, `DO_WORK_*` env options documented in the header, and long comments that explain the *trap* each guard closes (the SIGPIPE note on `grep -q`, the pathspec-vs-content-grep note). Header comments carry the "why" — that is the house pattern to match.

**`_dev/tests/contract-regressions.sh`** — `hardened_check_scripts` at :191-208 asserts each script is `-x` *and* that its basename appears in `actions/work.md`. The redaction probe at :439-458 is the only test that executes code: `mktemp -d`, guarded on `command -v jq`, `|| true` on the invocation, explicit `rm -rf`. No git fixture exists anywhere in the repo — `grep -rn "git init" tools/ _dev/ hooks/` returns nothing. This REQ introduces the first one.

**No auto-discovery of test files.** The suite is a single script run by hand; a new `_dev/tests/*.sh` file is dead unless `contract-regressions.sh` invokes it.

**Archived REQs quote the schema.** Several files under `do-work/archive/` contain `commit:` inside fenced YAML samples in their bodies — concrete proof that a file-wide `sed` would corrupt real files, and the reason the awk must be frontmatter-scoped.

**Concerns carried into implementation:** `${#var}` counts locale characters, so byte arithmetic needs `LC_ALL=C`; `mktemp` creates 0600, so the temp must be `cp -p`-seeded or the REQ loses its mode through the rename; awk that buffers frontmatter becomes the truncation bug itself unless `END` flushes.

*Generated during exploration*

---

## Scope

**Files I will touch:**
- `tools/checks/record-commit-hash.sh` (new) — the guarded write-back
- `actions/work-reference.md` (modify) — Commit & Metadata-Commit Procedure, lines 860/862-873/875
- `actions/work.md` (modify) — Step 9 one-line summaries at :627 and :629
- `_dev/tests/record-commit-hash-guards.sh` (new) — git-fixture guard probes
- `_dev/tests/contract-regressions.sh` (modify) — register the script; invoke the new probe file
- `CHANGELOG.md` (modify) — release entry
- `actions/version.md` (modify) — version bump

**Files I will NOT touch:** `actions/forensics.md` and `actions/cleanup.md` (REQ-063/REQ-064 own those), `tools/checks/blanked-req-scan.sh` (REQ-063 creates it), `SKILL.md` (word budget — no new action or routing row), `tools/queue-kanban/**` (no board contract changes here).

**Acceptance criteria (restated from REQ):**
- [ ] A blanked (0-byte) REQ passed to the script is refused, the file is untouched, and the message names the recovery command
- [ ] A pre-truncated REQ trips the HEAD size floor and is refused
- [ ] A correct write-back yields exactly `1 1` on `git diff --numstat` (or `1 0` when inserting)
- [ ] Running twice yields `OK:` then `NOOP:`; a stranded edit yields `OK:`, not `NOOP:`
- [ ] The frontmatter-scoped edit cannot touch a `commit:` occurrence in body prose
- [ ] `--verify` catches a commit whose content differs from what was verified
- [ ] `actions/work.md:627`/`:629` and `actions/work-reference.md` tell the same story
- [ ] `bash _dev/tests/contract-regressions.sh` passes clean with the new probes

---

## Implementation Summary

**Files changed:**
- `tools/checks/record-commit-hash.sh` (new)
- `_dev/tests/record-commit-hash-guards.sh` (new)
- `actions/work-reference.md` (modified)
- `actions/work.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `CHANGELOG.md` (modified)
- `actions/version.md` (modified)

**What was done:** Replaced the free-form `commit:` write-back with a shipped guard script that rewrites only the frontmatter `commit:` line and refuses to write unless closed-form byte/line arithmetic proves exactly that one line changed, plus a read-only `--verify` mode that reads the committed blob back. Rewrote the Step 9 procedure in `actions/work-reference.md` to invoke it (collapsing the serial path to a single command block) and mirrored the summaries into `actions/work.md`. Added the repo's first behavioral git fixture — 13 probes covering every guard — and wired it into the contract suite.

---

## Decisions

- **D-01**: The script edits and guards but does **not** stage or commit. Reasoning: keeps the blast radius to one file and keeps git history entirely in the caller's hands. Value: a tripped guard is a stop signal on an untouched file rather than a rollback of history. Risk: the refusal is a hard stop the procedure mandates rather than a physical impossibility — a caller that ignores exit 1 can still commit damage. Mitigated by `--verify` and by the prose stating the stop explicitly. (ESCALATE tier — settled with the user before implementation.)
- **D-02**: Guard on `status:` being **unchanged and non-empty**, not on a terminal-*success* value as the upstream report asked. Reasoning: the schema stamps `commit:` on every terminal flip and failed requests get committed too, so a terminal-success guard would hard-FAIL every legitimate `failed`/`cancelled` REQ. A non-terminal status is a WARN. (DECIDE & STATE — the upstream ask contradicts the schema; no reasonable reading prefers it.)
- **D-03**: Added a `--verify` mode beyond the upstream ask. Reasoning: a content-mutating pre-commit hook is a sufficient cause of the original incident and is invisible to every pre-commit guard — only a post-commit read-back sees it. Value: closes the one failure path the guards structurally cannot. Risk: low; the mode is read-only and additive. (DECIDE & STATE.)
- **D-04**: The numstat allowance carries an explicit trailing-newline term rather than a fixed `+1`. Reasoning: adding a missing final newline also rewrites the last line, so git reports `+2/-2` on an otherwise-correct one-line edit. Found by the RED probe, which failed on a legitimate file before this term existed. (DECIDE & STATE.)
- **D-05**: `id:` is accepted as `REQ-NNN` **or** `UR-NNN`. Reasoning: the detection work in REQ-063 covers `UR-*.md` files too, and rejecting them here would make the shared restore path in REQ-064 unable to reuse this script. (DECIDE & STATE.)

---

## Qualification

Passed — 7 files verified, all Detailed Requirements traced, P-A-U confirmed.

`tools/checks/qualify.sh` exits 0. `tools/checks/scope-drift.sh` reports `OK: Implementation Summary matches the Scope declaration` — no touched-but-undeclared or declared-but-untouched files. Substantiveness: both new files are real logic, not placeholders. Requirements traced one by one — deterministic frontmatter-scoped edit, atomic swap, the pre-edit guard set, the post-edit arithmetic, idempotency plus the stranded-edit exception, `--verify`, graceful degradation, WARN/INFO classification, and all four wiring targets each map to specific lines. Check 6 (data flows) is not applicable — there is no fetch path to stub — and the probes assert the file content actually changes rather than trusting the script's own report.

---

## Testing

**Tests run:** `bash _dev/tests/record-commit-hash-guards.sh`, then `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ All passing (13 guard probes; full contract suite clean)

**Red-green validation:**
- `_dev/tests/record-commit-hash-guards.sh` (whole file): ✗ before implementation (`FAIL: tools/checks/record-commit-hash.sh must exist and be executable`, exit 1) → ✓ after
- `no trailing newline: exits 0`: ✗ mid-implementation → ✓ after D-04. This probe caught a real defect: git reports `+2/-2` on a one-line edit to a file lacking a final newline, and the first numstat allowance rejected that legitimate case.

Traced to the REQ's `## Red-Green Proof`, and confirmed end-to-end in a scratch repo beyond the probe suite:
- **Upstream acceptance #1** (blanked body is caught, nothing committed): a 12,001-byte archived REQ blanked to 0 bytes → exit 2, file left at 0 bytes, **0 files staged**, message naming both recovery commands.
- **Upstream acceptance #3** (normal write-back stays tiny): same fixture intact → exit 0 and `git diff --numstat` reads exactly `1	1`.

**New tests added:**
- `_dev/tests/record-commit-hash-guards.sh` — 13 probes: replace, insert-after-`completed_at:`, body-`commit:`-untouched, idempotent rerun, stranded-edit exception, blanked file, truncated file (HEAD size floor), duplicate `commit:` keys, unterminated frontmatter, CRLF, argument hygiene (placeholder hash / wrong arg count / missing file), unresolvable hash, missing trailing newline, `--verify` clean and tampered, and non-git degradation.

**Existing tests updated (cross-REQ impact):** none — `_dev/tests/contract-regressions.sh` gained a registration and an invocation, no existing assertion changed.

*Verified by work action*

---

## Review

**Acceptance: Pass — 95%**

**Requirements check.** All five upstream items are satisfied in their corrected form. (1) The edit is deterministic and frontmatter-scoped. (2) The guards run before `git add` and abort on any anomaly — with the two documented corrections: the status guard is *unchanged*, not terminal-*success* (D-02), and byte/line arithmetic is primary while numstat is complementary, because numstat is blind when the REQ is untracked. (3) Idempotency holds, with the stranded-edit refinement so a hook-rejected run isn't stranded forever. (4) Worktree-dispatch mode is covered by construction — the hash is an argument, so both paths run identical guards, and a merge hash is detected and reported. (5) `actions/work.md:627`/`:629` mirror the reference.

**Code review.** Guard ordering is the load-bearing property and it is correct: every rejection path precedes the first write, the temp file is `cp -p`-seeded in the target's own directory, and the `mv` is gated on exact arithmetic. The awk `END` flush is the right defense against the script becoming the bug it prevents. `grep -c … || true` is used on both `git show` pipes, avoiding the SIGPIPE-under-pipefail trap. `LC_ALL=C` makes `${#var}` byte arithmetic sound. No `--no-verify`, `--amend`, or `git add -A` anywhere.

**Acceptance testing.** 13 probes plus a separate end-to-end run in a scratch repo. Upstream acceptance #1 (blanked body caught, nothing committed) and #3 (`1	1` numstat on a correct write-back) both demonstrated with real output. Four adversarial cases run beyond the suite: a path containing a space, frontmatter with no `completed_at:` (insert falls back to the last frontmatter line), an uppercase hash (rejected), and `commit: <hash>  # comment` (parsed-value comparison yields NOOP rather than a pointless rewrite).

**Restatement Sweep (MUST).** Swept every restatement of the write-back contract.
- **Important — fixed:** `actions/work.md:686`, the Orchestrator Checklist, still read "write hash to REQ in separate metadata commit" — an agent following only the checklist would have done the free-form edit the script exists to replace. Now names the script. In declared Scope, so fixed rather than reported.
- **Verified still accurate, no change needed:** `actions/work-reference.md:130-138` (schema comment pointing at the Commit Phase write-back), `actions/commit.md` (deconflicts commit pathways; does not describe this edit).
- **Minor — reported, not fixed:** `do-work/HANDOFF.md:35` tells this repo's future sessions to "write `commit:` hashes directly into archived REQs." Now sub-optimal advice — the script handles a git-excluded `do-work/` gracefully and should be used here too. Outside declared Scope; filed as a Discovered Task.

**Findings:** 1 Important (fixed inline as part of the sweep), 1 Minor (Discovered Task). No follow-up REQs required.

*Reviewed in pipeline mode*

---

## Discovered Tasks

- [low] `do-work/HANDOFF.md:35` advises writing `commit:` hashes directly into archived REQs and skipping metadata commits, which predates `tools/checks/record-commit-hash.sh`. The script degrades correctly when `do-work/` is git-excluded, so this repo's own sessions should use it. Local maintainer note, not shipped.

---

## Lessons Learned

**What worked:** Writing the probe suite first paid for itself immediately — the trailing-newline case (D-04) was a real defect in the numstat allowance that no amount of re-reading the script would have surfaced, because the failing input is a *legitimate* file, not a malformed one. Modeling the guard order on `tools/checks/qualify.sh` meant the reviewable question was only "is this order right," not "what shape should this be."

**What didn't:** The first numstat guard used a flat `baseline + 1` allowance. That silently encodes an assumption — one changed line means one insertion and one deletion — which is false whenever the file lacks a trailing newline, because adding it rewrites the final line too. The fix was to derive the allowance from the same `trailing_newline_added` term the byte arithmetic already carried, rather than to loosen the bound. Loosening would have kept the guard passing while weakening exactly the check the incident calls for.

**Worth knowing:** Three traps are load-bearing and easy to reintroduce. (1) The awk buffers the frontmatter, so removing the `END` flush turns this script into the truncation it prevents — a file with an unterminated block would emit nothing. (2) `git diff --numstat` is *silent*, not zero, when the path is untracked, which is this repo's own configuration; a guard built on it alone would pass vacuously here. (3) `grep -q` on a `git show` pipe under `pipefail` reports failure via SIGPIPE — the same trap `qualify.sh` documents — so both read-back greps use `grep -c … || true`.

---

## Orientation

The Step 9 commit-hash write-back is now a shipped guard script rather than prose: `tools/checks/record-commit-hash.sh` owns the `commit:` frontmatter edit for the whole work pipeline, on both the serial and worktree-dispatch paths. `[MAP CHANGED]` — `tools/checks/` gains its first *mutating* member (the other four only read and report), and `_dev/tests/` gains its first behavioral git fixture, which REQ-063 and REQ-064 both extend. No `prime_files` on this REQ and no prime covers `tools/checks/`, so there is nothing to spot-check for staleness.

---
*Source: upstream bug report from the `game-find-the-difference` consumer repo.*

Think carefully before answering.
