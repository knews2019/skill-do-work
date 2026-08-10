# Abandon Action

> **Part of the do-work skill.** Marks a REQ as won't-do: sets `status: cancelled`, records the reason, and archives the file — so the decision shows up with finished work on the board instead of haunting the queue as a warning. It is also how an already-archived `failed` REQ is resolved: cancelling it *in place* lets its User Request reach closure — or clears a legacy failure already inside a closed UR folder from the active board — without erasing the failure record.

Cancelling is a first-class outcome, not a deletion. The REQ file survives with its full trail of intent plus a `## Cancelled` section explaining why — six months later, "we decided not to do this, and here's why" is exactly as valuable as "we did this." The canonical status vocabulary (including where `cancelled` sits relative to `failed` and the terminal-success set) is defined in `actions/work-reference.md`'s Schema Read Contract → Terminal-resolved status set.

## When to Use

**Use when:**
- A pending REQ is no longer wanted — priorities changed, the need evaporated, or another REQ superseded it
- A `pending-answers` or blocked REQ isn't worth unblocking — the open questions aren't worth answering
- A `failed` REQ needs resolving so its owning UR can close — cancelling it is the *only* transition out of `failed`, so this applies whether a follow-up REQ already did the recovery work (a completed follow-up never flips the original out of `failed`) or no follow-up is wanted at all. A root failure holds its UR open until it is cancelled this way; a legacy failure already inside a closed `archive/UR-NNN/` folder still needs the same explicit cancellation to leave the board's active columns
- The user says "abandon", "cancel", "won't do", "drop", or "we're not doing this" about a specific REQ

**Do NOT use when:**
- A REQ *just* failed — at classification time the work action writes `failed` (not `cancelled`) and `actions/work.md` Step 8 spawns any follow-up; don't reach for abandon to pre-empt that. (Cancelling an *already-archived* `failed` REQ afterward — once its follow-up has done the work, or when no follow-up is wanted — is a supported use, and the only way to close its UR: see the `failed` bullet under "Use when" above.)
- The user wants to defer a REQ for later — leave it `pending`; the queue is the backlog, and sitting in it costs nothing

## Input

`$ARGUMENTS`: one or more `REQ-NNN` or `UR-NNN` ids (token shapes and UR→REQ expansion per the **Target ID Resolution** contract in `actions/work-reference.md`), optionally followed by free-text — the cancellation reason.

- `do-work abandon REQ-042` — cancel one REQ; the action asks for a one-line reason
- `do-work abandon UR-011` — cancel every cancellable REQ under UR-011 (its non-terminal members, plus any `failed` members still holding it open — see Step 1); the reason applies to all of them
- `do-work abandon REQ-042 superseded by REQ-051` — cancel with the reason inline (everything after the last id token is the reason, applied to every resolved REQ)
- `do-work abandon` (no ID) — list cancellable REQs and ask which; never guess a target. Two groups: non-terminal REQs in `do-work/queue/` and `do-work/working/`, plus already-archived `failed` REQs at `do-work/archive/` root **and `do-work/archive/legacy/`** (both non-recursive — never descend into `archive/UR-NNN/` folders, whose contents are already-closed URs). A legacy nested failure is an explicit-REQ-only repair: the user must name its exact `REQ-NNN`; no-ID discovery never proposes closed history. Show the failed group distinctly ("resolve a failure"): cancelling one closes out a failed attempt so its UR can reach closure, or clears a legacy nested failure from the active board, which is a different act from dropping unstarted work

## Steps

### Step 1: Locate and Gate

