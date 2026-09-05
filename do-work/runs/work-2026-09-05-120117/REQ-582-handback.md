# REQ-582 hand-back — Detect the Arrow-Form Section Citation in the Shipped-Package Reference Contract

## Branch

- Branch: `worktree-agent-REQ-582-arrow-citations`
- Head commit: `87956ed1`
- Base commit: `961697bb`
- Two commits: `04f02f10` (the check, the fixtures, the two live fixes) and `87956ed1` (one false-positive shape found by the sweep, narrowed and pinned).

## File manifest

| Verb | Path | What |
|---|---|---|
| modify | `_dev/tests/shipped-package-reference-contract.sh` | Resolve the section name in the `` `path.md` `` → **Named Section** form; new `run_section_citation_fixtures()` suite beside the anchor suites |
| modify | `skills/do-work/actions/cleanup.md` | Line 48: the dangling `**In-Progress Record (Step 2)**` now names `**In-Progress Record (Step 1)**` |
| modify | `skills/do-work/CHANGELOG.md` | Line 392 (0.266.5 entry): keeps its historical section name as prose and gains a live pointer to the section that content now lives under |

Tests touched: the same file. The new coverage is in-script (`run_section_citation_fixtures`), as the REQ asked — no new test file.

Working tree is clean at the head commit. `CHANGELOG.md` (root) is unmodified and is handed back as an integration seam below.

## P-A-U

**[PLAN]**

1. Sweep every arrow-form citation in the four shipped packages and classify what each one names, so the "heading only or also bold label" decision is made against the corpus rather than against what is easy to implement.
2. Parse the section name out of the raw text at the position `citation_candidate_tokens` already yields, so the new check rides the existing citation surface instead of adding a second scanner.
3. Have `cross_package_target_resolves` and `same_package_target_resolves` return the source target they already compute, so the section check runs at the point where the path is known to resolve — no second resolution path to drift.
4. Match a cited name against the target's ATX headings and bold runs.
5. Add an in-script fixture suite for both halves (what is read as a name, what a name may resolve against).
6. Fix the two live dangling citations; hand back the changelog mirror line.

**[APPLY]** Coded as planned, inside the declared write set only. One addition the plan did not foresee, made because the sweep found it: the bold name may wrap across a line break (three shipped citations do), and a bold phrase that closes *before* the arrow must not be read as a section name (one shipped line does that).

**[UNIFY]**

- `git diff --stat` against the base: `_dev/tests/shipped-package-reference-contract.sh` +284/-8, `skills/do-work/CHANGELOG.md` 1 line, `skills/do-work/actions/cleanup.md` 1 line.
- `bash -n _dev/tests/shipped-package-reference-contract.sh` — exit 0.
- `shellcheck -S style _dev/tests/shipped-package-reference-contract.sh` — clean, no output.
- `bash _dev/tests/action-shell-blocks.sh` — exit 0, "73 fenced blocks and 33 shipped shell files; ShellCheck enabled."
- No Go and no client JavaScript changed, so `gofmt`/`go vet`/`node --check` do not apply.
- Debug-artifact scan over added lines (`git diff -U0 | grep '^+' | grep -i 'print(|breakpoint|TODO|FIXME|pdb'`) — no hits. The only `print` calls in the file are the pre-existing PASS/FAIL reporting.
- File-by-file review: the test script (new helpers are additive; the two resolver return-tuples grew by one element and both are consumed positionally only inside `citation_messages`; the `same_package_check_enabled` gate moved past the success branch so a package-root CHANGELOG.md still reports no path failures); `cleanup.md` (one digit, no other text); `skills/do-work/CHANGELOG.md` (one line, historical name preserved).

## Test evidence

