## Review: REQ-659

**Approve**: both commands do what UR-145 item 2 asked, every refusal leaves the file bytes unchanged, and the `advance` refactor behaves exactly as before. Three small gaps are recorded below. None of them blocks.
Route B | merge `a56cd9c4` (range `1660977c..a56cd9c4`, 8 files, +542 / -22)

### What's built
- `do-work-cli frontmatter set <file|REQ-N> <field> (<value> | --at now)` writes one scalar on a working or queue REQ. It refuses to overwrite a non-empty `*_at` field, accepts a stamp only in `YYYY-MM-DDTHH:MM:SSZ` form, and refuses list or nested fields.
- `do-work-cli req append-section REQ-N --section <name> --from <file>` inserts a lifecycle section in canonical order. It is a no-op (exit 0) when the same body is already there. It refuses a different body, a second copy, and a name that is not in the canonical list.
- The canonical section order is now one list, `requestmodel.CanonicalSectionOrder`. `advance` reads its twelve "later sections" lists from it through `SectionsAfter`.

### Decisions / risks for you
- None needed. Only F1 (the `--from` heading gap) is likely to come up in real use, and its fix is about 5 lines.

### Findings

**Important:** None.

**Minor:**
- F1. `req append-section` does not check the `--from` body for a column-0 `## ` heading (`internal/corehelpers/request_writers.go:109`, `:158-166`). The shipped **Testing Section Template** (`actions/work-reference.md:529`) starts with `## Testing`, and `work.md:342` now points that template at this command. If an agent writes the template verbatim into the `--from` file, the command refuses with a misleading reason: "the inserted ## Testing would not be a visible section". In fact there are two visible copies. A body that holds a *different* canonical heading is not checked at all: `--section Qualification` with a body containing `## Review` exited 0 and wrote a stray `## Review` section (reproduced), which `advance` refuses later. Fix: before inserting, run `VisibleSections` on `sectionBody` and refuse any column-0 heading with a code such as `SECTION-BODY-HAS-HEADING` ("the --from file holds the body only; remove its ## heading"). This also replaces the misleading reason. impact-user-visible → report only
- F2. The containment check for the path form is lexical only (`request_writers.go:178-180`). If a directory under `do-work/working/` is a symlink, `frontmatter set` writes outside the tree. Reproduced: `do-work/working/esc -> <scratch>/outside` let `frontmatter set do-work/working/esc/REQ-900-out.md title pwned` rewrite the outside file, exit 0. D-16 says "Symlinks are refused by ReplaceExisting", but that covers only the last path component. The same lexical check also wrongly refuses an absolute path through the macOS `/tmp` → `/private/tmp` alias (reproduced, `REQUEST-NOT-ACTIVE`). The `..` escape is correctly refused (checked). Fix: call `filepath.EvalSymlinks` on the target's directory and on the repository root before `filepath.Rel`. impact-negligible → report only
- F3. `frontmatter set` writes any non-`*_at` scalar without checking the value against the field's allowed values. `frontmatter set REQ-801 status bogus-status` exited 0 (reproduced). `status` and `id` already have owning commands (`requeststate` transitions, `unblock`, `finalize`), so a success here looks like validation but is not. The capture's Assumption ("add when absent, replace when present") allows this behaviour, so it is not a requirement miss. Fix: refuse `status` and `id`, and name the owning command in the refusal. impact-rule-change → report only
- F4 (restatement sweep). `actions/work-reference.md:61` says "This paragraph is the only place in `actions/` that spells a command for obtaining one". That is now false, because `:73` spells `frontmatter set REQ-NNN <field> --at now`, which gets an instant while writing it. This is the same text as the integrator's first Discovered Task. Fix: change it to "obtaining an instant to hold", or name `:73` as the second sanctioned writer. impact-rule-change → report only
- F5 (restatement sweep). `actions/work-reference.md:73` says the four exceptions "keep their owning commands". Two of them are written by hand in prose, with no command: `abandon.md:62` re-stamps `completed_at` on the failed → cancelled path, and `clarify.md:138` and `:226` stamp `status_changed_at` on a flip. `frontmatter set` refuses these writes. The refusal is safe and the hand path stays valid. Fix: change the text to "keep their owning commands or the hand write at their defining site". impact-negligible → report only

