# Clarify Questions Action

> **Part of the do-work skill.** Batch-reviews pending questions from completed work — the user confirms, overrides, or discards builder decisions.

This is the second human-attention window in the pipeline. After actions/work.md processes requests autonomously, any ambiguities the builder encountered are surfaced here as a batch for efficient review.

## When to Use

**Use when:**
- A work run just finished and left `pending-answers` REQs in the queue.
- The user asks "what's blocked?", "show me pending questions", or similar.
- The pipeline can't advance because builder-decided questions need sign-off.
- `blocked` REQs are waiting on a **human-confirmable** external condition (a designer answered, a service is now up) and you want to confirm the condition is met so they re-enter the queue.

**Do NOT use when:**
- No `pending-answers` **and** no `blocked` REQs exist — tell the user the queue is clear and stop.
- The user wants to answer a *specific* open question by editing the REQ directly — that's just a file edit, not a batch review.
- The queue has only `pending` REQs — those need `do-work run`, not clarify.
- The `blocked` condition is **machine-checkable** (it carries a `blocked_check` probe) — `do-work run` auto-probes and unblocks those; clarify is for the human-confirmable ones.

## Input

Triggered by `do-work clarify` (also: `answers`, `questions`, `pending`, `what's blocked`). No arguments needed.

## Steps

### Step 1: Scan the queue

Find all `REQ-*.md` files in `do-work/queue/` with `status: pending-answers`. Also collect REQs with `status: blocked` (waiting on an external condition) for Step 5.5.

### Step 2: Check for pending questions

If neither any `pending-answers` REQ nor any `blocked` REQ is found: report "No pending questions or blocked REQs — queue is clear" and exit. If only `blocked` REQs exist (no `pending-answers`), skip Steps 3–5 and go straight to Step 5.5.

### Step 3: Present questions

**Load `crew-members/clear-questions.md` first** — it is the contract for every question you're about to show. Then **rewrite** each REQ's question to that contract instead of rendering the stored `## Open Questions` text verbatim: that text was authored mid-implementation by a builder with the whole spec in its head, and is presumed too dense for a cold reader. Gloss every coined label or section reference, and state why the decision was escalated to the user (Principle 7). The rewrite applies to the question and option *wording* — the Decision Brief structure below stays as-is.

Always lead with the builder's decision and its default — confirming is the intended fast path — and, for builder-decided follow-ups, show the **value and risk** so the user can judge in seconds instead of spelunking. This is the **DECISIONS FOR YOU** section of the Decision Brief (`actions/work-reference.md` → **Decision Brief (hand-back format)**). For each `pending-answers` REQ, show:

```
REQ-025 — Review fix: dark mode sidebar
(follow-up to REQ-003, from review)

1. [ ] Should the sidebar use the same dark palette as the main content?
   Decision: Yes, match main content palette   ·   Default if you say nothing: same
   Value: one consistent dark surface; nothing to re-theme later
   Risk:  if you wanted a distinct sidebar, this is a quick CSS revert (low, reversible)
   Also:  Separate sidebar palette, User-configurable

2. [ ] Should dark mode persist across sessions?
   Decision: Yes, save to localStorage   ·   Default if you say nothing: same
   Value: returning users keep their choice
   Risk:  a stale stored value can mask the OS-preference path (medium, reversible)
   Also:  Reset on refresh, Follow OS preference
```

**Fallback (mandatory).** Many `pending-answers` REQs come from templates that don't carry Value/Risk — capture, verify-requests, review-work follow-ups, and discovered tasks all emit `Recommended:`/`Also:` only. When a question has no `Value:`/`Risk:` lines, render it in that older form (`Recommended:` + `Also:`) — never block on the missing fields.

Builder-marked `- [~]` decisions reflect the "Think Before Coding" guardrail (`crew-members/coding-guardrails.md`) — surface tradeoffs early, not late.

### Step 4: Collect answers

If your environment has a structured question prompt (multi-question UI), batch questions in groups of **at most 4 per prompt** — chunk by question count, not by REQ. A REQ with 6 questions needs 2 prompts.

For each question, the user can:

- **Answer it** → update to `- [x] [question] → [user's answer]`
- **Confirm builder's choice** → update to `- [x] [question] → Confirmed: [builder's choice]`. Then check the REQ type:
  - *Discovered-task REQ* (has a "Should I process this as a new task?" question with recommended "Yes, add to queue"): flip `status` to `pending` so the task enters the work queue — see "Approved Discovered Task" below
  - *All other REQs* (builder-decision follow-ups): mark `status: completed` (no implementation needed — see "Builder Was Right" below)
- **Pick a different option** → update to `- [x] [question] → [user's chosen option]`
- **Skip for now** → leave as `- [ ]`, REQ stays `pending-answers`
- **Discard it** → update to `- [x] [question] → Discarded`, then mark the REQ `status: cancelled`, `completed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`), and archive it directly (same fast-path as "Builder Was Right", but with the honest won't-do status — no work happened and none is wanted; see "Discarded" below)

### Step 5: Activate answered REQs

For each REQ that wasn't already completed or discarded: if all questions are now `[x]` or `[~]`, flip `status` from `pending-answers` to `pending` and stamp `status_changed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`; the board's state timer reads it, so the card shows time since the answers landed rather than time since capture). These enter the queue for the next `do-work run`.

