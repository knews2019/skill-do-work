---
id: UR-136
title: 'A REQ file the association walk cannot parse claims no paths, and do-work commit continues'
created_at: 2026-10-06T15:27:35Z
requests: [REQ-634]
word_count: 599
---
# A REQ File the Association Walk Cannot Parse Claims No Paths, and do-work Commit Continues

## Summary
A consumer bug report against do-work 0.305.68, pasted into `do-work validate-feedback` and triaged against the code on 2026-10-06. One archived REQ with an odd number of backticks on one Implementation Summary line makes `do-work commit` Step 3 exit 2 with `PARSE-FAILED` in every later commit of a ~2,000-REQ repository, because the association walk (`AssociateProjectPaths`, `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go`) returns a single file's parse error as the result of the whole walk. The report names `qualificationSummaryEntries` as the parser; the walk's parser is `allBacktickedPaths`, same per-line rule.

Triage verdicts: accepted the root cause (function name corrected), the walk-level abort, the per-file degrade with exit 0, the docs update (plus `skills/do-work-toolbox/actions/inspect.md`, which carries the same sentence), and the tests reduced to the behaviour kept. Pushed back on per-bullet joining and the first-span-only rule (no new grammar; the multi-path bullet contract is tested on purpose). Qualify and scope-drift keep their strict parse. The optional in-flight marker is dropped.

Decisions recorded with the maintainer (2026-10-06):
- D1 The span rule does not change. No bullet-joining or wrapped-line grammar.
- D2 The only fix: in `AssociateProjectPaths`, a parse error on one REQ file means that file claims no paths; the walk continues; exit 0. The skip is reported as a stderr line naming the file. Owner rows on stdout stay exactly as they are. No new row type, no second call.
- D3 Flip the two tests that pin `PARSE-FAILED`; update the commit.md and inspect.md sentence so `PARSE-FAILED` is no longer an exit-2 case; qualify and scope-drift stay strict.
- D4 Drop the in-flight marker. One REQ.
- Deviation accepted at plan approval: the stderr line is produced by the runtime's existing finding printer, so it is a warning finding with a code and is also visible under `--format json`.

Maintainer's words during triage: "we should not kill the commit, why would it even trip on this formatting issue? I don't think it's a good idea to implement all variations of path representations that non-deterministcs llm's can generate just to do the commit."

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-634 | A REQ file the association walk cannot parse claims no paths, and do-work commit continues |

## Batch Constraints
- No parser grammar changes anywhere (D1).
- In-flight `working/` claims stay included as documented; no marker (D4).

## Full Verbatim Input
> ```
> Make REQ Implementation Summary parsing tolerant of wrapped backtick spans, and make
> `do-work commit` association degrade per file instead of aborting.
> 
> ## Observed failure (do-work v0.305.68)
> 
> `do-work commit` Step 3 runs `scripts/protected-inventory.sh associate`, which calls
> `do-work-cli protected-inventory associate`. In a repository with ~2,000 archived REQs it
> exits 2 with:
> 
>     PARSE-FAILED: unmatched backtick in Implementation Summary
> 
> and the whole association step is lost. Every future `do-work commit` in that repository
> fails the same way, because the offending files are immutable archive records.
> 
> ## Root cause
> 
> `internal/corehelpers/checks.go`, `qualificationSummaryEntries`:
> 
>     for _, line := range lines {
>         if !strings.HasPrefix(strings.TrimSpace(line), "- `") { continue }
>         parts := strings.Split(line, "`")
>         if len(parts)%2 == 0 {
>             return nil, fmt.Errorf("unmatched backtick in Implementation Summary")
>         }
>         ...
> 
> It counts backticks per physical line. A Summary bullet that an author or a formatter wrapped
> at ~100 columns puts the closing backtick on the next line, e.g. (archived REQ, verbatim):
> 
>     - `tests/unit/level-transfer.test.js` (modified) — new `describe("harden a malformed imported theme
>       (REQ-930, UR-238)")` with a `buildMalformedThemeImportZip` helper and 2 import cases: a wrong-type
> 
> Line 1 has 3 backticks, so the parser returns an error for the file. `AssociateProjectPaths`
> (`internal/corehelpers/inventory.go`) propagates that single-file error as the result of the
> whole walk, and the shim turns it into `PARSE-FAILED` exit 2. One malformed line in one
> archived REQ therefore disables association for the entire repository, permanently, because
> `do-work/archive/` is immutable by the skill's own rule.
> 
> ## Required behavior
> 
> 1. Parse backtick spans per Summary bullet, not per line. A bullet is the `- ` line plus its
>    indented continuation lines (same rule Markdown uses). Join the bullet's lines, then split on
>    backticks. The first backtick span is the path; later spans are prose and are ignored for
>    path extraction (today they are all treated as paths, which also makes
>    `buildMalformedThemeImportZip` above a "path").
> 2. A bullet that is still unbalanced after joining is a per-file finding, never a walk-level
>    error. `AssociateProjectPaths` must skip that REQ's Summary, continue the walk, and return
>    the skipped files alongside the associations. Exit 0 with an `ASSOCIATION-SUMMARY-UNPARSED`
>    warning finding per skipped file (path + line number), so `actions/commit.md` can report
>    them and still associate everything else.
> 3. Keep the hard `unmatched backtick` error only where the REQ being parsed is the one being
>    finalized (the qualification check at `checks.go:321`), where the author can fix it. Archive
>    reads are lookups, not validation.
> 4. `actions/commit.md` Step 3 and `docs/prescribed-shell-primitives.md` → "Protected inventory
>    fallbacks": replace "any other exit 2 such as PARSE-FAILED is an error to report" with the
>    new contract: association always completes; unparsed Summaries are reported by path.
> 
> ## Tests
> 
> - `qualificationSummaryEntries` on a bullet whose closing backtick is on the continuation line
>   returns exactly one entry, the path, with the verb from `(modified)`.
> - A Summary with one unbalanced bullet and two good bullets yields the two good paths plus one
>   unparsed-bullet finding; it does not return an error.
> - `AssociateProjectPaths` over a fixture archive with one unparseable REQ and one claiming REQ
>   still associates the claiming REQ's paths and reports the unparseable file once.
> - Shim mode: `protected-inventory associate` exits 0 and prints the `<owner>\t<path>` rows for
>   the good files; `PARSE-FAILED` is no longer emitted.
> 
> ## Optional, same area
> 
> Association currently treats an in-flight `do-work/working/` REQ's Summary as a claim on a path
> regardless of where the uncommitted change came from. In my case the script would have assigned
> another session's edit of `apps/game/daily-core-loop-state.js` to the REQ that happened to list
> that file. Consider marking matches from `working/` REQs as `ASSOCIATION-FOUND-INFLIGHT` so
> `actions/commit.md` asks the agent to confirm against the diff before committing under that
> REQ's id.
> ```