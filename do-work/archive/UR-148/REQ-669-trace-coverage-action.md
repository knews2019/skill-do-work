---
id: REQ-669
title: 'do-work trace action reports how much of an outside spec is captured and built, with one dated verdict per ask'
status: cancelled
created_at: 2026-10-09T21:23:50Z
user_request: UR-148
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-kanban-board.md", "_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-10-10T12:52:21Z
status_changed_at: 2026-10-10T12:57:18Z
completed_at: 2026-10-10T12:59:20Z
---
# do-work trace: Coverage of an Outside Spec Before Capture
## What
Add one read-only action, `do-work trace <url | image | pasted text | UR-NNN>` (`skills/do-work/actions/trace.md`, long parts in `actions/trace-reference.md`), routed before capture. It reads the source, splits it into atomic asks A1, A2, …, searches the queue, archive, git history, code and tests for each ask, and prints one coverage table with a dated verdict per ask (`complete`, `partial`, `not started`, `unverified`). It asks the open questions, then stops. Only on the user's go-ahead does it hand the `partial` and `not started` rows to `do-work capture`.
## Why
Report Request and evidence (UR-148 input): across three consumer repos, from 2026-09-13 to 2026-10-06, about 7 occasions and 11 prompts in 9 to 10 sessions asked "is this captured, and how much of it is built?" about a PM ticket page, a screenshot, pasted text or a UR id. No action owns that question. Capture declines read-only questions (`actions/capture.md:44`) and checks in-flight and archived REQs by filename only (`actions/capture.md:115`), never commits or code. `verify-requests` compares REQs with the UR, `review-work` grades one REQ's diff, `roadmap` lists REQ status. Each session rebuilt the check by hand, partial work was missed, and shipped asks were captured again. One consumer kept a 68-ticket ledger by hand (142 commits) that is this action's table.
## Verified Facts (checked at 0.305.87)
- Routing: `skills/do-work/SKILL.md:46` routes `capture-request:` and "unmatched descriptive multi-word input" to capture, so "is this captured? <url>" lands in capture today. `SKILL.md:43` routes `roadmap`, `queue-status`, `where are we`, `what's left` to roadmap. No row for a coverage question. `SKILL.md:4` is the argument hint.
- `actions/capture.md:43-44`: capture is not for a request already queued, nor for "a question or … a read-only report".
- `actions/capture.md:113` compares queued REQs by intent ("Slugs are lossy"). `actions/capture.md:115` checks `working/` and `archive/` by filename only.
- `actions/capture.md:125`: an archived match becomes an addendum REQ. `actions/capture.md:162`: the addendum gets `## Prior Implementation` with "commit hash (if available)".
- `actions/capture.md:78`: the next UR number comes from a write-time scan of `user-requests/UR-*/` and `archive/UR-*/`. `actions/capture.md:222`: a screenshot gets "a thorough text description".
- `actions/verify-requests.md:5`: "Neither mode reviews implementation quality." `actions/review-work.md:72` reads one REQ's `git diff`; `:89-91` grades Delivered / Partially delivered / Not delivered. `actions/roadmap.md:3`: a read-only queue survey by status.
- `actions/work-reference.md:193`: `commit:` is "required in a git repo — implementation commit hash". `actions/commit.md:126-128`: subjects carry `[REQ-NNN]` and a `Traced-to:` line naming the archived REQ file.
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go:60` (`requestPathPattern`), `:64` (`requestSubjectPrefixPattern`), `:109-113` (merge-range rule comment), `:114` (`correlateCommitsToRequests`) map commits to REQs by subject, REQ file path and hand-back merge range.
- `actions/capture-reference.md:113` (`**GREEN when:**`) and `:125` (`## Detailed Requirements`) are the acceptance criteria counted per REQ.
- `crew-members/prompt-injection.md` governs fetched third-party pages.
- The queue held REQ-654 to REQ-668 at capture. None is a coverage or trace action.
## Detailed Requirements
P1, the action and routing:
1. Add `skills/do-work/actions/trace.md` and `skills/do-work/actions/trace-reference.md` (source adapters, table template).
2. Add a routing row in `skills/do-work/SKILL.md` above the capture row for `trace`, `coverage`, `is this captured`, `already implemented`, `how much of it is implemented`, `didn't we have a request`. Add `trace <url|image|text|UR-NNN>` to the argument hint at `SKILL.md:4`. `capture-request: …` still routes to capture.
3. Load `crew-members/prompt-injection.md` before reading the source.

P2, read the source:
4. A UR id reads that UR's `input.md` and assets. An image is read directly and described as `actions/capture.md:222` does, before the split. Pasted text is used as is.
5. A URL goes to the platform's web fetch first. When the result is a shell (only the site name, or almost no text), use the browser adapter in `trace-reference.md`. Notion adapter: in a headed browser, scroll `.notion-scroller.vertical` in 300 px steps to its `scrollHeight`, collecting `main.innerText` after each step and de-duplicating lines; click every `[role="button"][aria-expanded="false"]` once (scroll into view, then click); sweep again for the toggle bodies; list the titles of toggles that did not open.
6. When no browser tool is available, say so and ask for a paste or a screenshot. Never judge coverage from an empty fetch.
7. Keep the transcription as a file (for example `pm-ticket-source.md`). It becomes a capture asset only if rows are captured in P7.

