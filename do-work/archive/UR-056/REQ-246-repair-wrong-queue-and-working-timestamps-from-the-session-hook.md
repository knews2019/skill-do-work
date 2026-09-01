---
id: REQ-246
title: Repair detectably wrong queue and working timestamps from the session hook
status: completed
created_at: 2026-08-18T12:38:26Z
claimed_at: 2026-08-18T16:09:27Z
completed_at: 2026-08-18T17:49:50Z
commit: 270a2d0
kb_status: promoted
kb_entry: REQ-246-repair-detectably-wrong-queue-and-workin.md
route: C
user_request: UR-056
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-247, REQ-244, REQ-245]
batch: timestamp-stamping-integrity
estimate:
  p50_active_minutes: 45
  confidence: medium
  calculated_at: 2026-08-18T16:10:30Z
  basis:
    - Route C
    - 4-file write set
    - 2 new files
    - 2 subsystems involved
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- skills/do-work/hooks/session-start.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
- _dev/tests/session-start-hook-behavior.sh
---

# Repair Detectably Wrong Queue and Working Timestamps From the Session Hook

## What

A core script that scans REQ files in `do-work/queue/` and `do-work/working/` for detectably wrong `*_at` stamps and rewrites them with a mechanically derived correct value — no agent judgment anywhere in the path. Wired into the SessionStart hook the way `scripts/cleanup-req-reservations.sh` already is, so repair happens before any agent or board render sees the file.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md`, four crew files, exemplar `record-commit-hash.sh`, the cleanup script + hook + both test suites, and the board's skew-allowance code. Approach: one script, string-comparable canonical stamp keys, two mechanical defect predicates, file-state-decided replacement source, record-commit-hash-style guards, hook block cloned from the cleanup wiring. (Transcribed from builder hand-back.)
- [x] **[APPLY]:** All writes confined to the four declared write-set files; board tool, action markdown and archive scope untouched. (Transcribed from builder hand-back.)
- [x] **[UNIFY]:** `git diff --stat 67dae6b..HEAD` = 4 files, +716/-0, each reviewed; `shellcheck --severity=warning` clean on all four; no debug artifacts; working tree clean; no bare `go build` run. (Transcribed from builder hand-back.)

## Why (if provided)

Detection already exists read-side (board generate/verify, forensics Check 11) but repair depended on the same agent that wrote the bad stamp. This closes that loop mechanically: a fabricated future stamp gets corrected regardless of agent compliance.

## Detailed Requirements

- **Detect only provable wrongness:** stamps later than now plus the 2-minute skew allowance (matching `actions/forensics.md` Check 11 and the board's `futureTimestampSkewAllowance` in `model.go`), and impossible orderings (`completed_at` earlier than `claimed_at`, `claimed_at` earlier than `created_at`).
- **Replacement value, by file state:** file dirty against HEAD → file mtime (the actual write instant); committed content → author time of the commit that introduced the stamp line. Clamp so the repaired set satisfies `created_at ≤ claimed_at ≤ completed_at ≤ now`.
- **Guard style of `tools/checks/record-commit-hash.sh`:** verify before replace, atomic write, a tripped guard leaves the file byte-identical, nonzero exit on any failure.
- **Audit trail:** print each correction (file, field, old value, new value, replacement source) to the hook output. No new frontmatter fields. (Log-only was the capture recommendation; the user did not override it — builder may surface it for confirmation if the hand-back warrants.)
- **Hook wiring:** invoked from the SessionStart hook (`skills/do-work/hooks/`) alongside `cleanup-req-reservations.sh`; also invocable directly.
- **Scope is queue + working only** (user-confirmed). Archived files are REQ-247's territory.

## Constraints

- Not in the board tool: its frontmatter surface is documented read-only (`frontmatter_cli.go` ~line 38) and its write-surface count is pinned in CLAUDE.md § Kanban Board Write Surfaces.
- No LLM, no heuristic text parsing beyond frontmatter field extraction.
- Platform floor applies: POSIX shell, no Go build in the repair path.
- Provenance: replaces the Discuss verdict on validate-feedback Finding 2 (UR-055 triage); Surface-cost — earned by the fabricated-stamp incident report, and cheaper than the alternative (per-site guard prose scattered across write sites, which the Timestamp rule's centralization clause forbids).

## Red-Green Proof

**RED prompt/case:** A `_dev/tests/` lock-in test builds a fixture REQ with `created_at` two hours in the future in a scratch queue, runs the repair script, and asserts the stamp is rewritten to the fixture file's mtime and the correction is logged. Fails today because the script does not exist.
**Why RED now:** Nothing repairs a detected-wrong stamp; the board only warns.
**GREEN when:** The test passes, an ordering-violation fixture is likewise repaired, a clean fixture passes through byte-identical, and `bash _dev/tests/maintainer-verify.sh` exits 0.
**Validation:** User confirmed (scope via ask tool; repair-window and mtime sources are the user's own design)

## Full Context

See `do-work/user-requests/UR-056/input.md` for complete verbatim input.

---
*Source: "basically if a wrong timestamp is detected, it should be automatically mechanically no-llm corrected by the tool" + ask-tool answer "Queue + working (Recommended)"*

---

## Triage

**Route: C** - Complex

**Reasoning:** A new shipped script with its own detection, derivation and clamping logic, plus hook wiring and a lock-in test suite — multiple components, many precise requirements, and a guard-style contract to match.

**Planning:** Required

## Plan

**Approach**

1. **One script, two entry conditions.** `skills/do-work/scripts/repair-req-timestamps.sh [project-root]` scans `do-work/queue/` and `do-work/working/` only. It runs from the SessionStart hook and stands alone as a direct invocation; the hook guards it exactly the way it guards `cleanup-req-reservations.sh`, so a partial install or a repairer failure can never break the status banner.

2. **Detection is arithmetic, never judgment.** Two predicates, both mechanical: a stamp later than `now + 120s` (the same 2-minute allowance `actions/forensics.md` Check 11 and the board's `futureTimestampSkewAllowance` use), and an impossible ordering among `created_at`, `claimed_at`, `completed_at`. Nothing else is a defect — a stamp that is merely surprising is left alone.

3. **Replacement source is decided by file state, not by preference.** Dirty against HEAD → the file's mtime, which is the actual write instant. Committed content → the author time of the commit that introduced that stamp line. Then clamp the repaired set so `created_at ≤ claimed_at ≤ completed_at ≤ now`.

4. **Guard style is `tools/checks/record-commit-hash.sh`'s, and that is a hard requirement, not a stylistic nod.** Verify before replace; write atomically; a tripped guard leaves the file byte-identical; any failure exits nonzero. Free-form frontmatter edits truncated six archived REQs to 0 bytes in a consumer repo — that is the incident this style exists to prevent.

5. **Tests are lock-in, not decoration.** The `## Red-Green Proof` cases go into `_dev/tests/prescribed-shell-scripts-behavior.sh` (future stamp repaired to mtime, ordering violation repaired and clamped, clean fixture byte-identical, tripped guard byte-identical and nonzero) and the wiring probe into `_dev/tests/session-start-hook-behavior.sh`. Each case names the failure it pins.