| # | Command | Exit | Observation |
|---|---|---|---|
| 1 | `bash _dev/tests/shipped-package-reference-contract.sh` (base tree) | 0 | **RED.** `shipped package reference contract: PASS` while two live citations name sections that do not exist. The failure is the silence: the script contains no arrow character and no bold parsing at all. Runtime 2.1 s. |
| 2 | same, after the check was added, before the two prose fixes | 1 | **RED, now visible.** Exactly two findings, no others: `skills/do-work/CHANGELOG.md:392: cited section is not a heading or bold label in skills/do-work/actions/work-reference.md: Recovery Refusals (Step 1)` and `skills/do-work/actions/cleanup.md:48: cited section is not a heading or bold label in skills/do-work/actions/work-reference.md: In-Progress Record (Step 2)`. |
| 3 | same, after both prose fixes, mirror seam applied | 0 | **GREEN.** `shipped package reference contract: PASS`. Runtime 1.6 s — same range as the base. |
| 4 | same, after both prose fixes, mirror seam **not** applied | 1 | `0 broken reference(s), changelog mirror differs`. Zero citation findings; the only failure is the byte-identical mirror, which the seam below closes. |
| 5 | planted third dangler (`` `actions/work-reference.md` `` → **No Such Section Anywhere** appended to `cleanup.md`, then reverted) | 1 | `skills/do-work/actions/cleanup.md:336: cited section is not a heading or bold label in skills/do-work/actions/work-reference.md: No Such Section Anywhere` — the REQ's third GREEN condition. |
| 6 | negative control on the wrapped form (`**Worktree Dispatch Mode\n(Step 1)**` → `(Step 9)` in `crew-members/background-agents.md`, then reverted) | 1 | Reported at `background-agents.md:187`. Proves the three hard-wrapped citations are genuinely checked and not silently skipped. |
| 7 | `bash _dev/tests/action-shell-blocks.sh` | 0 | ShellCheck lint over shipped shell still passes. |

Runs 3 and 4 differ only in whether root `CHANGELOG.md` carries the seam. I verified GREEN by copying `skills/do-work/CHANGELOG.md` over root `CHANGELOG.md`, running, then restoring root byte-for-byte — root `CHANGELOG.md` is unmodified in the committed tree and in the working tree (`git status` clean). Every temporary probe in runs 5 and 6 was reverted the same way and confirmed with `git status --porcelain`.

The existing `#fragment` anchor resolution is untouched: `run_anchor_slug_fixtures` and `run_anchor_topology_fixtures` run before the new suite and pass unchanged in every run above.

I did **not** run `_dev/tests/maintainer-verify.sh`, per the brief.

## Lesson evidence

Read:

