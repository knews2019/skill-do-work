# Validated Runtime Boundaries

Three runtime replays exposed one shared mistake: the code verified a convenient proxy while the real boundary lived one level wider. Keep these lessons together because timeout ownership, directory publication, artifact freshness, and backend authority all depend on naming the complete boundary.

## A timeout owns the process tree it starts

Killing a launcher is not a timeout when its descendants can keep touching services. A portable fallback must establish and verify an isolated process group before running probe code, signal that group on timeout, escalate when needed, and reap its leader. If isolation cannot be proved, fail closed without ever signaling the caller's group.

The lock-in fixture forces the stock-Bash path, records wrapper and descendant PIDs, and keeps an unrelated process alive as evidence that cleanup is complete and scoped.

## A directory is healthy only as a complete payload

A sentinel file is not an installation. Define one completeness predicate over the minimum runnable payload and use it for source validation, staged validation, `check`, and install/no-op decisions. Copy the full version into private adjacent staging, validate it there, and publish by same-filesystem rename; never merge files from different versions. Hold an incomplete prior tree as a private backup until the replacement is live and restore it on publication failure.

## Presence is not current-invocation provenance

An old non-empty artifact can survive a failed producer. Each invocation writes to its own private adjacent file and publishes only after that producer exits successfully and the staged artifact passes validation. Parallel callers retain every PID and wait status, wait all jobs, and use status—not target presence—to decide which outputs are current.

## Opt-in grants authority; a private cwd does not contain it

`DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=1` explicitly authorizes a sandbox-bypassed process with repository, credential, network, and external-side-effect capability. A locked temporary working directory reduces accidental output spread but is not a sandbox. Keep safe SVG/Mermaid completion sufficient, the agentic path default-off, its prompt sanitized, and its output behind the same current-invocation publication boundary.
