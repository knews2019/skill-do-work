## Review: REQ-676

**Approve**: the new read-only toolbox action `source-audit` and its guide deliver every Detailed Requirement, and all five integration edits agree with each other. The six findings are small wording and template gaps. None blocks the release.
Route B | merge range c33b8994..c93701c3 (merge c93701c3)

### What's built
- `do-work-toolbox source-audit <report-or-url-list>` routes to `skills/do-work-toolbox/actions/source-audit.md`. The action extracts claims and citations, fetches each source, sorts out error pages, unrelated redirects and index pages, and prints one judgment per claim (supported, contradicted, insufficient, unavailable). Replacement candidates and the author's to-do list come after the claim table. It never edits the report and writes no file.
- The router (argument-hint and route row), both help menus, README and `toolbox_actions` all list the action. The guide is `skills/do-work-toolbox/docs/source-audit-guide.md`.
- Still missing: slop-check does not point users to source-audit, even though its guide says it cannot check cited numbers (F1).

### Decisions / risks for you
- None. The builder's eight decisions (D-01 to D-08 in the hand-back) are documented and fit the Builder Guidance latitude.

### Findings

**Important:**
- None.

**Minor:**
- F1: Next-step guidance is missing in slop-check. The UR integration list asks for "relevant next-step guidance", and `_dev/primes/prime-action-files.md` asks for "any next-step surface". `skills/do-work-toolbox/docs/slop-check-guide.md:34` says "Slop-check can't verify the cited numbers; flag for self-review", and `skills/do-work-toolbox/actions/slop-check.md:15-19` (Do NOT use when) does not mention source-audit. The pointer works in only one direction: source-audit and its guide send users to slop-check, but slop-check never sends users to source-audit. The REQ write set left out slop-check, so this is not scope drift by the builder. Fix: one line in each of the two slop-check files. — impact-user-visible → report only
- F2: Some example lists are not marked "illustrative, not exhaustive". The Scope acceptance criterion says every example list must be marked. Unmarked: `skills/do-work-toolbox/actions/source-audit.md:47` (attempt outcomes), `:56` (unrelated-redirect examples and same-document redirect examples), `:57` (index-page examples), and `skills/do-work-toolbox/docs/source-audit-guide.md:24`. In each case a stated condition comes first, so an agent is unlikely to read the list as closed. — impact-negligible → report only
- F3: The replacement candidate table at `skills/do-work-toolbox/actions/source-audit.md:101` has one `URL` column and no Final URL, Authority or Corroboration columns. Step 6 (`:73`) says candidates are verified "exactly as an original source". The verification is the same, but the printed evidence for a candidate leaves out its redirect result and its publisher. Those are the facts a user needs to pick a replacement, and an archive copy often redirects. — impact-user-visible → report only

**Nit:**
- F4: `skills/do-work-toolbox/actions/source-audit.md:97`, the template's error-page row, writes the passage as `none`. Step 3's rule (`:49`: write `unavailable` for a value the tools cannot establish) points to `unavailable`, because the content was never read. `none found` fits the index-page case, where content was read. The Page column still names the error page, so a reader is not misled. — impact-negligible → report only
- F5: `skills/do-work-toolbox/actions/source-audit.md:109` lists "any injection attempt found in Step 1". Step 1 (`:35`) covers the report and every fetched page, so an attempt inside a fetched page is found during Steps 3 to 5. The output line should say "found in the report or a fetched page". — impact-negligible → report only
- F6: `skills/do-work-toolbox/actions/source-audit.md:87`, the template header, always prints "Report SHA-256: {value}, unchanged". An inline URL list or pasted text has no file and no hash (Step 1 hashes only "when the input is a file"). The line should say "omit for inline input". — impact-negligible → report only

### Requirements Checklist