**Ordering:** script first with its cases RED, then the hook wiring, then the probe.

**Watch for:** the previous session's recurring finding — a mechanism that looks like it closes a class and closes only the instance. Here the instance is *a future `created_at` in the queue*; the class is *every detectably wrong `*_at` in either directory, under both replacement sources*. Test the derivation you did not reach for first.

*Written inline by the orchestrator (no separate Plan agent) — Route C.*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (new) — the mechanical repairer: detection, replacement derivation, ordering clamp, atomic guarded write, audit trail
- `skills/do-work/hooks/session-start.sh` (modify) — invoke the repairer alongside `cleanup-req-reservations.sh`, guarded the same way
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — the lock-in cases named in `## Red-Green Proof`
- `_dev/tests/session-start-hook-behavior.sh` (modify) — the hook-wiring probe

**Files I will NOT touch:**
- `skills/do-work-board/tools/queue-kanban/**` — the board's frontmatter surface is read-only and its write-surface count is pinned (REQ Constraints).
- Any shipped action markdown under `skills/*/actions/` — REQ-249 is sweeping that tree in this same wave. If wiring needs documenting there, hand the exact lines back as an integration seam.
- `do-work/archive/**` — REQ-247's territory.

**Acceptance criteria (restated from REQ):**
- [ ] A fixture REQ with `created_at` two hours in the future is rewritten to the fixture file's mtime, and the correction is logged.
- [ ] An ordering-violation fixture (`completed_at` < `claimed_at`, or `claimed_at` < `created_at`) is likewise repaired, clamped so `created_at ≤ claimed_at ≤ completed_at ≤ now`.
- [ ] A clean fixture passes through byte-identical.
- [ ] A tripped guard leaves the file byte-identical and exits nonzero.
- [ ] The repairer runs from the SessionStart hook and can also be invoked directly.
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Implementation Summary

