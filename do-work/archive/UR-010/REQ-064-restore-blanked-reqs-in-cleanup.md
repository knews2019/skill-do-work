---
id: REQ-064
title: Restore blanked archived REQs from git history in cleanup
status: completed
claimed_at: 2026-07-30T22:33:10Z
completed_at: 2026-07-31T09:28:23Z
commit: 069c943
kb_status: pending
route: B
created_at: 2026-07-30T21:57:34Z
user_request: UR-010
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-062, REQ-063]
maintenance: false
related: [REQ-062, REQ-063]
batch: commit-hash-writeback-hardening
write_set: [tools/checks/blanked-req-scan.sh, actions/cleanup.md, docs/cleanup-guide.md, _dev/tests/record-commit-hash-guards.sh, _dev/tests/contract-regressions.sh, CHANGELOG.md, actions/version.md]
---

# Restore blanked archived REQs from git history in cleanup

## What

Give `tools/checks/blanked-req-scan.sh` a `--restore` mode that recovers each blanked REQ's content
from git history and re-applies its recorded `commit:` hash, and expose it as a consent-gated
`### Pass 6` in `actions/cleanup.md`. This automates the recovery that had to be done by hand for six
files in the consumer repo.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `prime_files` is empty; loaded `crew-members/general.md`, `coding-guardrails.md`, and
  `testing.md` (`tdd: true`). Approach: `--restore` consumes what the REQ-063 detector already
  resolves (recovery sha + recorded hash) rather than re-walking history, writes via temp-file +
  atomic rename, and delegates the `commit:` edit to `record-commit-hash.sh`. Pass 6 is modeled on
  Pass 5's consent shape and delegates the algorithm to the scanner. The six→seven pass-count
  restatements are treated as one closed-enumeration sweep across `actions/cleanup.md` and
  `docs/cleanup-guide.md`.
- [x] **[APPLY]:** 7 files, all declared (`_dev/tests/contract-regressions.sh` added mid-build and
  declared before the write — D-01). RED first: the restore probes failed against the pre-`--restore`
  script before the flag existed.
- [x] **[UNIFY]:** `git diff --stat` = 7 files, +239/−18. Verified each:
  `tools/checks/blanked-req-scan.sh` (arg parser, `restore_one_file`, exit-code branch — shellcheck
  clean at info+); `_dev/tests/record-commit-hash-guards.sh` (7 restore probes, suite green);
  `_dev/tests/contract-regressions.sh` (one pin line); `actions/cleanup.md` (Pass 6 + 9 restatement
  sites); `docs/cleanup-guide.md` (pass count, new section, key rules); `CHANGELOG.md` /
  `actions/version.md` (0.153.0). No debug artifacts, no `TODO`, no stray `echo` — the script's
  `echo`s are its reporting surface. `shellcheck -S warning tools/checks/*.sh` clean.

## Why

Detection without repair leaves the operator to reconstruct the git archaeology by hand, per file,
under time pressure — and the window is finite: the content survives only until a `git gc` collects
the unreferenced objects.

## Context

`actions/cleanup.md` already writes to the archive and already has the precedent for a consent-gated
pass: Pass 5 removes orphaned worktrees mechanically only when already merged, and otherwise only
with the user's explicit consent. Its `## When to Use` states "cleanup only reorganizes work items…
never deletes them" with two narrow exceptions listed inline.

## Detailed Requirements

**`tools/checks/blanked-req-scan.sh --restore`.** For each blanked file: write
`git show <recovery_sha>:<full-name>` to a temp file in the target's own directory and atomically
rename it into place, then apply the recorded hash by calling `tools/checks/record-commit-hash.sh` —
**reuse it, do not re-implement the frontmatter edit**. `--dry-run` prints the plan without writing.

Refuse to restore over a **non**-empty file unless the size floor says it is truncated; never touch a
file with no resolvable recovery source. Report each restore with the file, the recovery-source
commit, the byte count restored, and the hash written.

**`actions/cleanup.md` gains `### Pass 6: Restore Blanked Archived REQs (consent-gated)`**, modeled
on Pass 5's consent shape. Also update:
- `## When to Use` — the "never deletes" carve-out list gains a third narrow exception, stated as
  *restores* content rather than reorganizing it;
- `## Steps` — the "Six passes, in order" preamble becomes seven;
- `## Commit (Git repos only)` — stage restored paths;
- `## Reporting`, `## Archive Structure After Cleanup`, and `## What This Action Does NOT Do` as each
  applies.

