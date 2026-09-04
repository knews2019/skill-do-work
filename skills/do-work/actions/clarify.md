# Clarify Questions Action

> **Part of the do-work skill.** Batch-reviews pending questions from completed work — the user confirms, overrides, or discards builder decisions.

This is the second human-attention window in the pipeline. After actions/work.md processes requests autonomously, any ambiguities the builder encountered are surfaced here as a batch for efficient review.

## When to Use

**Use when:**
- A work run just finished and left `pending-answers` REQs in the queue.
- A work run left `pending-heavy-testing` REQs whose lanes did not run at exhaustion (skipped lane, plan drift, or historical plan) and a maintainer wants to run them by hand.
- The user asks "what's blocked?", "show me pending questions", or similar.
- The pipeline can't advance because builder-decided questions need sign-off.
- `blocked` REQs are waiting on a **human-confirmable** external condition (a designer answered, a service is now up) and you want to confirm the condition is met so they re-enter the queue.

**Do NOT use when:**
- No `pending-answers`, `pending-heavy-testing`, **and** no `blocked` REQs exist — tell the user the queue is clear and stop.
- The user wants to answer a *specific* open question by editing the REQ directly — that's just a file edit, not a batch review.
- The queue has only `pending` REQs — those need `do-work run`, not clarify.
- The `blocked` condition is **machine-checkable** (it carries a `blocked_check` probe) — `do-work run` auto-probes and unblocks those; clarify is for the human-confirmable ones.
- The reply in hand is an **outside stakeholder's** answer set — that is `do-work stakeholder-answers` (`actions/stakeholder-answers.md`). Clarify routes stakeholder-questions REQs (Step 5.5) but never ingests their answers.

## Input

Triggered by `do-work clarify` (also: `answers`, `questions`, `pending`, `what's blocked`). No arguments needed.

## Steps

### Step 1: Scan the queue

Find all `REQ-*.md` files in `do-work/queue/` with `status: pending-answers` or `status: pending-heavy-testing`. Also collect REQs with `status: blocked` (waiting on an external condition) for Step 5.5 — those split into two shapes there by the presence of a `stakeholder:` field (`actions/work-reference.md` → Request File Schema): plain external conditions, and stakeholder-questions REQs.

### Step 2: Check for pending questions

If no `pending-answers`, `pending-heavy-testing`, or `blocked` REQ is found: report "No pending questions, heavy tests, or blocked REQs — queue is clear" and exit. If only `blocked` REQs exist, skip Steps 2.5–5 and go straight to Step 5.5.

### Step 2.5: Run held heavy lanes by hand

List every `pending-heavy-testing` REQ with its title, exact full `commit:` revision, and stored Heavy Verification Plan. Refuse malformed entries whose hash is absent or does not resolve, whose plan is incomplete, or whose target differs from `commit:`; they stay held.

Recompute each stored plan in an isolated detached worktree. An ordinary exact-revision plan uses `plan-heavy-verification` at its stored target. A compatibility plan for already-landed pre-manifest work uses `plan-heavy-revalidation` with every stored historical source range and the exact stored execution revision; it reads the lane manifest from that execution revision and requires every historical target in its ancestry. Refuse any difference from the stored mode, ranges, manifest/execution revision, selected ids, argv, or reasons. REQs in one approved batch may share one lane execution only when their stored execution revision is identical. Run the deduplicated lanes there with the same runner the exhaustion drain uses, which reports each lane's exit status, skip state, and wall seconds:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root "<project-root>" --format json \
  run-heavy-verification --manifest _dev/tests/heavy-lanes.json --lane "<lane-id>"