**Resolve targeting tokens first (Target ID Resolution contract, `actions/work-reference.md`).** Expand each `UR-NNN` to its member REQs by scanning `user_request:` frontmatter across the locations this action already searches for cancellable work: `do-work/queue/`, `do-work/working/`, and — because a `failed` REQ is cancelled *in place* to close its UR (see the `failed` bullet under "Use when") — `failed` REQs at `do-work/archive/` root and `do-work/archive/legacy/`. That archived-`failed` reach is deliberate: it is exactly the case a UR argument should serve. **Never descend into a closed `do-work/archive/UR-NNN/` folder for UR expansion** (this Input's no-ID bullet already forbids it, and that holds for expansion); the narrow legacy repair below requires an explicitly named `REQ-NNN`. Dedupe the union with any explicitly-named `REQ-NNN`; a token that resolves to nothing is reported and skipped, and an argument list that resolves to an empty set stops the action. The per-REQ gates below then apply **unchanged** to every resolved member, whether named directly or reached through a UR.

For each resolved REQ id, glob `do-work/queue/REQ-NNN-*.md`, `do-work/queue/REQ-NNN.md`, `do-work/working/REQ-NNN-*.md`, and `do-work/archive/**/REQ-NNN*.md`. Normalize its `status` under the Schema Read Contract before applying the gates below, while retaining the original status for the report. Then gate on what you find:

- **Not found anywhere** → report `REQ-NNN: not found` and skip it.
- **Same REQ id found at more than one path** (any two of `do-work/queue/`, `do-work/working/`, `do-work/archive/` root, `archive/legacy/`, or an `archive/UR-NNN/` folder) → refuse and point at `actions/cleanup.md`'s duplicate handling; do not cancel any copy. This mirrors the duplicate that `actions/cleanup.md` Pass 1 already flags (`⚠ Duplicate: REQ-NNN found in both ...`), and it forestalls the split verdict the location/status arms below would otherwise give the same id.
- **Only in archive** → branch on status (the current Failure Classification sends a `failed` REQ to `do-work/archive/` root, or cleanup may place a no-UR failure in `do-work/archive/legacy/`; older skill versions could consolidate one inside `do-work/archive/UR-NNN/`. This location gate is where the failed path is decided, *before* the status rows below are ever reached):
  - `status: failed` at `do-work/archive/` **root or `do-work/archive/legacy/`** → **cancellable in place.** Continue to Step 2; Step 3 flips it to `cancelled` while preserving the failure record, and Step 5 leaves the file exactly where it is (no move). For a root failure this is what lets its held-open UR reach closure; a `legacy/` failure has no `user_request` and holds no UR open, so cancelling it is only housekeeping.
  - An explicitly named `status: failed` REQ found inside an already-archived `do-work/archive/UR-NNN/` folder → **cancellable in place.** This is a narrow legacy repair, never reached by no-ID discovery or `UR-NNN` expansion. Continue to Step 2; Step 3 flips it to `cancelled`, and Step 5 leaves the exact file and its closed UR folder where they are. Do not re-open, move, or re-consolidate the folder. The closed UR stays closed; the board stops treating the nested failure as active.
  - any other archived status (`completed`, `completed-with-issues`, or already-`cancelled`) → report its archive path and status; it's already terminally resolved — nothing to cancel.
- **Status `completed` or `completed-with-issues`** → refuse: finished work is history, not a cancellation target. If the user wants it undone, that's a new capture.
- **Status `failed`** (found in `do-work/queue/` or `do-work/working/`, not archive — unusual, e.g. a hand-edit before cleanup swept it) → **cancellable,** same as the archived-`failed` root case: continue to Step 2, and Step 3 preserves the failure record. (Step 5 will move it out of the queue normally, since it is not yet archived.)
- **Status `claimed`** → warn that a work loop may be mid-flight on it (one orchestrator per queue) and require an explicit extra confirmation before proceeding.
- **Any other status** (`pending`, `pending-answers`, `blocked`, `blocked-*`, or unrecognized) → cancellable; continue.

### Step 2: Confirm the Decision

Show the user what's about to be cancelled — ID, title, current status, and owning UR for **every** resolved target, with a total count (e.g. `Cancel 6 REQs?`), in one prompt (use your environment's ask-user prompt). A `UR-NNN` argument can resolve to many REQs, so this itemized enumeration is the safety property — bulk cancel is the most destructive thing a UR argument enables, and `crew-members/clear-questions.md` (loaded before any interactive question) already requires the prompt state its consequence. **Never a bare count without the per-REQ breakdown, and never a per-member confirmation loop or a `--yes`-style bypass** — one prompt that lists every member is the whole safety mechanism. If no reason was given in `$ARGUMENTS`, ask for a one-line reason in the same prompt; accept "no reason" but never invent one. Do not write anything until the user confirms.

**For a `failed` target the prompt must state the consequence** (per `crew-members/clear-questions.md`, loaded before any interactive question — options state what they buy and cost): cancelling flips it to `cancelled` so a root failure's UR can reach closure, or so a legacy nested failure leaves the active board while its UR folder remains closed. The failure record is *preserved* (`error`/`error_type` stay in frontmatter, and the body's failure sections plus a new `## Cancelled` note remain), and any auto-created follow-up REQ is left untouched — cancel that separately if it too is unwanted.

### Step 3: Write the Cancellation

For each confirmed REQ:

