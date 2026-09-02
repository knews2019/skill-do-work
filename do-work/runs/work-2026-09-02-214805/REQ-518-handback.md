# REQ-518 Builder Handback

## Outcome

REQ-518's source implementation is complete for review. The final integration correction replaces every newly prescribed bare gate-evidence CLI invocation in `work.md` and `work-reference.md` with the shipped launcher form:

`<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json ...`

The existing repository-gate contract predicate and mutation now require that form at the baseline check/record seam and the mandatory final-gate record seam. `_dev/tests/contract-regressions.sh` remains line-neutral at 8,440 lines.

## Files Changed

- `skills/do-work/actions/work.md` — checks typed green-gate evidence through the shipped launcher before the baseline gate, records direct green baseline/final runs through the same launcher, and keeps the final direct gate mandatory.
- `skills/do-work/actions/work-reference.md` — defines the same launcher, typed `gate_evidence`, `baseline_revision`, Git-private record, and late-attribution behavior in the canonical reference.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` — registers the two gate-evidence handlers at the public command boundary.
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/gate_evidence_integration_test.go` — exercises record, exact match, project-change invalidation, and gate-log-only descendant matching through the built public CLI.
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence.go` — stores private exact-argv evidence under the canonical Git common directory and validates repository identity, revision existence/ancestry, and every intervening commit path.
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_evidence_test.go` — covers identity/history matching, divergent/foreign records, malformed or unsafe targets, atomic replacement, permissions, and linked-worktree reuse.
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands.go` — parses the two commands, returns typed results, refuses non-green recording without mutation, and emits non-self-referential recovery evidence.
- `skills/do-work/tools/do-work-cli/internal/gateevidence/gate_commands_test.go` — covers handler registration, exact argv parsing, non-green refusal/non-mutation, and text/JSON parity.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` — adds the typed `gate_evidence` projection, stable states, normalization, and text rendering.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` — proves gate-evidence text/JSON parity and non-null argv normalization.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — indexes the new package owner and focused test command.
- `_dev/tests/contract-regressions.sh` — strengthens the existing baseline/final predicates and mutation without adding lines or a new sentence predicate.

No REQ, checkpoint, run manifest, release/version file, unrelated report, or file outside the declared Scope was edited by this builder pass.

## Implementation

- `record-green-gate` accepts the exact gate argv after `--` and only writes when `--gate-exit-status` is `0`. It records the current full `HEAD`, canonical Git common-directory identity, exact argv plus SHA-256 digest, truthful reported provenance, and zero exit status in a private `0600` JSON file under a private Git-directory folder.
- `check-green-gate` is read-only. It returns top-level success for a valid match or valid non-match and exposes the decision in typed `gate_evidence`. Exact `HEAD` matches. A descendant matches only when every intervening commit changes exclusively paths below `_dev/gate-runs/`. Missing, foreign, different-argv, missing-revision, divergent, project-changing, malformed, or unsafe evidence does not match; unverifiable evidence fails closed.
- The work action consumes only typed evidence. `matches: true` saves `baseline_revision` and skips the duplicate pre-build gate. `matches: false` retains the direct baseline. Every direct green baseline and final gate records evidence. The Step 6.5 gate remains direct and unconditional, and existing fingerprint/deferral/repair attribution remains unchanged.
- The action and canonical reference now invoke the real shipped launcher with explicit `--repo-root <project-root>` and JSON output, so execution does not depend on a globally installed `do-work-cli` or the caller's current directory.

## RED / GREEN Evidence

### RED

The builder brief required the public CLI test to be run before registration and retained with an assertion-level `UNKNOWN-COMMAND` result. The test and planned failure are present, but no verbatim pre-implementation output was retained in the run directory or REQ. Therefore the historical test-first RED is an evidence gap and is not claimed as observed.

The durable pre-change defect statement remains factual: the prior action text directly reran the canonical baseline and had no check/record command or public gate-evidence handler. The current Git diff shows those additions, but that diff is not a substitute for retained RED test output.

### GREEN

- `go test -count=1 ./cmd/do-work-cli -run TestGreenGateEvidenceLifecycle` — PASS. The public CLI lifecycle covers exact match, project-change non-match, and `_dev/gate-runs/`-only descendant match.
- `go test -count=1 ./internal/gateevidence ./internal/resultmodel ./cmd/do-work-cli` — PASS for all three focused packages.
- `go test -race -count=1 ./internal/gateevidence` — PASS.
- `go vet ./internal/gateevidence ./internal/resultmodel ./cmd/do-work-cli` — PASS, exit 0 with no output.
- `bash _dev/tests/contract-regressions.sh` — PASS, ending with `Contract regression checks passed.` Its launcher behavior, shipped-reference, shell-block, mutation, and repository-gate contract stages all completed successfully.
- `skills/do-work/tools/do-work-cli.sh --repo-root /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 --format json check-green-gate -- bash _dev/tests/maintainer-verify.sh` — PASS, `outcome: success`, typed `state: missing`, `matches: false`, and the expected repository root. This proves the prescribed launcher reaches the new public command from the project root.
- `gofmt -d` over all eight changed/new Go files — PASS with no output.
- `git diff --check` — PASS with no output.

The full Go module was not rerun because the captured pre-flight baseline already records unrelated failures in `internal/corehelpers`, and this handback was explicitly limited to focused tests. The canonical `bash _dev/tests/maintainer-verify.sh` final gate was also not rerun in this builder pass; the REQ pre-flight records it green at `bab2198d59e891ec23d98a209c3c03187bc1741d`, and the queue owner retains responsibility for the mandatory integrating gate.

## P-A-U

- **[PLAN]:** Followed the saved plan's package/result/action/reference seams and narrowed this pass to the identified launcher integration defect plus its existing line-neutral contract coverage.
- **[APPLY]:** Replaced all four new check/record prescriptions with the established `<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json` form and updated the existing two predicates plus existing mutation string. No scope expansion occurred.
- **[UNIFY]:** Reviewed the complete current source change across all 12 files above, checked action/reference agreement, checked every gate-success consumer, confirmed no debug markers or formatting drift, kept the contract at 8,440 lines with `4` additions and `4` deletions, and ran the focused checks listed above.

## Decisions

- **D1:** Keep the explicit launcher inline in both action owners. This is the established executable contract and makes the required project root and JSON rendering unambiguous.
- **D2:** Strengthen the existing `structured direct baseline` and `late base attribution` predicates to include the launcher and `--repo-root`; do not add a sentence predicate or grow the contract file.
- **D3:** Do not broaden transparent history beyond `_dev/gate-runs/`, exempt later lifecycle commits, change gate fingerprinting, or repair unrelated `internal/corehelpers` baseline failures. Those changes are outside REQ-518.

## Lessons Consulted

- `_dev/primes/prime-action-files.md` and `_dev/primes/prime-shell-commands.md` — kept the two action writers synchronized, used structured argv, and prescribed the real shipped launcher instead of relying on shell/session state.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and `lessons-do-work-cli.md` — kept typed evidence and public registration at their existing ownership seams, used Git-private atomic state, and preserved exact evidence rather than an opaque boolean cache.
- `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`, and `testing.md` — held scope, reviewed each changed file, and reported the missing historical RED instead of inferring it from current GREEN results.

## Discovered Tasks

None.

## Blockers

No implementation or focused-test blocker remains. The only evidence limitation is the missing retained pre-implementation `UNKNOWN-COMMAND` RED output; independent review should treat that as a TDD-history gap, not as a current GREEN failure.
