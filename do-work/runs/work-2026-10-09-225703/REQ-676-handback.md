# REQ-676 hand-back: do-work-toolbox source-audit

- Branch: `worktree-agent-REQ-676-source-audit-action`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-676-source-audit-action`
- Base: `e313e870`. One commit: `bf37f4ff` `[REQ-676] add do-work-toolbox source-audit action and guide`.
- Worktree clean after the commit. Nothing under `do-work/` touched in the worktree.

## File manifest

- `skills/do-work-toolbox/actions/source-audit.md` (new): the action. Blockquote links the guide and justifies toolbox ownership. When to Use (redirects to `slop-check` and `validate-feedback`), Input (usage line plus two examples), seven numbered Steps, Output Format. No Rules, Common Rationalizations, Red Flags or Verification Checklist (none earned).
- `skills/do-work-toolbox/docs/source-audit-guide.md` (new): short user guide with a "Not to be confused with slop-check" note, judgment table, input block, output order, key rules.
- `skills/do-work-toolbox/SKILL.md` (modified): `source-audit` added to `argument-hint` after `slop-check`; one route row after `slop-check` (`source-audit`, `audit sources`, `check citations`).
- `skills/do-work-toolbox/actions/help.md` (modified): one menu line after `slop-check`, description at the existing column.
- `skills/do-work/actions/help.md` (modified): ` · source-audit` appended to the `present-video · slop-check` row of the toolbox list.
- `README.md` (modified): one new paragraph after the "Common extension calls also include ..." line.
- `_dev/tests/staged-skills-contract.sh` (modified): `source-audit` in `toolbox_actions` after `slop-check`.

## P-A-U

- [x] **[PLAN]:** One action file whose spine is REQ requirements 1-7; four facts kept apart structurally by using two tables (claim table carries judgment and corroboration; source evidence table carries retrieval, page check, passage, date, authority). Error page and unrelated redirect classify as `unavailable` (the cited content was not read); index page as `insufficient`. Judgment precedence for multi-source claims fixed in one sentence. Help served by the router from When to Use and Input, no help code. Integration edits one line or row each, placed next to `slop-check`.
- [x] **[APPLY]:** Seven files exactly as the Scope list. One correction during apply: example path `docs/research/market-scan.md` was read by the staged-skills runtime-reference scan as a shipped `docs/` citation and failed; changed to `research/market-scan.md` in both files.
- [x] **[UNIFY]:** `git diff --stat e313e870..HEAD`: 7 files changed, 166 insertions(+), 2 deletions(-). `git diff --check`: clean (exit 0). No em-dashes in the two new files (grep). Checks on the final tree:

| Check | Exit | Wall |
|---|---|---|
| `REQ-676-probe.sh` from the worktree root | 0 | 1 s |
| `_dev/tests/shipped-package-reference-contract.sh` | 0 | 1 s |
| `_dev/tests/contracts/core-checks.sh` (runs standalone) | 0 | 6 s |
| `DO_WORK_MAINTAINER_TIER=heavy _dev/tests/staged-skills-contract.sh` | 0 | 32 s |
| `_dev/tests/contract-regressions.sh` | 1 | 32 s |

contract-regressions exit 1 is base-caused, not from this branch: its only FAIL lines are the quiet-grep pipeline audit flagging `do-work/runs/work-2026-10-09-225703/REQ-676-probe.sh:7` and `REQ-677-probe.sh:9` (`sed -n '/^argument-hint:/p' ... | grep -q ...`). Those files arrived in base commit `e313e870` (`[UR-151] run artifacts`). Reproduced on a `git clone --shared` of base `e313e870` in a temp dir: `quiet-grep-pipeline-audit.sh` reports the same 2 failures. Every other contract file in the runner passed (it accumulates failures and does not stop early). See Discovered Tasks.

Files reviewed: all seven above, plus the installed copies of the two new files (byte-identical to source).

## Behavioral exercise

Fixtures in a `mktemp -d` dir outside both trees (deleted afterwards). Server: a small `http.server.BaseHTTPRequestHandler` subclass on 127.0.0.1 port 51247, no network. `/a` 200 with body "Page not found"; `/b` 302 to `/home` ("Welcome to Garden Tools Co."); `/c` a "Lumen Labs news" link list; `/d` "In our tests the Lumen 3 battery lasted 12 hours on its lowest setting." with no date. Report `lumen-report.md`: four claims, numbered citations to `/a` to `/d`.

- SHA-256 before: `e5bdee3fbdbdd9bcddeef894ee694b75e57340a348a8d36a9dda1faa55e35c43`. After the audit: identical.
- Routing: `do-work-toolbox source-audit <report>` matches `SKILL.md:27` route row to `./actions/source-audit.md`. RED: on base `e313e870`, `grep -c source-audit skills/do-work-toolbox/SKILL.md` = 0 and the action file is absent; under `SKILL.md:38` an unknown single word prints help, so the request does not route.
- Followed the action step by step. Step 1 loaded `crew-members/prompt-injection.md`. Fetch tool: `curl -sS -L -D - -w '%{url_effective} %{http_code} %{num_redirects}'`. Step 7 loaded `crew-members/anti-slop.md`.

Output produced:

```
# Source Audit: lumen-report.md