### Step 5.5: Confirm blocked conditions

For each `status: blocked` REQ collected in Step 1, present its condition as one lightweight yes/no — no rewrite-contract machinery needed (the condition is a single line of `blocked_by` text, not a builder question):

```
REQ-042 — Wire up local translation
Blocked by: LM Studio running locally (since 3d ago)
Is this condition now satisfied?
  1. Yes — unblock it        2. Not yet — leave it        3. Abandon this REQ
```

Note for the user which blocked REQs carry a `blocked_check` probe — those unblock automatically on the next `do-work run`, so confirming them by hand here is optional. Present only the human-confirmable ones prominently.

- **Yes → unblock:** set `status: pending`, stamp `status_changed_at: <timestamp>` (blocked_at is removed on this flip, so this is the only trace of when it happened), **remove `blocked_by` and `blocked_at`** (keep any `blocked_check`), and append a history line to a `## Blocked` body section — `- [<date>] blocked on "<condition>" — cleared by user via clarify`. The REQ re-enters the queue for the next `do-work run`.
- **Not yet:** leave it `blocked`, unchanged.
- **Abandon:** hand off to `do-work abandon REQ-NNN` (marks `cancelled`, archives) — same as discarding a question.

### Step 6: Report

Summary of what was resolved and what's still pending — include any `blocked` REQs unblocked (now `pending`) or left waiting, alongside the answered/confirmed/discarded questions.

## Builder Was Right / Discarded

When the user reviews a `pending-answers` follow-up and confirms that the builder's original choice was correct (i.e., no implementation change needed):

1. Update the question to `- [x] [question] → Confirmed: [builder's choice]`
2. Update frontmatter: `status: completed`, `completed_at: <timestamp>` (current UTC instant — `date -u +%Y-%m-%dT%H:%M:%SZ`; Timestamp rule, `actions/work-reference.md`)
3. Archive the follow-up REQ directly (skip the work loop — there's nothing to build)
4. Append a brief note: `## Implementation\n\n**No changes needed.** User confirmed builder's choice from [original REQ].\n\n*Resolved via clarify questions*`

**Discarded** (questions or discovered tasks the user declines): the same fast-path applies, but the status is `cancelled`, not `completed` — nothing was built and nothing is wanted, and `cancelled` is the canonical won't-do terminal status (`actions/work-reference.md` → Terminal-resolved status set; it closes URs and shows with done work on the board). Mark `status: cancelled`, `completed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`), archive directly, and append:

```markdown
## Cancelled

- **When:** <timestamp>
- **Why:** user discarded this during clarify — [the question/task, one line]
- **Decided by:** user, via `do-work clarify`
```

## Approved Discovered Task

When the user reviews a discovered-task follow-up (one whose question is "Should I process this as a new task?" with recommended "Yes, add to queue") and confirms the recommendation:

1. Update the question to `- [x] [question] → Confirmed: Yes, add to queue`
2. Update frontmatter: `status: pending` (NOT `completed` — this task needs to be built), plus `status_changed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`)
3. **Do not archive.** The REQ stays in `do-work/queue/` and enters the normal work queue for the next `do-work run`

This is distinct from "Builder Was Right" because confirming a discovered task means the user wants it *executed*, not signed off. The task has no prior implementation to confirm — it's a new piece of work that needs a full work cycle.

## Rules

- This action avoids wasting a work cycle on a REQ that just needs sign-off or rejection, while correctly routing approved discovered tasks into the build queue
- Never block the user — if they skip all questions, exit gracefully
- Always show the builder's recommended choice prominently so confirming is the fast path

## Red Flags

- A question shown to the user still contains unglossed builder shorthand (a coined label, a spec §-reference, a finding number) — the stored text was rendered verbatim instead of rewritten per `crew-members/clear-questions.md`.
- A `pending-answers` REQ with no `## Open Questions` section — the marker and the body disagree; investigate before presenting nothing.
- User confirms every builder choice without reading — they may be rubber-stamping; ask once if they want a summary first.
- A discovered-task follow-up's `status` flipped to `completed` instead of `pending` after user confirmed "Yes, add to queue" — that's the wrong route (the task never gets built).
- `pending-answers` REQs pile up across multiple clarify runs without resolution — users are skipping; ask whether to discard the stale ones.

## Verification Checklist

- [ ] `crew-members/clear-questions.md` was loaded before the first question was presented, and stored question text was rewritten to its contract (not rendered verbatim).
- [ ] Every REQ presented had `status: pending-answers` in its frontmatter before the session started.
- [ ] Each question shown included the builder's recommended choice (confirming is the fast path).
- [ ] Answered REQs with all questions resolved flipped to `status: pending` (or `completed` for builder-was-right, `cancelled` for discarded).
- [ ] Approved discovered-task REQs flipped to `pending` and stayed in `do-work/queue/` — not archived.
- [ ] Skipped REQs remained `pending-answers` — nothing lost.
- [ ] `blocked` REQs the user confirmed satisfied flipped to `pending` with `blocked_by`/`blocked_at` removed and a `## Blocked` history line appended; unconfirmed ones stayed `blocked`.
- [ ] The final report names each REQ by id and what happened to it.
