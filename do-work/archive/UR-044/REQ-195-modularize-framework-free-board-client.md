---
id: REQ-195
title: Modularize the framework-free queue board client
status: completed
created_at: 2026-08-15T09:12:23Z
claimed_at: 2026-08-15T13:46:00Z
completed_at: 2026-08-15T14:24:38Z
user_request: UR-044
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: refactor
route: C
kb_status: pending
kb_entry:
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-185]
write_set:
  - skills/do-work-board/tools/queue-kanban/dependency_test.go
  - skills/do-work-board/tools/queue-kanban/filementions.go
  - skills/do-work-board/tools/queue-kanban/generate.go
  - skills/do-work-board/tools/queue-kanban/generate_test.go
  - skills/do-work-board/tools/queue-kanban/model.go
  - skills/do-work-board/tools/queue-kanban/model_test.go
  - skills/do-work-board/tools/queue-kanban/notes_test.go
  - skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
  - skills/do-work-board/tools/queue-kanban/serve.go
  - skills/do-work-board/tools/queue-kanban/serve_test.go
  - skills/do-work-board/tools/queue-kanban/web/board-calendar.js
  - skills/do-work-board/tools/queue-kanban/web/board-cards.js
  - skills/do-work-board/tools/queue-kanban/web/board-clipboard.js
  - skills/do-work-board/tools/queue-kanban/web/board-controls.js
  - skills/do-work-board/tools/queue-kanban/web/board-core.js
  - skills/do-work-board/tools/queue-kanban/web/board-detail.js
  - skills/do-work-board/tools/queue-kanban/web/board-filters.js
  - skills/do-work-board/tools/queue-kanban/web/board-testing.js
  - skills/do-work-board/tools/queue-kanban/web/board.css
  - skills/do-work-board/tools/queue-kanban/web/board.js
  - skills/do-work-board/tools/queue-kanban/web/template.html
  - _dev/tests/contract-regressions.sh
---

# Modularize the Framework-Free Queue Board Client

## What

