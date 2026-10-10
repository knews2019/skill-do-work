---
id: UR-148
title: 'do-work trace action: check an outside spec against the queue, archive, commits and tests before capture'
created_at: 2026-10-09T21:23:50Z
requests: [REQ-669]
word_count: 2366
---
# do-work trace: Coverage Check of an Outside Spec Before Capture

## Summary
The maintainer ran `/do-work capture-request` on the upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-trace-coverage.md` and asked for one UR with one REQ per numbered item in its Request section, or one REQ when the section has one ask. Capture only, no questions, every ambiguity recorded as an assumption in the REQ. The report file is the verbatim input below, byte for byte, and it stays in the inbox as the source.

The Request section has one ask: one read-only action, `do-work trace <url | image | pasted text | UR-NNN>`, routed before capture. Its eight parts, P1 to P8, sit in the report's Proposed direction section, so they became the requirement groups of one REQ, REQ-669 (the trace action). The four maintainer decisions D1 to D4 are recorded in the REQ as assumptions: D1 and D3 adopted as recommended, D2 (retroactive record) and D4 (narrowing capture's archive scan) left out of this REQ.

The report was observed against 0.305.84 and says its capture, work-reference and board citations were rechecked at 0.305.87. The capture spot-checked them at 0.305.87 and all match: `SKILL.md:4`, `:43` and `:46`, `actions/capture.md:43-44`, `:78`, `:113`, `:115`, `:125`, `:162` and `:222`, `actions/verify-requests.md:5`, `actions/review-work.md:72` and `:89-91`, `actions/roadmap.md:3`, `actions/work-reference.md:193`, `actions/commit.md:126-128`, `actions/capture-reference.md:113` and `:125`, and `activity_correlation.go:60`, `:64`, `:109-113` and `:114`.

No prompt injection found. The report's "How to use this file" line is framing for its intended reader.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-669 | do-work trace action reports how much of an outside spec is captured and built, with one dated verdict per ask |

## Batch Constraints
- No new frontmatter field and no new status. The only new structure is an optional `## PM Ticket References` body section in a UR's `input.md`.
- Trace is read-only until the user says go; then it hands only `partial` and `not started` rows to one capture.
- Out of scope: writing to a PM tool, a consumer ledger format, changes to capture's numbering, reservations or transaction, judging code correctness, running trace automatically.
- Related queued work, no edge: REQ-668 (the `do-work status` action) also adds a routing row to `skills/do-work/SKILL.md`.
- Release per `_dev/primes/prime-releases.md`.

