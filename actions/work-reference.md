# Work Action — Reference

> Companion file to `work.md`. Holds the heavy templates, tables, and sub-procedures the orchestrator references by name — extracted to keep `actions/work.md` focused on the ten-step skeleton. Each section below is pointed to from the matching step in `actions/work.md`. Loading this file is only necessary when you reach the step that references it — and read only the named section. If this file is already in context from an earlier step this session, reuse it; don't re-read it at every reference site.

---

## Architecture

```
work action (orchestrator - lightweight, stays in loop)
  │
  ├── Read CHECKPOINT.md if exists (resume context from previous session)
  │
  ├── For each pending request (skip pending-answers):
  │     │
  │     ├── TRIAGE: Assess complexity (no agent, just read & categorize)
  │     │
  │     ├── OPEN QUESTIONS? ── - [ ] items exist ──► Mark - [~], builder decides
  │     │                      (none / all resolved) ──► continue
  │     │     │
  │     │     ├── Route A (Simple) ──────────────────┐
  │     │     │   Skip plan/explore, direct to build │
  │     │     │                                      │
  │     │     ├── Route B (Medium) ───────┐          │
  │     │     │   Explore, scope declare  │          │
  │     │     │                           ▼          │
  │     │     └── Route C (Complex) ──► Plan ──► Explore ──► Scope declare
  │     │                                            │
  │     │                                            ▼
  │     │                                     Implementation agent
  │     │                                            │
  │     │                                            ▼
  │     │                                  Implementation Summary
  │     │                                            │
  │     │                                            ▼
  │     │                              Qualify (orchestrator verifies)
  │     │                                            │
  │     │                                            ▼
  │     │                                        Testing
  │     │                                            │
  │     │                                            ▼
  │     │                                  Review ◄─── Fail? ──► Remediate ──► Re-review
  │     │                                            │
  │     │                                            ▼
  │     ├── Archive ──► classify discovered tasks ──► queue follow-ups
  │     │                                            │
  │     │                                            ▼
  │     └── Commit (git repos only)
  │
  └── Context wipe → Loop | Write CHECKPOINT.md → cleanup → report
```

## Folder Structure

```
do-work/
├── queue/                         # Pending REQ files (the work queue)
│   └── REQ-018-pending-task.md
├── user-requests/                 # UR folders (verbatim input + assets)
│   └── UR-003/
│       ├── input.md
│       └── assets/
├── working/                       # Currently being processed
│   └── REQ-020-in-progress.md
└── archive/                       # Completed work
    ├── UR-001/                    # Archived as self-contained unit
    │   ├── input.md
    │   ├── REQ-013-feature.md
    │   └── assets/
    ├── REQ-010-legacy-task.md     # Legacy REQs (no UR) archive directly
    └── legacy/                    # Consolidated legacy items
```

- **`queue/`**: The queue — only pending `REQ-*.md` files
- **`working/`**: Claimed requests. Immutable to all actions except the work pipeline.
- **`archive/`**: Completed UR folders (self-contained) and legacy REQs/CONTEXT docs
- **`user-requests/`**: Active UR folders. Moved to `archive/` when all REQs complete.

This tree exists in the **main working tree only**. Worktree dispatch mode's builder checkouts live outside the repo entirely and never carry their own copy of it — see **Worktree Dispatch Mode (Step 1)**.

## Request File Schema — Full Frontmatter

**Timestamp rule — every `*_at` field in this schema, and any timestamp a future field adds:** write the **current UTC instant** as ISO-8601 `YYYY-MM-DDTHH:MM:SSZ`, obtained with `date -u +%Y-%m-%dT%H:%M:%SZ`. Never stamp local wall-clock time with a `Z` suffix appended — in any zone east of UTC that produces a *future* instant, which silently corrupts every elapsed-time reading (queue wait, claim stopwatch) and gets the REQ flagged by `do-work board` (a "future stamp" card badge plus a data warning, allowing 2 minutes of clock skew). Write sites that say `<timestamp>` or `<now>` mean exactly this rule.

```yaml
---
# Set by capture action
id: REQ-001
title: Short descriptive title
status: pending
domain: frontend  # choose one: frontend, backend, ui-design, general, security, or testing
tdd: false       # optional — set true when test-first applies (per capture's TDD heuristic); drives Step 6 testing-crew loading and RED/GREEN mode
caveman: false   # optional — `true` or intensity `lite` | `full` | `ultra`; loads crew-members/caveman.md to compress agent prose
maintenance: false  # optional — set true by capture for a removal/narrowing finding on the skill's OWN instructions (agent/action/crew/prime file); loads crew-members/maintenance.md (delete-before-you-add) in Step 6 alongside coding-guardrails. Not for ordinary app-source dead-code removal.
prime_files: []  # list paths to relevant prime-*.md files, or leave empty
created_at: 2025-01-26T10:00:00Z
user_request: UR-001          # May be absent on legacy REQs
addendum_to: REQ-NNN          # optional — present only when this REQ amends an in-flight or completed REQ; set by capture, or by review when creating follow-ups. **Legacy alias:** every read site (Step 8 upstream walk, Step 8 cycle detection, Step 8 follow-up generation, and roadmap Blocked classification) also recognizes `amends:`, `parent:`, and `amendment_to:` as synonyms when `addendum_to` is absent so natural-English glosses don't silently drop the parent linkage; `addendum_to` wins when multiple are present. Capture and follow-up REQs always emit `addendum_to:` — never propagate the alias.
depends_on: []                # optional list of REQ IDs that must reach `completed` or `completed-with-issues` before this REQ runs. Semantically distinct from `addendum_to` ("amends that REQ"): depends_on is "requires that REQ to be done first." A REQ can have both. Honored by Step 1's selection scan and by Step 8's upstream-failure classification. **Legacy alias:** every read site (Step 1 selection, Step 1 cycle detection, Step 1 `--wave` depth, Step 8 upstream walk, roadmap classification) also recognizes a `dependencies:` key as a synonym so muscle-memory typos from Python/Node/Cargo conventions don't silently bypass gating; `depends_on` wins when both are present. Capture and follow-up REQs always emit `depends_on:` — never propagate the alias.
write_set: []                 # optional list of repo-relative paths/globs this REQ expects to write. **Absent or empty ⇒ unknown ⇒ treat as overlapping every other REQ** — the safe default that keeps every existing REQ on today's serial behavior. Seeded by capture when the request names files (`actions/capture-reference.md` → Populating `write_set`), then firmed by Step 5.5: the `## Scope` section's "Files I will touch" list is the source and this field is its mirror, never the reverse. Two other writers exist and neither contradicts that: **Step 1's parallel-dispatch gate writes a partition** here at dispatch — the narrowed subset a co-dispatched builder owns, so the later re-checks see the partition instead of the set it was carved from (`actions/work.md` Step 1); the write is narrowing-only, and that subset becomes the REQ's `## Scope` list at Step 5.5, so the mirror confirms it. **Crash recovery clears it**, but only when the REQ has a `## Scope` for the mirror to have come from (**Crash Recovery (Step 1)**, substep 1). Read by Step 1's optional parallel-dispatch gate (`actions/work.md`) and parsed for display by `tools/queue-kanban/model.go`. A **scheduling input, not a safety guarantee** — nothing enforces it at the filesystem, so an overlap means serialize, not "write carefully." Globs use `path.Match` semantics: `*` never crosses `/`, `**` is not recursive, a malformed pattern is no-match for that direction though literal equality short-circuits first (so two REQs declaring the identical malformed pattern still badge each other), and an entry naming a directory never matches a file inside it (`actions/` vs `actions/board.md`) — those board miss-classes are illustrative, not a closed list, and the gate still treats unexpandable globs as overlapping.

# Set by work action when claimed
claimed_at: 2025-01-26T10:30:00Z
route: A | B | C

# Set by reserve action (do-work reserve — allocation to a DIFFERENT worktree/cloud session; see actions/reserve.md)
status: reserved              # holding state; REQ stays in do-work/queue/ — the default scan skips it, only targeted `do-work run REQ-NNN` claims it (clearing these fields)
reserved_for: "cloud-alpha"   # free-text owner label (always YAML-quoted — raw user text)
reserved_at: 2025-01-26T10:15:00Z   # staleness anchor: older than 24h ⇒ readers flag the reservation as stale and suggest recategorizing

# Set by capture (external-condition task) or by the work pipeline's mid-run blocked flip (Step 8's blocked-flip procedure). Holding state — the REQ stays in do-work/queue/ and the default scan walks past it, exactly like pending-answers/reserved.
status: blocked               # waiting on an EXTERNAL condition — not user answers (that's pending-answers), not another REQ (that's depends_on). Cleared to `pending` by a passing blocked_check probe (work Step 1), a `do-work clarify` confirmation, or a manual edit.
blocked_by: "LM Studio running locally"   # free text naming the condition (always YAML-quoted — raw user text). Legacy note for the board: an old id-LIST value renders joined for display and is NOT a dependency edge — dependency gating is `depends_on` only.
blocked_at: 2026-07-18T10:00:00Z          # stamped on every flip to blocked — the age anchor the exit summary, board drawer, and forensics read (no enforcement threshold; external conditions legitimately take weeks)
blocked_check: "curl -sf http://localhost:1234/v1/models"   # OPTIONAL shell probe (always YAML-quoted). User-authored content, run VERBATIM by work Step 1 (exit 0 ⇒ unblock to pending; any non-zero / timeout / unreadable ⇒ stays blocked). Absent ⇒ manual/clarify unblock only.

# Set on ANY status flip that has no dedicated *_at stamp of its own — that
# condition is the rule, the writers are illustrative: answered → pending
# (clarify Step 5), unblock → pending (clarify Step 5.5, work Step 1 probe —
# both REMOVE blocked_at, so this is the only trace of when the flip happened),
# manual/stuck resets back to pending. Flips with a dedicated stamp (claim →
# claimed_at, reserve → reserved_at, blocked → blocked_at, terminal →
# completed_at) do NOT write it. Display-only: the board's state timer prefers
# it over created_at/file-mtime for pending-tier cards ("updated … · 3m"); no
# pipeline logic reads it. Timestamp rule applies (current UTC instant).
status_changed_at: 2026-07-22T20:38:00Z

# Set by work action when finished. STAMPING RULE: every flip to a terminal
# status (completed / completed-with-issues / failed / cancelled) MUST stamp
# completed_at with a UTC ISO instant, plus commit with the implementation
# hash in a git repo. These two fields are the ONLY sources the board resolves
# a terminal REQ's completion instant from (no file-mtime fallback); a
# terminal REQ missing both — or carrying an unparseable completed_at or a
# hash git can't resolve — is flagged as a completion anomaly by do-work board
# (all three modes: serve, static, summary).
completed_at: 2025-01-26T10:45:00Z   # required on every terminal flip — UTC ISO instant
status: completed | completed-with-issues | failed
commit: abc1234               # required in a git repo — implementation commit hash (see work.md's Commit Phase write-back)
error: "Description"          # Set when a REQ failed; RETAINED verbatim if that failed REQ is later cancelled via do-work abandon — the surviving failure signal on a status: cancelled REQ, NOT drift to strip
error_type: intent|spec|code|environment   # Set with `error` on failure; likewise retained on a failed→cancelled flip

# Set by abandon action (do-work abandon — user-directed won't-do decision).
# Two entry paths: a not-yet-finished REQ (pending / pending-answers / blocked / ...), and an
# already-archived `failed` REQ resolved after the fact — the latter keeps its `error`/`error_type`
# (above) alongside status: cancelled, so error-on-cancelled is valid data, not corruption.
status: cancelled             # terminal, NOT successful; the reason lives in the REQ body's `## Cancelled` section
completed_at: 2025-01-26T10:45:00Z  # stamped (or, on the failed→cancelled path, re-stamped to the cancellation instant) — the terminal timestamp the board's recently-done window reads

# Set by kb-lessons handoff (work.md's Lessons-Capture Phase in pipeline mode / review-work.md's Self-Validation & Lessons Learned step standalone). Optional; absent on REQs that predate the handoff.
kb_status: promoted | pending | declined | skipped
kb_entry: REQ-042-lesson-slug.md   # filename only (survives bkb moves from inbox/ to capture/ to processed/); present only when kb_status: promoted

