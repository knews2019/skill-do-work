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
error: "Description"          # Only if failed

# Set by abandon action (do-work abandon — user-directed won't-do decision)
status: cancelled             # terminal, NOT successful; the reason lives in the REQ body's `## Cancelled` section
completed_at: 2025-01-26T10:45:00Z  # stamped at cancellation — the terminal timestamp the board's recently-done window reads

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

### Terminal-success status set

**A REQ counts as *terminally successful* when its `status` is `completed` or `completed-with-issues`.** This is the canonical set every reader that selects "completed work" must honor — `completed-with-issues` is terminal and counts toward UR completion (it just carries known follow-ups, per `actions/work.md` Step 8), so a filter that accepts only the literal `completed` silently drops remediated-with-issues work. `failed` is terminal but **not** successful — success-readers exclude it.

The trigger is the *condition above*, not the caller list: **any reader that filters for "the completed/most-recent work" inherits this contract.** The known consumers are illustrative, not exhaustive — `actions/cleanup.md` (UR close), `actions/ai-report.md` (report target), `actions/review-work.md` (standalone target), and `actions/commit.md` (REQ association); `actions/forensics.md` and `actions/roadmap.md` already honor both. When adding a new reader, accept both values and point back here — hand-enumerated caller lists go stale silently, which is why the condition, not the list, is the contract.

### Terminal-resolved status set

**A REQ counts as *terminally resolved* when its `status` is `completed`, `completed-with-issues`, or `cancelled`.** This is the set archive-sweep and UR-closure readers honor (`actions/cleanup.md` Pass 0 + Pass 1, `actions/work.md` Step 8's UR-final check, this file's Composed Exit Summary): `cancelled` records a deliberate won't-do decision — made via `do-work abandon` — so it archives like finished work and must never hold a UR open the way `failed` does. Three boundaries keep the sets honest:

- `cancelled` is **not** successful. Success-readers (the Terminal-success set above) exclude it — a cancelled REQ is never a review-work target, an ai-report subject, or a commit association.
- `cancelled` does **not** satisfy `depends_on` gating. A dependent presumably needed the cancelled REQ's output; the abandon action surfaces dependents at cancellation time so the user can cascade the cancellation or re-point `depends_on`.
- `failed` stays outside this set: it is terminal and unsuccessful, but it signals work that *should* have happened — Step 8's failure classification spawns follow-ups, and a UR with a `failed` REQ needs those follow-ups before it can close. Cancelling is the explicit way to say "no follow-up wanted."

## Crash Recovery (Step 1)

**Crash Recovery:** Before checking the queue, look inside `do-work/working/` for any `REQ-*.md` files. If any exist, a previous run may have been interrupted — but a live coexisting session's claimed file looks identical to a crashed one, so a concurrency gate runs first.

