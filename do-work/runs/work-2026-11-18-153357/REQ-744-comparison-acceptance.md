# REQ-744 request-outcomes acceptance

**Passed on 2026-11-18.** For REQ-001 (request outcomes across two test windows), successful requests rose by 70, from 90 to 160. The success rate fell by 10 percentage points, from 90% to 80%, as completed requests doubled from 100 to 200. The report uses synthetic fixture data and the installed action’s non-visual evidence mode. [Report](/private/tmp/req-744-acceptance-8u9ze8c1/request-outcomes/ai-reports/2026-11-18_1834_REQ-001-request-outcomes/index.html)

## Install and evidence receipts

- Archived candidate `5bdc38b2dd2c86e92658f73ffc73d86fb97881f8` from the supplied source checkout. The real `install-do-work-suite.sh --project-root <request-outcomes> --archive <candidate archive>` returned `install-suite: success`, four verified modules at v0.325.0, and `rollback: not_needed`. Byte comparisons passed for the installed action, reference, kit, gallery, presets and CLI launcher. The consumer was clean after install. [Installer log](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/request-outcomes-install.log)
- Read the installed toolbox route and action, then applied the completed-work reference before reading the archive. REQ-001 is `completed`. The archived review records the ID, window and arithmetic checks. No issues are recorded. Archive `95affa8d7aeffacc9e4a441a78a72450c28fc587` has reachable implementation `31f4eccee016c46cb391cc2404b282171bef91c9` as its parent. The installed merge-aware helper returned `show-commit-diff: success`. The current CSVs and check script matched the implementation commit. [Archive](/private/tmp/req-744-acceptance-8u9ze8c1/request-outcomes/do-work/archive/UR-001/REQ-001-record-request-outcomes.md) · [Diff receipt](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/request-outcomes-implementation.diff)
- Replayed `python3 scripts/check-request-outcomes.py` from the consumer root. Exit 0. The stdout below matched the archived testing receipt byte for byte. A separate recount from both CSVs also matched both chart endpoint pairs. [Replay](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/request-outcomes-replay.txt)

```text
TEST FIXTURE only (synthetic records)
| Window | UTC start | UTC end (exclusive) | Completed requests | Successful requests | Success rate |
|---|---|---|---:|---:|---:|
| Window A | 2026-11-01T00:00:00Z | 2026-11-02T00:00:00Z | 100 | 90 | 90% |
| Window B | 2026-11-02T00:00:00Z | 2026-11-03T00:00:00Z | 200 | 160 | 80% |
```

## Report and render receipts

After collecting that evidence and before creating the bundle, ran the installed `ai-report-gallery --list`. It returned `parts-over-time`, `rhythm-over-cycle`, `two-measures-side-by-side`, and `verdict-with-checks` in that order. Opened only the winning and runner-up sample records. The page head has exactly one `ai-report-gallery-sample` meta with content `two-measures-side-by-side`. This sample fits because the count and rate move in opposite directions across the same windows. Both charts show the same UTC windows in A/B order, put the unit in the heading, label each endpoint, and cite the archived output. One window selector highlights the same window in both charts and displays its calculation. [Actual gallery list](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/request-outcomes-gallery-list.txt)

The installed render judge ran twice because the report contains a text-bearing SVG. Both runs exited 0 with `verdict: pass` and `findings: []`. Visually inspected both sets of Chrome captures. The charts sit side by side at 1440 px and stack in order at 390 px. The dates, values and both closing choices are visible in both themes. Judge output and captures were cleaned up after inspection. [Pass 1 receipt](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/judge-captures-1.json) · [Pass 2 receipt](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/judge-captures-2.json)

| Capture | Viewport | Final scroll width / inner width |
|---|---|---|
| `wide-light.png` | 1440 × 1000 | 1440 / 1440 |
| `wide-dark.png` | 1440 × 1000 | 1440 / 1440 |
| `phone-light.png` | 390 × 844 | 390 / 390 |
| `phone-dark.png` | 390 × 844 | 390 / 390 |