**Nit:**
- F6. D-12 says a repeat on a CRLF file refuses as a conflict. On a fixture, the repeat was a no-op (exit 0), because the inserted section is LF. The file then has mixed line endings. This is harmless and nothing is duplicated, but the decision record is inaccurate. impact-negligible → report only
- F7 (anti-bloat). Bloat is low. Counted additions: 1 new source file and 1 new test file, 1 constant, 2 handlers, 2 exported names (`CanonicalSectionOrder`, `SectionsAfter`), 5 unexported helpers plus `requestIDPattern`, 0 new flags or options, and 5 tests. Each test names the failure it pins, and none is decorative. Two small extras the REQ did not need: (a) `requestIDPattern` (`request_writers.go:19`) duplicates the inline `^REQ-[0-9]+$` at `commands.go:159`, (b) the "no byte changed, no write" guard in `replaceRequestFile` fires only for a `set` that writes the value already there. Both are cheap. Also, in `--format json`, `set` reports `changes: []` and does not include the instant it wrote. This matches `frontmatter get`, and the prose uses `--format text`. impact-negligible → report only

### Requirements Checklist
- [x] 1. `frontmatter set <file|REQ-N> <field> <value>` and `--at now` for `*_at`: delivered (fixture: `--at now` printed `2026-10-10T14:50:04Z`, exit 0. `--at now` on `title` and `--at later` refused `FRONTMATTER-AT-INVALID`).
- [x] 2. Refuses to overwrite an existing `*_at` field, with identical bytes: delivered (second `review_at` call and `claimed_at` call exit 1 `FRONTMATTER-STAMP-EXISTS`. `cmp` of the before and after files was identical).
- [x] 3. Timestamp rule and Quoting contract: delivered (only `SetScalar` encodes. `title "it's: #x"` was written as `'it''s: #x'`, and a two-line value as `|-`. Rejected stamps: `.5Z` fractional, `2026-02-30`, `+00:00`, and a space form).
- [x] 4. `req append-section`: canonical order, conflict refusal, same-body no-op: delivered (Testing went in before Review. Scope went in before Testing and after the non-canonical `Why`. Orientation went at the end. A repeat with extra blank lines was a no-op. A different body was refused with `SECTION-CONFLICT`, a doubled Testing with `SECTION-DUPLICATE`, and `Notes` with `SECTION-NOT-CANONICAL`. An unclosed fence at the end of the file was refused by the visible-copy re-check, bytes identical). See F1 for the body-heading gap.
- [x] 5. Resolution: delivered (`REQ-999` → `REQUEST-NOT-FOUND`. The test covers a working plus queue duplicate → `REQUEST-AMBIGUOUS`. An archived id → `REQUEST-NOT-ACTIVE`. `do-work/working/../archive/...` and `../../../etc/hosts` were refused). See F2 for symlinked directories.
- [x] 6. Action lines name the commands and keep the hand path: delivered (`work.md:21`, `:201`, `:288`, `:309`, `:342`, `:375`. `work-reference.md:73` holds the one `frontmatter set` argv).
- [ ] 7. Release: pending. This is the integrator's job at finalization and is not part of this diff.
- Constraints: no new field or status. queue-kanban is untouched (`frontmatter_cli.go:33` still says there is no `set` verb, and `prime-kanban-board.md` is unchanged). The release was not folded into another REQ.

### Acceptance Testing

**Result: Pass**
- UR-145 acceptance check 2 was reproduced on a real binary through `skills/do-work/tools/do-work-cli.sh --repo-root <fixture> --format text|json` in `scratchpad/int-659/review/fx`. The first `frontmatter set REQ-N review_at --at now` added the field once. The second call exited 1 with identical bytes. `req append-section` run twice with the same file left one section. A different file refused. A `## Testing` added to a REQ that already had `## Review` was placed before it.
- `go test -count=1 ./internal/corehelpers/ ./internal/requestmodel/`: exit 0 (8.6 s, 0.4 s). `go test -count=1 -run 'TestAdvanceCommandPhaseMatrix|TestAdvanceCommandRefusesMalformedAmbiguousAndImpossibleStates' ./internal/lifecycleadvance/`: exit 0.
- The `advance` refactor behaves the same: a script compared each of the 12 removed `hasAnySection` literal lists with `SectionsAfter(<name>)` over `CanonicalSectionOrder`. All 12 are equal in content and order (4 × Triage, 2 × Plan, and one each for Exploration, Scope, Implementation Summary, Qualification, Testing, Review). The two lists left as literals (`:174` Route A set check, `:253` single Orientation check) do not express a tail of the order.