4 claims: 1 supported, 0 contradicted, 1 insufficient, 2 unavailable.
Fetch tool: curl. Report SHA-256: e5bdee3f...5e35c43, unchanged.

## Claim-to-source table

| Claim | Judgment | Sources | Corroboration | Why |
|---|---|---|---|---|
| C1 "draws 4 watts at full brightness" (line 3) | unavailable | S1 | none | S1 is an error page served with 200 ("Page not found") |
| C2 "founded in 2019" (line 5) | unavailable | S2 | none | S2 redirects (302) to /home, an unrelated Garden Tools Co. home page |
| C3 "two-year warranty" (line 7) | insufficient | S3 | none | S3 is a news index page; no item mentions a warranty |
| C4 "battery lasts 12 hours on its lowest setting" (line 9) | supported | S4 | single source | S4 states it directly |

| Source | Requested URL | Final URL | Retrieval (every attempt) | Page | Passage | Published | Authority |
|---|---|---|---|---|---|---|---|
| S1 | http://127.0.0.1:51247/a | http://127.0.0.1:51247/a | 200 | error page served with 200 | none | unknown | unknown |
| S2 | http://127.0.0.1:51247/b | http://127.0.0.1:51247/home | 302 -> /home, 200 | unrelated redirect (home page) | none | unknown | unknown |
| S3 | http://127.0.0.1:51247/c | http://127.0.0.1:51247/c | 200 | index or listing page | none found | unknown | Lumen Labs (publisher's own news list) |
| S4 | http://127.0.0.1:51247/d | http://127.0.0.1:51247/d | 200 | content page | "In our tests the Lumen 3 battery lasted 12 hours on its lowest setting." | unknown | Lumen Labs test team (manufacturer's own test) |

## Replacement candidates

| Candidate | For claim | URL | Retrieval | Page | Passage | Published | Would judge |
|---|---|---|---|---|---|---|---|
| R1 | C3 | http://127.0.0.1:51247/news/3 ("Lumen 3 preview", linked from S3) | 404 | error page | none | unknown | unavailable |

No candidate sought for C1 or C2: no search tool for this site and no archive reachable offline.

## Needs the author's attention

- C1: the cited page is gone (error page served with 200); find a working source for the 4 W figure.
- C2: the citation redirects to an unrelated home page; replace it.
- C3: the citation is an index page; cite the article that states the warranty. The linked preview returned 404.
```

GREEN against the REQ: A (C1) and B (C2) not supported, error page and redirect named; C (C3) insufficient; D (C4) supported with publication date `unknown`; the only replacement candidate sits in its own section and did not change C3's row; report SHA-256 unchanged.

`do-work-toolbox source-audit help` under the router rule (`SKILL.md:38`, format `skills/do-work/actions/help.md:62`), built from When to Use and Input, 12 lines, no fetch run:

```
source-audit: checks that the sources a report cites say what it claims; prints a claim-to-source table and never edits the report.
Usage: do-work-toolbox source-audit <report-or-url-list>
Arguments:
  <report file>   Markdown, plain text or HTML report with citations
  <URL list>      file or inline list, one URL per line, optional claim after each
  <pasted text>   report prose, audited as given
  (none)          asks for the report path or URL list
