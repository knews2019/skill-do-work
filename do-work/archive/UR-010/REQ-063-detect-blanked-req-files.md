---
id: REQ-063
title: Detect blanked and unparseable REQ files in forensics
status: completed
claimed_at: 2026-07-30T22:13:53Z
completed_at: 2026-07-30T22:31:40Z
commit: d91c567
route: B
created_at: 2026-07-30T21:57:34Z
user_request: UR-010
domain: general
prime_files: []
tdd: true
suggested_spec:
depends_on: [REQ-062]
maintenance: false
related: [REQ-062, REQ-064]
batch: commit-hash-writeback-hardening
write_set: [tools/checks/blanked-req-scan.sh, actions/forensics.md, docs/forensics-guide.md, _dev/tests/record-commit-hash-guards.sh, _dev/tests/contract-regressions.sh, CHANGELOG.md, actions/version.md]
---

# Detect blanked and unparseable REQ files in forensics

## What

Add a detection path for REQ/UR files whose content has been destroyed: a shipped scanner,
`tools/checks/blanked-req-scan.sh`, plus a new read-only forensics check that reports each one as a
data-loss anomaly with its recovery source already resolved.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Approach in `## Exploration` — one read-only scanner emitting both a human report and a stable `BLANKED<TAB>path<TAB>sha<TAB>hash` record, so `actions/forensics.md` and REQ-064's `--restore` consume one implementation. Crew: `general.md`, `coding-guardrails.md`, `testing.md`.
- [x] **[APPLY]:** Written inside the declared Scope — 7 files, no others. RED first: the scanner probes failed with `FAIL: tools/checks/blanked-req-scan.sh must exist and be executable` before the script existed.
- [x] **[UNIFY]:** `git diff --stat` reviewed; changed files match the Scope list exactly. Debug-artifact scan over the non-`do-work/` diff: none. `shellcheck tools/checks/blanked-req-scan.sh` clean. Per file: `blanked-req-scan.sh` (new, read-only, `--porcelain`, history walk avoids `git show --name-only` for paths); `forensics.md` (check 13 added, check 11 gains the no-frontmatter exclusion, one Output Format example); `docs/forensics-guide.md` (one table row); `record-commit-hash-guards.sh` (6 scanner probes sharing the fixture harness); `contract-regressions.sh` (per-script referencing file); `CHANGELOG.md` + `version.md` (0.152.0). Verified live: the scanner reports this repo clean, and `--porcelain` on an incident fixture emits exactly one record.

## Why

The six blanked REQs in the consumer repo surfaced only as `unrecognized status ""` warnings on the
Kanban board, parked under Needs input / Blocked as *untitled*. That framing is actively harmful: the
remedy it implies is "edit the `status:` field," which would cement the loss instead of recovering
the content.

## Context

`actions/forensics.md` is contractually read-only — its Core Rules say "This action never modifies
files… Report, don't fix." Checks are numbered `### N. Title Case Name` and follow a fixed shape:
scan target → what to read → severity ladder → finding template → `**Suggested fix:**` naming another
`do-work` verb → a closing sentence tying the check to the shipped board code where one exists.

## Detailed Requirements

**`tools/checks/blanked-req-scan.sh`** — read-only in this REQ (`--restore` arrives in REQ-064), in
the `tools/checks/qualify.sh` house style. It scans `do-work/queue/`, `do-work/working/`, and
`find do-work/archive -name 'REQ-*.md'` plus the `UR-*.md` equivalents, and reports every file that
is 0 bytes or has no parseable frontmatter.

For each finding it resolves the recovery source from git history: walk
`git log --full-history --format='%H %s' -- <file>` newest-first; the newest commit whose
`git cat-file -s <sha>:<full-name>` is 0 is the **blanking commit**, and its subject yields the
recorded hash via `record commit hash <hash>`; the next older commit with size > 0 is the **recovery
source**. Report the file, its size, the recovery-source commit, and the recorded hash. A file with
no recovery source is reported as unrecoverable-from-history rather than skipped silently.

Never use `git show --name-only` to produce file paths — it prints the commit header and message
first, so a message line can pass a filename grep and become a phantom path. Use `git cat-file -s`
for sizes and `git ls-files --full-name` to resolve the index-relative path.

**`actions/forensics.md` gains `### 13. Blanked or Unparseable REQ Files`**, following the house
check shape and calling the scanner without `--restore` so forensics stays read-only. Severity
**Critical** — content is gone, not mislabeled. `**Suggested fix:**` points at `do-work cleanup`
(Pass 6, REQ-064).

