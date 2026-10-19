---
id: REQ-701
title: 'A failed second-opinion model check keeps the finished page and reports it; a subagent model never makes the backend model ambiguous'
status: pending
created_at: 2026-10-19T14:43:26Z
user_request: UR-158
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-releases.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-700]
batch: UR-158
---
# A Failed Second-Opinion Model Check Keeps The Finished Page And Reports It; A Subagent Model Never Makes The Backend Model Ambiguous
## What
After a second-opinion backend finishes writing its page, the launcher's post-run model verification must never delete that page. When the actual model cannot be confirmed, keep the staged page and reply bytes in the bundle under an explicit unverified name, and report a typed finding naming the models seen. For the Claude backend, a subagent's model in the reply's `modelUsage` map is not ambiguity: the backend's own model is the one that was requested, or the one that answered the main turn, and extra entries from subagents are recorded, not treated as a second candidate.
## Why
A 12-minute generation was discarded because the editor used an Opus-based Explore subagent. The user had to recover the page from the transcript and rerun with the Agent tool denied through a temporary local setting. Destroying finished work on a probe verdict is the wrong failure mode; a page with an uncertain model label is still useful and the user can judge it.
## Context
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` line 394: `defer func() { _ = bundleRoot.Remove(stageName) }()` removes the staging page on every early return of `secondOpinionRun`, including a metadata error.
- Same file, `secondOpinionMetadata` for `claude`: `models = envelope.ModelUsage`, then `if len(models) != 1 { return ... "actual model is absent or ambiguous (%d models)" }`. Claude Code's result envelope lists every model that ran, subagents included, so any Explore or Agent dispatch inside the editorial pass trips this.
- The follow-up screenshot (see Assets) shows the agent's own diagnosis: "Found it: the editor spawned an Opus-based Explore subagent, which confused the launcher into thinking two models were ambiguous, causing it to delete the finished page."
## Requirements
- No code path in `secondOpinionRun` deletes a staged page or reply that the backend finished writing. Cleanup of the staging file is limited to the case where the backend never wrote it.
- When model verification fails after a finished page exists, publish the page and reply under a name that says the model is unverified (the builder picks the exact form, consistent with `index-<cli>-<model>[-N].html`), insert no banner, and emit a typed finding listing the requested model and every model seen.
- For `claude`, the actual model is the requested `--model` when it appears in `modelUsage`; otherwise the single entry when there is exactly one. Extra subagent entries are recorded in the run row, not treated as ambiguity. Only a map that is empty, or that lacks the requested model and has several entries, is a verification failure.
- Lock both behaviors in `report_second_opinion_test.go`: a fake backend that writes the page and returns a two-model envelope ends with the page still on disk and a finding; a fake that returns requested model plus one subagent model ends with the normal verified publication.
## Red-Green Proof
**RED prompt/case:** In `report_second_opinion_test.go`, a fake `claude` executable writes the staging page and prints a success envelope whose `modelUsage` has two entries, `claude-fable-5-1` (requested) and `claude-opus-5-5`. Assert the page exists after the run and the run row reports the requested model.
**Why RED now:** `len(models) != 1` returns an error, the deferred `Remove(stageName)` deletes the page, and the run row says "actual model is absent or ambiguous (2 models)".
**GREEN when:** The page is published under the requested model's name with the normal judge and banner flow, the subagent model is recorded in the row, and a separate test with an envelope lacking the requested model and holding two unrelated models ends with the page kept under an unverified name, no banner, and a typed finding.
**Validation:** Inferred during capture from the screenshot and the code at the cited lines.
## Assets
- `do-work/user-requests/UR-158/assets/REQ-701-screenshot-2-second-opinion-probe-deleted-page.png` — Claude Code transcript continuing the REQ-700 run: "Background command 'Run the second-opinion generation with the claude backend' completed (exit code 0)", "The launcher failed because its model probe found two models from the claude CLI and treated that as ambiguous.", "Still diagnosing: the generation itself ran for 12 minutes and then was discarded.", "Found it: the editor spawned an Opus-based Explore subagent, which confused the launcher into thinking two models were ambiguous, causing it to delete the finished page. I'll recover the page from the transcript as a fallback and rerun with the Agent tool blocked for the child process.", "Recovered page saved to the scratchpad as a fallback. Now adding a temporary Agent deny to the gitignored local settings so the child run cannot spawn a second-model subagent, then rerunning.", then the rerun command with `--backend claude --model claude-fable-5-1`.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 22085 tokens, over the 2000 budget; `slugged: partial` so no targeted entry is legal. Matched on native backend identity and rollback/deletion identity in do-work-cli internals (families native-backend-identity, final-boundary-identity, silent-skip-reads-as-red).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: [Image #3] <- see additional feedback*