# Set by the board's Testing view (do-work board serve — actions/board.md Step 6). Optional; the testing track is orthogonal to `status`: the board never writes `status`, and the work pipeline never writes these. Absent = not tested yet.
testing_status: in-testing | tested | returned   # who-tested-what tracking for finished REQs
tested_by: "Alice"                # tester profile from do-work/testers.md (raw user text, always YAML-quoted)
testing_updated_at: 2026-07-17T10:00:00Z   # stamped by the board server on every transition
testing_feedback: "…"             # present only while testing_status is returned (one-line double-quoted scalar; newlines as \n escapes)
---
```

## Schema Read Contract

Nine fields above are enum-or-boolean-valued, and an audit of `0.76.2`'s `dependencies:` → `depends_on` patch surfaced that several silently swallow natural typo variants from sister conventions (snake_case vs kebab-case YAML, `done`/`finished`/`closed` as English glosses of `completed`, lowercase route letters, etc.). Pure silent-alias is risky for enum values because an unknown value should not be quietly remapped — it should leave a footprint. Every read site in this file (and in `actions/roadmap.md`) honors a uniform **normalize-and-warn contract** for these fields:

1. **Normalize first.** Apply the per-field alias map below. If a canonical match results, use it silently.
2. **Warn-on-fallback.** If after normalization the value still doesn't match the canonical enum, emit:

   ```
   ⚠ {field}: '{value}' not recognized — expected one of [{enum}]. Treating as '{default}'.
   ```

   and proceed with the documented default.
3. **Never silently drop.** The warning is the missing feedback channel that allowed `dependencies:` to go unnoticed pre-0.76.2. Warnings render in the queue-status summary block (Step 1) or, for fields read outside Step 1, alongside the operation that triggered the read.

| Field (read sites) | Canonical enum | Normalization | Default on unknown |
|---|---|---|---|
| `domain` (Step 4 Route C plan-agent spawn, Step 6 crew load, Step 7 review-work spawn) | `frontend`, `backend`, `ui-design`, `general`, `security`, `testing` | `back-end`/`back_end` → `backend`; `front-end`/`front_end` → `frontend`; `ui_design` → `ui-design`; `sec` → `security`; `test` → `testing` | `general` |
| `status` (Step 1 scan + categorization, Step 8 archive trigger, abandon action, reserve action) | `pending`, `claimed`, `reserved`, `completed`, `completed-with-issues`, `failed`, `cancelled`, `pending-answers`, `blocked`, `blocked-archive-collision`, `blocked-dependency-cycle` | `done`/`finished`/`closed` → `completed`; `canceled`/`abandoned`/`wont-do`/`wontfix` → `cancelled` | skip REQ at Step 1 with the warning text — never claim or archive an unrecognized status silently |
| `route` (Step 3 dispatch, Step 5.5 scope declaration, Step 7 scope-drift comparison) | `A`, `B`, `C` | lowercase `a`/`b`/`c` → uppercase | treat as needing re-triage in Step 3 |
| `caveman` (Step 6 crew load) | `false`, `true`, `lite`, `full`, `ultra` | truthy strings (`yes`/`on`) → `true`; `light` → `lite` | `false` |
| `maintenance` (Step 6 crew load) | `true`, `false` (YAML boolean) | truthy strings (`yes`/`on`/`t`) → `true`; `no`/`off`/`f` → `false` | `false` (Step 6 maintenance crew not loaded) |
| `tdd` (Step 6 testing-crew load, Step 6.5 TDD-evidence gate; emission validated in `actions/capture.md`) | `true`, `false` (YAML boolean) | `test_first`/`yes`/`on`/`t` → `true`; `no`/`off`/`f` → `false` | `false` (Step 6 testing crew not loaded; Step 6.5 gate not enforced) |
| `error_type` (Step 8 failure classification, Step 8 upstream-failure short-circuit, forensics) | `intent`, `spec`, `code`, `environment` | (no common typo aliases identified) | `code` |
| `kb_status` (kb-lessons handoff — work.md's Lessons-Capture Phase / review-work.md's Self-Validation & Lessons Learned step; roadmap lessons rollup) | `promoted`, `pending`, `declined`, `skipped` | `skip` → `skipped`; `rejected` → `declined` | `pending` |
| `testing_status` (board Testing view — `tools/queue-kanban` parser + `/api/testing/status` writes; no work-pipeline read sites) | `in-testing`, `tested`, `returned` | `in_testing`/`in testing`/`testing`/`selected-for-testing` → `in-testing`; `returned-with-feedback`/`returned_with_feedback` → `returned` | treat as not-tested (Ready to test) with an invalid flag + data warning |

**Write paths are unaffected.** Step 2 claim, Step 8 archive, Step 8 follow-up generation, the kb-lessons handoff, and capture emission always write the canonical key and canonical enum value — never an alias, never the typo'd input. The normalize-and-warn contract is read-only.

**List-valued path fields are outside this contract and are read verbatim.** `prime_files`, `write_set`, and any future path list have no canonical vocabulary to normalize against — no alias map, no case folding, no path canonicalization, no warning. A reader takes the strings as written. (`depends_on`, `related`, and `blocked_by` are likewise row-less here; `depends_on`/`related` carry alias keys, documented on their schema lines above, but their *values* are also read verbatim.)

### Terminal-success status set

**A REQ counts as *terminally successful* when its `status` is `completed` or `completed-with-issues`.** This is the canonical set every reader that selects "completed work" must honor — `completed-with-issues` is terminal and counts toward UR completion (it just carries known follow-ups, per `actions/work.md` Step 8), so a filter that accepts only the literal `completed` silently drops remediated-with-issues work. `failed` is terminal but **not** successful — success-readers exclude it.

The trigger is the *condition above*, not the caller list: **any reader that filters for "the completed/most-recent work" inherits this contract.** The known consumers are illustrative, not exhaustive — `actions/cleanup.md` (UR close), `actions/ai-report.md` (report target), `actions/review-work.md` (standalone target), and `actions/commit.md` (REQ association); `actions/forensics.md` and `actions/roadmap.md` already honor both. When adding a new reader, accept both values and point back here — hand-enumerated caller lists go stale silently, which is why the condition, not the list, is the contract.

### Terminal-resolved status set

**A REQ counts as *terminally resolved* when its `status` is `completed`, `completed-with-issues`, or `cancelled`.** This is the set archive-sweep and UR-closure readers honor — any reader deciding whether a REQ still needs work inherits it, so the list that follows is illustrative, not exhaustive: `actions/cleanup.md` Pass 0 + Pass 1, `actions/work.md` Step 8's UR-final check, `actions/forensics.md` Check 4, and this file's Composed Exit Summary. `cancelled` records a deliberate won't-do decision — made via `do-work abandon` — so it archives like finished work and must never hold a UR open the way `failed` does. Three boundaries keep the sets honest:

- `cancelled` is **not** successful. Success-readers (the Terminal-success set above) exclude it — a cancelled REQ is never a review-work target, an ai-report subject, or a commit association.
- `cancelled` does **not** satisfy `depends_on` gating. A dependent presumably needed the cancelled REQ's output; the abandon action surfaces dependents at cancellation time so the user can cascade the cancellation or re-point `depends_on`.
- `failed` stays outside this set: it is terminal and unsuccessful, but it signals work that *should* have happened — Step 8's failure classification may spawn a follow-up REQ. A follow-up does the recovery work but **never flips the original out of `failed`** (nothing does so automatically), so a UR holding a `failed` REQ stays open even after that follow-up completes. The one transition out of `failed` is `do-work abandon REQ-NNN` (`actions/abandon.md`), which flips it to `cancelled` — and therefore into this set — while preserving the failure record (`error`/`error_type` and a `## Cancelled` note). So a UR held open by a `failed` REQ closes only once that REQ is cancelled: after its follow-up has done the needed work, or when no follow-up is wanted at all. `failed` itself never counts as resolved. This is the canonical statement of that resolution **rule**: any reader that decides whether a `failed` REQ still holds its UR open cites this set by reference and must not restate or fork the set, or re-derive the rule, as a competing definition. The condition is the trigger, not a caller list — the known such readers are illustrative, not exhaustive (`actions/cleanup.md` Pass 0 + Pass 1, `actions/forensics.md` Check 4, `actions/work.md` Step 8's UR-final check, and this file's Composed Exit Summary), and `actions/abandon.md` is the rule's sole *writer*. (A user-facing report or finding line that *points* the user at `do-work abandon` as the remedy — and may state the one-line reason why, e.g. that a completed follow-up never resolves the original — is a pointer, not a competing definition, and is expected; `actions/cleanup.md` Pass 1's open-UR report and `actions/forensics.md` Check 6 do exactly this. That is not what this prohibition forbids; what it forbids is a reader re-deriving *when a REQ counts as resolved* as its own competing definition.)

### Target ID Resolution

**The shared token grammar for every action that takes id arguments** — `do-work run`, `do-work abandon`, `do-work reserve`/`release`, and `do-work roadmap`. A caller cites this contract for token shapes and UR expansion instead of restating them; the condition, not a caller list, is the trigger, so any future id-taking action inherits it.

- **Token shapes.** `REQ-` + digits and `UR-` + digits, **case-insensitive** (`req-42`, `UR-011`, `Ur-11` all resolve).
- **A `UR-NNN` token expands to its member REQs** by scanning `user_request:` frontmatter across the live locations *the calling action already searches* — which locations those are is the caller's business; the scan key never is. It is **never** the UR's own `requests:` array, which is a capture-time record and explicitly not a membership predicate (`actions/capture.md` → "The `requests:` array is the capture-time record only"). Every prior bug here — REQ-048, REQ-058, REQ-059 — came from reading that array.
- **An id that resolves to nothing** is reported by id and skipped, never silently dropped.
- **An argument list that resolves to an empty set stops the action.** It never falls through to a whole-queue default (a full run, a full survey). Expansion adds a recognized shape; it never makes an unrecognized or empty argument permissive.
- **Expansion widens *which* REQs an action reaches; it never relaxes how any one is treated.** Each caller applies its own per-REQ gates — dependency-readiness, status refusals, reservation skips, confirmations — to every expanded member exactly as if that REQ had been named directly.

## Crash Recovery (Step 1)

**Crash Recovery:** Before checking the queue, look inside `do-work/working/` for any `REQ-*.md` files. If any exist, a previous run may have been interrupted — but a live coexisting session's claimed file looks identical to a crashed one, so a concurrency gate runs first.

**Concurrency gate (apply every time this runs, including on a Step 10 → Step 1 loop iteration within the same session — it re-reads the lock file fresh each time, so it never depends on this session's own memory of an earlier lock-guard choice).** Take the `working/` listing **first, then** read `do-work/orchestrator-lock.json` — this order is load-bearing, not stylistic. Step 2 writes a REQ's claim into the lock *before* moving the file into `working/` (`actions/work.md` Step 2 substep 1), so any file already present when the listing is taken is guaranteed to have its claim visible in a lock read taken after it; reading the lock first reopens the race, because a REQ claimed-and-moved between that stale lock read and the scan would look unclaimed and get re-queued out from under a live session. For each `REQ-*.md` found in the listing, extract its id and check whether **any** fresh claim set in the lock — the holder's or a coexisting entry's, **including this session's own** — currently claims it: the id appears in the top-level holder's `claimed_reqs` with the holder's `heartbeat_at` ≤45 minutes old, or it appears in any `coexisting_sessions[].claimed_reqs` with that entry's `heartbeat_at` ≤45 minutes old. Read every entry's `claimed_reqs` *list* (the legacy scalar `claimed_req` names only the first of a possibly multi-claim set — never gate on it alone; the one exception is an *older single-claim lock* carrying only `claimed_req` with no `claimed_reqs` field at all, where that scalar is the sole claim record and is the fallback — matching `actions/cleanup.md`'s Pass 0 gate, so both gates tell the same mixed-version story). **Skip that file** if so — it belongs to a live claim, this session's or another's. **Freshness alone gates, deliberately including this session's own claims:** `session_id` is generated fresh at acquisition/join, so a dead predecessor's claim carries a *different* id and ages past 45 minutes into recovery, while this session's own heartbeat is fresh by construction — so on a Step 10 → Step 1 loop iteration this session's still-building co-dispatched siblings are skipped rather than stripped and re-queued out from under themselves. Recover every other file below, including files whose claiming session has itself gone stale past 45 minutes — that's the self-heal path. Full lock file shape: `actions/work-reference.md` → **Concurrent-Orchestrator Lock Guard**.

