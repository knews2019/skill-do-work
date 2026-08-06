---
id: REQ-125
title: Disposition gate + effort_estimate label on automatic follow-up REQs, with board chip
status: completed
created_at: 2026-08-06T15:48:11Z
claimed_at: 2026-08-06T16:21:10Z
completed_at: 2026-08-06T16:36:10Z
commit: ea76c11
route: B
user_request: UR-027
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
related: [REQ-126, REQ-127]
batch: follow-up-runaway-fix
write_set: [actions/review-work.md, actions/work.md, actions/work-reference.md, actions/capture-reference.md, docs/review-work-guide.md, CLAUDE.md, tools/queue-kanban/model.go, tools/queue-kanban/generate.go, tools/queue-kanban/web/board.js, tools/queue-kanban/web/board.css, tools/queue-kanban/model_test.go]
---

# Disposition Gate + effort_estimate Label + Board Chip

## What

Replace the unconditional one-REQ-per-Important-finding reflex with a recorded disposition gate, and make every automatically created follow-up REQ carry an `effort_estimate: trivial | normal` frontmatter field that the queue-kanban board renders as a visible chip. Nothing is suppressed at this stage — every Important finding still becomes a REQ; it just arrives wearing an honest price tag.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Follow the `domain` field's exact pattern for the board (present-value-only parse guard + resolveSchemaField + warning row + display-only chip), and land the gate as a Step 10 preamble in review-work.md with the three restatement sites in work.md moving in the same edit. Approach mapped in `## Exploration`; crew loaded: general, coding-guardrails, maintenance.
- [x] **[APPLY]:** All eleven Scope files edited, nothing outside the declared list (scope-drift.sh: OK). Mid-run merge of origin/main (PR #136) forced an ID renumber (REQ-122→125, UR-026→UR-027) — bookkeeping, not scope.
- [x] **[UNIFY]:** `git diff --stat` reviewed per file (11 files, prose + Go + web). `gofmt -l` clean after formatting model.go/generate.go; `go build` clean; `go test ./...` ok; `bash _dev/tests/contract-regressions.sh` at the 7-failure root-runner baseline (counted every FAIL line — 7 probes + 1 sub-suite summary), no new failures, none naming touched files. No debug artifacts in diff (checked for stray prints/TODOs).

## Why (if provided)

UR-489: a one-hour feature (REQ-1305) cascaded into sixteen follow-up REQs over two days, fifteen of them trivial facets of one root cause, and the user had to invest their own time to discover the triviality. The user's stated most-important fix is the label: "that way I can easily decide if I want to stop or not the process."

## Detailed Requirements

**The gate.** Before any automatic follow-up REQ is created (review step, `actions/review-work.md` Step 10 and its `actions/work.md` Step 7 callers; also the Step 8 Discovered Tasks flow when it creates REQs), each Important finding gets a recorded disposition token in its report line:

- `gate: user-visible` — a user or developer would actually notice this issue in real use.
- `gate: rule-change` — fixing it establishes or changes a rule that applies in several places (a genuine maintainability rule, not a one-spot patch).
- `gate: trivial` — neither of the above.

The token is mandatory and auditable: it appears in the finding's line in the `## Review` section (or Discovered Tasks classification), so a skipped gate is visible after the fact. The original rule failed precisely because nothing recorded a checkable decision.

**The gate routes; it never re-scores.** Severity vocabulary (Important/Minor/Nit) and severity judgment are untouched. A finding can be genuinely Important ("the guard is blind to rgb() notation") while its disposition is `trivial` (current state is fine to ship). State this explicitly in the shipped text so agents don't resolve the tension by downgrading severities — that would corrupt the severity axis.

**The field.** `effort_estimate: trivial | normal` in REQ frontmatter:

- Closed two-value enum, deliberately — a triage bit, not an estimation system. Document the vocabulary as pinned in the schema comment.
- Absent or unrecognized reads as `normal` (normalize-and-warn class per the Schema Read Contract) — zero migration for existing REQs.
- Automatic follow-ups MUST set it from the gate token (`gate: trivial` → `trivial`; `user-visible`/`rule-change` → `normal`). Capture MAY set it. "Automatic follow-ups" means every REQ the pipeline creates without the user typing it: review follow-ups (Step 7 / review-work Step 10) AND Discovered-Tasks follow-ups (Step 8 substep 4) — both flows stamp the field at creation.
- In this REQ, every Important finding still becomes a `status: pending` REQ exactly as today — suppression/rerouting/consolidation are REQ-126 and REQ-127.

**The chip.** `tools/queue-kanban` renders `effort_estimate: trivial` as a visible chip on the card (and a drawer row), so trivial mechanical fixes are distinguishable from real work at a glance. Display only — no column logic, no scheduling.

**Lock-step obligations, all in the same commit:**

- `tools/queue-kanban/model.go` parses the field → add `effort_estimate` to the board-parsed-fields enumeration in this repo's CLAUDE.md (the enumeration is load-bearing: it's what attaches the mirroring obligation).
- Schema documentation: `actions/work-reference.md` (Full Frontmatter + Schema Read Contract, including the field's normalize-and-warn entry) and `actions/capture-reference.md` (template comment, "capture MAY set it").
- Per Closed Enumerations Go Stale (CLAUDE.md): grep every enumeration of the normalize-and-warn field set (e.g. the list in `actions/capture-reference.md` § Schema Aliases) and update each.

**Text sites to update** (line numbers as of capture — re-grep, don't trust them):

- `actions/review-work.md` ~:335 (Step 10 creation template: add `effort_estimate` + require the gate token), ~:466 (Common Rationalizations row "This finding is minor, not worth a follow-up REQ" — rewrite around the gate, do NOT delete: the failure it guards, silently dropping real findings, is still real), ~:493 (Verification Checklist item "Each Important finding has a follow-up REQ drafted" — update to require a recorded gate disposition and an `effort_estimate` on each created REQ), ~:450 (the existing anti-loop warning — cite the gate as its mechanism).
- `actions/work.md` ~:495, ~:501, ~:505 (all three restatements of one-REQ-per-Important gain the gate + label language together).

**Restatement sweep before done:** grep `"Important finding"` across `actions/`, `crew-members/`, `SKILL.md` — all restatements of the follow-up-creation contract move in this commit.

## Constraints

- See-something-say-something preserved: this REQ changes labeling only, never whether a finding is recorded.
- Severity vocabulary untouched.
- Inline-fix-at-review-resolution is out of scope (deferred by user decision — see UR-027 decision record).
- Board changes are display-only; the chip must not influence bucketing or scheduling.
- Chip goes in the shipped tool source (`tools/queue-kanban/`); never commit the built binary.

## Dependencies

None — this is the root of the batch. REQ-126 and REQ-127 depend on it (both consume the gate token and the label).

## Builder Guidance

Certainty: Firm — the design was discussed and confirmed with the user in detail. The gate token names (`user-visible` / `rule-change` / `trivial`) may be adjusted for clarity if needed, but the three-way semantics and the audit requirement are fixed.

## Red-Green Proof

**RED prompt/case:** Grep `effort_estimate` in `actions/review-work.md` and `tools/queue-kanban/model.go` — no hits. A review-created follow-up REQ today carries no effort marker and no recorded gate disposition, so a UR-489-style cascade produces REQs indistinguishable from real work on the board.
**Why RED now:** `actions/work.md:505` mandates one `status: pending` REQ per Important finding with no relevance check and nothing auditable.
**GREEN when:** A review-created follow-up REQ file contains `effort_estimate` set from a `gate:` token recorded in the reviewed REQ's `## Review` section, and `do-work board` renders a visible chip for `effort_estimate: trivial`. `queue-kanban` parses the field and the CLAUDE.md board-field enumeration lists it.
**Validation:** User confirmed (design discussion preceding capture, 2026-08-06)

## Full Context

See `do-work/user-requests/UR-027/input.md` for complete verbatim input and the decision record.

---

## Triage

**Route: B** - Medium

**Reasoning:** Files and changes are precisely enumerated, but the board chip needs pattern discovery in `tools/queue-kanban/` (how existing badges parse and render across model.go and the embedded web frontend) before code is written. The prose edits span five files that restate one contract — exploration confirms every restatement site.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Board tool (`tools/queue-kanban/`) — the new field follows the `domain` pattern exactly, because `domain` is the existing normalized-enum-with-default field:

- **Parse:** `model.go:638-645` — present-value-only guard (`resolveSchemaField` only when non-empty; absence stays empty so the frontend truthiness gate doesn't badge every card). Struct fields `~:146` (`EffortEstimate` + `OriginalEffortEstimate` + `EffortEstimateUnrecognized`), contract-table row in `schemaReadContractFields` (`model.go:853-940`, `defaultValue: "normal"`), struct-literal wiring `:663-699`, warning row list `model.go:1047-1054`.
- **Serialize:** `generate.go:95-145` `generatedRequest` projection + `buildGeneratedBoardData` `:259-301`; one edit covers both `generate` and `serve` (`serve.go:315` shares it).
- **Badge:** `web/board.js` `makeBadge` `:432`, conditional badge blocks in `makeRequestCard` (`assigned` at `:532-553` is the closest template); gate on `request.effortEstimate === "trivial"`. Drawer row via `appendMetaRow` `~:1735`. CSS modifier class after `.badge-assigned` (`board.css:832-838`).
- **Tests:** `model_test.go` — normalize table `:755-789`, resolve table `:795-828`, clone `TestUnrecognizedDomainFlagsAndWarns` (`:905`) + `TestRecognizedDomainRaisesNoWarning` (`:1047`) + `TestAssignedToNeverAffectsColumnPlacement` (`:679`, the display-only guarantee).
- **No hand-listed field enumerations elsewhere in Go:** `frontmatter_cli.go` is table-driven off `schemaReadContractFields` (gets `effort_estimate` for free); `verify.go` doesn't enumerate fields. Frontend is embedded via `go:embed`; `go build` is the whole build.

Prose restatement sites for the one-REQ-per-Important contract (re-grepped this run): `actions/work.md:495,501,505`; `actions/review-work.md:335` (Step 10), `:450` (anti-loop note), `:466` (rationalization row), `:493` (checklist); `docs/review-work-guide.md:55` (user-facing gloss — one extra site the REQ didn't list). Schema homes: `actions/work-reference.md` Full Frontmatter (~:130s) + Schema Read Contract table (~:202-214); `actions/capture-reference.md:130` normalize-and-warn field list + Simple REQ template comment; CLAUDE.md board-parsed-fields enumeration (~:161).

*Generated by Explore agent (board tool) + orchestrator grep (prose sites)*

## Scope

**Files I will touch:**
- `actions/review-work.md` (modify) — gate definition + tokens in Step 10, findings format, rationalization row, checklist, anti-loop note
- `actions/work.md` (modify) — Step 7 restatements (:495/:501/:505), Step 8 substep 4 effort_estimate stamping
- `actions/work-reference.md` (modify) — Full Frontmatter line, Schema Read Contract row, Discovered Tasks Classification stamping
- `actions/capture-reference.md` (modify) — template comment + normalize-and-warn field list
- `docs/review-work-guide.md` (modify) — user-facing gloss of follow-up creation
- `CLAUDE.md` (modify) — board-parsed-fields enumeration
- `tools/queue-kanban/model.go` (modify) — struct fields, contract row, parse, wiring, warning list
- `tools/queue-kanban/generate.go` (modify) — JSON projection
- `tools/queue-kanban/web/board.js` (modify) — trivial badge + drawer row
- `tools/queue-kanban/web/board.css` (modify) — `.badge-effort-estimate`
- `tools/queue-kanban/model_test.go` (modify) — normalize/resolve/warn/no-bucketing tests

**Files I will NOT touch:** `tools/queue-kanban/verify.go`, `frontmatter_cli.go`, `serve.go` (table-driven, no field enumeration); `actions/capture.md` (Step 1 assessments unchanged — capture MAY set the field, documented in capture-reference); `SKILL.md` (no routing change).

**Acceptance criteria (restated from REQ):**
- [ ] Every automatic follow-up creation site requires a recorded `gate:` token (user-visible / rule-change / trivial) per Important finding
- [ ] Gate routes, never re-scores — stated explicitly in shipped text
- [ ] `effort_estimate: trivial | normal` documented: schema, read contract (absent/unrecognized → normal), capture-reference; automatic follow-ups MUST stamp it (review + Discovered Tasks flows), capture MAY
- [ ] Board renders `effort_estimate: trivial` as a display-only chip + drawer row; unrecognized value warns like domain; no column logic
- [ ] CLAUDE.md enumeration + model.go change in this same commit
- [ ] All restatement sites of one-REQ-per-Important move together (grep verified)

## Implementation Summary

**Files changed:**
- `actions/review-work.md` (modified) — disposition gate defined in Step 10 (two questions, three recorded `gate:` tokens, mandatory + auditable, routes-never-re-scores paragraph); `effort_estimate` added to the follow-up template frontmatter; Findings template's Important line carries the token; anti-loop note cites the gate as its mechanism; Common Rationalizations row rewritten around the gate; Verification Checklist item now requires token + stamp
- `actions/work.md` (modified) — Step 7 restatements at :495/:501/:505 gained gate + `effort_estimate` language together; Step 8 substep 4 stamps `effort_estimate` on Discovered Tasks follow-ups
- `actions/work-reference.md` (modified) — `effort_estimate` schema line in Full Frontmatter (closed enum, absent/unrecognized → normal, lock-step note); Schema Read Contract table row; contract intro generalized off the hand-counted "Nine fields"; Discovered Tasks Classification gained the stamping paragraph (severity ≠ effort)
- `actions/capture-reference.md` (modified) — optional `effort_estimate` template comment ("capture MAY, never invent trivial"); field added to the normalize-and-warn enum list
- `docs/review-work-guide.md` (modified) — user-facing Follow-ups gloss now describes the gate and the chip
- `CLAUDE.md` (modified) — `effort_estimate` added to the load-bearing board-parsed-fields enumeration
- `tools/queue-kanban/model.go` (modified) — `EffortEstimate`/`OriginalEffortEstimate`/`EffortEstimateUnrecognized` struct fields; `effort_estimate` contract-table row (trivial|normal, default normal); domain-pattern parse block (present-value-only guard + resolveSchemaField); struct-literal wiring; warning-row entry in collectSchemaFieldWarnings
- `tools/queue-kanban/generate.go` (modified) — `EffortEstimate` JSON projection (`effortEstimate,omitempty`) + buildGeneratedBoardData wiring (covers generate and serve)
- `tools/queue-kanban/web/board.js` (modified) — `badge-effort-estimate` chip rendered only when `trivial` (tooltip states display-only), drawer "Effort estimate" row
- `tools/queue-kanban/web/board.css` (modified) — `.badge-effort-estimate` neutral/muted class with explanatory comment
- `tools/queue-kanban/model_test.go` (modified) — normalize + resolve table rows; `TestUnrecognizedEffortEstimateFlagsAndWarns`, `TestRecognizedEffortEstimateRaisesNoWarning`, `TestEffortEstimateNeverAffectsColumnPlacement`; stale "seven fields" comment generalized

**What was done:** Replaced the unconditional one-REQ-per-Important reflex with a recorded three-token disposition gate at every automatic follow-up creation site, and gave every such follow-up an `effort_estimate: trivial|normal` frontmatter field that the queue-kanban board renders as a display-only chip (trivial only) with domain-style normalize-and-warn handling — nothing is suppressed, every finding still becomes a REQ wearing an honest price tag.

## Qualification

Passed — 11 files verified against the diff, 6 acceptance criteria traced, P-A-U confirmed. Judgment checks: changes substantive (gate prose + parse logic + tests, no placeholders); `effort_estimate` flows parse → JSON projection → chip/drawer with no hardcoded stub; mechanical script OK (`tools/checks/qualify.sh`), scope-drift OK.

## Testing

**Tests run:** `go test ./...` (tools/queue-kanban); `gofmt -l .`; `go build`; `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ Go suite passing; gofmt clean; build clean; contract suite at its pre-existing 7-failure root-runner baseline (no new FAIL lines, none naming touched files — every FAIL line counted per the 0.176-session lesson)

**Red-green validation:**
- Captured RED (grep `effort_estimate` in `actions/review-work.md` / `tools/queue-kanban/model.go` → no hits) is now GREEN: Step 10 defines the gate + stamp, model.go parses the field, and `do-work board` renders the trivial chip
- `TestUnrecognizedEffortEstimateFlagsAndWarns` / `TestRecognizedEffortEstimateRaisesNoWarning` / `TestEffortEstimateNeverAffectsColumnPlacement`: ✗ before implementation (fields did not exist) → ✓ after

**New tests added:**
- `TestUnrecognizedEffortEstimateFlagsAndWarns` (default + verbatim original + board warning)
- `TestRecognizedEffortEstimateRaisesNoWarning` (case-folded canonical + absent stays empty, both silent)
- `TestEffortEstimateNeverAffectsColumnPlacement` (display-only guarantee)
- normalize/resolve table rows for `effort_estimate` in the two contract-table tests

*Verified by work action*

## Review

**Overall: 95%** | 2026-08-06T17:05:00Z

**Approve** — the gate, the stamp, and the chip land at every site the REQ named, with the lock-step obligations honored in one commit.
Route B | uncommitted (hash written back at Step 9)

**Findings:**

**Important:** None.

**Minor:**
- The two other automatic follow-up creators — Step 8's Failure Classification (Intent/Spec/Code follow-ups) and substep 3's builder-decided questions — don't explicitly stamp `effort_estimate`. Absent reads as `normal` per the contract, which is the honest default for both flows (failure recovery and UX questions are real work), so behavior is correct; stating it would be consistency polish. Report-only.
- No JS test harness exists, so the chip/drawer rendering has no automated coverage (pre-existing condition for every badge; Go tests pin the data layer).

**Requirements Checklist:**
- [x] Recorded `gate:` token required per Important finding at every named creation site — delivered (review-work Step 10 + work.md :495/:501/:505)
- [x] Gate routes, never re-scores — stated explicitly with the downgrade warning
- [x] `effort_estimate` schema + read contract + capture-reference — delivered (closed enum, absent/unrecognized → normal)
- [x] Automatic follow-ups MUST stamp (review + Discovered Tasks flows), capture MAY — delivered
- [x] Board chip display-only with domain-style warning — delivered, pinned by three new tests
- [x] CLAUDE.md enumeration + model.go in the same commit — staged together at Step 9
- [x] Restatement sweep — re-grepped; all sites moved (docs/review-work-guide.md included; actions/code-review.md's "Important findings" is a different action's severity rubric, not this contract)

**Acceptance Testing — Result: Pass**
- `go test ./...` green including the three new end-to-end tests (unrecognized warns + defaults, recognized/absent silent, no column effect)
- `queue-kanban summary` runs against the real tree (125 tickets) with the new parser
- Scope-drift comparison: OK (script)

**Scores:** Requirements 95% | Code Quality 95% | Test Adequacy 90% | Scope 100% | Acceptance Pass

*Reviewed in pipeline mode by work action (Step 7)*

## Lessons Learned

**What worked:** Cloning the `domain` field's complete pattern (present-value parse guard, contract-table row, warning row, the three test shapes) made the board half of a schema addition almost mechanical — the Explore agent's file:line map was accurate at every site. Committing the queue-file fixes from Codex's PR review *before* starting implementation kept the pre-flight tree clean.

**What didn't:** The captured REQ prescribed `grep -l` against a bare directory (exits 2, returns nothing) — caught only by an external review; prescribed commands in REQs deserve the same must-emit-what-consumers-read scrutiny as action files. REQ/UR numbers raced with a concurrent session (its REQ-122/UR-026 merged mid-flight), forcing a renumber to REQ-125/UR-027 — capture's "numbers are not reserved" caveat is real, not theoretical.

**Worth knowing:** `gofmt` realigns the whole struct literal when one long field name lands — run `gofmt -w` before reviewing the diff or the alignment noise buries the real change. The contract-regressions suite's root-runner baseline is 7 probe FAILs plus a sub-suite summary line — count every FAIL line and compare against that, per the 0.176-session miscount.

## Orientation

Reviews now price their follow-ups: every Important finding carries a recorded `gate:` disposition token, and the follow-up REQ it spawns arrives with an `effort_estimate` chip on the board — so trivial mechanical fixes are tellable from real work at a glance. Lives in the review/follow-up machinery (`actions/review-work.md` Step 10, `actions/work.md` Steps 7–8) and the queue-kanban board. **[MAP CHANGED]** — new schema field (`effort_estimate`) + a new gate in the follow-up-creation contract; REQ-126/127 build the reroute and sweep layers on top of it.

---
*Source: "do-work capture-request Ship priorities 1 through 3" — priority 1 of the agreed design*
