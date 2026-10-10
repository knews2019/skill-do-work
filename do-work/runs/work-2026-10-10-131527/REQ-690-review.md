## Review: REQ-690

**Approve** — the status action and `run-status` work end to end on the real queue and match the Plan. One remedy is wrong: C4 sends `failed` and other parked statuses to clarify, which cannot resolve them.
Route C | diff `bd8f3c9e..ada0ccab` (merge `ada0ccab`, includes integrator commit `d6945bef`)

### What's built
- `queue-kanban open-work --format json` prints every open ticket with its column, a placement reason set in the same `bucketColumns` arm, and the board's own last activity. The text digest is byte-identical to the one before this change.
- `do-work-cli run-status` prints one row per claimed, needs-input/blocked, waiting or earmarked REQ: class C1 to C8, ETA (never negative), remedy, run-local files. `do-work status` (`actions/status.md`) wraps it and adds the checks that only this session can do. Routing, help, guide and doc pointers are in place.
- Missing from the REQ: nothing in the core feature. Three narrow gaps: C4's remedy for statuses that are not questions (F1), the builder-branch "age" (F7), and the text remedy lines that do not quote the board's reason (F6).

### Decisions / risks for you
- F3: the D-19 ruling removed the `--run` stat. A typo in `--run` now silently changes classes. On this repo, `--run do-work/runs/typo` turned three C3 rows into C7 rows, and the text named `run-with-recovery` (the command that resets a claim) for live claims. `next_argv` stays the read-only `git log`, so no automation acts on it. Value of keeping D-19: one less guard. Risk: a person may read the C7 text and run the recovery command by hand. The optional fix below does not bring back a refusal.
- F4: C3 also fires while an integrator is running right now (REQ-690 itself: phase `integration`, last activity 6 min ago, `next: do-work run`). This follows the REQ's own table. Risk: someone pastes `do-work run` in a second session and starts a second orchestrator. Spec-level. No fix proposed here.

### Findings

**Important:**
- F1. `skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go:329`: C4 is "any row in the Needs input · Blocked column". That column also holds `failed`, `blocked-archive-collision`, `blocked-dependency-cycle` and unrecognized statuses (`queue-kanban/model.go:1099` `isNeedsInputOrBlockedStatus`, plus the `default:` arm). Those rows get `next_argv ["do-work","clarify"]`, but `actions/clarify.md` does not handle any of those statuses (grep finds no match). The REQ table limits C4's clarify remedy to "`pending-answers` or `blocked` with every dependency met". The Plan widened C4 to the whole column without recording that as a decision. — impact-user-visible → report only
  Fix (verbatim). Replace lines 329-331:
  ```go
  	case request.Column == "needs-input-or-blocked":
  		record.Class, finding.NextArgv = "C4", []string{"do-work", "clarify"}
  		row.remedy = "it waits on a person; run `do-work clarify`"
  ```
  with:
  ```go
  	case request.Column == "needs-input-or-blocked" && (request.Status == "pending-answers" || request.Status == "blocked"):
  		record.Class, finding.NextArgv = "C4", []string{"do-work", "clarify"}
  		row.remedy = "it waits on a person; run `do-work clarify`"
  	case request.Column == "needs-input-or-blocked":
  		record.Class, finding.NextArgv = "C4", []string{"do-work", "forensics"}
  		row.remedy = "status " + request.Status + " is not a question clarify answers; run `do-work forensics` to see what holds it"
  ```
  Pin it in `run_status_test.go` (append):
  ```go
  // A failed REQ also sits in Needs input · Blocked, but clarify answers
  // questions and never resolves a failure, so its remedy is the read-only
  // diagnostic instead.
  func TestRunStatusFailedRequestNamesForensicsNotClarify(t *testing.T) {
  	repositoryRoot, boardFacts := writeRunStatusFixture(t, fixtureRequest{id: "REQ-809", status: "failed", column: "needs-input-or-blocked", placementReason: "Needs input · Blocked"})
  	finding := findingFor(t, runFixture(t, repositoryRoot, fixtureNow, "--board-facts", boardFacts), "REQ-809")
  	if finding.Code != "C4" || !reflect.DeepEqual(finding.NextArgv, []string{"do-work", "forensics"}) {
  		t.Fatalf("finding = %+v", finding)
  	}
  }
  ```
  Then change `skills/do-work/actions/status.md:47`'s last cell from `` `do-work clarify` `` to `` `do-work clarify` for `pending-answers` or `blocked`; `do-work forensics` for any other status ``, and `skills/do-work/docs/status-guide.md:27`'s last cell from `` `do-work clarify` `` to `` `do-work clarify`; `do-work forensics` when the status is failed or a collision ``.