```

Configure Chromium explicitly for a selected browser lane: a skip is incomplete evidence, never green. An explicit user request for `bash _dev/tests/maintainer-verify.sh --heavy` is force-all and bypasses narrowing.

For each REQ, build one canonical `answer` manifest with `expected_status: pending-heavy-testing`, `mode: heavy-testing`, and typed `heavy_testing` evidence: exact request target, exact execution revision, and every applicable lane's id, argv, exit status, skip state, and wall seconds. Use `confirmed` only when every selected command exits zero with no skip. Green records `heavy_verified_at` and `heavy_verified_revision`, clears `claimed_at`, and returns the REQ to runnable `pending` without archiving it; `next` projects `resume_phase: review`, and the work loop claims normally before fresh review, lessons, release judgment, and canonical finalization. Use `answered` with every failing or skipped lane for red; it clears the same claim and green-evidence fields and returns the REQ to ordinary remediation. Stale, mismatched, or non-ancestor evidence refuses or cannot produce the review resume. Invoke the ordinary `answer --manifest ... --at ... --commit` transaction; there is no hand edit or free-form fallback. The ticked hold line, answer note, and generated Heavy Verification Result are the durable result evidence.

### Step 3: Present questions

**Load `crew-members/clear-questions.md` first** — it is the contract for every question you're about to show. Then **rewrite** each REQ's question to that contract instead of rendering the stored `## Open Questions` text verbatim: that text was authored mid-implementation by a builder with the whole spec in its head, and is presumed too dense for a cold reader. Gloss every coined label or section reference, and state why the decision was escalated to the user (Principle 7). The rewrite applies to the question and option *wording* — the Decision Brief structure below stays as-is.

Always lead with the builder's decision and its default — confirming is the intended fast path — and, for builder-decided follow-ups, show the **value and risk** so the user can judge in seconds instead of spelunking. This is the **DECISIONS FOR YOU** section of the Decision Brief (`actions/work-reference.md` → **Decision Brief (hand-back format)**).

Present each `pending-answers` REQ in three layers — story, then decisions, then detail on request. `pending-heavy-testing` was already handled by Step 2.5 and does not pass through this builder-decision renderer. The story layer is a clarify-local addition: the decision block underneath it is the Decision Brief's format, unchanged.

**Layer 1 — the story (1–4 sentences, always present).** Write for a reader who knows what a UR and a REQ *are* and has forgotten everything about *this* one — the work may be days old. In plain sentences, answer three things: what the user originally asked for, what the builder ran into while building it, and why it couldn't be settled without them. That third part carries Principle 7 of `crew-members/clear-questions.md` — the escalation reason belongs in the story, not wedged into the question.

Source the facts from the REQ's `## What`, the parent UR's request text, the `## Prior Attempt` or review finding that spawned the follow-up, and the escalated decision record's rationale. For a discovered-task REQ, say what was found and what the builder was actually working on when it turned up.

**Hard cap: four sentences, and never longer than the block it introduces.** If the question is self-evident from its header, write one sentence. A story that only restates the question in different words is padding — cut it to the single sentence that adds the fact the question is missing.

**Every question gets a story, not every REQ.** When a REQ holds several questions that all arise from the same situation, one story above the block covers them — that's the common case. But when two questions came from *different* situations (one about how a file is encoded, one about which rows it contains), the shared story can't explain both: give each its own one-sentence lead-in directly above its numbered entry, and let the REQ-level story carry only what they genuinely share. The test is per-question, not per-REQ — a reader must be able to answer question 3 without having inferred anything from question 1.

**Layer 2 — the decision block.** The question, the builder's decision, the default, and the trade-off:

```
REQ-025 — Review fix: dark mode sidebar
(a fix that came out of reviewing the dark-mode work)

The story: You asked for dark mode across the app. While building it the sidebar
turned out to carry its own colour set, separate from the main page, and nothing in
the request said which way it should go. Changing it later is cheap, but picking the
wrong one silently would leave you with a look you never chose.

1. [ ] Should the sidebar use the same dark palette as the main content?
   Decision: Yes, match main content palette   ·   Default if you say nothing: same
   Value: one consistent dark surface; nothing to re-theme later
   Risk:  if you wanted a distinct sidebar, this is a quick CSS revert (low, reversible)
   Also:  Separate sidebar palette, User-configurable

Separately: nothing said whether the choice should survive closing the tab.

2. [ ] Should dark mode be remembered the next time you open the app?
   Decision: Yes, save it in the browser   ·   Default if you say nothing: same
   Value: returning users keep their choice
   Risk:  a stale saved value can hide the follow-your-system-setting path (medium, reversible)
   Also:  Forget on refresh, Follow the system setting

Want the builder's original wording, the files touched, or the full decision record for
either of these? Ask and I'll show it.
```

