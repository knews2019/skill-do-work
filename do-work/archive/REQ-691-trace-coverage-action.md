---
id: REQ-691
title: 'do-work trace action reports how much of an outside spec is captured and built, with one dated verdict per ask'
status: completed
route: C
estimate:
  p50_active_minutes: 40
  confidence: medium
  basis:
  - Route C
  - 9-file write set
  - 2 subsystems involved
  - 7 acceptance criteria
  calculated_at: 2026-10-10T13:19:25Z
created_at: 2026-10-10T13:05:57Z
user_request: UR-153
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md", "_dev/primes/prime-kanban-board.md", "_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
claimed_at: 2026-10-10T12:52:21Z
status_changed_at: 2026-10-10T13:14:55Z
planning_at: 2026-10-10T13:21:17Z
related: [REQ-690, REQ-692]
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/actions/trace.md", "skills/do-work/actions/trace-reference.md", "skills/do-work/SKILL.md", "skills/do-work/actions/help.md", "skills/do-work-board/tools/queue-kanban/activity_correlation.go", "skills/do-work-board/tools/queue-kanban/request_commits.go", "skills/do-work-board/tools/queue-kanban/request_commits_test.go", "skills/do-work-board/tools/queue-kanban/main.go", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md"]
builder_handback_at: 2026-10-10T13:32:48Z
integration_at: 2026-10-10T18:04:45Z
review_at: 2026-10-10T18:18:41Z
kb_status: pending
commit: a270e71294b3c1050d21547ecc0e35ab20ced76f
heavy_verified_at: 2026-10-10T18:18:53Z
heavy_verified_revision: a270e71294b3c1050d21547ecc0e35ab20ced76f
completed_at: 2026-10-10T18:19:28Z
release_at: 2026-10-10T18:19:28Z
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
1. Add `skills/do-work/actions/trace.md` and `skills/do-work/actions/trace-reference.md` (table template, source handling).
2. Add a routing row in `skills/do-work/SKILL.md` above the capture row for `trace`, `coverage`, `is this captured`, `already implemented`, `how much of it is implemented`, `didn't we have a request`. Add `trace <url|image|text|UR-NNN>` to the argument hint at `SKILL.md:4`. `capture-request: …` still routes to capture.
3. Load `crew-members/prompt-injection.md` before reading the source.

P2, read the source:
4. A UR id reads that UR's `input.md` and assets. An image is read directly and described as `actions/capture.md:222` does, before the split. Pasted text is used as is.
5. A URL goes to the platform's web fetch first. When the result is a shell (only the site name, or almost no text), say so and ask for a paste or a screenshot. Never judge coverage from an empty fetch. (The source report's Notion browser adapter, with site-specific selectors and scroll steps, is dropped by the maintainer's 2026-10-10 decision: a platform-agnostic skill does not carry one site's DOM.)
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
20. When the source has an external ticket id, the UR title carries it, and the transcription asset's first lines record the ticket id, source URL and fetch date. No new UR body section is added.
21. Trace never picks the UR number. Capture scans it at write time (`actions/capture.md:78`), because a browser sweep can take minutes while another session captures.

P8, done work is not captured again:
22. By default the table cites the archived REQ and commit for a `complete` row and nothing is written.

