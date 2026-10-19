---
id: UR-157
title: 'ai-report: commit the bundle as the project rules say, and expect one judge link finding per href'
created_at: 2026-10-19T14:29:39Z
requests: [REQ-699]
word_count: 1154
---
## Summary

Captured by `do-work validate-feedback --capture --run` from the consumer-authored upstream suggestion file `do-work/inbox/2026-10-19_do-work-upstream-suggestion-ai-report-revise-commit-rule.md` (the session paste was that file's lower half, from "## What happened" down; the file is the complete document and is the verbatim input below). Triage verdicts, in the order the paste listed them: Finding 1 (revise step 7 "Do not commit" is absolute) Accept; Finding 2 (Step 8 creation path silent on committing) Accept; Finding 3 (judge wording "that one link" but the judge reports one `AI-REPORT-JUDGE-BROKEN-LINK` finding per href) Accept; Finding 4 (checklist line names no destination for the catalog pair) Accept. No Already done, Push back or Discuss items, so no Discuss answers. All four findings edit one file, two of them the same line, so they are captured as one REQ instead of validate-feedback's default of one REQ per finding; the four findings are listed separately under Detailed Requirements with their provenance.

## Extracted Requests

| REQ | Title |
|-----|-------|
| REQ-699 | ai-report: commit the bundle as the project rules say, and expect one judge link finding per href |

## Batch Constraints

- Shipped prose only: `skills/do-work-toolbox/actions/ai-report.md` and `skills/do-work-toolbox/docs/ai-report-guide.md`; no Go change, the judge already behaves as the new wording says.
- Unchanged: the prior bundle's bytes, catalog regeneration in revise step 6, the run's ownership of commits inside `do-work run`, and no push in any form.
- A change to shipped files is a release (`_dev/primes/prime-releases.md`).

## Full Verbatim Input
> ```
> # Upstream suggestion for `knews2019/skill-do-work` — the ai-report revise form must not forbid committing
> 
> **How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
> clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
> consumer repo it was authored in. Observed against **v0.306.0**; all line numbers are from that tag. Still
> present in v0.306.1.
> 
> ---
> 
> ## Request
> 
> The `ai-report` revise form ends with "Do not commit. Committing stays the user's or the run's decision."
> That is right inside a `do-work run`, where the run owns the commit. It is wrong when the user invokes
> `ai-report revise` directly in a project whose own instructions say to commit every verified increment: the
> agent then stops with a verified bundle and a printed commit command, the user deploys, and the report is
> missing from the deployed site because the deploy ships git content. The ordinary creation path (Step 8)
> says nothing about committing, so the two forms of the same action also disagree with each other.
> 
> Please make the commit step defer to the project's rules, in both forms:
> 
> 1. **Revise form, step 7** (`skills/do-work-toolbox/actions/ai-report.md:59`): replace the absolute "Do not
>    commit" with a rule that commits when the project asks for it (bundle plus the regenerated catalog pair,
>    by exact path), prints the command otherwise, and leaves the commit to the run inside a run.
> 2. **Step 8, Verify and Report the Result** (`skills/do-work-toolbox/actions/ai-report.md:176-178`): say the
>    same for the creation path, and have the summary name either the commit hash or the exact command.
> 3. **Revise step 5 and Step 7 render wording** (`ai-report.md:57` and `:170`): the judge reports one
>    `AI-REPORT-JUDGE-BROKEN-LINK` finding per `href` into the prior bundle, not "that one link"; a report that
>    links the prior bundle twice (rev block plus an anchor link in the body) gets two expected findings.
> 
> Nothing here changes the revise contract: the prior bundle's bytes still never change, the catalog is still
> regenerated in step 6, and pushing stays with the user.
> 
> ## What happened, in one consumer repo (timeline)
> 
> A consumer repo whose `CLAUDE.md` says "commit completed increments frequently" and whose public site builds
> its report index from the `ai-reports/` folder of a checkout updated by `git pull`. Suite v0.306.0 throughout.
> 
> | When (UTC) | Event |
> | --- | --- |
> | Oct 19, ~13:30 | A session finished an `ai-report revise` bundle (47 browser assertions, render judge clean apart from the prior-bundle links), regenerated `catalog.json` and `index.html`, and handed back a commit command instead of committing, citing step 7. |
> | Oct 19, ~14:00 | The user deployed. The live report index showed the 18 October original and not the revision, because the bundle and the catalog pair were untracked. |
> | Oct 19, 14:0x | The user asked: "what made it not commit the report?" The session answered that it had let step 7 outrank the project's `CLAUDE.md`, and committed. |
> 
> Cost: one deploy round trip, one extra commit, and a correction the user should not have had to make. The suite
> already held evidence that consumers expect the commit: the 2026-10-09 upstream report on report kinds and
> the catalog (folder `2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index`) quotes a user asking,
> on Oct 2, "if yes, do so, then give me the link after you commit it".
> 
> ## Where the behaviour lives today
> 
> - `skills/do-work-toolbox/actions/ai-report.md:59` (Revise form, step 7): "7. **Do not commit.** Committing
>   stays the user's or the run's decision." An absolute instruction; an agent following the action literally
>   cannot commit even when the consumer's instructions require it.
> - `skills/do-work-toolbox/actions/ai-report.md:176-178` (Step 8, Verify and Report the Result): confirms the
>   bundle, removes the judge directory and prints a summary. Nothing about committing, so the creation form is
>   silent where the revise form forbids.
> - `skills/do-work-toolbox/actions/ai-report.md:201` (Verification Checklist): "A revise wrote a new `-rev<N>`
>   sibling with the rev block and supersedes meta, regenerated the catalog, and left the prior bundle
>   byte-identical." No line on where the regenerated catalog pair ends up.
> - `skills/do-work-toolbox/actions/ai-report.md:57` (Revise form, step 5): "so it reports that one link as
>   `AI-REPORT-JUDGE-BROKEN-LINK`: this finding is expected". `:170` (Step 7): "(for a revise, until only the
>   expected prior-bundle link finding remains)". Observed: two findings for two hrefs into the same prior
>   bundle, one with an `#anchor`.
> - Neighbouring rule that already handles the run case: the core `work` action commits a REQ's output as part
>   of the run, which is what "the run's decision" in step 7 means. That stays as it is.
> 
> ## Proposed direction
> 
> **`skills/do-work-toolbox/actions/ai-report.md:59`, replace step 7:**
> 
> > 7. **Commit as the project's rules say.** Inside a `do-work run`, the run commits the bundle with its other
> > output; do nothing here. Invoked directly, follow the consumer project's instructions on committing (its
> > `CLAUDE.md` or equivalent): when they ask for completed increments to be committed, commit the new bundle
> > and the regenerated `ai-reports/catalog.json` and `ai-reports/index.html` by exact path in the same turn and
> > report the hash; otherwise print that exact commit command. Never stage anything outside the bundle and the
> > catalog pair, and never push.
> 
> **`skills/do-work-toolbox/actions/ai-report.md:178`, append to the Step 8 summary sentence:**
> 
> > The summary also names the commit hash when the bundle was committed under the project's rules, or the
> > exact commit command (bundle by exact path, plus the catalog pair for a revise) when the project leaves the
> > commit to the user. Inside a run, the run commits.
> 
> **`skills/do-work-toolbox/actions/ai-report.md:57`, replace "so it reports that one link as
> `AI-REPORT-JUDGE-BROKEN-LINK`: this finding is expected" with:**
> 
> > so it reports each `href` into the prior bundle as `AI-REPORT-JUDGE-BROKEN-LINK`, one finding per `href`
> > (an anchor link counts separately): these findings are expected
> 
> **`skills/do-work-toolbox/actions/ai-report.md:170`, replace "(for a revise, until only the expected
> prior-bundle link finding remains)" with:**
> 
> > (for a revise, until only the expected prior-bundle link findings remain)
> 
> **`skills/do-work-toolbox/actions/ai-report.md:201`, append to the checklist line:**
> 
> > …left the prior bundle byte-identical, and committed the bundle and the catalog pair or printed the commit
> > command, as the project's rules say.
> 
> ## Acceptance check
> 
> - `grep -n "Do not commit" skills/do-work-toolbox/actions/ai-report.md` prints nothing.
> - `grep -n "Commit as the project's rules say" skills/do-work-toolbox/actions/ai-report.md` prints the step 7
>   line.
> - `grep -n "one finding per" skills/do-work-toolbox/actions/ai-report.md` prints the revise step 5 line.
> - In a consumer whose instructions say to commit increments, a direct `ai-report revise` ends with a commit
>   hash in its summary and `git status` shows neither the bundle nor `ai-reports/catalog.json` as untracked.
> - Unchanged: the prior bundle's bytes, the catalog regeneration in step 6, no push, and the run's ownership
>   of commits inside `do-work run`.
> 
> ## Out of scope
> 
> - Pushing. It stays with the user in every form.
> - The consumer-side fix (a memory note that the project's commit rule wins) is already done in the consumer.
> - Serving the bundle over HTTP for browser checks when `file://` is blocked; consumer tooling, not the suite.
> ```
