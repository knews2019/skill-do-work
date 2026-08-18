# REQ-243 Hand-Back — Check That Shipped Markdown Pointers Actually Resolve

## Branch

`worktree-agent-REQ-243-check-that-shipped-markdown-pointers-resolve`

Implementation commit: `5c395ab`
Merge base: `2432f450913d01e56f3685a6ca6a53c7424bd25d`

**Headline: half of this REQ was already built.** The relative-path half is covered today by
`_dev/tests/shipped-package-reference-contract.sh`. Only the `#anchor` half was a real gap, and
that is what shipped. Details in **## Pushback**.

---

## P-A-U

- [x] **[PLAN]**

Read `_dev/primes/prime-shell-commands.md`, CLAUDE.md, and the crew members
(`general.md`, `coding-guardrails.md`, `communication-style.md`, `maintenance.md`).

`maintenance: true` plus CLAUDE.md's **Delete before you add** made the first move a sweep, not a
build. The brief named the sweep explicitly: does `shipped-package-reference-contract.sh` already
walk these links? I read it in full (1009 lines) before writing anything.

It does. Lines 957–991 of the pre-change file walk every markdown target in every module source
tree from `suite/modules.tsv`, resolve each relative target from the citing file's own directory
(`markdown_path.parent / relative_target`), *and* resolve it a second time in installed topology
(`.claude/skills/...`), reporting `{markdown_path}:{line_number}: relative target is missing in
{source and/or installed} topology: {target}`. That is requirement 1 and requirement 3 of this REQ,
already shipped, with a ~500-line CommonMark-grade parser behind it that masks fenced code,
indented code, inline code, HTML comments, and escaped link syntax so that links inside examples
are not walked.

What it does not do is read the fragment. `urllib.parse.urlsplit(target).path` drops `#anchor`, so
`prescribed-shell-primitives.md#does-not-exist` resolves to a real file and passes.

So the plan became: **do not add a walker; add the fragment to the existing one.**

- Generate heading anchors the way the anchors were generated from the headings — GitHub's rule:
  rendered heading text, lowercased, every non-word/non-hyphen/non-space character dropped, spaces
  to hyphens, repeated slugs numbered `-1`, `-2` in document order.
- Take heading text from the **raw** line but gate the line on the **code-masked** line. Masking
  blanks inline code, which would turn ``## What `run` does`` into `## What does` and produce the
  wrong slug; but masking is exactly what makes a `#` inside a fence not a heading. Using masked to
  decide and raw to read gets both. `strip_markdown_code` preserves byte and line offsets — there
  are existing fixtures asserting precisely that — so the two line arrays index in lock step.
- Key the skip rule on a **condition**, per Closed Enumerations Go Stale, and write it into the file.
- Check the anchor in each *distinct* resolved target across the two topologies.

Planned write: `_dev/tests/shipped-package-reference-contract.sh` (not the declared write-set file).

- [x] **[APPLY]**

Written exactly as planned, one file, `+153 / -0`.

- `inline_link_label_pattern`, `atx_heading_pattern`, `heading_slug_cache` — module-level.
- `heading_rendered_text()` — collapses `[label](dest)` / `![alt](dest)` / `[label][ref]` to the
  label, to fixpoint, so a link's destination never pollutes the slug.
- `heading_anchor_slug()` — the GitHub rule above.
- `heading_anchor_slugs_from_text()` — the masked-gate/raw-read heading walk with duplicate numbering.
- `heading_anchor_slugs()` — path-keyed cache over it, `None` on unreadable.
- `run_anchor_slug_fixtures()` — six fixtures; registered next to `run_parser_fixtures()`.
- The check itself, inside the existing relative-target loop, after a new `continue` on the
  missing-target branch (so a link at a missing file is not also reported as a missing anchor).

No other file touched. No debug prints left; the instrumentation used in UNIFY was a copy under
`/tmp` and a temporary `_dev/tests/zz-req243-probe.sh` that was deleted in the same command that ran it.

- [x] **[UNIFY]**

```
 _dev/tests/shipped-package-reference-contract.sh | 153 +++++++++++++++++++++++
 1 file changed, 153 insertions(+)
```

One file changed. What I verified on it:

