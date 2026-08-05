---
id: REQ-112
title: Give frontmatter.go a CLI surface so prose can stop reimplementing it
status: completed
created_at: 2026-08-05T15:53:39Z
claimed_at: 2026-08-05T19:20:17Z
completed_at: 2026-08-05T19:24:57Z
route: A
user_request: UR-021
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: true
depends_on: [REQ-111]
write_set: [tools/queue-kanban/main.go, tools/queue-kanban/frontmatter_cli.go, tools/queue-kanban/frontmatter_cli_test.go]
maintenance: false
related: [REQ-111]
batch: census-durable-findings
---

# Give `frontmatter.go` a CLI Surface So Prose Can Stop Reimplementing It

## What

Add a `queue-kanban frontmatter` subcommand that reads one field from one REQ/UR file, optionally normalizing it per the Schema Read Contract and optionally testing set membership. Today `main.go` L60–76 exposes exactly seven subcommands — `summary | generate | serve | next-req | next-version | verify | now` — and **none takes a file-and-field argument**, so `splitFrontmatter` / `parseFrontmatterFields` / `lenientFrontmatterFields` (`frontmatter.go` L28, L82, L118) are unreachable from any action file. Every prose frontmatter read is therefore a hand reimplementation *by construction*, not by oversight.

Proposed surface (the exact flag names are the builder's to settle):

```
queue-kanban frontmatter get <file> <field> [--normalize] [--in-set terminal-success|terminal-resolved]
```

- `get` prints the raw value and exits 0; a missing field exits non-zero with nothing on stdout.
- `--normalize` applies the Schema Read Contract alias map and emits the contract's warning to **stderr** on an unrecognized value, so stdout stays a clean single value a caller can capture.
- `--in-set` exits 0/1 for membership in the Terminal-success or Terminal-resolved set (`actions/work-reference.md` L216–228), printing nothing. This is the check ~35 prose sites perform by hand.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** New `frontmatter_cli.go` holding a `get` verb over the existing parser, registered in `main.go`'s dispatch. Inject the two writers so the command is testable without a subprocess. Reuse `resolveSchemaField` from REQ-111 for `--normalize`, and delegate `--in-set` to the board's own `isCompletedStatus`/`isTerminalResolvedStatus` predicates so this cannot become a second definition of "finished". Verify: `go test ./...`, then a real binary smoke test — unit tests with injected writers do not prove the dispatch wiring.
- [x] **[APPLY]:** Added `tools/queue-kanban/frontmatter_cli.go` + `frontmatter_cli_test.go`, one dispatch case in `main.go`, and `schemaFieldWarningText` in `model.go` (deferred out of REQ-111 precisely because it had no consumer then; this REQ is the consumer). No action prose touched, per D-01.
- [x] **[UNIFY]:** `git diff --stat` → 4 files. `gofmt -l .` clean; `go vet ./...` clean; `go test ./...` passes. Grepped added lines for debug artifacts (`console.log`, `debugger`, `TODO`, `FIXME`, stray `fmt.Print` outside the writer-injected paths): none — every `fmt.Fprintf` targets an injected writer. Built the binary and exercised all six paths end-to-end: `get`, absent field, `--in-set` hit, `--normalize` alias, `--normalize` warning, and the `V=$(...)` capture idiom. Confirmed the shell floor returns the same answer via `awk` on the same fixture. `_dev/tests/contract-regressions.sh`: 7 failures, unchanged from the `main` baseline.

## Why (if provided)

Three frontmatter parsers already ship — the Go one, plus awk implementations inside `tools/checks/record-commit-hash.sh` (`frontmatter_line_for` L108–121) and `tools/checks/blanked-req-scan.sh` (`has_parseable_frontmatter` L88). Every prose read is a fourth-and-onward copy. The `status` vocabulary alone is read at ~35 prose sites, and five separate Red Flags already document the same resulting bug — filtering on the literal `completed` and silently dropping `completed-with-issues` (`actions/cleanup.md` L306, `actions/commit.md` L224, `actions/review-work.md` L479, `actions/ai-report.md` L341, `actions/present-work.md` L527). Five documented instances of one bug class is the evidence that the contract is fine and its enforcement is 35 hand copies.

## Detailed Requirements

- New subcommand registered in `main.go`'s dispatch switch, rejecting leftover tokens like every existing subcommand does (`exitOnLeftoverArguments` — silently discarding an argument is how `next-version` shipped bumping the wrong repo).
- Reuse `frontmatter.go` as-is. Do **not** fork or reimplement the parser; its CRLF handling, duplicate-top-level-key recovery (L70–81), and lenient block-list recovery (L109–117) are exactly the behaviours prose cannot replicate and are the reason this REQ exists.
- `--normalize` delegates to the normalizers, including the seven REQ-111 adds. Do not add a second normalization implementation here.
- Read-only. This must not become an eighth write surface: `CLAUDE.md` → Shipped Tooling states the tool has exactly two write surfaces and that adding a third means amending that sentence in the same commit. This REQ adds none, so that sentence stays untouched — and a `frontmatter set` verb is explicitly **out of scope**.

## Constraints

- **The compiled-tooling exception is the hard constraint.** `actions/board.md` is the only capability allowed to *need* a compiler (ADR-016; `CLAUDE.md` → Shipped Tooling → "Toolchain exception to design for the floor"). This subcommand is therefore permitted only in the **accelerator** shape that `next-req`, `next-version`, and `now` already use: named as the *preferred* source for something an action already obtains a shell-portable way, **gated on the binary already being built**, with the prose procedure documented as the floor and never a `go build` to obtain the value. An action that would compile the tool, or that has no floor path, is the prohibited shape.
- **This REQ ships the surface only — it rewrites no action prose.** Migrating any of the ~95 prose read sites to call it is separate, per-action work with its own review, and doing it here would turn a bounded tooling change into a sweep across most of `actions/`.

## Dependencies

`depends_on: [REQ-111]` — `--normalize` would silently no-op on seven of the nine fields until REQ-111's normalizers exist, which would ship a flag that lies about what it does.

## Builder Guidance

**Mixed.** Firm on: reusing `frontmatter.go`, read-only, the accelerator gating, and no prose migration. Exploratory on: flag naming, whether `--in-set` belongs on `get` or its own verb, and the exit-code convention for a missing field versus an unparseable file. Prefer the smallest surface that lets a prose step replace a hand-rolled awk read with one call.

## Open Questions

- [~] Should the first consumer land in the same REQ — one action's read site migrated to prove the surface works end-to-end — or should the subcommand ship with tests only? → **D-01**: Builder chose: tests only, no prose migration. Reasoning: migrating a read site means editing an action file, and every action file is a shell-floor document — the migration's real question is not "does the command work" but "what does this action do when the binary isn't built," which is a per-action prose judgment needing that action's own review. Shipping the surface with tests keeps this diff inside `tools/queue-kanban/` where the Go suite can prove it. Value: the CLI becomes available for any action to adopt deliberately, and this REQ stays reviewable as a bounded tooling change. Risk: a surface with no in-repo consumer can drift from what callers actually need, and the census's own top finding stays unfixed until someone migrates a site — reversible, since adopting it later costs one prose edit per action and nothing has to be undone.

<!-- D-XX counter: last used D-01. Next decision: D-02. -->

## Red-Green Proof

**RED prompt/case:** Run `queue-kanban frontmatter get do-work/archive/UR-020/REQ-110-census-completeness-floor-note.md status`. It fails today with `unknown subcommand "frontmatter"` and exit 2, from `main.go`'s dispatch default.
**Why RED now:** The seven registered subcommands are the only entry points, and none accepts a file-and-field pair, so the shipped parser has no caller outside the board's own walk.
**GREEN when:** That command prints `completed` and exits 0; `--normalize` on a `domain: back-end` fixture prints `backend` with the contract warning on stderr; and `--in-set terminal-success` exits 0 for `completed-with-issues` and 1 for `failed`.
**Validation:** Inferred during capture — the surface shape is the census's candidate #1 proposal, not a user-specified API.

## Full Context

See `do-work/user-requests/UR-021/input.md` for the complete verbatim input and batch constraints.

---
*Source: census finding — `frontmatter.go` has no CLI surface, so all ~95 prose frontmatter reads are hand reimplementations by construction (`decisions/audits/2026-08-05-shell-logic-in-prose-census.md` §1 structural fact 1, §4 candidate 1)*

---

## Triage

**Route: A** - Simple

**Reasoning:** The surface was specified in the REQ, the parser to reuse exists, and the subcommand conventions to match are in the same package. Scope is one new file plus one dispatch case.

## Implementation Summary

**What was done:** Added a `queue-kanban frontmatter get <file> <field>` subcommand so the shipped parser is reachable from prose for the first time, with `--normalize` for the Schema Read Contract and `--in-set` for the terminal-status membership test.

**Files changed:**
- `tools/queue-kanban/frontmatter_cli.go` (new) — `readFrontmatterField`, `parseFrontmatterCommandArguments`, `runFrontmatterCommand`, and the `--in-set` set map
- `tools/queue-kanban/frontmatter_cli_test.go` (new) — 9 tests covering raw read, absent field, no-frontmatter error, duplicate-key lenient recovery, stdout/stderr separation, the warning path, all 9 `--in-set` cases, and 6 usage-error shapes
- `tools/queue-kanban/main.go` (modified) — one dispatch case plus the unknown-subcommand message
- `tools/queue-kanban/model.go` (modified) — added `schemaFieldWarningText`, the contract's warning formatter

**Interface, as shipped:**

| Invocation | Behaviour |
|---|---|
| `frontmatter get <file> <field>` | prints the raw value, exit 0 |
| `… --normalize` | prints the canonical value; contract warning to **stderr** on an unrecognized value, default to stdout |
| `… --in-set terminal-success\|terminal-resolved` | prints nothing; exit 0 member / 1 non-member |
| absent field | exit 1, empty stdout, diagnostic on stderr |
| usage error / unreadable file | exit 2 |

## Decisions

- **D-01** (Step 3.5, above): ship the surface with tests only, no prose migration.
- **D-02: `--in-set` normalizes before testing membership.** DECIDE & STATE. A membership test on a raw `done` must resolve the alias first, or the contract's alias map would apply to `get` and silently not to the set check — which is the same class of half-applied contract the census exists to document. Verified end-to-end: `status: done --in-set terminal-success` exits 0.
- **D-03: `--in-set` delegates to `isCompletedStatus`/`isTerminalResolvedStatus`.** DECIDE & STATE. Those are the predicates the board already buckets with, so the command cannot drift into a second definition of "finished." Re-listing the set names here would have been the exact second copy this REQ is meant to eliminate.
- **D-04: a file with no frontmatter is exit 2, not an absent field.** DECIDE & STATE. A caller must be able to distinguish "this REQ does not set `domain`" (routine, exit 1) from "this is not a REQ" (broken input, exit 2). Collapsing them is how a typo'd path reads as a queue full of unset fields.
- **D-05: `schemaFieldWarningText` added in this REQ, not REQ-111.** DECIDE & STATE. REQ-111 deliberately left it out because it had no consumer and would have been dead code; this REQ is the consumer, so it lands with its caller.

## Qualification

Passed — 4 files verified in the diff, all requirements traced: subcommand registered and rejecting leftover tokens, parser reused rather than forked, `--normalize` delegating to REQ-111's normalizers, no write surface added, no prose changed. Judgment checks clean: the new file is 200 lines of substantive logic, not boilerplate, and the data path was exercised against real files rather than only mocks.

## Testing

**Tests run:** `go test ./...`, `gofmt -l .`, `go vet ./...`, an end-to-end binary smoke test, a shell-floor equivalence check, and `_dev/tests/contract-regressions.sh`.

**Result:** Go suite passes (9 new tests). `gofmt`/`vet` clean. Contract suite at its 7 pre-existing failures, unchanged. Smoke test against real REQ files confirmed every documented path, including `V=$(queue-kanban frontmatter get … title)` capturing a clean single value.

**Red-green validation:**
- RED: all 9 tests failed to compile — `undefined: readFrontmatterField`, `undefined: runFrontmatterCommand`. Confirmed before implementing.
- GREEN: all pass. The REQ's stated RED case (`frontmatter get … status` failing with `unknown subcommand`) now prints `completed` and exits 0.

**Floor verified, not assumed.** The REQ's binding constraint is that no action may lose its shell path. Checked directly: `awk 'NR==1&&/^---/{i=1;next} i&&/^---/{exit} i&&/^status:/{print $2}'` on the same fixture returns the same field the subcommand returns. The accelerator is preferred, never required.

**What the tests do not cover:** no action currently calls this command, so there is no integration test proving a prose site works with it. That is D-01's accepted cost, and the reason the follow-up below exists.

## Review

**Approve** — the surface is reachable, tested, floor-preserving, and adds no write path.

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 2 minor
**Minor:** (1) `--normalize` special-cases `status`/`testing_status` inside `runFrontmatterCommand` because they are not rows in the contract table — correct, but it puts a second piece of dispatch knowledge next to `normalizeSchemaField`'s own. If a third such field appears, that branch should move into the table layer. (2) The command has no in-repo consumer, so nothing yet proves the interface fits a real call site.
**Acceptance:** Pass — verified end-to-end against real REQ files, not only unit fixtures.
**Follow-ups created:** REQ-113 (`pending-answers`) — migrate the first prose read site, which is D-01's deferred half.

## Lessons Learned

**What worked:** Injecting the two writers instead of writing to `os.Stdout` directly. It made all nine tests plain function calls with no subprocess, and it forced the stdout/stderr split — the value on stdout, every diagnostic on stderr — to be an explicit design decision rather than an accident of `fmt.Printf` placement.

**What didn't:** I nearly shipped `schemaFieldWarningText` in REQ-111, where it would have been dead code with no caller. Deferring it was right, but I only noticed the problem when writing REQ-112's implementation and finding I needed it — the two-REQ split had made a shared helper's home ambiguous. A helper belongs in the REQ that first consumes it, not the REQ that first imagines it.

**Worth knowing:** Unit tests with injected writers prove the command's logic and nothing about its wiring. The `main.go` dispatch case is untested by the Go suite — a missing `case "frontmatter"` would leave every test green while the binary reported `unknown subcommand`. The binary smoke test is what covers that gap, and it is the step to repeat for any future subcommand.

## Orientation

`frontmatter.go`'s parser is now reachable from the command line: `queue-kanban frontmatter get <file> <field> [--normalize] [--in-set …]`. That closes the structural gap the census identified — prose could not call the shipped parser because no subcommand accepted a file and a field — so an action can now read and normalize a REQ field through one tested implementation instead of hand-rolled awk. Lives in `tools/queue-kanban/`, the board tool. `[MAP CHANGED]`: the tool gains an eighth subcommand and its first non-board-oriented read surface, while keeping its two write surfaces untouched. No action prose consumes it yet — that is REQ-113.