**Fallback (mandatory).** Many `pending-answers` REQs come from templates that don't carry Value/Risk — capture, verify-requests, review-work follow-ups, and discovered tasks all emit `Recommended:`/`Also:` only. When a question has no `Value:`/`Risk:` lines, render it in that older form (`Recommended:` + `Also:`) — never block on the missing fields. The story layer still applies; it doesn't depend on those fields.

**Layer 3 — detail, on request only.** Close each REQ with one line offering what was held back: the builder's original question wording, the files touched, the parent REQ or UR, and the full decision record. Never render it unprompted — the offer *is* the disclosure mechanism.

**Speakable-first rendering.** Layer 1 has to survive being read aloud end to end:

- Name the request in plain words after the header line; don't repeat bare ids mid-sentence.
- Paraphrase every technical token on first use ("saved in the browser so it survives a refresh" for `localStorage`). Layer 2 may then use the token freely — the constraint is that no term the user needs *in order to choose* appears there cold.
- No file paths, no section marks, no arrows or interpuncts or slashes standing in for words, no CamelCase or snake_case, no abbreviation the user didn't introduce.
- One idea per sentence, and no nested parentheticals — a listener can't see brackets.
- State options as alternatives with their consequence, builder's choice first (it is the default).

Builder-marked `- [~]` decisions reflect the "Think Before Coding" guardrail (`crew-members/coding-guardrails.md`) — surface tradeoffs early, not late.

### Step 4: Collect answers

Collect the user's semantic choices first, then build one strict `answer` manifest per REQ. The outcome bullets below describe manifest content; do not edit the question line, note, frontmatter, or location directly. Supply raw answer text through a regular payload file so the command derives containment from the source bytes.

> **Named entry point — Canonical answered-question format.** The durable record is the `- [x] [question] → [answer]` form below **together with a dated note carrying the reasoning, including anything the answer put out of scope** (`crew-members/clear-questions.md` Principle 8, which defines both halves). The answer line alone is not the record: it says what was chosen and not why, so the next reader re-derives the decision from the stored `Recommended:` rationale — which the answer may have just overruled. Both halves are the durable record for **any** caller that obtains a user answer to a REQ question, not just clarify — an orchestrator that asked and got an answer mid-run writes the same form before dispatch (`actions/work.md` Step 3.5). Callers cite it by this name, not by step number. What is clarify-local is Step 5's `pending-answers` → `pending` flip: a REQ already in flight has no such status to leave.
>
> **Named sub-contract — Outside-text containment.** Part of this format is a wider write rule: text from outside the file is neutralized before it is written into **any do-work Markdown record body**, whether that record is a REQ answer, a UR's Full Verbatim Input, an addendum, or a cancellation note. Those sites are illustrative, not a writer checklist; the condition is what gives a future body writer the rule. First require text the record can carry byte-identically: if it contains any **C0 control or DEL except LF and TAB**, refuse the write and report the offending character — do not normalize it, drop it, or invent a hand-authored escape table. Once that preflight passes, the text stays **byte-identical apart from containment bytes** added around it.
>
> Then judge every line against one condition: **could this line be read as one of this file's own delimiters?** The shapes already proved are illustrative, not a checklist — a line-leading `- [ ]` or `- [x]` becomes a real open question that Step 5 then pins the REQ on, a bare `---` line reads as a frontmatter fence, a `## ` line invents a section every prose grep obeys, an unbalanced code fence swallows every section below it on the board while those same greps still see them. The destination decides which containment branch applies:
>
> - **The answer is one line** and it cannot be a delimiter where it lands → write it inline after the `→`. Position neutralizes it: a delimiter has to start a line, and this text never does.
> - **Anything written as its own body passage, or any answer with a line break or delimiter-shaped line** → put the outside text inside a blockquote whose lines open a code fence longer than the longest backtick run anywhere in the text. For an answered question, keep the one-line answer summary after the `→` and put this contained text in the dated note. The `> ` prefix takes every line start away from a line-based scan; the longer fence keeps those same lines literal for a Markdown reader and cannot be closed from inside. Prefix every physical line, including blank lines and the two fence lines. Only the prefix and fence are added — nothing is edited or dropped.
>
> Placement answers the same condition: contained outside text lands inside the body section that owns it; an answer touches its own question line and the note beneath it. It is never the file's first line — the frontmatter's opening fence is a delimiter too, and one byte written above it empties `status`, `title` and `user_request` together.

