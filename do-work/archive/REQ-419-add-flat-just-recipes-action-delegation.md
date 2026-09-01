---
id: REQ-419
title: 'Add flat Just recipes, collision validation, action delegation, and compatibility aliases'
status: completed-with-issues
created_at: 2026-08-29T20:28:26Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md, _dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: [REQ-418]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-420]
batch: go-no-llm-command-platform
claimed_at: 2026-09-01T12:53:31Z
completed_at: 2026-09-01T15:06:50Z
route: C
kb_status: pending
kb_entry:
estimate:
  p50_active_minutes: 105
  confidence: low
  calculated_at: 2026-09-01T12:54:34Z
  basis:
    - Route C
    - 24-file write set
    - 2 new files
    - 6 subsystems involved
    - 8 acceptance criteria
    - dependency depth 14
    - cross-route regression gates
    - full-suite verification
write_set:
  - justfile
  - skills/do-work-board/justfile.template
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go
  - skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go
  - skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go
  - skills/do-work-toolbox/actions/note.md
  - skills/do-work-toolbox/actions/architecture-report.md
  - skills/do-work-toolbox/actions/ai-report-reference.md
  - skills/do-work-toolbox/actions/present-work.md
  - skills/do-work-toolbox/actions/install.md
  - skills/do-work-toolbox/actions/maintainability-audit.md
  - skills/do-work-toolbox/actions/maintainability-audit-reference.md
  - skills/do-work/actions/version.md
  - README.md
  - skills/do-work/actions/help.md
  - skills/do-work-knowledge/actions/help.md
  - skills/do-work-toolbox/actions/help.md
  - skills/do-work/docs/version-guide.md
  - skills/do-work/docs/command-line-guide.md
  - skills/do-work/docs/prescribed-shell-primitives.md
  - skills/do-work/tools/prime-do-work-update.md
  - skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md
  - skills/do-work-board/actions/board.md
  - skills/do-work-board/docs/board-guide.md
  - _dev/tests/flat-just-recipes-behavior.sh
  - _dev/tests/contract-regressions.sh
  - _dev/tests/staged-skills-contract.sh
  - _dev/tests/install-suite-behavior.sh
  - _dev/tests/update-script-behavior.sh
---

# Add Flat Just Recipes, Collision Validation, Action Delegation, and Compatibility Aliases

## What
Expose every public command through flat Just recipes and make the existing skill aliases delegate to the canonical CLI.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reopened TDD from the failed 50% review. Reproduce command substitution and splitting through the actual shipped template, replace raw variadic interpolation with per-recipe Just positional argv, make publication recovery execute that same shipped seam, add exact NUL/sentinel assertions for every hostile class, sweep the five newly scoped ownership restatements plus stale action/reference wording, then rerun scope, focused, aggregate, and maintainer gates.
- [x] **[APPLY]:** Confirmed the actual template executed a `$(touch …)` payload and reparsed a newline as a second shell command. Applied `[positional-arguments]` plus quoted `"$@"` to every non-board variadic recipe in both managed surfaces (preserving fixed-parameter CLI ordering with local shift variables), made publication recipes carry their canonical flags, replaced the invented caller seam, and aligned all mandatory restatements within the expanded 30-file scope.
- [x] **[UNIFY]:** Inspected the complete declared 30-file project change, including both 33-recipe positional-argv rewrites, actual-template hostile tests, publication protocol/flags, all canonical action/help/guidance changes, template-derived installer logic, five expanded-scope restatements, and every shell contract. `git diff --stat`/`git diff --check`, scope-drift, gofmt, shell syntax/ShellCheck, focused Go ordinary/race/vet, exact 40/33/33 recipe invariants, actual hostile NUL/sentinel probes, install/update/staged/aggregate suites, and the final-state canonical maintainer gate all passed. Added-line scans found no debug/TODO artifacts or hidden Unicode controls.

