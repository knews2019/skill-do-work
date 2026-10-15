# REQ-698 builder handback

Branch: `worktree-agent-REQ-698-ai-report-second-opinion`
Worktree: `/private/tmp/worktree-agent-REQ-698-ai-report-second-opinion`
Final HEAD: `2abd7939ac52c9813616909e388cfedfb8b78b6a` (clean worktree; includes the review remediation addendum).
Base: `61ca1e9c`.
Commits, in order:
- `406ff3c4fda1d083dc144c01929b453360583ae5` — implementation, prompt/docs, focused contracts.
- `9fa76500a7177cf00b77c36defb0637a9f252283` — source/reply link portability.
- `2f5cd87a` — equivalent opening backlink validation and two-word typed Go result fields (JSON keys unchanged).
- `69a2a998cfbf21391fd2a9c7c3cff4c7718da6d3` — live artifacts only.

Implementation and scoped verification are complete. Parent owns review, full maintainer gate and release/lifecycle. After explicit user approval, exactly one further AGY proof completed and correctly refused a candidate with navigation before the required body-opening backlink; see the final addendum. No AGY artifact was published. Do not represent AGY or unavailable Gemini as successful live acceptance.

## P-A-U
- [x] **[PLAN]:** Read durable brief, scoped plan/exploration, CLAUDE.md, full five primes; general, coding-guardrails, shared-principles, communication-style, backend, testing, security, anti-slop, prompt-injection and background-agent crews. Implement explicit native CLI adapters, protected bundle publication, shared judge, typed evidence and install-safe prompt discovery. RED cases target the real handler registry seam.
- [x] **[APPLY]:** Changes remain inside the declared source/test/docs and live bundle scope. Parent owns lifecycle, release and main REQ. No consumer edits.
- [x] **[UNIFY]:** Reviewed every changed source/docs/test file and `git diff --stat`/`git diff --check`; gofmt clean, focused race tests, vet, Windows compile, source/install launcher and references pass. Full maintainer suite is parent-owned. Live proof follows below.

## Decisions
- D-05 — DECIDE & STATE: Use one short native read-only model-resolution probe (90-second bound) before generation, then require identical native model metadata on generation. AGY uses its zero-turn `/model` command. Codex native text startup header, Claude modelUsage and Gemini stats.models provide actual identity; absent/ambiguous/mismatched identity refuses. A default CLI model is never represented as this chat's active model.
- D-06 — DECIDE & STATE: Models fill an unpredictable launcher-reserved staging leaf. The launcher publishes the verified model-named final page/reply exclusively and maps that exact staging leaf to the final page name inside the saved reply. New assets remain directly in the bundle under the verified backend/model prefix. Captures use disposable scratch outside the repository and typed findings survive cleanup.
- D-07 — DECIDE & STATE: Reject symlinks/special files in report inputs and new outputs; snapshot every pre-existing bundle file's inode, complete mode and bytes. Do not auto-revert a backend violation. Preserve strict pre/post Git status, including both rename paths, and report named violations.

- D-08 — DECIDE & STATE: Save relative published page links in replies by replacing the exact absolute reserved stage path before its leaf name. Prompt source-map references use repository-relative text unless a canonical remote link is independently known to be published; a local commit and configured remote are insufficient. Applied only to new generated artifacts, with disclosure and a separate post-edit render judge.
- D-09 — DECIDE & STATE: Backlink validation accepts reordered anchor attributes, unquoted hrefs and harmless inline formatting, but rejects duplicate/wrong hrefs, misleading href text inside other attributes, wrong labels and backlinks following main content. A focused regression proved the old parser too strict; the original AGY candidate was not retained, so its failure cannot be confidently attributed to model noncompliance.

## Required lesson evidence
Read `_dev/primes/lessons-releases.md` whole-satellite. No missing required entry.
Touch-conditional reads: `_dev/primes/lessons-action-files.md` Template, Agent Compatibility, alternate-writer-contract-drift/restated-mechanism families; `_dev/primes/lessons-shell-commands.md` prescribed commands, unchecked exit, SIGPIPE, every-flag caller sections; `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` Package routing/direction, Scoped command contracts, Verify, final-boundary-identity, opaque-evidence-projection, reaped-by-its-own-parent, smoke-vs-characterization, shipped-module-test-self-containment families. Large initial combined outputs were truncated; the relevant sections/families were subsequently read in bounded reads. Also read `_dev/lessons/validated-runtime-boundaries.md` fully. Do not claim entire large satellites were read.