P3, split:
8. Split into atomic asks A1, A2, …, one observable sentence each. Keep the source's own numbering when it has one.

P4, search per ask, recording every candidate, not only the first:
9. UR `input.md` files in `user-requests/` and `archive/UR-*/`.
10. REQ title, heading, `## What` and `## Detailed Requirements` in `queue/`, `working/` and `archive/`, by intent, not by filename.
11. The matched REQ's Implementation Summary and its `commit:` value.
12. `git log`, matched the way the board does (subject `[REQ-NNN]`, a touched REQ file path, plus the merge range of a matched hand-back merge), and by `Traced-to:` lines naming the REQ file. Reuse the board's matching; do not rebuild it (see Assumptions).
13. Code and tests, by the names the ask uses. A test whose assertion matches the REQ's "GREEN when" line counts as evidence for that criterion.

P5, one table:
14. Columns: Ask | REQ (plain title, status) | Commit | Criteria evidenced | Verdict (as of) | Gap | Open question. The report's three-row example is the shape.
15. Verdicts: `complete` (every acceptance criterion of the known scope has evidence: commit, test, or file and line); `partial` (some criteria have evidence, the unfinished ones are named); `not started` (no code evidence, and a pending REQ or no REQ; the search terms are the evidence); `unverified` (evidence not enough either way, including any row built from an incomplete fetch).
16. Every verdict carries the date it was checked and one pointer to its evidence. "Criteria evidenced" counts the matched REQ's Detailed Requirements bullets plus its "GREEN when" line.

P6, questions, then stop:
17. Ask the open questions with the ask-user tool, one per row that has one, following `crew-members/clear-questions.md`. Then stop. Trace writes nothing to the repo unless the user says go.

P7, hand off only the gaps:
18. On go, call `do-work capture` once, with the table and the transcription as the verbatim input, for `not started` and `partial` rows only.
19. A partial row that matches an archived REQ uses the existing addendum path (`actions/capture.md:125`), and P4's commit hash fills `## Prior Implementation` (`actions/capture.md:162`).
20. When the source has an external ticket id, the UR title carries it and `input.md` gets an optional `## PM Ticket References` body section: ticket id, source URL, fetch date, transcription asset path. Add that section to the UR `input.md` template in `actions/capture-reference.md`.
21. Trace never picks the UR number. Capture scans it at write time (`actions/capture.md:78`), because a browser sweep can take minutes while another session captures.

P8, done work is not captured again:
22. By default the table cites the archived REQ and commit for a `complete` row and nothing is written.

