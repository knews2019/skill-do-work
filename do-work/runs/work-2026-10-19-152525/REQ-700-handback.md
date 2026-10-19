# REQ-700 hand-back: second opinion from the session model runs in-session, not as a hidden shell

- Branch: `worktree-agent-REQ-700-second-opinion-session-prepare-finalize`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-700-second-opinion-session-prepare-finalize`
- Base commit (before the first edit): `5f237933` (main after REQ-701's release 0.306.3)
- Commits (both subjects start with `[REQ-700]`):
  - `87a93660` tests plus Go code: `--prepare` and `--finalize` phases, digest sidecar, shared publication tail, result-model fields
  - `7ab469f3` prose: action § Second-opinion form and the user guide § A Second Opinion
- Working tree after the second commit: clean (`git status --short` empty). The hand-back file is the only main-tree write; nothing under `do-work/` was staged or committed.

## File manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified): `--prepare`/`--finalize` parsing and validation, `secondOpinionUsage` names both phases, `secondOpinionStagePattern`, `secondOpinionRecord`/`secondOpinionDigest`/`secondOpinionDigests` (the sidecar), `secondOpinionFile` gains `mode` and `digest` filled by `secondOpinionSnapshot`, `secondOpinionPreserved(bundle, before, allowed []string, prefix)` with the nil-`info` digest branch, `secondOpinionOutsideCheck` lifted from the closure, `secondOpinionPrepare`, `secondOpinionFinalize`, `secondOpinionPublication` extracted from `secondOpinionRun` (carries REQ-701's unverified publication and removes the stage leaves after both publishes), `row.ModelEvidence = "native"` on the exec path. `secondOpinionArgv`, `secondOpinionMetadata`, `secondOpinionCapture` and the supervisor are untouched.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified): `secondOpinionFakeLogger`, `TestSecondOpinionSessionPrepareNeverExecsOwnBackend`, `TestSecondOpinionSessionFinalizePublishesAgentWrittenPage` with the four subtests.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified): `StagingPage`, `StagingReply`, `AssetPrefix`, `PromptText`, `ModelEvidence` on `AIReportSecondOpinionRun` after `OutsidePaths`, all `omitempty`; normalizer unchanged.
- `skills/do-work-toolbox/actions/ai-report.md` (modified): Session bullet rewritten to the prepare, do-the-prompt, finalize sequence with the five row fields; Targeted bullet routes a target equal to the agent's own CLI through the same sequence; native-metadata sentence qualified to launcher-run backends. Heading text `## Second-opinion form` unchanged; no new shell fence; no new citation.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): three plain-words sentences in the second paragraph of § A Second Opinion (in-session pass, what the CLI prepares and finalizes, the declared model label); the native-metadata sentence qualified to launcher-run backends. REQ-701's `-unverified` sentences in both files are kept.

## P-A-U

- [x] **[PLAN]:** Read the REQ in full (What through Scope), the brief with its Coordinator note, UR-158's input, the archived REQ-701, the crew rules (general, coding-guardrails § 5, shared-principles, communication-style, testing, anti-slop), the three primes (action-files § Cross-Referencing and the alternate-writer-contract-drift trap; shell-commands § Closed Enumerations Go Stale; releases), and the `native-backend-identity` and `reaped-by-its-own-parent` bullets of `lessons-do-work-cli.md`. Followed the Plan's five tasks in order with the decided mechanism; no change to the mechanism.
- [x] **[APPLY]:** Task 1 RED prepare test; task 2 result model, parser, handler branch, `secondOpinionOutsideCheck`, `secondOpinionPrepare` (GREEN); task 3 refactor under green (digest-aware `secondOpinionPreserved`, `secondOpinionPublication` extraction, pre-flight probe rerun); task 4 RED finalize test, then `--finalize` parsing, handler bundle derivation with the supersession guard, `secondOpinionFinalize` (GREEN); task 5 prose. Scope held to the five write-boundary files.
- [x] **[UNIFY]:** `git diff 5f237933 --stat`:
  ```
   skills/do-work-toolbox/actions/ai-report.md        |   6 +-
   skills/do-work-toolbox/docs/ai-report-guide.md     |   4 +-
   .../internal/resultmodel/result_model.go           |   6 +
   .../toolboxcommands/report_second_opinion.go       | 329 +++++++++++++++++----
   .../toolboxcommands/report_second_opinion_test.go  | 135 +++++++++
   5 files changed, 415 insertions(+), 65 deletions(-)
  ```
  Checks (worktree root or module root, wall times): `git diff --check` clean before each commit; `gofmt -l ./internal` prints nothing (one blank-line fix applied by `gofmt -w` before the Go commit); `go vet ./internal/...` exit 0; `go test -count=1 ./internal/toolboxcommands ./internal/resultmodel` both `ok` (17 s, 0.6 s); REQ-700 probe exit 0 (1 s); pre-flight probe exit 0 (10 s; also 15 s after task 2 and 10 s after task 3); `action-shell-blocks.sh` exit 0 (2 s, 81 fences, ShellCheck enabled); `prescribed-shell-canonicalization.sh` exit 0 (1 s); `shipped-package-reference-contract.sh` exit 0 (1 s); `audit-lockins.sh` exit 0 (1 s). Files checked by reading the full diff: the Go file (every hunk reviewed for the exec path staying unchanged: argv, metadata, capture, supervisor, probe calls), the test file, the result model, both prose files. No debug artifacts; the full `go test` runs were on the two touched packages only, never the repository gate.