**Minor:**
- F2. `skills/do-work/actions/status.md:33-36`: the block ends with `rm -f`, so it always exits 0. Line 39 still says "Exit 0 means the report was produced" and treats exit 1 as a refusal. That exit status can never be seen (prime-shell-commands § Unchecked Exit Status Reads as Content). An agent still sees the error text, so nothing unsafe follows. — impact-negligible → report only
  Fix (verbatim). Replace the block with:
  ```bash
  board_facts="$(mktemp)"
  "<skill-root>/../do-work-board/tools/queue-kanban/queue-kanban" open-work --format json --repo-root "<project-root>" >"$board_facts" \
    && "<skill-root>/tools/do-work-cli.sh" --repo-root "<project-root>" run-status --board-facts "$board_facts"
  report_status=$?
  rm -f "$board_facts"
  (exit "$report_status")
  ```
- F3. Result of D-19 (`run_status.go:152`, `resolveRunDirectory`): a `--run` path that does not exist reads as "no run". Hand-backs then count as absent, so live claims older than 3 h drop from C3 to C7, and the text never says which run directory it read. — impact-user-visible → report only
  Optional fix with no refusal. After line 152 (`report.RunDirectory = displayPath(repositoryRoot, runDirectory)`), insert:
  ```go
  		if info, statError := os.Stat(runDirectory); statError != nil || !info.IsDir() {
  			report.RunDirectory += " (not found; run fields absent)"
  		}
  ```
  and in `run_status_render.go`, just before the final `output.WriteString(nextLine + "\n")`, insert:
  ```go
  	if report.RunDirectory != "" {
  		fmt.Fprintf(&output, "run: %s\n", report.RunDirectory)
  	}
  ```
  (The `--watch` branch returns earlier and is unchanged.)
- F4. C3 cannot tell "hand-back landed, waiting" from "integrator running now". It recommends `do-work run` for a REQ that is being integrated (seen live on REQ-690). This matches the REQ's table. A later REQ could use `last_activity_phase == "integration"` to say "integrating". — impact-user-visible → report only
- F5. Test gaps. Nothing pins C2 (including "no activity record goes to C2, never C1"), the C1/C2 20-minute boundary, C4, C8, or the precedence order (for example C3 over C7, which the live run depends on). The nine tests that exist do pin the GREEN items named in the REQ, and none is decorative. — impact-negligible → report only
- F6. `run_status.go:337`: the C6 remedy line in the text never quotes the board's Earmarked reason and never names the session it is earmarked for. Only the JSON evidence does. The REQ's remedy cell says "quote the board reason". — impact-negligible → report only
  Fix (verbatim). Replace line 337 with:
  ```go
  		row.remedy = request.PlacementReason + " Run `do-work run " + request.ID + "` to run it by name"
  ```
- F7. Item 6 asked for the builder branch tip *age*. `builderBranches` stores the tip instant in the local offset (`2026-10-10T20:41:54+03:00`), computes no age, and the text table does not show it. The value is for display only, and no class reads it. — impact-negligible → report only
- F8. Restatement: `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md:3` (subcommand list) and `:7` (`open_work.go` "the headless in-flight digest") do not mention `--format json`. The builder already listed this in Discovered Tasks. — impact-negligible → report only
  Fix: on `:7`, after "shows nothing terminal", append `; `--format json` prints the same open tickets with column, placement reason and last activity for do-work-cli run-status`.
- F9. Restatement: `stuck` now routes to status, but `skills/do-work/actions/forensics.md:10` ("User suspects something is stuck, broken, or producing confusing results") and `README.md:176` ("Run `do-work forensics` to diagnose stuck or failed work.") still send stuck questions to forensics. — impact-negligible → report only
  Fix: `forensics.md:10` becomes `- User suspects something is broken or producing confusing results (for how long in-flight work has been quiet, see `actions/status.md`)`. In `README.md:176`, the last sentence becomes `Run `do-work status` to see how long in-flight work has been quiet, and `do-work forensics` to diagnose failed or broken work.`
- F10. Anti-bloat: `verificationArgv` (`run_status.go:487`, called at `:143`, D-17) emits an argv that names the `--board-facts` temp file. The only caller (`actions/status.md`) deletes that file in the same block, so the argv can never be re-run. The REQ does not name it. — impact-negligible → report only
  Fix: delete line 143 and the `verificationArgv` function (lines 487-493).

