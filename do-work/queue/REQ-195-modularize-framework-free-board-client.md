---
id: REQ-195
title: Modularize the framework-free queue board client
status: pending
created_at: 2026-08-15T09:12:23Z
user_request: UR-044
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: refactor
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
