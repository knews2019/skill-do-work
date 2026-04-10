# Specification Templates

Reusable templates defining output structure, quality standards, and implementation patterns for common task types.

**How to use:** Specs are referenced during `capture` (optional `--spec` hint in `$ARGUMENTS`) or loaded automatically during `work` when the REQ's domain or task type matches a template. The work action checks `specs/` for a matching file after triage and before planning.

**How they work:** A spec informs the implementation checklist order, quality standards to verify against, and common pitfalls to watch for. The spec is guidance, not override — the REQ's specific requirements always take priority.

**How to create new ones:** Follow the existing template structure: Output Structure, Quality Standards, Implementation Checklist, Evolution Path, Common Pitfalls. Name the file after the task type (e.g., `migration.md`, `cli-command.md`). Keep it concise — a spec that's longer than the implementation it guides has failed.

**Available specs:** `api-endpoint.md`, `ui-component.md`, `refactor.md`, `bug-fix.md`