All:
23. Document the action in the user-facing docs the way sibling actions are documented (help, guides, any action list a contract test checks), and release per `_dev/primes/prime-releases.md`.
## Constraints
- Out of scope (report, Out of scope): writing to a PM tool or changing a ticket's status; any consumer's ledger file format; new frontmatter fields or statuses; changes to capture's numbering, reservations or transaction; judging whether the found code is correct (that is `review-work`); running trace automatically in the background or on every capture.
- No new structure: no new frontmatter field, no new status, no new UR body section.
- Read-only until go: after a trace with no go-ahead, `git status --porcelain` is empty.
- No row from a fetched source is `complete` without fetched text.
## Assumptions (recorded at capture, no questions asked)
- **One REQ.** The source report had one ask (the `trace` action) in eight parts; the maintainer kept it as one REQ on 2026-10-09 and on 2026-10-10 dropped the Notion browser adapter and the PM Ticket body section.
- **D1, name: separate `trace` action.** Adopted as the report recommends, over `capture --check-coverage`, because `actions/capture.md:44` keeps read-only reports out of capture. `trace` is a single-word invocation verb, the same exemption `do-work run` uses under the naming rule.
- **D2, retroactive record: not in this REQ.** A UR written straight into `archive/` with `Traced-to:` lines and no REQ has no path today, and the report says leaving it out of the first version is fine. If wanted later, it is its own capture.
- **D4, narrowing capture's filename-only archive scan: not in this REQ.** The report marks it optional. It would add per-capture reading cost to every capture, which is a separate trade-off. If the maintainer wants it, it is its own small REQ against `actions/capture.md:115`.
- **Programs beat prose for P4 step 12.** Matching commits to a REQ by subject, path, merge range and `Traced-to:` is mechanical and already exists in Go (`correlateCommitsToRequests`, in the `queue-kanban` module). Prose `git log` recipes in the action would be a second hand-maintained copy of that predicate (lessons family `paired-predicate-drift`). Recommended: expose the existing correlation through one read-only machine-output command (in the board tool or `do-work-cli`, builder's choice, one source of truth) that the action calls, and add `Traced-to:` matching there. If the builder finds this too large, the action may prescribe the `git log` matching in prose, and the REQ's Implementation Summary records why. If the command touches `queue-kanban`, `_dev/primes/prime-kanban-board.md` versioning applies. That command gets focused tests in its own harness even though this REQ is `tdd: false`.
- **Activity correlation is settled.** The maintainer decided on 2026-10-10 (shipped 0.305.88): the merge-range scaffold was deleted and nothing stamps commits; the board's commit-to-REQ correlation stays as code. Trace reads whatever the board computes and adds no second correlation.
- **"Read-only until go" vs "keep the transcription as a file".** The two conflict if the file lands in the repo. Resolution: before go, the transcription lives in a temporary directory outside the working tree, and the action prints its path. On go it becomes the capture's raw input and asset. This keeps `git status --porcelain` empty after a no-go trace.
- **Routing phrase `coverage`.** The bare word can also mean test coverage. The row routes `coverage` to trace when it comes with a spec, URL, image or UR target. The builder settles the exact wording and keeps the bare unknown single word rule (`SKILL.md`, after the table) intact.
- **One UR on go.** The hand-off is one capture invocation, so new asks and addendum REQs share one new UR. If `capture-files` or the commit message format cannot hold a mixed batch, the builder reports that and uses the smallest change that keeps one capture per trace.
- **Verdict date.** "As of" is the UTC date the trace ran. The verdict vocabulary is exactly the four values above, taken from the consumer ledger so results stay comparable over time.
- **Search scale.** The archive holds hundreds of REQs. The builder may narrow candidates by a fixed-string search over titles and `## What` before reading full files, as long as every candidate found is listed, not only the first.
## Dependencies
No `depends_on` edge. Related queued work: REQ-690 (do-work status) and REQ-692 (validate-feedback --capture) also add rows to `skills/do-work/SKILL.md`, so the three touch the same table; whichever lands later rebases its row.
## Builder Guidance
Medium certainty on what to show; the report gives the table, verdict vocabulary and acceptance checks. Latitude: the evidence-command design (see Assumptions), routing wording, table formatting, and how long parts split between `trace.md` and `trace-reference.md`. Keep `trace.md` short: read, split, search, print, ask, stop.
## Red-Green Proof
**RED prompt/case:** In this repo, ask a session "is this captured? <a spec with three asks>" or run `do-work trace UR-12`.
**Why RED now:** `skills/do-work/actions/trace.md` does not exist and `SKILL.md` has no row for a coverage question. The phrase falls through to capture (`SKILL.md:46`), which by its own rule (`actions/capture.md:44`) is not for read-only questions.
**GREEN when:** (from the report's Acceptance check)
- "is this captured? <url>", "do-work trace UR-12" and "how much of it is implemented" route to `actions/trace.md`. "capture-request: …" still routes to capture.
- In a fixture repo with one archived REQ (commit subject `[REQ-007]`, tests covering 2 of its 3 Detailed Requirements) and a spec with three asks (one done, one half done, one new), trace prints `complete` with the hash, `partial` with 2 of 3, and `not started`, each dated.
- After a trace with no go-ahead, `git status --porcelain` is empty.
- After go-ahead, one capture runs. The new UR holds only the partial and not-started asks. The partial ask becomes an addendum REQ with `addendum_to` and a Prior Implementation commit hash. A ticket id appears in the UR title and in the transcription asset.
- A URL whose web fetch returns only the site name makes trace ask for a paste or a screenshot. No row from that source is ever `complete` without fetched text.
- A screenshot source is read directly and described before the split.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
Refreshed at claim (pre-dispatch, 2026-10-10) against `do-work/lessons-index.md`. Kept: `_dev/primes/lessons-releases.md` (666 tokens; matching reason: `trace.md` cites the board tool across a package boundary, which the shipped-package reference contract checks). Dropped:
- `_dev/primes/lessons-action-files.md` as a whole satellite (7756 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: adds an action and a routing row.
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` as a whole satellite (8815 tokens, over budget; `slugged: partial`). Matching reason: reuses the board's commit-to-REQ correlation (families `paired-predicate-drift`, `git-history-evidence`).
- `_dev/primes/lessons-kanban-board.md` as a whole satellite (6122 tokens, over budget; `slugged: partial`). Matching reason: adds a queue-kanban subcommand.
- `_dev/primes/lessons-shell-commands.md` as a whole satellite (8480 tokens, over budget; `slugged: partial`). Matching reason: the action prescribes search and build commands.
No longer a match: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (the Plan's D-01 puts the new command in the board, not in `do-work-cli`).
## Full Context
See `do-work/user-requests/UR-153/input.md` for the 2026-10-10 decision record. The cancelled original, `do-work/archive/UR-148/REQ-669-trace-coverage-action.md`, carries the complete body; the source report stays at `do-work/inbox/2026-10-09_do-work-upstream-suggestion-trace-coverage.md`, with verbatim input in `do-work/archive/UR-148/input.md`.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Followed the REQ's Plan D-01 to D-08 without change. Board: extract one attribution helper, add a read-only `request-commits` with the frontmatter command's testable shape (writers in, exit code out). The git runner is passed in so the tests stub it. Prose: keep `trace.md` short and put the long parts in `trace-reference.md`. Routing row above the capture fallback. One help line. Tests first for the Go RED. Then a hand dry run in a scratch fixture repo for the GREEN cases, including the go-ahead capture through the real `capture-files`.
- [x] **[APPLY]:** Wrote the two tests and ran them against a stub. Then wrote the helper, the command and the wiring, then the two action files, the routing row, the hint and the help line. The dry run showed one case the reference did not cover: an ask narrower than its REQ. I added one clause for it (D-10). No file outside Scope was touched.
- [x] **[UNIFY]:** `git diff bd56c4b0 --stat` shows exactly the nine Scope files. Checked files: all nine. I read each diff. There is no debug output and no build artifact; the built `queue-kanban` binary is gitignored and was removed after the dry run. Check results (probe, focused `go test`, `gofmt -l`, `go vet`, `git diff --check`, `contract-regressions.sh`, all green) and the nine-file stat are in the hand-back, `do-work/runs/work-2026-10-10-131527/REQ-691-handback.md`.
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-trace-coverage.md`, Request: "Please add one read-only action, `do-work trace <url | image | pasted text | UR-NNN>`, routed before capture."*

## Triage

**Route: C** - Complex

**Reasoning:** A new user-facing action spanning two packages: core prose (`actions/trace.md`, `actions/trace-reference.md`, the `SKILL.md` routing row and argument hint, help and next-steps) and one read-only Go subcommand in the board module that exposes the existing commit-to-REQ correlation. Twenty-three requirements and one design call (where the commit evidence comes from) need a plan before the build.

**Planning:** Required

## Plan

Written by the pre-dispatch agent, which already held the exploration context, in place of a separate Plan agent. Base: main at bd56c4b0 (0.305.101).

**Decisions made at pre-dispatch (D-01 onwards; the builder continues from D-09):**

- **D-01, where the commit evidence comes from: one read-only board subcommand, `queue-kanban request-commits`.** The commit-to-REQ predicate lives only in the board (`skills/do-work-board/tools/queue-kanban/activity_correlation.go`, `requestPathPattern` :60, `requestSubjectPrefixPattern` :64, `correlateCommitsToRequests` :114). The core CLI is a separate Go module, so a `do-work-cli` command would be a second copy of that predicate (lessons family `paired-predicate-drift`), and prose `git log` recipes in `trace.md` would be a third. The board command reuses the predicate by extracting one helper that both callers use. Value: one source of truth, testable. Risk: trace needs a Go toolchain for the git half of the search; D-04 covers its absence. Reversible.
- **D-02, `Traced-to:` matching needs no code.** Supplementary commits written by `actions/commit.md` (:126-128, format at :137) start their subject with `[{REQ id}]`, so the subject token already credits them. Requirement 12's `Traced-to:` clause is met by the subject match; the Implementation Summary says so. No commit-body parsing.
- **D-03, no merge range.** Requirement 12 says "plus the merge range of a matched hand-back merge". That expansion was deleted from the board at 0.305.88 (comment at `activity_correlation.go:104-113`; maintainer decision 2026-10-10). The REQ's own Assumptions say trace reads whatever the board computes, so the command credits by subject token and touched REQ path only. Builder commits carry the `[REQ-NNN]` prefix by the branch contract, so they are found directly.
- **D-04, no Go toolchain or a failed build.** Trace says the git-history search was skipped, uses each matched REQ's `commit:` value as its commit evidence, and marks a row `unverified` (never `not started`) when its verdict would rest on the skipped search. Same shape as `actions/forensics.md` Check 14.
- **D-05, output shape.** Tab-separated, with a header line: `request_id`, `commit` (full hash), `committed_at` (RFC 3339 as git prints `%cI`), `paths_outside_do_work` (count of touched paths not under `do-work/`), `subject`. One row per (requested id, commit), in `git log` order (newest first). The count tells trace whether a commit is lifecycle bookkeeping (0: claim, archive, run artifacts, and every merge, because `git log --name-only` lists no paths for a merge) or code and docs evidence (above 0). Full history, no `--since` window: about 1.4 s on this repo's 4,181 commits. Exit 0 with the header alone when nothing matches; exit 2 for no ids or an argument that is not `REQ-` plus digits; exit 1 with a stderr line when git fails. Read-only, so the board's write-surface count in `_dev/primes/prime-kanban-board.md` does not change.
- **D-06, docs surface.** The action's own files, the `SKILL.md` routing row and argument hint, one line in the `actions/help.md` menu, the board's subcommand list in `prime-do-kanban.md:3` and the unknown-subcommand message in `main.go`. No new `docs/*-guide.md` (clarify, abandon and stakeholder-answers ship without one; the action's When to Use carries the triggers) and no `next-steps.md` row (trace owns its own hand-off). `_dev/tests/staged-skills-contract.sh`'s `core_files` list is not exhaustive (run-with-recovery and stakeholder-answers are absent), so it is not edited.
- **D-07, routing row position and wording.** Directly above the capture fallback row (`SKILL.md:46`), so every earlier row still wins and `capture-request:` still routes to capture. Triggers: `trace`, `is this captured`, `already implemented`, `how much of it is implemented`, `didn't we have a request`, and `coverage` only together with a spec, URL, image or UR target (the bare word stays on the unknown-single-word rule after the table). REQ-692 (validate-feedback capture route) adds its own row above the same fallback; the triggers do not overlap, so their order between the two rows does not matter.
- **D-08, transcription before go.** The source transcription is written under a `mktemp -d` directory outside the working tree and its path is printed; on go it becomes capture's raw input and asset. `git status --porcelain` stays empty after a no-go trace (the board binary that `go build` writes is gitignored by the board tool's own `.gitignore`).

**Tasks (in order):**

1. **Board command.** In `activity_correlation.go`, extract the per-commit attribution (touched REQ path or bracketed subject token) from `correlateCommitsToRequests` into one helper that returns the REQ ids a commit credits; `correlateCommitsToRequests` calls it unchanged in behaviour. Add `request_commits.go` (the subcommand: flag parsing, one `git log --format=%H%x00%cI%x00%s --name-only` through `gitCommandRunner`, `parseCorrelationLog`, the helper, TSV output) and `request_commits_test.go` with two focused tests on canned log output: `TestRequestCommitsListsOnlyCommitsTheBoardCredits` (a subject-token match, a REQ-path-only match, an unbracketed id that must not match, a lifecycle commit with 0 paths outside `do-work/`) and `TestRequestCommitsRefusesAnArgumentThatIsNotARequestId`. Wire the subcommand in `main.go` (switch case and the unknown-subcommand list) and add it to `prime-do-kanban.md:3`'s subcommand list.
2. **`skills/do-work/actions/trace.md`.** Short (read, split, search, print, ask, stop, hand off on go), following the action template in `_dev/primes/lessons-action-files.md#template`, with the package-justification sentence in its description blockquote (core: it completes capture's duplicate check and hands off to capture). Loads `crew-members/prompt-injection.md` before reading the source and `crew-members/clear-questions.md` before asking.
3. **`skills/do-work/actions/trace-reference.md`.** Source handling (UR id, image, pasted text, URL with the shell-fetch rule), the per-ask search order (requirements 9-13, every candidate listed), the `request-commits` invocation and the D-04 fallback, the table template with a three-row example, the four verdict definitions, the date rule, and the go-ahead capture payload (requirements 18-21).
4. **Routing and help.** `SKILL.md` row (D-07) and `trace <url|image|text|UR-NNN>` in the argument hint at `SKILL.md:4`; one line in the `actions/help.md` menu beside `roadmap`.
5. **Proof.** The two Go tests plus the existing `activity_correlation_test.go` tests; a hand dry run of trace against a scratch fixture repo for the GREEN cases (see the builder brief).

**Plan validation:**
- Requirement coverage: 1 → T2/T3; 2 → T4; 3 → T2; 4-5, 7 → T3; 8 → T2; 9-11, 13 → T3; 12 → T1 + T3 (D-02, D-03); 14-16 → T3; 17 → T2; 18-21 → T3; 22 → T2/T3; 23 → T4 (docs) and the integrator (release). Requirement 6 does not exist in this REQ's numbering (it was the dropped browser adapter); nothing to cover.
- No orphan tasks: T1 serves requirement 12 and the Assumptions' "programs beat prose".
- Scope sanity: 5 tasks, which is past the 3-task comfort line. Accepted: T4 is two small edits and T5 is proof; the maintainer kept this as one REQ (UR-153), so no split.
- Consumer field contract: the go-ahead hand-off is the one mutation (capture). Per row it must carry: the ask id (A1…), the verdict, the matched REQ id and its archive path (for `addendum_to`), the commit hash (for `## Prior Implementation`), and, when present, the external ticket id, source URL and fetch date (for the UR title and the transcription asset's first lines). `request-commits` output is read for judgment only and drives no mutation.

*Generated by Plan agent*

## Exploration

Pre-dispatch exploration, 2026-10-10, read-only, at bd56c4b0 (0.305.101).

- **The predicate to reuse.** `skills/do-work-board/tools/queue-kanban/activity_correlation.go`: `requestPathPattern` (:60) matches a REQ file in queue, working or archive (any UR nesting) and run artifacts named for a REQ; `requestSubjectPrefixPattern` (:64) matches every bracketed `[REQ-NNN]` in a subject; `parseCorrelationLog` (:78) reads `git log --format=%H%x00%cI%x00%s --name-only`; `correlateCommitsToRequests` (:114-131) attributes by those two rules only and returns instants, not hashes. `gitCommandRunner` (:27) is the injectable runner tests already stub. Existing tests in `activity_correlation_test.go` (:93 path-only match, :125 every prefix token, :318 the path pattern) must stay green after the helper extraction.
- **Merge range is gone.** The comment at `activity_correlation.go:104-113` records the 0.305.88 deletion of the second-parent range expansion. Lesson (`lessons-do-kanban.md`, family `git-history-evidence`, 0.305.67): `git log --name-only` lists no paths for a merge commit, and a `do-work/` pathspec drops merges entirely, so the command must log without a pathspec.
- **Traced-to.** `skills/do-work/actions/commit.md:126-128` and `:137`: supplementary commits use the subject `[{REQ id}] {REQ title} — additional changes` plus a `Traced-to:` body line, so the subject token already credits them.
- **Board CLI shape.** `main.go:66-88` is the subcommand switch; `:86` is the unknown-subcommand message that lists every subcommand; `main.go:11-31` is the header synopsis listing each subcommand's usage. `runFrontmatterCommand(args, standardOut, standardErr) int` (`frontmatter_cli.go:172`) is the testable shape (writers in, exit code out) to copy. `exitOnLeftoverArguments` is at `main.go:284`. `prime-do-kanban.md:3` lists the subcommands. Write surfaces are counted in `_dev/primes/prime-kanban-board.md` § Conventions and pinned by `_dev/tests/contract-regressions.sh`; a read-only command leaves the count alone. The board has no version of its own; the integrator releases with the suite version.
- **How a core action runs the board tool.** `skills/do-work/actions/forensics.md:72-78` (Check 14): build with `go build -o queue-kanban .` inside the board tool directory, run with `--repo-root <project-root>`, and if `go` is absent or the build fails, skip and report the check as unverified. The built binary is gitignored by `skills/do-work-board/tools/queue-kanban/.gitignore` (`/queue-kanban`). Full-history `git log --format=%H%x00%cI%x00%s --name-only` on this repo: 27,914 lines, 4,181 commits, 1.4 s wall.
- **Routing.** `skills/do-work/SKILL.md:4` argument hint; routing table `:26-46`, first match wins; `:43` roadmap; `:46` capture fallback (`capture-request:` and unmatched descriptive multi-word input); `:48` the unknown-single-word rule. Earlier rows that could catch a trace phrase: `:31` (`verify`, `check`, …) and `:32` (`review`, …) match only as the command word, so "is this captured?" and "how much of it is implemented" reach the new row. `_dev/tests/staged-skills-contract.sh:288-333` reads the core routing section (Routing to Dispatch) and rejects only pipeline/full and moved-command rows; no test enumerates core routes, so the new row needs no contract edit. `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` holds no trace or coverage trigger.
- **Help and docs.** `skills/do-work/actions/help.md:7-38` is the core menu (roadmap at :25). `_dev/primes/prime-action-files.md:11`: every new action updates the owning `SKILL.md`, help, and any next-step surface. `skills/do-work/next-steps.md` lists only non-obvious cases, not every action. `skills/do-work/docs/` has guides for capture, cleanup, commit, forensics, review-work, roadmap, verify-requests, version and work only.
- **Template and citations.** The action template is `_dev/primes/lessons-action-files.md` § Template (:80-148: When to Use, Input, Steps, Output Format, Rules, Common Rationalizations, Red Flags, Verification Checklist; only earned sections). Cross-package citations follow `_dev/primes/prime-action-files.md` § Cross-Referencing and are checked by `_dev/tests/shipped-package-reference-contract.sh`. Lesson family `example-path-read-as-citation` (REQ-676): an example path whose first segment is a package content directory (`actions`, `docs`, `tools`, …) is read as a citation and fails; keep example file names like `pm-ticket-source.md` bare. Lesson family `read-only-action-tool-side-effects` (REQ-677): a read-only action that drives a tool says where it runs and what it may leave behind.
- **Capture hooks trace hands off to.** `actions/capture.md:44` (no read-only reports), `:113` (queued intent compare), `:115` (filename-only scan of working and archive), `:125` (archived match becomes an addendum REQ with `addendum_to`), `:162` (`## Prior Implementation` with the commit hash), `:222` (screenshot description). The UR-number scan is in capture's numbering step near `:78`; trace never picks a number. `actions/capture-reference.md:113` (`**GREEN when:**`) and `:125` (`## Detailed Requirements`) are the criteria trace counts.
- **Seams with siblings.** REQ-690 (do-work status) adds a `SKILL.md` row before the clarify and roadmap rows and may touch the argument hint and `help.md`. REQ-692 (validate-feedback capture route) adds a `SKILL.md` row above the capture fallback, the same neighbourhood as this REQ's row. REQ-689 (run coordinate mode) touches `SKILL.md` routing per the run table. No sibling touches `skills/do-work-board/tools/queue-kanban/` Go files except REQ-659 (`frontmatter_cli.go`, not touched here); `main.go` is not in any sibling's known set.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/trace.md` (new): the action: read, split, search, print, ask, stop, hand off on go
- `skills/do-work/actions/trace-reference.md` (new): source handling, per-ask search order, the commit-evidence call and its fallback, table template, verdicts, go-ahead payload
- `skills/do-work/SKILL.md` (modify): routing row above the capture fallback and the argument hint
- `skills/do-work/actions/help.md` (modify): one menu line
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (modify): extract the per-commit attribution helper, behaviour unchanged
- `skills/do-work-board/tools/queue-kanban/request_commits.go` (new): the read-only request-commits subcommand
- `skills/do-work-board/tools/queue-kanban/request_commits_test.go` (new): two focused tests on canned log output
- `skills/do-work-board/tools/queue-kanban/main.go` (modify): dispatch case, header synopsis line, unknown-subcommand list
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modify): subcommand list

**Files I will NOT touch:** `skills/do-work/actions/capture.md` and `capture-reference.md` (no change to capture's numbering, duplicate scan or UR template; REQ-688 owns the capture-reference fence fix), `skills/do-work/next-steps.md`, any `skills/do-work/docs/` guide, `_dev/tests/staged-skills-contract.sh`, `skills/do-work/tools/do-work-cli/` (no core CLI command; Plan D-01), every other queue-kanban file including `frontmatter_cli.go` and `web/`, `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION` and the version mirrors (the integrator releases), anything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] "is this captured? <url>", "do-work trace UR-12" and "how much of it is implemented" route to `actions/trace.md`; "capture-request: …" still routes to capture.
- [ ] In a fixture repo with one archived REQ (commit subject `[REQ-007]`, tests covering 2 of its 3 Detailed Requirements) and a spec with three asks (one done, one half done, one new), trace prints `complete` with the hash, `partial` with 2 of 3, and `not started`, each dated.
- [ ] After a trace with no go-ahead, `git status --porcelain` is empty.
- [ ] After go-ahead, one capture runs; the new UR holds only the partial and not-started asks; the partial ask becomes an addendum REQ with `addendum_to` and a Prior Implementation commit hash; a ticket id appears in the UR title and in the transcription asset.
- [ ] A URL whose web fetch returns only the site name makes trace ask for a paste or a screenshot; no row from that source is ever `complete` without fetched text.
- [ ] A screenshot source is read directly and described before the split.
- [ ] `queue-kanban request-commits` lists exactly the commits the board's correlation credits to each named REQ (subject token or touched REQ path), with hash, date, paths-outside-do-work count and subject; the board's activity correlation keeps its behaviour (existing `activity_correlation_test.go` tests pass).

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/trace.md` (new)
- `skills/do-work/actions/trace-reference.md` (new)
- `skills/do-work/SKILL.md` (modified)
- `skills/do-work/actions/help.md` (modified)
- `skills/do-work-board/tools/queue-kanban/activity_correlation.go` (modified)
- `skills/do-work-board/tools/queue-kanban/request_commits.go` (new)
- `skills/do-work-board/tools/queue-kanban/request_commits_test.go` (new)
- `skills/do-work-board/tools/queue-kanban/main.go` (modified)
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified)

**What was done:** The new core action `do-work trace` (`actions/trace.md`, six steps and three rules; every detail in `actions/trace-reference.md`) reads a UR, image, pasted text or URL after loading the prompt-injection guardrail, writes the transcription to a `mktemp -d` directory outside the tree, splits the source into asks A1, A2, ..., searches URs, REQs, credited commits, code and tests per ask, prints one table with a dated verdict per ask (`complete`, `partial`, `not started`, `unverified`), asks the open questions, and stops. Only on go does it hand the `partial` and `not started` rows to capture in one `capture-request:` call. The board gained the read-only `queue-kanban request-commits [--repo-root DIR] REQ-NNN...` (`request_commits.go`): one full-history `git log` with no pathspec, the per-commit attribution now extracted into `requestIdsCreditedByCommit` (`activity_correlation.go`) and shared with the board's activity lines, TSV output with a count of paths outside `do-work/`, exit 0/1/2. `SKILL.md` routes the trace phrases to the new action from a row directly above the capture fallback and adds `trace <url|image|text|UR-NNN>` to the argument hint; `help.md` and `prime-do-kanban.md` list it. Integrator: the hand-back merge had one conflict, the `SKILL.md` argument hint, resolved as the union of REQ-690's `status [REQ] [--watch]` and this REQ's `trace` token. Review fixes (integrator, on the builder branch, `775e24f9`, re-merged as `a270e712`): the `main.go` header comment and `prime-do-kanban.md` name `request-commits` (F1, F2, with a Read-first line for `request_commits.go`), the Commit Evidence gloss says the newest row above 0 can be a release commit (F3), and the go-ahead payload tells capture to start the UR title with the ticket id because capture has no ticket rule of its own (F5).

## Decisions

*(from the builder hand-back)*

- **D-09, DECIDE & STATE: wording of the routing row.** The row reads: `` | `trace`, `is this captured`, `already implemented`, `how much of it is implemented`, `didn't we have a request`, or `coverage` together with a spec, URL, image or UR target | `./actions/trace.md` | ``. It is placed directly above the capture fallback. The sentence about an unknown single word, after the table, is unchanged. A bare `coverage` therefore still gets the "captured or another command?" question.
- **D-10, DECIDE & STATE: an ask narrower than its REQ.** The acceptance check says one archived REQ yields both `complete` and `partial 2 of 3`. That is only possible when the done ask covers part of REQ-007. `trace-reference.md` now says: "When an ask covers only part of its REQ, count only the criteria inside the ask and name them (`1 of 1, bullet 1 only`)." This follows requirement 15's "known scope" wording. Reversible: one sentence.
- **D-11, DECIDE & STATE: the format of the "Criteria evidenced" cell is `N of M + GREEN yes|no`.** Requirement 16 counts the Detailed Requirements bullets plus the GREEN-when line. The acceptance text says "partial with 2 of 3". This format keeps both: the bullets as `2 of 3`, the GREEN line as its own yes or no. I changed the source report's example table to this format (`3 of 3 + GREEN yes`, `2 of 4 + GREEN no`). A REQ with neither section counts its `## What` as the one criterion.
- **D-12, DECIDE & STATE: the git runner is a fourth parameter.** `runRequestCommitsCommand(args, standardOut, standardErr, runner gitCommandRunner) int` is the frontmatter command's shape plus the existing injectable runner. `main.go` passes `runGitCommand`. This avoids a package-level variable just for tests.
- **D-13, DECIDE & STATE: flags first, then ids.** The command takes the synopsis order `[--repo-root DIR] REQ-NNN…`. A flag placed after an id is read as an id. It fails the `REQ-` plus digits check and exits 2, so it is never silently dropped. That follows the board's leftover-argument rule.
- **D-14, DECIDE & STATE: "read-only until go" is stated as unchanged status.** It is not stated as an empty status. The Rules line says `git status --porcelain` is the same after a no-go trace as before it, because a user's tree can already be dirty. In the clean fixture, that means empty, which is the acceptance check.
- **D-15, DECIDE & STATE: the Commit column cites the newest `request-commits` row above 0.** A row of 0 is bookkeeping or a merge. When no row is above 0, use the REQ's `commit:` value. Any non-zero exit, a build failure included, takes the D-04 path, so a git failure never reads as "no commits".
- **D-16, DECIDE & STATE: how the content is split.** `trace.md` (49 lines) holds the triggers, the six steps and three rules. Each step points to a named section of `trace-reference.md` (Source Handling, Search Order, Commit Evidence, Coverage Table, Go-Ahead Payload), which holds every detail. The verdict definitions and the table template live only in the reference, so they are stated once. `trace.md` names the four verdict words and does not define them.
- **D-17, DECIDE & STATE: the go-ahead hand-off is one `capture-request:` invocation.** Its verbatim input is the partial and not-started rows plus the transcription file as an asset. The fixture showed that `capture-files` holds the mixed batch (one addendum REQ and one new REQ under one UR) in one transaction. No change was needed.

## Discovered Tasks

*(from the builder hand-back)*

- `skills/do-work/actions/capture.md` Step 5 never shows the manifest's `operation` value. My first hand-built manifest used `"capture"` and was refused with `PUBLICATION-MANIFEST-INVALID`; `"capture-files"` is required. REQ-688 (`capture-files --example`) covers this gap. impact-negligible → report only
- The narrative header comment in `skills/do-work-board/tools/queue-kanban/main.go` (the paragraph before the synopsis) lists the subcommands in prose and was already missing `now`. I updated only the synopsis line and the unknown-subcommand list, as the brief says. impact-negligible → report only
- `_dev/tests/staged-skills-contract.sh` is heavy-only, so it could not run here. The integrator's gate runs it. The Plan's exploration found that no test lists core routes, so no contract edit should be needed. impact-negligible → report only

## Qualification

**Gate records (`advance --diff-range 00116871..d6aacf29`):** `qualify` satisfied, `scope-drift` satisfied. One `QUALIFY-NEW-FILE-UNWIRED` warning on `request_commits_test.go` is a false positive: it is a Go test file in package `main`, found by `go test`, not by a static reference. No debug artifacts, no leftover output primitives.

**Requirement trace against the diff (`git diff 00116871..d6aacf29 --stat`: 9 files, +295/-13) and the merged files:**
1. `skills/do-work/actions/trace.md` (49 lines) and `skills/do-work/actions/trace-reference.md` (Source Handling, Search Order, Commit Evidence, Coverage Table, Go-Ahead Payload) exist on main.
2. `skills/do-work/SKILL.md`: the trace row sits directly above the capture fallback with `trace`, `is this captured`, `already implemented`, `how much of it is implemented`, `didn't we have a request`, and `coverage` only with a target (D-09); the capture row stays last, so `capture-request:` still routes to capture. Argument hint carries `trace <url|image|text|UR-NNN>` beside REQ-690's `status [REQ] [--watch]`.
3. `trace.md` Step 1 loads `crew-members/prompt-injection.md` before reading the source.
4-5. Source Handling: UR id reads `input.md` and `assets/`; image is described before the split (as capture Step 4 does); pasted text as is; URL fetch first, a shell result asks for a paste or a screenshot and is never judged (`trace.md` Step 1 and Rules).
7. Transcription goes to a `mktemp -d` directory outside the tree with `Ticket:` / `Source:` / `Fetched:` first lines; it becomes a capture asset only on go.
8. `trace.md` Step 2: atomic asks A1, A2, ..., source numbering kept.
9-11, 13. Search Order items 1-3 and 5: UR inputs, REQs by intent in queue/working/archive with every candidate listed, Implementation Summary and `commit:`, code and tests by name with the GREEN-when rule.
12. Commit Evidence calls `queue-kanban request-commits`, which uses the board's own `requestIdsCreditedByCommit` (extracted from `correlateCommitsToRequests`, behaviour unchanged). `Traced-to:` is covered by the subject token (Plan D-02); the merge range is gone from the board since 0.305.88 (D-03). A failed build or git call takes the D-04 fallback and marks the row `unverified`.
14-16. Coverage Table: the seven columns, a three-row example, the `N of M + GREEN yes|no` count (D-11), the as-of date, and the four verdict definitions.
17. `trace.md` Step 5: loads `crew-members/clear-questions.md`, asks one question per row that has one, then stops unless go.
18-21. Go-Ahead Payload: one `capture-request:` with the partial and not-started rows and the transcription as asset; a partial archived match carries the REQ id and archive path for `addendum_to` and the commit for `## Prior Implementation`; the ticket id goes in the UR title; trace never picks a number (`trace.md` Rules).
22. `complete` rows stay in the table only (`trace.md` Step 6).
23. Docs: `help.md` menu entry and `prime-do-kanban.md` subcommand list; no guide per Plan D-06. Release: finalization.
- Constraint check: no new frontmatter field, status or UR body section in the diff. "Read-only until go" is stated as unchanged `git status --porcelain` (D-14); the builder's fixture showed empty status before and after a no-go trace.

**Merge seams (integrator):** one conflict, the `SKILL.md` argument hint, resolved as the union. `main.go` auto-merged with REQ-690's `open-work` case: both cases are present and the unknown-subcommand list ends `| open-work | request-commits)`. `help.md` auto-merged under REQ-690's `status` line. I read every merged hunk; none landed in the wrong section. `trace.md` and `trace-reference.md` cite `actions/capture.md` by step name (Step 2, Step 4), not by line, so REQ-688 and REQ-692's capture edits moved nothing they cite.

**Scope:** declared `write_set` = 9 files; touched = the same 9. No drift. Anti-bloat: the builder lists one helper (`requestIdsCreditedByCommit`, the extraction the Plan asked for), one regexp (`requestIdArgumentPattern`), the injected runner parameter (D-12), and one sentence for asks narrower than their REQ (D-10); no extra flag, file or test.

**Re-merge after review (`00116871..a270e712`, 9 files, +301/-17):** the review-fix delta (`775e24f9`) changes `main.go` (comment only), `prime-do-kanban.md` and `trace-reference.md`, all inside the declared write set. No conflict. The classifier had moved to the review phase, so the re-merge gate was recorded with `record-green-gate`.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge `d6aacf29` (1-minute load 1.97 before the run, no other gate running), then `advance REQ-691 --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-10-131527/REQ-691-probe.sh`
**Result:** ✓ Gate exit 0, wall 115 s (queue-kanban 423 Go tests, do-work-cli 913 Go tests, contract and shell suites including `contract-regressions.sh`). GREEN probe exit 0 (`BLOCKED-PROBE-SUCCEEDED`, raw status 0); `scope-drift`, `run-blocked-check` and `green-gate` satisfied; green record at `d6aacf29`.

**Red-green validation:** *(from the builder hand-back proof record; tdd: false, so this is regression and acceptance evidence)*
- RED at base `bd56c4b0`: `test -f skills/do-work/actions/trace.md` exit 1; no `actions/trace.md` row in `SKILL.md`; the probe printed `FAIL: skills/do-work/actions/trace.md is missing`, exit 1 → ✓ probe exit 0 in the worktree and again at merge `d6aacf29`.
- `TestRequestCommitsListsOnlyCommitsTheBoardCredits`: ✗ against a stub (stdout empty, git arguments `[]`) → ✓ PASS. Pins the subject-token and REQ-path credits, the unbracketed id that must not match, the 0 count for bookkeeping, and the exact `git log` argv (no `--since`, no pathspec).
- `TestRequestCommitsRefusesAnArgumentThatIsNotARequestId`: ✗ exit 0 against the stub → ✓ exit 2, nothing on stdout, git not run.
- The existing correlation tests the extraction touches (`...WithoutAPrefix`, `...EveryPrefixTokenInASubject`, `TestRequestPathPatternMatchesOnlyRequestFilesAndRunArtifacts`) still pass.
- `## Red-Green Proof` GREEN cases: routing walk-through (first matching row) for "is this captured? <url>", "do-work trace UR-12", "how much of it is implemented" → trace and "capture-request: add X" → capture; a fixture repo with one archived REQ-007 (2 of 3 criteria tested) and the three-ask spec `SHOP-42` printed `complete` (`170e7ac`), `partial` (2 of 3 + GREEN no) and `not started`, each dated 2026-10-10; `git status --porcelain` empty before and after the no-go stop; the go-ahead capture ran once through `capture-files --commit` and produced UR-002 with the ticket id in its title, REQ-008 (`addendum_to: REQ-007`, Prior Implementation `170e7aca...`) and REQ-009, plus the transcription asset with `Ticket:` / `Source:` / `Fetched:` first lines. The URL-shell and screenshot rules are prose walk-throughs (`trace.md` Step 1 and Rules; `trace-reference.md` Source Handling).

**Repository gate after the review fix:** same argv at re-merge `a270e712` (1-minute load 2.61 before, no other gate running): exit 0, wall 122 s (queue-kanban 423 Go tests, do-work-cli 913); GREEN probe exit 0 (3 s); green record written with `record-green-gate --gate-exit-status 0 -- bash _dev/tests/maintainer-verify.sh` (recorded revision `a270e712`).

**New tests added:**
- `skills/do-work-board/tools/queue-kanban/request_commits_test.go` — the two request-commits tests above.

**Heavy verification plan:**
- Range: 00116871f205cfc8a03b3b8bd57b03dbd2c47b6b..a270e71294b3c1050d21547ecc0e35ab20ced76f (the plan at the first merge `d6aacf29` selected the same three lanes)
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — the five queue-kanban files matched subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — the same five queue-kanban files
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all 9 changed files matched subtree `skills`

*Verified by work action*

## Review

**Overall: 93%** | 2026-10-10T18:18:41Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 93% |
| Test Adequacy | 88% |
| Scope | 97% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. `skills/do-work-board/tools/queue-kanban/main.go:11-19` narrative header comment omits `request-commits` (and already missed `now`); synopsis and unknown-subcommand list are correct — impact-negligible → report only
- F2. `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md:3` "the latter group reads frontmatter or supports the release ritual" no longer fits `request-commits` (it reads git history), and `## Read first` lacks `request_commits.go` — impact-negligible → report only
- F3. `skills/do-work/actions/trace-reference.md:41` "above 0 is code or docs evidence" overstates: the newest row above 0 can be a release commit (REQ-690 → `8ab3d327 ... complete: (0.305.111)`) — impact-negligible → report only
- F4. `skills/do-work/SKILL.md` trace row sits last-but-one, so variants that start with `check`, `review` or `status` route to those earlier actions; the GREEN phrases route correctly; optional move to just below the `validate-feedback` row — impact-user-visible → report only
- F5. `skills/do-work/actions/trace-reference.md:67` "(capture puts it in the UR title)" claims a capture rule `actions/capture.md` does not have; requirement 20's UR-title part rests on the capture agent reading the payload — impact-user-visible → report only
- F7. `do-work/working/REQ-691-trace-coverage-action.md:249,251` Discovered Tasks use `impact-low`, which is not an impact enum value; use `impact-negligible` — impact-negligible → report only
- Nit F6. `skills/do-work-board/tools/queue-kanban/request_commits.go:70-74` a repeated id prints each row twice; trace passes distinct ids — impact-negligible → report only

**Acceptance:** Pass — implementation and integration stages: focused Go tests and vet green, real-repo `request-commits` run (exit 0/1/2 paths) correct, routing walked for the GREEN phrases; trace prose fixture verified from the builder's proof record, not re-run.
**Restatement sweep:** redefined the queue-kanban subcommand set (new `request-commits`) and the core action set (new `trace` in the `SKILL.md` routing, argument hint and help menu); the attribution helper extraction changes no meaning. Swept `main.go` (synopsis, unknown list, narrative comment), `prime-do-kanban.md`, `_dev/primes/prime-kanban-board.md` (write-surface count unaffected, read-only), `skills/do-work-board/actions/board.md`, `docs/board-guide.md`, `justfile.template`, `skills/do-work/actions/help.md`, `README.md`, `docs/work-guide.md` trigger aliases, `_dev/tests/staged-skills-contract.sh` `core_files` (not exhaustive, Plan D-06). Stale: F1, F2.
**Suggested testing:** 5 items
**Follow-ups created:** None (7 findings report only)

**Integrator disposition:** F1, F2, F3 and F5 fixed with the reviewer's replacement text on the builder branch (`775e24f9`), re-merged as `a270e712`, gate and heavy lanes re-run green; no second reviewer, because the delta is the reviewer's own text. F7 fixed in this REQ file (`impact-low` → `impact-negligible` in Discovered Tasks). F4 (row order; the Plan's D-07 placement stays) and F6 (repeated id prints twice; trace passes distinct ids) stay report only. Full report: `do-work/runs/work-2026-10-10-131527/REQ-691-review.md`.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Extracting the attribution rule into one helper (`requestIdsCreditedByCommit`) before adding the second reader, and pinning the exact `git log` argv in the new test. The board's activity lines and `request-commits` cannot credit different commits, and a narrowing `--since` window or `do-work/` pathspec would fail the test.
**What didn't:** The first Commit Evidence gloss called every row above 0 "code or docs evidence"; on this repo the newest such row for a released REQ is the release commit (review F3). Payload wording that says what another action "does" (capture puts the ticket id in the UR title) claimed a rule capture does not have (F5); a hand-off must instruct, not describe.
**Worth knowing:** `request-commits` prints every row twice for a repeated id (F6, report only). A trace row placed just above the capture fallback loses phrases that start with an earlier command word (`check`, `review`, `status`) to those rows (F4, report only). The board binary that trace builds is gitignored, so a no-go trace leaves `git status --porcelain` unchanged.

## Orientation

Now you can ask "is this captured?" or `do-work trace <url|image|text|UR-NNN>` and get one dated coverage verdict per ask before anything is captured again; lives in the core action layer (`actions/trace.md` + `trace-reference.md`, `_dev/primes/prime-action-files.md`) with commit evidence from a new read-only board subcommand (`queue-kanban request-commits`, `_dev/primes/prime-kanban-board.md`). Prime staleness: `prime-do-kanban.md` now lists the subcommand and its file; the other primes' paths still exist.

## Heavy Verification Plan

- Base: 00116871f205cfc8a03b3b8bd57b03dbd2c47b6b
- Target: a270e71294b3c1050d21547ecc0e35ab20ced76f
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — the five queue-kanban files matched subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — the same five queue-kanban files
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all 9 changed files matched subtree `skills`

## Heavy Verification Result

- Target: a270e71294b3c1050d21547ecc0e35ab20ced76f
- Execution revision: a270e71294b3c1050d21547ecc0e35ab20ced76f (detached checkout `.git/work-run-work-2026-10-10-131527/drain-head-REQ-691`, `QUEUE_KANBAN_BROWSER` set to Google Chrome)
- queue-kanban-javascript: executed, exit 0, 7 s
- queue-kanban-browser: executed, exit 0, 74 s
- staged-skills: executed, exit 0, 33 s
- No lane skipped. The first drain at `d6aacf29` also passed every lane (9 s, 86 s, 39 s), all executed.

## Timing

Observed 2026-10-10T18:03:59Z to 2026-10-10T18:18:40Z: 14m 41s total, 15m 23s attributed across 7 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 9m 02s | 4 |
| review | 5m 17s | 1 |
| handback-merge | 1m 04s | 2 |

Slowest stage: review / review agent, REQ-691 00116871..d6aacf29, 5m 17s, outcome success.

Note: no builder-work event was recorded. The hand-back had landed long before this integrator started (dispatch 2026-10-10T13:26:03Z, builder hand-back 13:32:48Z), so per fan-out-reference "Landed hand-back" the builder wait is not timed here.
