# REQ-255 builder brief

**Route B.** Estimated 35 active minutes (P50, medium confidence). **This is the priority fix of the whole follow-up queue: instance 1 is the one path where the shipped SessionStart hook makes a file worse than it found it.**

Context you inherit (your worktree's `do-work/` is stale):
- `repair-req-timestamps.sh` is now a **sourceable library** (REQ-247, 0.212.14): `audit-archive-timestamps.sh` sources it with `timestamp_repair_git_only=1` and report-only switches. Your fix in the shared extractor/`comparison_key_for` lands in BOTH tools — do not edit the auditor unless a switch contract must change.
- The REQ's acceptance criteria (updated since capture) require each widened/refused shape pinned through **both scan scopes** (queue/working via the repairer; archive via the auditor) — the behavior suite already has fixture-repo helpers for both (see the `audit-archive-timestamps:` cases).
- The suite is at 55 named cases. ShellCheck 0.11.0 warning-level runs over every tracked .sh in the gate.
- REQ-247's review noted (Minor, for your awareness): the library switches are `${{var:-default}}` env-overridable; a symlinked archive REQ file is silently skipped; mixed `--fix` semantics unpinned. None of those is your scope — but don't make them worse.
- Two of your six instances came from PR #145's external review with reproduction transcripts in the REQ body; the orchestrator reproduced both before recording. The calendar-impossible case (`9999-99-99T99:99:99Z` erased to mtime) and the duplicate-key case (first-vs-last occurrence divergence) are exactly as described.

Design guidance, not prescription: the root cause is one primitive — `extract_timestamp_fields`' first-whitespace tokenization feeding `comparison_key_for`'s shape patterns. Fix the class at that primitive (e.g. extract the full value with trailing-comment awareness, validate calendar components, honor last-occurrence-wins) rather than patching six symptoms. Where you refuse a shape, refuse it whole and document beside D-04's entry.

## How this build runs

You are a **worktree builder** for the do-work pipeline. Work only inside `/home/user/skill-do-work-worktrees/worktree-agent-REQ-255-give-the-timestamp-repairer-shape-parity-with-the-read-side-detectors` on branch `worktree-agent-REQ-255-give-the-timestamp-repairer-shape-parity-with-the-read-side-detectors`, cut from integration tip `7bb73c3` (0.212.17).

- Never write anything under `/home/user/skill-do-work` except your hand-back file (below). Never read/write `do-work/` in your worktree (stale snapshot; your REQ body is inlined below).
- Commit on your branch in small increments. Never touch `VERSION`, `CHANGELOG.md`, `skills/do-work/actions/version.md` (serial-only).
- One-line needs outside your write set are integration seams (hand back the exact line); larger needs stop-and-report. Out-of-scope finds go to `## Discovered Tasks`.
- Crew rules from your own worktree first: `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`. This REQ is `tdd: true` — also read `crew-members/testing.md`; RED first, per shape. Read every `prime_files` path.
- P-A-U phasing mandatory — work the block in your REQ body, record evidence in the hand-back. D-XX decisions with reasoning.

## Environment notes

- `bash _dev/tests/maintainer-verify.sh` exits 0 at your branch point — baseline and gate; exit code only, never piped.
- Toolchain: Go 1.26.1, ShellCheck 0.11.0, `just`, Node 22, Chromium. Never bare `go build` in the board tool dir.
- Clock: `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment of use. Fixtures in scratch space only.

## Hand-back

Write to exactly this absolute path (the one main-tree write; never stage or commit it):

```
/home/user/skill-do-work/do-work/runs/work-2026-08-18-200845/REQ-255-handback.md
```

Structure: `# REQ-NNN hand-back` — **Branch**, **Commits**, `## What I built`, `## File manifest`, `## P-A-U evidence`, `## Testing evidence` (real RED/GREEN output; observed maintainer-verify exit code), `## Decisions (D-XX)`, `## Integration seams`, `## Discovered Tasks`, `## Pushback`.

---

# Your REQ (verbatim copy — the live one lives in the main tree)

---
id: REQ-255
title: Give the timestamp repairer shape parity with the read-side detectors
status: claimed
created_at: 2026-08-18T17:47:52Z
claimed_at: 2026-08-18T20:08:45Z
route: B
user_request: UR-056
addendum_to: REQ-246
domain: general
review_generated: true
sweep: true
sweep_key: repairer-detector-shape-parity
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-18T20:08:45Z
  basis:
    - Route B
    - 2-file write set
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Give the Timestamp Repairer Shape Parity With the Read-Side Detectors

## What

`repair-req-timestamps.sh`'s hand-rolled shape recognition is strictly narrower than the read-side detectors it claims parity with (the board's `parseTimestamp` / `splitFrontmatter`), and in one shape it actively corrupts: an unquoted space-separated instant is half-rewritten into an unparseable value. Every instance below is a shape the board detects and the repairer either mangles or silently skips. For each: either repair it, or refuse it and document the refusal next to D-04's numeric-offset entry — never half-rewrite. Lock-in cases either way.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Both findings come from REQ-246's independent review, reproduced by execution against the shipped script (review I1 + I2, findings and fixtures quoted in REQ-246's archived Review section). The shared root cause is the D-01 argument — "the board flags it, so it must be repairable here" — applied to the shapes D-01 missed. **I1 first:** until the mangle is fixed, it is the one path where running the shipped SessionStart hook makes a file worse than it found it.

## Instances

- [ ] **Unquoted space-separated instant is mangled (I1, corrupting).** `created_at: 2093-01-01 00:00:00` → `created_at: 2026-08-10T12:00:00Z 00:00:00`: the extractor tokenizes at the first whitespace, the date fragment matches the date-only pattern and is "repaired", the time-of-day survives as a phantom suffix — unparseable YAML from a board-detectable input (`model.go` includes layout `2006-01-02 15:04:05`). The header's "a space may stand in for the T" fold in `comparison_key_for` is unreachable dead logic. Fix or refuse whole; never half-rewrite. The audit line must report the full old value.
- [ ] **Quoted space-separated instant is silently skipped** — truncates to an unmatched-quote token and passes through with no output.
- [ ] **CRLF-fenced files are invisible (I2).** The awk fence match `$0 == "---"` fails on `---\r`; the board's `frontmatter.go` handles CRLF explicitly. Windows agents — the likeliest source of both CRLF files and wrong local-time stamps — get detection without repair, the exact gap REQ-246 exists to close.
- [ ] **BOM-prefixed files are invisible (I2).** Same fence mismatch; the board strips the BOM.
- [ ] **Shape-valid but calendar-impossible values are erased instead of left for diagnosis (external review, reproduced).** `created_at: 9999-99-99T99:99:99Z` matches the glob shape, compares lexicographically later than the horizon, and gets rewritten to the file mtime — contrary to the script's guarantee that unparseable values stay untouched, and it destroys the malformed evidence the board would have diagnosed. Range-check calendar components before using a string as a comparison key. (PR #145 Codex finding, orchestrator-reproduced 2026-08-18.)
- [ ] **Duplicate anchor keys read the first occurrence while the board reads the last (external review, reproduced).** With `claimed_at: <later>` followed by a duplicate `claimed_at: <earlier>`, the board's effective ordering is reversed (defect) while the repairer checks the first occurrence and reports the file clean. Use the last occurrence to match the read side, or refuse duplicate anchors without claiming clean. (PR #145 Codex finding, orchestrator-reproduced 2026-08-18.)
- [ ] **Report-only riders (fold in only if touching the same lines anyway):** dead `frontmatter_value_for` helper; a lock-in grep tying `future_stamp_skew_seconds=120` to the board constant.

## Requirements

- No board-detectable shape is ever half-rewritten: each instance is either repaired to the canonical form or refused byte-identical, and every refusal is documented in the script header next to the D-04 entry.
- Lock-in cases in `_dev/tests/prescribed-shell-scripts-behavior.sh` for each instance, in whichever direction it resolves — the unquoted space-separated mangle case is mandatory.
- Each widened or refused shape gets a fixture through BOTH scan scopes (queue/working via the repairer, archive via `audit-archive-timestamps.sh`), so the shared-fix-reaches-both-tools property is pinned rather than assumed (REQ-247 review).
- `bash _dev/tests/maintainer-verify.sh` exits 0.

---

## Triage

**Route: B** - Medium

**Reasoning:** Six instances of one root cause in a shared shell library with a corrupting case first; the fix direction per shape (repair vs documented refusal) is builder judgment inside two files, with both-scan-scopes pinning required.

**Planning:** Not required

## Plan

**Planning not required** - Route B

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify) — shape-parity fixes in the shared library (reaches both tools by sourcing)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — lock-ins per shape, through both scan scopes

**Files I will NOT touch:**
- `skills/do-work/scripts/audit-archive-timestamps.sh` — it sources the library; parity arrives without editing it (edit only if a switch's contract must change, and then say so).
- `skills/do-work/hooks/session-start.sh` — wiring unchanged.

**Acceptance criteria (restated from REQ):**
- [ ] No board-detectable shape is ever half-rewritten: each instance repaired to canonical form or refused byte-identical with the refusal documented next to D-04's entry.
- [ ] Lock-in cases per instance; the unquoted space-separated mangle case is mandatory; each widened/refused shape pinned through BOTH scan scopes.
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (0.212.17 tip)
**Dependencies:** ✓ toolchain present

*Checked by work action*
