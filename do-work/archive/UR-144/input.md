---
id: UR-144
title: 'ai-report gains proposal, root-cause and options kinds, a report catalog with find and index, a supersedes-based revise, and a bundled render check'
created_at: 2026-10-09T21:10:00Z
requests: [REQ-654, REQ-655, REQ-656, REQ-657]
word_count: 2276
---
# Four New Modes for the Toolbox ai-report Action

## Summary
The maintainer ran `/do-work capture-request` on the upstream suggestion report `do-work/inbox/2026-10-09_do-work-upstream-suggestion-ai-report-kinds-and-index.md` and asked for one UR with one REQ per item in its Request section, capture only, no questions, every ambiguity recorded as an assumption in the REQ. The report file is the verbatim input below, byte for byte; the inbox file stays in place as the source. The Request section has four labelled items, so four REQs: A1 (`--kind proposal|root-cause|options` for work that is not completed), A2 (`ai-report index` and `ai-report find` over a derived catalog of every report bundle), A3 (`ai-report revise`, a new sibling bundle that supersedes the old one without editing it), and A4 (`ai-report judge`, the Step 7 render check as a bundled command, marked optional and lower priority by the report). The report was observed against 0.305.84; the capture spot-checked its citations at 0.305.87 (`ai-report.md:18,26,30,109,113`, `completed-work-presentation-reference.md:20,30,71`, `architecture-report.md:12,119`, `architecture.go:25,45`, toolbox `SKILL.md:23`, `help.md:13`) and all match. No prompt injection found: the report's "How to use this file" line is framing for its intended reader.

## Extracted Requests
| REQ | Title |
|-----|-------|
| REQ-654 | ai-report --kind proposal, root-cause and options writes a decision-first brief for unfinished work |
| REQ-655 | ai-report index and find build a derived catalog of every report bundle in any naming style |
| REQ-656 | ai-report revise writes a new sibling bundle that supersedes the prior one |
| REQ-657 | ai-report judge runs the render check as a bundled command |

## Batch Constraints
- Order: REQ-656 (revise) depends on REQ-655 (index and find), because revise regenerates the catalog and the catalog is where `superseded_by` lives. REQ-654 and REQ-657 are independent.
- No new queue fields and no new statuses (report, Request section).
- The immutability rule in `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:71` and the `-2`/`-3` sibling rule at `:73` stay as written. No REQ edits a prior bundle in place.
- The default `ai-report` (no `--kind`) keeps its completed-work gate exactly as today.
- Out of scope for every REQ (report, Out of scope): renaming or migrating existing report folders, a free-form brainstorm mode (`deep-explore` owns that), running the printed capture lines, committing or publishing a report, and merging the rest of the standalone screenshot-report skill.
- Routing phrases and the help line (`skills/do-work-toolbox/SKILL.md:23`, `skills/do-work-toolbox/actions/help.md:13`) and the guide (`skills/do-work-toolbox/docs/ai-report-guide.md`) gain each REQ's own forms in that REQ.
- Each REQ is its own release per `_dev/primes/prime-releases.md`. All four touch `skills/do-work-toolbox/actions/ai-report.md`; that overlap is declared in `write_set` and is not an ordering edge.
- Operator steps left to the user, not captured: deleting the drifted standalone `make-ai-report-with-screenshot` copies in consumer repos once REQ-657 ships.