**Test** in the fixture repo introduced by REQ-062: commit a full REQ, commit a blanking commit whose
message is `[REQ-1282] record commit hash 9617040b`, run `--restore`, and assert the restored file is
**byte-identical** to the pre-blanking blob except for the one `commit:` line, which equals
`9617040b`. Also assert `--dry-run` writes nothing, and that a non-empty file is left alone.

## Constraints

- The restore must go through `tools/checks/record-commit-hash.sh` for the `commit:` field so there
  is exactly one implementation of the frontmatter edit and it carries the same guards.
- Never `git show --name-only` for file paths; and note that `git show` on a merge commit prints a
  usually-empty combined diff, so history walking must use `git cat-file -s <sha>:<path>` for sizes.
- Consent is required before any write. `--dry-run` must be the safe default path the action shows
  the user first.
- Add `Common Rationalizations` rows only if they name a specific failure and where it happened.
- `SKILL.md` must not grow: cleanup keeps `(none needed)` for arguments.
- Shipped files must not cite this repo's `CLAUDE.md`/`AGENTS.md` — both are `export-ignore`d.

## Dependencies

`depends_on: [REQ-062, REQ-063]` — reuses REQ-062's write-back script and test harness, and extends
REQ-063's scanner with the `--restore` mode.

## Builder Guidance

Certainty level: **Firm** on reusing the write-back script and on the consent gate; **Mixed** on Pass
6's exact prose. Keep the pass short — it delegates to the scanner, so it should read like Pass 5
(what it finds, what it asks, what it does), not like a re-specification of the recovery algorithm.

## Red-Green Proof

**RED prompt/case:** In a scratch repo, commit a ~9 KB archived REQ, then commit a blanking commit
titled `[REQ-1282] record commit hash 9617040b` that truncates it to 0 bytes. Run `do-work cleanup`.
Today the file is left blanked — cleanup has no notion of destroyed content.
**Why RED now:** there is no repair path anywhere in the skill; the recovery was performed by hand.
**GREEN when:** Pass 6 reports the blanked file, asks for consent, and on approval restores content
byte-identical to the pre-blanking blob with `commit: 9617040b` in the frontmatter — and `--dry-run`
produces the same report while writing nothing.
**Validation:** Inferred during capture from the upstream acceptance criteria.

## Full Context

See `do-work/user-requests/UR-010/input.md` for complete verbatim input.

---
*Source: upstream bug report from the `game-find-the-difference` consumer repo.*

Think carefully before answering.

## Triage

**Route: B** (Explore then Build)

The "what" is fully specified by the Detailed Requirements — a `--restore` mode on an existing
script plus a new consent-gated pass in an existing action file. The "where" needed discovery:
Pass 5's consent shape is the model for Pass 6, and cleanup.md's six-pass framing is restated in
five places (`## When to Use`, `## Steps` preamble, `## Reporting`, Common Rationalizations, the
Verification Checklist) plus `docs/cleanup-guide.md`. No planning phase — the design was settled
at capture.

*Re-derived on session resume 2026-07-31; the original triage was lost when the prior session
ended mid-REQ without writing this section. `route: B` in frontmatter was already set by that
session and is unchanged.*

## Exploration

Key findings that shaped the build:

- `tools/checks/blanked-req-scan.sh` (REQ-063) already resolves, per damaged file, both the
  recovery-source commit and the hash recorded in the blanking commit's message, and emits them
  as a `BLANKED<TAB>path<TAB>sha<TAB>hash` record. `--restore` therefore needs no new archaeology
  — it consumes what the detector already computed.
- `tools/checks/record-commit-hash.sh` (REQ-062) is the sole guarded implementation of the
  `commit:` frontmatter edit. The restore calls it rather than re-implementing the edit.
- `actions/cleanup.md` Pass 5 is the consent precedent: interactivity test → mechanical path when
  safe → ask before anything destructive → non-interactive runs report only. Pass 6 mirrors that
  shape but inverts the risk (restoring content is safe; the gate exists because it writes into
  the archive at all).
- The "six passes" count is restated in six places across two files. All of them move together.

*Re-derived on session resume; the original exploration output was lost with the prior session.*

## Scope

**Files I will touch:**