**Concurrency gate (apply every time this runs, including on a Step 10 → Step 1 loop iteration within the same session — it re-reads the lock file fresh each time, so it never depends on this session's own memory of an earlier lock-guard choice).** Take the `working/` listing **first, then** read `do-work/orchestrator-lock.json` — this order is load-bearing, not stylistic. Step 2 writes a REQ's claim into the lock *before* moving the file into `working/` (`actions/work.md` Step 2 substep 1), so any file already present when the listing is taken is guaranteed to have its claim visible in a lock read taken after it; reading the lock first reopens the race, because a REQ claimed-and-moved between that stale lock read and the scan would look unclaimed and get re-queued out from under a live session. For each `REQ-*.md` found in the listing, extract its id and check whether a session *other than this one* currently and freshly claims it: the id matches the top-level holder's `claimed_req` (holder `session_id` ≠ this session's) with `heartbeat_at` ≤45 minutes old, or it matches any `coexisting_sessions[].claimed_req` (entry `session_id` ≠ this session's) with that entry's `heartbeat_at` ≤45 minutes old. **Skip that file** if so — it belongs to a live session, not an interrupted one. Recover every other file below, including files whose claiming session has itself gone stale past 45 minutes — that's the self-heal path. Full lock file shape: `actions/work-reference.md` → **Concurrent-Orchestrator Lock Guard**.

For each REQ file that passes the gate (not actively claimed by another live session):
1. Reset frontmatter: set `status` to `pending`, **unless** the REQ file contains a `## Open Questions` section with at least one unresolved `- [ ]` item — in that case, restore `status` to `pending-answers`. (If the `## Open Questions` section exists but all items are already `[x]` or `[~]`, or if no `## Open Questions` section exists at all, set `status` to `pending`.) **Exception — a recovered REQ that already carries `status: blocked` with a `blocked_by` condition stays `blocked`** (the mid-run blocked flip completed its frontmatter write before the crash; its condition is unchanged and it must not be silently promoted to runnable). Remove `claimed_at` and `route`; leave `blocked_by`/`blocked_at`/`blocked_check` intact.
2. Strip sections generated during the interrupted run: remove `## Triage`, `## Exploration`, `## Plan`, `## Scope`, `## Pre-Flight`, `## Implementation Summary`, `## Qualification`, `## Testing`, `## Review`, `## Lessons Learned`, `## Orientation`, `## Decisions`, and `## Discovered Tasks` sections (and their content) if present — these may be incomplete or stale from the crash. Leave `## Open Questions` and user-authored content intact.
3. Move the REQ back to `do-work/queue/`

Once every file this session is allowed to touch has been recovered (a live session's file may remain in `working/` — that's expected, not an error), proceed with finding the next request.

## Concurrent-Orchestrator Lock Guard (Step 1)

Guards the one-orchestrator-per-queue convention: a second `do-work run` starting on the same working tree must detect the first before claiming or processing anything. Not a distributed lock — but every update to the lock file is a read-modify-write on a file with more than one writer, so it must be serialized (**Serialized lock updates**, below). Plain unguarded reads and writes are *not* sufficient. (Origin incident: two sessions ran `do-work run` on the same tree on 2026-07-01; the second committed and archived the first's in-flight REQ with a hollow paper trail. The collision merged by luck — nothing detected it.) **Design note (REQ-018 remediation):** the lock records exactly one HOLDER plus zero or more `coexisting_sessions` — a session that picks "proceed anyway" gets its own tracked entry instead of displacing the holder, so the holder's `session_id` never stops matching and the holder is never disrupted.

**Lock file:** `do-work/orchestrator-lock.json`. Re-derive this path deterministically every time you touch it — do not carry a shell variable across separate command blocks. `mkdir -p do-work` first if the directory doesn't exist yet.

```json
{
  "session_id": "do-work-20260727T083332Z-51234",
  "started_at": "2026-07-27T08:33:32Z",
  "heartbeat_at": "2026-07-27T08:41:10Z",
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
  "claimed_req": "REQ-018",
  "coexisting_sessions": [
    {
      "session_id": "do-work-20260727T085501Z-52890",
      "joined_at": "2026-07-27T08:55:01Z",
      "heartbeat_at": "2026-07-27T09:02:44Z",
      "claimed_req": "REQ-042"
    }
  ]
}
```

`session_id` is `do-work-$(date -u +%Y%m%dT%H%M%SZ)-$$` — a UTC timestamp plus the current shell's process id, generated once at acquisition (or at joining, for a coexisting session). The orchestrating agent already knows this string from its own context for the rest of the run (it isn't a shell export that needs to survive between command blocks) — every later touch just re-types the same value into a fresh, deterministically-derived command. The top-level `session_id`/`started_at`/`heartbeat_at`/`claimed_req` describe the recorded **holder**; `claimed_req` names whatever REQ currently sits in `do-work/working/` under the holder's claim, or `null` between REQs. `coexisting_sessions` lists other sessions running alongside the holder via **proceed anyway** (below) — each entry is `{session_id, joined_at, heartbeat_at, claimed_req}`, refreshed independently of the holder's own fields; empty (`[]`) when no session has chosen to coexist. All timestamps follow the Timestamp rule above (current UTC instant).

**Serialized lock updates (every write, no exceptions).** More than one session writes this file — a coexisting session refreshes its own entry while the holder refreshes the holder fields, and both the heartbeat prune and a take-over rewrite *other* sessions' entries. So "each writer only changes its own logical fields" describes intent, not atomicity: if the holder reads at T0, a coexisting session records its `claimed_req` at T1, and the holder writes its T0 image back at T2, that claim is silently lost. The Crash Recovery per-file gate then sees no live claim on a REQ that is actively being built and re-queues it — the exact 2026-07-01 failure this guard exists to prevent.

Every path that writes the lock — acquisition, take-over, the **proceed anyway** append, each heartbeat refresh, the opportunistic prune, and release — wraps its read-modify-write in this mutex, and re-reads the file *inside* the critical section rather than reusing an earlier read:

```bash
lock_path="do-work/orchestrator-lock.json"
mutex_path="do-work/orchestrator-lock.json.mutex"

# `mkdir` is atomic on POSIX: exactly one contender wins, the rest fail.
mutex_attempts=0
until mkdir "$mutex_path" 2>/dev/null; do
  # A session that died mid-update must not wedge the queue. The critical section is
  # sub-second, so a mutex older than a minute is abandoned by definition — break it.
  if [ -d "$mutex_path" ] && [ -z "$(find "$mutex_path" -maxdepth 0 -mmin -1 2>/dev/null)" ]; then
    rmdir "$mutex_path" 2>/dev/null || true
    continue
  fi
  mutex_attempts=$((mutex_attempts + 1))
  if [ "$mutex_attempts" -ge 30 ]; then     # ~15s of contention on a sub-second section
    rmdir "$mutex_path" 2>/dev/null || true # treat as abandoned and take it
    continue
  fi
  sleep 0.5
done

# --- critical section: re-read, modify, write atomically ---
# (read $lock_path fresh here, apply ONLY this session's own field changes.
#  An ACQUIRING session whose fresh read finds a lock file that now exists lost
#  the acquisition race: write nothing, release the mutex, and re-enter the
#  procedure's existing-lock branch — see Procedure step 1.)
lock_tmp="$lock_path.$$.tmp"
printf '%s\n' "$updated_lock_json" > "$lock_tmp" && mv -f "$lock_tmp" "$lock_path" || rm -f "$lock_tmp"

rmdir "$mutex_path" 2>/dev/null || true
```

The temp-file-plus-`mv` is not decoration: `mv` within a directory is an atomic rename, so a reader outside the mutex can never observe a half-written lock file. Never write `$lock_path` in place. **Never block the pipeline on this** — the mutex always yields within ~15 seconds by design; a wedged queue is worse than a rare unserialized update.

**Never committed.** At acquisition (first write), append the local-ignore snippet from `crew-members/background-agents.md` → **Local-ignore snippet (for genuinely-transient paths)**, substituting `do-work/orchestrator-lock.json` for `<path>`, and again for `do-work/orchestrator-lock.json.*` — the glob (not an exact path) is required because it must cover both the mutex directory and the PID-suffixed temp files: a `do-work/orchestrator-lock.json.<pid>.tmp` orphaned by a session killed between its write and its rename carries a PID no later session repeats, so only a pattern keeps it out of `git status` and off the committable surface. The mutex directory and the `.tmp` files are the same kind of transient session state as the lock itself. Same mechanism `do-work/pipeline.json` already uses (`actions/pipeline.md` Step 2 substep 4) — a lock is transient session state, not part of the Trail of Intent, so it must stay out of git regardless of install layout.

**Stale threshold: 45 minutes.** Generous on purpose — a Route C build can legitimately run that long between two of the heartbeat touchpoints below, and a false "stale" verdict on a genuinely live session is worse than a slightly slower self-heal on a genuinely dead one. Applies identically to the holder's `heartbeat_at` and to every `coexisting_sessions[].heartbeat_at`.

**Interactivity test (concrete, not self-assessed).** A session is interactive only if BOTH hold — if either can't be confirmed, it is non-interactive:
1. **This is the top-level session the user is directly conversing with, not a subagent spawned via Task/Agent dispatch.** A subagent has no user turn to pause on — it returns its final message to whatever dispatched it, never to a human. This is the same distinction `actions/code-review.md` already draws for its own consent gate ("running non-interactively (e.g., via subagent)") and the harness-availability check in `actions/kb-lessons-handoff.md`, applied here to the identical failure mode: `actions/pipeline.md`'s foreground dispatch of the `run` step explicitly may run as a subagent ("subagent if available, inline otherwise" — `actions/pipeline.md` Step 4 substep 2), and a subagent that wrongly assumes it can prompt hangs waiting for an answer nobody can give.
2. **An ask-user / blocking-prompt mechanism is actually available and usable in this environment right now** — not merely documented or importable.

**Failure-safe default:** if either check is unconfirmed or ambiguous, treat this session as non-interactive and take the refusal branch below. A refusal is a clean, resumable no-op — re-run `do-work run` directly (not via subagent dispatch) to get the interactive prompt. A hung ask-user call blocks the entire pipeline with no way to resume, which is strictly worse than a refusal.

**Procedure — run once, the first time this session reaches Step 1** (session start, or resuming via CHECKPOINT.md); skip it on a Step 10 → Step 1 loop iteration once this session's own `session_id` already appears in the lock file (as holder or in `coexisting_sessions`) — a durable, file-backed fact, not something this session needs to remember on its own.

1. **No lock file exists** → acquire: generate `session_id` and write the fresh holder object (`started_at`/`heartbeat_at` = now, `claimed_req: null`, `coexisting_sessions: []`) through the serialized mutex — **re-validating existence inside the critical section**: the branch was chosen on an unserialized existence check, so two sessions starting on an empty queue can both select it. If the fresh read inside the mutex finds a lock file that now exists, another session won the race — write nothing, release the mutex, and re-enter this procedure at step 2 with that file (its holder's heartbeat is seconds old, so the fresh-holder path applies). Only when the fresh read still finds no lock does the write proceed. Append the local-ignore snippet. Continue to the CHECKPOINT.md check and Crash Recovery normally.
2. **Lock file exists** → read it, compute `age = now − heartbeat_at` in minutes against the holder's `heartbeat_at` (same style as Step 1's Stale-reservation check comparing `reserved_at` against now).
   - **`age > 45` (holder stale):** warn `Lock held by <session_id> looks abandoned (heartbeat <age>m old; was building <claimed_req, or "nothing">). Treating as stale — taking over.` Rewrite only the top-level holder fields — freshly generated `session_id`, `started_at`/`heartbeat_at` = now, `claimed_req: null`. **Preserve `coexisting_sessions`**, pruning only entries whose own `heartbeat_at` is also past 45 minutes — a live coexisting session keeps running, undisturbed by a holder take-over. Continue to the CHECKPOINT.md check and Crash Recovery normally — the per-file concurrency gate (`actions/work-reference.md` → **Crash Recovery (Step 1)**) decides which `working/` file this take-over actually reclaims.
   - **`age ≤ 45` (holder fresh — another orchestrator may be live):** apply the Interactivity test above.
     - **Interactive:** load `crew-members/clear-questions.md`, then present:

       ```
       Another do-work session appears active on this queue.
         Holder: <session_id>
         Started: <started_at>
         Heartbeat: <age>m ago (fresher than the 45m stale threshold)
         Currently building: <claimed_req, or "nothing claimed yet">
         Also coexisting: <N other session(s)> (omit this line if coexisting_sessions is empty)

       A second orchestrator on the same queue can double-claim or double-commit a REQ (this happened on 2026-07-01) — pick one:
         (a) Proceed anyway — run alongside it as a tracked coexisting session. I will scan for other pending REQs but will NOT touch whatever the holder (or another coexisting session) has actively claimed in do-work/working/. The holder is untouched — it keeps running exactly as if I weren't here.
         (b) Take over — assume the holder is gone even though its heartbeat looks fresh. I become the new holder and run normal crash recovery on its claimed do-work/working/ file (reset and re-queued); any other live coexisting session is left alone.
         (c) Abort — stop now. I will not write anything; the existing lock is left exactly as it was.
       ```

       - **(a) Proceed anyway** — append a new entry to `coexisting_sessions`: `{session_id: <new>, joined_at: now, heartbeat_at: now, claimed_req: null}`. Do **not** touch the holder's `session_id`/`started_at`/`heartbeat_at`/`claimed_req` — the holder's own heartbeat-refresh touchpoints keep matching, so it runs to completion undisturbed. Continue to the CHECKPOINT.md check and Crash Recovery normally: the per-file concurrency gate skips only files actively claimed by another live session (the holder's, or another coexisting session's) — it is not an all-or-nothing skip.
       - **(b) Take over** — same rewrite as the stale path above (new holder `session_id`, `claimed_req: null`; preserve/prune `coexisting_sessions`).
       - **(c) Abort** — make no writes at all; report the holder info above and stop before the CHECKPOINT.md check, Crash Recovery, or the queue scan.
     - **Non-interactive** (Interactivity test above resolved "no," including "couldn't confirm"): refuse.

       ```
       Refusing — another do-work session appears active on this queue and no interactive prompt is available.
         Holder: <session_id>
         Started: <started_at>
         Heartbeat: <age>m ago
         Currently building: <claimed_req, or "nothing claimed yet">
         Also coexisting: <N other session(s)> (omit this line if coexisting_sessions is empty)
       Re-run interactively to choose proceed/take-over/abort, or wait for the heartbeat to pass the 45-minute stale threshold.
       ```

       Then stop without creating, modifying, or deleting any file — lock, queue, or `working/`.

**Heartbeat refresh (every subsequent touch — not the guard above).** At Step 2 (claim), Step 6.25 (Implementation Summary), Step 8 (Archive), Step 9 (Commit Phase), and Step 10 (loop point): re-derive the lock path and read it.
- If this session's `session_id` matches the top-level holder: rewrite the holder's `heartbeat_at` = now and `claimed_req` = whatever's currently in `do-work/working/` under this session's claim (or `null`).
- Else if this session's `session_id` matches an entry in `coexisting_sessions`: rewrite *that entry's* `heartbeat_at` and `claimed_req` in place. Leave the holder and every other entry untouched.
- Else (this session's `session_id` appears nowhere in the lock file): it was superseded — another session's take-over overwrote the holder slot, or this session's own coexisting entry was pruned as stale. Stop and report rather than silently re-adding yourself — this is what keeps "file-based" from turning into "two sessions both believing they hold it."

While performing any of the checks above, opportunistically prune `coexisting_sessions` entries (other than the one just refreshed) whose `heartbeat_at` is past the 45-minute threshold — this is what lets an abandoned coexisting session's claim self-heal via the Crash Recovery per-file gate on a later Step 1 entry, with no separate cleanup pass.

**Release last — after every repo mutation this session will make, including `actions/cleanup.md` and its commit.** Releasing is what lets the next `do-work run` start claiming REQs and moving files, so anything this session does afterward races that new session over the same working tree. Cleanup is the trap: its Pass 0 sweeps `do-work/queue/` and `do-work/working/` and stages those paths (`actions/cleanup.md`), which is exactly where an arriving session has just put its freshly claimed REQ — release-then-cleanup lets a departing session sweep and commit the newcomer's in-flight work. The lock costs nothing to hold through a cleanup pass; only Step 10's exit runs cleanup, so this ordering matters there in particular.

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

**Validation check (successful REQs only):** Before committing, compare the `## Implementation Summary` file list against the staged files (excluding `do-work/` paths). If the Implementation Summary lists files that aren't staged, or if the only staged files are `do-work/` metadata, `CHANGELOG.md`, and/or the version file it bumped (the changelog entry and the version bump describe the implementation, they aren't the implementation), flag the mismatch — the commit may not contain the actual implementation. Fix the staging or update the Implementation Summary before proceeding. Design-artifact files placed outside `do-work/` satisfy this check — they are project deliverables. **Skip this check for failed REQs** — they may have no Implementation Summary or no project files staged, and that's expected.

**Write commit hash back to the archived REQ.** After the commit succeeds, retrieve the hash with `git rev-parse --short HEAD` and update the archived REQ's frontmatter `commit:` field with the actual value. Then create a **separate metadata commit** (do not amend — amending changes the hash and invalidates what you just wrote):

```bash
# After the initial commit succeeds:
HASH=$(git rev-parse --short HEAD)
# Update the commit: field in the archived REQ frontmatter
# (replace "commit:" line or add it if missing)
git add do-work/archive/UR-NNN/REQ-NNN-slug.md
git commit -m "[REQ-NNN] record commit hash ${HASH}"
```

This ensures the `commit:` field in the archived REQ contains the real implementation commit hash, which the review-work and present-work actions depend on for traceability. The metadata commit is a lightweight bookkeeping entry — it does not contain implementation changes.

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