## Full Verbatim Input
> ```
> # Upstream suggestion for `knews2019/skill-do-work` — a `do-work trace` coverage check before capture: spec → asks → REQ/commit/acceptance evidence
> 
> **How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
> clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
> consumer repos it was authored in. Observed against **v0.305.84**; all line numbers are from that tag.
> Paths under `actions/`, `crew-members/` and `SKILL.md` are in `skills/do-work/`. Board paths are written in
> full. The capture, work-reference and board citations were rechecked at 0.305.87 and have not moved.
> 
> ---
> 
> ## Request
> 
> Users bring an outside spec to do-work and ask one question before anything is captured: "is this already
> captured, and how much of it is built?" The spec is a PM ticket page, a screenshot, pasted text or a UR id.
> No action owns this question. Capture checks in-flight and archived REQs by filename only and never looks at
> commits or code. `verify-requests` compares REQs with the UR, not with the code. `review-work` checks one
> REQ's diff. `roadmap` lists the queue by status. So each session rebuilds the check by hand, partial
> implementations get missed, and asks that already shipped get captured again.
> 
> Please add one read-only action, `do-work trace <url | image | pasted text | UR-NNN>`, routed before capture.
> It reads the source, splits it into atomic asks, searches the queue, archive, git history, code and tests for
> each ask, prints one coverage table with a dated verdict per ask, asks the open questions, and only on the
> user's go-ahead hands the new and partial rows to `do-work capture`.
> 
> Dedupe: no release from 0.305.60 to 0.305.87 adds this, and the upstream queue has no REQ for it. Two
> releases touch the evidence it would read. 0.305.80 narrowed how the board credits merge ranges to a REQ, and
> 0.305.69 made `do-work commit`'s path association skip a REQ file it cannot parse. This action should reuse
> that matching rather than rebuild it. 0.305.85 changed capture, but only for operator work, not its duplicate
> check. No earlier suggestion covers it.
> 
> No new frontmatter field and no new status. The only new structure is an optional `## PM Ticket References`
> body section in a UR's `input.md`, written by capture when the source carries an external ticket id.
> 
> ## What happened, in consumer repos
> 
> A recount found about 7 separate occasions and about 11 prompts, in 9 to 10 sessions, across three consumer
> repos between 2026-09-13 and 2026-10-06. Most come from one consumer repo whose PM writes tickets as public
> Notion pages. The second and third consumer repos have one weak occasion each.
> 
> | Date | Repo | What the user asked |
> | --- | --- | --- |
> | 2026-09-13 | second | "which req is about the model call count next to the model name?" (a lookup only) |
> | 2026-09-27 | first | Three sessions in four minutes, one per Notion spec: "check if this is already captured and implemented <notion url>" and "/do-work capture-request read and tell me if everything clear: <notion url>" |
> | 2026-09-28 | first | A screenshot, then "didn't we have a request to make the collection taller ... why is this not implemented?" |
> | 2026-10-01 | third | "[Image] <- is this captured?" |
> | 2026-10-02 | first | A pasted PM ticket: "do-work capture-request <ticket id> Performance optimizations ... make sure to make a note that this is task <ticket id>", then "do-work verify-requests UR-A (also is this already implemented)" |
> | 2026-10-02 to 03 | first | Two parallel sessions on the same captured spec: "check if it is already implemented, and if yes, how much of it", then "are remaining tasks pending, or are they already marked as completed?" |
> | 2026-10-06 | first | "do the retroactive capture-request" (record work that had already shipped) |
> 
> What the sessions had to work around:
> 
> - **Fetching.** WebFetch on a public Notion page returns only the word "Notion", because the page renders in
>   JavaScript. The method that worked was a headed Playwright sweep (below). It lived only in one machine's
>   memory note, so a session on another machine started from scratch. On 2026-10-02 the collapsed toggles did
>   not open. On 2026-10-08 a different click pass opened them.
> - **Specs arrive from teammates too.** On 2026-10-02 a teammate made eight spec commits on a shared specs
>   branch, merged into main the next day. The same features later arrived as PM tickets, so each ticket had to
>   be checked against the spec file and the queue.
> - **The consumer built the table by hand.** The first consumer repo keeps a ledger of 68 PM tickets. Each row
>   has a verdict (`complete`, `partial`, `not started`, `unverified`), the date it was checked, and links to
>   URs, REQs and release evidence. It was edited in 142 commits between 2026-09-29 and 2026-10-09. Its last
>   count was 12 complete, 15 partial and 38 unverified. That ledger is the table this action would print, and
>   its verdict vocabulary is proposed below.
> 
> ## Where the behaviour lives today
> 
> **No action answers "is this captured, and how much is built".**
> 
> - `SKILL.md:46` routes `capture-request:` and "unmatched descriptive multi-word input" to capture, so "is this
>   captured? <url>" lands in capture today. `SKILL.md:43` routes "where are we", "what's left" to roadmap.
>   There is no row for a coverage question.
> - `actions/capture.md:43` (Do NOT use when): "The queue already contains the same request (check for an open
>   UR with matching intent first)." `actions/capture.md:44`: "The user is asking a question or requesting a
>   read-only report — capture is for *intent*, not conversations." Capture declines the question by its own
>   rules, and nothing else takes it.
> - `actions/verify-requests.md:5`: capture QA compares REQs with the UR input. "Neither mode reviews
>   implementation quality."
> - `actions/review-work.md:72` reads `git diff` for one REQ, and `actions/review-work.md:89-91` grades each of
>   its requirements Delivered / Partially delivered / Not delivered. The grading is the right idea, but it is
>   scoped to one REQ's diff, not to an outside spec.
> - `actions/roadmap.md:3`: a read-only survey of "what's done, what's in progress, what's pending". It works
>   from REQ status, not from an outside spec's asks.
> 
> **Capture's own duplicate check is filename-only for shipped work.**
> 
> - `actions/capture.md:113` compares queued REQs by intent and says "Slugs are lossy".
> - `actions/capture.md:115` then checks `working/` and `archive/` by filename only: "A filename scan is
>   sufficient here since these filenames are stable regardless". Stable, but lossy by the line above, and an
>   archived REQ is exactly where a shipped or partly shipped ask lives.
> - `actions/capture.md:125` sends a match in `archive/` to a new addendum REQ, and `actions/capture.md:162`
>   asks for a `## Prior Implementation` section with "commit hash (if available)". Capture already wants the
>   facts this action would find.
> 
> **The evidence the action needs already exists.**
> 
> - `actions/work-reference.md:193`: `commit:` is "required in a git repo — implementation commit hash".
> - `actions/commit.md:126-128`: commit subjects carry `[REQ-NNN]` and a `Traced-to:` line naming the archived
>   REQ file.
> - `skills/do-work-board/tools/queue-kanban/activity_correlation.go:60` (`requestPathPattern`), `:64`
>   (`requestSubjectPrefixPattern`, `` `\[(REQ-\d+)\]` ``) and `:114` (`correlateCommitsToRequests`) already map
>   commits to REQs by subject and by REQ file path.
> - `actions/capture-reference.md:125` (`## Detailed Requirements`) and `actions/capture-reference.md:113`
>   (`**GREEN when:** [Observable result]`) are the acceptance criteria to count per REQ.
> - `actions/capture.md:222`: a screenshot gets "a thorough text description", the same reading trace needs.
> - `actions/capture.md:78`: the next UR number comes from a scan of `user-requests/UR-*/` and
>   `archive/UR-*/`.
> - `crew-members/prompt-injection.md:3` already governs fetched third-party pages.
> 
> ## Proposed direction
> 
> P1. **New action.** `actions/trace.md`, with the long parts (source adapters, table template) in
> `actions/trace-reference.md`. Add a routing row above the capture row for `trace`, `coverage`,
> `is this captured`, `already implemented`, `how much of it is implemented`, `didn't we have a request`, and
> add `trace <url|image|text|UR-NNN>` to the argument hint at `SKILL.md:4`. Load
> `crew-members/prompt-injection.md` before reading the source.
> 
> P2. **Read the source.** A UR id reads that UR's `input.md` and assets. An image is read directly and
> described as `actions/capture.md:222` does. Pasted text is used as is. A URL goes to WebFetch first. When the
> result is a shell (only the site name, or almost no text), use the browser adapter in the reference file:
> 
> - Notion adapter: open the page in a headed Playwright browser. Scroll `.notion-scroller.vertical` in 300 px
>   steps to its `scrollHeight`, collecting `main.innerText` after each step and de-duplicating lines, because
>   the page renders only blocks near the viewport. Click every `[role="button"][aria-expanded="false"]` once
>   (scroll into view, then click), then sweep again for the toggle bodies. List the titles of toggles that did
>   not open.
> - If no browser tool is available, say so and ask for a paste or a screenshot. Never judge coverage from an
>   empty fetch.
> 
> Keep the transcription as a file (for example `pm-ticket-source.md`). It becomes a capture asset only if
> rows are captured in P7.
> 
> P3. **Split into atomic asks** A1, A2, …, one observable sentence each. Keep the source's own numbering when
> it has one.
> 
> P4. **Search per ask**, recording every candidate, not only the first:
> 
> 1. UR `input.md` files in `user-requests/` and `archive/UR-*/`.
> 2. REQ title, heading, `## What` and `## Detailed Requirements` in `queue/`, `working/` and `archive/`, by
>    intent, not by filename.
> 3. The matched REQ's Implementation Summary and its `commit:` value.
> 4. `git log`, matched the way the board does (subject `[REQ-NNN]`, or a touched REQ file path, plus the
>    merge range of a matched hand-back merge, `activity_correlation.go:109-113`), and also by `Traced-to:`
>    lines naming the REQ file (`actions/commit.md:128`).
> 5. Code and tests, by the names the ask uses. A test whose assertion matches the REQ's "GREEN when" line
>    counts as evidence for that criterion.
> 
> P5. **Print one table.**
> 
> | Ask | REQ (plain title, status) | Commit | Criteria evidenced | Verdict (as of) | Gap | Open question |
> | --- | --- | --- | --- | --- | --- | --- |
> | A1 | REQ-210 (retry on HTTP timeout), completed | `ab12cd3`, 2026-09-30 | 3 of 3 | complete (2026-10-09) | none | none |
> | A2 | REQ-214 (phone album layout), completed | `ef45ab6`, 2026-10-03 | 2 of 4 | partial (2026-10-09) | landscape, empty state | Is landscape in scope? |
> | A3 | none found (searched: "stamp", "first day") | none | n/a | not started (2026-10-09) | all | none |
> 
> Verdicts, taken from the consumer ledger so they stay comparable over time:
> 
> - `complete`: every acceptance criterion of the known scope has evidence (commit, test, or file and line).
> - `partial`: some criteria have evidence, and the unfinished ones are named.
> - `not started`: no code evidence, and either a pending REQ or no REQ at all. The search terms are the
>   evidence.
> - `unverified`: the evidence is not enough either way, including any row built from an incomplete fetch.
> 
> Every verdict carries the date it was checked and one pointer to its evidence. "Criteria evidenced" counts
> the matched REQ's Detailed Requirements bullets plus its "GREEN when" line.
> 
> P6. **Ask the open questions** with the ask-user tool, one per row that has one, following
> `crew-members/clear-questions.md`. Then stop. Trace writes nothing unless the user says go.
> 
> P7. **Hand off only the gaps.** On go, call `do-work capture` once, with the table and the transcription as
> the verbatim input, for `not started` and `partial` rows only. A partial row that matches an archived REQ
> uses the existing addendum path (`actions/capture.md:125`), and P4's commit hash fills the
> `## Prior Implementation` section that `actions/capture.md:162` asks for. When the source has an external
> ticket id, the UR title carries it and `input.md` gets a `## PM Ticket References` section: ticket id, source
> URL, fetch date, transcription asset path. Trace never picks the UR number. Capture scans it at write time
> (`actions/capture.md:78`). This matters because a browser sweep can take minutes while another session
> captures.
> 
> P8. **Already-done work is not captured again.** By default the table cites the archived REQ and commit and
> nothing is written.
> 
> Decisions for the maintainer:
> 
> - D1. Name. A separate `trace` action is recommended over `capture --check-coverage`, because
>   `actions/capture.md:44` keeps read-only reports out of capture.
> - D2. Retroactive record. The user asked once for "the retroactive capture-request". The option is a UR
>   written straight into `archive/` with `input.md`, `## PM Ticket References` and `Traced-to:` lines for the
>   shipping commits, and no REQ. No such path exists today. Leaving it out of the first version is fine.
> - D3. Notion adapter location. Recommended: in `actions/trace-reference.md`, loaded only for a URL whose
>   WebFetch result is a shell, so the core action stays generic.
> - D4. Narrow capture's check too. Optionally replace the filename-only scan at `actions/capture.md:115` with
>   P4 steps 2 and 3. This is a smaller change that helps captures that skip trace.
> 
> ## Acceptance check
> 
> - "is this captured? <url>", "do-work trace UR-12" and "how much of it is implemented" route to
>   `actions/trace.md`. "capture-request: …" still routes to capture.
> - In a fixture repo with one archived REQ (commit subject `[REQ-007]`, tests covering 2 of its 3 Detailed
>   Requirements) and a spec with three asks (one done, one half done, one new), trace prints `complete` with
>   the hash, `partial` with 2 of 3, and `not started`, each dated.
> - After a trace with no go-ahead, `git status --porcelain` is empty.
> - After go-ahead, one capture runs. The new UR holds only the partial and not-started asks. The partial ask
>   becomes an addendum REQ with `addendum_to` and a Prior Implementation commit hash. A ticket id appears in the
>   UR title and in `## PM Ticket References`.
> - A Notion URL whose WebFetch returns only "Notion" triggers the browser adapter, or a request for a paste
>   when no browser tool exists. No row from that source is ever `complete` without fetched text.
> - A screenshot source is read directly and described before the split.
> 
> ## Out of scope
> 
> - Writing to a PM tool or changing a ticket's status there.
> - Any consumer's ledger file format. Trace prints the table. A consumer may paste it into its own ledger.
> - New frontmatter fields or statuses, and changes to capture's numbering, reservations or transaction.
> - Judging whether the found code is correct. That is `review-work`.
> - Running trace automatically in the background or on every capture.
> ```