## Proof record

1. **RED prepare** (`go test -count=1 -run '^TestSecondOpinionSessionPrepareNeverExecsOwnBackend$' ./internal/toolboxcommands`, 1 s):
   ```
   --- FAIL: TestSecondOpinionSessionPrepareNeverExecsOwnBackend (0.07s)
       report_second_opinion_test.go:577: missing typed backend summary: {"schema_version":0,"command":"","outcome":"failure",... "observed_evidence":["unexpected argument \"--prepare\"; ai-report-second-opinion [--backend codex|gemini|agy|claude [--model <slug>] | --all --exclude <own-cli>] [<bundle-dir>|latest]"] ...
   FAIL
   ```
2. **GREEN prepare** (same command, 1 s): `--- PASS: TestSecondOpinionSessionPrepareNeverExecsOwnBackend (0.08s)`. Pre-flight probe after task 2: `REQ-700 pre-flight: existing second-opinion tests, handler pin and typed-summary round trip passed` (15 s). Pre-flight probe after the task 3 refactor: passed again (10 s).
3. **RED finalize** (`go test -count=1 -run '^TestSecondOpinionSessionFinalizePublishesAgentWrittenPage$' ./internal/toolboxcommands`): all four subtests fail, each with `unexpected argument "--finalize"` inside the usage refusal:
   ```
   --- FAIL: TestSecondOpinionSessionFinalizePublishesAgentWrittenPage/publishes
   --- FAIL: TestSecondOpinionSessionFinalizePublishesAgentWrittenPage/judge_pass_inserts_banner
   --- FAIL: TestSecondOpinionSessionFinalizePublishesAgentWrittenPage/original_modified_keeps_stage
   --- FAIL: TestSecondOpinionSessionFinalizePublishesAgentWrittenPage/never_written_cleans_up
   FAIL
   ```
4. **GREEN finalize** (`go test -count=1 -v -run '^(TestSecondOpinionSessionPrepareNeverExecsOwnBackend|TestSecondOpinionSessionFinalizePublishesAgentWrittenPage)$' ./internal/toolboxcommands`, 1 s, no skips):
   ```
   --- PASS: TestSecondOpinionSessionPrepareNeverExecsOwnBackend (0.06s)
   --- PASS: TestSecondOpinionSessionFinalizePublishesAgentWrittenPage (0.32s)
       --- PASS: TestSecondOpinionSessionFinalizePublishesAgentWrittenPage/publishes (0.09s)
       --- PASS: TestSecondOpinionSessionFinalizePublishesAgentWrittenPage/judge_pass_inserts_banner (0.09s)
       --- PASS: TestSecondOpinionSessionFinalizePublishesAgentWrittenPage/original_modified_keeps_stage (0.07s)
       --- PASS: TestSecondOpinionSessionFinalizePublishesAgentWrittenPage/never_written_cleans_up (0.07s)
   ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	0.823s
   ```