## Detailed Requirements
- Add flat recipes for all core, knowledge, toolbox, interview, memory, Dream, BKB, media, last30days, and absorbed audit-metrics public commands named by UR-081.
- Preserve existing board recipes and `run-do-work-update` as compatibility aliases.
- Implement dynamic reserved-recipe collision validation for the expanded generated section.
- Update natural-language actions, help, guides, install/update behavior, and upgrade contracts so deterministic phases delegate to `do-work-cli`.
- Missing or failed canonical tooling must stop with actionable output and never fall back to free-form mutation.

## Constraints
- The public aliases remain unchanged; the deterministic implementation becomes singular.
- Just recipes must remain flat and directly runnable without an LLM.

## Dependencies
Depends on REQ-418 (all command families must exist before publishing the complete interface).

## Builder Guidance
Certainty level: Firm. Derive the reserved recipe set from the actual managed template rather than duplicating a hard-coded list.

## Red-Green Proof
**RED prompt/case:** List and invoke every advertised recipe in a fresh install, test a custom collision for each reserved recipe source, and invoke an action with the CLI missing.
**Why RED now:** Only a handful of board/update recipes exist and natural-language actions still own deterministic mutations.
**GREEN when:** Every advertised command runs mechanically via Just, compatibility aliases work, collisions are refused dynamically, and action delegation stops actionably on canonical-tool failure.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

---

## Triage

**Route: C** - Complex

**Reasoning:** The request publishes a cross-suite command interface, changes natural-language action delegation and installer/update contracts, and requires hostile-path shell-safety plus whole-suite compatibility evidence.

**Planning:** Required

## Review Addendum — Shell-Safe Publication Recipe Arguments

REQ-413's fresh re-review found that publication results render manifest paths with Go double-quoted strings. When pasted into a shell-backed Just recipe, valid paths containing `$HOME`, `$(...)`, or backticks can expand or execute even though machine `next_argv` remains exact.

### Additional Requirements

- Render every generated publication recipe argument with shell-safe, byte-preserving quoting suitable for direct execution by the emitted flat Just recipe.
- Keep `next_argv` exact and ensure the human-runnable recipe represents the same argument bytes without expansion.
- Cover spaces, single and double quotes, dollar signs, command substitutions, backticks, tabs, and newlines where the interface permits them.

### Additional Red-Green Proof

**RED prompt/case:** Generate a publication recovery recipe for manifest paths containing shell expansions and quoting metacharacters, execute it through the advertised Just/shell boundary, and compare received argv bytes.
**Why RED now:** `strconv.Quote` emits Go double quotes, which are not shell-literal for `$`, command substitution, or backticks.
**GREEN when:** The generated recipe executes without expansion and passes byte-identical arguments for every supported hostile path shape.

## Plan

1. **Publish the complete flat recipe surface and derive collision ownership dynamically.** Expand `skills/do-work-board/justfile.template` and the repository `justfile`, preserve the board/update aliases, make installer collision validation derive names from the managed template, and add install/update/contract coverage for every derived recipe.
2. **Delegate remaining natural-language deterministic phases.** Point note, architecture-report, media, portfolio, last30days, and audit-metrics action phases at the canonical core launcher; retain judgment and consent in prose; stop actionably when the canonical tool is missing or fails.
3. **Make publication recovery recipes shell-literal.** Replace Go-string quoting with byte-preserving shell quoting while keeping `next_argv` exact, then execute hostile values through the emitted shell/Just boundary and compare received argument bytes.
4. **Update discoverability and upgrade contracts.** Document the expanded deterministic Just interface, canonical `do-work-update`, compatibility alias, and the boundary between mechanical recipes and natural-language workflows without duplicating a closed inventory in every guide.

**Architectural decisions:** The managed Just template is the recipe-namespace authority; natural-language routes remain unchanged; compatibility shells and the separate audit-metrics tree remain until REQ-420; publication recovery recipes use POSIX single-quote encoding; no `queue-kanban` Go source changes are planned.

