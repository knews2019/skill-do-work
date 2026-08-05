---
id: REQ-104
title: Label-less checkpoint entries — "locally modified" is not evidence of authorship
status: completed
completed_at: 2026-08-05T11:37:40Z
claimed_at: 2026-08-05T11:19:30Z
created_at: 2026-08-04T21:15:00Z
kb_status: pending
user_request: UR-018
addendum_to: REQ-094
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set: [actions/work-reference.md, _dev/tests/contract-regressions.sh, docs/work-guide.md, decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md]
related: [REQ-094, REQ-095, REQ-096, REQ-103]
batch: parallel-building
---

# Label-Less Checkpoint Entries — "Locally Modified" Is Not Evidence of Authorship

## What

`actions/work-reference.md` → **Crash Recovery (Step 1)**, the label-less bullet, treats a locally
modified `do-work/CHECKPOINT.md` as evidence that *this* checkout wrote the entries in it:

> **Named there with no `writer:` label at all** (an entry written before the label existed) → **own
> only where `do-work/CHECKPOINT.md` is locally modified or otherwise uncommitted in this checkout**,
> which is evidence this checkout wrote it and has not shared it; recover it as an own crash.

REQ-095's two-clone acceptance run demonstrated that the premise fails under the claim-anywhere model
this batch is building. Once a second checkout can claim, **every** concurrent claim forces a
`CHECKPOINT.md` merge conflict — F-06 of REQ-095's `## Testing` shows `CONFLICT (add/add)` even for two
fully disjoint claims, because two single-line appends land at the same position. A checkout that
resolved that conflict is holding a modified checkpoint for a reason that has nothing to do with who
wrote which entry. So the heuristic fires on a *foreign* label-less entry, classifies it as an own
crash, and strips a live claim — the 2026-07-01 incident, reachable again through the label-less door.

Evidence: REQ-095 `## Testing` → *Defect found: the label-less bullet is unsound under claim-anywhere*
(run 6 transcript, `R6-3`/`R6-4`).

## Detailed Requirements