If your environment has a structured question prompt (multi-question UI), batch questions in groups of **at most 4 per prompt** — chunk by question count, not by REQ. A REQ with 6 questions needs 2 prompts.

Before writing answers, initialize one deduplicated `overturned_decision_sources` set for this clarify session. A source enters it only when its REQ carries `builder_decided: true` and the user's answer is semantically different from the stored `Recommended:` builder choice. Preserve the recommendation before replacing the question line so the comparison and later decision-revalidation source both retain old → new provenance. Confirmation, discard, discovery approval, and an ambiguous same-choice paraphrase never enter the set.

**A per-question verb records that question's outcome and nothing else. It never sets the REQ's status and never archives.** A REQ can hold several questions and the user can answer them differently — discard one and skip the next — so a verb that flipped whole-file state would make two branches unfollowable at once, and could archive a REQ still holding an open `- [ ]`. Step 5 computes the REQ's status once, from the whole set of outcomes.

Every answered line is written with its dated note (**Canonical answered-question format**, above): the answer, then the reasoning and anything the answer put out of scope. Obtain the date per the Timestamp rule's date-only paragraph (`actions/work-reference.md`) — cite it, never spell a clock command here.

For each question, the user can:

- **Answer it** → update to `- [x] [question] → [user's answer]`, plus the dated note; when this is a `builder_decided: true` follow-up and the answer differs from `Recommended:`, add its REQ id to `overturned_decision_sources`
- **Confirm builder's choice** → update to `- [x] [question] → Confirmed: [builder's choice]`, plus the dated note
- **Pick a different option** → update to `- [x] [question] → [user's chosen option]`, plus the dated note; when this is a `builder_decided: true` follow-up, add its REQ id to `overturned_decision_sources`
- **Skip for now** → leave as `- [ ]`
- **Discard it** → update to `- [x] [question] → Discarded`, plus the dated note saying why

### Step 5: Resolve each REQ's status from its whole question set

Step 4 wrote answers; this step decides what each REQ now *is*. Compute it **once per REQ, from every question's outcome together** — never from the last verb the user happened to use.

**Any remaining `- [ ]` wins.** If even one question is still open, the REQ stays `pending-answers` and stays in `do-work/queue/`, whatever the other answers were. A REQ is never archived while it holds an unanswered question: the answers already given are recorded and the REQ comes back to the next clarify session for the rest.

Otherwise every question is `[x]` or `[~]`, and the REQ's disposition follows from the answers:

- **Every question discarded** → `status: cancelled`, `completed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`), archive directly. The honest won't-do status: no work happened and none is wanted — see "Discarded" below.
- **Otherwise, a `builder_decided: true` follow-up whose questions were all confirmed** → `status: completed` (no implementation needed — see "Builder Was Right" below). **The marker is the entire test.** It is stamped at creation by the Builder-Decided Follow-up Template (`actions/work-reference.md` → Step 8) and by nothing else, so routing cannot be changed by rewording a question. Never infer this branch from question prose: a follow-up that still needs building, confirmed down this path, is archived `completed` without anyone building it.
- **Everything else** → flip `status` from `pending-answers` to `pending` so the work loop picks it up — an approved discovered task, an answered review follow-up, a partially-discarded REQ with real answers left in it. See "Approved Discovered Task" below.