For each REQ file that passes the gate (not freshly claimed by any live session — another session's *or this session's own*, per the gate above):
1. Reset frontmatter: set `status` to `pending`, **unless** the REQ file contains a `## Open Questions` section with at least one unresolved `- [ ]` item — in that case, restore `status` to `pending-answers`. (If the `## Open Questions` section exists but all items are already `[x]` or `[~]`, or if no `## Open Questions` section exists at all, set `status` to `pending`.) **Exception — a recovered REQ that already carries `status: blocked` with a `blocked_by` condition stays `blocked`** (the mid-run blocked flip completed its frontmatter write before the crash; its condition is unchanged and it must not be silently promoted to runnable). Remove `claimed_at` and `route`; leave `blocked_by`/`blocked_at`/`blocked_check` intact. **Clear `write_set` only when this REQ actually has a `## Scope` section** — check for it here, while it still exists (substep 2 strips it next). `## Scope` is `write_set`'s only source (**Scope Declaration Template (Step 5.5)**, below), so a mirror that outlives its source would feed the parallel-dispatch gate a set with nothing left to re-derive it from; absent reads as *overlaps everything ⇒ serialize*, the safe post-recovery default. **With no `## Scope`, preserve `write_set`** — it is capture-seeded, not a mirror (the REQ crashed before Step 5.5, or it is a Route A REQ that never runs Step 5.5 at all — `actions/work.md` Step 5.5). Nothing downstream ever re-seeds that field, so clearing it destroys user- and capture-authored frontmatter and serializes the REQ for good.
2. Strip sections generated during the interrupted run: remove `## Triage`, `## Exploration`, `## Plan`, `## Scope`, `## Pre-Flight`, `## Implementation Summary`, `## Qualification`, `## Testing`, `## Review`, `## Lessons Learned`, `## Orientation`, `## Decisions`, and `## Discovered Tasks` sections (and their content) if present — these may be incomplete or stale from the crash. Leave `## Open Questions` and user-authored content intact. (Stripping `## Scope` here is what makes substep 1's `write_set` decision conditional on it — read that decision before this substep runs, never after.)
3. Move the REQ back to `do-work/queue/`

**Worktree sweep — runs once, not per file in the loop above.** Only applies where the run used worktree dispatch (**Worktree Dispatch Mode (Step 1)**, below); skip it entirely when the repo has no `git worktree` support or no `worktree-agent-*` names exist. A leftover branch can outlive its `working/` file (the REQ archived, the branch didn't), which is why this sweeps names rather than iterating the files above. Run `git worktree prune` first so already-deleted directories don't surface as ghosts, then enumerate `git worktree list --porcelain` and `git branch --list 'worktree-agent-*'` — a branch with no worktree and a worktree with no branch each count. For each leftover, read the REQ id out of its `worktree-agent-REQ-NNN-…` name and:

- **Exempt it** if that id is exempted by the same live-claim gate applied above — it belongs to a builder that is still running, not to an interrupted one. Report the exemption in the same words the gate uses. **Where the exempting claim is this session's own**, exempt only the name this session recorded as that REQ's operative name (**Worktree Dispatch Mode (Step 1)**, below): a *different* name carrying the same id is the predecessor a collision variant dispatched around, not the live builder's worktree, so it goes through the branches below instead of being silently skipped. Another session's operative name isn't knowable from here, so its claim still exempts by id alone.
- **Merged into the integration branch** (`git worktree remove <path>` succeeds, then `git branch -d <branch>` succeeds — both run from the integration branch, because `-d` is HEAD-relative; see **Worktree Dispatch Mode (Step 1)**, *Cleanup — happy path*): pure residue. Remove it mechanically, run `git worktree prune` again, and report `Removed merged worktree <name>`.
- **Unmerged or dirty** (either command refuses): **report it and move on — never `-D`, never `--force`.** The branch may hold the only copy of a builder's work; deleting it belongs to `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)**, which asks first. A reported unmerged leftover does **not** block re-dispatch of its REQ — the Naming rule's collision variant (**Worktree Dispatch Mode (Step 1)**, below) dispatches the recovered REQ under a fresh unique variant, so the two coexist until Pass 5 resolves the leftover.
- **Any other worktree name**: not ours. Leave it alone.

The gate above now exempts **every** fresh claim — other live sessions' *and this session's own* — because it reads each entry's `claimed_reqs` list and freshness alone gates (**Crash Recovery (Step 1)**, above). So a co-dispatched batch's still-building siblings are exempt on a Step 10 → Step 1 sweep exactly as another session's claims are; a leftover branch is swept only once its REQ's claim has aged out or been released. Co-dispatching several worktree builders from one session is now representable — each builder's REQ is a separate id in this session's `claimed_reqs` — and is governed by Step 1's parallel-dispatch gate (`actions/work.md`).

Once every file this session is allowed to touch has been recovered (a live session's file may remain in `working/` — that's expected, not an error), proceed with finding the next request.

## Worktree Dispatch Mode (Step 1)

**Optional, advanced harnesses only.** Each builder runs in its own git worktree on its own branch; the orchestrator stays in the main tree, merges builder branches in dependency order, and remains the only writer of `do-work/` state. It is the isolation layer beneath Step 1's parallel-dispatch gate (`actions/work.md`): worktrees change the *consequence* of a write-set overlap — merge cost instead of interleaved writes — not the gate itself, which still decides what may be co-dispatched. **Nothing here requires concurrency — a single builder in a worktree is the same one-claim loop**, and that is how the mode ships by default: the isolation is worth having even for one builder. Real co-dispatch is now representable — each concurrent builder's REQ is a separate id in the orchestrator's `claimed_reqs` list — and is governed by Step 1's parallel-dispatch gate (`actions/work.md`).

**Precondition, then degrade silently.** Probe `git worktree list` (a non-zero exit, or an unrecognized subcommand, means no worktree support) and confirm the harness can run an agent against a working directory you choose. If either is missing, run the serial loop exactly as documented and say nothing further. Unlike `actions/board.md`'s Go check — which reports and stops, because there the toolchain *is* the capability — this mode's absence is not an error, so it must never surface as one.

**Naming.** The worktree directory's basename and the branch name are the **same string**: `worktree-agent-REQ-NNN-<suffix>`. Embedding the REQ id is what lets any later sweep correlate a leftover with its REQ by name alone; sharing one string is what makes a single grep find both the directory and the branch. Derive `<suffix>` from the REQ's filename slug, lowercased and reduced to `[a-z0-9-]`, **as a text operation before you compose any shell command** — never pipe raw REQ text through `tr`/`sed` inside a quoted command line, where an apostrophe breaks the quoting and becomes an injection vector. **On a name collision at creation** — `git worktree add` or branch creation fails because `worktree-agent-REQ-NNN-<suffix>` already exists as a leftover (typically a crash-recovered REQ re-dispatching into the name its own reported-but-not-deleted leftover still holds) — do **not** delete or force. Dispatch under a fresh unique variant: append an **incrementing numeric token** to the suffix — `-2`, then `-3` if that collides too — keeping the `worktree-agent-REQ-NNN-` prefix intact so both names still correlate to the REQ by name (this is what the sweeps grep for). **One scheme, not a choice:** a free pick between a counter and a timestamp token lets two runs shape the same collision differently, so the name stops being readable evidence of what happened. Report the coexistence and leave the original leftover to its owners — the crash sweep above if it turns out merged, `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)** if unmerged.

**The name actually created is this REQ's operative name.** Whatever `git worktree add` succeeded with — the derived name in the common case, the variant after a collision — is the one string every later worktree/branch operation uses: the hand-back merge's `git merge` argument (below), Step 8's `git worktree remove` and `git branch -d` (*Cleanup — happy path*, below), the crash sweep's own-session bookkeeping (**Crash Recovery (Step 1)**, above), and anything reported back to the user. This reference calls it **`<operative_name>`** wherever a command needs it; the worktree path is that same name under the worktrees parent directory. Hold it exactly the way `<pre>`/`<merge_hash>` and the lock's `session_id` are held — known from this session's own context and re-typed as a literal into each fresh command, never a shell variable (**Hold both endpoints as re-typed literals**, below). Nothing persists it, because nothing outside this session's run consumes it: the sweeps and `actions/cleanup.md` Pass 5 discover leftover names by enumerating git, never by re-deriving them. **Re-deriving the name from the slug at cleanup time is the failure this closes** — after a variant dispatch the derived string names the *leftover*, so `git worktree remove`/`git branch -d` target unmerged work, refuse, and halt the run on a false "the merge was skipped or lost" alarm while the variant worktree is never cleaned at all. **With no collision the operative name *is* the derived name**, so the common path and serial mode behave exactly as before.

**Where worktrees live: outside the repo working tree.** A sibling directory (`../<repo>-worktrees/worktree-agent-REQ-NNN-…`) or a scratch directory — never nested inside the repo. Inside, it is a second checkout of the repo sitting in the repo: `actions/cleanup.md` Pass 3a scans the filesystem for any `do-work/` directory outside the project root, and where the consumer commits `do-work/` it would find the builder's checkout of one and try to relocate it into the canonical queue. The extra tree also reads as stray untracked residue to every status and stray-file check downstream. That is a corruption path, not just untidiness.

**State stays home.** `do-work/` — the queue, `working/`, the orchestrator lock, `CHECKPOINT.md` — exists in the main tree only. Builders receive their REQ brief in the dispatch prompt and never read or write queue state from inside a worktree. Untracked files do not propagate into a new worktree, so on the common install (where `do-work/` is untracked) a builder simply finds nothing there. The trap is the other install: where a consumer **commits** `do-work/`, the worktree carries a *stale snapshot* of the queue as it stood at the branch point. Treat that snapshot as absent — never read a status from it, never write to it. Every claim, status flip, and archive move happens in the main tree, by the orchestrator.

**Sole integrator.** Builders never write the main tree or its branch. Each commits on its own branch and hands back its file manifest; the orchestrator merges, in dependency order (`depends_on` first). A shared file needing one line of wiring — a `<link>` in a shared template, a registry entry, an export barrel — is an **integration seam**: the builder hands back the exact line and where it goes, and the orchestrator applies it in the main tree **inside the merge commit** — step 3 of the hand-back sequence below, which is what keeps the seam inside this REQ's merge range. A builder that edits the seam itself converts a scheduling decision into a conflict for every sibling.

**Merge, never rebase.** Integrate with `git merge --no-ff <branch>` (the full invocation, which adds `--no-commit` so the integration seam has somewhere to go, is step 2 of the hand-back sequence below). Rebasing rewrites the builder's commits, so `git branch -d` no longer recognizes the work as merged and the free merged-ness assertion below is destroyed. `--no-ff` also preserves the merge commit as the "integrated by orchestrator" provenance record.

**When to merge, and the range every evidence step reads.** The orchestrator merges each builder branch **at hand-back — end of Step 6, before Step 6.25 (Implementation Summary)** — so every downstream evidence step (6.25, 6.3, 6.5, 7, 8, 9) observes one integrated main tree, not a clean one. Any position after 6.25 leaves qualify (Step 6.3) and review (Step 7) reading a clean main tree with nothing to check. Run this four-step hand-back sequence on the integration branch:

1. **Capture `<pre>`** — run `git rev-parse --short HEAD` and read the printed hash. It is the integration tip before this REQ's **first** merge and the lower bound of the merge range. Capture it **once per REQ**; a remediation re-merge does not re-capture it (below). Never recover it afterwards as `HEAD^1` or live `HEAD` — both move as soon as the orchestrator commits the changelog (Step 9) or merges the next sibling.
2. **Merge without committing** — `git merge --no-ff --no-commit <operative_name>` (this REQ's operative name, *Naming* above — the branch the builder was actually dispatched on, which is the collision variant where there was one), then resolve any conflict. `--no-ff` forces a merge commit even where the branch could fast-forward (Merge, never rebase, above); `--no-commit` is what leaves the integration seam somewhere to go. If git answers `Already up to date.` the builder committed nothing: no `MERGE_HEAD` is set, so a `git commit` here would either fail or fabricate a non-merge commit — stop and treat the hand-back as empty instead.
3. **Apply the integration seams, then commit** — stage the handed-back seam lines (Sole integrator, above) and `git commit`. Folding the seam into the merge commit is the only placement that puts it inside the merge range: a seam committed *after* the merge is that merge commit's **child**, hence outside `<pre>..<merge_hash>`, and qualify, review, and Step 9's validation would never see it.
4. **Capture `<merge_hash>`** — `git rev-parse --short HEAD` on the commit just made. It is the upper bound of the merge range and the hash Step 9 writes into the REQ's `commit:` field.

**Hold both endpoints as re-typed literals, never as shell variables.** Shell state does not survive between prescribed command blocks, and the consumers sit in later blocks with model round-trips in between — a `"$pre..$merge_hash"` composed in a fresh shell expands to `".."`, which git rejects. Hold both hashes exactly the way the lock's `session_id` is held (**Concurrent-Orchestrator Lock Guard**, below): known from this session's own context and re-typed into each fresh command. `tools/checks/qualify.sh` hard-FAILs on a range it cannot resolve rather than reading an empty diff, so a lost endpoint surfaces as a qualification failure naming the range instead of a vacuous pass.

**The merge range is `<pre>..<merge_hash>`**, and in worktree dispatch mode it replaces the working diff wherever an evidence step reads changes: **Step 6.3** (`tools/checks/qualify.sh` via `DO_WORK_DIFF_RANGE="<pre>..<merge_hash>"`), **Step 7** (review's Get-the-Diff), **Step 8** (post-merge verification, below), and **Step 9's** staged-list validation all consume it. Use the *captured* `<merge_hash>` as the upper bound, never live `HEAD`: HEAD moves the moment the orchestrator commits the changelog (Step 9) or merges the next sibling, and under **per-batch** verification every sibling is merged *before* this REQ is verified, so `<pre>..HEAD` would sweep in unrelated commits and misattribute them to this REQ. `<pre>` is `<merge_hash>`'s first-parent ancestor — its direct first parent after a single merge — so `merge-base(<pre>, <merge_hash>) == <pre>` and this is already the merge-base form. Serial mode is unchanged: no range is set and every step reads the working/staged diff exactly as before.

**Remediation re-merges: the range is cumulative.** A failed review sends the REQ back to Step 6 (`actions/work.md`) and the builder's fix branch is merged again. Repeat steps 2–4 above but **not step 1**: keep the first `<pre>` and re-capture only `<merge_hash>` from the newest merge commit, making the range `<pre₁>..<merge_hash₂>` — first pre-merge tip to latest merge — which covers the original work *and* the fix. Re-capturing `<pre>` would instead give `<pre₂>..<merge_hash₂>`, covering only the fix delta: review would read the fix in isolation and every originally-touched file would WARN as listed-but-not-in-the-diff. **Step 9 records the latest `<merge_hash>`** in `commit:`. The cumulative range's one cost is over-inclusion — an orchestrator commit that landed on the integration branch between the two merges (a sibling REQ's changelog commit) falls inside it and surfaces as an undeclared touch for your judgment. That is the safe direction of error: it shows up and gets judged, where under-inclusion silently hides the REQ's own work.

**Post-merge verification, before archive.** Each builder verified its own branch; nobody has verified the merged result. Re-run the REQ's acceptance checks — the `## Red-Green Proof` GREEN condition, the project's test command, whatever `## Scope` named — against the merged main tree **before** Step 8 archives anything (`actions/work.md`). Four individually-green builders must not compose into a red main that archives as done. Pick the granularity by the rollback you are prepared to make: **the unit you verify is the unit you roll back.** **Per-merge is the default whenever more than one REQ is in flight** — Step 8 archives and commits each REQ as it verifies (`actions/work.md`), so once a co-dispatched batch's earlier REQs are archived the batch has no revertible unit left; per-batch is permitted only if every REQ's Step 8 is held until the batch's last merge verifies. A red merged state is not an archive-plus-follow-up: stop, revert to the last verified state, and re-dispatch.

**Cleanup — happy path (Step 8).** After the archive move, remove the builder's worktree and branch **by this REQ's operative name** (*Naming*, above — never re-derived from the slug here): `git worktree remove <path>` (no `--force`, where `<path>` is the worktree whose basename is `<operative_name>`), then `git branch -d <operative_name>`, then `git worktree prune`. Run `branch -d` **from the integration branch you merged into**: `-d` tests merged-ness against the current HEAD (or the branch's configured upstream), so from anywhere else a perfectly-merged branch refuses and an unmerged one can pass — "refusal = unmerged" silently becomes "refusal = wrong branch." Both refusals are signal, not friction: `worktree remove` refuses on a dirty worktree (uncommitted builder work you are about to lose), `branch -d` refuses on an unmerged branch (a merge that was skipped or lost). **Never `-D`, never `--force`.** Report the refusal and stop — forcing destroys the only evidence that the integration didn't happen.

**Cleanup — crash path.** Leftovers from an interrupted run are swept by **Crash Recovery (Step 1)** above: already-merged ones are removed mechanically, unmerged ones are reported and never auto-deleted. Discarding unmerged builder work belongs to `actions/cleanup.md` → **Pass 5: Orphaned Worktrees (consent-gated)**, which asks first and only acts when a human can answer.

## Concurrent-Orchestrator Lock Guard (Step 1)

Guards the one-orchestrator-per-queue convention: a second `do-work run` starting on the same working tree must detect the first before claiming or processing anything. Not a distributed lock — but every update to the lock file is a read-modify-write on a file with more than one writer, so it must be serialized (**Serialized lock updates**, below). Plain unguarded reads and writes are *not* sufficient. (Origin incident: two sessions ran `do-work run` on the same tree on 2026-07-01; the second committed and archived the first's in-flight REQ with a hollow paper trail. The collision merged by luck — nothing detected it.) **Design note (REQ-018 remediation):** the lock records exactly one HOLDER plus zero or more `coexisting_sessions` — a session that picks "proceed anyway" gets its own tracked entry instead of displacing the holder, so the holder's `session_id` never stops matching and the holder is never disrupted.

**Lock file:** `do-work/orchestrator-lock.json`. Re-derive this path deterministically every time you touch it — do not carry a shell variable across separate command blocks. `mkdir -p do-work` first if the directory doesn't exist yet.

```json
{
  "session_id": "do-work-20260727T083332Z-51234",
  "started_at": "2026-07-27T08:33:32Z",
  "heartbeat_at": "2026-07-27T08:41:10Z",
  "claimed_reqs": ["REQ-018"],
  "claimed_req": "REQ-018",
  "coexisting_sessions": []
}
```

With one session coexisting — this is the shape you edit when refreshing or removing **your own** entry, so match by `session_id` and change only that element, never the holder's fields and never a sibling entry:

```json
{
  "session_id": "do-work-20260727T083332Z-51234",
  "started_at": "2026-07-27T08:33:32Z",
  "heartbeat_at": "2026-07-27T08:41:10Z",
  "claimed_reqs": ["REQ-018"],
  "claimed_req": "REQ-018",
  "coexisting_sessions": [
    {
      "session_id": "do-work-20260727T085501Z-52890",
      "joined_at": "2026-07-27T08:55:01Z",
      "heartbeat_at": "2026-07-27T09:02:44Z",
      "claimed_reqs": ["REQ-042"],
      "claimed_req": "REQ-042"
    }
  ]
}
```

`session_id` is `do-work-$(date -u +%Y%m%dT%H%M%SZ)-$$` — a UTC timestamp plus the current shell's process id, generated once at acquisition (or at joining, for a coexisting session). The orchestrating agent already knows this string from its own context for the rest of the run (it isn't a shell export that needs to survive between command blocks) — every later touch just re-types the same value into a fresh, deterministically-derived command. The top-level `session_id`/`started_at`/`heartbeat_at`/`claimed_reqs` describe the recorded **holder**. **`claimed_reqs` is the canonical claim field** — a JSON list of every REQ id this holder currently holds a claim on: the ids in its **dispatch record** (**Dispatch bookkeeping is the recompute's source**, below), which is what it has in `do-work/working/` plus, in the window Step 2 opens, one it has claimed but not yet moved there (one id in the serial default, several under parallel dispatch, empty `[]` between REQs). **`claimed_req` is retained as a derived legacy mirror** — always `claimed_reqs[0]`, or `null` when the list is empty — so an older reader that only knows the single-claim field keeps working. It is additive back-compat and is *never* itself made array-shaped (a reader comparing a string to an array would match nothing and exempt nothing — catastrophic); every writer that changes `claimed_reqs` re-derives `claimed_req` in the same write. `coexisting_sessions` lists other sessions running alongside the holder via **proceed anyway** (below) — each entry is `{session_id, joined_at, heartbeat_at, claimed_reqs, claimed_req}` with the same canonical-list/derived-mirror relationship, refreshed independently of the holder's own fields; empty (`[]`) when no session has chosen to coexist. All timestamps follow the Timestamp rule above (current UTC instant).

**Serialized lock updates (every write, no exceptions).** More than one session writes this file — a coexisting session refreshes its own entry while the holder refreshes the holder fields, and both the heartbeat prune and a take-over rewrite *other* sessions' entries. So "each writer only changes its own logical fields" describes intent, not atomicity: if the holder reads at T0, a coexisting session records its `claimed_reqs` at T1, and the holder writes its T0 image back at T2, that claim is silently lost. The Crash Recovery per-file gate then sees no live claim on a REQ that is actively being built and re-queues it — the exact 2026-07-01 failure this guard exists to prevent.

Every path that writes the lock — acquisition, take-over, the **proceed anyway** append, each heartbeat refresh, the opportunistic prune, and release — wraps its read-modify-write in this mutex, and re-reads the file *inside* the critical section rather than reusing an earlier read:

```bash
lock_path="do-work/orchestrator-lock.json"
mutex_path="do-work/orchestrator-lock.json.mutex"

# `mkdir` is atomic on POSIX: exactly one contender wins, the rest fail.
until mkdir "$mutex_path" 2>/dev/null; do
  if [ ! -d "$mutex_path" ]; then
    # mkdir failed with no contender holding the path — permissions, read-only fs,
    # missing parent. Retrying cannot succeed; report and stop instead of spinning
    # forever, and NEVER fall through to an unserialized lock write.
    echo "orchestrator-lock mutex: cannot create $mutex_path — fix the environment, then re-run" >&2
    exit 1
  fi
  # A session that died mid-update must not wedge the queue. The mutex directory's mtime
  # advances only when an entry is created inside it (the owner token at acquisition, then
  # the staged lock image below), so `-mmin -1` measures time since the owner's last
  # filesystem progress — it proves the mutex's AGE, not that its owner is dead. The
  # critical section can legitimately span a model round-trip, so a slow-but-live owner is
  # indistinguishable from a dead one, and an abandoned mutex is still guaranteed
  # reclaimable a minute after its owner's last write. Reclaiming on age alone is safe
  # because the staged image lives INSIDE the mutex directory: the `rm -rf` here destroys
  # it, so an evicted owner's publishing rename can only fail — it cannot land a stale
  # image over its successor's write. This is the ONLY break path. There is deliberately no
  # attempt-count override on a younger mutex: eviction is now lossless but not free, since
  # every early eviction of a live owner discards a completed read-modify-write cycle.
  if [ -z "$(find "$mutex_path" -maxdepth 0 -mmin -1 2>/dev/null)" ]; then
    rm -rf "$mutex_path"
    continue
  fi
  sleep 0.5
done
# $session_id below: re-type this session's literal id — shell state does not survive
# between command blocks. The token lets release tell "still mine" from "reclaimed".
printf '%s\n' "$session_id" > "$mutex_path/owner"

# --- critical section: re-read, modify, write atomically ---
# (read $lock_path fresh here, apply ONLY this session's own field changes.
#  A branch chosen on an unserialized read must re-validate its premise against
#  this fresh read before writing — the fresh read exists to be re-judged, not
#  just to avoid clobbering other sessions' fields. An ACQUIRING session whose
#  fresh read finds a lock file that now exists lost the acquisition race, and a
#  STALE-TAKEOVER whose fresh read shows a changed holder or a recomputed age
#  ≤ 45m is aiming at a live session: in both cases write nothing, release the
#  mutex, and re-enter the procedure at the branch the fresh state selects.)
# Stage INSIDE the mutex directory, never beside $lock_path. This is what makes the
# age-based reclaim above lossless rather than merely narrow: the evicting `rm -rf
# "$mutex_path"` deletes this staged file too, and an evictor cannot `mkdir` its
# replacement mutex until that removal has completed. So by the time a successor can read
# the lock, this image is already gone and the rename below can only fail — there is no
# interleaving in which an evicted owner publishes over its successor's write. The $$
# suffix keeps the source unique to this shell, so the rename can never pick up a
# successor's staged image and publish it early.
lock_tmp="$mutex_path/staged.$$.json"
if ! printf '%s\n' "$updated_lock_json" > "$lock_tmp"; then
  rm -f "$lock_tmp"
  echo "orchestrator-lock mutex: could not stage $lock_tmp — lock NOT updated" >&2
  exit 1
fi

# Pre-write ownership re-check — the diagnostic, not the guarantee (the staging path above
# is). The age-based reclaim can evict a slow-but-live owner, so confirm the mutex is still
# ours before attempting to publish. A mismatch means this session was evicted
# mid-critical-section and its staged image is stale: discard the temp file, write NOTHING
# to $lock_path, leave the mutex alone (it is the successor's now, so exiting here
# deliberately skips the release below), and re-run this whole block from the top —
# re-deriving $updated_lock_json from a fresh in-mutex read, which re-judges the branch
# rather than just re-writing it. Never retry by re-`mv`ing the image staged before the
# eviction: that publishes the very lost update this catches.
if [ "$(cat "$mutex_path/owner" 2>/dev/null)" != "$session_id" ]; then
  rm -f "$lock_tmp"
  echo "orchestrator-lock mutex: evicted mid-critical-section — write discarded, re-acquire and redo the read-modify-write" >&2
  exit 3
fi

# The publishing rename fails closed. A failure here means the staged image is gone, which
# means this session was evicted between the re-check and now — the exact window the
# in-mutex staging closes — so it takes the same retry path as the check above rather than
# being swallowed. Never treat a failed publish as a completed update.
if ! mv -f "$lock_tmp" "$lock_path"; then
  rm -f "$lock_tmp"
  echo "orchestrator-lock mutex: publish failed (staged image evicted) — nothing written, re-acquire and redo the read-modify-write" >&2
  exit 3
fi

# Release only what this session still owns. If the critical section overran the
# one-minute reclaim bound and another session took the mutex, an unconditional
# removal here would delete the NEW owner's mutex and cascade the collapse to a
# third contender.
if [ "$(cat "$mutex_path/owner" 2>/dev/null)" = "$session_id" ]; then
  rm -rf "$mutex_path"
fi
```

The temp-file-plus-`mv` is not decoration: `mv` within a directory is an atomic rename, so a reader outside the mutex can never observe a half-written lock file. Never write `$lock_path` in place. **Bounded wait, never an override.** Worst-case wait is ~one minute — the stale-owner **age** bound (the mtime check; the mutex directory's mtime advances only when an entry is created inside it, so it tracks the owner's last write and an abandoned mutex is guaranteed reclaimable a minute after that write). Age is not death: the reclaim can evict an owner that is merely slow. Two things make that safe, and it matters which is which. The **pre-write ownership re-check** is the diagnostic — it catches the common case early and gives the evicted session a clear exit 3. The guarantee is **staging inside the mutex directory**: `$lock_tmp` is `$mutex_path/staged.$$.json`, so the evicting `rm -rf "$mutex_path"` destroys the staged image, and an evictor cannot `mkdir` its replacement mutex until that removal completes — by the time any successor can read the lock, the evicted owner's rename has no source left and can only fail. There is no interleaving in which a stale image lands over a successor's write, so the lost-update window is **closed**, not merely narrowed. That is also why the publishing `mv` must fail closed (exit 3, same retry path) rather than swallow its error. **An exit 3 is a retry instruction, not a failure:** re-run the block from the top and re-derive `$updated_lock_json` from the fresh read inside the *new* critical section; re-`mv`ing the discarded image is the one way back to the defect this design prevents. There is still deliberately no shorter attempt-count break on a younger mutex — eviction is lossless but not free, since each early eviction of a live owner throws away a completed read-modify-write cycle. If the loop exits with the `mkdir`-failure report instead, the environment is broken (permissions, read-only checkout) — surface that to the user and stop; an unserialized lock write is never the fallback.

**Never committed.** At acquisition (first write), append the local-ignore snippet from `crew-members/background-agents.md` → **Local-ignore snippet (for genuinely-transient paths)**, substituting `do-work/orchestrator-lock.json` for `<path>`, and again for `do-work/orchestrator-lock.json.*` — the glob (not an exact path) is required because the mutex directory is not a fixed single file to name: a session killed mid-critical-section leaves `do-work/orchestrator-lock.json.mutex/` behind holding an `owner` token and a `staged.<pid>.json` image whose PID no later session repeats, so only a pattern keeps the whole directory out of `git status` and off the committable surface. The mutex directory and everything staged inside it are the same kind of transient session state as the lock itself. Same mechanism `do-work/pipeline.json` already uses (`actions/pipeline.md` Step 2 substep 4) — a lock is transient session state, not part of the Trail of Intent, so it must stay out of git regardless of install layout. **The snippet cannot rescue an already-tracked lock** — ignore rules have no effect on indexed files, and `git check-ignore` passes on them, so the append silently changes nothing (`crew-members/background-agents.md` → the snippet's tracked-path caveat). The lock is a must-never-commit path, so acquisition ALSO runs `git ls-files -- do-work/orchestrator-lock.json` (cwd-relative from wherever the lock path resolves): a non-empty result means some earlier session committed the lock before this guard existed, every heartbeat is now dirtying a tracked file, and a committed lock propagates stale session state to every other clone. Report the remedy to the user — `git rm --cached do-work/orchestrator-lock.json` — never untrack on their behalf, and continue the run: a tracked lock is a hygiene failure, not a correctness one, and the Commit Phase's explicit-file staging keeps it out of this session's own commits.

**Stale threshold: 45 minutes.** Generous on purpose — a single dispatched agent (a Route C exploration, plan, or build attempt) can legitimately run tens of minutes between the dispatch-time and return-time refreshes below, and a false "stale" verdict on a genuinely live session is worse than a slightly slower self-heal on a genuinely dead one: the stale path takes over *without asking* and Crash Recovery then strips and re-queues whatever the "stale" session was actually mid-way through building. That is why the refresh schedule below is not just phase boundaries — the threshold only stays safe if no live session lets 45 minutes pass between lock writes. Applies identically to the holder's `heartbeat_at` and to every `coexisting_sessions[].heartbeat_at`.

**Interactivity test (concrete, not self-assessed).** A session is interactive only if BOTH hold — if either can't be confirmed, it is non-interactive:
1. **This is the top-level session the user is directly conversing with, not a subagent spawned via Task/Agent dispatch.** A subagent has no user turn to pause on — it returns its final message to whatever dispatched it, never to a human. This is the same distinction `actions/code-review.md` already draws for its own consent gate ("running non-interactively (e.g., via subagent)") and the harness-availability check in `actions/kb-lessons-handoff.md`, applied here to the identical failure mode: `actions/pipeline.md`'s foreground dispatch of the `run` step explicitly may run as a subagent ("subagent if available, inline otherwise" — `actions/pipeline.md` Step 4 substep 2), and a subagent that wrongly assumes it can prompt hangs waiting for an answer nobody can give.
2. **An ask-user / blocking-prompt mechanism is actually available and usable in this environment right now** — not merely documented or importable.

**Failure-safe default:** if either check is unconfirmed or ambiguous, treat this session as non-interactive and take the refusal branch below. A refusal is a clean, resumable no-op — re-run `do-work run` directly (not via subagent dispatch) to get the interactive prompt. A hung ask-user call blocks the entire pipeline with no way to resume, which is strictly worse than a refusal.

**Procedure — run once, the first time this session reaches Step 1** (session start, or resuming via CHECKPOINT.md); skip it on a Step 10 → Step 1 loop iteration once this session's own `session_id` already appears in the lock file (as holder or in `coexisting_sessions`) — a durable, file-backed fact, not something this session needs to remember on its own.

1. **No lock file exists** → acquire: generate `session_id` and write the fresh holder object (`started_at`/`heartbeat_at` = now, `claimed_reqs: []`, the derived `claimed_req: null`, `coexisting_sessions: []`) through the serialized mutex — **re-validating existence inside the critical section**: the branch was chosen on an unserialized existence check, so two sessions starting on an empty queue can both select it. If the fresh read inside the mutex finds a lock file that now exists, another session won the race — write nothing, release the mutex, and re-enter this procedure at step 2 with that file (its holder's heartbeat is seconds old, so the fresh-holder path applies). Only when the fresh read still finds no lock does the write proceed. Append the local-ignore snippet. Continue to the CHECKPOINT.md check and Crash Recovery normally.
2. **Lock file exists** → read it, compute `age = now − heartbeat_at` in minutes against the holder's `heartbeat_at` (same style as Step 1's Stale-reservation check comparing `reserved_at` against now).
   - **`age > 45` (holder stale):** warn `Lock held by <session_id> looks abandoned (heartbeat <age>m old; was building <all ids in claimed_reqs, or "nothing">). Treating as stale — taking over.` **Re-validate staleness inside the mutex before writing** — the stale verdict came from an unserialized read, and the holder may have refreshed in the gap (a session legitimately goes quiet for tens of minutes during a long dispatch, then heartbeats on return): from the fresh read inside the critical section, confirm the holder `session_id` is still the one just warned about and **recompute the holder's age from the fresh `heartbeat_at`**. If the holder changed or the recomputed age is ≤ 45 minutes, the holder is live — write nothing, release the mutex, and re-enter this procedure at step 2 with the fresh read (the fresh-holder path applies). Only when the fresh read still shows the same holder still past 45 minutes does the take-over proceed: rewrite only the top-level holder fields — freshly generated `session_id`, `started_at`/`heartbeat_at` = now, `claimed_reqs: []` and the derived `claimed_req: null`. **Preserve `coexisting_sessions`**, pruning only entries whose own `heartbeat_at` — as read inside the mutex, not from the earlier unserialized read — is also past 45 minutes: a live coexisting session keeps running, undisturbed by a holder take-over. Continue to the CHECKPOINT.md check and Crash Recovery normally — the per-file concurrency gate (`actions/work-reference.md` → **Crash Recovery (Step 1)**) decides which `working/` file this take-over actually reclaims.
   - **`age ≤ 45` (holder fresh — another orchestrator may be live):** apply the Interactivity test above.
     - **Interactive:** load `crew-members/clear-questions.md`, then present:

       ```
       Another do-work session appears active on this queue.
         Holder: <session_id>
         Started: <started_at>
         Heartbeat: <age>m ago (fresher than the 45m stale threshold)
         Currently building: <all ids in claimed_reqs, or "nothing claimed yet">
         Also coexisting: <N other session(s)> (omit this line if coexisting_sessions is empty)

       A second orchestrator on the same queue can double-claim or double-commit a REQ (this happened on 2026-07-01) — pick one:
         (a) Proceed anyway — run alongside it as a tracked coexisting session. I will scan for other pending REQs but will NOT touch whatever the holder (or another coexisting session) has actively claimed in do-work/working/. The holder is untouched — it keeps running exactly as if I weren't here.
         (b) Take over — assume the holder is gone even though its heartbeat looks fresh. I become the new holder and run normal crash recovery on its claimed do-work/working/ file (reset and re-queued); any other live coexisting session is left alone.
         (c) Abort — stop now. I will not write anything; the existing lock is left exactly as it was.
       ```

       - **(a) Proceed anyway** — append a new entry to `coexisting_sessions`: `{session_id: <new>, joined_at: now, heartbeat_at: now, claimed_reqs: [], claimed_req: null}`. Do **not** touch the holder's `session_id`/`started_at`/`heartbeat_at`/`claimed_reqs`/`claimed_req` — the holder's own heartbeat-refresh touchpoints keep matching, so it runs to completion undisturbed. Continue to the CHECKPOINT.md check and Crash Recovery normally: the per-file concurrency gate skips every file under a *fresh* claim — the holder's, another coexisting session's, and (once this session starts claiming) **including this session's own** co-dispatched siblings, because freshness alone gates (**Crash Recovery (Step 1)**, above) — and recovers every other file, so it is not an all-or-nothing skip.
       - **(b) Take over** — same rewrite as the stale path above (new holder `session_id`, `claimed_reqs: []` and the derived `claimed_req: null`; preserve/prune `coexisting_sessions`, prune ages from the in-mutex read). The in-mutex identity check still applies: the user authorized taking over the holder they were *shown* — if the fresh read inside the mutex shows a different holder `session_id`, write nothing, release the mutex, and re-enter this procedure at step 2 to re-present. (The age recheck does not apply — the user knowingly chose to take over a fresh holder.)
       - **(c) Abort** — make no writes at all; report the holder info above and stop before the CHECKPOINT.md check, Crash Recovery, or the queue scan.
     - **Non-interactive** (Interactivity test above resolved "no," including "couldn't confirm"): refuse.

       ```
       Refusing — another do-work session appears active on this queue and no interactive prompt is available.
         Holder: <session_id>
         Started: <started_at>
         Heartbeat: <age>m ago
         Currently building: <all ids in claimed_reqs, or "nothing claimed yet">
         Also coexisting: <N other session(s)> (omit this line if coexisting_sessions is empty)
       Re-run interactively to choose proceed/take-over/abort, or wait for the heartbeat to pass the 45-minute stale threshold.
       ```

       Then stop without creating, modifying, or deleting any file — lock, queue, or `working/`.

**Heartbeat refresh (every subsequent touch — not the guard above).** At Step 2 (claim), Step 6.25 (Implementation Summary), Step 8 (Archive), Step 9 (Commit Phase), and Step 10 (loop point): re-derive the lock path and read it. **These five are the floor, not the schedule.** Phase boundaries alone leave the whole triage → plan → build → review stretch dark — long enough for a live session to cross the 45-minute threshold and get taken over mid-build (the take-over then re-queues and strips the REQ it was building). So additionally refresh: **immediately before dispatching any long-running agent** (explore, plan, build, review — Steps 3–7) and **immediately when each one returns**; and at any natural pause the orchestrator passes through (between build attempts, test-fix iterations, review rounds, remediation loops) whenever **more than 15 minutes** have passed since this session's last lock write. Each refresh is the same serialized read-modify-write already specified — cheap, and no new state. A floor-capability session that must block synchronously on one agent for over ~40 minutes cannot refresh mid-wait; prefer splitting such a dispatch, and accept that the residual risk is confined to that single window rather than to every phase.
- If this session's `session_id` matches the top-level holder: rewrite the holder's `heartbeat_at` = now and recompute `claimed_reqs` = **the ids in this session's own dispatch record** (**Dispatch bookkeeping is the recompute's source**, below) — never a bare listing of every REQ present in `do-work/working/`, which under co-dispatch would sweep a coexisting session's file into this holder's claim (the field-stealing the serialized-mutex rule above forbids), and never the lock's own previous `claimed_reqs`. Re-derive the legacy `claimed_req` mirror as `claimed_reqs[0]` (or `null` when the recompute is empty).
- Else if this session's `session_id` matches an entry in `coexisting_sessions`: rewrite *that entry's* `heartbeat_at` = now and recompute *that entry's* `claimed_reqs` the same way — from this session's own dispatch record, never a listing of `working/` and never the entry's previous `claimed_reqs` — then re-derive that entry's `claimed_req` mirror. Leave the holder and every other entry untouched.
- Else (this session's `session_id` appears nowhere in the lock file): it was superseded — another session's take-over overwrote the holder slot, or this session's own coexisting entry was pruned as stale. Stop and report rather than silently re-adding yourself — this is what keeps "file-based" from turning into "two sessions both believing they hold it."

**Dispatch bookkeeping is the recompute's source.** "Under this session's own claim" is a fact this session records as it dispatches — not something re-read out of the lock, and not something inferrable from a `working/` listing (a listing cannot tell whose claim a file is under, which is the whole reason the bullets above forbid it). An id **enters** the record when Step 2 appends it to `claimed_reqs` (`actions/work.md` Step 2 substep 1). That append happens *before* Step 2 moves the file into `working/`, so at Step 2's own refresh the freshly appended id is already part of the set even though nothing in `working/` bears it yet — a recompute that consulted the directory there would erase the claim the same write just made and reopen the abandoned-crash-artifact window the claim-before-move ordering exists to close. An id **leaves** the record on exactly two events:

- **Step 8 released it** — the archive, blocked, or failed path removed this id from `claimed_reqs` right after moving the file out of `working/` (`actions/work.md` Step 8); the record drops it in the same beat, so the next recompute cannot re-assert a claim on a REQ that has already left. The ordinary exit.
- **Its builder is known-dead with no Step 8 left to run** — the dispatch returned aborted or errored, or the harness reports the agent gone, so no orchestrator step will ever archive or release that REQ. Drop the id. The next refresh then omits it and re-derives the shortened `claimed_req` mirror, the REQ's still-present `working/` file stops being exempt from the per-file gate, and this session's own next Crash Recovery pass (Step 10 → Step 1) resets and re-queues it. Without this, the claim is re-asserted on every heartbeat and the REQ is unreclaimable until 45 minutes after the entire session ends. **An ordinary failed build is not a dead builder:** an attempt that *returns*, however badly, keeps its claim through retries, review remediation, and Step 8's failure classification, which owns the removal on that path — dropping a claim there would invite another live session to strip a REQ mid-remediation.

Sourcing the recompute from the lock's previous `claimed_reqs` is what would make a dead builder's claim self-reinforcing: the per-file gate skips whatever `claimed_reqs` names, so a recompute reading `claimed_reqs` keeps naming it.

While performing any of the checks above, opportunistically prune `coexisting_sessions` entries (other than the one just refreshed) whose `heartbeat_at` is past the 45-minute threshold — this is what lets an abandoned coexisting session's claim self-heal via the Crash Recovery per-file gate on a later Step 1 entry, with no separate cleanup pass.

**Release last — after every repo mutation this session will make, including `actions/cleanup.md` and its commit.** Releasing is what lets the next `do-work run` start claiming REQs and moving files, so anything this session does afterward races that new session over the same working tree. Cleanup is the trap: its Pass 0 sweeps `do-work/queue/` and `do-work/working/` and stages those paths (`actions/cleanup.md`), which is exactly where an arriving session has just put its freshly claimed REQ — release-then-cleanup lets a departing session sweep and commit the newcomer's in-flight work. The lock costs nothing to hold through a cleanup pass; only Step 10's exit runs cleanup, so this ordering matters there in particular. Holding the lock fences off *arriving* sessions only — a session already coexisting keeps running regardless, which is why cleanup's Pass 0 carries its own live-claim gate (`actions/cleanup.md`) exempting every freshly claimed id in any live session's `claimed_reqs` from its sweeps.

**Release (every normal exit — Step 1's composed no-pending exit, Step 1's targeted-mode stop, Step 10's exit, and an unrecoverable-error stop):** re-derive the path and read it.
- **This session is the holder:** prune any `coexisting_sessions` entries past the 45-minute threshold first. If none remain, delete the whole lock file. **If any live entry remains, do NOT delete or modify the file** — leave it exactly as-is. This session's heartbeat simply stops refreshing and ages past 45 minutes on its own; the next session to reach Step 1 takes the stale-holder path above (warn, take over, preserving/pruning `coexisting_sessions`) without disrupting whichever coexisting session(s) are still live. Trades a bounded (≤45m) delay in freeing the holder slot for not inventing a hand-off mechanism — YAGNI, not a distributed lock.
- **This session is a coexisting session** (its `session_id` matches an entry, not the holder): remove only that entry from `coexisting_sessions`. Never delete the whole file or touch the holder's fields — the holder or other coexisting sessions may still be active.
- **This session's `session_id` matches neither:** already superseded — leave the file completely alone, it's no longer yours to release.

## Composed Exit Summary (Step 1)

**Exit paths when no `pending` REQs found:**

The exit report is **composed**, not picked from disjoint branches. Whenever the scan finds no dependency-ready `pending` REQ, lead with the headline that matches the actual queue state — `No pending REQs in queue.` when the queue holds no `pending` REQs at all, or `No dependency-ready pending REQs.` when `pending` REQs exist but every one is dependency-blocked (the blocked-by-dependencies section below then enumerates them, so the headline never strands the user). Then append every section that has at least one REQ. Six sections may apply, in this order:

1. **Completed/done section** — applies if any REQ in `do-work/queue/` has status `completed`, `completed-with-issues`, `cancelled`, or `done`. Read the `user_request` frontmatter field from each to group by UR. Render:

   ```
   ⚠ N finished REQs awaiting archive (UR-137: 3 REQs, UR-138: 1 REQ, ...):
     REQ-351 — [title] (done)
     REQ-352 — [title] (completed)
     REQ-353 — [title] (cancelled)
     ...

   Run `do-work cleanup` to archive completed work, then `do-work recap` to see full history.
   ```

2. **Pending-answers section** — applies if any REQ has status `pending-answers`. Render from frontmatter only — do not open the REQ body to count `## Open Questions` items at this stage (Step 1 reads frontmatter per the queue scan). The count is deferred to `do-work clarify`, which is the action that reads Open Questions sections:

   ```
   ⚠ N REQs awaiting clarification:
     REQ-NNN — [title] (pending-answers)
     ...

   Run `do-work clarify` to batch-review the open questions; resolved REQs flip to `pending` and re-enter the queue.
   ```

3. **Blocked-on-external-condition section** — applies if any REQ has status `blocked` (waiting on an external condition named in `blocked_by` — a service being up, a person answering, credentials provisioned — not user answers and not another REQ). Render from frontmatter: the `blocked_by` condition, the age from `blocked_at` (now − `blocked_at`), and whether an auto-probe is configured or failed this run. Step 1 already re-ran each `blocked_check` probe before composing this summary, so a REQ that still appears here either has no probe or its probe did not pass this run:

   ```
   ⚠ N REQs blocked on external conditions:
     REQ-NNN — [title] (blocked by: <condition>, since <age>) [probe failed this run | no auto-probe]
     ...

   When a condition is satisfied, re-run `do-work run` (REQs with a `blocked_check` are re-probed automatically and unblock on exit 0),
   or confirm a human-checkable one via `do-work clarify`. To give up on one, `do-work abandon REQ-NNN`.
   ```

4. **Blocked-archive-collision section** — applies if any REQ has status `blocked-archive-collision`. Read the matching archive path from each blocked REQ's frontmatter if recorded; otherwise re-run the Step 2.0 glob (`do-work/archive/**/REQ-NNN-*.md` and `do-work/archive/**/REQ-NNN.md`) to find it. Render:

   ```
   ⚠ N REQs held by archive-collision guard:
     REQ-NNN — [title] (queue file: do-work/queue/REQ-NNN-slug.md)
       already archived at <archive-path>
       recover: rename the queue file (if this is an intentional re-do) or delete it (if it's a stale duplicate), then flip status back to `pending`
     ...
   ```

5. **Blocked-by-dependencies section** — applies if any `pending` REQ has an unmet `depends_on` reference (dependency-blocked) or any REQ has status `blocked-dependency-cycle`. Pending REQs stay `pending` (the gating is dynamic — they become ready as upstream REQs complete); only cycle-detected REQs are flipped to a held status. Render both groups under one heading:

   ```
   ⚠ N REQs blocked by unmet dependencies:
     REQ-NNN — [title] (pending; depends on REQ-MMM, status: <pending|claimed|reserved|pending-answers|failed|cancelled>)
     REQ-PPP — [title] (blocked-dependency-cycle; chain: REQ-PPP → REQ-QQQ → REQ-PPP)
     ...

   Resolve the blocking REQs first, then re-run. To force a scoped run that ignores dependency gating for a specific REQ, use `do-work run REQ-NNN`. To break a dependency cycle, edit the REQ's `depends_on` and flip its status back to `pending`. A dependency on a `cancelled` (or `failed`) REQ never self-resolves — re-point the dependent's `depends_on`, or abandon it too (`do-work abandon REQ-NNN`).
   ```

6. **Reserved section** — applies if any REQ has status `reserved` (allocated to another worktree/cloud session via `do-work reserve`, `actions/reserve.md`). Render each with its `reserved_for` label and age (now − `reserved_at`); a reservation older than **24 hours** is stale and gets the recategorize suggestion:

   ```
   N REQs reserved for other sessions:
     REQ-NNN — [title] (reserved for: <label>, <age> ago)
     REQ-MMM — [title] (reserved for: <label>, <age> ago) ⚠ STALE
     ...

   ⚠ Reservations older than 24h may belong to dead sessions. Recategorize each: `do-work release REQ-MMM`
   to return it to the queue, `do-work run REQ-MMM` to claim it in this session, or leave it if the owning
   session is still active.
   ```

**After rendering all applicable sections, exit the work loop** — do not proceed to Step 2.0 or beyond. There is no `pending` REQ to claim. Step 1's contract on the no-pending path is "render the composed summary, then stop"; the only path that continues is the one where Step 1 finds at least one dependency-ready `pending` REQ.

If **no section applies** (no REQs at all in `do-work/queue/`), report completion and exit. Never silently exit when any of the six sections applies — every non-pending or non-ready REQ in the queue is something the user needs to see.

**Composition is deliberate.** A queue with both `pending-answers` and `blocked-archive-collision` REQs (and no completed/done) renders both sections back-to-back. A queue with all six categories renders all six. The user sees the full picture in one report instead of a single branch's slice.

## Triage Section Template (Step 3)

```markdown
---

## Triage

**Route: [A/B/C]** - [Simple/Medium/Complex]

**Reasoning:** [1-2 sentences]

**Planning:** [Required/Not required]
```

## Plan Template — Route C (Step 4)

```markdown
## Plan

[Plan agent output]

*Generated by Plan agent*
```

## Plan Skip Note — Routes A/B (Step 4)

```markdown
## Plan

**Planning not required** - [Route A: Direct implementation / Route B: Exploration-guided implementation]

*Skipped by work action*
```

## Scope Declaration Template (Step 5.5)

```markdown
## Scope

**Files I will touch:**
- `src/stores/theme-store.ts` (new) — theme state management
- `src/components/settings/SettingsPanel.tsx` (modify) — add toggle
- `tests/theme-store.test.js` (new) — unit tests

**Files I will NOT touch:** [any files that seem related but are out of scope]

**Acceptance criteria (restated from REQ):**
- [ ] Dark mode toggle visible in settings
- [ ] Theme persists across page reload
- [ ] OS preference respected on first visit
```

**"Files I will touch" is the source of the `write_set` frontmatter field.** After writing this section, the orchestrator mirrors the list into `write_set:` — one direction only, so the prose and the field cannot drift. Never edit `write_set` and expect the Scope list to follow.

**A dispatch-time partition directive survives the mirror (parallel dispatch only).** If a builder was handed a narrowed partition — a subset of a shared file set — at dispatch (`actions/work.md` Step 1 parallel-dispatch gate), that narrowed subset *is* its "Files I will touch," so mirroring it writes the partition, not the full shared list. The mirror never widens a partitioned set back to the un-partitioned one. The gate also wrote that subset into the REQ's `write_set` frontmatter as it dispatched — the one write to the field that does not originate in a `## Scope` section, and narrowing-only — which is what lets a sibling's Step 5.5 re-check see the partition rather than flag the orchestrator's own directive as an overlap. So the mirror here confirms a field that already holds the partition. A serial run has no partition and no other in-flight REQ, so this never applies — the plain one-directional mirror above is the whole story.

## Pre-Flight Template (Step 5.75)

```markdown
## Pre-Flight

**Git:** ⚠ 3 uncommitted files (src/temp.ts, .env.local, notes.md)
**Tests baseline:** ✓ All passing (47 tests)
**Dependencies:** ✓ Installed

*Checked by work action*
```

## Implementation Summary Template (Step 6.25)

```markdown
## Implementation Summary

**Files changed:**
- `src/stores/theme-store.ts` (new)
- `src/components/settings/SettingsPanel.tsx` (modified)
- `tests/theme-store.test.js` (new)

**What was done:** [1-2 sentences — what the implementation actually did]
```

## Qualification Anti-Rationalization Table (Step 6.3)

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "The summary says files changed" | Check the file system | The summary is a claim, not evidence |
| "Tests pass so requirements are met" | Compare requirements to diff, word by word | Tests can be incomplete |
| "The builder checked the UNIFY box" | Read the actual diff for debug artifacts | A checked box is a claim, not a fact |
| "This works on my test case" | Test at least 2 additional cases including an edge case | One test case proves nothing about generality |
| "The existing code was already like this" | Flag it in Discovered Tasks | Pre-existing problems are still problems |
| "It's just a small deviation from the plan" | Log it as a Decision (D-XX) | Unlogged deviations break traceability |

## Testing Section Template (Step 6.5)

```markdown
## Testing

**Tests run:** [command]
**Result:** ✓ All passing (X tests)

**Red-green validation:** *(for bug fixes and new features)*
- [test name/file]: ✗ before implementation → ✓ after
- [test name/file]: ✗ before implementation → ✓ after

**New tests added:**
- [list]

**Existing tests updated (cross-REQ impact):**
- [test file] (from REQ-NNN): [what changed and why — intentional behavior change]

*Verified by work action*
```

## Deferred Prime-Link Path Computation (Step 7.5)

**Path computation rule (for use in Step 8):** the link path must be relative to the prime file's location, not the repo root. Count how many directories deep the prime file sits (i.e., the number of path components before the filename). Prepend that many `../` steps to the REQ's repo-root-relative archive path. Examples:
- Prime at `prime-auth.md` (0 dirs deep) → `do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`
- Prime at `src/utils/prime-auth.md` (2 dirs deep: `src/` and `utils/`) → `../../do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`
- Prime at `web/src/auth/prime-auth.md` (3 dirs deep) → `../../../do-work/archive/UR-005/REQ-042-auth-fix.md#lessons-learned`

The existence-verify check on the resolved path runs in Step 8 (post-move) — that's the whole reason for deferring.

## Builder-Decided Follow-up Template (Step 8)

   ```markdown
   ---
   id: REQ-NNN
   title: "Confirm: [brief description of the choice]"
   status: pending-answers
   created_at: [timestamp]
   user_request: [same UR as the original REQ]
   addendum_to: [original REQ id]
   builder_decided: true
   ---

   # Confirm: [Brief Description]

   ## What
   During [REQ-id], the builder chose [choice] for [question]. This follow-up
   confirms whether that choice matches your intent or if you'd prefer a different approach.

   ## What the Builder Chose
   [Description of the choice and its impact on the implementation]

   ## What Would Change
   [If the user picks a different option, what would need to change]

   ## Open Questions
   - [ ] [Original question]
     Recommended: [builder's choice — already implemented]
     Value: [what this choice buys — copied from the D-NN record]
     Risk: [what breaks if it's wrong, and how reversible — copied from the D-NN record]
     Also: [other alternatives]
   ```

   The `Value:`/`Risk:` lines come from the escalated `D-NN` entry's record (work.md Step 3.5/6). They let `do-work clarify` render the **DECISIONS FOR YOU** section of the Decision Brief so the user can judge in seconds. If the original decision was logged without them (older REQ), omit both lines — `clarify`'s fallback renders `Recommended:`/`Also:` alone.

   All template text the user will read — the question, `Recommended:`/`Also:`/`Value:`/`Risk:`, and the `## What` section — must satisfy `crew-members/clear-questions.md`: self-contained, no spec-internal shorthand or coined labels without a one-line gloss, and the why-this-was-escalated stated in `## What` or the question itself (Principle 7). The builder writing this template has the full spec in context; the user answering it in a later clarify session does not.

## Discovered Tasks Classification (Step 8)

   The builder should classify each discovered task when appending them:
   ```
   ## Discovered Tasks
   - **[critical]** SQL injection vulnerability in user search endpoint
   - **[normal]** Dead code in utils/legacy-parser.ts can be removed
   - **[low]** Variable naming inconsistency in auth module
   ```

   If the builder did not classify them, the orchestrator classifies based on:
   - **critical**: Security vulnerability, data loss risk, broken functionality in production paths
   - **normal**: Technical debt, missing tests, minor bugs in non-critical paths
   - **low**: Style issues, naming, dead code, documentation gaps

   **For `[critical]` discoveries:** Create follow-up REQ with `status: pending` (not `pending-answers`) — these skip user confirmation and go straight into the work queue. Add a note in Open Questions: `- [x] Auto-approved: critical severity (security/data/production risk). → Added to queue immediately.` Report prominently: `⚠ CRITICAL discovered: [description] — auto-queued as REQ-NNN`

   **Test-hygiene carve-out:** A `[normal]` or `[low]` discovery ALSO auto-queues with `status: pending` (same as critical) when ALL three hold:
   - The fix touches **only test files** (`tests/**`, `*.test.*`, `*.spec.*`, test helpers/fixtures) — zero production-source changes.
   - It is **mechanical hygiene** — silencing warnings/console noise, deflaking, lint/format cleanup in tests — with no behavior or assertion-meaning changes.
   - It is **small** — a single file or a couple of files, no new infrastructure.

   Failing ANY bullet keeps the `pending-answers` flow below. The paper trail mirrors the critical flow: add a note in Open Questions: `- [x] Auto-approved: test-only mechanical hygiene ([severity]). → Added to queue.` Report visibly: `↺ test-hygiene discovery auto-queued as REQ-NNN`. The user can still discard the REQ from the queue afterwards — that stays the escape hatch.

   **For `[normal]` and `[low]` discoveries (when the test-hygiene carve-out does not apply):** Use the existing `pending-answers` flow:
   - Set frontmatter: `status: pending-answers`, `user_request: [same UR]`, `addendum_to: [current REQ id]`, `domain: [same domain as current REQ]`.
   - Add an `## Open Questions` section with this checkbox format:
     `- [ ] I discovered this out-of-scope task while working on [current REQ]: [Task Description]. Should I process this as a new task?`
     `  Recommended: Yes, add to queue (will flip to 'pending').`
     `  Also: No, discard it.`
   This ensures non-critical discoveries — other than qualifying test-only hygiene — require the user's explicit permission via `do-work clarify` before execution.

## Failure Classification (Step 8)


Before classifying via the symptom table below, **check for upstream failure**. Cascades from a failed prerequisite often present as plausible-looking `code` or `spec` symptoms in the downstream REQ; misclassifying them sends the builder chasing phantom bugs in the wrong domain.

**Upstream-failure short-circuit:**

Read the frontmatter of every REQ this one depends on:
- `addendum_to` (single parent, if set; or `amends`/`parent`/`amendment_to` as the legacy alias if `addendum_to` is absent — same back-compat shape as the `depends_on`/`dependencies:` pair; `addendum_to` wins when multiple are present)
- every entry in `depends_on` (if set, or every entry in the legacy `dependencies:` alias if `depends_on` is absent — same back-compat rule as Step 1; `depends_on` wins when both present)

Resolve each ID by globbing `do-work/archive/**/REQ-NNN-*.md`, `do-work/archive/**/REQ-NNN.md`, `do-work/queue/REQ-NNN-*.md`, and `do-work/working/REQ-NNN-*.md`. If any referenced REQ has `status: failed`, skip the symptom table and short-circuit classification:

- `status: failed`
- `error_type: spec` (the local approach is downstream-correct only if the upstream is correct; with the upstream broken, the local spec is implicitly unsound)
- `error: "Upstream REQ-NNN failed (error_type: <ancestor.error_type>); downstream blocked. Original error: <original error message>"`

Create the follow-up REQ per the Spec row below. It inherits `addendum_to: <this failed REQ>`; the cascade is now visible in the addendum chain and the follow-up's error description names the upstream root cause. The follow-up should also carry the original dependency list so it re-blocks on the same upstream until the upstream's own follow-up lands — and it always emits the canonical `depends_on:` key, even if the failed REQ used the legacy `dependencies:` alias. Don't propagate the alias on follow-ups.

If no upstream REQ is `failed`, fall through to the symptom-based classification table:

| Type | Symptoms | Recovery |
|------|----------|----------|
| **Intent** | Requirements are ambiguous or contradictory; builder couldn't determine what to build | Create a follow-up REQ with `status: pending-answers` containing the specific ambiguities as Open Questions. Archive original as `failed` with `error_type: intent`. |
| **Spec** | Requirements are clear but the technical approach was wrong (wrong files, wrong pattern, wrong architecture) | Create a follow-up REQ with a `## Prior Attempt` section summarizing what was tried and why it failed. Set `status: pending`. Archive original with `error_type: spec`. |
| **Code** | Approach was right but implementation has bugs (tests fail, runtime errors, logic errors) | Create a follow-up REQ targeting the specific code issue. Set `status: pending`. Archive original with `error_type: code`. |
| **Environment** | External dependency unavailable, permissions issue, tooling broken | **First apply the blocked-flip test** (see `actions/work.md`'s mid-run blocked-flip procedure): if *no substantive implementation edits landed this attempt* AND the missing thing is a precondition expected to become available on its own (a service comes up, a person answers, credentials get provisioned), do **not** fail — flip the REQ to `status: blocked` with `blocked_by` + `blocked_at` (non-terminal, stays in the queue) instead. Reserve `error_type: environment` for post-work breakage, a broken/permission-denied environment the user must repair, or a precondition that will not self-resolve. Then: no follow-up REQ — archive with `error_type: environment` and a clear description of what's needed. |

**Anti-rationalization addition.** When checking the symptom table:

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "This REQ failed on a code bug" | Check whether any `addendum_to` or `depends_on` ancestor (or `dependencies:` alias) is also `failed` first | Downstream failures often inherit upstream rot; misclassifying as `code` chases phantom bugs in the wrong domain |

**Procedure:**
1. Run the upstream-failure short-circuit. If it fires, jump to step 3.
2. Otherwise classify using the symptom table above.
3. Update frontmatter: `status: failed`, `error: "description"`, `error_type: [intent|spec|code|environment]`
4. For Intent/Spec/Code failures: create the appropriate follow-up REQ (details above). Set `addendum_to` to the failed REQ's ID so context chains. Preserve the original dependency list on the follow-up when the failure was upstream-driven — always emit it under the canonical `depends_on:` key, even if the failed REQ used the legacy `dependencies:` alias.
5. Move to `archive/` (failed REQs always go to archive root, not into UR folders).
6. Report to user: `[REQ-NNN] failed ([type]): [description]. Follow-up: [REQ-NNN] / None.` When the short-circuit fired, prefix the report with `(upstream cascade — original failure at REQ-NNN)`.

## Changelog Entry Procedure (Step 9)

Every successful REQ (`completed` / `completed-with-issues`) gets an entry in the target repo's root `CHANGELOG.md`, written **before** the commit so it ships inside it. Failed and cancelled REQs get no entry — the changelog records delivered change, not attempts. A changelog entry is a human-facing artifact: load `crew-members/anti-slop.md` before writing it (its JIT_CONTEXT condition already covers this — noted here so it isn't skipped).

**Precedence check first.** If the repo already has a `CHANGELOG.md` whose entries follow a different convention (keep-a-changelog categories, generated conventional-commit logs, plain dated lists), **match the existing format** — never impose the house voice on a repo with its own. Everything below applies when there is no changelog yet or the existing one already follows this format.

**Bootstrap.** If no root `CHANGELOG.md` exists, create one:

```markdown
# Changelog

What's new, what's better, what's different. Most recent stuff on top.

---
```

**Entry key.** Always `## X.Y.Z — [Short Descriptive Title] (YYYY-MM-DD)` — every entry carries both a version and a date. The title must say what was delivered so a reader scanning only headings knows what changed ("Board View Filters", not a whimsical codename). It must be unique against every existing entry in the file (grep before writing — duplicates have occurred), and the new `X.Y.Z` must be **strictly greater** than the version in the file's first existing entry (duplicate version numbers have occurred).

**Where `X.Y.Z` comes from.** Resolve the version source once per entry, in this order:

1. **A version in a repo file** — `package.json`'s `"version"`, `Cargo.toml`, `pyproject.toml`, a `VERSION` file, or the like (this list is illustrative; any file the project maintains a version line in qualifies). Bump that line by the rule below, write the bumped value into the file, and **stage the file with the REQ's commit**. The repo's version and the changelog header stay in lock-step.
2. **Version only in release tags** (no version line in any file) — take the highest release tag as the current version and bump from it, but write the result **only into the changelog header**. do-work never creates a git tag: a tag is a release announcement, and only a human decides when one happens.
3. **No version anywhere** — the changelog is the source of truth. Take the highest `X.Y.Z` across the file's existing entries and bump from it. If there are no entries yet (bootstrap), seed the first entry at `0.1.0`. Nothing outside `CHANGELOG.md` is touched — an unversioned repo stays unversioned, and the header number is a changelog fact, not a claim that a release was cut.

Two guards on source resolution. If **two or more version files disagree** with each other, do not guess which one the release process uses: leave every file untouched, fall back to the changelog counter (source 3), and say so in the Step 9 report. If the resolved source is **behind** the newest changelog entry (someone released or edited out of band), bump from whichever is higher — never emit a version below one already in the file.

**Bump size.** Read the change the REQ actually delivered, not its wording:

| Bump      | When                                                                                                                                       |
| --------- | ------------------------------------------------------------------------------------------------------------------------------------------ |
| **major** | An existing consumer breaks: a public API or CLI flag removed or renamed, an on-disk or wire format changed, a documented default reversed. |
| **minor** | A user-invocable capability exists that didn't before, and nothing existing breaks.                                                          |
| **patch** | Everything else — bug fixes, performance, refactors, tests, docs, internal-only changes.                                                    |

Tie-breakers, in order: a breaking change outranks an additive one in the same REQ (bump major, not minor); when genuinely torn between two levels, pick the **smaller** one. **Below `1.0.0`, a breaking change bumps the minor, not the major** — `0.x` is unstable by semver's own definition, so the first breaking change in a seeded repo must not silently promote it to a `1.0.0` release. A `completed-with-issues` REQ is bumped on what it delivered, exactly like `completed`.

The `CHANGELOG.md` change — and the version file, when source 1 applied — are part of the REQ's lifecycle files. Stage them in the commit below.

**Voice contract (house style).** 1–2 casual sentences leading with *why it matters* — the situation that prompted the change and what's better now — then bullets for the specifics. Lead with value, not implementation; file paths and flags belong in the bullets, not the lead. Keep it brief. Newest on top, one entry per REQ (this matches one-commit-per-request).

```markdown
## 0.4.0 — Clear Questions for Interactive Prompts (2026-07-07)

Agents kept asking questions only they could parse — codenames coined
mid-analysis, options with no stated consequence. Question wording is
now a contract, not a hope.

- New `crew-members/clear-questions.md`, loaded before any interactive ask
- Six principles: one decision per question, decode your own shorthand, say the consequence…
```

## Commit & Metadata-Commit Procedure (Step 9)

```bash
# Stage implementation files + archived REQ + the changelog entry
git add src/stores/theme-store.ts src/components/settings/SettingsPanel.tsx \
  do-work/archive/UR-002/REQ-003-dark-mode.md CHANGELOG.md

# Stage the bumped version file — only when the changelog resolved to a repo
# version file (source 1). Tag-versioned and unversioned repos have none.
git add package.json

# Stage follow-up REQs created in Step 8 (if any)
git add do-work/queue/REQ-025-confirm-sidebar-palette.md

# Stage UR-folder move (if this was the last REQ and the UR moved to archive/)
# Both the old path (deletion) and new path (addition) must be staged.
# Exception: if the UR was never committed (capture's commit step was skipped,
# or the repo wasn't git at capture time), the old path matches nothing and
# `git add` exits 128 — stage only the new archive path in that case.
git add do-work/user-requests/UR-002/ do-work/archive/UR-002/

git commit -m "$(cat <<'EOF'
[REQ-003] Dark Mode (Route C)

Implements: do-work/archive/UR-002/REQ-003-dark-mode.md

- Created src/stores/theme-store.ts
- Modified src/components/settings/SettingsPanel.tsx

EOF
)"
```

**Format:** `[{id}] {title} (Route {route})` + `Implements:` line + summary bullets. Add a co-author trailer if your platform convention calls for one (e.g., `Co-Authored-By: Agent <agent@example.com>`), otherwise omit.

One commit per request. Stage all files created, modified, moved, or deleted during this request's lifecycle: implementation files (listed in the Implementation Summary), the archived REQ file, the `CHANGELOG.md` entry and the version file it bumped, if any (successful REQs only — see the Changelog Entry Procedure above), any follow-up REQs created in Step 8 (`pending-answers` files in `do-work/queue/`), and any UR-folder moves to `archive/`. If Step 8 substep 7 wrote prime-file lessons links, the modified prime files must also be staged — they are part of the REQ's lifecycle changes even though they aren't listed in the Implementation Summary's `Files changed`. Do not use `git add -A` or `git add .`, and never bypass a commit hook (see `actions/commit.md` § Rules for the full guard). Failed requests get committed too.

**In worktree dispatch mode** the builder's implementation is already committed and merged (Step 6's `--no-ff` merge), so do **not** stage implementation files here — stage only the archived REQ, the `CHANGELOG.md` entry and any bumped version file, follow-up REQs, UR-folder moves, and prime-file lessons links. The `commit:` field gets the `--no-ff` merge commit's hash (`<merge_hash>`, captured in Step 6 — the **latest** merge if remediation re-merged; **Worktree Dispatch Mode (Step 1)** above), not this changelog commit's hash.

**Validation check (successful REQs only):** Before committing, compare the `## Implementation Summary` file list against the staged files (excluding `do-work/` paths). If the Implementation Summary lists files that aren't staged, or if the only staged files are `do-work/` metadata, `CHANGELOG.md`, and/or the version file it bumped (the changelog entry and the version bump describe the implementation, they aren't the implementation), flag the mismatch — the commit may not contain the actual implementation. Fix the staging or update the Implementation Summary before proceeding. Design-artifact files placed outside `do-work/` satisfy this check — they are project deliverables. **Skip this check for failed REQs** — they may have no Implementation Summary or no project files staged, and that's expected. **In worktree dispatch mode** the implementation files live in the merge commit, not this commit's stage, so validate the `## Implementation Summary` file list against `git diff --name-only <pre>..<merge_hash>` (the merge range, excluding `do-work/` paths) instead of the staged set — a stage of only the changelog/version/`do-work/` metadata is correct here, not a mismatch.

**Write commit hash back to the archived REQ.** After the commit succeeds, resolve the implementation hash — **serially** it is the commit you just made, so read it with `git rev-parse --short HEAD`; **in worktree dispatch mode do NOT rev-parse `HEAD` here**, because HEAD is the changelog commit and the implementation lives in Step 6's `--no-ff` merge — use the `<merge_hash>` literal held since Step 6 (the latest merge, if remediation re-merged). Hand that hash to the shipped guard script, which makes the one-line edit and verifies it, then create a **separate metadata commit** (never amend — amending changes the hash and invalidates what was just written):

```bash
# Serial mode — the hash is resolved and consumed inside the SAME command block, so nothing
# has to survive between blocks and nothing is re-typed.
<skill-root>/tools/checks/record-commit-hash.sh do-work/archive/UR-NNN/REQ-NNN-slug.md "$(git rev-parse --short HEAD)"

# Worktree dispatch mode — pass the <merge_hash> literal held since Step 6, NOT
# `git rev-parse --short HEAD`, which here names the changelog commit you just made.
<skill-root>/tools/checks/record-commit-hash.sh do-work/archive/UR-NNN/REQ-NNN-slug.md <merge_hash>

# Then stage and commit ONLY that file — the script prints these two lines back to you.
git add -- do-work/archive/UR-NNN/REQ-NNN-slug.md
git commit -m "[REQ-NNN] record commit hash <hash>" -- do-work/archive/UR-NNN/REQ-NNN-slug.md

# Finally, confirm the commit landed what was verified. This is the only check that catches a
# content-mutating pre-commit hook rewriting the file after every pre-commit guard passed.
<skill-root>/tools/checks/record-commit-hash.sh --verify do-work/archive/UR-NNN/REQ-NNN-slug.md <hash>
```

The script edits only the `commit:` line **inside the frontmatter block** (adding it after `completed_at:` when absent), and refuses to write at all unless the rewrite changed exactly that one line. **Read its output before staging anything:**

- `OK: …` — the field records the hash; stage and commit as printed.
- `NOOP: …` — the field already records this hash **and** that is already committed. Make **no** metadata commit.
- `FAIL: …` (exit 1) — **stop.** The REQ file is left as it was found. Do not retry the command, do not hand-edit around it, and do not commit anyway; the message names the specific defect and the recovery command.
- Exit 2 is a usage error — a bad hash shape (including passing the literal `<hash>`), a missing file, or a file that is not a usable REQ.

**Why a script and not prose:** this write-back was free-form until a consumer repo's metadata commits truncated six archived REQ files to 0 bytes — 9 KB to 26 KB of decision trail each — with commit messages that claimed success. The guards run **before** `git add`, so a mass-deletion diff can never become a commit. A hand edit here has already destroyed archived REQ content; the script is the sanctioned path.

**If the script is missing** (a consumer on an older tarball), do the edit by hand and run its two cheapest guards before committing, in this order: `git diff --numstat HEAD -- <req-file>` must read `1	1` (or `1	0` when the field was added) — a deletion count above that means content was destroyed, so **stop and recover instead of committing** — and `test -s <req-file>` must pass. Then stage and commit as above, re-typing the hash as a literal: shell state does not survive between command blocks.

This ensures the `commit:` field in the archived REQ contains the real implementation commit hash — the commit just made serially, the `--no-ff` merge commit in worktree dispatch mode — which the review-work and present-work actions depend on for traceability. The metadata commit is a lightweight bookkeeping entry — it does not contain implementation changes.

## Session Checkpoint Template (Step 10)

```markdown
---
session_ended: [timestamp]
last_completed: REQ-NNN
queue_state: [N pending, N pending-answers, N blocked, N blocked-archive-collision, N in-progress]
reqs_processed_this_session: N
session_depth: light | moderate | heavy
---

# Session Checkpoint

## Completed This Session
- REQ-NNN: [title] (Route [X], [score]%)
- REQ-NNN: [title] (Route [X], [score]%)

## In Progress (interrupted)
- REQ-NNN: [title] — stopped at [phase: triage/planning/exploring/implementing/testing/reviewing]
  Last known state: [1-2 sentences]
  Key files being modified: [list]
  Known issues: [any blockers or concerns]

## Still Queued
- REQ-NNN: [title] (pending)
- REQ-NNN: [title] (pending-answers — [N] questions)

## Session Notes
[Environment issues, user preferences expressed, patterns discovered, decisions made outside REQ files]

## Context Summary (heavy sessions only)
[Recap of key decisions (D-XX references), architectural patterns encountered, and prime files
that the next session should re-read before starting. Include this section when 6+ REQs were
processed — at that volume, carried-over assumptions are unreliable.]
```

## Progress Reporting Example

```
Processing REQ-003-dark-mode.md...
  Triage: Complex (Route C)
  Open Questions: 2 found → builder decided (follow-ups queued)
  Planning...     [done]
  Scope...        [done] 4 files declared
  Exploring...    [done]
  Implementing... [done]
  Summary...      [done] 3 files changed
  Qualifying...   [done] ✓ files verified, requirements traced
  Testing...      [done] ✓ 12 tests passing
  Reviewing...    [done] 92% — 0 follow-ups
  Archiving...    [done]
  Committing...   [done] → abc1234

Processing REQ-004-fix-typo.md...
  Triage: Simple (Route A)
  Implementing... [done]
  Summary...      [done] 1 file changed
  Qualifying...   [done] ✓ verified
  Testing...      [done] ✓ 3 tests passing
  Reviewing...    [done] 88% — 0 follow-ups
  Archiving...    [done]
  Committing...   [done] → def5678

All 2 requests completed:
  - REQ-003 (Route C) → abc1234 [review: 92%]
  - REQ-004 (Route A) → def5678 [review: 88%]
```

## Decision Brief (hand-back format)

The canonical shape for handing work back to the user. Used by the end-of-run completion hand-back (work.md Step 10 / Progress Reporting), `actions/clarify.md`'s question presentation, and `actions/review-work.md`'s Step 9 report. Lead with what was built and what needs the user; **never lead with a self-grade**. Applies `crew-members/anti-slop.md` § 8 (lead with the decision) and the decide-vs-escalate gate in `crew-members/coding-guardrails.md` § Think Before Coding. Render only the sections that have content.

```
WHAT'S BEING BUILT            (feature + subsystem altitude — the value)
  • Now you can X — lives in <subsystem>
  • [MAP CHANGED] <new data flow / renamed concept / new contract>   ← only if true

DECISIONS FOR YOU             (escalated by exception — each carries value + risk)
  ▸ <decision, one line>
      Value:  what you gain / why it matters
      Risk:   cost, reversibility, what breaks if the choice is wrong
      → recommend <X>; default if you say nothing: <X>

HANDLED  (FYI — spot-check, don't ratify)
  • decided <Y> because <Z>     ← reversible calls made without asking
```

- **WHAT'S BEING BUILT** renders each REQ's `## Orientation` block (work.md Step 7.5) at feature/subsystem altitude — not a file list. Anchor to the touched `prime_files`; flag `[MAP CHANGED]` only when the change alters the system's shape.
- **DECISIONS FOR YOU** renders the **ESCALATE**-tier decisions — the `- [~]` / `D-NN` entries that became `pending-answers` follow-ups — each with the Value/Risk carried from the decision record. Source Value/Risk from the touched prime's `## Stakes` when present, else builder-derive.
- **HANDLED** lists the **DECIDE & STATE** decisions (reversible `D-NN` entries) so the user can spot-check without being asked to ratify. Omit if empty.
- **Scale context to reach.** A leaf REQ collapses to a single WHAT'S BEING BUILT line with no DECISIONS and a short HANDLED list; a map-changing REQ earns a short paragraph and a why-it-matters. Review scores never lead — they live under the decision (review-work Step 9) or in the per-REQ progress lines above.