All:
23. Document the action in the user-facing docs the way sibling actions are documented (help, guides, any action list a contract test checks), and release per `_dev/primes/prime-releases.md`.
## Constraints
- Out of scope (report, Out of scope): writing to a PM tool or changing a ticket's status; any consumer's ledger file format; new frontmatter fields or statuses; changes to capture's numbering, reservations or transaction; judging whether the found code is correct (that is `review-work`); running trace automatically in the background or on every capture.
- The only new structure is the optional `## PM Ticket References` body section in a UR's `input.md`. No new frontmatter field, no new status.
- Read-only until go: after a trace with no go-ahead, `git status --porcelain` is empty.
- No row from a fetched source is `complete` without fetched text.
## Assumptions (recorded at capture, no questions asked)
- **One REQ.** The maintainer asked for one REQ per numbered item in the Request section, and one REQ when it has one ask. The Request section has one ask (the `trace` action); P1 to P8 sit in Proposed direction and are this REQ's parts. Capture's challenge, recorded and not acted on: this is a large REQ. P2's browser adapter and P7's hand-off are the most separable parts. If the maintainer wants a split, `do-work abandon` this REQ and recapture.
- **D1, name: separate `trace` action.** Adopted as the report recommends, over `capture --check-coverage`, because `actions/capture.md:44` keeps read-only reports out of capture. `trace` is a single-word invocation verb, the same exemption `do-work run` uses under the naming rule.
- **D2, retroactive record: not in this REQ.** A UR written straight into `archive/` with `Traced-to:` lines and no REQ has no path today, and the report says leaving it out of the first version is fine. If wanted later, it is its own capture.
- **D3, Notion adapter location: `actions/trace-reference.md`.** Adopted as recommended. Loaded only for a URL whose web fetch result is a shell, so the core action stays generic.
- **D4, narrowing capture's filename-only archive scan: not in this REQ.** The report marks it optional. It would add per-capture reading cost to every capture, which is a separate trade-off. If the maintainer wants it, it is its own small REQ against `actions/capture.md:115`.
- **Programs beat prose for P4 step 12.** Matching commits to a REQ by subject, path, merge range and `Traced-to:` is mechanical and already exists in Go (`correlateCommitsToRequests`, in the `queue-kanban` module). Prose `git log` recipes in the action would be a second hand-maintained copy of that predicate (lessons family `paired-predicate-drift`). Recommended: expose the existing correlation through one read-only machine-output command (in the board tool or `do-work-cli`, builder's choice, one source of truth) that the action calls, and add `Traced-to:` matching there. If the builder finds this too large, the action may prescribe the `git log` matching in prose, and the REQ's Implementation Summary records why. If the command touches `queue-kanban`, `_dev/primes/prime-kanban-board.md` versioning applies. That command gets focused tests in its own harness even though this REQ is `tdd: false`.
- **Activity correlation decision is open.** The maintainer has not decided whether to keep the board's activity correlation code or replace part of it (the ai-report for REQ-632, activity correlation determinism, holds the two options). Trace reads whatever the board computes and adds no second correlation, so either outcome flows through.
- **"Read-only until go" vs "keep the transcription as a file".** The two conflict if the file lands in the repo. Resolution: before go, the transcription lives in a temporary directory outside the working tree, and the action prints its path. On go it becomes the capture's raw input and asset. This keeps `git status --porcelain` empty after a no-go trace.
- **Browser tool is platform-specific.** The skill is platform-agnostic. The adapter describes the method (headed browser, scroll sweep, toggle clicks) and names Playwright only as an example tool. No browser tool means ask for a paste or screenshot (P2 item 6), never a guess.
- **Routing phrase `coverage`.** The bare word can also mean test coverage. The row routes `coverage` to trace when it comes with a spec, URL, image or UR target. The builder settles the exact wording and keeps the bare unknown single word rule (`SKILL.md`, after the table) intact.
- **One UR on go.** The hand-off is one capture invocation, so new asks and addendum REQs share one new UR. If `capture-files` or the commit message format cannot hold a mixed batch, the builder reports that and uses the smallest change that keeps one capture per trace.
- **Verdict date.** "As of" is the UTC date the trace ran. The verdict vocabulary is exactly the four values above, taken from the consumer ledger so results stay comparable over time.
- **Search scale.** The archive holds hundreds of REQs. The builder may narrow candidates by a fixed-string search over titles and `## What` before reading full files, as long as every candidate found is listed, not only the first.
## Dependencies
No `depends_on` edge. Related queued work: REQ-668 (the `do-work status` action) also adds a routing row to `skills/do-work/SKILL.md`, so the two touch the same table. Neither blocks the other.
## Builder Guidance
Medium certainty on what to show; the report gives the table, verdict vocabulary and acceptance checks. Latitude: the evidence-command design (see Assumptions), routing wording, table formatting, and how long parts split between `trace.md` and `trace-reference.md`. Keep `trace.md` short: read, split, search, print, ask, stop.
## Red-Green Proof
**RED prompt/case:** In this repo, ask a session "is this captured? <a spec with three asks>" or run `do-work trace UR-12`.
**Why RED now:** `skills/do-work/actions/trace.md` does not exist and `SKILL.md` has no row for a coverage question. The phrase falls through to capture (`SKILL.md:46`), which by its own rule (`actions/capture.md:44`) is not for read-only questions.
**GREEN when:** (from the report's Acceptance check)
- "is this captured? <url>", "do-work trace UR-12" and "how much of it is implemented" route to `actions/trace.md`. "capture-request: …" still routes to capture.
- In a fixture repo with one archived REQ (commit subject `[REQ-007]`, tests covering 2 of its 3 Detailed Requirements) and a spec with three asks (one done, one half done, one new), trace prints `complete` with the hash, `partial` with 2 of 3, and `not started`, each dated.
- After a trace with no go-ahead, `git status --porcelain` is empty.
- After go-ahead, one capture runs. The new UR holds only the partial and not-started asks. The partial ask becomes an addendum REQ with `addendum_to` and a Prior Implementation commit hash. A ticket id appears in the UR title and in `## PM Ticket References`.
- A Notion URL whose web fetch returns only "Notion" triggers the browser adapter, or a request for a paste when no browser tool exists. No row from that source is ever `complete` without fetched text.
- A screenshot source is read directly and described before the split.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: adds an action and a routing row.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8334 tokens, over budget; `slugged: partial`). Matching reason: reuses the board's commit-to-REQ correlation (families `paired-predicate-drift`, `git-history-evidence`).
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: the action may prescribe `git log` and search commands.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over budget; `slugged: partial`). Matching reason: a new read-only evidence command may land in `do-work-cli`.
## Full Context
See `do-work/user-requests/UR-148/input.md` for complete verbatim input. The source report stays at `do-work/inbox/2026-10-09_do-work-upstream-suggestion-trace-coverage.md`. No queued or archived REQ shares this intent (queue REQ-654 to REQ-668 read by intent; archive filenames scanned for trace, coverage and duplicate-check work).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-trace-coverage.md`, Request: "Please add one read-only action, `do-work trace <url | image | pasted text | UR-NNN>`, routed before capture."*

## Cancelled

- **When:** 2026-10-10T12:59:20Z
- **Why:** folded into the simplified recapture of 2026-10-10 (maintainer chose the aggressive simplification)
- **Decided by:** user, via `do-work abandon`