- `tools/checks/blanked-req-scan.sh` — add `--restore` / `--dry-run`
- `actions/cleanup.md` — Pass 6 plus the six→seven restatements
- `docs/cleanup-guide.md` — user-facing paragraph for the new pass
- `_dev/tests/record-commit-hash-guards.sh` — restore probes
- `_dev/tests/contract-regressions.sh` — extend the scanner's `hardened_check_scripts` pin to the
  second referencing action file (added mid-build, see D-01)
- `CHANGELOG.md` — release entry
- `actions/version.md` — version bump

**Acceptance criteria (restated from the REQ):**

1. `--restore` recovers each blanked file's content from its recovery commit and re-applies the
   recorded `commit:` hash **by calling `tools/checks/record-commit-hash.sh`**, never by editing
   frontmatter directly.
2. `--dry-run` prints the plan and writes nothing.
3. A healthy (non-blanked) file is never touched.
4. `actions/cleanup.md` gains `### Pass 6: Restore Blanked Archived REQs (consent-gated)`, and
   every restatement of the pass count and the never-deletes carve-out is updated with it.
5. A fixture test asserts the restored file is byte-identical to the pre-blanking blob except the
   one `commit:` line.

## Implementation Summary

**Files changed:**

- `tools/checks/blanked-req-scan.sh` (modified) — `--restore` / `--dry-run` modes: a `while`-loop
  arg parser replacing the single-argument `case`, mutual-exclusion guards, an up-front refusal
  when `record-commit-hash.sh` is absent, and `restore_one_file()` doing the per-file repair.
- `_dev/tests/record-commit-hash-guards.sh` (modified) — 7 restore probes in a throwaway git repo:
  dry-run names the file and writes nothing, restore exits 0, restored content is byte-identical to
  the pre-blanking blob except the `commit:` line, the healthy neighbour is untouched, and a re-run
  finds nothing left.
- `_dev/tests/contract-regressions.sh` (modified) — second `hardened_check_scripts` pin asserting
  `actions/cleanup.md` also references the scanner.
- `actions/cleanup.md` (modified) — `### Pass 6: Restore Blanked Archived REQs (consent-gated)`,
  plus the nine sites restating the pass count and the never-deletes carve-out.
- `docs/cleanup-guide.md` (modified) — pass count, a user-facing Pass 6 section, key-rules line.
- `CHANGELOG.md` (modified) — 0.153.0 entry.
- `actions/version.md` (modified) — 0.152.0 → 0.153.0.

**What was done:** `blanked-req-scan.sh` gained a repair mode that reads each damaged file's
content out of the recovery commit it already resolves, writes it back through a temp file in the
target's own directory plus an atomic rename, refuses to write recovered content that is itself
empty, and then re-applies the recorded `commit:` hash by *calling* `record-commit-hash.sh` rather
than editing frontmatter — so the one guarded implementation of that edit stays the only one.
`--dry-run` reports the same plan and writes nothing. `actions/cleanup.md` exposes it as Pass 6,
gated on the user's consent, report-only when non-interactive, and running last so it repairs files
at the paths the earlier passes just moved them to.

## Decisions

**D-01 — Extend `hardened_check_scripts` with a second pin rather than leaving or weakening the
existing one.** DECIDE & STATE. `_dev/tests/contract-regressions.sh` pinned
`blanked-req-scan.sh` to `actions/forensics.md`; Pass 6 makes `actions/cleanup.md` a second
referencing file. The three options were leave it (still valid, but silent if cleanup's pointer is
deleted), generalize the entry to "any action file" (weaker — the whole value of the pair format is
naming *which* file must keep the pointer), or add a second pin. Added the second pin: one line, and
it now fails if either caller's reference is removed. This put a file outside the declared
`write_set`, so the Scope list and `write_set` were extended before the write, per the
write-boundary rule — no other REQ was in flight to contend for it.

**D-02 — No separate "size floor" check before restoring.** DECIDE & STATE. The REQ asked to
"refuse to restore over a non-empty file unless the size floor says it is truncated." That refusal
is satisfied one layer up instead: only files the detector already classified as damaged (0 bytes,
or non-empty with no parseable frontmatter) ever reach `restore_one_file`, so a healthy file is
never a candidate and there is nothing for a floor to arbitrate. Adding a second, independent size
test inside the restore would create a way for the two gates to disagree about what counts as
damaged. The `restore: the healthy neighbour is untouched` probe asserts the property the
requirement was after.