- Fix the label-less bullet so a merge-resolution-dirtied checkpoint cannot be read as authorship.
  Two candidate shapes, and the choice is the point of this REQ:
  - **Narrow the heuristic** — require *modified* **and** no merge in progress (`git rev-parse
    --verify -q MERGE_HEAD` fails; note `^`/quoting is not involved here but the same
    re-derive-don't-carry rule applies). Keeps auto-recovery of genuine pre-0.170.0 own crashes.
  - **Drop the heuristic** — a label-less entry is always report-only, never recovered. Strictly safer
    and shorter; costs auto-recovery for checkpoints written before 0.170.0, which a human can still
    reset by hand via `actions/forensics.md` Check 1.
  Recommend dropping it: the population it serves is checkouts that have not run the pipeline since
  0.170.0, which shrinks to nothing, while the failure it enables is unbounded and silent. Prefer the
  narrowing only if the wider recovery is worth carrying a second condition that readers must keep true.
- Whatever the choice, keep the surrounding pinned phrases intact — `absent checkpoint is ambiguous`,
  `foreign claim`, `Crash Recovery's input`, `claim held by`, and the label format string are all
  asserted by `_dev/tests/contract-regressions.sh`. Reword around them.
- Mirror the decision at every site that restates the label-less rule. `actions/work.md` Step 1 and
  Step 10's session-start note both carry a version of it; grep for the condition rather than trusting
  this list (per the Closed-Enumerations rule).
- Add a suite assertion pinning whichever rule lands, so the next re-grain cannot quietly widen it back.

## Constraints

- No liveness machinery. This is a change to *how an entry is attributed*, not a check on whether
  anything is still running — refresh intervals, staleness checks and liveness probes stay banned by
  name (`actions/work-reference.md` → In-Progress Record, the `never grow into one` paragraph).
- `maintenance: true` — the candidate fix is removing or narrowing a rule in the skill's own operating
  instructions, so `crew-members/maintenance.md`'s delete-before-you-add discipline applies.

## Red-Green Proof

**RED prompt/case:** REQ-095 run 6 — a label-less foreign entry plus a merge-resolved (therefore
modified) `CHECKPOINT.md` classifies as an own crash under the shipped bullet.
**GREEN when:** the same input state classifies as report-only (or as own only with no merge in
progress, if the narrowing is chosen), and a suite assertion pins it.
**Validation:** Evidence-backed from REQ-095's acceptance run, not reasoning.

## Full Context

See `do-work/user-requests/UR-018/input.md`, `assets/approved-plan.md` (Phase 1), and REQ-095's
`## Testing` record.

---
*Source: REQ-095 acceptance-run finding F-07 (critical discovered task)*

---

## Triage

**Route: B** - Medium

**Reasoning:** The core edit is well-specified (drop or narrow one bullet in `actions/work-reference.md`'s Crash Recovery), but the REQ itself mandates discovery: every site restating the label-less rule must be found by grep (the Closed-Enumerations rule forbids trusting the listed two), and the suite assertion must follow `_dev/tests/contract-regressions.sh`'s existing idiom.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Canonical rule:** `actions/work-reference.md:248` — the label-less bullet in Crash Recovery (Step 1)'s four-case ladder (L244–249). L253 ("the label-less-report-only cases above") is the only in-block forward reference; it stays correct under either fix shape.

**Restating sites (grep-swept, per Closed-Enumerations):**
- `actions/work.md:125` (Step 1) — already narrower: "runs only on a REQ the checkpoint records under this checkout's own `writer:` label". Consistent with DROP as-is.
- `actions/work.md:655` (Step 10 session-start) — says "unlabeled, or labeled for another checkout" is a foreign claim, unconditionally. **Already states the DROP behavior — currently contradicts work-reference L248.**
- `docs/work-guide.md:122` — user-facing; describes own-label-recovers/everything-else-reported. Already reads as DROP. (REQ-101's archived note flagged this paragraph for a post-REQ-104 re-read.)
- `decisions/records/adr-018-*.md:51` (restates four cases) and `:79` (Consequences — documents the bullet as known-unsound, cites the REQ's now-stale `do-work/queue/` path).
- Preservation-rule sites (`actions/work.md:646`, work-reference Session Checkpoint Template L843–850) mention label-less entries but govern *carrying entries through verbatim*, not classification — untouched by this REQ, and two of them are suite-pinned.
- `actions/forensics.md:39` (Check 1) and `actions/cleanup.md:40` (Pass 0) are positively scoped to own-label entries — label-less is out of scope by construction; no edit needed.
- `tools/queue-kanban/verify.go` — no Go code reads `writer:` or implements the classification; parser-transparent. No lock-step change.

**Suite:** `_dev/tests/contract-regressions.sh` extracts `crash_recovery_block` via `sed -n '/^## Crash Recovery (Step 1)/,/^## Worktree Dispatch Mode/p'` (L456) and pins ~10 phrases in it (`foreign claim`, `absent checkpoint is ambiguous`, `claim held by`, etc.). **No existing assertion pins the label-less bullet** — REQ-104 adds one. Idiom: `assert_block_contains "$block" 'pinned phrase' 'message ending in (REQ-NNN).'`; there is also `assert_block_not_contains` for negative pins. Section headers `## Crash Recovery (Step 1)` and `## Worktree Dispatch Mode` are sed-range boundaries — must not be renamed.

**Evidence base:** REQ-095 `## Testing` R6-3/R6-4: a label-less foreign entry plus a merge-resolved (uncommitted) checkpoint classifies as OWN CRASH under the shipped bullet → live claim stripped. F-06 shows the checkpoint conflicts on *every* concurrent claim, so "locally modified" is routine, not evidence. Note: the NARROW option (`MERGE_HEAD` probe) closes the reproduced case but not a committed-merge-then-edited checkpoint.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `actions/work-reference.md` (modify) — replace the label-less bullet (L248) with report-only semantics; verify L253 stays coherent
- `_dev/tests/contract-regressions.sh` (modify) — add assertions pinning the new rule (positive pin on report-only wording; negative pin on the retired "locally modified ⇒ own" inference)
- `docs/work-guide.md` (modify if needed) — REQ-101 flagged one crash-recovery sentence as going wrong when REQ-104 lands; verify and align
- `decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md` (modify, Consequences only) — mark the known-unsound edge resolved by this REQ; fix the stale `do-work/queue/REQ-104` citation

**Files I will NOT touch:** `actions/work.md` (both restating sites already state the drop semantics — verify, don't edit), `actions/forensics.md` / `actions/cleanup.md` (positively scoped to own-label entries; label-less out of scope by construction), `tools/queue-kanban/` (no classification logic exists there), preservation-rule sites (`actions/work.md:646`, Session Checkpoint Template) — different rule, suite-pinned, correct as-is.

**Acceptance criteria (restated from REQ):**
- [ ] RED case goes GREEN: a label-less foreign entry + merge-resolution-dirtied checkpoint classifies as report-only, never an own crash
- [ ] All suite-pinned phrases survive the reword (`absent checkpoint is ambiguous`, `foreign claim`, `Crash Recovery's input`, `claim held by`, label format string)
- [ ] Every site restating the label-less rule mirrors the decision (grep-verified sweep, not the REQ's two-item list)
- [ ] A suite assertion pins whichever rule lands so a re-grain cannot quietly widen it back
- [ ] No liveness machinery (no staleness checks, refresh intervals, or probes)

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified) — Crash Recovery's label-less bullet rewritten: a label-less checkpoint entry is now a **claim of unknown origin, always report-only**; no local state of `do-work/CHECKPOINT.md` upgrades it to an own crash. Points reclaiming at the takeover ladder or `actions/forensics.md` Check 1's manual reset.
- `_dev/tests/contract-regressions.sh` (modified) — new pin pair in the crash-recovery cluster: positive pin on `claim of unknown origin, always report-only`, negative pin (`assert_block_not_contains`) on the retired `locally modified or otherwise uncommitted` fingerprint, with a house-style comment citing REQ-095 F-06/F-07.
- `decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md` (modified) — since-revised aside at the four-case restatement; Consequences marked the known-unsound edge resolved by REQ-104; stale `do-work/queue/` citation de-staled.

**What was done:** Dropped the label-less authorship heuristic (D-01: DROP over NARROW). `docs/work-guide.md` and `actions/work.md` L125/L655 were verified already consistent with the drop semantics and left unedited; the grep sweep found exactly one shipped site asserting the retired inference (the bullet itself). Contract-regression suite exits 0 with the new pins in place.

## Qualification

Passed — 3 files verified in the diff (bullet rewrite substantive, +19-line pin pair in the suite's crash-recovery cluster, ADR-018 since-revised aside + Consequences resolution), 4 requirements traced (bullet fix, suite assertion, mirror sweep grep-verified with zero edits needed, pinned phrases preserved), P-A-U confirmed against the diff — no debug artifacts. Declared-but-untouched `docs/work-guide.md` was conditional in Scope ("modify if needed") and the verify-only outcome is recorded. Orchestrator re-ran the suite independently: exit 0.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh` (run by builder post-edit and independently by orchestrator)
**Result:** ✓ All passing — `Contract regression checks passed.`, exit 0 (chained probes `record-commit-hash-guards.sh` and `update-script-behavior.sh` also green)

**Red-green validation:** (traced to `## Red-Green Proof`)
- Negative pin `locally modified or otherwise uncommitted` in crash_recovery_block: present pre-edit (assertion would FAIL — RED) → absent post-edit (✓ GREEN)
- Positive pin `claim of unknown origin, always report-only`: absent pre-edit (assertion would FAIL — RED) → present post-edit (✓ GREEN)
- Baseline before edits was green (pre-flight `baseline.json`: exit 0), so both failures are the new assertions' RED, not pre-existing breakage

**New tests added:**
- `_dev/tests/contract-regressions.sh`: positive + negative pin pair for the label-less report-only rule (REQ-104)

*Verified by work action*

## Review

**Overall: 90%** | 2026-08-05T11:36:39Z

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 85% |
| Test Adequacy | 90% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 2 important, 3 minor, 2 nit

- **Important 1:** `actions/work-reference.md:456` (In-Progress Record) restates the classification as a closed two-item set ("unnamed, or named under another checkout's label"); under the drop the non-own set is three cases (label-less included). Behavior-correct today, but the Closed-Enumerations failure shape. → follow-up REQ-108
- **Important 2:** the new bullet routes reclaim to the takeover ladder or `actions/forensics.md` Check 1, but neither path's checkpoint-entry removal rule reaches a *label-less* entry (both are scoped to own-label); a reclaimed label-less REQ leaves a phantom entry that `actions/work.md` Step 10's delete gate then blocks on forever. → follow-up REQ-108
- **Minor:** `decisions/log.md:106` still records the edge as open (contradicts ADR-018's amended Consequences); ADR-018 `updated:` not bumped for the 2026-08-05 edits; trailing clause "the same path a foreign claim takes" collides with the sibling bullet's "foreign claim" (which never enters the ladder). First two fold into REQ-108; the clause is report-only.
- **Nit:** ADR L79 names REQ-104 twice in one sentence; the negative pin is exact-phrase brittle (house idiom, no action).

**Acceptance:** Pass — contract-regression suite exits 0; positive pin present and retired fingerprint absent in the extracted crash_recovery_block; RED state independently reproduced from HEAD (both new assertions would have failed pre-edit); both sed-range headers unmoved and all seven other pinned phrases still present.
**Scope drift:** `docs/work-guide.md` declared-but-untouched — judged acceptable (Scope entry was conditional "modify if needed"; the guide already states drop semantics, an edit would have been a no-op).
**Suggested testing:** walk the label-less reclaim path by hand once (reproduces Important 2); re-run the REQ-095 two-clone scenario with a label-less entry end-to-end; at next re-grain confirm a reworded reintroduction is caught by something.
**Follow-ups created:** REQ-108 — In-Progress Record's case list and the label-less removal rule (covers both Important findings + two Minor bookkeeping items)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Pre-exploration that mapped the suite's sed-range boundaries and every pinned phrase before dispatch — the builder rewrote a heavily-pinned bullet without breaking a single existing assertion. Red-green demonstrated with the suite's own extraction idiom (run the sed + grep against HEAD to show both new pins would have failed pre-edit) made the proof mechanical instead of rhetorical.

**What didn't:** The builder's mirror sweep grepped for the *retired inference* ("locally modified ⇒ own") and found exactly one site — but a sweep for *restatements of the classification itself* would have caught In-Progress Record's two-case enumeration in the very file being edited. Sweeping for the deleted phrase is not the same as sweeping for the rule.

**Worth knowing:** Dropping a classification case can orphan its downstream lifecycle rules — the label-less case lost its authorship heuristic, which silently disconnected it from every own-entry removal rule (checkpoint entries now have no documented exit for that case; REQ-108). When deleting a case from a ladder, walk what *used to happen after* that case classified, not just the classification.

## Orientation

Crash recovery in the work pipeline is now strictly label-gated: a checkpoint claim entry with no `writer:` label is never auto-recovered, whatever state the local `CHECKPOINT.md` is in — closing the last door to the cross-checkout claim-strip that the two-clone acceptance run (REQ-095) reproduced. Lives in the work pipeline's crash-recovery contract (`actions/work-reference.md`), pinned both ways by the contract-regression suite. No system-shape change — one classification case narrowed to report-only.

*kb/ handoff: deferred (unattended run — no KB write without consent); kb_status: pending*

## Decisions

**D-01 — Drop the label-less recovery heuristic entirely (Route DROP, not NARROW).** A label-less
`## In Progress (interrupted)` entry is now *always* report-only: no state of `do-work/CHECKPOINT.md`
can upgrade it to an own crash. Reasoning: (a) the NARROW option (`MERGE_HEAD` probe) closes only the
reproduced case — a checkout that *committed* the conflict resolution and then edited the checkpoint
still reads as "modified with no merge in progress", so the strip stays reachable and the rule now
carries a second condition readers must keep true; (b) `crew-members/maintenance.md` § 1 (delete before
you add) and § 2's "is a tool too broad?" both point at removal, and the drift's cause here is the
inference itself, not a stale source around it; (c) the exploration found `actions/work.md:655` and
`docs/work-guide.md:122` already state DROP semantics unconditionally, so DROP *removes* a live
contradiction rather than adding a third variant; (d) the population DROP costs is checkouts that have
not run the pipeline since 0.170.0 — shrinking to nothing — while the failure it prevents (stripping a
live foreign claim, the 2026-07-01 incident) is unbounded and silent; the human path (`actions/forensics.md`
Check 1's manual reset, or the takeover ladder) still recovers a genuinely-own pre-0.170.0 entry.

Tier: **DECIDE & STATE** per `crew-members/coding-guardrails.md` — reversible (one bullet), evidence-backed
(REQ-095 run 6, F-06/F-07), and explicitly recommended by the REQ's own Detailed Requirements. Not
ESCALATE: the REQ author already made the call and no user taste or irreversible cost is involved.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]** Rewrite `actions/work-reference.md`'s label-less bullet (L248) so the case classifies
  as a claim of unknown origin that is always report-only, keeping the sibling bullets' bold
  classification-case opening, the "Never guess-strip" rationale, and a pointer to the human path
  (`actions/forensics.md` Check 1 / the takeover ladder). No suite-pinned phrase moves and neither
  sed-range section header is renamed; L253's "label-less-report-only cases above" stays literally true
  because *all* label-less cases are now report-only. Pin the new rule in
  `_dev/tests/contract-regressions.sh`'s crash-recovery cluster with a positive pin on
  `claim of unknown origin, always report-only` and an `assert_block_not_contains` negative pin on
  `locally modified or otherwise uncommitted` (the exact phrase that authorized the strip). Mark the
  edge resolved in ADR-018 (Consequences, plus a bracketed since-revised aside on the four-case
  restatement) and de-stale its `do-work/queue/REQ-104-…` citation to a bare id. `docs/work-guide.md`
  L122 already describes DROP semantics — verify, expect no edit. No liveness machinery of any kind.
- [x] **[APPLY]** Edits confined to the declared `## Scope` files plus this REQ file.
- [x] **[UNIFY]** `git diff --stat` clean of stray paths; no debug artifacts in added lines;
  `bash _dev/tests/contract-regressions.sh` exits 0. Files checked: `actions/work-reference.md`
  (edited — bullet rewritten; all eight suite-pinned phrases in the block re-verified present, both
  sed-range headers unmoved at L240/L285, L253's "label-less-report-only cases above" still literally
  true), `_dev/tests/contract-regressions.sh` (edited — one positive + one negative pin added to the
  crash-recovery cluster), `decisions/records/adr-018-…md` (edited — L51 bracketed since-revised aside,
  L79 Consequences marked resolved and the stale `do-work/queue/` citation replaced with the bare id),
  `docs/work-guide.md` (read, **no edit** — L122 already states own-label-recovers / everything-else-
  reported, which is exactly the DROP semantics), `actions/work.md` L125/L646/L655 (read, no edit —
  L655's unconditional "unlabeled … is a foreign claim" no longer contradicts work-reference).
