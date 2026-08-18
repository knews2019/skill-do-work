---
session_ended: 2026-08-18T14:07:35Z
last_completed: REQ-244
queue_state: 7 pending, 2 pending-answers, 0 blocked, 0 blocked-archive-collision, 0 blocked-dependency-cycle, 0 in-progress
reqs_processed_this_session: 5
session_depth: heavy
---

# Session Checkpoint

## In Progress (interrupted)

- REQ-248: Anchor the Durations day buckets to UTC midnight so Panel B stays on canvas — claimed 2026-08-18T16:09:27Z — writer: vm:/home/user/skill-do-work
- REQ-249: Decide the cross-package citation path form and sweep to match — claimed 2026-08-18T16:09:27Z — writer: vm:/home/user/skill-do-work

## Completed This Session

- REQ-241: Reconcile the Durations label metrics with the rendered face (Route B) — commit `90c74b7`, shipped as **0.212.3**
- REQ-243: Check that shipped markdown pointers actually resolve (Route B) — commit `37d7729`, shipped as **0.212.4**
- REQ-245: Name fabricated stamps in the board's future-stamp warnings (Route A) — commit `23bad9d`, shipped as **0.212.5**
- REQ-242: Stop Panel B's slowest-day annotation colliding with its title (Route B) — commit `48263dd`, shipped as **0.212.6**
- REQ-244: Cite the Timestamp rule at every timestamp write site (Route C) — commit `f733365`, shipped as **0.212.7**

Every hash was confirmed with `record-commit-hash.sh --verify`. `maintainer-verify.sh` exits 0 at every commit boundary and on the final tree. All five ran as worktree builders in two waves; every worktree and `worktree-agent-*` branch was removed with `git worktree remove` / `git branch -d` (never `-D`), the worktrees parent directory is gone, and `git worktree list` shows only the main tree.

**All five were reviewed independently and all five came back PASS-WITH-FINDINGS** (80%, 88%, 88%, 88%, 66%). Every one was remediated on its builder's own branch and re-merged before shipping. No REQ shipped on its first pass.

## Still Queued

Nine, none claimed:

- **REQ-246** / **REQ-247** (pending): mechanical no-LLM timestamp repair — the SessionStart hook repairer and the git-commit-time archive auditor. REQ-247 depends on REQ-246. Captured before this session and deliberately left out of its assignment.
- **REQ-248** (pending): anchor the Durations day buckets to UTC midnight. Panel B's leftmost bar sits in the axis gutter on this repository's own board, and a one- or two-day board renders Panel B **entirely off-canvas** (measured `x=-3330`). The highest-value item in the queue.
- **REQ-249** (pending-answers): decide the cross-package citation path form. Two incompatible readings coexist; REQ-244 added eleven more in the prescribed form while the question stayed open.
- **REQ-250** (pending): the four remaining markdown link-checker limits.
- **REQ-251** (pending): two stale copies of the future-stamp message outside REQ-245's write set.
- **REQ-252** (pending): record the browser with every measured-face number.
- **REQ-253** (pending-answers): the two clock-write shapes the Timestamp rule governs with neither paragraph.
- **REQ-254** (pending): let `qualify.sh` tell a check's own output from leftover instrumentation.

## Session Notes

- **Every REQ this session shipped a mechanism that looked like it closed a class and closed only the instance.** This is the session's single finding, and it recurred five times independently. REQ-241's width constant was calibrated from a sweep that held one of two unbounded parameters fixed. REQ-243's heading walk relied on an invariant stronger than the one its fixture locked. REQ-242's "x-free" test sampled six x values, and a mutant banded on x reproduced the original defect **at the real board's own slowest-day position** while the suite stayed green. REQ-244's checker keyed on a literal bracket span containing one word, so it locked in exactly the drift it had removed. REQ-245's badge guard asserted a phrase was *absent*, and swapping the whole string for a wrong one passed. In every case the code was right and the net had a hole, and in three of the five the hole was precisely where the real data lives.

- **The reviewers earned their cost this session; the builders' own evidence did not catch any of the five.** All five reviews were adversarial and all five found something real. Two findings were reproduced by mutation against the shipped code rather than argued from reading, which is what made them undeniable.

- **A builder disproved the orchestrator's preferred fix with measurement, and was right.** REQ-241 was told to clamp the hour count so the label space would be finite. It measured that with hours clamped to `999h+`, a label carrying a five-digit id still exceeds the constant, because the id is the *other* unbounded parameter. (The measured worst-case labels are quoted in REQ-241's archived Lessons Learned; they are deliberately not repeated here, because a synthetic five-digit id in this file reads to `queue-kanban verify` as a REQ that does not exist.) It shipped an amortization argument instead — only digits repeat without limit, every wide fixed character is amortized away, so per-character width converges to a pure digit run's 7.14 and cannot pass it. **"Complete sweep" is a property of the argument, not the sample size**: its first sweep was 10 000 strings and wrong, its second 280 800 and still would have been.

- **The predicted merge conflict never happened, and a real one arrived where nothing was watching.** `generate_test.go` merged clean every time — three conflicts occurred last session and zero this one, because REQ-245 was steered to `timestamp_test.go` and REQ-242 was given sole ownership of the file for the wave. Instead REQ-241 and REQ-242 each measured the same 12px axis-title face on **different Chromium builds** (12.0372 vs 11.2300), rounded up, and declared the same constant in different files of one package. Git saw nothing — the edits never touched adjacent lines — and it surfaced as `redeclared in this block`. Resolved to the larger value inside the merge commit, and **both sides' tests were run by name** after resolving, not merely compiled.