**D-03 — Pass 6 runs last, after the moving passes.** DECIDE & STATE. A blanked REQ that Pass 1 or
Pass 2 consolidates into a UR folder would otherwise be restored at its old path and then moved,
producing a delete-plus-add pair in the cleanup commit instead of one content edit, and giving the
repoint step a moved file whose content changed in the same run. Stated explicitly in the pass so a
future reordering has to argue with it.

## Qualification

Passed — 7 files verified on disk and in `git diff`, all 5 acceptance criteria traced to a specific
change, P-A-U boxes confirmed against the diff.

- **Files exist / show in diff:** all 7, `+239/−18`. No `(new)` files, so no dead-code exception to
  evaluate.
- **Substantive:** the scanner's `+120` is one arg parser, one 50-line function, and one exit-code
  branch — no placeholder bodies. `actions/cleanup.md`'s `+54` is a full pass plus 9 restatements.
- **Requirements traced:** (1) restore-via-`record-commit-hash.sh` → `blanked-req-scan.sh:190`;
  (2) `--dry-run` writes nothing → `:154-158` + probe; (3) healthy file untouched → detector gate +
  probe; (4) Pass 6 and its restatements → `actions/cleanup.md` (verified by grep: no `six passes` /
  `all 6 passes` / `Passes 0–5` survives in a shipped file); (5) byte-identity probe →
  `record-commit-hash-guards.sh:414`.
- **Flowing:** not a data-path change; the one external call (`record-commit-hash.sh`) is executed,
  not stubbed, and its failure is reported rather than swallowed.
