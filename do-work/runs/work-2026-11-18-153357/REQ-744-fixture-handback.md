# REQ-744 fixture handback

Both disposable Git consumers are ready for the writer. Their checked-in evidence is synthetic test data. Each has a completed REQ and parent UR input, an implementation commit, and an archive commit that records the implementation SHA and actual verification output. The private parent directory has mode 0700. Both worktrees are clean and have no remotes or report bundles.

| Consumer root | Ordinary target | Implementation commit | Archive commit |
|---|---|---|---|
| `/private/tmp/req-744-acceptance-8u9ze8c1/request-outcomes` | `REQ-001` or `UR-001` | `31f4eccee016c46cb391cc2404b282171bef91c9` | `95affa8d7aeffacc9e4a441a78a72450c28fc587` |
| `/private/tmp/req-744-acceptance-8u9ze8c1/label-correction` | `REQ-002` or `UR-002` | `da70efd554643081d3c2b5ef97c19cfa1e961df4` | `4376515ee7d0ef33dbec8cc0674de28540304e20` |

The real installer installed the REQ-743 baseline suite from `82a0eefe05f5cefee388151560ec56aeb378d6de`. The installed CLI's gallery list returned four rows in both consumers and every sample record path resolved. The consumers keep the installer outputs in local Git excludes so the committed histories stay small and second-opinion prepare can start with a clean worktree. The per-consumer install logs, gallery lists, and merge-aware implementation diffs are under the private parent directory.

## Install the writer candidate

From the shared workspace, pass the full committed candidate SHA to this subshell's `revision` variable in place of `CANDIDATE_SHA`. Use the actual integrated revision so the installed action, reference, kit, presets and CLI are from one tree.

```bash
(
  set -euo pipefail
  source_repo=/Users/t2/.t3/worktrees/skill-do-work2/t3-a5de53e3-5b37f50fbf62596b
  scratch=/private/tmp/req-744-acceptance-8u9ze8c1
  revision=CANDIDATE_SHA
  setup=$(mktemp -d "$scratch/candidate.XXXXXXXX")
  git -C "$source_repo" archive --format=tar.gz --prefix=suite/ -o "$setup/suite.tar.gz" "$revision"
  tar -xzf "$setup/suite.tar.gz" -C "$setup"
  for name in request-outcomes label-correction; do
    repo="$scratch/$name"
    bash "$setup/suite/skills/do-work/tools/install-do-work-suite.sh" \
      --project-root "$repo" --archive "$setup/suite.tar.gz" <<< y > "$setup/$name-install.log" 2>&1
    toolbox=.claude/skills/do-work-toolbox
    cmp "$setup/suite/skills/do-work-toolbox/actions/ai-report.md" "$repo/$toolbox/actions/ai-report.md"
    cmp "$setup/suite/skills/do-work-toolbox/actions/ai-report-reference.md" "$repo/$toolbox/actions/ai-report-reference.md"
    cmp "$setup/suite/skills/do-work-toolbox/docs/ai-report-kit/skeleton.html" "$repo/$toolbox/docs/ai-report-kit/skeleton.html"
    test -z "$(git -C "$repo" status --porcelain)"
    bash "$repo/.claude/skills/do-work/tools/do-work-cli.sh" --repo-root "$repo" \
      --format text ai-report-gallery --list
  done
)
```

## Writer invocations and expected results

Run these as agent actions from the indicated project root. Route through installed `.claude/skills/do-work-toolbox/SKILL.md` to `actions/ai-report.md`. The shell CLI is for the action's gallery, provenance, and prepare steps. Give the writer the invocation and question only, with the oracles reserved for review.

- **Request outcomes.** From `/private/tmp/req-744-acceptance-8u9ze8c1/request-outcomes`, invoke `do-work-toolbox ai-report REQ-001` and ask “We served more successful requests. Did the success rate improve?” Replay `python3 scripts/check-request-outcomes.py` from that root. The script validates and aggregates 300 committed outcome rows into two adjacent 24-hour UTC windows. Window A spans Nov 1 to Nov 2 and has 90 successes out of 100 (90%). Window B spans Nov 2 to Nov 3 and has 160 out of 200 (80%). The verdict is no. The count rose by 70 and the rate fell by 10 percentage points.
- **Request-outcomes presentation.** The page head has one `ai-report-gallery-sample` meta with content `two-measures-side-by-side`. The two charts sit side by side when width permits and use the same windows in the same order, with titles identifying successful request count and success rate in percent. Start and end values are visible and sourced to the fixture. The report labels the synthetic data and uses the kit with the sample's presentation moves. The writer's handback names the sample and explains the count-versus-rate fit. Check the action trace for the real `ai-report-gallery --list` call after evidence collection and before visual construction, followed by reading the selected sample records.
- **Runner-up and teaching.** Expect at least one and at most two fitting gallery layout prompts, with at least one real preset. `verdict-with-checks` is the expected runner-up because the reader can check the denominator, windows, units and arithmetic. Look for a one- or two-sentence rationale per alternative and the existing teaching suggestion. The writer prints every closing command and launches no second-opinion model.
- **Label correction.** From `/private/tmp/req-744-acceptance-8u9ze8c1/label-correction`, invoke `do-work-toolbox ai-report REQ-002`. Verify the visible verdict and the locale evidence say “Account settings.” The history shows one string correction from “Acount settings.” The page head has a single `ai-report-gallery-sample` meta with content `none`, and the handback says no sample fit. Look for the existing alternate layout and teaching suggestions, with one runnable line for each. This fixture has no measurements, judgment to audit, cycle or UI captures.
- **Alternate routes.** In the metrics consumer, use `do-work-toolbox ai-report --kind proposal "Test fixture: how to present successful request count and rate together"` and `do-work-toolbox ai-report --kind root-cause "Test fixture: why success rate fell while successful request count rose"`. Cite the same committed CSVs and replay. For root cause, keep the claim to the arithmetic demonstrated by the fixture. For proposal, compare display options using that evidence. Both new heads must carry their kind and gallery decision. After the ordinary metrics report exists, invoke `do-work-toolbox ai-report revise <ordinary-bundle-dir> "Clarify the denominator while keeping both fixture windows"`. Verify the new head and handback record a fresh fit decision, closing lines target the new bundle, and the prior bundle stays byte-identical. For these alternate routes use the actual evidence fit, preserve the route's own template and confirm the gallery decision is reached before building visuals.