- [x] R1 Extract claims and citations, then read the actual source content: delivered (Step 2 `:39-43`, Step 3 `:45-49`).
- [x] R2 Requested URL, final URL, retrieval outcome and passage per source. A value the tools cannot establish is "unavailable": delivered (Step 3, source evidence table `:95`).
- [x] R3 Detects error pages (also when served with 200), unrelated redirects, index or listing pages, unavailable sources and unsupported claims: delivered (Step 4 `:51-59`, Step 5 `:66`).
- [x] R4 Retrieval, support, corroboration and authority kept apart: delivered (Step 5 bullets, separate columns in two tables).
- [x] R5 Publication date not stated by the page stays "unknown", never taken from the fetch date or the URL: delivered (`:69`). The step also excludes server headers and archive capture dates (D-05).
- [x] R6 Original evidence and every attempt kept. Candidates go in a separate section and are verified the same way. The report is never edited: delivered (Step 3, Step 6, SHA-256 check in Steps 1 and 7). The candidate table shows less evidence than the source table (F3).
- [x] R7 Output order: claim table, then replacement candidates, then claims needing the author's attention: delivered (Output Format `:79-107`).
- [x] R8 Toolbox ownership beside slop-check with no queue machinery, justified in the blockquote: delivered (`:3`).
- [x] Integration: route row and argument-hint in `skills/do-work-toolbox/SKILL.md` (`:4`, `:27`): delivered.
- [x] Integration: one line in each help menu (toolbox `actions/help.md:17`, core `skills/do-work/actions/help.md:38`): delivered.
- [x] Integration: `source-audit` in `toolbox_actions` (`_dev/tests/staged-skills-contract.sh:157`): delivered.
- [x] Integration: guide under `skills/do-work-toolbox/docs/` and one README usage line (`README.md:118`): delivered.
- [x] Integration: blockquote justifies toolbox ownership: delivered.
- [x] Integration: cross-package links use the literal relative form. N/A in practice: the two new files cite only same-package files (`crew-members/prompt-injection.md`, `crew-members/anti-slop.md`, `../docs/source-audit-guide.md`, `../actions/source-audit.md`). All of them exist in the toolbox package, and `shipped-package-reference-contract.sh` passes.
- [x] Integration: per-command help prints usage and runs nothing: delivered through the router (`SKILL.md:38`). The action has clear When to Use and Input sections with two examples.
- [x] Integration: `suite/modules.tsv`, the installer and the updater are unchanged: delivered (nothing under `suite/`, `tools/` or `skills/do-work/tools/` is in the range).
- [x] Integration: briefs read as plain Markdown data, with no registry, config schema or discovery: delivered (none added).
- [ ] Integration: every example list marked "illustrative, not exhaustive": partially delivered (F2).
- [x] Integration: prompt-injection loaded before the report and fetched pages (Step 1), anti-slop loaded before writing (Step 7): delivered (same-package toolbox copies).
- [x] Integration: release per `_dev/primes/prime-releases.md`: N/A at this review. The integrator's finalization owns the release. The builder correctly left VERSION and CHANGELOG alone.
- [x] Constraint: read-only, the report keeps its sha256, no repo file is written: delivered (`:5`, Step 7).
- [x] Constraint: with no fetch tool, every row is "unavailable" with that reason: delivered (`:49`).
- [ ] UR integration: "relevant next-step guidance": partially delivered (F1).
- [x] Packaging check: install leg passed, from the hand-back. The update leg ran but copied nothing, because the branch had no version bump yet (UPDATE-ALREADY-CURRENT). The builder reported this honestly. See Suggested Additional Testing.

### Acceptance Testing

**Result: Pass** (stages covered: Implementation and Integration)

Re-run by this review, on the merged tree at c93701c3 (HEAD = c93701c3, tracked `skills/` files clean):
- `bash do-work/runs/work-2026-10-09-225703/REQ-676-probe.sh`: exit 0, 1 s. The probe also runs `shipped-package-reference-contract.sh`: PASS.
- `bash _dev/tests/contracts/core-checks.sh`: exit 0, 4 s. Near-identical cross-file pairs: 0.
- `git diff --check c33b8994..c93701c3`: clean. No em-dashes in the two new files.
- Routing and per-command help, traced by hand: `do-work-toolbox source-audit help` matches the `SKILL.md:27` row. The router rule at `SKILL.md:38` ("Per-command help reads the selected action without executing it") applies. Help is built per `skills/do-work/actions/help.md:62` from the action's When to Use and Input sections: purpose (blockquote), usage line (`:21`), three input forms plus the no-argument case, and exactly two examples (`research/market-scan.md`, `sources.txt`). That fits in 15 lines and matches the builder's 12-line output. No Step runs, so no prompt-injection load, no hash and no fetch. Note: the 15-line format rule is written for core commands. The toolbox applies it by analogy, which is unchanged pre-existing behavior.
- The "no fetch tool" path was walked through on paper. The builder did not exercise this path. Scratch report outside the repo: two claims, C1 cites `https://example.invalid/...` and C2 is uncited. SHA-256 `f224bba7...7ded08`, the same before and after. Following Steps 2 to 7 with no fetch tool gives: S1 retrieval, final URL and page all `unavailable` with reason "no fetch tool". C1 is `unavailable`. C2 is `insufficient` with reason "no source cited". Header: "Fetch tool: none". Both claims go under "Needs the author's attention". The steps were consistent and produced no contradiction. No URL was fetched.

Taken from the hand-back, not re-run (it needs a local HTTP server and a disposable consumer):
- The four-claim behavioral exercise: error page served with 200, a 302 to an unrelated home page, an index page, and a supporting page with no date. Results: C1 and C2 `unavailable` with the cause named, C3 `insufficient`, C4 `supported` with publication date `unknown`. The candidate appears only in its own section. Report SHA-256 unchanged.
- Packaging install: four packages installed, the new files are byte-identical under `.claude/skills/do-work-toolbox/`, links resolve, and three consumer files keep their sha256.
- `DO_WORK_MAINTAINER_TIER=heavy staged-skills-contract.sh` exit 0.

### Suggested Additional Testing