5. **Contract scripts** after the prose: `action-shell-blocks.sh` 0, `prescribed-shell-canonicalization.sh` 0, `shipped-package-reference-contract.sh` 0, `audit-lockins.sh` 0.
6. **Verify set**: REQ-700 probe exit 0 (`REQ-700 probe: both session-form tests passed`); pre-flight probe exit 0; both packages `ok`; gofmt clean; vet clean; diff stat shows exactly the five files; `git diff --check` clean.

REQ-701's tests stay green by name inside the pre-flight probe: `TestSecondOpinionSubagentModelIsNotAmbiguous`, `TestSecondOpinionUnverifiedModelKeepsFinishedPage`, `TestSecondOpinionFreshGenerationModelMismatchPublishesUnverified`, `TestSecondOpinionUnverifiedModelWithoutPageFails`, plus the kept-staging assertion in `TestSecondOpinionPreservesAssetsAndRejectsUnprefixedSymlinkAndMissingBacklink`.

## Decisions

- **D-01 (DECIDE & STATE): the sidecar is a dotfile inside the bundle**, `<stem>.json` beside the stage, mode 0600, holding backend, model, prefix and the before-snapshot as mode plus sha256 digests. An agent with repository write access could edit it, but the in-session agent is already the trusted editor, the exec path never reads a sidecar, and the finalize preservation still refuses any change to a pre-existing path the record names.
- **D-02 (DECIDE & STATE): the reply stage is agent-written, not pre-created.** Prepare returns `staging_reply` as `<stem>-reply.md` without creating it; finalize fails with "no final editorial reply" (files kept) when it is missing or blank. Pre-creating it would make an empty reply indistinguishable from an unwritten one.
- **D-03 (DECIDE & STATE): the extracted helper's REQ-701 input is `unverifiedModel string`**, the requested model whose identity native evidence did not confirm; empty means verified. The exec path passes `resolved` when unverified; finalize passes `""` ("none"). The seen list comes from `row.SeenModels`, so no second parameter.
- **D-04 (DECIDE & STATE): a prepare that fails after reserving its files removes them.** The deferred cleanup in `secondOpinionPrepare` removes the stage and sidecar when the row did not reach `prepared` (for example an unresolved placeholder). No agent has written to them at that point, so this is the never-written case, not a finished page.
- **D-05 (DECIDE & STATE): in the exec path, a reply publish failure after the page publish now keeps the staging page** and names it in `reason`, where before REQ-700 the deferred cleanup removed it (the old `stagePublished` flag was set after the page publish). The stage leaves are removed only after both publishes succeed, so the never-delete-a-finished-page rule now covers this path too. No test covered the old behavior; the kept bytes are harmless and visible.
- **D-06 (DECIDE & STATE): finalize derives the bundle from the staging page's parent and refuses a supersession jump** (`staging page must sit in a current report bundle`) so a page prepared in an older revision is never published into a newer bundle.

No ESCALATE decisions.

## Discovered Tasks