**Nit:**
- F11. Anti-bloat inventory of other items the REQ did not name. Each is small, and none breaks a Constraint: `run_status_render.go` (D-05/D-18, needed by the Plan's size rule); rows sorted by class precedence (D-09); `requestIDPattern` validation of `--req`; the `RunStatusResult.RunDirectory` field; the `classNames` map; placement sentences for Ready, Claimed, Needs input and unrecognized (D-13; the Ready sentence has no reader yet); `--watch` header class counts (D-16); the 160-character cap on a run-local file's first line; the `refusal()` helper. Keep all of them. — impact-negligible → report only

Constraint check, all clean: no lock, heartbeat, lease, PID check or process registry in the Go code (`grep -nE 'os\.Remove|syscall\.Kill|FindProcess'` over `internal/runstatus/*.go` exits 1, and the only `exec.Command` is `git for-each-ref`). `nextselection` changes only by renaming `frozenEstimate` to `FrozenEstimate`. No scheduling, `advance` or threshold change. No new REQ field or status. No `web/` or payload change: `RequestTicket.PlacementReason` is not serialized into `generatedBoardData`.

Item checks the coordinator asked for:
1. `classifyRow` precedence is C8, C3, C4, C5, C6, C7, C2, C1, as the Plan says. The ETA cannot go negative: `estimateRemaining` returns nil plus "over estimate by N" when elapsed > p50, and `p50 - elapsed >= 0` otherwise. C7's `next_argv` is `git log --full-history -- <path>`. `finding.AffectedPaths[0]` cannot panic: C7 needs `MinutesSinceClaim != nil`, which needs `record.ClaimedAt`, which is set only when the REQ file was found, and that is also when `AffectedPaths` is set. Confirmed with a ghost REQ that has no file (classed C1, no panic). After D-19: a facts file with no threshold reads as 0, so every claimed row without a hand-back becomes C7 (read-only argv). The board always writes 180. A missing `--run` is F3.
2. `PlacementReason` is set in the same `bucketColumns` arm as the column, in all six open arms. The text digest is byte-identical, checked with `cmp` of `bd8f3c9e` vs `ada0ccab` binaries on the real queue, and `--format text` matches the default. `--format yaml` exits 2 with a message. Activity comes from `attachRequestActivity` over one `readWorktreeAgentGitState`, the same call `generate.go:683-685` makes. There is no second correlation.
3. `status.md` follows the template: blockquote with the package reason, When to Use, Input, Steps, and an earned Rules section (destructive-next-argv REQ-629; lock may belong to another checkout). The cross-package path depth is right. The shell block has the exit-status issue (F2). The routing row sits directly above clarify, so also above roadmap. `blocked` still routes to clarify. The argument hint keeps REQ-692's untouched hint and adds `status [REQ] [--watch]`. REQ-691 has not merged yet.

### Requirements Checklist
- [x] 1 newest `work-*` run unless `--run`; `--req` filter — delivered
- [x] 2 one record per claimed / needs-input / blocked / waiting REQ with status, claimed_at, blocked_by, blocked_at, depends_on, assigned_to, open questions — delivered
- [x] 3 column and placement reason from the board's own partition — delivered
- [x] 4 last activity, kind, phase, minutes since, from the board's correlation — delivered
- [x] 5 manifest dispatch instant, hand-back presence, landed status (as the verbatim manifest row) — delivered
- [~] 6 builder branch tip age — partially delivered (tip instant, no age, F7)
- [x] 7 ETA = p50 - elapsed, "over estimate by N min", never negative — delivered
- [x] 8 run-local files with path, age, first line, label, never deleted — delivered
- [~] 9 one class + `next_argv` per row in `CommandFinding` shape — delivered, except C4's remedy for non-question statuses (F1)
- [x] 10-13 action: CLI first, session-local checks labelled, `ps` lookup and printed `rm`, table + remedy lines + last `next:` line — delivered
- [x] 14 routing row above clarify and roadmap, `blocked` stays on clarify — delivered
- [x] 15 `--watch` compact, `/loop 15m /do-work status --watch` documented, no state — delivered (14 lines for 3 rows live; test pins ≤20 for 8)
- [x] 17 help, guide, board.md, fan-out-reference pointers — delivered; the release belongs to the integrator (N/A here)
- [x] Constraints (no PID/lock Go, no scheduling change, no new field, no UI change) — met

### Acceptance Testing

**Result: Pass** (implementation and integration stages)
- Built both binaries from `git archive ada0ccab` into a scratch dir. `go test -count=1` passes for `internal/runstatus`, `internal/resultmodel`, `internal/nextselection`, `cmd/do-work-cli`, and the queue-kanban `TestOpenWork` subset. `go vet` and `gofmt -l` are clean.
- Real queue, read-only: `open-work --format json` → `run-status` gives 3 rows (REQ-689, 690, 691, all C3), 3 run-local lines, `next: do-work run`, exit 0. `--watch` gives 5 lines. The JSON `--req REQ-690` output has the `CommandFinding` shape (code `C3`, `affected_paths`, evidence with the manifest row, the builder branch, and the read-only `next_argv`). `git status` is unchanged after every run.
- Edge cases: an empty facts file (`{}`) gives "No claimed, blocked, waiting or earmarked REQs." with exit 0. A missing facts file is refused with exit 1. A ghost REQ with no file is C1 with no panic. A typo `--run` → C7 rows (F3).

### Suggested Additional Testing
- Deployment: unassessed. Run the `staged-skills` and `installer` heavy lanes listed in `## Testing` before release.
- Live acceptance: unassessed. In a consumer install with Go, say "status and ETA please" and check that the router picks `actions/status.md`, the board builds, and the last line is a `next:` command.
- Manual: run Step 3 with a real spawned agent that has been quiet for over 20 min. Check that the transcript age and the nudge line appear, labelled "this machine, this session, now".
- Manual: put a `full-gate.lock` whose owner pid is dead in a run dir. Check that the action prints `rm -- '<path>'` and does not run it.
- Edge: one REQ with `status: failed` in the queue. Check the C4 remedy after F1.

### Scores (on the record, not the headline)

**Overall: 87%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 90% | Items 6 and 9 partial (F7, F1) |
| Code Quality | 85% | Wrong C4 remedy for failed and collision statuses; masked exit status in the action shell |
| Test Adequacy | 80% | GREEN items pinned; C2, C4, C8 and precedence not pinned |
| Scope | 95% | One documented extra file (D-18); one dead helper (F10) |
| Risk | Low | Read-only; every `next_argv` is safe to follow literally |
| Acceptance | Pass | Implementation and integration stages exercised |

### Follow-ups created
None (11 findings report only)

## Review

**Overall: 87%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 90% |
| Code Quality | 85% |
| Test Adequacy | 80% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 `run_status.go:329` C4 gives `do-work clarify` to every Needs input · Blocked row, including `failed`, `blocked-archive-collision`, `blocked-dependency-cycle` and unrecognized statuses, which clarify does not handle. The REQ limits that remedy to pending-answers or blocked. Verbatim fix in the report: forensics for the other statuses, plus one test. — impact-user-visible → report only

**Minor findings:**
- F2 `actions/status.md:33-36` ends with `rm -f`, so the block always exits 0 while line 39 keys on the exit status. Fix: capture `$?` before `rm` and end with `(exit "$report_status")`. — impact-negligible → report only
- F3 D-19 result: a mistyped `--run` silently drops run fields. Live claims over 3 h move from C3 to C7, and the text names `run-with-recovery`. The text never says which run directory was read. Optional non-refusal fix in the report. — impact-user-visible → report only
- F4 C3 recommends `do-work run` while an integrator is running now (seen on REQ-690). This is per the REQ table. — impact-user-visible → report only
- F5 No tests for C2 (including "no activity record goes to C2"), the 20-minute boundary, C4, C8, or the class precedence. — impact-negligible → report only
- F6 `run_status.go:337`: the C6 text remedy does not quote the board's Earmarked reason or its owner. — impact-negligible → report only
- F7 Item 6: builder branch shows a tip instant in the local offset, not an age, and is not in the text table. — impact-negligible → report only
- F8 `queue-kanban/prime-do-kanban.md:3,7` do not mention `open-work --format json`. — impact-negligible → report only
- F9 `actions/forensics.md:10` and `README.md:176` still send "stuck" questions to forensics. — impact-negligible → report only
- F10 Anti-bloat: `verificationArgv` (`run_status.go:143,487-493`, D-17) names a temp file the action deletes. Delete it. — impact-negligible → report only
- F11 (Nit) Anti-bloat inventory of small unrequested items, all kept: render file split, class-ordered rows, `--req` regex, `RunDirectory` field, `classNames`, the four extra placement sentences, watch class counts, the 160-character cap, `refusal()`. — impact-negligible → report only

**Acceptance:** Pass. Implementation and integration stages: focused Go tests, vet and gofmt pass. On the real queue, the board JSON → run-status flow gives 3 C3 rows, `--watch` gives 5 lines, the text digest is byte-identical, and the tree stays unchanged.
**Restatement sweep:** redefined: the `status`/`stuck`/`eta` routing words (stale: `forensics.md:10`, `README.md:176`, F9); the `queue-kanban open-work` flag set (stale: `prime-do-kanban.md:3,7`, F8; `board.md`, `board-guide.md` and `justfile.template` agree); the fan-out progress-question rule (no other restatement). Claim staleness is not redefined: the threshold is read from `staleClaimThreshold`.
**Suggested testing:** 5 items
**Follow-ups created:** None (11 findings report only)

*Reviewed by review-work action*