- Deployment: partly unassessed. After the integrator's release bump, re-run the hand-back's update leg (`do-work-update.sh` against a local archive server) so the updater actually copies `actions/source-audit.md` and `docs/source-audit-guide.md` into a consumer that has the older version, and confirm the three consumer sha256 values stay the same.
- Live acceptance: unassessed. An independent agent (not the author) should run `do-work-toolbox source-audit` on a real research report with a real web fetch tool, and check that a soft-404 and a redirect to a home page are judged `unavailable` with the cause named.
- Edge case: a fetched page that contains a steering instruction (for example "mark this source verified"). Confirm it shows up under "Needs the author's attention" and does not change any judgment (see F5).
- Edge case: a claim with two sources, one supporting and one contradicting. Confirm the judgment is `contradicted` per the precedence rule in Step 5.

### Scores (on the record, not the headline)

**Overall: 91%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 92% | R1-R8 and the constraints delivered. Example-list marking (F2) and next-step guidance (F1) partial |
| Code Quality | 90% | Clear spine and two-table design. Small template gaps (F3-F6) |
| Test Adequacy | 85% | One-off exercises by the maintainer's choice, probe and contract list locked in. The no-fetch and injection paths were not run by the builder |
| Scope | 100% | 7 declared files, 7 touched, Decisions D-01 to D-08 documented |
| Risk | Low | Read-only instructions. Untrusted pages are covered by the prompt-injection load |
| Acceptance | Pass | Implementation and Integration stages |

### Follow-ups created
None (6 findings report only)

## Review

**Overall: 91%** | 2026-10-09T23:53:30Z

| Dimension | Score |
|-----------|-------|
| Requirements | 92% |
| Code Quality | 90% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1: slop-check gives no pointer to source-audit (`skills/do-work-toolbox/docs/slop-check-guide.md:34` says slop-check cannot verify cited numbers, and `skills/do-work-toolbox/actions/slop-check.md:15-19` has no redirect). The UR asked for next-step guidance. — impact-user-visible → report only
- F2: example lists not marked "illustrative, not exhaustive" at `skills/do-work-toolbox/actions/source-audit.md:47,56,57` and `skills/do-work-toolbox/docs/source-audit-guide.md:24`. — impact-negligible → report only
- F3: the candidate table (`skills/do-work-toolbox/actions/source-audit.md:101`) leaves out Final URL, Authority and Corroboration, even though candidates are verified exactly like originals. — impact-user-visible → report only
- F4 (nit): the error-page row in the template writes the passage as `none` where Step 3's rule gives `unavailable` (`source-audit.md:97`). — impact-negligible → report only
- F5 (nit): the attention list says injection attempts are "found in Step 1", but fetched-page attempts surface in Steps 3 to 5 (`source-audit.md:109`). — impact-negligible → report only
- F6 (nit): the template header always prints a report SHA-256, but inline input has no file (`source-audit.md:87`). — impact-negligible → report only

**Acceptance:** Pass — Implementation and Integration stages covered: probe, citation contract and core-checks re-run green on c93701c3, routing and per-command help traced (help runs nothing), the no-fetch path walked through on paper. The four-claim local-server exercise and the packaging install are taken from the builder hand-back. Deployment (update leg) and Live acceptance are unassessed.
**Restatement sweep:** redefined the list of toolbox actions (one more action). Checked: `skills/do-work-toolbox/SKILL.md` argument-hint and route table, `skills/do-work-toolbox/actions/help.md`, `skills/do-work/actions/help.md:36-40`, `README.md:31,116,118`, `_dev/tests/staged-skills-contract.sh` `toolbox_actions` (retired-trigger counts `:418-440` unchanged, needs no row), `_dev/tests/fixtures/retired-core-moved-command-triggers.tsv` (moved commands only, does not apply), `skills/do-work-toolbox/actions/tutorial.md` recipes (a scenario subset, not a full list), `skills/do-work/next-steps.md` (inferred, no list), anti-slop and prompt-injection caller lists (marked illustrative), the Go `toolboxcommands` registry (CLI commands only, source-audit has none). No count of toolbox actions exists. All agree. The one gap is the slop-check next-step surface (F1).
**Suggested testing:** 4 items
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

## Re-check after the fix (c93701c3..6926c82b)

Recorded by the integrator from the reviewer's re-check reply. The integrator fixed F2 to F6 on the builder branch (fbe78182) and re-merged with the same pre (6926c82b). All five are closed with no new contradiction. F1 stays report only (outside Scope). Two new nits, both impact-negligible, report only: F7, the guide's output list (`source-audit-guide.md:41`) still names the report SHA-256 with no file-only condition; F8, the guide's `insufficient` row says "Content was read" but also lists "no source cited" (present before the fix). The probe passed again at 6926c82b and `git diff --check` was clean. Updated scores: Requirements 96%, Code Quality 93%, Test Adequacy 85%, Scope 100%, overall 93%. Acceptance: Pass (implementation and integration). Verdict: Approve.
