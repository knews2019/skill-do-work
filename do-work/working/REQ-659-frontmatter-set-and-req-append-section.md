---
id: REQ-659
title: 'frontmatter set and req append-section write REQ stamps and sections with the schema rules enforced'
status: claimed
created_at: 2026-10-09T21:14:04Z
user_request: UR-145
domain: backend
prime_files: ["_dev/primes/prime-shell-commands.md", "_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-658, REQ-660, REQ-661]
batch: cli-ergonomics
claimed_at: 2026-10-10T12:52:19Z
---
# frontmatter set and req append-section Write REQ Stamps and Sections With the Schema Rules Enforced
## What
`do-work-cli frontmatter set <file|REQ-N> <field> (<value> | --at now)` sits next to the existing `frontmatter get`. It refuses to overwrite an existing `*_at` field, formats the instant per the Timestamp rule, and keeps the Frontmatter Quoting contract. `do-work-cli req append-section REQ-N --section <name> --from <file>` inserts a REQ body section in canonical order, refuses a second copy, and is a no-op with exit 0 when the same bytes are already there.
## Why
Report item 2 (UR-145 input, "What happened"): 112 python heredocs that stamp `review_at:` or `integration_at:` into a REQ file, in about 24 sessions over 13 days. One example (2026-10-02, this checkout) asserted an anchor and ran `s.replace("\nintegration_at:","\nreview_at: 2026-10-02T18:47:04Z\nintegration_at:",1)`. Separately, `date -u +%Y-%m-%dT...` was used 356 times against 40 uses of `do-work-cli now`. The append-only and timestamp rules are ones a CLI can enforce and a heredoc cannot.
## Verified Facts (checked at 0.305.87)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:77-79`: `handleFrontmatter` accepts only `get` ("usage: frontmatter get <file> <field> [--normalize] [--in-set SET]").
- `skills/do-work/actions/work.md:309` (`integration_at: <now>`, "only if the field is absent") and `:375` (`review_at`, `remediation_at`, `re_review_at`, same absent-only rule). The report cites `:317` and `:383` at 0.305.84.
- `skills/do-work/actions/work-reference.md:59` (Timestamp rule: UTC `YYYY-MM-DDTHH:MM:SSZ`) and `:73` ("**Stamps are append-only.**", with four documented current-state exceptions: `status_changed_at`, `completed_at` on `failed` → `cancelled`, `blocked_at`, and `heavy_verified_at` / `heavy_verified_revision`).
- `skills/do-work/tools/replace-text-section.sh:10`: `replace-section` works on marker-delimited managed blocks, not REQ `## Testing` or `## Review` sections.
- `skills/do-work/actions/sample-archived-req.md`: the de facto section order (What, Red-Green Proof, AI Execution State, Triage, Plan, Exploration, Scope, Implementation Summary, Qualification, Testing, Review, Lessons Learned, Orientation). No command checks it. `advance` refuses an out-of-order write later ("later lifecycle evidence exists before ...", `tools/do-work-cli/internal/lifecycleadvance/advance_commands.go:328`).
- The archived request that gave frontmatter its read surface (`do-work/archive/UR-021/REQ-112-frontmatter-read-cli-surface.md:51`) put a `set` verb out of scope. That rule counted the write surfaces of the queue-kanban tool (`_dev/primes/prime-kanban-board.md:14`, `skills/do-work-board/tools/queue-kanban/frontmatter_cli.go:29`). The report asks to revisit it because the write happens anyway, by hand, with no validation.
## Detailed Requirements
1. `frontmatter set <file|REQ-N> <field> <value>`, or `--at now` in place of `<value>` for `*_at` fields.
2. `set` refuses to overwrite an existing `*_at` field (append-only rule) and leaves the file bytes identical on refusal.
3. `set` formats an `--at now` instant per the Timestamp rule and writes every value per the Frontmatter Quoting contract (`actions/work-reference.md` → Request File Schema).
4. `req append-section REQ-N --section <name> --from <file>` inserts the section in canonical order, refuses a different body for a section that already exists, and exits 0 with no change when the same bytes are already there.
5. Both commands resolve `REQ-N` to its one file in `do-work/working/` or `do-work/queue/` and refuse on zero or several matches.
6. Change the action lines that describe a hand stamp or hand section write (`actions/work.md:309`, `:375` and any other stamp or section write site the builder finds) to name the commands, keeping the hand path valid.
7. Release per `_dev/primes/prime-releases.md`.
## Constraints
- No new frontmatter field and no new status.
- The queue-kanban tool's `frontmatter` command stays read-only; its write-surface count (`_dev/primes/prime-kanban-board.md:14`) is unchanged.
- This REQ is its own release; do not fold REQ-658, REQ-660 or REQ-661 into it.
## Assumptions (recorded at capture, no questions asked)
- The new write surface lives in do-work-cli (`corehelpers`), not in queue-kanban. REQ-112's out-of-scope line was about queue-kanban's counted write surfaces, so this REQ does not reopen that count. By capturing this item the maintainer accepts the report's request to revisit REQ-112's "no `set`" decision for do-work-cli.
- `set` refuses to overwrite every existing `*_at` field, including the four documented current-state exceptions. Those exceptions keep their owning commands (`advance`, `unblock` and the others); `set` is only for append-only stamps. The builder records in Decisions if a caller needs otherwise.
- `--at` accepts only `now`, and only on `*_at` fields. A non-`*_at` field given `--at` refuses.
- For a field that is not `*_at`, `set` adds it when absent and replaces it when present.
- A new field is appended at the end of the frontmatter block. If the schema already defines a field order the CLI uses elsewhere, use that order instead.
- The canonical section order is the `sample-archived-req.md` order. A section name outside that order refuses rather than being guessed into place. If the builder finds the order already encoded in Go (for example in `advance`), reuse that one list and do not add a second.
- Archived REQs are not targets. `review-work`'s archived `## Review` append is out of scope here.
## Dependencies
None. Independent of REQ-658, REQ-660 and REQ-661.
## Builder Guidance
Certainty is high on the rules being enforced (they already exist in prose). Latitude: whether `req append-section` is a new top-level `req` command or a sibling of `frontmatter`, the help text, and how a REQ id resolves to a file (reuse any existing resolver).
## Red-Green Proof
**RED prompt/case:** On a fixture REQ with no `review_at`, run `frontmatter set REQ-N review_at --at now`, then run it again.
**Why RED now:** `commands.go:77-79` rejects any verb other than `get`.
**GREEN when:** The first call adds `review_at` once in Timestamp-rule form. The second call exits non-zero and the file bytes are identical. `req append-section` run twice with the same file leaves one section. A different file for an existing section refuses. A `## Testing` appended to a REQ that already has `## Review` lands before `## Review`.
**Validation:** Inferred during capture (from the report's Acceptance check).
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (18987 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its index row covers alternate stored-format writers and lifecycle section evidence; families `alternate-writer-contract-drift` and `lifecycle-section-evidence` fit a second writer of REQ frontmatter and body sections.
- `_dev/primes/lessons-action-files.md` as a whole satellite (7189 tokens, over budget; `slugged: partial`). Matching reason: family `alternate-writer-contract-drift` and its index row on alternate artifact writers.
## Full Context
See `do-work/user-requests/UR-145/input.md` for complete verbatim input (sections Request item 2, What happened, Where the behaviour lives today Item 2, Proposed direction 2, Acceptance check). No queued candidate shares this root cause (the queue held only REQ-654 to REQ-657, the ai-report modes, at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-cli-ergonomics.md`, Request item 2: "`frontmatter set <file|REQ-N> <field> (<value> | --at now)` next to the existing `frontmatter get`, plus `req append-section REQ-N --section <name> --from <file>`."*