1. Frontmatter: set `status: cancelled` and stamp `completed_at: <now>` (current UTC instant — Timestamp rule, `actions/work-reference.md`) — that timestamp is what places the card in the board's recently-done window. Leave `claimed_at`/`route` and every other field untouched; they're history.

   **Failed → cancelled path:** a failure archived by the work action carries `completed_at` (the failure instant), `error`, and `error_type` — but a legacy or hand-edited failure may be missing any of them (that is exactly the unclassified case `actions/forensics.md` Check 6 flags and now routes here). Handle each field **by presence, never by assumption:**
   - **`completed_at`:** if present, read and hold it (the original failure instant, for the `## Cancelled` line below), then re-stamp `completed_at` to now; if absent, just stamp now. Re-stamping (not preserving) is deliberate: it satisfies the terminal-flip STAMPING RULE (`actions/work-reference.md`) with no new schema field and keeps the board's completion-anomaly check clean, while any recorded failure instant survives in the body.
   - **`error`/`error_type`:** part of "every other field" — leave whatever is present verbatim; that retained pair is the surviving failure signal. **Never fabricate a value for a field that was absent** — in particular, do not fill a missing `error_type` with the Schema Read Contract's `code` default, which is a normalize-a-*present*-value rule, not a write-path fill-in; recording an unclassified failure as `code` invents a failure mode that never happened.
2. Append to the body:

   ```markdown
   ## Cancelled

   - **When:** 2026-07-06T16:45:00Z
   - **Why:** [the user's reason, verbatim — or "no reason given"]
   - **Decided by:** user, via `do-work abandon`
   - **Previously:** failed[ (`error_type: <type>`)][ — failed at <original completed_at>] — resolved by decision not to retry
   ```

   Include the **Previously** line **only** when the prior status was `failed`; omit it entirely for an ordinary cancellation. Fill the bracketed clauses by presence: include `(`error_type: <type>`)` only if that field was set (never invent one); and for the failure instant, if `completed_at` was set write `— failed at <original completed_at>`, otherwise write the literal `— failure instant unrecorded` (one treatment, not a choice — never guess a timestamp). The line, together with whatever `error`/`error_type` frontmatter is retained, preserves the failure signal for body-reading tools (review, ai-report, present-work) that don't surface frontmatter.

Always write the canonical value `cancelled` — never `canceled`, `abandoned`, or `wont-do` (those are read-side aliases only; write paths emit canonical values per the Schema Read Contract).

### Step 4: Surface Dependents

Grep `do-work/queue/` and `do-work/working/` for REQs whose `depends_on` (or legacy `dependencies:`) lists a cancelled ID. A cancelled REQ does **not** satisfy dependency gating, so each dependent would sit blocked forever. For each dependent, ask the user to pick one:

- **Cascade** — abandon the dependent too (loop it back through Steps 1–3)
- **Re-point** — edit its `depends_on` to drop or replace the cancelled ID
- **Leave** — keep it blocked deliberately; it will show under blocked-by-dependencies until edited

Never cascade silently.

### Step 5: Archive

**Already-archived target (the archived-`failed` path from Step 1, including an explicitly named legacy nested failure):** the cancellation was written *in place* — do **not** move the file, and **skip the collision guard below** (the only `do-work/archive/**/REQ-NNN*.md` it would match is the file itself, so the guard would fire against its own target and its "leave it in `do-work/queue/`" remedy is incoherent for a file that was never in the queue). Never relocate it into, out of, or between `do-work/archive/UR-NNN/` folders: consolidating a now-resolved root REQ into its UR is `actions/cleanup.md` Pass 2's job, and a UR folder already sitting in `archive/` is closed. Report the exact in-place path and continue.

Otherwise (the target came from `do-work/queue/` or `do-work/working/`), move each cancelled REQ file out of the queue:

