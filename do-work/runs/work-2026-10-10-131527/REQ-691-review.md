# Review: REQ-691 (do-work trace action and the read-only queue-kanban request-commits subcommand)

**Approve**: the build delivers the trace action, its routing and the board command as the Plan said. Seven small findings, none critical, all report only.
Route C | range `00116871..d6aacf29` (builder commit `e8e4e9c8`, merge `d6aacf29`)

Mode: orchestrated, independent reviewer. Read-only on the repo except this file. The working tree held only the integrator's edit to the REQ file, before and after this review.

## What's built

- `do-work trace <url | image | pasted text | UR-NNN>` exists (`skills/do-work/actions/trace.md`, 49 lines, six steps and three rules; detail in `skills/do-work/actions/trace-reference.md`). It loads the prompt-injection guardrail, transcribes the source to a `mktemp -d` file outside the tree, splits it into asks, searches per ask, prints one seven-column table with a dated verdict per ask, asks one question per open row, and stops. Only on go does it hand the `partial` and `not started` rows to one `capture-request:`.
- `queue-kanban request-commits [--repo-root DIR] REQ-NNN...` exists. It runs one full-history `git log` with no pathspec, credits commits through the one shared helper `requestIdsCreditedByCommit`, and prints TSV with a paths-outside-`do-work/` count. Exit codes are 0, 1 and 2.
- Still missing: the release (the integrator's finalization does it, per requirement 23 and the Plan).

## Decisions / risks for you

- None that need a choice. F4 (the routing row sits last-but-one, so a phrase that starts with `check`, `review` or `status` goes to that earlier action) is the one behaviour a user could notice. The fix is optional and is described under F4.

## Findings

**Important:** None.

**Minor:**
- F1. `skills/do-work-board/tools/queue-kanban/main.go:11-19`: the narrative header comment lists the subcommands in prose and does not name `request-commits` (it also already missed `now`, a gap that existed before this REQ). The synopsis and the unknown-subcommand list below it are correct. Replace lines 16-19 (from `// request) — plus the read-only` through `// invariants otherwise checked by hand.`) with the block below. — impact-negligible → report only

  ```go
  // request) — plus the read-only `frontmatter` field reader, the read-only
  // `request-commits` commit evidence for do-work trace, `now` (the Timestamp
  // rule's stamp), and three release-ritual subcommands: `next-req` atomically
  // reserves a number, `next-version` allocates a version, and `verify` checks
  // the cross-file invariants otherwise checked by hand.
  ```
- F2. `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md:3`: `request-commits` was appended to "the latter group", which the next sentence says "reads frontmatter or supports the release ritual and Timestamp rule". `request-commits` does neither; it reads git history. The `## Read first` list also names a file for every other non-board subcommand but not `request_commits.go`. Replacement for the sentence on line 3: `The latter group reads frontmatter or git history, or supports the release ritual and Timestamp rule, rather than rendering the board.` New line after line 8 (`frontmatter_cli.go`): `` - `request_commits.go` — `request-commits`: read-only TSV of every commit the board's attribution helper (`requestIdsCreditedByCommit`, `activity_correlation.go`) credits to each named REQ; the commit evidence for `../../../do-work/actions/trace-reference.md` `` — impact-negligible → report only
- F3. `skills/do-work/actions/trace-reference.md:41` says a `paths_outside_do_work` above 0 "is code or docs evidence" and to cite the newest such row. In a repo whose completion commit also releases, the newest row above 0 is the release commit. On this repo, `request-commits REQ-690` puts `8ab3d327 [REQ-690] complete: ... (0.305.111)` first (CHANGELOG, VERSION, version.md), not the implementation commit. The cited hash is still a commit credited to the REQ, so the evidence stays valid. Only the gloss overstates. Replacement for `above 0 is code or docs evidence.`: `above 0 touched files outside \`do-work/\`: code, docs, or, where the completion commit also releases, the changelog and version files.` — impact-negligible → report only
- F4. Routing: the trace row in `skills/do-work/SKILL.md` sits directly above the capture fallback (Plan D-07), so every earlier row wins first. The three GREEN phrases ("is this captured? <url>", "do-work trace UR-12", "how much of it is implemented") and `capture-request:` route correctly, and REQ-692's merged `validate-feedback` row does not collide. A natural variant that starts with an earlier command word does go elsewhere: "check if this ticket is already implemented <url>" → `verify-requests` (`check`), "review this spec, is it captured?" → `review-work`, "status of this ticket <url>" → `status`. The REQ only required the row to be above capture, so this is not a miss. Optional fix if wanted: move the trace row up to directly below the `validate-feedback` row (above `verify`). Its triggers are `trace` plus multi-word phrases that no earlier row uses as its command word. Risk of the move: low. A bare `check` or `review` still reaches its own row, because none of trace's phrases is a single command word except `trace`. — impact-user-visible → report only
- F5. `skills/do-work/actions/trace-reference.md:67`: "the external ticket id (capture puts it in the UR title)" describes capture behaviour that `actions/capture.md` does not prescribe. Capture has no ticket-id rule, so requirement 20's "the UR title carries it" holds only if the capturing agent acts on this payload. It did in the builder's fixture (`UR-002` title `SHOP-42: ...`). Replacement: `- when the source has one: the external ticket id, with the instruction that the UR title starts with it (capture has no ticket rule of its own), and the source URL and fetch date (already the transcription's first lines).` — impact-user-visible → report only
- F7. `do-work/working/REQ-691-trace-coverage-action.md:249` and `:251` (Discovered Tasks, copied from the hand-back) end with `impact-low → report only`. `impact-low` is not in the impact enum (`impact-critical | impact-user-visible | impact-rule-change | impact-negligible`, `actions/work-reference.md` Request File Schema). An unrecognized value reads as `impact-user-visible`. Replace `impact-low` with `impact-negligible` on both lines. — impact-negligible → report only

**Nit:**
- F6. `skills/do-work-board/tools/queue-kanban/request_commits.go:70-74`: a repeated id (`REQ-7 REQ-7`) prints every row twice. Trace passes distinct matched ids, so this does not affect trace. — impact-negligible → report only

## Requirements Checklist

Cross-referenced against Plan D-01..D-08, builder Decisions D-09..D-17, and UR-153 (the 2026-10-10 decision record: keep one REQ, drop the Notion adapter and the PM Ticket body section).

- [x] 1. `actions/trace.md` + `actions/trace-reference.md`: delivered.
- [x] 2. Routing row above capture, the listed triggers, `coverage` only with a target (D-09), hint `trace <url|image|text|UR-NNN>`, `capture-request:` still on capture: delivered. The hint union with REQ-690's `status [REQ] [--watch]` is intact (F4 covers the row-order side note).
- [x] 3. Prompt-injection loaded before the source (`trace.md` Step 1): delivered.
- [x] 4. UR id (input.md + assets), image described before the split (cites `actions/capture.md` Step 4, which is still correct on main after REQ-688 and REQ-692), pasted text as is: delivered.
- [x] 5. URL: web fetch first, a shell result asks for a paste or a screenshot, never judged: delivered (Step 1, Source Handling, Rules).
- [x] 7. Transcription kept as a file outside the tree, becomes an asset only on go: delivered (D-08).
- [x] 8. Atomic asks A1..., source numbering kept: delivered.
- [x] 9-11. UR inputs, REQs by intent across queue/working/archive with every candidate listed, Implementation Summary and `commit:`: delivered.
- [x] 12. Board matching reused through `request-commits`. `Traced-to:` is covered by the subject token (D-02). The merge range is gone since 0.305.88 (D-03). The D-04 fallback marks a row `unverified`: delivered.
- [x] 13. Code and tests by name, GREEN-when assertion as evidence: delivered.
- [x] 14-16. Seven columns, a three-row example, the four verdicts, the as-of date, one evidence pointer, `N of M + GREEN yes|no` (D-11): delivered. See F3 for the Commit-column gloss.
- [x] 17. clear-questions loaded, one question per open row, then stop: delivered.
- [x] 18. One capture on go, partial and not-started rows only, transcription as asset: delivered (D-17).
- [x] 19. A partial archived match takes the addendum path with the commit for `## Prior Implementation` (cites `actions/capture.md` Step 2, still correct): delivered.
- [~] 20. Ticket id in the UR title and in the asset's first lines; no new UR body section: the asset header is delivered. The UR title depends on the capture agent following the payload (F5).
- [x] 21. Trace never picks a number: delivered (Rules + payload).
- [x] 22. `complete` rows cited, never captured: delivered (Step 6).
- [x] 23. Docs: help menu + board prime, with no guide by Plan D-06. Release: N/A here (integrator finalization).
- [x] Constraints: no new frontmatter field, status or UR body section in the diff. Read-only until go is stated as unchanged `git status --porcelain` (D-14), which equals empty in a clean fixture. No fetched row is `complete` without fetched text.
- [x] Acceptance criterion 7 (request-commits credits exactly the board's rule; correlation unchanged): delivered and run (below).

## Code Review

- `activity_correlation.go`: the extraction moves the two loops verbatim into `requestIdsCreditedByCommit`. `correlateCommitsToRequests` iterates the same map. Behaviour is unchanged: the existing `RequestActivity*`, `RequestPathPattern*` and `Correlat*` tests pass.
- `request_commits.go`: flags first, then ids validated as `^REQ-\d+$` before git runs. The repo root is resolved through the existing `resolveRepoRootOrDefault`. Git failure gives exit 1 with a stderr line. `committed_at` is re-printed with `time.RFC3339`, which keeps git's offset. Everything follows the `frontmatter` command's writers-in, exit-code-out shape. The comment cites `../../../do-work/actions/trace-reference.md`, which resolves.
- `main.go`: the switch case, synopsis line and unknown-subcommand list are correct, and the auto-merge with REQ-690's `open-work` is clean. The narrative comment is stale (F1).
- Prose: `trace.md` is short and points to named reference sections. The Commit Evidence block reuses forensics Check 14's `<suite-root>` build-and-run shape. The verdicts are defined once (reference only). Every citation into `actions/capture.md` is by step name, and Step 2 and Step 4 still hold the cited content on main.
- Anti-bloat count (things the REQ did not name): one regexp `requestIdArgumentPattern` (needed for D-05's exit 2), the injected runner parameter (D-12, so tests stub git without a package variable), and one sentence for an ask narrower than its REQ (D-10, needed for the GREEN fixture). There is no extra flag, no extra file and no decorative test. The two tests are the Plan's two, and each comment names the failure it pins. None is a finding.

## Acceptance Testing

**Result: Pass.** Stages covered: implementation and integration (the merged tree on main).

- `go test -count=1 -run 'RequestCommits|RequestActivity|RequestPathPattern|Correlat' .` → ok. 10 tests pass (both new ones plus the existing correlation tests). One JavaScript browser test skipped because it is a heavy lane.
- `go vet` → clean. `gofmt -l` → nothing.
- A binary built into the session scratchpad (outside ROOT, removed afterwards) and run against ROOT:
  - `request-commits REQ-690 REQ-691` → exit 0, 0.77 s wall, newest first. It credits the subject-token commits and the path-only `[UR-153]` and `[UR-145]` bookkeeping commits with 0. It shows `e8e4e9c8` (the REQ-691 builder commit) at 9 paths and the `[REQ-690]` code commits above 0.
  - No ids → exit 2. A flag after an id → exit 2 naming `"--repo-root"`. A repo root with no git → exit 1 `git log failed`. An unknown id → header only, exit 0. An unknown subcommand → the list includes `request-commits`.
- Routing walked by first match on the merged `SKILL.md` for the three GREEN phrases and `capture-request:` (correct). See F4 for variants.
- The trace GREEN fixture (complete / partial 2 of 3 / not started, an empty `git status --porcelain` after no-go, the go-ahead capture producing an addendum REQ with Prior Implementation and the ticket id) was not re-run here. I read the builder's proof record (hand-back §3) and checked it against the prose it exercised.

## Suggested Additional Testing

- Live acceptance: unassessed. In a consumer repo, run `do-work trace <a real PM ticket URL>`. Check that a shell fetch triggers the paste or screenshot ask and that a full fetch yields a dated table.
- Deployment: unassessed. After release, confirm the installed suite carries `trace.md`, `trace-reference.md` and the board's `request_commits.go`, and that `go build` succeeds from the installed `do-work-board/tools/queue-kanban`.
- Edge: answer the Step 5 go question as a plain typed "go" on a platform without an ask-user tool. Confirm it is taken as the trace answer, not as the `go` trigger of the run row.
- Edge: a trace whose ask matches an in-flight (`working/`) REQ. The payload names an "archive path" only. Confirm capture's Step 2 working-tier addendum path still fires.
- Regression: the board's activity lines on a served board (heavy `queue-kanban-browser` lane) after the helper extraction.

## Scores (on the record — not the headline)

**Overall: 93%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | All delivered. Requirement 20's UR title rests on payload wording (F5). Release is the integrator's. |
| Code Quality | 93% | Clean extraction and command. Stale comment (F1). Nit F6. |
| Test Adequacy | 88% | Two focused tests plus the existing correlation tests. The exit-1 git-failure path is untested. Prose GREEN is by walk-through and the builder's fixture. |
| Scope | 97% | Exactly the 9 declared files. Decisions D-09..D-17 cover every visible choice. |
| Risk | Low | Read-only command. The extraction does not change behaviour. |
| Acceptance | Pass | Implementation and integration stages. |

### Follow-ups created
- None (7 findings report only)

## Self-validation

I re-checked the two places most likely to be wrong. (1) The capture citations: confirmed against the post-REQ-688/692 `capture.md` headings (Step 2 Check for Duplicates, Step 4 Handle Screenshots). (2) The extraction: diffed line by line, and the existing correlation tests were run, not assumed. No new issue was found. I did not re-run the trace prose fixture (see Acceptance).

## Review

**Overall: 93%** | <REVIEW_AT>

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

*Reviewed by review-work action*
