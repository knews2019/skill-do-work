# Reserve Action

> **Part of the do-work skill.** Invoked when the user wants to reserve pending REQs for a *different* worktree or cloud session working the same repo, or to release such a reservation. A reservation is a queue-resident holding state — the REQ stays in `do-work/queue/` with `status: reserved`, so the local default work loop walks past it while the other session (in its own checkout, after a git sync) picks it up by naming it explicitly. This action belongs in do-work rather than a sibling skill because it manipulates the do-work queue's own status vocabulary (`actions/work-reference.md` → Schema Read Contract) — it is inseparable from the queue schema it writes.

**Reservation ≠ claim.** A claim (`status: claimed`, `actions/work.md` Step 2) moves the REQ into `do-work/working/` and means *this* session is building it right now; crash recovery treats anything in `working/` as interrupted work and re-queues it. A reservation never touches `working/`, never survives into crash recovery, and carries an owner label — it means "allocated to someone else; hands off by default." Design intent, cross-file contract, and traps: `actions/prime-req-reservation.md`.

## When to Use

**Use when:**
- You're about to spin up (or have already spun up) another worktree or cloud session to work specific REQs, and local `do-work run` must not claim them.
- Splitting one queue across parallel sessions: reserve each session's slice under its own label.
- Reviewing current allocations (`do-work reserve` with no arguments lists them, with staleness flags).

**Do NOT use when:**
- The REQ won't ever be done — that's `do-work abandon` (`actions/abandon.md`).
- The REQ must wait for *other REQs* — that's `depends_on` gating, not a reservation.
- The REQ is waiting on user answers — that's `pending-answers` (`do-work clarify`).

## Input

`$ARGUMENTS` selects the mode:

| Arguments | Mode |
|---|---|
| `REQ-NNN [REQ-NNN ...] for <label>` | **reserve** — mark each REQ reserved for `<label>` |
| `release REQ-NNN [REQ-NNN ...]` | **release** — return the named REQs to `pending` |
| `release <label>` | **release** — return every REQ reserved under `<label>` to `pending` |
| (empty) | **list** — show current reservations with age + staleness |

The label is free text naming the owning session ("cloud-alpha", "worktree feature-auth"). If REQ IDs were given without a label, load `crew-members/clear-questions.md` and ask for one with your environment's ask-user prompt — offer a suggested default derived from context (branch name, session name) as the recommended option. Never invent a label silently; the label is how a human later tells whose reservation this is.

## Steps

### Mode: reserve

1. For each named REQ, glob `do-work/queue/REQ-NNN-*.md` (and `do-work/queue/REQ-NNN.md`). Missing → report and skip.
2. Read frontmatter and normalize `status` per the Schema Read Contract (`actions/work-reference.md`). Only `pending` REQs are reservable. Anything else → report why and skip: `claimed`/in `working/` means a session already owns it; `pending-answers` needs `do-work clarify` first; `blocked` needs its external condition cleared first (`do-work run` re-probes it, or `do-work clarify` confirms it); `reserved` is already allocated (report the existing `reserved_for` and stop — re-labeling an existing reservation requires an explicit release first).
3. Update frontmatter on each reservable REQ: `status: reserved`, `reserved_for: "<label>"` (always YAML-quoted — the label is raw user text; treat it as data, never interpolate it into a shell command), `reserved_at: <timestamp>` (current UTC instant — `date -u +%Y-%m-%dT%H:%M:%SZ`; Timestamp rule, `actions/work-reference.md` — a future-dated stamp breaks the board's staleness math).
4. Report per the Output Format, and remind: the reservation only protects sibling sessions **after it syncs** — commit and push the queue edit (or let the user's normal flow do it). Files are the only channel between checkouts.

### Mode: release

1. Resolve targets: REQ IDs → glob as above; a label → scan `do-work/queue/REQ-*.md` frontmatter for `status: reserved` with matching `reserved_for`. No matches → report and stop.
2. For each target with `status: reserved`: set `status: pending`, stamp `status_changed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`), remove `reserved_for` and `reserved_at`. A target that is not `reserved` is reported and left untouched — release never rewrites other statuses.
3. Report which REQs re-entered the queue.

### Mode: list

1. Scan `do-work/queue/REQ-*.md` frontmatter for `status: reserved`.
2. Render each with label and age (now − `reserved_at`). A reservation **older than 24 hours is stale** — flag it and suggest recategorizing (see Output Format). Staleness is a *suggestion trigger*, not an auto-release: the other session may legitimately still be working; only a human knows.

## Output Format

```
Reserved 2 REQs for "cloud-alpha":
  REQ-042 — [title]
  REQ-043 — [title]
Sync reminder: commit & push do-work/queue/ so other sessions see the reservation.
```

List mode (and the stale flag, wherever reservations render):

```
Reservations:
  REQ-042 — [title] (reserved for: cloud-alpha, 3h ago)
  REQ-051 — [title] (reserved for: worktree feature-auth, 31h ago) ⚠ STALE
⚠ 1 reservation older than 24h. Recategorize: `do-work release REQ-051` to return it to the queue,
  `do-work run REQ-051` to claim it in this session, or leave it if the other session is still active.
```

## Rules

- **Only `pending` REQs are reservable.** Release always restores exactly `pending` — reserve never captures, and release never invents, any other status.
- **Never move the file.** Reservations live in `do-work/queue/`; `do-work/working/` belongs exclusively to the work pipeline's claim.
- **Targeted runs override.** `do-work run REQ-NNN` claims a reserved REQ (clearing the reservation) — that is the designed pickup path for the owning session and the human override for everyone else. Only the *default* full-queue scan honors reservations. Don't add extra guards to targeted mode.
- **Quote the label.** `reserved_for` is raw user text — write it as a quoted YAML scalar and never substitute it into shell commands unquoted.
- **24-hour staleness is advisory.** Flag and suggest; never auto-release.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "I'll reserve it by moving it into `working/` — that's what claimed means" | Set `status: reserved` in place in `do-work/queue/` | Crash recovery re-queues everything in `working/` at the next session start — the reservation would be silently destroyed |
| "The reservation is set locally, job done" | Remind the user to sync (commit/push) the queue edit | Worktrees and cloud sessions have separate checkouts; an unsynced reservation protects nothing |
| "This reservation is >24h old, I'll release it automatically" | Flag it and suggest the three recategorize options | The other session may still be mid-build; auto-release invites double-claiming |
| "The label has an apostrophe, I'll inline it into the sed/yq command" | Edit the frontmatter as a file operation with the label as quoted data | Raw text inside shell quoting is a breakage and injection vector (see CLAUDE.md's prescribed-command traps) |

## Red Flags

- A REQ with `status: reserved` sitting in `do-work/working/` — reservations never live there; something conflated reserve with claim.
- `reserved_for` missing or empty on a `reserved` REQ — the owner is unknowable; ask the user rather than guessing.
- Reservations accumulating past 24h across multiple lists without recategorization — sessions died without releasing; prompt the user to sweep them.

## Verification Checklist

- [ ] Every reserved REQ still sits in `do-work/queue/` with `status: reserved`, a quoted `reserved_for`, and a `reserved_at` timestamp.
- [ ] Every released REQ is back to `status: pending` with both reservation fields removed.
- [ ] No file was moved into or out of `do-work/working/` by this action.
- [ ] Stale (>24h) reservations were flagged with the recategorize suggestions, not auto-released.