- Foreign-backend exec path (codex, gemini, agy) still shows no progress while the subprocess runs; the REQ says this is a separate request. impact-minor → report only
- `skills/do-work-toolbox/actions/help.md` has no second-opinion line, so `help` never surfaces the form or its two phases. impact-minor → report only
- `actual_model` still shows the probe model on a non-identity generation failure in the exec path (left by REQ-701's builder). impact-minor → report only
- Kept staging files (finished page after a refusal, in both the exec and session paths) are never auto-cleaned; a later prepare or finalize does not sweep them. impact-minor → report only
- The handler's finding text for a finalize that finds no sidecar reads `: no prepare record for this staging page` with an empty backend name, because the backend is only known from the sidecar. impact-negligible → report only

## Lessons read

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: the two `[family: native-backend-identity]` bullets (0.306.0 REQ-698, 0.306.3 REQ-701) and the two `[family: reaped-by-its-own-parent]` bullets (0.306.0 REQ-698, REQ-543). They explain why the exec path reads identity from the native envelope and why `secondOpinionCapture` keeps its supervisor; both stay as they are. The `background-worker-self-bound` bullet was read for the proposed lesson's family choice.
- `_dev/primes/lessons-shell-commands.md` and `_dev/primes/lessons-action-files.md`: not read in full (dropped for budget by the capture); the prime sections the brief named were read instead (Closed Enumerations Go Stale; Cross-Referencing; the alternate-writer-contract-drift trap line).

## Anti-bloat check

`git diff --stat` is in [UNIFY] above. Added beyond what the REQ or Plan named:

- `secondOpinionBacklinkRequirement` (package const): the backlink sentence that fills `{{BACKLINK_REQUIREMENT}}` was a literal inside `secondOpinionRun`; prepare needs the same bytes, so one const keeps the two prompts identical instead of a second copy of the literal.

Nothing else: every function, type, flag, field, pattern and test helper is one the brief's named list contains (`--prepare`, `--finalize`, `secondOpinionPrepare`, `secondOpinionFinalize`, `secondOpinionPublication`, `secondOpinionOutsideCheck`, `secondOpinionDigests`, `secondOpinionRecord`, `secondOpinionDigest`, `secondOpinionStagePattern`, the five row fields, the sidecar file, `secondOpinionFakeLogger`). The `slices` and `crypto/sha256` imports are standard library. No new file, option, config or result-level type. Tests: exactly the two named tests and the named helper.

## Proposed CHANGELOG entry (integrator writes it with the version)

**The session's own second opinion runs where the user can watch it.** When an agent asks for a second opinion from its own model, the CLI no longer execs that agent's CLI as a hidden subprocess that sits silent for minutes; the editorial pass runs in the session, or in a background subagent or workflow the agent starts, and the CLI keeps every launcher guarantee around it.

- `ai-report-second-opinion --prepare --backend <cli> --model <slug> [<bundle-dir>|latest]` checks the bundle and outside-bundle Git status, reserves the dotfile staging page and the asset prefix, and returns one `runs[]` row with `state: prepared`, `staging_page`, `staging_reply`, `asset_prefix` and the filled `prompt` (`{{OUTPUT_FILE}}` is the absolute staging path). It never looks up or runs a backend.
- `ai-report-second-opinion --finalize <staging-page>` reads the prepare record, re-checks outside-bundle Git status and preservation, requires the backlink and a non-blank reply, publishes `index-<cli>-<model>[-N].html` with `<cli>-<model>[-N]-reply.md`, runs the render judge and inserts the pass-gated banner, exactly as the exec path does. Only a never-written staging page is removed; any other refusal keeps the page and names it.
- The prepare record is a `<stem>.json` sidecar beside the staging page holding mode and sha256 digests of every bundle file, never bytes, so preservation is proven across the process boundary.
- Session rows carry `model_evidence: session-declared` (the model the session declared, never CLI-verified); exec rows carry `model_evidence: native`.
- The exec path for codex, gemini, agy, and for `claude` when another agent targets it, is unchanged; its publication tail is now shared with finalize.
- The action's § Second-opinion form and the user guide describe the prepare, do-the-prompt, finalize sequence; the prompt template is unchanged.

## Proposed lesson bullet

Satellite: `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, new family `session-is-the-backend`.

- [family: session-is-the-backend] <version> ([REQ-700](<canonical archive link>), 2026-10-19): **when the calling agent is the backend, split the launcher at the generation call instead of spawning the agent's own CLI.** A hidden `claude -p` under the supervisor showed the user ten silent minutes and made the child's model envelope the only identity evidence. Prepare and finalize keep every mechanical guarantee (outside-bundle Git status, preservation, backlink, names, judge, banner) by carrying the before-snapshot across the process boundary as digests in a sidecar, and the model is recorded as the session's declared one, never as verified. What didn't work in review: nothing; the exec path's `os.SameFile` proof survives because the digest comparison applies only to snapshot entries rebuilt from the record.

## Integration seams and wall times

- Integration seams: none. REQ-701 was on the base before the branch existed; no other REQ in this run touches these files.
- Wall times: RED prepare 1 s; GREEN prepare 1 s; pre-flight probe 15 s, 10 s, 10 s (three runs); RED finalize under 1 s; GREEN both tests 1 s; package tests 17 s (toolboxcommands) and 0.6 s (resultmodel); REQ-700 probe 1 s; contract scripts 2 s, 1 s, 1 s, 1 s. No budget or browser probe failed, so no reruns under load were needed.