Split the framework-free queue-kanban browser client from one approximately 2,524-line `web/board.js` source unit into a private shell and eight ordered closure fragments. Preserve the existing browser runtime and every static/live behavior while making source ownership and review boundaries smaller and explicit.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`, its linked REQ-185 and REQ-194 lessons, and the always-on general/coding guardrails. Preserve the current 105,549-byte browser runtime by cutting eight contiguous raw closure fragments, assembling them through one private shell seam in an explicit order, and pinning authored-source inventory, order, boundaries, syntax, and static/live byte parity before retargeting factual ownership comments.
- [x] **[APPLY]:** Added the structural RED first, cut the current browser client into the private shell and eight raw closure fragments, wired the fixed-order assembler into the shared static/live HTML path, and added the planned source-map, ownership, contract, syntax, and byte-parity coverage within the exact 22-file implementation scope.
- [x] **[UNIFY]:** Reviewed every scoped diff and the final raw-marker remediation; verified shell/fragment cut boundaries, first-assembly byte parity, explicit execution order, one classic inline runtime, no fragment routes, strict Node-probe accounting, factual comment ownership, and contract retargeting. Focused/full Go tests, vet, contract regressions, static/live browser journeys, gofmt, and `git diff --check` pass without debug or temporary migration artifacts.

## Why

The current client combines core helpers, filters, cards, calendar, testing, detail, controls, clipboard, and bootstrap behavior in one low-churn but high-coupling source unit. The consumer implementation demonstrated that these responsibilities can be separated without changing the runtime, reducing navigation, ownership, and review cost.

## Finding Provenance

- Source: consumer maintainability audit and `/Users/t2/Desktop/llm-tinker2/do-work/archive/UR-045/REQ-197-modularize-framework-free-board-client.md`
- Source commit: `d9237469478bd65e3574f2e80d7b57aac9148dfe`
- Original severity: P2, verified, low churn
- Upstream validation verdict: Accept
- Validated upstream baseline at capture: `web/board.js` is 2,524 lines and 105,430 bytes; SHA-256 `9bda99f4f71a7f1a81f70d1db629598c59d31b39602cc97307067e43ff12aa74`

## Detailed Requirements

- Retain `web/board.js` as the private IIFE/bootstrap shell.
- Divide existing implementation into these ordered classic-script closure fragments:
  1. `board-core.js`
  2. `board-filters.js`
  3. `board-cards.js`
  4. `board-calendar.js`
  5. `board-testing.js`
  6. `board-detail.js`
  7. `board-controls.js`
  8. `board-clipboard.js`
- Keep fragments private to the existing closure. They must not become ES modules or expose new browser globals.
- Make `generate.go` own one explicit ordered fragment manifest. A wildcard embed may include authored JavaScript files but must never determine execution order.
- Assemble through exactly one internal placeholder, with the shell and every fragment included exactly once.
- Preserve one classic-script runtime with no additional browser request, bundler, framework, module loader, or network dependency.
- Keep static generation and the live server on the same assembled client bytes.
- Establish a fresh upstream pre-change byte count and SHA-256 before splitting. Prove the first assembled result is byte-identical to that baseline, treating the hash as one-time migration evidence rather than a permanent golden test.
- Retarget factual source-owner comments to their exact fragment or to the assembled client. Update `prime-do-kanban.md` with the ordered source map and describe one assembled client instead of one hand-authored source.
- Retarget the upstream-only `_dev/tests/contract-regressions.sh` write-set tooltip checks from raw `web/board.js` to the owning `web/board-cards.js` fragment so they cannot become false-green.
- Regenerate the next upstream suite version and changelog entry from the integration-time upstream version; do not copy consumer release history.

## Required Characterization and Regression Tests

- RED must begin with structural tests that fail against the current monolith: the embedded authored-JavaScript set does not yet equal one shell plus the eight-file manifest, and there is no fragment assembler seam to validate.
- Assert the embedded authored-JavaScript set equals the shell plus manifest with no omissions or duplicates.
- Assert assembly follows manifest order and every fragment occurs exactly once.
- Assert the shell has exactly one internal fragment placeholder and assembled output retains none.
- Assert assembler-owned fragment boundaries/newlines are deterministic.
- Run `node --check` against the assembled client when Node is available, composing with REQ-185's maintainer-strict JavaScript lane rather than weakening or duplicating it.
- Prove static generation and live serving contain the identical assembled client.
- Preserve existing Go and Node behavior probes, including lazy Markdown/Copy, drawer behavior, filters, Testing state, persistence, and generated-asset contracts.
- Run the complete queue-kanban `go test ./...` and `go vet ./...` suites.
- Perform one-time static `file://` and loopback-live browser characterization during migration. Do not copy the consumer's repository-specific browser test path as the permanent upstream harness.

## Constraints

- No intended user-visible behavior or private runtime API change.
- Do not copy `e2e/playwright-scripts/queue-kanban.spec.ts` into upstream.
- Do not copy `cmd/check-maintainability-limits.sh` or its application response, TSX, or backend ceilings.
- Do not add the consumer's `>2524` line ceiling; it permits the original monolith and does not replay the maintenance incident.
- Do not make the consumer hash a permanent golden assertion.
- Do not copy consumer REQ-197 as the upstream request identity.
- Do not add the consumer's internal-architecture paragraph to `actions/board.md`; the prime file owns implementation architecture.
- Do not touch `main.go` unless implementation actually moves shell state out of `board.js`; the accepted structure retains that state.
- Keep release/version files serially owned by the integrating work action.

## Surface-Cost

**Earned.** The replayed incident is the verified 2,524-line mixed-responsibility client. The lasting surface is eight cohesive source files, one explicit assembler manifest/seam, and structural tests. That cost is lower than retaining the monolith because it preserves one runtime while reducing review and ownership coupling; the tests directly cover the omission, duplication, reordering, seam, syntax, and static/live-divergence risks introduced by the split.

