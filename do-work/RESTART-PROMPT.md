```
do-work run --fan-out 2

This command is sufficient; everything below it is context.

You are picking up a do-work queue in the skill-do-work2 checkout. Nine requests
are pending, six of them ready now. Nothing is claimed, nothing is half-merged,
`do-work/working/` holds no request, and the canonical gate is green.

Three things to know while you build, none of which change the command above:

1. Other sessions write to this same checkout. Before every merge, re-read the
   tip; never revert, stash or commit another session's bytes. Run the canonical
   gate and the heavy lanes from a detached worktree whenever the shared tree is
   dirty, and seed `do-work/test-durations.tsv` with its header line there or
   the staged-skills lane fails instantly with an invalid-header error.

2. The browser heavy lane reports SKIPPED, not failed, when it finds no browser,
   and a skipped lane is not a pass. Export
   QUEUE_KANBAN_BROWSER="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
   before running heavy lanes, then check each lane's wall_seconds — a lane that
   took 0s did not run. Never pipe the heavy runner through `tail`; the per-lane
   JSON is the only record of what actually executed.

3. Closing a user request moves its members into `do-work/archive/UR-NNN/`, which
   silently breaks any shipped file linking one of those records by its old flat
   path, and the canonical gate then goes red on a request that has nothing to do
   with it. That fired twice on 2026-09-05. Before each release, run:

   for f in $(grep -rl 'blob/main/do-work/archive/' skills/); do grep -o 'blob/main/do-work/archive/[^)#]*' "$f" | sed 's#blob/main/##' | sort -u | while read p; do [ -e "$p" ] || echo "MISSING in $f: $p"; done; done

Leave the worktree `worktree-agent-REQ-573-activity-drawer` alone. It is an
abandoned remediation, deliberately unmerged, and must not be force-deleted.
```

---

# Reference

Written 2026-09-05 at the close of run `work-2026-09-05-120117`. This section is for humans and debugging; the paste block above stands alone.

## State at handoff

- **Integration branch:** `main`, tip `64b78ea8`.
- **Version:** 0.303.1. Canonical gate `bash _dev/tests/maintainer-verify.sh` green at `87097144`, recorded as a reusable green-gate record.
- **`do-work/working/`:** empty of requests. `do-work-cli recover` returns `finalization_passed: true`, no claims, findings `FINALIZATION-NONE` / `RECOVERY-NONE`.
- **Queue:** 9 pending, 6 ready. No request is `pending-answers` or `blocked`, so `do-work clarify` has nothing to do.

## What this run finished

Seven requests, each merged, qualified, independently reviewed, heavy-lane verified, archived and released:

| Request | Release | Merge range |
|---|---|---|
| REQ-572 — every lifecycle transition as its own Activity row | 0.296.0 | `7ad53bff..fbdcd35e` |
| REQ-573 — click an Activity row to open the drawer | 0.300.0 | `45a9010d..2d3981f4` |
| REQ-575 — timestamps are append-only across transitions | 0.295.0 | `cd686ed7..a2c6f4cf` |
| REQ-578 — the Verify Findings strip steps aside on Activity | 0.297.0 | `7dbb2756..09aaa9a4` |
| REQ-579 — verify findings as compact rows | 0.301.0 | `4362ac0d..b169396e` |
| REQ-581 — a cleanup test that can actually fail | 0.302.0 | `1a06c3bc..92339213` |
| REQ-582 — a citation's section name has to resolve | 0.303.0 | `43be2af6..7b2673b6` |

Two standalone maintainer commits outside the pipeline, both unblocking a red gate caused by stale archive links: `09a13839` (0.294.1) and `87097144` (0.303.1). The second was caused by this run's own UR-115 closure.

Run artifacts: `do-work/runs/work-2026-09-05-120117/` — manifest, seven briefs, seven hand-backs, all committed.

## Worktree verdicts

The survey found two worktrees.

