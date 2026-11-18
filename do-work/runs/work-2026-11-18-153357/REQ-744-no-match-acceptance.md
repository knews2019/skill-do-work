# REQ-744 no-match acceptance

**Pass for the ordinary text-only report.** Candidate `5bdc38b2dd2c86e92658f73ffc73d86fb97881f8` produced an installed-action report with the corrected label, a single gallery marker of `none`, and both closing suggestions. Chrome judge exited 0 with four captures and zero findings. The revision also met the action's acceptance rule with one expected prior-bundle link finding.

## Install and provenance

- Consumer was only `/private/tmp/req-744-acceptance-8u9ze8c1/label-correction`. Archived the exact candidate from the source worktree named in the fixture handback and ran its real `install-do-work-suite.sh --project-root <consumer> --archive <candidate archive>` with the handback's `y` input. The installed action, reference, and skeleton each passed `cmp` against the candidate archive. The consumer's Git status was clean after installation.
- Routed `do-work-toolbox ai-report REQ-002` through installed `.claude/skills/do-work-toolbox/SKILL.md` and `actions/ai-report.md` in this session. It is an agent action. The installed launcher provided the gallery, judge, and catalog operations.
- Read the parent UR input and completed REQ. Archive commit `4376515ee7d0ef33dbec8cc0674de28540304e20` directly follows implementation `da70efd554643081d3c2b5ef97c19cfa1e961df4`. The installed merge-aware helper showed one edit in `locales/en.json:1`, from `"account_settings": "Acount settings"` to `"account_settings": "Account settings"`. Current code agrees. The archived test replay exited 0 and printed exactly `Account settings`. No completion issues were recorded.

## Ordinary report and writer summary

**REQ-002, English Account settings label correction.** The label now reads “Account settings.” Non-visual evidence came from the archived request and review, the one-line implementation diff, the current locale file, and the passing exact-text check. The verdict and proof show the corrected text, and the report states “UI captures were not expected for this work.” No diagram was needed to explain a single string. The bundle contains its HTML and the two unmodified local kit assets.

**Gallery summary.** No sample fit. The evidence is a spelling correction with an exact-value check. It has no changing total, repeating cycle, paired measures, or judgment to audit. The head contains exactly one `<meta name="ai-report-gallery-sample" content="none">`. The installed CLI was called with `ai-report-gallery --list` after collecting and replaying evidence and before building the report. All four rows were available. There was no Borrow record to open.

**Spelling table.** Put prior text first, corrected text second, and the check immediately below. That keeps both spellings and their proof together in reading order.

**Predict and check.** Ask the reader to identify the changed word and the stable key, then reveal the actual value so the reader can explain the scope of the change.

```text
ai-report second-opinion --editor-direction 'Lead with a two-row spelling table, prior text then corrected text, and place the exact locale check directly below it' ai-reports/2026-11-18_1832_REQ-002-account-settings-label
ai-report second-opinion --editor-direction 'Add a predict-and-check exercise that asks the reader which word changes and which translation key stays the same, then reveals the exact file value' ai-reports/2026-11-18_1832_REQ-002-account-settings-label
```

Each line was extracted from `#another-layout` and checked using `shlex.split`. Both produced five tokens with exactly one nonempty direction argument and one positional bundle target resolving to the report. Both are intentionally preset-free. The judge captured 1440×1000 and 390×844 in light and dark, with scroll widths equal to viewport widths. Visual inspection confirmed wide alignment, phone stacking, and readable closing lines. The organize check printed nothing. The catalog check recommended `do-work-toolbox ai-report index`, as expected before the revision builds the catalog.

Consumer commit command, per the action's handback rule for this fixture.

```bash
git add -- ai-reports/2026-11-18_1832_REQ-002-account-settings-label && git commit -m "Report Account settings label correction" -- ai-reports/2026-11-18_1832_REQ-002-account-settings-label
```

Ordinary report: [index.html](/private/tmp/req-744-acceptance-8u9ze8c1/label-correction/ai-reports/2026-11-18_1832_REQ-002-account-settings-label/index.html)

## Installed route inspection