## Dependencies

No hard dependency. Pending REQ-185 overlaps `generate_test.go` and `web/board.js` but owns a separate maintainer-strict JavaScript behavior lane. Reconcile its final test helper/API and integrate serially if it lands first; do not merge the intents or weaken either request.

## Builder Guidance

Firm intent. Treat the consumer commit as external evidence, not a patch to apply wholesale. The eight fragment bodies are portable only after the builder records and verifies a fresh upstream baseline; all existing-file edits, tests, comments, and release integration must be adapted to current upstream.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Add the package-owned assembly characterization that requires one embedded shell, the eight ordered fragment paths, one exact internal seam, no missing or duplicate authored JavaScript, and identical static/live assembled client bytes; run it against the current monolithic upstream source.
**Why RED now:** Current `generate.go` embeds and reads only `web/board.js`; the eight fragments and ordered assembler contract do not exist, so the structural characterization cannot pass.
**GREEN when:** A freshly baselined shell-plus-fragment assembly is byte-identical to the pre-change upstream runtime, every structural assertion passes, `node --check`, queue-kanban Go tests/vet, existing behavior probes, and one-time static/live browser journeys pass with no new globals, modules, requests, or user-visible changes.
**Validation:** User confirmed by invoking `do-work capture-request` on the immediately preceding validated handoff.

## Assets

None.

## Full Context

See `do-work/user-requests/UR-044/input.md` for the verbatim invocation and validated handoff summary. The authoritative external evidence remains the named consumer commit and archived consumer REQ path above.

---
*Source: "do-work capture-request" — using the immediately preceding provenance-preserving validated handoff.*

## Triage

**Route C** — this is a high-reach, behavior-preserving source migration across an embedded browser runtime, static/live assembly, ordered fragments, strict JavaScript probes, package documentation, and upstream-only contracts; it requires an explicit cut map, parity plan, and independent exploration before Scope is finalized.

## Plan

1. Add a structural RED against the current monolith: the embedded authored-JavaScript inventory must equal the retained shell plus the prescribed eight-fragment manifest, and the shell must contain the one internal assembly marker.
2. Cut the current source byte-for-byte at the verified responsibility boundaries. Keep the IIFE state/guard and bootstrap in `board.js`; move only raw statements and declarations into the eight classic-script fragments, with separator blank lines owned by the assembler.
3. Add one injectable assembler seam in `generate.go`. Embed authored JavaScript with `web/*.js`, execute only the explicit manifest order, validate the single shell marker and canonical fragment endings, and have the existing shared static/live HTML path consume the assembled client.
4. Immediately compare the first assembled client with the fresh pre-change baseline (2,531 lines, 105,549 bytes, SHA-256 `f7211abfbf17e873146287f6790049b4e8861cf72207889ab4dd76cd459488a5`). Record that result as migration evidence, then remove the temporary hash assertion so no golden survives.
5. Add permanent mutation-resistant tests for exact authored-source inventory, literal manifest order, one-time fragment inclusion, marker removal, deterministic boundary bytes, static/live equality, and assembled-client syntax. Reuse REQ-185's Node lookup while keeping syntax outside its counted behavior-probe invariant.
6. Retarget only factual ownership comments, document the ordered source map in `prime-do-kanban.md`, and move the raw tooltip contract assertions to `board-cards.js`. Preserve fixture strings where `web/board.js` is merely example data.
7. Run focused structural/Node tests, full Go tests/vet, contract regressions, one-time static `file://` and live-loopback browser journeys, and the canonical maintainer gate; inspect every scoped file and exclude temporary artifacts.

## Exploration

