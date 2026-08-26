---
id: REQ-375
title: 'Copy carries titles and a referenced-requests glossary'
status: pending
created_at: 2026-08-26T13:02:24Z
user_request: UR-074
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-374]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-374, REQ-376]
batch: ticket-id-autocomplete
write_set:
  - skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go
  - _dev/primes/prime-kanban-board.md
---

# Copy Carries Titles And A Referenced-Requests Glossary

## What

The clipboard payload gains the same treatment REQ-374 gives the drawer: the first mention of each
REQ or UR id in the body is expanded with its title, and a glossary of every referenced ticket is
appended at the end. The frontmatter fence is never touched, so a paste still parses as a REQ file.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Copy is the surface where a bare number hurts most. The whole point of the button is to paste a
ticket into a fresh agent session, and that session has no board to click through — it gets
`REQ-1679/REQ-1108` as literal text with nothing attached. This REQ is the half of the user's ask
that says "even on copy… it should always autocomplete".

## Context

All copy behaviour is in `web/board-clipboard.js` (300 lines). The payload is produced on the Go
side and shipped lazily:

- `generate.go:804` `buildGeneratedBoardMarkdownData` stores `FrontmatterMarkdown + BodyMarkdown`
  per ticket — the file's original bytes, never re-serialized.
- It is written to `board-markdown.js`, deliberately not referenced by `index.html`, and injected
  by `loadBoardMarkdownData` (`board-clipboard.js:9`) on the first copy click.

Three handlers consume it: the drawer `Copy` (line 170, with a `copyTextWithHeading` fallback for a
stale bundle that lacks the markdown sibling), the UR `Copy all` (line 246), and the per-column
`Copy all` (line 278). `rawMarkdownForRequests` (line 222) concatenates with cat semantics and
throws rather than publishing a partial clipboard.

**The verbatim contract this REQ ends.** The comment block at `board-clipboard.js:50-65` states
that the primary path copies the ticket file "VERBATIM" because "verbatim has to mean verbatim or
the paste stops round-tripping back into a valid file". Annotation ends that as written. What
survives, and what makes the trade acceptable:

- The **Go payload stays byte-exact**. `TestBuildGeneratedBoardMarkdownDataKeepsExactSources` and
  `TestBuildGeneratedBoardMarkdownDataRoundTripsTheWholeFile` must remain green **and unmodified** —
  they are the pin that the annotation is a client-side presentation step, not a payload change.
- The **frontmatter fence is never annotated**, so a paste still parses as a REQ file with a valid
  `depends_on`, `related` and `user_request`.
- The additions are **marked as additions** by the glossary's own heading, so a reader can see what
  the board added.

## Detailed Requirements

- **`annotateTicketMentions(markdownText)`** — expands the first mention of each id in the body,
  using the resolver REQ-374 moves into `board-core.js`. Three exclusions, all mandatory:
  1. **The frontmatter fence is skipped entirely.** `depends_on: [REQ-1679]` and
     `user_request: UR-389` must stay parseable YAML. Split on the same leading `---` fence the Go
     side splits on.
  2. Fenced code blocks are skipped.
  3. Inline code spans are skipped.
- **`buildReferencedTicketsGlossary(ids, excludedIds)`** — one appendix at the very end of the
  payload, listing each referenced ticket with its **full, untruncated** title and status:

  ```
  ---
  ## Referenced requests (added by the board — not part of the file)
  - REQ-1679 — Admin can delete a card — any card, admin-only, mapped level assets deleted too (completed)
  - UR-389 — <the user request's title>
  ```

- **`excludedIds` drops ids whose full text is already in the payload.** A column or UR `Copy all`
  must not gloss the tickets it already contains — their titles are right there in their own
  frontmatter.
- **Wire all three handlers.** Drawer `Copy` gets the same treatment on both its primary and its
  `copyTextWithHeading` fallback path, so the two shapes agree. UR `Copy all` and column `Copy all`
  annotate the joined result; `rawMarkdownForRequests`'s cat semantics and its throw-on-missing
  behaviour are unchanged.