- **No debug artifacts** in the diff.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh` (includes the git-fixture guard probes),
`shellcheck tools/checks/*.sh`
**Result:** ✓ All passing — 57 assertions in the guard fixture, contract regressions clean,
`shellcheck` clean at info+ on `blanked-req-scan.sh` and at warning+ across `tools/checks/`.

**Red-green validation** (`tdd: true`, tracing the REQ's `## Red-Green Proof`):

- **RED:** the restore probes were written against the pre-`--restore` scanner. `--restore` was
  rejected by the old single-argument `case` (`usage: … [--porcelain]`, exit 2), the dry-run probe
  found no `WOULD RESTORE` line, and the byte-identity assertion compared the pre-blanking blob
  against a 0-byte file. The captured RED case is reproduced exactly: a ~9 KB archived REQ committed,
  then a blanking commit titled `[REQ-1282] record commit hash <hash>` truncating it to 0 bytes.
- **GREEN:** `--restore` exits 0 and the restored file is byte-identical to the pre-blanking blob
  except the one `commit:` line, which carries the hash from the blanking commit's message;
  `--dry-run` produces the same report with a 0-byte file left on disk; the healthy neighbour is
  unchanged; a re-run reports nothing left to restore.
- **Deviation from the captured proof, with reason:** the fixture uses a *real* commit hash rather
  than the literal `9617040b` from the REQ text. `record-commit-hash.sh` correctly refuses a hash
  the repo cannot resolve, so a fictional one would fail the write-back for the wrong reason and
  the probe would assert nothing about the restore. Noted inline in the fixture.

**Existing tests updated:** none — the pre-existing scan probes still pass unchanged, which is the
regression evidence that `--restore` didn't disturb detection.

**Baseline:** `do-work/working/baseline.json` recorded `bash _dev/tests/contract-regressions.sh`
green at pre-flight, so there were no pre-existing failures to net out.

## Review

**Overall: 92%** | 2026-07-31T09:28:23Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 85% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 0 important, 3 minor

- **Minor — two `restore_one_file` branches are untested.** The no-recovery-source `SKIP:` path and
  the missing-`record-commit-hash.sh` refusal are both correct on inspection but have no probe. The
  REQ's Test paragraph named three assertions and all three exist, so this is beyond what was asked;
  recorded rather than fixed to keep the diff at the REQ's boundary.
- **Minor — `--restore` prints the detector's human block (including "Restore with: `do-work
  cleanup`") before each repair line.** Harmless and arguably useful context during a dry run, but
  the suggestion is circular when Pass 6 is the thing running. Left alone: changing it would touch
  the shared detection output that forensics also renders.
- **Minor — `actions/help.md:82` still labels cleanup "Consolidate the archive."** Considered and
  deliberately left: it's a one-line menu label describing the action's purpose, and Pass 6 is a
  safety net rather than the headline. Flagged so the choice is visible rather than an oversight.

**Restatement sweep (MUST — performed):** this diff redefines the cleanup pass count, the
"never deletes / only reorganizes" carve-out, and `blanked-req-scan.sh`'s CLI contract.

- Pass count: 9 sites in `actions/cleanup.md` (When to Use, Steps preamble, Reporting example, the
  "Archive is clean" gate, Archive Structure note, Commit staging comment + template + staging
  paragraph, Does-NOT-Do, Common Rationalizations, Red Flags, Verification Checklist ×2) plus 3 in
  `docs/cleanup-guide.md` — all updated. Grep confirms no `six passes` / `all 6 passes` /
  `Passes 0–5` survives in any shipped file.
- `actions/forensics.md` check 13 already forward-referenced "`actions/cleanup.md` Pass 6" and
  "Never pass `--restore` here" (written by REQ-063 in anticipation). Those pointers now resolve —
  verified, not changed.
- Historical restatements in `CHANGELOG.md` and archived REQs under `do-work/archive/` were left
  alone by design: they record what was true at their release.
- No file outside `actions/cleanup.md` restates the never-deletes carve-out;
  `actions/stray-check.md` and `docs/dream-guide.md` only delimit cleanup's *territory*, which is
  unchanged.

**Acceptance:** Pass — all five declared acceptance criteria met, with D-02 recording how criterion
3 is satisfied structurally rather than by the mechanism the REQ sketched.
**Suggested testing:** 1 item — run `do-work cleanup` interactively in a repo with a real blanked
REQ and confirm the ask fires before any write, then again dispatched as a subagent and confirm it
reports without writing.
**Follow-ups created:** None — all three findings are Minor.

*Reviewed by review-work action (pipeline mode)*

## Discovered Tasks

- `[low]` `_dev/tests/record-commit-hash-guards.sh` carries two pre-existing shellcheck warnings
  from the REQ-063 scan block — SC2120 (`run_scan_script` takes args that are never passed) and
  SC2034 (`intact_bytes` assigned but unused, at line 321). Both predate this REQ (confirmed against
  `HEAD`) and are test-file hygiene with no behavioral effect. Left in place per the
  surgical-changes rule.

## Lessons Learned

**What worked:** Building `--restore` on top of what the REQ-063 detector already resolves — the
recovery sha and the recorded hash were computed and printed; the repair mode just consumes them.
No second history walk exists to drift from the first. Delegating the `commit:` edit to
`record-commit-hash.sh` rather than re-implementing a one-line frontmatter write was the constraint
that made the whole thing small.

**What didn't:** The first instinct on "refuse to restore over a non-empty file unless the size
floor says it is truncated" was to add a size check inside the restore. That would have been a
second, independently-drifting definition of "damaged" sitting under the detector's — the refusal
already falls out of the detector gate (D-02). Also: the pass-count sweep is bigger than it looks.
"Six passes" appears once, but the *count* is restated nine times in `actions/cleanup.md` under
different phrasings ("all 6 passes", "Passes 0–5", "two narrow exceptions") — grepping the literal
string finds one of them.

**Worth knowing:**
- The restore fixture must use a **real** commit hash. `record-commit-hash.sh` refuses a hash git
  can't resolve, so a made-up one fails the write-back for an unrelated reason and the probe proves
  nothing about the restore.
- `restore_one_file` prints its count to stdout and everything human-readable to stderr, because
  the caller consumes it as `$(...)`. Adding a plain `echo` for the operator inside that function
  corrupts the count.
- Pass 6 deliberately runs last (D-03): the earlier passes move files, and restoring before the
  move would split one content edit into a delete-plus-add in the cleanup commit.

## Orientation

`do-work cleanup` can now repair archived REQs whose content was destroyed by an unguarded
`commit:` write-back — it detects them, shows a dry run, asks, and restores from git history with
the recorded hash re-applied through the guarded write-back. This closes the loop UR-010 opened:
`record-commit-hash.sh` (REQ-062) prevents the damage, `forensics` check 13 (REQ-063) detects it,
and cleanup Pass 6 repairs it. Lives in the cleanup action and the `tools/checks/` family.

**[MAP CHANGED]** — cleanup gains a seventh pass and, for the first time, a path that *writes into*
a work item's content rather than relocating it. That is a change to the action's contract, not just
its step list: "cleanup never modifies file contents" is no longer true, and the two places that
stated it now carve Pass 6 out explicitly.

No `prime_files` on this REQ, and no prime covers `actions/cleanup.md` or `tools/checks/` — no
staleness spot-check applicable, no prime-link writes pending.