**Testing approach:** Run TDD at the real caller seams: publication shell/Just execution, dynamically derived collision fixtures, fresh install/update publication, canonical-delegation semantic contracts, focused Go tests with race/vet, and the canonical maintainer gate.

**Plan validation:** Passed — all Detailed Requirements, Constraints, and Review Addendum requirements map to one or more tasks; all four tasks trace to captured requirements; task count remains below the five-task split warning.

*Generated by Plan agent*

## Exploration

- The managed template and root `justfile` currently expose 21 definitions: four board recipes, `run-do-work-update`, and 16 knowledge recipes. The complete surface becomes 40 definitions after adding 12 core and seven toolbox recipes.
- `managedsection.JustDefinitionNames` already provides the canonical parser. The stale duplication is `suiteinstall.reservedRecipeNames` plus hard-coded shell-test lists.
- All required handlers are registered. Important mappings include `do-work-update` → `update-suite`, `architecture-report-preflight --scan|--publish`, `install-last30days check|install`, and `audit-metrics inventory|folders|churn|hotspots`.
- `publicationProtocol` is the only unsafe recovery-recipe site using `strconv.Quote`; `knowledgecommands.quoteRecipeArgument` provides the existing POSIX single-quote pattern. `next_argv` already stores exact argument bytes and should remain unchanged.
- Existing canonical-delegation patterns live in the BKB, Dream, interview, memory, cleanup, and forensics actions: invoke the core launcher, consume typed output, and stop without prose fallback on missing, failed, or malformed canonical tooling.
- Retained toolbox scripts, standalone audit-metrics sources, CLI registration/result-model code, managed-section parsing code, and `queue-kanban` Go sources remain outside this REQ.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `justfile` (modify) — source-checkout flat recipe surface
- `skills/do-work-board/justfile.template` (modify) — installed managed recipe authority
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modify) — derive collision/completeness names from template bytes
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modify) — dynamic collision and candidate-completeness coverage
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go` (modify) — shell-literal recovery recipe arguments
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go` (modify) — hostile-byte shell/Just execution coverage
- `skills/do-work-toolbox/actions/note.md` (modify) — canonical note delegation
- `skills/do-work-toolbox/actions/architecture-report.md` (modify) — scan/publication delegation
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modify) — media lifecycle delegation
- `skills/do-work-toolbox/actions/present-work.md` (modify) — portfolio publication delegation
- `skills/do-work-toolbox/actions/install.md` (modify) — last30days delegation
- `skills/do-work-toolbox/actions/maintainability-audit.md` (modify) — canonical audit-metrics delegation
- `skills/do-work-toolbox/actions/maintainability-audit-reference.md` (modify) — remove deterministic fallback contract
- `skills/do-work/actions/version.md` (modify) — canonical update delegation
- `README.md` (modify) — deterministic command discoverability
- `skills/do-work/actions/help.md` (modify) — core recipe entry-point guidance
- `skills/do-work-knowledge/actions/help.md` (modify) — knowledge recipe guidance
- `skills/do-work-toolbox/actions/help.md` (modify) — toolbox recipe guidance
- `skills/do-work/docs/version-guide.md` (modify) — canonical and compatibility update recipes
- `skills/do-work/docs/command-line-guide.md` (new) — canonical non-LLM command/recipe guide
- `skills/do-work/docs/prescribed-shell-primitives.md` (modify) — align runtime authority with canonical CLI ownership
- `skills/do-work/tools/prime-do-work-update.md` (modify) — align update-recipe ownership with the canonical command
- `skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md` (modify) — mark retained source as compatibility/parity surface rather than active fallback
- `skills/do-work-board/actions/board.md` (modify) — describe the expanded managed recipe surface and canonical update recipe
- `skills/do-work-board/docs/board-guide.md` (modify) — remove the legacy-only managed-interface gloss
- `_dev/tests/flat-just-recipes-behavior.sh` (new) — complete recipe execution matrix
- `_dev/tests/contract-regressions.sh` (modify) — action-delegation and inventory contracts
- `_dev/tests/staged-skills-contract.sh` (modify) — replace legacy-script/action and compatibility-recipe assertions that contradict canonical CLI delegation
- `_dev/tests/install-suite-behavior.sh` (modify) — fresh-install managed recipe coverage
- `_dev/tests/update-script-behavior.sh` (modify) — update replacement and compatibility coverage

