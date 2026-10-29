# Periodic checks

`do-work-toolbox periodic-check` runs the project's full checklist. Add `quick` for its explicit subset or `status` to read both profiles' history without running checks or changing files. Add `--capture` to create repair requests, `--report-only` to override all capture settings, or `--unattended` to defer questions instead of prompting. Capture never runs repairs.

## Project setup

Create `prime-periodic-check.md` at the project root using the [security example](prime-periodic-check.example.md). Replace its example commands with the project's actual diagnostic commands and document prerequisites, timeouts, and interpretation. No setup or scanner installation is implicit.

The prime has exactly one fenced block labelled `periodic-check` containing JSON:

- `version`: nonempty checklist compatibility version. Change it when the meaning/scope of the checklist changes.
- `capture`: explicit standing authorization, default false. This must be deliberately configured by the maintainer.
- `checks`: nonempty list, each with a unique lowercase hyphenated `id`, nonempty compatibility `version`, positive `cadence_seconds` (up to ten years), `quick` boolean, diagnostic `command`, `prerequisites`, and `interpretation`.

Commands are read and run by the agent, not evaluated by the history tool. `quick: true` opts a check into quick mode; no marked checks means quick is undefined. Versions define comparison compatibility rather than software dependency versions. Keep advisory/component identities stable while changing summaries only for substantive changes. A project's security tooling may download current advisory data, but must not install updates or change deployment state.

### Optional slash shortcut

For Claude Code clients supporting project command files, copy the shipped [shortcut template](../shortcuts/periodic-check.md) to the project's `.claude/commands/periodic-check.md` during explicit setup. Preserve an existing command rather than overwriting it. The template forwards all arguments to the installed toolbox action. This is a client adapter, not another required suite package. For other clients, use their supported alias mechanism if available; the portable invocation is always `do-work-toolbox periodic-check`.

## Mechanical helper contract

The [core CLI](../../do-work/tools/do-work-cli.sh) exposes `periodic-check-history`. Global `--repo-root` and optional `--format json` precede that command. The JSON command envelope retains its normal exact-text output field; `record.json` is the structured evidence. These verbs never run checks, capture requests, commit, or push:

| Verb | Arguments and behavior |
|---|---|
| `status` | Read current prime and both profiles' available history; no writes |
| `start` | `--profile full|quick`, optional `--capture`, `--report-only`, `--unattended`; show prior history and return `Run: <relative-path>` |
| `record` | `--run <relative-path> --input <JSON-file>`; record one selected check and derive finding changes |
| `link` | Same arguments; record capture disposition for one active finding |
| `finish` | `--run <relative-path>`; publish report and freeze completed attempt, including incomplete coverage |

A `record` input:

```json
{"check_id":"go-dependencies","completed":true,"evidence":"govulncheck completed; advisory X affects package Y","findings":[{"id":"advisory-X/package-Y","summary":"Package Y is affected by advisory X","evidence":"Observed affected version and scanner result"}]}
```

Use `completed: false` and an explanation for missing tools, credentials, timeouts, or uncertain results. An empty findings array alone never proves coverage. Do not supply timestamps or derived finding fields; the helper owns them. To record a clean check, use `completed: true` with actual evidence and `findings: []`.

A `link` input:

```json
{"check_id":"go-dependencies","finding_id":"advisory-X/package-Y","request_id":"REQ-123","disposition":"linked","note":"Existing pending request covers the same package and advisory"}
```

Dispositions are `linked`, `created`, `verify-resolution`, `deferred`, or `report-only`. The first three require a request ID; `created` also requires effective capture authorization. Every disposition needs a note. Links are evidence of the canonical capture outcome, never instructions to mutate a request.

Each unique attempt directory contains `record.json` and, once finished, `report.md`. The record preserves initial revision/dirt, configuration snapshot, profile, timestamps, per-check findings/coverage and request links. Save record files through the helper only. Reports use relative paths and can be read after a normal Git synchronization. Do not place credentials in evidence.

Last completed and clean runs remain visible even when a newer attempt fails. Due checks derive from their last equivalent successful execution in the same profile. Older incompatible results remain historical. A successful check within an incomplete attempt can supply evidence for that check, but cannot make the overall attempt completed or cover missing checks. A new or changed finding stays unresolved until an equivalent check actually establishes absence.

## Hosted execution

Configure a hosted machine with the repository, scanner prerequisites, and an unattended agent separately. Give that agent this same prompt:

```text
Read .claude/skills/do-work-toolbox/SKILL.md, then run:
do-work-toolbox periodic-check --unattended --report-only
Use the project prime and canonical history helper. Report missing prerequisites as incomplete. Do not prompt, install updates, fix files, push, or deploy.
```

For authorized capture, replace `--report-only` with `--capture`. The agent must use core capture, including its ordinary exact-path local commit; separately granted Git write access is necessary for that mode. Its scanner commands and result interpretation remain identical to a manual run. Agent CLI syntax is selected by the host configuration; the skill is not itself a shell executable or scheduler.

To schedule weekly execution, the operator configures the hosted scheduler, agent credentials, scanner credentials, timeout and notification destination. Serialize jobs per repository. Each job synchronizes shared history before running. A separately authorized publication step then commits only this attempt's records and any already-authorized capture artifacts, publishes through the chosen branch/PR process, and reports publication failure separately from scan results. Never force-push on conflict. A job interrupted mid-scan should preserve/share its partial record if its configured publication stage can run.

A local manual run does not share itself. Another checkout sees records only after normal commit/push/fetch synchronization. Without shared records, the tool cannot know about remote attempts or guarantee cross-checkout capture deduplication. Schedules, notifications, credentials, write permissions, and the consumer's actual security prime are outside suite installation.
