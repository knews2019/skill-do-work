# Project periodic security checks

Copy to the project root as `prime-periodic-check.md` and replace the example commands with existing project-owned scanner entrypoints before use. The commands below are examples, not tools supplied by the suite. They must perform diagnostics only. Keep credentials in the runner's configured secret store, never in this file or saved evidence.

```periodic-check
{
  "version": "security-checklist-1",
  "capture": false,
  "checks": [
    {"id":"go-dependencies","version":"1","cadence_seconds":604800,"quick":true,"command":"./project-checks go-dependencies","prerequisites":"Go scanner installed; advisory service reachable","interpretation":"Complete only after the scanner finishes; report affected advisory/package pairs. Timeout after 10 minutes means incomplete."},
    {"id":"node-dependencies","version":"1","cadence_seconds":604800,"quick":true,"command":"./project-checks node-dependencies","prerequisites":"Node audit tool and lockfile available","interpretation":"A findings exit code is completed-with-findings; registry/authentication failure is incomplete. Timeout after 10 minutes."},
    {"id":"runtime-support","version":"1","cadence_seconds":604800,"quick":false,"command":"./project-checks runtime-support","prerequisites":"Declared/deployed runtime evidence and current vendor lifecycle source available","interpretation":"Compare actual runtimes with vendor support/security information. Missing deployment evidence is incomplete. Timeout after 10 minutes."},
    {"id":"deployment-images","version":"1","cadence_seconds":604800,"quick":false,"command":"./project-checks deployment-images","prerequisites":"Read-only registry/deployment credentials and image scanner available","interpretation":"Record deployed image digest and affected advisory/component pairs. Never rebuild or redeploy. Missing registry access is incomplete. Timeout after 15 minutes."}
  ]
}
```

Use stable advisory/component IDs for dependency/image findings, and stable runtime/component IDs for unsupported-runtime findings. Put changing severity or affected versions in the summary; put timestamps, command output and source URLs in evidence. Amend this example to match the project's actual technologies, commands and cadence; do not invent unrelated audits. Bump checklist/check compatibility versions when changing coverage or interpretation, not merely when an advisory database updates.