- `_dev/primes/prime-shell-commands.md` (the REQ's own prime, in full). None of its traps are touched by this change — no git plumbing, no `curl`, no signal handling, no exit-status-shaped lane, no arbitrary-file walk. Its **Closed Enumerations Go Stale** rule is why the new code keys on the *condition* (a delimited name after the arrow, resolved against whatever the target declares) instead of a list of section names or of files.
- `_dev/primes/lessons-shell-commands.md` — the satellite named in the REQ's "Required Lessons — Dropped for Budget". Its one directly relevant entry is REQ-312 (the previous widening of this same checker): *replay the original silent mutations, pair source and installed topology in the resolver, and let the whole-repo gate expose stale assertions outside the focused checker.* Applied: runs 5 and 6 replay the silent mutation rather than trusting the untouched corpus, and I grepped for every other file that asserts on this checker (`_dev/tests/contracts/probe-lanes.sh` is the only one and it keys on exit status alone, so there is no stale assertion outside the focused test — but see the orchestrator note below, since I may not run the gate).
- `_dev/primes/prime-releases.md` — the changelog house rules, read before touching `skills/do-work/CHANGELOG.md`. Its mirror rule ("the installed changelog mirror is byte-identical") is what produces the integration seam below.

No listed path was missing.

## Decisions

**D-01 — A cited "section" is an ATX heading *or* a bold label the target declares. ESCALATE.**

This is the decision the request asked to settle, and I settled it against the 75 live arrow-form citations rather than against what is cheapest to implement. Of them, 60 name an ATX heading exactly. The rest do not, and they are not sloppy:

- `actions/work-reference.md` → **Frontmatter Quoting** names `**Named contract — Frontmatter Quoting.**`, a bold paragraph lead. The file itself calls it a "Named contract", and five other lines in that same file cite it as "the **Frontmatter Quoting** contract above".
- `actions/work-reference.md` → **Stakeholder REQ terminal semantics** names `> **Named entry point — Stakeholder REQ terminal semantics.**`, whose own line says "`actions/clarify.md` Step 5.5's reclaim branch cites this by name".
- `actions/capture-reference.md` → **Populating `write_set`** and → **Populating `depends_on`** name bold paragraph leads.
- `crew-members/background-agents.md` → **Worktree isolation is a separate axis** names a bold paragraph lead.

A heading-only rule would report all five of those correct, deliberately-declared citation targets as broken, and the REQ's own constraint is that the fix widens what the checker sees rather than shrinking what shipped prose may write. So both forms count. **Value:** the check answers honestly for the whole live corpus instead of for 60 of 75 citations. **Risk:** bold is also ordinary emphasis, so a name that merely happens to be bolded somewhere in the target resolves. That is a false *pass*, never a false failure, and it is bounded by what the rule is actually for — a reader must be able to find the named thing in the target file. Reversible in one function (`section_names_from_text`).

**D-02 — Matching is whole-word containment, not equality. DECIDE & STATE.**

Two live shapes need it: a citation may name a heading without its parenthetical qualifier (`**Hook Install Internals**` for `## Hook Install Internals (used by actions/setup-memory.md → memory-module)`; `**Worktree Dispatch Mode**` for `## Worktree Dispatch Mode (Step 1)`, four sites), and a label without its declaration prefix (`**Frontmatter Quoting**` for `**Named contract — Frontmatter Quoting.**`). A reader follows both. Equality would report seven correct citations. Containment still catches what this check exists for: `**In-Progress Record (Step 2)**` is not contained in `In-Progress Record (Step 1)`.

**D-03 — Only the delimited bold form is read. DECIDE & STATE.**

45 live citations write an undelimited name (`` `actions/work-reference.md` `` → Request File Schema). There is no closing marker, so where the name stops and the sentence resumes is a guess, and a guess reports ordinary prose as a broken pointer. Those stay unchecked, which is the pre-existing state, not a regression. Reported as a discovered task (DT-1) rather than solved here.

**D-04 — The section half is checked inside a package-root CHANGELOG.md, where the path half is not. DECIDE & STATE.**

`should_check_same_package` switches off same-package *path* resolution for a package-root `CHANGELOG.md`, because a historical entry may name a file that has since gone. That reason does not reach a citation whose file is still present and only whose section name is stale, and the REQ requires the changelog citation to be reported. Implemented by moving the `same_package_check_enabled` gate to *after* the resolution success branch, so a changelog citation whose path does not resolve still produces nothing. The existing `package_changelog_messages` / `nested_changelog_messages` fixtures cover that and still pass.

**D-05 — Changelog remedy: keep the historical name, add a live pointer. Do not retarget. ESCALATE.**

The 0.266.5 entry announced a section that 0.266.6 renamed the same day. Retargeting the bold name at today's heading would make 0.266.5 claim it shipped a section under a name it did not have, and would make 0.266.6's "section renamed to **Stuck Runs Hand Off to Judgment (any step)**" read as a rename to the name it already had. So the line now reads:

> New `actions/work-reference.md` section **Recovery Refusals (Step 1)**, renamed in 0.266.6 and now `actions/work-reference.md` → **Stuck Runs Hand Off to Judgment (any step)**: judge each blocked path, …

The historical name survives as prose (no arrow, so it asserts no live pointer, which is exactly what it is — a past name), and the entry gains one machine-checked citation to where that content lives now. **Value:** history stays true and the reader can still follow the pointer. **Risk:** it edits a shipped historical entry at all; the alternative readings are (a) retarget, which falsifies the record, and (b) leave it and exempt changelogs from section checking, which drops one of the REQ's two named fixes. Reversible — one line, both copies.

**D-06 — The bold name may wrap across a line break, but not across a blank line. DECIDE & STATE.**

Three shipped citations write `**Worktree Dispatch Mode\n(Step 1)**` (`crew-members/background-agents.md:187` in `do-work` and its two mirrored copies). Stopping at the newline would leave exactly the citations a hard-wrapped file writes unchecked — the failure family this REQ closes. Whitespace is collapsed before matching, and the blank-line bound keeps an unclosed `**` inside its own paragraph. Run 6 is the negative control.

**D-07 — Only closing punctuation that cannot itself be emphasis may sit between the path and the arrow. DECIDE & STATE.**

`actions/abandon.md:39` writes ``- `status: failed` at `do-work/archive/legacy/`** → **cancellable in place.**`` — the arrow means "becomes", and the `**` before it closes a phrase. My first pass stepped over that `**` and read "cancellable in place." as a cited section. It caused no failure today only because the token is a consumer-queue `do-work/` path the checker skips for unrelated reasons. Dropping `*` and `_` from the leading punctuation class removes the whole class; a fixture pins the live shape. This is commit `87956ed1`.

## Discovered Tasks

- **DT-1 — 45 live arrow citations write an undelimited section name and stay unchecked** (`` `actions/work-reference.md` `` → Request File Schema, and 44 like it). Closing this needs a different mechanism from delimiter parsing: match the text after the arrow against the target's known section names longest-first, and report only when no known name starts there. That is a real design with real false-positive risk, not a tweak. `impact-rule-change` → report only.
- **DT-2 — `skills/do-work-toolbox/actions/code-review.md:327` cites `.claude/skills/do-work/actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**, and no part of this checker sees it.** The `.claude/` prefix is an installed-topology root the citation surface deliberately treats as "not from here" (pinned by `run_citation_surface_fixtures`), so neither its path nor its section is resolved. The cited section does exist (`work-reference.md:633`), so nothing is broken today; the point is that this one file addresses a sibling package through the installed tree while every other cross-package citation uses `../../`. `impact-negligible` → report only.
- **DT-3 — five live citations name a bold contract label rather than a heading, and three name a heading minus its parenthetical.** They are all followable and all pass under D-01/D-02, so this is not a defect list. It is worth knowing that "cite the bold Named contract line" is now load-carrying prose in this repo: `**Named contract — …**` and `**Named entry point — …**` are declaration markers other files depend on. If the suite ever wants section citations to be exactly resolvable, promoting those labels to real headings is the cheaper half of that job. `impact-negligible` → report only.
- **DT-4 — the arrow is heavily overloaded in shipped prose.** Beyond the 130 occurrences whose left side is a `.md` path, roughly 150 more use `→` for alias mappings (`` `back_end` `` → `` `backend` ``), state transitions (`pending` → `claimed`), pipeline steps, and table cells. Any future work that widens arrow parsing further must key on the left side being a resolvable path, exactly as this change does. `impact-negligible` → report only.

## Integration seams

One. The root `CHANGELOG.md` must stay byte-identical to `skills/do-work/CHANGELOG.md` — the very check this REQ changes enforces it — and root `CHANGELOG.md` is not in my write set, so I did not touch it.

**File:** `CHANGELOG.md` (repository root), **line 392.**

Replace this exact line:

```
- New `actions/work-reference.md` → **Recovery Refusals (Step 1)**: judge each blocked path, take the least destructive clearing action (delete or locally exclude non-project files, revert or finish this session's own abandoned write, commit the user's unrelated work on its own with the hash reported), re-run the exact `verification_argv`, and stop only for shared state whose owner the orchestrator cannot decide, naming the verb that resolves it.
```

with this exact line:

```
- New `actions/work-reference.md` section **Recovery Refusals (Step 1)**, renamed in 0.266.6 and now `actions/work-reference.md` → **Stuck Runs Hand Off to Judgment (any step)**: judge each blocked path, take the least destructive clearing action (delete or locally exclude non-project files, revert or finish this session's own abandoned write, commit the user's unrelated work on its own with the hash reported), re-run the exact `verification_argv`, and stop only for shared state whose owner the orchestrator cannot decide, naming the verb that resolves it.
```

Equivalently: `cp skills/do-work/CHANGELOG.md CHANGELOG.md` after the merge, which is what the house rules prescribe anyway.

**Until this seam lands, `_dev/tests/shipped-package-reference-contract.sh` fails with `changelog mirror differs` and the repository gate fails with it.** Zero citation findings — the mirror is the only complaint. This is the one thing on this REQ that needs the orchestrator's hand.

## Exploration

The sweep, which is this REQ's delegated exploration. Numbers are from the four shipped source roots in `suite/modules.tsv`, measured on the branch head.

**Census.** 130 arrow occurrences have a backticked `.md` path on the left: 102 in `skills/do-work`, 14 in `do-work-toolbox`, 10 in `do-work-knowledge`, 4 in `do-work-board`. By what follows the arrow:

| Form | Count | Checked now |
|---|---|---|
| `**Bold Section**` | 76 | 75 reach the citation scanner, 74 are resolved (the 76th is DT-2's `.claude/` path, which the scanner never sees as a candidate) |
| plain undelimited name | 45 | no — DT-1 |
| another `` `path` `` (a rename or mapping) | 7 | no, and correctly so: an arrow to a file names no section inside one |
| a quoted phrase | 2 | no |

A further ~150 arrows have a non-path left side entirely (alias tables, status transitions, command pipelines) — DT-4.

**What the 76 bold citations actually name**, resolved against their targets:

- 60 name an ATX heading exactly.
- 5 name a bold label that is not a heading: `**Frontmatter Quoting**`, `**Stakeholder REQ terminal semantics**`, `**Populating `write_set`**`, `**Populating `depends_on`**`, `**Worktree isolation is a separate axis**`. Every one of these is a *declared* target — three of the five sit on lines that say "Named contract —", "Named entry point —", or are cited by name from another file. This is D-01's evidence.
- 5 name a heading minus its parenthetical qualifier: `**Worktree Dispatch Mode**` (four sites, for `## Worktree Dispatch Mode (Step 1)`) and `**Hook Install Internals**` (for `## Hook Install Internals (used by actions/setup-memory.md → memory-module)`). This is D-02's evidence.
- 3 wrap the name across a line break — `crew-members/background-agents.md:187` and its mirrored copies in `do-work-knowledge` and `do-work-toolbox`. Invisible to a newline-bounded parser. This is D-06's evidence.
- 1 is DT-2's installed-topology path.
- 2 are genuinely dangling: the two the REQ named. **The honest first pass surfaced no third dangling citation.**

That last line is the answer to the REQ's expectation that more danglers would appear. They did appear on a strict reading — 11 more citations fail if "section" means "ATX heading, exact match". I checked each of the 11 by hand and none is a broken pointer: 5 name declared bold labels, 5 name a heading minus its qualifier, 1 is the `.claude/` path. The extra findings of this sweep are therefore not dangling citations but two blind spots in my own first implementation (the wrapped name, D-06; the bold phrase closing before the arrow, D-07) and the two report-only gaps DT-1 and DT-2.

**Where a renamed heading would still escape today:** the 45 undelimited citations (DT-1), the one `.claude/`-rooted citation (DT-2), and any citation whose target is not a `.md` file. Everything else in the arrow form is now covered.