- `_dev/tests/shipped-package-reference-contract.sh` — read the complete diff. Additions only; no
  existing line altered except the `if missing_locations:` branch gaining a `continue`. Confirmed
  that `continue` cannot skip anything, because that branch was previously the last statement in
  the loop body. Confirmed the new code runs *after* every existing skip (scheme'd URL, dynamic
  target, root-absolute, empty path) so it can never widen what the contract walks.
- ShellCheck warning-level over 50 tracked shell files: clean (the Python lives in a quoted
  heredoc, so the added lines are data to ShellCheck; the surrounding shell is unchanged).
- Working tree clean apart from the committed file — `git status --porcelain --untracked-files=all`
  shows nothing after commit. No stray scratch files in the worktree or the repo root.
- Queue guard `git diff --name-only <pre>...HEAD -- do-work/` → empty.
- Serial-file guard over `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`,
  `CHANGELOG.md` → empty.
- Full `maintainer-verify.sh` → exit 0 (pasted below).

---

## Files Changed

```
 _dev/tests/shipped-package-reference-contract.sh | 153 +++++++++++++++++++++++
 1 file changed, 153 insertions(+)
```

- **`_dev/tests/shipped-package-reference-contract.sh`** — added heading-anchor resolution to the
  markdown-target walk that already lived here, plus the anchor-slug generator it needs and six
  fixtures pinning that generator. Chosen over the declared write-set file because this contract
  already walks exactly the set in question, already resolves each target from the citing file's
  own directory in both topologies, and already owns the markdown parser required to tell a real
  link from one inside a code fence. A second walker in
  `prescribed-shell-canonicalization.sh` would have duplicated all of it.

**Declared write set was `_dev/tests/prescribed-shell-canonicalization.sh`; that file was not
touched.** See **## Pushback**.

---

## Red-Green Evidence

### Baseline — both suites green before any change

```
=== BASELINE reference contract ===
shipped package reference contract: PASS
exit=0
=== BASELINE canonicalization ===
Prescribed shell primitive canonicalization checks passed.
exit=0
```

### Pre-change probe — proving the brief's "Why RED now" is wrong

Before writing a line of code, both mutations were run against the **unmodified** tree.

Wrong relative depth (`../../do-work/...` → `../do-work/...`), pre-change code:

```
=== WRONG-DEPTH: reference contract ===
FAIL: skills/do-work-toolbox/actions/present-work.md:130: relative target is missing in source and installed topology: ../do-work/docs/prescribed-shell-primitives.md#portfolio-summary-publication
FAIL: skills/do-work-toolbox/actions/present-work.md:136: relative target is missing in source and installed topology: ../do-work/docs/prescribed-shell-primitives.md#portfolio-summary-publication
shipped package reference contract: FAIL (2 broken reference(s), changelog mirror matches)
exit=1
```

Already caught. Wrong anchor (`#portfolio-summary-publication` → `#portfolio-summary-publications`),
pre-change code:

```
=== WRONG-ANCHOR: reference contract ===
shipped package reference contract: PASS
exit=0
=== WRONG-ANCHOR: canonicalization ===
Prescribed shell primitive canonicalization checks passed.
exit=0
```

Not caught anywhere. **The anchor is the entire real gap.**

### Non-vacuity — the new check does work on the live corpus

A green run proves nothing if every link is skipped. Instrumented copy, run from the repo root:

```
CHECKED skills/do-work-board/actions/board.md:85 -> skills/do-work/docs/prescribed-shell-primitives.md#local-git-ignore
CHECKED skills/do-work-knowledge/actions/memory-reference.md:67 -> skills/do-work/docs/prescribed-shell-primitives.md#raw-text-before-shell-quoting
CHECKED skills/do-work-toolbox/actions/present-work.md:130 -> skills/do-work/docs/prescribed-shell-primitives.md#portfolio-summary-publication
CHECKED skills/do-work/actions/work-reference.md:349 -> skills/do-work/docs/prescribed-shell-primitives.md#state-across-command-blocks
... (27 total)
--- total anchors checked ---
      27
```

27 anchors across 4 packages actually resolve through the new code path, all green.

> Note on method: an earlier run of the same probe from `/tmp` reported 0 checked and looked like a
> vacuous pass. It was not — the script derives `repo_root` from `BASH_SOURCE`, so from `/tmp` it
> failed at `cannot read /suite/modules.tsv` before reaching any link. Re-run from inside the repo.
> Worth recording: a "0 checks ran" result from a relocated script is a harness artifact, not a finding.

