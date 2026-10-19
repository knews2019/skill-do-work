## Review: REQ-700

**Approve** — the session's own second opinion no longer runs as a hidden shell: prepare hands the filled prompt and reserved paths to the agent, the agent does the editorial pass where the user can watch it, and finalize keeps every launcher guarantee. The prose names exactly the three places the user asked for. Six minor findings and three nits, all report only.
Route C | merge 4362eb61b6cbd037bee8f81adc1888a6b16d4a7b (range e11d3fca..4362eb61)

### What's built
- `ai-report-second-opinion --prepare --backend <cli> --model <slug> [<bundle>|latest]` reserves a dotfile staging page plus a digest sidecar and returns one `prepared` row with `staging_page`, `staging_reply`, `asset_prefix`, `prompt` and `model_evidence: session-declared`; it never looks up or runs a backend.
- `ai-report-second-opinion --finalize <staging-page>` re-checks outside-bundle Git status, preservation against the sidecar digests, backlink and reply, then publishes, judges and inserts the pass-gated banner through `secondOpinionPublication`, the tail now shared with the exec path. Only a never-written stage is removed.
- `actions/ai-report.md` § Second-opinion form and `docs/ai-report-guide.md` § A Second Opinion prescribe prepare, do the prompt in your own turns or one background subagent or workflow, finalize. Codex, gemini, agy, and `claude` targeted from another agent keep the single-call exec path with `model_evidence: native`.
- Still missing (by the REQ's own scope): progress feedback for the foreign-backend exec path; a `help.md` line for the form.

### Decisions / risks for you
- D-05 (builder): in the exec path a reply-publish failure after the page publish now keeps the staging page and names it, where before it was removed. Value: the never-delete-a-finished-page rule covers one more branch. Risk: none observed; no test pinned the old removal.
- D-01 (builder): the sidecar is a 0600 dotfile inside the bundle. The in-session agent is already the trusted editor, and finalize still refuses any change to a pre-existing path the record names. Risk: an agent could edit the sidecar to hide a change it made to an existing asset; that is the same trust the editor already has over the staging page.

### Findings

**Important:**
None

**Minor:**
- M1 `skills/do-work-toolbox/actions/ai-report.md:65` step (c): "telling it to write the page to `staging_page` and its FINAL REPLY to `staging_reply`" is attached grammatically to the subagent-or-workflow branch. An agent that does the prompt in its own turns can finish without writing the reply file, and finalize then refuses with "no final editorial reply" (files kept, a rerun after writing the reply succeeds). One clause making the two destinations apply to both branches would close it. — impact-user-visible → report only
- M2 `skills/do-work-toolbox/actions/ai-report.md:63`: the unchanged intro sentence "delegates all execution and publication to the canonical CLI" is now false for the Session form, where the CLI hands execution back to the session and owns only prepare and finalize. Stale restatement inside the rewritten section. — impact-negligible → report only
- M3 `skills/do-work-toolbox/docs/ai-report-guide.md:109`: "The page's model label is the one the session declared, and the reply says so." Nothing in the code or the action makes the editor's reply say the model was declared; the typed row carries `model_evidence: session-declared` and the action's step (e) reads the row. The sentence drifts from the code. — impact-negligible → report only
- M4 `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md` OUTPUT paragraph (file the REQ declared unchanged): "The launcher assigned the staging filename below and will publish it under the verified backend/model name." On the session path the published name carries the session-declared model, not a verified one. Stale restatement found by the sweep; not scope drift. — impact-negligible → report only
- M5 `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`: two decided behaviors have no pinning test: D-02 (page written, reply stage missing or blank → `failed` "no final editorial reply", files kept) and D-06 (a stage inside a superseded bundle → "staging page must sit in a current report bundle"). Both are one-assertion subtests of the existing finalize test. — impact-negligible → report only
- M6 `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:5` (inherited from REQ-701's sweep, its finding M2): the 0.306.0 `native-backend-identity` bullet still says identity is read "during both the availability probe and generation" and refused "before naming artifacts"; after REQ-700 the session path has no probe at all and records a declared model. Reconciled only when the integrator adds the hand-back's proposed `session-is-the-backend` bullet. — impact-negligible → report only

**Nit:**
- N1 `skills/do-work-toolbox/actions/ai-report.md:69`: "Only a plain terminal without an active agent omits both `--backend` and `--all`" — `--finalize` also omits both; the handler routes finalize before the terminal default, so the sentence is loose, not wrong. — impact-negligible → report only
- N2 `report_second_opinion.go:135` (handler finding text): a finalize whose sidecar is missing reports `: no prepare record for this staging page` with an empty backend name, because the backend is only known from the sidecar (the builder's own discovered task). — impact-negligible → report only
- N3 `report_second_opinion.go:480-490`: the never-written cleanup also removes a reply stage the agent may already have written. Per the Plan, and a reply without a page has no use, but the row does not say the reply was removed. — impact-negligible → report only

**Coding-guardrails spot-check (informational):** nothing the dimensions missed. `secondOpinionPublication` takes ten parameters, which is heavy but is the Plan's named signature; names are two-word and greppable; no speculative surface (the only unnamed addition deduplicates a literal). Surgical: the `if/else` chain in the handler became a `switch` only because two cases were added.

**Anti-bloat count: 1.** The hand-back's claim holds against the diff: the package const `secondOpinionBacklinkRequirement` replaces the backlink literal that was inline in `secondOpinionRun` (the literal is deleted at `report_second_opinion.go` old line 410, the const fills `{{BACKLINK_REQUIREMENT}}` at new lines 437 and 603). Every other package-level addition is on the brief's named list, checked one by one against `grep '^func \|^const \|^var \|^type '` at 4362eb61: `secondOpinionStagePattern`, `secondOpinionRecord`, `secondOpinionDigest`, `secondOpinionDigests`, `secondOpinionOutsideCheck`, `secondOpinionPrepare`, `secondOpinionFinalize`, `secondOpinionPublication`, the `prepare`/`finalize` fields of `secondOpinionOptions`, the `--prepare`/`--finalize` flags and the usage string, `secondOpinionFile.mode`/`.digest`, the `secondOpinionPreserved(bundle, before, allowed, prefix)` signature, the five `omitempty` row fields, the `<stem>.json` sidecar, `secondOpinionFakeLogger`, `TestSecondOpinionSessionPrepareNeverExecsOwnBackend`, `TestSecondOpinionSessionFinalizePublishesAgentWrittenPage` with its four named subtests. Unnamed items without reach, not counted: the `unverifiedModel` parameter name (D-03), test-local closures `prepare` and `writePage`, handler locals `stagePath`/`stageName`/`template`, and the finding-text source change from `backend` to `row.BackendName` (needed because finalize learns the backend from the sidecar). Imports `crypto/sha256` and `slices` are standard library. No new file, command, option, config or result-level type; `TestHandlersRegisterCanonicalToolboxCommands` still pins 10 handlers.

**Editorial pass against the user's words.** UR-158 verbatim: "this should have just triggered a normal prompt inside the ai-coder that way I can see the progress, alternatively it's fine to start a background subagent, or a workflow, but not a shell like this".
- Three places, exactly: action step (c) "do the `prompt` yourself in your own turns, or hand it verbatim to one background subagent or workflow"; guide "runs inside the session itself, where you can watch it read the report and write the page, or in a background subagent or workflow the agent starts". No fourth place is offered. "one" background subagent is a tightening the brief added, not a drift from the user's words.
- Hidden shell forbidden: Session bullet "never as a subprocess of your own CLI"; Targeted bullet "A target equal to your own CLI takes the Session sequence above instead of a subprocess"; the invoke paragraph's `claude` sentence now only covers another agent targeting it.
- No sentence in either file still instructs the agent to run the editorial pass as a subprocess. The remaining CLI calls are prepare and finalize, both short.
- Progress visibility: the guide's "where you can watch it read the report and write the page" matches "that way I can see the progress"; the action states the reason ("where the user can watch it").
- Drifts noted: M1 (reply destination attached to one branch only), M2 (intro sentence still says the CLI executes everything), M3 ("the reply says so" has no code behind it).

### Requirements Checklist

REQ `## Requirements` and the Scope's acceptance criteria, cross-checked against the Plan and the UR:
- [x] Session form and own-CLI target never launch a subprocess of the agent's CLI; work is visible in the session — delivered (handler `switch` runs finalize or prepare before any `exec.LookPath`; prepare calls no capture, argv or metadata helper; both tests assert the fake `claude` on PATH was never invoked; prose step (c))
- [x] Prepare prints the filled prompt, reserved staging page, reserved reply path and asset prefix — delivered (`secondOpinionPrepare` lines 401-450; `{{OUTPUT_FILE}}` is the absolute stage path; the prepare test asserts all four fields and no `{{`)
- [x] Finalize keeps every launcher guarantee — delivered (outside check 477, never-written cleanup 480-490, digest preservation with the three leaves allowed 496, reply required 502, `secondOpinionPublication` 651: backlink, `secondOpinionNames`, two publishes, judge, outside re-check, inode and byte identity of `index.html`, banner; subtests `publishes`, `judge pass inserts banner`, `original modified keeps stage`)
- [x] A finished staging page is never deleted by finalize; only an empty stage is cleaned up with its sidecar and reply stage — delivered (lines 480-490 are the only removal; subtests `original modified keeps stage`, `never written cleans up`)
- [x] In-session model recorded as `model_evidence: session-declared`, never CLI-verified; exec rows say `native` — delivered (handler line 108, finalize line 476, exec line 625; both prose files)
- [x] Final reply lists page and reply paths, backend and model, judge verdict, and the 3 to 5 clarity changes, corrections and open questions — delivered by prose (step (e) and the unchanged Finish paragraph); the prompt template's FINAL REPLY rule is unchanged and finalize requires a non-blank reply
- [x] Foreign backends keep the single-call launcher path; their tests stay green — delivered (`secondOpinionArgv`, `secondOpinionMetadata`, `secondOpinionCapture` and the supervisor are absent from the diff; exec callers pass `nil, ""` and `[]string{stageName}`; family run 19+ tests green at the merged tree)
- [x] Action § Second-opinion form and guide § A Second Opinion describe the in-session sequence; CLI, tests and changelog ship as one release — delivered for the builder's part; the changelog, version bump and mirrors are the integrator's finalization (Scope says so; the `scope-drift` warnings in Qualification are those files)
- [x] UR: "not a shell like this" — delivered (see editorial pass)
- [x] Plan § Consumer field contract: prepare rows carry `staging_page`, `staging_reply`, `prompt`, empty `page`/`reply`, `judge_verdict: not-run`, `banner_state: not-inserted`; finalize rows carry `page`, `reply`, `judge_verdict`, `banner_state` with staging fields empty; `mode: session` for both — delivered (handler lines 80-81 and 103-109; `publishes` subtest)
- [x] REQ-701 behaviors carried into the extracted tail: `-unverified` names, no deletion of a finished page, `seen_models` in the reason — delivered (lines 654-659; REQ-701's four tests green by name in the family run)

Implementation Summary claims verified against the diff: all five files and every listed change appear in `git diff e11d3fca..4362eb61`. P-A-U: all three boxes `[x]` with the hand-back's text.

**Existing tests updated (cross-REQ):** none modified; REQ-701's tests and REQ-698's tests run unchanged. The behavior change D-05 (exec reply-publish failure keeps the stage) is documented in the hand-back and the REQ's Qualification; no test covered the old behavior, so no test update was needed. Traceability intact.

### Acceptance Testing

**Result: Pass** — Implementation and Integration stages: the diff was read in full; at the merged tree 4362eb61, `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands` from `skills/do-work/tools/do-work-cli` printed `ok` in 12.9 s, exit 0 (the two new tests, four subtests, REQ-701's four tests and REQ-698's family); `go test -count=1 ./internal/resultmodel` printed `ok` in 0.2 s, exit 0. Not exercised: Deployment (installed consumer copy) and Live acceptance (a real Claude session running the Session form); see Suggested Additional Testing.

**Red-green trace** (hand-back § Proof record against `## Red-Green Proof`):
- RED prepare: `TestSecondOpinionSessionPrepareNeverExecsOwnBackend` failed at base with `unexpected argument "--prepare"` inside "missing typed backend summary". This is the REQ's RED (the CLI had only the exec path for the session's own model); the test is the Red-Green Proof's named Go test (fake `claude` on PATH records invocations; the session path must record zero and hand back the prompt and reserved path). GREEN: PASS (0.08 s).
- RED finalize: all four subtests of `TestSecondOpinionSessionFinalizePublishesAgentWrittenPage` failed with `unexpected argument "--finalize"`. GREEN: all four PASS. This covers the Red-Green Proof's GREEN clause "the CLI then finalizes it with the same preservation, Git, judge and banner guarantees".
- The GREEN clause "visible work in the session" and "the final reply lists page, reply, model, judge verdict" are prose, verified by the editorial pass above, not by a Go test.
- Reproduced by this review: the family run above at the merged tree; the fake backend's never-invoked assertion holds in both tests.

Red flags checked: none present (Implementation Summary files all appear in the diff; new helpers are all called; assertions are substantive; `tdd: true` with red-green evidence; Scope 5 builder files = 5 diff files).

### Suggested Additional Testing

- Live acceptance: unassessed. In a real Claude Code session run `/do-work-toolbox ai-report second-opinion ai-reports/<bundle>/index.html` and confirm the only shells are the prepare and finalize calls (seconds each), the editorial pass is visible as the agent's own tool calls, the agent writes `staging_reply`, and finalize publishes with the banner on a judge `pass`.
- Deployment: unassessed. Install the release into a consumer project through the updater and confirm the installed `actions/ai-report.md` carries the Session sequence and the installed CLI accepts `--prepare` and `--finalize` (the gate's `updater` and `installer` heavy lanes cover the mechanics).
- Edge (M1): do the prompt in your own turns without writing the reply file; confirm finalize says "no final editorial reply" with "staged page kept at", then write the reply and rerun finalize to a successful publication.
- Edge (D-06): run `--finalize` on a stage inside a bundle that `catalog.json` marks superseded; expect "staging page must sit in a current report bundle" and no publication.
- Edge: `--finalize` with a repository-relative stage path (the tests use absolute paths only).
- Regression: from a non-Claude agent, `--backend claude --model <slug>` still takes the exec path and its row says `model_evidence: native`.

### Scores (on the record — not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | Every REQ requirement, acceptance criterion and the user's three places delivered |
| Code Quality | 92% | Follows the file's patterns; comments explain why; ten-parameter tail is the Plan's shape; M2/M4 prose drift |
| Test Adequacy | 88% | Red-green traced for both phases; D-02 and D-06 unpinned (M5) |
| Scope | 100% | Exactly the five builder files; integrator-only files untouched by the builder |
| Risk | Low | Exec path unchanged by diff; finalize validates the stage name and bundle; sidecar trust equals the editor's existing trust |
| Acceptance | Pass | Implementation and Integration stages |

(100 + 92 + 88 + 100) / 4 = 95; no modifiers.

### Follow-ups created
- None (9 findings report only)

## Review

**Overall: 95%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** M1 `actions/ai-report.md:65` the FINAL-REPLY-to-`staging_reply` instruction is attached only to the subagent branch, so an own-turns pass can skip the reply file and finalize refuses (files kept) — impact-user-visible → report only; M2 `actions/ai-report.md:63` "delegates all execution and publication to the canonical CLI" is stale for the Session form — impact-negligible → report only; M3 `docs/ai-report-guide.md:109` "and the reply says so" has no code or action rule behind it — impact-negligible → report only; M4 `docs/ai-report-second-opinion-prompt.md` says the launcher publishes "under the verified backend/model name", untrue on the session path — impact-negligible → report only; M5 D-02 (reply missing, page kept) and D-06 (supersession guard) have no pinning test — impact-negligible → report only; M6 `lessons-do-work-cli.md:5` (inherited, REQ-701's M2) still describes probe-and-generation identity for every run; reconciled only by the proposed `session-is-the-backend` bullet — impact-negligible → report only. Nits: N1 `actions/ai-report.md:69` terminal-default sentence ignores `--finalize` — impact-negligible → report only; N2 finalize-without-sidecar finding text has an empty backend prefix — impact-negligible → report only; N3 never-written cleanup also removes an agent-written reply stage without saying so — impact-negligible → report only. Anti-bloat count: 1 (`secondOpinionBacklinkRequirement`, deduplicates the backlink literal); every other addition is on the brief's named list.
**Acceptance:** Pass — Implementation and Integration stages: `go test -count=1 -run '^TestSecondOpinion' ./internal/toolboxcommands` ok 12.9 s and `go test -count=1 ./internal/resultmodel` ok at 4362eb61; Deployment and Live acceptance unassessed
**Restatement sweep:** redefined the `ai-report-second-opinion` argument contract (`--prepare`, `--finalize`, usage string), the run-row schema (`staging_page`, `staging_reply`, `asset_prefix`, `prompt`, `model_evidence`, state `prepared`), the meaning of `mode: session` (prepare and finalize phases, no exec of the agent's own CLI), the model-identity rule (session rows declared, exec rows native), the `{{OUTPUT_FILE}}` fill (absolute stage path in-session), and the staging-leaf cleanup rule (leaves removed after both publishes, D-05); wave end, inherited from REQ-701's recorded line: the `-unverified` filename contract, the row schema (`seen_models`, state `unverified`), the `AI-REPORT-SECOND-OPINION-<STATE>` code set, and the identity rule (mismatch publishes unverified, a subagent model is not ambiguity). Swept `skills/` and `_dev/` (no `docs/` at the root): agreeing — `actions/ai-report-reference.md:148` (finalize is the launcher's phase, so "only the launcher may add a pass-gated banner" holds; `#second-opinion-form` anchor kept), `docs/ai-report-guide.md:72` usage line, `SKILL.md:3,24` trigger words, `actions/help.md` (no second-opinion line), `tools/do-work-cli.sh:10` (prepare reads the same env var), `_dev/tests/do-work-cli-launcher-behavior.sh:52,252` (export only), `lessons-do-work-cli.md:7` (REQ-701 bullet still true), CHANGELOG 0.306.0 and 0.306.3 (history), `commands_test.go` 10-handler pin, `result_model_test.go` round trip, and the inherited `-unverified`/`seen_models`/`-UNVERIFIED` contract (code at `secondOpinionPublication` 654-659, action line 73, guide line 107, REQ-701's tests green); stale — M2 (`ai-report.md:63`), M4 (prompt template OUTPUT paragraph), M6 (`lessons-do-work-cli.md:5`), N1 (`ai-report.md:69`)
**Suggested testing:** 6 items
**Follow-ups created:** None (9 findings report only)

*Reviewed by review-work action*

### Delta re-review (4362eb61..c70b9b72)

Read the delta only (`git diff 4362eb61..c70b9b72`: one file, `skills/do-work-toolbox/actions/ai-report.md`, two sentences, builder commit 2ba229e4, re-merged as c70b9b72). No tests or gate run for this delta, as the integrator asked.

- **M1 resolved.** Step (c) now reads "do the `prompt` yourself in your own turns, writing the page to `staging_page` and your FINAL REPLY to `staging_reply`, or hand it verbatim to one background subagent or workflow told to write to those same two paths". Both destinations now apply to the own-turns branch and to the delegated branch, which matches finalize (`report_second_opinion.go:496-502`: the reply stage is required non-blank, the page is the stage).
- **M2 resolved.** The intro sentence now reads "delegates every check and publication to the canonical CLI, which also runs the editorial pass for every backend except the session's own". That matches the handler at c70b9b72: the `switch` routes `--prepare` and `--finalize` to the session phases and every other backend (codex, gemini, agy, and `claude` targeted from another agent) to `secondOpinionRun`'s exec path.
- **User's three places still exact:** own turns, one background subagent, a workflow; no fourth place added, and the "never as a subprocess of your own CLI" sentence is unchanged. The Targeted and Batch bullets are untouched.
- **New defect or bloat:** none. Prose only; no shell fence, no new citation, heading text unchanged (the `#second-opinion-form` anchor still resolves).

**Updated scores:** Requirements 100% (unchanged), Code Quality 92% → 95% (the two prose drifts in the action file are gone; M3 and M4 remain in other files), Test Adequacy 88% (unchanged, M5 stands), Scope 100% (unchanged). Risk Low, Acceptance Pass (Implementation and Integration stages, from the 4362eb61 run; the prose delta needs no test). **Overall 95% → 96%** ((100 + 95 + 88 + 100) / 4 = 95.75).

**Remaining findings (all report only):** M3 `docs/ai-report-guide.md:109` "and the reply says so" — impact-negligible; M4 prompt template "verified backend/model name" — impact-negligible; M5 D-02 and D-06 unpinned by tests — impact-negligible; M6 `lessons-do-work-cli.md:5` inherited stale identity bullet — impact-negligible; N1 `ai-report.md:69` terminal-default sentence ignores `--finalize` — impact-negligible; N2 empty backend prefix in the no-sidecar finding text — impact-negligible; N3 never-written cleanup removes an agent-written reply stage silently — impact-negligible. **Follow-ups created:** None (7 findings report only).