On any flip, stamp `status_changed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`; the board's state timer reads it, so the card shows time since the answers landed rather than time since capture). REQs flipped to `pending` enter the queue for the next `do-work run`.

Pass the whole question set, expected status/path, canonical timestamp, and any terminal archive/UR-closure evidence to one invocation:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root "<project-root>" --format json answer --manifest "<manifest-path>" --at "<answer-timestamp>" --commit
```

The command uniquely matches each plain or Q-ID question, writes every answer and dated note, derives the disposition from the whole record, and couples any archive/UR move atomically. A missing or refused command leaves the REQ byte-identical and stops this REQ. There is no hand-edit, helper, or manual Git fallback.

### Step 5.25: Revalidate queued work after reversals

If `overturned_decision_sources` is empty, skip this step. Otherwise invoke the named **Decision Revalidation Workflow** in `actions/verify-requests.md` once with every source id in the set — equivalent to one `--against REQ-NNN` pair per source. The sources were activated in Step 5 and remain in `do-work/queue/`, but that workflow excludes source follow-ups from their own candidate scan.

Use that workflow's inventory before reading queue bodies:

- At **10,000 queued words or fewer**, run the full semantic scan automatically.
- Above **10,000 queued words**, show the queued-file count, word count, approximate 1.3–1.6-tokens-per-word input range, and any claimed REQ ids excluded from v1. Ask one choice: `Run the decision-revalidation scan now?` The `crew-members/clear-questions.md` contract already loaded for this session governs that cost question too.
- If the user declines, finish clarify without scanning and print one copyable combined command with repeated flags, preserving set order: `do-work verify-requests --against REQ-NNN --against REQ-MMM`.

The threshold applies only to this automatic clarify trigger. An explicit `do-work verify-requests --against ...` invocation displays its estimate and proceeds. Revalidation itself is read-only: aside from clarify's answer/status writes already completed above, it changes no REQ body, frontmatter, status, or location.

### Step 5.5: Confirm blocked conditions

**Stakeholder-questions REQs first (`stakeholder:` present) — routed, never yes/no-confirmed or ingested here.** Before printing a REQ's line, check its `blocked_by:` — when it says the report is pending regeneration, or names a bundle path that no longer exists, regenerate now: follow `../../do-work-toolbox/actions/stakeholder-report.md`, then update `blocked_by:` and append the `## Reports` history line in the same edit (the standing regeneration condition, `actions/stakeholder-answers.md` Step 5). `blocked_by:` carries the person's own name, so the rewrite follows the **Frontmatter Quoting** contract (`actions/work-reference.md` → Request File Schema) — the frontmatter sibling of Step 4's neutralization branch. Then print the routing summary, one line per REQ, so a question mis-routed to a stakeholder is visible and recoverable:

```
Stakeholder questions open:
  REQ-NNN — 4 questions for Priya (design), 1 irreversible, open 3d — report: <latest bundle path from blocked_by>
```

Then offer per REQ:

1. **I have their reply** → hand off to `do-work stakeholder-answers REQ-NNN` (`actions/stakeholder-answers.md` — it loads the prompt-injection guardrail before reading the paste and routes each answer by its printed Q-ID; clarify never ingests a stakeholder reply).
2. **A question is really mine** → reclaim it: mark that line `- [x] **Q-NN** — [question] → Reclaimed by user via clarify — [the user's inline answer | moved to REQ-MMM]`. An inline answer is handled exactly like a stakeholder answer — the delta against `Assumed:` may mint a change REQ (`actions/stakeholder-answers.md` Step 4); without one, mint or fold a `pending-answers` follow-up carrying the question per the **Builder-Decided Follow-up Template (Step 8)** (`actions/work-reference.md`) so it re-enters the user's own clarify flow as `REQ-MMM`. If the reclaim resolves the REQ's last open question, apply the **Stakeholder REQ terminal semantics** (`actions/stakeholder-answers.md` Step 5); while open questions remain, report regeneration follows the same standing condition stated there.
3. **Leave it** → unchanged.
4. **Give up on the answers** → `do-work abandon REQ-NNN`.