**Files I will NOT touch:** CLI command registration/result-model code, retained compatibility scripts, standalone audit-metrics sources, managed-section parser code, `queue-kanban` Go sources, or release/version/changelog files during builder implementation.

**Acceptance criteria (restated from REQ):**
- [x] All 40 managed recipe definitions expose the required core, knowledge, toolbox, interview, memory, Dream, BKB, media, last30days, audit, board, and compatibility surfaces through flat Just entry points.
- [x] Existing board recipes and `run-do-work-update` remain compatible while `do-work-update` is canonical.
- [x] Reserved-name collisions and candidate completeness derive from the actual managed template rather than a hard-coded list.
- [x] Named natural-language deterministic phases delegate once to `do-work-cli`, preserve action-owned judgment, and stop actionably without fallback on missing/failed/malformed tooling.
- [x] Publication `next_argv` remains exact and the generated Just/shell recipe passes byte-identical hostile arguments without expansion.
- [x] Help, guides, install/update behavior, and upgrade contracts distinguish deterministic recipes from natural-language workflows and point to `just --list` as the live inventory.
- [x] Fresh install and update publish the complete managed interface while preserving exterior custom bytes and rejecting collisions before mutation.
- [x] Focused Go, race, vet, recipe, install, update, contract, and canonical maintainer gates pass.

## Decisions

### D-01 — Extend Scope to the staged shipped-package contract

**Decision: DECIDE & STATE.** Add `_dev/tests/staged-skills-contract.sh` to Scope and `write_set` because its active assertions require the exact legacy-script calls and update-recipe body that REQ-419 replaces. Leaving it untouched would make the repository's canonical gate reject the requested behavior. This is requirement-forced regression coverage, not adjacent cleanup.

- **D-02:** Treat the complete board-managed template as the only recipe namespace authority. Installer completeness now sorts and checks every definition parsed from those bytes; collision tests derive the same set instead of maintaining another inventory.
- **D-03:** Encode publication recovery arguments with POSIX single-quote literals and the standard close-quote/double-quote/reopen splice for embedded apostrophes. This preserves every supported byte through shell and Just while leaving `next_argv` untouched.
- **D-04:** Keep natural-language routes responsible for judgment, drafting, and consent, but make missing, failed, or malformed canonical tooling terminal for each deterministic phase. Retained compatibility scripts are not action fallbacks and remain for REQ-420's shim cutover.

- **D-05:** Extend remediation Scope to the five stale contract files named by the mandatory Restatement Sweep. Each actively assigns canonical ownership or fallback behavior contradicted by REQ-419; correcting them is required to make the singular-authority/no-fallback requirement true for future agents, not adjacent documentation cleanup.
- **D-06:** Use Just's per-recipe `[positional-arguments]` attribute rather than a file-global setting. Quoted `"$@"` preserves the original argument vector without changing board or consumer-owned exterior recipes; named leading parameters are captured and shifted before reconstructing their canonical CLI flag order.
- **D-07:** Make publication recovery recipes spell `--manifest` and `--at`, matching the exact canonical `next_argv` semantics after the flat recipe contributes only launcher/root resolution. Shell-literal encoding protects the outer invocation; positional argv protects the inner managed-recipe invocation.

## Discovered Tasks

None. The stale staged-package assertions were directly required regression work and were resolved by the D-01 scope expansion rather than deferred.

## Implementation Summary