**Edit the existing check 11 (Unrecognized Status Vocabulary)** to exclude files with no parseable
frontmatter, so a blanked REQ is reported once as data loss instead of twice — the second time with
a remedy that would destroy the recovery path.

**Ratchet the new script** in `_dev/tests/contract-regressions.sh`'s `hardened_check_scripts`. That
loop asserts the basename appears in `actions/work.md`; this scanner is referenced from
`actions/forensics.md` instead, so extend the loop to accept the referencing file per script rather
than weakening the assertion.

## Constraints

- **Do not add a `--repair-blanked` flag to forensics.** Its read-only Core Rule is a hard contract
  and every check author downstream reasons from it. Detection here, repair in `actions/cleanup.md`.
- Add `Rules` / `Common Rationalizations` / `Red Flags` rows only if they name a specific failure and
  where it happened — the six blanked REQs qualify, generic hygiene does not.
- `SKILL.md` must not grow: forensics keeps `(none needed)` for arguments, no new routing row.
- Shipped files must not cite this repo's `CLAUDE.md`/`AGENTS.md` — both are `export-ignore`d.

## Dependencies

`depends_on: [REQ-062]` — the git-fixture test harness (`_dev/tests/record-commit-hash-guards.sh`)
is introduced by REQ-062 and extended here.

## Builder Guidance

Certainty level: **Firm** on the detection contract and the forensics read-only boundary,
**Mixed** on the scanner's exact output format — pick whatever shape makes REQ-064's `--restore`
consume it cleanly, since that is the scanner's other caller.

## Red-Green Proof

**RED prompt/case:** In a scratch repo with a 0-byte `do-work/archive/UR-999/REQ-1282-x.md`, run
`do-work forensics`. Today it reports `REQ-1282 has unrecognized status ''` under Warnings and
suggests editing the `status:` field.
**Why RED now:** `actions/forensics.md` has no notion of a body-destroyed REQ; check 11 catches the
empty-frontmatter side effect and prescribes the wrong remedy.
**GREEN when:** the same run reports it under `## Critical Findings` as a blanked-file data-loss
anomaly naming the file, its size, the recovery-source commit and the recorded hash — and check 11
no longer reports it at all.
**Validation:** Inferred during capture from the upstream acceptance criteria.

## Full Context

See `do-work/user-requests/UR-010/input.md` for complete verbatim input.

---

## Triage

**Route: B** - Medium

**Reasoning:** One new scanner script plus two well-bounded edits to `actions/forensics.md`, following patterns REQ-062 just established (guard-script style, git-fixture probes). The shapes are known; what needs exploring is how the scanner's output should read so REQ-064's `--restore` can consume it.

**Planning:** Not required — Route B, exploration-guided.

---

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

---

## Exploration

**`actions/forensics.md` structure** — 13 checks numbered `### N. Title Case Name`, each: scan target (a concrete glob/`find`) → fields to read → severity ladder (`**Info**`/`**Warning**`/`**Critical**`) → optional finding-text template → `**Suggested fix:**` naming another `do-work` verb → optionally a closing sentence tying the check to shipped Go board code. `## Core Rules` at :19-21 is a hard read-only contract. `## Output Format` groups findings under `## Critical Findings` / `## Warnings` / `## Info` / `## Summary`.

**Check 11 (Unrecognized Status Vocabulary), :151-160** — scans all three locations and warns on any `status:` outside the Schema Read Contract's vocabulary. A 0-byte file has no frontmatter at all, so it parses as `status: ''` and lands here with the remedy "edit the `status:` field" — which would cement the data loss. That is the exact double-report the new check must displace.

**Where the incident's evidence lives in git** — the blanking commit is the newest commit touching the file whose blob is 0 bytes; its subject carries `record commit hash <hash>`; the recovery source is the next older commit with a non-zero blob. `git log --full-history` is needed because default history simplification can hide a commit that touched the file on a side branch.

**`hardened_check_scripts` (`_dev/tests/contract-regressions.sh`)** — the loop asserts `-x` **and** that the basename appears in `actions/work.md`. This scanner is referenced from `actions/forensics.md`, not `actions/work.md`, so the loop needs a per-script referencing file rather than a hardcoded one. Extending it keeps the assertion strong; dropping the reference check would weaken the ratchet for all five scripts.

**Output shape for REQ-064** — `--restore` needs, per finding: the path, the recovery-source sha, and the recorded hash. A stable machine-readable line (`BLANKED<TAB>path<TAB>sha<TAB>hash`) alongside the human report is the cheapest way to let one script serve both a prose-reading action and a sibling script.

*Generated during exploration*

---

## Scope