**Plain blocked REQs** (no `stakeholder:`): present each condition as one lightweight yes/no — no rewrite-contract machinery needed (the condition is a single line of `blocked_by` text, not a builder question). It still gets one sentence of story, because a REQ that has been sitting blocked for days is exactly the one the user no longer remembers:

```
REQ-042 — Wire up local translation
What it was for: you wanted translation to run on your own machine instead of a paid API.
Blocked by: LM Studio running locally (waiting 3 days)
Is this condition now satisfied?
  1. Yes — unblock it        2. Not yet — leave it        3. Abandon this REQ
```

Note for the user which blocked REQs carry a `blocked_check` probe — those unblock automatically on the next `do-work run`, so confirming them by hand here is optional. Present only the human-confirmable ones prominently.

- **Yes → unblock:** invoke `unblock REQ-NNN --request-path <exact-path> --original-status blocked --source user-via-clarify --confirmed`. The canonical transaction validates the stale preimage, stamps `pending`/`status_changed_at`, removes `blocked_by`/`blocked_at`, retains `blocked_check`, and appends the user-cleared history. A refusal leaves bytes unchanged; do not reproduce the mutation free-form. The REQ re-enters the queue for the next `do-work run`.
- **Not yet:** leave it `blocked`, unchanged.
- **Abandon:** hand off to `do-work abandon REQ-NNN` (marks `cancelled`, archives) — same as discarding a question.

### Step 6: Report

Summary of what was resolved and what's still pending — include every heavy-test REQ returned to `pending` for review-resume or remediation, or left held; any `blocked` REQs unblocked (now `pending`) or left waiting; and any stakeholder-questions REQs handed off, reclaimed from, or left open, alongside the answered/confirmed/discarded questions. When Step 5.25 ran, append its evidence-backed candidate report. When the user deferred an over-threshold scan, say it was not run and include the combined explicit command.

## Builder Was Right / Discarded

When the user reviews a `pending-answers` follow-up and confirms that the builder's original choice was correct (i.e., no implementation change needed):

1. Update the question to `- [x] [question] → Confirmed: [builder's choice]`
2. Update frontmatter: `status: completed`, `completed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`)
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

When a REQ comes out of Step 5 on the flip-to-`pending` branch — the ordinary case for a discovered-task follow-up, which carries no `builder_decided: true` marker — and the user confirmed the recommendation:

1. Update the question to `- [x] [question] → Confirmed: Yes, add to queue`
2. Update frontmatter: `status: pending` (NOT `completed` — this task needs to be built), plus `status_changed_at: <timestamp>` (current UTC instant — Timestamp rule, `actions/work-reference.md`)
3. **Do not archive.** The REQ stays in `do-work/queue/` and enters the normal work queue for the next `do-work run`

This is distinct from "Builder Was Right" because confirming a discovered task means the user wants it *executed*, not signed off. The task has no prior implementation to confirm — it's a new piece of work that needs a full work cycle.

## Rules

- This action avoids wasting a work cycle on a REQ that just needs sign-off or rejection, while correctly routing approved discovered tasks into the build queue
- Never block the user — if they skip all questions, exit gracefully
- Always show the builder's recommended choice prominently so confirming is the fast path
- The story exists to make the question answerable cold — it is not a summary of the work. If it adds no fact the question itself lacks, it's padding; cut it back
- Batch every overturned builder decision into one revalidation scan; never reread the whole queue once per answer

## Red Flags