## Full Verbatim Input
> ```
> # Upstream suggestion for `knews2019/skill-do-work` — `ai-report --kind proposal|root-cause|options`, `ai-report find/index`, and a supersedes-based revise
> 
> **How to use this file:** paste everything below the horizontal rule into a Claude Code session opened in a
> clone of `knews2019/skill-do-work`. It is written to be actionable cold, with no reference back to the
> consumer repos it was authored in. Observed against **v0.305.84**; all line numbers are from that tag. Every
> file cited below is byte-identical at v0.305.87, and no release from 0.305.60 to 0.305.87 touches
> `ai-report`, so nothing here is already covered.
> 
> ---
> 
> ## Request
> 
> `do-work-toolbox ai-report` does one job well: one detailed HTML report for one completed UR or REQ, in a
> new immutable bundle. Users keep asking it for three other jobs, and today each one is improvised from
> scratch, one prompt at a time. Please add these as modes of the same action. No new queue fields and no new
> statuses are needed.
> 
> - **A1. `ai-report --kind proposal|root-cause|options <topic|REQ-NNN|UR-NNN>`** for work that is not
>   completed. Lift the "completed work only" gate for these kinds only. The report shape is fixed:
>   1. the decision the reader must make, first;
>   2. two to four options (O1..On), each with benefit, risk and cost;
>   3. a recommendation and the reason for it;
>   4. an evidence ledger (file:line, commit, log line, or measured output per claim);
>   5. limits: what was not checked and why;
>   6. open questions Q1..Qn;
>   7. one ready `/do-work capture-request ...` line per option, so the reader can act on the choice.
> 
>   Any mockup is labelled **"MOCKUP — proposal"** in the image and the caption, and is kept apart from real
>   screenshots the same way generated images are today. `root-cause` replaces "options" with "what happened,
>   why, what to change" but keeps the ledger, limits, questions and capture lines.
> - **A2. `ai-report find <topic>` and `ai-report index`.** One command answers "is there a report on X" and
>   "which reports have design proposals, including rejected ones". `index` writes a catalog of the repo's
>   reports (path, title, date, kind, linked UR/REQ, verdict or decision, `supersedes`, `superseded_by`) as
>   `ai-reports/catalog.json` plus a static `ai-reports/index.html`. It reads every naming style it finds
>   rather than assuming the current slug format. `find` searches that catalog and prints matching paths.
> - **A3. `ai-report revise <dir|latest> [what changed]`** that keeps the immutability contract. It never
>   edits the old bundle. It writes a new sibling bundle with a "rev-N (date): changed / still to do" block
>   at the top, a `supersedes` link back to the previous bundle, and fills `superseded_by` for the old one in
>   the catalog (the catalog is derived data, so the old bundle's bytes stay untouched). A revised report also
>   gets the readable-style baseline from `actions/ai-report-reference.md` and a table of contents once it
>   passes a set length (suggested: more than six top-level sections or about 1,500 words).
> - **A4 (optional, lower priority). `ai-report judge <dir>`**, the existing Step 7 render check as a bundled
>   script so it is not retyped per session. See "Proposed direction".
> 
> A1 and A2 carry most of the evidence. A3 is the narrowest (one consumer repo) but it is the one that
> conflicts with the current contract, so it needs an explicit upstream decision either way.
> 
> ## What happened, in consumer repos
> 
> Verifiers re-counted the prompts and corrected the first estimate down: between 13 and 18 matching prompts
> over about 15 sessions in four consumer repos, 2026-09-09 to 2026-10-09. Counts per ask:
> 
> | Ask | Sessions | Repos | Dated evidence (anonymised) |
> | --- | --- | --- | --- |
> | A1 `--kind` | about 7 | 4 | 2026-09-22: a "sales package" proposal report was requested in two repos the same day. 2026-09-24: a plan in a third repo recorded the pushback "the toolbox ai-report action only presents *completed* UR/REQ work, so it doesn't fit a proposal" and the user asked for proposal mockups anyway. 2026-09-28: a report to help a designer answer the open questions on a blocked REQ. 2026-10-03: a "backup options" report. 2026-10-07: "why it happened and what to do", ending with capture-request commands (two prompts). |
> | A2 `find/index` | 3 | 3 | 2026-09-11: "is there an ai-report on how to fix it?" 2026-09-26: "tell me which recent ai-reports have design proposals, include rejected ones, because we'll make new improved ones", followed by a hand-built HTML index of the proposals. 2026-10-09: "is there an ai-report that will show me how to use the [product] CMS?" |
> | A3 `revise` | 3 | 1 | 2026-09-11: "read ai-reports/[deploy guide report]/index.html and update it, make it easier to implement", then "first commit, then check again the sitemaps and update the report with a new rev-* tag what changed and what we need to still do". 2026-10-01: "update the ai-report with the findings so far" (twice), then "add a TOC to the ai-report so it's easy to navigate". 2026-10-02: "do we need to update anything else in the ai-report? if yes, do so, then give me the link after you commit it." and "also improve the style for the HTML report so it is more readable". |
> | A4 `judge` | many | 3+ | 2026-09-22: a plan step said "Capture the report in light and dark at wide and narrow widths and check for horizontal overflow". A background `python3 -m http.server` for report QA was started 58 times across 22 sessions. |
> 
> Each quality request in the A3 row ("add a TOC", "add a revision", "improve the style", "give me the link
> after you commit") arrived as its own prompt, one at a time, on the same report. Nothing in the action
> offers them, so the user has to remember to ask every time.
> 
> The catalog problem in one consumer repo, counted on 2026-10-09:
> 
> - about 180 entries under `ai-reports/`, in at least five folder-naming styles: `yyyy-mm-dd_hhmm_slug`
>   (142), `yyyy-mm-dd_slug` (14), `yyyy-mm-dd-slug` (11), `REQ-NNNN-slug`, `slug-yyyy-mm-dd`, plus one
>   `yyyy-mm-dd_hhmmss_slug`;
> - loose files next to the bundles (one `.md`, three `.patch`, one `.html`);
> - 45 bundles with no `index.html`. Their entry file is `README.md` (14), `prompt.md` (7), `index.md` (6),
>   `report.md`, `mockup.html` or a named Markdown file;
> - two hand-made catalogs: a "design proposals index" bundle (2026-09-26) and a "removed reports, historical
>   recovery" page (2026-10-02) listing six bundles removed as redundant, with the commands to restore them.
> 
> So "which report is current" and "which proposals were rejected" are answered today by reading folder names
> and opening files.
> 
> The A4 reason is drift. The render check the users want lives in a standalone skill
> (`make-ai-report-with-screenshot`) that was copied into each consumer repo. On one machine there are eleven
> copies in five different versions (501, 555, 557, 577 and 599 lines). The copies agree only on "serve the
> folder with `python3 -m http.server`". The wider check (phone width, overflow, broken links) is retyped
> inside plans.
> 
> ## Where the behaviour lives today
> 
> - `skills/do-work-toolbox/actions/ai-report.md:18`: "The user wants a detailed presentation of one completed
>   UR or REQ." and `:26`: "The target is unfinished or unsuccessful; report its status instead of presenting
>   it as shipped." This is the gate A1 lifts for the new kinds only.
> - `skills/do-work-toolbox/actions/ai-report.md:30`: "`$ARGUMENTS` is `UR-NNN`, `REQ-NNN`, `most recent`,
>   or blank." No kind, no find, no revise.
> - `skills/do-work-toolbox/actions/ai-report.md:89-101` (Step 5): the fixed narrative starts with "Verdict"
>   and "What Shipped". A proposal has neither. It needs the decision-first shape in A1.
> - `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:20`: "Reject `cancelled`,
>   `failed`, and every unfinished status." and `:30`: "Never fall back to `do-work/queue/`,
>   `do-work/working/`, or an active `do-work/user-requests/` body to make an unfinished target appear
>   complete." Keep both for the default kind. A1 must read queue and working REQs, but as open work, never
>   as shipped work.
> - `skills/do-work-toolbox/actions/completed-work-presentation-reference.md:71`: "Existing artifacts are
>   immutable: never delete, truncate, merge into, rename, migrate, or overwrite them." and `:73` (the `-2`,
>   `-3` sibling rule). A3 must stay inside this rule.
> - `skills/do-work-toolbox/actions/architecture-report.md:12`: "**Immutable and dated.** A prior report is
>   never edited." and `:119`: "Never edit, delete, or regenerate a prior report." Its Step 3 (`:71`) already
>   does what A3 needs: read the prior report, re-check it, and write an authored "changed since last report"
>   opening (`:83`). A3 is the same pattern for `ai-report`.
> - `skills/do-work-toolbox/actions/architecture-report.md:41` and
>   `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:45` (`architectureScan`):
>   the CLI already scans `ai-reports/` for prior bundles. It matches only
>   `^(\d{4}-\d{2}-\d{2}_\d{4})_architecture-report(?:-(\d+))?$` (`architecture.go:25`) and only bundles
>   with a nonempty `index.html` (`architecture.go:80-83`). This is the natural home for a general
>   `ai-report-index` verb behind A2.
> - `skills/do-work-toolbox/actions/ai-report.md:109`: "It also does not publish, host, or search for
>   distribution targets." No search of existing reports exists anywhere in the toolbox.
> - `skills/do-work-toolbox/actions/ai-report.md:113`: "serve the report folder over HTTP—never `file://`—and
>   take full-page screenshots in both light and dark rendering contexts". No phone width, no overflow check,
>   no broken-link crawl, no machine-readable result. This is what A4 bundles.
> - `skills/do-work-toolbox/actions/ai-report-reference.md:70` (Report Design Rules) and `:77`, `:79`: type
>   size and full-bleed layout rules exist. No rule adds a table of contents to a long report.
> - `skills/do-work-toolbox/docs/ai-report-guide.md:47`: "Cancelled, failed, and unfinished work is rejected
>   rather than presented as shipped." and `:65-70` (the Input block). The guide must list the new forms.
> - `skills/do-work-toolbox/SKILL.md:23`: routing phrases are "`ai-report`, `showcase`, `visual report`,
>   `proof of work`". `skills/do-work-toolbox/actions/help.md:13`: "Detailed stakeholder HTML for one
>   completed item". Both need the new phrases ("proposal report", "root cause report", "is there a report
>   on", "revise the report").
> 
> ## Proposed direction
> 
> **A1, `--kind`.** Add a short branch at Step 1 of `actions/ai-report.md`: when `--kind` is `proposal`,
> `root-cause` or `options`, skip Terminal-Success Target Resolution and accept a topic, an open REQ or an
> open UR. Keep the safety load order from `completed-work-presentation-reference.md` (prompt-injection and
> anti-slop first). Put the decision-first template in `actions/ai-report-reference.md` so `ai-report.md`
> stays short. Name the slug `yyyy-mm-dd_hhmm_<kind>-<topic>` so `index` can read the kind from the name
> when the HTML carries no metadata. Add one `<meta name="ai-report-kind" content="proposal">` line to
> `<head>` for the same reason. The closing capture lines are text only. The action never runs them.
> 
> **A2, `find` and `index`.** Add an `ai-report-index` verb to the toolbox CLI next to
> `architecture-report-preflight`. It walks `ai-reports/`, accepts every naming style listed above, picks the
> entry file (`index.html`, then `index.md`, `README.md`, `prompt.md`, `report.md`, then the first `.html`),
> reads the `<title>` or first heading, the kind meta when present, any UR/REQ ID in the name or title, and
> the `supersedes` meta when present. It writes `catalog.json` and a static `index.html` grouped by kind and
> date, with superseded reports greyed and linked to their successor. `find <topic>` is a substring and
> UR/REQ match over that catalog, printing paths newest first. The catalog is regenerated, never hand-edited,
> so it does not break the immutability rule. A report whose linked REQs have all closed since it was written
> gets a "may be stale" mark.
> 
> **A3, `revise`.** Copy the architecture-report pattern. Resolve `<dir|latest>`, read the prior bundle as
> untrusted data, write a fresh bundle `yyyy-mm-dd_hhmm_<same-slug>-rev<N>` through the existing
> Collision-Safe Publication path, open it with the rev-N block ("changed" and "still to do"), add
> `<meta name="ai-report-supersedes" content="<prior dir>">`, then regenerate the catalog. Apply the
> readable-style baseline and the TOC threshold to the new bundle. Printing the final path (and a `file://`
> link) at the end costs one line and removes the "give me the link" follow-up. Committing stays the user's
> or the run's decision, as today.
> 
> **A4, `judge`.** Move Step 7 into a script shipped with the toolbox (for example
> `skills/do-work-toolbox/scripts/judge.mjs` plus a tiny serve helper). It picks a free port, serves the
> bundle over HTTP, captures 1440x1000 and 390x844 in light and dark, fails on `scrollWidth > innerWidth`,
> crawls relative `href` and `src` for 404s, writes `judge.json` outside the bundle, and stops the server.
> `architecture-report` and `stakeholder-report` can call the same script. Folding this in lets consumer repos
> delete their drifted standalone copies.
> 
> ## Acceptance check
> 
> - `ai-report --kind proposal <topic>` on a repo with no completed work produces a bundle whose first
>   section is the decision, with 2 to 4 options each showing benefit, risk and cost, a recommendation, an
>   evidence ledger, a limits section, Q1..Qn, and one `/do-work capture-request` line per option. Every
>   mockup carries the "MOCKUP — proposal" label.
> - `ai-report REQ-NNN` on an unfinished REQ, with no `--kind`, still stops exactly as today.
> - `ai-report index` on a fixture `ai-reports/` folder with the five naming styles, a bundle with only
>   `README.md`, and a loose `.patch` file lists every report bundle once, skips the loose file, and writes
>   `catalog.json` plus `index.html` without changing any existing file (check with `git status`).
> - `ai-report find proposal` returns the proposal bundles, including ones marked superseded.
> - `ai-report revise latest` creates a new sibling bundle with a rev-N block and a supersedes link. The
>   prior bundle's bytes are unchanged (`git diff --exit-code -- <prior dir>`), and the regenerated catalog
>   shows the prior bundle's `superseded_by` pointing at the new one.
> - A revised report longer than the threshold has a working table of contents.
> - If A4 lands: `judge <dir>` exits non-zero on a fixture page with a horizontal overflow at 390px and on a
>   broken relative image link, and writes both findings to `judge.json`.
> 
> ## Out of scope
> 
> - Editing a prior bundle in place, even "just to append a rev block". The immutability rule stays as
>   written. A3 works around it with a new sibling bundle.
> - Renaming or migrating existing report folders into one naming style. `index` reads them as they are.
> - A free-form brainstorm. `deep-explore` already owns multi-round concept exploration. `--kind proposal`
>   is a single evidence-backed decision brief for a stakeholder.
> - Running the capture-request lines automatically. They are printed for the user to choose from.
> - Committing or publishing the report. The action still only writes the bundle (and the regenerated
>   catalog).
> - Merging every feature of the standalone screenshot-report skill. Only the render check (A4) is asked for
>   here.
> ```