## Tests
- RED: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run SecondOpinion -count=1`: all initial cases failed at `ai-report-second-opinion command is absent` (assertion failure, not compile failure), `/private/tmp/req698-evidence/builder-red.log`.
- GREEN implementation focused PASS7.084s (`builder-publish-test.log`), reply portability focused PASS7.595s (`builder-reply-portability-test.log`).
- Backlink regression RED: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run TestSecondOpinionBacklinkAllowsSemantic -count=1` failed three valid-equivalent examples; `builder-backlink-red.log`. GREEN: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands ./internal/resultmodel -run SecondOpinion -count=1` PASS7.375s/0.363s (`builder-backlink-green.log`).
- Final race package checks: `go -C skills/do-work/tools/do-work-cli test -race ./internal/toolboxcommands ./internal/resultmodel`, PASS toolbox11.485s/resultmodel2.010s (`builder-final-recheck.log`), after backlink and naming edits. Prior race checks also passed.
- `go -C skills/do-work/tools/do-work-cli vet ./...`: final exit0 (`builder-final-vet.log`).
- `GOOS=windows GOARCH=amd64 go -C skills/do-work/tools/do-work-cli test -c ./internal/toolboxcommands -o /private/tmp/req698-evidence/toolbox-windows.test.exe`: exit0.
- `bash _dev/tests/do-work-cli-launcher-behavior.sh`: exit0, source and installed prompt topology fixtures pass.
- `bash _dev/tests/shipped-package-reference-contract.sh`: exit0, PASS.
- `bash _dev/tests/contract-regressions.sh`: first restricted-sandbox attempt exit1 with unavailable cached executables and `ps: Operation not permitted`; exact retry with escalation exit0, `Contract regression checks passed`. No cache implementation changes.
- Fake executable coverage: all4native adapters; terminal chooses once; no backend/model substitution; generation model switch; filename/reply collisions; latest/direct-child/symlink refusal; pre/post outside dirt and rename source; asset/original preservation; required backlink; pass/fail/skipped/error banner gating; serial batch missing/excluded/failed states; typed JSON findings; cancellation; actual-body banner and per-filename idempotency; required headless flags.

## Integration seams
None deferred. Scoped handler registration/count, typed result normalizer, launcher template path, routing and source/install contract test are included. Parent should run canonical maintainer gate and release/archive workflow against the merge range.

## Discovered Tasks
None outside this request.

## Live proof
All commands use `skills/do-work/tools/do-work-cli.sh --repo-root /private/tmp/worktree-agent-REQ-698-ai-report-second-opinion --format json ai-report-second-opinion` and bundle `ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism`.

| Mode / arguments | Native identity | Outcome | Judge | Banner | Outside paths |
| --- | --- | --- | --- | --- | --- |
| Targeted `--backend codex` | `gpt-6-astra` | Exit0; page/reply published | pass, findings[] | inserted | [] |
| Pinned `--backend claude --model claude-opus-5-5` | `claude-opus-5-5` | Exit0; page/reply published | pass, findings[] | inserted | [] |
| Batch `--all --exclude codex`: codex | none (excluded) | excluded calling backend | not-run | not-inserted | [] |
| Same batch: gemini | none | failed native result, exit55 | not-run | not-inserted | [] |
| Same batch: agy | `gemini-3.8-flash-low` | backlink validation failure; no published output | not-run | not-inserted | [] |

Evidence files are `/private/tmp/req698-evidence/live-codex.json`, `live-claude.json`, `live-batch.json` (the batch exits1 with findings). Gemini exact reason: `backend failed: exit status 55; gemini returned no successful native result; inspect backend authentication/availability`. Earlier standalone native probe stderr records `UNSUPPORTED_CLIENT`; no authentication or account settings changed. No fallback backend was substituted.

AGY exact reason: `new page must open its body with Original report: index.html linking to index.html`. Native actual identity passed before the refusal. Its original generation stream lived only in the capture buffer, and its reserved stage was removed on return; retained evidence cannot distinguish omitted/misplaced backlink from the now-fixed parser false rejection. The independent parent native success fixture remains `/private/tmp/req698-evidence/agy-probe.jsonl` and confirms the installed init/result schema. Do not call the original batch refusal proven model noncompliance.

Parent requested one same-model targeted AGY proof after fixing the parser. Automatic approval review rejected that command twice before any process started. First reason: transmitting the bundle to AGY lacked explicit trusted-user authorization. The retry supplied UR-156 lines80/82 and the REQ exact-bundle requirement; the reviewer again rejected it, treating the request file as untrusted and the user skill invocation as insufficient egress authorization. No diagnostic wrapper was created and no new AGY call occurred. Parent subsequently obtained explicit user approval for exactly one run; its result is recorded in the final addendum. The two rejected attempts started no model process, and no bypass occurred.

This active chat's exact model identity is unavailable, and it is a Codex session rather than Claude Code. Therefore the Claude proof is an explicit pinned-backend equivalent, not a claim of running from an actual Claude session. Batch truthfully excluded Codex. Actual-Claude-session and complete all-installed-backend live success remain unproven here; fake executable tests cover the scoped mode contract.

### Artifact provenance and inspection
- Codex initially produced `index-codex-gpt-6-astra.html` and `codex-gpt-6-astra-reply.md`; native launcher judge passed. Parent authorized a verifier edit removing 22 source hyperlinks to unpublished local commit406ff3c4 while preserving their path/code text, disclosing the source-map change, making the exact worktree output link relative, and separating the model's unavailable preview note from launcher verification. Original model bytes remain `/private/tmp/req698-evidence/codex-model-original.html` and `codex-model-original-reply.md`. Do not attribute amended bytes to the model.
- The amended Codex page independently passed `ai-report-judge` in a fresh scratch bundle, wide/phone and light/dark, findings[]; `/private/tmp/req698-evidence/codex-postedit-judge/judge.json` plus captures and `codex-postedit-judge-result.json`. The temporary scratch bundle was removed. Two historical remote source links are retained; their commits were verified present in `origin/main`.
- Claude output needed no verifier content changes; the launcher normalized the assigned stage path to its published relative link. Its reply's statement that the model could not run browser tools is model-only; the launcher independently returned judge pass. Page uses relative evidence links/plain repository source paths.
- Parent observed Codex in a desktop browser: readable layout, original backlink and navigation. Toggling “Include the matching commit at minute70” off changed the largest gap60→120 minutes (dispatch→builder handback).
- Parent observed Claude in a desktop browser: legible title/content, original backlink. Clicking “Report option1” set aria-selected=true and switched to the stamp-then-delete explanation. These are parent-observed interactions, not builder claims of independent browser interaction.
- Final byte comparison against61ca1e9c proves original `index.html` differs only by exactly two marker-wrapped banners; all three pre-existing evidence files are byte-identical. Reply file links resolve within the bundle and contain no worktree/staging path. Evidence: `/private/tmp/req698-evidence/artifact-invariants.json` (also includes final artifact SHA256s).
- Fresh invocations intentionally use numeric suffixes and get distinct banners. Per-filename banner idempotency is covered by focused tests; no second expensive fresh run was misrepresented as “already present.” Parent D-01 resolved that request contradiction.

## Full file manifest and review
All18 changed files were reviewed; no unrelated or lifecycle files appear. Docs were checked for official trigger/dispatcher, full prompt preservation, canonical source-link wording and relative sibling topology. CLI files were checked for native identity, process ownership, guarded publication, typed JSON, result normalization and launcher template lookup. Tests cover the scoped behavior; live HTML/replies and original-only banners have the verification above.

```text
M	_dev/tests/do-work-cli-launcher-behavior.sh
A	ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/claude-claude-opus-5-5-reply.md
A	ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/codex-gpt-6-astra-reply.md
A	ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-claude-claude-opus-5-5.html
A	ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index-codex-gpt-6-astra.html
M	ai-reports/2026-10-09_1916_REQ-632-activity-correlation-determinism/index.html
M	skills/do-work-toolbox/SKILL.md
M	skills/do-work-toolbox/actions/ai-report-reference.md
M	skills/do-work-toolbox/actions/ai-report.md
M	skills/do-work-toolbox/docs/ai-report-guide.md
A	skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md
M	skills/do-work/tools/do-work-cli.sh
M	skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go
M	skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go
M	skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go
M	skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go
A	skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go
A	skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go
```

## Parent next steps
The original implementation/live range has been handed back. Run the canonical maintainer gate and formal review, then release/lifecycle if appropriate. Preserve Gemini unavailability, the actual-session-context limitation and the final AGY refusal in the user report. No additional AGY attempt remains authorized or pending from this builder.


## Final addendum — one explicitly approved Antigravity verification

Trusted authorization relayed by parent: user answered **“Allow this one Antigravity verification run”** to **“May I rerun Antigravity (model gemini-3.8-flash-low) once to verify the backlink fix, sending the tracked REQ-632 report bundle and relevant repository source to that model?”**

Executed exactly one invocation at clean HEAD `69a2a998cfbf21391fd2a9c7c3cff4c7718da6d3`, with the same wrapper command/bundle above and `--backend agy --model gemini-3.8-flash-low`. A transparent scratch PATH wrapper called the installed `/Users/t2/.local/bin/agy`, preserving stdout, stderr and the reserved candidate before launcher cleanup. It changed no tracked source. Native probe and generation both resolved the explicitly requested model; generation emitted terminal `SUCCESS` for conversation `39bcf718-8719-4f38-8788-9543c0320b35`.

**Launcher result: exit2 / failure; state failed, actual_model gemini-3.8-flash-low, judge not-run, banner not-inserted, page/reply empty, outside_paths[].** Exact reason remains `new page must open its body with Original report: index.html linking to index.html`.

This candidate is now diagnosable: immediately after `<body>` it opens a container and sidebar headed `Report Navigation`, followed by a full table of contents. Only after that sidebar does `<main class="main-content">` begin with the correct `<p>Original report: <a href="index.html">index.html</a></p>`. The hyperlink is correct, but preceding navigation violates the prompt's required body-opening backlink. This is a real candidate ordering violation for the approved run; it does not retrospectively establish the cause of the earlier unretained batch candidate.

Evidence retained:
- `/private/tmp/req698-evidence/live-agy-approved.json` — typed launcher result.
- `/private/tmp/req698-evidence/agy-approved-probe.stdout` and `.stderr` — native zero-turn model resolution.
- `/private/tmp/req698-evidence/agy-approved-generation.stdout` and `.stderr` — full native generation stream (60,472 bytes stdout, empty stderr).
- `/private/tmp/req698-evidence/agy-approved-candidate.html` — exact model candidate, never published or amended.
- `/private/tmp/req698-evidence/agy-approved-before.json`, `agy-approved-original-before.html`, `agy-approved-invariants.json` — byte comparison and opening markup.

All eight pre-existing bundle files are byte-identical before/after the run, including the original with its two existing banners. Git status remains completely clean; no candidate, reply, asset or new banner was published. No new commit was needed; final HEAD remains69a2a998. No source changes and no second model retry were made. This completes the explicitly authorized single verification attempt.

## Review remediation addendum — F1 process ownership

Remediation commit: `2abd7939ac52c9813616909e388cfedfb8b78b6a` on the retained builder branch, based on `69a2a998`. Exactly two files changed: `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` and `report_second_opinion_test.go`. Parent owns the re-review, merged gate, lesson/prime update, lifecycle and release. No external model calls occurred in this remediation.

Read the complete `REQ-698-review.md` before production patching; reread `skills/do-work/crew-members/debugging.md`, the whole CLI prime and the exact `reaped-by-its-own-parent` satellite entry, including its `ESRCH` trap. Parent reserved D-12 for its lesson/prime update.

**D-11 — DECIDE & STATE:** Keep a small POSIX-shell supervisor as the isolated group owner until the existing `ownedprocess.TerminateGroup` cleanup completes. Backend executable and arguments remain positional (`"$@"`), never interpolated/evaluated as shell code. The supervisor waits for the backend, writes its exact exit status to a private inherited pipe, closes that pipe and waits on parent-owned stdin. The backend closes fd3 and receives `/dev/null` stdin, so it cannot inherit either completion/control channel. Parent status parsing is bounded to32bytes and rejects malformed/noncanonical statuses. Cleanup still proves ownership through the shared helper, runs before `Cmd.Wait`, retains cleanup errors, and then reaps the supervisor/drains captured output. No unchecked numeric process-group signal or shared ownedprocess change was introduced. Windows remains fail-closed before process start.

### Structured diagnosis and evidence
- Observation: independent review probe returned success, then the redirected child wrote `late-write`300ms later.
- Hypothesis: `Cmd.Wait` reaped the isolated backend leader before `TerminateGroup` could prove its PGID; its ESRCH fallback missed live descendants. Inherited stdout/stderr instead kept `Cmd.Wait` blocked while children could still write.
- Prediction: a redirected child will write after successful capture; an inherited-output child will delay capture until it closes the pipes.
- RED confirmed both: `remediation-red.log` reports `capture returned while a descendant could still write` in both variants and inherited-pipe elapsed919.212041ms. Test command: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run 'TestSecondOpinionCapture(DrainsExitedLeaderChildren|PreservesBackendStatusAndChannels)$' -count=1`, exit1.
- Root cause confirmed; one production approach was implemented. Initial strengthened cancellation fixture asserted no write even during cleanup and failed. Measurement showed its `late` file already existed at capture return (killing its waited-for sleep resumes the shell's next statement during teardown). The corrected assertion compares the snapshot at capture return with the later state; it verifies the promised post-return boundary. Parent explicitly agreed this distinction. The successful-leader cases continue to require that the delayed write never occurs at all.
- Focused GREEN: `go -C skills/do-work/tools/do-work-cli test ./internal/toolboxcommands -run 'TestSecondOpinionCapture' -count=1 -v`, pass3.137s, `remediation-green.log`. Covers redirected/inherited-pipe children, cancellation with no post-return write, exact stdout/stderr and literal argv, truthful status0/7, closed backend stdin and closed status fd3.
- Reviewer's unchanged original probe GREEN: `go -C skills/do-work/tools/do-work-cli test -overlay /private/tmp/req698-evidence/remediation-review-overlay.json ./internal/toolboxcommands -run '^TestReviewSecondOpinionLeaderExitLeavesWriter$' -count=1`, pass2.047s, `remediation-review-green.log`. Overlay only points the retained reviewer test into this worktree.
- Final race packages: `go -C skills/do-work/tools/do-work-cli test -race ./internal/toolboxcommands ./internal/resultmodel`, pass toolbox15.804s/resultmodel cached, `remediation-race.log`. This includes native adapter/metadata contracts and all existing scoped command tests after the runner change.
- `go -C skills/do-work/tools/do-work-cli vet ./...`: exit0, `remediation-vet.log`.
- `GOOS=windows GOARCH=amd64 go -C skills/do-work/tools/do-work-cli test -c ./internal/toolboxcommands -o /private/tmp/req698-evidence/remediation-toolbox-windows.test.exe`: exit0, `remediation-windows.log`.
- gofmt and `git diff --check` clean; reviewed both changed files and119insertions/14deletions. Worktree clean after commit. Live artifact bytes untouched.

### Lifetime limits and lesson handoff
The supervisor is our child and is reaped by our `Cmd.Wait`. Backend descendants are terminated through the existing bottom-up shared group boundary. Descendants orphaned by a backend that exits itself can be proved stopped, but their eventual reap belongs to init; this patch does not claim to reap another process's child. A deliberate `setsid` escape remains outside group ownership and is not signalled, as the governing lesson already states. The reviewed in-group late-writer and inherited-pipe defects are covered by the new behavioral regressions.

Proposed parent lesson wording: **Symptom:** success returned before a redirected child wrote, while inherited-output children delayed completion. **Root cause:** backend completion was conflated with ownership-lifetime completion, and reaping the sole group leader erased PGID proof. **Fix:** a supervisor retains isolated ownership, reports backend status separately and stays alive until proven bottom-up cleanup; only then wait/reap and publish. **Discovery:** reviewer minimal delayed-write probe plus differential real-child RED/GREEN tests. Parent owns updating the existing generalized prime trap/satellite/index rather than adding duplicate lesson machinery.