- `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` — the main checkout, branch `main`. Clean.
- `.git/work-run-20260905-1201/worktree-agent-REQ-573-activity-drawer`, head `9de57c92`, **clean, branch unmerged**. This is mine and it is **deliberately abandoned**, not removable. It fails the first REMOVABLE condition: the branch is not in the survey's merged list. It holds a remediation for REQ-573's F1 finding that was superseded mid-flight when another session fixed the same defect on main at `7443fe11` by a better route (the drawer's own setter became the single writer and both document click listeners were deleted). The maintainer chose main's fix. `git branch -d` will refuse, and forcing it would destroy the only evidence that the integration did not happen. If you want it gone, that is a maintainer decision to make explicitly:

  ```
  git worktree remove .git/work-run-20260905-1201/worktree-agent-REQ-573-activity-drawer
  git branch -D worktree-agent-REQ-573-activity-drawer   # -D only, because it is unmerged on purpose
  ```

No foreign worktree was touched and none was removed.

## Uncommitted files at handoff

None. Two files were untracked when the survey ran and are now committed:

- `do-work/runs/work-2026-09-05-124800/REQ-585-handback.md`
- `do-work/runs/work-2026-09-05-124800/REQ-586-handback.md`

They belonged to another session's run, which had already finished — both records archived under UR-120 and UR-121, both released as 0.298.0 and 0.299.0, that run's manifest closed at `d460759f`, and the rest of its run directory already tracked. The session that owned them had ended. While they sat untracked, `do-work-cli recover` refused with a typed `FINALIZATION-DISCOVERY-AMBIGUOUS` naming exactly those two paths, and because `do-work run` runs `recover` first and treats a typed refusal as a whole-run stop, **every session in this checkout was blocked behind two orphaned files.** Committed unmodified on their own as `64b78ea8`, which is the least destructive clearing action `actions/work.md` prescribes for unrelated work with no live owner.

## Parallelism analysis

`--fan-out 2`, and start with the two `priority: next` requests.

**Safe to run concurrently.** REQ-583 (pinning the evidence-gate remedy redirection) declares exactly one file, `internal/lifecycleadvance/evidence_gates_test.go`, which nothing else in the queue names. It is disjoint from every other pending request and pairs safely with anything.

**Expect a conflict, but do not gate on it.** Six requests — REQ-552, REQ-554, REQ-555, REQ-556, REQ-557, REQ-558 — every one of them declares `_dev/tests/audit-lockins.sh`. That file takes an appended lock-in per request, so two builders touching it produce an ordinary append/append conflict at integration, cheap to resolve by keeping both. This is a "you will merge by hand", not a "these must never run together", so it is deliberately **not** encoded as a dependency gate: inventing gates the maintainer did not write would distort the queue's real shape to describe a two-line merge. The merge is the non-interference proof, which is the design.

**Undeclared write sets — treat as unknown, not safe.** REQ-486 (collapsible UR groups with progress summaries) and REQ-587 (giving the Timeline view one scroll surface) declare no `write_set`. Both are queue-kanban board client work and will land in some overlap of `web/template.html`, `web/board.css` and the `web/board-*.js` fragments. Three separate sessions collided in exactly those files on 2026-09-05. An absent declaration reads as *unknown*, never as safe, so do not run these two as a pair; put one of them beside REQ-583 instead.

**Already serialised by dependency gates, no action needed:** REQ-552 → REQ-557 → REQ-558, and REQ-554 → REQ-555. Those chains are the queue's critical path for the audit batch; the selector already honours them.

**Critical path.** REQ-552 → REQ-557 → REQ-558 is three deep and REQ-557 alone declares nine Go files, so start that chain early rather than leaving it for last. REQ-583 and one board request fill the other builder while it runs.

## Findings left open, and what it would take to promote them

Both are recorded in their archived requests' `## Review` sections. Neither auto-queued, because the rule reserves automatic follow-ups for `impact-critical` work, and this project's standing preference is to keep noncritical findings in the report rather than mint requests for them. Promote either with one `do-work capture-request:` if you want it built.

- **Append-only timestamps kept three exceptions** (`do-work/archive/UR-116/REQ-575-*.md`). A verification stamp and a blocked stamp are withdrawn together with the state they describe, and one terminal stamp is re-written on the cancel path. The reviewer judged that faithful to the intent, and separately measured that **only the claim stamp is structurally write-once** — the other lifecycle stamps are written by prose in `actions/work.md` that carries no "only when absent" condition, so a recovered request that is re-planned still overwrites its planning stamp. That is the same evidence loss the request set out to stop.
- **The citation checker's bold-label rule is wider than intended** (`do-work/archive/REQ-582-*.md`). Measured: one reference file declares 42 headings but yields 288 accepted names, 246 bold-only; renaming three real headings leaves the check green because each old name survives in bold prose. The reviewer verified a tighter rule — accept only a *paragraph-leading* bold run — which keeps all 74 live citations green and makes all three simulated renames report. Two further findings ride along: nothing pins the wiring between the check's two halves (mutating the driver call leaves the suite green, the same "control that cannot fail" class REQ-581 was raised for), and two live citations using the ASCII `->` arrow are invisible to the check entirely.

## Heads-up list — what will bite in the first ten minutes

- **The heavy browser lane skips silently.** It exits 0 with `wall_seconds: 0` and a `HEAVY-RUN-LANE-SKIPPED` finding when no browser is on PATH; Chrome is installed here but not under a name the probe looks for. Export `QUEUE_KANBAN_BROWSER` as the paste block shows. *Next session.*
- **Never pipe the heavy runner through `tail`.** The per-lane JSON is the only record of which lanes ran; truncating it costs a full re-run. *Next session.*
- **Detached-worktree lane runs need a seeded durations header.** `printf 'run_id\tfile\tseconds\tother_gate_processes\n' > do-work/test-durations.tsv` in the detached worktree, or `staged-skills` fails in 0s on an invalid header that looks like a real failure. *Next session.*
- **Closing a user request breaks archive links in shipped files.** Run the sweep in the paste block before every release. Both occurrences and the rule are recorded in `_dev/primes/lessons-releases.md` under `[family: canonical-link-outlives-its-target]`. *Next session.*
- **Other sessions are live in this checkout.** On 2026-09-05 three sessions wrote here at once; two independently fixed the same review finding twenty minutes apart, because a report-only finding is not a claim and nothing makes a second reader visible to the first. If you pick up a finding from a report, expect that someone else may already be on it. *Maintainer, if the intake brake is ever tuned.*
- **The finalization manifest has no dry-run.** `advance --finalization-manifest` iterates on typed refusals; the release manifest must be wrapped as `{"operation":"release", "release":{...}}` with `maintainer_release: true`, and `commit_paths` must be a superset of everything the planner plans — closing a user request pulls in every sibling record and asset on both the archive and user-requests sides. Submit, read the refusal, add the named paths. *Next session.*
