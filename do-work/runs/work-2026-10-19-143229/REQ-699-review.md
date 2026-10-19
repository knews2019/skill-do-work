## Review: REQ-699

**Approve** — the revise form no longer forbids committing, both forms now state one commit rule that defers to the project's instructions, and the five "one expected link finding" sentences now match what the judge code does.
Route A | merged at 3fccc496 (range 6f6f8655..3fccc496, builder commit 07e075d1)

### What's built
- `ai-report.md` Step 8 holds the commit rule once (project's instructions decide; inside a run the run commits; never stage outside the bundle; never push). Revise step 7 applies it after the catalog regenerates and takes `catalog.json` and `index.html` with the bundle. The Step 8 summary names the commit hash or the exact commit command.
- Revise step 5, Step 7, Step 8 and two checklist lines now say the judge reports one `AI-REPORT-JUDGE-BROKEN-LINK` finding per `href` (an anchor link counts separately). This matches `report_judge.go:212-228`, which appends one finding per `[href]`/`[src]` attribute with status 400 or above and strips the fragment only for the fetch, with no de-duplication.
- `ai-report-guide.md:99` no longer says committing is left to the user.
- The release (changelog entry, version bump) is not in this range. That is expected: the finalize transaction in work-reference Step 9 owns it, not the builder.

### Decisions / risks for you
- None. D-01 (rule stated once in Step 8, applied by revise step 7) and D-02 (keep the "follow the project's instructions" form, not "commit by default") are recorded in the REQ and match the diff.

### Findings

**Important:**
- None

**Minor:**
- M1 `skills/do-work-toolbox/actions/ai-report.md:178` vs `:59` and `:180`. The Step 8 rule sentence ends its "otherwise" branch with "leave the commit to the user", while revise step 7 and the Step 8 summary paragraph end theirs with "print that exact commit command". The outcome is the same because the summary paragraph two lines later names the command, but a reader who stops at the rule sentence can take "leave it" as "print nothing". Replacing "leave the commit to the user" with "print that exact commit command" at `:178` would make the three sites identical. — impact-negligible → report only

**Nit:**
- N1 `skills/do-work-toolbox/actions/ai-report.md:180`. In "either the commit hash or the exact commit command (...) when the project leaves the commit to the user", the "when" clause can be read as qualifying both alternatives. Moving "when the project leaves the commit to the user" directly after "command" (before the parenthesis) removes the ambiguity. — impact-negligible → report only

### Requirements Checklist

- [x] Finding 1: revise step 7 "Do not commit" replaced by "Commit as the project's rules say" with the commit-or-print rule, run ownership, catalog pair by exact path, never push — delivered (`ai-report.md:59`)
- [x] Finding 1: guide sentence "committing is left to you" replaced — delivered (`ai-report-guide.md:99`)
- [x] Finding 2: Step 8 states the same rule for the creation path and the summary names the hash or the exact command — delivered (`ai-report.md:178,180`)
- [x] Finding 3: singular "one link" wording fixed at revise step 5, Step 7, Step 8, and both checklist lines — delivered (`:57,170,176,200,203`); the REQ named five sites and five changed
- [x] Finding 4: checklist revise line names where the catalog pair ends up — delivered (`:203`)
- [x] Constraints: prior bundle bytes unchanged (`:57` still ends "The prior bundle's bytes never change"); catalog still regenerated in revise step 6 (`:58` untouched); run owns commits inside a run (`:59`, `:178`); "never push" at both sites; no Go file in the diff
- [x] Shipped file never cites this repo's own `CLAUDE.md`: `:178` says "its `CLAUDE.md` or similar", the consumer's own file, same form as `skills/do-work/actions/capture.md:102`
- [x] Cross-package citations: the diff adds no path citation; `ai-reports/catalog.json` and `ai-reports/index.html` are consumer-project paths, not package citations
- [x] P-A-U: all three boxes `[x]`
- [ ] Release entry — N/A for this range; owned by finalization (work-reference Step 9 finalize transaction)

### Acceptance Testing

**Result: Pass** (implementation stage from my own greps; integration stage from the builder's recorded green runs of `_dev/tests/contract-regressions.sh` and `_dev/tests/shipped-package-reference-contract.sh`, not rerun here per the orchestrator; the full gate runs separately)
- `grep -n "Do not commit"` on `ai-report.md`: no output.
- `grep -n "Commit as the project's rules say"`: lines 59 and 178 (the REQ's GREEN allows the builder's restructuring and the Qualification records both lines).
- `grep -n "one finding per"`: line 57.
- `grep -n -i "link finding remains\|only finding is the expected link\|only the expected prior-bundle link finding"`: no output.
- `grep -n "left to you"` on `ai-report-guide.md`: no output.
- Read `report_judge.go:200-230`: one finding appended per attribute with status 400 or above, fragment stripped before fetch, no de-duplication. The new prose is correct.
- Restatement sweep over `skills/` and `_dev/primes` for "do not commit", "committing stays", "left to you", "left to the user", "prior-bundle link finding", "that one link", "expected link to the prior", "this finding is expected": only the two rewritten sites in `ai-report.md` and one unrelated hit in `capture.md:108` (operator steps, not ai-report). `help.md:17` and `ai-report-reference.md` say nothing about committing. The CHANGELOG is historical and not flagged.

### Suggested Additional Testing

- Live acceptance: unassessed. In a consumer whose instructions ask for completed increments to be committed, run a direct `ai-report revise` and confirm the summary ends with a commit hash and `git status` shows neither the new bundle nor `ai-reports/catalog.json` and `ai-reports/index.html` as untracked. This repository cannot run that.
- Deployment: unassessed. After finalization, confirm the changelog entry title says what shipped and the version mirrors match, per `_dev/primes/prime-releases.md`.

### Scores (on the record — not the headline)

**Overall: 98%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All four findings and every constraint delivered; release is finalization's |
| Code Quality | 95% | One "otherwise" branch worded differently from its two siblings (M1) |
| Test Adequacy | N/A | Prose-only Route A change; acceptance greps and the two contract scripts are the proof |
| Scope | 100% | Two files, both in `write_set`, +10/-8, every hunk traces to a finding |
| Risk | None | Prose in a toolbox action; no code consumer changed |
| Acceptance | Pass | Implementation and integration stages; live acceptance and deployment unassessed |

### Follow-ups created
- None (2 findings report only)

---

## Proposed `## Review` block for the REQ file

```markdown
## Review

**Overall: 98%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** M1 `skills/do-work-toolbox/actions/ai-report.md:178` ends its "otherwise" branch with "leave the commit to the user" while `:59` and `:180` end theirs with "print that exact commit command"; same outcome, different words — impact-negligible → report only. N1 (nit) `:180` "when the project leaves the commit to the user" reads as qualifying both the hash and the command — impact-negligible → report only
**Acceptance:** Pass — implementation stage by the REQ's acceptance greps, integration stage by the builder's green `contract-regressions.sh` and `shipped-package-reference-contract.sh`; live acceptance (consumer-side direct revise ends with a hash) and deployment (finalization release) unassessed
**Restatement sweep:** redefined the ai-report commit rule (revise step 7, Step 8 rule and summary) and the expected-judge-finding count (one per `href`); grepped `skills/` and `_dev/primes` for the old phrasings; no stale restatement outside the diff; CHANGELOG historical, not flagged
**Suggested testing:** 2 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*
```

## Re-check

Fix range `3fccc496..88f8c686` (builder commit e97d4aa6, merge 88f8c686). Cumulative range `6f6f8655..88f8c686` is still two files, +10/-8.

- M1 resolved. `ai-report.md:178` now ends its "otherwise" branch with "and report the hash; otherwise print that exact commit command", the same words as revise step 7 at `:59`.
- N1 resolved. `ai-report.md:180` now reads "and the commit hash, or, when the project leaves the commit to the user, the exact commit command (...)", so the user-owned case scopes only the command.
- Nothing else changed: the fix diff touches only those two lines, and the cumulative diff against `6f6f8655` has the same two files and the same line counts as the first review.

**Overall: 100%** | <timestamp>. Code Quality 95% rises to 100%; Requirements 100%, Scope 100%, Test Adequacy N/A, Risk None, Acceptance Pass (same stages as above). Verdict stays Approve. Follow-ups created: None (0 findings open).