- The current clean upstream source is LF-terminated at 2,531 lines and 105,549 bytes with SHA-256 `f7211abfbf17e873146287f6790049b4e8861cf72207889ab4dd76cd459488a5`. The captured consumer figures are stale because REQ-185 subsequently extracted two JavaScript helpers.
- Exact source cuts from the pre-change `board.js` are: prefix lines 1–45; core 46–311; filters 313–478; cards 480–1049; calendar 1051–1128; testing 1130–1583; detail 1585–2165; controls 2167–2318; clipboard 2320–2508; suffix 2510–2531. The eight omitted lines are single blank separators and become assembler-owned.
- Static generation and live `GET /` already converge through `assembleStaticPage`; the implementation needs no new route or browser request. The template retains one classic inline script after `board-data.js`, and `board-markdown.js` remains the existing Copy-only lazy request.
- Execution order must come only from the fixed manifest, never wildcard enumeration. The wildcard embed exists to make the complete authored-source inventory testable.
- The permanent Node syntax test belongs under the `TestJavaScriptBehavior` prefix so the strict lane runs it, but it must not increment `javaScriptBehaviorProbeCount`; otherwise syntax alone could mask deletion of every executable behavior probe.
- `main.go` remains truthful because `windowHours` stays in the shell. Historical changelog/report references and `model_test.go` fixture literals that use `web/board.js` as arbitrary write-set data are not source-ownership claims and remain unchanged.
- Independent planning and exploration both challenged the seeded write set and found all 22 files necessary and sufficient. No sibling or consumer repository was inspected, and no external patch will be copied.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/dependency_test.go` — factual source-owner comment update.
- `skills/do-work-board/tools/queue-kanban/filementions.go` — factual source-owner comment update.
- `skills/do-work-board/tools/queue-kanban/generate.go` — ordered embedded assembly and exact shell-marker validation.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` — inventory, order, boundary, marker, syntax, and behavior-lane tests.
- `skills/do-work-board/tools/queue-kanban/model.go` — factual source-owner comment update.
- `skills/do-work-board/tools/queue-kanban/model_test.go` — factual source-owner comment update.
- `skills/do-work-board/tools/queue-kanban/notes_test.go` — factual source-owner comment update.
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` — one assembled classic client and ordered source map.
- `skills/do-work-board/tools/queue-kanban/serve.go` — factual assembled-client ownership update.
- `skills/do-work-board/tools/queue-kanban/serve_test.go` — static/live assembled-client byte parity and request-shape coverage.
- `skills/do-work-board/tools/queue-kanban/web/board-calendar.js` — extracted calendar closure fragment.
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` — extracted cards, columns, recent-window, warning, and By-UR closure fragment.
- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` — extracted lazy Markdown/clipboard closure fragment.
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` — extracted view/filter control closure fragment.
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` — extracted core helper closure fragment.
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` — extracted drawer, linkification, and resize closure fragment.
- `skills/do-work-board/tools/queue-kanban/web/board-filters.js` — extracted filter predicate closure fragment.
- `skills/do-work-board/tools/queue-kanban/web/board-testing.js` — extracted Testing profile/status closure fragment.
- `skills/do-work-board/tools/queue-kanban/web/board.css` — factual source-owner comment update.
- `skills/do-work-board/tools/queue-kanban/web/board.js` — retained private IIFE/state/bootstrap shell and one canonical internal marker line.
- `skills/do-work-board/tools/queue-kanban/web/template.html` — factual assembled-client ownership update.
- `_dev/tests/contract-regressions.sh` — cards-owned write-set tooltip contract retarget.

**Explicit exclusions:** `main.go`, browser harnesses, maintainability ceilings, consumer paths, action architecture prose, release files during implementation, and every file outside the 22-path frontmatter `write_set`.

**Acceptance criteria:**