### Suggested Additional Testing
- Edge case: a `--from` body containing `## <name>` or another canonical heading (F1). A focused test should be added with the fix.
- Edge case: `frontmatter set` through a symlinked directory and through an aliased absolute root (F2).
- Not separately tested: D-15's visible-copy re-check has no unit test. It was checked here by hand on a fixture with an unclosed fence, and it refused with identical bytes.
- Regression: an orchestrated run that uses `frontmatter set` for `dispatch_at` with a held instant and `req append-section` for `## Testing`. The integrator already did both live on REQ-659 itself.

### Scores (on the record, not the headline)

**Overall: 91%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | 1-6 delivered. Release pending at finalization |
| Code Quality | 88% | Reuses the existing parser, writer and resolver. F1 and F2 gaps |
| Test Adequacy | 88% | Five focused RED→GREEN tests. Re-check and heading cases untested |
| Scope | 95% | `commands_test.go` out-of-boundary edit was needed and recorded (D-11/D-17) |
| Risk | Low | Local CLI. Refusals never write. The symlink write needs a planted directory symlink |
| Acceptance | Pass | UR acceptance check 2 reproduced on the real binary |

### Follow-ups created
- None (7 findings report only)

## Review

**Overall: 91%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 88% |
| Test Adequacy | 88% |
| Scope | 95% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token, the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. `req append-section` does not check the `--from` body for a column-0 `## ` heading (`internal/corehelpers/request_writers.go:109`, `:158-166`). The Testing Section Template's own `## Testing` line gets the misleading reason "would not be a visible section". A different canonical heading (for example `## Review`) is written unchecked. Fix: refuse any visible heading in the body. impact-user-visible → report only
- F2. The path-form containment check is lexical (`request_writers.go:178-180`). A symlinked directory under `do-work/working/` lets `frontmatter set` write outside the tree (reproduced), and an absolute path through the macOS `/tmp` alias falsely refuses. Fix: `filepath.EvalSymlinks` before `Rel`. impact-negligible → report only
- F3. `frontmatter set` replaces any non-`*_at` scalar with no check on the value (`status bogus-status` exit 0). Fix: refuse `status` and `id`, and name their owning commands. impact-rule-change → report only
- F4. The Timestamp rule says it is "the only place in `actions/` that spells a command for obtaining one" (`actions/work-reference.md:61`). `:73` now spells `frontmatter set ... --at now`. impact-rule-change → report only
- F5. `actions/work-reference.md:73` says the four exceptions "keep their owning commands", but `abandon.md:62` (`completed_at` re-stamp) and `clarify.md:138` and `:226` (`status_changed_at`) are hand writes. The refusal is safe. impact-negligible → report only
- F6 (nit). D-12 says a CRLF repeat refuses as a conflict. It is actually a no-op, and the file ends with mixed line endings. impact-negligible → report only
- F7 (nit, anti-bloat). Additions are low. `requestIDPattern` duplicates `commands.go:159`. The no-change write guard only fires on a same-value `set`. JSON `set` does not report the written instant (this matches `get`). impact-negligible → report only

**Acceptance:** Pass. UR-145 check 2 was reproduced on the real binary. Refusals leave the bytes identical. The 12 `advance` lists equal `SectionsAfter` tails exactly.
**Restatement sweep:** redefined (E1) how REQ `*_at` stamps are written (`frontmatter set`, append-only enforced), (E2) how lifecycle sections are written (`req append-section`), (E3) the section order as `requestmodel.CanonicalSectionOrder`. Stale: `work-reference.md:61` "only place … spells a command" (F4). `work-reference.md:73` "owning commands" vs `abandon.md:62`, `clarify.md:138`, `:226` (F5). Consistent: `work-reference.md:147`, `:155`, `:178`, `:182`, `:199`, `:209` (schema comments citing the rule). `fan-out-reference.md:163`, `:167` (absent-only `builder_handback_at` / `dispatch_at`, hand path valid, D-09). `work.md:21`, `:201`, `:288`, `:309`, `:342`, `:375`, `:409`, `:497`. `sample-archived-req.md` section order (matches. It omits Pre-Flight, which sits between Scope and Implementation Summary in the list and in `work-reference.md:494`). `review-work.md` Append to REQ File (`## Review` added at the end). `queue-kanban/frontmatter_cli.go:33` and `prime-do-kanban.md:3` (about queue-kanban, which stays read-only). Go: no other hard-coded section-order list. `advance_commands.go:174` and `:253` are a set check and a single check, not order tails.
**Suggested testing:** 4 items
**Follow-ups created:** None (7 findings report only)

*Reviewed by review-work action*