- **A measured face is per-browser, and the numbers we record are not portable.** REQ-241's D-03 budget of 1.364 units above Panel B's title measures **0.185 units** on Chromium 146 — still positive, still zero intersections, but roughly 7× thinner than recorded. That is REQ-252.

- **A builder corrected its own record unprompted, and the correction was the point.** REQ-244's original hand-back quoted a GREEN transcript from a prototype under a sentence asserting it came from the committed checker. The committed checker printed nothing on success, so it could not have produced that line. In a REQ about agents writing values they did not read, that is the REQ's own failure mode appearing inside its evidence. The fix — make the check print its counts — then tripped `qualify.sh`'s debug-artifact scan, which is REQ-254.

- **`qualify.sh` was overridden once, deliberately and on the record.** Its debug-artifact heuristic cannot distinguish a check's own success output from leftover instrumentation. The override is written into REQ-244 with the diff quoted, because an unrecorded override teaches the next reader that a qualify FAIL can be waved through.

- **The environment failed mid-session and it was not the code.** `maintainer-verify.sh` failed with 36 `No space left on device` errors inside the installer probe. The volume is at 100% (263 MiB free of 123 GiB). Freed roughly 350 MB of this and prior sessions' scratch — including a 168 MB Chromium a reviewer downloaded and an **11 MB compiled `queue-kanban` binary a builder left in the source tree** (gitignored, so no git-hygiene failure, but it inflated every install-probe copy). Verify then exited 0 with zero space errors. **The machine still needs real space.**

- **A builder can write outside its worktree without ever issuing a write.** The Playwright MCP drops a `.playwright-mcp/` directory into whatever it considers its working root, which was the main tree. It is gitignored and already held 36 files from sibling sessions. REQ-245's builder removed only its own two files by timestamp and left the rest — the correct call. The brief template now states this as an explicit exception rather than leaving builders to assume it.

- **`queue-kanban verify` reports a live builder's worktree as `merged-worktree-leftover [fixable]`.** True and dangerous to act on: the branch tip is contained in HEAD because the orchestrator merged it, while the builder is still adding commits to it. Running `do-work cleanup` on that advice mid-remediation would delete live work.

- **Estimator calibration this session:** five REQs, estimate vs actual active minutes — 35/28, 20/20, 5/48, 25/34, 55/32. Appended to `do-work/calibration-log.tsv`. The Route C estimate overshot again (55 vs 32), matching last session's pattern of both Route C estimates overshooting by roughly 2–3×. **REQ-245's 5/48 is the outlier and it is not an estimator failure** — it was estimated as trivial correctly, then had its scope widened twice by the orchestrator on the builder's own findings. A recalibration pass should exclude it or record why.

## Context Summary (heavy sessions only)

**Read these fresh before starting; five REQs of carried-over assumptions are not reliable.**

- `_dev/primes/prime-kanban-board.md` — entry point for anything under `skills/do-work-board/tools/queue-kanban/`. Gained REQ-241, REQ-242 and REQ-245's lessons plus **one new convention**: a measured face is per-browser, so record the build beside the number and take the larger where two disagree.
- `_dev/primes/prime-shell-commands.md` — gained REQ-243 and REQ-244.
- `_dev/tests/contract-regressions.sh` — now carries the shape-keyed timestamp citation check, which prints `54 instant write sites cited, 17 date-only sites recognized` on every verify. A change in that count is signal.
- `_dev/tests/shipped-package-reference-contract.sh` — now resolves markdown `#anchor` targets as well as paths, for 27 anchors across four packages.

**Decisions with reach beyond their own REQ:**

- **REQ-241 D-07: this is now the pattern for any constant modelling a rendered face.** State the supremum over the space, not the worst case of a sample; say what makes the sweep complete, next to the number; pin from both sides. A one-sided pin cannot distinguish "correct" from "equal to the last measurement" — which is exactly how 6.75 passed against a 6.71 reference.
- **REQ-242: an exact structural check beats a bigger sweep when the claim is that an input is unused.** Its assertion reads the shipped function out of the generated page and requires the baseline expression to mention neither `dayCentreX` nor `medianMinutes` — that is what makes one measurement at one x a statement about every x. A sweep is still a sample; a band narrower than its spacing slips through.
- **REQ-244 D-08: recognition broad, requirement narrow.** The detector recognizes any placeholder shaped like a value assigned to something that names a clock value or sits under an `*_at` key; only two spellings satisfy it. That split is what makes a checker both complete and enforcing, and it is why names and date-only sites are excluded **by shape rather than by exception** — no path list, no exempt-file list.
- **REQ-245 D-04/D-07: the strict JavaScript behavior lane selects by name pattern, not by file or registry.** A `TestJavaScriptBehavior…` probe participates from any file in the package, so a "you may not touch `generate_test.go`" constraint never blocks JS behavior coverage.
- **REQ-245: a message rendered in two languages needs the pairing made mechanical.** The Go constant and the JS constant are held together by a verbatim comparison that runs without Node; the rendered text is checked by a probe that drives the real code path. Asserting a phrase is *absent* is not a guard — it passed when the whole string was replaced.

**Architectural note:** the board's Durations view now has three separately-pinned geometric guarantees — same-row label separation, label/mark clearance, and the annotation's baseline independence — and they are pinned by three different *kinds* of test (measured-constant assertion, live-DOM intersection, structural source inspection). Adding a fourth guarantee means choosing which kind fits, not copying whichever was written last.