**Files changed:**
- `README.md` (modified)
- `justfile` (modified)
- `skills/do-work-board/justfile.template` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/publication_commands_test.go` (modified)
- `skills/do-work-toolbox/actions/note.md` (modified)
- `skills/do-work-toolbox/actions/architecture-report.md` (modified)
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified)
- `skills/do-work-toolbox/actions/present-work.md` (modified)
- `skills/do-work-toolbox/actions/install.md` (modified)
- `skills/do-work-toolbox/actions/maintainability-audit.md` (modified)
- `skills/do-work-toolbox/actions/maintainability-audit-reference.md` (modified)
- `skills/do-work/actions/version.md` (modified)
- `skills/do-work/actions/help.md` (modified)
- `skills/do-work-knowledge/actions/help.md` (modified)
- `skills/do-work-toolbox/actions/help.md` (modified)
- `skills/do-work/docs/version-guide.md` (modified)
- `skills/do-work/docs/command-line-guide.md` (new)
- `skills/do-work/docs/prescribed-shell-primitives.md` (modified)
- `skills/do-work/tools/prime-do-work-update.md` (modified)
- `skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md` (modified)
- `skills/do-work-board/actions/board.md` (modified)
- `skills/do-work-board/docs/board-guide.md` (modified)
- `_dev/tests/flat-just-recipes-behavior.sh` (new)
- `_dev/tests/contract-regressions.sh` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)

**What was done:** Published a 40-definition flat Just interface backed by one managed template, made installer collision/completeness validation derive from that template, and delegated the named deterministic action phases to `do-work-cli` with fail-closed behavior. Review remediation replaced raw `{{args}}` shell-source interpolation in all 33 non-board variadic recipes per surface with Just positional argv and quoted `"$@"`, made recovery recipes include the canonical manifest/time flags, moved hostile tests onto the actual shipped recipe, and aligned action/reference/prime/board restatements with singular CLI ownership.

## Qualification

Passed — all 30 declared project files were inspected; scope-drift reports exact Scope/Implementation Summary parity; the actual advertised recipe is hostile-byte safe; canonical ownership restatements are aligned; P-A-U is complete; and no debug artifacts, hidden Unicode controls, hollow tests, or undeclared project touches were found.

## Testing

**Tests run:** `go test ./internal/suiteinstall ./internal/publication -count=1`; `go test -race ./internal/suiteinstall ./internal/publication -count=1`; `go vet ./internal/suiteinstall ./internal/publication`; `_dev/tests/flat-just-recipes-behavior.sh`; scope-drift; staged-skills, update, install, and aggregate contract suites; shell syntax/ShellCheck/gofmt checks; and two successful full `bash _dev/tests/maintainer-verify.sh` final-remediation runs (the second after the last prose restatement edits).

**Result:** ✓ All passing. The actual shipped recipe forwarded NUL-exact hostile argv without sentinel side effects. The final-state canonical maintainer gate exited 0 after aggregate contracts, queue-kanban vet/uncached tests/strict JavaScript, do-work-cli vet/uncached tests, and audit-metrics vet/tests. Its strict browser lane was explicitly skipped because no browser was available.

**Red-green validation:**
- Template-derived candidate completeness: failed when a managed definition was omitted → passed after completeness derived from template bytes.
- Reserved collision coverage: prior subset missed expanded names → passed for all 40 dynamically parsed definitions with unchanged target bytes.
- Publication recovery argv initial RED used an invented safe recipe and missed the shipped seam. Review RED against the actual template executed the `$(touch …)` sentinel and reparsed an embedded newline as a command; GREEN now passes NUL-exact spaces, quotes, dollars, command substitutions, backticks, tabs, and newlines through the shipped recipe with no sentinel files while `next_argv` stays exact.
- Flat interface: 21 managed definitions found → passed with 40 definitions invoked through real Just.
- Action delegation: seven deterministic phases lacked the full canonical fail-closed contract → passed with canonical command, typed-result, missing/failed/malformed, and no-fallback assertions.

**New tests added:**
- `_dev/tests/flat-just-recipes-behavior.sh`
- Dynamic installer completeness/collision cases in `install_transaction_test.go`
- Hostile shell/Just argument round trip in `publication_commands_test.go`
- Canonical delegation and expanded managed-interface contracts across the shell suites

## Lessons Learned

**What worked:** Reproducing the failure through the actual managed template with command-substitution sentinels exposed the real shell boundary; Just's per-recipe `[positional-arguments]` plus `"$@"` preserved every original argument byte.

**What didn't:** POSIX-quoting the outer generated command was insufficient while inner recipes interpolated `{{args}}` as shell source. The first test used an invented safer recipe and therefore proved the wrong seam.

**Worth knowing:** Publication recipes must preserve two boundaries together: shell-literal outer arguments and positional inner Just arguments. Named leading parameters need capture-and-shift before canonical CLI flags are rebuilt.

*Verified by work action*

## Orientation

[MAP CHANGED] Deterministic core, knowledge, toolbox, and board operations now share one flat managed Just interface, with action prose retaining judgment while delegating mutation to the canonical command platform. The action-files, shell-command, and Kanban-board maps remain structurally current; their referenced paths all resolve.

## Review — Initial

**Overall: 50%** | 2026-09-01T14:09:49Z

| Dimension | Score |
|-----------|-------|
| Requirements | 63% |
| Code Quality | 55% |
| Test Adequacy | 50% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Fail |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Shipped variadic Just recipes interpolate `{{args}}` into shell, allowing command substitution and corrupting hostile argument bytes; generated publication recovery recipes therefore fail their real execution boundary — impact-critical → remediation required
- Canonical CLI ownership and fail-closed semantics remain contradicted by live action, reference, prime, prescribed-primitive, and board-guide restatements — impact-rule-change → remediation required
- Hostile-argv and flat-recipe tests exercise a safer invented Just seam or never assert captured argv, so the required caller-seam RED→GREEN proof is absent — impact-rule-change → remediation required

**Minor findings:** 1 (legacy updater usage wording in an internal error path; report only)
**Acceptance:** Fail — normal suites pass, but the real advertised recipe executes hostile shell substitutions and does not preserve argv bytes
**Suggested testing:** 3 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Review

**Overall: 50%** | 2026-09-01T15:02:03Z

| Dimension | Score |
|-----------|-------|
| Requirements | 94% |
| Code Quality | 94% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- The retired board-only compiler/fallback contract remains active in `_dev/primes/prime-kanban-board.md:16`, `skills/do-work/actions/work-reference.md:105,113`, `skills/do-work-board/tools/queue-kanban/timestamp.go:19-22`, and `frontmatter_cli.go:24-32`; these instructions still say only the board may require Go and core needs shell fallbacks, contradicting the canonical build-on-demand `do-work-cli` and leaving the prior restatement finding unclosed — impact-rule-change → folded into REQ-420

**Minor findings:** 1 (the internal empty-root error at `update_transaction.go:169` still advertises legacy `do-work-update.sh --project-root` usage; report only)
**Acceptance:** Fail — actual-template hostile argv preservation, exact `next_argv`, focused Go ordinary/race/vet, staged contracts, scope parity, and the canonical maintainer gate pass, but the finding-closure ratchet fails on the remaining live toolchain restatements
**Suggested testing:** 2 items
**Follow-ups created:** None; **sweeps appended to:** REQ-420

*Reviewed by review-work action*

## Remediation

The single allowed remediation closed the executable shell-injection/argv-corruption defect, replaced the invented test seam with the shipped template, and corrected the scoped ownership restatements. Fresh review confirmed those fixes and every acceptance gate, but found four additional live restatements outside the expanded 30-file implementation scope. Because the remediation allowance is exhausted, this REQ archives as `completed-with-issues`; the remaining rule-change is folded into REQ-420's whole-suite shim/parity sweep.
