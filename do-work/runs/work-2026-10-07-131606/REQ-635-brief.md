# Builder brief — REQ-635 (VisibleSections: CommonMark 0-3 space rule for headings and fences, fence before comment, and no comment opener inside an inline code span)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-635-visible-sections-indented-heading-and-fence-comment-rules
- Branch: worktree-agent-REQ-635-visible-sections-indented-heading-and-fence-comment-rules (commit here; never merge, never rebase, never push)
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-635-visible-sections-indented-heading-and-fence-comment-rules.md — read it fully: What, Detailed Requirements, Constraints, Builder Guidance, Red-Green Proof.
- Hand-back file (the one main-tree path you may write, never commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-07-131606/REQ-635-handback.md

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md (tdd: true), backend.md. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/prime-do-work-cli.md. Lesson families to honour (the do-work-cli satellite is dropped for budget; read these entries by grep in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/lessons-do-work-cli.md): closed-enumeration-for-a-condition (REQ-460: "a trim that runs before the test can destroy the very bytes that constitute the structure"; this REQ is that trap firing again), lifecycle-section-evidence (hidden text and quoted text are separate classes).

## The change (decided at triage, do not reopen)
All in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-635-visible-sections-indented-heading-and-fence-comment-rules/skills/do-work/tools/do-work-cli/internal/requestmodel/visible_sections.go.

1. F1, heading rule: a `## ` heading is recognised only after 0 to 3 leading spaces. A tab anywhere in the indentation, or 4+ spaces, means "not a heading". Replace the `strings.TrimLeft(line, " \t")` at line 60.
2. Fence rule, same bound: `leadingPunctuationRun` (line 80) trims leading whitespace without bound. Apply the same 0-3 space bound (a tab or 4+ spaces means no fence opener or closer). A small helper shared by both checks is the obvious shape (builder latitude on its name; two words minimum per coding-guardrails § 5).
3. F5, fence before comment: on a line that does not start inside a comment, evaluate the fence-opener rule before the comment scan. A line that opens a fence is not scanned for `<!--` (CommonMark: the rest of a fence opener line is its info string). Lines that start inside a comment keep today's behaviour (a fence-looking line inside a comment opens nothing).
4. NEW at triage (D-01 below): a `<!--` inside an inline code span on the same line is not a comment opener. CommonMark: a backtick run of length n opens a code span that closes at the next backtick run of exactly length n on that line; text between them is literal. An unmatched backtick run is literal text, so a `<!--` after it still opens a comment (conservative; do not try multi-line code spans). This applies only to the opener scan on lines outside a comment; the `-->` close scan inside a comment ignores backticks (inside an HTML comment backticks mean nothing). Why this is in scope: the REQ-635 file itself has `<!--` inside backticks at Detailed Requirements item 4, so `VisibleSections` hides every section after it and `do-work-cli advance REQ-635` cannot see the `## Triage` section the orchestrator already wrote. Same function, same defect family (the comment scan opens on bytes that are not a comment). After your fix, `bash <worktree>/skills/do-work/tools/do-work-cli.sh --repo-root /Users/t2/Desktop/e1-experimental-repos/skill-do-work2 --format json advance REQ-635 --request-path do-work/working/REQ-635-visible-sections-indented-heading-and-fence-comment-rules.md` must report a phase past `agent judgment: triage and open questions` (it is read-only without gate args; you may run it, nothing else under the main tree).
5. Do NOT change: `## Plan <!-- note` with an unclosed comment still yields a zero-length section. Offsets and line endings untouched for `\n` and `\r\n`. Heading names stay raw (no normalisation).
6. Update the doc comments on `VisibleSections` and `leadingPunctuationRun` so they state the conditions (0-3 spaces, fence before comment, code spans), not examples.

## Tests (TDD: write first, run, record the failing message, then fix)
In visible_sections_test.go (requestmodel), each over both `\n` and `\r\n`:
- F1 RED body from the REQ: `"# Request\nExample:\n\n    ## Plan\n    user sample\n\nMUST keep.\n## Timing\nt\n"` → only `Timing`.
- 4 spaces and a leading tab (`"\t## Plan"`, also `" \t## Plan"`) → no `Plan`. 1-3 spaces still a heading.
- A 4-space-indented ```` ``` ```` line does not open a fence: `"# Request\n    ```\n## Plan\np\n"` → `Plan` visible.
- F5 body `"# Request\n``` <!--\n## X\n```\n## Plan\np\n## Timing\nt\n"` → exactly `Plan`, `Timing`.
- Code span: `"# Request\nUse `<!--` here.\n## Plan\np\n"` → `Plan` visible; a double-backtick span containing `<!--` too; an unmatched single backtick before `<!--` keeps today's hidden behaviour.
- `## Plan <!-- note` unclosed → still a zero-length `Plan` section (pin it if no test pins it today).
In /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-635-visible-sections-indented-heading-and-fence-comment-rules/skills/do-work/tools/do-work-cli/internal/requeststate/recovery_markdown_test.go: the F1 body wrapped in a REQ document through `stripGeneratedRecoverySections` keeps `"    ## Plan\n    user sample\n\nMUST keep.\n"` byte for byte (generated `## Timing` removed). Keep the 1/2/3-space cases green.
In /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-635-visible-sections-indented-heading-and-fence-comment-rules/skills/do-work/tools/do-work-cli/internal/lifecycletiming/lifecycle_timing_test.go: `replaceTimingSection` on a body holding an indented `    ## Timing` sample appends a real `## Timing` section and leaves the sample untouched; and on the F5 body it replaces the existing `## Timing` instead of appending a duplicate.
Each test file must stay under 30 seconds.

## Write boundary
visible_sections.go, visible_sections_test.go, requeststate/recovery_markdown_test.go, lifecycletiming/lifecycle_timing_test.go. Never touch VERSION, CHANGELOG.md, skills/do-work/VERSION, skills/do-work/CHANGELOG.md, skills/do-work/actions/version.md, do-work/lessons-index.md, skills/do-work/tools/do-work-cli/lessons-do-work-cli.md, or anything under do-work/ other than your hand-back file. The orchestrator writes the release and the lessons entry. Need another file: stop and say so in the hand-back.

## Verify before hand-back
- gofmt -l and go vet ./... in the worktree's skills/do-work/tools/do-work-cli.
- go test -C <worktree>/skills/do-work/tools/do-work-cli -count=1 ./... (the full module, REQ requirement 6). Report total wall and any test file near the 30s budget.
- The `advance REQ-635` read-only call above, showing the phase it reports.
- git diff --stat on your branch.

## Hand-back (/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-07-131606/REQ-635-handback.md)
Branch and commit hashes; file manifest (new/modified); red-green evidence (test, failing message before, pass after); P-A-U ([PLAN], [APPLY], [UNIFY] with git diff --stat and linters); ## Decisions (continue from D-02: D-01 is the code-span scope addition, D-02 is the orchestrator running estimate/plan with the fixed CLI; your first is D-03; each DECIDE & STATE or ESCALATE with Value/Risk); ## Discovered Tasks; lessons read; integration seams (none expected). Final chat reply under 3000 characters.