Use when: a report cites sources; you have URLs to check; research is about to be published or acted on.
Not for: prose quality (slop-check), reviewer feedback (validate-feedback), correcting the report.
Examples:
  do-work-toolbox source-audit research/market-scan.md
  do-work-toolbox source-audit sources.txt
```

Not exercised: the "no fetch tool" path (curl was available) and an injection attempt inside a fetched page. Both are prose rules in Steps 1 and 3 only. The exercise was run by the same agent that wrote the action, so it shows the procedure is followable, not that an independent agent reads it the same way.

## Packaging check

Run after the last commit against branch tip `bf37f4ffe5d3ae5d5cdf5e1baf79f25abb1a877e`, in a `mktemp -d` dir (deleted afterwards).

1. `git archive --format=tar.gz --prefix=skill-do-work-main/ worktree-agent-REQ-676-source-audit-action`: ok.
2. Consumer seeded and committed. SHA-256 before: `docs/project-brief.md` 8491f671147b37e73b1d78c332387dbfc5fae915f392f414935014999fb07c49; `do-work/queue/REQ-001-sample.md` ceb101922ebc17570577cd8800814340a0bc967b3d3f55c06ca127d593117375; `src/app.js` f9444510dc7403e41049deb133f6892aa6a63c05591b2b59e4ee5b234d7bbd99.
3. Install (`tools/install-do-work-suite.sh --archive`): exit 0, four packages created at v0.305.89, `rollback: not_needed`. Three SHA-256 values unchanged. Install committed in the consumer.
4. Update (`do-work-update.sh` with `DO_WORK_UPSTREAM_URL=http://127.0.0.1:51321/archive/refs/heads/main.tar.gz`): exit 0, the server logged `GET /archive/refs/heads/main.tar.gz 200`, then `skipped UPDATE-ALREADY-CURRENT: the installed suite is already v0.305.89`, `rollback: not_needed`. Three SHA-256 values unchanged. Honest partial: the branch has no version bump (release is the integrator's), so the updater fetched but did not reconcile files. I stopped there as the brief says. The update leg that copies files is exercised only after the integrator's release bump.
5. Found under `.claude/skills/do-work-toolbox/`: `actions/source-audit.md` (8058 bytes, byte-identical to source), `docs/source-audit-guide.md` (3199 bytes); installed `SKILL.md` has the argument-hint entry and the route row. Links from each file's own directory: `actions/` -> `../docs/source-audit-guide.md` ok; `docs/` -> `../actions/source-audit.md` ok. Same-package citations from the package root: `crew-members/prompt-injection.md`, `crew-members/anti-slop.md`, `docs/source-audit-guide.md`, `actions/source-audit.md` all resolve.

Servers stopped; no listeners left on either port.

## Decisions

- D-01 DECIDE & STATE: two tables in the claim-to-source section (one row per claim with the judgment, one row per source with the evidence). This keeps one judgment per claim and makes the four facts separate columns instead of prose. Builder Guidance gave latitude on columns.
- D-02 DECIDE & STATE: error page and unrelated redirect judge `unavailable` (the cited content was never read); an index page judges `insufficient` (the cited URL is what loaded, it just does not state the claim). Matches the REQ's GREEN wording.
- D-03 DECIDE & STATE: multi-source precedence is contradicted > supported > insufficient > unavailable, so a single contradicting source is never hidden by a supporting one.
- D-04 DECIDE & STATE: uncited claims and URL-list entries with no claim are `insufficient` with the reason named, so every row still has one of the four judgments.
- D-05 DECIDE & STATE: publication date excludes server headers and archive capture dates as well as fetch date and URL; an archive snapshot is a likely replacement candidate, and its capture date is the easy wrong answer.
- D-06 DECIDE & STATE: the audit is printed only; no saved-file option (the REQ says the action writes no repo file).
- D-07 DECIDE & STATE: README gets a new one-sentence paragraph after the "Common extension calls" line rather than editing that line, to keep the edit local. It still sits next to REQ-677's likely edit, so expect a textual conflict there.
- D-08 DECIDE & STATE: did not add the guide to `toolbox_files` (optional per Exploration; one fewer seam line).