### RED 1 — wrong relative depth, post-change

```
FAIL: skills/do-work-toolbox/actions/present-work.md:130: relative target is missing in source and installed topology: ../do-work/docs/prescribed-shell-primitives.md#portfolio-summary-publication
FAIL: skills/do-work-toolbox/actions/present-work.md:136: relative target is missing in source and installed topology: ../do-work/docs/prescribed-shell-primitives.md#portfolio-summary-publication
shipped package reference contract: FAIL (2 broken reference(s), changelog mirror matches)
exit=1
```

Reverted → GREEN:

```
--- GREEN after reverting RED 1 ---
shipped package reference contract: PASS
exit=0
```

### RED 2 — wrong anchor, post-change (this is the new coverage)

```
FAIL: skills/do-work-toolbox/actions/present-work.md:130: anchor #portfolio-summary-publications is not a heading in source topology target skills/do-work/docs/prescribed-shell-primitives.md: ../../do-work/docs/prescribed-shell-primitives.md#portfolio-summary-publications
FAIL: skills/do-work-toolbox/actions/present-work.md:136: anchor #portfolio-summary-publications is not a heading in source topology target skills/do-work/docs/prescribed-shell-primitives.md: ../../do-work/docs/prescribed-shell-primitives.md#portfolio-summary-publications
shipped package reference contract: FAIL (2 broken reference(s), changelog mirror matches)
exit=1
```

The message carries the citing file, the line, the anchor, the resolved target, and the raw link.
Reverted → GREEN (see the final verification run).

### RED 3 — mutating the slug rule itself, to prove the fixtures are not decorative

Removed `.lower()` from `heading_anchor_slug`:

```
FAIL: anchor fixture 'case folding and punctuation removal': expected ['portfolio-summary-publication', 'raw-text-before-shell-quoting'], got ['Portfolio-summary-publication', 'Raw-text-before-shell-quoting']
FAIL: anchor fixture 'inline code contributes its content, not its backticks': expected ['what-run-does-and-does-not-do'], got ['What-run-does-and-does-not-do']
FAIL: anchor fixture 'link labels replace their destinations': expected ['see-the-guide-first'], got ['See-the-guide-first']
FAIL: anchor fixture 'code blocks and comments hold no headings': expected ['real-heading'], got ['Real-Heading']
FAIL: anchor fixture 'repeated headings take numbered suffixes': expected ['same', 'same-1', 'same-2'], got ['Same', 'Same-1', 'Same-2']
FAIL: anchor fixture 'closing hash sequences are not part of the anchor': expected ['trailing-hashes'], got ['Trailing-Hashes']
exit=1
```

Restored → `shipped package reference contract: PASS`, `exit=0`.

---

## Verification

Run from the worktree root, unpiped:

```
maintainer-verify: checking Go go1.26.1
maintainer-verify: checking ShellCheck 0.11.0
maintainer-verify: ShellCheck warning-level lint (50 tracked files)
maintainer-verify: aggregate contract suite
Maintainer verification self-test passed.
Suite manifest contract probes passed.
shipped package reference contract: PASS
Shell-block lint self-test passed.
Shell-block lint passed: 74 fenced blocks and 31 shipped shell files; ShellCheck enabled.
SessionStart hook behavior probes passed.
Prescribed shell primitive canonicalization checks passed.
Defensive-surface exact deletion regressions passed.
record-commit-hash and blanked-req-scan guard probes passed.
update-script behavior probes passed.
Prescribed shell script behavior probes passed (42 named script cases).
staged skills contract: PASS
suite installer behavior probes passed.
p50 estimator suite: all probes passed.
Contract regression checks passed.
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	20.153s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (5.36s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	5.602s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.619s
Maintainer verification passed.
0
```

`echo $?` → `0`.

---

## Integration Seams

**None.** One self-contained file, no shared registry, no cross-REQ text, no new suite to register
(`contract-regressions.sh` already invokes `shipped-package-reference-contract.sh` at line 3719).

One thing for the integrator to be aware of rather than apply: this check now runs against every
shipped markdown file on every verify. If a sibling REQ in this wave lands a markdown pointer with a
bad anchor, `maintainer-verify` will fail at the merge, naming the file and line. That is the check
working, not a conflict with this REQ.