**What was done:** Built `skills/do-work/scripts/repair-req-timestamps.sh`, a POSIX-floor mechanical repairer for detectably wrong `*_at` stamps in `do-work/queue/` and `do-work/working/`: future stamps beyond the shared 2-minute skew allowance (any top-level `*_at` key — the suffix is the rule, not a field list) and impossible orderings among `created_at`/`claimed_at`/`completed_at`. Replacement source is file mtime when the file is dirty against HEAD (or untracked / no git), otherwise the introducing commit's author time via `git blame --line-porcelain`; derived values are clamped to `created_at ≤ claimed_at ≤ completed_at ≤ now`. Guard style follows `tools/checks/record-commit-hash.sh` (verify-before-replace, atomic same-directory rename, tripped guard leaves the file byte-identical, nonzero exit, one audit line per correction). Wired into `hooks/session-start.sh` alongside the reservation cleanup, presence-guarded so a partial install still prints the banner; also directly invocable.

**Files changed:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (new) — detection, replacement derivation, ordering clamp, guarded atomic write, audit trail (573 lines)
- `skills/do-work/hooks/session-start.sh` (modified) — invokes the repairer after the reservation cleanup, guarded the same way
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified) — five lock-in cases: future→mtime with logged correction, ordering clamp on both later fields, committed-file→blame author time incl. quoted stamps, clean fixture byte-identical, tripped guard byte-identical+nonzero
- `_dev/tests/session-start-hook-behavior.sh` (modified) — hook-wiring probe: repair runs at session start, audit lines join the banner, partial-install cases unchanged