- [x] The first assembled client is byte-identical to the fresh pre-change baseline, with the result recorded as one-time evidence and no permanent hash golden.
- [x] The authored-JavaScript set is exactly one shell plus eight uniquely ordered fragments; each fragment occurs once, the shell marker occurs once before assembly and never afterward, and boundary bytes are deterministic.
- [x] Static and live pages contain identical assembled client bytes in one classic inline script, with no fragment requests, modules, new globals, bundler, framework, or dependency.
- [x] REQ-185's counted behavior lane still requires executable behavior probes; assembled syntax uses the same optional-Node policy without satisfying that count.
- [x] Existing Go/JavaScript behavior, including filters, cards, calendar, Testing, detail, resizing, By-UR/recent-window behavior, persistence, and lazy Copy, remains green.
- [x] Static and live browser characterization completes without console/reference errors or user-visible behavior differences, and only the existing asset requests occur.
- [x] Factual ownership comments, the prime source map, and the cards tooltip contract identify their exact post-split owners without rewriting sample or historical data.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/dependency_test.go` (modified) — retargets generated-annotation ownership to the cards fragment
- `skills/do-work-board/tools/queue-kanban/filementions.go` (modified) — retargets body-mention ownership to the detail fragment
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified) — embeds authored JavaScript inventory and assembles the exact manifest through one validated shell marker
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — locks inventory, order, boundaries, raw/canonical markers, syntax, and strict-lane composition
- `skills/do-work-board/tools/queue-kanban/model.go` (modified) — updates source-owner guidance for shell state and fragment gates
- `skills/do-work-board/tools/queue-kanban/model_test.go` (modified) — updates factual renderer-owner references while retaining sample write-set literals
- `skills/do-work-board/tools/queue-kanban/notes_test.go` (modified) — retargets note rendering ownership to the cards fragment
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified) — documents the private shell, fixed fragment order, and one assembled classic client
- `skills/do-work-board/tools/queue-kanban/serve.go` (modified) — identifies the clipboard fragment as the lazy Markdown owner
- `skills/do-work-board/tools/queue-kanban/serve_test.go` (modified) — proves generated and live pages contain identical assembled bytes with no fragment/module requests
- `skills/do-work-board/tools/queue-kanban/web/board-calendar.js` (new) — owns calendar rendering inside the retained private closure
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (new) — owns cards, columns, recent-window, warnings, notes, and By-UR rendering
- `skills/do-work-board/tools/queue-kanban/web/board-clipboard.js` (new) — owns lazy Markdown loading and clipboard behavior
- `skills/do-work-board/tools/queue-kanban/web/board-controls.js` (new) — owns view switching, filters, search, and window controls
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` (new) — owns shared escaping, time, path, and display helpers
- `skills/do-work-board/tools/queue-kanban/web/board-detail.js` (new) — owns linkification, detail drawer behavior, and persisted resizing
- `skills/do-work-board/tools/queue-kanban/web/board-filters.js` (new) — owns request and User Request filtering predicates
- `skills/do-work-board/tools/queue-kanban/web/board-testing.js` (new) — owns testing profiles, status transitions, and feedback UI
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified) — retargets renderer-owner comments to exact fragments
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modified) — retains the IIFE, shared state, bootstrap, and single canonical fragment marker line
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified) — retargets UI-owner comments while retaining one inline classic-script slot
- `_dev/tests/contract-regressions.sh` (modified) — moves cards-owned tooltip assertions to `web/board-cards.js`

**Behavior:** The browser still executes one private classic IIFE with byte-identical pre-migration behavior, but its source is now maintained as a small shell plus eight cohesive fragments assembled only by the fixed Go manifest. Static generation and live serving share those exact assembled bytes; fragments add no routes, requests, globals, modules, frameworks, or dependencies.

## Testing

**RED:** `TestEmbeddedAuthoredJavaScriptInventory` failed against the monolith because the embedded authored set contained only `web/board.js`, with none of the required fragment files or assembly seam.