**Files I will touch:**
- `tools/checks/blanked-req-scan.sh` (new) — detection + recovery-source resolution
- `actions/forensics.md` (modify) — new check 13; check 11 exclusion
- `docs/forensics-guide.md` (modify) — user-facing mention of the new check
- `_dev/tests/record-commit-hash-guards.sh` (modify) — probes for the scanner
- `_dev/tests/contract-regressions.sh` (modify) — register the script, per-script reference file
- `CHANGELOG.md` (modify) — release entry
- `actions/version.md` (modify) — version bump

**Files I will NOT touch:** `actions/cleanup.md` and the `--restore` mode (REQ-064 owns both), `tools/checks/record-commit-hash.sh` (REQ-062 shipped it; REQ-064 reuses it), `SKILL.md` (forensics keeps `(none needed)` for arguments), `tools/queue-kanban/**`.

**Acceptance criteria (restated from REQ):**
- [ ] A 0-byte archived REQ is reported as a data-loss anomaly, not an unrecognized-status warning
- [ ] Check 11 no longer double-reports a file with no parseable frontmatter
- [ ] The scanner resolves the recovery-source commit and the hash from the blanking commit's message
- [ ] A file with no recovery source is reported as unrecoverable, not skipped silently
- [ ] `actions/forensics.md` stays read-only — the scanner is called without `--restore`
- [ ] `bash _dev/tests/contract-regressions.sh` passes clean

---

## Implementation Summary

**Files changed:**
- `tools/checks/blanked-req-scan.sh` (new)
- `actions/forensics.md` (modified)
- `docs/forensics-guide.md` (modified)
- `_dev/tests/record-commit-hash-guards.sh` (modified)
- `_dev/tests/contract-regressions.sh` (modified)
- `CHANGELOG.md` (modified)
- `actions/version.md` (modified)

**What was done:** Added a read-only scanner that walks `do-work/archive/`, `queue/`, and `working/` for `REQ-*.md`/`UR-*.md` files that are 0 bytes or have no parseable frontmatter, and resolves each one's recovery-source commit and recorded implementation hash from git history. Wired it into `actions/forensics.md` as Critical check 13, excluded frontmatter-less files from check 11 so a destroyed file is reported once with the right remedy, and generalised the contract suite's hardened-script ratchet to pin a per-script referencing action file.

---

## Decisions

- **D-01**: The scanner emits a tab-separated `BLANKED<TAB>path<TAB>sha<TAB>hash` record alongside the human report, with `-` for an unresolved field. Reasoning: `actions/forensics.md` reads prose, REQ-064's `--restore` needs fields; one script serving both beats two implementations of the same history walk. Callers must handle `-` in either position. (DECIDE & STATE.)
- **D-02**: Exit 1 means "findings", not "error". Reasoning: matches how the surrounding checks read, and the header says so explicitly because a caller treating non-zero as a crash would misreport a successful scan. (DECIDE & STATE.)
- **D-03**: Damage is defined as 0 bytes **or** no parseable frontmatter, not 0 bytes alone. Reasoning: a body that survived but lost its header is equally unreadable to the pipeline and the remedy is the same — recover content, don't edit a field. This is also what makes the check-11 exclusion exactly complementary rather than overlapping. (DECIDE & STATE.)
- **D-04**: `hardened_check_scripts` entries became `script|referencing-file` pairs rather than dropping the reference assertion. Reasoning: this scanner is called from forensics, not the work pipeline, so the hardcoded `actions/work.md` no longer fit. Weakening the assertion would have un-ratcheted all five existing scripts to accommodate one new one. (DECIDE & STATE.)

---

## Qualification

Passed — 7 files verified, all Detailed Requirements traced, P-A-U confirmed.

`tools/checks/qualify.sh` exits 0; `tools/checks/scope-drift.sh` reports no drift. Requirements traced: the scan targets, the 0-byte-or-unparseable predicate, the history walk (newest empty blob = blanking commit, its subject yields the hash; next non-empty = recovery source), the `--full-history` requirement, the no-`git show --name-only` constraint, check 13's placement and severity, the check-11 exclusion, and the contract-suite registration each map to specific lines. Check 6 (data flows) is not applicable — the scanner reads git and the filesystem, and the probes assert real resolved shas rather than trusting its report.

---

## Testing