**Deliberately untouched:** numeric-offset / fractional-second values (not provable without timezone arithmetic — D-04), nested keys like `estimate.calculated_at`, unparseable values, symlinks, and everything under `do-work/archive/` (REQ-247's territory).

*Integrated by orchestrator from builder hand-back; merge range `e427aa1..270a2d0`.*

## Qualification

Passed — 4 files verified in the merge range `e427aa1..270a2d0` (573-line script is substantive, not placeholder), all six requirement clusters traced (skew-120s detection, file-state-decided replacement, record-commit-hash guard style, log-only audit trail, hook + direct invocation, queue+working scope), P-A-U audited against the diff (no debug artifacts; the script's success/audit lines are its contract output, not instrumentation). Orchestrator spot-check: the hook's `2>/dev/null` does not lose failure lines — the repairer prints FAILED lines to stdout by design (D-03), only the cannot-read-clock abort uses stderr.

## Decisions

Transcribed by the orchestrator from the builder hand-back (the run artifact is uncommitted by contract, so the durable record lives here — review Minor finding 2).

- **D-01 (DECIDE & STATE): quoted stamps are comparable and repairable.** The schema's YAML readers unquote values, so a quoted future stamp is board-detectable; one matching pair of wrapping quotes is stripped for comparison and the replacement is written canonical-unquoted. Locked in by the committed-file case.
- **D-02 (DECIDE & STATE): the `_at` suffix is the detection rule, not a field list.** Ordering constraints still name the three schema anchors; the future check covers every top-level `*_at`, matching forensics Check 11 and § Closed Enumerations Go Stale.
- **D-03 (DECIDE & STATE): the hook keeps repairer output on failure (`|| true`) instead of discarding.** The repairer's failure lines ARE the audit trail for a tripped guard; the banner still can never break. (Verified by execution in review.)
- **D-04 (DECIDE & STATE): non-comparable shapes are left alone.** Numeric offsets and fractional seconds need timezone arithmetic to judge; the conservative direction for an unattended hook is not to touch what it cannot prove.
- **D-05 (DECIDE & STATE): ordering repairs rewrite the later field of the pair; the earlier anchors.** The earlier stamp has the stronger provenance.
- **D-06 (DECIDE & STATE): epoch-to-stamp conversion is probed, not guessed.** GNU `date -d @EPOCH` vs BSD `date -r EPOCH` verified against epoch 0 before use.

## Review

**Overall: 78%** | 2026-08-18T17:44:44Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 85% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Verdict: Approve with follow-ups** — the mechanism is real and well-guarded (every stated acceptance criterion reproduced by hand, boundary-exact, and the guard architecture held under every mutant thrown at it), but the repair-side parser recognizes a strictly narrower set of shapes than the read-side detectors it claims parity with, and in one reproduced case it rewrites a detectable stamp into a *worse*, unparseable value.

### Requirements Checklist (each walked against the diff and re-run, not taken from the hand-back)

- [x] **Future `created_at` rewritten to file mtime, correction logged** — reproduced directly (`created_at: 2093-01-01T00:00:00Z -> 2026-08-10T12:00:00Z (file mtime)`), plus suite case re-run green.
- [x] **Detection is exactly now+120s, matching board/forensics** — boundary probed by execution: now+119s passed byte-identical; now+121s repaired. Strictly-greater-than-horizon semantics match `futureTimestampSkewAllowance`.
- [x] **`_at`-suffix class, not a field list** — reproduced: a future `status_changed_at` (non-anchor field) was repaired.
- [x] **Ordering violations repaired, clamped `created ≤ claimed ≤ completed ≤ now`** — reproduced: both later fields clamp to the floor; `completed_at < created_at` with `claimed_at` absent repairs against the right anchor; a pass-1 future repair landing earlier than `created_at` is caught by pass 2 and clamped up.
- [x] **Replacement source decided by file state** — reproduced both branches beyond the suite: a dirty *tracked* file in a git repo uses mtime (the suite's mtime cases are git-less, so that path had never actually been exercised); a committed file uses the *introducing* commit's author time even when a later commit edited other lines.
- [x] **Never writes a worse instant from clamping's point of view** — a fabricated future mtime (`touch -t 2090…`) is clamped to now via `earlier_of`.
- [x] **Guard style of `record-commit-hash.sh`** — reproduced: truncation-floor trip leaves the file byte-identical, exits 1, names the reason.
- [x] **Audit trail log-only, no new frontmatter fields** — delivered (caveat under Finding I1).
- [x] **Hook wiring + direct invocation** — reproduced end-to-end: a tripped guard's `FAILED` line survives the hook's `2>/dev/null` and reaches the banner while the hook still exits 0 (D-03 verified by execution); all adversarial runs were direct invocations with explicit project root.
- [x] **Queue + working only; archive untouched; symlinks refused** — reproduced: archive fixture byte-identical; a symlinked REQ file with a future stamp left untouched and still a symlink.
- [x] **`bash _dev/tests/maintainer-verify.sh` exits 0** — re-run un-piped by the reviewer: exit 0.
- [x] **Scope check** — `scope-drift.sh` exit 0; diff touches exactly the four declared write-set files.

### Findings

**Important:**

- **I1 (reproduced by execution): a space-separated instant is mangled, not repaired — the one case found where the script writes a *worse* value.** Fixture `created_at: 2093-01-01 00:00:00` (unquoted) came out as `created_at: 2026-08-10T12:00:00Z 00:00:00` — the extractor tokenizes at the first whitespace, the leading date fragment matches the date-only pattern and gets "repaired", and the time-of-day survives as a phantom suffix, producing an unparseable YAML value. The board detects the original shape (`model.go` `parseTimestamp` includes layout `"2006-01-02 15:04:05"`), so this is detectable-wrong made worse. The script's header claim "a space may stand in for the `T`" is unreachable dead logic in `comparison_key_for`; the quoted form is silently skipped. The audit line also under-reports the old value. File: `skills/do-work/scripts/repair-req-timestamps.sh`. — gate: **user-visible** → sweep REQ-255 (shared root cause below).
- **I2 (reproduced by execution): CRLF- and BOM-fenced REQ files are invisible to the repairer while fully visible to the board.** The awk fence match `$0 == "---"` fails on `---\r` and on a BOM, while the board's `frontmatter.go` handles both. The environment most likely to produce CRLF files and wrong local-time stamps (Windows agents) is exactly where the shipped repair silently no-ops while the board keeps warning — the detection-without-repair gap this REQ exists to close. Undocumented (unlike D-04's numeric offsets). Conservative direction held (no corruption). — gate: **rule-change** → sweep REQ-255.

  **Shared root cause (one sweep, not two REQs):** the repairer's hand-rolled shape recognition is narrower than the read-side detectors it claims parity with — the D-01 argument applied to the shapes D-01 missed. `sweep_key: repairer-detector-shape-parity`. Instances: unquoted space-separated instants (mangled — fix or refuse, never half-rewrite), quoted space-separated instants (skipped), CRLF fences (skipped), BOM fences (skipped).

- **I3 (restatement sweep, argued from reading): README's hook description is stale in a way the diff widens.** `README.md:188` — "SessionStart hook that injects the installed version and pending REQ count" — now omits that the hook *writes repairs into the user's queue/working files* at session start. A consumer auditing "what writes to my repo at session start" is misled. Routed to follow-up REQ-256, not builder scope drift (README was correctly outside the declared write set). — gate: **user-visible**.

**Minor:** 3 (report only)
- `frontmatter_value_for` (script line 177) is defined and never called — dead code.
- The D-01…D-06 decision records live only in the untracked run artifact; the durable REQ cites D-03/D-04 without defining them. (Orchestrator note: the Decisions are now transcribed below.)
- The audit line's old-value truncation (subsumed by I1's fix).

**Nit:** 1 — the 120-second skew constant now has a fourth hand-kept copy; a trivial lock-in grep tying `future_stamp_skew_seconds=120` to the board const would match how the repo pins cause-clause pairs.

### Acceptance Testing

**Result: Partial** — every acceptance criterion passes end-to-end, reproduced against the real script and the real hook (both suites re-run green; maintainer-verify exit 0; twelve hand-built adversarial fixtures under the session scratchpad, none in the repo's `do-work/`). Partial, not Pass, because a detectable edge shape actively misbehaves under the shipped code (I1 writes a corrupted value). Red-green evidence in the hand-back is consistent with re-running; GREEN independently confirmed. Cross-REQ test updates: none (additions only). P-A-U consistent with the diff.

**Not tested (environment limits):** the BSD `date -r` branch (Linux box), Windows/PowerShell behavior, two concurrent sessions racing one REQ file (argued safe from mktemp-per-run + atomic rename, not run).

### Suggested Additional Testing

- Lock-in cases for whichever way the sweep resolves I1/I2 (repair or documented refusal) — the space-separated mangle especially.
- A BSD/macOS run of the behavior suite before the next release.
- A two-session concurrency smoke if consumer reports ever hint at temp-file debris.

### Reviewer-Recommended Disposition

**Approve with follow-ups.** Merge stands. One sweep follow-up (REQ-255, `sweep_key: repairer-detector-shape-parity`) carrying I1+I2, with I1 fixed first — until then, the mangle is the one path where running the shipped hook makes a file worse than it found it. I3 → REQ-256 (doc line). Minor findings report-only per Step 10's threshold.

**Important findings (audit record):**
- I1 space-separated instant mangled into unparseable value — gate: user-visible → sweep REQ-255
- I2 CRLF/BOM-fenced files board-detectable but silently unrepairable, undocumented — gate: rule-change → sweep REQ-255
- I3 README hook description omits that the hook now writes to user files — gate: user-visible → REQ-256

**Acceptance:** Partial — all stated criteria pass reproduced; one edge shape corrupts (I1)
**Follow-ups created:** REQ-255 (sweep), REQ-256 (doc) — created by orchestrator per orchestrated mode

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Cloning the guard architecture of `record-commit-hash.sh` wholesale (verify-before-replace, atomic rename, byte-identical on trip) survived every adversarial mutant the review threw at it — including a fabricated future mtime. Small commits (script+RED, GREEN, hook wiring) made two transport-level interruptions nearly free to resume.

**What didn't:** The repair-side parser was hand-rolled instead of derived from the read-side detectors it claims parity with — so it recognizes strictly fewer shapes than the board (space-separated instants, CRLF/BOM fences), and in one shape (unquoted space-separated) it half-rewrites and corrupts. The session's standing class-vs-instance warning fired anyway, one layer deeper than the builder looked: D-01 closed quoted stamps and the review found the shapes D-01 missed (REQ-255).

**Worth knowing:** The hook's `2>/dev/null` is safe only because the repairer deliberately prints failure lines to stdout (D-03) — anyone adding stderr output to the script will silently lose it in the banner. `comparison_key_for`'s space-fold is dead code until REQ-255 resolves it. The 120s skew constant now has a fourth hand-kept copy.

## Orientation

Now the queue repairs itself: detectably wrong `*_at` stamps in `do-work/queue/` and `do-work/working/` are mechanically corrected at session start (SessionStart hook), with no agent judgment in the path — replacement values derive from file mtime or the introducing commit's author time. Lives in core's `scripts/` + `hooks/` subsystem alongside the reservation cleaner. [MAP CHANGED] — the SessionStart hook is now a *write* surface on consumer queue files, not just a banner (which is exactly what follow-up REQ-256 documents). Prime staleness spot-check: `_dev/primes/prime-shell-commands.md` paths still resolve; not stale.

## Pre-Flight

**Git:** ✓ clean outside `do-work/`
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 (recorded in `do-work/working/baseline.json`)
**Dependencies:** ⚠ this checkout needed Go 1.26.1, ShellCheck 0.11.0 and `just` installed before the baseline could run at all, and one pre-existing Linux-only test failure had to be fixed first (0.212.8) — see the REQ brief.

*Checked by work action*