**Migration evidence:** The first production-assembled client matched the fresh pre-change source exactly: 105,549 bytes, 2,531 LF-terminated lines, SHA-256 `f7211abfbf17e873146287f6790049b4e8861cf72207889ab4dd76cd459488a5`. The temporary comparison was removed; no golden hash or baseline copy remains.

**GREEN:**
- focused inventory, assembly, invalid-structure, Node syntax, and static/live tests — PASS
- raw-marker remediation mutations (duplicate EOF marker, noncanonical line, and fragment-retained marker) — correctly rejected
- maintainer-strict JavaScript behavior lane and zero-probe rejection self-test — PASS
- `go test -count=1 ./...` — PASS
- `go vet ./...` — PASS
- `bash _dev/tests/contract-regressions.sh` — PASS
- `gofmt`, `bash -n`, cut/reassembly comparisons, scope audit, and `git diff --check` — PASS

**Browser characterization:** A disposable 194-REQ/44-UR fixture passed static `file://` and loopback-live journeys for board/search, By-UR and 7-day lenses, calendar, Testing, detail/linkification, Copy, and persisted drawer resizing. Static Testing remained read-only. Live Testing created a disposable profile, transitioned REQ-194 to in-testing and returned-with-feedback through HTTP 200 writes, and persisted the expected temporary frontmatter. Both surfaces used one inline classic client; no fragment request occurred, Markdown loaded exactly once on first Copy, and no JavaScript/reference warning appeared. The live server's sole console error was an unrelated missing favicon. The disposable fixture was moved to Trash after verification.

## Qualification

- **Scope:** PASS — the 22-path Implementation Summary exactly matches the frozen Scope; orchestration state, release integration, and foreign REQ-189–192 edits remain separate.
- **Mechanical checks:** PASS — all implementation paths exist in the diff, P-A-U is complete, and no debug or temporary migration artifacts remain.
- **Substance and traceability:** PASS — every requirement maps to the fixed manifest/assembler, structural and behavior tests, exact first-assembly evidence, shared static/live caller seam, browser characterization, or factual owner retarget.
- **Wiring/data flow:** PASS — wildcard embed exposes inventory only; the literal manifest owns order; one canonical shell line is replaced; the existing shared HTML assembler supplies identical bytes to static and live callers.

## Review

**Overall: 99%** | 2026-08-15T14:24:38Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 99% |
| Test Adequacy | 99% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** The initial independent review found that marker uniqueness counted only the canonical marker-plus-LF shape, allowing an additional raw marker at EOF. Remediation separated raw-token uniqueness from canonical-line validation, rejects any retained marker after assembly, and added mutation cases. The reviewer replayed the EOF, noncanonical-line, and fragment-retained-marker mutations and confirmed the finding closed.
**Minor findings:** None.
**Acceptance:** Pass — the runtime split is byte-preserving, behaviorally characterized, structurally mutation-resistant, and confined to the exact source surface.
**Follow-ups created:** None; **sweeps appended to:** None.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Recording the fresh source hash before cutting and comparing the first production assembly before changing factual comments isolated the migration itself from later non-runtime edits. A wildcard embed makes authored inventory observable while a separate literal manifest keeps execution order reviewable and deterministic.

**What didn't:** Counting only the full placeholder line was not equivalent to proving one raw marker token; a second marker without the canonical newline survived both replacement and the first tests. Raw-token uniqueness, canonical placement, and post-assembly absence are three separate invariants and need separate assertions.

**Worth knowing:** Keep fragment files as raw statements inside the shell's existing IIFE. Separator blank lines belong to the assembler, so fragment endings and manifest joins are part of the byte contract even though the browser would tolerate many equivalent layouts.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Orientation

**[MAP CHANGED]** The queue-kanban browser runtime is still one private, framework-free classic client, but its source map is now a retained `board.js` shell plus eight ordered responsibility fragments assembled by `generate.go`. Static generation and live serving consume the same assembled bytes, making source ownership smaller without adding browser-visible loading or API seams.
