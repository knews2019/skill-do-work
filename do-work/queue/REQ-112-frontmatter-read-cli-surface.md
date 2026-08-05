---
id: REQ-112
title: Give frontmatter.go a CLI surface so prose can stop reimplementing it
status: pending
created_at: 2026-08-05T15:53:39Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

- [ ] Should the first consumer land in the same REQ — one action's read site migrated to prove the surface works end-to-end — or should the subcommand ship with tests only?
  Recommended: tests only. Keeps this REQ's diff inside `tools/queue-kanban/`, and a migration needs the target action's own review.
  Also: migrate exactly one low-risk read site (e.g. `actions/commit.md` Step 3's terminal-success check) as a worked example.

## Red-Green Proof

**RED prompt/case:** Run `queue-kanban frontmatter get do-work/archive/UR-020/REQ-110-census-completeness-floor-note.md status`. It fails today with `unknown subcommand "frontmatter"` and exit 2, from `main.go`'s dispatch default.
**Why RED now:** The seven registered subcommands are the only entry points, and none accepts a file-and-field pair, so the shipped parser has no caller outside the board's own walk.
**GREEN when:** That command prints `completed` and exits 0; `--normalize` on a `domain: back-end` fixture prints `backend` with the contract warning on stderr; and `--in-set terminal-success` exits 0 for `completed-with-issues` and 1 for `failed`.
**Validation:** Inferred during capture — the surface shape is the census's candidate #1 proposal, not a user-specified API.

## Full Context

See `do-work/user-requests/UR-021/input.md` for the complete verbatim input and batch constraints.

---
*Source: census finding — `frontmatter.go` has no CLI surface, so all ~95 prose frontmatter reads are hand reimplementations by construction (`decisions/audits/2026-08-05-shell-logic-in-prose-census.md` §1 structural fact 1, §4 candidate 1)*