## Verify displayed commands through prepare

For every closing line, compare the displayed text with the writer's handback and tokenize with `shlex.split`. Expect `ai-report second-opinion`, the declared options and one positional target resolving inside the current bundle location, with that bundle's exact path. If present, `--editor-direction` has one nonempty argument without a single quote. Resolve each `--preset` by looking first for `<bundle>/second-opinion-presets/<name>.md`, then for `.claude/skills/do-work-toolbox/docs/second-opinion-presets/<name>.md` inside the consumer. Review how the preset helps the proposed alternative.

Save the actual directed runner-up command as a single line in a scratch text file. From the metrics consumer, pass the file to the following prepare-only probe by replacing `PATH_TO_COMMAND_FILE`. It forwards the split command arguments without evaluating shell text.

```bash
cd /private/tmp/req-744-acceptance-8u9ze8c1/request-outcomes
python3 - PATH_TO_COMMAND_FILE <<'PY'
import json
from pathlib import Path
import shlex
import subprocess
import sys

repo = Path.cwd()
tokens = shlex.split(Path(sys.argv[1]).read_text())
assert tokens[:2] == ["ai-report", "second-opinion"]
assert len(tokens) == 7 and set(tokens[2:-1:2]) == {"--preset", "--editor-direction"}
options = dict(zip(tokens[2:-1:2], tokens[3:-1:2]))
direction = options["--editor-direction"]
assert direction and "'" not in direction
bundle = (repo / tokens[-1]).resolve()
assert bundle.parent == repo / "ai-reports" and (bundle / "index.html").is_file()
preset = options["--preset"]
assert Path(preset).name == preset
local = bundle / "second-opinion-presets" / (preset + ".md")
installed = repo / ".claude/skills/do-work-toolbox/docs/second-opinion-presets" / (preset + ".md")
expected = local if local.is_file() else installed
assert expected.is_file()
original = (bundle / "index.html").read_bytes()
launcher = repo / ".claude/skills/do-work/tools/do-work-cli.sh"
raw = subprocess.check_output([
    "bash", str(launcher), "--repo-root", str(repo), "--format", "json",
    "ai-report-second-opinion", "--prepare", "--backend", "codex",
    "--model", "gpt-6.2-sol", *tokens[2:],
], text=True)
print(raw)
row, = json.loads(raw)["ai_report_second_opinion"]["runs"]
assert row["state"] == "prepared" and row["judge_verdict"] == "not-run"
assert row["editor_direction"] == direction and row["preset"] == preset
assert (repo / row["preset_file"]).resolve() == expected
assert row["prompt"] and direction in row["prompt"]
assert row["staging_page"] and row["staging_reply"]
assert (bundle / "index.html").read_bytes() == original
PY
```

This prepare-only check exercises the existing-page dispatcher and verifies that the selected preset and direction survive into its existing editor prompt. For an explicit existing-page target, use the same launcher with `--prepare`, the declared backend and model, an existing `<bundle>/index.html`, and a single `--editor-direction` argument. Expect the original page as the base and the supplied direction in the prepared prompt, with no gallery reranking. Stop at prepare for these checks.

## Provenance replay

In either consumer, run this from the repository root. The archive's `commit` value must resolve to the reachable implementation parent of the archive commit. The receipt below the archived `## Testing` heading gives the fixture command and exact output.

```bash
req_file=$(find do-work/archive -name 'REQ-*.md' -type f)
implementation=$(sed -n 's/^commit: //p' "$req_file")
archive_commit=$(git log -1 --format=%H -- "$req_file")
test "$(git rev-parse "$archive_commit^")" = "$implementation"
git cat-file -e "$implementation^{commit}"
bash .claude/skills/do-work/scripts/show-commit-diff.sh "$implementation"
```

Fixture preparation verified both test receipts, both clean Git histories and the real baseline gallery list. Writer output, render checks, alternate routes and prepare are pending the REQ-744 candidate.