| Route | Result |
|---|---|
| `--kind proposal\|root-cause` | Static inspection of installed action Step 1 and reference § Proposal and Root-Cause Kinds. It takes the open-work entry after safety loading, keeps its kind-specific decision-first template and `ai-report-kind` meta, and inherits Steps 3–4. Step 3 runs gallery selection after evidence and before visuals. The shared design rules add the gallery marker and Step 8 records the decision. Other kind values or a missing target stop with usage. No kind brief was generated in this fixture. |
| `revise` | Executed below. The installed route catalogs and resolves the prior bundle, safely rereads it, rechecks claims, then joins Steps 2–8. Step 3 explicitly ranks afresh. |
| `second-opinion` | The installed dispatcher runs before ordinary Steps 1–8. It takes a current existing bundle or explicit base page and delegates prepare/finalize to `ai-report-second-opinion`. The launch path skips gallery reranking. A separate prepare-only probe on the ordinary `index.html` returned `prepared` / `not-run`, the exact explicit `base_page`, both staging paths, and the supplied direction in the prompt. The report hashes were unchanged. Neither displayed line was run and no editor model was launched. |

## Revision and writer summary

Ran `do-work-toolbox ai-report revise ai-reports/2026-11-18_1832_REQ-002-account-settings-label "Name the inserted letter and clarify which file the test reads"`. After resolving the prior bundle through the catalog, reread the prior report as source, checked its claims against the archive and current locale file, and replayed the diff and test. The fresh gallery listing again returned four candidates. **No sample fit the rechecked single spelling change**, so the revision independently set one head marker to `none`. It names the inserted second c and the file read by the test. The revision block follows the heading and links to the prior page. The head points `ai-report-supersedes` to the bare prior folder. The rebuilt catalog links both directions and all three prior file SHA-256 checks still pass.

The revision judge exited 1 with one `AI-REPORT-JUDGE-BROKEN-LINK` (404 for `../2026-11-18_1832_REQ-002-account-settings-label/index.html`). This is the installed action's expected finding when the judge serves only the revised bundle. The relative link resolves to the prior report on disk. All four captures rendered without overflow and were visually inspected. The spelling-table layout and predict-and-check method still fit the same evidence.

```text
ai-report second-opinion --editor-direction 'Lead with a two-row spelling table, prior text then corrected text, and place the exact locale check directly below it' ai-reports/2026-11-18_1833_REQ-002-account-settings-label-rev1
ai-report second-opinion --editor-direction 'Add a predict-and-check exercise that asks the reader which word changes and which translation key stays the same, then reveals the exact file value' ai-reports/2026-11-18_1833_REQ-002-account-settings-label-rev1
```

Both closing lines passed the same token checks and target the exact revised bundle. The organize check printed nothing.

```bash
git add -- ai-reports/2026-11-18_1833_REQ-002-account-settings-label-rev1 ai-reports/catalog.json ai-reports/index.html && git commit -m "Clarify Account settings label report" -- ai-reports/2026-11-18_1833_REQ-002-account-settings-label-rev1 ai-reports/catalog.json ai-reports/index.html
```

Revised report: [index.html](/private/tmp/req-744-acceptance-8u9ze8c1/label-correction/ai-reports/2026-11-18_1833_REQ-002-account-settings-label-rev1/index.html)

## Receipts and deviations

Scratch receipts are in `/private/tmp/req-744-acceptance-8u9ze8c1/candidate.UjZwsMLP`.

- Installation: `label-correction-install.log`, `label-correction-install-gallery.txt`.
- Ordinary provenance and gallery: `label-provenance.txt`, `label-implementation.diff`, `label-test.txt`, `label-ordinary-gallery.txt`.
- Judge and closing lines: `ordinary-judge.json`, `revise-judge.json`, `revise-judge-exit.txt`, `ordinary-closing-lines.txt`, `revise-closing-lines.txt`, `bundle-assertions.txt`. Render directories and captures were removed after inspection per Step 8. Judge JSON was retained.
- Revision and preparation: `label-revise-gallery.txt`, `revise-test.txt`, `revise-implementation.diff`, `revise-initial-index.txt`, `revise-final-index.txt`, `prior-byte-check.txt`, `explicit-page-prepare.json`. Prepare created two hidden staging files, which were moved to the receipt directory. No finalize ran.
- The first revision judge shell wrapper stopped before invoking the judge because zsh reserves `status`. The retry used `judge_exit` and ran the installed judge normally. The fixture handback's metrics consumer and metrics-only prepare probe were outside this task's scope. The consumer has no new commits. `final-status.txt` lists only the two report bundles and the catalog pair. No child agents, push, or PR.