- If `do-work/archive/UR-NNN/` exists for its `user_request` → move it there.
- Otherwise → move it to `do-work/archive/` root (cleanup's Pass 2 consolidates later).
- **Collision guard:** if any `do-work/archive/**/REQ-NNN*.md` already exists, do NOT overwrite — leave the cancelled file in `do-work/queue/`, report the collision with both paths, and let the user resolve it (mirrors `actions/cleanup.md`'s duplicate handling).

### Step 6: Report

Summarize per REQ, note dependents and how each was dispositioned, and check the owning UR: if every sibling REQ is now terminally resolved (`completed`, `completed-with-issues`, or `cancelled`) and its UR is still live, say that `do-work cleanup` will close the UR. For a legacy nested failure, report that the UR folder was already closed and remained in place.

## Output Format

```
Cancelled REQ-042 — [title]
  reason: superseded by REQ-051
  archived: do-work/archive/UR-012/REQ-042-slug.md
  dependents: REQ-047 re-pointed (depends_on: REQ-042 removed)

UR-012: all 3 REQs terminally resolved — `do-work cleanup` will close it.
```

For an already-archived `failed` target, the report shows the in-place path and confirms the failure record survived — no move happened:

```
Cancelled REQ-031 — [title]  (was: failed, error_type: code)
  reason: environment unavailable, not pursuing
  in place: do-work/archive/REQ-031-slug.md  (not moved — already archived)
  failure record: preserved (error_type: code; ## Cancelled notes prior failure)

UR-009: last unresolved REQ resolved — `do-work cleanup` will close it.
```

## Rules

- **Never delete the REQ file.** Cancel + archive preserves the trail of intent; deletion destroys it.
- **Never cancel without confirmation** of the specific REQ IDs — this action removes items from the queue, and the queue is user intent.
- **Only the REQs the user named.** No opportunistic cancelling of stale-looking neighbors.
- **Write canonical `cancelled` only** — aliases are for reading hand-edited files, never for writing.
- Touch nothing beyond the target REQ files, their dependents' `depends_on` (when the user picks re-point), and the archive move — or, for an already-archived `failed` target, the in-place cancellation edit. Never relocate any *other* archived REQ, and never move, re-open, or re-consolidate a UR folder: that is `actions/cleanup.md`'s job, and a UR folder already in `do-work/archive/` stays closed even when an explicitly named nested failure is cancelled in place.

## Common Rationalizations

| If you're thinking...                                   | STOP. Instead...                                            | Because...                                                                 |
| ------------------------------------------------------- | ----------------------------------------------------------- | -------------------------------------------------------------------------- |
| "This attempt failed, so I'll write `cancelled` to skip failure classification" | Let the work action write `failed`; classify and (if wanted) spawn a follow-up | At *classification* time `failed` is correct — it records that work was attempted and drives recovery. `cancelled` is the *after-the-fact* resolution of an already-archived `failed` REQ (Step 1's archived-`failed` path) — applicable whether or not a follow-up ran, since a follow-up never flips the original out of `failed` — never a way to dodge the failure record |
| "This failed REQ is stale — I'll just flip it to `cancelled` to unblock its UR" | Confirm with the user (Step 2) and preserve the failure record | The failure signal is *why* the UR was held open; cancel it as a deliberate decision with `error`/`error_type` and the `## Cancelled` prior-status note intact, not as silent queue hygiene |
| "Deleting the file is cleaner than archiving it"         | Set `cancelled`, append the reason, archive                  | The skill's primary value is the trail of intent — a recorded "no" included |
| "It's claimed but probably stale — cancel it quietly"    | Warn and get explicit confirmation first                     | Another orchestrator may be mid-flight on it; cancelling under it corrupts the run |
| "The queue is long — I'll cancel other stale REQs too"   | Cancel only the named REQs, mention candidates in the report | Staleness is the user's call; the queue is their backlog, not yours         |

## Red Flags

- A REQ file is gone from the repo after an abandon run — deletion happened instead of archival
- Frontmatter says `abandoned`, `canceled`, or `wont-do` — an alias leaked into a write path
- A cancelled REQ still sits in `do-work/queue/` with no reported archive collision — Step 5 was skipped
- A dependent REQ flipped to `cancelled` without the user choosing cascade
- The board shows the cancelled REQ under Needs input / Blocked — `completed_at` wasn't stamped or the status value drifted
- A cancelled-from-`failed` REQ dropped an `error`/`error_type` that it carried *before* the flip, or its `## Cancelled` section has no **Previously** line — the failure signal was erased (the whole point of the failed path is to preserve what was there; a field that was already absent is not "lost")
- A `do-work/archive/UR-NNN/` folder reappeared in `do-work/user-requests/`, or a file moved into or out of a closed UR folder, after an abandon run — the explicit nested-failure path must only cancel the named REQ in place

## Verification Checklist

- [ ] Each cancelled REQ file lives under `do-work/archive/` (UR folder, root, or `legacy/`) — evidence: final file path in the report
- [ ] Frontmatter has `status: cancelled` + `completed_at`; body has a `## Cancelled` section carrying the user's reason verbatim
- [ ] Every dependent found in Step 4 was dispositioned by the user (cascade / re-point / leave) — evidence: one line per dependent in the report
- [ ] For a cancelled-from-`failed` REQ: any `error`/`error_type` present before the flip survive unchanged (none were fabricated for a field that was absent), and the `## Cancelled` section carries the **Previously: failed** line — evidence: frontmatter + body of the archived file
- [ ] For an already-archived target: the flip was in place — no file or folder was moved into or out of `do-work/archive/` — evidence: the report shows the in-place path, not a move
- [ ] No file was deleted, and no file outside the named REQs (plus user-approved dependent edits) was modified