## Displayed commands and prepare receipt

**Verdict with checks** is the one fitting runner-up. It puts the result first and makes the dates, denominator, units and arithmetic easy to audit. The installed `80-20` preset keeps those checks up front and moves supporting detail behind a disclosure.

```sh
ai-report second-opinion --preset 80-20 --editor-direction 'Lead with the 10 percentage point drop, show checks for the UTC windows, denominator, units and arithmetic, and fold the fixture rows and calculation into one disclosure.' ai-reports/2026-11-18_1834_REQ-001-request-outcomes
```

**Feynman technique** helps the reader explain the result plainly and check that the explanation includes both denominators.

```sh
ai-report second-opinion --preset feynman ai-reports/2026-11-18_1834_REQ-001-request-outcomes
```

Decoded both displayed `<pre><code>` lines, matched them to the handback, and split them using Python `shlex.split`. They produced seven and five tokens, with one exact current-bundle target each. The runner-up direction was one nonempty argument with no internal single quote. Both presets resolved to installed Markdown files after the bundle lookup. Forwarded the runner-up’s option tokens and bundle to the actual launcher with `ai-report-second-opinion --prepare --backend codex --model gpt-6.2-sol`. [Token receipt](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/closing-command-tokens.json)

| Prepare field | Actual value |
|---|---|
| Exit / outcome | `0` / `success` |
| State / judge verdict | `prepared` / `not-run` |
| Model evidence | `session-declared` |
| Preset file | `.claude/skills/do-work-toolbox/docs/second-opinion-presets/80-20.md` |
| Base page | `ai-reports/2026-11-18_1834_REQ-001-request-outcomes/index.html` |
| Staging page | `.codex-gpt-6.2-sol-staging-1208801270.html`, 0 bytes |
| Staging reply | `.codex-gpt-6.2-sol-staging-1208801270-reply.md`, absent |
| Outside-bundle paths / published page and reply | `[]` / empty |

The prepared prompt contained both the full installed preset and the exact direction. Its bundle-relative `index.html` resolved to the same file as the typed base-page field. All three original files matched their pre-prepare SHA-256 digests, including `index.html` at `81669473d7596aa4d72d35e1a60a6b457f6bd1bb4244aa09af74d57811c55017`. Stopped at prepare and removed the empty staging reservation and its record after preserving the receipt. [Raw CLI receipt](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/prepare-run.json) · [Verified prepare receipt](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/prepare-verified.json)

## Problems and final state

- The first install check tried to compare `verdict-with-checks.md` as a preset. `verdict-with-checks` is a gallery sample, and the corrected directory comparisons passed. The installer had already succeeded.
- Visual review found tight spacing above the rate chart labels. Added 12 px of plot headroom and inspected the second four-capture pass.
- The first verification assertion expected an absolute page path in the editor prompt. The prompt uses bundle-relative `index.html`. Resolving it against the explicit bundle path verified the original prepare result.
- The consumer’s HEAD remains the archive commit. The only reported Git changes are the three untracked bundle files. The bundle has the original kit CSS/JS and no `FILL:` markers. The organize check was silent. The catalog check printed that the new bundle is absent from `ai-reports/catalog.json` and named `do-work-toolbox ai-report index`. [Final Git status](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/status-after-prepare.txt) · [Catalog check](/private/tmp/req-744-acceptance-8u9ze8c1/candidate.YxJcU1mG/catalog-check.txt)

The consumer’s instructions leave the bundle uncommitted. The installed action’s exact-path commit line, from the consumer root, is:

```sh
git add -- ai-reports/2026-11-18_1834_REQ-001-request-outcomes && git commit -m 'Report REQ-001 request outcomes' -- ai-reports/2026-11-18_1834_REQ-001-request-outcomes
```

ai-reports/2026-11-18_1834_REQ-001-request-outcomes/index.html · [Open local report](file:///private/tmp/req-744-acceptance-8u9ze8c1/request-outcomes/ai-reports/2026-11-18_1834_REQ-001-request-outcomes/index.html)