---

## Decisions

**D-01 — The check lives in `shipped-package-reference-contract.sh`, not in the declared write-set
file.** That contract already walks the same corpus, already resolves relative targets from the
citing file's directory in both source and installed topology, and already owns the code-masking
markdown parser needed to distinguish a live link from one inside a fence. Adding an anchor walker
to `prescribed-shell-canonicalization.sh` would have meant a second parser over the same files.
Reach: any future "check something about shipped markdown links" REQ should extend this same walk.

**D-02 — Anchor slugs follow GitHub's generation rule, because that is the rule the anchors in this
repo were generated under.** The REQ says "compared the way the anchor was generated from it".
Concretely: rendered text, lowercased, characters outside `[\w\- ]` dropped, spaces to hyphens,
duplicates numbered `-1`/`-2` in document order. Reach: if the repo's docs are ever published by a
renderer with a different slug rule, this function is the single place to change.

**D-03 — Headings are read from the raw line and gated on the code-masked line.** Reading masked
text alone would corrupt any heading containing inline code (``## What `run` does`` → `what-does`);
reading raw text alone would treat `#` comments inside fenced shell blocks as headings, and this
repo's markdown is full of those. The gate works only because `strip_markdown_code` preserves line
offsets — an existing invariant with its own fixtures. Reach: anything else that needs "the real
lines of this markdown file" should use the same two-array technique rather than re-deriving fences.

**D-04 — Anchors resolve against ATX headings only, and that limitation is stated in the file.**
The repo has no setext headings. Supporting them is speculative machinery per YAGNI, and the failure
mode if someone writes one and links to it is a loud FAIL naming the file and anchor — not a silent
miss. Recorded as a decision because the honest framing matters: this is a stated limitation with a
visible failure, not an undetected hole.

**D-05 — The skip rule is stated as a condition in the file, per Closed Enumerations Go Stale.**
A link is anchor-checked exactly when it survives the pre-existing skips and carries both a fragment
and a Markdown target. Everything skipped is skipped for a reason the comment gives — a scheme'd URL
is not fetched, a bare fragment names no file, a root-absolute or templated target has no single
on-disk meaning, and a heading only exists in Markdown. No list of files or link forms to go stale.

---

## Lessons Learned

- **The gap a REQ names and the gap a REQ has are not always the same gap.** REQ-243's Red-Green
  Proof asserted "no check reads markdown links at all, so the wrong depth passes the entire suite
  today." Running the stated mutation against the unmodified tree took one command and disproved it:
  `shipped-package-reference-contract.sh` has resolved relative markdown targets from the citing
  file's own directory, in two topologies, for some time. Run the RED against pre-change code
  *before* building — not to satisfy a ritual, but because a "RED" that was already red is the
  cheapest possible signal that the work is half-done already.
- **Two hand-fixes of the same class is good evidence a checker is needed; it is not evidence that
  no checker exists.** REQ-230 and REQ-238 both verified the path by hand and the anchor by hand.
  The path half was machine-checked the whole time and nobody knew. Before writing the third
  instance of a manual check, grep for the check as well as for the defect.
- **"0 checks ran" from a relocated script is usually the harness.** The instrumented probe reported
  zero anchors checked when copied to `/tmp`, which reads exactly like a vacuous pass. The cause was
  `repo_root` derived from `BASH_SOURCE`. Any suite that locates itself relative to its own file
  cannot be probed from outside the tree — copy it *into* the tree, run, delete.
- **Where masking helps and where it hurts, in one file.** `strip_markdown_code` is what makes link
  extraction trustworthy and what makes heading extraction wrong, for the same reason. The
  offset-preserving property — already fixture-locked — is what lets both be right at once.

---

## Pushback

**1. The REQ's core premise is half wrong, and I proceeded on the corrected version.**

> **Why RED now:** no check reads markdown links at all, so the wrong depth passes the entire suite
> today.

This is false. `_dev/tests/shipped-package-reference-contract.sh` reads markdown links, resolves
them from the citing file's own directory, and additionally resolves them in installed topology. The
brief's own suggested mutation fails on the unmodified tree — evidence quoted above under
*Pre-change probe*. Requirements 1 and 3 of this REQ ("every relative `.md` link resolves from the
citing file's own directory", "a broken link names the citing file, the line, the raw link, and what
it resolved to") were satisfied before I started.

The genuinely uncovered requirement was requirement 2, the `#anchor`. I built that and nothing else.
This is the **Delete before you add** sweep returning a real deletion: roughly half the specified
work was already present, and the right output was ~150 lines inside an existing walk rather than a
new walker.

**2. The declared write set was wrong, and I deviated as the brief permitted.**

Declared: `_dev/tests/prescribed-shell-canonicalization.sh`. Touched:
`_dev/tests/shipped-package-reference-contract.sh`. The REQ's own Constraints anticipated this
("check whether the shipped-package reference contract already walks these links for a different
reason and could carry this one instead of a second walker") — it does, and it now does. The
canonicalization suite proves phrases are *absent*; it has no parser and no notion of a link. Putting
a link checker there would have meant a second CommonMark-grade parser over the same file set.

No collision risk from the deviation: no sibling REQ in this wave declares
`shipped-package-reference-contract.sh`, and the file is maintainer-side (`_dev/` is export-ignored).

**3. One scope question I did *not* decide unilaterally.**

The REQ lists "bare anchors into the same file" among the things to skip, and I skipped them as
specified. Worth flagging for a future pass: now that heading slugs are computed and cached,
validating a same-file `#anchor` is about three lines and closes a hole of exactly the same class —
a `[see below](#some-section)` pointing at a heading that was renamed is silently broken today.
I left it out because the REQ named it a skip deliberately rather than by omission, and narrowing or
widening a scope decision the author made on purpose is the author's call, not mine. If you want it,
it is a very small follow-up.

---

## Addendum — Review Remediation

Second commit on the same branch, addressing findings 1 and 2 from the adversarial review.

**Commit:** `7b6d8fc` — *split heading lines on newlines only, and cover the installed topology*

```
 _dev/tests/shipped-package-reference-contract.sh | 154 ++++++++++++++++++-----
 1 file changed, 126 insertions(+), 28 deletions(-)
```

Still one file. Queue guard and serial-file guard over the full two-commit range: both empty.

### Finding 1 — `splitlines()` desync

Confirmed on the real functions before changing anything. Both of the reviewer's documents
reproduce exactly:

```
=== A: formfeed inside a fence ===
raw_lines = 8  masked_lines = 7
slugs = []

=== B: formfeed in inline code ===
raw_lines = 5  masked_lines = 4
slugs = ['fake']

=== control: split('\n') counts agree in both cases ===
A raw: 8 masked: 8
B raw: 5 masked: 5
```

Fix is `.split("\n")` in the two places that built the arrays, so the code depends on exactly the
property the existing parser fixtures lock — length and `\n` positions — and nothing stronger. The
comment that claimed "preserves byte and line offsets" is replaced with one that names the nine code
points, says masking turns each into a space, and states the consequence. The overclaim was half the
defect and it is gone.

**One correction to the finding's characterization, which changed what I asserted.** Document B is
described as a false negative *plus* a false positive — "a link to `target.md#fake` PASSES though no
such anchor exists". `#fake` does exist there. The backtick opening line two never closes, so it is
literal text and `# Fake` on line three is a genuine ATX heading. Verified against the fixed code:

```
masked repr: '     \n`open\n# Fake\n# Real Heading\n'
slugs B (fixed): ['fake', 'real-heading']
```

So the defect in both directions is the same one: the desynced gate **drops real headings**, and a
valid anchor gets reported broken. I also tried to construct a genuine false accept — a fenced
`# heading` promoted by the shift — and could not:

```
C fenced heading promoted?
   buggy=[]
   fixed=['live']
```

I mention it because I nearly wrote a fixture asserting `{"real-heading"}` for document B on the
strength of the finding's wording. That fixture would have been green on broken code and red on
correct code. The fixture asserts the verified truth instead.

New fixtures in `run_anchor_slug_fixtures` (now 8 cases):

```python
(
    "only newlines divide lines, so a masked form feed cannot shift the gate",
    "```sh\nprintf 'a\x0cb'\n```\n\n# Real One\n\n# Real Two\n",
    {"real-one", "real-two"},
),
(
    # The unclosed backtick on line two is literal text, so both headings below
    # it are genuine; a desynced gate finds only the first.
    "a masked form feed leaves the headings after it aligned",
    "`a\x0cb`\n`open\n# One\n# Two\n",
    {"one", "two"},
),
```

**RED A** — `.splitlines()` put back in both places:

```
FAIL: anchor fixture 'only newlines divide lines, so a masked form feed cannot shift the gate': expected ['real-one', 'real-two'], got []
FAIL: anchor fixture 'a masked form feed leaves the headings after it aligned': expected ['one', 'two'], got ['one']
exit=1
```

Both fail for the reason they exist — headings lost to the shift, not a reference error. Restored →
`shipped package reference contract: PASS`, `exit=0`.

### Finding 2 — installed topology never executed

Confirmed: the dedup guard suppressed the second pass on all 27 live anchors. Extracted two units so
the branch is reachable without a synthetic dual-topology tree:

- `link_topology_targets(source_target, installed_source_target)` — the pair a pointer must satisfy.
- `anchor_failure_messages(anchor, topology_targets, read_anchor_slugs)` — the per-topology walk,
  returning messages instead of calling `fail` directly.

Four fixtures in `run_anchor_topology_fixtures`: both topologies resolving to one target read once
(the live case), a distinct installed target read on its own, an unreadable installed target
reported rather than skipped, and an anchor missing from both reported twice.

**RED B** — first attempt did not go red, and that is worth recording. Dropping the installed
topology from the *call site* left the suite green: the fixtures called
`anchor_failure_messages` with hand-built tuples, so they covered the function but not the wiring,
and the live corpus cannot tell the two topologies apart. Rather than report that as a residual gap
I moved the pairing into `link_topology_targets` and had every fixture build its pair through it.

**RED B1** — installed dropped from the pairing:

```
FAIL: anchor topology fixture 'a distinct installed target is read on its own': expected ['anchor #present is not a heading in installed topology target installed.md'], got []
FAIL: anchor topology fixture 'an unreadable installed target is reported, not skipped': expected ['cannot read installed topology target unreadable.md'], got []
FAIL: anchor topology fixture 'an anchor missing from both topologies is reported twice': expected [... 'installed topology target installed.md'], got ['anchor #absent is not a heading in source topology target source.md']
exit=1
```

Restored → PASS, `exit=0`.

### Verification

```
maintainer-verify: ShellCheck warning-level lint (50 tracked files)
maintainer-verify: aggregate contract suite
Maintainer verification self-test passed.
Suite manifest contract probes passed.
shipped package reference contract: PASS
...
Contract regression checks passed.
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	16.054s
maintainer-verify: queue-kanban strict JavaScript behavior lane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (4.87s)
ok  	github.com/knews2019/skill-do-work/queue-kanban	5.085s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.523s
Maintainer verification passed.
0
```

`echo $?` → `0`.

Working tree clean; all Python extractions used for probing deleted from `/tmp` and none written
into the repo.

### New Decisions

**D-06 — Depend on the invariant that is locked, not the one that is true today.** The masked and
raw arrays are split on `"\n"` alone. `strip_markdown_code` happens to preserve more than that, but
only length and `\n` positions are fixture-locked, and `splitlines()` silently needed more. Where a
consumer can be written against the weaker guarantee, write it against the weaker guarantee — then
the fixture that exists is the fixture that protects it. Reach: anything else pairing masked and raw
markdown must split on `"\n"`, and the comment there says so.

**D-07 — A branch the corpus cannot reach becomes a named unit.** `link_topology_targets` and
`anchor_failure_messages` exist because the source and installed topologies resolve to one file for
every link in this repository, so the installed branch had no execution path. Extracting the pairing
as well as the walk is what makes the mutation "drop a topology" fail; extracting only the walk left
the wiring untested and a mutation survived. Reach: when coverage for a branch is impossible from
real data, the fix is a seam at the point the *data* is constructed, not only at the point it is
consumed.

### Note on the follow-ups

Not touched, as instructed: HTML tags/entities in heading text, blockquoted ATX headings, same-file
bare `#anchor` validation, and the `os.path.normpath` path-escape. Agreed on the last one being the
only real hole — and worth flagging that my change does escalate it, since an escaping `..` target
now reaches `read_text()` rather than only `stat()`. It is bounded by the missing-target `continue`
(the path must exist in both topologies to be read at all), but that is a weaker bound than clamping.
