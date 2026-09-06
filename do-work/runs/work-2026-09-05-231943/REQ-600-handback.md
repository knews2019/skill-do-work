# REQ-600 hand-back — put the SIGPIPE trap in the shell prime, fix the one shipped block

**Done.** Branch `worktree-agent-REQ-600-sigpipe-prime`, head `cb606f287609be7bf8e1f4657276216307b6f27c`, two commits on top of `dc5d818`. Gate run from the worktree: exit 0, last line `Maintainer verification passed.`

**Write set widened by three files under `_dev/tests/`** (below, "Guard decision"). The request's `write_set` names only the prime and `memory-reference.md`; the finalizer's allowlist needs `_dev/tests/action-shell-blocks.sh`, `_dev/tests/quiet-grep-pipeline-audit.sh` and the new `_dev/tests/quiet-grep-pipeline-scanner.sh` added. `memory-reference.md` is a shipped file, so this is a release; the changelog and version mirrors were left to the finalizer's release commit, as every recent `[REQ-NNN] release:` commit shows.

## Sweep totals

166 `.md` files under `skills/`, 32 with at least one shell fence, 74 fences (all ```bash; zero ```sh, zero ```shell, deepest indent 3, none unterminated). `action-shell-blocks.sh` counts the same 74. Hits: **1**, `skills/do-work-knowledge/actions/memory-reference.md:88`, a block agents copy and run (first of three backend probes under "Semantic Recall (Layer 2)"). A supplementary grep over the 74 extracted blocks for the readers the scanner cannot see (`rg -q`, `head`, `sed -n … q`, `awk … exit`, `read`, `grep -m` after a pipe) found zero. The only other `grep -q` in any block is `memory-reference.md:142`, a file-argument grep with no writer, not a hit.

## Block changed (commit 22a1ea4)

`skills/do-work-knowledge/actions/memory-reference.md:88-91`, before:

```
ollama list 2>/dev/null | grep -qiE 'embed'   # a local embedding model is pulled
command -v embed >/dev/null 2>&1               # a standalone embed CLI
[ -n "${OPENAI_API_KEY:-}${VOYAGE_API_KEY:-}" ] # an embeddings API key is exported
```

after:

```
ollama_models="$(ollama list 2>/dev/null || true)"  # no ollama, or a stopped daemon, means no model: the empty listing is the answer
grep -qiE 'embed' <<<"$ollama_models"               # a local embedding model is pulled
command -v embed >/dev/null 2>&1                    # a standalone embed CLI
[ -n "${OPENAI_API_KEY:-}${VOYAGE_API_KEY:-}" ]     # an embeddings API key is exported
```

Not a live failure: `ollama list` prints under 2 KB on any realistic install (about 10 KB at 100 models), far below the ~36 KB window, so the old line could not misfire in practice. This is about what shipped guidance teaches. The producer's failure (no ollama, stopped daemon) is the "no backend" answer the surrounding prose already prescribes ("silently proceed lexical-only"), so its status is collapsed with `|| true` and the comment says why in place, the sibling prime section's third way out. Extracted block passes `bash -n` and shellcheck at the lint's settings; checked by hand that a missing `ollama` yields exit 1 from the grep and a stubbed listing containing `nomic-embed-text` yields a hit.

## Prime section (commit 22a1ea4)

`_dev/primes/prime-shell-commands.md`, new `## A Writer's SIGPIPE Death Reads as the Reader's Verdict`, placed directly after `## Unchecked Exit Status Reads as Content` and before `## Closed Enumerations Go Stale`. Four paragraphs: the condition (a writer piped into a reader that can leave before the writer is done, under `pipefail`), with `grep -q`, `-m 1`, `--quiet`, `--silent`, `rg -q`, `head -1`, `sed -n '1p;q'`, `awk '/x/{exit}'`, `read` as an illustrative list; wrong in both directions, the negative matcher named as the dangerous half, the window from REQ-593's record (0 of 50 at 36 KB, 50 of 50 at 200 KB); the two-half fix (capture and herestring, plus a separate producer-status assertion because the capture discards the status, with REQ-593's truncated-archive listing that still carried the marker first, and the one-line `listing="$(…)" || fail` form), and "say so in place" when failure is the answer; the guard pointer, stating plainly that `quiet-grep-pipeline-audit.sh` walks tracked `*.sh` only with a grep-only reader set, that Markdown is not its input, and that `action-shell-blocks.sh` now runs the same scanner over the shipped fences.

No Lessons index line: the prime's `## Lessons` section is a one-line pointer to `lessons-shell-commands.md`, which the work pipeline appends on archive.

## Guard decision (commit cb606f2)

The sweep judged `action-shell-blocks.sh` a natural fit, so the scan was added there. `lint_shell_source` runs `quiet_grep_pipeline_offenders` on every extracted unit (74 fences and the 33 shipped `.sh` under `skills/`) after `bash -n` and before the shellcheck gate, so it runs without shellcheck too, and prints `FAIL: path:line: quiet grep fed from a pipeline (…): <command>` at the Markdown line. A one-shape wiring fixture (`run_quiet_grep_wiring_fixture`) runs on every default invocation, because the gate never passes `--self-test` and a pin behind that flag could not fail when the scan is removed.

Widening: the scanner function was defined inside `quiet-grep-pipeline-audit.sh` next to its top-level run, so sourcing that file runs the audit. The function and its contract comment moved byte-for-byte (diffed) into `_dev/tests/quiet-grep-pipeline-scanner.sh`, sourced by both scripts. Not copied: two scanner bodies is the drift REQ-594 measured against. The request's constraint "no change to the guard" is read as its fixture and behaviour; the file changed by a `source` line, the 19+7 fixture did not, and the audit's own bytes still scan clean.

Mutation evidence, each restored afterwards and the restore diffed:

- M1 scanner call removed from `lint_shell_source`: `action-shell-blocks.sh` exits 1, `FAIL: the fence walk no longer flags a quiet grep fed from a pipeline at its Markdown line; the shared scanner is not wired in.`
- M2 original `memory-reference.md` block restored: exits 1, `FAIL: skills/do-work-knowledge/actions/memory-reference.md:88: quiet grep fed from a pipeline (…): ollama list 2>/dev/null | grep -qiE 'embed' …`
- M3 `--quiet` dropped from the shared body: the audit exits 1, `no longer caught: grep --quiet`, so the lifted body is still pinned by shape name.

Nothing already pinned is weaker: `action-shell-blocks.sh` still reports 74 fenced blocks and 33 shipped shell files with ShellCheck enabled, `--self-test` passes, and the audit passes over 95 tracked shell files (94 plus the helper) with 19 must-flag and 7 must-not-flag shapes.

## Gate

`DO_WORK_GATE_ROOT=<worktree> bash <scratch>/gate.sh` at head `cb606f2`: exit 0, `Maintainer verification passed.` Zero FAIL lines in the log.

## Found and not fixed

- `lessons-shell-commands.md` carries no entry for REQ-593 or REQ-594 although both are archived; the pipeline's archive step (work.md Step 8 substep 7) should have appended them. Outside this write set.
- `action-shell-blocks.sh`'s fence regex accepts 0-3 leading blanks and `bash|sh` only. Nothing is lost today (74 = 74, no ```shell fences, max indent 3), but the new scan inherits that enumeration. Left as is; changing it changes what the lint already pins.
- The scanner's reader set is grep/egrep/fgrep. A fence feeding `rg -q`, `head`, `sed -n '1p;q'`, `awk '/x/{exit}'` or `read` from a pipe would pass both scripts; the sweep found zero today, and the prime now says so.