## Discovered Tasks

- impact-user-visible: base commit `e313e870` (`[UR-151] run artifacts`) makes `_dev/tests/contract-regressions.sh` exit 1, because `do-work/runs/work-2026-10-09-225703/REQ-676-probe.sh:7` and `REQ-677-probe.sh:9` pipe `sed` into `grep -q`, which `_dev/tests/quiet-grep-pipeline-audit.sh` refuses. Fix in the orchestrator's run artifacts: `grep -q 'source-audit' <(sed -n '/^argument-hint:/p' skills/do-work-toolbox/SKILL.md)` or `grep -Eq '^argument-hint:.*source-audit' skills/do-work-toolbox/SKILL.md`. The integrator's gate will hit this on every branch of the wave. → report only
- The staged-skills runtime-reference scan treats any `docs/...` token in a shipped file as a package citation, including example user paths inside a usage example. Not a bug to fix here; it is the lesson below. → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole; families canonical-link-outlives-its-target, manifest-ownership-vs-edit-content).
- `_dev/primes/lessons-action-files.md` § Template (lines 76-143) and every `[family: alternate-writer-contract-drift]` bullet (REQ-477, 498, 513, 461, 531, 566, 640, 641, 642, 644, 648, 647, 652). Applied: swept every reader of the toolbox command list named in Exploration (router, two menus, README, contract list).
- Primes: `prime-action-files.md`, `prime-releases.md`. Crew members: general, coding-guardrails, shared-principles, communication-style, anti-slop (toolbox copy), prompt-injection (toolbox copy).

## Proposed CHANGELOG entry

**Source Audit for Research Reports**

`do-work-toolbox source-audit <report-or-url-list>` checks whether the pages a report cites actually say what it claims. A link that loads is not proof: error pages served with success, redirects to a home page and index pages now show up as what they are.

- New read-only toolbox action `source-audit` with a user guide (`docs/source-audit-guide.md`).
- Returns a claim-to-source table with one judgment per claim (supported, contradicted, insufficient, unavailable), the evidence for every source, replacement candidates in their own section, and the claims that need the author's attention.
- Keeps retrieval, support, corroboration and authority apart; an unstated publication date stays "unknown"; the report is never edited.
- Listed in the toolbox router, both help menus and the README.

## Proposed lesson bullet

Satellite: `_dev/primes/lessons-action-files.md`

- [family: example-path-read-as-citation] [REQ-676: a usage example's user path that starts with a shipped-package directory name (`docs/research/market-scan.md`) is read by the staged-skills runtime-reference scan as a citation and fails as unresolved; pick example paths whose first segment is not `actions`, `docs`, `tools`, `hooks`, `crew-members` or `specs`](../../do-work/archive/UR-151/REQ-676-source-audit-action.md#lessons-learned)

## Integration seams

REQ-677 (journey-qa) edits the same lines: `SKILL.md:4` argument-hint (certain conflict; keep both names), the route table (my row is after `slop-check`, theirs is expected after `ui-review`), toolbox `actions/help.md` (my line after `slop-check`), core `actions/help.md` (I appended to the `present-video · slop-check` row; if REQ-677 also appends there, keep both and rewrap if over width), `README.md` (my new paragraph after line 116), and `toolbox_actions` (my entry after `slop-check`).

Test wall times: probe 1 s, shipped-package-reference 1 s, core-checks 6 s, staged-skills heavy 32 s (58 s on the first, failing run), contract-regressions 31-32 s.