- A question shown to the user still contains unglossed builder shorthand (a coined label, a spec §-reference, a finding number) — the stored text was rendered verbatim instead of rewritten per `crew-members/clear-questions.md`.
- A story that restates the question in different words — it must add what was asked for, what the builder hit, or why the call is the user's.
- A story longer than the decision block it introduces.
- A file path, bare identifier, or unglossed technical token inside a story — unreadable aloud and unanswerable cold.
- Layer 3 detail rendered unprompted — the layering bought nothing, and the user is back to scanning.
- A REQ-level story that only explains the first of several questions — the later ones are being answered on inference.
- A question presented here is a fragment of an answer recorded earlier — that answer went in as live markup instead of the format's contained form, and the REQ is now pinned on text nobody asked.
- A `pending-answers` REQ with no `## Open Questions` section — the marker and the body disagree; investigate before presenting nothing.
- User confirms every builder choice without reading — they may be rubber-stamping; ask once if they want a summary first.
- A discovered-task follow-up's `status` flipped to `completed` instead of `pending` after user confirmed "Yes, add to queue" — that's the wrong route (the task never gets built).
- A confirmed, discarded, or discovered-task choice triggered decision revalidation — no decision was reversed.
- Revalidation changed a candidate REQ or moved it to `pending-answers` — the sweep is evidence-only in v1.
- `pending-answers` REQs pile up across multiple clarify runs without resolution — users are skipping; ask whether to discard the stale ones.

## Verification Checklist

- [ ] `crew-members/clear-questions.md` was loaded before the first question was presented, and stored question text was rewritten to its contract (not rendered verbatim).
- [ ] Every REQ presented had `status: pending-answers` in its frontmatter before the session started.
- [ ] Each question shown included the builder's recommended choice (confirming is the fast path).
- [ ] Every presented REQ carried a 1–4 sentence story naming what was asked for, what the builder hit, and why the call is the user's.
- [ ] Every *question* was answerable from its own story — questions arising from a different situation than their REQ-mates got their own one-sentence lead-in.
- [ ] No story contained a file path, bare identifier, or unglossed technical token.
- [ ] Layer 3 detail was offered, not rendered, unless the user asked for it.
- [ ] Each `blocked` REQ presented in Step 5.5 carried a one-sentence "what it was for" line.
- [ ] Every answer written into a REQ used the **Canonical answered-question format** — the answer line **and** its dated note carrying the reasoning and anything the answer put out of scope. The same form any mid-run caller must use, so the decision survives the session that heard it rather than being re-derived from `Recommended:`. Its neutralization branch applies too: no answer's text landed as live markup.
- [ ] Every dated note **cited** the Timestamp rule's date-only paragraph for its date. No clock command was spelled anywhere in this file.
- [ ] No per-question verb set a REQ-level status or archived anything. Status was computed once per REQ, in Step 5, from every question's outcome together.
- [ ] **A REQ holding even one remaining `- [ ]` stayed `pending-answers` in `do-work/queue/`** — including one where another question was discarded. Nothing with an open question was archived.
- [ ] REQs with every question resolved flipped to `status: pending` — except all-discarded ones (`cancelled`) and confirmed `builder_decided: true` follow-ups (`completed`).
- [ ] The `completed` fast path was taken **only** on the `builder_decided: true` marker, never inferred from a question's wording.
- [ ] Approved discovered-task REQs flipped to `pending` and stayed in `do-work/queue/` — not archived.
- [ ] Skipped REQs remained `pending-answers` — nothing lost.
- [ ] Every genuinely different answer on a `builder_decided: true` follow-up entered one deduplicated reversal set; confirmations, discards, and discovered-task approvals did not.
- [ ] Reversal sources shared one queue scan; a scan above 10,000 queued words required confirmation, and a declined scan emitted one combined explicit command.
- [ ] Decision revalidation itself changed no candidate REQ or queue status.
- [ ] `blocked` REQs the user confirmed satisfied flipped to `pending` with `blocked_by`/`blocked_at` removed and a `## Blocked` history line appended; unconfirmed ones stayed `blocked`.
- [ ] Stakeholder-questions REQs were routed or reclaimed, never yes/no-confirmed and never ingested here — no stakeholder answer text was written by clarify.
- [ ] The final report names each REQ by id and what happened to it.
