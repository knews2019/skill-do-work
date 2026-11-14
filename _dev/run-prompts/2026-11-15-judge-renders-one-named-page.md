# Capture prompt: ai-report-judge can render one named page of a bundle

## How to run this prompt

Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.

## The request

Two changes to `ai-report-judge`, both found while building report pages on 2026-11-14.

1. **Render one named page.** Today the judge renders only `index.html` (`skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go`, lines 25, 111 and 168). To check or capture a second-opinion page such as `index-v1c01-<name>.html`, an agent must copy the whole bundle and write that page over the copy's `index.html`. Add an option, for example `--page <file.html>`, that renders that page of the bundle in place. `secondOpinionRender` in `report_second_opinion.go` (from line 1706) already does the copy-and-swap internally; reuse or replace it as the builder judges.
2. **A relative `--out`.** A relative `--out` is resolved against `--repo-root`, not against the current folder, and the usage text does not say so. With `--repo-root` pointing at a scratch copy, the captures landed inside the copy. Either resolve a relative `--out` against the current folder, or state the rule in the usage text and in the action. Park the choice as a pending answer if existing callers depend on today's behavior.

### Files to change

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go` and its test file.
- `skills/do-work-toolbox/actions/ai-report.md` (the judge form) and `skills/do-work-toolbox/docs/ai-report-guide.md`, where the judge is described.
- `CHANGELOG.md` and `VERSION`, as `_dev/primes/prime-releases.md` requires.

### What must not change

- Without the new option the judge behaves exactly as today: same captures, same `judge.json`, same exit codes.
- The judge never writes inside the bundle. Bundles stay byte-identical.
- The second-opinion finalize path keeps its verdicts for the same pages.

### Verification

RED first: a Go test runs the judge with `--page index-v1c01-x.html` on a fixture bundle with two pages and fails on the unknown option.

GREEN: the test passes; the result's served address ends with the named page and the captures show that page; a named page that is missing or outside the bundle is refused with a typed finding; the fixture bundle's digest is unchanged; `go test ./...` passes in the do-work-cli module.