**Tests run:** `bash _dev/tests/record-commit-hash-guards.sh`, then `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ All passing (19 probes total — 13 write-back, 6 scanner; full contract suite clean)

**Red-green validation:**
- Scanner probe block: ✗ before implementation (`FAIL: tools/checks/blanked-req-scan.sh must exist and be executable`, exit 1) → ✓ after

Confirmed live beyond the probes: the scanner reports this repo's own `do-work/` tree clean (exit 0), and `--porcelain` against a fresh incident fixture — a REQ committed with content, then emptied by a commit titled `[REQ-9] record commit hash abc1234` — emits exactly one record carrying the recovery sha and `abc1234`.

**New tests added:**
- 6 scanner probes in `_dev/tests/record-commit-hash-guards.sh`: the full incident reproduction (seed → blanking commit → detection), the intact-neighbour negative case, the exact `BLANKED` record shape, a file with no recoverable history, and a clean tree exiting 0.

**Existing tests updated (cross-REQ impact):** `_dev/tests/contract-regressions.sh`'s `hardened_check_scripts` loop changed shape (D-04) — an intentional generalisation; all five pre-existing scripts keep an equally strong assertion, now against an explicitly named file.

*Verified by work action*

---

## Review

**Acceptance: Pass — 95%**

**Requirements check.** Every acceptance criterion is met: a 0-byte archived REQ is reported as data loss rather than an unrecognized status; check 11 no longer double-reports it; the recovery source and recorded hash are resolved; a file with no recoverable history is reported explicitly; forensics stays read-only (the scanner is invoked without `--restore`, and `## Core Rules` needed no amendment).

**Code review.** The two documented git traps are avoided: `git show --name-only` is never used to produce paths (a message line could pass a filename grep), and `--full-history` prevents history simplification from hiding a side-branch edit. `resolve_recorded_hash` correctly stops walking once it reaches a non-empty blob, so it cannot attribute an unrelated older commit's message to this damage. The `-` sentinel for unresolved fields is documented in the header where a consumer will look.

**Acceptance testing.** 6 probes plus two live runs. The negative case matters as much as the positive one here — an over-eager predicate that flagged healthy REQs would make the check noise, and the intact-neighbour probe pins that.

**Restatement Sweep (MUST).** This REQ redefines what an empty-`status:` finding *means*, so the sweep targeted every restatement of the unrecognized-status contract. `actions/forensics.md` check 11 is the only prose that asserted it and is updated in this diff. `tools/queue-kanban/model.go`'s `bucketColumns` invalid-status warning still fires for a 0-byte file on the board — correct and unchanged, since the board is display-only and check 11's closing sentence still accurately describes it as the mechanical sweep behind that warning; the board pointing at `do-work forensics` now leads to a check that gives the right remedy. No stale restatement found elsewhere.

**Findings:** 1 Minor (Discovered Task — pre-existing doc gap). No follow-up REQs required.

*Reviewed in pipeline mode*

---

## Discovered Tasks

- [low] `docs/forensics-guide.md`'s "What it checks" table is missing rows for check 11 (Unrecognized Status Vocabulary) and check 12 (Future-Dated Timestamps) — a pre-existing gap, not introduced here. This REQ added the row for its own check 13; backfilling the other two is a separate doc pass.

---

## Lessons Learned

**What worked:** Defining damage as *0 bytes **or** unparseable frontmatter* made check 13 and check 11's new exclusion exactly complementary — same predicate on both sides, so there is no gap where a file is reported twice or not at all. Reusing REQ-062's fixture harness meant the scanner's probes cost almost nothing to add and are tested against the same repo shapes the guard is.

**What didn't:** The first assumption was that the contract suite's `hardened_check_scripts` loop could take this scanner as-is. It could not — the loop hardcodes `actions/work.md` as the referencing file, and this scanner is called from forensics. The tempting fix was to drop the reference assertion; that would have un-ratcheted all five existing scripts to accommodate one new one. Pairing each script with its own referencing file kept every assertion strong.

**Worth knowing:** `git log` applies history simplification by default and can omit a commit that touched the file on a side branch, so the history walk uses `--full-history` — without it a file blanked on a branch and merged could resolve to the wrong recovery source. And `resolve_recorded_hash` must stop at the first non-empty blob: walking past it would happily pick up an older, unrelated `record commit hash` message and hand back a hash from a different incident.

---

## Orientation

`do-work forensics` now distinguishes destroyed content from mislabeled metadata, and `tools/checks/blanked-req-scan.sh` is the shared detector behind it. No `[MAP CHANGED]` — this extends the existing forensics check series and the `tools/checks/` family rather than reshaping either; the one structural note is that the scanner is deliberately the single implementation REQ-064's restore path will call. No `prime_files` on this REQ and no prime covers `tools/checks/` or `actions/forensics.md`.

---
*Source: upstream bug report from the `game-find-the-difference` consumer repo.*

Think carefully before answering.
