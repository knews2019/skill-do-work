---
id: REQ-183
title: Static board generation can publish a mixed three-file bundle
status: completed
created_at: 2026-08-15T07:13:20Z
claimed_at: 2026-08-15T10:14:36Z
completed_at: 2026-08-15T10:27:19Z
commit: 803e4e7
user_request: UR-041
domain: backend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-181, REQ-182, REQ-184, REQ-185, REQ-186, REQ-187, REQ-188]
batch: audit-findings-2026-08-14
write_set: [skills/do-work-board/tools/queue-kanban/generate.go, skills/do-work-board/tools/queue-kanban/generate_test.go]
route: A
kb_status: skipped
kb_entry:
---

# Static Board Generation Can Publish a Mixed Three-File Bundle

## What

Make static board publication all-or-recover across `board-data.js`, `board-markdown.js`, and `index.html`, so a handled write failure cannot expose files from different generations.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route A direct bug fix: add a seeded second-publication failure regression first, then stage all three private payloads and publish them through one bounded backup/restore helper that preserves pre-existing bytes, unrelated output entries, and cleanup invariants.
- [x] **[APPLY]:** Added the deterministic second-publication failure test first, then replaced direct writes with one three-output private stage, unique backup, sequential publish, and rollback helper inside `generate.go`.
- [x] **[UNIFY]:** Reviewed the complete two-file diff and similar direct-writer scan. Verified first generation, successful replacement, byte-identical rollback of all three targets, unrelated-entry preservation, private cleanup, full `go test ./...`, and `git diff --check`; no debug artifacts remain.

## Why

`generateStaticSite` writes the public files directly in sequence. A forced second-write failure exits nonzero after publishing new rendered data beside the old copy payload and shell, leaving a plausible but internally inconsistent shareable artifact.

## Context

- Audit priority: P2; impact 3; effort normal.
- Root-cause key: `static-board-all-or-recover-publication`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 3.
- Reproduce: `cd skills/do-work-board/tools/queue-kanban && probe_root="$(mktemp -d /tmp/queue-kanban-audit.XXXXXX)" && fixture_repo="$probe_root/repo" && probe_out="$probe_root/out" && mkdir -p "$fixture_repo/do-work/queue" && printf '%s\n' '---' 'id: REQ-001' 'title: Old title' 'status: pending' '---' '' '# Old body' > "$fixture_repo/do-work/queue/REQ-001-fixture.md" && go run . generate --repo-root "$fixture_repo" --out "$probe_out" && printf '%s\n' '---' 'id: REQ-001' 'title: New title' 'status: pending' '---' '' '# New body' > "$fixture_repo/do-work/queue/REQ-001-fixture.md" && chmod 0444 "$probe_out/board-markdown.js" && go run . generate --repo-root "$fixture_repo" --out "$probe_out"; status=$?; printf 'status=%s\n' "$status"; grep -F 'New title' "$probe_out/board-data.js"; grep -F 'Old body' "$probe_out/board-markdown.js"`

## Detailed Requirements

- Build all three payloads privately before publishing any public target.
- Publish only `board-data.js`, `board-markdown.js`, and `index.html` as one bounded all-or-recover operation.
- Use unique backups/restoration so a handled publication failure restores all pre-invocation target bytes.
- Preserve unrelated entries already present under `--out`.
- Cover seeded failure, first generation, successful replacement, rollback, and private-residue cleanup.
- Preserve the existing static, zero-network board output contract.

## Constraints

- Define no stronger promise than all-or-recover for these three files; do not claim one filesystem-atomic rename across them.
- Do not introduce a general transaction framework.
- Lock-in limit: after a handled failure, zero public target bytes differ from their pre-invocation state and no private residue remains.

## Dependencies

None. The `generate_test.go` overlap with REQ-185 is declared for coordination, not a semantic prerequisite.

## Builder Guidance

Firm intent. The audit says one bounded three-file staging/backup helper is earned by the reproduced failure and `_dev/lessons/validated-runtime-boundaries.md`.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Generate an initial fixture, change the ticket, make the second target unwritable, and generate again.
**Why RED now:** The command exits 1 after `board-data.js` contains the new title while `board-markdown.js` still contains the old body.
**GREEN when:** The seeded second-write failure returns nonzero while all three public targets remain byte-identical to their pre-invocation state, unrelated output entries remain, and no staging or backup residue survives; first-generation and success paths also pass.
**Validation:** Confirmed by the user during verification on 2026-08-15.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 03, labeled P2, impact 3, normal effort.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 3 in the canonical audit for the complete batch constraints and validated evidence record.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*

## Triage

**Route A** — the request identifies the exact two files, supplies a deterministic reproduction, and defines the required all-or-recover behavior without an architectural choice outside the focused REQ text.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified) — builds all payload bytes before mutation, stages the bounded bundle privately, backs up existing targets uniquely, and restores pre-invocation targets on handled publication failure
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — covers first generation, successful replacement, deterministic second-publication rollback, unrelated entries, exact target bytes, and residue cleanup

**Behavior:** The three static outputs are not claimed to publish atomically, but any handled staged-publication failure returns the public target set to its pre-invocation state. Successful runs still produce a self-contained zero-network board.

## Testing

**RED:** The seeded second-publication regression failed against direct writes: the injected publisher received zero calls, generation returned nil, and all three public files differed from their snapshots.

**GREEN:**
- Focused static-publication tests — PASS
- `cd skills/do-work-board/tools/queue-kanban && go test ./...` — PASS (5.727s)
- `git diff --check` — PASS
- Similar-pattern scan — no remaining direct static-target writer

## Qualification

- **Boundary:** PASS — Route A stayed inside the two files named by the focused REQ text; unrelated queue edits remain excluded.
- **Mechanical checks:** PASS — `qualify.sh` found both modified files in the diff, all P-A-U phases complete, and no debug artifacts.
- **Substance and traceability:** PASS — the helper privately stages exactly three outputs, uniquely relocates existing targets before any publication, removes partial replacements on failure, restores every backup, and removes private state after success or a completed rollback.
- **Wiring/data flow:** PASS — `generateStaticSite` always calls the bounded publisher with `os.Rename`; the test-only publisher seam injects a deterministic second rename failure without changing the real CLI path.

## Review

**Result:** Approve — Acceptance: Pass  
**Overall score:** 100%

- All three payloads are assembled and staged before any public target changes.
- Existing targets move into one uniquely named private directory; a failed publication removes partial replacements and restores every prior target byte-for-byte.
- Cleanup errors join the returned error. If restoration itself fails, the private directory remains so recovery evidence/backups are not destroyed; the operation does not falsely claim recovery.
- First generation, successful replacement, deterministic second-publication rollback, unrelated output, cleanup, zero-network output, and non-atomicity wording all pass acceptance.

**Important findings:** None.  
**Minor findings:** None.  
**Explicit remediation:** None.

**Independent checks:** `go test -count=1 ./...` PASS (5.370s), `go vet ./...` PASS, `git diff --check` PASS.

## Lessons Learned

- The existing validated-runtime-boundary lesson already captured the reusable rule: complete private staging plus held backups and restoration, with no cross-file atomicity overclaim. This REQ adds regression coverage for the board but no new knowledge-base rule.

**Knowledge handoff:** Skipped; the applicable lesson already exists in `_dev/lessons/validated-runtime-boundaries.md`.
