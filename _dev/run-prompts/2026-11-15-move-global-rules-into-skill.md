# Capture prompt: move the maintainer's global rules into the shipped skill

## How to run this prompt

Capture the request below with `do-work capture-request`. Do not build it in the capture session. Run unattended: when something is unclear, park it as a pending answer and continue. If nothing moves for 20 minutes, stop and report what blocked. Operator steps are left out on purpose.

## The request

The maintainer keeps 35 lines of standing rules in `~/.claude/CLAUDE.md`. That file exists on one machine only, so sessions on other machines behave differently. The maintainer wants these rules in the do-work skill, which is checked into every repository. Move the rules that the skill does not carry yet, with the smallest addition. The global file is deleted afterwards, outside this request.

### The source text (copy of `~/.claude/CLAUDE.md`, 2026-11-15)

Treat this block as the content to move. The build session may not have the file.

```
# Standing rules (all projects)

## Explaining state or a system
- Big picture first (one short paragraph), then the exact picture.
- Define every status word in plain words on first use (degraded, in-flight, earmarked, claimed, blocked, parked, pending).
- Close with the next command to type.

## Decisions
- Any decision for the user goes through the AskUserQuestion tool: recommended option first, value and risk for each option. Never a prose list of questions.
- Decide reversible trivia yourself and report it in one line.

## Shell and commands for the user to run
- The user's terminal is zsh (login shell /bin/zsh). Remote hosts are bash: wrap remote payloads as ssh host 'bash -s' <<'EOF' … EOF.
- Fill in every value. No M=/<path> placeholders in a command the user will paste.
- Logic longer than about 5 lines goes into a committed script with a one-line launcher.

## Scope
- A bug report means fix the code, not write a memory note.
- Never mint a REQ or task for operator work (deploy, publish, verify on live hosts, approve content). Write a runbook report under ai-reports/ and capture only the tooling gaps the AI must build.
- Numeric targets are approximate unless the user says strict.
- Never remove existing UI while redesigning unless asked.

## Prompts for another session or tool
- Write the full prompt to docs/run-prompts/<date>-<slug>.md, commit it by exact path, and return a paste block of 3 lines or fewer.
- The prompt must run unattended: park questions as pending-answers, include a 20-minute stall check, leave operator steps out.

## Concurrent sessions in one checkout
- Commit by exact path only (never git add -A or git commit -a).
- Re-read the highest REQ/UR id at write time; never stage or revert files you did not change.
- The checkout-lease SessionStart hook prints the other live sessions; prefer EnterWorktree for build work when it fires.

## Subagents
- Pick the model per Agent call by the work: lookups, searches, file sweeps and Explore agents get model: haiku; planning, review and build agents inherit the session model. Never leave a pure lookup on the default.
```

### What a text search found on 2026-11-15 (re-check each line at HEAD)

| Code | Rule group | State in `skills/` |
|---|---|---|
| G1 | Explaining state | Mostly present in `skills/do-work/crew-members/communication-style.md` ("Answer first, next action last"). Missing: big picture before the exact picture, and a plain definition of every status word on first use. |
| G2 | Decisions through the ask tool | Not found. `skills/do-work/crew-members/clear-questions.md` covers how to word a question, not "use the ask tool, recommended option first, value and risk per option, decide reversible trivia yourself". |
| G3 | Commands for the user | Not found. |
| G4 | Scope | The operator-work rule is shipped (0.305.85). The other three lines are not found. |
| G5 | Prompts for another session | Only the folder name, in `skills/do-work-toolbox/actions/folder-structure.md`. The unattended shape, the stall check and the short paste block are not found. |
| G6 | Several sessions in one checkout | Exact-path commits are shipped. The `checkout-lease` hook lives in `~/.claude/settings.json` and is out of scope here. |
| G7 | Model per subagent | Not found. |

### Rules for the move

- Delete before you add. For each source line, first look for text that already says it. Add a line only when none does. Do not restate a shipped rule in a second file.
- Always-on rules must sit in an always-on file. `communication-style.md` is the only file the installer links from a consumer project's agent instructions, so G1, G2 and G3 apply outside a work loop only if they live there, or in a file it tells the agent to load on a named condition. Do not add a second always-on file.
- Write conditions, not this maintainer's machine. "zsh" becomes "write commands for the user's login shell, and wrap a payload for a remote host in that host's shell". "AskUserQuestion" becomes "the agent's ask-user tool, when it has one". "haiku" becomes "the cheapest model tier that can do the lookup". No vendor, tool or host name in shipped text.
- Suggested homes, which the build session may change with a stated reason: G1 to G3 in `communication-style.md`; the option-shape half of G2 in `clear-questions.md`; G4 in `skills/do-work/crew-members/coding-guardrails.md`; G7 in `skills/do-work/crew-members/background-agents.md`; G5 wherever the skill already describes a handoff to another session. If G5 has no natural home, park that as a pending answer.
- G5 names a folder. Use the project's run-prompts folder as `folder-structure.md` defines it, and do not write a second path. In this repository the folder question is still open (see `_dev/run-prompts/2026-11-15-run-prompts-folder-rule.md`).
- Keep the addition small. About 25 shipped lines in total is the target, and it is approximate.

### Files to change

- `skills/do-work/crew-members/communication-style.md`
- `skills/do-work/crew-members/clear-questions.md`
- `skills/do-work/crew-members/coding-guardrails.md`
- `skills/do-work/crew-members/background-agents.md`
- The changelog and version, as `_dev/primes/prime-releases.md` requires for a change to shipped files.

### What must not change

- `~/.claude/CLAUDE.md` and `~/.claude/settings.json`. This request never reads, edits or deletes a file outside the repository.
- No new crew-member file and no new always-on file.
- No rule that is already shipped is reworded.
- The Aliases section and the Examples section of `communication-style.md`.

### Verification

RED first: record, for each of G1, G2, G3, G4, G5 and G7, the search that shows the rule is absent from `skills/` at the starting commit.

GREEN: the same six searches each find exactly one shipped file. A search for `zsh`, `AskUserQuestion`, `haiku` and `EnterWorktree` in the four changed crew-member files finds nothing. The fast gate passes. No new test file is added, because these are prose rules and a test of their wording would pin nothing real.

### Open questions to park as pending answers

- Q1: G6's session list comes from a hook that only this machine has. Should the skill's installer ship that hook, or is the exact-path commit rule enough?
- Q2: Should the notes in the maintainer's local memory folder get the same treatment in a later request?