- **Both id kinds**, matching REQ-374: UR ids are expanded and glossed exactly as REQ ids are.
- **An unresolved id is named as not-found, not left silent.** It is not expanded (there is no
  title), but it gets a glossary line reading `- REQ-9999 — not found in this queue`. Plain text
  has no red, so the glossary line is the paste's equivalent of REQ-374's blocked-accent span:
  the reader of the paste learns the reference is dead instead of hunting for it. An ambiguous
  REQ segment is not a dead reference and gets no line.
- **Rewrite the contract comment at `board-clipboard.js:50-65`** to state the new rule: the Go
  payload is verbatim, the clipboard payload is those bytes plus marked annotations, and the
  frontmatter fence is never touched. Record the same split in `_dev/primes/prime-kanban-board.md`
  so the next reader of that prime is not working from the old invariant.

## Constraints

- **Never rewrite a REQ file.** The annotation exists only in the clipboard payload.
- **Never annotate frontmatter.** This is the constraint that turns a convenience into a
  correctness requirement: a pasted file whose `depends_on` no longer parses is worse than a
  cryptic number.
- **Do not modify the two Go round-trip tests.** If either needs changing, the annotation has
  leaked into the payload and the approach is wrong.
- **Write-set overlap with REQ-374** on `generate_test.go` is real and deliberate — it is why this
  REQ carries `depends_on: [REQ-374]` rather than being dispatched in parallel. Do not fan these
  two out to concurrent builders.

## Dependencies

`depends_on: [REQ-374]` — this REQ consumes the shared resolver (`resolveTicketMention`,
`ticketTitleFor`) that REQ-374 moves into `board-core.js`, and shares a test file with it.

## Builder Guidance

**Certainty level: Firm.** The copy shape (inline expansion **and** glossary) was chosen by the user
from three presented options during capture, with the verbatim-contract cost stated in the option
they picked. It is a taken decision, not an open one — but say in the REQ's own record how you kept
the frontmatter safe, because that is the part a reviewer should check hardest.

Both new functions are pure string transforms and should be tested as such under the Node harness
(`TestJavaScriptBehaviorDrawerHeadingDeduplication`, `generate_test.go:905`, is the pattern), then
confirmed end-to-end in the existing Chromium clipboard lane
(`clipboard_browser_probe_test.go:55` `TestBrowserBehaviorBoardColumnCopyAll`).

Actually paste the result into an editor and read it. A payload that passes every assertion and
pastes as a broken file is the failure this REQ is most exposed to.

## Red-Green Proof

**RED prompt/case:** Open REQ-1685 on the board and press `Copy`, then paste. Today the clipboard
holds the file's raw bytes: the PLAN line still reads `REQ-1679/REQ-1108 lessons` with no titles,
and there is no reference list anywhere in the paste.

**Why RED now:** The drawer `Copy` handler (`board-clipboard.js:170`) returns `rawMarkdown`
untouched whenever `board-markdown.js` resolved, by explicit design — the comment at line 50-65
says so.

**GREEN when:** The pasted text's frontmatter block is byte-identical to the source file (checked
by diffing that region), the body's first `REQ-1679` carries its title and the later one does not,
a backticked id is untouched, the glossary appears once at the end with full titles and statuses,
an id with no board record appears there as `not found in this queue`,
and a column `Copy all` glosses only ids that are not themselves in the payload.

**Validation:** User confirmed — "Both: glossary + inline" was selected from three options with the
round-trip cost stated in the option text, and reaffirmed in the follow-up: "it's good to have a
glossary and an inline mention when they appear for the first time."

## Full Context

See `do-work/user-requests/UR-074/input.md` for complete verbatim input.

---
*Source: user request in session, 2026-08-26, prompted by REQ-1685's body citing REQ-1679/REQ-1108.*

## Addendum (2026-08-26)

User added, mid-run on REQ-374:

> ````text
> address the **The "Silent Failure" Gap:** <- now show the placeholder as broken (red with tooltip)
> The text explicitly notes that if a REQ ID is mentioned but doesn't actually exist in the system, it simply renders as ordinary prose. There is no warning or "broken link" indicator for dead references, meaning typos in IDs remain invisible until someone manually tries to find them.
> ````

- The directive names the *display* surface, which is REQ-374's. This REQ carries the copy-surface
  equivalent, since a paste has no colour: an unresolved id earns a `not found in this queue`
  glossary line instead of the silence the original requirement specified.
- The reversed requirement is edited in place above rather than left contradicted; this section is
  the record that it changed and why.
