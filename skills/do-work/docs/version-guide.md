# Version, Updates & Recap

Handles version reporting, update checks, and recent work summaries.

## Version info

Shows the current installed version and last 5 changelog entries.

```
do-work version
do-work what's new
do-work release notes
```

## Update check

Fetches the upstream version and compares. If an update is available, runs the update automatically. Warns if you've modified shipped skill files locally.

```
do-work update
do-work check for updates
just do-work-update
```

`just do-work-update` is the canonical no-agent entry point. Existing installations may continue using `just run-do-work-update`; it is a compatibility alias over the same `update-suite` transaction.

## Recap

Shows the last 5 completed user requests with their REQs and current status.

```
do-work recap
```

## Usage

```
do-work version              # current version + last 5 releases
do-work update               # natural-language route to the canonical update transaction
just do-work-update          # canonical no-agent recipe; reviews diff and asks before overwriting
just run-do-work-update      # compatibility alias for the same transaction
do-work what's new           # same as version
do-work release notes        # same as version
do-work history              # same as version
do-work recap                # last 5 completed URs with their REQs
```
